package beehive

import (
	"github.com/gin-gonic/gin"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	l "skygo_detection/logic/beehive"
)

type LogController struct{}

/**
 * apiType http
 * @api {get} /api/v1/beehive/log/:task_id 获取任务的日志
 * @apiVersion 1.0.0
 * @apiName GetAll
 * @apiGroup beehive
 * @apiDescription 获取任务的日志
 *
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X GET -d  http://localhost/api/v1/beehive/log/23
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *   "code": 0,
 *   "data": [
 *       {
 *           "id": 1,
 *           "task_id": 23,
 *           "create_time": 22,
 *           "content": "fgbfg"
 *       },
 *       {
 *           "id": 2,
 *           "task_id": 23,
 *           "create_time": 345,
 *           "content": "sgb"
 *       }
 *   ],
 *   "msg": ""
 * }
 */
func (this LogController) GetAll(ctx *gin.Context) {
	tid := request.ParamInt(ctx, "task_id")
	queryParams := ctx.Request.URL.RawQuery
	res := l.GetBeehiveTaskLog(tid, queryParams)
	response.RenderSuccess(ctx, res)
}
