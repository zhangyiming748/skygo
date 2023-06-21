package controller

import (
	"errors"
	"fmt"
	"time"

	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/http/transformer"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/log"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/license"
	_ "skygo_detection/lib/license"
	"skygo_detection/mongo_model_tmp"
	"skygo_detection/mysql_model"
)

type TaskController struct{}

// 查询所有任务
func (this TaskController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	s := mysql.GetSession()
	s.Where("create_user_id=?", http_ctx.GetUserId(ctx))
	// 查询组键
	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	widget.AndWhereEqual("create_user_id", http_ctx.GetUserId(ctx))
	widget.AddSorter(*(orm.NewSorter("create_time", orm.DESCENDING)))
	widget.SetTransformer(&transformer.TaskTransformer{})
	all := widget.PaginatorFind(s, &[]mysql_model.Task{})
	response.RenderSuccess(ctx, all)
}

// 查询单个任务
func (this TaskController) GetOne(ctx *gin.Context) {
	id := request.ParamString(ctx, "id")
	s := mysql.GetSession()
	s.Where("id=?", id)
	s.And("create_user_id=?", http_ctx.GetUserId(ctx))
	widget := orm.PWidget{}
	widget.SetTransformer(&transformer.TaskTransformer{})
	result, err := widget.One(s, &mysql_model.Task{})

	if err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, errors.New("任务不存在"))
	}
}

// 创建任务
func (this TaskController) Create(ctx *gin.Context) {
	if result := license.VerifyMenu(license.TEST_TASK); !result {
		response.RenderFailure(ctx, errors.New("授权证书无效，创建任务失败"))
		return
	}
	form := mysql_model.TaskCreateForm{}
	form.Name = request.MustString(ctx, "name")
	form.AssetVehicleId = request.MustInt(ctx, "asset_vehicle_id")
	form.PieceId = request.MustInt(ctx, "piece_id")
	form.PieceVersionId = request.MustInt(ctx, "piece_version_id")
	form.FirmwareTemplateId = request.Int(ctx, "firmware_template_id")
	form.NeedConnected = request.MustInt(ctx, "need_connected")
	form.Describe = request.MustString(ctx, "describe")
	if form.IsToolTask = request.MustInt(ctx, "is_tool_task"); form.IsToolTask == common.IS_TOOL_TASK {
		form.ToolId = request.MustString(ctx, "tool_id")
		form.Tool = request.MustString(ctx, "tool")
		form.Category = request.MustString(ctx, "category")
	} else {
		form.ScenarioId = request.MustInt(ctx, "scenario_id")
	}

	if form.IsToolTask == common.NOT_TOOL_TASK {
		if total, err := sys_service.NewSessionWithCond(qmap.QM{"e_scenario_id": form.ScenarioId}).Count(new(mysql_model.KnowledgeTestCase)); err == nil {
			if total <= 0 {
				response.RenderFailure(ctx, errors.New("任务创建失败，所选择场景中未包含测试用例"))
				return
			}
		}
		new(mysql_model.KnowledgeTestCase).KnowledgeTestCaseFindByScenarioIds(form.ScenarioId)
	}
	if taskModel, err := mysql_model.TaskCreate(&form, ctx); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		// 根据任务中的字段判断不同的子任务
		if form.IsToolTask == common.IS_TOOL_TASK {
			// 创建固件任务
			if form.Tool == common.TOOL_FIRMWARE_SCANNER {
				// 创建固件扫描任务
				pieceVersion, err := new(mysql_model.AssetTestPieceVersion).FindById(form.PieceVersionId)
				if err != nil {
					panic(err)
				}
				firmwareTask := new(mysql_model.FirmwareTask)
				firmwareTask.TaskId = taskModel.Id
				firmwareTask.Name = fmt.Sprintf("固件子任务_%s_%d", taskModel.Name, custom_util.GetCurrentMilliSecond())
				firmwareTask.FileId = pieceVersion.FirmwareFileUuid
				firmwareTask.TemplateId = form.FirmwareTemplateId
				firmwareTask.CreateTime = int(time.Now().Unix())
				firmwareTask.UpdateTime = int(time.Now().Unix())
				firmwareTask.Status = common.FIRMWARE_STATUS_PROJECT_CREATE
				if _, err := firmwareTask.Create(); err == nil {
					if err := new(mysql_model.ScannerTask).TaskInsert(firmwareTask.Id, firmwareTask.Name, common.TOOL_FIRMWARE_SCANNER); err != nil {
						response.RenderFailure(ctx, err)
					}
					taskModel.ToolTaskId = fmt.Sprintf("%d", firmwareTask.Id)
					taskModel.Update("tool_task_id")
					response.RenderSuccess(ctx, gin.H{"id": taskModel.Id, "tool_task_id": taskModel.ToolTaskId})
				} else {
					response.RenderFailure(ctx, err)
				}
			} else if form.Tool == common.TOOL_VUL_SCANNER {
				// 创建车机漏扫任务
				req := &qmap.QM{}
				*req = req.Merge(*request.GetRequestBody(ctx))
				*req = req.Merge(*request.GetRequestQueryParams(ctx))
				(*req)["parent_task_id"] = taskModel.Id
				if result, err := new(mongo_model_tmp.EvaluateVulTask).Create(taskModel.TaskUuid, *req); err == nil {
					taskModel.ToolTaskId = result.TaskID
					taskModel.Update("tool_task_id")
					// todo mysql里也创建任务，目前这个任务还没有串联起来，因为客户端的上传结果还是存mongo库里边
					name := fmt.Sprintf("固件子任务_%s_%d", taskModel.Name, custom_util.GetCurrentMilliSecond())
					_, err = new(mysql_model.VulTask).Create(name, taskModel.Id, taskModel.ToolTaskId)
					if err != nil {
						fmt.Println(err)
					}
					response.RenderSuccess(ctx, gin.H{"id": taskModel.Id, "tool_task_id": taskModel.ToolTaskId})
				} else {
					response.RenderFailure(ctx, err)
				}
			} else {
				// todo 工具不合规 只创建父类任务
				response.RenderSuccess(ctx, gin.H{"id": taskModel.Id})
			}
		} else if form.IsToolTask == common.NOT_TOOL_TASK {
			// 创建场景的任务
			this.CreateSubTask(form.ScenarioId, taskModel)
			response.RenderSuccess(ctx, gin.H{"id": taskModel.Id})
		} else {
			response.RenderFailure(ctx, errors.New("无法创建任务"))
		}
	}
}

// 更新任务
func (this TaskController) Update(ctx *gin.Context) {
	id := request.ParamInt(ctx, "id")
	name := request.MustString(ctx, "name")
	describe := request.MustString(ctx, "describe")
	task := new(mysql_model.Task)
	task.Id = id
	task.Name = name
	task.Describe = describe
	if _, err := task.Update("name", "describe"); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, nil)
	}
}

// 删除任务
func (this TaskController) BulkDelete(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	userId := int(http_ctx.GetUserId(ctx))
	userName := http_ctx.GetUserName(ctx)
	successNum := 0
	failureNum := 0
	if _, has := req.TrySlice("ids"); has {
		ids := req.SliceInt("ids")

		for _, id := range ids {
			task := mysql_model.Task{}
			if has, err := sys_service.NewSession().Session.Where("id=?", id).Get(&task); err == nil && has {
				if task.Status != common.TASK_STATUS_FAILURE && task.Status != common.TASK_STATUS_SUCCESS {
					// 任务只有在测试失败或者测试成功的状态下才能删除
					failureNum++
					continue
				}
				task.Status = common.TASK_STATUS_REMOVE
				new(mysql_model.TaskLog).Insert(userId, userName, &task)
				_, err := sys_service.NewSession().Session.ID(task.Id).Delete(new(mysql_model.Task))
				if err != nil {
					log.GetHttpLogLogger().Error(fmt.Sprintf("%v", err))
					failureNum++
				} else {
					// 任务删除后也要删除任务关联的漏洞
					sys_service.NewSession().Session.Where("task_id=?", task.Id).Delete(new(mysql_model.Vulnerability))
					// 删除任务关联的漏洞扫描报告中的漏洞关系
					new(mysql_model.VulnerabilityScannerVulRelation).Delete(task.Id, 0, 0)
					successNum++
				}
			} else {
				failureNum++
			}
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"success_num": successNum, "failure_num": failureNum})
}

func (this TaskController) CreateSubTask(scenarioId int, parentTask *mysql_model.Task) {
	// 1. 根基场景找出所有测试用例
	testCases := new(mysql_model.KnowledgeTestCase).KnowledgeTestCaseFindByScenarioIds(scenarioId)
	// 1.1 复制测试用例
	// 1.2 根据不同组创建子任务
	// 1.3 根据子任务去分发到不同的地方测试
	mysql_model.CopyTestCases(testCases, parentTask)
}

// 查询任务中的测试件
func (this TaskController) GetAssetTestPieces(ctx *gin.Context) {
	task := new(mysql_model.Task)
	id := request.ParamString(ctx, "id")
	s := mysql.GetSession()
	s.Where("id=?", id)

	widget := orm.PWidget{}
	if _, err := widget.One(s, task); err != nil {
		response.RenderFailure(ctx, errors.New("任务不存在"))
		return
	}
	// 查询task 里的测试件，返回测试件的信息
	piece := new(mysql_model.AssetTestPiece)
	s = mysql.GetSession()
	s.Where("id=?", task.PieceId)
	widget = orm.PWidget{}
	_, err := widget.One(s, piece)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	// 返回测试件中 hg ，固件
	pieceVersion := new(mysql_model.AssetTestPieceVersion)
	s = mysql.GetSession().Where("id=?", task.PieceVersionId)
	widget = orm.PWidget{}
	_, err = widget.One(s, pieceVersion)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	// 查询车机信息
	vehicle := new(mysql_model.AssetVehicle)
	s = mysql.GetSession().Where("id=?", piece.AssetVehicleId)
	widget = orm.PWidget{}
	_, err = widget.One(s, vehicle)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	// 返回测试件中 文件
	pieceVersionFile := new(mysql_model.AssetTestPieceVersionFile)
	s = mysql.GetSession().Where("version_id=?", pieceVersion.Id)
	widget = orm.PWidget{}
	widget.One(s, pieceVersionFile)
	// if err != nil {
	// 	response.RenderFailure(ctx, err)
	// 	return
	// }
	var result = make(map[string]interface{}, 0)
	pieceInfo := map[string]interface{}{
		"name":          piece.Name,
		"piece_type":    piece.PieceType,
		"piece_version": pieceVersion.Version,
		"brand":         vehicle.Brand,
		"code":          vehicle.Code,
	}
	result["piece_info"] = pieceInfo
	result["need_connected"] = task.NeedConnected
	firmwareInfo := map[string]interface{}{
		"name":        pieceVersion.FirmwareName,
		"version":     pieceVersion.Version,
		"device_type": pieceVersion.FirmwareDeviceType,
		"file":        pieceVersionFile.FileName,
		"size":        pieceVersionFile.FileSize,
	}
	result["firmware_info"] = firmwareInfo

	if err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this TaskController) GetTaskCases(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	id := request.ParamString(ctx, "id")
	body := request.GetRequestBody(ctx)
	s := mysql.GetSession()
	s.Where("task_id=?", id)
	if _, has := body.TryInterface("module_id"); has {
		moduleArr := body.SliceString("module_id")
		s.In("module_id", moduleArr)
	}
	//判断是否采用demand_id查询
	if _, has := body.TryInterface("demand_chapter_id"); has {
		demandChapterArr := body.SliceInt("demand_chapter_id")
		testCaseChapters := new(mysql_model.KnowledgeTestCaseChapter).GetByDemandChapterIds(demandChapterArr)
		var testCaseArr = make([]int, 0)
		for _, tmptestCase := range testCaseChapters {
			testCaseArr = append(testCaseArr, tmptestCase.TestCaseId)
		}
		s.In("test_case_id", testCaseArr)
	}
	widget := orm.PWidget{}

	widget.SetQueryStr(queryParams)

	widget.SetTransformer(&transformer.TaskToolTransformer{})
	result := widget.PaginatorFind(s, &[]mysql_model.TaskTestCase{})
	response.RenderSuccess(ctx, result)
}

/**
 * apiType http
 * @api {get} /api/v1/tasks/:id/tool_task_info 获取工具任务终端信息
 * @apiVersion 1.0.0
 * @apiName GetToolTaskInfo
 * @apiGroup Task
 *
 * @apiDescription 根据任务id,获取工具任务信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}		id 	任务id
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "task_list": [
 *                 {
 *                     "cpu": "64",
 *                     "create_time": 1635567091,
 *                     "id": 4,
 *                     "last_connect_time": 1635572690,
 *                     "name": "合规子任务_1030车载合规测试任务2_1635567091280",
 *                     "os_type": "android",
 *                     "os_version": "9.0",
 *                     "status": "auto_test",
 *                     "task_uuid": "G3RD5D1Y"
 *                 }
 *             ],
 *             "tool": "hg_scanner",
 *             "tool_name": "车机检测工具"
 *         },
 *         {
 *             "task_list": [
 *                 {
 *                     "create_time": 1635567222,
 *                     "id": 4,
 *                     "name": "合规子任务_1030车载合规测试任务2_1635567091334",
 *                     "parent_id": 7,
 *                     "search_content": "合规子任务_1030车载合规测试任务2_1635567091334_G3RD5D1Y",
 *                     "status": 2,
 *                     "system_info": {
 *                         "_id": "617cc676aee3d1849aaec456",
 *                         "brand": "Android",
 *                         "car_mode": "AOSP on crosshatch",
 *                         "company": "Google",
 *                         "cpu_mode": "AArch64 Processor rev 13 (aarch64)",
 *                         "cpu_version": "64",
 *                         "platform": "android",
 *                         "sys_sdk_ver": 28,
 *                         "sys_version": "9",
 *                         "task_id": "G3RD5D1Y"
 *                     },
 *                     "task_id": "G3RD5D1Y",
 *                     "test_time": 0,
 *                     "vul_scanner_id": ""
 *                 }
 *             ],
 *             "tool": "vul_scanner",
 *             "tool_name": "车机漏扫检测工具"
 *         }
 *     ],
 *     "msg": ""
 * }
 */
func (this TaskController) GetToolTaskInfo(ctx *gin.Context) {
	id := request.ParamString(ctx, "id")
	// 查询任务信息
	task := new(mysql_model.Task)
	has, err := sys_service.NewSession().Session.Where("id=?", id).Get(task)
	if err != nil {
		response.RenderFailure(ctx, err)
	} else if !has {
		response.RenderFailure(ctx, errors.New("Task not found"))
	}
	hgTasks := []map[string]interface{}{}
	vulTasks := []map[string]interface{}{}
	// 查询任务关联的合规任务信息
	if has, task := sys_service.NewSessionWithCond(qmap.QM{"e_task_uuid": task.TaskUuid}).GetOne(new(mysql_model.HgTestTask)); has {
		hgTasks = append(hgTasks, *task)
	}
	// 查询任务关联的漏洞扫描任务信息
	if list, err := sys_service.NewSessionWithCond(qmap.QM{"e_parent_id": task.Id}).Get(&[]mysql_model.VulTask{}); err == nil {
		for _, item := range *list {
			params := qmap.QM{
				"e_task_id": task.TaskUuid,
			}
			ormSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_DEVICE_INFO, params)
			if res, err := ormSession.GetOne(); err == nil {
				item["system_info"] = res
			} else {
				item["system_info"] = qmap.QM{}
			}
			vulTasks = append(vulTasks, item)
		}
	}
	result := []qmap.QM{
		{
			"tool":      common.TOOL_HG_ANDROID_SCANNER,
			"tool_name": common.TOOL_HG_ANDROID_SCANNER_NAME,
			"task_list": hgTasks,
		}, {
			"tool":      common.TOOL_VUL_SCANNER,
			"tool_name": common.TOOL_VUL_SCANNER_NAME,
			"task_list": vulTasks,
		},
	}
	response.RenderSuccess(ctx, result)
}

// 获取任务中的漏洞
func (this TaskController) GetVul(ctx *gin.Context) {
	// todo 获取任务中的漏洞
	response.RenderSuccess(ctx, "vul")
}

// 获取任务中的测试结果
func (this TaskController) GetTestResult(ctx *gin.Context) {
	// todo 获取任务中的测试结果
	response.RenderSuccess(ctx, "result")
}

func (this TaskController) GetAllScenarios(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	s := mysql.GetSession()

	// 查询组键
	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	all, err := widget.All(s, &[]mysql_model.KnowledgeScenario{})
	results := []map[string]interface{}{}
	for _, one := range all {
		result := map[string]interface{}{}
		result["id"] = one["id"]
		result["name"] = one["name"]
		result["describe"] = one["describe"]
		if a, b := one["tag"]; b {
			result["tag"] = a
		} else {
			result["tag"] = ""
		}
		if a, b := one["tasking"]; b {
			result["tasking"] = a
		} else {
			result["tasking"] = ""
		}
		results = append(results, result)
	}

	if err != nil {
		response.RenderFailure(ctx, err)
	}
	response.RenderSuccess(ctx, results)
}

func (this TaskController) GetAllTools(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	// UserID := int(session.GetUserId(http_ctx.GetHttpCtx(ctx)))
	// UserName := session.GetUserName(ctx)
	search := req.String("search")
	categoryID := req.Int("category_id")
	queryParams := qmap.QM{
		"e_status": 1,
	}
	if search != "" {
		queryParams.Merge(map[string]interface{}{"l_search": search})
	}

	if categoryID > 0 {
		queryParams.Merge(map[string]interface{}{"e_category_id": categoryID})
	}

	mgoSession := mongo.NewMgoSession(common.MC_TOOL).AddCondition(queryParams).AddUrlQueryCondition(req.String("query_params"))
	res := mgoSession.All()
	results := []map[string]interface{}{}
	for _, one := range res {
		result := map[string]interface{}{}
		id := one["_id"]
		a := id.(bson.ObjectId)
		result["id"] = a.Hex()
		result["name"] = one["tool_name"]
		result["describe"] = one["tool_detail"]
		if a, b := one["tag"]; b {
			result["tag"] = a
		} else {
			result["tag"] = ""
		}
		if a, b := one["tasking"]; b {
			result["tasking"] = a
		} else {
			result["tasking"] = ""
		}
		results = append(results, result)
	}
	response.RenderSuccess(ctx, results)
}

// 查询任务中的测试用例,需要人工接入的部分
func (this TaskController) GetTaskCasesByAuto(ctx *gin.Context) {
	id := request.ParamString(ctx, "id")
	s := mysql.GetSession()
	s.Where("task_id=?", id)
	s.Where("action_status=?", common.CASE_STATUS_TESTING)
	s.Where("auto_test_level<>?", common.IS_TASK_CASE_AUTO)
	s.OrderBy("auto_test_level desc")
	widget := orm.PWidget{}
	widget.SetTransformer(&transformer.TaskToolTransformer{})
	result, err := widget.All(s, &[]mysql_model.TaskTestCase{})

	if err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

// 查询任务中的场景
func (this TaskController) GetTaskDemand(ctx *gin.Context) {
	id := request.ParamString(ctx, "id")
	s := mysql.GetSession()
	s.Where("task_id=?", id)

	widget := orm.PWidget{}
	result, err := widget.All(s, &[]mysql_model.TaskTestCase{})
	if err == nil {
		demandList := make(map[int]interface{}, 0)
		var count = 0
		for _, TestCase := range result {
			demandId := TestCase["demand_id"]
			tmp := demandId.(int)
			if _, ok := demandList[tmp]; !ok {
				count++
				if tmpDemand, has := mysql_model.KnowledgeDemandFindById(tmp); has {
					demandList[tmp] = tmpDemand
				} else {
					// todo 如果demand库中不存在这个id，先不处理
					tmpDemand.Id = tmp
					tmpDemand.Name = "未知"
					demandList[tmp] = tmpDemand
				}
			} else {
				continue
			}
		}
		result := make([]interface{}, 0)
		for _, tmp := range demandList {
			result = append(result, tmp)
		}
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

// 完成任务
func (this TaskController) Finish(ctx *gin.Context) {
	id := request.ParamInt(ctx, "id")
	tastStatus := request.MustInt(ctx, "status")
	info := qmap.QM{
		"status": tastStatus,
	}
	task := new(mysql_model.Task)
	if task, err := task.UpdateTaskById(id, info, int(http_ctx.GetUserId(ctx)), http_ctx.GetUserName(ctx)); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		// 添加一条报告任务记录 report_task
		reportTask := new(mysql_model.ReportTask)
		reportTask.TaskId = id
		// reportTask.Status = 0
		reportTask.ReportType = 4 // 新增的一种类型，报告任务
		reportTask.CreateTime = time.Now().Format("2006-01-02 15:04:05")
		reportTask.Name = task.Name
		nid, err := reportTask.Create()
		if err != nil {
			log.GetHttpLogLogger().Error(fmt.Sprintf("生成报告任务,id=%v,err=%v", nid, err))
		}

		// 如果任务的状态是修改已完成，那么测试用例所有的状态都修改为，测试失败
		err = new(mysql_model.TaskTestCase).SetCaseStatus(id)
		if err != nil {
			response.RenderFailure(ctx, err)
		}
		response.RenderSuccess(ctx, task)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/task_logs 分页查询任务日志
 * @apiVersion 0.1.0
 * @apiName GetLogAll
 * @apiGroup Task
 *
 * @apiDescription 分页查询任务日志
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
* {
*     "code": 0,
*     "data": {
*         "list": [
*             {
*                 "create_time": 1636013380,
*                 "id": 5,
*                 "level": 1,
*                 "message": "",
*                 "status": 2,
*                 "task_id": 13,
*                 "task_name": "12asdfadf1",
*                 "user_id": 0,
*                 "user_name": "任务监控服务"
*             }
*         ],
*         "pagination": {
*             "current_page": 1,
*             "per_page": 20,
*             "total": 5,
*             "total_pages": 1
*         }
*     },
*     "msg": ""
* }
*/
func (this TaskController) GetLogAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	s := mysql.GetSession()
	s.Where("user_id=?", http_ctx.GetUserId(ctx))
	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	widget.AddSorter(*(orm.NewSorter("id", 1)))
	all := widget.PaginatorFind(s, &[]mysql_model.TaskLog{})
	response.RenderSuccess(ctx, all)
}

/**
 * apiType http
 * @api {get} /api/v1/task/vehicle_group 分组查询任务列表
 * @apiVersion 0.1.0
 * @apiName GetTasksByGroup
 * @apiGroup Task
 *
 * @apiDescription 分组查询任务列表
 *
 * @apiUse authHeader
 *
 * @apiParam {int}		min_create_time 	查询任务起始时间(秒时间戳)
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "tasks": [
 *                 {
 *                     "asset_vehicle_id": 1,
 *                     "category": "固件扫描工具",
 *                     "client_info_time": 0,
 *                     "complete_time": 0,
 *                     "create_time": 1635501272,
 *                     "create_user_id": 592,
 *                     "describe": "222",
 *                     "firmware_template_id": 100,
 *                     "hg_client_info": "",
 *                     "hg_file_uuid": "",
 *                     "id": 2,
 *                     "is_tool_task": 1,
 *                     "last_connect_time": 0,
 *                     "last_op_id": 592,
 *                     "name": "1029固件检测任务2",
 *                     "need_connected": 1,
 *                     "piece_id": 1,
 *                     "piece_version_id": 1,
 *                     "scenario_id": 0,
 *                     "status": 2,
 *                     "task_uuid": "G3QJCK1W",
 *                     "tool": "firmware_scanner",
 *                     "tool_id": "612b3b33e6ed5c7e3a316623",
 *                     "tool_task_id": "1",
 *                     "update_time": 0,
 *                     "vulnerability_total": 51
 *                 }
 *             ],
 *             "vehicle_brand": "比亚迪",
 *             "vehicle_code": "秦",
 *             "vehicle_id": 1,
 *             "piece_id": 1,
 *             "vehicle_brand": "",
 *             "vehicle_code": ""
 *         }
 *     ],
 *     "msg": ""
 * }
 */
func (this TaskController) GetTasksByGroup(ctx *gin.Context) {
	minTaskCreateTime := request.QueryInt(ctx, "min_create_time")
	// 获取车型品牌列表
	testPieces := []*mysql_model.AssetTestPiece{}
	result := []interface{}{}
	if err := sys_service.NewSession().Session.Desc("id").Find(&testPieces); err == nil {
		for _, piece := range testPieces {
			tasks := []mysql_model.Task{}
			widget := orm.PWidget{}
			widget.AndWhereEqual("piece_id", piece.Id)
			widget.AndWhereEqual("create_user_id", http_ctx.GetUserId(ctx))
			widget.AndWhereGte("create_time", minTaskCreateTime)
			widget.SetTransformerFunc(this.getTasksByGroupTaskTransformer)
			if all, err := widget.All(mysql.GetSession().Desc("id"), &tasks); err == nil {
				if len(tasks) > 0 {
					vehicleItem := qmap.QM{
						"piece_id":      piece.Id,
						"piece_name":    piece.Name,
						"vehicle_id":    piece.AssetVehicleId,
						"vehicle_brand": "",
						"vehicle_code":  "",
						"tasks":         all,
					}
					vehicle := new(mysql_model.AssetVehicle)
					if has, err := sys_service.NewSession().Session.Where("id=?", piece.AssetVehicleId).Get(vehicle); err == nil && has {
						vehicleItem["vehicle_brand"] = vehicle.Brand
						vehicleItem["vehicle_code"] = vehicle.Code
					}
					result = append(result, vehicleItem)
				}
			} else {
				response.RenderFailure(ctx, err)
			}
		}
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this TaskController) getTasksByGroupTaskTransformer(data qmap.QM) qmap.QM {
	switch data.String("tool") {
	case "":
		// 如果是场景任务，则去查询测试用例和漏洞的统计信息
		taskId := data.Int("id")
		{
			// 查询"未测试"的测试用例数量
			if total, err := sys_service.NewSession().Session.Table(new(mysql_model.TaskTestCase)).Where("task_id=?", taskId).In("action_status", []int{common.CASE_STATUS_READY, common.CASE_STATUS_QUEUING, common.CASE_STATUS_TESTING, common.CASE_STATUS_ANALYSIS}).Count(); err == nil {
				data["untest_test_case_total"] = total
			} else {
				panic(err)
			}
		}
		{
			// 查询"测试完成"的测试用例数量
			if total, err := sys_service.NewSession().Session.Table(new(mysql_model.TaskTestCase)).Where("task_id=?", taskId).In("action_status", []int{common.CASE_STATUS_COMPLETED, common.CASE_STATUS_FAIL}).Count(); err == nil {
				data["tested_test_case_total"] = total
			} else {
				panic(err)
			}
		}
		{
			// 查询任务中的漏洞数量
			if total, err := sys_service.NewSession().Session.Table(new(mysql_model.Vulnerability)).Where("task_id=?", taskId).Count(); err == nil {
				data["vulnerability_total"] = total
			} else {
				panic(err)
			}
		}
		{
			// 如果任务"已完成",则去查询报告下载的文件id
			data["report_file_id"] = ""
			data["report_file_name"] = ""
			if data.Int("status") == common.TASK_STATUS_SUCCESS {
				reportTask := new(mysql_model.ReportTask)
				if has, err := sys_service.NewSession().Session.Where("task_id=?", taskId).Get(reportTask); err == nil {
					if has {
						data["report_file_id"] = reportTask.FileId
						data["report_file_name"] = reportTask.ReportName
					}
				} else {
					panic(err)
				}
			}
		}
	case common.TOOL_VUL_SCANNER:
		// 如果是漏洞扫描任务并且任务已经测试完成，则去查询漏洞中的漏洞统计信息
		if data.Int("status") == common.TASK_STATUS_SUCCESS {
			params := qmap.QM{
				"e_task_id": data.String("task_uuid"),
			}
			if total, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_INFO, params).Count(); err == nil {
				data["vulnerability_total"] = total
			} else {
				data["vulnerability_total"] = 0
			}
		} else {
			data["vulnerability_total"] = 0
		}
	case common.TOOL_FIRMWARE_SCANNER:
		// 如果是固件，则去查询固件报告中的漏洞统计信息
		// 如果是漏洞扫描任务并且任务已经测试完成，则去查询漏洞中的漏洞统计信息
		if data.Int("status") == common.TASK_STATUS_SUCCESS {
			// 查询对应的固件扫描任务信息
			firmwareTask := new(mysql_model.FirmwareTask)
			if has, err := sys_service.NewSession().Session.Where("task_id=?", data.Int("id")).Get(firmwareTask); err == nil {
				if has {
					if total, err := sys_service.NewSession().Session.Table(new(mysql_model.FirmwareReportRtsCve)).Where("scanner_id=?", firmwareTask.Id).Count(); err == nil {
						data["vulnerability_total"] = total
					} else {
						panic(err)
					}
				} else {
					data["vulnerability_total"] = 0
				}
			} else {
				panic(err)
			}
		} else {
			data["vulnerability_total"] = 0
		}
	}
	return data
}

// 获取任务中的测试用例，用例的状态数据
func (this TaskController) GetTaskCasesStatus(ctx *gin.Context) {
	id := request.ParamString(ctx, "id")
	s := mysql.GetSession()
	s.Where("task_id=?", id)
	s.And("action_status<>?", common.CASE_STATUS_COMPLETED)
	s.And("action_status<>?", common.CASE_STATUS_FAIL)
	testingCount, err := s.FindAndCount(&[]mysql_model.TaskTestCase{})
	result := qmap.QM{
		"testing": testingCount,
	}
	if err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

// 获取任务中的测试用例，用例的状态数据
func (this TaskController) GetCategory(ctx *gin.Context) {
	UserId := int(http_ctx.GetUserId(ctx))
	lists := mysql_model.GetTaskCategoryList(UserId)
	response.RenderSuccess(ctx, lists)
}

/**
 * apiType http
 * @api {get} /api/v1/task/long_connection_scanners 查询任务下长连接扫描任务列表
 * @apiVersion 0.1.0
 * @apiName GetTasksLongConnectionScanner
 * @apiGroup Task
 *
 * @apiDescription 查询任务下长连接扫描任务列表
 *
 * @apiUse authHeader
 *
 * @apiParam {int}		task_id 	任务id
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *			"hg_scanner"
 *		]
 *     "msg": ""
 * }
 */
func (this TaskController) GetTasksLongConnectionScanner(ctx *gin.Context) {
	scannerList := []string{}
	// 查询是否存在安卓合规检测任务
	{
		session := sys_service.NewSession().Session.Where("task_id=?", request.QueryInt(ctx, "task_id"))
		session.And("test_tool=?", common.TOOL_HG_ANDROID_SCANNER)
		session.And("action_status <>?", common.CASE_STATUS_INVALID)
		if has, err := session.Get(new(mysql_model.TaskTestCase)); err == nil {
			if has {
				scannerList = append(scannerList, common.TOOL_HG_ANDROID_SCANNER)
			}
		} else {
			panic(err)
		}
	}
	// 查询是否存在Linux合规检测任务
	response.RenderSuccess(ctx, qmap.QM{"data": scannerList})
}
