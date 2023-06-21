package http_ctx

import (
	"context"

	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/metadata"
	"skygo_detection/lib/common_lib/session"
)

const HTTP_CONTEXT = "Http-Context"
const SESSION_ID = "SID"

var (
	CookieMaxAge = 0 //cookie有效时间
)

// 初始化上下文变量
func InitContext(ctx *gin.Context) context.Context {
	sessionId := GetSessionId(ctx)
	md := qmap.QM{
		session.SESSION_ID: sessionId,
		session.REQUEST_ID: session.GenerateRequestId(),
	}
	newCtx := metadata.NewContext(context.Background(), md)
	ctx.Set(HTTP_CONTEXT, newCtx)
	return newCtx
}

func GetSessionId(ctx *gin.Context) string {
	if sessionId, err := ctx.Cookie(SESSION_ID); err == nil {
		return sessionId
	} else {
		sessionId = session.GenerateSessionId()
		ctx.SetCookie(SESSION_ID, sessionId, CookieMaxAge, "", "", false, false)
		return sessionId
	}
}

func GetHttpCtx(ctx *gin.Context) context.Context {
	if tmpCtx := ctx.Value(HTTP_CONTEXT); tmpCtx != nil {
		switch tmpCtx.(type) {
		case context.Context:
			return tmpCtx.(context.Context)
		}
	}
	newCtx := metadata.NewContext(context.Background(), qmap.QM{})
	ctx.Set(HTTP_CONTEXT, newCtx)
	return newCtx
}

func NewOutputContext(ctx *gin.Context) context.Context {
	return session.NewOutputContext(GetHttpCtx(ctx))
}

func GetString(ctx *gin.Context, key string) string {
	return metadata.String(GetHttpCtx(ctx), key)
}

func GetInt64(ctx *gin.Context, key string) int64 {
	return metadata.Int64(GetHttpCtx(ctx), key)
}

func Get(ctx *gin.Context, key string) interface{} {
	return metadata.Value(GetHttpCtx(ctx), key)
}

func Set(ctx *gin.Context, key string, val interface{}) bool {
	return metadata.Set(GetHttpCtx(ctx), key, val)
}

func GetChannelId(ctx *gin.Context) string {
	return session.GetChannelId(GetHttpCtx(ctx))
}

func GetAccountType(ctx *gin.Context) string {
	return session.GetAccountType(GetHttpCtx(ctx))
}

func GetQueryChannelId(ctx *gin.Context) string {
	return session.GetQueryChannelId(GetHttpCtx(ctx))
}

func GetUserId(ctx *gin.Context) int64 {
	return session.GetUserId(GetHttpCtx(ctx))
}

func GetRoleId(ctx *gin.Context) int64 {
	return session.GetRoleId(GetHttpCtx(ctx))
}

func GetUserName(ctx *gin.Context) string {
	return session.GetUserName(GetHttpCtx(ctx))
}

// 获取请求全局唯一id
func GetRequestId(ctx *gin.Context) string {
	return session.GetRequestId(GetHttpCtx(ctx))
}

// 获取当前请求内部id
func GetRequestInnerId(ctx *gin.Context) string {
	return session.GetRequestInnerId(GetHttpCtx(ctx))
}

// 设置日志flag
func SetLogOutputFlag(ctx *gin.Context, probability float32) bool {
	return session.SetLogOutputFlag(GetHttpCtx(ctx), probability)
}

// 设置日志flag
func GetLogOutputFlag(ctx *gin.Context) bool {
	return session.GetLogOutputFlag(GetHttpCtx(ctx))
}
