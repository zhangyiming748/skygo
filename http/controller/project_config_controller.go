package controller

import (
	"errors"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mongo_model"
)

type ProjectConfigController struct{}

/**
 * apiType http
 * @api {get} /api/v1/project_config/all_vul_type 查询所有漏洞类型
 * @apiVersion 1.0.0
 * @apiName GetAllVulType
 * @apiGroup ProjectConfig
 *
 * @apiDescription 查询所有漏洞类型
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json}  请求成功示例:
 * {
 *         "code": 0,
 *         "data": [
 *                 {
 *                         "_id": 2,
 *                         "create_time": 1609992391,
 *                         "name": "test2",
 *                         "status": 0,
 *                         "update_time": 1609992391
 *                 }
 *         ]
 * }
 */
func (this ProjectConfigController) GetAllVulType(ctx *gin.Context) {
	// params := &qmap.QM{
	//	"query_params": ctx.Request.URL.RawQuery,
	// }
	// *params = params.Merge(*request.GetRequestQueryParams(ctx))
	// *params = params.Merge(*request.GetRequestBody(ctx))
	if res, err := new(mongo_model.EvaluateVulType).GetAll(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/project_config/upsert_vul_type 更新/新增漏洞类型
 * @apiVersion 1.0.0
 * @apiName UpsertVulType
 * @apiGroup ProjectConfig
 *
 * @apiDescription 更新/新增漏洞类型
 *
 * @apiUse authHeader
 *
 * @apiParam {int}           	[id]					漏洞类型id
 * @apiParam {string}           name					漏洞类型名称
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *         "id":1,
 *         "name":"漏洞类型名称"
 *     }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0
 * }
 */
func (this ProjectConfigController) UpsertVulType(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	_, err := new(mongo_model.EvaluateVulType).Upsert(*params)
	if err == nil {
		response.RenderSuccess(ctx, nil)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/project_config/bulk_delete_vul_type 批量删除漏洞类型
 * @apiVersion 1.0.0
 * @apiName BulkDeleteVulType
 * @apiGroup ProjectConfig
 *
 * @apiDescription 批量删除漏洞类型
 *
 * @apiUse authHeader
 *
 * @apiParam {[]int]}   ids  	漏洞类型id
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":{
 *				"number":2
 *			}
 *      }
 */
func (this ProjectConfigController) BulkDeleteVulType(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	ids := params.SliceInt("ids")
	successNum := 0
	successIds := []int{}
	for _, id := range ids {
		params := qmap.QM{
			"e_risk_type": id,
		}
		if _, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VULNERABILITY, params).GetOne(); err == nil {
			// 如果漏洞类型已经关联了漏洞，则不允许删除该漏洞类型
			response.RenderFailure(ctx, errors.New("漏洞类型已经关联了漏洞，无法删除"))
		} else {
			successIds = append(successIds, id)
		}
	}
	for _, id := range successIds {
		match := bson.M{
			"_id": bson.M{"$eq": id},
		}
		if _, err := mongo.NewMgoSession(common.MC_EVALUATE_VUL_TYPE).RemoveAll(match); err == nil {
			successNum++
		} else {
			response.RenderFailure(ctx, err)
		}
	}
	response.RenderSuccess(ctx, &qmap.QM{"number": successNum})
}
