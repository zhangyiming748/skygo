package mysql_model

type FirmwareReportRtsApkVul struct {
	Id          int    `xorm:"not null pk autoincr INT(10)"`
	ScannerId   int    `xorm:"comment('固件扫描任务id') INT(10)"`
	Type        string `xorm:"VARCHAR(255)"`
	VulDesc     string `xorm:"not null VARCHAR(512)"`
	Reference   string `xorm:"VARCHAR(1024)"`
	ParentType  string `xorm:"VARCHAR(255)"`
	Objective   string `xorm:"VARCHAR(1024)"`
	RiskLevel   string `xorm:"VARCHAR(255)"`
	CheckResult string `xorm:"VARCHAR(255)"`
	Solution    string `xorm:"VARCHAR(1023)"`
	Name        string `xorm:"VARCHAR(255)"`
	VulEffect   string `xorm:"VARCHAR(1023)"`
	PkgName     string `xorm:"VARCHAR(255)"`
	Detail      string `xorm:"TEXT"`
	CreateTime  int    `xorm:"INT(10)"`
}
