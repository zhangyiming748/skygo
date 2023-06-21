package mysql_model

import (
	"errors"
	"fmt"
	"skygo_detection/lib/common_lib/mysql"
	"time"
)

type GpsTask struct {
	Id         int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	TaskId     int    `xorm:"not null comment('任务id') INT(11)" json:"task_id"`
	Status     int    `xorm:"not null comment('状态') INT(10)" json:"status"`
	Time       string `xorm:"not null comment('时间') VARCHAR(20)" json:"time"`
	CreateTime string `xorm:"not null comment('创建时间') varchar(20)" json:"create_time"`
	UpdateTime string `xorm:"not null comment('修改时间') varchar(20)" json:"update_time"`
}

func (g *GpsTask) Create() (int64, error) {
	return mysql.GetSession().InsertOne(g)
}

func (g *GpsTask) Update(cols ...string) (int64, error) {
	return mysql.GetSession().Table(g).ID(g.Id).Cols(cols...).Update(g)
}

func (g *GpsTask) FindGpsByTaskId(taskId int) (bool, error) {
	return mysql.GetSession().Table(g).Where("task_id=?", taskId).Get(g)
}

func (g *GpsTask) UpdateStatusByTaskId(taskId, status int) (bool, error) {
	has, err := mysql.GetSession().Where("task_id = ?", taskId).Get(g)
	if err != nil {
		return has, err
	}
	if !has {
		return has, errors.New("not found gps task")
	}
	g.Status = status
	g.UpdateTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))

	_, err = g.Update("status", "update_time")
	if err != nil {
		return has, err
	}
	return true, nil
}
