package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mongo_model"
)

type EvaluateMaterialController struct{}

//@auto_generated_api_begin
/**
 * apiType http
 * @api {get} /api/v1/evaluate_materials 查询所有物料列表
 * @apiVersion 0.1.0
 * @apiName GetAll
 * @apiGroup EvaluateMaterial
 *
 * @apiDescription 查询所有物料列表
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "list": [
 *             {
 *                 "_id": "5feae0ef5d50260465e4718a",
 *                 "asset_name": "as",
 *                 "comment": "as",
 *                 "create_time": 1609228527927,
 *                 "image": "5fe86035b1bebf0007852f23",
 *                 "name": "as",
 *                 "number": 1,
 *                 "project_id": "5fd7218624b64712a27f47e8"
 *             },
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
func (this EvaluateMaterialController) GetAll(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	queryParams := ctx.Request.URL.RawQuery

	mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_MATERIEL).AddUrlQueryCondition(queryParams)
	mgoSession.SetTransformFunc(EvaluateMaterialTransformer)
	if res, err := mgoSession.GetPage(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func EvaluateMaterialTransformer(data qmap.QM) qmap.QM {
	assetId := data.MustString("asset_id")
	asset, _ := new(mongo_model.EvaluateAsset).GetOne(assetId)
	assetName := asset.String("name")
	evaluateType := asset.String("evaluate_type")
	data["asset_name"] = assetName
	data["evaluate_type"] = evaluateType
	return data
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_materials/:id 查询某一物料信息
 * @apiVersion 0.1.0
 * @apiName GetOne
 * @apiGroup EvaluateMaterial
 *
 * @apiDescription 根据物料id查询某一物料信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   id  		物料id
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *    "code": 0,
 *    "data": {
 *        "_id": "5feae0ef5d50260465e4718a",
 *        "asset_name": "as",
 *        "comment": "as",
 *        "create_time": 1609228527927,
 *        "image": "5fe86035b1bebf0007852f23",
 *        "name": "as",
 *        "number": 1,
 *        "project_id": "5fd7218624b64712a27f47e8"
 *    }
 * }
 */
func (this EvaluateMaterialController) GetOne(ctx *gin.Context) {
	id := ctx.Param("id")

	params := qmap.QM{
		"e__id": bson.ObjectIdHex(id),
	}
	ormSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MATERIEL, params)
	ormSession.SetTransformFunc(EvaluateMaterialTransformer)
	data, _ := ormSession.GetOne()
	response.RenderSuccess(ctx, data)
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_materials 添加物料
 * @apiVersion 0.1.0
 * @apiName Create
 * @apiGroup EvaluateMaterial
 *
 * @apiDescription 添加物料
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           project_id    			项目id
 * @apiParam {string}           name    				物料名称
 * @apiParam {int}           	number           		设备数量
 * @apiParam {int}           	asset_name              资产名称
 * @apiParam {string}           comment           		备注
 * @apiParam {string}           image				 	图片ID
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "project_id":"5e61f95024b64748d37d8cc6",
 *          "name":"物料名称",
 *          "number": 2,
 *          "asset_name": "资产名称",
 *          "comment":"v1123",
 *          "image":"as2h123h1h2j31jj23h13"
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *          "id":"5e61f95024b64748d3711111",
 *          "project_id":"5e61f95024b64748d37d8cc6",
 *          "name":"物料名称",
 *          "number": 2,
 *          "asset_name": "资产名称",
 *          "comment":"v1123",
 *          "image":"as2h123h1h2j31jj23h13"
 *     }
 * }
 */
func (this EvaluateMaterialController) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	//project记录必须存在
	projectId := req.MustString("project_id")
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(projectId),
	}
	project, err := mongo.NewMgoSessionWithCond(common.MC_PROJECT, params).GetOne()
	if err != nil || project == nil {
		panic("project不能为空")
	}

	//上传的二进制图片，要通过GridFs存，image字段就是一个string

	if evaluateMateriel, err := new(mongo_model.EvaluateMateriel).Create(req); err == nil {
		retCols := map[string]bool{
			"Name":      true,
			"Number":    true,
			"AssetId":   true,
			"Comment":   true,
			"Image":     true,
			"ProjectId": true,
		}
		ret := custom_util.StructToMapWithColumns(*evaluateMateriel, retCols)
		response.RenderSuccess(ctx, ret)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {put} /api/v1/evaluate_materials/:id  更新物料
 * @apiVersion 0.1.0
 * @apiName Update
 * @apiGroup EvaluateMaterial
 *
 * @apiDescription 根据资产id,更新物料信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           project_id    			项目id
 * @apiParam {string}           name    				物料名称
 * @apiParam {int}           	number           		设备数量
 * @apiParam {int}           	asset_name              资产名称
 * @apiParam {string}           comment           		备注
 * @apiParam {string}           image				 	图片ID
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "id":"5e61f95024b64748d37d8cc6",
 *          "name":"物料名称",
 *          "number": 2,
 *          "asset_name": "资产名称",
 *          "comment":"备注",
 *          "image":"as2h123h1h2j31jj23h13"
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *          "id":"5e61f95024b64748d37d8cc6",
 *          "project_id":"5e61f95024b64748d37d8cc6",
 *          "name":"物料名称",
 *          "number": 2,
 *          "asset_name": "资产名称",
 *          "comment":"备注",
 *          "image":"as2h123h1h2j31jj23h13"
 *     }
 * }
 */
func (this EvaluateMaterialController) Update(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	id := ctx.Param("id")

	updateCols := map[string]string{
		"name":     "string",
		"number":   "int",
		"asset_id": "string",
		"comment":  "string",
		"image":    "string",
	}
	rawInfo := custom_util.CopyMapColumns(*req, updateCols)
	if Materiel, err := new(mongo_model.EvaluateMateriel).Update(id, rawInfo); err == nil {
		//retCols := map[string]bool{
		//	"Name":                true,
		//	"Company":             true,
		//	"StartTime":           true,
		//	"EndTime":             true,
		//	"EvaluateRequirement": true,
		//	"Description":         true,
		//	"ManagerId":           true,
		//	"CreateTime":          true,
		//}
		//ret := util.StructToMapWithColumns(*project, retCols)
		response.RenderSuccess(ctx, Materiel)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/evaluate_materials 批量删除测试物料
 * @apiVersion 0.1.0
 * @apiName BulkDelete
 * @apiGroup EvaluateMaterial
 *
 * @apiDescription 批量删除测试物料
 *
 * @apiUse authHeader
 *
 * @apiParam {[]string}   ids  测试物料id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  	"ids":[
 *		    "5e688f7a24b6476b74bb3548"
 * 		]
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
func (this EvaluateMaterialController) BulkDelete(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	effectNum := 0
	if rawIds, has := req.TrySlice("ids"); has {
		ids := []bson.ObjectId{}
		for _, id := range rawIds {
			ids = append(ids, bson.ObjectIdHex(id.(string)))
		}
		if len(ids) > 0 {
			match := bson.M{
				"_id": bson.M{"$in": ids},
			}
			if changeInfo, err := mongo.NewMgoSession(common.MC_EVALUATE_MATERIEL).RemoveAll(match); err == nil {
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

/**
 * apiType http
 * @api {get} /api/v1/evaluate_material/task_id/:id 查询任务下所有物料列表
 * @apiVersion 0.1.0
 * @apiName GetAll
 * @apiGroup EvaluateMaterial
 *
 * @apiDescription 查询所有物料列表
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "list": [
 *             {
 *                 "_id": "5feae0ef5d50260465e4718a",
 *                 "asset_name": "as",
 *                 "comment": "as",
 *                 "create_time": 1609228527927,
 *                 "image": "5fe86035b1bebf0007852f23",
 *                 "name": "as",
 *                 "number": 1,
 *                 "project_id": "5fd7218624b64712a27f47e8"
 *             },
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
func (this EvaluateMaterialController) GetAllWithTaskId(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	(*req)["id"] = ctx.Param("id")

	evaluateMateriel := new(mongo_model.EvaluateMateriel)
	data, err := evaluateMateriel.GetAllWithTaskId(req)
	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, data)
	}
}
