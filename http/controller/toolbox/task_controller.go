package toolbox

import (
	"errors"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mysql_model"

	"skygo_detection/logic/toolbox"

	"github.com/gin-gonic/gin"
)

type TaskController struct{}

/**
 * apiType http
 * @api {post} /api/v1/beehive/task 新增任务
 * @apiVersion 1.0.0
 * @apiName Create
 * @apiGroup Task任务相关
 *
 * @apiDescription  新增任务
 *
 * @apiParam {string}   name   任务名称
 * @apiParam {string}   tool_type   工具类型 gsm-sniffer、gsm-system、lte-system 三选一
 * @apiParam {string}   [describe]   任务描述
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *     "name":"任务名称",
 *     "tool_type":"gsm-sniffer|gsm-system|lte-system"
 *     "describe":"任务描述",
 * }
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST -d band=800 http://localhost/api/v1/beehive/task
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {}
 *     "msg":""
 * }
 */
func (t TaskController) Create(ctx *gin.Context) {
	name := request.MustString(ctx, "name")
	toolType := "privacy"
	// toolType := request.MustString(ctx, "tool_type")
	TaskLogic := new(toolbox.Task)
	// has := TaskLogic.CheckToolType(toolType)
	// if !has {
	// 	response.RenderFailure(ctx, errors.New("任务类型有误"))
	// 	return
	// }
	describe := request.String(ctx, "describe")
	data := map[string]int64{
		"id": 0,
	}
	id, err := TaskLogic.Create(ctx, name, toolType, describe)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	data["id"] = id
	response.RenderSuccess(ctx, data)
}

/**
 * apiType http
 * @api {put} /api/v1/beehive/task1/:task_id 修改任务
 * @apiVersion 1.0.0
 * @apiName Update
 * @apiGroup Task任务相关
 *
 * @apiDescription  修改任务
 *
 * @apiParam {string}   name   任务名称
 * @apiParam {string}   tool_type   工具类型 gsm-sniffer、gsm-system、lte-system 三选一
 * @apiParam {string}   [describe]   任务描述
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *     "name":"任务名称",
 *     "tool_type":"gsm-sniffer",
 *     "describe":"任务描述",
 * }
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST -d band=800 http://localhost/api/v1/beehive/task/1
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {}
 *     "msg":""
 * }
 *
 */
func (t TaskController) Update(ctx *gin.Context) {
	id := request.ParamInt(ctx, "task_id")
	name := request.MustString(ctx, "name")
	toolType := "privacy"
	// toolType := request.MustString(ctx, "tool_type")
	TaskLogic := new(toolbox.Task)
	// has := TaskLogic.CheckToolType(toolType)
	// if !has {
	// 	response.RenderFailure(ctx, errors.New("任务类型有误"))
	// 	return
	// }
	describe := request.String(ctx, "describe")
	if _, err := TaskLogic.Update(ctx, id, name, toolType, describe); err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, nil)
}

/**
 * apiType http
 * @api {get} /api/v1/beehive/task/:task_id 获取任务基础信息
 * @apiVersion 1.0.0
 * @apiName GetOne
 * @apiGroup Task任务相关
 *
 * @apiDescription 蜂窝网安全 获取任务基础信息
 *
 * @apiExample {curl} 请求示例:
 * curl http://localhost/api/v1/beehive/task/6
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *          "name": "任务名",
 *          "task_uuid":"G75WDAAM",
 *          "describe":"任务描述",
 *          "test_type":"GSM模拟",
 *          "tool_type":"gsm-system",
 *          "user":"创建人",
 *          "create_time":"2022-02-18 15:29:17",
 *     }
 *     "msg":""
 * }
 */
func (t TaskController) GetOne(ctx *gin.Context) {
	id := request.ParamInt(ctx, "task_id")
	TaskLogic := new(toolbox.Task)
	task := TaskLogic.GetOne(ctx, id)
	response.RenderSuccess(ctx, task)
}

/**
 * apiType http
 * @api {put} /api/v1/beehive/complete/task/:task_id 任务测试完成
 * @apiVersion 1.0.0
 * @apiName Complete
 * @apiGroup Task任务相关
 *
 * @apiDescription 蜂窝网安全 任务测试完成
 *
 * @apiExample {curl} 请求示例:
 * curl http://localhost/api/v1/beehive/complete/task/6
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {}
 *     "msg":""
 * }
 */
func (t TaskController) Complete(ctx *gin.Context) {
	id := request.ParamInt(ctx, "task_id")
	taskLogic := new(toolbox.Task)
	if _, err := taskLogic.Complete(ctx, id); err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, nil)
}
func (t TaskController) Start(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	task_uuid := req.MustString("task_uuid")
	task, err := mysql_model.GetTaskByUuid(task_uuid)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	if task.Id == 0 {
		response.RenderFailure(ctx, errors.New("任务不存在"))
		return
	}

	if task.Status == 2 {
		response.RenderFailure(ctx, errors.New("任务已经完成"))
		return
	}

	_ = new(mysql_model.PrivacyAnalysisLog).SetPrivacyLog(task.Id, "隐私任务开始测试")

	err = toolbox.Start(task_uuid)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, nil)
}
func (t TaskController) Stop(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	task_uuid := req.MustString("task_uuid")
	task, err := mysql_model.GetTaskByUuid(task_uuid)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	if task.Id == 0 {
		response.RenderFailure(ctx, errors.New("任务不存在"))
		return
	}

	if task.Status == 2 {
		response.RenderFailure(ctx, errors.New("任务已经完成"))
		return
	}

	_ = new(mysql_model.PrivacyAnalysisLog).SetPrivacyLog(task.Id, "隐私任务结束测试")

	err = toolbox.Stop(task_uuid)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, nil)
}
