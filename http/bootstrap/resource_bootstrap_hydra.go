package bootstrap

import (
	"github.com/gin-gonic/gin"
	"skygo_detection/http/controller"
)

func InitResourceHydra(engine *gin.Engine) {
	routeGroup := engine.Group("/api/v1")
	{
		svr := new(controller.HydraController)
		// 创建任务
		routeGroup.POST("/hydra/create", svr.Create)
		// 获取全部任务
		routeGroup.GET("/hydra/task", svr.GetAll)
		// 可选协议列表
		routeGroup.GET("/hydra/protocol/list", svr.ProtocolList)
		// 上传用户名字典
		routeGroup.POST("/hydra/upload/username", svr.UploadUsername)
		// 上传密码字典
		routeGroup.POST("/hydra/upload/password", svr.UploadPassword)
		// 任务详情
		routeGroup.GET("/hydra/task/:task_id", svr.Detail)
		// 接收并更新任务状态
		routeGroup.POST("/hydra/receive", svr.Recv)
		// 手动取消任务
		routeGroup.POST("/hydra/abort", svr.Abort)
		// 删除任务
		routeGroup.DELETE("/hydra/task", svr.Delete)
		// 编辑任务名称
		routeGroup.PUT("/hydra/task/:task_id", svr.Edit)
	}
}
