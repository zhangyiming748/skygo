package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/http/transformer"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/log"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/common_lib/session"
	"skygo_detection/mysql_model"
)

type AssetVehicleController struct{}

/**
 * apiType http
 * @api {get} /api/v1/asset_vehicles 车型列表查询
 * @apiVersion 0.1.0
 * @apiName GetAll
 * @apiGroup AssetVehicle
 *
 * @apiDescription 车型列表查询
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "list": [
 *             {
 *                 "brand": "b1111112",
 *                 "code": "c111",
 *                 "create_time": 1628674339,
 *                 "create_user_id": 0,
 *                 "detail": "d111",
 *                 "id": 4,
 *                 "serial_number": "",
 *                 "update_time": 0
 *             }
 *         ],
 *         "pagination": {
 *             "current_page": 1,
 *             "per_page": 20,
 *             "total": 5,
 *             "total_pages": 1
 *         }
 *     },
 *     "msg": ""
 * }
 */
func (this AssetVehicleController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	s := mysql.GetSession()

	// 查询组键
	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	widget.SetTransformer(&transformer.AssetVehicleTransformer{})
	all := widget.PaginatorFind(s, &[]mysql_model.AssetVehicle{})
	response.RenderSuccess(ctx, all)
}

/**
 * apiType http
 * @api {post} /api/v1/asset_vehicles 创建车型记录
 * @apiVersion 0.1.0
 * @apiName Create
 * @apiGroup AssetVehicle
 *
 * @apiDescription 车型记录详情
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 4,
 *         "brand": "b1111112",
 *         "code": "c111",
 *         "detail": "d111",
 *         "create_user_id": 0,
 *         "create_user_name": 0,
 *         "update_time": 0,
 *         "create_time": 1628674339
 *     },
 *     "msg": ""
 * }
 */
func (this AssetVehicleController) Create(ctx *gin.Context) {
	// 表单
	form := &mysql_model.AssetVehicleCreateForm{}
	form.Brand = request.MustString(ctx, "brand")
	form.Code = request.MustString(ctx, "code")
	form.Detail = request.MustString(ctx, "detail")

	uid := session.GetUserId(http_ctx.GetHttpCtx(ctx))
	model, err := mysql_model.AssetVehicleCreateFromForm(form, int(uid))

	if err == nil {
		response.RenderSuccess(ctx, model)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/asset_vehicles/:id 车型记录详情
 * @apiVersion 0.1.0
 * @apiName GetOne
 * @apiGroup AssetVehicle
 *
 * @apiDescription 车型记录详情
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 4,
 *         "brand": "b1111112",
 *         "code": "c111",
 *         "detail": "d111",
 *         "create_user_id": 0,
 *         "create_user_name": 0,
 *         "update_time": 0,
 *         "create_time": 1628674339
 *     },
 *     "msg": ""
 * }
 */
func (this AssetVehicleController) GetOne(ctx *gin.Context) {
	id := request.ParamString(ctx, "id")
	s := mysql.GetSession()
	s.Where("id=?", id)

	w := orm.PWidget{}
	result, err := w.One(s, &mysql_model.AssetVehicle{})

	if err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {put} /api/v1/asset_vehicles/:id 车型记录更新
 * @apiVersion 0.1.0
 * @apiName Update
 * @apiGroup AssetVehicle
 *
 * @apiDescription 车型记录更新
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 4,
 *         "serial_number": "",
 *         "brand": "b1111112",
 *         "code": "c111",
 *         "detail": "d111",
 *         "create_user_id": 0,
 *         "update_time": 0,
 *         "create_time": 1628674339
 *     },
 *     "msg": ""
 * }
 */
func (this AssetVehicleController) Update(ctx *gin.Context) {
	data := request.GetRequestBody(ctx)
	id := request.ParamInt(ctx, "id")

	if model, err := mysql_model.AssetVehicleUpdateById(id, *data); err == nil {
		response.RenderSuccess(ctx, model)
		return
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/asset_test_pieces 车型资产批量删除
 * @apiVersion 0.1.0
 * @apiName BulkDelete
 * @apiGroup AssetVehicle
 *
 * @apiDescription 车型资产批量删除
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "number": 1
 *     },
 *     "msg": ""
 * }
 */
func (this AssetVehicleController) BulkDelete(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)

	successNum := 0
	if _, has := req.TrySlice("ids"); has {
		ids := req.SliceInt("ids")

		s := mysql.GetSession()
		for _, id := range ids {
			// todo 检查是否可以删除

			_, err := s.ID(id).Delete(&mysql_model.AssetVehicle{})
			if err != nil {
				log.GetHttpLogLogger().Error(fmt.Sprintf("%v", err))
				break
			} else {
				successNum++
			}
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"number": successNum})
}

/**
 * apiType http
 * @api {get} /api/v1/asset_vehicle/select_list 车型品牌代号列表
 * @apiVersion 0.1.0
 * @apiName SelectList
 * @apiGroup AssetVehicle
 *
 * @apiDescription 车型品牌代号列表
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 *  {
 *      "code": 0,
 *      "data": [
 *          {
 *              "brand": "b1111112",
 *              "codes": [
 *                  {
 *                      "id": 4,
 *                      "code": "c111"
 *                  },
 *                  {
 *                      "id": 5,
 *                      "code": "c111"
 *                  }
 *              ]
 *          }
 *      ],
 *      "msg": ""
 *  }
 */
func (this AssetVehicleController) SelectList(ctx *gin.Context) {
	data := mysql_model.AssetVehicleSelectList()
	response.RenderSuccess(ctx, data)
	return
}
