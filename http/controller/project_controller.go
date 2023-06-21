package controller

import (
	"context"
	"errors"
	"strings"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/orm_mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/common_lib/session"
	"skygo_detection/logic"
	"skygo_detection/mongo_model"
)

type ProjectController struct{}

/**
 * apiType http
 * @api {get} /api/v1/projects 项目列表
 * @apiVersion 1.0.0
 * @apiName GetAll
 * @apiGroup Project
 *
 * @apiDescription 查询车机列表接口
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/projects
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "list": [
 *             {
 *                 "_id": "5e61f95024b64748d37d8cc6",
 *                 "company": "项目所属车厂",
 *                 "create_time": 1583479120091,
 *                 "description": "项目描述",
 *                 "end_time": 1583478184,
 *                 "evaluate_requirement": "",
 *                 "manager_id": 1,
 *                 "member_ids": [],
 *                 "name": "项目名称",
 *                 "start_time": 1583478180,
 *                 "update_time": 1583479120091
 *             }
 *         ],
 *         "pagination": {
 *             "count": 5,
 *             "current_page": 1,
 *             "per_page": 20,
 *             "total": 5,
 *             "total_pages": 1
 *         }
 *     }
 * }
 */
func (this ProjectController) GetAll(ctx *gin.Context) {
	params := qmap.QM{
		"e_all_users": session.GetUserId(http_ctx.GetHttpCtx(ctx)),
	}

	widget := orm_mongo.NewWidgetWithParams(common.MC_PROJECT, params).SetQueryStr(ctx.Request.URL.RawQuery)
	widget.SetTransformerFunc(this.ProjectTransformer)
	if res, err := widget.PaginatorFind(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/projects/:id 查询某一个项目信息
 * @apiVersion 1.0.0
 * @apiName GetOne
 * @apiGroup Project
 *
 * @apiDescription 查询某一个项目信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}       id        项目id
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/projects/:id
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "company": "项目所属车厂",
 *         "create_time": 1583479120091,
 *         "description": "项目描述",
 *         "end_time": 1583478184,
 *         "evaluate_requirement": "",
 *         "id": "5e61f95024b64748d37d8cc6",
 *         "manager_id": 1,
 *         "member_ids": null,
 *         "name": "项目名称",
 *         "start_time": 1583478180,
 *         "update_time": 1583479120091
 *     }
 * }
 */
func (this ProjectController) GetOne(ctx *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(request.ParamString(ctx, "id"))
	params := qmap.QM{
		"e__id":        id,
		"e_all_users":  session.GetUserId(http_ctx.GetHttpCtx(ctx)),
		"e_is_deleted": common.STATUS_DELETE,
	}

	widget := orm_mongo.NewWidgetWithCollectionName(common.MC_PROJECT).SetQueryStr(ctx.Request.URL.RawQuery).SetParams(params)
	widget.SetTransformerFunc(this.ProjectTransformer)
	if res, err := widget.Get(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this ProjectController) ProjectTransformer(data qmap.QM) qmap.QM {
	_id, _ := primitive.ObjectIDFromHex(data.MustString("company"))
	params := qmap.QM{
		"e__id": _id,
	}
	if factory, err := orm_mongo.NewWidgetWithCollectionName(common.MC_FACTORY).SetParams(params).Get(); err == nil {
		data["company_name"] = factory.String("name")
	} else {
		data["company_name"] = ""
	}

	id := data["_id"].(primitive.ObjectID)

	// 获取任务数
	taskParams := qmap.QM{
		"e_project_id": id.Hex(),
	}
	if total, err := orm_mongo.NewWidgetWithCollectionName(common.MC_EVALUATE_TASK).SetParams(taskParams).Count(); err == nil {
		data["task_total"] = total
	}

	// 获取测试用例数
	itemParams := qmap.QM{
		"e_project_id": id.Hex(),
	}
	if total, err := orm_mongo.NewWidgetWithCollectionName(common.MC_EVALUATE_ITEM).SetParams(itemParams).Count(); err == nil {
		data["case_total"] = total
	}

	// 获取漏洞数
	params2 := qmap.QM{
		"e_project_id": id.Hex(),
	}

	if total, err := orm_mongo.NewWidgetWithCollectionName(common.MC_EVALUATE_VULNERABILITY).SetParams(params2).Count(); err == nil {
		data["vulnerability_total"] = total
	} else {
		panic(err)
	}
	return data
}

/**
 * apiType http
 * @api {post} /api/v1/projects 创建项目
 * @apiVersion 1.0.0
 * @apiName Create
 * @apiGroup Project
 *
 * @apiDescription 创建新项目接口
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           name                    项目名称
 * @apiParam {int}              start_time              项目开始时间
 * @apiParam {int}              end_time                项目结束时间
 * @apiParam {string}           company                 项目所属车厂
 * @apiParam {string}           description             项目描述
 * @apiParam {string}           evaluate_requirement    项目要求
 * @apiParam {int}              manager_id              项目经理id
 * @apiParam {int}              amount                  项目金额（单位:万元）
 * @apiParam {int}              [status]                项目状态(-1:异常,0:未开始,1:初测中,2:初测完成,3:复测开始,4:复测完成,项目结束:5)
 * @apiParam {json}             evaluate_targets        项目评估对象信息
 * @apiParam {[]string]}        contracts               项目合同文件
 * @apiParam {[]string]}        biddings                项目标书
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "name":"项目名称",
 *          "start_time":1583478180,
 *          "end_time": 1583478184,
 *          "company":"项目所属车厂",
 *          "description":"项目描述",
 *          "member_ids":[2,3],
 *          "amount":1000,
 *          "code_name": "改款2020"
 *          "brand":"长安"
 *          "contracts":["12313sadf547asd121", "1231cxvaoi04Z1"]
 *          "biddings":["12313sadf547asd121", "1231cxvaoi04Z1"]
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "project_id": "5fcb713ae138231067f457d9"   //项目ID"
 *     }
 * }
 */
func (this ProjectController) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	// 创建session
	sess, err := orm_mongo.GetMongoClient().StartSession()
	if err != nil {
		panic(err)
	}
	defer sess.EndSession(context.TODO())
	sessCtx := mongo.NewSessionContext(context.TODO(), sess)
	// 启动事务
	if err = sess.StartTransaction(); err != nil {
		panic(err)
	}

	var errHere error
	if project, err := new(mongo_model.Project).Create(sessCtx, req, session.GetUserId(http_ctx.GetHttpCtx(ctx))); err != nil {
		errHere = err
		goto AbortTransaction
	} else {
		// 添加项目评估对象
		if evaluateTargets, has := req.TrySlice("evaluate_targets"); has && len(evaluateTargets) > 0 {
			for _, item := range evaluateTargets {
				if _, errHere = new(mongo_model.EvaluateAsset).Create(sessCtx, project.Id.Hex(), session.GetUserId(http_ctx.GetHttpCtx(ctx)), item.(map[string]interface{})); errHere != nil {
					goto AbortTransaction
				}
			}
		}
		// 创建`项目标书`文档
		biddings, err := new(mongo_model.PMFile).Create(sessCtx, project.Id.Hex(), "", "项目标书", mongo_model.FILE_TYPE_DIR, "", 0, int(session.GetUserId(http_ctx.GetHttpCtx(ctx))))
		if err != nil {
			errHere = err
			goto AbortTransaction
		}
		biddingFiles := req.SliceString("biddings")
		for _, item := range biddingFiles {
			_id, _ := primitive.ObjectIDFromHex(item)
			if fi, err := orm_mongo.GridFSOpenId(common.MC_File, _id); err == nil {
				if _, err := new(mongo_model.PMFile).Create(sessCtx, project.Id.Hex(), item, fi.GetFile().Name, mongo_model.FILE_TYPE_DOC, biddings.Id.Hex(), int(fi.GetFile().Length), int(session.GetUserId(http_ctx.GetHttpCtx(ctx)))); err != nil {
					errHere = err
					goto AbortTransaction
				}
				fi.Close()
			} else {
				errHere = err
				goto AbortTransaction
			}
		}
		// 创建`项目合同`文档
		contracts, _ := new(mongo_model.PMFile).Create(sessCtx, project.Id.Hex(), "", "项目合同", mongo_model.FILE_TYPE_DIR, "", 0, int(session.GetUserId(http_ctx.GetHttpCtx(ctx))))
		contractFiles := req.SliceString("contracts")
		for _, item := range contractFiles {
			_id, _ := primitive.ObjectIDFromHex(item)
			if fi, err := orm_mongo.GridFSOpenId(common.MC_File, _id); err == nil {
				if _, err := new(mongo_model.PMFile).Create(sessCtx, project.Id.Hex(), item, fi.GetFile().Name, mongo_model.FILE_TYPE_DOC, contracts.Id.Hex(), int(fi.GetFile().Length), int(session.GetUserId(http_ctx.GetHttpCtx(ctx)))); err != nil {
					errHere = err
					goto AbortTransaction
				}
				fi.Close()
			} else {
				errHere = err
				goto AbortTransaction
			}
		}

		if err = sess.CommitTransaction(context.Background()); err != nil {
			panic(err)
		} else {
			response.RenderSuccess(ctx, map[string]string{"project_id": project.Id.Hex()})
			return
		}
	}

AbortTransaction:
	// 结束事务
	_ = sess.AbortTransaction(context.Background())
	response.RenderFailure(ctx, errHere)
	return
}

/**
 * apiType http
 * @api {put} /api/v1/projects/:id 更新项目
 * @apiVersion 1.0.0
 * @apiName Update
 * @apiGroup Project
 *
 * @apiDescription 更新项目接口
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           id                      项目id
 * @apiParam {string}           name    				项目名称
 * @apiParam {int}           	start_time           	项目开始时间
 * @apiParam {int}           	end_time           		项目结束时间
 * @apiParam {string}           company           		项目所属车厂
 * @apiParam {string}           description        	 	项目描述
 * @apiParam {string}           evaluate_requirement 	项目要求
 * @apiParam {int}          	manager_id         		项目经理id
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X PUT http://localhost/api/v1/projects/:id
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "name":"项目名称",
 *          "start_time":1583478180,
 *          "end_time": 1583478184,
 *          "company":"项目所属车厂",
 *          "description":"项目描述",
 *          "member_ids":[2,3],
 *          "brand":"长安",
 *          "code_name":"改款2020",
 *          "evaluate_requirement":"开发项目1的项目要求"，
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *        "result": true
 *     }
 * }
 */
func (this ProjectController) Update(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	id := request.ParamString(ctx, "id")
	if err := mongo_model.CheckIsProjectManager(id, session.GetUserId(http_ctx.GetHttpCtx(ctx))); err != nil {
		response.RenderFailure(ctx, err)
		return
	}

	updateCols := map[string]string{
		"name":                 "string",
		"company":              "string",
		"brand":                "string",
		"code_name":            "string",
		"start_time":           "int",
		"end_time":             "int",
		"status":               "int",
		"evaluate_requirement": "string",
		"description":          "string",
		"manager_id":           "int",
		"member_ids":           "interface",
		"amount":               "int",
	}
	rawInfo := custom_util.CopyMapColumns(*req, updateCols)
	if _, err := new(mongo_model.Project).Update(id, rawInfo); err == nil {
		response.RenderSuccess(ctx, map[string]bool{"result": true})
		return
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/projects 批量删除项目
 * @apiVersion 1.0.0
 * @apiName BulkDelete
 * @apiGroup Project
 *
 * @apiDescription 批量删除项目
 *
 * @apiUse authHeader
 *
 * @apiParam {string]}      ids      模块id,多个模块id用"\\|"链接(如:"1\\|2\\|3")
 *
 * @apiSuccessExample {json} 请求成功示例:
 *       {
 *            "code": 0,
 *			  "data":{
 *				"number":1
 *			}
 *       }
 */
func (this ProjectController) BulkDelete(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	ids := strings.Split(req.MustString("ids"), "|")
	rawInfo := qmap.QM{"is_deleted": common.PSD_DELETE}
	effectNum := 0
	for _, id := range ids {
		_id, _ := primitive.ObjectIDFromHex(id)
		info, err := orm_mongo.NewWidgetWithCollectionName(common.MC_PROJECT).SetParams(qmap.QM{"e__id": _id}).Get()
		if err != nil {
			err := errors.New("Project not found")
			response.RenderFailure(ctx, err)
			return
		}
		// 项目经理不是当前用户不能删除
		if info["manager_id"] != int(session.GetUserId(http_ctx.GetHttpCtx(ctx))) {
			err := errors.New("您没有操作该项目的权限")
			response.RenderFailure(ctx, err)
			return
		}

		if info["is_deleted"] == common.PSD_DELETE {
			err := errors.New("该项目已被删除")
			response.RenderFailure(ctx, err)
			return
		}
		// 项目创建了任务的情况下，不可以删除
		_, err = orm_mongo.NewWidgetWithCollectionName(common.MC_EVALUATE_TASK).SetParams(qmap.QM{"e_project_id": id}).Get()
		if err == nil {
			err := errors.New("该项目已创建任务，不可删除")
			response.RenderFailure(ctx, err)
			return
		}

		if _, err := new(mongo_model.Project).Update(id, rawInfo); err == nil {
			effectNum++
			// 软删除项目，同时软删除项目任务
			selector := bson.M{"project_id": id}
			updateItem := bson.M{
				"$set": qmap.QM{
					"is_deleted": common.PSD_DELETE,
				},
			}
			_, err := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_TASK).UpdateMany(context.TODO(), selector, updateItem)
			if err != nil {
				response.RenderFailure(ctx, err)
				return
			}
		} else {
			response.RenderFailure(ctx, err)
			return
		}
	}

	response.RenderSuccess(ctx, qmap.QM{"number": effectNum})
	return
}

/**
 * apiType http
 * @api {get} /api/v1/project/configs 获取配置参数
 * @apiVersion 1.0.0
 * @apiName GetConfigs
 * @apiGroup Project
 *
 * @apiDescription 获取配置参数接口
 *
 * @apiUse authHeader
 *
 * @apiExample {curl} 请求示例:
 * curl -i  http://localhost/api/v1/project/configs
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":{
 *				"companies": [
 *					"车厂1",
 *					"车厂2"
 *				],
 *			}
 *      }
 */
func (this ProjectController) GetConfigs(ctx *gin.Context) {
	result := qmap.QM{}

	// 查询项目所有项目经理
	if rsp, err := new(logic.AuthLogic).GetSpecifiedServiceUsers(common.PM_SERVICE, common.ROLE_PM); err == nil {
		result["users"] = rsp
	} else {
		response.RenderFailure(ctx, err)
		return
	}

	// 查询项目所有测试人员
	if rsp, err := new(logic.AuthLogic).GetSpecifiedServiceUsers(common.PM_SERVICE, common.ROLE_TEST); err == nil {
		result["testers"] = rsp
	} else {
		response.RenderFailure(ctx, err)
		return
	}

	result["current_role_id"] = session.GetRoleId(http_ctx.GetHttpCtx(ctx))

	// 查询所有车厂
	widget := orm_mongo.NewWidgetWithCollectionName(common.MC_FACTORY)
	widget.SetLimit(10000)
	factories, factoryErr := widget.Find()
	custom_util.CheckErr(factoryErr)
	result["companies"] = factories
	// 查询所有测试类型
	if res, err := orm_mongo.NewWidgetWithCollectionName(common.MC_EVALUATE_TYPE).Find(); err == nil {
		result["evaluate_type"] = res
	} else {
		response.RenderFailure(ctx, err)
		return
	}
	// 项目状态
	result["project_status"] = []qmap.QM{
		{"code": 0, "name": "创建"},
		{"code": 1, "name": "测试中"},
		{"code": 9, "name": "项目完成"},
	}
	// 测试用例状态
	result["evaluate_item_status"] = []qmap.QM{
		{"code": 0, "name": "待初测"},
		{"code": 1, "name": "测试完成"},
		{"code": 2, "name": "待补充"},
		{"code": 3, "name": "审核通过"},
	}
	// 评估漏洞状态
	result["evaluate_vul_status"] = []qmap.QM{
		{"code": 0, "name": "未修复"},
		{"code": 1, "name": "已修复"},
		{"code": 2, "name": "重打开"},
	}
	// 评估漏洞状态
	result["asset_status"] = []qmap.QM{
		{"code": 0, "name": "使用中"},
		{"code": 1, "name": "已归还"},
	}
	// 报告类型
	result["report_type"] = []qmap.QM{
		// {"type": "week", "name": "周报"},
		{"type": "test", "name": "初测报告"},
		{"type": "retest", "name": "复测报告"},
	}
	// 漏洞类型
	if vulTypes, err := new(mongo_model.EvaluateVulType).GetAll(); err == nil {
		result["vul_type"] = vulTypes
	} else {
		result["vul_type"] = []qmap.QM{}
	}

	response.RenderSuccess(ctx, result)
	return
}

//@auto_generated_api_end
/**
 * apiType http
 * @api {get} /api/v1/project/dashboard 项目DashBoard信息
 * @apiVersion 1.0.0
 * @apiName DashBoard
 * @apiGroup Project
 *
 * @apiDescription 项目DashBoard信息
 *
 * @apiUse authHeader
 *
 * @apiExample {curl} 请求示例:
 * curl -i  http://localhost/api/v1/project/dashboard
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":{
 *				"companies": [
 *					"车厂1",
 *					"车厂2"
 *				],
 *			}
 *      }
 */
func (this ProjectController) DashBoard(ctx *gin.Context) {
	uid := int(session.GetUserId(http_ctx.GetHttpCtx(ctx)))

	myProjectCount := 0
	allProjectCount := 0
	allRunningProjectCount := 0
	allFinishProjectCount := 0
	allAbnormalProjectCount := 0

	myItemCount := 0
	allItemCount := 0

	myVulnerabilityCount := 0
	allVulnerabilityCount := 0
	var err error

	// 项目
	projectIds := []string{}
	myProjectIds := []string{}
	{
		coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_PROJECT)
		// 查询语句举例 db.project.find({$or : [{"member_ids" : 3}, {"manager_id" : 5}]})
		// member_ids是一个数组存多个用户id，我们采用模糊匹配， manager_id
		filter := bson.M{
			"$or": []bson.M{
				bson.M{"member_ids": uid},
				bson.M{"manager_id": uid},
			},
			"status": bson.M{"$ne": -1},
		}
		myResult := []map[string]interface{}{}

		// 查询
		c, _ := coll.Find(context.TODO(), filter)
		c.All(context.TODO(), &myResult)

		for _, m := range myResult {
			id := m["_id"].(primitive.ObjectID).Hex()
			myProjectIds = append(myProjectIds, id)
		}
		myProjectCount = len(myResult)
		// myProjectCount, err = mgoSession.Session.Find(filter).Select(bson.M{"status": 1}).Count()
		// if err != nil {
		//	panic(err)
		// }

		result := []qmap.QM{}

		c, _ = coll.Find(context.TODO(), bson.M{"status": bson.M{"$ne": -1}})
		err := c.All(context.TODO(), &result)
		if err != nil {
			panic(err)
		}

		for _, m := range result {
			id := m["_id"].(primitive.ObjectID).Hex()
			projectIds = append(projectIds, id)

			if isDelete, has := m.TryInt("is_deleted"); has {
				// 只要状态不是删除就统计到总项目数
				if isDelete != common.PSD_DELETE {
					allProjectCount++
				}
			}

			if _v, has := m["status"]; has {
				v := _v.(int)

				list := []int{
					common.PS_NEW,
					common.PS_TEST,
					common.PS_COMPLETE,
				}
				if custom_util.InIntSlice(v, list) == true {
					allRunningProjectCount++
				}

				if v == common.PS_COMPLETE {
					allFinishProjectCount++
				}
				if v == common.PS_ABNORMAL {
					allAbnormalProjectCount++
				}
			}

		}
	}

	{
		coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_ITEM)
		// 查询语句举例 db.project.find({$or : [{"member_ids" : 3}, {"manager_id" : 5}]})
		// member_ids是一个数组存多个用户id，我们采用模糊匹配， manager_id
		// filter := bson.M{
		//	"$or": []bson.M{
		//		bson.M{"last_update_op_id": uid},
		//		bson.M{"op_id": uid},
		//	},
		// }
		// result := []map[string]interface{}{}
		// err = mgoSession.Session.Find(filter).All(&result)
		// if err != nil {
		//	panic(err)
		// }
		// myItemCount = len(result)
		myResult := []map[string]interface{}{}
		filter := []bson.M{}
		for _, projectId := range myProjectIds {
			tmp := bson.M{"project_id": projectId}
			filter = append(filter, tmp)
		}
		if len(filter) != 0 {

		}
		if len(filter) != 0 {
			c, _ := coll.Find(context.TODO(), bson.M{"$or": filter})
			err = c.All(context.TODO(), &myResult)
			if err != nil {
				panic(err)
			}
		}
		myItemCount = len(myResult)

		// 通过未删除的项目，筛选测试项
		result := []map[string]interface{}{}
		query := []bson.M{}
		for _, projectId := range projectIds {
			tmp := bson.M{"project_id": projectId}
			query = append(query, tmp)
		}
		filter1 := bson.M{
			"$or": query,
		}
		if len(query) != 0 {
			c, _ := coll.Find(context.TODO(), filter1)
			err = c.All(context.TODO(), &myResult)
			if err != nil {
				panic(err)
			}
		}
		allItemCount = len(result)
	}

	{
		// 查询漏洞数量
		coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_VULNERABILITY)
		myResult := []map[string]interface{}{}
		filter := []bson.M{}
		for _, projectId := range myProjectIds {
			tmp := bson.M{"project_id": projectId}
			filter = append(filter, tmp)
		}
		if len(filter) != 0 {
			c, _ := coll.Find(context.TODO(), bson.M{"$or": filter})
			err = c.All(context.TODO(), &myResult)
			if err != nil {
				panic(err)
			}
		}
		myVulnerabilityCount = len(myResult)

		// 通过未删除的项目，筛选测试项
		result := []map[string]interface{}{}
		query := []bson.M{}
		for _, projectId := range projectIds {
			tmp := bson.M{"project_id": projectId}
			query = append(query, tmp)
		}
		filter1 := bson.M{
			"$or": query,
		}
		if len(query) != 0 {
			c, _ := coll.Find(context.TODO(), filter1)
			err = c.All(context.TODO(), &myResult)
			if err != nil {
				panic(err)
			}
		}
		allVulnerabilityCount = len(result)
	}

	r := qmap.QM{
		"my_project_count":           myProjectCount,
		"all_project_count":          allProjectCount,
		"all_abnormal_project_count": allAbnormalProjectCount,
		"all_running_project_count":  allRunningProjectCount,
		"all_finish_project_count":   allFinishProjectCount,
		"my_item_count":              myItemCount,
		"all_item_count":             allItemCount,
		"my_vulnerability_count":     myVulnerabilityCount,
		"all_vulnerability_count":    allVulnerabilityCount,
	}

	response.RenderSuccess(ctx, r)
}

/**
 * apiType http
 * @api {get} /api/v1/project/select_list_project_asset 根据用户uid，得到下拉列表，当前用户所属项目列表以及每个项目关联的资产
 * @apiVersion 1.0.0
 * @apiName SelectListProjectAsset
 * @apiGroup Project
 *
 * @apiDescription 根据用户uid，得到下拉列表，当前用户所属项目列表以及每个项目关联的资产
 *
 * @apiUse authHeader
 *
 * @apiExample {curl} 请求示例:
 * curl -i  http://localhost/api/v1/project/select_list_project_asset
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "asset_list": [
 *                 {
 *                     "asset_id": "4b5131514252444b",
 *                     "name": "sww2-资产1"
 *                 }
 *             ],
 *             "name": "sww2_a1_测试项目",
 *             "project_id": "60cc071de830c66f2e18ad34"
 *         },
 *         {
 *             "asset_list": [
 *                 {
 *                     "asset_id": "4b51305245563231",
 *                     "name": "sww的资产"
 *                 }
 *             ],
 *             "name": "sww_a5_漏洞检测",
 *             "project_id": "60cb2387e830c66f2e18ad2c"
 *         }
 *     ]
 * }
 */
func (this ProjectController) SelectListProjectAsset(ctx *gin.Context) {
	uid := int(session.GetUserId(http_ctx.GetHttpCtx(ctx)))
	selectList := new(mongo_model.Project).SelectListProjectAsset(uid)

	data := qmap.QM{
		"data": selectList,
	}

	response.RenderSuccess(ctx, data)
}

/**
 * apiType rpc
 * @api {post} /project_manage.Project/GetMyProjectSummaryInfo 当前用户管理项目概要信息查询接口
 * @apiVersion 0.1.0
 * @apiName GetMyProjectSummaryInfo
 * @apiGroup Project
 *
 * @apiDescription 当前用户管理项目概要信息查询接口
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *      "code": 0,
 *		"data":{
 *			"complete_project_num": 5,//已结束项目总数
 *			"create_project_num": 38, //创建项目总数
 *			"create_task_num": 0, //待分配任务总数
 *			"project_item_num": 1326,//测试用例总数
 *			"project_member_num": 0, //项目成员总数
 *			"project_tester_num": 1, //测试人员总数
 *			"project_vehicle_model_num": 0, //车型总数
 *			"project_vul_num": 47, //漏洞总数
 *			"testing_project_num": 25, //进行中项目总数
 *			"testing_task_num": 0 //执行中任务总数
 *		}
 * }
 */
func (this ProjectController) GetMyProjectSummaryInfo(ctx *gin.Context) {
	result := qmap.QM{}
	newProjects := new(mongo_model.Project).GetManageProjects(session.GetUserId(http_ctx.GetHttpCtx(ctx)), common.PS_NEW, -1)
	testingProjects := new(mongo_model.Project).GetManageProjects(session.GetUserId(http_ctx.GetHttpCtx(ctx)), common.PS_TEST, -1)
	projects := append(newProjects, testingProjects...)
	// 创建中项目数量
	result["create_project_num"] = len(newProjects)
	// 进行中项目数量
	result["testing_project_num"] = len(testingProjects)
	// 已结束项目数量
	result["complete_project_num"] = len(new(mongo_model.Project).GetManageProjects(session.GetUserId(http_ctx.GetHttpCtx(ctx)), common.PS_COMPLETE, -1))
	// 待分配任务数量
	result["create_task_num"] = new(mongo_model.EvaluateTask).CountProjectTask(projects, []string{common.PTS_SIGN_CREATE, common.PTS_SIGN_TASK_AUDIT})
	// 执行中任务数量
	result["testing_task_num"] = new(mongo_model.EvaluateTask).CountProjectTask(testingProjects, []string{common.PTS_SIGN_TEST, common.PTS_SIGN_REPORT_AUDIT})
	// 已结束任务数量
	result["complete_task_num"] = new(mongo_model.EvaluateTask).CountProjectTask(testingProjects, []string{common.PTS_SIGN_FINISH})
	// 项目人员数量
	result["project_member_num"] = new(mongo_model.Project).CountProjectsMembers(session.GetUserId(http_ctx.GetHttpCtx(ctx)), []int{common.PS_NEW, common.PS_TEST})
	// 测试人员
	result["project_tester_num"] = new(mongo_model.EvaluateTask).CountProjectsTaskTester(testingProjects)
	// 项目漏洞总数
	result["project_vul_num"] = new(mongo_model.EvaluateVulnerability).CountProjectVuls(testingProjects)
	// 项目车型数
	result["project_vehicle_model_num"] = new(mongo_model.Project).CountProjectsVehicleCode(session.GetUserId(http_ctx.GetHttpCtx(ctx)), []int{common.PS_NEW, common.PS_TEST})
	// 项目测试用例总数
	result["project_item_num"] = new(mongo_model.EvaluateItem).CountItem(testingProjects)

	response.RenderSuccess(ctx, result)
}

/**
 * apiType rpc
 * @api {post} /project_manage.Project/GetMyProjectList 获取当前用户管理项目列表
 * @apiVersion 0.1.0
 * @apiName GetMyProjectList
 * @apiGroup Project
 *
 * @apiDescription 获取当前用户管理项目列表
 *
 * @apiSuccessExample {json} 请求成功示例:
 *	{
 *      "code": 0,
 *		"data": [
 *			{
 *				"id": "60d05070e830c66f2e18adff",
 *				"name": "伟伟的大项目",
 *				"process": 9
 *			},
 *			{
 *				"id": "60d03936e830c66f2e18ade8",
 *				"name": "sww5_领克100_漏洞扫描任务",
 *				"process": 87
 *			}
 *		]
 *	}
 */
func (this ProjectController) GetMyProjectList(ctx *gin.Context) {
	list := new(mongo_model.Project).GetManageProjectInfo(session.GetUserId(http_ctx.GetHttpCtx(ctx)), []int{common.PS_TEST})
	response.RenderSuccess(ctx, list)
}

/**
 * apiType rpc
 * @api {post} /project_manage.Project/GetProjectSummaryInfo 获取某一项目概要信息
 * @apiVersion 0.1.0
 * @apiName GetProjectSummaryInfo
 * @apiGroup Project
 *
 * @apiDescription 获取某一项目概要信息
 *
 * @apiParam {string}		id		项目id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *     "id" : "5fe464b1f98f923e40e8dd5f"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *	{
 *      "code": 0,
 *		"data": {
 *			"asset_vul": [ //资产漏洞分布
 *				{
 *					"asset_id": "KQ6DDYW2",
 *					"asset_name": "非法",
 *					"number": 30
 *				},
 *				{
 *					"asset_id": "KQ6DCTRU",
 *					"asset_name": "大资产",
 *					"number": 73
 *				}
 *			],
 *			"test_case": [ //测试用例完成情况
 *				{
 *					"name": "已测试",
 *					"number": 1
 *				},
 *				{
 *					"name": "测试中",
 *					"number": 0
 *				},
 *				{
 *					"name": "未测试",
 *					"number": 102
 *				}
 *			],
 *			"vul_repair": [ //漏洞修复情况
 *				{
 *					"name": "已修复",
 *					"number": 2
 *				},
 *				{
 *					"name": "未修复",
 *					"number": 101
 *				}
 *			]
 *		}
 *	}
 */
func (this ProjectController) GetProjectSummaryInfo(ctx *gin.Context) {
	projectId := request.QueryString(ctx, "id")
	result := qmap.QM{}
	// 测试用例完成情况
	result["test_case"] = new(mongo_model.EvaluateItem).StatisticItemInfo(projectId)
	// 资产漏洞分布情况
	result["asset_vul"] = new(mongo_model.EvaluateVulnerability).StatisticAssetVul(projectId)
	// 项目漏洞修复情况
	result["vul_repair"] = new(mongo_model.EvaluateVulnerability).StatisticProjectVulPatch(projectId)

	response.RenderSuccess(ctx, result)
}

/**
 * apiType rpc
 * @api {post} /project_manage.Project/GetProjectTaskSeries 获取某一项目的任务时间序列信息
 * @apiVersion 0.1.0
 * @apiName GetProjectTaskSeries
 * @apiGroup Project
 *
 * @apiDescription 获取某一项目的任务时间序列信息
 *
 * @apiParam {string}		id		项目id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *     "id" : "5fe464b1f98f923e40e8dd5f"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "complete": [ //项目任务已结束
 *             1,
 *             1,
 *             1
 *         ],
 *         "create": [ //项目任务未分配
 *             102,
 *             102,
 *             102
 *         ],
 *         "deadline": "2021-06-21 16:40:16", //项目deadline
 *         "testing": [ //项目任务创建中
 *             0,
 *             0,
 *             0
 *         ],
 *         "x_axis": [ //项目任务x时间轴
 *             "2021-06-25 01:55:00",
 *             "2021-06-25 01:55:10",
 *             "2021-06-25 01:55:20"
 *         ]
 *     }
 * }
 */
func (this ProjectController) GetProjectTaskSeries(ctx *gin.Context) {
	projectId := request.QueryString(ctx, "id")
	result := new(mongo_model.ProjectTaskInfo).GetProjectTaskSeries(projectId)
	_id, _ := primitive.ObjectIDFromHex(projectId)
	params := qmap.QM{
		"e__id": _id,
	}

	if project, err := orm_mongo.NewWidgetWithCollectionName(common.MC_PROJECT).SetParams(params).Get(); err == nil {
		result["deadline"] = custom_util.TimestampToString(project.Int64("create_time") / 1000)
	} else {
		result["deadline"] = ""
	}

	response.RenderSuccess(ctx, result)
}

/**
 * apiType rpc
 * @api {post} /project_manage.Project/GetMyBacklogSummary 获取当前用户的待办事项信息
 * @apiVersion 0.1.0
 * @apiName GetMyBacklogSummary
 * @apiGroup Project
 *
 * @apiDescription 获取当前用户的待办事项信息
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "backlog_task_number": 0, //待办任务数
 *         "need_audit_report_number": 0,//待审核报告数
 *         "week_published_report_number": 0, //本周发布报告数
 *         "week_complete_task_number": 0 //本周完成任务数
 *     }
 * }
 */
func (this ProjectController) GetMyBacklogSummary(ctx *gin.Context) {
	weekStartTime := custom_util.GetThisWeekStartTime().UnixNano() / 1000000
	projectIds := new(mongo_model.Project).GetManageProjects(session.GetUserId(http_ctx.GetHttpCtx(ctx)), common.PS_TEST, -1)
	completeProjects := new(mongo_model.Project).GetManageProjects(session.GetUserId(http_ctx.GetHttpCtx(ctx)), common.PS_TEST, weekStartTime)
	projectIds = append(projectIds, completeProjects...)
	result := qmap.QM{}
	// 待办任务数量
	result["backlog_task_number"] = new(mongo_model.EvaluateTask).CountBacklogTask(session.GetUserId(http_ctx.GetHttpCtx(ctx)))
	// 本周完成任务数量
	result["week_complete_task_number"] = new(mongo_model.EvaluateTask).CountCompleteTask(projectIds, weekStartTime)
	// 待审核报告数量
	result["need_audit_report_number"] = new(mongo_model.Report).CountOperateReport(session.GetUserId(http_ctx.GetHttpCtx(ctx)))
	// 本周发布的报告数量
	result["week_published_report_number"] = new(mongo_model.Report).CountPublishedReport(projectIds, weekStartTime)

	response.RenderSuccess(ctx, result)
}
