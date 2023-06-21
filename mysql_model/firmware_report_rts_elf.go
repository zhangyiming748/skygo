package mysql_model

type FirmwareReportRtsElf struct {
	Id         int    `xorm:"not null pk autoincr INT(10)"`
	ScannerId  int    `xorm:"comment('固件扫描任务id') INT(10)"`
	Executable string `xorm:"VARCHAR(1024)"`
	Type       string `xorm:"VARCHAR(255)"`
	CreateTime int    `xorm:"INT(10)"`
}
