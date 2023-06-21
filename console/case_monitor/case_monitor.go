package case_monitor

import (
	"fmt"
	"runtime"
	"sync"
	"time"

	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/mysql_model"
	"skygo_detection/service"
)

type none struct{}

/*
*
用于协程之间告知“定时任务协程”平滑停止完毕
*/
var schedulerCloser = make(chan none)

var running bool
var runtimeErrors = 0
var stopChan = make(chan int)

func Run() {
	running = true
	// interrupt信号捕捉后调用函数
	service.InterruptHandleAddFunc(func() {
		fmt.Println("received interrupt signal, scheduler stop: case_monitor")
		running = false
		<-schedulerCloser
	})
	run()
	close(schedulerCloser)
}

func initRecover() func() {
	return func() {
		if err := recover(); err != nil {
			runtimeErrors++
			var stacktrace string
			for i := 1; ; i++ {
				_, f, l, got := runtime.Caller(i)
				if !got {
					break
				}
				stacktrace += fmt.Sprintf("%s:%d\n", f, l)
			}
			// when stack finishes
			logMessage := fmt.Sprintf("Trace: %s\n", err)
			logMessage += fmt.Sprintf("\n%s", stacktrace)
			service.GetDefaultLogger("case_monitor").Info(logMessage)
		}
	}
}

func run() {
	defer initRecover()
	for running {
		if runtimeErrors > 5 {
			break
		}
		updateTaskList()
		wg := sync.WaitGroup{}
		for _, task := range currentTaskList {
			wg.Add(1)
			go monitorCase(task, &wg)
		}
		wg.Wait()
		<-time.After(time.Second * 2)
	}
}

func monitorCase(task *mysql_model.Task, wg *sync.WaitGroup) {
	defer initRecover()
	defer wg.Add(-1)
	updateCaseResult := false // 任务中的测试用例状态(true：表示任务中还有测试用例未结束，false：表示任务中已经没有需要测试的测试用例了)
	if task.Tool != "" {
		// 如果任务是工具任务
		updateCaseResult = toolTaskMonitor(task)
	} else {
		updateCaseResult = caseTaskMonitor(task)
	}
	if updateCaseResult == false {
		// 如果没有可以更新的测试用例，则尝试将该任务从任务监控列表中删除
		tryRemoveTask(task.Id)
	}
}

// 工具任务监控
func toolTaskMonitor(task *mysql_model.Task) bool {
	switch task.Tool {
	case common.TOOL_FIRMWARE_SCANNER:
		if result := isFirmwareScannerTaskComplete(task.Id); result != 1 {
			completeTask(task)
			return false
		}
	case common.TOOL_VUL_SCANNER:
		if result := isVulScannerTaskCompleted(task.Id); result != 1 {
			// 如果任务关联的漏洞扫描已经结束，则开始更新任务中的测试用例状态，确保至少有一个测试用例处于测试中
			completeTask(task)
			return false
		}
	}
	return true
}

func completeTask(task *mysql_model.Task) {
	info := qmap.QM{
		"status": common.TASK_STATUS_SUCCESS,
	}
	task.UpdateTaskById(task.Id, info, 0, "任务监控服务")
}

// 场景任务监控
func caseTaskMonitor(task *mysql_model.Task) bool {
	updateCaseResult := false // 任务中的测试用例状态(true：表示任务中还有测试用例未结束，false：表示任务中已经没有需要测试的测试用例了)
	// 确保有一个没有关联检测工具的手动测试用例正在测试中
	if result := updateTaskTestCase(task, []int{common.IS_TASK_CASE_MAN}, ""); result == true {
		updateCaseResult = true
	}
	// 查询主任务关联的固件检测任务是否都已经完成（0:未查询到任务 1:未完成 2:完成）
	if isCompleted := isFirmwareScannerTaskComplete(task.Id); isCompleted == 2 {
		// 如果任务关联的固件扫描已经结束，则开始更新任务中的测试用例状态，确保至少有一个测试用例处于测试中
		if result := updateTaskTestCase(task, []int{common.IS_TASK_CASE_MAN, common.IS_TASK_CASE_SEMI}, common.TOOL_FIRMWARE_SCANNER); result == true {
			updateCaseResult = true
		}
	} else if isCompleted == 1 {
		updateCaseResult = true
	}
	// 查询主任务关联的漏洞检测任务是否都已经完成（0:未查询到任务 1:未完成 2:完成）
	if isCompleted := isVulScannerTaskCompleted(task.Id); isCompleted == 2 {
		// 如果任务关联的漏洞扫描已经结束，则开始更新任务中的测试用例状态，确保至少有一个测试用例处于测试中
		if result := updateTaskTestCase(task, []int{common.IS_TASK_CASE_MAN, common.IS_TASK_CASE_SEMI}, common.TOOL_VUL_SCANNER); result == true {
			updateCaseResult = true
		}
	} else if isCompleted == 1 {
		updateCaseResult = true
	}
	return updateCaseResult
}

// 查询主任务关联的固件检测任务是否都已经完成（0:未查询到任务 1:未完成 2:完成）
func isFirmwareScannerTaskComplete(taskId int) int {
	params := qmap.QM{
		"e_task_id": taskId,
	}
	firmwareTasks := []mysql_model.FirmwareTask{}
	if list, err := sys_service.NewSessionWithCond(params).Get(&firmwareTasks); err == nil {
		if len(*list) > 0 {
			for _, item := range *list {
				itemQM := qmap.QM(item)
				if itemQM.Int("status") == common.FIRMWARE_STATUS_SCAN_FAILURE {
					// 如果固件任务执行失败，则将该固件任务关联的测试用例状态置为"测试失败"
					new(mysql_model.TaskTestCase).SetFirmwareCaseStatusToFailure(itemQM.Int("id"))
				} else if itemQM.Int("status") != common.FIRMWARE_STATUS_SCAN_SUCCESS {
					return 1
				}
			}
			return 2
		} else {
			return 0
		}
	} else {
		return 0
	}
}

// 查询主任务关联的漏洞检测任务是否都已经完成（0:未查询到任务 1:未完成 2:完成）
func isVulScannerTaskCompleted(taskId int) int {
	params := qmap.QM{
		"e_parent_id": taskId,
	}
	vulTasks := []mysql_model.VulTask{}
	if list, err := sys_service.NewSessionWithCond(params).Get(&vulTasks); err == nil {
		if len(*list) > 0 {
			for _, item := range *list {
				itemQM := qmap.QM(item)
				if itemQM.Int("status") == common.VUL_UNSTART || itemQM.Int("status") == common.VUL_PRELIMINARY_BEGIN {
					return 1
				}
			}
			return 2
		} else {
			return 0
		}
	} else {
		return 0
	}
}

// 确保有一个符合条件的测试用例处于测试中
func updateTaskTestCase(task *mysql_model.Task, autoTestLevel []int, testTool string) bool {
	if exist := isExistRunningCase(task.Id, autoTestLevel, testTool); !exist {
		return setOneCaseToTesting(task.Id, autoTestLevel, testTool)
	}
	return true
}

// 获取一个待测试的用例，并让其进入测试中
func setOneCaseToTesting(taskId int, autoTestLevel []int, testTool string) bool {
	// 更新测试用例状态前先确保任务状态为"进行中"
	if status := mysql_model.TaskGetStatusById(taskId); status != common.TASK_STATUS_RUNNING {
		return false
	}
	params := qmap.QM{
		"e_task_id":       taskId,
		"e_action_status": common.CASE_STATUS_QUEUING,
	}
	if len(autoTestLevel) > 0 {
		params["in_auto_test_level"] = autoTestLevel
	}
	if testTool != "" {
		params["e_test_tool"] = testTool
	}
	testCase := new(mysql_model.TaskTestCase)
	if has, _ := sys_service.NewSessionWithCond(params).GetOne(testCase); has {
		testCase.ActionStatus = common.CASE_STATUS_TESTING
		_, err := sys_service.NewOrm().Table(testCase).ID(testCase.Id).AllCols().Update(testCase)
		if err != nil {
			panic(err)
		}
		return true
	} else {
		params["e_action_status"] = common.CASE_STATUS_READY
		if has, _ := sys_service.NewSessionWithCond(params).GetOne(testCase); has {
			testCase.ActionStatus = common.CASE_STATUS_TESTING
			_, err := sys_service.NewOrm().Table(testCase).ID(testCase.Id).AllCols().Update(testCase)
			if err != nil {
				panic(err)
			}
			return true
		}
	}
	return false
}

// 判断是否存在正在测试的测试用例
func isExistRunningCase(taskId int, autoTestLevel []int, testTool string) bool {
	params := qmap.QM{
		"e_task_id":       taskId,
		"e_action_status": common.CASE_STATUS_TESTING,
	}
	if len(autoTestLevel) > 0 {
		params["in_auto_test_level"] = autoTestLevel
	}
	if testTool != "" {
		params["e_test_tool"] = testTool
	}
	has, _ := sys_service.NewSessionWithCond(params).GetOne(new(mysql_model.TaskTestCase))
	return has
}
