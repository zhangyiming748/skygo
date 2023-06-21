package gps

import (
	"encoding/json"
	"errors"
	"skygo_detection/common"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mysql_model"

	"skygo_detection/logic/gps"

	"github.com/gin-gonic/gin"
)

type MotionController struct{}

// 设置gps时间
func (t MotionController) SetTime(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	time := request.MustString(ctx, "time")

	data := map[string]int{
		"id": 0,
	}
	_, err := gps.GetTaskById(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}

	id, err := gps.SetTime(taskId, time)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	data["id"] = id

	response.RenderSuccess(ctx, data)
}

// 启动设备
func (t MotionController) Start(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	gt, err := gps.GetGpsTaskByTaskId(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	if gt.Status == gps.DEVICE_RUNNING {
		response.RenderFailure(ctx, errors.New("任务已经启动了"))
		return
	}
	if gt.Time == "" || gt.Time == "0000-00-00 00:00:00" {
		response.RenderFailure(ctx, errors.New("请设置GPS时间"))
		return
	}
	gps.Start(taskId, gt.Time)
	response.RenderSuccess(ctx, nil)
}

// 关闭系统
func (t MotionController) Close(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	_, err := gps.GetGpsTaskByTaskId(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	gps.Close(taskId)
	response.RenderSuccess(ctx, nil)
}

// 查询设备启动结果
func (t MotionController) Status(ctx *gin.Context) {
	taskId := request.ParamInt(ctx, "task_id")
	gt, err := gps.GetGpsTaskByTaskId(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	data := map[string]int{
		"status": 0,
	}
	data["status"] = gt.Status
	response.RenderSuccess(ctx, data)
}

func (t MotionController) Online(ctx *gin.Context) {
	taskId := request.ParamInt(ctx, "task_id")
	_, err := gps.GetGpsTaskByTaskId(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}

	data := map[string]bool{
		"status": false,
	}
	status, err := gps.Online()
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}

	data["status"] = status
	response.RenderSuccess(ctx, data)
}

// 生成路线
func (t MotionController) Line(ctx *gin.Context) {
	gs := mysql_model.GpsSearch{}
	taskId := request.MustInt(ctx, "task_id")
	templateId := request.MustInt(ctx, "template_id")
	start := request.MustString(ctx, "start")
	mid := request.MustSlice(ctx, "middle") // json
	if len(mid) > 0 {
		middle, err := json.Marshal(mid)
		if err != nil {
			response.RenderFailure(ctx, err)
			return
		}
		gs.Middle = string(middle)
	}
	end := request.MustString(ctx, "end")
	gs.TaskId = taskId
	gs.Type = common.GPS_TYPE_MOTION
	gs.TemplateId = templateId
	gs.Start = start
	gs.End = end
	req := request.MustSlice(ctx, "req") // json
	if len(req) > 0 {
		reqstr, err := json.Marshal(req)
		if err != nil {
			response.RenderFailure(ctx, err)
			return
		}
		gs.Req = string(reqstr)
	}

	id, err := gps.CreateLine(gs)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	data := map[string]int{
		"id": id,
	}
	response.RenderSuccess(ctx, data)
}

func (t MotionController) LineHistory(ctx *gin.Context) {
	taskId := request.QueryInt(ctx, "task_id")
	templateId := request.QueryInt(ctx, "template_id")
	data, err := gps.LineHistory(taskId, templateId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, data)
}

// 开始欺骗
func (t MotionController) Cheat(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	searchId := request.MustInt(ctx, "search_id") //路线id
	_, err := gps.Cheat(taskId, searchId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, nil)
}

// 获取最新的一条欺骗记录
func (t MotionController) GetLatestCheat(ctx *gin.Context) {
	taskId := request.ParamInt(ctx, "task_id")
	gc, err := gps.GetLatestCheatByTaskId(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, gc)
}
