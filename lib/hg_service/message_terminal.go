package hg_service

import (
	"fmt"
	"time"

	"skygo_detection/guardian/src/net/qmap"

	"github.com/gorilla/websocket"

	"skygo_detection/lib/common_lib/redis"
)

const (
	MAX_RECONNECT_TIMES = 3 // 最大重连次数
	RECONNECT_TIME      = 3 // 失败重连时间间隔(秒)

	// terminal_status 终端状态
	TERMINAL_STATUS_IDLE     = 0 // 闲置中
	TERMINAL_STATUS_SCANNING = 1 // 扫描中

	// 优先级队列级别
	PRIORITY_QUEUE_LOW  = 1 // 低优先级队列
	PRIORITY_QUEUE_HIGH = 5 // 高优先级队列
)

func NewScanMsgTerminal(terminalType, sid, taskId string, ws *websocket.Conn, isReconnect bool) *ScanMsgTerminal {
	scanMsgTerminal := ScanMsgTerminal{
		LowPriorityReceiveChan:  make(chan string),
		HighPriorityReceiveChan: make(chan string),
		ReplyChan:               make(chan bool),
		CloseChan:               make(chan bool),
		Sid:                     sid,
		TerminalType:            terminalType,
		TaskId:                  taskId,
		Websocket:               ws,
		LowPriorityQueue:        generateTaskQueue(terminalType, sid, PRIORITY_QUEUE_LOW),
		HighPriorityQueue:       generateTaskQueue(terminalType, sid, PRIORITY_QUEUE_HIGH),
		IsReconnect:             isReconnect,
		TerminalStatus:          TERMINAL_STATUS_IDLE,
	}
	SubTerminal(&scanMsgTerminal)
	return &scanMsgTerminal
}

type ScanMsgTerminal struct {
	Websocket               *websocket.Conn
	LowPriorityReceiveChan  chan string
	HighPriorityReceiveChan chan string
	ReplyChan               chan bool
	CloseChan               chan bool
	Sid                     string // 长连接通讯的会话id，通过该id实现会话重连
	TaskId                  string // 检测任务id
	TaskStatus              string // 检测任务状态（创建:create  获取信息:client_info   普配用例:choose_test_case  测试:auto_test  完成:complete）
	TerminalType            string // 终端类型(client/web)
	LowPriorityQueue        string // 低优先级消息队列,存储实时性要求低的消息，但是缓存时间可以较长(例如：下发执行的命令)
	HighPriorityQueue       string // 高优先级消息队列，存储实时性要求高的消息(例如：心跳)
	TerminalStatus          int    // 终端状态（1:闲置中，2:扫描中）,如果终端处于繁忙中，那么对下发给终端的命令有一定的限制判断，避免终端进行不断的状态机切换
	IsReconnect             bool   // 是否是重连客户端
	IsDownloadedCase        bool   // 是否下载了测试用例
}

// 开始进行长连接通讯
func (this *ScanMsgTerminal) Run() {
	// 在启动前，需要进行消息初始化
	this.initMessage()
	// 更新任务状态信息
	go this.retrieveTaskStatus()
	// 监控从终端推送过来的消息
	go this.monitorReceiveMsg()
	// 监控推送给当前端的消息(从redis队列中监控)
	go this.monitorPushMsg()
	// 将消息推送给终端
	go this.pushMsgToTerminal()
}

// 停止长连接
func (this *ScanMsgTerminal) Stop() {
	close(this.CloseChan)
}

// 监控推送给当前端的消息(从redis队列中监控)
func (this *ScanMsgTerminal) monitorPushMsg() {
	ack := false
	redis := new(redis.Redis_service)
	for {
		select {
		case <-this.CloseChan:
			return
		default:
			if len := redis.LLen(this.HighPriorityQueue); len.Val() > 0 {
				msg := redis.LIndex(this.HighPriorityQueue, 0).Val()
				this.HighPriorityReceiveChan <- msg
				if ack = <-this.ReplyChan; ack {
					redis.LPop(this.HighPriorityQueue)
				}
			} else if len := redis.LLen(this.LowPriorityQueue); len.Val() > 0 {
				msg := redis.LIndex(this.LowPriorityQueue, 0).Val()
				this.LowPriorityReceiveChan <- msg
				if ack = <-this.ReplyChan; ack {
					redis.LPop(this.LowPriorityQueue)
				}
			} else {
				<-time.After(time.Second)
			}
		}
	}
}

// 监控从终端推送过来的消息
func (this *ScanMsgTerminal) monitorReceiveMsg() {
	for {
		select {
		case <-this.CloseChan:
			return
		default:
			// 循环监听从客户端推送过来的消息
			msgType, msg, err := this.Websocket.ReadMessage()
			if err != nil {
				// 如果接收消息异常，则停止长连接通讯
				this.Stop()
				return
			}
			this.receiveMsgHandle(msgType, string(msg))
		}
	}
}

// 向终端推送消息
func (this *ScanMsgTerminal) pushMsgToTerminal() {
	for {
		select {
		case <-this.CloseChan:
			this.Websocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "close"))
			close(this.ReplyChan)
			return
		case msg := <-this.HighPriorityReceiveChan:
			if this.TerminalType == TERMINAL_TYPE_CLIENT {
				println("push low", this.TerminalType, this.TaskId, msg)
			}
			this.sendMessage(msg)
		case msg := <-this.LowPriorityReceiveChan:
			if this.TerminalType == TERMINAL_TYPE_CLIENT {
				println("push low", this.TerminalType, this.TaskId, msg)
			}
			this.sendMessage(msg)
		}
	}
}

func (this *ScanMsgTerminal) sendMessage(msg string) {
	failTimes := 0 // 推送失败次数，推送失败超过次数限制，任务停止
	for {
		if err := this.Websocket.WriteMessage(websocket.TextMessage, []byte(msg)); err == nil {
			TerminalSendLog("", this.TaskId, msg, "")
			break
		} else {
			failTimes++
			// 如果重试次数超过最大限度，则停止任务
			if failTimes >= MAX_RECONNECT_TIMES {
				this.Stop()
			}
			select {
			case <-this.CloseChan:
				close(this.ReplyChan)
				return
			}
			// 推送失败，等会儿再重试
			<-time.After(time.Second * RECONNECT_TIME)
		}
	}
	this.ReplyChan <- true
}

// 添加低优先级消息到当前终端的消息推送队列
func (this *ScanMsgTerminal) addLowPriorityMsgToCurrentTerminal(msg qmap.QM) {
	PushMsgToTerminal(this.TerminalType, this.TaskId, PRIORITY_QUEUE_LOW, msg)
}

// 添加低优先级消息到对方终端的消息推送队列
// 例如:当前终端是client，则对方终端是:web
func (this *ScanMsgTerminal) addLowPriorityMsgToOppositeTerminal(msg qmap.QM) {
	toTerminal := TERMINAL_TYPE_WEB
	if this.TerminalType == TERMINAL_TYPE_WEB {
		toTerminal = TERMINAL_TYPE_CLIENT
	}
	PushMsgToTerminal(toTerminal, this.TaskId, PRIORITY_QUEUE_LOW, msg)
}

// 添加高优先级消息到当前终端的消息推送队列
func (this *ScanMsgTerminal) addHighPriorityMsgToCurrentTerminal(msg qmap.QM) {
	PushMsgToTerminal(this.TerminalType, this.TaskId, PRIORITY_QUEUE_HIGH, msg)
}

// 添加高优先级消息到对方终端的消息推送队列
// 例如:当前终端是client，则对方终端是:web
func (this *ScanMsgTerminal) addHighPriorityMsgToOppositeTerminal(msg qmap.QM) {
	toTerminal := TERMINAL_TYPE_WEB
	if this.TerminalType == TERMINAL_TYPE_WEB {
		toTerminal = TERMINAL_TYPE_CLIENT
	}
	PushMsgToTerminal(toTerminal, this.TaskId, PRIORITY_QUEUE_HIGH, msg)
}

func generateTaskQueue(terminalType, sid string, priorityLevel int) string {
	return fmt.Sprintf("hg_scanner:%s:%s:%d", terminalType, sid, priorityLevel)
}
