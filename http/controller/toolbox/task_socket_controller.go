package toolbox

import (
	"errors"
	"fmt"
	"net/http"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/hg_service"
	"skygo_detection/logic/toolbox"
	"skygo_detection/mysql_model"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type TaskSocketController struct{}

var WEB_SOCKET = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 新增如下代码,解决跨域问题,即403错误
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// 启动检测系统
func (t TaskSocketController) Terminal(ctx *gin.Context) {
	isReconnect := true
	sessionId := ctx.GetHeader(hg_service.HEADER_SESSION)
	if sessionId == "" {
		isReconnect = false
		sessionId = t.generateSessionId(hg_service.TERMINAL_TYPE_WEB)
	}
	taskId := request.QueryString(ctx, "task_id")
	if taskId == "" {
		return
	}
	fmt.Println("taskId:", taskId)

	task, err := mysql_model.GetTaskByUuid(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	if task.Id == 0 {
		_ = new(mysql_model.PrivacyAnalysisLog).SetPrivacyLog(task.Id, "隐私任务任务不存在")
		response.RenderFailure(ctx, errors.New("任务不存在"))
		return
	}

	if task.Status == 2 {
		_ = new(mysql_model.PrivacyAnalysisLog).SetPrivacyLog(task.Id, "隐私任务已完成")
		response.RenderFailure(ctx, errors.New("任务已经完成"))
		return
	}

	// 升级为长连接
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	ws, err := WEB_SOCKET.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		_ = new(mysql_model.PrivacyAnalysisLog).SetPrivacyLog(task.Id, "隐私任务建立长连接失败"+err.Error())
		fmt.Println("err:", err)
		return
	}
	_ = new(mysql_model.PrivacyAnalysisLog).SetPrivacyLog(task.Id, "隐私任务建立长连接成功")

	defer ws.Close()

	// 隐私任务客户端
	taskTerminal := toolbox.NewTerminal(sessionId, taskId, ws, isReconnect)
	taskTerminal.Run()

	<-taskTerminal.CloseChan
}

func (t TaskSocketController) generateSessionId(terminalType string) string {
	return fmt.Sprintf("%s-%d", terminalType, custom_util.GetCurrentMilliSecond())
}
