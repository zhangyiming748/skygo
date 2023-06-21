package mysql_model

type FirmwareReportRtsCveTotal struct {
	Id         int `xorm:"not null pk autoincr INT(11)"`
	ScannerId  int `xorm:"index INT(11)"`
	No         int `xorm:"INT(11)"`
	Low        int `xorm:"INT(11)"`
	Mid        int `xorm:"INT(11)"`
	High       int `xorm:"INT(11)"`
	Heavy      int `xorm:"INT(11)"`
	CreateTime int `xorm:"INT(11)"`
}
