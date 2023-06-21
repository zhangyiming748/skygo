package mysql_model

type FirmwareReportElf struct {
	Id           int `xorm:"not null pk autoincr INT(10)"`
	ScannerId    int `xorm:"comment('固件扫描任务id') INT(10)"`
	Executable   int `xorm:"not null INT(11)"`
	SharedLib    int `xorm:"not null INT(11)"`
	KernelModule int `xorm:"not null INT(11)"`
	CreateTime   int `xorm:"INT(10)"`
}
