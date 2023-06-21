package controller

import (
	"errors"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"
	"xorm.io/builder"

	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mysql_model"
)

type Module struct{}

/**
 * apiType http
 * @api {get} /api/v1/modules 查询模块列表
 * @apiVersion 1.0.0
 * @apiName GetAll
 * @apiGroup Module
 *
 * @apiDescription 分页查询模块列表
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/modules
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "list": [
 *             {
 *                 "enable": 0,
 *                 "foreign_link": "/tbox/event",
 *                 "icon_name": "",
 *                 "id": 59,
 *                 "is_menu": 1,
 *                 "name": "TBOX告警",
 *                 "parent_id": 57,
 *                 "rank": 10
 *             }
 *         ],
 *         "meta": {
 *             "count": 20,
 *             "current_page": 1,
 *             "per_page": 20,
 *             "total": 40,
 *             "total_pages": 2
 *         }
 *     }
 * }
 */
func (this Module) GetAll(ctx *gin.Context) {
	session := mysql.GetSession()
	widget := orm.PWidget{}
	widget.SetQueryStr(ctx.Request.URL.RawQuery)

	// all查询，对于每条记录需要查询父类名称
	transformerFunc := func(data qmap.QM) qmap.QM {
		session := mysql.GetSession().Where("id = ?", data.Int("parent_id"))
		res := mysql_model.SysModule{}
		if has, _ := session.Get(&res); has {
			data["parent_name"] = res.Name
		} else {
			data["parent_name"] = ""
		}

		return data
	}

	widget.SetTransformerFunc(transformerFunc)
	all := widget.PaginatorFind(session, &[]mysql_model.SysModule{})
	response.RenderSuccess(ctx, all)
}

/**
 * apiType http
 * @api {get} /api/v1/modules/:id 查询模块信息
 * @apiVersion 1.0.0
 * @apiName GetOne
 * @apiGroup Module
 *
 * @apiDescription 根据id查询某一模块信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   id  模块id
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/modules/1
 *
 * @apiSuccessExample {json} 请求成功示例:
 *            {
 *                "code": 0,
 *                "module": {
 *                      "id": 29,
 *                      "name": "终端日志",
 *                      "rank": 60,
 *                      "icon_name": "",
 *                      "enable": 0,
 *                      "parent_id": 5,
 *                      "foreign_link": "/terminal/log",
 *                      "is_menu": 1,
 *                      "parent_name": "IVI管理"
 *                }
 *             }
 */
func (this Module) GetOne(ctx *gin.Context) {
	session := mysql.GetSession().Where("id = ?", request.ParamInt(ctx, "id"))
	res := new(mysql_model.SysModule)
	if has, _ := session.Get(&res); has {
		response.RenderSuccess(ctx, res)
	} else {
		panic(errors.New("Item not found"))
	}
}

/**
 * apiType http
 * @api {post} /api/v1/modules 添加模块
 * @apiVersion 1.0.0
 * @apiName Create
 * @apiGroup Module
 *
 * @apiDescription 添加一个模块
 *
 * @apiUse authHeader
 *
 * @apiParam {string}       name                模块名称
 * @apiParam {int}          [rank=0]            排序
 * @apiParam {string}       [icon_name]         模块图标名称
 * @apiParam {int=0,1}      [enable=0]          禁用状态(0:正常 1:禁用)
 * @apiParam {int}          [parent_id=0]       父模块ID
 * @apiParam {string}       [foreign_link]      外部链接
 * @apiParam {int=0,1}      [is_menu=0]         是否为菜单(0:是 1:否)
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST -d name=模块名称&rank=50 http://localhost/api/v1/modules
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "name":"模块名称",
 *          "rank":50,
 *          "parent_id":10,
 *          "is_menu":1
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *                {
 *                    "code": 0,
 *                    "data": {
 *                      "id": 29,
 *                      "name": "终端日志",
 *                      "rank": 60,
 *                      "icon_name": "",
 *                      "enable": 0,
 *                      "parent_id": 5,
 *                      "foreign_link": "/terminal/log",
 *                      "is_menu": 1,
 *                      "parent_name": "IVI管理"
 *                    }
 *                }
 */

func (this Module) Create(ctx *gin.Context) {
	newModule := mysql_model.SysModule{
		Name:        request.String(ctx, "name"),
		Rank:        request.Int(ctx, "rank"),
		IconName:    request.String(ctx, "icon_name"),
		ParentId:    request.Int(ctx, "parent_id"),
		ForeignLink: request.String(ctx, "foreign_link"),
		IsMenu:      request.Int(ctx, "is_menu"),
	}
	session := mysql.GetSession()
	if _, err := session.InsertOne(newModule); err == nil {
		if has, res := session.Get(&newModule); has {
			response.RenderSuccess(ctx, res)
		} else {
			panic("Create failure")
		}
	} else {
		panic(err)
	}
}

/**
 * apiType http
 * @api {put} /api/v1/modules/:id 更新模块
 * @apiVersion 1.0.0
 * @apiName Update
 * @apiGroup Module
 *
 * @apiDescription 根据模块id,更新模块信息
 *
 * @apiUse authHeader
 *
 * @apiParam {int}          id                  模块id
 * @apiParam {string}       name                模块名称
 * @apiParam {int}          rank                排序
 * @apiParam {string}       [icon_name]         模块图标名称
 * @apiParam {int=0,1}      [enable=0]          禁用状态(0:正常 1:禁用)
 * @apiParam {int}          [parent_id=0]       父模块ID
 * @apiParam {string}       [foreign_link]      外部链接
 * @apiParam {int=0,1}      [is_menu=0]         是否为菜单(0:是 1:否)
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X PUT -d name=模块名称&rank=50 http://localhost/api/modules/1
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "name":"模块名称",
 *          "rank":50,
 *          "parent_id":10,
 *          "is_menu":1
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *                {
 *                    "code": 0,
 *                    "data": {
 *                      "id": 29,
 *                      "name": "终端日志",
 *                      "rank": 60,
 *                      "icon_name": "",
 *                      "enable": 0,
 *                      "parent_id": 5,
 *                      "foreign_link": "/terminal/log",
 *                      "is_menu": 1,
 *                      "parent_name": "IVI管理"
 *                    }
 *                }
 */
func (this Module) Update(ctx *gin.Context) {
	id := request.ParamInt(ctx, "id")
	requestBody := request.GetRequestBody(ctx)
	columns := map[string]string{
		"name":         "string",
		"rank":         "int",
		"icon_name":    "string",
		"parent_id":    "int",
		"foreign_link": "string",
		"is_menu":      "int",
	}
	updateData := custom_util.CopyMapColumns(*requestBody, columns)

	session := mysql.GetSession()
	if _, err := session.ID(id).Table(new(mysql_model.SysModule)).Update(updateData); err == nil {
		res := new(mysql_model.SysModule)
		if has, res := session.Where("id = ?", id).Get(res); has {
			response.RenderSuccess(ctx, res)
		} else {
			panic("Item not found")
		}
	} else {
		panic(err)
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/modules 删除模块
 * @apiVersion 1.0.0
 * @apiName DeleteBulk
 * @apiGroup Module
 *
 * @apiDescription 批量删除模块
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   ids 模块id,多个模块id用"\\|"链接(如:"1\\|2\\|3")
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X DELETE http://localhost/api/v1/modules?ids=1|2|3
 *
 * @apiSuccessExample {json} 请求成功示例:
 *       {
 *            "code": 0,
 *			  "data":{
 *				"number":1
 *			}
 *       }
 */
func (this Module) DeleteBulk(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	idStr := strings.Split(req.MustString("ids"), "|")
	ids := []int{}
	for _, val := range idStr {
		if id, convErr := strconv.Atoi(val); convErr == nil {
			ids = append(ids, id)
		} else {
			panic(convErr)
		}
	}

	if effectNum, err := mysql.GetSession().Where(builder.In("id", ids)).Delete(new(mysql_model.SysModule)); err == nil {
		response.RenderSuccess(ctx, &qmap.QM{"number": effectNum})
	} else {
		panic(err)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/module/all 查询所有菜单列表
 * @apiVersion 1.0.0
 * @apiName GetAllModules
 * @apiGroup Module
 *
 * @apiDescription 查询所有菜单列表
 *
 * @apiUse authHeader
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/module/all
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "foreign_link": "/dashboard",
 *             "icon_name": "dashboard-o",
 *             "id": 1,
 *             "name": "安全总览",
 *             "rank": 0
 *         }
 *     ]
 * }
 */
func (this Module) GetAllModules(ctx *gin.Context) {
	response.RenderSuccess(ctx, new(mysql_model.SysModule).GetAllModules())
}

/**
 * apiType http
 * @api {get} /api/v1/module/all_menus 查询菜单列表
 * @apiVersion 1.0.0
 * @apiName GetAllMenus
 * @apiGroup Module
 *
 * @apiDescription 查询菜单列表
 *
 * @apiUse authHeader
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/module/all_menus
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "children": [
 *                 {
 *                     "foreign_link": "/event/list",
 *                     "icon_name": "",
 *                     "id": 3,
 *                     "name": "告警事件",
 *                     "rank": 10
 *                 },
 *                 {
 *                     "foreign_link": "/event/tickets",
 *                     "icon_name": "",
 *                     "id": 4,
 *                     "name": "工单管理",
 *                     "rank": 11
 *                 }
 *             ],
 *             "foreign_link": "/event",
 *             "icon_name": "medicine-box-o",
 *             "id": 2,
 *             "name": "安全管理",
 *             "rank": 5
 *         }
 *     ]
 * }
 */
func (this Module) GetAllMenus(ctx *gin.Context) {
	response.RenderSuccess(ctx, new(mysql_model.SysModule).GetModuleTree())
}

/**
 * apiType http
 * @api {get} /api/v1/module/menus 查询当前用户菜单列表
 * @apiVersion 1.0.0
 * @apiName GetCurrentUserMenus
 * @apiGroup Module
 *
 * @apiDescription 查询当前用户菜单列表
 *
 * @apiUse authHeader
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/module/menus
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "children": [],
 *             "foreign_link": "/dashboard",
 *             "icon_name": "dashboard-o",
 *             "id": 1,
 *             "name": "安全总览",
 *             "rank": 0
 *         },
 *         {
 *             "children": [
 *                 {
 *                     "foreign_link": "/event/list",
 *                     "icon_name": "",
 *                     "id": 3,
 *                     "name": "告警事件",
 *                     "rank": 10
 *                 },
 *                 {
 *                     "foreign_link": "/event/tickets",
 *                     "icon_name": "",
 *                     "id": 4,
 *                     "name": "工单管理",
 *                     "rank": 11
 *                 }
 *             ],
 *             "foreign_link": "/event",
 *             "icon_name": "medicine-box-o",
 *             "id": 2,
 *             "name": "安全管理",
 *             "rank": 5
 *         }
 *     ]
 * }
 */
func (this Module) GetCurrentUserMenus(ctx *gin.Context) {
	menus := new(mysql_model.SysModule).GetMenusByRoleID(int(http_ctx.GetRoleId(ctx)))
	response.RenderSuccess(ctx, menus)
}
