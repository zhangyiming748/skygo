package mysql_model

import (
	"skygo_detection/lib/common_lib/mysql"
)

type PrivacyAnalysisLog struct {
	Id         int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	TaskId     int    `xorm:"not null comment('任务id') INT(10)" json:"task_id"`
	CreateTime string `xorm:"created not null comment('创建时间') DATETIME" json:"create_time"`
	Content    string `xorm:"comment('日志内容') VARCHAR(512)" json:"content"`
}

func (this PrivacyAnalysisLog) SetPrivacyLog(tid int, content string) error {
	this.TaskId = tid
	this.Content = content
	_, err := mysql.GetSession().InsertOne(this)
	if err != nil {
		return err
	}
	return nil
}

func GetPrivacyLog(tid int) []PrivacyAnalysisLog {
	logs := make([]PrivacyAnalysisLog, 0)
	s := mysql.GetSession()
	s.Where("task_id = ?", tid)
	s.Desc("create_time")
	s.Find(&logs)
	return logs
}
