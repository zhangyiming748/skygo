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
	"skygo_detection/mongo_model"
)

type EvaluateAssetController struct{}

//@auto_generated_api_begin
/**
 * apiType http
 * @api {get} /api/v1/evaluate_assets 查询所有资产列表
 * @apiVersion 0.1.0
 * @apiName GetAll
 * @apiGroup EvaluateAsset
 *
 * @apiDescription 查询所有资产列表
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 * 	"code": 0,
 * 	"data": [
 * 		{
 * 			"_id": "5e6f5a03e85e913a09907011",
 * 			"attributes": {
 * 				"confidential": false,
 * 				"count": 223,
 * 				"date": 1584979200000,
 * 				"price": 3,
 * 				"version": "dasdasd"
 * 			},
 * 			"create_time": 1584355843638,
 * 			"evaluate_type": "IVI",
 * 		    "evaluate_item_total": 2,
 * 			"name": "fasfasf",
 * 			"project_id": "5e6f5a0324b6471290457cad",
 * 			"update_time": 1584355843638
 * 		}
 * 	]
 * }
 */
func (this EvaluateAssetController) GetAll(ctx *gin.Context) {
	// 接收参数 ?where={"project_id":{"e":"60f677e3e830c649001ef13c"}}
	queryParams := ctx.Request.URL.RawQuery
	widget := orm_mongo.NewWidgetWithCollectionName(common.MC_EVALUATE_ASSET).
		SetQueryStr(queryParams).
		SetLimit(50000).
		SetTransformerFunc(EvaluateAssetTransformer)
	res, _ := widget.Find()
	response.RenderSuccess(ctx, res)
}

/**
    * apiType http
    * @api {get} /api/v1/evaluate_assets/:id 查询某一资产信息
    * @apiVersion 0.1.0
    * @apiName GetOne
    * @apiGroup EvaluateAsset
    *
    * @apiDescription 根据id查询某一资产信息
    *
    * @apiUse authHeader
    *
    * @apiParam {string}   id  		资产id
    *
    * @apiSuccessExample {json} 请求成功示例:
	* {
	* 	"code": 0,
	* 	"data": {
	* 		"_id": "5e6f5a03e85e913a09907011",
	* 		"attributes": {
	* 			"confidential": false,
	* 			"count": 223,
	* 			"date": 1584979200000,
	* 			"price": 3,
	* 			"version": "dasdasd"
	* 		},
	* 		"create_time": 1584355843638,
	* 		"evaluate_type": "IVI",
	* 		"evaluate_item_total": 2,
	* 		"name": "fasfasf",
	* 		"project_id": "5e6f5a0324b6471290457cad",
	* 		"update_time": 1584355843638
	* 	}
	* }
*/
func (this EvaluateAssetController) GetOne(ctx *gin.Context) {
	id := ctx.Param("id")

	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if taskId, has := req.TryString("task_id"); has { // 查询任务资产信息
		info, err := new(mongo_model.EvaluateAsset).GetTaskAssetInfo(id, taskId)
		if err == nil {
			checkErr := mongo_model.CheckIsProjectUser((info)["project_id"].(string), session.GetUserId(http_ctx.GetHttpCtx(ctx)))
			if checkErr == nil {
				response.RenderSuccess(ctx, info)
				return
			}
			response.RenderFailure(ctx, checkErr)
			return
		}
		response.RenderFailure(ctx, err)
		return
	} else { // 查询资产信息
		info, err := new(mongo_model.EvaluateAsset).GetOne(id)
		if err == nil {
			checkErr := mongo_model.CheckIsProjectUser((info)["project_id"].(string), session.GetUserId(http_ctx.GetHttpCtx(ctx)))
			if checkErr == nil {
				response.RenderSuccess(ctx, info)
				return
			}
			response.RenderFailure(ctx, checkErr)
			return
		}
		response.RenderFailure(ctx, err)
		return
	}
}

/**
    * apiType http
    * @api {post} /api/v1/evaluate_assets 添加资产
    * @apiVersion 0.1.0
    * @apiName Create
    * @apiGroup EvaluateAsset
    *
    * @apiDescription 添加资产
    *
    * @apiUse authHeader
    *
    * @apiUse authHeader
    *
    * @apiParam {string}      	                        name    			    资产名称
    * @apiParam {string}      	                       	project_id    		    项目id
    * @apiParam {string}      	                       	evaluate_type    		资产类型
    *
    * @apiParamExample {json}  请求参数示例:
	* {
	*		"project_id":"5e6f5a0324b6471290457cad",
	*		"name":"test",
	*		"evaluate_type":"ivi",
	*		"count":1,
	*		"version": "1.1"
	* }
	*
	* @apiSuccessExample {json} 请求成功示例:
	* {
	* 	"code": 0,
	* 	"data": {
	* 		"attributes": {
	* 			"confidential": false,
	* 			"count": 1,
	* 			"date": 0,
	* 			"price": 0,
	* 			"version": "1.1"
	* 		},
	* 		"create_time": 1584429025499,
	* 		"evaluate_type": "ivi",
	* 		"id": "5e7077e124b64714ba3ed666",
	* 		"name": "test",
	* 		"op_id": 0,
	* 		"project_id": "5e6f5a0324b6471290457cad",
	* 		"update_time": 1584429025499
	* 	}
	* }
*/
func (this EvaluateAssetController) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if evaluateAsset, err := new(mongo_model.EvaluateAsset).Create(context.TODO(), req.MustString("project_id"), session.GetUserId(http_ctx.GetHttpCtx(ctx)), *req); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, custom_util.StructToMap2(*evaluateAsset))
	}
}

/**
    * apiType http
    * @api {put} /api/v1/evaluate_assets/:id  更新资产
    * @apiVersion 0.1.0
    * @apiName Update
    * @apiGroup EvaluateTarget
    *
    * @apiDescription 根据项目id,更新资产信息
    *
    * @apiUse authHeader
    *
    * @apiParam {string}    							id                   	资产id
    * @apiParam {string}      	                        name    			    资产名称
    *
    * @apiParamExample {json}  请求参数示例:
	 * {
	 * 		"id": "5e7077e124b64714ba3ed666",
	 * 		"name": "资产名称",
	 *		"count":1
	 * }
	 *
	 * @apiSuccessExample {json} 请求成功示例:
	 * {
	 * 	"code": 0,
	 * 	"data": {
	 * 		"attributes": {
	 * 			"confidential": false,
	 * 			"count": 1,
	 * 			"date": 0,
	 * 			"price": 0,
	 * 			"version": "1.1"
	 * 		},
	 * 		"create_time": 1584429025499,
	 * 		"evaluate_type": "ivi",
	 * 		"id": "5e7077e124b64714ba3ed666",
	 * 		"name": "test",
	 * 		"op_id": 0,
	 * 		"project_id": "5e6f5a0324b6471290457cad",
	 * 		"update_time": 1584429025499
	 * 	}
	 * }
*/
func (this EvaluateAssetController) Update(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	id := ctx.Param("id")

	if evaluateAsset, err := new(mongo_model.EvaluateAsset).Update(context.TODO(), id, session.GetUserId(http_ctx.GetHttpCtx(ctx)), *req); err == nil {
		response.RenderSuccess(ctx, custom_util.StructToMap2(*evaluateAsset))
	} else {
		panic(err)
	}

}

/**
 * apiType http
 * @api {delete} /api/v1/evaluate_asset 批量删除资产
 * @apiVersion 0.1.0
 * @apiName BulkDelete
 * @apiGroup EvaluateAsset
 *
 * @apiDescription 批量删除资产
 *
 * @apiUse authHeader
 *
 * @apiParam {[]string}   ids  资产id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  "ids":[
 *		    "5e688f7a24b6476b74bb3548"
 * 	]
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *          "code": 0
 *			 "data":{
 *				"number":1
 *			}
 *      }
 */
func (this EvaluateAssetController) BulkDelete(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	effectNum := 0
	if rawIds, has := req.TrySlice("ids"); has {
		ids := []string{}
		for _, id := range rawIds {
			ids = append(ids, id.(string))
			//检查每个任务是否可以删除
			params := qmap.QM{
				"e_asset_id": id,
			}
			if info, err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_ITEM, params).Get(); err == nil {
				if info["pre_bind"] != "" || info["test_count"].(int32) > 0 {
					response.RenderFailure(ctx, errors.New("存在已分配任务的测试用例，不可删除"))
					return
				}
			}
			if _, err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_MATERIEL, params).Get(); err == nil {
				response.RenderFailure(ctx, errors.New(fmt.Sprintf("存在相关的物料，无法删除当前资产： %s", id)))
				return
			}
		}

		// 删除资产
		if len(ids) > 0 {
			match := bson.M{
				"_id": bson.M{"$in": ids},
			}
			coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_ASSET)
			if changeInfo, err := coll.DeleteMany(context.Background(), match); err == nil {
				effectNum = int(changeInfo.DeletedCount)
			} else {
				response.RenderFailure(ctx, err)
				return
			}

			// 删除对应的测试用例
			match2 := bson.M{
				"asset_id": bson.M{"$in": ids},
			}
			collI := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_ITEM)
			if changeInfo, err := collI.DeleteMany(context.Background(), match2); err == nil {
				effectNum = int(changeInfo.DeletedCount)
			} else {
				response.RenderFailure(ctx, err)
				return
			}
		}
	}

	response.RenderSuccess(ctx, &qmap.QM{"number": effectNum})
	return
}

/**
 * apiType http
 * @api {delete} /api/v1/evaluate_asset/module_type 删除资产测试组件分类
 * @apiVersion 0.1.0
 * @apiName DeleteModuleType
 * @apiGroup EvaluateAsset
 *
 * @apiDescription 删除资产测试组件分类
 *
 * @apiUse authHeader
 *
 * @apiParam {[]string}   ids  资产id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  "ids":[
 *		    "5e688f7a24b6476b74bb3548"
 * 	]
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *          "code": 0
 *			 "data":{
 *				"number":1
 *			}
 *      }
 */
func (this EvaluateAssetController) DeleteModuleType(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	assetId := req.MustString("id")
	moduleTypeIds := req.SliceString("module_type_id")
	// 判断相关的测试用例是否已经绑定任务，如果已经绑定则不能删除
	if _, err := new(mongo_model.EvaluateAsset).CheckModuleTypeIdCanBeDelete(assetId, moduleTypeIds); err != nil {
		response.RenderFailure(ctx, err)
		return
	}

	if _, err := new(mongo_model.EvaluateAsset).DeleteModuleType(assetId, *req); err == nil {
		// 删除对应的测试用例
		if _, err := new(mongo_model.EvaluateAsset).DeleteItemByModuleTypeIds(assetId, moduleTypeIds); err != nil {
			response.RenderFailure(ctx, err)
			return
		}
		response.RenderSuccess(ctx, qmap.QM{"result": true})
		return
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_asset/type_asset/:project_id 获取资产分类及资产名称
 * @apiVersion 0.1.0
 * @apiName TypeAsset
 * @apiGroup EvaluateAsset
 *
 * @apiDescription 获取资产分类及资产
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   project_id  项目id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *       "project_id":"5e688f7a24b6476b74bb3548"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *          "code": 0
 *			 "data":{
 *				"number":1
 *			}
 *      }
 */
func (this EvaluateAssetController) TypeAsset(ctx *gin.Context) {
	projectId := ctx.Param("project_id")

	result, err := new(mongo_model.EvaluateAsset).TypeAsset(projectId)
	if err != nil {
		panic(err)
	} else {
		response.RenderSuccess(ctx, result)
	}
}

//@auto_generated_api_end
/**
 * apiType http
 * @api {get} /api/v1/evaluate_asset/task_asset/:task_id 获取任务资产
 * @apiVersion 0.1.0
 * @apiName TaskAsset
 * @apiGroup EvaluateAsset
 *
 * @apiDescription 获取资产分类及资产
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   task_id  任务id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *       "task_id":"US393L24"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *          "code": 0
 *			 "data":{
 *				"number":1
 *			}
 *      }
 */
func (this EvaluateAssetController) TaskAsset(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	taskId := ctx.Param("task_id")

	params := qmap.QM{
		"e__id": bson.ObjectIdHex(taskId),
	}
	assetId, hasAssetId := req.TryString("asset_id")
	result := map[interface{}]interface{}{}
	sorted := []interface{}{}
	info, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params).GetOne()
	// 统计资产用例数
	itemCompleteTotal := map[string]int{}
	itemTotal := map[string]int{}
	moduleIds := map[string][]string{}
	for _, itemId := range (*info)["evaluate_item_ids"].([]interface{}) {
		params := qmap.QM{
			"e__id": itemId,
		}
		itemInfo, _ := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params).GetOne()

		assetModuleIds := moduleIds[(*itemInfo)["asset_id"].(string)]
		itemModuleTypeId := (*itemInfo)["module_type_id"].(string)
		if !custom_util.IndexOfSlice(itemModuleTypeId, assetModuleIds) {
			moduleIds[(*itemInfo)["asset_id"].(string)] = append(assetModuleIds, itemModuleTypeId)
		}

		params = qmap.QM{
			"e_item_id":          itemId,
			"e_evaluate_task_id": taskId,
		}
		taskItemInfo, _ := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK_ITEM, params).GetOne()

		if (*taskItemInfo)["status"] == 1 || (*taskItemInfo)["status"] == 3 {
			itemCompleteTotal[(*itemInfo)["asset_id"].(string)]++
		}
		itemTotal[(*itemInfo)["asset_id"].(string)]++
	}

	moduleMap := GetModuleMap()
	moduleName := map[string][]string{}
	moduleId := map[string][]string{}
	for assetId, itemModuleTypeIds := range moduleIds {
		for _, id := range itemModuleTypeIds {
			if moduleMap[id] != nil && !custom_util.IndexOfSlice(moduleMap[id].(string), moduleName[assetId]) {
				moduleName[assetId] = append(moduleName[assetId], moduleMap[id].(string))
				moduleId[assetId] = append(moduleId[assetId], id)
			}
		}
	}
	if err == nil {
		for id, version := range (*info)["asset_versions"].(map[string]interface{}) {
			if hasAssetId && assetId != id {
				continue
			}

			assetParams := qmap.QM{
				"e__id": id,
			}
			if evaluateType, has := req.TryString("evaluate_type"); has {
				assetParams["e_evaluate_type"] = evaluateType
			}
			mgoSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, assetParams)
			asset, _ := mgoSession.GetOne()
			if (*asset) != nil {
				(*asset)["version"] = version
				(*asset)["id"] = id
				(*asset)["module_type_id"] = moduleIds[(*asset)["_id"].(string)]

				selectStruct := [][]interface{}{}
				for _, selectItem := range (*asset)["select_struct"].([]interface{}) {
					item := selectItem.([]interface{})
					if custom_util.IndexOfSlice(item[1].(string), moduleIds[(*asset)["_id"].(string)]) {
						selectStruct = append(selectStruct, selectItem.([]interface{}))
					}
				}
				(*asset)["select_struct"] = selectStruct

				//已完成测试用例数
				(*asset)["test_case_complete_total"] = itemCompleteTotal[(*asset)["_id"].(string)]

				//总测试用例数
				(*asset)["test_case_total"] = itemTotal[(*asset)["_id"].(string)]
				(*asset)["module_total"] = len(moduleName[(*asset)["_id"].(string)])
				(*asset)["module_type_total"] = len(moduleIds[(*asset)["_id"].(string)])

				//如果指定了资产ID，只返回指定的资产信息
				moduleTypeSlice := []interface{}{}
				if hasAssetId {
					for _, modutypeID := range moduleId[(*asset)["_id"].(string)] {
						moduleType, _ := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, qmap.QM{"e__id": bson.ObjectIdHex(modutypeID)}).GetOne()
						moduleTypeSlice = append(moduleTypeSlice, moduleType)
					}

					(*asset)["module_type"] = moduleTypeSlice
					response.RenderSuccess(ctx, asset)
					return
				} else {
					(*asset)["module_type"] = moduleName[(*asset)["_id"].(string)]
				}

				//按照创建时间排序
				key := int((*asset)["create_time"].(int64))
				if result[key] == nil {
					result[key] = (*asset)
				}
				sorted = custom_util.KSort(result)
			}
		}

		response.RenderSuccess(ctx, custom_util.ArrayReverse(sorted))
		return
	}
	panic(err)
}

func EvaluateAssetTransformer(data qmap.QM) qmap.QM {
	assetId := data.MustString("_id")
	projectId := data.MustString("project_id")
	params := qmap.QM{
		"e_project_id": projectId,
		"e_asset_id":   assetId,
	}

	itemTotal := 0
	itemCompleteTotal := 0
	moduleTypeIds := []string{}
	if itemList, err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).AddCondition(params).SetLimit(50000).Get(); err == nil {
		for _, item := range *itemList {
			itemTotal++
			if item["test_status"] == 1 || item["test_status"] == 3 {
				itemCompleteTotal++
			}
			//统计测试组件和分类
			if !custom_util.InArray(item["module_type_id"].(string), moduleTypeIds) {
				moduleTypeIds = append(moduleTypeIds, item["module_type_id"].(string))
			}
		}
	}

	//已完成测试用例数
	data["test_case_complete_total"] = itemCompleteTotal

	//总测试用例数
	data["test_case_total"] = itemTotal

	//统计实际测试项所述的测试组件和分类数
	if len(moduleTypeIds) > 0 {
		moduleMap := GetModuleMap()
		module := []string{}
		for _, id := range moduleTypeIds {
			if moduleMap[id] != nil && !custom_util.IndexOfSlice(moduleMap[id].(string), module) {
				module = append(module, moduleMap[id].(string))
			}
		}
		data["module_type"] = module
		data["module_total"] = len(module)
		data["module_type_total"] = len(moduleTypeIds)

	} else {
		data["module_total"] = 0
		data["module_type_total"] = 0
	}

	return data
}
