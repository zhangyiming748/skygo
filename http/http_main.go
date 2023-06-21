package http

import (
	"strconv"

	"skygo_detection/guardian/app/sys_service"

	"github.com/gin-gonic/gin"

	"skygo_detection/common"
	"skygo_detection/http/bootstrap"
	"skygo_detection/http/middleware"
	"skygo_detection/lib/hg_service"
	"skygo_detection/service"
)

// Main gateway入口
func Main() {
	// 依赖初始化
	bootstrap.InitService()

	// interrupt信号监听，启动后，可以注册要执行的func
	service.InitInterruptHandler()

	// term信号监听，启动后，可以注册要执行的func
	service.InitTermHandler()

	// gin服务
	engine := gin.New()
	engine.Use(middleware.Recover())

	// 根据模式来决定
	switch common.CliFlagEnv {
	// 如果是开发模式，不校验授权
	case "dev":
		// do nothing
	default:
		engine.Use(middleware.CheckReferer())
		engine.Use(middleware.Authentication())
	}

	// 请求的body数据都存到ctx中
	engine.Use(middleware.JsonDecode())

	bootstrap.InitResource(engine)
	bootstrap.InitResourceHg(engine)
	bootstrap.InitResourceScreen(engine)
	bootstrap.InitResourceLive(engine)
	bootstrap.InitResourceGps(engine)
	bootstrap.InitResourceBeehive(engine)
	bootstrap.InitResourceTool(engine)
	bootstrap.InitResourceHydra(engine)

	// 启动扫描任务消息管理
	scanManager := hg_service.GetScanManager()
	scanManager.Run()
	defer scanManager.Close()
	msgAnalysis := hg_service.GetMessageAnalysis()
	msgAnalysis.Run()
	defer msgAnalysis.Close()

	// 启动http服务
	httpConfig := sys_service.GetHttpConfig()
	addr := httpConfig.Host + ":" + strconv.Itoa(httpConfig.Port)
	engine.Run(addr)
}
