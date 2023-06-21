package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/common_lib/session"
	"skygo_detection/mongo_model"
)

type ToolTaskController struct{}

//@auto_generated_api_begin
/**
 * apiType http
 * @api {get} /api/v1/tool/task/list 工具任务列表
 * @apiVersion 0.1.0
 * @apiName TaskList
 * @apiGroup ToolTask
 *
 * @apiDescription 查看工具任务列表
 *
 * @apiParam {string}      			pro_task_id      					项目任务ID
 * @apiParam {string}      			search      					项目任务ID
 * @apiParam {string}      			tool_category_name      					工具任务名称
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *    "pro_task_id":"adshfkjasdgkhsg",
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *   {
 *       "code": 0,
 *       "data": {
 *           "list": [
 *               {
 *                   "_id": "5fe1eaf589e90f9b4eaa0bd6",
 *                   "assets_cate_name": "TSP",
 *                   "assets_id": "XFE20201223",
 *                   "create_time": 1608641269,
 *                   "create_user_id": 130,
 *                   "create_user_name": "zhupenghui",
 *                   "pro_id": "5fd97feef98f923b505be841",
 *                   "status": 2,
 *                   "task_name": "task_20201222204749",
 *                   "task_number": 20201222204749,
 *                   "tool_category_name": "无线电安全",
 *                   "tool_id": "5fcf4adef98f9239b6721811",
 *                   "tool_name": "zhuph20201222",
 *                   "update_time": 1608641269
 *               },
 *               {
 *                   "_id": "5fe1e09289e90f883fc531de",
 *                   "assets_cate_name": "TSP",
 *                   "assets_id": "XFE20201223",
 *                   "create_time": 1608638610,
 *                   "create_user_id": 130,
 *                   "create_user_name": "zhupenghui",
 *                   "pro_id": "5fd97feef98f923b505be841",
 *                   "status": 2,
 *                   "task_name": "task_20201222200330",
 *                   "task_number": 20201222200330,
 *                   "tool_category_name": "无线电安全",
 *                   "tool_id": "5fcf4adef98f9239b6721811",
 *                   "tool_name": "zhuph20201222",
 *                   "update_time": 1608638610
 *               },
 *               {
 *                   "_id": "5fe1e06789e90f883fc531dd",
 *                   "assets_cate_name": "TSP",
 *                   "assets_id": "XFE20201224",
 *                   "create_time": 1608638567,
 *                   "create_user_id": 130,
 *                   "create_user_name": "zhupenghui",
 *                   "pro_id": "5fd97feef98f923b505be841",
 *                   "status": 2,
 *                   "task_name": "task_20201222200247",
 *                   "task_number": 20201222200247,
 *                   "tool_category_name": "无线电安全",
 *                   "tool_id": "5fcf4adef98f9239b6721811",
 *                   "tool_name": "zhuph20201221",
 *                   "update_time": 1608638567
 *               }
 *           ],
 *           "pagination": {
 *               "count": 3,
 *               "current_page": 1,
 *               "per_page": 20,
 *               "total": 3,
 *               "total_pages": 1
 *           }
 *       }
 *   }
 */
func (this ToolTaskController) TaskList(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	//UserID := int(session.GetUserId(http_ctx.GetHttpCtx(ctx)))
	//UserName := session.GetUserName(ctx)
	search := req.String("search")
	categoryName := req.String("category_name")
	ProTaskId := req.MustString("pro_task_id")
	queryParams := qmap.QM{
		"ne_status":     0,
		"e_pro_task_id": ProTaskId,
	}
	if search != "" {
		search = custom_util.SpeCharsAddBackslash(search)
		queryParams.Merge(map[string]interface{}{"l_tool_name": search})
	}

	if categoryName != "" {
		queryParams.Merge(map[string]interface{}{"e_tool_category_name": categoryName})
	}

	mgoSession := mongo.NewMgoSession(common.MC_TOOL_TASK).AddCondition(queryParams).AddUrlQueryCondition(req.String("query_params"))
	mgoSession.SetTransformFunc(TaskToolFormat)
	if res, err := mgoSession.GetPage(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func TaskToolFormat(data qmap.QM) qmap.QM {

	toolId := data["tool_id"].(string)
	queryParams := qmap.QM{
		"e__id": bson.ObjectIdHex(toolId),
	}
	mgoSession := mongo.NewMgoSession(common.MC_TOOL).AddCondition(queryParams)
	if res, err := mgoSession.GetOne(); err == nil {
		data["tool_number"] = res.String("tool_number")
		return data
	} else {
		return nil
	}

}

/**
 * apiType http
 * @api {post} /api/v1/tool/task/create 创建工具任务
 * @apiVersion 0.1.0
 * @apiName CreateTask
 * @apiGroup ToolTask
 *
 * @apiDescription 创建工具任务
 *
 * @apiParam {string}               pro_task_id                         项目任务ID
 * @apiParam {string}               tool_name                           工具名称
 * @apiParam {string}               tool_id                             工具ID
 * @apiParam {string}               tool_number                         工具编号
 * @apiParam {string}               category_name                       工具类型名称
 * @apiParam {string}               assets_id                           资产ID
 * @apiParam {string}               assets_cate_name                    资产类别名称
 * @apiParam {string}               file_id                             上传文件GSF ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *    "pro_task_id":"5fd97feef98f923b505be841",
 *    "tool_name":"zhuph20201223",
 *    "tool_id":"5fcf4adef98f9239b6721811",
 *    "tool_number":"20201208174358",
 *    "category_name": "主机安全检测",
 *    "assets_id":"5fd97feef98f923b505be841",
 *    "assets_cate_name":"TSP",
 *    "file_id":"5e709c9f24b647174f6d05e9"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *    "code": 0,
 *    "data": {
 *        "Id": "5fe4520989e90fc46517cb5f",
 *        "ProTaskId": "5fd97feef98f923b505be841",
 *        "TaskName": "task_20201224163209",
 *        "TaskNumber": 20201224163209
 *    }
 * }
 */
func (this ToolTaskController) CreateTask(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	UserID := int(session.GetUserId(http_ctx.GetHttpCtx(ctx)))
	UserName := session.GetUserName(ctx)
	if rts, err := new(mongo_model.ToolTaskData).Create(req, UserID, UserName); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, rts)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/tool/task/del 删除工具任务
 * @apiVersion 0.1.0
 * @apiName DelTask
 * @apiGroup ToolTask
 *
 * @apiDescription 根据工具任务id删除工具任务
 *
 * @apiParam {string}               pro_task_id                         项目任务ID
 * @apiParam {string}               id                                  工具任务ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *    "pro_task_id":"5fd97feef98f923b505be841",
 *    "id":"5fe4520989e90fc46517cb5f",
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *    "code": 0,
 *    "data": {
 *        "Id": "5fe4520989e90fc46517cb5f",
 *        "ProTaskId": "5fd97feef98f923b505be841",
 *        "TaskName": "task_20201224163209",
 *        "TaskNumber": 20201224163209
 *    }
 * }
 */
func (this ToolTaskController) DelTask(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	masterId := req.MustString("id")
	ProTaskId := req.String("pro_task_id")
	delStatus := 0
	if rts, err := new(mongo_model.ToolTaskData).UpdateStatus(masterId, ProTaskId, delStatus); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, rts)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/tool/task/stop 暂停工具任务
 * @apiVersion 0.1.0
 * @apiName StopTask
 * @apiGroup ToolTask
 *
 * @apiDescription 暂停工具任务
 *
 * @apiParam {string}               pro_task_id                         项目任务ID
 * @apiParam {string}               id                                  工具任务ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *    "pro_task_id":"5fd97feef98f923b505be841",
 *    "id":"5fe4520989e90fc46517cb5f",
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *    "code": 0,
 *    "data": {
 *        "Id": "5fe4520989e90fc46517cb5f",
 *        "ProTaskId": "5fd97feef98f923b505be841",
 *        "TaskName": "task_20201224163209",
 *        "TaskNumber": 20201224163209
 *    }
 * }
 */
func (this ToolTaskController) StopTask(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	masterId := req.MustString("id")
	ProTaskId := req.MustString("pro_task_id")
	stopStatus := 3
	if rts, err := new(mongo_model.ToolTaskData).UpdateStatus(masterId, ProTaskId, stopStatus); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, rts)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/tool/task/detail 查看工具任务详情
 * @apiVersion 0.1.0
 * @apiName TaskDetail
 * @apiGroup ToolTask
 *
 * @apiDescription 查看工具任务详情
 *
 * @apiParam {string}               id                                  工具任务ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *    "id":"5fe4520989e90fc46517cb5f",
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *    "code": 0,
 *        "data": {
 *           "_id": ObjectId("5fe45c8989e90fdabc3501e8"),
 *           "pro_task_id": "5fd97feef98f923b505be841",
 *           "task_name": "task_20201224171657",
 *           "task_number": "D3EZYK",
 *           "tool_name": "zhuph20201223",
 *           "tool_id": "5fcf4adef98f9239b6721811",
 *           "tool_category_name": "主机安全检测",
 *           "assets_id": "5fd97feef98f923b505be841",
 *           "file_id": "5e709c9f24b647174f6d05e9",
 *           "assets_cate_name": "TSP",
 *           "create_user_id": NumberInt("0"),
 *           "create_user_name": "",
 *           "create_time": NumberInt("1608801417"),
 *           "update_time": NumberInt("1608801417"),
 *           "status": NumberInt("2")
 *       }
 * }
 */
func (this ToolTaskController) TaskDetail(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	masterID := bson.ObjectIdHex(req.MustString("id"))
	queryParams := qmap.QM{
		"ne_status": 0,
		"e__id":     masterID,
	}

	var assetVersion string
	mgoSession := mongo.NewMgoSession(common.MC_TOOL_TASK).AddCondition(queryParams).AddUrlQueryCondition(req.String("query_params"))
	if res, err := mgoSession.GetOne(); err == nil {
		assetID := res.MustString("assets_id")
		taskRts := *res
		if assetData, assetErr := new(mongo_model.EvaluateAsset).GetOne(assetID); assetErr == nil {
			assetVersion = assetData.MustString("version")
		}
		taskRts["assets_version"] = assetVersion
		response.RenderSuccess(ctx, taskRts)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/tool/task/result_list 查看工具任务解析结果
 * @apiVersion 0.1.0
 * @apiName TaskResultList
 * @apiGroup ToolTask
 *
 * @apiDescription 查看工具任务解析结果
 *
 * @apiParam {string}       id        工具任务ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *    "id":"5fe4520989e90fc46517cb5f",
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *    "code": 0,
 *        "data": {
 *           "list":[
 *               {
 *                   "_id": ObjectId("5fe0664baee3d1849a6a28db"),
 *                   "involve_module": "Qualcomm RF 组件",
 *                   "task_id": "123456",
 *                   "cve_type": "Eop",
 *                   "date_bulletin": "2016-04-02",
 *                   "date_exposure": "2016-01-06 05:16:34",
 *                   "description": "Android是美国谷歌（Google）公司和开放手持设备联盟（简称OHA）共同开发的一套以Linux为基础的开源操作系统。Qualcomm RF是使用在其中的一个美国高通（Qualcomm）公司开发的前端解决方案（包含功率放大器、天线开关、天线调谐器以及包络跟踪器等一系列芯片技术）组件。Android的Qualcomm RF组件中存在提权漏洞，该漏洞源于程序没有正确限制使用套接字ioctl调用。本地攻击者可借助特制的应用程序利用该漏洞获取权限，在内核上下文中执行任意代码。以下版本受到影响：Android 4.4.4之前版本，5.0.2之前版本，5.1.1之前版本，6.0之前版本和6.0.1之前版本。",
 *                   "google_severity_level": 1,
 *                   "cve_id": "CVE-2016-0844",
 *                   "fix_status": 4,
 *                   "sketch": "Qualcomm RF 组件中的提权漏洞",
 *                   "search_content": "CVE-2016-0844"
 *               },
 *               {
 *                   "_id": ObjectId("5fe0664baee3d1849a6a28e2"),
 *                   "sketch": "Qualcomm RF 组件中的提权漏洞",
 *                   "cve_id": "CVE-2016-0844",
 *                   "cve_type": "Eop",
 *                   "date_bulletin": "2016-04-02",
 *                   "date_exposure": "2016-01-06 05:16:34",
 *                   "description": "Android是美国谷歌（Google）公司和开放手持设备联盟（简称OHA）共同开发的一套以Linux为基础的开源操作系统。Qualcomm RF是使用在其中的一个美国高通（Qualcomm）公司开发的前端解决方案（包含功率放大器、天线开关、天线调谐器以及包络跟踪器等一系列芯片技术）组件。Android的Qualcomm RF组件中存在提权漏洞，该漏洞源于程序没有正确限制使用套接字ioctl调用。本地攻击者可借助特制的应用程序利用该漏洞获取权限，在内核上下文中执行任意代码。以下版本受到影响：Android 4.4.4之前版本，5.0.2之前版本，5.1.1之前版本，6.0之前版本和6.0.1之前版本。",
 *                   "google_severity_level": 1,
 *                   "involve_module": "Qualcomm RF 组件1",
 *                   "task_id": "123456",
 *                   "fix_status": 4,
 *                   "search_content": "CVE-2016-0844"
 *               }
 *            ]
 *       }
 * }
 */
func (this ToolTaskController) TaskResultList(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	taskId := req.MustString("task_id")
	getToolTaskParam := qmap.QM{
		"e__id":     bson.ObjectIdHex(taskId),
		"in_status": []int{1, 2, 3},
	}
	//获取当前工具任务信息
	if toolTaskRes, err := mongo.NewMgoSession(common.MC_TOOL_TASK).AddCondition(getToolTaskParam).GetOne(); err == nil {

		taskNumber := toolTaskRes.MustString("task_number")
		//获取此工具任务ID下已经绑定的cve主键ID
		queryBindParams := qmap.QM{
			"e_tool_task_id": taskId,
			"e_status":       1,
		}
		notInId := []bson.ObjectId{}
		if bindRts, bindErr := mongo.NewMgoSession(common.MC_TOOL_TASK_RESULT_BIND_TEST).AddCondition(queryBindParams).Get(); bindErr == nil {
			for _, val := range *bindRts {
				notInId = append(notInId, bson.ObjectIdHex(val["result_id"].(string)))
			}
		}
		//自定义sql not in
		//此处task_id 对应工具任务表中task_number

		customOperations := []bson.M{
			{"$match": bson.M{"task_number": taskNumber, "_id": bson.M{"$nin": notInId}}},
		}
		if notInRts, err := mongo.NewMgoSession(common.MC_TOOl_TASK_LOOPHOLE).QueryGet(customOperations); err == nil {
			response.RenderSuccess(ctx, notInRts)
		} else {
			response.RenderFailure(ctx, err)
		}
	} else {
		response.RenderFailure(ctx, errors.New("未发现此工具任务"))
	}
}

/**
 * apiType http
 * @api {get} /api/v1/tool/task/get_bind_result 根据测试用例ID获取绑定结果信息
 * @apiVersion 0.1.0
 * @apiName GetBindTaskResultForTestId
 * @apiGroup ToolTask
 *
 * @apiDescription 根据测试用例ID获取绑定结果信息
 *
 * @apiParam {string}       test_id        测试用例ID
 * @apiParam {string}       task_id        工具任务ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *    "test_id":"5fe4520989e90fc46517cb5f",
 *    "task_id":"5fe4520989e90fc46517cb5f",
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *    "code": 0,
 *        "data": {
 *           "list":[
 *               {
 *                   "_id": ObjectId("5fe0664baee3d1849a6a28db"),
 *                   "involve_module": "Qualcomm RF 组件",
 *                   "task_id": "123456",
 *                   "cve_type": "Eop",
 *                   "date_bulletin": "2016-04-02",
 *                   "date_exposure": "2016-01-06 05:16:34",
 *                   "description": "Android是美国谷歌（Google）公司和开放手持设备联盟（简称OHA）共同开发的一套以Linux为基础的开源操作系统。Qualcomm RF是使用在其中的一个美国高通（Qualcomm）公司开发的前端解决方案（包含功率放大器、天线开关、天线调谐器以及包络跟踪器等一系列芯片技术）组件。Android的Qualcomm RF组件中存在提权漏洞，该漏洞源于程序没有正确限制使用套接字ioctl调用。本地攻击者可借助特制的应用程序利用该漏洞获取权限，在内核上下文中执行任意代码。以下版本受到影响：Android 4.4.4之前版本，5.0.2之前版本，5.1.1之前版本，6.0之前版本和6.0.1之前版本。",
 *                   "google_severity_level": 1,
 *                   "cve_id": "CVE-2016-0844",
 *                   "fix_status": 4,
 *                   "sketch": "Qualcomm RF 组件中的提权漏洞",
 *                   "search_content": "CVE-2016-0844"
 *               },
 *               {
 *                   "_id": ObjectId("5fe0664baee3d1849a6a28e2"),
 *                   "sketch": "Qualcomm RF 组件中的提权漏洞",
 *                   "cve_id": "CVE-2016-0844",
 *                   "cve_type": "Eop",
 *                   "date_bulletin": "2016-04-02",
 *                   "date_exposure": "2016-01-06 05:16:34",
 *                   "description": "Android是美国谷歌（Google）公司和开放手持设备联盟（简称OHA）共同开发的一套以Linux为基础的开源操作系统。Qualcomm RF是使用在其中的一个美国高通（Qualcomm）公司开发的前端解决方案（包含功率放大器、天线开关、天线调谐器以及包络跟踪器等一系列芯片技术）组件。Android的Qualcomm RF组件中存在提权漏洞，该漏洞源于程序没有正确限制使用套接字ioctl调用。本地攻击者可借助特制的应用程序利用该漏洞获取权限，在内核上下文中执行任意代码。以下版本受到影响：Android 4.4.4之前版本，5.0.2之前版本，5.1.1之前版本，6.0之前版本和6.0.1之前版本。",
 *                   "google_severity_level": 1,
 *                   "involve_module": "Qualcomm RF 组件1",
 *                   "task_id": "123456",
 *                   "fix_status": 4,
 *                   "search_content": "CVE-2016-0844"
 *               }
 *            ]
 *       }
 * }
 */
func (this ToolTaskController) GetBindTaskResultForTestId(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	testId := req.MustString("test_id")
	toolTaskId := req.MustString("tool_task_id")
	queryParams := qmap.QM{
		"e_test_id":      testId,
		"e_tool_task_id": toolTaskId,
		"e_status":       1,
	}

	mgoSession := mongo.NewMgoSession(common.MC_TOOL_TASK_RESULT_BIND_TEST).AddCondition(queryParams).AddUrlQueryCondition(req.String("query_params"))
	mgoSession.SetTransformFunc(BindTaskResultFormat)
	if res, err := mgoSession.Get(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func BindTaskResultFormat(data qmap.QM) qmap.QM {

	resultId := data["result_id"].(string)
	queryParams := qmap.QM{
		"e__id": bson.ObjectIdHex(resultId),
	}
	mgoSession := mongo.NewMgoSession(common.MC_TOOl_TASK_LOOPHOLE).AddCondition(queryParams)
	if res, err := mgoSession.GetOne(); err == nil {
		data["cve"] = res
		return data
		//return data["cve"] = res
	} else {
		return nil
	}

}

/**
 * apiType http
 * @api {post} /api/v1/tool/task/result_link_test 工具任务结果与测试用例关联
 * @apiVersion 0.1.0
 * @apiName ResultLinkTest
 * @apiGroup ToolTask
 *
 * @apiDescription 工具任务结果与测试用例关联
 *
 * @apiParam {string}               pro_task_id                         项目任务ID
 * @apiParam {string}               tool_task_id                        工具任务ID
 * @apiParam {string}               result_id                           解析结果ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *    "pro_task_id":"5fd97feef98f923b505be841",
 *    "tool_task_id":"5fe4520989e90fc46517cb5f",
 *    "result_id":"5fe4520989e90fc46517cb5f",
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *    "code": 0,
 *    "data": {
 *        "Id": "5fe4520989e90fc46517cxtr"
 *    }
 * }
 */
func (this ToolTaskController) ResultLinkTest(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	UserID := int(session.GetUserId(http_ctx.GetHttpCtx(ctx)))
	UserName := session.GetUserName(ctx)
	if ret, err := new(mongo_model.ToolTaskResultBindTest).BindTestResult(req, UserID, UserName); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, ret)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/tool/task/result_unlink_test 解除工具任务结果与测试用例关联
 * @apiVersion 0.1.0
 * @apiName ResultUnLinkTest
 * @apiGroup ToolTask
 *
 * @apiDescription 解除工具任务结果与测试用例关联
 *
 * @apiParam {string}               id                           绑定结果ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *    "id":"5fe4520989e90fc46517cxtr",
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *    "code": 0,
 *    "data": {
 *        "Id": "5fe4520989e90fc46517cxtr"
 *    }
 * }
 */
func (this ToolTaskController) ResultUnLinkTest(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	UserID := int(session.GetUserId(http_ctx.GetHttpCtx(ctx)))
	UserName := session.GetUserName(ctx)
	if ret, err := new(mongo_model.ToolTaskResultBindTest).UnBindTestResult(req, UserID, UserName); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, ret)
	}
}

//@auto_generated_api_end
