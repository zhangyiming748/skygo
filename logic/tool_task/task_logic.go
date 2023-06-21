package tool_task

import (
	"errors"
	"skygo_detection/mysql_model"
	"time"
)

const (
	TASK_DEAFAULT = 1 // 待测试
	TASK_TESTING  = 2 // 测试中
	TASK_COMPLETE = 3 // 完成
	TASK_FAIL     = 4 // 失败
	TASK_PAUSE    = 5 // 暂停---测过，但还没手动点完成
)

type Task struct{}

// 创建任务
func (t Task) Create(userId int, taskUuid string, subTaskId int, subDetailUrl, name,
	toolType, describe, tool, category string,
) (int64, error) {
	task := new(mysql_model.Task)
	task.SubTaskId = subTaskId
	task.SubDetailUrl = subDetailUrl
	task.Name = name
	task.Status = TASK_DEAFAULT
	task.ToolType = toolType
	task.Describe = describe
	task.Category = category
	task.Tool = tool
	task.TaskUuid = taskUuid
	//task.TaskUuid = mysql.GetTaskId()
	task.CreateUserId = userId
	task.CreateTime = int(time.Now().Unix())
	if _, err := task.Create(); err != nil {
		return 0, err
	}
	return int64(task.Id), nil
}

// 更新状态
func (t Task) UpdateStatus(subTaskId int, toolType string, taskUuid string, status int) (res int64, err error) {
	taskList, err := mysql_model.GetTaskBySubInfo(subTaskId, toolType, taskUuid)
	if err != nil {
		return
	}
	if taskList.Id <= 0 {
		return 0, errors.New("任务不存在")
	}
	nowTime := int(time.Now().Unix())
	data := map[string]interface{}{
		"status":      status,
		"update_time": nowTime,
	}
	// 完成状态，更新完成时间
	if status == TASK_COMPLETE {
		data = map[string]interface{}{
			"status":        status,
			"update_time":   nowTime,
			"complete_time": nowTime,
		}
	}
	res, err = mysql_model.UpdateStatusBySubInfo(taskList.Id, data)
	if err != nil {
		return
	}
	return
}

// 删除任务
func (t Task) Delete(subTaskId int, toolType string, taskUuid string) (res int64, err error) {
	taskList, err := mysql_model.GetTaskBySubInfo(subTaskId, toolType, taskUuid)
	if err != nil {
		return
	}
	if taskList.Id <= 0 {
		return 0, errors.New("任务不存在")
	}
	res, err = mysql_model.DeleteTaskByTaskId(taskList.Id)
	if err != nil {
		return
	}
	return
}
