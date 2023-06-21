package controller

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"
	"xorm.io/builder"

	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mysql_model"
)

type MService struct{}

/**
 * apiType http
 * @api {get} /api/v1/mservices 查询微服务列表
 * @apiVersion 1.0.0
 * @apiName GetAll
 * @apiGroup MService
 *
 * @apiDescription 分页查询微服务列表
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/mservices
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "create_time": 1572350048,
 *             "id": 6,
 *             "service": "skygo.admin",
 *             "service_name": "运营平台网关服务"
 *         }
 *     ]
 * }
 */
func (this MService) GetAll(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	s := mysql.GetSession().OrderBy("id desc").Limit(1000)

	widget := orm.PWidget{}
	widget.SetQueryStr(ctx.Request.URL.RawQuery)
	widget.SetTransformerFunc(userTransform)
	all, _ := widget.All(s, &[]mysql_model.MicroService{})
	response.RenderSuccess(ctx, all)
}

/**
 * apiType http
 * @api {get} /api/v1/mservices/:id 查询微服务信息
 * @apiVersion 1.0.0
 * @apiName GetOne
 * @apiGroup MService
 *
 * @apiDescription 根据id查询某一微服务信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   id  用户id
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/mservices/{id}
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "create_time": 1572349900,
 *         "id": 1,
 *         "service": "skygo.auth",
 *         "service_name": "授权服务"
 *     }
 * }
 */
func (this MService) GetOne(ctx *gin.Context) {
	id := ctx.Param("id")

	model := mysql_model.MicroService{}

	if has, _ := mysql.GetSession().ID(id).Get(&model); has {
		response.RenderSuccess(ctx, model)
	} else {
		err := errors.New("Item not found")
		response.RenderFailure(ctx, err)
	}

}

/**
 * apiType http
 * @api {post} /api/v1/mservices 添加微服务
 * @apiVersion 1.0.0
 * @apiName Create
 * @apiGroup MService
 *
 * @apiDescription 添加微服务
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           service                服务
 * @apiParam {string}           service_name           服务名称
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST -d service=service&service_name=服务名称 http://localhost/api/v1/mservices
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "service":"service",
 *          "service_name":"服务名称"
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "create_time": 1572349900,
 *         "id": 1,
 *         "service": "skygo.auth",
 *         "service_name": "授权服务"
 *     }
 * }
 */
func (this MService) Create(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)

	service := mysql_model.MicroService{
		Service:     req.MustString("service"),
		ServiceName: req.MustString("service_name"),
	}

	if _, err := mysql.GetSession().InsertOne(&service); err == nil {
		response.RenderSuccess(ctx, custom_util.StructToMap2(service))
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {put} /api/v1/mservices/:id 更新微服务
 * @apiVersion 1.0.0
 * @apiName Update
 * @apiGroup MService
 *
 * @apiDescription 根据id,更新微服务
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           [service]                  	微服务
 * @apiParam {string}           [service_name]              微服务名称
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X PUT -d service=skygo.test&service_name=测试服务 http://localhost/api/v1/mservices/1
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "service":"skygo.test",
 *          "service_name":"测试服务"
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "create_time": 1572349900,
 *         "id": 1,
 *         "service": "skygo.auth",
 *         "service_name": "授权服务"
 *     }
 * }
 */
func (this MService) Update(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))

	id := ctx.Param("id")
	columns := map[string]string{
		"service":      "string",
		"service_name": "string",
	}
	updateData := custom_util.CopyMapColumns(*req, columns)
	if _, err := mysql.GetSession().ID(id).Table(mysql_model.MicroService{}).Update(updateData); err == nil {
		model := mysql_model.MicroService{}
		if has, _ := mysql.GetSession().ID(id).Get(&model); has {
			response.RenderSuccess(ctx, model)
			return
		} else {
			err := errors.New("Item not found")
			response.RenderFailure(ctx, err)
			return
		}
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/mservices 批量删除微服务
 * @apiVersion 1.0.0
 * @apiName BulkDelete
 * @apiGroup MService
 *
 * @apiDescription 批量删除微服务
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   ids  用户id,多个微服务id之间用"\\|"连接(如:"1\\|2\\|3")
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X DELETE http://localhost/api/v1/mservices?ids=1|2|3
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":{
 *				"number":2
 *			}
 *      }
 */
func (this MService) BulkDelete(ctx *gin.Context) {

	idStr := strings.Split(request.MustString(ctx, "ids"), "|")

	ids := []int{}
	for _, val := range idStr {
		if id, convErr := strconv.Atoi(val); convErr == nil {
			ids = append(ids, id)
		} else {
			response.RenderFailure(ctx, errors.New("convErr")) // todo err
			return
		}
	}
	if effectNum, err := mysql.GetSession().Where(builder.In("id", ids)).Delete(new(mysql_model.MicroService)); err == nil {
		response.RenderSuccess(ctx, qmap.QM{"number": effectNum})
		return
	} else {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, gin.H{})
}
