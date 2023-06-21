package controller

import (
	"context"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/common_lib/session"
	"skygo_detection/logic"
	"skygo_detection/mongo_model"
)

type EvaluateTaskController struct{}

//@auto_generated_api_begin
/**
 * apiType http
 * @api {get} /api/v1/evaluate_tasks 分页查询项目任务列表
 * @apiVersion 0.1.0
 * @apiName GetAll
 * @apiGroup EvaluateTask
 *
 * @apiDescription 分页查询项目任务列表
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": {
 *                 "list": [
 *                         {
 *                                 "_id": "5fe01f8224b647330dfc3e77",
 *                                 "create_time": 1608523650207,
 *                                 "evaluate_item_ids": [
 *                                         "3",
 *                                         "2"
 *                                 ],
 *                                 "name": "项目任务名称",
 *                                 "op_id": 0,
 *                                 "project_id": "5fd7218624b64712a27f47e8",
 *                                 "status": 1,
 *                                 "update_time": 1608523650207
 *                         }
 *                 ],
 *                 "pagination": {
 *                         "count": 1,
 *                         "current_page": 1,
 *                         "per_page": 20,
 *                         "total": 1,
 *                         "total_pages": 1
 *                 }
 *         }
 * }
 */
func (this EvaluateTaskController) GetAll(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["query_params"] = ctx.Request.URL.RawQuery

	if res, err := new(mongo_model.EvaluateTask).GetAll(*req, ctx); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_tasks/:id 查询某一项目任务信息
 * @apiVersion 0.1.0
 * @apiName GetOne
 * @apiGroup EvaluateTask
 *
 * @apiDescription 根据id查询某一项目任务信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   id  		项目任务id
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": {
 *                 "_id": "5fe01f8224b647330dfc3e77",
 *                 "create_time": 1608523650207,
 *                 "evaluate_item_ids": [
 *                         "3",
 *                         "2"
 *                 ],
 *                 "name": "项目任务名称",
 *                 "op_id": 0,
 *                 "project_id": "5fd7218624b64712a27f47e8",
 *                 "status": 1,
 *                 "update_time": 1608523650207
 *         }
 * }
 */
func (this EvaluateTaskController) GetOne(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["query_params"] = ctx.Request.URL.RawQuery

	id := request.ParamString(ctx, "id")

	params := qmap.QM{
		"e__id": bson.ObjectIdHex(id),
	}
	ormSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params)
	ormSession.SetTransformFunc(TaskNode)
	result, err := ormSession.GetOne()

	if err == nil {
		//查询项目信息
		info, _ := mongo.NewMgoSession(common.MC_PROJECT).AddCondition(qmap.QM{"e__id": bson.ObjectIdHex((*result)["project_id"].(string))}).GetOne()
		fmt.Println((*info)["manager_id"].(int))
		(*result)["manager_name"] = mongo_model.GetUserName((*info)["manager_id"].(int), ctx)
		(*result)["project_name"] = (*info)["name"]
		(*result)["op_name"] = mongo_model.GetUserName((*result)["op_id"].(int), ctx)
		(*result)["task_auditor_name"] = mongo_model.GetUserName((*result)["task_auditor_id"].(int), ctx)
		(*result)["report_auditor_name"] = mongo_model.GetUserName((*result)["report_auditor_id"].(int), ctx)
		(*result)["tester_name"] = mongo_model.GetUserName((*result)["tester_id"].(int), ctx)

		list := (*result)["node"].(qmap.QM)["list"].([]map[string]interface{})
		for index, item := range list {
			history := item["history"].([]interface{})
			for subIndex, historyItem := range history {
				opId := historyItem.(map[string]interface{})["op_id"].(int)
				list[index]["history"].([]interface{})[subIndex].(map[string]interface{})["op_name"] = mongo_model.GetUserName(opId, ctx)
			}

		}
		(*result)["node"].(qmap.QM)["list"] = list

	}

	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, result)
	}
}

func TaskNode(data qmap.QM) qmap.QM {
	if taskNode, err := new(mongo_model.EvaluateTaskNode).GetTaskNode(data["_id"].(bson.ObjectId).Hex(), data["status"].(int)); err == nil {
		data["node"] = taskNode
	} else {
		data["node"] = nil
	}
	return data
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_tasks 添加项目任务
 * @apiVersion 0.1.0
 * @apiName Create
 * @apiGroup EvaluateTask
 *
 * @apiDescription 添加项目任务
 *
 * @apiUse authHeader
 *
 * @apiParam {string}      	                        name    			    项目任务名称
 * @apiParam {string}      	                       	project_id    		    项目id
 * @apiParam {string}      	                       	evaluate_item_ids    	测试用例id
 * @apiParam {map[string]string]}      	            [asset_versions]    	资产对应的版本信息（默认关联最新资产版本）
 *
 * @apiParamExample {json}  请求参数示例:
 *	{
 *		"project_id": "5fd7218624b64712a27f47e8",
 *		"name": "项目任务名称",
 *		"evaluate_item_ids": ["3","2"],
 *		"test_phase": "初测",
 *		"asset_versions":{
 *			"1as123a":"1.0",
 *			"1as123b":"2.0"
 *		}
 *	}
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": {
 *                 "create_time": 1608371080128,
 *                 "evaluate_item_ids": [
 *                         "3",
 *                         "2"
 *                 ],
 *                 "id": "5fddcb8824b64731079c2765",
 *                 "name": "项目任务名称",
 *                 "op_id": 0,
 *                 "project_id": "5fd7218624b64712a27f47e8",
 *                 "status": 1,
 *                 "update_time": 1608371080128
 *         }
 * }
 */
func (this EvaluateTaskController) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	task, err := logic.EvaluateTaskCreate(*req, session.GetUserId(http_ctx.GetHttpCtx(ctx)))
	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, qmap.QM{"data": custom_util.StructToMap2(*task)})
	}
}

/**
 * apiType http
 * @api {put} /api/v1/evaluate_tasks/:id 更新项目任务
 * @apiVersion 0.1.0
 * @apiName Update
 * @apiGroup EvaluateTask
 *
 * @apiDescription 根据id,更新项目任务信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   							id  					项目任务id
 * @apiParam {string}      	                        name    			    项目任务名称
 * @apiParam {string}      	                       	project_id    		    项目id
 * @apiParam {string}      	                       	evaluate_item_ids    	测试用例id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *		"id":"5fddcb8824b64731079c2765",
 *		"project_id": "5fd7218624b64712a27f47e8",
 *		"name": "项目任务名称11",
 *		"evaluate_item_ids": ["3","2"]
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": {
 *                 "create_time": 1608371080128,
 *                 "evaluate_item_ids": [],
 *                 "id": "5fddcb8824b64731079c2765",
 *                 "name": "项目任务名称11",
 *                 "op_id": 0,
 *                 "project_id": "5fd7218624b64712a27f47e8",
 *                 "status": 1,
 *                 "update_time": 1608371080128
 *         }
 * }
 */
func (this EvaluateTaskController) Update(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")

	if evaluateTask, err := logic.EvaluateTaskUpdate(req.MustString("id"), *req); err == nil {
		response.RenderSuccess(ctx, qmap.QM{"data": custom_util.StructToMap2(*evaluateTask)})
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/evaluate_tasks 批量删除项目任务
 * @apiVersion 0.1.0
 * @apiName BulkDelete
 * @apiGroup EvaluateTask
 *
 * @apiDescription 批量删除项目任务
 *
 * @apiUse authHeader
 *
 * @apiParam {[]string}   ids  项目任务id
 *
 * @apiParamExample {json}  请求参数示例:
 *  {
 *      "ids":[
 *		    "5fddcb8824b64731079c2765"
 * 	    ]
 *  }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":{
 *				"number":1
 *			}
 *      }
 */
func (this EvaluateTaskController) BulkDelete(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	effectNum := 0
	if rawIds, has := req.TrySlice("ids"); has {
		idsObj := []bson.ObjectId{}
		ids := []string{}
		for _, id := range rawIds {
			//检查每个任务是否可以删除
			params := qmap.QM{
				"e_id": bson.ObjectIdHex(id.(string)),
			}
			if info, err := orm_mongo.NewWidgetWithCollectionName(common.MC_EVALUATE_TASK).SetParams(params).Get(); err == nil && info["status"] != common.PTS_CREATE {
				response.RenderFailure(ctx, errors.New("只有创建阶段可以删除任务"))
				return
			}
			idsObj = append(idsObj, bson.ObjectIdHex(id.(string)))
			ids = append(ids, id.(string))
		}

		// 删除项目任务
		if len(idsObj) > 0 {
			match := bson.M{
				"_id": bson.M{"$in": idsObj},
			}
			collEt := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_TASK)
			if deleteResult, err := collEt.DeleteMany(context.Background(), match); err == nil {
				effectNum = int(deleteResult.DeletedCount)
			} else {
				response.RenderFailure(ctx, err)
				return
			}

			// 第一步，清除任务用例关系表
			match = bson.M{
				"evaluate_task_id": bson.M{"$in": ids},
			}
			collEti := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_TASK_ITEM)
			if _, err := collEti.DeleteMany(context.Background(), match); err != nil {
				response.RenderFailure(ctx, err)
				return
			}

			// 第二步，重置测试用例状态，且接触预绑定
			selector := bson.M{
				"pre_bind": bson.M{"$in": ids},
			}
			updateItem := bson.M{
				"$set": qmap.QM{
					"status":   common.EIS_FREE,
					"pre_bind": "",
				},
			}
			collEi := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_ITEM)
			if _, err := collEi.UpdateMany(context.Background(), selector, updateItem); err != nil {
				response.RenderFailure(ctx, err)
				return
			}
		}
	}

	response.RenderSuccess(ctx, qmap.QM{"number": effectNum})
	return
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_task/report/:id 查询项目任务报告
 * @apiVersion 0.1.0
 * @apiName GetTaskReport
 * @apiGroup EvaluateTask
 *
 * @apiDescription 查询项目任务报告
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   	id  							项目任务id
 * @apiParam {string}   	[item_id]  						测试用例id
 * @apiParam {string}   	[item_name]  					测试用例名称
 * @apiParam {int}			[record_audit_status]			审核状态(1:通过 0:待审核 -1:驳回)
 *
 * @apiParamExample {json}  请求参数示例:
 *     {
 *  	    "id":"5fe464b1f98f923e40e8dd5e"
 *     }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "children": [
 *                 {
 *                     "children": [
 *                         {
 *                             "children": [
 *                                 {
 *                                     "_id": "12314123", //测试用例编号
 *                                     "asset_id": "US4VY61L", //测试资产id
 *                                     "auto_test_level": "自动化",自动化测试程度
 *                                     "children": [],
 *                                     "create_time": 1608715443791,
 *                                     "evaluate_task_id": "5fe464b1f98f923e40e8dd5e",
 *                                     "external_input": "", //外部输入
 *                                     "label": "测试用例名称", //测试用例名称
 *                                     "level": 1,//测试难度（1:低、2:中、3:高）
 *                                     "module_name": "android", //测试组件
 *                                     "module_type": "应用系统管理测试", //测试分类
 *                                     "module_type_id": "5fd87548aee3d1849a56f112",//测试分类id
 *                                     "name": "测试用例名称", //测试用例名称
 *                                     "objective": "测试目的",//测试目的
 *                                     "op_id": 64,
 *                                     "project_id": "5fd7218624b64712a27f47e8",
 *                                     "record_id": "5fe464b1f98f923e40e8dd5f",
 *                                     "records": [ //测试记录
 *                                         {
 *                                             "_id": "5fe464b1f98f923e40e8dd5f",
 *                                             "asset_id": "US4VY61L",//资产id
 *                                             "asset_version": "1.0",//资产版本
 *                                             "attachment": null,//测试附件
 *                                             "conclude": "测试结论",//测试结论
 *                                             "create_time": 1608877759403,//测试时间
 *                                             "evaluate_task_id": "5fe464b1f98f923e40e8dd5e",
 *                                             "item_id": "12314123",//测试用例id
 *                                             "op_id": 0,
 *                                             "project_id": "5fd7218624b64712a27f47e8",
 *                                             "risk_type": "设计",//风险根源类型(设计/配置/代码/其他)
 *                                             "test_phase": 1,//测试阶段 （1:初测、2:复测1、3:复测2、4:复测3 ...）
 *                                             "test_procedure": "测试过程"//测试过程
 *                                         }
 *                                     ],
 *                                     "status": 0,//测试状态 （0:待测试 1:测试完成）
 *                                     "tag": [
 *                                         "a",
 *                                         "b",
 *                                         "c",
 *                                         "d"
 *                                     ],
 *                                     "test_case_level": "基础测试",//测试用例级别（基础测试、全面测试、提高测试、专家测试）
 *                                     "test_count": 0,//测试次数
 *                                     "test_method": "黑盒",//测试方法（黑盒、白盒）
 *                                     "test_phase": 1,//测试阶段 （1:初测、2:复测1、3:复测2、4:复测3 ...）
 *                                     "test_procedure": "测试步骤",//测试过程
 *                                    "test_script": [ //测试脚本
 *                                           {
 *                                                   "name": "脚本",
 *                                                   "value": "dadfa123"
 *                                           }
 *                                   ],
 *                                   "test_sketch_map": [ //测试框架示意图
 *                                           {
 *                                                   "name": "图片",
 *                                                   "value": "dadfa124"
 *                                           }
 *                                   ],
 *                                   "test_standard": "测试标准",//测试标准
 *                                   "test_time": 0,
 *                                   "update_time": 1609157090435,
 *                                   "vul_number": 0
 *                                 }
 *                             ],
 *                             "id": "5fd87548aee3d1849a56f112",
 *                             "label": "应用系统管理测试"
 *                         }
 *                     ],
 *                     "label": "android"
 *                 }
 *             ],
 *             "id": "US4VY61L",
 *             "label": "s"
 *         }
 *     ]
 * }
 */
func (this EvaluateTaskController) GetTaskReport(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")

	id := req.MustString("id")
	if report, err := new(mongo_model.EvaluateTaskItem).GetEvaluateTaskReport(id, *req); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, report)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_task/asset_versions 查询项目任务关联的资产版本信息
 * @apiVersion 0.1.0
 * @apiName GetAssetVersions
 * @apiGroup EvaluateTask
 *
 * @apiDescription 查询项目任务关联的资产版本信息
 *
 * @apiUse authHeader
 *
 * @apiParam 	{string}     task_id  		项目任务id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  	"task_id":"5fd7218624b64712a27f47e8"
 * }
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": [
 *                 {
 *                         "id": "US4VX1KR",
 *                         "name": "test",
 *                         "versions": [
 *                                 "1.0"
 *                         ]
 *                 }
 *         ]
 * }
 */
func (this EvaluateTaskController) GetAssetVersions(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if assetIds := new(mongo_model.EvaluateItem).GetTaskRelatedAssets(req.MustString("task_id")); len(assetIds) >= 0 {
		assetList := []qmap.QM{}
		for _, assetId := range assetIds {
			if assetInfo, err := new(mongo_model.EvaluateAsset).GetAssetAllVersions(assetId); err == nil {
				assetList = append(assetList, assetInfo)
			}
		}
		response.RenderSuccess(ctx, assetList)
	} else {
		response.RenderSuccess(ctx, gin.H{})
	}
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_task/project_tasks 查询项目中的所有项目任务
 * @apiVersion 0.1.0
 * @apiName GetProjectTasks
 * @apiGroup EvaluateTask
 *
 * @apiDescription 查询项目中的所有项目任务
 *
 * @apiUse authHeader
 *
 * @apiParam 	{string}   	 project_id  			项目id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  	"project_id":"5ffd03f14806351af1b696f8"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": [
 *                 {
 *                         "_id": "5ffd64bd4806351b6c77da8c",
 *                         "asset_versions": {
 *                                 "KJTNLBV1": "1.1"
 *                         },
 *                         "audit_record": [],
 *                         "audit_status": 0,
 *                         "create_time": 1610441917517,
 *                         "evaluate_item_ids": [
 *                                 "KJTNLBV1TC070250A005",
 *                                 "KJTNLBV1TC070250A004"
 *                         ],
 *                         "name": "任务3",
 *                         "op_id": 1,
 *                         "project_id": "5ffd03f14806351af1b696f8",
 *                         "report_audit_record": [],
 *                         "report_audit_status": 0,
 *                         "report_auditor_id": 0,
 *                         "status": 3,
 *                         "task_auditor_id": 0,
 *                         "test_phase": 1,
 *                         "tester_id": 142,
 *                         "update_time": 1610441917517
 *                 }
 *         ]
 * }
 */
func (this EvaluateTaskController) GetProjectTasks(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	params := qmap.QM{
		"e_project_id": req.MustString("project_id"),
	}
	if list, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params).SetLimit(100000).Get(); err == nil {
		response.RenderSuccess(ctx, list)
	} else {
		response.RenderFailure(ctx, err)
	}
}

//@auto_generated_api_end
/**
 * apiType http
 * @api {get} v1/evaluate_task/phase 获取阶段列表
 * @apiVersion 0.1.0
 * @apiName phase
 * @apiGroup EvaluateTask
 *
 * @apiDescription 获取阶段列表
 *
 * @apiParam {string}   project_id  项目ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  "project_id":"5fddcb8824b64731079c2765"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":{
 *				"复测",
 *      		"初测"
 *			}
 *      }
 */
func (this EvaluateTaskController) Phase(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["project_id"] = ctx.Param("project_id")

	phase := map[int]qmap.QM{
		1: {"id": 1, "name": "初测"},
		2: {"id": 2, "name": "复测1"},
		3: {"id": 3, "name": "复测2"},
		4: {"id": 4, "name": "复测3"},
		5: {"id": 5, "name": "复测4"},
	}
	projectId, has := req.TryString("project_id")
	if has == false || projectId == "" {
		response.RenderSuccess(ctx, qmap.QM{"data": custom_util.MapToSlice(phase)})
		return
	}

	if list, err := new(mongo_model.EvaluateTask).GetPhase(projectId); err != nil {
		response.RenderFailure(ctx, err)
		return
	} else {
		result := []qmap.QM{}
		for _, id := range *list {
			if phase[id] != nil {
				result = append(result, phase[id])
			}
		}
		response.RenderSuccess(ctx, result)
		return
	}
}

/**
 * apiType http
 * @api {get} v1/evaluate_task/status 获取状态列表
 * @apiVersion 0.1.0
 * @apiName status
 * @apiGroup EvaluateTask
 *
 * @apiDescription 获取状态列表
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":[
 *				 {
 *					"name": "创建",
 *					"status": 1
 *				}
 *			]
 *      }
 */
func (this EvaluateTaskController) Status(ctx *gin.Context) {
	status := []qmap.QM{
		{"status": 1, "name": "创建"},
		{"status": 2, "name": "任务审核"},
		{"status": 3, "name": "测试"},
		{"status": 4, "name": "报告审核"},
		{"status": 5, "name": "任务完成"},
	}

	response.RenderSuccess(ctx, qmap.QM{"data": status})
}

/**
 * apiType http
 * @api {get} v1/evaluate_task/auditor 获取审核人员列表
 * @apiVersion 0.1.0
 * @apiName auditor
 * @apiGroup EvaluateTask
 *
 * @apiDescription 获取审核人员列表
 *
 * @apiParam {string}   project_id  项目ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  "project_id":"5fddcb8824b64731079c2765"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":[
 *				 {
 *					"id": 117,
 *					"real_name": "王一一"
 *				}
 *			]
 *      }
 */
func (this EvaluateTaskController) Auditor(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["project_id"] = ctx.Param("project_id")

	result := []qmap.QM{}
	if projectId, has := req.TryString("project_id"); has && projectId != "" {
		params := qmap.QM{
			"e__id": bson.ObjectIdHex(req.MustString("project_id")),
		}
		if project, err := mongo.NewMgoSessionWithCond(common.MC_PROJECT, params).GetOne(); err != nil {
			response.RenderFailure(ctx, err)
			return
		} else {
			result = getAuthNames(ctx, (*project)["member_ids"].([]interface{}), common.ROLE_PM)
		}
	} else {
		if userList := getAuthUsers(ctx, common.ROLE_PM); userList != nil {
			for _, item := range userList {
				id := item["id"]
				if int(id.(int)) == int(session.GetUserId(http_ctx.GetHttpCtx(ctx))) {
					continue
				}
				user := qmap.QM{
					"id":        id,
					"real_name": item["realname"],
				}
				result = append(result, user)
			}
		}
	}
	response.RenderSuccess(ctx, result)
	return
}

// 获取人员真实名称
func getAuthNames(ctx *gin.Context, ids []interface{}, roleId int) []qmap.QM {
	result := []qmap.QM{}
	//查询项目所有成员
	users := getAuthUsers(ctx, roleId)
	if users != nil {
		for _, id := range ids {
			for _, user := range users {
				userId := user["id"]
				if int(userId.(float64)) == id {
					item := qmap.QM{
						"id":        id,
						"real_name": user["realname"],
					}
					result = append(result, item)
				}
			}
		}
	} else {
		return nil
	}
	return result
}

// TODO 涉及到auth了
func getAuthUsers(ctx *gin.Context, roleId int) []map[string]interface{} {
	if rsp, err := new(logic.AuthLogic).GetSpecifiedServiceUsers(common.PM_SERVICE, roleId); err == nil {
		return (rsp)
	}
	return nil
}

/**
 * apiType http
 * @api {get} v1/evaluate_task/tester 获取测试团队列表
 * @apiVersion 0.1.0
 * @apiName tester
 * @apiGroup EvaluateTask
 *
 * @apiDescription 获取测试团队列表
 *
 * @apiParam {string}   project_id  项目ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  "project_id":"5fddcb8824b64731079c2765"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":[
 *				 {
 *					"id": 117,
 *					"real_name": "王一一"
 *				}
 *			]
 *      }
 */
func (this EvaluateTaskController) Tester(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["project_id"] = ctx.Param("project_id")

	result := []qmap.QM{}
	if projectId, has := req.TryString("project_id"); has && projectId != "" {
		if list, err := new(mongo_model.EvaluateTask).GetTester(projectId); err == nil {
			result = getAuthNames(ctx, *list, common.ROLE_TEST)
		}
	} else {
		if userList := getAuthUsers(ctx, common.ROLE_TEST); userList != nil {
			for _, item := range userList {
				user := qmap.QM{
					"id":        item["id"],
					"real_name": item["realname"],
				}
				result = append(result, user)
			}
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"data": result})
}

/**
 * apiType http
 * @api {post} v1/evaluate_task/assign 指派审核人或测试团队
 * @apiVersion 0.1.0
 * @apiName Assgin
 * @apiGroup EvaluateTask
 *
 * @apiDescription 指派审核人或测试团队
 *
 * @apiParam {string}   id  任务ID
 * @apiParam {string}   type  指派类型：task_auditor_id任务审核人，report_auditor_id报告审核人,tester_id测试团队
 * @apiParam {int}   user_id  用户ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  "project_id":"5fddcb8824b64731079c2765"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data": true
 *      }
 */
func (this EvaluateTaskController) Assign(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	userList := GetUserList(ctx)
	if result, err := mongo_model.LogicEvaluateTaskAssign(*req, int(session.GetUserId(http_ctx.GetHttpCtx(ctx))), userList); err == nil {
		response.RenderSuccess(ctx, result)
		return
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}

func GetUserList(ctx *gin.Context) map[int]string {
	result := map[int]string{}
	//查询项目所有项目经理
	if rsp, err := new(logic.AuthLogic).GetSpecifiedServiceUsers(common.PM_SERVICE, common.ROLE_PM); err == nil {
		users := rsp
		for _, item := range users {
			result[int(item["id"].(int))] = item["realname"].(string)
		}
	} else {
	}

	//查询项目所有项目经理
	if rsp, err := new(logic.AuthLogic).GetSpecifiedServiceUsers(common.PM_SERVICE, common.ROLE_TEST); err == nil {
		users := rsp
		for _, item := range users {
			result[int(item["id"].(int))] = item["realname"].(string)
		}
	} else {
		fmt.Println(err)
	}

	return result
}

/**
 * apiType http
 * @api {post} v1/evaluate_task/list 获取任务数组
 * @apiVersion 0.1.0
 * @apiName Assgin
 * @apiGroup EvaluateTask
 *
 * @apiDescription 获取任务数组
 *
 * @apiParam {string}   project_id  项目ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  "project_id":"5fddcb8824b64731079c2765"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data": true
 *      }
 */
func (this EvaluateTaskController) GetTaskSlice(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["project_id"] = ctx.Param("project_id")

	projectId, has := req.TryString("project_id")
	result := []qmap.QM{}
	if has == false || projectId == "" {
		uid := session.GetUserId(http_ctx.GetHttpCtx(ctx))
		customOperations := []bson.M{
			{"$match": bson.M{"$or": []bson.M{{"op_id": uid}, {"task_auditor_id": uid}, {"report_auditor_id": uid}, {"tester_id": uid}}}},
		}
		list, err := mongo.NewMgoSession(common.MC_EVALUATE_TASK).QueryGet(customOperations)
		if err == nil {
			for _, item := range *list {
				task := qmap.QM{
					"id":   item["_id"],
					"name": item["name"],
				}
				result = append(result, task)
			}
			response.RenderSuccess(ctx, result)
			return
		}
		response.RenderFailure(ctx, err)
		return
	} else {
		params := qmap.QM{
			"e_project_id": req.MustString("project_id"),
		}
		list, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params).SetLimit(5000).Get()
		if err == nil {
			for _, item := range *list {
				task := qmap.QM{
					"id":   item["_id"],
					"name": item["name"],
				}
				result = append(result, task)
			}
			response.RenderSuccess(ctx, result)
			return
		}
		response.RenderFailure(ctx, err)
		return
	}
}

/**
 * apiType http
 * @api {post} v1/evaluate_task/audit 任务审核报告审核
 * @apiVersion 0.1.0
 * @apiName Assgin
 * @apiGroup EvaluateTask
 *
 * @apiDescription 任务审核报告审核
 *
 * @apiParam {string}   id  任务ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  	"id":"5fddcb8824b64731079c2765"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data": true
 *      }
 */
func (this EvaluateTaskController) Audit(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if result, err := new(mongo_model.EvaluateTask).Audit(*req, int(session.GetUserId(http_ctx.GetHttpCtx(ctx)))); err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {post} v1/evaluate_task/audit 任务审核报告审核
 * @apiVersion 0.1.0
 * @apiName Assgin
 * @apiGroup EvaluateTask
 *
 * @apiDescription 任务审核报告审核
 *
 * @apiParam {string}   id  任务ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  	"id":"5fddcb8824b64731079c2765"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data": true
 *      }
 */
func (this EvaluateTaskController) Submit(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if result, err := new(mongo_model.EvaluateTask).Submit(*req, int(session.GetUserId(http_ctx.GetHttpCtx(ctx)))); err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {post} v1/evaluate_task/task_project 用户分配任务的项目列表
 * @apiVersion 0.1.0
 * @apiName Assgin
 * @apiGroup TaskProject
 *
 * @apiDescription 任务审核报告审核
 *
 * @apiParam {string}   id  任务ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  	"id":"5fddcb8824b64731079c2765"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data": true
 *      }
 */
func (this EvaluateTaskController) TaskProject(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	uid := session.GetUserId(http_ctx.GetHttpCtx(ctx))
	match := bson.M{
		"$or": []bson.M{{"op_id": uid}, {"task_auditor_id": uid}, {"report_auditor_id": uid}, {"tester_id": uid}},
	}
	if isDeleted, has := req.TryInt("is_deleted"); has {
		match["is_deleted"] = isDeleted
	}
	Operations := []bson.M{
		{"$match": match},
	}
	list, err := mongo.NewMgoSession(common.MC_EVALUATE_TASK).QueryGet(Operations)
	projectSlice := mongo_model.GetProjectSlice()
	filterList := qmap.QM{}
	if err == nil {
		result := []qmap.QM{}
		for _, item := range *list {
			name := item["project_id"]
			if projectSlice[item["project_id"].(string)] != nil {
				name = projectSlice[item["project_id"].(string)]
			}

			if filterList[item["project_id"].(string)] != nil {
				continue
			} else {
				filterList[item["project_id"].(string)] = name
			}

			task := qmap.QM{
				"id":   item["project_id"],
				"name": name,
			}
			result = append(result, task)
		}
		response.RenderSuccess(ctx, result)
		return
	}
	response.RenderFailure(ctx, err)
}

/**
 * apiType http
 * @api {post} v1/evaluate_task/my_list 获取用户相关任务列表
 * @apiVersion 0.1.0
 * @apiName Assgin
 * @apiGroup TaskProject
 *
 * @apiDescription 获取用户相关任务列表
 *
 * @apiParam {string}   id  任务ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  	"id":"5fddcb8824b64731079c2765"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data": true
 *      }
 */
func (this EvaluateTaskController) GetList(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if res, err := new(mongo_model.EvaluateTask).GetList(*req, int(session.GetUserId(http_ctx.GetHttpCtx(ctx))), ctx); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, res)
	}
}

/**
 * apiType http
 * @api {post} v1/evaluate_task/task_item_list 获取任务测试用例列表
 * @apiVersion 0.1.0
 * @apiName task_item_list
 * @apiGroup TaskProject
 *
 * @apiDescription 获取任务测试用例列表
 *
 * @apiParam {string}   id  任务ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  	"task_id":"5fddcb8824b64731079c2765"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data": true
 *      }
 */
func (this EvaluateTaskController) GetTaskItemList(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if res, err := new(mongo_model.EvaluateTask).GetTaskItemList(*req, int(session.GetUserId(http_ctx.GetHttpCtx(ctx))), ctx); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, res)
	}
}

/**
 * apiType http
 * @api {post} v1/evaluate_task/task_item_info 获取任务测试用例详情
 * @apiVersion 0.1.0
 * @apiName task_item_info
 * @apiGroup TaskProject
 *
 * @apiDescription 获取任务测试用例详情
 *
 * @apiParam {string}   id  任务ID
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  	"test_id":"5fddcb8824b64731079c2765"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data": true
 *      }
 */
func (this EvaluateTaskController) GetTaskItemInfo(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	testId, hasTestId := req.TryString("id")
	taskId, hasTaskId := req.TryString("evaluate_task_id")
	itemId, hasItemId := req.TryString("item_id")

	taskItemInfo := &qmap.QM{}
	// 如果上传了test_id，直接查出 taskItemInfo
	if hasTestId {
		params := qmap.QM{
			"e__id": bson.ObjectIdHex(testId),
		}
		ormSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK_ITEM, params)
		if info, err := ormSession.GetOne(); err == nil {
			taskItemInfo = info
		}
	} else if hasTaskId && hasItemId { //如果上传的是 task_id和item_id，则使用他们查出taskItemInfo
		params := qmap.QM{
			"e_item_id":          itemId,
			"e_evaluate_task_id": taskId,
		}
		ormSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK_ITEM, params)
		if info, err := ormSession.GetOne(); err == nil {
			taskItemInfo = info
		}
	}

	//taskItemInfo存在的时候，将数据写入结果
	if (*taskItemInfo)["_id"] != nil {
		//调用controller，此处参数为 id
		params := qmap.QM{
			"e__id": (*taskItemInfo)["item_id"],
		}
		mgoSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params)
		mgoSession.SetTransformFunc(EvaluateItemTransformer)
		result, err := mgoSession.GetOne()

		if err == nil {
			(*result)["record_id"] = (*taskItemInfo)["record_id"]
			(*result)["test_id"] = (*taskItemInfo)["_id"]
			(*result)["evaluate_task_id"] = (*taskItemInfo)["evaluate_task_id"]
			(*result)["item_id"] = (*taskItemInfo)["item_id"]
			(*result)["test_status"] = (*taskItemInfo)["status"]
			(*result)["test_phase"] = (*taskItemInfo)["test_phase"]
			(*result)["audit_status"] = (*taskItemInfo)["audit_status"]
			(*result)["record_audit_status"] = (*taskItemInfo)["record_audit_status"]
			response.RenderSuccess(ctx, result)
			return
		}
		response.RenderFailure(ctx, err)
		return
	} else if hasTaskId && hasItemId { //taskItemInfo不存在的时候，将初始化数据写入结果
		params := qmap.QM{
			"e_pre_bind": taskId,
			"e__id":      itemId,
		}
		ormSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params)
		ormSession.SetTransformFunc(EvaluateItemTransformer)
		itemInfo, err := ormSession.GetOne()

		if err == nil {
			(*itemInfo)["test_id"] = ""
			(*itemInfo)["record_id"] = ""
			(*itemInfo)["evaluate_task_id"] = taskId
			(*itemInfo)["item_id"] = (*itemInfo)["_id"]
			(*itemInfo)["test_status"] = common.TIS_READY
			(*itemInfo)["audit_status"] = common.EIAS_DEFAULT
			(*itemInfo)["record_audit_status"] = common.IRAS_DEFAULT
			response.RenderSuccess(ctx, itemInfo)
			return
		}
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderFailure(ctx, errors.New("参数传递错误！"))
}
