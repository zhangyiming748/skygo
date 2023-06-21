package hydra

import (
	"go.uber.org/zap"
	"mime/multipart"
	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/custom_util/clog"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/mysql_model"
	"skygo_detection/service"
	"strings"
)

type Params struct {
	TaskId       int    `json:"task_id"`       // 任务唯一标识
	TaskName     string `json:"task_name"`     //字符型taskid
	Address      string `json:"address"`       // 服务器地址
	Port         string `json:"port"`          // 手动输入端口号
	Protocol     string `json:"protocol"`      // 协议类型
	Path         string `json:"path"`          // 登录接口地址
	Form         string `json:"form"`          // 登录表单
	UserName     string `json:"username"`      // 用户手动输入用户名字典
	UserNameFile string `json:"username_file"` // 用户上传的用户名字典文件
	UserNameType int    `json:"username_type"` // 用户名类型 0:默认字典 1:手动录入 2:自定义上传
	Passwd       string `json:"passwd"`        // 用户手动输入密码字典
	PasswdFile   string `json:"passwd_file"`   // 用户上传的密码字典文件
	PasswdType   int    `json:"passwd_type"`   // 密码类型 0:默认字典 1:手动录入 2:自定义上传
	Sid          string `json:"sid"`           // Oracle Sid
	RequestHost  string `json:"request_host"`  // 发出查询请求的host
}

const (
	CATEGORY  = "密码爆破工具"
	TOOL      = "hydra"
	TOOL_TYPE = "hydra"
)

// 导入文件
func ImportFile(typeName string, taskName string, file multipart.File,
	fileHeader *multipart.FileHeader) (fileId string, err error) {
	// 上传到mongo
	fileId, err = SaveUsernameFile(taskName, file, fileHeader)
	if err != nil {
		return
	}
	// 上传给密码破解服务器
	urlPath := strings.Join([]string{service.LoadHydraConfig().Server, "hydra", "upload", typeName}, "/")
	_, err = custom_util.HttpProxyFileUploadCustom(fileHeader, "file", taskName,
		nil, nil, urlPath)
	if err != nil {
		return
	}
	return
}

// 任务创建
func CreateTask(params *Params, uid int64) (res bool, err error) {
	clog.Debug("CreateTask Params", zap.Any("params ", params), zap.Any("uid ", uid))
	// 创建主任务获取ID
	taskId, err := create(params.TaskName, uid)
	if err != nil {
		clog.Error("CreateTask create", zap.Any("taskId ", taskId), zap.Any("err ", err))
		return
	}
	params.TaskId = taskId
	clog.Debug("CreateTask Params", zap.Any("params TaskId", taskId))
	// 创建子任务
	err = createSubTask(params, int(uid), params.TaskName)
	if err != nil {
		clog.Error("CreateTask createSubTask", zap.Any("err ", err))
		return
	}
	// 发送命令
	urlPath := strings.Join([]string{service.LoadHydraConfig().Server, "hydra"}, "/")
	_, err = custom_util.HttpPostJson(nil, params, urlPath)
	if err != nil {
		clog.Error("CreateTask HttpPostJson", zap.Any("urlPath ", urlPath), zap.Any("err ", err))
		return
	}
	clog.Debug("CreateTask Response", zap.Any("urlPath ", urlPath))
	return
}

// 创建任务并获取ID
func create(taskName string, uid int64) (int, error) {
	task := new(mysql_model.Task)
	task.Name = taskName
	task.Category = CATEGORY
	task.Tool = TOOL
	task.ToolType = TOOL_TYPE
	task.TaskUuid = mysql.GetTaskId()
	task.CreateUserId = int(uid)
	_, err := task.Create()
	if err != nil {
		return 0, err
	}
	return task.Id, nil
}

func createSubTask(h *Params, uid int, tname string) error {
	this := new(mysql_model.HydraTask)
	this.TaskId = h.TaskId
	this.TaskName = tname
	this.Address = h.Address
	this.Port = h.Port
	this.Protocol = h.Protocol
	this.Path = h.Path
	this.Form = h.Form
	this.UserName = h.UserName
	this.UserNameFile = h.UserNameFile
	this.UserNameType = h.UserNameType
	this.Passwd = h.Passwd
	this.PasswdFile = h.PasswdFile
	this.PasswdType = h.PasswdType
	this.RequestHost = h.RequestHost
	this.Status = 1
	this.UserId = uid
	return this.Create()
}

func Update(tid int, success int, origin string, results string) {
	this := new(mysql_model.HydraTask)
	this.Results = results
	if success == 1 {
		this.Status = 2
		this.Success = 1
	} else {
		this.Status = 3
		this.Success = 2
	}
	this.Results = results
	this.OriginResults = origin

	err := this.UpdateByTaskId(tid)
	if err != nil {
		clog.Info("回传数据写入数据库", zap.Any("err", err))
	}
}

func SaveUsernameFile(tid string, fi multipart.File, header *multipart.FileHeader) (string, error) {
	if tid == "" && header != nil {
		tid = header.Filename
	}
	fileContent := make([]byte, header.Size)
	_, err := fi.Read(fileContent)
	if err != nil {
		clog.Info("read file error", zap.Any("err", err))
		return "", err
	}
	if fileId, err := mongo.GridFSUpload(common.Hydra, tid, fileContent); err == nil {
		return fileId, nil
	} else {
		return "", err
	}
}
func SavePasswordFile(tid string, fi multipart.File, header *multipart.FileHeader) (string, error) {
	if tid == "" && header != nil {
		tid = header.Filename
	}
	fileContent := make([]byte, header.Size)
	_, err := fi.Read(fileContent)
	if err != nil {
		clog.Info("read file error", zap.Any("err", err))
		return "", err
	}
	if fileId, err := mongo.GridFSUpload(common.Hydra, tid, fileContent); err == nil {
		return fileId, nil
	} else {
		return "", err
	}
}
func SaveUsernameFileId(fid string, tid int) error {
	this := new(mysql_model.HydraTask)
	this.UserNameFile = fid
	return this.UpdateByTaskId(tid)

}
func SavePasswordFileId(fid string, tid int) error {
	this := new(mysql_model.HydraTask)
	this.PasswdFile = fid
	return this.UpdateByTaskId(tid)
}

func GetOne(tid int) (mysql_model.HydraTask, error) {
	this, err := new(mysql_model.HydraTask).FindByTaskId(tid)
	if err != nil {
		return mysql_model.HydraTask{}, err
	}
	return this, err
}
func Abort(tid int) error {
	this := new(mysql_model.HydraTask)
	this.Status = 4
	return this.UpdateByTaskId(tid)
}
func Fail(tid int) error {
	this := new(mysql_model.HydraTask)
	this.Status = 3
	return this.UpdateByTaskId(tid)
}
func SyncTask(tid int) (int64, error) {
	this := new(mysql_model.Task)
	this.Id = tid
	return this.Update()
}
func Delete(tids []int) int {
	success := 0
	for _, tid := range tids {
		this := new(mysql_model.HydraTask)
		_, err := this.DeleteByTaskId(tid)
		if err != nil {
			continue
		}
		task := new(mysql_model.Task)
		task.Id = tid
		_, err = task.Remove()
		if err != nil {
			continue
		}
		success++
	}
	return success
}
func Edit(tid int, name string) (int64, error) {
	this := new(mysql_model.Task)
	this.Id = tid
	this.Name = name
	return this.Update("name")
}
