package gps

import (
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"

	"skygo_detection/logic/gps"

	"github.com/gin-gonic/gin"
)

type TaskController struct{}

func (t TaskController) Create(ctx *gin.Context) {
	name := request.MustString(ctx, "name")
	TaskLogic := new(gps.Task)
	describe := request.String(ctx, "describe")
	data := map[string]int64{
		"id": 0,
	}
	id, err := TaskLogic.Create(ctx, name, describe)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	data["id"] = id
	response.RenderSuccess(ctx, data)
}

func (t TaskController) Update(ctx *gin.Context) {
	id := request.ParamInt(ctx, "task_id")
	name := request.MustString(ctx, "name")
	TaskLogic := new(gps.Task)
	describe := request.String(ctx, "describe")
	if _, err := TaskLogic.Update(ctx, id, name, describe); err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, nil)
}

func (t TaskController) GetOne(ctx *gin.Context) {
	id := request.ParamInt(ctx, "task_id")
	TaskLogic := new(gps.Task)
	task := TaskLogic.GetOne(ctx, id)
	response.RenderSuccess(ctx, task)
}

func (t TaskController) Complete(ctx *gin.Context) {
	id := request.ParamInt(ctx, "task_id")
	taskLogic := new(gps.Task)
	if _, err := taskLogic.Complete(ctx, id); err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, nil)
}
