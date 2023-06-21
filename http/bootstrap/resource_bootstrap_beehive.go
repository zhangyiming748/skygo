package bootstrap

import (
	"github.com/gin-gonic/gin"

	"skygo_detection/http/controller/beehive"
)

func InitResourceBeehive(engine *gin.Engine) {
	routeGroup := engine.Group("/api")

	// task 相关
	{
		svr := new(beehive.TaskController)
		routeGroup.POST("v1/beehive/task", svr.Create)
		routeGroup.PUT("v1/beehive/task/:task_id", svr.Update)
		routeGroup.GET("v1/beehive/task/:task_id", svr.GetOne)
		routeGroup.PUT("v1/beehive/complete/task/:task_id", svr.Complete)
	}

	// memo 相关
	{
		svr := new(beehive.MemoController)
		routeGroup.POST("v1/beehive/memo", svr.Create)
		routeGroup.GET("v1/beehive/memo/:task_id", svr.View)
	}

	// sniffer 相关
	{
		svr := new(beehive.SnifferController)
		routeGroup.POST("v1/beehive/gsm_sniffer/scan", svr.Scan)
		routeGroup.POST("v1/beehive/gsm_sniffer/rescan", svr.ReScan)
		routeGroup.POST("v1/beehive/gsm_sniffer/getfreq", svr.GetFreq)
		routeGroup.POST("v1/beehive/gsm_sniffer/gettimer", svr.GetTimer)
		routeGroup.POST("v1/beehive/gsm_sniffer/sniff", svr.StartSniff)
		routeGroup.POST("v1/beehive/gsm_sniffer/getimsi", svr.GetImsi)
		routeGroup.POST("v1/beehive/gsm_sniffer/getsms", svr.GetSms)
		routeGroup.POST("v1/beehive/gsm_sniffer/count", svr.ImsiSmsCount)
		routeGroup.POST("v1/beehive/gsm_sniffer/stopscan", svr.StopScan)
		routeGroup.POST("v1/beehive/gsm_sniffer/stopsniff", svr.StopSniff)
		routeGroup.POST("v1/beehive/gsm_sniffer/close", svr.Close)
		routeGroup.POST("v1/beehive/gsm_sniffer/delimsi", svr.DelImsi)
		routeGroup.POST("v1/beehive/gsm_sniffer/delsms", svr.DelSms)
		routeGroup.GET("v1/beehive/gsm_sniffer/imsi_list", svr.ImsiList)
		routeGroup.GET("v1/beehive/gsm_sniffer/sms_list", svr.SmsList)
		routeGroup.GET("v1/beehive/gsm_sniffer/detail", svr.Detail)
	}

	// 伪基站 相关
	{
		svr := new(beehive.GsmSystemController)
		// 设置启动参数 不会被调用
		routeGroup.POST("v1/beehive/gsm_system/set_config", svr.SetConfig)
		// 启动系统 不会被调用
		routeGroup.GET("v1/beehive/gsm_system/start/:task_id", svr.StartSystem)
		// 启动流程
		routeGroup.POST("v1/beehive/gsm_system/start", svr.Start)
		// 获取设备信息
		routeGroup.GET("v1/beehive/gsm_system/get_device/:task_id", svr.GetDevices)
		// 获取短信数量角标
		routeGroup.GET("v1/beehive/gsm_system/sms_count/:task_id", svr.GetSMSNum)
		// 任务详情主页获取短信按钮(从设备中获取短信填充数据库)
		routeGroup.GET("v1/beehive/gsm_system/get_sms/:task_id", svr.GetSmsButton)
		// 批量模拟短信
		routeGroup.POST("v1/beehive/gsm_system/send_sms", svr.GsmSystemSendSMS)
		// 批量删除短信
		routeGroup.DELETE("v1/beehive/gsm_system/del_sms", svr.BulkDeleteSMS)
		// 根据收件人发件人短信内容模糊搜索
		routeGroup.GET("v1/beehive/gsm_system/search_sms/:key", svr.GsmSystemSearch)
		// 获取短信填充列表
		routeGroup.GET("v1/beehive/gsm_system/sms/:task_id", svr.GetSMS)
		// 关闭系统
		routeGroup.GET("v1/beehive/gsm_system/stop/:task_id", svr.StopSystem)

		// 批量模拟短信收件人下拉列表
		routeGroup.GET("v1/beehive/gsm_system/receiver/:task_id", svr.DeviceList)
		//获取任务状态
		routeGroup.GET("v1/beehive/gsm_system/detail/:task_id", svr.Detail)

	}

	{
		svr := new(beehive.LogController)
		// 获取日志
		routeGroup.GET("v1/beehive/log/:task_id", svr.GetAll)
	}
	// LTE System
	{
		lsc := new(beehive.LteSystemController)
		// 密码破解
		routeGroup.POST("v1/beehive/crack_apn", lsc.CrackApn)
		// 密码配置
		routeGroup.GET("v1/beehive/get_apn/:task_id", lsc.GetApn)
		// 获取破解后的用户名和密码
		routeGroup.GET("v1/beehive/get_list", lsc.GetList)
		// 更新密码表
		routeGroup.POST("v1/beehive/update_password", lsc.UpdatePassword)
		// 写卡操作
		routeGroup.POST("v1/beehive/lte_system", lsc.Create)
		// 启动LTE系统
		routeGroup.GET("v1/beehive/lte_system/start/:task_id", lsc.StartSystem)
		// 获取设备信息
		routeGroup.GET("v1/beehive/lte_system/get_equipment_info/:task_id", lsc.GetOne)
		// 关闭LTE系统
		routeGroup.GET("v1/beehive/lte_system/stop/:task_id", lsc.StopSystem)
		// 开始抓包
		routeGroup.GET("v1/beehive/lte_system/package/:task_id", lsc.GetPackage)
		// 删除包
		routeGroup.DELETE("v1/beehive/lte_system/package", lsc.Delete)
		// 获取包信息
		routeGroup.GET("v1/beehive/lte_system/package", lsc.GetALL)
		// 获取系统状态
		routeGroup.GET("v1/beehive/lte_system/get_system_state/:task_id", lsc.GetSystemState)
		// 获取设备信息列表
		routeGroup.GET("v1/beehive/lte_system/get_equipment_info_list", lsc.GetEquipmentAll)
	}
}
