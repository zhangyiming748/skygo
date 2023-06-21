package middleware

import (
	"errors"
	"fmt"
	"os"
	"runtime/debug"
	"time"

	"skygo_detection/guardian/app/http/response"
	"skygo_detection/guardian/util"

	"skygo_detection/lib/common_lib/common_const"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/service"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Recover() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		defer func() {
			errMsg := ""
			logLevel := zapcore.InfoLevel
			if err := recover(); err != nil {
				logLevel = zapcore.ErrorLevel
				// 记录panic级别日志
				errMsg = fmt.Sprintf("%v", err)

				// body := gin.H{
				// 	"code": -1,
				// 	"msg":  "系统错误",
				// }
				body := gin.H{
					"code": -1,
					"msg":  errMsg,
				}

				if common_const.CliFlagDebug == true {
					body["stack"] = fmt.Sprintf("%s", debug.Stack())
					body["error"] = errMsg
				}
				ctx.AbortWithStatusJSON(200, body)
			}
			logRequest(errMsg, logLevel, ctx, startTime)
		}()
		ctx.Next()
	}
}

func Recover2() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()
		defer func() {
			errMsg := ""
			logLevel := zapcore.InfoLevel
			if recoverErr := recover(); recoverErr != nil {
				logLevel = zapcore.ErrorLevel
				errMsg = fmt.Sprintf("%v", recoverErr)

				switch recoverErr.(type) {
				case string:
					response.RenderFailure(ctx, errors.New(recoverErr.(string)))
				default:
					response.RenderFailure(ctx, recoverErr.(error))
				}
			}
			if logLevel != zapcore.InfoLevel || http_ctx.GetLogOutputFlag(ctx) {
				logRequest(errMsg, logLevel, ctx, startTime)
			}
		}()

		ctx.Next()
	}
}

func logRequest(errMsg string, level zapcore.Level, ctx *gin.Context, startTime time.Time) {
	logger := service.GetDefaultLogger("http")
	defer logger.Sync()

	if logger.Core().Enabled(level) {
		hostname, _ := os.Hostname()
		logger.Check(level, errMsg).Write(
			zap.String("sn", request.GetHeaderSn(ctx)),
			zap.String("url", ctx.Request.URL.Path),
			zap.String("request_id", http_ctx.GetRequestId(ctx)),
			zap.String("request_body", request.GetRequestBody(ctx).ToString()),
			zap.String("response_body", response.GetResponseBody(ctx)),
			zap.String("hostname", hostname),
			zap.String("ip", util.GetIPAddr()),
			zap.String("start_time", startTime.Format(time.RFC3339)),
			zap.Float32("cost_ms", util.DurationToMilliseconds(time.Since(startTime))),
		)
	}
}
