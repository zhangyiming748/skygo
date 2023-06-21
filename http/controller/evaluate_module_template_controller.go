package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mongo_model"
)

type EvaluateModuleTemplateController struct{}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_module/all 查询测试组件
 * @apiVersion 0.1.0
 * @apiName GetAll
 * @apiGroup EvaluateModuleTemplate
 *
 * @apiDescription 查询测试组件
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "list": [
 *             {
 *                 "_id": "5f65c6675d502677ed23e5fa",
 *                 "item_name": "测试项3",
 *                 "level": 1,
 *                 "module_name": "测试组件",
 *                 "module_type": "测试分类",
 *                 "objective": "测试目标"
 *             },
 *             ...
 *         ],
 *         "pagination": {
 *             "count": 4,
 *             "current_page": 1,
 *             "per_page": 20,
 *             "total": 4,
 *             "total_pages": 1
 *         }
 *     }
 * }
 */
func (this EvaluateModuleTemplateController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery

	mgoSession := mongo.NewMgoSession(common.MC_EvaluateModuleTemplate).AddUrlQueryCondition(queryParams)

	if res, err := mgoSession.GetPage(); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, res)
	}
}

func (this EvaluateModuleTemplateController) GetOne(ctx *gin.Context) {
	id := ctx.Param("id")

	params := qmap.QM{
		"e__id": bson.ObjectIdHex(id),
	}
	ormSession := mongo.NewMgoSessionWithCond(common.MC_EvaluateModuleTemplate, params)
	if res, err := ormSession.GetOne(); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, res)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_templates 创建测试模板
 * @apiVersion 0.1.0
 * @apiName Create
 * @apiGroup EvaluateModuleTemplate
 *
 * @apiDescription 创建测试模板
 *
 * curl http://10.16.133.118:3001/api/v1/evaluate_templates
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           module_name                      测试组件
 *
 * @apiParamExample {json}  请求参数示例:
 *	{
 *		"module_name":"测试组件",
 *		"module_type":"测试分类",
 *		"item_name":"测试项",
 *		"objective":"测试目标",
 *		"level":1
 *	}
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": "5f65c5ac5d502677ed23e5f7",
 *         "item_name": "测试项",
 *         "level": 1,
 *         "module_name": "测试组件",
 *         "module_type": "测试分类",
 *         "objective": "测试目标"
 *     }
 * }
 */
func (this EvaluateModuleTemplateController) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if testItem, err := new(mongo_model.EvaluateModuleTemplate).Create(*req); err == nil {
		response.RenderSuccess(ctx, testItem)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this EvaluateModuleTemplateController) BulkDelete(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if testItem, err := new(mongo_model.EvaluateModuleTemplate).BulkDelete(*req); err == nil {
		response.RenderSuccess(ctx, testItem)
	} else {
		response.RenderFailure(ctx, err)
	}
}
