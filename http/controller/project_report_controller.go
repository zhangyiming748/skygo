package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/common_lib/session"
	"skygo_detection/mongo_model"
)

type ProjectReportController struct{}

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
func (this ProjectReportController) GetAll(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	queryParams := params.String("query_params")
	orm_mongo.NewWidgetWithCollectionName(common.MC_REPORT).SetQueryStr(queryParams).PaginatorFind()
	mgoSession := orm_mongo.NewWidgetWithCollectionName(common.MC_REPORT).SetQueryStr(queryParams)
	if res, err := mgoSession.PaginatorFind(); err == nil {
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
func (this ProjectReportController) GetOne(ctx *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(request.ParamString(ctx, "id"))
	params := qmap.QM{
		"e__id": id,
	}
	ormSession := orm_mongo.NewWidgetWithParams(common.MC_REPORT, params)
	if result, err := ormSession.Get(); err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
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
func (this ProjectReportController) Create(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	if report, err := new(mongo_model.Report).Create(int(session.GetUserId(ctx)), params); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		result, _ := custom_util.StructToMap(*report)
		response.RenderSuccess(ctx, result)
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
func (this ProjectReportController) BulkDelete(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	effectNum := 0
	if rawIds, has := params.TrySlice("ids"); has {
		ids := []bson.ObjectId{}
		for _, id := range rawIds {
			bsonId := bson.ObjectIdHex(id.(string))
			if reportInfo, err := mongo.NewMgoSession(common.MC_REPORT).AddCondition(qmap.QM{"e__id": bsonId}).GetOne(); err == nil {
				if (*reportInfo)["status"] != common.RS_NEW {
					response.RenderFailure(ctx, errors.New("非创建阶段，不可删除"))
				}
				ids = append(ids, bsonId)
			}
		}
		if len(ids) > 0 {
			match := bson.M{
				"_id": bson.M{"$in": ids},
			}
			if changeInfo, err := mongo.NewMgoSession(common.MC_REPORT).RemoveAll(match); err == nil {
				effectNum = changeInfo.Removed
			} else {
				response.RenderFailure(ctx, err)
			}
		}
	}
	result := &qmap.QM{"number": effectNum}
	response.RenderSuccess(ctx, result)
}

/**
 * apiType http
 * @api {POST} /api/v1/project_report/export 报告导出
 * @apiVersion 1.0.0
 * @apiName Export
 * @apiGroup ProjectReport
 *
 * @apiDescription 报告导出
 *
 * @apiUse authHeader
 *
 * @apiParam {string}  		project_id  	项目id
 * @apiParam {[]string}   	item_ids  		测试对象id
 * @apiParam {string}   	report_type  	报告类型(周报:week,  初测报告:test, 复测报告:retest)
 * @apiParam {string}   	report_all  	是否导出全部测试项的报告
 *
 * @apiSuccessExample {json} 请求成功示例:
 *       {
 *            "code": 0
 *       }
 */
func (this ProjectReportController) Export(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	projectId := params.MustString("project_id")
	reportType := params.MustString("report_type")
	rawIds, _ := params.TrySlice("item_ids")
	evaluateItems := []string{}
	for _, id := range rawIds {
		evaluateItems = append(evaluateItems, id.(string))
	}
	if _, err := new(mongo_model.Report).ExportReport(int(session.GetUserId(ctx)), projectId, reportType, evaluateItems, ctx); err != nil {
		response.RenderFailure(ctx, err)
	}
	response.RenderSuccess(ctx, "")
}

//@auto_generated_api_end
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
func (this ProjectReportController) Status(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	data := []qmap.QM{
		{"status": common.RS_NEW, "name": "创建"},
		{"status": common.RS_AUDIT, "name": "审核"},
		{"status": common.RS_SUCCESS, "name": "发布成功"},
		{"status": common.RS_FAILED, "name": "发布失败"},
	}

	response.RenderSuccess(ctx, data)
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
func (this ProjectReportController) Phase(ctx *gin.Context) {
	params := qmap.QM{
		"e_status": common.ENABLED,
	}
	mgoSession := orm_mongo.NewWidgetWithParams(common.MC_REPORT_PHASE, params)
	mgoSession.AddSorter("_id", 0)
	if res, err := mgoSession.SetLimit(10000).Find(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {POST} /api/v1/project_report/create_phase 创建报告审核阶段
 * @apiVersion 1.0.0
 * @apiName CreatePhase
 * @apiGroup ProjectReport
 *
 * @apiDescription 创建报告审核阶段
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
func (this ProjectReportController) CreatePhase(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	if reportPhase, err := new(mongo_model.ReportPhase).Create(int(session.GetUserId(ctx)), params); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		result, _ := custom_util.StructToMap(*reportPhase)
		response.RenderSuccess(ctx, result)
	}

}

/**
 * apiType http
 * @api {POST} /api/v1/project_report/create_node 创建审核节点
 * @apiVersion 1.0.0
 * @apiName CreateNode
 * @apiGroup ProjectReport
 *
 * @apiDescription 创建审核节点
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 *	{
 *		"code": 0,
 *		"data": true
 *	}
 */
func (this ProjectReportController) CreateNode(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	nodeInfo := qmap.QM{
		"project_id": params.MustString("project_id"),
		"report_id":  params.MustString("report_id"),
		"auditor_id": params.MustInt("auditor_id"),
		"name":       params.MustString("name"),
		"history": mongo_model.History{
			Result:    true,
			Operation: "添加审核节点",
			Comment:   "",
			OpId:      int(session.GetUserId(ctx)),
			OpTime:    custom_util.GetCurrentMilliSecond(),
		},
	}

	if _, err := new(mongo_model.ReportNode).Create(nodeInfo); err == nil {
		// 更新报告信息
		updateItem := qmap.QM{
			"current_operator_id": params.MustInt("auditor_id"),
			"current_phase":       params.MustString("name"),
			"update_time":         custom_util.GetCurrentMilliSecond(),
		}
		if err := new(mongo_model.Report).Update(params.MustString("project_id"), updateItem); err != nil {
			response.RenderFailure(ctx, err)
		}

		client := mongo.NewMgoSession(common.MC_REPORT).AddCondition(qmap.QM{"e__id": bson.ObjectIdHex(nodeInfo["report_id"].(string))})
		if report, err := client.GetOne(); err == nil {
			// 创建节点成功后，将审核人添加到报告相关人员列表中
			selector := bson.M{
				"_id": bson.ObjectIdHex(nodeInfo["report_id"].(string)),
			}
			var relativeId []int
			if (*report)["relative_id"] == nil {
				relativeId = append(relativeId, nodeInfo["auditor_id"].(int))
			} else {
				relativeId = custom_util.InterfaceToInt((*report)["relative_id"].([]interface{}))
				if !custom_util.InIntSlice(nodeInfo["auditor_id"].(int), relativeId) {
					relativeId = append(relativeId, nodeInfo["auditor_id"].(int))
				}
			}

			data := qmap.QM{
				"relative_id": relativeId,
			}

			if err := mongo.NewMgoSession(common.MC_REPORT).Update(selector, qmap.QM{"$set": data}); err != nil {
				response.RenderFailure(ctx, err)
			}
		}
		response.RenderSuccess(ctx, &qmap.QM{"data": true})
	}
	response.RenderSuccess(ctx, &qmap.QM{"data": false})
}

/**
 * apiType http
 * @api {POST} /api/v1/project_report/audit 审核节点
 * @apiVersion 1.0.0
 * @apiName Audit
 * @apiGroup ProjectReport
 *
 * @apiDescription 审核节点
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 *	{
 *		"code": 0,
 *		"data": true
 *	}
 */
func (this ProjectReportController) Audit(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	opId := int(session.GetUserId(ctx))
	reportId := params.MustString("report_id")
	nodeId := params.MustString("node_id")
	result := params.MustInt("result")
	comment := params.MustString("comment")

	params1 := qmap.QM{
		"e__id":       bson.ObjectIdHex(nodeId),
		"e_report_id": reportId,
	}
	mgoSession := mongo.NewMgoSession(common.MC_REPORT_NODE).AddCondition(params1)
	if nodeInfo, err := mgoSession.GetOne(); err == nil {
		if (*nodeInfo)["auditor_id"] == opId {
			RES := false
			operationTag := "驳回"
			if result == common.RAS_SUCCESS {
				RES = true
				operationTag = "通过"
			} else {
				// 驳回，将报告状态改为 发布失败
				selector := bson.M{"_id": bson.ObjectIdHex(reportId)}
				updateItem := bson.M{
					"$set": qmap.QM{
						"status": common.RS_FAILED,
					},
				}
				if err := mongo.NewMgoSession(common.MC_REPORT).Update(selector, updateItem); err != nil {
					response.RenderFailure(ctx, err)
				}
			}
			history := (*nodeInfo)["history"]
			history = append(history.([]interface{}), mongo_model.History{
				Result:    RES,
				Operation: "节点审核 " + operationTag,
				Comment:   comment,
				OpId:      opId,
				OpTime:    custom_util.GetCurrentMilliSecond(),
			})
			// 修改节点数据
			selector := bson.M{"_id": bson.ObjectIdHex(nodeId)}
			updateItem := bson.M{
				"$set": qmap.QM{
					"result":     result,
					"history":    history,
					"audit_time": custom_util.GetCurrentMilliSecond(),
				},
			}
			if err := mongo.NewMgoSession(common.MC_REPORT_NODE).Update(selector, updateItem); err == nil {
				response.RenderSuccess(ctx, &qmap.QM{"data": true})
			} else {
				response.RenderFailure(ctx, err)
			}
		} else {
			response.RenderFailure(ctx, errors.New("当前用户没有审核权限"))
		}
	} else {
		//特殊处理，之前的代码即返回错误，又返回data数据
		m := gin.H{
			"code": 0,
			"msg":  err.Error(),
			"data": qmap.QM{"data": false},
		}
		ctx.AbortWithStatusJSON(200, m)
	}
}

/**
 * apiType http
 * @api {POST} /api/v1/project_report/node 获取节点列表
 * @apiVersion 1.0.0
 * @apiName Node
 * @apiGroup ProjectReport
 *
 * @apiDescription 获取节点列表
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 *	{
 *		"code": 0,
 *		"data": true
 *	}
 */
func (this ProjectReportController) Node(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	reportId := params.MustString("report_id")
	// 查询报告信息
	params1 := qmap.QM{
		"e__id": bson.ObjectIdHex(reportId),
	}

	info, err := mongo.NewMgoSessionWithCond(common.MC_REPORT, params1).GetOne()
	if err != nil {
		response.RenderFailure(ctx, err)
	}

	params1 = qmap.QM{
		"e_report_id": reportId,
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_REPORT_NODE, params1)
	mongoClient.AddSorter("_id", 0)
	if result, err := mongoClient.Get(); err == nil {
		for index, item := range *result {
			item["auditor_name"] = mongo_model.GetUserName(item["auditor_id"].(int), ctx)
			for index2, history := range item["history"].([]interface{}) {
				item["history"].([]interface{})[index2].(map[string]interface{})["op_name"] = mongo_model.GetUserName(history.(map[string]interface{})["op_id"].(int), ctx)
			}
			(*result)[index] = item
		}
		(*info)["node"] = result
		response.RenderSuccess(ctx, qmap.QM{"data": info})
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {POST} /api/v1/project_report/publish 发布报告
 * @apiVersion 1.0.0
 * @apiName Publish
 * @apiGroup ProjectReport
 *
 * @apiDescription 发布报告
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 *	{
 *		"code": 0,
 *		"data": true
 *	}
 */
func (this ProjectReportController) Publish(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	reportId := params.MustString("report_id")
	params1 := qmap.QM{
		"e__id": bson.ObjectIdHex(reportId),
	}
	reportInfo, err := mongo.NewMgoSessionWithCond(common.MC_REPORT, params1).GetOne()
	if err != nil {
		response.RenderFailure(ctx, err)
	}
	if (*reportInfo)["status"] == common.RS_SUCCESS || (*reportInfo)["status"] == common.RS_FAILED {
		response.RenderFailure(ctx, errors.New("当前状态不能发布"))
	}

	updateItem := qmap.QM{
		"status":              common.RS_SUCCESS,
		"current_operator_id": 0,
		"current_phase":       "发布成功",
		"update_time":         custom_util.GetCurrentMilliSecond(),
	}
	if err := new(mongo_model.Report).Update(reportId, updateItem); err == nil {
		response.RenderSuccess(ctx, qmap.QM{"data": true})
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {POST} /api/v1/project_report/getList 获取用户相关报告列表
 * @apiVersion 1.0.0
 * @apiName getList
 * @apiGroup ProjectReport
 *
 * @apiDescription 发布报告
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 *	{
 *		"code": 0,
 *		"data": true
 *	}
 */
func (this ProjectReportController) GetList(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	if res, err := new(mongo_model.Report).GetList(*params, int(session.GetUserId(ctx)), ctx); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}
