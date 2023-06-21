package mysql_model

import (
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
)

type BeehiveLog struct {
	Id         int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	TaskId     int    `xorm:"not null comment('任务id') INT(10)" json:"task_id"`
	Title      string `xorm:"comment('操作说明') VARCHAR(100)" json:"title"`
	CreateTime string `xorm:"created not null comment('创建时间') DATETIME" json:"create_time"`
	Content    string `xorm:"comment('操作结果') VARCHAR(100)" json:"content"`
}

func (this BeehiveLog) SetLog() error {
	_, err := mysql.GetSession().InsertOne(this)
	if err != nil {
		return err
	}
	return nil
}

func GetLog(tid int) ([]BeehiveLog, error) {
	logs := make([]BeehiveLog, 0)
	s := mysql.GetSession().Table("beehive_log")
	s.Where("task_id = ?", tid)
	s.Desc("create_time")
	err := s.Find(&logs)
	if err != nil {
		return nil, err
	}
	return logs, nil
}
func GetBeehiveLog(tid int) map[string]interface{} {
	all := mysql.GetSession().
		Table("beehive_log").
		Where("task_id = ?", tid)

	widget := orm.PWidget{}
	widget.AddSorter(*(orm.NewSorter("create_time", 1)))
	res := widget.PaginatorFind(all, &[]BeehiveLog{})
	return res
}
