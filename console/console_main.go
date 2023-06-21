package console

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"skygo_detection/console/case_monitor"
	"skygo_detection/console/scanner"
)

var (
	taskType string // 任务类型
	taskName string // 任务名称
)

func init() {
	flag.StringVar(&taskType, "t", "", "task type")
	flag.StringVar(&taskName, "n", "", "task name")
}

/*
 * 启动测试用例监控：go run main.go console -t case_monitor -c ./config/dev/config.tml
 * 启动扫描任务分发器：go run main.go console -t scanner -c ./config/dev/config.tml
 */
func Main(taskType string, taskName string) {
	defer initRecover()
	switch taskType {
	case "case_monitor":
		case_monitor.Run()
	case "scanner":
		scanner.Run()
	}
}

func initRecover() func() {
	return func() {
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
			logMessage := fmt.Sprintf("Trace: %s\n", err)
			logMessage += fmt.Sprintf("\n%s", stacktrace)
			println(logMessage)

			// 程序运行是“定时任务”模式，则捕捉到异常时，推出状态为1。这样HULK定时任务平台才会知道任务执行状态失败
			os.Exit(1)
		}
	}
}
