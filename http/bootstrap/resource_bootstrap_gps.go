package bootstrap

import (
	"github.com/gin-gonic/gin"

	"skygo_detection/http/controller/gps"
)

func InitResourceGps(engine *gin.Engine) {
	routeGroup := engine.Group("/api/v1/gps/")
	// task
	{
		svr := new(gps.TaskController)
		routeGroup.POST("task", svr.Create)
		routeGroup.PUT("task/:task_id", svr.Update)
		routeGroup.GET("task/:task_id", svr.GetOne)
		routeGroup.PUT("complete/task/:task_id", svr.Complete)

	}
	// gps 模板
	{
		stc := new(gps.SteerTemplateController)
		routeGroup.GET("template", stc.GetAll)
		routeGroup.POST("template", stc.Create)
		routeGroup.PUT("template/:template_id", stc.Update)
		routeGroup.DELETE("template", stc.Delete)
		routeGroup.POST("template/upload", stc.Upload)
		routeGroup.GET("template/download/:template_id", stc.Download)
	}
	// 实时位置欺骗
	{
		rpd := new(gps.RealtimePositionDeceiveController)
		routeGroup.GET("realtime_position_deceive/start_deceive_task/:task_id", rpd.StartDeceiveTask)
		routeGroup.GET("realtime_position_deceive/get_five/:task_id", rpd.GetFive)
		routeGroup.POST("realtime_position_deceive", rpd.Create)
		routeGroup.POST("realtime_position_deceive/start", rpd.StartDeceive)
		routeGroup.PUT("realtime_position_deceive/:task_id", rpd.StopDeceive)
		// GPS欺骗记录
		routeGroup.GET("realtime_position_deceive", rpd.GetAll)
	}

	{
		svr := new(gps.MotionController)
		routeGroup.POST("set_time", svr.SetTime)
		routeGroup.POST("start", svr.Start)
		routeGroup.POST("close", svr.Close)
		routeGroup.GET("status/:task_id", svr.Status)          // 轮询 设备是否启动成功
		routeGroup.POST("motion_line", svr.Line)               // 生成线路
		routeGroup.GET("motion_line_history", svr.LineHistory) // 线路记录
		routeGroup.POST("cheat", svr.Cheat)                    //运动轨迹欺骗
		routeGroup.GET("cheat/:task_id", svr.GetLatestCheat)   // 获取任务最新的一条欺骗记录

		routeGroup.GET("online/:task_id", svr.Online) // 轮询 设备是否健康

	}
}
