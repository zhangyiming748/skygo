package controller

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/mysql_model"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mongo_model"
)

type EvaluateModuleController struct{}

// @auto_generated_api_begin
/**
 * apiType http
 * @api {get} /api/v1/evaluate_module/all 查询测试组件树
 * @apiVersion 0.1.0
 * @apiName GetAllModuleTree
 * @apiGroup EvaluateModule
 *
 * @apiDescription 查询测试组件树
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "list": [
 *                 {
 *                     "_id": "5fd87545aee3d1849a56efa7",
 *                     "module_name": "无线电",
 *                     "module_name_code": "090",
 *                     "module_type": "GNSS",
 *                     "module_type_code": "520"
 *                 },
 *                 {
 *                     "_id": "5fd87545aee3d1849a56ef9d",
 *                     "module_name": "无线电",
 *                     "module_name_code": "090",
 *                     "module_type": "蜂窝网络",
 *                     "module_type_code": "510"
 *                 }
 *             ],
 *             "module_name": "无线电",
 *             "module_name_code": 90
 *         }
 *     ]
 * }
 */
func (this EvaluateModuleController) GetAllModuleTree(ctx *gin.Context) {
	list := []interface{}{}
	// moduleNames := []string{}
	collM := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_MODULE)
	if moduleNames, err := collM.Distinct(context.Background(), "module_name", bson.M{}); err == nil {
		w := orm_mongo.NewWidgetWithCollectionName(common.MC_EVALUATE_MODULE).SetLimit(100000)
		for _, moduleName := range moduleNames {
			w.SetParams(qmap.QM{"e_module_name": moduleName})
			w.SetTransformerFunc(UsedTransformer)
			if res, err := w.Find(); err == nil {
				if len(res) > 0 {
					var tmp qmap.QM = res[0]
					item := qmap.QM{
						"module_name":      moduleName,
						"module_name_code": tmp.Int("module_name_code"),
						"list":             res,
					}
					list = append(list, item)
				}
			} else {
				response.RenderFailure(ctx, err)
				return
			}
		}
	} else {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, list)
}

func UsedTransformer(data qmap.QM) qmap.QM {
	data["used"] = 0
	id := data["_id"].(primitive.ObjectID).Hex()
	models := make([]mysql_model.KnowledgeTestCase, 0)
	mysql.GetSession().Where("module_id = ?", id).Find(&models)
	if len(models) > 0 {
		data["used"] = 1
	}
	return data
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_modules 查询所有测试组件
 * @apiVersion 0.1.0
 * @apiName GetAll
 * @apiGroup EvaluateModule
 *
 * @apiDescription 查询所有测试组件
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *          {
 *                     "_id": "5fd87545aee3d1849a56efa7",
 *                     "module_name": "无线电",
 *                     "module_name_code": "090",
 *                     "module_type": "GNSS",
 *                     "module_type_code": "520"
 *                 },
 *                 {
 *                     "_id": "5fd87545aee3d1849a56ef9d",
 *                     "module_name": "无线电",
 *                     "module_name_code": "090",
 *                     "module_type": "蜂窝网络",
 *                     "module_type_code": "510"
 *                 }
 *     ]
 * }
 */
func (this EvaluateModuleController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery

	w := orm_mongo.NewWidgetWithCollectionName(common.MC_EVALUATE_MODULE).SetQueryStr(queryParams).SetLimit(1000000)
	if res, err := w.Find(); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, res)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_modules 添加测试组件
 * @apiVersion 0.1.0
 * @apiName Create
 * @apiGroup EvaluateModule
 *
 * @apiDescription 添加测试组件
 *
 * @apiUse authHeader
 *
 * @apiParam {string}		module_name				测试组件名称
 * @apiParam {string}		module_name_code		测试组件编号
 * @apiParam {string}		module_type				测试分类名称
 * @apiParam {string}		module_type_code		测试分类编号
 *
 * @apiParamExample {json}  请求参数示例:
 *     {
 *         "module_name" : "李清测试组件2",
 *         "module_name_code" : "011",
 *         "module_type" : "李清测试分类",
 *         "module_type_code" : "022"
 *     }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0
 * }
 */
func (this EvaluateModuleController) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if eType, err := new(mongo_model.EvaluateModule).Create(*req); err == nil {
		response.RenderSuccess(ctx, eType)
		return
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}

/**
 * apiType http
 * @api {put} /api/v1/evaluate_modules 更新测试组件
 * @apiVersion 0.1.0
 * @apiName Update
 * @apiGroup EvaluateModule
 *
 * @apiDescription 更新测试组件
 *
 * @apiUse authHeader
 *
 * @apiParam {string}		id						测试组件id
 * @apiParam {string}		module_name				测试组件名称
 * @apiParam {string}		module_name_code		测试组件编号
 * @apiParam {string}		module_type				测试分类名称
 * @apiParam {string}		module_type_code		测试分类编号
 *
 * @apiParamExample {json}  请求参数示例:
 *     {
 *         "id" : "sjj12030ajemnf",
 *         "module_name" : "李清测试组件2",
 *         "module_name_code" : "011",
 *         "module_type" : "李清测试分类",
 *         "module_type_code" : "022"
 *     }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0
 * }
 */
func (this EvaluateModuleController) Update(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")
	id := req.MustString("id")
	models := make([]mysql_model.KnowledgeTestCase, 0)
	mysql.GetSession().Where("module_id = ?", id).Find(&models)
	var reason string
	if len(models) > 0 {
		for _, model := range models {
			reason = fmt.Sprintf("已被测试用例%v引用,无法修改", model.Name)
		}
		response.RenderFailure(ctx, errors.New(reason))
		return
	}
	if eType, err := new(mongo_model.EvaluateModule).Update(id, *req); err == nil {
		response.RenderSuccess(ctx, eType)
		return
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/evaluate_modules 批量删除测试组件
 * @apiVersion 0.1.0
 * @apiName BulkDelete
 * @apiGroup EvaluateModule
 *
 * @apiDescription 批量删除测试组件
 *
 * @apiUse authHeader
 *
 * @apiParam {[]string}		ids		测试组件id
 *
 * @apiParamExample {json}  请求参数示例:
 *     {
 *         "ids" : ["5fe464b1f98f923e40e8dd5f"]
 *     }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *			"failure_number":0,
 *			"success_number":1
 *     }
 * }
 */
func (this EvaluateModuleController) BulkDelete(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")
	successNum := 0
	if _, has := req.TrySlice("ids"); has {
		ids := req.SliceString("ids")
		idsObject := []bson.ObjectId{}
		for _, id := range ids {
			params := qmap.QM{
				"e_module_type_id": id,
			}
			// 如果测试分类已经关联了测试用例，则不允许删除该测试分类
			models := make([]mysql_model.KnowledgeTestCase, 0)
			mysql.GetSession().Where("module_id = ?", id).Find(&models)
			//reason := make([]string, 0)
			if len(models) > 0 {
				//for _, model := range models {
				//	reason = append(reason, fmt.Sprintf("要删除的项目正在被测试用例\"%v\"使用,无法删除", model.Name))
				//}
				//response.RenderSuccess(ctx, reason)
				response.RenderFailure(ctx, errors.New("组件或测试分类正在被使用,无法删除"))
				return
			}
			if _, err := sys_service.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params).GetOne(); err == nil {
				response.RenderFailure(ctx, errors.New("没有查到这条测试分类"))
				return
			} else {
				idsObject = append(idsObject, bson.ObjectIdHex(id))
			}
		}
		// 删除测试记录
		match := bson.M{
			"_id": bson.M{"$in": idsObject},
		}
		if changeInfo, err := sys_service.NewMgoSession(common.MC_EVALUATE_MODULE).RemoveAll(match); err == nil {
			successNum = changeInfo.Removed
		} else {
			response.RenderFailure(ctx, err)
			return
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"number": successNum})
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_module/module_name_list 查询测试组件列表
 * @apiVersion 0.1.0
 * @apiName GetModuleNameList
 * @apiGroup EvaluateModule
 *
 * @apiDescription 查询测试组件列表
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "module_name": "无线电",
 *             "module_name_code": "090"
 *         },
 *         {
 *             "module_name": "通信SOC",
 *             "module_name_code": "020"
 *         }
 *     ]
 * }
 */
func (this EvaluateModuleController) GetModuleNameList(ctx *gin.Context) {
	list := []interface{}{}
	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_MODULE)

	if moduleNames, err := coll.Distinct(context.Background(), "module_name", bson.M{}); err == nil {
		w := orm_mongo.NewWidgetWithCollectionName(common.MC_EVALUATE_MODULE)
		for _, moduleName := range moduleNames {
			w.SetParams(qmap.QM{"e_module_name": moduleName})
			moduleItem := new(mongo_model.EvaluateModule)
			if err := w.One(moduleItem); err == nil {
				item := qmap.QM{
					"module_name":      moduleName,
					"module_name_code": moduleItem.ModuleNameCode,
				}
				list = append(list, item)
			} else {
				response.RenderFailure(ctx, err)
				return
			}
		}
	} else {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, list)
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_module/module_type_list 查询测试分类列表
 * @apiVersion 0.1.0
 * @apiName GetModuleTypeList
 * @apiGroup EvaluateModule
 *
 * @apiDescription 查询测试分类列表
 *
 * @apiUse authHeader
 *
 * @apiParam {string}		[module_name]		测试组件（排除某一组件）
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "module_type": "IVI蓝牙",
 *             "module_type_code": "450"
 *         },
 *         {
 *             "module_type": "蓝牙钥匙",
 *             "module_type_code": "460"
 *         }
 *     ]
 * }
 */
func (this EvaluateModuleController) GetModuleTypeList(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	list := []interface{}{}
	query := bson.M{}
	if moduleName, has := req.TryString("module_name"); has {
		query["module_name"] = bson.M{"$ne": moduleName}
	}

	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_MODULE)
	if moduleTypes, err := coll.Distinct(context.Background(), "module_type", bson.M{}); err == nil {
		w := orm_mongo.NewWidgetWithCollectionName(common.MC_EVALUATE_MODULE)
		for _, moduleType := range moduleTypes {
			w.SetParams(qmap.QM{"e_module_type": moduleType})
			moduleItem := new(mongo_model.EvaluateModule)
			if err := w.One(moduleItem); err == nil {
				item := qmap.QM{
					"module_type":      moduleType,
					"module_type_code": moduleItem.ModuleTypeCode,
				}
				list = append(list, item)
			} else {
				response.RenderFailure(ctx, err)
				return
			}
		}
	} else {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, list)
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_module/recommend_code 获取推荐测试组件/测试分类编号
 * @apiVersion 0.1.0
 * @apiName GetRecommendCode
 * @apiGroup EvaluateModule
 *
 * @apiDescription 获取推荐测试组件/测试分类编号
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *              "module_name_code": "110",
 *              "module_type_code": "540"
 *     }
 * }
 */
func (this EvaluateModuleController) GetRecommendCode(ctx *gin.Context) {
	result := qmap.QM{
		"module_name_code": new(mongo_model.EvaluateModule).GetRecommendModuleNameCode(),
		"module_type_code": new(mongo_model.EvaluateModule).GetRecommendModuleTypeCode(),
	}
	response.RenderSuccess(ctx, result)
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_module/rename_module_name 重命名组件名称
 * @apiVersion 0.1.0
 * @apiName RenameModuleName
 * @apiGroup EvaluateModule
 *
 * @apiDescription 重命名组件名称
 *
 * @apiUse authHeader
 *
 * @apiParam {string}		old_module_name		旧组件名称
 * @apiParam {string}		new_module_name		新组件名称
 *
 * @apiParamExample {json}  请求参数示例:
 *     {
 *              "old_module_name": "旧组件名称",
 *              "new_module_name": "新组件名称"
 *     }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *              "success_number": 12
 *     }
 * }
 */
func (this EvaluateModuleController) RenameModuleName(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	query := bson.M{
		"module_name": bson.M{"$eq": req.MustString("old_module_name")},
	}
	update := qmap.QM{
		"$set": qmap.QM{
			"module_name": req.MustString("new_module_name"),
		},
	}
	number := 0
	if change, err := mongo.NewMgoSession(common.MC_EVALUATE_MODULE).UpdateAll(query, update); err == nil {
		number = change.Updated
	} else {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, qmap.QM{"success_number": number})
}

// @auto_generated_api_end
/**
 * apiType http
 * @api {get} /api/v1/evaluate_module/project 查询项目测试组件树
 * @apiVersion 0.1.0
 * @apiName GetProjectModuleTree
 * @apiGroup EvaluateModule
 *
 * @apiDescription 查询项目下测试组件树
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "list": [
 *                 {
 *                     "_id": "5fd87545aee3d1849a56efa7",
 *                     "module_name": "无线电",
 *                     "module_name_code": "090",
 *                     "module_type": "GNSS",
 *                     "module_type_code": "520"
 *                 },
 *                 {
 *                     "_id": "5fd87545aee3d1849a56ef9d",
 *                     "module_name": "无线电",
 *                     "module_name_code": "090",
 *                     "module_type": "蜂窝网络",
 *                     "module_type_code": "510"
 *                 }
 *             ],
 *             "module_name": "无线电",
 *             "module_name_code": 90
 *         }
 *     ]
 * }
 */
func (this EvaluateModuleController) GetProjectModuleTree(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	// 查询项目下所有的 module_type_id
	projectId := req.MustString("project_id")
	moduleTypeIds := []string{}
	assets, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, qmap.QM{"e_project_id": projectId}).SetLimit(10000).Get()
	if err == nil {
		for _, asset := range *assets {
			if asset["module_type_id"].([]interface{}) != nil {
				for _, id := range asset["module_type_id"].([]interface{}) {
					if !custom_util.IndexOfSlice(id.(string), moduleTypeIds) {
						moduleTypeIds = append(moduleTypeIds, id.(string))
					}

				}
			}
		}
	}

	list := []interface{}{}
	moduleNames := []string{}
	if err := mongo.NewMgoSession(common.MC_EVALUATE_MODULE).Session.Find(nil).Distinct("module_name", &moduleNames); err == nil {
		mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_MODULE).SetLimit(100000)
		for _, moduleName := range moduleNames {
			mgoSession.AddCondition(qmap.QM{"e_module_name": moduleName})
			if res, err := mgoSession.Get(); err == nil {
				if len(*res) > 0 {
					for _, item := range *res {
						if custom_util.IndexOfSlice(item["_id"].(bson.ObjectId).Hex(), moduleTypeIds) {
							var tmp qmap.QM = (*res)[0]
							item := qmap.QM{
								"module_name":      moduleName,
								"module_name_code": tmp.Int("module_name_code"),
								"list":             res,
							}
							list = append(list, item)

						}
					}
				}
			} else {
				response.RenderFailure(ctx, err)
				return
			}
		}
	} else {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, list)

}

func GetModuleMap() qmap.QM {
	moduleList, err := mongo.NewMgoSession(common.MC_EVALUATE_MODULE).SetLimit(50000).Get()
	if err != nil {
		return nil
	}
	moduleMap := map[string]interface{}{}
	for _, item := range *moduleList {
		id := item["_id"]
		itType := reflect.TypeOf(id)
		switch itType.Name() {
		case "bson.ObjectId":
			moduleMap[item["_id"].(bson.ObjectId).Hex()] = item["module_name"]
		default:
			continue
		}
	}
	return moduleMap
}
