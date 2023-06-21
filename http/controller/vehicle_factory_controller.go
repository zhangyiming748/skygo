package controller

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"
	"xorm.io/builder"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/common_lib/session"
	"skygo_detection/mysql_model"
)

type VehicleFactory struct{}

/**
 * apiType http
 * @api {get} /api/v1/vehicle_factories 查询车厂列表
 * @apiVersion 1.0.0
 * @apiName GetAll
 * @apiGroup VehicleFactory
 *
 * @apiDescription 分页查询车厂列表
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/vehicle_factories
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "list": [
 *             {
 *                 "account_type": "user",
 *                 "channel_id": "T56205",
 *                 "create_time": 1571367647,
 *                 "email": "",
 *                 "head_pic": "",
 *                 "id": 45,
 *                 "mobile": 0,
 *                 "nickname": "",
 *                 "realname": "test",
 *                 "sex": 0,
 *                 "status": 2,
 *                 "username": "test"
 *             }
 *         ],
 *         "meta": {
 *             "count": 4,
 *             "current_page": 1,
 *             "per_page": 20,
 *             "total": 9,
 *             "total_pages": 1
 *         }
 *     }
 * }
 */
func (this VehicleFactory) GetAll(ctx *gin.Context) {
	s := mysql.GetSession()

	// 超级管理员对用户的查询不做渠道限制
	// if session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx)) != common.SUPER_ADMINISTRATE_ROLE_ID {
	// 	s.Where("channel_id = ?", session.GetQueryChannelId(http_ctx.GetHttpCtx(ctx)))
	// }

	widget := orm.PWidget{}
	widget.SetQueryStr(ctx.Request.URL.RawQuery)
	all := widget.PaginatorFind(s, &[]mysql_model.SysVehicleFactory{})
	response.RenderSuccess(ctx, all)
}

/**
 * apiType http
 * @api {get} /api/v1/vehicle_factories/:id 查询车厂信息
 * @apiVersion 1.0.0
 * @apiName GetOne
 * @apiGroup VehicleFactory
 *
 * @apiDescription 根据id查询某一车厂信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   id  车厂id
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/vehicle_factories/1
 *
 * @apiSuccessExample {json} 请求成功示例:
 *            {
 *                "code": 0,
 *                "data": {
 *                      "id": 1,
 *                      "name": "test",
 *                      "type": 0,
 *                      "channel_id": "C02001",
 *                      "status": 1,
 *                      "update_time": 0,
 *                      "create_time": 0
 *                }
 *             }
 */
func (this VehicleFactory) GetOne(ctx *gin.Context) {
	s := mysql.GetSession().Table(new(mysql_model.SysVehicleFactory)).Where("id = ?", ctx.Param("id"))

	// 超级管理员对用户的查询不做渠道限制
	if session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx)) != common.SUPER_ADMINISTRATE_ROLE_ID {
		s.Where("channel_id = ?", session.GetQueryChannelId(http_ctx.GetHttpCtx(ctx)))
	}

	session := mysql.GetSession()
	widget := orm.PWidget{}
	widget.SetQueryStr(ctx.Request.URL.RawQuery)
	if one, err := widget.Get(session); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, one)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/vehicle_factories 添加车厂
 * @apiVersion 1.0.0
* @apiName Create
* @apiGroup VehicleFactory
*
* @apiDescription 添加一个车厂
*
* @apiUse authHeader
*
* @apiParam {string}           name                      车厂名称
* @apiParam {int=0,1}          type                      车厂类型(0:车厂，1:供应商)
* @apiParam {status=0,1}       [status=1]                状态(0:禁用， 1:启用)
*
* @apiExample {curl} 请求示例:
* curl -i -X POST -d name=车厂名称&type=1 http://localhost/api/v1/vehicle_factories
*
* @apiParamExample {json}  请求参数示例:
*              {
*                      "name": "test",
*                      "type": 0,
*                      "status": 1
*              }
*
* @apiSuccessExample {json} 请求成功示例:
*                {
*                    "code": 0,
*                	 "data": {
*                      "id": 1,
*                      "name": "test",
*                      "type": 0,
*                      "channel_id": "C02001",
*                      "status": 1,
*                      "update_time": 0,
*                      "create_time": 0
*                	}
*                }
*/
func (this VehicleFactory) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	// (*req)["id"] = ctx.Param("id")
	// (*req)["query_params"] = ctx.Request.URL.RawQuery

	channelId := session.GetQueryChannelId(http_ctx.GetHttpCtx(ctx))
	if channelId != "" && session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx)) != common.SUPER_ADMINISTRATE_ROLE_ID {
		err := errors.New("This is not permitted")
		response.RenderFailure(ctx, err)
		return
	}
	newItem := mysql_model.SysVehicleFactory{
		Name:   req.String("name"),
		Type:   req.MustInt("type"),
		Status: req.Int("status"),
	}
	newItem.ChannelId = newItem.GenerateChannelId(newItem.Type)
	ormSession := mysql.GetSession()
	if _, err := ormSession.InsertOne(newItem); err == nil {
		if has, _ := ormSession.Get(&newItem); has {
			response.RenderSuccess(ctx, newItem)
			return
		} else {
			response.RenderFailure(ctx, errors.New("Create failure"))
			return
		}
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}

/**
 * apiType http
 * @api {put} /api/v1/vehicle_factories/:id 更新车厂
 * @apiVersion 1.0.0
 * @apiName Update
 * @apiGroup VehicleFactory
 *
 * @apiDescription 根据车厂id,更新车厂信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   		id  						车厂id
 * @apiParam {string}           [name]                      车厂名称
 * @apiParam {int=0,1}          [type]                      车厂类型(0:车厂，1:供应商)
 * @apiParam {status=0,1}       [status]                    状态(0:禁用， 1:启用)
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X PUT -d name=车厂名称&type=1&status=1 http://localhost/api/vehicle_factories/11
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *                      "id": 1,
 *                      "name": "test",
 *                      "type": 0,
 *                      "status": 1
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *                {
 *                      "code": 0,
 *                      "data": {
 *                          "id": 1,
 *                          "name": "test",
 *                          "type": 0,
 *                          "channel_id": "C02001",
 *                          "status": 1,
 *                          "update_time": 0,
 *                          "create_time": 0
 *                      }
 *                }
 *
 */
func (this VehicleFactory) Update(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")
	// (*req)["query_params"] = ctx.Request.URL.RawQuery

	id := req.MustInt("id")
	columns := map[string]string{
		"name":   "string",
		"type":   "int",
		"status": "int",
	}
	updateData := custom_util.CopyMapColumns(*req, columns)

	s := mysql.GetSession().Where("id = ?", id).And("channel_id = ?", session.GetQueryChannelId(http_ctx.GetHttpCtx(ctx)))

	if _, err := s.ID(id).Table(new(mysql_model.SysVehicleFactory)).Update(updateData); err == nil {
		res := new(mysql_model.SysVehicleFactory)
		if has, res := mysql.GetSession().ID(id).Get(res); has {
			response.RenderSuccess(ctx, res)
			return
		} else {
			response.RenderFailure(ctx, errors.New("Item not found"))
			return
		}
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/vehicle_factories 批量删除车厂
 * @apiVersion 1.0.0
 * @apiName BulkDelete
 * @apiGroup VehicleFactory
 *
 * @apiDescription 批量删除车厂
 *
 * @apiUse authHeader
 *
 * @apiParam {string}       ids         名单id,id之间用"\\|"连接,如"1\\|2\\|3"
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X DELETE http://localhost/api/v1/vehicle_factories?ids=1|2|3
 *
 * @apiSuccessExample {json} 请求成功示例:
 *       {
 *           "code": 0,
 *			 "data":{
 *				"number":2
 *		    }
 *       }
 */
func (this VehicleFactory) BulkDelete(ctx *gin.Context) {
	idStr := strings.Split(request.MustQueryString(ctx, "ids"), "|")
	ids := []int{}
	for _, val := range idStr {
		if id, convErr := strconv.Atoi(val); convErr == nil {
			ids = append(ids, id)
		} else {
			response.RenderFailure(ctx, convErr)
			return
		}
	}

	s := mysql.GetSession().Where("channel_id = ?", session.GetQueryChannelId(http_ctx.GetHttpCtx(ctx)))
	s.And(builder.In("id", ids))
	if effectNum, err := s.Delete(new(mysql_model.SysVehicleFactory)); err == nil {
		response.RenderSuccess(ctx, &qmap.QM{"number": effectNum})
		return
	} else {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, gin.H{})
}
