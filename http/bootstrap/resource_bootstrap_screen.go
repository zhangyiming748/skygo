package bootstrap

import (
	"github.com/gin-gonic/gin"
	"skygo_detection/http/controller"
)

func InitResourceScreen(engine *gin.Engine) {

	routeGroup := engine.Group("/api/v1")
	{
		routeGroup.GET("/dashboard/project/all", new(controller.ScreenProjectController).All)
		routeGroup.GET("/dashboard/project/test_cases", new(controller.ScreenProjectController).TestCase)
		routeGroup.GET("/dashboard/tasks", new(controller.Dashboard).Task)

		routeGroup.GET("/dashboard/vulnerabilities", new(controller.Dashboard).Vul)
		routeGroup.GET("/dashboard/vulnerability/top", new(controller.Dashboard).Top)
		routeGroup.GET("/dashboard/test_case", new(controller.Dashboard).TestCase)

		routeGroup.GET("/dashboard/vehicle", new(controller.Dashboard).Vehicle)
		routeGroup.GET("/dashboard/asset_test_pieces", new(controller.Dashboard).AssetTestPieces)
		routeGroup.GET("/dashboard/asset_test_pieces/vul", new(controller.Dashboard).AssetTestPiecesVul)
		routeGroup.GET("/dashboard/case_daily_total", new(controller.Dashboard).CaseDailyTotal)

		routeGroup.POST("/dashboard/upload", new(controller.Dashboard).Upload)

		//新大屏地址
		routeGroup.GET("/screen/vehicle_test_progress", new(controller.Screen).GetScreenVehicleTestProgress)
		routeGroup.POST("/screen/vehicle_test_progress", new(controller.Screen).CreateScreenVehicleTestProgress)
		routeGroup.DELETE("/screen/vehicle_test_progress", new(controller.Screen).DeleteScreenVehicleTestProgress)
		routeGroup.GET("/screen/piece_test_progress", new(controller.Screen).GetScreenPieceTestProgress)
		routeGroup.POST("/screen/piece_test_progress", new(controller.Screen).CreateScreenPieceTestProgress)
		routeGroup.DELETE("/screen/piece_test_progress", new(controller.Screen).DeleteScreenPieceTestProgress)
		routeGroup.GET("/screen/test_case", new(controller.Screen).GetScreenTestCase)
		routeGroup.POST("/screen/test_case", new(controller.Screen).CreateScreenTestCase)
		routeGroup.DELETE("/screen/test_case", new(controller.Screen).DeleteScreenTestCase)
		routeGroup.GET("/screen/info", new(controller.Screen).GetScreenInfo)
		routeGroup.POST("/screen/info", new(controller.Screen).CreateScreenInfo)
		routeGroup.DELETE("/screen/info", new(controller.Screen).DeleteScreenInfo)
		routeGroup.GET("/screen/vehicle_info", new(controller.Screen).GetScreenVehicleInfo)
		routeGroup.POST("/screen/vehicle_info", new(controller.Screen).CreateScreenVehicleInfo)
		routeGroup.DELETE("/screen/vehicle_info", new(controller.Screen).DeleteScreenVehicleInfo)
		routeGroup.GET("/screen/task_info", new(controller.Screen).GetScreenTaskInfo)
		routeGroup.POST("/screen/task_info", new(controller.Screen).CreateScreenTaskInfo)
		routeGroup.DELETE("/screen/task_info", new(controller.Screen).DeleteScreenTaskInfo)
	}
}
