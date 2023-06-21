package common

func GetFieldValue(m map[int]string, i int) string {
	if v, has := m[i]; has {
		return v
	}
	return ""
}

const StoreTypeEs = "es"
const StoreTypeKafka = "kafka"

// es默认日期的格式
const DefaultTimeFormat = "strict_date_optional_time||epoch_millis"

// 属性分类， 1 预定义 2 自定义
// 所有列表通用的字段，都有属性分类
const (
	CategoryPredefined = 1
	CategoryCustomized = 2
)

var CategoryMap = map[int]string{ // 前端展示
	CategoryPredefined: "预定义",
	CategoryCustomized: "自定义",
}

// es分区类型，1不分区， 2按天分区
const (
	EsPartitionModeNo  = 1
	EsPartitionModeDay = 2
)

var EsPartitionModeMap = map[int]string{ // 前端展示
	EsPartitionModeNo:  "不分区",
	EsPartitionModeDay: "按天分区",
}

// 对象配置 - kafka集群基础配置
// kafka认证类型, ，0无 1kerberos
const (
	KafkaAuthTypeNo       = 1
	KafkaAuthTypeKerberos = 2
)

var KafkaAuthTypeMap = map[int]string{ // 前端展示
	KafkaAuthTypeNo:       "无",
	KafkaAuthTypeKerberos: "kerberos",
}
var KafkaAuthTypeMap2 = map[int]string{ // 给etl_engine提供时使用的映射
	KafkaAuthTypeNo:       "",
	KafkaAuthTypeKerberos: "kerberos",
}

// 对象配置 - es集群基础配置
// es认证方式，1用户密码
const (
	EsAuthUserPass = 1
)

var EsAuthTypeMap = map[int]string{ // 前端展示
	EsAuthUserPass: "账户密码",
}

const (
	TemplateIfKafkaMemoryYes = 1 // 数据模板，支持kafka存储使用
	TemplateIfKafkaMemoryNo  = 0 // 数据模板，不支持kafka存储使用

	TemplateIfEsMemoryYes = 1 // 数据模板，支持es存储使用
	TemplateIfEsMemoryNo  = 0 // 数据模板，不支持es存储使用

	TemplateEsSearchableYes = 1 // 数据模板，支持做为【检索日志】
	TemplateEsSearchableNo  = 0 // 数据模板，不支持做为【检索日志】
)

// BasicLogGroupOriginNodeId 全部节点的id
const BasicLogGroupOriginNodeId = 0

const BasicLogGroupWithoutName = "未分组"
const BasicLogGroupWithoutId = -1 // “未分组”的id

// 任务状态
// 状态，1草稿，2未运行 3运行中 4异常
const (
	TaskStatusDraft    = 1
	TaskStatusNotRun   = 2
	TaskStatusRunning  = 3
	TaskStatusAbnormal = 4
)

var TaskStatusMap = map[int]string{ // 前端展示
	TaskStatusDraft:    "草稿",
	TaskStatusNotRun:   "未运行",
	TaskStatusRunning:  "运行中",
	TaskStatusAbnormal: "异常",
}

// 规则
// 配置方式，1界面，2高级指令
const (
	RuleConfigTypeWindow = 1
	RuleConfigTypeCode   = 2
)

var RuleConfigTypeMap = map[int]string{ // 前端展示
	RuleConfigTypeWindow: "界面",
	RuleConfigTypeCode:   "高级",
}

// 规则状态，1草稿  2正常
const (
	RuleStatueDraft  = 1
	RuleStatusNormal = 2
)

var RuleStatusMap = map[int]string{
	RuleStatueDraft:  "草稿",
	RuleStatusNormal: "正常",
}

// 界面配置，解析类型 1json 2 grok
const (
	RuleWdParsingTypeJson    = 1 // json
	RuleWdParsingTypeGrok    = 2 // grok
	RuleWdParsingTypeSplit   = 3 // 分隔符，逗号分割
	RuleWdParsingTypeNoParse = 4 // 不解析
)

var RuleWdParsingTypeMap = map[int]string{
	RuleWdParsingTypeJson: "json",
	// RuleWdParsingTypeGrok : "grok",
	// RuleWdParsingTypeSplit : "split",
}

// 输出数据类型
const (
	OutputDataFormatJson = 1
)

var OutputDataFormatMap = map[int]string{
	OutputDataFormatJson: "json",
}

// 字段类型
const (
	FieldTypeKeyWord = iota + 1
	FieldTypeText
	FieldBoolean
	FieldTypeInteger
	FieldTypeLong
	FieldTypeFloat
	FieldTypeDouble
	FieldTypeDate
	FieldTypeIp
	FieldTypeStruct
)

var FieldTypeMap = map[int]string{
	FieldTypeKeyWord: "字符串",
	FieldTypeText:    "文本",
	FieldBoolean:     "布尔",
	FieldTypeInteger: "整型",
	FieldTypeLong:    "长整型",
	FieldTypeFloat:   "浮点数4字节",
	FieldTypeDouble:  "浮点数8字节",
	FieldTypeDate:    "日期",
	FieldTypeIp:      "IP",
	FieldTypeStruct:  "结构体",
}

var FieldTypeMap2 = map[int]string{
	FieldTypeKeyWord: "string",
	FieldTypeText:    "string",
	FieldBoolean:     "bool",
	FieldTypeInteger: "int4",
	FieldTypeLong:    "int8",
	FieldTypeFloat:   "float4",
	FieldTypeDouble:  "float8",
	FieldTypeDate:    "date",
	FieldTypeIp:      "ip",
	FieldTypeStruct:  "struct",
}

// field字段类型和es中字段类型的对应关系
var FieldTypeEsMap = map[int]string{
	FieldTypeKeyWord: "keyword",
	FieldTypeText:    "text",
	FieldBoolean:     "boolean",
	FieldTypeInteger: "integer",
	FieldTypeLong:    "long",
	FieldTypeFloat:   "float",
	FieldTypeDouble:  "double",
	FieldTypeDate:    "date",
	FieldTypeIp:      "ip",
	// 注意，没有struct
}
