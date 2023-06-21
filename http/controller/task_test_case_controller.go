package controller

import (
	"github.com/gin-gonic/gin"

	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mysql_model"
)

type TaskTestCaseController struct{}

// 任务的测试用例分页展示
// 必须传task_id, 查询指定任务id的测试案例分页列表
func (this TaskTestCaseController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	s := mysql.GetSession()

	taskId := ctx.Query("task_id")
	s.Where("task_id = ?", taskId)

	// 查询组键
	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	widget.AddSorter(*(orm.NewSorter("action_status", 1)))
	all := widget.PaginatorFind(s, &[]mysql_model.TaskTestCaseView{})
	response.RenderSuccess(ctx, all)
}

/**
 * apiType http
 * @api {get} /api/v1/task/test_cases/:id 查询任务测试用例
 * @apiVersion 0.1.0
 * @apiName GetOne
 * @apiGroup TaskTestCase
 *
 * @apiDescription 根据id,查询任务测试用例
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "action_status": 2,
 *         "auto_test_level": 0,
 *         "case_priority": 0,
 *         "case_result": "",
 *         "case_result_file": "616fc3e224b6474e38e3d91b|616fc96424b6474e8df02fcb",
 *         "case_uuid": "hg_CASE_6",
 *         "complete_time": 0,
 *         "create_time": 1634713198,
 *         "demand_id": 0,
 *         "file_id": "",
 *         "id": 370,
 *         "log_result": 0,
 *         "task_id": 50,
 *         "task_name": "基金",
 *         "task_uuid": "G3GGVJ5W",
 *         "template_id": 0,
 *         "test_attachment": "",
 *         "test_case_id": 266,
 *         "test_case_name": "通信传输安全",
 *         "test_procedure": "",
 *         "test_result_status": 0,
 *         "test_tool_name": "车机检测工具",
 *         "test_tools": "hg_scanner",
 *         "update_time": 1634713198
 *     },
 *     "msg": ""
 * }
 */
func (this TaskTestCaseController) GetOne(ctx *gin.Context) {
	id := request.ParamString(ctx, "id")
	s := mysql.GetSession()
	s.Where("id=?", id)

	w := orm.PWidget{}
	result, err := w.One(s, &mysql_model.TaskTestCase{})

	if err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/task/test_cases/:id 更新任务测试用例
 * @apiVersion 0.1.0
 * @apiName Update
 * @apiGroup TaskTestCase
 *
 * @apiDescription 根据id,更新任务测试用例
 *
 * @apiUse authHeader
 *
 * @apiParam {int} 			action_status  			测试用例状态(待测试:1,队列中:2,测试中:3,分析中:4,测试完成:5)
 * @apiParam {int} 			test_result_status  	测试结果状态(通过:1,未通过:2)
 * @apiParam {string}       case_result_file    	测试文件
 * @apiParam {string}		test_attachment			测试用例附件
 * @apiParam {string} 		test_procedure			测试过程
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *		"action_status":3,
 *		"test_result_status":1,
 *		"case_result_file": "5fd7218624b64712a27f47e8",
 *		"test_attachment": "",
 *		"test_procedure": "测试过程"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 * 		"code": 0
 * }
 */
func (this TaskTestCaseController) Update(ctx *gin.Context) {
	id := request.ParamInt(ctx, "id")
	taskCase := new(mysql_model.TaskTestCase)
	if _, err := taskCase.UpdateCaseById(id, *(request.GetRequestBody(ctx))); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, nil)
	}
}
