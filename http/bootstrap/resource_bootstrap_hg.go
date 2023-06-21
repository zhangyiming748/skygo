package bootstrap

import (
	"github.com/gin-gonic/gin"

	"skygo_detection/http/controller"
)

func InitResourceHg(engine *gin.Engine) {
	routeGroup := engine.Group("/api")
	// 车型
	{
		svr := new(controller.AssetVehicleController)
		routeGroup.GET("v1/asset_vehicles", svr.GetAll)
		routeGroup.POST("v1/asset_vehicles", svr.Create)
		routeGroup.GET("v1/asset_vehicles/:id", svr.GetOne)
		routeGroup.PUT("v1/asset_vehicles/:id", svr.Update)
		routeGroup.DELETE("v1/asset_vehicles", svr.BulkDelete)
		routeGroup.GET("v1/asset_vehicle/select_list", svr.SelectList)
	}

	// 测试组件
	{
		svr := new(controller.AssetTestPieceController)
		routeGroup.GET("v1/asset_test_pieces", svr.GetAll)
		routeGroup.POST("v1/asset_test_pieces", svr.Create)
		routeGroup.GET("v1/asset_test_pieces/:id", svr.GetOne)
		routeGroup.PUT("v1/asset_test_pieces/:id", svr.Update)
		routeGroup.DELETE("v1/asset_test_pieces", svr.BulkDelete)

		routeGroup.GET("v1/asset_test_piece/get_by_version_id/:id", svr.GetByVersionId)
		//routeGroup.PUT("v1/asset_test_piece/update_by_version_id/:id", svr.UpdateByVersionId)
		// 文件上传
		routeGroup.POST("v1/asset_test_piece/upload_file", svr.UpdateFile)
		// 上传固件
		routeGroup.POST("v1/asset_test_piece/upload_firmware", svr.UpdateFirmware)
		// 测试组件不分页
		routeGroup.GET("v1/asset_test_piece/all", svr.GetAllWithNotPage)
		// 添加测试件版本
		routeGroup.POST("v1/asset_test_pieces/:id/version", svr.CreateVersion)
		// 更新测试件版本号
		routeGroup.PUT("v1/asset_test_pieces/:id/version", svr.UpdateVersion)
		routeGroup.DELETE("v1/asset_test_pieces/:id/version", svr.BulkDeleteVersion)
		// 删除测试件里版本里的文件
		routeGroup.DELETE("v1/asset_test_piece/version_file", svr.DeleteVersionFile)
		// 固件/文件下载
		routeGroup.POST("/v1/asset_test_piece/download_file", svr.DownloadFirmware)

	}

	// 知识库-安全需求
	{
		svr := new(controller.KnowledgeDemandController)
		routeGroup.GET("v1/knowledge_demands", svr.GetAll)
		routeGroup.POST("v1/knowledge_demands", svr.Create)
		routeGroup.GET("v1/knowledge_demands/:id", svr.GetOne)
		routeGroup.PUT("v1/knowledge_demands/:id", svr.Update)
		routeGroup.DELETE("v1/knowledge_demands", svr.BulkDelete)
		routeGroup.GET("v1/knowledge_demand/select_list", svr.SelectList)

		// 章节级联列表
		routeGroup.GET("v1/knowledge_demand/chapter_tree/:id", svr.ChapterTree)
		// 某个需求的章节分页查询
		routeGroup.GET("v1/knowledge_demand/chapter_all", svr.ChapterAll)
		// 某个需求的章节详情
		routeGroup.POST("v1/knowledge_demand/chapter_one", svr.ChapterOne)
		// 某个需求的章节添加
		routeGroup.PUT("v1/knowledge_demand/chapter_create", svr.ChapterCreate)
		// 某个需求的章节修改
		routeGroup.POST("v1/knowledge_demand/chapter_update", svr.ChapterUpdate)
		// 某个需求的章节删除
		routeGroup.DELETE("v1/knowledge_demand/chapter_delete", svr.ChapterDelete)
		// 某个需求章节下拉列表
		routeGroup.GET("v1/knowledge_demand/chapter_select_list", svr.ChapterSelectList)
		// 某个需求的法规标准编号
		routeGroup.GET("v1/knowledge_demand/code_list", svr.CodetList)
		// 父章节编号下拉列表 todo
		// 章节编号下拉列表 todo

		// 查看需求下的测试用例
		routeGroup.POST("/v1/knowledge_demand/knowledge_test_cases", svr.GetTestCases)
	}

	// 任务测试用例
	{
		svr := new(controller.TaskTestCaseController)
		routeGroup.GET("v1/task/test_cases", svr.GetAll)
		routeGroup.GET("v1/task/test_cases/:id", svr.GetOne)
		routeGroup.PUT("v1/task/test_cases/:id", svr.Update)
	}
}
