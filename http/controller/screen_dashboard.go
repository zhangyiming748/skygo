package controller

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"

	common "skygo_detection/lib/common_lib/response"
	"skygo_detection/service"
)

type Dashboard struct{}

func (this *Dashboard) Vehicle(ctx *gin.Context) {
	result := new(MyFileContreller).Open(service.LoadConfig().VehicleScreen.FilePath + "中间-车型.txt")
	response := map[string]interface{}{}
	json.Unmarshal(result, &response)
	common.RenderSuccess(ctx, response)

}

func (this *Dashboard) Task(ctx *gin.Context) {
	result := new(MyFileContreller).Open(service.LoadConfig().VehicleScreen.FilePath + "左侧-测试任务.txt")
	response := map[string]interface{}{}
	json.Unmarshal(result, &response)
	common.RenderSuccess(ctx, response)

}

func (this *Dashboard) Vul(ctx *gin.Context) {
	result := new(MyFileContreller).Open(service.LoadConfig().VehicleScreen.FilePath + "右侧-车型未修复漏洞分布.txt")
	response := map[string]interface{}{}
	json.Unmarshal(result, &response)
	common.RenderSuccess(ctx, response)
}

func (this *Dashboard) Top(ctx *gin.Context) {
	result := new(MyFileContreller).Open(service.LoadConfig().VehicleScreen.FilePath + "右侧-漏洞排行榜.txt")
	response := map[string]interface{}{}
	json.Unmarshal(result, &response)
	common.RenderSuccess(ctx, response)
}

func (this *Dashboard) TestCase(ctx *gin.Context) {
	result := new(MyFileContreller).Open(service.LoadConfig().VehicleScreen.FilePath + "右侧-测试用例完成情况.txt")
	response := map[string]interface{}{}
	json.Unmarshal(result, &response)
	common.RenderSuccess(ctx, response)
}

func (this *Dashboard) AssetTestPieces(ctx *gin.Context) {
	result := new(MyFileContreller).Open(service.LoadConfig().VehicleScreen.FilePath + "中间-测试件分布接口.txt")
	response := map[string]interface{}{}
	json.Unmarshal(result, &response)
	common.RenderSuccess(ctx, response)
}

func (this *Dashboard) AssetTestPiecesVul(ctx *gin.Context) {
	result := new(MyFileContreller).Open(service.LoadConfig().VehicleScreen.FilePath + "中间-测试组件漏洞的分布.txt")
	response := map[string]interface{}{}
	json.Unmarshal(result, &response)
	common.RenderSuccess(ctx, response)
}

func (this *Dashboard) CaseDailyTotal(ctx *gin.Context) {
	result := new(MyFileContreller).Open(service.LoadConfig().VehicleScreen.FilePath + "中间-测试用例实时图.txt")
	response := map[string]interface{}{}
	json.Unmarshal(result, &response)
	common.RenderSuccess(ctx, response)
}

func (this *Dashboard) Upload(ctx *gin.Context) {
	file, err := ctx.FormFile("file")
	if err != nil {
		ctx.String(500, "上传错误")
	}
	// c.JSON(200, gin.H{"message": file.Header.Context})

	src, _ := file.Open()
	defer src.Close()
	tmp, _ := ioutil.ReadAll(src)
	var a = map[string]interface{}{}
	if err := json.Unmarshal(tmp, &a); err != nil {
		fmt.Println(err.Error())
		common.RenderFailure(ctx, errors.New("json 存在问题"))
		return
	}
	ctx.SaveUploadedFile(file, service.LoadConfig().VehicleScreen.FilePath+file.Filename)
	ctx.String(http.StatusOK, file.Filename)
}
