package mysql_model

type FirmwareReportRtsLinuxTotal struct {
	Id            int `xorm:"not null pk autoincr INT(10)"`
	ScannerId     int `xorm:"comment('固件扫描任务id') index INT(10)"`
	AbnormalCount int `xorm:"not null default 0 INT(10)"`
	CreateTime    int `xorm:"INT(10)"`
}
