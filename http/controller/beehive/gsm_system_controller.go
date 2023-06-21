package beehive

import (
	"errors"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
	"skygo_detection/custom_util/blog"
	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	l "skygo_detection/logic/beehive"
	"skygo_detection/mysql_model"
)

type GsmSystemController struct{}

/**
 * apiType http
 * @api {post} /api/v1/beehive/gsm_system/config 选择参数
 * @apiVersion 1.0.0
 * @apiName config
 * @apiGroup gsm_system
 *
 * @apiDescription 选择设备参数配置
 *
 * @apiParam {int}           task_id                任务id
 * @apiParam {int}           config_id              默认配置1~4
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST -d task_id=任务id&config_id=配置 http://localhost/api/v1/beehive/gsm_system/config
 *
 * @apiParamExample {json}  请求参数示例:
 *      {
 *          "task_id":1,
 *          "config_id":1,
 *      }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "config_status":     "true",
 *         "config_message_id": 1
 *     },
 *     "msg":""
 * }
 */
func (this GsmSystemController) SetConfig(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	tid := req.MustInt("task_id")
	cid := req.MustInt("config_id")
	res, err := l.GsmSystemSetConfig(tid, cid)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, res)
}

/**
 * apiType http
 * @api {get} /api/v1/beehive/gsm_system/get_devices_info/:taskid 获取设备信息
 * @apiVersion 1.0.0
 * @apiName get_devices_info
 * @apiGroup gsm_system
 *
 * @apiDescription 获取扫描到的设备信息
 *
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X GET http://localhost/api/v1/beehive/gsm_system/get_devices_info/23
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "devices_status":     "true",
 *         "devices_message_id": 1,
 *         "devices": [
 *              [
 *                  "001010123456780",
 *                  "351615087961130",
 *                  "10000001",
 *                  "none"
 *              ],
 *              [
 *                  "001012333333333",
 *                  "355754071347990",
 *                  "10000002",
 *                  "192.168.99.2"
 *              ]
 *          ],
 *     "msg":""
 * }
 */
//获取设备信息
type terminal struct {
	Imei        string `json:"imei"`
	Imsi        string `json:"imsi"`
	PhoneNumber string `json:"phone_number"`
	Ip          string `json:"ip"`
}

func (this GsmSystemController) GetDevices(ctx *gin.Context) {
	tid := request.ParamInt(ctx, "task_id")
	q := ctx.Request.URL.RawQuery
	_, err := l.GsmSystemGetDevices(tid)
	if err != nil {
		response.RenderFailure(ctx, nil)
		return
	}

	s := mysql.GetSession()
	// 查询组键
	widget := orm.PWidget{}
	widget.SetQueryStr(q)
	all := widget.PaginatorFind(s, &[]mysql_model.BeehiveGsmSystemTty{})

	response.RenderSuccess(ctx, all)

}

/**
 * apiType http
 * @api {get} /api/v1/beehive/gsm_system/start/:task_id 启动系统
 * @apiVersion 1.0.0
 * @apiName start
 * @apiGroup gsm_system
 *
 * @apiDescription 启动系统
 *
 *
 * @apiParam {int}           task_id                任务id
 *
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X GET http://localhost/api/v1/beehive/gsm_system/start/23
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "start_status":     "true",
 *         "start_message_id": 1
 *     },
 *     "msg":""
 * }
 */
func (this GsmSystemController) StartSystem(ctx *gin.Context) {
	tid := request.ParamInt(ctx, "task_id")
	start, err := l.GsmSystemStart(tid)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, start)
}
func (this GsmSystemController) Start(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	tid := req.MustInt("task_id")
	cid := req.MustInt("config_id")
	start, err := l.GsmSystemSetConfig(tid, cid)
	if err != nil {
		response.RenderFailure(ctx, err)
		blog.Info("config error", zap.Any("error info", err))
		return
	}
	start, err = l.GsmSystemStart(tid)
	if err != nil {
		response.RenderFailure(ctx, err)
		blog.Info("start error", zap.Any("error info", err))
		return
	}

	response.RenderSuccess(ctx, start)
}

/**
 * apiType http
 * @api {get} /api/v1/beehive/gsm_system/get_sms/:task_id 获取短信
 * @apiVersion 1.0.0
 * @apiName get_sms
 * @apiGroup gsm_system
 *
 * @apiDescription 获取短信(分页)
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X GET http://localhost/api/v1/beehive/gsm_system/get_sms/23
 *
 * @apiSuccessExample {json} 请求成功示例:
 *  {
 *     "code": 0,
 *     "data": {
 *         "list": [
 *             {
 *                 "id": 3,
 *                 "recv_imsi": 0,
 *                 "recv_num": 132,
 *                 "recv_time": "22",
 *                 "send_num": 231,
 *                 "sms_content": "",
 *                 "task_id": 23
 *             },
 *             {
 *                 "id": 4,
 *                 "recv_imsi": 0,
 *                 "recv_num": 142,
 *                 "recv_time": "22",
 *                 "send_num": 0,
 *                 "sms_content": "",
 *                 "task_id": 23
 *             }
 *         ],
 *         "pagination": {
 *             "current_page": 1,
 *             "per_page": 20,
 *             "total": 2,
 *             "total_pages": 1
 *         }
 *     },
 *     "msg": ""
 * }
 */
// 短信选项卡中短信分页列表
func (this GsmSystemController) GetSMS(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	taskId := request.ParamInt(ctx, "task_id")
	s := mysql.GetSession()
	s.Where("task_id = ?", taskId)

	// 查询组键
	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	all := widget.PaginatorFind(s, &[]mysql_model.BeehiveGsmSystemSms{})
	response.RenderSuccess(ctx, all)
}

/**
* apiType http
* @api {get} api/v1/beehive/gsm_system/get_sms_num/:task_id 获取短信数量
* @apiVersion 1.0.0
* @apiName get_sms_num
* @apiGroup gsm_system
*
* @apiDescription 获取收到的短信总数填充角标
*
*
*
* @apiExample {curl} 请求示例:
* curl -i -X GET http://localhost/api/v1/beehive/gsm_system/get_sms_num/23
*
* @apiSuccessExample {json} 请求成功示例:
* {
*     "code": 0,
*     "data": {
*         "count": 32
*     },
*     "msg":""
* }
 */
// 任务详情页短信数量角标
func (this GsmSystemController) GetSMSNum(ctx *gin.Context) {
	tid := request.ParamInt(ctx, "task_id")
	log.Println(tid)
	num, err := l.GsmSystemGetSmsNum(tid)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, num)
}

/**
 * apiType http
 * @api {delete} /api/v1/beehive/gsm_system/del_sms 批量删除短信
 * @apiVersion 1.0.0
 * @apiName BulkDelete
 * @apiGroup gsm_system
 *
 * @apiDescription 批量删除短信
 *
 *
 *
 * @apiParam {array}           ids                短信id列表
 * @apiParamExample {json}  请求参数示例:
 * {
 *     "ids":[3,4,5]
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "number":3,
 *     },
 *     "msg":""
 * }
 */
// 批量删除短信
func (this GsmSystemController) BulkDeleteSMS(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	ids := req.SliceInt("ids")
	fail := l.GmsSystemBulkDeleteSMS(ids)
	response.RenderSuccess(ctx, qmap.QM{"success": len(ids) - fail})
}

/**
 * apiType http
 * @api {get} /api/v1/beehive/gsm_system/get_sms/:key 搜索短信
 * @apiVersion 1.0.0
 * @apiName search
 * @apiGroup gsm_system
 *
 * @apiDescription 搜索短信
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X GET http://localhost/api/v1/beehive/gsm_system/get_sms/23
 *
 * @apiSuccessExample {json} 请求成功示例:
 *  {
 *     "code": 0,
 *     "data": {
 *         "list": [
 *             {
 *                 "id": 3,
 *                 "recv_imsi": 0,
 *                 "recv_num": 132,
 *                 "recv_time": "22",
 *                 "send_num": 231,
 *                 "sms_content": "",
 *                 "task_id": 23
 *             },
 *             {
 *                 "id": 4,
 *                 "recv_imsi": 0,
 *                 "recv_num": 142,
 *                 "recv_time": "22",
 *                 "send_num": 0,
 *                 "sms_content": "",
 *                 "task_id": 23
 *             }
 *         ],
 *         "pagination": {
 *             "current_page": 1,
 *             "per_page": 20,
 *             "total": 2,
 *             "total_pages": 1
 *         }
 *     },
 *     "msg": ""
 * }
 */
// 搜索短信
func (this *GsmSystemController) GsmSystemSearch(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	key := request.QueryString(ctx, "key")
	tid := request.QueryInt(ctx, "task_id")
	if tid == 0 {
		response.RenderFailure(ctx, errors.New("must_tid"))
	}
	res := l.GsmSystemSearchSms(tid, key, queryParams)
	response.RenderSuccess(ctx, res)
}

/**
* apiType http
* @api {post} api/v1/beehive/gsm_system/send_sms 批量模拟短信
* @apiVersion 1.0.0
* @apiName send_sms
* @apiGroup gsm_system
*
* @apiDescription 批量模拟短信
*
*
*
* @apiParam {int}              task_id                任务id
* @apiParam {int}              send_num               发件人手机号
* @apiParam {array}            recv_nums              收件人手机号
* @apiParam {string}           sms_content            短信内容
*
* @apiExample {curl} 请求示例:
* curl -i -X POST -d task_id=任务id&send_num=发件人手机号&recv_nums=收件人手机号&sms_content=短信内容 http://localhost/api/v1/beehive/gsm_system/send_sms
*
* @apiParamExample {json}  请求参数示例:
*      {
*          "task_id":     "任务id",
*          "send_num":    "发件人手机号",
*          "recv_nums":   ["收件人imsi"],
*          "sms_content": "短信内容"
*      }
*
* @apiSuccessExample {json} 请求成功示例:
* {
*         "code":    0,
*         "data":{
*             "success": 2,
*             "failed":  1
*                   },
*         "msg":""
*
* }
 */

// 批量发送短信
func (GsmSystemController) GsmSystemSendSMS(ctx *gin.Context) {
	res := request.GetRequestBody(ctx)
	var sms l.Sms

	sms.TaskId = res.MustInt("task_id")
	// 自定义发件人手机号
	sms.Send = res.MustString("send_num")
	// 收件人的imsi列表
	sms.Recv = res.SliceString("recv_nums")
	// 短信内容
	sms.Content = res.MustString("sms_content")

	success := l.GsmSystemSimulateSms(sms)

	response.RenderSuccess(ctx, success)

}

/**
 * apiType http
 * @api {get} /api/v1/beehive/gsm_system/stop/:task_id 停止系统
 * @apiVersion 1.0.0
 * @apiName stop
 * @apiGroup gsm_system
 *
 * @apiDescription 停止系统
 *
 *
 * @apiParam {int}           task_id                任务id
 *
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X GET http://localhost/api/v1/beehive/gsm_system/stop/23
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "stop_status":     "true",
 *         "stop_message_id": 1
 *     },
 *     "msg":""
 * }
 */
// 停止系统
func (GsmSystemController) StopSystem(ctx *gin.Context) {
	tid := request.ParamInt(ctx, "task_id")
	stop, err := l.GsmSystemStop(tid)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, stop)
}

// 模拟短信时直接从数据库查找最后一次刷新的设备列表
func (GsmSystemController) DeviceList(ctx *gin.Context) {
	tid := request.ParamInt(ctx, "task_id")
	list, err := l.GsmSystemGetReceiveList(tid)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, list)
}

// 任务首页获取短信按钮
func (GsmSystemController) GetSmsButton(ctx *gin.Context) {
	tid := request.ParamInt(ctx, "task_id")
	log.Println(tid)
	sms, err := l.GsmSystemGetSms(tid)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, sms)
}

// 任务详情
func (GsmSystemController) Detail(ctx *gin.Context) {
	tid := request.ParamInt(ctx, "task_id")
	task, err := l.Detail(tid)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, task)
}
