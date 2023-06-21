package gps

import (
	"errors"
	"fmt"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/mysql_model"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	TASK_DEAFAULT = iota // 待测试
	TASK_TESTING         // 测试中
	TASK_COMPLETE        // 完成
	TASK_FAIL            // 失败
	TASK_PAUSE           // 暂停---测过，但还没手动点完成
)

type Task struct{}

func (t Task) Create(ctx *gin.Context, name, describe string) (int64, error) {
	task := new(mysql_model.Task)
	task.Name = name
	task.ToolType = "gps"
	task.Describe = describe
	task.Category = "GPS测试工具"
	task.Tool = "gps"
	task.TaskUuid = mysql.GetTaskId()
	task.CreateUserId = int(http_ctx.GetUserId(ctx))
	task.CreateTime = int(time.Now().Unix())
	if _, err := task.Create(); err != nil {
		return 0, err
	}
	return int64(task.Id), nil
}

// 编辑任务基础信息
func (t Task) Update(ctx *gin.Context, id int, name, describe string) (int64, error) {
	task, bool := mysql_model.TaskFindById(id)
	if !bool {
		return 0, errors.New("没找到记录")
	}
	if task.Status > 0 {
		return 0, errors.New("该任务已执行过，不能修改了")
	}
	task.Name = name
	task.Describe = describe
	return task.Update()
}

// 获取任务基础信息
func (t Task) GetOne(ctx *gin.Context, id int) map[string]interface{} {
	info := make(map[string]interface{})
	task, bool := mysql_model.TaskFindById(id)
	if !bool {
		return nil
	}
	info["name"] = task.Name
	info["describe"] = task.Describe
	info["user"] = ""
	if userModel, err := mysql_model.SysUserFindById(task.CreateUserId); err == nil {
		info["user"] = userModel.Realname
	}

	info["create_time"] = fmt.Sprint(time.Unix(int64(task.CreateTime), 0).Format("2006-01-02 15:04:05"))
	info["task_uuid"] = task.TaskUuid
	info["status"] = task.Status
	info["complete_time"] = ""
	if task.CompleteTime > 0 {
		info["complete_time"] = fmt.Sprint(time.Unix(int64(task.CompleteTime), 0).Format("2006-01-02 15:04:05"))
	}
	return info
}

// 完成任务
func (t Task) Complete(ctx *gin.Context, id int) (int64, error) {
	task, bool := mysql_model.TaskFindById(id)
	if !bool {
		return 0, errors.New("没找到记录")
	}
	gc, err := GetLatestCheatByTaskId(id)
	if err != nil {
		return 0, err
	}
	if gc.Status == CHEAT_STATUS_ING {
		return 0, errors.New("正在欺骗中，请先停止欺骗，并且关闭系统")
	}
	gt, err := GetGpsTaskByTaskId(id)
	if err != nil {
		return 0, err
	}
	if gt.Status == DEVICE_RUNNING {
		return 0, errors.New("请先关闭系统")
	}
	task.Status = TASK_COMPLETE
	task.CompleteTime = int(time.Now().Unix())
	return task.Update()
}

// 开始任务
func (t Task) Start(id int) (int64, error) {
	task, bool := mysql_model.TaskFindById(id)
	if !bool {
		return 0, errors.New("没找到记录")
	}
	task.Status = TASK_TESTING
	return task.Update()
}

// 关闭系统
func (t Task) Close(id int) (int64, error) {
	task, bool := mysql_model.TaskFindById(id)
	if !bool {
		return 0, errors.New("没找到记录")
	}
	task.Status = TASK_PAUSE
	return task.Update()
}
