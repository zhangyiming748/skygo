package mysql_model

type FirmwareReportRtsApkSensitive struct {
	Id         int    `xorm:"not null pk autoincr INT(10)"`
	ScannerId  int    `xorm:"comment('固件扫描任务id') INT(10)"`
	Type       string `xorm:"VARCHAR(255)"`
	Content    string `xorm:"VARCHAR(1024)"`
	PkgName    string `xorm:"VARCHAR(1024)"`
	CreateTime int    `xorm:"INT(10)"`
}
