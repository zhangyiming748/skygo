package mysql_model

type FirmwareReportRtsLinux struct {
	Id         int    `xorm:"not null pk autoincr INT(10)"`
	ScannerId  int    `xorm:"comment('固件扫描任务id') INT(10)"`
	Type       string `xorm:"VARCHAR(255)"`
	FullPath   string `xorm:"VARCHAR(1024)"`
	Detail     string `xorm:"TEXT"`
	CreateTime int    `xorm:"INT(10)"`
}
