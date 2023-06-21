package middleware

import (
	"errors"
	"strings"

	"github.com/gin-gonic/gin"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/session"
	"skygo_detection/logic"
	"skygo_detection/mysql_model"
	"skygo_detection/service"
)

type LoginUserInfo struct {
	Id             int      `json:"id"`
	PermissionList []string `json:"permissionList"`
}

type LoginUserInfoResponse struct {
	Code int           `json:"code"`
	Data LoginUserInfo `json:"data"`
}

var whiteRequestList = []string{
	"/api/v1/user/authenticate",
	"/api/v1/captcha",
	"/api/v1/project_file/image",
	"/api/v1/firmware/download",
	"/api/v1/user/logout",
	"/api/v1/project_file/download",
	"/api/v1/project_file/upload",
	"/message/v1/hg_scanner/terminal",
	"/message/v1/hg_scanner/web",
	"/message/v1/hg_scanner/download_case",
	"/message/v1/hg_scanner/upload",
	"/api/v1/evaluate_vul_scanners",
	"/api/v1/evaluate_vul_scanner/check_auth",
	"/message/v1/privacy/terminal",
	"/message/v1/privacy/analysis_record",
	"/api/user/client/authenticate",
	"/api/v1/hydra/receive",
}

// 用户信息认证
func Authentication() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 白名单检查
		if inWhiteList := isInList(ctx.Request.URL.Path, whiteRequestList); inWhiteList {
			ctx.Next()
			return
		}

		// 不在白名单的uri访问，验证token
		if err := ParseToken(ctx); err != nil {
			panic(err)
		}

		// 验证api访问权限
		if err := CheckPrivilege(ctx); err != nil {
			panic(err)
		}

		ctx.Next()
	}
}

// token处理，解析后把用户信息存到cxt中
func ParseToken(ctx *gin.Context) error {
	// 获取token
	token := request.GetAuthToken(ctx)
	if token == "" {
		return errors.New("TokenNotFound") // todo
	}

	// 检验用户token
	res, err := new(logic.AuthLogic).VerifyToken(token)
	if err != nil {
		return errors.New("TokenFailed") // todo
	}

	// 检验成功，用户信息存到ctx中
	http_ctx.Set(ctx, session.USER_INFO_ID, res.UserId)
	http_ctx.Set(ctx, session.USER_GLOBAL_ROLE_ID, res.GlobalRoleId)
	http_ctx.Set(ctx, session.USER_ROLE_ID, res.RoleId)
	http_ctx.Set(ctx, session.USER_INFO_NAME, res.Username)
	http_ctx.Set(ctx, session.USER_INFO_CHANNEL_ID, res.ChannelId)
	http_ctx.Set(ctx, session.USER_INFO_HMAC, res.Hmd5)
	http_ctx.Set(ctx, session.USER_ACCOUNT_TYPE, res.AccountType)
	http_ctx.SetLogOutputFlag(ctx, service.LoadConfig().Log.OutputProbability)
	return nil
}

// 检查用户访问权限
func CheckPrivilege(ctx *gin.Context) error {
	saasRoleId := session.GetRoleId(http_ctx.GetHttpCtx(ctx)) // 这个获取的是用户的saas_role_id

	// 对于权限校验，实际上都是在校验后台程序用户角色有啥权限，因此要拿到这个role_id
	roleId := new(mysql_model.SysSaasRoleDetail).GetMServiceRoleId(saasRoleId, common.ADMIN_SERVICE)

	// 超级管理员不校验任何uri检查
	if roleId == common.SUPER_ADMINISTRATE_ROLE_ID {
		return nil
	}

	if has := new(mysql_model.SysRoleApi).CheckPrivilege(roleId, formatUrl(ctx, ctx.Request.URL.Path), ctx.Request.Method); !has {
		return errors.New("PermissionDeny") // todo
	}
	return nil
}

func isInList(dest string, list []string) bool {
	for _, elem := range list {
		if elem == dest {
			return true
		}
	}
	return false
}

func formatUrl(ctx *gin.Context, url string) string {
	for _, param := range ctx.Params {
		url = strings.Replace(url, "/"+param.Value, "/"+":"+param.Key, -1)
	}
	return url
}
