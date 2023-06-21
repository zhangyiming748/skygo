package controller

import (
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"
	"xorm.io/builder"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	mycaptcha "skygo_detection/lib/common_lib/captcha"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/common_lib/session"
	"skygo_detection/mysql_model"
	"skygo_detection/service"
)

type UserController struct{}

/**
 * apiType http
 * @api {post} api/v1/user/authenticate 用户登录接口
 * @apiVersion 1.0.0
 * @apiName Authenticate
 * @apiGroup User
 *
 * @apiDescription 用户登录接口
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   captcha_text  	验证码
 * @apiParam {string}   username  		用户名
 * @apiParam {string}   password  		密码
 *
 * @apiExample {curl} 请求示例:
 * curl -i -x POST http://localhost/api/v1/user/authenticate
 *
 * @apiSuccessExample {json} 返回值示例：
 *	{
 *		"code": 0,
 *		"data": {
 *			"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1NzA2MDMzMjIsImp0aSI6IjEiLCJpYXQiOjE1NzA1MTY5MjJ9.SWjkxMjeH2ATq76aoHEKGk6upBQ6dkWxzcEfc5EQD9M"
 *		}
 *	}
 */
func (this UserController) Authenticate(ctx *gin.Context) {
	username, password := request.GetBasicAuthAccount(ctx)
	req := qmap.QM{
		"username": username,
		"password": password,
	}

	if _, e := ctx.Cookie("captcha_"); e != nil {
		response.RenderFailure(ctx, errors.New("验证码已过期"))
		return
	}

	captchaText := ctx.PostForm("captcha_text")
	isCaptcha := mycaptcha.Verify(ctx, captchaText)

	if isCaptcha == false {
		response.RenderFailure(ctx, errors.New("验证码错误"))
		return
	}

	if user := new(mysql_model.SysUser).GetUserFindByUsername(req.MustString("username")); user != nil {
		if err := service.CheckPassword(user.Password, req.MustString("password")); err == nil {
			if token, err := service.GenerateJWT(fmt.Sprintf("%d", user.Id), "", user.AuthorizeTime); err == nil {
				userInfo := &qmap.QM{
					"id":         user.Id,
					"username":   user.Username,
					"realname":   user.Realname,
					"channel_id": user.ChannelId,
					"nickname":   user.Nickname,
					"email":      user.Email,
					"status":     user.Status,
				}
				resp := qmap.QM{
					"token": token,
					"user":  userInfo,
				}
				response.RenderSuccess(ctx, resp)
				return
			} else {
				panic(err)
			}
		} else {
			response.RenderFailure(ctx, err)
			return
		}
	}
	// response.RenderFailure(ctx, errors.New("AccountAuthenticateError"))
	response.RenderFailure(ctx, errors.New("用户密码不正确"))
}

/**
* apiType http
* @api {get} /api/v1/users 查询用户列表
* @apiVersion 1.0.0
* @apiName GetAll
* @apiGroup User
*
* @apiDescription 分页查询用户列表
*
* @apiUse authHeader
*
* @apiUse urlQueryParams
*
* @apiExample {curl} 请求示例:
* curl -i http://localhost/api/v1/users
*
* @apiSuccessExample {json} 请求成功示例:
 * {
*     "code": 0,
*     "data": {
*         "list": [
*             {
*                 "account_type": "user",
*                 "channel_id": "T56205",
*                 "create_time": 1557372989,
*                 "email": "",
*                 "head_pic": "",
*                 "id": 20,
*                 "mobile": 0,
*                 "nickname": "亿咖通mno管理",
*                 "realname": "亿咖通",
*				   "role_id": 1,
*                 "sex": 0,
*                 "status": 2,
*                 "username": "ecarx_mno_admin"
*             }
*         ],
*         "meta": {
*             "count": 3,
*             "current_page": 1,
*             "per_page": 20,
*             "total": 8,
*             "total_pages": 1
*         }
*     }
* }
*/
func (this UserController) GetAll(ctx *gin.Context) {
	s := mysql.GetSession()
	s.Where(builder.In("role_id"), mysql_model.GetSubRoleIds(int(session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx)))))

	// 超级管理员对用户的查询不做渠道限制
	if session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx)) != common.SUPER_ADMINISTRATE_ROLE_ID {
		s.Where("channel_id = ?", session.GetQueryChannelId(http_ctx.GetHttpCtx(ctx)))
	}

	widget := orm.PWidget{}
	widget.SetQueryStr(ctx.Request.URL.RawQuery)
	widget.SetTransformerFunc(userTransform)
	all := widget.PaginatorFind(s, &[]mysql_model.SysUser{})
	response.RenderSuccess(ctx, all)
}

func userTransform(data qmap.QM) qmap.QM {
	data["channel_name"] = new(mysql_model.SysVehicleFactory).GetChannelName(data.String("channel_id"))
	data["role_name"] = new(mysql_model.SysSaasRole).GetSaasRoleName(data.Int("role_id"))
	delete(data, "password")
	return data
}

/**
 * apiType http
 * @api {get} /api/v1/users/:id 查询用户信息
 * @apiVersion 1.0.0
 * @apiName GetOne
 * @apiGroup User
 *
 * @apiDescription 根据id查询某一用户信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   id  用户id
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/users/1
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "account_type": "user",
 *         "channel_id": "T56205",
 *         "create_time": 1551865814,
 *         "email": "",
 *         "head_pic": "",
 *         "id": 1,
 *         "mobile": 0,
 *         "nickname": "test",
 *         "realname": "test",
 *		   "role_id": 1,
 *         "sex": 0,
 *         "status": 2,
 *         "username": "test"
 *     }
 * }
 */
func (this UserController) GetOne(ctx *gin.Context) {
	id := ctx.Param("id")

	s := mysql.GetSession().Table(mysql_model.SysUser{})
	s.Where("id = ?", id)
	s.Where(builder.In("role_id"), mysql_model.GetSubRoleIds(int(session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx)))))

	// 超级管理员对用户的查询不做渠道限制
	if session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx)) != common.SUPER_ADMINISTRATE_ROLE_ID {
		s.Where("channel_id = ?", session.GetQueryChannelId(ctx))
	}

	widget := orm.PWidget{}
	widget.SetQueryStr(ctx.Request.URL.RawQuery)
	widget.SetTransformerFunc(userTransform)
	result, _ := widget.Get(s)
	response.RenderSuccess(ctx, result)
}

/**
 * apiType http
 * @api {post} /api/v1/users 添加用户
 * @apiVersion 1.0.0
 * @apiName Create
 * @apiGroup User
 *
 * @apiDescription 添加用户
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           username                用户名
 * @apiParam {string}           channel_id              渠道号
 * @apiParam {string}           role_id	              	角色id
 * @apiParam {string}           password                密码
 * @apiParam {string}           [nickname]              昵称
 * @apiParam {string}           realname                真实姓名
 * @apiParam {string}           [email]                 邮件
 * @apiParam {string}           [mobile]                手机号
 * @apiParam {string}           [head_pic]              头像地址
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST -d username=用户名&password=123456&realname=李明&role_id=2 http://localhost/api/users
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "username":"用户名",
 *          "password":"123456",
 *          "realname":"李明",
 *          "role_id":2,
 *          "nickname":"车联网",
 *          "email":"liming@gmail.com",
 *          "channel_id":"T56205"
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "channel_id": "Q00001",
 *         "email": "",
 *         "id": 37,
 *         "nickname": "",
 *         "realname": "test",
 *         "username": "test4121"
 *     }
 * }
 */
func (this UserController) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	channelId := session.GetQueryChannelId(ctx)
	createChannelId := req.MustString("channel_id")
	if createChannelId == "" {
		response.RenderFailure(ctx, errors.New("channel id cannot be null"))
		return
	} else if session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx)) != common.SUPER_ADMINISTRATE_ROLE_ID && channelId != "" && channelId != createChannelId {
		// 只有超级管理员管理员或者奇虎用户才能创建其他渠道的用户
		response.RenderFailure(ctx, errors.New("PermissionDeny"))
	}
	roleId := req.MustInt("role_id")
	// 只能创建不高于当前用户角色权限的用户
	if !mysql_model.HasRoleManagePrivilege(int(session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx))), roleId) {
		response.RenderFailure(ctx, errors.New("UltraViresError"))
		return
	}

	username := req.MustString("username")
	newUser := mysql_model.SysUser{
		Username:    username,
		ChannelId:   createChannelId,
		RoleId:      roleId,
		Realname:    req.MustString("realname"),
		Password:    service.HashPassword(req.MustString("password")),
		AccountType: req.DefaultString("account_type", "user"),
		Nickname:    req.String("nickname"),
		Email:       req.String("email"),
		Status:      mysql_model.USER_STATUS_NORMAL,
	}
	if newUser.AccountType == "" {
		newUser.AccountType = "user"
	}

	s := mysql.GetSession().Table(mysql_model.SysUser{})
	if _, err := s.InsertOne(newUser); err == nil {
		widget := orm.PWidget{}
		widget.SetTransformerFunc(func(qm qmap.QM) qmap.QM {
			delete(qm, "password")
			return qm
		})
		s.Table(mysql_model.SysUser{})
		result, _ := widget.Get(s)
		response.RenderSuccess(ctx, result)
		return
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {put} /api/v1/users/:id 更新用户
 * @apiVersion 1.0.0
 * @apiName Update
 * @apiGroup User
 *
 * @apiDescription 根据用户id,更新用户信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           [username]                  用户名
 * @apiParam {string}           [nickname]                  昵称
 * @apiParam {string}           [realname]                  真实姓名
 * @apiParam {string}           [email]                     邮件
 * @apiParam {string}           [mobile]                    手机号
 * @apiParam {string}           [head_pic]                  头像地址
 * @apiParam {int=0,1,2,99}     [status=2]                  状态(0:待审核 1:禁用 2:正常 99:删除)
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X PUT -d username=用户名&password=123456&realname=李明 http://localhost/api/v1/users/{2
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "username":"用户名",
 *          "realname":"李明",
 *          "nickname":"车联网",
 *          "email":"liming@gmail.com"
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *                {
 *                    "code": 0,
 *                    "data": {
 *                      "id": 116,
 *                      "username": "qihoo",
 *                      "nickname": "",
 *                      "realname": "奇虎360",
 *                      "email": "",
 *                      "mobile": 0,
 *                      "head_pic": "http://qa.admin.v2.adlab.com/assets/images/head_default.png",
 *                      "sex": 0,
 *                      "status": 2,
 *                      "create_time": 1487242836,
 *                    }
 *                }
 */
func (this UserController) Update(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")
	(*req)["query_params"] = ctx.Request.URL.RawQuery

	id := req.MustInt("id")
	// 只能修改同等及以下权限的角色用户信息
	if !mysql_model.HasRoleManagePrivilege(int(session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx))), new(mysql_model.SysUser).GetUserRoleId(id)) {
		response.RenderFailure(ctx, errors.New("UltraViresError"))
		return
	}

	columns := map[string]string{
		"realname": "string",
		"nickname": "string",
		"email":    "string",
		"status":   "int",
	}
	updateData := custom_util.CopyMapColumns(*req, columns)

	if updateChannelId, has := req.TryString("channel_id"); has {
		if updateChannelId == session.GetChannelId(http_ctx.GetHttpCtx(ctx)) {
			updateData["channel_id"] = updateChannelId
		} else if session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx)) == common.SUPER_ADMINISTRATE_ROLE_ID || session.GetQueryChannelId(http_ctx.GetHttpCtx(ctx)) == "" {
			// 只有超级管理员或者奇虎用户才能创建其他渠道的用户
			updateData["channel_id"] = updateChannelId
		}
	}
	if roleId, has := req.TryInt("role_id"); has {
		// 新更新的角色权限不能高于当前用户角色权限
		if !mysql_model.HasRoleManagePrivilege(int(session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx))), roleId) {
			response.RenderFailure(ctx, errors.New("UltraViresError"))
			return
		} else {
			updateData["role_id"] = roleId
		}
	}

	s := mysql.GetSession()
	s = s.Where("id =?", id)

	// 超级管理员对用户的查询不做渠道限制
	if session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx)) != common.SUPER_ADMINISTRATE_ROLE_ID {
		s.And("channel_id = ?", session.GetQueryChannelId(http_ctx.GetHttpCtx(ctx)))
	}

	if _, err := mysql.GetSession().Table(new(mysql_model.SysUser)).ID(id).Update(updateData); err == nil {

		s := mysql.GetSession().Table(mysql_model.SysUser{})

		widget := orm.PWidget{}
		widget.SetQueryStr(ctx.Request.URL.RawQuery)
		widget.SetTransformerFunc(func(qm qmap.QM) qmap.QM {
			delete(qm, "password")
			return qm
		})
		one, _ := widget.Get(s)
		response.RenderSuccess(ctx, one)
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/users 批量删除用户
 * @apiVersion 1.0.0
 * @apiName BulkDelete
 * @apiGroup User
 *
 * @apiDescription 批量删除用户
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   ids  用户id,多个用户id之间用"\\|"连接(如:"1\\|2\\|3")
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X DELETE http://localhost/api/v1/users?ids=1|2|3
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":{
 *				"number":2
 *			}
 *      }
 */
func (this UserController) BulkDelete(ctx *gin.Context) {
	req := &qmap.QM{
		"ids": strings.Split(request.MustString(ctx, "ids"), "|"),
	}

	ids := req.SliceInt("ids")

	s := mysql.GetSession().Where(builder.In("id", ids)).And(builder.In("role_id", mysql_model.GetSubRoleIds(int(session.GetGlobalRoleId(http_ctx.GetHttpCtx(ctx)))))).
		And(builder.NotIn("id", session.GetUserId(http_ctx.GetHttpCtx(ctx))))
	if total, _ := s.Count(new(mysql_model.SysUser)); total < int64(len(ids)) {
		response.RenderFailure(ctx, errors.New("对所选用户没有删除权限"))
		return
	}

	if effectNum, err := mysql.GetSession().Where(builder.In("id", ids)).Delete(new(mysql_model.SysUser)); err == nil {
		response.RenderSuccess(ctx, qmap.QM{"number": effectNum})
		return
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}

/**
 * apiType http
 * @api {post} /api/v1/user/change_password 修改密码
 * @apiVersion 1.0.0
 * @apiName ChangePassword
 * @apiGroup User
 *
 * @apiDescription 密码修改
 *
 * @apiHeader {string}      Authorization       用户名和密码以HTTP Basic认证的方式放入http请求中
 *
 * @apiParam {string}       old_password        旧密码
 * @apiParam {string}       new_password        新密码
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST -d old_password=sfaw&new_password=asdfad http://localhost/api/user/change_password
 *
 * @apiParamExample {json} 请求参数示例:
 *      {
 *         "old_password":"oldpassword",
 *         "new_password":"newpassword",
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *          "code": 0,
 *      }
 */
func (this UserController) ChangePassword(ctx *gin.Context) {
	req := &qmap.QM{
		"old_password": request.MustString(ctx, "old_password"),
		"new_password": request.MustString(ctx, "new_password"),
	}

	oldPassword := req.MustString("old_password")
	newPassword := req.MustString("new_password")

	if err := new(mysql_model.SysUser).ChangePassword(session.GetQueryChannelId(http_ctx.GetHttpCtx(ctx)), int(session.GetUserId(http_ctx.GetHttpCtx(ctx))), oldPassword, newPassword); err != nil {
		response.RenderFailure(ctx, err)
	}
	response.RenderSuccess(ctx, gin.H{})
}

/**
 * apiType http
 * @api {get} /api/v1/user/me 查询当前用户信息
 * @apiVersion 1.0.0
 * @apiName GetCurrentUserInfo
 * @apiGroup User
 *
 * @apiDescription 查询当前用户信息
 *
 * @apiUse authHeader
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/user/me
 *
 * @apiSuccessExample {json} 请求成功示例:
 *            {
 *                "code": 0,
 *                "data": {
 *                      "id": 116,
 *                      "username": "qihoo",
 *                      "nickname": "",
 *                      "realname": "奇虎360",
 *                      "email": "",
 *                      "mobile": 0,
 *                      "status": 2,
 *                      "create_time": 1487242836,
 *                }
 *             }
 */
func (this UserController) GetCurrentUserInfo(ctx *gin.Context) {
	if user, has := new(mysql_model.SysUser).FindById(int(session.GetUserId(http_ctx.GetHttpCtx(ctx)))); has {
		res := &qmap.QM{
			"id":         user.Id,
			"role_id":    user.RoleId,
			"username":   user.Username,
			"realname":   user.Realname,
			"channel_id": user.ChannelId,
			"nickname":   user.Nickname,
			"email":      user.Email,
			"status":     user.Status,
		}
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, errors.New("AccountNotFound"))
	}
}

/**
 * apiType http
 * @api {post} /api/v1/user/logout 注销
 * @apiVersion 1.0.0
 * @apiName Logout
 * @apiGroup User
 *
 * @apiDescription 注销
 *
 * @apiUse authHeader
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST http://localhost/api/v1/user/logout
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *          "code": 0
 *      }
 */
func (this UserController) Logout(ctx *gin.Context) {
	response.RenderSuccess(ctx, nil)
}
