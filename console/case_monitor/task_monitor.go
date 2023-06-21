package case_monitor

import (
	"sync"

	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/mysql_model"
)

var lastMaxTaskCreateTime int64 // 最大任务创建时间
var currentTaskList = []*mysql_model.Task{}
var taskLock sync.Mutex

// 查询获取最大任务创建时间
func GetMaxTaskCreateTime() int64 {
	var result struct {
		MaxCreateTime int64
	}
	if has, _ := sys_service.NewSession().Table(new(mysql_model.Task)).Select("max(create_time) as max_create_time").GetOne(&result); has {
		return result.MaxCreateTime
	}
	return 0
}

// 更新场景任务列表
func updateTaskList() {
	currentTaskCreateTime := GetMaxTaskCreateTime()
	if lastMaxTaskCreateTime < currentTaskCreateTime {
		// 如果当前最大任务创建时间大于上次最大任务创建时间，说明有新创建的任务
		taskList := make([]*mysql_model.Task, 0)
		if err := sys_service.NewSession().Session.Table(new(mysql_model.Task)).Where("create_time>?", lastMaxTaskCreateTime).And("create_time<=?", currentTaskCreateTime).And("status=?", common.TASK_STATUS_RUNNING).Find(&taskList); err != nil {
			panic(err)
		} else {
			taskLock.Lock()
			currentTaskList = append(currentTaskList, taskList...)
			taskLock.Unlock()
		}
		lastMaxTaskCreateTime = currentTaskCreateTime
	}
}

// 尝试移除场景任务
func tryRemoveTask(taskId int) {
	params := qmap.QM{
		"e_id": taskId,
	}
	task := new(mysql_model.Task)
	if has, _ := sys_service.NewSessionWithCond(params).GetOne(task); has {
		if task.Status == common.TASK_STATUS_SUCCESS {
			taskLock.Lock()
			for index, item := range currentTaskList {
				if item.Id == task.Id {
					currentTaskList = append(currentTaskList[0:index], currentTaskList[index+1:]...)
				}
			}
			taskLock.Unlock()
		}
	}
}
