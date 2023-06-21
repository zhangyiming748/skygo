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
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mysql_model"
)

type Role struct{}

/**
 * apiType http
 * @api {get} /api/v1/roles/:service 查询某个微服务角色列表
 * @apiVersion 1.0.0
 * @apiName GetAll
 * @apiGroup Role
 *
 * @apiDescription 分页查询某个微服务角色列表
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/roles/skygo.auth
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "id": 2,
 *             "name": "管理员"
 *         },
 *         {
 *             "id": 1,
 *             "name": "test22"
 *         }
 *     ]
 * }
 */
func (this Role) GetAll(ctx *gin.Context) {
	serviceName := request.ParamString(ctx, "service")

	switch serviceName {
	case common.PM_SERVICE:
		s := mysql.GetSession().OrderBy("id desc")
		widget := orm.PWidget{}
		widget.SetQueryStr(ctx.Request.URL.RawQuery)
		all, _ := widget.All(s, &[]mysql_model.SysRoleBusiness{})
		response.RenderSuccess(ctx, all)
	case common.ADMIN_SERVICE:
		s := mysql.GetSession().OrderBy("id desc")
		widget := orm.PWidget{}
		widget.SetQueryStr(ctx.Request.URL.RawQuery)
		all, _ := widget.All(s, &[]mysql_model.SysRole{})
		response.RenderSuccess(ctx, all)
	default:
		panic("Unknown service name")
	}
}

/**
 * apiType http
 * @api {get} /api/v1/roles/:service/:id 查询某个微服务的某一角色信息
 * @apiVersion 1.0.0
 * @apiName GetOne
 * @apiGroup Role
 *
 * @apiDescription 根据id查询某一微服务的某一角色信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   service  	服务名称
 * @apiParam {string}   id  		角色id
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/roles/skygo.auth/1
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 1,
 *         "name": "test"
 *     }
 * }
 */
func (this Role) GetOne(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")
	// (*req)["query_params"] = ctx.Request.URL.RawQuery

	id := req.MustInt("id")

	serviceName := request.ParamString(ctx, "service")
	switch serviceName {
	case common.PM_SERVICE:
		s := mysql.GetSession().Where("id = ?", id)

		model := mysql_model.SysRoleBusiness{}
		if has, _ := s.Get(&model); has {
			response.RenderSuccess(ctx, model)
		} else {
			err := errors.New("Item not found")
			response.RenderFailure(ctx, err)
		}
	case common.ADMIN_SERVICE:
		s := mysql.GetSession().Where("id = ?", id)

		model := mysql_model.SysRole{}
		if has, _ := s.Get(&model); has {
			response.RenderSuccess(ctx, model)
		} else {
			err := errors.New("Item not found")
			response.RenderFailure(ctx, err)
		}
	default:
		panic("Unknown service name")
	}
}

/**
 * apiType http
 * @api {post} /api/v1/roles/:service 向某个微服务添加角色
 * @apiVersion 1.0.0
 * @apiName Create
 * @apiGroup Role
 *
 * @apiDescription 向某个微服务添加角色
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           [name]    		角色名称
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST -d api_type=http&method=POST&url=/api/v1/test&resource=test&description=接口描述 http://localhost/api/v1/roles
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "name":"test"
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 1,
 *         "name": "test"
 *     }
 * }
 */
func (this Role) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	serviceName := request.ParamString(ctx, "service")
	switch serviceName {
	// 之前的逻辑，就是在pm项目的mysql表中存的两个角色，项目管理者和项目测试者, 在不大改的前提下，把表贴过来，叫做sys_role_business
	case common.PM_SERVICE:
		newItem := mysql_model.SysRoleBusiness{
			Name: req.MustString("name"),
		}

		if _, err := mysql.GetSession().InsertOne(&newItem); err == nil {
			response.RenderSuccess(ctx, custom_util.StructToMap2(newItem))
		} else {
			response.RenderFailure(ctx, err)
		}
	case common.ADMIN_SERVICE:
		newItem := mysql_model.SysRole{
			Name: req.MustString("name"),
		}

		if _, err := mysql.GetSession().InsertOne(&newItem); err == nil {
			response.RenderSuccess(ctx, custom_util.StructToMap2(newItem))
		} else {
			response.RenderFailure(ctx, err)
		}
	default:
		panic("Unknown service name")
	}
}

/**
 * apiType http
 * @api {put} /api/v1/roles/:service/:id 更新微服务
 * @apiVersion 1.0.0
 * @apiName Update
 * @apiGroup Role
 *
 * @apiDescription 根据id,更新微服务
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           name    		角色名称
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X PUT -d api_type=http&method=POST&url=/api/v1/test&resource=test&description=接口描述 http://localhost/api/v1/roles/skygo.auth/1
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "name":"test"
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 1,
 *         "name": "test"
 *     }
 * }
 */
func (this Role) Update(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")
	// (*req)["query_params"] = ctx.Request.URL.RawQuery

	serviceName := request.ParamString(ctx, "service")
	switch serviceName {
	case common.PM_SERVICE:
		id := req.MustInt("id")
		columns := map[string]string{
			"name": "string",
		}

		updateData := custom_util.CopyMapColumns(*req, columns)
		if _, err := mysql.GetSession().ID(id).Table(new(mysql_model.SysRoleBusiness)).Update(updateData); err == nil {
			model := mysql_model.SysSaasRole{}
			if has, res := mysql.FindById(id, &model); has {
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
	case common.ADMIN_SERVICE:
		id := req.MustInt("id")
		columns := map[string]string{
			"name": "string",
		}

		updateData := custom_util.CopyMapColumns(*req, columns)
		if _, err := mysql.GetSession().ID(id).Table(new(mysql_model.SysRole)).Update(updateData); err == nil {
			model := mysql_model.SysRole{}
			if has, res := mysql.FindById(id, &model); has {
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
	default:
		panic("Unknown service name")
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/roles/:service 批量删除某一微服务角色
 * @apiVersion 1.0.0
 * @apiName BulkDelete
 * @apiGroup Role
 *
 * @apiDescription 批量删除某一微服务角色
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   ids  用户id,多个微服务id之间用"\\|"连接(如:"1\\|2\\|3")
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X DELETE http://localhost/api/v1/roles/skygo.auth?ids=1|2|3
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":{
 *				"number":2
 *			}
 *      }
 */
func (this Role) BulkDelete(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	// (*req)["id"] = ctx.Query("id")
	// (*req)["query_params"] = ctx.Request.URL.RawQuery

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
		if effectNum, err := mysql.GetSession().Where(builder.In("id", ids)).Delete(new(mysql_model.SysRoleBusiness)); err == nil {
			response.RenderSuccess(ctx, qmap.QM{"number": effectNum})
		} else {
			response.RenderFailure(ctx, err)
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
		if effectNum, err := mysql.GetSession().Where(builder.In("id", ids)).Delete(new(mysql_model.SysRole)); err == nil {
			response.RenderSuccess(ctx, qmap.QM{"number": effectNum})
		} else {
			response.RenderFailure(ctx, err)
		}
	default:
		panic("Unknown service name")
	}
}

/**
 * apiType http
 * @api {get} /api/v1/role_apis/:service 查询某一服务的角色权限
 * @apiVersion 1.0.0
 * @apiName GetAllRoleApi
 * @apiGroup Role
 *
 * @apiDescription 查询某一服务的角色权限
 *
 * @apiUse authHeader
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X DELETE http://localhost/api/v1/role_apis/skygo.auth
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "id": 2,
 *             "name": "管理员",
 *             "privilege": "1|2"
 *         },
 *         {
 *             "id": 1,
 *             "name": "test22",
 *             "privilege": "1|2|3"
 *         }
 *     ]
 * }
 */
func (this Role) GetAllRoleApi(ctx *gin.Context) {
	serviceName := request.ParamString(ctx, "service")
	switch serviceName {
	case common.PM_SERVICE:
		s := mysql.GetSession()
		widget := orm.PWidget{}
		widget.SetQueryStr(ctx.Request.URL.RawQuery)
		all, _ := widget.All(s, &[]mysql_model.SysRoleBusiness{})
		roleApis := []interface{}{}
		for _, role := range all {
			rolePrivilege := qmap.QM{}
			id := role["id"]
			idid := id.(int)
			rolePrivilege["id"] = id
			rolePrivilege["name"] = role["name"]
			rolePrivilege["privilege"] = custom_util.Join(new(mysql_model.SysRoleApiBusiness).GetRoleApi(idid), "|")
			roleApis = append(roleApis, rolePrivilege)
		}
		response.RenderSuccess(ctx, qmap.QM{"data": roleApis})
	case common.ADMIN_SERVICE:
		s := mysql.GetSession()
		widget := orm.PWidget{}
		widget.SetQueryStr(ctx.Request.URL.RawQuery)
		all, _ := widget.All(s, &[]mysql_model.SysRole{})
		roleApis := []interface{}{}
		for _, role := range all {
			rolePrivilege := qmap.QM{}
			id := role["id"]
			idid := id.(int)
			rolePrivilege["id"] = id
			rolePrivilege["name"] = role["name"]
			rolePrivilege["privilege"] = custom_util.Join(new(mysql_model.SysRoleApi).GetRoleApi(idid), "|")
			roleApis = append(roleApis, rolePrivilege)
		}
		response.RenderSuccess(ctx, qmap.QM{"data": roleApis})
	default:
		panic("Unknown service name")
	}
}

/**
 * apiType http
 * @api {get} /api/v1/role_apis/:service/:role_id 查询某一服务某一角色的权限
 * @apiVersion 1.0.0
 * @apiName GetRoleApi
 * @apiGroup Role
 *
 * @apiDescription 查询某一服务某一角色的权限
 *
 * @apiUse authHeader
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X DELETE http://localhost/api/v1/role_apis/skygo.auth/1
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *             "id": 1,
 *             "name": "test22",
 *             "privilege": "1|2|3"
 *     }
 * }
 */
func (this Role) GetRoleApi(ctx *gin.Context) {
	roleId := request.ParamInt(ctx, "role_id")
	s := mysql.GetSession().Where("id = ?", roleId)

	serviceName := request.ParamString(ctx, "service")

	switch serviceName {
	case common.PM_SERVICE:
		model := mysql_model.SysRoleBusiness{}
		if has, _ := s.Get(&model); has {
			res := qmap.QM{}
			res["id"] = model.Id
			res["name"] = model.Name
			res["privilege"] = custom_util.Join(new(mysql_model.SysRoleApiBusiness).GetRoleApi(roleId), "|")
			response.RenderSuccess(ctx, &res)
		} else {
			err := errors.New("Item not found")
			response.RenderFailure(ctx, err)
		}
	case common.ADMIN_SERVICE:
		model := mysql_model.SysRole{}
		if has, _ := s.Get(&model); has {
			res := qmap.QM{}
			res["id"] = model.Id
			res["name"] = model.Name
			res["privilege"] = custom_util.Join(new(mysql_model.SysRoleApi).GetRoleApi(roleId), "|")
			response.RenderSuccess(ctx, &res)
		} else {
			err := errors.New("Item not found")
			response.RenderFailure(ctx, err)
		}
	default:
		panic("Unknown service name")
	}
}

/**
 * apiType http
 * @api {post} /api/v1/role_apis/:service/:role_id 更新某一服务某一角色的接口权限
 * @apiVersion 1.0.0
 * @apiName UpdateRoleApi
 * @apiGroup Role
 *
 * @apiDescription 更新某一服务某一角色的接口权限
 *
 * @apiUse authHeader
 *
 * @apiParam {int}      role_id     角色id
 * @apiParam {string}   privilege   由模块id组成的菜单权限,,多个模块id由"\\|"连接(如"1\\|2\\|3")
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST http://localhost/api/v1/role_apis/skygo.auth/1
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "privilege":"1|2|3"
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0
 * }
 */
func (this Role) UpdateRoleApi(ctx *gin.Context) {
	roleId := request.ParamInt(ctx, "role_id")
	privileges := custom_util.SplitInt(request.MustString(ctx, "privilege"), "|")

	serviceName := request.ParamString(ctx, "service")

	switch serviceName {
	case common.PM_SERVICE:
		if err := new(mysql_model.SysRoleApiBusiness).UpdateRoleApi(roleId, privileges); err == nil {
			response.RenderSuccess(ctx, new(qmap.QM))
		} else {
			panic(err)
		}
	case common.ADMIN_SERVICE:
		if err := new(mysql_model.SysRoleApi).UpdateRoleApi(roleId, privileges); err == nil {
			response.RenderSuccess(ctx, new(qmap.QM))
		} else {
			panic(err)
		}
	default:
		panic("Unknown service name")
	}
}

/**
 * apiType http
 * @api {get} /api/v1/role_modules 查询所有角色的菜单权限
 * @apiVersion 1.0.0
 * @apiName GetAllRoleModule
 * @apiGroup Role
 *
 * @apiDescription 查询所有角色的菜单权限
 *
 * @apiUse authHeader
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X DELETE http://localhost/api/v1/role_modules
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "id": 2,
 *             "name": "管理员",
 *             "modules": "1|2"
 *         },
 *         {
 *             "id": 1,
 *             "name": "test22",
 *             "modules": "1|2|3"
 *         }
 *     ]
 * }
 */
func (this Role) GetAllRoleModule(ctx *gin.Context) {
	s := mysql.GetSession().OrderBy("id desc")

	widget := orm.PWidget{}
	widget.SetQueryStr(ctx.Request.URL.RawQuery)
	all, _ := widget.All(s, &[]mysql_model.SysSaasRole{})
	roleModules := []interface{}{}
	for _, role := range all {
		rolePrivilege := qmap.QM{}
		id := role["id"]
		idid := id.(int)
		rolePrivilege["id"] = id
		rolePrivilege["name"] = role["name"]
		rolePrivilege["modules"] = custom_util.Join(new(mysql_model.SysRoleModule).GetRoleModule(idid), "|")
		roleModules = append(roleModules, rolePrivilege)
	}
	response.RenderSuccess(ctx, roleModules)
}

/**
 * apiType http
 * @api {get} /api/v1/role_modules/:role_id 查询某一角色的菜单权限
 * @apiVersion 1.0.0
 * @apiName GetRoleModule
 * @apiGroup Role
 *
 * @apiDescription 查询某一角色的菜单权限
 *
 * @apiUse authHeader
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X DELETE http://localhost/api/v1/role_modules/1
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *             "id": 2,
 *             "name": "管理员",
 *             "modules": "1|2"
 *     }
 * }
 */
func (this Role) GetRoleModule(ctx *gin.Context) {
	roleId := request.ParamInt(ctx, "role_id")

	s := mysql.GetSession().Where("id = ?", roleId)

	model := mysql_model.SysSaasRole{}
	if has, _ := s.Get(&model); has {
		res := qmap.QM{}
		res["id"] = model.Id
		res["name"] = model.Name
		res["modules"] = custom_util.Join(new(mysql_model.SysRoleModule).GetRoleModule(roleId), "|")
		response.RenderSuccess(ctx, res)
	} else {
		err := errors.New("Item not found")
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/role_modules/:role_id 更新某一角色的菜单权限
 * @apiVersion 1.0.0
 * @apiName UpdateRoleModule
 * @apiGroup Role
 *
 * @apiDescription 更新某一角色的菜单权限
 *
 * @apiUse authHeader
 *
 * @apiParam {int}      role_id     	角色id
 * @apiParam {string}   [modules]   	由模块id组成的菜单权限,,多个模块id由"\\|"连接(如"1\\|2\\|3")
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST http://localhost/api/v1/role_modules/1
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "modules":"1|2|3"
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0
 * }
 */
func (this Role) UpdateRoleModule(ctx *gin.Context) {
	roleId := request.ParamInt(ctx, "role_id")
	modules := custom_util.SplitInt(request.MustString(ctx, "modules"), "|")
	if err := new(mysql_model.SysRoleModule).UpdateRoleModule(roleId, modules); err == nil {
		response.RenderSuccess(ctx, nil)
	} else {
		panic(err)
	}
}
