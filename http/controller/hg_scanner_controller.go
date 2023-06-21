package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"net/http"
	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/custom_util/clog"
	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/hg_service"
	"skygo_detection/logic/hg_scanner_logic"
	"skygo_detection/mysql_model"
)

type ScanController struct{}

var WEB_SOCKET = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// 新增如下代码,解决跨域问题,即403错误
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
} // use default options
var connCount = 0

/**
 * apiType http
 * @api {get} /message/v1/hg_scanner/terminal 合规扫描-扫描器长连接
 * @apiVersion 1.0.0
 * @apiName Terminal
 * @apiGroup HgScanner
 *
 * @apiDescription 合规扫描-扫描器长连接
 *
 * @apiHeader {string} 		[X-SESSION-ID]   	会话id（如果是重连，需要传重连前的会话id）
 *
 * @apiParam {string}   	task_id  			扫描任务id
 *
 * @apiExample {curl} 请求示例:
 * curl -i ws://qa.vadmin.car.qihoo.net:8100/message/v1/hg_scanner/terminal?task_id=test
 *
 */
func (this ScanController) Terminal(ctx *gin.Context) {
	isReconnect := false
	if sid := ctx.GetHeader(hg_service.HEADER_SESSION); sid != "" {
		isReconnect = true
	}
	sessionId := this.generateSessionId(hg_service.TERMINAL_TYPE_CLIENT)
	taskId := request.QueryString(ctx, "task_id")
	if taskId == "" {
		return
	}
	// 升级为长连接
	ws, err := WEB_SOCKET.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()
	// 新建扫描任务消息客户端
	scanMsgClient := hg_service.NewScanMsgTerminal(hg_service.TERMINAL_TYPE_CLIENT, sessionId, taskId, ws, isReconnect)
	// 开始进行长连接消息的接受和推送
	go scanMsgClient.Run()
	<-scanMsgClient.CloseChan
}

/**
 * apiType http
 * @api {get} /message/v1/hg_scanner/web 合规扫描-Web长连接
 * @apiVersion 1.0.0
 * @apiName Web
 * @apiGroup HgScanner
 *
 * @apiDescription 合规扫描-Web长连接
 *
 * @apiHeader {string} 		[X-SESSION-ID]   	会话id（如果是重连，需要传重连前的会话id）
 *
 * @apiParam {string}   	task_id  			扫描任务id
 *
 * @apiExample {curl} 请求示例:
 * curl -i ws://qa.vadmin.car.qihoo.net:8100/message/v1/hg_scanner/web?task_id=test
 *
 */
func (this ScanController) Web(ctx *gin.Context) {
	isReconnect := true
	sessionId := ctx.GetHeader(hg_service.HEADER_SESSION)
	if sessionId == "" {
		isReconnect = false
		sessionId = this.generateSessionId(hg_service.TERMINAL_TYPE_WEB)
	}
	taskId := request.QueryString(ctx, "task_id")
	if taskId == "" {
		return
	}
	// 升级为长连接
	ctx.Writer.Header().Set("Access-Control-Allow-Origin", "*")
	ctx.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
	ctx.Writer.Header().Set("Access-Control-Allow-Headers", "*")
	ws, err := WEB_SOCKET.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		return
	}
	defer ws.Close()

	// 新建扫描任务消息客户端
	scanMsgClient := hg_service.NewScanMsgTerminal(hg_service.TERMINAL_TYPE_WEB, sessionId, taskId, ws, isReconnect)
	scanMsgClient.Run()

	<-scanMsgClient.CloseChan
}

func (this ScanController) generateSessionId(terminalType string) string {
	return fmt.Sprintf("%s-%d", terminalType, custom_util.GetCurrentMilliSecond())
}

func (this ScanController) DownloadCase(ctx *gin.Context) {
	taskUuid := request.QueryString(ctx, "task_id")

	// 获取脚本
	testScriptSlice, err := hg_scanner_logic.GetTestScript(taskUuid)
	if err != nil {
		clog.Error("GetTestScript  Err", zap.Any("error", err))
		response.RenderFailure(ctx, err)
		return
	}

	// 打包文件
	b, err := hg_scanner_logic.ArchiveFile(testScriptSlice)
	if err != nil {
		clog.Error("ArchiveFile  Err", zap.Any("error", err))
		response.RenderFailure(ctx, err)
		return
	}

	// 文件名
	zipFileName := fmt.Sprintf("%s_%d.zip", taskUuid, len(testScriptSlice))

	// 文件传输
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", zipFileName))
	ctx.Header("Content-Type", "*")
	ctx.Header("Accept-Length", fmt.Sprintf("%d", len(b.Bytes())))
	_, err = ctx.Writer.Write(b.Bytes())
	if err != nil {
		clog.Error("Ctx Writer Err", zap.Any("error", err))
		response.RenderFailure(ctx, err)
		return
	}
}

/**
 * apiType http
 * @api {post} /message/v1/hg_scanner/upload 合规检测文件上传
 * @apiVersion 1.0.0
 * @apiName Upload
 * @apiGroup HgScanner
 *
 * @apiDescription 合规检测文件上传
 *
 * @apiUse authHeader
 *
 * @apiParam {string} 	[file_name]       	文件名称
 * @apiParam {file}		file 				文件
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/message/v1/hg_scanner/upload
 *
 * @apiSuccessExample {json} 请求成功示例:
 *  {
 *		"code":0,
 *		"msg":"",
 *		"data":{
 *			"file_id":"a834qafmcxvadfq1123"
 *		}
 *  }
 */
func (this ScanController) Upload(ctx *gin.Context) {
	fileName := ctx.Request.FormValue("file_name")
	if file, header, err := ctx.Request.FormFile("file"); err == nil {
		if fileName == "" {
			fileName = header.Filename
		}
		// 文件内容
		fileContent := make([]byte, header.Size)
		num, err := file.Read(fileContent)

		if err != nil {
			response.RenderFailure(ctx, errors.New("文件处理失败"))
			return
		}
		if num <= 0 {
			response.RenderFailure(ctx, errors.New("您上传的文件为空文件！"))
			return
		}

		if fileId, err := mongo.GridFSUpload(common.MC_PROJECT, fileName, fileContent); err != nil {
			panic(err)
		} else {
			m := gin.H{
				"code": 0,
				"data": gin.H{"file_id": fileId},
			}
			ctx.AbortWithStatusJSON(200, m)
			return
		}
	} else {
		panic(err)
	}
}

/**
 * apiType http
 * @api {get} /message/v1/hg_scanner/terminal_info 获取终端信息
 * @apiVersion 1.0.0
 * @apiName GetTerminalInfo
 * @apiGroup HgScanner
 *
 * @apiDescription 获取终端信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}		task_id 	任务id
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/message/v1/hg_scanner/terminal_info?task_id=xxxxc
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "cpu": "64",
 *         "id": 12,
 *         "last_connect_time": 1635143320,
 *         "name": "合规子任务_bhjkjk55",
 *         "os_type": "android",
 *         "os_version": "9.0",
 *         "status": "create",
 *         "task_uuid": "G3JWRRCT"
 *     },
 *     "msg": ""
 * }
 */
func (this ScanController) GetTerminalInfo(ctx *gin.Context) {
	taskId := request.QueryString(ctx, "task_id")
	if task, err := new(mysql_model.HgTestTask).FindOne(taskId); err == nil {
		result := qmap.QM{
			"id":                task.Id,
			"cpu":               task.Cpu,
			"last_connect_time": task.LastConnectTime,
			"name":              task.Name,
			"os_type":           task.OsType,
			"os_version":        task.OsVersion,
			"status":            task.Status,
			"task_uuid":         task.TaskUuid,
		}
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}
