package mysql_model

import "skygo_detection/lib/common_lib/mysql"

type BeehiveGsmSnifferImsi struct {
	Id         int    `xorm:"not null pk comment('自增主键id') INT(11)" json:"id"`
	TaskId     int    `xorm:"not null comment('任务id') INT(11)" json:"task_id"`
	Date       string `xorm:"not null comment('创建时间') VARCHAR(20)" json:"date"`
	Content    string `xorm:"not null comment('内容') VARCHAR(2048)" json:"content"`
	CreateTime string `xorm:"not null comment('创建时间') VARCHAR(20)" json:"create_time"`
}

func (b *BeehiveGsmSnifferImsi) Create() (int64, error) {
	return mysql.GetSession().InsertOne(b)
}

func (b *BeehiveGsmSnifferImsi) Update(cols ...string) (int64, error) {
	return mysql.GetSession().Table(b).ID(b.Id).Cols(cols...).Update(b)
}
