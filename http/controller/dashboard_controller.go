package controller

import (
	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mysql_model"
)

type DashboardController struct{}

/**
 * apiType http
 * @api {get} /api/v1/dashboard/summary_info 概要信息
 * @apiVersion 0.1.0
 * @apiName GetSummaryInfo
 * @apiGroup Dashboard
 *
 * @apiDescription 查询概要信息
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "completed_task_number": 1,
 *         "running_task_number": 7,
 *         "test_piece_number": 2,
 *         "vehicle_model_number": 1,
 *         "vul_number": 4
 *     },
 *     "msg": ""
 * }
 */
func (this DashboardController) GetSummaryInfo(ctx *gin.Context) {
	result := qmap.QM{}
	userId := http_ctx.GetUserId(ctx)
	// 查询车型数量
	{
		if total, err := sys_service.NewSession().Count(new(mysql_model.AssetVehicle)); err == nil {
			result["vehicle_model_number"] = total
		} else {
			panic(err)
		}
	}
	// 查询测试件数量
	{
		if total, err := sys_service.NewSession().Count(new(mysql_model.AssetTestPiece)); err == nil {
			result["test_piece_number"] = total
		} else {
			panic(err)
		}
	}
	// 查询"执行中"任务数量
	{
		params := qmap.QM{
			"e_create_user_id": userId,
			"e_status":         common.TASK_STATUS_RUNNING,
		}
		if total, err := sys_service.NewSessionWithCond(params).Count(new(mysql_model.Task)); err == nil {
			result["running_task_number"] = total
		} else {
			panic(err)
		}
	}
	// 查询"已结束"任务数量
	{
		params := qmap.QM{
			"e_create_user_id": userId,
			"in_status":        []int{common.TASK_STATUS_SUCCESS, common.TASK_STATUS_FAILURE},
		}
		if total, err := sys_service.NewSessionWithCond(params).Count(new(mysql_model.Task)); err == nil {
			result["completed_task_number"] = total
		} else {
			panic(err)
		}
	}
	// 查询漏洞数量
	{
		if total, err := sys_service.NewSessionWithCond(qmap.QM{"e_create_user_id": userId}).Count(new(mysql_model.Vulnerability)); err == nil {
			result["vul_number"] = total
		} else {
			panic(err)
		}
	}
	response.RenderSuccess(ctx, result)
}
