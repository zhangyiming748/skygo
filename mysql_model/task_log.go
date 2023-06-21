package mysql_model

import (
	"time"

	"skygo_detection/guardian/app/sys_service"

	"skygo_detection/common"
)

type TaskLog struct {
	Id                int    `xorm:"not null pk autoincr INT(11)"`
	TaskId            int    `xorm:"comment('测试任务id') INT(11)"`
	TaskName          string `xorm:"comment('任务名称') VARCHAR(255)"`
	Tool              string `xorm:"comment('工具') VARCHAR(255)"`
	Level             int    `xorm:"comment('日志级别(1:info, 2:warning 3:error)') TINYINT(11)"`
	Status            int    `xorm:"comment('任务状态') TINYINT(11)"`
	CurrentTaskStatus int    `xorm:"comment('当前任务状态(1:存在，-1:删除)') TINYINT(11)"`
	Message           string `xorm:"comment('描述') VARCHAR(255)"`
	UserId            int    `xorm:"comment('用户id') INT(11)"`
	UserName          string `xorm:"comment('用户名称') VARCHAR(255)"`
	CreateTime        int    `xorm:"comment('创建时间') INT(11)"`
}

func (this *TaskLog) Create() (int64, error) {
	return sys_service.NewSession().Session.InsertOne(this)
}

func (this *TaskLog) Update(cols ...string) (int64, error) {
	return sys_service.NewSession().Session.Table(this).ID(this.Id).Cols(cols...).Update(this)
}

func (this *TaskLog) Remove() (int64, error) {
	return sys_service.NewSession().Session.ID(this.Id).Delete(this)
}

func (this *TaskLog) Insert(userId int, username string, task *Task) (int64, error) {
	this.TaskId = task.Id
	this.TaskName = task.Name
	this.Tool = task.Tool
	this.Status = task.Status
	this.Level = common.LOG_LEVEL_INFO
	this.Message = ""
	this.UserId = userId
	this.UserName = username
	this.CreateTime = int(time.Now().Unix())
	num, err := this.Create()
	if err == nil {
		if this.Status == common.TASK_STATUS_REMOVE {
			new(TaskLog).UpdateTaskLogStatusToRemove(task.Id)
		}
	}
	return num, err
}

const (
	TASK_LOG_STATUS_EXIST  = 1
	TASK_LOG_STATUS_DELETE = -1
)

// 标记任务状态为已删除(1:存在，-1:删除)
func (this *TaskLog) UpdateTaskLogStatusToRemove(taskId int) {
	sql := "update task_log set current_task_status=? where task_id=?"
	if _, err := sys_service.NewSession().Session.Exec(sql, TASK_LOG_STATUS_DELETE, taskId); err != nil {
		panic(err)
	}
}
