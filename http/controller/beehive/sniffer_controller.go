package beehive

import (
	"errors"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/logic/beehive"
	"skygo_detection/mysql_model"
	"strings"

	"github.com/gin-gonic/gin"
)

type SnifferController struct{}

// 开始扫描
func (s SnifferController) Scan(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	snifferLogic := new(beehive.Sniffer)
	_, err := snifferLogic.GetTaskById(ctx, taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	channel := request.MustInt(ctx, "channel")
	if err := snifferLogic.Scan(ctx, taskId, channel); err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, nil)
}

// 重新扫描
func (s SnifferController) ReScan(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	snifferLogic := new(beehive.Sniffer)
	_, err := snifferLogic.GetTaskById(ctx, taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	if err := snifferLogic.ReScan(ctx, taskId); err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, nil)
}

// 轮询获取频点，定时器每10s请求一次
func (s SnifferController) GetFreq(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	snifferLogic := new(beehive.Sniffer)
	data, err := snifferLogic.GetFreq(ctx, taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, data)
}

// 取倒计时的剩余时间
func (s SnifferController) GetTimer(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	snifferLogic := new(beehive.Sniffer)
	data, err := snifferLogic.GetTimer(ctx, taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, data)
}

// 开始嗅探
func (s SnifferController) StartSniff(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	freq := request.MustString(ctx, "freq")
	mode := request.MustString(ctx, "mode")
	snifferLogic := new(beehive.Sniffer)
	if mode == "" {
		response.RenderFailure(ctx, errors.New("mode参数的值不能为空"))
		return
	}
	if freq == "" {
		response.RenderFailure(ctx, errors.New("freq参数的值不能为空"))
		return
	}

	freq = strings.Replace(freq, "M", "", -1)
	freq = strings.Replace(freq, "m", "", -1)

	bool, err := snifferLogic.StartSniff(ctx, taskId, freq, mode)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	if !bool {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, nil)
}

// 取imsi数据 10s一次
func (s SnifferController) GetImsi(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	snifferLogic := new(beehive.Sniffer)
	num, err := snifferLogic.GetImsi(ctx, taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	stop := snifferLogic.IfStopScan(ctx, taskId)
	data := map[string]interface{}{"num": num, "stop": stop}
	response.RenderSuccess(ctx, data)
}

// 取sms数据，20s一次
func (s SnifferController) GetSms(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	snifferLogic := new(beehive.Sniffer)
	num, err := snifferLogic.GetSms(ctx, taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	stop := snifferLogic.IfStopScan(ctx, taskId)
	data := map[string]interface{}{"num": num, "stop": stop}
	response.RenderSuccess(ctx, data)
}

func (s SnifferController) ImsiSmsCount(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	snifferLogic := new(beehive.Sniffer)
	imsiCount, err := snifferLogic.GetImsiCount(ctx, taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	smsCount, err := snifferLogic.GetSmsCount(ctx, taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	data := map[string]int{"imsi_count": imsiCount, "sms_count": smsCount}
	response.RenderSuccess(ctx, data)
}

func (s SnifferController) DelImsi(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	ids := request.MustSlice(ctx, "ids")
	if len(ids) > 0 {
		snifferLogic := new(beehive.Sniffer)
		_, err := snifferLogic.DelImsi(ctx, taskId, ids)
		if err != nil {
			response.RenderFailure(ctx, err)
			return
		}
	}
	response.RenderSuccess(ctx, nil)
}

func (s SnifferController) DelSms(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	ids := request.MustSlice(ctx, "ids")
	if len(ids) > 0 {
		snifferLogic := new(beehive.Sniffer)
		_, err := snifferLogic.DelSms(ctx, taskId, ids)
		if err != nil {
			response.RenderFailure(ctx, err)
			return
		}
	}
	response.RenderSuccess(ctx, nil)
}

func (s SnifferController) ImsiList(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	taskId := request.QueryInt(ctx, "task_id")
	if taskId < 1 {
		response.RenderFailure(ctx, errors.New("task_id不能为空"))
		return
	}
	session := mysql.GetSession()
	session = session.Where("task_id=?", taskId).Where("delete_time=0")

	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	widget.AddSorter(*(orm.NewSorter("id", 1)))
	all := widget.PaginatorFind(session, &[]mysql_model.BeehiveGsmSnifferImsi{})
	response.RenderSuccess(ctx, all)
}

func (s SnifferController) SmsList(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	taskId := request.QueryInt(ctx, "task_id")
	if taskId < 1 {
		response.RenderFailure(ctx, errors.New("task_id不能为空"))
		return
	}
	session := mysql.GetSession()
	session = session.Where("task_id=?", taskId).Where("delete_time=0")

	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	widget.AddSorter(*(orm.NewSorter("id", 1)))
	all := widget.PaginatorFind(session, &[]mysql_model.BeehiveGsmSnifferSms{})
	response.RenderSuccess(ctx, all)
}

// 停止扫描
func (s SnifferController) StopScan(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	snifferLogic := new(beehive.Sniffer)
	_, err := snifferLogic.GetTaskById(ctx, taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	if err := snifferLogic.StopScan(ctx, taskId); err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, nil)
}

// 停止嗅探
func (s SnifferController) StopSniff(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	snifferLogic := new(beehive.Sniffer)
	_, err := snifferLogic.GetTaskById(ctx, taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	if err := snifferLogic.StopSniff(ctx, taskId); err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, nil)
}

// 关闭系统
func (s SnifferController) Close(ctx *gin.Context) {
	taskId := request.MustInt(ctx, "task_id")
	snifferLogic := new(beehive.Sniffer)
	_, err := snifferLogic.GetTaskById(ctx, taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	if err := snifferLogic.Close(ctx, taskId); err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, nil)
}

func (s SnifferController) Detail(ctx *gin.Context) {
	taskId := request.QueryInt(ctx, "task_id")
	if taskId < 1 {
		response.RenderFailure(ctx, errors.New("task_id不能为空"))
		return
	}
	snifferModel := mysql_model.BeehiveGsmSniffer{}
	_, err := snifferModel.FindByTaskId(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, snifferModel)
}
