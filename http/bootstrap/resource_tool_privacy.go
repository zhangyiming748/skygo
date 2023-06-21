package bootstrap

import (
	"github.com/gin-gonic/gin"
	"skygo_detection/http/controller/toolbox"
)

func InitResourceTool(engine *gin.Engine) {
	routeGroup := engine.Group("/api/v1/privacy/")
	{
		pc := new(toolbox.PrivacyController)
		// 应用隐私调用统计
		routeGroup.GET("app_list", pc.AppList)
		// 应用调用频次统计
		routeGroup.GET("app_per_list", pc.AppPerList)
		// 应用权限请求列表
		routeGroup.GET("per_count_list", pc.PerCountList)
		// 应用请求记录
		routeGroup.GET("app_record", pc.RecordList)
		// 应用请求权限和次数统计
		routeGroup.GET("app_count", pc.AppCount)
		// 获取日志
		routeGroup.GET("log/:task_id", pc.GetLog)
	}

	// task 相关
	{
		svr := new(toolbox.TaskController)
		routeGroup.POST("task", svr.Create)
		routeGroup.PUT("task/:task_id", svr.Update)
		routeGroup.GET("task/:task_id", svr.GetOne)
		routeGroup.PUT("complete/task/:task_id", svr.Complete)
		routeGroup.POST("task/start", svr.Start)
		routeGroup.POST("task/stop", svr.Stop)
	}

	// memo 相关
	{
		svr := new(toolbox.MemoController)
		routeGroup.POST("memo", svr.Create)
		routeGroup.GET("memo/:task_id", svr.View)
	}

	// 版本对比
	{
		svr := new(toolbox.AppVersionController)
		// 获取任务选中的应用
		routeGroup.GET("app/task/:task_id", svr.GetAllApp)
		// 获取选中应用的所有版本
		routeGroup.GET("version/app/:app_name", svr.GetAppAllVersion)
		// 根据选中应用版本返回列表
		routeGroup.POST("task_app/permission", svr.GetAppVersionCompare)
	}

	routeGroup = engine.Group("/message")
	// 长连接
	{
		pc := new(toolbox.PrivacyController)
		// 用于日志分析处理
		routeGroup.POST("/v1/privacy/analysis_record", pc.AnalysisRecord)
		svr := new(toolbox.TaskSocketController)
		routeGroup.GET("/v1/privacy/terminal", svr.Terminal)
	}
}
