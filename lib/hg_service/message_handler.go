package hg_service

import (
	"fmt"
	"runtime"
	"time"

	"github.com/gorilla/websocket"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/redis"
	"skygo_detection/mysql_model"
)

const (
	MAIN_CASE_BLOCK           = "runCaseBlock"  // 默认主block
	CAPTURE_SCREEN_CASE_BLOCK = "captureScreen" // 截屏

	CASE_TEST_LEVEL_LABOUR    = 1
	CASE_TEST_LEVEL_SEMI_AUTO = 2
	CASE_TEST_LEVEL_AUTO      = 3

	// 消息类型
	MSG_TYPE_SESSION                     = "session"                     // 推送会话
	MSG_TYPE_RESESSION                   = "reconnect_session"           // 重连会话
	MSG_TYPE_ACK                         = "ack"                         // 信息确认送达
	MSG_TYPE_UPDATE_TASK_STATUS          = "update_task_status"          // 更新任务状态
	MSG_TYPE_UPDATE_CASE_STATUS          = "update_case_status"          // 更新测试用例状态
	MSG_TYPE_HEART_BEAT                  = "heart_beat"                  // 心跳
	MSG_TYPE_SCANNER_STATUS_SYNC         = "scanner_status_sync"         // 扫描器状态同步
	MSG_TYPE_SCANNER_STATUS_SYNC_REPLY   = "scanner_status_sync_reply"   // 扫描器状态同步应答
	MSG_TYPE_TERMINAL_INFO               = "terminal_info"               // 终端信息上传
	MSG_TYPE_TERMINAL_INFO_REPLY         = "terminal_info_reply"         // 终端信息上传应答
	MSG_TYPE_TERMINAL_UPDATE_CASES       = "terminal_update_cases"       // 终端更新测试用例
	MSG_TYPE_TERMINAL_UPDATE_CASES_REPLY = "terminal_update_cases_reply" // 终端更新测试用例应答
	MSG_TYPE_START_CASE                  = "start_case"                  // 开始扫描测试用例
	MSG_TYPE_START_CASE_REPLY            = "start_case_reply"            // 开始扫描测试用例应答
	MSG_TYPE_END_CASE                    = "end_case"                    // 结束扫描测试用例
	MSG_TYPE_END_CASE_REPLY              = "end_case_reply"              // 结束扫描测试用例应答
	MSG_TYPE_START_CASE_BLOCK            = "start_case_block"            // 开始执行测试用例block
	MSG_TYPE_START_CASE_BLOCK_REPLY      = "start_case_block_reply"      // 开始执行测试用例block应答
)

// 初始化消息
func (this *ScanMsgTerminal) initMessage() error {
	// 长连接连接成功后，推送会话消息给终端
	if err := this.pushSessionInfo(); err != nil {
		return err
	}
	// 对重连的扫描器，需要和扫描器同步状态
	if this.TerminalType == TERMINAL_TYPE_CLIENT {
		// 如果扫描器重连，则向扫描器发送状态更新命令
		// if this.IsReconnect == true {
		// 	if err := this.pushScannerStatusSync(); err != nil {
		// 		return err
		// 	}
		// }
		new(redis.Redis_service).LRemoveAll(this.LowPriorityQueue)
		new(redis.Redis_service).LRemoveAll(this.HighPriorityQueue)
	}
	return nil
}

// 获取任务状态信息
// 如果任务状态获取为空，则定时抓取任务状态信息，不为空，则不再定时抓取任务状态信息
// 如果任务状态为"创建"或者"获取信息"阶段，则发送"请求终端信息"命令
func (this *ScanMsgTerminal) retrieveTaskStatus() {
	if this.TerminalType == TERMINAL_TYPE_CLIENT {
		sendTerminal := false
		for {
			select {
			case <-this.CloseChan:
				return
			default:
				// 查询任务状态，保证终端完成测试用例下载
				if this.IsDownloadedCase {
					println("IsDownloadedCase true")
					go this.pushTestCase()
					return
				} else if !sendTerminal {
					println("IsDownloadedCase false")
					// 给终端发一次下载测试用例的命令
					this.ensureTerminalDownloadTestCase()
					sendTerminal = true
				}
				<-time.After(time.Millisecond * 200)
			}
		}
	}
}

func (this *ScanMsgTerminal) ensureTerminalDownloadTestCase() {
	if taskStatus, err := this.getTaskStatus(); err == nil {
		this.TaskStatus = taskStatus
		if this.TaskStatus == common.HG_TEST_TASK_STATUS_CREATE || this.TaskStatus == common.HG_TEST_TASK_STATUS_CLIENT_INFO || this.TaskStatus == common.HG_TEST_TASK_STATUS_CHOOSE_TEST_CASE {
			// 如果扫描任务处于"创建"或者"获取信息"阶段，则发送"请求终端信息"命令
			this.pushTerminalInfo()
		} else if this.TaskStatus == common.HG_TEST_TASK_STATUS_AUTO_TEST {
			// 如果测试任务处于"测试中"阶段，判断如果设备属于重连，则不用重新下发测试用例下载命令
			if this.IsReconnect {
				this.IsDownloadedCase = true
			} else {
				// 否则下发"请求终端信息"命令
				this.pushTerminalInfo()
			}
		}
	}
}

// 查询测试任务状态
func (this *ScanMsgTerminal) getTaskStatus() (string, error) {
	if task, err := new(mysql_model.HgTestTask).FindOne(this.TaskId); err == nil {
		return task.Status, err
	} else {
		return "", err
	}
}

// 如果测试任务处于"测试中"，则不断监听测试用例，将所有待测试的测试用例按照顺序串行下发到终端
func (this *ScanMsgTerminal) pushTestCase() {
	if this.TerminalType == TERMINAL_TYPE_CLIENT {
		for {
			select {
			case <-this.CloseChan:
				return
			default:
				// 获取下一个待执行的合规测试用例
				// {
				// 	"type":"start_case",
				// 	"sequence":12313,
				// 	"data":{
				// 		"case_id":"123131",
				// 		"case_type":"so",//测试用例类型：jar/so/apk
				// 	}
				// }
				if this.TerminalStatus == TERMINAL_STATUS_IDLE {
					if testCase, has := new(mysql_model.TaskTestCase).GetNextHgExecTestCase(this.TaskId); has {
						caseType := ""
						if taskParam, err := qmap.NewWithString(testCase.String("task_param")); err == nil {
							caseType = taskParam.String("case_type")
						}
						caseInfo := qmap.QM{
							"type":     MSG_TYPE_START_CASE,
							"sequence": custom_util.GetCurrentMilliSecond(),
							"data": qmap.QM{
								"case_id":   testCase.String("case_uuid"),
								"case_type": caseType, // 测试用例类型：jar/so/apk
							},
						}
						this.TerminalStatus = TERMINAL_STATUS_SCANNING
						this.addLowPriorityMsgToCurrentTerminal(caseInfo)
					}
					<-time.After(time.Second * 4)
				} else {
					<-time.After(time.Second)
				}
			}
		}
	}
}

// 推送会话消息
func (this *ScanMsgTerminal) pushSessionInfo() error {
	msg := qmap.QM{
		"sequence": custom_util.GetCurrentMilliSecond(),
		"data": qmap.QM{
			"session_id": this.Sid,
		},
	}
	if this.IsReconnect {
		msg["type"] = MSG_TYPE_RESESSION
	} else {
		msg["type"] = MSG_TYPE_SESSION
	}

	return this.Websocket.WriteMessage(websocket.TextMessage, []byte(msg.ToString()))
}

// 消息确认
func (this *ScanMsgTerminal) receiveConfirm(sequence int64) {
	msg := qmap.QM{
		"type":     MSG_TYPE_ACK,
		"sequence": sequence,
	}
	this.addHighPriorityMsgToCurrentTerminal(msg)
}

// 推送状态同步消息
func (this *ScanMsgTerminal) pushScannerStatusSync() error {
	msg := qmap.QM{
		"type":     MSG_TYPE_SCANNER_STATUS_SYNC,
		"sequence": custom_util.GetCurrentMilliSecond(),
	}
	return this.Websocket.WriteMessage(websocket.TextMessage, []byte(msg.ToString()))
}

// 推送终端信息获取命令
func (this *ScanMsgTerminal) pushTerminalInfo() {
	msg := qmap.QM{
		"type":     MSG_TYPE_TERMINAL_INFO,
		"sequence": custom_util.GetCurrentMilliSecond(),
	}
	this.addHighPriorityMsgToCurrentTerminal(msg)
}

// 对从终端接收的消息进行统一处置
func (this *ScanMsgTerminal) receiveMsgHandle(msgType int, msgStr string) {
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
		TerminalReceiveLog("terminal_message_handle", this.TaskId, msgStr, logMessage)
	}()
	msg, err := qmap.NewWithString(msgStr)
	if err != nil {
		panic(err)
	}
	sequence := msg.Int64("sequence")
	messageType := msg.String("type")
	switch messageType {
	case MSG_TYPE_HEART_BEAT:
		// 心跳
		this.heartBeatMSH(msg)
	case MSG_TYPE_SCANNER_STATUS_SYNC_REPLY:
		// 扫描器状态同步回应
		this.scannerStatusSyncReplyMSH(msg)
	case MSG_TYPE_TERMINAL_INFO_REPLY:
		// 上传终端信息
		this.terminalInfoReplyMSH(msg)
	case MSG_TYPE_TERMINAL_UPDATE_CASES_REPLY:
		// 终端测试用例更新应答
		this.terminalUpdateCaseReplyMSH(msg)
	case MSG_TYPE_START_CASE_REPLY:
		// 测试用例初始化响应命令
		this.startCaseReplyMSH(msg)
	case MSG_TYPE_END_CASE:
		// 测试用例结束命令
		this.endCaseMSH(msg)
	case MSG_TYPE_END_CASE_REPLY:
		// 测试用例结束响应命令
		this.endCaseReplyMSH(msg)
	case MSG_TYPE_START_CASE_BLOCK:
		// 测试用例block命令
		this.startCaseBlockMSH(msg)
	case MSG_TYPE_START_CASE_BLOCK_REPLY:
		// 测试用例block响应命令
		this.startCaseBlockReplyMSH(msg)
	default:

	}
	// 对收到的所有消息进行确认
	if messageType != MSG_TYPE_ACK && messageType != MSG_TYPE_HEART_BEAT {
		this.receiveConfirm(sequence)
	}
}

// 心跳消息处理
func (this *ScanMsgTerminal) heartBeatMSH(msg qmap.QM) {
	if this.TerminalType == TERMINAL_TYPE_CLIENT {
		// 更新终端的在线时间
		if err := new(mysql_model.HgTestTask).UpdateTerminalConnectionTime(this.TaskId); err != nil {
			panic(err)
		}
		// 将来自扫描器的心跳消息转发给web端
		this.addHighPriorityMsgToOppositeTerminal(msg)
		this.addHighPriorityMsgToCurrentTerminal(msg)
	}
}

// 扫描器状态同步回应消息处理
/*
	{
		"type":"scanner_status_sync_reply",
		"sequence":12313,
		"data":{
			 "status":1 //（1：busy,0:idle）
		}
	}
*/
func (this *ScanMsgTerminal) scannerStatusSyncReplyMSH(msg qmap.QM) {
	// terminalInfo := msg.Map("data")
	// this.TerminalStatus = terminalInfo.Int("status")
}

// 终端信息上传响应消息处理
/*
	{
		"type":"terminal_info",
		"sequence":12313,
		"data":{
			"cpu": "64",
			"os_type": "linux",
			"os_version": "1.0.0"
		}
	}
*/
func (this *ScanMsgTerminal) terminalInfoReplyMSH(msg qmap.QM) {
	terminalInfo := msg.Map("data")
	(terminalInfo)["task_id"] = this.TaskId
	info := msg.Map("data")
	info["status"] = common.HG_TEST_TASK_STATUS_CHOOSE_TEST_CASE
	// 更新任务关联的终端信息
	if _, err := new(mysql_model.HgTestTask).Update(this.TaskId, info); err == nil {
		// 硬件信息上传成功后，将测试用例下载地址发给终端，让终端更新测试用例
		terminalUpdateCaseMsg := qmap.QM{
			"type":     MSG_TYPE_TERMINAL_UPDATE_CASES,
			"sequence": custom_util.GetCurrentMilliSecond(),
			"data":     qmap.QM{"url": "/message/v1/hg_scanner/download_case"},
		}
		this.addLowPriorityMsgToCurrentTerminal(terminalUpdateCaseMsg)
		// // 硬件信息上传成功后，通知web端更新扫描任务状态
		// updateTaskStatusMsg := qmap.QM{
		// 	"type":     MSG_TYPE_UPDATE_TASK_STATUS,
		// 	"sequence": custom_util.GetCurrentMilliSecond(),
		// }
		// this.addHighPriorityMsgToOppositeTerminal(updateTaskStatusMsg)
	} else {
		panic(err)
	}
}

// 终端更新测试用例应答信息处理
/*
	{
		"type":"terminal_update_cases_reply",
		"sequence":12313,
		"data":{
		  "status":"success/fail",
		  "reason":""
		}
	}
*/
func (this *ScanMsgTerminal) terminalUpdateCaseReplyMSH(msg qmap.QM) {
	status := msg.Map("data").String("status")
	if status == "success" {
		this.IsDownloadedCase = true
		info := qmap.QM{
			"status": common.HG_TEST_TASK_STATUS_AUTO_TEST,
		}
		// 更新任务信息
		if _, err := new(mysql_model.HgTestTask).Update(this.TaskId, info); err != nil {
			panic(err)
		}
		// 测试用例更新成功后，通知web端更新扫描任务状态
		updateTaskStatusMsg := qmap.QM{
			"type":     MSG_TYPE_UPDATE_TASK_STATUS,
			"sequence": custom_util.GetCurrentMilliSecond(),
		}
		this.addHighPriorityMsgToOppositeTerminal(updateTaskStatusMsg)
	}
}

// 测试用例初始化命令
/*
	{
		"type":"start_case",
		"sequence":12313,
		"data":{
		  "case_id":"123131",
		  "case_type":"so",//测试用例类型：jar/so/apk
		}
	}
*/
func (this *ScanMsgTerminal) startCaseMSH(msg qmap.QM) {
	// 更新测试
	this.addHighPriorityMsgToOppositeTerminal(msg)
}

// 测试用例初始化响应命令
/*
	{
		"type":"start_case_reply",
		"sequence":12313,
		"data":{
		  "case_id":"123131",
		  "status":"success/fail",
		  "reason":""
		}
	}
*/
func (this *ScanMsgTerminal) startCaseReplyMSH(msg qmap.QM) {
	data := msg.Map("data")
	caseUuid := data.MustString("case_id")
	updateInfo := qmap.QM{}
	if status := data.String("status"); status == "success" {
		updateInfo["action_status"] = common.CASE_STATUS_TESTING
		// 命令下发成功，扫描器处于繁忙状态
		this.TerminalStatus = TERMINAL_STATUS_SCANNING
		blockParams := qmap.QM{
			"case_id":          caseUuid,
			"case_type":        "", // 测试用例类型：jar/so/apk
			"block_name":       MAIN_CASE_BLOCK,
			"time_out":         3600000,
			"test_level":       1,
			"ret_of_crash":     "",
			"upload_crash_log": true,
			"case_parameter":   "{}",
		}
		if caseParam, err := new(mysql_model.TaskTestCase).GetCaseParams(this.TaskId, caseUuid); err == nil {
			blockParams["case_type"] = caseParam.String("case_type")
			blockParams["time_out"] = caseParam.Int("time_out")
			blockParams["test_level"] = caseParam.Int("test_level")
			blockParams["ret_of_crash"] = caseParam.String("ret_of_crash")
			blockParams["case_parameter"] = caseParam.String("case_parameter")
			blockParams["upload_crash_log"] = caseParam.Bool("upload_crash_log")
		}
		caseInfo := qmap.QM{
			"type":     MSG_TYPE_START_CASE_BLOCK,
			"sequence": custom_util.GetCurrentMilliSecond(),
			"data":     blockParams,
		}
		this.addHighPriorityMsgToCurrentTerminal(caseInfo)
	} else {
		updateInfo["action_status"] = common.CASE_STATUS_COMPLETED
		updateInfo["test_result_status"] = common.CASE_TEST_STATUS_UNPASS

		this.TerminalStatus = TERMINAL_STATUS_IDLE
	}
	if _, err := new(mysql_model.TaskTestCase).UpdateCase(this.TaskId, caseUuid, updateInfo); err != nil {
		panic(err)
	}
}

func (this *ScanMsgTerminal) startCaptureScreenBlockMSH(caseUuid string) {
	caseInfo := qmap.QM{
		"type":     MSG_TYPE_START_CASE_BLOCK,
		"sequence": custom_util.GetCurrentMilliSecond(),
		"data": qmap.QM{
			"case_id":          caseUuid,
			"case_type":        "so", // 测试用例类型：jar/so/apk
			"block_name":       CAPTURE_SCREEN_CASE_BLOCK,
			"time_out":         36000000,
			"test_level":       1,
			"ret_of_crash":     "",
			"upload_crash_log": false,
		},
	}
	<-time.After(time.Second)
	this.addLowPriorityMsgToCurrentTerminal(caseInfo)
}

// 测试用例结束命令
/*
	{
		"type":"end_case",
		"sequence":12313,
		"data":{
		  "case_id":"123131"
		}
	}
*/
func (this *ScanMsgTerminal) endCaseMSH(msg qmap.QM) {
	this.addLowPriorityMsgToOppositeTerminal(msg)
}

// 测试用例结束响应命令
/*
	{
		"type":"end_case_reply",
		"sequence":12313,
		"data":{
		  "case_id":"123131",
		  "status":"success/fail",
		  "reason":""
		}
	}
*/
func (this *ScanMsgTerminal) endCaseReplyMSH(msg qmap.QM) {
	// 命令下发成功，扫描器处于繁忙状态
	this.TerminalStatus = TERMINAL_STATUS_IDLE
}

// 测试用例block命令
/*
	{
		"type":"start_case_block",
		"sequence":12313,
		"data":{
			"case_id":"1212313",
			"case_type":"so",//测试用例类型：jar/so/apk
			"block_name":"runBlock",
			"time_out":10//block执行超时(单位:秒)
		}
	}
*/
func (this *ScanMsgTerminal) startCaseBlockMSH(msg qmap.QM) {
	this.addLowPriorityMsgToOppositeTerminal(msg)
}

// 测试用例block响应命令
/*
{
  "sequence": 1633919661746,
  "type": "start_case_block_reply",
  "data": {
    "status": "success",
    "reason": "",
    "case_id": "hg_CASE_6",
    "block_name": "runCaseBlock",
    "result": "{\n   \"attachment\" : {\n      \"file_id\" : \"6163a2ac24b64710652d0647\",\n      \"name\" : \"output.zip\"\n   },\n   \"result_detail\" : \"success\"\n}\n"
  }
}"

{
    "sequence": 1633919662914,
    "type": "start_case_block_reply",
    "data": {
        "status": "success",
        "reason": "",
        "case_id": "hg_CASE_6",
        "block_name": "captureScreen",
        "result": "{\"result_detail\":\"true\",\"attachment\":{\"name\":\"1633919662111.png\",\"file_id\":\"6163a2ad24b64710652d0649\"}}"
    }
}
*/
func (this *ScanMsgTerminal) startCaseBlockReplyMSH(msg qmap.QM) {
	data := msg.Map("data")
	caseUuid := data.String("case_id")
	blockName := data.String("block_name")
	// 更新测试用例block结果
	if blockName == MAIN_CASE_BLOCK {
		// 如果block是默认主block,则去获取该测试用例的测试级别
		testCase, err := new(mysql_model.TaskTestCase).FindOne(this.TaskId, caseUuid)
		if err != nil {
			panic(err)
		}
		update := qmap.QM{
			"case_result": data.String("result"),
		}
		TerminalReceiveLog("startCaseBlockReplyMSH", this.TaskId, fmt.Sprintf("case_test_level:%d, status:%s", testCase.AutoTestLevel, data.String("status")), "")
		if testCase.AutoTestLevel == CASE_TEST_LEVEL_AUTO {
			if data.String("status") == "success" {
				// 如果该测试用例是自动测试用例，则将结果丢给安全分析团队进行结果分析
				(data)["task_id"] = this.TaskId
				update["action_status"] = common.CASE_STATUS_ANALYSIS
				AddAnalysisMessage(msg)
			} else {
				update["action_status"] = common.CASE_STATUS_FAIL
			}
		}

		// 更新测试用例信息
		if _, err := new(mysql_model.TaskTestCase).UpdateCase(this.TaskId, caseUuid, update); err != nil {
			panic(err)
		}
		if testCase.AutoTestLevel == CASE_TEST_LEVEL_AUTO {
			// 如果是自动测试用例，则直接向终端发送结束测试用例命令
			this.sendEndCase(caseUuid)
		}
	} else if blockName == CAPTURE_SCREEN_CASE_BLOCK {
		if data.String("result") == "" {
			return
		}
		if result, err := qmap.NewWithString(data.String("result")); err == nil {
			if result.String("result_detail") == "true" {
				update := qmap.QM{
					"case_result_file": result.Map("attachment"),
				}
				if _, err := new(mysql_model.TaskTestCase).UpdateCase(this.TaskId, caseUuid, update); err != nil {
					panic(err)
				}
			}
		} else {
			panic(err)
		}
		this.addLowPriorityMsgToOppositeTerminal(msg)
	}
}

// 结束测试用例
func (this *ScanMsgTerminal) sendEndCase(caseId string) {
	caseType := ""
	if caseParam, err := new(mysql_model.TaskTestCase).GetCaseParams(this.TaskId, caseId); err == nil {
		caseType = caseParam.String("case_type")
	}
	caseInfo := qmap.QM{
		"type":     MSG_TYPE_END_CASE,
		"sequence": custom_util.GetCurrentMilliSecond(),
		"data": qmap.QM{
			"case_id":   caseId,
			"case_type": caseType, // 测试用例类型：jar/so/apk
		},
	}
	<-time.After(time.Second)
	this.addLowPriorityMsgToCurrentTerminal(caseInfo)
}
