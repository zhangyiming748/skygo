package license

const (
	ASSET_MANAGE   = iota + 1 // 资产管理
	TEST_TASK                 // 测试任务
	REPORT_MANAGE             // 报告管理
	VUL_MANAGE                // 漏洞管理
	KNOWLEDGE_CASE            // 检测知识库
)

var MenuMap = map[int]string{
	ASSET_MANAGE:   "资产管理",
	TEST_TASK:      "测试任务",
	REPORT_MANAGE:  "报告管理",
	VUL_MANAGE:     "漏洞管理",
	KNOWLEDGE_CASE: "检测知识库",
}
