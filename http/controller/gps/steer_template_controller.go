package gps

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"skygo_detection/common"
	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/logic/project_file_logic"
	"skygo_detection/mysql_model"
)

const (
	STATE          = 5
	TEMPLATE_STATE = 2
)

type SteerTemplateController struct{}

func (this SteerTemplateController) GetAll(ctx *gin.Context) {
	res := mysql_model.GetAllTemplate()
	response.RenderSuccess(ctx, res)
}

func (this SteerTemplateController) Create(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	module := new(mysql_model.GpsSteerTemplate)
	module.Name = req.MustString("name")
	module.MaxLatAcc = req.Float32("max_lat_acc")
	module.MaxLongAcc = req.Float32("max_long_acc")
	module.MaxJerk = req.MustInt("max_jerk")
	module.MaxSpeed = req.Float32("max_speed")
	module.StationaryPeriod = req.Float32("stationary_period")
	module.StationaryPeriodEnd = req.Float32("stationary_period_end")
	module.PositionSmoothingFactor = req.MustInt("position_smoothing_factor")
	module.SpeedSmoothingFactor = req.MustInt("speed_smoothing_factor")
	module.PictureName = req.MustString("file_name")
	userModule := new(mysql_model.SysUser)
	userId := int(http_ctx.GetUserId(ctx))
	user, has := userModule.FindById(userId)
	if has {
		module.Creator = user.Realname
	}
	module.Type = STATE
	module.TemplateState = TEMPLATE_STATE
	module.FileId = req.MustString("file_id")

	if _, err := module.Create(); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, "success")
	}
}

func (this SteerTemplateController) Update(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	id := request.ParamInt(ctx, "template_id")
	module := new(mysql_model.GpsSteerTemplate)
	// 系统行驶模板，不可修改、不可删除
	templateModule, _ := module.FindById(id)
	if templateModule.TemplateState == 1 {
		response.RenderFailure(ctx, errors.New("系统模板不允许修改"))
		return
	}

	if request.IsExist(ctx, "name") {
		module.Name = req.MustString("name")
	}
	if request.IsExist(ctx, "max_lat_acc") {
		module.MaxLatAcc = req.Float32("max_lat_acc")
	}
	if request.IsExist(ctx, "max_long_acc") {
		module.MaxLongAcc = req.Float32("max_long_acc")
	}
	if request.IsExist(ctx, "max_jerk") {
		module.MaxJerk = req.MustInt("max_jerk")
	}
	if request.IsExist(ctx, "max_speed") {
		module.MaxSpeed = req.Float32("max_speed")
	}
	if request.IsExist(ctx, "stationary_period") {
		module.StationaryPeriod = req.Float32("stationary_period")
	}
	if request.IsExist(ctx, "stationary_period_end") {
		module.StationaryPeriodEnd = req.Float32("stationary_period_end")
	}
	if request.IsExist(ctx, "position_smoothing_factor") {
		module.PositionSmoothingFactor = req.MustInt("position_smoothing_factor")
	}
	if request.IsExist(ctx, "speed_smoothing_factor") {
		module.SpeedSmoothingFactor = req.MustInt("speed_smoothing_factor")
	}
	if request.IsExist(ctx, "file_id") {
		module.FileId = req.MustString("file_id")
	}
	if request.IsExist(ctx, "file_name") {
		module.PictureName = req.MustString("file_name")
	}
	if _, err := module.Update(id); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, "success")
	}
}

func (this SteerTemplateController) Delete(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	id := req.MustInt("id")
	module := new(mysql_model.GpsSteerTemplate)

	rep, _ := module.FindById(id)
	if rep.TemplateState == 1 {
		response.RenderFailure(ctx, errors.New("系统模板不允许删除"))
		return
	}
	if _, err := module.Remove(id); err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, "success")
}

// 上传
func (this SteerTemplateController) Upload(ctx *gin.Context) {
	fileName := ctx.Request.FormValue("file_name")
	file, header, _ := ctx.Request.FormFile("file")
	if fileName == "" && header != nil {
		fileName = header.Filename
	}
	fileContent := make([]byte, header.Size)
	file.Read(fileContent)
	if fileId, err := mongo.GridFSUpload(common.MC_File, fileName, fileContent); err == nil {
		res := &qmap.QM{"file_id": fileId, "file_name": fileName}
		response.RenderSuccess(ctx, res)
	} else {
		panic(err)
	}
}

// 下载
func (this SteerTemplateController) Download(ctx *gin.Context) {
	id := request.ParamInt(ctx, "template_id")
	module := new(mysql_model.GpsSteerTemplate)
	resp, has := module.FindById(id)
	if !has {
		response.RenderFailure(ctx, errors.New("文件没有上传"))
		return
	}
	fi, fileContent, err := project_file_logic.FindFileByFileID(resp.FileId)
	if err != nil {
		response.RenderFailure(ctx, errors.New("file_id不正确"))
		return
	}
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fi.Name()))
	ctx.Header("Content-Type", "*")
	ctx.Header("Accept-Length", fmt.Sprintf("%d", len(fileContent)))
	ctx.Writer.Write(fileContent)
}
