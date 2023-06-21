package common

const (
	// 任务状态
	TASK_STATUS_RUNNING = 1 // 进行中
	TASK_STATUS_SUCCESS = 2 // 成功
	TASK_STATUS_FAILURE = 3 // 失败
	TASK_STATUS_REMOVE  = 4 // 删除
)

const (
	TOOL_FIRMWARE_SCANNER_NAME   = "固件扫描工具"
	TOOL_VUL_SCANNER_NAME        = "车机漏扫检测工具"
	TOOL_HG_ANDROID_SCANNER_NAME = "车机检测工具" // 车机检测工具安卓版
	TOOL_HG_LINUX_SCANNER_NAME   = "车机检测工具Linux版"
)

// 测试工具 test_tools
const (
	TOOL_FIRMWARE_SCANNER   = "firmware_scanner"
	TOOL_VUL_SCANNER        = "vul_scanner"
	TOOL_HG_ANDROID_SCANNER = "hg_scanner" // hg_android_scanner
	TOOL_HG_Linux_SCANNER   = "hg_linux_scanner"
)

const (
	HG_TEST_TASK_STATUS_CREATE           = "create"           // 任务状态 - 创建
	HG_TEST_TASK_STATUS_CLIENT_INFO      = "client_info"      // 任务状态 - 获取信息
	HG_TEST_TASK_STATUS_CHOOSE_TEST_CASE = "choose_test_case" // 任务状态 - 适配用例
	HG_TEST_TASK_STATUS_AUTO_TEST        = "auto_test"        // 任务状态 - 测试中
	HG_TEST_TASK_STATUS_COMPLETE         = "complete"         // 任务状态 - 完成

	HG_TEST_TASK_CONNECT_STATUS_NEVER = 1 // 连接状态 - 未连接
	HG_TEST_TASK_CONNECT_STATUS_YES   = 2 // 连接状态 - 已连接
	HG_TEST_TASK_CONNECT_STATUS_NO    = 3 // 连接状态 - 连接断开

	HG_TEST_TASK_FLOW_NO  = 1 // 状态流程节点状态 - 未完成
	HG_TEST_TASK_FLOW_YES = 2 // 状态流程节点状态 - 完成

	HG_TEST_TASK_CASE_BLOCK_STATUS_SUCCESS = 1 // 任务中测试案例Block的状态 - 失败
	HG_TEST_TASK_CASE_BLOCK_STATUS_FAIL    = 2 // 任务中测试案例Block的状态 - 成功
	HG_TEST_TASK_CASE_BLOCK_STATUS_RUNNING = 3 // 任务中测试案例Block的状态 - 运行中

	// 测试用例状态
	CASE_STATUS_READY     = 1 // 待测试
	CASE_STATUS_QUEUING   = 2 // 队列中
	CASE_STATUS_TESTING   = 3 // 测试中
	CASE_STATUS_ANALYSIS  = 4 // 分析中
	CASE_STATUS_COMPLETED = 5 // 测试完成
	CASE_STATUS_FAIL      = 6 // 测试失败
	CASE_STATUS_INVALID   = 7 // 无效
	CASE_STATUS_CLOSED    = 8 // 已关闭
	CASE_STATUS_CANCELED  = 9 // 已取消

	// 测试用例测试结果状态
	CASE_TEST_STATUS_PASS   = 1 // 通过
	CASE_TEST_STATUS_UNPASS = 2 // 未通过

	// 测试block状态
	BLOCK_STATUS_READY    = 1 // 待测试
	BLOCK_STATUS_TESTING  = 2 // 测试中
	BLOCK_STATUS_ANALYSIS = 3 // 分析中
	BLOCK_STATUS_PASS     = 4 // 通过
	BLOCK_STATUS_UNPASS   = 5 // 未通过
	BLOCK_STATUS_FAILURE  = 6 // 测试失败
)

var TaskTestCaseStatusList = map[int]string{
	CASE_STATUS_READY:     "待测试",
	CASE_STATUS_TESTING:   "测试中",
	CASE_STATUS_ANALYSIS:  "分析中",
	CASE_STATUS_COMPLETED: "测试完成",
	CASE_STATUS_FAIL:      "测试失败",
}

const (
	CASE_HG_PRE       = "hg_"
	CASE_FIRMWARE_PRE = "fw_"
	CASE_VUL_PRE      = "vul_"
)

// 优先级
const (
	CASE_PRIORITY_DEFAULT = 0
	CASE_PRIORITY_HG      = 1
	CASE_PRIORITY_MD      = 2
	CASE_PRIORITY_LW      = 3
)

// 是否是工具任务
const (
	IS_TOOL_TASK  = 1
	NOT_TOOL_TASK = 0
)

// 测试任务是否是人工和自动
const (
	IS_TASK_CASE_MAN  = 1 // 人工
	IS_TASK_CASE_SEMI = 2 // 半人工
	IS_TASK_CASE_AUTO = 3 // 自动
)

// 固件扫描任务状态
// 1 待上传 2 上传完成 3 上传失败 4 取消上传 5 (下载完成) 扫描中 6 (创建任务) 扫描中 7 取消扫描 8 扫描完成 9 扫描失败 10 已解析 0 已删除
const (
	FIRMWARE_STATUS_PROJECT_CREATE    = 1 // 项目创建
	FIRMWARE_STATUS_FIRMWARE_DOWNLOAD = 2 // 固件下载
	FIRMWARE_STATUS_TASK_CREATE       = 3 // 任务创建
	FIRMWARE_STATUS_TASK_START        = 4 // 任务启动
	FIRMWARE_STATUS_TASK_CANCEL       = 5 // 任务取消
	FIRMWARE_STATUS_TASK_SCANNING     = 6 // 任务扫描中
	FIRMWARE_STATUS_SCAN_FAILURE      = 7 // 扫描失败
	FIRMWARE_STATUS_REPORT_ANALYSIS   = 8 // 扫描完成(报告解析)
	FIRMWARE_STATUS_SCAN_SUCCESS      = 9 // 扫描成功
)

const (
	FIRMWARE_TEMPLATE_ID_100 = 100 // 通用IoT固件检测模板
	FIRMWARE_TEMPLATE_ID_101 = 101 // fs
	FIRMWARE_TEMPLATE_ID_102 = 102 // APK扫描
	FIRMWARE_TEMPLATE_ID_103 = 103 // 固件检测_附带二进制安全检测
	FIRMWARE_TEMPLATE_ID_104 = 104 // 二进制单文件扫描new
)

const (
	DEVICE_TYPE_GW  = 1 // 汽车网关(GW)
	DEVICE_TYPE_ECU = 2 // 远程通信单元(ECU)
	DEVICE_TYPE_IVI = 3 // 信息娱乐单元(IVI)

	DEVICE_TYPE_GW_NAME  = "汽车网关(GW)"
	DEVICE_TYPE_ECU_NAME = "远程通信单元(ECU)"
	DEVICE_TYPE_IVI_NAME = "信息娱乐单元(IVI)"
)

const (
	REPORT_NOT_START = 0 // 报告任务未开始
	REPORT_START     = 1 // 报告任务开始
	REPORT_SUCCESS   = 2 // 报告任务结束
)

const (
	LOG_LEVEL_INFO    = iota + 1 // 操作日志:信息
	LOG_LEVEL_WARNING            // 操作日志:警告
	LOG_LEVEL_ERROR              // 操作日志:错误
)

const (
	SCREEN_TESTING = 1 // 大屏状态 测试中
	SCREEN_FINISH  = 2 // 大屏状态，已完成
)

const (
	TASK_CONNECT_STATUS_CONNECTIED  = 1 // 任务是否连接设备:是
	TASK_CONNECT_STATUS_UNCONNECTED = 2 // 任务是否连接设备:否
)

const (
	KNOWLEDGE_TEST_CASE_LEVEL_LOW    = 1 // 测试用例测试难度:低
	KNOWLEDGE_TEST_CASE_LEVEL_MIDDLE = 2 // 测试用例测试难度:中
	KNOWLEDGE_TEST_CASE_LEVEL_HIGH   = 3 // 测试用例测试难度:高
)

const (
	KNOWLEDGE_TEST_CASE_LEVEL_BASIC    = 1 // 测试用例级别:基础测试
	KNOWLEDGE_TEST_CASE_LEVEL_COMPLETE = 2 // 测试用例级别:全面测试
	KNOWLEDGE_TEST_CASE_LEVEL_IMPROVE  = 3 // 测试用例级别:提高测试
	KNOWLEDGE_TEST_CASE_LEVEL_EXPERT   = 4 // 测试用例级别:专家模式
)
const (
	KNOWLEDGE_TEST_CASE_METHOD_BLACK = 1 // 测试方式:黑盒
	KNOWLEDGE_TEST_CASE_METHOD_GRAY  = 2 // 测试方式:灰盒
	KNOWLEDGE_TEST_CASE_METHOD_WHITE = 3 // 测试方式:白盒
)
