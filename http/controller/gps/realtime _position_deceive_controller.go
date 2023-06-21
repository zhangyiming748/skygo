package gps

import (
	"github.com/gin-gonic/gin"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/logic/gps"
	"skygo_detection/mysql_model"
	"strconv"
	"time"
)

const (
	TASK_TYPE_REALTIME = 1
)

type RealtimePositionDeceiveController struct{}

// 开启实时位置欺骗任务
func (this RealtimePositionDeceiveController) StartDeceiveTask(ctx *gin.Context) {
	taskId := request.ParamInt(ctx, "task_id")
	gps.SetLogs(taskId, "开启实时位置欺骗任务")
	response.RenderSuccess(ctx, "success")
}

// 搜索（保存数据）
func (this RealtimePositionDeceiveController) Create(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	module := new(mysql_model.GpsSearch)
	module.TaskId = req.MustInt("task_id")
	module.Type = TASK_TYPE_REALTIME
	module.Start = req.MustString("start")
	module.Lng = req.Float32("lng")
	module.Lat = req.Float32("lat")
	module.CreateTime = time.Now().Format("2006-01-02 15:04:05")

	sLng := strconv.FormatFloat(float64(module.Lng), 'f', 6, 64)
	sLat := strconv.FormatFloat(float64(module.Lat), 'f', 6, 64)
	gps.SetLogs(module.TaskId, "设置实时位置为："+module.Start+"，经纬度："+
		sLng+", "+sLat+"。")
	if _, err := module.Create(); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, "success")
	}
}

// 开始欺骗
func (this RealtimePositionDeceiveController) StartDeceive(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	taskId := req.MustInt("task_id")
	start := req.MustString("start")
	lng := req.Float32("lng")
	lat := req.Float32("lat")
	res, err := gps.StartRealtimePositionDeceive(taskId, start, lng, lat)
	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, res)
	}
}

// 停止欺骗
func (this RealtimePositionDeceiveController) StopDeceive(ctx *gin.Context) {
	taskId := request.ParamInt(ctx, "task_id")
	res, err := gps.StopRealtimePositionDeceive(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, res)
	}
}

// 查询最近5条数据
func (this RealtimePositionDeceiveController) GetFive(ctx *gin.Context) {
	tid := ctx.Param("task_id")
	module := new(mysql_model.GpsSearch)
	res, err := module.GetFive(tid)
	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, res)
	}
}

// GPS欺骗记录
func (this RealtimePositionDeceiveController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	taskId := request.QueryInt(ctx, "task_id")

	s := mysql.GetSession()
	s.Where("task_id=?", taskId)

	// 查询组键
	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	all := widget.PaginatorFind(s, &[]mysql_model.GpsCheat{})
	response.RenderSuccess(ctx, all)
}
