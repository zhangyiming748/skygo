package mysql_model

import "skygo_detection/lib/common_lib/mysql"

type BeehiveGsmSniffer struct {
	Id         int    `xorm:"not null pk comment('自增主键id') INT(11)" json:"id"`
	TaskId     int    `xorm:"not null comment('任务id') INT(11)" json:"task_id"`
	Status     int    `xorm:"not null comment('1:频点扫描 2:嗅探imsi 3:嗅探短信') INT(11)" json:"status"`
	Channel    int    `xorm:"not null comment('最近一次扫描所使用的频段 1:900M 2:1800M') INT(11)" json:"channel"`
	ScanTime   string `xorm:"not null comment('最近一次扫描的时间') varchar(20)" json:"scan_time"`
	Frequency  string `xorm:"not null comment('最近一次获取到的频点，多个时以逗号连起来') varchar(1024)" json:"frequency"`
	SniffTime  string `xorm:"not null comment('最近一次嗅探的时间') varchar(20)" json:"sniff_time"`
	SniffFreq  string `xorm:"not null comment('最近一次嗅探的频点') varchar(20)" json:"sniff_freq"`
	CreateTime string `xorm:"not null comment('创建时间') varchar(20)" json:"create_time"`
	UpdateTime string `xorm:"not null comment('修改时间') varchar(20)" json:"update_time"`
}

func (b *BeehiveGsmSniffer) Create() (int64, error) {
	return mysql.GetSession().InsertOne(b)
}

func (b *BeehiveGsmSniffer) Update(cols ...string) (int64, error) {
	return mysql.GetSession().Table(b).ID(b.Id).Cols(cols...).Update(b)
}

func (b *BeehiveGsmSniffer) FindByTaskId(taskId int) (bool, error) {
	return mysql.GetSession().Table(b).Where("task_id=?", taskId).Get(b)
}
