package common

const (
	// RedisReportResourceLockKey 上报资源的锁
	RedisReportResourceLockKey = "reporting_res_lock_{resource}"

	// RedisReportResourceLastTimeKey 上报资源的最后上报时间
	RedisReportResourceLastTimeKey = "reporting_res_last_time_{resource}"

	// RedisEtlIPAreaMappingKey ip对ipmin_ipmax的对应关系
	RedisEtlIPAreaMappingKey = "etl_ip_{ip}"
	// RedisEtlIPAreaDataKey ipmin_ipmax对应的区域信息
	RedisEtlIPAreaDataKey = "etl_iparea_{ipmin_ipmax}"
)
