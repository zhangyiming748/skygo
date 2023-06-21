package scanner

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"go.uber.org/zap"
	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util/clog"
	"skygo_detection/mysql_model"
)

const (
	MaxErrorRetryTimes = 3               // 最大错误重试次数
	IntervalTime       = 2 * time.Second // 任务轮询执行时间间隔
	MaxTaskNum         = 20              // 最大扫描任务并发数量
	TaskExecInterval   = int64(10)       // 任务执行间隔(单位:秒)
)

// go run main.go scanner -c ./config/dev/config.tml
func Run() {
	defer func() {
		if err := recover(); err != nil {
			clog.Error("console Main panic", zap.Any("error", err))
		}
	}()
	scanner := NewScanner()
	scanner.run()
}

type ScannerInterface interface {
	Run(qm qmap.QM) (int, error)
}

var ScannerList = map[string]ScannerInterface{
	common.TOOL_FIRMWARE_SCANNER: FirmwareScanner{},
}

func NewScanner() *Scanner {
	return &Scanner{Status: true}
}

type Scanner struct {
	CurrentTaskNum int32
	Status         bool // 扫描器状态（false:停止，true:运行）
	taskWG         sync.WaitGroup
}

/**
 * @Description: 启动扫描器
 */
func (this *Scanner) run() {
	this.handleInterrupt()
	for {
		if !this.Status {
			break
		}
		if rows, err := sys_service.NewOrm().In("status", []interface{}{common.SCANNER_STATUS_READY, common.SCANNER_STATUS_HANDING}).And("next_exec_time <= ?", time.Now().Unix()).Rows(new(mysql_model.ScannerTask)); err == nil {
			tasks := []*mysql_model.ScannerTask{}
			for rows.Next() {
				task := new(mysql_model.ScannerTask)
				err = rows.Scan(task)
				if err != nil {
					panic(err)
				}
				tasks = append(tasks, task)
			}
			rows.Close()
			// 从任务数组中读取任务进行分发
			this.dispatchTask(tasks)
			// 扫描器休眠
			this.sleeping()
		} else {
			panic(err)
		}
	}
}

/**
 * @Description:分发扫描任务
 * @param tasks
 */
func (this *Scanner) dispatchTask(tasks []*mysql_model.ScannerTask) {
	for _, task := range tasks {
		for {
			if this.CurrentTaskNum < MaxTaskNum {
				break
			} else {
				// 如果最大扫描任务数达到上限,则一直等待
				<-time.After(time.Millisecond * 200)
			}
		}
		atomic.AddInt32(&this.CurrentTaskNum, 1) // 原子操作
		taskInfo := qmap.QM{
			"id":           task.Id,
			"scanner_id":   task.ScannerId,
			"name":         task.Name,
			"scanner_type": task.ScannerType,
			"retry_times":  task.RetryTimes,
		}
		this.taskWG.Add(1)
		go this.runScannerTask(task.ScannerType, taskInfo)
	}
	// 等待所有已经分发出去的任务执行完成
	this.taskWG.Wait()
}

func (this *Scanner) runScannerTask(taskType string, taskInfo qmap.QM) {
	errMsg := ""
	sleepTime := TaskExecInterval
	defer func() {
		if err := recover(); err != nil {
			var stacktrace string
			for i := 1; ; i++ {
				_, f, l, got := runtime.Caller(i)
				if !got {
					break
				}
				stacktrace += fmt.Sprintf("%s:%d\n", f, l)
			}
			// when stack finishes
			errMsg = fmt.Sprintf("Trace: %s\n", err)
			errMsg += fmt.Sprintf("\n%s", stacktrace)
		}
		this.UpdateTaskStatus(taskInfo, errMsg, sleepTime)
		atomic.AddInt32(&this.CurrentTaskNum, -1)
		this.taskWG.Done()
	}()

	scanner, err := this.getScanner(taskType)
	if err != nil {
		errMsg = err.Error()
		return
	}
	stime, taskErr := scanner.Run(taskInfo)
	sleepTime = int64(stime)
	if taskErr != nil {
		errMsg = taskErr.Error()
	}
}

/**
 * @Description: 根据任务类型查询任务
 * @param scannerType
 * @return ScannerInterface
 * @return error
 */
func (this *Scanner) getScanner(scannerType string) (ScannerInterface, error) {
	if val, has := ScannerList[scannerType]; has {
		return val, nil
	} else {
		return nil, errors.New(fmt.Sprintf("未知的扫描任务类型:%s", scannerType))
	}
}

/**
 * @Description: 扫描器休眠
 */
func (this *Scanner) sleeping() {
	sleepInterval := time.Millisecond * 500
	var sleepTime time.Duration = 0
	for {
		if this.Status && sleepTime <= IntervalTime {
			<-time.After(sleepInterval)
			sleepTime += sleepInterval
		} else {
			break
		}
	}
}

/**
 * @Description: 更新任务状态
 */
func (this *Scanner) UpdateTaskStatus(taskInfo qmap.QM, errMsg string, sleepTime int64) {
	rawInfo := qmap.QM{
		"status": common.SCANNER_STATUS_HANDING,
	}
	logMsg := "更新任务状态为：扫描中"
	// 如果休眠时间小于0，则认为该任务不再需要执行
	if sleepTime < 0 {
		if errMsg != "" {
			// 如果不需要再次执行，并且此次执行失败
			rawInfo["retry_times"] = taskInfo.Int("retry_times") + 1
			rawInfo["status"] = common.SCANNER_STATUS_FAILURE
			logMsg = "更新任务状态为：扫描失败"
		} else {
			// 如果不需要再次执行，并且此次执行成功
			rawInfo["status"] = common.SCANNER_STATUS_SUCCESS
			logMsg = "更新任务状态为：扫描成功"
		}
	} else {
		if errMsg != "" {
			// 如果任务执行失败，并且超过最大重试次数，则将任务状态置为执行失败
			retryTimes := taskInfo.Int("retry_times") + 1
			rawInfo["retry_times"] = retryTimes
			if retryTimes >= MaxErrorRetryTimes {
				rawInfo["status"] = common.SCANNER_STATUS_FAILURE
				logMsg = "更新任务状态为：扫描失败"
			} else {
				rawInfo["next_exec_time"] = time.Now().Unix() + sleepTime
			}
		} else {
			rawInfo["next_exec_time"] = time.Now().Unix() + sleepTime
		}
	}
	// 更新任务状态
	new(mysql_model.ScannerTask).Update(taskInfo.MustInt("id"), rawInfo)
	// 插入任务执行日志
	new(mysql_model.ScannerTaskLog).Insert(taskInfo, logMsg, errMsg)
}

/**
 * @Description: 处理终端信号
 */
func (this *Scanner) handleInterrupt() {
	signChan := make(chan os.Signal, 1)
	signal.Notify(signChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	go func() {
		<-signChan
		this.Status = false
	}()
}
