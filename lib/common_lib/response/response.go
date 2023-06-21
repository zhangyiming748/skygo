package response

import (
	"fmt"
	"runtime/debug"

	"github.com/gin-gonic/gin"

	"skygo_detection/lib/common_lib/common_const"
)

const (
	HTTP_RESPONSE_BODY = "Http-Response-Body"

	RC_SUCCESS = 0  // 成功默认返回码
	RC_FAILURE = -1 // 失败默认返回码
)

func RenderInnerFailure(ctx *gin.Context, err error) {
	m := gin.H{
		"code": -1,
		"msg":  "系统错误",
	}
	if common_const.CliFlagDebug == true {
		m["stack"] = fmt.Sprintf("%s", debug.Stack())
		m["error"] = err.Error()
	}
	ctx.AbortWithStatusJSON(200, m)
}

func RenderFailure(ctx *gin.Context, err error) {
	m := gin.H{
		"code": -1,
		"msg":  err.Error(),
	}
	if common_const.CliFlagDebug == true {
		m["stack"] = fmt.Sprintf("%s", debug.Stack())
	}
	ctx.AbortWithStatusJSON(200, m)
}

func RenderSuccess(ctx *gin.Context, h interface{}) {
	m := gin.H{
		"code": 0,
		"msg":  "",
		"data": h,
	}
	ctx.AbortWithStatusJSON(200, m)
}

func Render(ctx *gin.Context, code int, msg string, h interface{}) {
	m := gin.H{
		"code": code,
		"msg":  msg,
		"data": h,
	}
	ctx.AbortWithStatusJSON(200, m)
}
