package beehive

import (
	"errors"
	"github.com/gin-gonic/gin"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/logic/beehive"
	"skygo_detection/mysql_model"
)

type LteSystemController struct{}

// 获取用户名和密码
func (this LteSystemController) GetApn(ctx *gin.Context) {
	taskId := request.ParamInt(ctx, "task_id")
	res, err := beehive.GetCrackResult(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, res)
	return
}

// 获取用户名和密码
func (this LteSystemController) GetList(ctx *gin.Context) {
	res, err := beehive.GetPassConfigData()
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, res)
	return
}

// 更新密码表
func (this LteSystemController) UpdatePassword(ctx *gin.Context) {
	// 获取上传的文件
	file, header, err := ctx.Request.FormFile("wordlist")
	if file == nil {
		response.RenderFailure(ctx, errors.New("请上传文件"))
		return
	}
	if err != nil {
		response.RenderFailure(ctx, errors.New("文件处理失败"))
		return
	}
	fileName := ctx.Request.FormValue("file_name")
	// 文件内容
	fileContent := make([]byte, header.Size)
	num, err := file.Read(fileContent)
	if num <= 0 {
		response.RenderFailure(ctx, errors.New("您上传的文件为空文件！"))
		return
	}
	res, err := beehive.UpdatePasswordConfig(fileName, fileContent, header)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, res)
	return
}

// 密码破解
func (this LteSystemController) CrackApn(ctx *gin.Context) {
	res, err := beehive.CrackApn()
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, res)
	return
}

// 写卡
func (this LteSystemController) Create(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	taskId := req.MustInt("task_id")
	imsl := req.MustString("imsi")
	res, err := beehive.SetImslLteEquipment(taskId, imsl)
	if err != nil {
		response.RenderFailure(ctx, err)
	}
	response.RenderSuccess(ctx, res)
}

// 启动LTE系统
func (this LteSystemController) StartSystem(ctx *gin.Context) {
	taskId := ctx.Param("task_id")
	res, err := beehive.StartSystem(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
	}
	response.RenderSuccess(ctx, res)
}

// 获取lte设备信息
func (this LteSystemController) GetOne(ctx *gin.Context) {
	tasId := ctx.Param("task_id")
	res, err := beehive.GetBasicInfo(tasId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, res)
}

// 停止lte设备
func (this LteSystemController) StopSystem(ctx *gin.Context) {
	taskId := ctx.Param("task_id")
	res, err := beehive.StopBasicInfo(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, res)
}

// 抓包
func (this LteSystemController) GetPackage(ctx *gin.Context) {
	taskId := ctx.Param("task_id")
	res, err := beehive.GetCapturePackage(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, res)
}

// 删除包
func (this LteSystemController) Delete(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	packageId := req.MustInt("package_id")
	res, err := beehive.DeleteLteSystemPackage(packageId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, res)
}

// 查询所有包
func (this LteSystemController) GetALL(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	taskId := request.QueryInt(ctx, "task_id")
	session := mysql.GetSession()
	session.Where("task_id=?", taskId)

	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	widget.AddSorter(*(orm.NewSorter("id", 1)))
	all := widget.PaginatorFind(session, &[]mysql_model.BeehiveLteSystemPackage{})
	response.RenderSuccess(ctx, all)
}

// 获取系统状态
func (this LteSystemController) GetSystemState(ctx *gin.Context) {
	taskId := request.ParamInt(ctx, "task_id")
	res, err := beehive.GetSystemState(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, res)
}

// 获取设备信息列表
func (this LteSystemController) GetEquipmentAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	taskId := request.QueryInt(ctx, "task_id")

	session := mysql.GetSession()
	session.Where("task_id=?", taskId)

	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	widget.AddSorter(*(orm.NewSorter("id", 1)))
	all := widget.PaginatorFind(session, &[]mysql_model.BeehiveLteSystem{})
	response.RenderSuccess(ctx, all)
}
