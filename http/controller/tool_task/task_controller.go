package tool_task

import (
	"errors"
	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/app/http/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/logic/tool_task"
	ump_model "skygo_detection/mysql_model/ump"
)

type ToolTaskController struct{}

// 创建任务
func (t ToolTaskController) Create(ctx *gin.Context) {
	umpUserId := request.MustInt(ctx, "ump_user_id")
	taskUuid := request.MustString(ctx, "task_uuid")
	subTaskId := request.MustInt(ctx, "sub_task_id")
	subDetailUrl := request.MustString(ctx, "sub_detail_url")
	taskName := request.MustString(ctx, "task_name")
	toolType := request.MustString(ctx, "tool_type")
	describe := request.MustString(ctx, "describe")
	category := request.MustString(ctx, "category")
	tool := request.MustString(ctx, "tool")
	taskLogic := new(tool_task.Task)
	// 获取映射关系用户ID
	umpModel := new(ump_model.UmpUser)
	umpInfo, err := umpModel.GetUmpInfo(umpUserId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	if umpInfo.UserId <= 0 {
		response.RenderFailure(ctx, errors.New("用户ID不存在"))
		return
	}
	id, err := taskLogic.Create(umpInfo.UserId, taskUuid, subTaskId, subDetailUrl, taskName,
		toolType, describe, tool, category)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	data := map[string]int64{
		"id": id,
	}
	response.RenderSuccess(ctx, data)
}

// 任务状态更新
func (t ToolTaskController) UpdateStatus(ctx *gin.Context) {
	subTaskId := request.MustInt(ctx, "sub_task_id")
	toolType := request.MustString(ctx, "tool_type")
	taskUuid := request.MustString(ctx, "task_uuid")
	status := request.MustInt(ctx, "status")
	taskLogic := new(tool_task.Task)
	id, err := taskLogic.UpdateStatus(subTaskId, toolType, taskUuid, status)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	data := map[string]int64{
		"id": id,
	}
	response.RenderSuccess(ctx, data)
}

// 删除任务
func (t ToolTaskController) Delete(ctx *gin.Context) {
	taskLogic := new(tool_task.Task)
	subTaskId := request.MustInt(ctx, "sub_task_id")
	toolType := request.MustString(ctx, "tool_type")
	taskUuid := request.MustString(ctx, "task_uuid")
	id, err := taskLogic.Delete(subTaskId, toolType, taskUuid)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	data := map[string]int64{
		"id": id,
	}
	response.RenderSuccess(ctx, data)
}
