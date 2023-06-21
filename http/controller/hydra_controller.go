package controller

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/custom_util/clog"
	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	hydra "skygo_detection/logic/hydra"
	"skygo_detection/mysql_model"
	"skygo_detection/service"
	"strconv"
	"strings"
)

type HydraController struct{}

const (
	// 1 手动录入 2 导入文件
	TYPE_HANDLE = 1
	TYPE_FILE   = 2
)

func (this HydraController) Create(ctx *gin.Context) {
	// 必填项
	taskName := request.MustString(ctx, "task_name")
	address := request.MustString(ctx, "address")
	port := request.MustString(ctx, "port")
	protocol := request.MustString(ctx, "protocol")
	username := request.MustString(ctx, "username")
	password := request.MustString(ctx, "password")

	uid := http_ctx.GetUserId(ctx)

	userNameType, err := strconv.Atoi(ctx.Request.FormValue("username_type"))
	passwdType, err := strconv.Atoi(ctx.Request.FormValue("password_type"))
	requestHost := strings.Join([]string{service.LoadHydraConfig().Client, "api", "v1",
		"hydra", "receive"}, "/")
	sid := ctx.Request.FormValue("sid")
	path := ctx.Request.FormValue("path")
	form := ctx.Request.FormValue("form")
	if path != "" && form != "" {
		if !strings.Contains(form, "^USER^") && !strings.Contains(form, "^PASS^") {
			response.RenderFailure(ctx, errors.New("表单格式有误"))
		}
	}
	// 密码为导入文件
	passwdFile := ""
	// 密码为导入文件（上传并保存密码字典、上传给密码破解服务）
	if passwdType == TYPE_FILE {
		file, fileHeader, err := ctx.Request.FormFile("pass_file")
		if err != nil {
			response.RenderFailure(ctx, err)
			return
		}
		passwdFile, err = hydra.ImportFile("password", taskName, file, fileHeader)
		if err != nil {
			response.RenderFailure(ctx, err)
			return
		}
	}
	// 用户名为导入文件（上传并保存用户名字典、上传给密码破解服务）
	userNameFile := ""
	if userNameType == TYPE_FILE {
		file, fileHeader, err := ctx.Request.FormFile("user_file")
		if err != nil {
			response.RenderFailure(ctx, err)
			return
		}
		userNameFile, err = hydra.ImportFile("username", taskName, file, fileHeader)
		if err != nil {
			response.RenderFailure(ctx, err)
			return
		}
	}
	params := &hydra.Params{
		TaskId:       0,
		TaskName:     taskName,
		Address:      address,
		Port:         port,
		Protocol:     common.UpperToLower[protocol],
		Path:         path,
		Form:         form,
		UserName:     username,
		UserNameFile: userNameFile,
		UserNameType: userNameType,
		Passwd:       password,
		PasswdFile:   passwdFile,
		PasswdType:   passwdType,
		Sid:          sid,
		RequestHost:  requestHost,
	}
	res, err := hydra.CreateTask(params, uid)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, res)
	return
}

func (this HydraController) Recv(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	success, has := req.TryBool("success")
	if !has {
		clog.Warn("hydra report error", zap.Any("request", req))
		response.RenderSuccess(ctx, nil)
		return
	}
	generator, _ := req.TryMap("generator")
	cmdline := generator.String("commandline")
	prefix := strings.Split(cmdline, "/root/result/")[1]
	suffix := strings.Split(prefix, ".json")[0]
	tid, _ := strconv.Atoi(suffix)

	results := fmt.Sprintln(req.Slice("results"))
	errormessages := fmt.Sprintln(req.Slice("errormessages"))
	origin := fmt.Sprintln(req)
	if success {
		hydra.Update(tid, 1, origin, results)
	} else {
		hydra.Update(tid, 2, origin, errormessages)
	}
	// 任务完成状态同步更新到task表
	//hydra.SyncTask(tid)
	response.RenderSuccess(ctx, req)
}
func (this HydraController) UploadUsername(ctx *gin.Context) {
	tid := ctx.Request.FormValue("task_id")
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		response.RenderFailure(ctx, err)
	}
	// 上传到mongo
	fileId, err := hydra.SaveUsernameFile(tid, file, header)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	// 上传给服务器
	url := strings.Join([]string{service.LoadHydraConfig().Server, "hydra", "upload", "username"}, "/")
	_, err = custom_util.HttpProxyFileUploadCustom(header, "file", tid, nil, nil, url)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	// 保存fileid到mysql
	taskId, _ := strconv.Atoi(tid)
	err = hydra.SaveUsernameFileId(fileId, taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}

	response.RenderSuccess(ctx, fileId)

}

func (this HydraController) UploadPassword(ctx *gin.Context) {
	tid := ctx.Request.FormValue("task_id")
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		response.RenderFailure(ctx, err)
	}
	// 上传到mongo
	fileId, err := hydra.SavePasswordFile(tid, file, header)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	// 上传给服务器
	url := strings.Join([]string{service.LoadHydraConfig().Server, "hydra", "upload", "password"}, "/")
	_, err = custom_util.HttpProxyFileUploadCustom(header, "file", tid, nil, nil, url)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}

	// 保存fileid到mysql
	taskId, _ := strconv.Atoi(tid)
	err = hydra.SavePasswordFileId(fileId, taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}

	response.RenderSuccess(ctx, fileId)

}

func (this HydraController) ProtocolList(ctx *gin.Context) {
	list := []string{}
	for k := range common.UpperToLower {
		list = append(list, k)
	}
	response.RenderSuccess(ctx, &list)
}

func (this HydraController) Detail(ctx *gin.Context) {
	tid := request.ParamInt(ctx, "task_id")
	one, err := hydra.GetOne(tid)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}

	response.RenderSuccess(ctx, one)
}
func (this HydraController) Abort(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	tid := req.MustInt("task_id")
	err := hydra.Abort(tid)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, qmap.QM{"task": tid, "status": "abort"})
}
func (this HydraController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	s := mysql.GetSession()

	// 查询组键
	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	widget.AddSorter(*(orm.NewSorter("create_time", orm.DESCENDING)))
	all := widget.PaginatorFind(s, &[]mysql_model.HydraTask{})
	response.RenderSuccess(ctx, all)
}

//func (this HydraController) GetAllWithTask(ctx *gin.Context) {
//	queryParams := ctx.Request.URL.RawQuery
//	s := mysql.GetSession().Table("hydra_task").
//}

func (this HydraController) Delete(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	tids := req.SliceInt("task_ids")
	success := hydra.Delete(tids)
	response.RenderSuccess(ctx, qmap.QM{"success": success})
}
func (this HydraController) Edit(ctx *gin.Context) {
	tid := request.ParamInt(ctx, "task_id")
	name := request.String(ctx, "name")
	_, err := hydra.Edit(tid, name)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, qmap.QM{"task_id": tid, "task_name": name})
}
