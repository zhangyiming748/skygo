package beehive

import (
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/logic/beehive"
	"skygo_detection/mysql_model"

	"github.com/gin-gonic/gin"
)

type MemoController struct{}

/**
 * apiType http
 * @api {post} /api/v1/beehive/memo 增加备注
 * @apiVersion 1.0.0
 * @apiName Create
 * @apiGroup 备注
 *
 * @apiDescription GSM嗅探类型任务 增加备注
 *
 * @apiParam {int}   task_id   任务id
 * @apiParam {string}   content   备注的内容
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *     "task_id":6,
 *     "content":"我是备注内容"
 * }
 *
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X POST -d band=800 http://localhost/api/v1/beehive/memo
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {}
 * 	   "msg":""
 * }
 */
func (m MemoController) Create(ctx *gin.Context) {
	task_id := request.MustInt(ctx, "task_id")
	content := request.MustString(ctx, "content")
	memoLogic := new(beehive.Memo)
	if _, err := memoLogic.Save(ctx, task_id, content); err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, nil)
}

/*
*
  - apiType http
  - @api {get} /api/v1/beehive/memo/:task_id 查看备注
  - @apiVersion 1.0.0
  - @apiName View
  - @apiGroup 备注
    *
  - @apiDescription 获取备注记录列表
    *
  - @apiExample {curl} 请求示例:
  - curl  http://localhost/api/v1/beehive/memo/1
    *
  - @apiSuccessExample {json} 请求成功示例:
  - {
  - "code": 0,
  - "data": {
  - "id": 6,
  - "task_id": 372,
  - "content": "发顺丰,ddd，人人",
  - "create_time": "2022-02-23 14:11:03"
  - },
  - "msg": ""
    }
*/
func (m MemoController) View(ctx *gin.Context) {
	task_id := request.ParamInt(ctx, "task_id")
	memo := mysql_model.BeehiveMemo{}
	_, err := mysql.GetSession().Where("task_id=?", task_id).Get(&memo)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, memo)
}
