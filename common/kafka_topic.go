package common

const (
	// KafkaTopicBasicAsset 关联资产信息topic
	KafkaTopicBasicAsset = "basic_asset"

	// KafkaTopicReportFileList 上报国家数据文件队列
	KafkaTopicReportFileList = "report_file_list_{hostname}"

	// 安全隐患-T：hd_loophole
	TopicHdLoophope = "hd_loophole"

	// 主机受控-T：hd_loophole
	TopicSeHostControlled = "se_host_controlled"

	// 数据泄露事件
	TopicSeDataBreach = "se_data_breach"

	// 07-事件篡改：se_information_tampering
	TopicInformationTampering = "se_information_tampering"

	// 09-网络攻击事件-T:se_cyber_attacks
	TopicSeCyberAttacks = "se_cyber_attacks"

	// 10-有害程序-T:se_armful_program
	TopicSeArmfulProgram = "se_armful_program"

	// 11-高级威胁-T:se_advanced_threats
	TopicSeAdvancedThreats = "se_advanced_threats"

	// 12-异常违规-T:se_bnormal_violations
	TopicSeBnormalViolations = "se_bnormal_violations"

	KafkaTopicBasicCompany = "basic_company"

	KafkaTopicBasicPlatform = "basic_platform"
	// 流量session日志
	KafkaTopicFlowSession = "flow_session"
)
