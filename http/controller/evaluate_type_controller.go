package controller

import (
	"errors"

	"skygo_detection/guardian/src/net/qmap"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mongo_model"
)

type EvaluateTypeController struct{}

//@auto_generated_api_begin
/**
 * apiType http
 * @api {get} /api/v1/evaluate_type/all 查询所有测试对象类型列表
 * @apiVersion 0.1.0
 * @apiName GetAll
 * @apiGroup EvaluateType
 *
 * @apiDescription 查询所有测试对象类型列表
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 * 	"code": 0,
 * 	"data": [
 * 		{
 * 			"_id": "5e68902d24b6476b74bb3549",
 * 			"attrs": [
 * 				{
 * 					"attr_key": "app_name1",
 * 					"attr_name": "属性名称1",
 * 					"attr_type": "string",
 * 					"is_required": 1
 * 				},
 * 				{
 * 					"attr_key": "app_name2",
 * 					"attr_name": "属性名称2",
 * 					"attr_type": "string",
 * 					"is_required": 0
 * 				}
 * 			],
 * 			"name": "测试项名称"
 * 		}
 * 	]
 * }
 */
func (this EvaluateTypeController) GetAll(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	queryParams := ctx.Request.URL.RawQuery

	mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_TYPE).AddUrlQueryCondition(queryParams)
	if res, err := mgoSession.Get(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_types 分页查询测试对象类型列表
 * @apiVersion 0.1.0
 * @apiName GetPagingAll
 * @apiGroup EvaluateType
 *
 * @apiDescription 分页查询测试对象类型列表
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 * 	"code": 0,
 * 	"data": {
 * 		"list": [
 * 			{
 * 				"_id": "5e68902d24b6476b74bb3549",
 * 				"attrs": [
 * 					{
 * 						"attr_key": "app_name1",
 * 						"attr_name": "属性名称1",
 * 						"attr_type": "string",
 * 						"is_required": 1
 * 					},
 * 					{
 * 						"attr_key": "app_name2",
 * 						"attr_name": "属性名称2",
 * 						"attr_type": "string",
 * 						"is_required": 0
 * 					}
 * 				],
 * 				"name": "测试项名称"
 * 			}
 * 		],
 * 		"pagination": {
 * 			"count": 2,
 * 			"current_page": 1,
 * 			"per_page": 20,
 * 			"total": 2,
 * 			"total_pages": 1
 * 		}
 * 	}
 * }
 */
func (this EvaluateTypeController) GetPagingAll(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	queryParams := ctx.Request.URL.RawQuery

	mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_TYPE).AddUrlQueryCondition(queryParams)
	if res, err := mgoSession.GetPage(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_types/:id 查询某一测试对象类型信息
 * @apiVersion 0.1.0
 * @apiName GetOne
 * @apiGroup EvaluateType
 *
 * @apiDescription 根据id查询某一测试对象类型信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   id  		测试对象类型id
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 * 	"code": 0,
 * 	"data": {
 * 		"_id": "5e68902d24b6476b74bb3549",
 * 		"attrs": [
 * 			{
 * 				"attr_key": "app_name1",
 * 				"attr_name": "属性名称1",
 * 				"attr_type": "string",
 * 				"is_required": 1
 * 			},
 * 			{
 * 				"attr_key": "app_name2",
 * 				"attr_name": "属性名称2",
 * 				"attr_type": "string",
 * 				"is_required": 0
 * 			}
 * 		],
 * 		"name": "测试项名称"
 * 	}
 * }
 */
func (this EvaluateTypeController) GetOne(ctx *gin.Context) {
	id := ctx.Param("id")

	params := qmap.QM{
		"e__id": bson.ObjectIdHex(id),
	}
	ormSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TYPE, params)
	data, _ := ormSession.GetOne()
	response.RenderSuccess(ctx, data)
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_types 添加测试对象类型
 * @apiVersion 0.1.0
 * @apiName Create
 * @apiGroup EvaluateType
 *
 * @apiDescription 添加测试对象类型
 *
 * @apiUse authHeader
 *
 * @apiUse authHeader
 *
 * @apiParam {string}      	                        name    			    测试对象类型名称
 * @apiParam {string}      	                        attrs.name    		    属性名称
 * @apiParam {string}      	                        attrs.attr_key          属性关键字
 * @apiParam {string=string,int,float,date,bool}		attrs.attr_type		    属性类型(字符串:string,整形:int, 浮点型:float, 日期:date, 布尔型:bool)
 * @apiParam {int=0,1}      	                        attrs.is_required  	    是否必填(0:否，1:是)
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 * 	"name":"测试项名称",
 * 	"attrs":[
 * 		{
 * 			"attr_name":"属性名称1",
 * 			"attr_key":"app_name1",
 * 			"attr_type":"string",
 * 			"is_required":1
 * 		}
 * 	]
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 * 	"code": 0,
 * 	"data": {
 * 		"attrs": [
 * 			{
 * 				"AttrKey": "app_name1",
 * 				"AttrName": "属性名称1",
 * 				"AttrType": "string",
 * 				"IsRequired": 1
 * 			}
 * 		],
 * 		"id": "5e688f7a24b6476b74bb3548",
 * 		"name": "测试项名称"
 * 	}
 * }
 */
func (this EvaluateTypeController) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	name := req.MustString("name")
	params := qmap.QM{
		"e_name": name,
	}
	ormSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TYPE, params)
	data, _ := ormSession.GetOne()
	if len(*data) > 0 {
		response.RenderFailure(ctx, errors.New("资产类型名不能重复"))
		return
	}

	if eType, err := new(mongo_model.EvaluateType).Create(req); err == nil {
		retCols := map[string]bool{
			"Name": true,
			"Id":   true,
		}
		ret := custom_util.StructToMapWithColumns(*eType, retCols)
		ret["attrs"] = eType.Attrs
		response.RenderSuccess(ctx, ret)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {put} /api/v1/evaluate_types/:id  更新测试对象类型
 * @apiVersion 0.1.0
 * @apiName Update
 * @apiGroup EvaluateType
 *
 * @apiDescription 根据id,更新测试对象类型信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}    	                        id                      测试对象类型id
 * @apiParam {string}      	                        name    			    测试项名称
 * @apiParam {string}      	                        attrs.name    		    属性名称
 * @apiParam {string}      	                        attrs.attr_key          属性关键字
 * @apiParam {string=string,int,float,date,bool}    attrs.attr_type		    属性类型(字符串:string,整形:int, 浮点型:float, 日期:date, 布尔型:bool)
 * @apiParam {int=0,1}      	                    attrs.is_required  	    是否必填(0:否，1:是)
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  "id":"5e688f7a24b6476b74bb3548",
 * 	"name":"测试项名称",
 * 	"attrs":[
 * 		{
 * 			"attr_name":"属性名称1",
 * 			"attr_key":"app_name1",
 * 			"attr_type":"string",
 * 			"is_required":1
 * 		}
 * 	]
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 * 	"code": 0,
 * 	"data": {
 * 		"attrs": [
 * 			{
 * 				"AttrKey": "app_name1",
 * 				"AttrName": "属性名称1",
 * 				"AttrType": "string",
 * 				"IsRequired": 1
 * 			}
 * 		],
 * 		"id": "5e688f7a24b6476b74bb3548",
 * 		"name": "测试项名称"
 * 	}
 * }
 */
func (this EvaluateTypeController) Update(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	id := ctx.Param("id")

	//检查evaluate_type是否被项目资产使用，如果使用不能修改
	typeInfo, err := mongo.NewMgoSession(common.MC_EVALUATE_TYPE).AddCondition(qmap.QM{"e__id": bson.ObjectIdHex(id)}).GetOne()
	if err != nil {
		response.RenderFailure(ctx, errors.New("该模板不存在"))
		return
	}
	cond := qmap.QM{
		"e_evaluate_type": (*typeInfo)["name"],
	}
	info, err := mongo.NewMgoSession(common.MC_EVALUATE_ASSET).AddCondition(cond).GetOne()
	if err == nil && (*info)["_id"] != nil {
		response.RenderFailure(ctx, errors.New("该模板已被项目资产使用"))
		return
	}

	updateCols := map[string]string{
		"name":  "string",
		"attrs": "interface",
	}
	rawInfo := custom_util.CopyMapColumns(*req, updateCols)
	if eType, err := new(mongo_model.EvaluateType).Update(id, rawInfo); err == nil {
		retCols := map[string]bool{
			"Name": true,
			"Id":   true,
		}
		ret := custom_util.StructToMapWithColumns(*eType, retCols)
		ret["attrs"] = eType.Attrs
		response.RenderSuccess(ctx, ret)
		return
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/evaluate_types 批量删除测试对象类型
 * @apiVersion 0.1.0
 * @apiName BulkDelete
 * @apiGroup EvaluateType
 *
 * @apiDescription 批量删除测试对象类型
 *
 * @apiUse authHeader
 *
 * @apiParam {[]string}   ids  测试项id
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
func (this EvaluateTypeController) BulkDelete(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	effectNum := 0
	if rawIds, has := req.TrySlice("ids"); has {
		ids := []bson.ObjectId{}
		for _, id := range rawIds {
			//检查evaluate_type是否被项目资产使用，如果使用不能删除
			typeInfo, err := mongo.NewMgoSession(common.MC_EVALUATE_TYPE).AddCondition(qmap.QM{"e__id": bson.ObjectIdHex(id.(string))}).GetOne()
			if err != nil {
				response.RenderFailure(ctx, errors.New("该模板不存在"))
				return
			}

			cond := qmap.QM{
				"e_evaluate_type": (*typeInfo)["name"],
			}
			info, err := mongo.NewMgoSession(common.MC_EVALUATE_ASSET).AddCondition(cond).GetOne()
			if err == nil && (*info)["_id"] != nil {
				response.RenderFailure(ctx, errors.New("该模板已被项目资产使用"))
				return
			}

			ids = append(ids, bson.ObjectIdHex(id.(string)))
		}
		if len(ids) > 0 {
			match := bson.M{
				"_id": bson.M{"$in": ids},
			}
			if changeInfo, err := mongo.NewMgoSession(common.MC_EVALUATE_TYPE).RemoveAll(match); err == nil {
				effectNum = changeInfo.Removed
			} else {
				response.RenderFailure(ctx, err)
				return
			}
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"number": effectNum})
}

//@auto_generated_api_end
