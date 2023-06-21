package main

import (
	"flag"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/websocket"
	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/guardian/util"

	"skygo_detection/lib/hg_service"
)

var (
	addr   = flag.String("addr", "127.0.0.1:82", "http service address")
	taskId = flag.String("tid", "G4KB2M2J", "task id")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	u := url.URL{Scheme: "ws", Host: *addr, Path: "/message/v1/hg_scanner/terminal", RawQuery: "task_id=" + *taskId}
	log.Printf("connecting to %s", u.String())
	header := http.Header{}
	// header.Add(hg_service.HEADER_SESSION, "client-1620620455891")
	c, _, err := websocket.DefaultDialer.Dial(u.String(), header)

	if err != nil {
		log.Fatal("dial:", err)
	}
	defer c.Close()

	done := make(chan struct{})

	go func() {
		for {
			select {
			case <-done:
				return
			default:
				_, message, err := c.ReadMessage()
				if err != nil {
					log.Println("read:", err)
				} else {
					println("receive:", string(message))
				}
				// if res, err := qmap.NewWithString(string(message)); err == nil {
				// 	receiveMSH(c, res)
				// }
			}
			<-time.After(time.Second * 3)
		}
	}()
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	msg := qmap.QM{
		"type": "heart_beat",
	}
	for {
		select {
		case <-interrupt:
			log.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.Println("write close:", err)
				return
			}
			close(done)
			return
		default:
			<-time.After(30 * time.Second)
			println("heart_beat")
			msg["sequence"] = time.Now().Unix()
			err := c.WriteMessage(websocket.TextMessage, []byte(msg.ToString()))
			if err != nil {
				log.Println("write:", err)
				return
			}
		}
	}
}

func receiveMSH(ws *websocket.Conn, msg qmap.QM) {
	sequence := msg.Int64("sequence")
	msgType := msg.String("type")
	var err error
	switch msgType {
	case hg_service.MSG_TYPE_TERMINAL_INFO:
		// 接收到获取终端信息命令
		err = pushTerminalInfoReply(ws)
	case hg_service.MSG_TYPE_TERMINAL_UPDATE_CASES:
		// 接收到更新测试用例命令
		err = pushTerminalUpdateCaseReply(ws)
	case hg_service.MSG_TYPE_START_CASE:
		// 接收到初始化测试用例命令
		<-time.After(time.Second * 10)
		err = pushStartCaseReply(ws)
	case hg_service.MSG_TYPE_START_CASE_BLOCK:
		// 接收到执行测试用例block命令
		<-time.After(time.Second * 10)
		err = pushStartCaseBlockReply(ws)
	case hg_service.MSG_TYPE_END_CASE:
		// 接收到结束测试用例命令
		<-time.After(time.Second * 10)
		err = pushEndCaseReply(ws)

	}
	if err != nil {
		panic(err)
	}
	if msgType != hg_service.MSG_TYPE_ACK {
		receiveConfirm(sequence, ws)
	}
}

func pushTerminalInfoReply(ws *websocket.Conn) error {
	msg := qmap.QM{
		"sequence": util.GetCurrentMilliSecond(),
		"type":     hg_service.MSG_TYPE_TERMINAL_INFO_REPLY,
		"data": qmap.QM{
			"cpu":        "64",
			"os_type":    "android",
			"os_version": "9.0",
		},
	}
	return ws.WriteMessage(websocket.TextMessage, []byte(msg.ToString()))
}

/*
	{
		"type":"terminal_update_cases_reply",
		"sequence":12313,
		"data":{

		}
	}
*/
func pushTerminalUpdateCaseReply(ws *websocket.Conn) error {
	msg := qmap.QM{
		"sequence": util.GetCurrentMilliSecond(),
		"type":     hg_service.MSG_TYPE_TERMINAL_UPDATE_CASES_REPLY,
		"data": qmap.QM{
			"status": "success",
			"reason": "",
		},
	}
	return ws.WriteMessage(websocket.TextMessage, []byte(msg.ToString()))
}

/*
	{
		"type":"start_case_reply",
		"sequence":12313,
		"data":{
		  "status":"success/fail",
		  "reason":""
		}
	}
*/
func pushStartCaseReply(ws *websocket.Conn) error {
	msg := qmap.QM{
		"sequence": util.GetCurrentMilliSecond(),
		"type":     hg_service.MSG_TYPE_START_CASE_REPLY,
		"data": qmap.QM{
			"status": "success",
			"reason": "",
		},
	}
	return ws.WriteMessage(websocket.TextMessage, []byte(msg.ToString()))
}

/*
	{
	    "type":"start_case_block_reply",
	    "sequence":12313,
	    "data":{
	      "status":"success/fail",
	      "reason":"",
	      "case_id":"1212313",
	      "block_name":"runBlock",
	      "result":""//自动化扫描测试结果为：true/false，如果是复杂的则该字符串交给测试结果解析引擎分析
	    }
	}
*/
func pushStartCaseBlockReply(ws *websocket.Conn) error {
	msg := qmap.QM{
		"sequence": util.GetCurrentMilliSecond(),
		"type":     hg_service.MSG_TYPE_START_CASE_BLOCK_REPLY,
		"data": qmap.QM{
			"status":     "success",
			"reason":     "",
			"case_id":    "1212313",
			"block_name": "runBlock",
			"result":     "", // 自动化扫描测试结果为：true/false，如果是复杂的则该字符串交给测试结果解析引擎分析
		},
	}
	return ws.WriteMessage(websocket.TextMessage, []byte(msg.ToString()))
}

/*
	{
		"type":"end_case_reply",
		"sequence":12313,
		"data":{
		  "status":"success/fail",
		  "reason":""
		}
	}
*/
func pushEndCaseReply(ws *websocket.Conn) error {
	msg := qmap.QM{
		"sequence": util.GetCurrentMilliSecond(),
		"type":     hg_service.MSG_TYPE_END_CASE_REPLY,
		"data": qmap.QM{
			"status": "success/fail",
			"reason": "",
		},
	}
	return ws.WriteMessage(websocket.TextMessage, []byte(msg.ToString()))
}

func receiveConfirm(sequence int64, ws *websocket.Conn) error {
	msg := qmap.QM{
		"type":     hg_service.MSG_TYPE_ACK,
		"sequence": sequence,
	}
	return ws.WriteMessage(websocket.TextMessage, []byte(msg.ToString()))
}
