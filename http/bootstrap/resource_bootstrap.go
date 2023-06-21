package bootstrap

import (
	"github.com/gin-gonic/gin"
	"skygo_detection/http/controller"
	"skygo_detection/http/controller/ump"
)

func InitResource(engine *gin.Engine) {
	routeGroup := engine.Group("/api")

	// 验证码
	{
		c := new(controller.CollectionController)
		routeGroup.GET("/v1/captcha", c.GetCaption)
	}

	// 授权证书
	{
		c := new(controller.LicenseController)
		routeGroup.GET("/v1/license/info", c.GetLicense)
		routeGroup.POST("/v1/license/import", c.ImportLicense)
		routeGroup.POST("/v1/license/generate", c.GenerateLicense)
	}
	// 统一登录
	{
		s := new(ump.UmpController)
		// 单点登录
		routeGroup.GET("/sso/login", s.Login)
		// 应用授权
		routeGroup.GET("/user/client/authenticate", s.Authenticate)

	}

	// 用户
	{
		c := new(controller.UserController)
		routeGroup.GET("/v1/users", c.GetAll)
		routeGroup.GET("/v1/users/:id", c.GetOne)
		routeGroup.GET("/v1/user/me", c.GetCurrentUserInfo)
		routeGroup.POST("/v1/user/logout", c.Logout)

		routeGroup.POST("/v1/users", c.Create)
		routeGroup.POST("/v1/user/change_password", c.ChangePassword)
		routeGroup.POST("/v1/user/authenticate", c.Authenticate)

		routeGroup.PUT("/v1/users/:id", c.Update)
		routeGroup.DELETE("/v1/users", c.BulkDelete)
	}
	// dashboard
	{
		c := new(controller.DashboardController)
		routeGroup.GET("/v1/dashboard/summary_info", c.GetSummaryInfo)
	}

	// 菜单
	{
		c := new(controller.Module)
		routeGroup.GET("/v1/modules", c.GetAll)
		routeGroup.GET("/v1/modules/:id", c.GetOne)
		routeGroup.GET("/v1/module/all", c.GetAllModules)
		routeGroup.GET("/v1/module/all_menus", c.GetAllMenus)
		routeGroup.GET("/v1/module/menus", c.GetCurrentUserMenus)

		routeGroup.POST("/v1/modules", c.Create)

		routeGroup.PUT("/v1/modules/:id", c.Update)

		routeGroup.DELETE("/v1/modules", c.DeleteBulk)
	}

	// 获取mservices
	{
		svr := new(controller.MService)
		routeGroup.GET("/v1/mservices", svr.GetAll)
		routeGroup.GET("/v1/mservices/:id", svr.GetOne)
		routeGroup.POST("/v1/mservices", svr.Create)
		routeGroup.PUT("/v1/mservices/:id", svr.Update)
		routeGroup.DELETE("/v1/mservices", svr.BulkDelete)
	}

	// 车厂
	{
		c := new(controller.VehicleFactory)
		routeGroup.GET("/v1/vehicle_factories", c.GetAll)
		routeGroup.GET("/v1/vehicle_factories/:id", c.GetOne)
		routeGroup.POST("/v1/vehicle_factories", c.Create)
		routeGroup.PUT("/v1/vehicle_factories/:id", c.Update)
		routeGroup.DELETE("/v1/vehicle_factories", c.BulkDelete)
	}

	// 系统角色	saas_role
	{
		svr := new(controller.SaasRoleController)
		routeGroup.GET("/v1/saas_roles", svr.GetAll)
		routeGroup.GET("/v1/saas_roles/:id", svr.GetOne)
		routeGroup.POST("/v1/saas_roles", svr.Create)
		routeGroup.PUT("/v1/saas_roles/:id", svr.Update)
		routeGroup.DELETE("/v1/saas_roles", svr.DeleteBulk)
		routeGroup.GET("/v1/saas_role/mservice_role_detail/:id", svr.GetMServiceRoles)
		routeGroup.POST("/v1/saas_role/mservice_role_detail/:id", svr.UpdateMServiceRoles)
	}

	// 各个服务角色	role
	{
		svr := new(controller.Role)
		routeGroup.GET("/v1/roles/:service", svr.GetAll)
		routeGroup.GET("/v1/roles/:service/:id", svr.GetOne)
		routeGroup.POST("/v1/roles/:service", svr.Create)
		routeGroup.PUT("/v1/roles/:service/:id", svr.Update)
		routeGroup.DELETE("/v1/roles/:service", svr.BulkDelete)

		routeGroup.GET("/v1/role_apis/:service", svr.GetAllRoleApi)
		routeGroup.GET("/v1/role_apis/:service/:role_id", svr.GetRoleApi)
		routeGroup.POST("/v1/role_apis/:service/:role_id", svr.UpdateRoleApi)

		routeGroup.GET("/v1/role_modules", svr.GetAllRoleModule)
		routeGroup.GET("/v1/role_modules/:role_id", svr.GetRoleModule)
		routeGroup.POST("/v1/role_modules/:role_id", svr.UpdateRoleModule)
	}

	// 各个服务的接口 apis
	{
		svr := new(controller.Api)
		routeGroup.GET("/v1/apis/:service", svr.GetAll)
		routeGroup.GET("/v1/apis/:service/:id", svr.GetOne)
		routeGroup.POST("/v1/apis/:service", svr.Create)
		routeGroup.PUT("/v1/apis/:service/:id", svr.Update)
		routeGroup.DELETE("/v1/apis/:service", svr.BulkDelete)
	}

	{
		svr := new(controller.ProjectController)
		routeGroup.GET("/v1/projects", svr.GetAll)
		routeGroup.GET("/v1/projects/:id", svr.GetOne)
		routeGroup.POST("/v1/projects", svr.Create)
		routeGroup.PUT("/v1/projects/:id", svr.Update)
		routeGroup.DELETE("/v1/projects", svr.BulkDelete)
		routeGroup.GET("/v1/project/configs", svr.GetConfigs)
		routeGroup.GET("/v1/project/my_project_summary_info", svr.GetMyProjectSummaryInfo)
		routeGroup.GET("/v1/project/my_project_list", svr.GetMyProjectList)
		routeGroup.GET("/v1/project/project_summary_info", svr.GetProjectSummaryInfo)
		routeGroup.GET("/v1/project/project_task_series", svr.GetProjectTaskSeries)
		routeGroup.GET("/v1/project/my_backlog_summary", svr.GetMyBacklogSummary)
		routeGroup.GET("/v1/project/dashboard", svr.DashBoard)
		routeGroup.GET("/v1/project/select_list_project_asset", svr.SelectListProjectAsset)
	}

	// 合规测试
	{
		c := new(controller.ProjectHgTestTaskController)
		routeGroup.GET("/v1/hg_test_tasks", c.GetAll)
		routeGroup.GET("/v1/hg_test_tasks/:id", c.GetOne)
		routeGroup.POST("/v1/hg_test_tasks", c.Create)
		routeGroup.PUT("/v1/hg_test_tasks/:id", c.Update)
		routeGroup.DELETE("/v1/hg_test_tasks/:id", c.Delete)
		routeGroup.DELETE("/v1/hg_test_tasks", c.BulkDelete)
		routeGroup.GET("/v1/hg_test_task/get_status_flow/:id", c.GetStatusFlow)
		routeGroup.GET("/v1/hg_test_task/get_test_case/:id", c.GetTestCase)
		routeGroup.POST("/v1/hg_test_task/complete/:id", c.Complete)
		routeGroup.POST("/v1/hg_test_task/update_test_case", c.UpdateTestCase)
	}

	// 工具
	{
		c := new(controller.ToolTaskController)
		routeGroup.GET("/v1/tool/task/list", c.TaskList)
		routeGroup.POST("/v1/tool/task/create", c.CreateTask)
		routeGroup.POST("/v1/tool/task/del", c.DelTask)
		routeGroup.POST("/v1/tool/task/stop", c.StopTask)
		routeGroup.GET("/v1/tool/task/detail", c.TaskDetail)
		routeGroup.GET("/v1/tool/task/result_list", c.TaskResultList)
		routeGroup.GET("/v1/tool/task/get_bind_result", c.GetBindTaskResultForTestId)
		routeGroup.POST("/v1/tool/task/result_link_test", c.ResultLinkTest)
		routeGroup.POST("/v1/tool/task/result_unlink_test", c.ResultUnLinkTest)
	}

	// 项目评估
	{
		c := new(controller.EvaluateTypeController)
		routeGroup.GET("/v1/evaluate_type/all", c.GetAll)
		routeGroup.GET("/v1/evaluate_types", c.GetPagingAll)
		routeGroup.GET("/v1/evaluate_types/:id", c.GetOne)
		routeGroup.POST("/v1/evaluate_types", c.Create)
		routeGroup.PUT("/v1/evaluate_types/:id", c.Update)
		routeGroup.DELETE("/v1/evaluate_types", c.BulkDelete)
	}

	// 项目评估 -- 资产
	{
		c := new(controller.EvaluateAssetController)
		routeGroup.GET("/v1/evaluate_assets", c.GetAll)
		routeGroup.GET("/v1/evaluate_assets/:id", c.GetOne)
		routeGroup.POST("/v1/evaluate_assets", c.Create)
		routeGroup.PUT("/v1/evaluate_assets/:id", c.Update)
		routeGroup.DELETE("/v1/evaluate_assets", c.BulkDelete)
		routeGroup.DELETE("/v1/evaluate_asset/module_type", c.DeleteModuleType)
		routeGroup.GET("/v1/evaluate_asset/type_asset/:project_id", c.TypeAsset)
		// @auto_generated_api_end
		routeGroup.GET("/v1/evaluate_asset/task_asset/:task_id", c.TaskAsset)
	}

	// 项目评估 -- 资产
	{
		c := new(controller.EvaluateMaterialController)
		routeGroup.GET("/v1/evaluate_materials", c.GetAll)
		routeGroup.GET("/v1/evaluate_materials/:id", c.GetOne)
		routeGroup.POST("/v1/evaluate_materials", c.Create)
		routeGroup.PUT("/v1/evaluate_materials/:id", c.Update)
		routeGroup.DELETE("/v1/evaluate_materials", c.BulkDelete)
		// @auto_generated_api_end
		routeGroup.GET("/v1/evaluate_material/task_id/:id", c.GetAllWithTaskId)
	}

	// 项目评估 -- items
	{
		c := new(controller.EvaluateItemController)
		routeGroup.POST("/v1/evaluate_items/all", c.GetAll)
		routeGroup.GET("/v1/evaluate_items/:id", c.GetOne)
		routeGroup.POST("/v1/evaluate_items", c.Create)
		routeGroup.PUT("/v1/evaluate_items/:id", c.Update)
		routeGroup.DELETE("/v1/evaluate_items", c.BulkDelete)
		routeGroup.GET("/v1/evaluate_item/tag", c.GetTag)
		routeGroup.POST("/v1/evaluate_item/tag", c.EditTag)
		routeGroup.GET("/v1/evaluate_item/navigation", c.GetNavigation)
		routeGroup.POST("/v1/evaluate_item/asset_versions", c.GetAssetVersions)
		routeGroup.POST("/v1/evaluate_item/upsert_record", c.UpsertRecord)
		routeGroup.DELETE("/v1/evaluate_item/bulk_delete_record", c.BulkDeleteRecord)
		routeGroup.POST("/v1/evaluate_item/audit_task_item", c.AuditTaskItem)
		routeGroup.GET("/v1/evaluate_item/audited_items", c.GetAuditedItems)
		routeGroup.GET("/v1/evaluate_item/records", c.GetItemRecords)
		routeGroup.POST("/v1/evaluate_item/audit_record_item", c.AuditRecordItem)
		routeGroup.GET("/v1/evaluate_item/audited_record_items", c.GetAuditedRecordItems)
		routeGroup.POST("/v1/evaluate_item/complete_test", c.CompleteItemTest)
		routeGroup.POST("/v1/evaluate_item/all_ids", c.GetAllIds)
		// @auto_generated_api_end
		routeGroup.GET("/v1/evaluate_item/export/:id", c.Export)
		routeGroup.POST("/v1/evaluate_item/import/:id", c.Import)
		routeGroup.POST("/v1/evaluate_item/templates", c.CreateItemFromTemplate)
		routeGroup.POST("/v1/evaluate_item/pre_bind", c.PreBind)
		routeGroup.POST("/v1/evaluate_item/clear_bind", c.ClearBind)
	}

	// 项目评估 -- module
	{
		c := new(controller.EvaluateModuleController)
		routeGroup.GET("/v1/evaluate_module/all", c.GetAllModuleTree)
		routeGroup.GET("/v1/evaluate_modules", c.GetAll)
		routeGroup.POST("/v1/evaluate_modules", c.Create)
		routeGroup.PUT("/v1/evaluate_modules/:id", c.Update)
		routeGroup.DELETE("/v1/evaluate_modules", c.BulkDelete)
		routeGroup.GET("/v1/evaluate_module/module_name_list", c.GetModuleNameList)
		routeGroup.GET("/v1/evaluate_module/module_type_list", c.GetModuleTypeList)
		routeGroup.GET("/v1/evaluate_module/recommend_code", c.GetRecommendCode)
		routeGroup.POST("/v1/evaluate_module/rename_module_name", c.RenameModuleName)
		// @auto_generated_api_end
		routeGroup.GET("/v1/evaluate_module/project", c.GetProjectModuleTree)
	}

	// 项目评估 -- vulnerabilities
	{
		c := new(controller.EvaluateVulnerabilityController)
		routeGroup.GET("/v1/evaluate_vulnerabilities", c.GetAll)
		routeGroup.GET("/v1/evaluate_vulnerabilities/:id", c.GetOne)
		routeGroup.GET("/v1/evaluate_vulnerability/evaluate_task", c.GetTaskVulAll)
		routeGroup.GET("/v1/evaluate_vulnerability/evaluate_task/:id", c.GetTaskVulOne)
		routeGroup.POST("/v1/evaluate_vulnerability/evaluate_task", c.TaskVulCreate)
		routeGroup.PUT("/v1/evaluate_vulnerability/evaluate_task/:id", c.TaskVulUpdate)
		routeGroup.DELETE("/v1/evaluate_vulnerability/evaluate_task", c.TaskVulBulkDelete)
		routeGroup.GET("/v1/evaluate_vulnerability/item_vulnerabilities", c.GetItemVulnerabilities)
		routeGroup.GET("/v1/evaluate_vulnerability/item_task_vulnerabilities", c.GetItemTaskVulnerabilities)
		// @auto_generated_api_end
		routeGroup.GET("/v1/evaluate_vulnerabilitys/export/:id", c.Export)
		routeGroup.POST("/v1/evaluate_vulnerability/tag", c.EidtTag)
		routeGroup.GET("/v1/evaluate_vulnerability/get_tag/:id", c.GetTag)
		routeGroup.GET("/v1/evaluate_vulnerability/tags", c.GetAllTags)
	}

	// 项目评估 -- vulnerabilities
	{
		c := new(controller.EvaluateModuleTemplateController)
		routeGroup.GET("/v1/evaluate_templates", c.GetAll)
		routeGroup.GET("/v1/evaluate_templates/:id", c.GetOne)
		routeGroup.POST("/v1/evaluate_templates", c.Create)
		routeGroup.DELETE("/v1/evaluate_templates", c.BulkDelete)
	}

	// 项目评估 -- TestCase
	{
		c := new(controller.EvaluateTestCaseController)
		routeGroup.GET("/v1/evaluate_test_cases", c.GetAll)
		routeGroup.GET("/v1/evaluate_test_cases/:id", c.GetOne)
		routeGroup.POST("/v1/evaluate_test_cases", c.Create)
		routeGroup.PUT("/v1/evaluate_test_cases/:id", c.Update)
		routeGroup.DELETE("/v1/evaluate_test_cases", c.BulkDelete)
		routeGroup.POST("/v1/evaluate_test_case/upload", c.Upload)
		routeGroup.POST("/v1/evaluate_test_case/download", c.Download)
	}

	// 项目评估 -- task
	{
		c := new(controller.EvaluateTaskController)
		routeGroup.GET("/v1/evaluate_tasks", c.GetAll)
		routeGroup.GET("/v1/evaluate_tasks/:id", c.GetOne)
		routeGroup.POST("/v1/evaluate_tasks", c.Create)
		routeGroup.PUT("/v1/evaluate_tasks/:id", c.Update)
		routeGroup.DELETE("/v1/evaluate_tasks", c.BulkDelete)
		routeGroup.GET("/v1/evaluate_task/report/:id", c.GetTaskReport)
		routeGroup.GET("/v1/evaluate_task/asset_versions", c.GetAssetVersions)
		routeGroup.GET("/v1/evaluate_task/project_tasks", c.GetProjectTasks)
		routeGroup.GET("/v1/evaluate_task/phase", c.Phase)
		routeGroup.GET("/v1/evaluate_task/status", c.Status)
		routeGroup.GET("/v1/evaluate_task/auditor", c.Auditor)
		routeGroup.GET("/v1/evaluate_task/tester", c.Tester)
		routeGroup.POST("/v1/evaluate_task/assign", c.Assign)
		routeGroup.GET("/v1/evaluate_task/list", c.GetTaskSlice)
		routeGroup.POST("/v1/evaluate_task/audit", c.Audit)
		routeGroup.POST("/v1/evaluate_task/submit", c.Submit)
		routeGroup.GET("/v1/evaluate_task/task_project", c.TaskProject)
		routeGroup.GET("/v1/evaluate_task/my_list", c.GetList)
		routeGroup.GET("/v1/evaluate_task/task_item_list", c.GetTaskItemList)
		routeGroup.GET("/v1/evaluate_task/task_item_info", c.GetTaskItemInfo)
	}

	// 项目评估 -- VulScanner
	{
		c := new(controller.EvaluateVulScannerController)
		routeGroup.GET("/v1/evaluate_vul_scanners", c.GetAll)
		routeGroup.GET("/v1/evaluate_vul_scanners/:id", c.GetOne)
		routeGroup.POST("/v1/evaluate_vul_scanners", c.Create)
		routeGroup.GET("/v1/evaluate_vul_scanner/distribution/:id", c.Distribution)
		routeGroup.GET("/v1/evaluate_vul_scanner/vul_numbers/:id", c.VulNumbers)
		routeGroup.GET("/v1/evaluate_vul_scanner/sys_info/:id", c.GetSysInfo)
		routeGroup.GET("/v1/evaluate_vul_scanner/vul_info/:id", c.GetVulInfo)
		routeGroup.POST("/v1/evaluate_vul_scanner/check_auth", c.CheckAuth)
	}

	// 项目评估 -- VulTask
	{
		c := new(controller.EvaluateVulTaskController)
		routeGroup.GET("/v1/evaluate_vul_tasks", c.GetAll)
		routeGroup.POST("/v1/evaluate_vul_tasks", c.Create)
		routeGroup.DELETE("/v1/evaluate_vul_tasks", c.BulkDelete)
		routeGroup.GET("/v1/evaluate_vul_task/download", c.DownloadTool)
	}

	// 车厂列表信息
	{
		c := new(controller.ProjectFactoryController)
		routeGroup.GET("/v1/project_factories", c.GetAll)
		routeGroup.GET("/v1/project_factories/:id", c.GetOne)
		routeGroup.POST("/v1/project_factories", c.Create)
		routeGroup.PUT("/v1/project_factories/:id", c.Update)
		routeGroup.DELETE("/v1/project_factories", c.BulkDelete)
	}

	// 工具
	{
		svr := new(controller.ToolController)
		routeGroup.GET("/v1/tool/list", svr.List)
		routeGroup.GET("/v1/tool/category", svr.Category)
		routeGroup.GET("/v1/tool/category_tool", svr.CategoryTool)
		routeGroup.GET("/v1/tool/detail", svr.Detail)
		routeGroup.POST("/v1/tool/add", svr.Add)
		routeGroup.POST("/v1/tool/edit", svr.Edit)
		routeGroup.POST("/v1/tool/del", svr.Delete)
		routeGroup.POST("/v1/tool/upload", svr.Upload)
		routeGroup.POST("/v1/tool/edit_tag", svr.EditTag)
		routeGroup.GET("/v1/tool/download_app", svr.DownloadApp)

	}

	// 文件上传
	{
		svr := new(controller.ProjectFileController)
		routeGroup.POST("v1/project_file/upload_project_file", svr.CreateProjectFile)
		routeGroup.POST("v1/project_file/upload", svr.Upload)
		routeGroup.GET("/v1/project_file/download", svr.Download)
		routeGroup.GET("/v1/project_file/image", svr.ViewImage)
		routeGroup.GET("v1/project_file/all", svr.GetProjectList)
		routeGroup.DELETE("v1/project_files", svr.BulkDeleteProjectFile)
		routeGroup.POST("v1/project_file/rename", svr.RenameProjectFile)
	}

	{
		svr := new(controller.ProjectReportController)
		routeGroup.GET("v1/project_reports", svr.GetAll)
		routeGroup.GET("v1/project_reports/:id", svr.GetOne)
		routeGroup.POST("v1/project_reports", svr.Create)
		routeGroup.DELETE("v1/project_reports", svr.BulkDelete)
		routeGroup.POST("v1/project_report/export", svr.Export)
		// @auto_generated_api_end
		routeGroup.GET("v1/project_report/status", svr.Status)
		routeGroup.GET("v1/project_report/phase", svr.Phase)
		routeGroup.POST("v1/project_report/create_phase", svr.CreatePhase)
		routeGroup.POST("v1/project_report/create_node", svr.CreateNode)
		routeGroup.POST("v1/project_report/audit", svr.Audit)
		routeGroup.GET("v1/project_report/node", svr.Node)
		routeGroup.POST("v1/project_report/publish", svr.Publish)
		routeGroup.GET("v1/project_report/list", svr.GetList)
	}

	// 任务
	{
		svr := new(controller.TaskController)
		routeGroup.GET("v1/tasks", svr.GetAll)
		routeGroup.POST("v1/tasks", svr.Create)
		routeGroup.GET("v1/tasks/:id", svr.GetOne)
		routeGroup.PUT("v1/tasks/:id", svr.Update)
		routeGroup.DELETE("v1/tasks", svr.BulkDelete)
		// 完成测试任务
		routeGroup.POST("v1/tasks/:id/finish", svr.Finish)
		// 任务日志
		routeGroup.GET("v1/task_logs", svr.GetLogAll)

		// 获取任务中的测试件
		routeGroup.GET("v1/tasks/:id/asset_test_pieces", svr.GetAssetTestPieces)
		// 获取任务中的测试用例
		routeGroup.POST("v1/tasks/:id/task_cases", svr.GetTaskCases)
		routeGroup.GET("v1/tasks/:id/tool_task_info", svr.GetToolTaskInfo)
		// 获取任务中的测试用例，未完成需要人工接入的部分
		routeGroup.GET("v1/tasks/:id/task_cases/auto", svr.GetTaskCasesByAuto)
		// 获取任务中的测试用例，用例的状态数据
		routeGroup.GET("v1/tasks/:id/task_case/status", svr.GetTaskCasesStatus)
		// 获取任务中的漏洞
		routeGroup.GET("v1/tasks/:id/vul", svr.GetVul)
		// 获取任务中的测试结果
		routeGroup.GET("v1/tasks/:id/test_result", svr.GetTestResult)
		// 创建工具类的任务
		routeGroup.POST("v1/tasks_tool", svr.GetTestResult)

		// 测试任务里 开始任务页面接口 场景
		routeGroup.GET("v1/task/knowledge_scenarios", svr.GetAllScenarios)
		// 测试任务里 开始任务页面接口 工具
		routeGroup.GET("v1/task/tools", svr.GetAllTools)
		// 获取任务中的所有的安全需求
		routeGroup.GET("v1/tasks/:id/demand/list", svr.GetTaskDemand)
		// 按照车型查询任务
		routeGroup.GET("v1/task/vehicle_group", svr.GetTasksByGroup)
		// 测试任务里 任务类型
		routeGroup.GET("v1/task/category", svr.GetCategory)
		routeGroup.GET("v1/task/long_connection_scanners", svr.GetTasksLongConnectionScanner)

		// routeGroup.PUT("v1/task/knowledge_scenarios_add_tag", svr.AddScenariosTag)
	}

	// 漏洞
	{
		svr := new(controller.VulnerabilityController)
		routeGroup.GET("v1/vulnerabilities", svr.GetAll)
		routeGroup.GET("v1/vulnerabilities/:id", svr.GetOne)
		routeGroup.POST("v1/vulnerabilities", svr.Create)
		routeGroup.PUT("v1/vulnerabilities/:id", svr.Update)
		routeGroup.DELETE("v1/vulnerabilities", svr.BulkDelete)
		routeGroup.POST("v1/vulnerability/tag", svr.AddTag)
		routeGroup.GET("v1/vulnerability/tags", svr.GetAdd)

		// 漏洞日志
		routeGroup.GET("v1/vulnerability_logs", svr.GetLogAll)
	}

	// 漏洞类型
	{
		svr := new(controller.ProjectConfigController)
		routeGroup.GET("v1/project_config/all_vul_type", svr.GetAllVulType)
		routeGroup.POST("v1/project_config/upsert_vul_type", svr.UpsertVulType)
		routeGroup.DELETE("v1/project_config/bulk_delete_vul_type", svr.BulkDeleteVulType)
	}

	// 场景
	{
		svr := new(controller.KnowledgeScenarioController)
		routeGroup.GET("v1/knowledge_scenarios", svr.GetAll)
		routeGroup.POST("v1/knowledge_scenarios", svr.Create)
		routeGroup.GET("v1/knowledge_scenarios/:id", svr.GetOne)
		routeGroup.GET("v1/knowledge_scenario/chapter_tree/:id", svr.ChapterTree)
		routeGroup.PUT("v1/knowledge_scenarios/:id", svr.Update)
		routeGroup.DELETE("v1/knowledge_scenarios", svr.BulkDelete)
		routeGroup.POST("v1/ks_tag/:id", svr.UpdateTag)
		// routeGroup.GET("v1/knowledge_scenarios_chapter/:id",svr.GetChapterByScenario)
	}

	// 测试用例
	{
		svr := new(controller.KnowledgeTestCaseController)
		routeGroup.GET("v1/knowledge_test_cases", svr.GetAll)
		routeGroup.POST("v1/knowledge_test_cases", svr.Create)
		routeGroup.GET("v1/knowledge_test_cases/:id", svr.GetOne)
		routeGroup.PUT("v1/knowledge_test_cases/:id", svr.Update)
		routeGroup.DELETE("v1/knowledge_test_cases", svr.BulkDelete)
		// 批量导入测试用例
		routeGroup.POST("/v1/knowledge_test_case/upload", svr.Upload)
		// 批量拷贝测试用例
		routeGroup.POST("/v1/knowledge_test_case/copy", svr.Copy)
		// 添加标签
		routeGroup.POST("v1/knowledge_test_cases/:id/tag", svr.AddTag)
		routeGroup.GET("v1/knowledge_test_this_tools/:id/this_tool_list", svr.ScenarioToolList)
	}

	// 测试用例类型
	{
		svr := new(controller.KnowledgeTestCaseModuleController)
		routeGroup.GET("v1/knowledge_test_case/modules", svr.GetAll)
		routeGroup.POST("v1/knowledge_test_case/modules", svr.Create)
		routeGroup.GET("v1/knowledge_test_case/modules/:id", svr.GetOne)
		routeGroup.PUT("v1/knowledge_test_case/modules/:id", svr.Update)
		routeGroup.DELETE("v1/knowledge_test_case/modules", svr.BulkDelete)
	}

	// 固件扫描
	{
		svr := new(controller.FirmwareController)
		routeGroup.POST("/v1/firmware/cancel_scanning", svr.CancelScanning)
		routeGroup.POST("/v1/firmware/start_scanning", svr.StartScanning)
		routeGroup.POST("/v1/firmware/del_task", svr.DelTask)
		routeGroup.GET("/v1/firmware/list", svr.List)
		routeGroup.GET("/v1/firmware/basic", svr.Basic)
		routeGroup.GET("/v1/firmware/category_detail", svr.CategoryDetail)
		routeGroup.GET("/v1/firmware/analysis_detail", svr.AnalysisDetail)
		routeGroup.GET("/v1/firmware/analysis_download", svr.AnalysisDownload)
		routeGroup.GET("/v1/firmware/analysis_detail_page", svr.AnalysisDetailPage)
		routeGroup.GET("/v1/firmware/analysis_tracker", svr.AnalysisTracker)
		routeGroup.POST("/v1/firmware/upload_firmware_msg", svr.UploadFirmWareMsg)
		routeGroup.POST("/v1/firmware/upload_firmware_file", svr.FirmWareUpload)
		routeGroup.POST("/v1/firmware/update_upload_status", svr.UpdateUploadStatus)
		routeGroup.GET("/v1/firmware/download", svr.DownloadFirmWare)
		routeGroup.GET("/v1/firmware/apk_basic_info", svr.ApkBasicInfo)
		routeGroup.GET("/v1/firmware/apk_common_vue", svr.ApkCommonVue)
		routeGroup.GET("/v1/firmware/apk_common_vue_detail", svr.ApkCommonVueDetail)
		// V2
		routeGroup.GET("/v2/firmware/basic", svr.BasicV2)
		routeGroup.GET("/v2/firmware/category_detail", svr.CategoryDetailV2)
		routeGroup.GET("/v2/firmware/analysis_tracker", svr.AnalysisTrackerV2)
		routeGroup.GET("/v2/firmware/analysis_detail_page", svr.AnalysisDetailPageV2)
		routeGroup.GET("/v2/firmware/apk_common_vue", svr.ApkCommonVueV2)
		routeGroup.GET("/v2/firmware/apk_common_vue_detail", svr.ApkCommonVueDetailV2)
		routeGroup.GET("/v2/firmware/apk_basic_info", svr.ApkBasicInfoV2)
		routeGroup.GET("/v2/firmware/templates", svr.Template)
	}

	// 测试结果
	{
		svr := new(controller.ReportController)
		// 测试报告
		routeGroup.GET("/v1/report/result/pass_rate", svr.GetResultPASSRate)
		routeGroup.GET("/v1/report/result/distribution", svr.GetResultDistribution)
		routeGroup.GET("/v1/report/result/view", svr.GetResultView)
		routeGroup.GET("/v1/report/result/detail", svr.GetResultDetail)
		routeGroup.GET("/v1/reports", svr.GetAll)
		routeGroup.GET("v1/reports_task_if_generate/:id", svr.IfGenerate)
		routeGroup.POST("v1/report", svr.Create) // 生成报告
	}

}
