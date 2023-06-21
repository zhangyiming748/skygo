package controller

import (
	"errors"

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

type SaasRoleController struct{}

//@auto_generated_api_begin
/**
 * apiType http
 * @api {get} /api/v1/saas_roles 查询所有系统角色列表
 * @apiVersion 1.0.0
 * @apiName GetAll
 * @apiGroup SaasRole
 *
 * @apiDescription 查询所有系统角色列表
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "channel_id": "Q00001",
 *             "id": 1,
 *             "name": "超级管理员",
 *             "parent_id": 0
 *         },
 *         {
 *             "channel_id": "T56205",
 *             "id": 2,
 *             "name": "亿咖通管理员1",
 *             "parent_id": 1
 *         }
 *     ]
 * }
 */
func (this SaasRoleController) GetAll(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	// (*req)["id"] = ctx.Param("id")
	// (*req)["query_params"] = ctx.Request.URL.RawQuery

	ids := mysql_model.GetSubRoleIds(int(session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx))))
	s := mysql.GetSession().OrderBy("id desc")
	s.Where(builder.In("id", ids))

	// 超级管理员对用户的查询不做渠道限制
	if session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx)) != common.SUPER_ADMINISTRATE_ROLE_ID {
		s.And("channel_id = ?", session.GetQueryChannelId(http_ctx.GetHttpCtx(ctx)))
	}

	widget := orm.PWidget{}
	widget.SetQueryStr(ctx.Request.URL.RawQuery)
	all, _ := widget.All(s, &[]mysql_model.SysSaasRole{})
	response.RenderSuccess(ctx, all)
}

/**
 * apiType http
 * @api {get} /api/v1/saas_roles/:id  查询某一系统角色信息
 * @apiVersion 1.0.0
 * @apiName GetOne
 * @apiGroup SaasRole
 *
 * @apiDescription 查询某一系统角色信息
 *
 * @apiUse authHeader
 *
 * @apiParam {int}          id          系统角色id
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "channel_id": "T56205",
 *         "id": 2,
 *         "name": "管理员",
 *         "parent_id": 1
 *     }
 * }
 */
func (this SaasRoleController) GetOne(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")
	// (*req)["query_params"] = ctx.Request.URL.RawQuery

	id := req.MustInt("id")
	ids := mysql_model.GetSubRoleIds(int(session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx))))

	s := mysql.GetSession().Where("id = ?", id).And(builder.In("id", ids))

	model := mysql_model.SysSaasRole{}
	if has, _ := s.Get(&model); has {
		response.RenderSuccess(ctx, model)
	} else {
		err := errors.New("Item not found")
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/saas_roles  添加系统角色
 * @apiVersion 1.0.0
 * @apiName Create
 * @apiGroup SaasRole
 *
 * @apiDescription 添加系统角色
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           [channel_id]    	渠道号
 * @apiParam {string}           [name]    		角色名称
 * @apiParam {int}           	[parent_id]    	父角色id
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "channel_id":"O68468",
 *          "name":"超级管理员",
 *          "parent_id":0
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "channel_id": "O68468",
 *         "id": 5,
 *         "name": "超级管理员",
 *         "parent_id": 1
 *     }
 * }
 */
func (this SaasRoleController) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	// (*req)["id"] = ctx.Param("id")
	(*req)["query_params"] = ctx.Request.URL.RawQuery

	createChannelId := session.GetQueryChannelId(http_ctx.GetHttpCtx(ctx))

	// 只有超级管理员管理员或者奇虎用户才能创建其他渠道的用户
	if session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx)) == common.SUPER_ADMINISTRATE_ROLE_ID || createChannelId == "" {
		createChannelId = req.String("channel_id")
	}
	newItem := mysql_model.SysSaasRole{
		ChannelId: createChannelId,
		Name:      req.MustString("name"),
		ParentId:  req.MustInt("parent_id"),
	}
	if _, err := mysql.GetSession().InsertOne(&newItem); err == nil {
		response.RenderSuccess(ctx, custom_util.StructToMap2(newItem))
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {put} /api/v1/saas_roles/:id  更新系统角色
 * @apiVersion 1.0.0
 * @apiName Update
 * @apiGroup SaasRole
 *
 * @apiDescription 更新系统角色
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           [channel_id]    	渠道号
 * @apiParam {string}           [name]    		    角色名称
 * @apiParam {int}           	[parent_id]    	    父角色id
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "channel_id":"O68468",
 *          "name":"超级管理员1",
 *          "parent_id":0
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *          "id":1,
 *          "channel_id":"O68468",
 *          "name":"超级管理员1",
 *          "parent_id":0
 *     }
 * }
 */
func (this SaasRoleController) Update(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")
	// (*req)["query_params"] = ctx.Request.URL.RawQuery

	id := req.MustInt("id")
	// 校验当前用户是否有管理该角色的权限
	if id == 0 || !mysql_model.HasRoleManagePrivilege(int(session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx))), id) {
		panic("没有管理该角色的权限")
	}
	createChannelId := session.GetQueryChannelId(http_ctx.GetHttpCtx(ctx))
	// 只有超级管理员管理员或者奇虎用户才能创建其他渠道的用户
	if session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx)) == common.SUPER_ADMINISTRATE_ROLE_ID || createChannelId == "" {
		createChannelId = req.String("channel_id")
	}
	updateData := qmap.QM{}
	if parentId, has := req.TryInt("parent_id"); has {
		// 只能将当前角色更新为自己管辖范围内的另一个角色
		if mysql_model.HasRoleManagePrivilege(int(session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx))), parentId) {
			updateData["parent_id"] = parentId
		} else {
			panic("不能将当前角色设为自己的上层角色")
		}
	}
	if name, has := req.TryString("name"); has {
		updateData["name"] = name
	}

	updateData["channel_id"] = createChannelId

	if _, err := mysql.GetSession().ID(id).Table(new(mysql_model.SysSaasRole)).Update(updateData); err == nil {
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
}

/**
 * apiType http
 * @api {delete} /api/v1/saas_roles 批量删除系统角色
 * @apiVersion 1.0.0
 * @apiName DeleteBulk
 * @apiGroup SaasRole
 *
 * @apiDescription 批量删除系统角色
 *
 * @apiUse authHeader
 *
 * @apiParam {[]int}   ids  用户id
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "ids":[10,11]
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":{
 *				"number":2
 *			}
 *      }
 */
func (this SaasRoleController) DeleteBulk(ctx *gin.Context) {
	// todo 这个删除逻辑，是有问题的，后续进行角色大改版的时候，一起改
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	// (*req)["id"] = ctx.Param("id")
	// (*req)["query_params"] = ctx.Request.URL.RawQuery

	ids := req.SliceInt("ids")

	if total, err := mysql.GetSession().Where(builder.In("role_id", ids)).Count(new(mysql_model.SysUser)); err != nil || total > 0 {
		response.RenderFailure(ctx, errors.New("将要删除的角色已经和用户绑定，请解除用户绑定关系后再进行删除操作！"))
		return
	}

	for _, id := range ids {
		// 获取要删除的角色的子角色
		ids2 := mysql_model.GetSubRoleIds(id)
		if _, err := mysql.GetSession().Where(builder.In("id", ids2)).Delete(new(mysql_model.SysSaasRole)); err == nil {
			continue
			// response.RenderSuccess(ctx, qmap.QM{"number": effectNum})
		} else {
			response.RenderFailure(ctx, err)
			return
		}
	}
	response.RenderSuccess(ctx, gin.H{})
}

/**
 * apiType http
 * @api {get} /api/v1/saas_role/mservice_role_detail/:id 查询某一系统角色关联的微服务角色
 * @apiVersion 1.0.0
 * @apiName GetMServiceRoles
 * @apiGroup SaasRole
 *
 * @apiDescription 查询某一系统角色关联的微服务角色
 *
 * @apiUse authHeader
 *
 * @apiParam {int}		id		角色id
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "create_time": 0,
 *             "id": 1,
 *             "op_id": 0,
 *             "saas_role_id": 1,
 *             "service": "mservice_vehicle",
 *             "service_name": "车机卫士",
 *             "service_role_id": 1,
 *             "service_role_name": "管理员",
 *             "update_time": 0
 *         }
 *     ]
 * }
 */
func (this SaasRoleController) GetMServiceRoles(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")
	// (*req)["query_params"] = ctx.Request.URL.RawQuery

	s := mysql.GetSession().Where("saas_role_id = ?", req.MustInt("id"))
	s.And(builder.In("saas_role_id", mysql_model.GetSubRoleIds(int(session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx))))))

	models := []mysql_model.SysSaasRoleDetail{}
	if err := s.Find(&models); err == nil {
		response.RenderSuccess(ctx, models)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/saas_role/mservice_role_detail/:id 更新某一系统角色关联的微服务角色
 * @apiVersion 1.0.0
 * @apiName UpdateMServiceRoles
 * @apiGroup SaasRole
 *
 * @apiDescription 更新某一系统角色关联的微服务角色
 *
 * @apiUse authHeader
 *
 * @apiParam {int}   	id  		        系统角色id
 * @apiParam {array}   	roles  				微服务角色
 * @apiParam {id}   	roles.role_id  		微服务角色id
 * @apiParam {string}   roles.role_name  	微服务角色名称
 * @apiParam {string}   roles.service  		微服务
 * @apiParam {string}   roles.service_name 	微服务名称
 *
 * @apiParamExample {json}  请求参数示例:
 *  {
 *    "roles":[
 *        {
 *            "service": "mservice_vehicle",
 *            "service_name": "车机卫士",
 *            "service_role_id": 2,
 *            "service_role_name": "管理员"
 *        }
 *     ]
 *  }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0
 * }
 */
func (this SaasRoleController) UpdateMServiceRoles(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")
	(*req)["query_params"] = ctx.Request.URL.RawQuery

	saasRoleId := req.MustInt("id")
	if !mysql_model.HasRoleManagePrivilege(int(session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx))), saasRoleId) {
		response.RenderFailure(ctx, errors.New("UltraViresError")) // TODO ERROR
		return
	}

	if _, err := mysql.GetSession().Where("saas_role_id = ?", saasRoleId).Delete(new(mysql_model.SysSaasRoleDetail)); err == nil {
		if roles := req.Slice("roles"); roles != nil {
			saasRoles := []mysql_model.SysSaasRoleDetail{}
			for _, item := range roles {
				temp := qmap.QM(item.(map[string]interface{}))
				saasRoleDetail := mysql_model.SysSaasRoleDetail{
					SaasRoleId:      saasRoleId,
					ServiceRoleId:   temp.MustInt("service_role_id"),
					ServiceRoleName: temp.String("service_role_name"),
					Service:         temp.MustString("service"),
					ServiceName:     temp.String("service_name"),
				}
				saasRoles = append(saasRoles, saasRoleDetail)
			}
			_, err := mysql.GetSession().Insert(saasRoles)
			if err != nil {
				response.RenderFailure(ctx, err)
				return
			} else {
				response.RenderSuccess(ctx, gin.H{})
				return
			}
		}
		response.RenderSuccess(ctx, gin.H{})
		return
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}

//@auto_generated_api_end
