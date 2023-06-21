package mysql_model

type HgTestTaskLog struct {
	Id         int `xorm:"not null pk autoincr INT(10)"`
	Status     int `xorm:"TINYINT(10)"`
	CreateTime int `xorm:"default 0 INT(10)"`
}
