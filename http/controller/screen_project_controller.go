package controller

import (
	"encoding/json"

	"github.com/gin-gonic/gin"

	common "skygo_detection/lib/common_lib/response"
	"skygo_detection/service"
)

type ScreenProjectController struct{}

func (this *ScreenProjectController) All(ctx *gin.Context) {
	result := new(MyFileContreller).Open(service.LoadConfig().VehicleScreen.FilePath + "左侧-项目信息.txt")
	response := map[string]interface{}{}
	json.Unmarshal(result, &response)
	common.RenderSuccess(ctx, response)
}

func (this *ScreenProjectController) TestCase(ctx *gin.Context) {
	result := new(MyFileContreller).Open(service.LoadConfig().VehicleScreen.FilePath + "左侧-项目测试用例分布.txt")
	response := map[string]interface{}{}
	json.Unmarshal(result, &response)
	common.RenderSuccess(ctx, response)
}
