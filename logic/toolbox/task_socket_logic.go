package toolbox

import (
	"time"

	"skygo_detection/guardian/src/net/qmap"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/redis"
	"skygo_detection/service"
)

const (
	MAX_RECONNECT_TIMES = 3 // 最大重连次数
	RECONNECT_TIME      = 3 // 失败重连时间间隔(秒)

	// 消息类型
	MSG_TYPE_SESSION   = "session"           // 推送会话
	MSG_TYPE_RESESSION = "reconnect_session" // 重连会话

	MSG_TYPE_START = "start_case_block" // 开始测试
	MSG_TYPE_STOP  = "error_handle"     // 结束测试
)

func NewTerminal(sid, taskId string, ws *websocket.Conn, isReconnect bool) *TaskTerminal {
	tt := TaskTerminal{
		StartChan:   make(chan bool),
		StopChan:    make(chan bool),
		CloseChan:   make(chan bool),
		Sid:         sid,
		TaskId:      taskId,
		Websocket:   ws,
		IsReconnect: isReconnect,
	}
	return &tt
}

type TaskTerminal struct {
	Websocket   *websocket.Conn
	StartChan   chan bool
	StopChan    chan bool
	CloseChan   chan bool
	Sid         string // 长连接通讯的会话id，通过该id实现会话重连
	TaskId      string // 检测任务id
	IsReconnect bool   // 是否是重连客户端
}

// 开始进行长连接通讯
func (t *TaskTerminal) Run() {
	t.pushSessionInfo()          // 在启动前，需要进行消息初始化
	go t.monitorStartFromRedis() // 监控redis 开始测试
	go t.monitorStopFromRedis()  // 监控redis 结束测试
	go t.monitorReceiveMsg()     // 监控从终端推送过来的消息
	go t.pushMsgToTerminal()     // 将消息推送给终端
}

// 推送会话消息
func (t *TaskTerminal) pushSessionInfo() error {
	msg := qmap.QM{
		"sequence": custom_util.GetCurrentMilliSecond(),
		"data": qmap.QM{
			"session_id": t.Sid,
		},
	}
	if t.IsReconnect {
		msg["type"] = MSG_TYPE_RESESSION
	} else {
		msg["type"] = MSG_TYPE_SESSION
	}

	return t.Websocket.WriteMessage(websocket.TextMessage, []byte(msg.ToString()))
}

// 监控从终端推送过来的消息
func (t *TaskTerminal) monitorReceiveMsg() {
	for {
		select {
		case <-t.CloseChan:
			t.Websocket.Close()
			return
		default:
			// 循环监听从客户端推送过来的消息
			_, _, err := t.Websocket.ReadMessage()
			if err != nil {
				// 如果接收消息异常，则停止长连接通讯
				t.Stop()
				t.Websocket.Close()
				return
			}
		}
	}
}

// 停止长连接
func (t *TaskTerminal) Stop() {
	close(t.CloseChan)
}

// 监控redis队列，开始测试
func (t *TaskTerminal) monitorStartFromRedis() {
	redis := new(redis.Redis_service)
	key := common.PRIVACY_TASK_START
	for {
		select {
		case <-t.CloseChan:
			return
		default:
			if len := redis.LLen(key); len.Val() > 0 {
				list := redis.LRange(key, 0, -1).Val()
				for _, v := range list {
					if v == t.TaskId {
						t.StartChan <- true
						if err := redis.LRem(key, 0, v).Err(); err != nil {
							TerminalSendLog("", t.TaskId, "开始测试，删除redis元素出错", err.Error())
						}
						break
					}
				}
			} else {
				<-time.After(time.Second)
			}
		}
	}
}

// 监控redis队列，结束测试
func (t *TaskTerminal) monitorStopFromRedis() {
	redis := new(redis.Redis_service)
	key := common.PRIVACY_TASK_STOP
	for {
		select {
		case <-t.CloseChan:
			return
		default:
			if len := redis.LLen(key); len.Val() > 0 {
				list := redis.LRange(key, 0, -1).Val()
				for _, v := range list {
					if v == t.TaskId {
						t.StopChan <- true
						if err := redis.LRem(key, 0, v).Err(); err != nil {
							TerminalSendLog("", t.TaskId, "停止测试，删除redis元素出错", err.Error())
						}
						break
					}
				}
			} else {
				<-time.After(time.Second)
			}
		}
	}
}

// 向终端推送消息
func (t *TaskTerminal) pushMsgToTerminal() {
	for {
		select {
		case <-t.CloseChan:
			t.Websocket.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "close"))
			close(t.StartChan)
			close(t.StopChan)
			return
		case <-t.StartChan:
			t.sendMessage(MSG_TYPE_START)
		case <-t.StopChan:
			t.sendMessage(MSG_TYPE_STOP)
		}
	}
}

func (t *TaskTerminal) sendMessage(msgType string) {
	msg := qmap.QM{
		"sequence": custom_util.GetCurrentMilliSecond(),
		"data": qmap.QM{
			"session_id": t.Sid,
		},
	}
	msg["type"] = msgType
	strMsg := msg.ToString()
	failTimes := 0 // 推送失败次数，推送失败超过次数限制，任务停止
	for {
		if err := t.Websocket.WriteMessage(websocket.TextMessage, []byte(strMsg)); err == nil {
			TerminalSendLog("", t.TaskId, strMsg, "")
			break
		} else {
			failTimes++
			// 如果重试次数超过最大限度，则停止任务
			if failTimes >= MAX_RECONNECT_TIMES {
				t.Stop()
			}
			if <-t.CloseChan {
				close(t.StartChan)
				close(t.StopChan)
				return
			}
			// 推送失败，等会儿再重试
			<-time.After(time.Second * RECONNECT_TIME)
		}
	}
}

func TerminalSendLog(messageType, taskId, rawMsg, errMsg string) {
	logger := service.GetDefaultLogger("terminal_message_send")
	defer logger.Sync()
	errLevel := zapcore.InfoLevel
	if errMsg != "" {
		errLevel = zapcore.ErrorLevel
	}
	if logger.Core().Enabled(errLevel) {
		logger.Check(errLevel, errMsg).Write(
			zap.String("task_id", taskId),
			zap.String("message_type", messageType),
			zap.String("raw_message", rawMsg),
			zap.String("create_time", time.Now().Format(time.RFC3339)),
		)
	}
}
