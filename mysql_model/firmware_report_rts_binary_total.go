package mysql_model

type FirmwareReportRtsBinaryTotal struct {
	Id         int `xorm:"not null pk autoincr INT(11)"`
	ScannerId  int `xorm:"index INT(11)"`
	Nx         int `xorm:"INT(11)"`
	Pie        int `xorm:"INT(11)"`
	Relro      int `xorm:"INT(11)"`
	Canary     int `xorm:"INT(11)"`
	Stripped   int `xorm:"INT(11)"`
	Count      int `xorm:"INT(11)"`
	CreateTime int `xorm:"INT(11)"`
}
