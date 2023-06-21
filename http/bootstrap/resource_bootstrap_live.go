package bootstrap

import (
	"github.com/gin-gonic/gin"

	"skygo_detection/http/controller"
)

func InitResourceLive(engine *gin.Engine) {
	routeGroup := engine.Group("/message")
	// 长连接
	{
		svr := new(controller.ScanController)
		routeGroup.GET("/v1/hg_scanner/terminal", svr.Terminal)
		routeGroup.GET("/v1/hg_scanner/web", svr.Web)
		routeGroup.GET("/v1/hg_scanner/download_case", svr.DownloadCase)
		routeGroup.POST("/v1/hg_scanner/upload", svr.Upload)
		routeGroup.GET("/v1/hg_scanner/terminal_info", svr.GetTerminalInfo)
	}
}
