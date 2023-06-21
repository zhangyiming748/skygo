package controller

import (
	"errors"

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

type ProjectFactoryController struct{}

//@auto_generated_api_begin
/**
 * apiType http
 * @api {get} /api/v1/project_factories 车厂列表
 * @apiVersion 1.0.0
 * @apiName GetAll
 * @apiGroup ProjectFactory
 *
 * @apiDescription 查询车机列表接口
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/project_factories
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "list": [
 *             {
 *                 "_id": "5e61f95024b64748d37d8cc6",
 *                 "name": "车厂名称",
 *                 "create_time": 1583479120091
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
func (this ProjectFactoryController) GetAll(ctx *gin.Context) {
	mgoSession := mongo.NewMgoSession(common.MC_FACTORY).AddUrlQueryCondition(ctx.Request.URL.RawQuery)
	if res, err := mgoSession.GetPage(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/project_factories/:id 查询某一个车厂信息
 * @apiVersion 1.0.0
 * @apiName GetOne
 * @apiGroup ProjectFactory
 *
 * @apiDescription 查询某一个车厂信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}       id        车厂id
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/project_factories/:id
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *                 "_id": "5e61f95024b64748d37d8cc6",
 *                 "name": "车厂名称",
 *                 "create_time": 1583479120091
 *     }
 * }
 */
func (this ProjectFactoryController) GetOne(ctx *gin.Context) {
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(request.ParamString(ctx, "id")),
	}
	ormSession := mongo.NewMgoSessionWithCond(common.MC_PROJECT, params)
	if res, err := ormSession.GetOne(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/project_factories 创建车厂
 * @apiVersion 1.0.0
 * @apiName Create
 * @apiGroup ProjectFactory
 *
 * @apiDescription 创建新车厂
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           name    				车厂名称
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST http://localhost/api/v1/project_factories
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *                 "name": "车厂名称",
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *                 "_id": "5e61f95024b64748d37d8cc6",
 *                 "name": "车厂名称",
 *                 "create_time": 1583479120091
 *     }
 * }
 */
func (this ProjectFactoryController) Create(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	if factory, err := new(mongo_model.Factory).Create(params); err == nil {
		if ff, err := custom_util.StructToMap(*factory); err == nil {
			response.RenderSuccess(ctx, ff)
		} else {
			response.RenderFailure(ctx, err)
		}
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {put} /api/v1/project_factories/:id 更新车厂
 * @apiVersion 1.0.0
 * @apiName Update
 * @apiGroup ProjectFactory
 *
 * @apiDescription 更新车厂接口
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           id                      车厂id
 * @apiParam {string}           name    				车厂名称
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X PUT http://localhost/api/v1/project_factories/:id
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *                 "name": "车厂名称",
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *                 "_id": "5e61f95024b64748d37d8cc6",
 *                 "name": "车厂名称",
 *                 "create_time": 1583479120091
 *     }
 * }
 */
func (this ProjectFactoryController) Update(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
		"id":           request.ParamString(ctx, "id"),
	}

	id := request.String(ctx, "id")
	updateCols := map[string]string{
		"name": "string",
	}

	rawInfo := custom_util.CopyMapColumns(*params, updateCols)
	if factory, err := new(mongo_model.Factory).Update(id, rawInfo); err == nil {
		if ff, err := custom_util.StructToMap(*factory); err == nil {
			response.RenderSuccess(ctx, ff)
		} else {
			response.RenderFailure(ctx, err)
		}
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/project_factories 批量删除车厂
 * @apiVersion 1.0.0
 * @apiName BulkDelete
 * @apiGroup ProjectFactory
 *
 * @apiDescription 批量删除车厂
 *
 * @apiUse authHeader
 *
 * @apiParam {[]string}   ids  车厂id
 *
 * @apiSuccessExample {json} 请求成功示例:
 *       {
 *            "code": 0,
 *			  "data":{
 *				"number":1
 *			}
 *       }
 */
func (this ProjectFactoryController) BulkDelete(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	effectNum := 0
	if rawIds, has := params.TrySlice("ids"); has {
		ids := []bson.ObjectId{}
		for _, id := range rawIds {
			// 车厂已经关联了项目，则无法删除
			_, err := mongo.NewMgoSession(common.MC_PROJECT).AddCondition(qmap.QM{"e_company": id.(string)}).GetOne()
			if err == nil {
				response.RenderFailure(ctx, errors.New("车厂已经关联测试项目，不可删除"))
			}
			ids = append(ids, bson.ObjectIdHex(id.(string)))
		}
		if len(ids) > 0 {
			match := bson.M{
				"_id": bson.M{"$in": ids},
			}
			if changeInfo, err := mongo.NewMgoSession(common.MC_FACTORY).RemoveAll(match); err == nil {
				effectNum = changeInfo.Removed
			} else {
				response.RenderFailure(ctx, err)
			}
		}
	}
	response.RenderSuccess(ctx, &qmap.QM{"number": effectNum})
}
