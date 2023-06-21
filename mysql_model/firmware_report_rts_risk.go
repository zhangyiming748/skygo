package mysql_model

type FirmwareReportRtsRisk struct {
	Id               int `xorm:"not null pk autoincr INT(10)"`
	ScannerId        int `xorm:"comment('固件扫描任务id') INT(10)"`
	RiskCount        int `xorm:"not null INT(10)"`
	PassRiskCount    int `xorm:"not null INT(10)"`
	BinaryCount      int `xorm:"not null INT(10)"`
	LinuxCount       int `xorm:"not null INT(10)"`
	OverCertCount    int `xorm:"not null INT(10)"`
	RiskSuspectCount int `xorm:"not null INT(10)"`
}
