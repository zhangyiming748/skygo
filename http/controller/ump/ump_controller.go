package ump

import (
	"github.com/gin-gonic/gin"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/license"
	"skygo_detection/logic/ump"
)

type UmpController struct{}

// 单点登录
func (this *UmpController) Login(ctx *gin.Context) {
	response.RenderSuccess(ctx, license.GetLicense())
}

// 应用授权
func (this *UmpController) Authenticate(ctx *gin.Context) {
	rawAccessToken := ctx.Request.FormValue("access_token")
	data, err := ump.Authenticate(rawAccessToken, ctx)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, data)
	return
}
