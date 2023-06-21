package hg_service

import (
	"fmt"
	"runtime"
	"time"

	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/redis"
)

var (
	messageQueue *MessageQueue
)

type Msg struct {
	TaskId   string
	Priority int
	Payload  qmap.QM
}

const (
	HEADER_SESSION = "X-SESSION-ID" // 会话id
	// 终端类型
	TERMINAL_TYPE_CLIENT = "client" // 客户端
	TERMINAL_TYPE_WEB    = "web"    // 网页端

	MSG_EXPIRE_TIME = time.Minute * 60 // 队列消息过期时间

	SUBSCRIBE_EXPIRE_TIME = 86400 // 订阅终端过期清理时间（秒）
)

func GetScanManager() *MessageQueue {
	if messageQueue == nil {
		messageQueue = &MessageQueue{
			WebSubscribes:    qmap.QM{},
			WebPushedMsg:     make(chan *Msg, 10),
			ClientSubscribes: qmap.QM{},
			ClientPushedMsg:  make(chan *Msg, 10),
			CloseChan:        make(chan bool),
		}
	}
	return messageQueue
}

// 订阅消息终端
func SubTerminal(terminal *ScanMsgTerminal) {
	scanManager := GetScanManager()
	if terminal.TerminalType == TERMINAL_TYPE_CLIENT {
		sub := &MessageSub{
			TaskId:            terminal.TaskId,
			LowPriorityQueue:  []string{terminal.LowPriorityQueue},
			HighPriorityQueue: []string{terminal.HighPriorityQueue},
			ActiveTime:        time.Now().Unix(),
		}
		scanManager.ClientSubscribes[terminal.TaskId] = sub
	} else if terminal.TerminalType == TERMINAL_TYPE_WEB {
		if val, has := scanManager.WebSubscribes.TryInterface(terminal.TaskId); has {
			sub := val.(*MessageSub)
			sub.LowPriorityQueue = append(sub.LowPriorityQueue, terminal.LowPriorityQueue)
			sub.HighPriorityQueue = append(sub.HighPriorityQueue, terminal.HighPriorityQueue)
			sub.ActiveTime = time.Now().Unix()
		} else {
			sub := &MessageSub{
				TaskId:            terminal.TaskId,
				LowPriorityQueue:  []string{terminal.LowPriorityQueue},
				HighPriorityQueue: []string{terminal.HighPriorityQueue},
				ActiveTime:        time.Now().Unix(),
			}
			scanManager.WebSubscribes[terminal.TaskId] = sub
		}
	}
}

// 接收需要派发给其他端的消息
func PushMsgToTerminal(toTerminalType, taskId string, priority int, payload qmap.QM) {
	scanManager := GetScanManager()
	msg := &Msg{
		TaskId:   taskId,
		Priority: priority,
		Payload:  payload,
	}
	if toTerminalType == TERMINAL_TYPE_CLIENT {
		scanManager.ClientPushedMsg <- msg
	} else if toTerminalType == TERMINAL_TYPE_WEB {
		scanManager.WebPushedMsg <- msg
	}
}

type MessageQueue struct {
	ClientSubscribes qmap.QM
	ClientPushedMsg  chan *Msg // 扫描器端推送的消息队列
	WebSubscribes    qmap.QM
	WebPushedMsg     chan *Msg // 网页端推送的消息队列
	CloseChan        chan bool
}

func (this *MessageQueue) Run() {
	go this.dispatchMsg()
	go this.regularClearSubTerminal()
}

func (this MessageQueue) Close() {
	close(this.CloseChan)
}

func (this *MessageQueue) dispatchMsg() {
	redis := new(redis.Redis_service)
	for {
		this.dispatchMsgCirculate(redis)
	}
}

func (this *MessageQueue) dispatchMsgCirculate(redis *redis.Redis_service) {
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
			logMessage := fmt.Sprintf("Trace: %s\n", err)
			logMessage += fmt.Sprintf("%s", stacktrace)
			TerminalReceiveLog("dispatch_message", "", "", logMessage)
		}
	}()
	select {
	case <-this.CloseChan:
		return
	case msg := <-this.WebPushedMsg:
		if val, has := this.WebSubscribes.TryInterface(msg.TaskId); has {
			sub := val.(*MessageSub)
			sub.ActiveTime = time.Now().Unix()
			if msg.Priority == PRIORITY_QUEUE_LOW {
				for _, queue := range sub.LowPriorityQueue {
					redis.LPush(queue, msg.Payload.ToString())
					redis.Expire(queue, MSG_EXPIRE_TIME)
				}
			} else {
				for _, queue := range sub.HighPriorityQueue {
					redis.LPush(queue, msg.Payload.ToString())
					redis.Expire(queue, MSG_EXPIRE_TIME)
				}
			}
		}
	case msg := <-this.ClientPushedMsg:
		if val, has := this.ClientSubscribes.TryInterface(msg.TaskId); has {
			sub := val.(*MessageSub)
			sub.ActiveTime = time.Now().Unix()
			if msg.Priority == PRIORITY_QUEUE_LOW {
				for _, queue := range sub.LowPriorityQueue {
					redis.LPush(queue, msg.Payload.ToString())
					redis.Expire(queue, MSG_EXPIRE_TIME)
				}
			} else {
				for _, queue := range sub.HighPriorityQueue {
					redis.LPush(queue, msg.Payload.ToString())
					redis.Expire(queue, MSG_EXPIRE_TIME)
				}
			}
		}
	}
}

type MessageSub struct {
	TaskId            string
	LowPriorityQueue  []string
	HighPriorityQueue []string
	ActiveTime        int64
}

// 定期清理订阅的客户端
func (this *MessageQueue) regularClearSubTerminal() {
	for {
		now := time.Now().Unix()
		for key, sub := range this.ClientSubscribes {
			if sub.(*MessageSub).ActiveTime+SUBSCRIBE_EXPIRE_TIME < now {
				delete(this.ClientSubscribes, key)
			}
		}
		for key, sub := range this.WebSubscribes {
			if sub.(*MessageSub).ActiveTime+SUBSCRIBE_EXPIRE_TIME < now {
				delete(this.WebSubscribes, key)
			}
		}
		<-time.After(time.Second * 60)
	}
}
