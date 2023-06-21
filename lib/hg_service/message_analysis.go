package hg_service

import (
	"errors"
	"fmt"
	"runtime"
	"time"

	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/redis"
	"skygo_detection/mysql_model"
	"skygo_detection/service"
)

const (
	HG_MESSAGE_ANALYSIS_RAW    = "hg_scanner:analysis:raw"    // 合规扫描需要进行结果分析的原始消息redis队列
	HG_MESSAGE_ANALYSIS_RESULT = "hg_scanner:analysis:result" // 合规扫描分析后的消息redis队列
)

func getHgMessageAnalysisRawKey() string {
	return fmt.Sprintf("%s:%d", HG_MESSAGE_ANALYSIS_RAW, service.LoadConfig().Http.Port)
}

func getHgMessageAnalysisResultKey() string {
	return fmt.Sprintf("%s:%d", HG_MESSAGE_ANALYSIS_RESULT, service.LoadConfig().Http.Port)
}

var messageAnalysis *MessageAnalysis

func GetMessageAnalysis() *MessageAnalysis {
	if messageAnalysis == nil {
		messageAnalysis = &MessageAnalysis{
			PushChan:    make(chan qmap.QM, 10),
			ReceiveChan: make(chan qmap.QM, 10),
			CloseChan:   make(chan bool),
		}
	}
	return messageAnalysis
}

type MessageAnalysis struct {
	PushChan    chan qmap.QM
	ReceiveChan chan qmap.QM
	CloseChan   chan bool
}

func AddAnalysisMessage(msg qmap.QM) {
	AnalysisRawLog("message_analysis", "", msg.ToString(), "")
	GetMessageAnalysis().PushChan <- msg
}

var testData = `
{
    "data": {
        "block_name": "runBlock",
        "case_id": "TC999030M108",
        "reason": "",
        "result": "{\n   \"attachment\" : {\n      \"file_id\" : \"60a4c8cde830c637640574d6\",\n      \"name\" : \"capturescreen.zip\"\n   },\n   \"result_detail\" : \"true\"\n}\n",
        "status": "success",
        "task_id": "ZIVMRM"
    },
    "sequence": 1621412045966,
    "type": "start_case_block_reply"
}`

func (this *MessageAnalysis) Run() {
	// 获取需要分析的消息，推动给安全分析团队
	go this.pushAnalysisMessage()
	// 接受分析后的消息，并入库
	go this.receiveAnalysisMessage()
}

func (this *MessageAnalysis) Close() {
	close(this.CloseChan)
}

// 获取需要分析的消息，并插入到redis队列中
func (this *MessageAnalysis) pushAnalysisMessage() {
	redis := new(redis.Redis_service)
	for {
		select {
		case <-this.CloseChan:
			return
		case msg := <-this.PushChan:
			println("Analysis Push msg", HG_MESSAGE_ANALYSIS_RAW, msg.ToString())
			redis.LPush(getHgMessageAnalysisRawKey(), msg.ToString())
		}
	}
}

// 获取需要分析的消息，并插入到redis队列中
func (this *MessageAnalysis) receiveAnalysisMessage() {
	redis := new(redis.Redis_service)
	for {
		select {
		case <-this.CloseChan:
			return
		default:
			if len := redis.LLen(getHgMessageAnalysisResultKey()); len.Val() > 0 {
				msgStr := redis.LPop(getHgMessageAnalysisResultKey()).Val()
				println("Analysis receive msg", HG_MESSAGE_ANALYSIS_RESULT, msgStr)
				this.analysisReceiveMessage(msgStr)
			} else {
				<-time.After(time.Second)
			}
		}
	}
}

// 分析处理收到的消息
/*
   {
		"type":"start_case_analyze_replay",
		"sequence":12313,
		"data":{
				"status":"success",//分析插件是否正确执行
				"reason":"",//分析插件启动失败原因
				"task_id": "1212313",
				"case_id":"1212313",
				"block_name":"runBlock",
				"analyze_result":{
					"case_result":"success", //"success","failed","error" 测试项结果
					"logs":"",  //脚本执行和分析结果时产生的日志, 预留
					"remark":"测试件中检测到adbd服务，未通过本项测试。"  //对测试结果的简单说明，可以展示在前端界面中
				}
       }
   }
*/
func (this *MessageAnalysis) analysisReceiveMessage(msgStr string) {
	defer func() {
		logMessage := ""
		if err := recover(); err != nil {
			var stacktrace string
			for i := 1; ; i++ {
				_, f, l, got := runtime.Caller(i)
				if !got {
					break
				}
				stacktrace += fmt.Sprintf("%s:%d\n", f, l)
			}
			logMessage = fmt.Sprintf("Trace: %s\n", err)
			logMessage += fmt.Sprintf("%s", stacktrace)

		}
		AnalysisResultLog("terminal_message_handle", "", msgStr, logMessage)
	}()

	if msg, err := qmap.NewWithString(msgStr); err == nil {
		this.addBlockAnalysisResult(msg.Map("data"))
	} else {
		panic(err)
	}
}

// 添加block分析结果
/*
 * {
 *     "status": "success",
 *     "reason": "",
 *     "task_id": "1212313",
 *     "case_id": "1212313",
 *     "block_name": "runBlock",
 *     "analyze_result": {
 *         "case_result": "success", //success failed error
 *         "logs": "",
 *         "remark": "测试件中检测到adbd服务，未通过本项测试。"
 *     }
 * }
 */
func (this *MessageAnalysis) addBlockAnalysisResult(result qmap.QM) {
	analyzeResult := result.Map("analyze_result")
	update := qmap.QM{
		"status": common.CASE_STATUS_COMPLETED,
	}
	if analyzeResult.String("case_result") == "success" {
		update["test_result_status"] = common.CASE_TEST_STATUS_PASS
		update["action_status"] = common.CASE_STATUS_COMPLETED
	} else if analyzeResult.String("case_result") == "failed" {
		update["test_result_status"] = common.CASE_TEST_STATUS_UNPASS
		update["action_status"] = common.CASE_STATUS_COMPLETED
	} else {
		// update["test_result_status"] = common.CASE_TEST_STATUS_UNPASS
		update["action_status"] = common.CASE_STATUS_FAIL
	}
	if remark := analyzeResult.String("remark"); remark != "" {
		update["test_procedure"] = remark
	}
	taskTestCase := new(mysql_model.TaskTestCase)
	if has, err := sys_service.NewSession().Session.Where("task_uuid = ? and case_uuid = ?", result.MustString("task_id"), result.MustString("case_id")).Get(taskTestCase); err == nil {
		if taskTestCase.ActionStatus != common.CASE_STATUS_ANALYSIS {
			panic(errors.New(`这个测试用例已经不处于"待分析"状态，不再接收分析结果`))
		}
		if has {
			// 更新测试用例信息
			if _, err := new(mysql_model.TaskTestCase).UpdateCase(result.MustString("task_id"), result.MustString("case_id"), update); err != nil {
				panic(err)
			}
		} else {
			panic(errors.New("测试用例不存在"))
		}
	} else {
		panic(err)
	}
}
