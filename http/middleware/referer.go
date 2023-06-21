package middleware

import (
	"net/url"

	"skygo_detection/custom_error"
	"skygo_detection/service"

	"github.com/gin-gonic/gin"
)

// 进行referer的校验
func CheckReferer() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		conf := service.LoadConfig().Referer
		if conf.Enable == true {
			referer := ctx.Request.Referer()
			if referer != "" {
				if rUrl, err := url.Parse(referer); err == nil {
					referer = rUrl.Host
					//从配置文件中读取到全部urls
					urls := conf.Url
					for _, v := range urls {
						if v == referer {
							ctx.Next()
							return
						}
					}
				}
			}
			panic(custom_error.ErrCodeInvalidRefer)
			ctx.Abort()
		} else {
			ctx.Next()
		}
	}
}
