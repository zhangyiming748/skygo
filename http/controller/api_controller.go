package controller

import (
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"
	"xorm.io/builder"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mysql_model"
)

type Api struct{}

/**
 * apiType http
 * @api {get} /api/v1/apis/:service 查询某个微服务接口列表
 * @apiVersion 1.0.0
 * @apiName GetAll
 * @apiGroup Api
 *
 * @apiDescription 查询某个微服务接口列表
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/apis/skygo.auth
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "list": [
 *             {
 *                 "api_type": "rpc",
 *                 "description": "",
 *                 "id": 1,
 *                 "method": "POST",
 *                 "resource": "",
 *                 "url": ""
 *             }
 *         ],
 *         "pagination": {
 *             "count": 1,
 *             "current_page": 1,
 *             "per_page": 20,
 *             "total": 1,
 *             "total_pages": 1
 *         }
 *     }
 * }
 */
func (this Api) GetAll(ctx *gin.Context) {
	serviceName := request.ParamString(ctx, "service")
	switch serviceName {
	case common.PM_SERVICE:
		if res, err := new(mysql_model.SysApiBusiness).GetApiTreeList(); err == nil {
			response.RenderSuccess(ctx, res)
		} else {
			panic(err)
		}
	case common.ADMIN_SERVICE:
		if res, err := new(mysql_model.SysApi).GetApiTreeList(); err == nil {
			response.RenderSuccess(ctx, res)
		} else {
			panic(err)
		}
	default:
		panic("Unknown service name")
	}
}

/**
 * apiType http
 * @api {get} /api/v1/mservices/:service/:id 查询某个微服务的某一接口信息
 * @apiVersion 1.0.0
 * @apiName GetOne
 * @apiGroup Api
 *
 * @apiDescription 根据id查询某一微服务的某一接口信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   service  	服务名称
 * @apiParam {string}   id  		接口id
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/mservices/skygo.auth/1
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "api_type": "rpc",
 *         "description": "",
 *         "id": 1,
 *         "method": "POST",
 *         "resource": "",
 *         "url": ""
 *     }
 * }
 */
func (this Api) GetOne(ctx *gin.Context) {
	id := request.ParamInt(ctx, "id")

	serviceName := request.ParamString(ctx, "service")
	switch serviceName {
	case common.PM_SERVICE:
		model := mysql_model.SysApiBusiness{}
		if has, _ := mysql.FindById(id, &model); has {
			response.RenderSuccess(ctx, model)
		} else {
			panic("Item not found")
		}
	case common.ADMIN_SERVICE:
		model := mysql_model.SysApi{}
		if has, _ := mysql.FindById(id, &model); has {
			response.RenderSuccess(ctx, model)
		} else {
			panic("Item not found")
		}
	default:
		panic("Unknown service name")
	}
}

/**
 * apiType http
 * @api {post} /api/v1/mservices 向某个微服务添加接口
 * @apiVersion 1.0.0
 * @apiName Create
 * @apiGroup Api
 *
 * @apiDescription 向某个微服务添加接口
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           api_type    		接口类型(http/rpc)
 * @apiParam {string}           method           	接口请求方式(GET/POST/DELETE/PUT)
 * @apiParam {string}           url           		统一资源定位符
 * @apiParam {string}           resource           	资源名称
 * @apiParam {string}           description         接口描述
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST -d api_type=http&method=POST&url=/api/v1/test&resource=test&description=接口描述 http://localhost/api/v1/mservices
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "api_type":"http",
 *          "method":"POST",
 *          "url":"/api/v1/test",
 *          "resource":"test",
 *          "description":"接口描述"
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "api_type": "rpc",
 *         "description": "",
 *         "id": 1,
 *         "method": "POST",
 *         "resource": "",
 *         "url": ""
 *     }
 * }
 */
func (this Api) Create(ctx *gin.Context) {
	serviceName := request.ParamString(ctx, "service")
	switch serviceName {
	case common.PM_SERVICE:
		newItem := mysql_model.SysApiBusiness{
			ApiType:     request.MustString(ctx, "api_type"),
			Method:      request.MustString(ctx, "method"),
			Url:         request.MustString(ctx, "url"),
			Resource:    request.MustString(ctx, "resource"),
			Description: request.MustString(ctx, "description"),
		}

		s := mysql.GetSession()
		if _, err := s.InsertOne(&newItem); err == nil {
			response.RenderSuccess(ctx, custom_util.StructToMap2(newItem))
		} else {
			panic(err)
		}
	case common.ADMIN_SERVICE:
		newItem := mysql_model.SysApi{
			ApiType:     request.MustString(ctx, "api_type"),
			Method:      request.MustString(ctx, "method"),
			Url:         request.MustString(ctx, "url"),
			Resource:    request.MustString(ctx, "resource"),
			Description: request.MustString(ctx, "description"),
		}

		s := mysql.GetSession()
		if _, err := s.InsertOne(&newItem); err == nil {
			response.RenderSuccess(ctx, custom_util.StructToMap2(newItem))
		} else {
			panic(err)
		}
	default:
		panic("Unknown service name")
	}
}

/**
 * apiType http
 * @api {put} /api/v1/mservices/:id 更新微服务
 * @apiVersion 1.0.0
 * @apiName Update
 * @apiGroup Api
 *
 * @apiDescription 根据id,更新微服务
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           api_type    		接口类型(http/rpc)
 * @apiParam {string}           method           	接口请求方式(GET/POST/DELETE/PUT)
 * @apiParam {string}           url           		统一资源定位符
 * @apiParam {string}           resource           	资源名称
 * @apiParam {string}           description         接口描述
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X PUT -d api_type=http&method=POST&url=/api/v1/test&resource=test&description=接口描述 http://localhost/api/v1/Apis/skygo.auth/1
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "api_type":"http",
 *          "method":"POST",
 *          "url":"/api/v1/test",
 *          "resource":"test",
 *          "description":"接口描述"
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "api_type": "rpc",
 *         "description": "",
 *         "id": 1,
 *         "method": "POST",
 *         "resource": "",
 *         "url": ""
 *     }
 * }
 */
func (this Api) Update(ctx *gin.Context) {
	post := request.GetRequestBody(ctx)
	(*post)["id"] = request.ParamInt(ctx, "id")

	serviceName := request.ParamString(ctx, "service")
	switch serviceName {
	case common.PM_SERVICE:
		req := request.GetRequestBody(ctx)
		id := request.MustInt(ctx, "id")
		columns := map[string]string{
			"api_type":    "string",
			"method":      "string",
			"url":         "string",
			"resource":    "string",
			"description": "string",
		}
		updateData := custom_util.CopyMapColumns(*req, columns)

		model := mysql_model.SysApiBusiness{}
		if _, err := mysql.GetSession().ID(id).Table(model).Update(updateData); err == nil {
			if has, _ := mysql.FindById(id, &model); has {
				response.RenderSuccess(ctx, model)
			} else {
				panic("Item not found")
			}
		} else {
			panic(err)
		}
	case common.ADMIN_SERVICE:
		req := request.GetRequestBody(ctx)
		id := request.MustInt(ctx, "id")
		columns := map[string]string{
			"api_type":    "string",
			"method":      "string",
			"url":         "string",
			"resource":    "string",
			"description": "string",
		}
		updateData := custom_util.CopyMapColumns(*req, columns)

		model := mysql_model.SysApi{}
		if _, err := mysql.GetSession().ID(id).Table(model).Update(updateData); err == nil {
			if has, _ := mysql.FindById(id, &model); has {
				response.RenderSuccess(ctx, model)
			} else {
				panic("Item not found")
			}
		} else {
			panic(err)
		}
	default:
		panic("Unknown service name")
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/apis/:service 批量删除某一微服务接口
 * @apiVersion 1.0.0
 * @apiName BulkDelete
 * @apiGroup Api
 *
 * @apiDescription 批量删除某一微服务接口
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   ids  用户id,多个微服务id之间用"\\|"连接(如:"1\\|2\\|3")
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X DELETE http://localhost/api/v1/apis/skygo.auth?ids=1|2|3
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":{
 *				"number":2
 *			}
 *      }
 */
func (this Api) BulkDelete(ctx *gin.Context) {
	serviceName := request.ParamString(ctx, "service")
	switch serviceName {
	case common.PM_SERVICE:
		idStr := strings.Split(request.MustString(ctx, "ids"), "|")
		ids := []int{}
		for _, val := range idStr {
			if id, convErr := strconv.Atoi(val); convErr == nil {
				ids = append(ids, id)
			} else {
				panic(convErr)
			}
		}
		if effectNum, err := mysql.GetSession().Where(builder.In("id", ids)).Delete(new(mysql_model.SysApiBusiness)); err == nil {
			response.RenderSuccess(ctx, &qmap.QM{"number": effectNum})
		} else {
			panic(err)
		}
	case common.ADMIN_SERVICE:
		idStr := strings.Split(request.MustString(ctx, "ids"), "|")
		ids := []int{}
		for _, val := range idStr {
			if id, convErr := strconv.Atoi(val); convErr == nil {
				ids = append(ids, id)
			} else {
				panic(convErr)
			}
		}
		if effectNum, err := mysql.GetSession().Where(builder.In("id", ids)).Delete(new(mysql_model.SysApi)); err == nil {
			response.RenderSuccess(ctx, &qmap.QM{"number": effectNum})
		} else {
			panic(err)
		}
	default:
		panic("Unknown service name")
	}
}
