package toolbox

import (
	"github.com/gin-gonic/gin"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	t "skygo_detection/logic/toolbox"
)

type AppVersionController struct{}

func (a AppVersionController) GetAllApp(ctx *gin.Context) {
	task_id := request.ParamInt(ctx, "task_id")
	Apps, err := t.GetAllTaskApp(task_id)
	if err != nil {
		response.RenderFailure(ctx, err)
	}
	response.RenderSuccess(ctx, Apps)
}
func (a AppVersionController) GetAppAllVersion(ctx *gin.Context) {
	app := request.ParamString(ctx, "app_name")
	versions, err := t.GetAppAllVersion(app)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, versions)
}

func (a AppVersionController) GetAppVersionCompare(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	app := req.MustString("app_name")
	res, err := t.UsePermission(app)
	// zh := t.Permission2PermissionZhWithMap(res)
	zh := t.Permission2PermissionZh(res)
	if err != nil {
		response.RenderFailure(ctx, err)
	}
	response.RenderSuccess(ctx, zh)
}
