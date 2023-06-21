package mysql_model

import (
	"fmt"
	"skygo_detection/lib/common_lib/mysql"
	"time"
)

type BeehiveGsmSystem struct {
	Id           int    `xorm:"not null pk autoincr comment('主键id') INT(10)" json:"id"`
	TaskId       int    `xorm:"not null comment('任务id') INT(10)" json:"task_id"`
	SystemStatus int    `xorm:"not null comment('1运行中2暂停') TINYINT(10)" json:"system_status"`
	ConfigId     int    `xorm:"not null comment('系统控制的配置参数') INT(1)" json:"config_id"`
	CreateTime   string `xorm:"not null comment('创建时间') DATETIME" json:"create_time"`
	UpdateTime   string `xorm:"updated not null comment('更新时间') DATETIME" json:"update_time"`
}

func CreateGsmTask(tid int) (int64, error) {
	this := new(BeehiveGsmSystem)
	this.TaskId = tid
	this.CreateTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
	one, err := mysql.GetSession().InsertOne(this)
	if err != nil {
		return 0, err
	}
	return one, err
}
func (this BeehiveGsmSystem) Update() (int64, error) {
	s := mysql.GetSession()
	i, err := s.Where("task_id = ?", this.TaskId).Update(this)
	if err != nil {
		return 0, err
	}
	return i, nil
}
func (this BeehiveGsmSystem) ForceUpdateConfigId() (int64, error) {
	s := mysql.GetSession()
	i, err := s.Where("task_id = ?", this.TaskId).
		Cols("config_id").Update(this)
	if err != nil {
		return 0, err
	}
	return i, nil
}
func (this BeehiveGsmSystem) ForceUpdateSystemStatus() (int64, error) {
	s := mysql.GetSession()
	i, err := s.Where("task_id = ?", this.TaskId).
		Cols("system_status").Update(this)
	if err != nil {
		return 0, err
	}
	return i, nil
}

func (this BeehiveGsmSystem) Stop() (int64, error) {
	s := mysql.GetSession()
	s.Where("task_id =?", this.TaskId)
	return s.Update(this)
}
func GsmTaskNotExist(tid int) bool {
	has, _ := mysql.GetSession().Where("task_id = ?", tid).Get(new(BeehiveGsmSystem))
	if has {
		return false
	}
	return true
}
func (this BeehiveGsmSystem) GetOne(tid int) (BeehiveGsmSystem, error) {
	_, err := mysql.GetSession().Where("task_id =?", tid).Get(&this)
	if err != nil {
		return this, err
	}
	return this, err

}
