package controller

import (
	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/http/transformer"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/common_lib/session"
	"skygo_detection/logic"
	"skygo_detection/mongo_model"
)

type ProjectHgTestTaskController struct{}

/**
 * apiType http
 * @api {get} /api/v1/project_reports 报告列表
 * @apiVersion 1.0.0
 * @apiName GetAll
 * @apiGroup ProjectReport
 *
 * @apiDescription 查询报告接口
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/project_reports
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "list": [
 *             {
 *                 "_id": "5e8f1b2624b64716658cf1d2",
 *                 "create_time": 1586436902729,
 *                 "file_id": "5e8f1b2624b64716658cf1ca",
 *                 "file_size": 1708567,
 *                 "name": "456789-周报告-1586436902.docx",
 *                 "operator_id": 0,
 *                 "project_id": "5e73a77c24b64720bdd8c9ea",
 *                 "report_type": "week"
 *             }
 *         ],
 *         "pagination": {
 *             "count": 6,
 *             "current_page": 1,
 *             "per_page": 20,
 *             "total": 6,
 *             "total_pages": 1
 *         }
 *     }
 * }
 */
func (this ProjectHgTestTaskController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	mgoSession := mongo.NewMgoSession(common.McHgTestTask).AddUrlQueryCondition(queryParams)
	mgoSession.SetTransformer(&transformer.HgTestTaskTransformer{})
	if res, err := mgoSession.GetPage(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/project_reports/:id 查询某一个报告信息
 * @apiVersion 1.0.0
 * @apiName GetOne
 * @apiGroup ProjectReport
 *
 * @apiDescription 查询某一个报告信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}       id        报告id
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/project_reports/:id
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "_id": "5e8f1b2624b64716658cf1d2",
 *         "create_time": 1586436902729,
 *         "file_id": "5e8f1b2624b64716658cf1ca",
 *         "file_size": 1708567,
 *         "name": "456789-周报告-1586436902.docx",
 *         "operator_id": 0,
 *         "project_id": "5e73a77c24b64720bdd8c9ea",
 *         "report_type": "week"
 *     }
 * }
 */
func (this ProjectHgTestTaskController) GetOne(ctx *gin.Context) {
	uuid := ctx.Param("id")
	result, _ := new(logic.HgTestTaskLogic).GetOne(uuid)
	response.RenderSuccess(ctx, result)
}

/**
 * apiType http
 * @api {post} /api/v1/project_reports 创建项目报告
 * @apiVersion 1.0.0
 * @apiName Create
 * @apiGroup ProjectReport
 *
 * @apiDescription 创建新项目报告
 *
 * @apiUse authHeader
 *
 * @apiParam {string}  		project_id  	项目id
 * @apiParam {string}   	report_type  	报告类型(周报:week,  初测报告:test, 复测报告:retest)
 * @apiParam {string}   	file_id  		文件id
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST http://localhost/api/v1/project_reports
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *         "file_id": "5e8f1b2624b64716658cf1ca",
 *         "project_id": "5e73a77c24b64720bdd8c9ea",
 *         "report_type": "week"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "_id": "5e8f1b2624b64716658cf1d2",
 *         "create_time": 1586436902729,
 *         "file_id": "5e8f1b2624b64716658cf1ca",
 *         "file_size": 1708567,
 *         "name": "456789-周报告-1586436902.docx",
 *         "operator_id": 0,
 *         "project_id": "5e73a77c24b64720bdd8c9ea",
 *         "report_type": "week"
 *     }
 * }
 */
func (this ProjectHgTestTaskController) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	name := req.MustString("name")
	createUserId := int(session.GetUserId(http_ctx.GetHttpCtx(ctx)))
	if model, err := new(logic.HgTestTaskLogic).CreateTask(createUserId, name); err == nil {
		response.RenderSuccess(ctx, model)
	} else {
		panic(err)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/project_reports 创建项目报告
 * @apiVersion 1.0.0
 * @apiName Create
 * @apiGroup ProjectReport
 *
 * @apiDescription 创建新项目报告
 *
 * @apiUse authHeader
 *
 * @apiParam {string}  		project_id  	项目id
 * @apiParam {string}   	report_type  	报告类型(周报:week,  初测报告:test, 复测报告:retest)
 * @apiParam {string}   	file_id  		文件id
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST http://localhost/api/v1/project_reports
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *         "file_id": "5e8f1b2624b64716658cf1ca",
 *         "project_id": "5e73a77c24b64720bdd8c9ea",
 *         "report_type": "week"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "_id": "5e8f1b2624b64716658cf1d2",
 *         "create_time": 1586436902729,
 *         "file_id": "5e8f1b2624b64716658cf1ca",
 *         "file_size": 1708567,
 *         "name": "456789-周报告-1586436902.docx",
 *         "operator_id": 0,
 *         "project_id": "5e73a77c24b64720bdd8c9ea",
 *         "report_type": "week"
 *     }
 * }
 */
func (this ProjectHgTestTaskController) Update(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	id := ctx.Param("id")
	if testCase, err := new(mongo_model.HgTestTask).Update(id, *req); err == nil {
		response.RenderSuccess(ctx, testCase)
	} else {
		panic(err)
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/project_reports 批量删除报告
 * @apiVersion 1.0.0
 * @apiName BulkDelete
 * @apiGroup ProjectReport
 *
 * @apiDescription 批量删除报告接口
 *
 * @apiUse authHeader
 *
 * @apiParam {[]string}   ids  报告id
 *
 * @apiSuccessExample {json} 请求成功示例:
 *       {
 *            "code": 0,
 *			  "data":{
 *				"number":1
 *			}
 *       }
 */
func (this ProjectHgTestTaskController) Delete(ctx *gin.Context) {
	taskId := ctx.Param("id")
	if err := new(logic.HgTestTaskLogic).DeleteAll(taskId); err == nil {
		response.RenderSuccess(ctx, nil)
	} else {
		panic(err)
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/project_reports 批量删除报告
 * @apiVersion 1.0.0
 * @apiName BulkDelete
 * @apiGroup ProjectReport
 *
 * @apiDescription 批量删除报告接口
 *
 * @apiUse authHeader
 *
 * @apiParam {[]string}   ids  报告id
 *
 * @apiSuccessExample {json} 请求成功示例:
 *       {
 *            "code": 0,
 *			  "data":{
 *				"number":1
 *			}
 *       }
 */
func (this ProjectHgTestTaskController) BulkDelete(ctx *gin.Context) {
	params := &qmap.QM{}
	req := params.Merge(*request.GetRequestBody(ctx))
	param := req.MustSlice("uuid_list")
	uuidList := make([]string, 0)
	for _, v := range param {
		if _v, ok := v.(string); ok {
			uuidList = append(uuidList, _v)
		}
	}
	err := new(logic.HgTestTaskLogic).DeleteAll(uuidList...)
	if err == nil {
		response.RenderSuccess(ctx, nil)
	} else {
		panic(err)
	}
}

/**
 * apiType http
 * @api {POST} /api/v1/project_report/status 获取报告状态列表
 * @apiVersion 1.0.0
 * @apiName Status
 * @apiGroup ProjectReport
 *
 * @apiDescription 获取报告状态列表
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 *	{
 *		"code": 0,
 *		"data": [
 *			{
 *				"status": 0,
 *				"name": "创建"
 *			},
 *		]
 *	}
 */
func (this ProjectHgTestTaskController) GetStatusFlow(ctx *gin.Context) {
	uuid := ctx.Param("id")
	result, _ := new(logic.HgTestTaskLogic).GetStatusFlow(uuid)
	response.RenderSuccess(ctx, result)
}

/**
 * apiType http
 * @api {POST} /api/v1/project_report/phase 获取报告审核阶段
 * @apiVersion 1.0.0
 * @apiName Phase
 * @apiGroup ProjectReport
 *
 * @apiDescription 获取报告审核阶段
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 *	{
 *		"code": 0,
 *		"data": [
 *			{
 *				"status": 0,
 *				"name": "创建"
 *			},
 *		]
 *	}
 */
func (this ProjectHgTestTaskController) GetTestCase(ctx *gin.Context) {
	uuid := ctx.Param("id")
	result, _ := new(logic.HgTestTaskLogic).GetTestCase(uuid)
	response.RenderSuccess(ctx, result)
}

/**
 * apiType http
 * @api {post} /api/v1/project_reports 合规测试任务完成
 * @apiVersion 1.0.0
 * @apiName Create
 * @apiGroup ProjectReport
 *
 * @apiDescription 合规测试任务完成
 *
 * @apiUse authHeader
 *
 * @apiParam {string}  		project_id  	项目id
 * @apiParam {string}   	report_type  	报告类型(周报:week,  初测报告:test, 复测报告:retest)
 * @apiParam {string}   	file_id  		文件id
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST http://localhost/api/v1/project_reports
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *         "file_id": "5e8f1b2624b64716658cf1ca",
 *         "project_id": "5e73a77c24b64720bdd8c9ea",
 *         "report_type": "week"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "_id": "5e8f1b2624b64716658cf1d2",
 *         "create_time": 1586436902729,
 *         "file_id": "5e8f1b2624b64716658cf1ca",
 *         "file_size": 1708567,
 *         "name": "456789-周报告-1586436902.docx",
 *         "operator_id": 0,
 *         "project_id": "5e73a77c24b64720bdd8c9ea",
 *         "report_type": "week"
 *     }
 * }
 */
func (this ProjectHgTestTaskController) Complete(ctx *gin.Context) {
	params := &qmap.QM{}
	req := params.Merge(*request.GetRequestBody(ctx))
	req["status"] = "complete"

	id := ctx.Param("id")
	if testCase, err := new(mongo_model.HgTestTask).Update(id, req); err == nil {
		response.RenderSuccess(ctx, testCase)
	} else {
		panic(err)
	}
}

/*
*
测试用例修改

	{
	  "test_case_id" : "xxx",
	  "task_id" : "xxx",
	  "status" : 1,
	  "man_made_result_desc" : "xxx"
	}
*/
func (this ProjectHgTestTaskController) UpdateTestCase(ctx *gin.Context) {
	// 请求内容参考
	// {
	//   "test_case_id" : "TC999030M109",
	//   "task_id" : "220E5J",
	//   "status" : 4,
	//   "man_made_result_desc" : "xxx12321355555",
	//   "delete_file_id" : "60a61f17e830c668b95e8f40"
	// }
	params := &qmap.QM{}
	req := params.Merge(*request.GetRequestBody(ctx))

	taskId := req.MustString("task_id")
	testCaseId := req.MustString("test_case_id")
	err := new(logic.HgTestTaskLogic).UpdateTestCase(taskId, testCaseId, req)
	if err == nil {
		response.RenderSuccess(ctx, nil)
	} else {
		panic(err)
	}
}

// todo 还要提供5个微服务
//这部分的微服务 已经修改到本地调用，scan_service下边
