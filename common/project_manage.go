package common

import (
	"skygo_detection/guardian/src/net/qmap"
	"strconv"
)

// 通用
const (
	ENABLED  = 1 //激活
	DISABLED = 0 //禁用
)

// 项目角色
const (
	ROLE_PM   = 10000 //项目经理
	ROLE_TEST = 11000 //测试人员
)

// 项目状态(Project Status)
const (
	PS_ABNORMAL = 100 //项目状态:异常
	PS_NEW      = 0   //新建
	PS_TEST     = 1   //测试中
	PS_COMPLETE = 9   //项目完成
)

// 项目软删除，独立字段标志
const (
	PSD_DEFAULT = 0 //项目删除：默认未删除
	PSD_DELETE  = 1 //项目删除：已删除
)

// 潜力项目状态(Potential Project Status)
const (
	PPS_DELETE     = -1 //潜力项目状态:删除
	PPS_UNAPPROVED = 0  //潜力项目状态:未立项
	PPS_APPROVAL   = 1  //潜力项目状态:已立项
)

// 潜力项目类型(Potential Project Type)
const (
	PPT_NORMAL   = 1 //潜力项目类型:普通项目
	PPT_STRATEGY = 2 //潜力项目类型:战略项目
)

// 测试项状态(Evaluate Item Status)
const (
	EIS_READY_PRELIMINARY = 0 //测试项状态:待初测
	EIS_PRELIMINARY_BEGIN = 1 //测试项状态:初测完成
	EIS_RETEST_BEGIN      = 2 //测试项状态:待复测
	EIS_RETEST_END        = 3 //测试项状态:复测完成
)

// 评估漏洞状态(Evaluate Vul Status)
const (
	EVS_UNREPAIRED = 0 //评估漏洞状态:未修复
	EVS_REPAIRED   = 1 //评估漏洞状态:已修复
	EVS_REOPEN     = 2 //评估漏洞状态:重打开
)

var EvaluateVulMap = qmap.QM{
	string(EVS_UNREPAIRED): "未修复",
	string(EVS_REPAIRED):   "已修复",
	string(EVS_REOPEN):     "重打开",
}

// 测试项难度(Evaluate Target Level)
const (
	ETL_LOW    = 1 //测试项难度:低
	ETL_MIDDLE = 2 //测试项难度:中
	ETL_HIGH   = 3 //测试下难度:高
)

var EvaluateTargetLevelMap = qmap.QM{
	strconv.Itoa(ETL_LOW):    "低",
	strconv.Itoa(ETL_MIDDLE): "中",
	strconv.Itoa(ETL_HIGH):   "高",
}

// 资产状态(Asset Status)
const (
	AS_USING  = 0 //资产状态:使用中
	AS_RETURN = 1 //资产状态:已归还
)

var AssetStatusMap = qmap.QM{
	strconv.Itoa(AS_USING):  "使用中",
	strconv.Itoa(AS_RETURN): "已归还",
}

// 报告类型(Report Type)
const (
	RT_WEEK   = "week"   //报告类型:周报
	RT_TEST   = "test"   //报告类型:初测
	RT_RETEST = "retest" //报告类型:复测
)

var ReportStatusMap = qmap.QM{
	RT_WEEK:   "周报",
	RT_TEST:   "初测报告",
	RT_RETEST: "复测报告",
}

// 报告状态
const (
	RS_NEW     = 0  //创建
	RS_AUDIT   = 1  //审核
	RS_SUCCESS = 2  //发布成功
	RS_FAILED  = -1 //发布失败
)

// 报告审核结果
const (
	RAS_NEW     = 0  //待审核
	RAS_SUCCESS = 1  //审核通过
	RAS_FAILED  = -1 //审核驳回
)

// 漏洞状态(Vul Status)
const (
	VS_UNREPAIR = 0 //未修复
	VS_REPAIRED = 1 //已修复
	VS_REOPEN   = 2 //重打开
)

var VulStatusMap = qmap.QM{
	strconv.Itoa(VS_UNREPAIR): "未修复",
	strconv.Itoa(VS_REPAIRED): "已修复",
	strconv.Itoa(VS_REOPEN):   "重打开",
}

// 漏洞级别(Vul Level)
const (
	VL_INFO    = 0 //提示
	VL_LOW     = 1 //低危
	VL_MIDDLE  = 2 //中危
	VL_HIGH    = 3 //高危
	VL_SERIOUS = 4 //严重
)

var VulLevelMap = qmap.QM{
	strconv.Itoa(VL_INFO):    "提示",
	strconv.Itoa(VL_LOW):     "低危",
	strconv.Itoa(VL_MIDDLE):  "中危",
	strconv.Itoa(VL_HIGH):    "高危",
	strconv.Itoa(VL_SERIOUS): "严重",
}

// 项目中标难度(Bid Acceptance Probability)
const (
	BAP_LOW    = 1 //测试项难度:低
	BAP_MIDDLE = 2 //测试项难度:中
	BAP_HIGH   = 3 //测试下难度:高
)

var BidAcceptanceProbabilityMap = qmap.QM{
	strconv.Itoa(BAP_LOW):    "低",
	strconv.Itoa(BAP_MIDDLE): "中",
	strconv.Itoa(BAP_HIGH):   "高",
}

// 测试用例删除状态
const (
	STATUS_DELETE = 0 //存在
	STATUS_EXIST  = 1 //已删除
)

// 项目任务状态(PTS=project task status )
const (
	PTS_CREATE       = 1 //创建
	PTS_TASK_AUDIT   = 2 //任务审核
	PTS_TEST         = 3 //任务测试
	PTS_REPORT_AUDIT = 4 //报告审核
	PTS_FINISH       = 5 //完成
)

// 项目任务审核状态
const (
	PTS_AUDIT_STATUS_NEW = 0  //未审核
	PTS_AUDIT_STATUS_YES = 1  //审核通过
	PTS_AUDIT_STATUS_NO  = -1 //审核驳回
)

// 任务状态标志
const (
	PTS_SIGN_CREATE       = "create"
	PTS_SIGN_TASK_AUDIT   = "task_audit"
	PTS_SIGN_TEST         = "test"
	PTS_SIGN_REPORT_AUDIT = "report_audit"
	PTS_SIGN_FINISH       = "finish"
)

// 任务用例是否预绑定
const (
	IS_PREBIND  = 1
	NOT_PREBIND = 0
)

// 任务用例是否预删除
const (
	IS_PREDEL  = 1
	NOT_PREDEL = 0
)

// 测试用例测试状态（EIS=evaluate item status）
const (
	EIS_FREE  = 0 //可创建任务
	EIS_INUSE = 1 //使用中，不可创建任务
)

// 任务测试状态（TIS = task item status）
const (
	TIS_READY              = 0 //待测试
	TIS_TEST_COMPLETE      = 1 //测试完成
	TIS_PART_TEST_COMPLETE = 2 //待补充
	TIS_COMPLETE           = 3 //审核通过
)

// 测试用例测试记录审核状态(IRAS = item record audit status)
const (
	IRAS_ACCEPT  = 1  //通过
	IRAS_DEFAULT = 0  //默认
	IRAS_REJECT  = -1 //驳回
)

// 测试用例审核状态（EIAS=evaluate item audit status）
const (
	EIAS_ACCEPT  = 1  //通过
	EIAS_DEFAULT = 0  //默认
	EIAS_REJECT  = -1 //驳回
)

// 漏洞检测状态
const (
	VUL_UNSTART           = 0 //漏洞检测状态:未开始
	VUL_PRELIMINARY_BEGIN = 1 //漏洞检测状态:测试中
	VUL_PRELIMINARY_END   = 2 //漏洞检测状态:测试完成
	VUL_UNSHOW            = 4 //漏洞检测状态:未展示
)

// 漏洞修复状态
const (
	VUL_FIX_REPAIR   = 2 //漏洞修复状态 已修复
	VUL_FIX_UNREPAIR = 1 //漏洞修复状态 未修复
	VUL_FIX_OTHER    = 3 //漏洞修复未涉及
)

// 谷歌漏洞级别
const (
	VUL_GOOGLE_LEVEL_LOW     = 0 //漏洞级别：低危
	VUL_GOOGLE_LEVEL_MIDDLE  = 1 //漏洞级别：中危
	VUL_GOOGLE_LEVEL_HIGHT   = 2 //漏洞级别：高危
	VUL_GOOGLE_LEVEL_SERIOUS = 3 //漏洞级别：严重
)

// 用例 测试难度（0:默认、1:低、2:中、3:高）
const (
	TEST_LEVEL_DEFAULT = 0 //默认
	TEST_LEVEL_LOW     = 1 //低
	TEST_LEVEL_MIDDLE  = 2 //中
	TEST_LEVEL_HIGH    = 3 //高
)

// 测试用例 手动创建 或者 自动创建（A:自动生成、M:手动生成）
const (
	TEST_AUTO   = "A" //自动生成
	TEST_MANUAL = "M" //手动生成
)

// 扫描任务状态
const (
	SCANNER_STATUS_READY   = 0 // 待执行
	SCANNER_STATUS_HANDING = 1 // 执行中
	SCANNER_STATUS_SUCCESS = 2 // 执行成功
	SCANNER_STATUS_FAILURE = 3 // 执行失败
)

// 报告类型
const (
	REPORT_TEST        = "test"         //项目初测报告
	REPORT_RETEST      = "retest"       //项目复测报告
	REPORT_ASSETTEST   = "asset_test"   //资产初测报告
	REPORT_ASSETRETEST = "asset_retest" //资产复测报告
	REPORT_ASSETPRE    = "asset"        //资产报告前缀
)
