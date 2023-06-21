package beehive

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/globalsign/mgo/bson"
	"go.uber.org/zap"
	"io/ioutil"
	"mime/multipart"
	"os"
	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/custom_util/blog"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/mysql_model"
	"skygo_detection/mysql_model/beehive/lte"
	"skygo_detection/service"
	"strconv"
	"strings"
	"time"
)

type ReturnValue struct {
	Status    bool   `json:"status"`
	MessageId int    `json:"message_id"`
	Message   string `json:"message"`
}
type BasicInfo struct {
	Status    bool   `json:"status"`
	MessageId int    `json:"message_id"`
	Message   string `json:"message"`
	Apn       string `json:"apn"`
	Imsi      string `json:"imsi"`
	Ip        string `json:"ip"`
}

const (
	IMSI_HEAD = "00101"

	START_LTE_PARAMETER_BAND      = "0"
	START_LTE_PARAMETER_APN       = "skygoapn"
	START_LTE_PARAMETER_MCC       = "001"
	START_LTE_PARAMETER_MNC       = "01"
	GET_PACKAGE_PARAMETER_FILE_ID = 0

	URL_WRITESIM  = "/writesim"
	URL_GETFILE   = "/getfile"
	URL_START     = "/start"
	URL_BASICINFO = "/basicinfo"
	URL_STOPS     = "/stop"

	PATH         = "/hg_scanner/"
	STATUS_START = 1
	STATUS_STOP  = 2
)

func SetLogs(taskId int, content string) {
	log := new(mysql_model.BeehiveLog)
	log.TaskId = taskId
	log.Content = content
	err := log.SetLog()
	if err != nil {
		return
	}
	return
}

type CrackApnRes struct {
	Status    bool   `json:"status"`
	MessageId int    `json:"message_id"`
	Message   string `json:"message"`
}

type CrackResultRes struct {
	Status    bool   `json:"status"`
	MessageId int    `json:"message_id"`
	Message   string `json:"message"`
	Apn       string `json:"apn"`
	Imsi      string `json:"imsi"`
	Ip        string `json:"ip"`
	UserName  string `json:"user_name"`
	Password  string `json:"password"`
}

type UploadPasswordRes struct {
	Status    bool   `json:"status"`
	MessageId int    `json:"message_id"`
	Message   string `json:"message"`
}

type PassConfig struct {
	Status int    `json:"status"`
	Name   string `json:"name"`
}

// 获取密码配置
func GetPassConfigData() (res PassConfig, err error) {
	passwordConfigModel := new(lteModel.BeehiveLteSystemPasswordConfig)
	data, err := passwordConfigModel.GetOne()
	if err != nil {
		return
	}
	res.Status = data.Status
	res.Name = data.Name
	return res, nil
}

// 上传文件并更新密码配置文件
func UpdatePasswordConfig(fileName string, fileContent []byte, file *multipart.FileHeader) (
	uploadPasswordRes UploadPasswordRes, err error) {
	// 上传mongo
	fileId, err := mongo.GridFSUpload(common.MC_File, fileName, fileContent)
	if err != nil {
		return uploadPasswordRes, err
	}
	fmt.Println("fileID", fileId)
	// 更新数据库表
	passwordConfigModel := new(lteModel.BeehiveLteSystemPasswordConfig)
	data, err := passwordConfigModel.GetOne()
	if err != nil {
		return uploadPasswordRes, err
	}
	if data.Id <= 0 {
		return uploadPasswordRes, errors.New("默认配置数据失败")
	}
	data.Status = lteModel.StatusUpload
	data.UploadFileId = fileId
	data.UpdateTime = fmt.Sprint(time.Unix(int64(time.Now().Unix()), 0).Format("2006-01-02 15:04:05"))
	_, err = data.Update("status", "upload_file_id", "update_time")
	if err != nil {
		return uploadPasswordRes, err
	}
	//_, err = UploadPassword(file)
	//if err != nil {
	//	return uploadPasswordRes, err
	//}
	return
}

// 上传密码表文件
func UploadPassword(file *multipart.FileHeader) (uploadPasswordRes UploadPasswordRes, err error) {
	url, err := GetHost(LteSystem)
	if err != nil {
		return
	}
	url += common.LteUrlPasswordUpload
	resp, err := custom_util.HttpProxyFileUpload(file, "wordlist", nil, nil, url)
	blog.Info("UploadPassword", zap.Any("request: ", file), zap.Any("url: ", url),
		zap.Any("responds: ", string(resp)),
		zap.Any("err: ", err),
	)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &uploadPasswordRes)
	if err != nil {
		return
	}
	if uploadPasswordRes.Status == false {
		return uploadPasswordRes, errors.New(uploadPasswordRes.Message)
	}
	return uploadPasswordRes, nil
}

// 获取破解后的用户名和密码
func GetCrackResult(taskId int) (crackResultRes CrackResultRes, err error) {
	url, err := GetHost(LteSystem)
	if err != nil {
		return
	}
	url += common.LteUrlGetCrackResult
	resp, err := custom_util.HttpPostJson(nil, nil, url)
	blog.Info("GetCrackResult", zap.Any("url: ", url),
		zap.Any("responds: ", string(resp)),
		zap.Any("err: ", err),
	)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &crackResultRes)
	if err != nil {
		return
	}
	if crackResultRes.Status == false {
		return crackResultRes, errors.New(crackResultRes.Message)
	}

	module := new(mysql_model.BeehiveLteSystem)
	imsi := crackResultRes.Imsi
	splitImsi := strings.Replace(imsi, IMSI_HEAD, "", 1)
	id := module.Get(taskId, splitImsi)
	if id != 0 {
		module.UserName = crackResultRes.UserName
		module.Password = crackResultRes.Password
		module.Update()
		return crackResultRes, nil
	}
	module.TaskId = taskId
	module.Status = STATUS_START
	module.Apn = crackResultRes.Apn
	module.Imsi = splitImsi
	module.Ip = crackResultRes.Ip
	module.UserName = crackResultRes.UserName
	module.Password = crackResultRes.Password
	module.Create()
	return crackResultRes, nil
}

// 密码破解
func CrackApn() (crackApnRes CrackApnRes, err error) {
	url, err := GetHost(LteSystem)
	if err != nil {
		return
	}
	url += common.LteUrlCrackApn
	resp, err := custom_util.HttpPostJson(nil, nil, url)
	blog.Info("CrackApn", zap.Any("url: ", url),
		zap.Any("responds: ", string(resp)),
		zap.Any("err: ", err),
	)
	if err != nil {
		return
	}
	err = json.Unmarshal(resp, &crackApnRes)
	if err != nil {
		return
	}
	if crackApnRes.Status == false {
		return crackApnRes, errors.New(crackApnRes.Message)
	}
	return crackApnRes, nil
}

// 写卡
func SetImslLteEquipment(taskId int, imsl string) (ReturnValue, error) {
	newImsi := IMSI_HEAD + imsl
	resultRes, err := setImsl(newImsi)
	if err != nil {
		SetLogs(taskId, "写卡失败")
		return resultRes, err
	}
	if resultRes.MessageId != 1 && resultRes.MessageId != 4 {
		SetLogs(taskId, "写卡失败")
		return resultRes, err
	}
	SetLogs(taskId, "写卡成功")
	return resultRes, nil
}
func setImsl(imsl string) (ReturnValue, error) {
	var r ReturnValue
	args := make(map[string]string)
	args["imsi"] = imsl
	url, _ := GetHost(LteSystem)
	path := url + URL_WRITESIM
	resp, err := custom_util.HttpPostJson(nil, args, path)
	blog.Info("SetImslLteEquipment", zap.Any("url: ", url), zap.Any("responds: ", string(resp)), zap.Any("err: ", err))
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return r, err
	}
	return r, nil
}

// 启动系统
func StartSystem(strTaskId string) (ReturnValue, error) {
	taskId, _ := strconv.Atoi(strTaskId)
	resultRes, err := startEquipment()
	if err != nil {
		SetLogs(taskId, "启动LTE设备失败")
		return resultRes, err
	}
	if resultRes.MessageId == 1 || resultRes.MessageId == 2 {
		SetLogs(taskId, "系统启动成功")
		return resultRes, nil
	}
	SetLogs(taskId, "系统启动失败")
	return resultRes, err
}
func startEquipment() (ReturnValue, error) {
	var r ReturnValue
	args := make(map[string]interface{})
	args["band"] = START_LTE_PARAMETER_BAND
	args["apn"] = START_LTE_PARAMETER_APN
	args["mcc"] = START_LTE_PARAMETER_MCC
	args["mnc"] = START_LTE_PARAMETER_MNC
	beehiveConfig := service.LoadConfig().Beehive
	args["network"] = beehiveConfig.LetSystemNetwork
	url, err := GetHost(LteSystem)
	url += URL_START
	resp, err := custom_util.HttpPostJson(nil, args, url)
	blog.Info("StartSystem", zap.Any("url: ", url), zap.Any("responds: ", string(resp)), zap.Any("err: ", err))
	if err != nil {
		return r, err
	}
	err = json.Unmarshal(resp, &r)
	if err != nil {
		return r, err
	}
	return r, nil
}

// 获取lte设备信息
func GetBasicInfo(strTaskId string) (BasicInfo, error) {
	taskId, _ := strconv.Atoi(strTaskId)
	result, err := getBasicInfo()
	if err != nil {
		return result, err
	}
	if result.MessageId != 1 {
		SetLogs(taskId, "获取设备信息失败")
		return result, err
	}
	err = saveBasicInfo(taskId, result)
	if err != nil {
		return result, err
	}
	SetLogs(taskId, "获取设备信息成功")
	return result, nil
}
func getBasicInfo() (BasicInfo, error) {
	var b BasicInfo
	url, err := GetHost(LteSystem)
	url += URL_BASICINFO
	resp, err := custom_util.HttpPostJson(nil, nil, url)
	blog.Info("GetBasicInfo", zap.Any("url: ", url), zap.Any("responds: ", string(resp)), zap.Any("err: ", err))
	if err != nil {
		return b, err
	}
	err = json.Unmarshal(resp, &b)
	if err != nil {
		return b, err
	}
	return b, nil
}
func saveBasicInfo(taskId int, b BasicInfo) error {
	splitImsi := strings.Replace(b.Imsi, IMSI_HEAD, "", 1)

	module := new(mysql_model.BeehiveLteSystem)
	id := module.Get(taskId, splitImsi)
	module.Imsi = splitImsi
	module.Ip = b.Ip
	module.TaskId = taskId
	module.Apn = b.Apn
	module.Status = STATUS_START
	if id == 0 {
		module.Create()
	} else {
		module.Id = id
		module.Update()
	}
	return nil
}

// 停止系统
func StopBasicInfo(strTaskId string) (BasicInfo, error) {
	taskId, _ := strconv.Atoi(strTaskId)
	result, err := getBasicInfo()
	if err != nil {
		return result, err
	}
	resultRes, err := stopEquipment(taskId)
	if err != nil {
		return resultRes, err
	}
	if resultRes.MessageId != 1 && resultRes.MessageId != 2 {
		SetLogs(taskId, "系统停止失败")
		return resultRes, err
	}
	err = updateStatus(taskId)
	if err != nil {
		return resultRes, err
	}
	SetLogs(taskId, "系统停止成功")
	return resultRes, nil
}
func stopEquipment(taskId int) (BasicInfo, error) {
	var b BasicInfo
	url, err := GetHost(LteSystem)
	url += URL_STOPS
	resp, err := custom_util.HttpPostJson(nil, nil, url)
	blog.Info("StopBasicInfo", zap.Any("url: ", url), zap.Any("responds: ", string(resp)),
		zap.Any("err: ", err))
	if err != nil {
		SetLogs(taskId, "系统停止失败")
		return b, err
	}
	err = json.Unmarshal(resp, &b)
	if err != nil {
		return b, err
	}
	return b, nil
}
func updateStatus(taskId int) error {
	module := new(mysql_model.BeehiveLteSystem)
	module.TaskId = taskId
	module.Status = STATUS_STOP
	module.ForceUpdateSystemStatus()
	return nil
}

// 抓包
func GetCapturePackage(strTaskId string) (string, error) {
	taskId, err := strconv.Atoi(strTaskId)
	if err != nil {
		return "", err
	}
	downloadConfig := service.LoadConfig().Download
	path := downloadConfig.DownloadPath + PATH
	_, err = os.Stat(path)
	if err != nil {
		os.MkdirAll(path, 0777)
	}
	task, bool := mysql_model.TaskFindById(taskId)
	if !bool {
		return "没有启动设备", nil
	}
	b, err := getBasicInfo()
	if err != nil {
		return "", err
	}
	time := strconv.FormatInt(time.Now().Unix(), 10)
	fileName := task.Name + b.Ip + time + ".pcap"
	filePathName := path + fileName
	err = downloadPackage(taskId, filePathName)
	if err != nil {
		return "", err
	}
	fileId, _ := Upload(fileName, path)
	fileSize, err := countFileSize(filePathName)
	if err != nil {
		return "", err
	}
	err = SaveBeehiveLteSystemPackage(taskId, fileName, fileSize, fileId)
	if err != nil {
		return "", err
	}
	err = os.Remove(filePathName)
	if err != nil {
		return "", err
	}
	return "success", nil
}
func downloadPackage(taskId int, filePathName string) error {
	args := make(map[string]int)
	args["fileid"] = GET_PACKAGE_PARAMETER_FILE_ID
	url, err := GetHost(LteSystem)
	url += URL_GETFILE
	err = custom_util.HttpPostJsoDownload(nil, args, url, filePathName)
	blog.Info("downloadPackage", zap.Any("request:", args), zap.Any("responds:", err))
	if err != nil {
		SetLogs(taskId, "抓包失败")
		return err
	}
	SetLogs(taskId, "抓包成功")
	return nil
}
func SaveBeehiveLteSystemPackage(taskId int, filName string, fileSize string, fileId string) error {
	module := new(mysql_model.BeehiveLteSystemPackage)
	module.TaskId = taskId
	module.Name = filName
	module.Size = fileSize
	module.FileId = fileId
	err := module.Create()
	if err != nil {
		return err
	}
	return nil
}
func countFileSize(file string) (string, error) {
	fi, err := os.Stat(file)
	if err != nil {
		return "", err
	}
	fileSize := fi.Size()
	newStrSize := strconv.FormatInt(fileSize, 10)
	var flatFileSize = ""
	if len(newStrSize) <= 6 {
		flatSize := fileSize / 1024
		strSize := strconv.FormatInt(flatSize, 10)
		flatFileSize = strSize + "KB"
	} else {
		flatSize := float64(fileSize) / float64(1024*1024)
		strSize2 := strconv.FormatFloat(flatSize, 'E', -1, 64)
		strSplit0 := strings.Split(strSize2, ".")[0]
		strSplit1 := strings.Split(strSize2, ".")[1]
		new_strSplit := SubstrByByte(strSplit1, 2)
		flatFileSize = strSplit0 + "." + new_strSplit + "MB"
	}
	return flatFileSize, nil
}
func SubstrByByte(str string, length int) string {
	if len([]byte(str)) <= length {
		return str
	}
	bs := []byte(str)[:length]
	bl := 0
	for i := len(bs) - 1; i >= 0; i-- {
		switch {
		case bs[i] >= 0 && bs[i] <= 127:
			return string(bs[:i+1])
		case bs[i] >= 128 && bs[i] <= 191:
			bl++
		case bs[i] >= 192 && bs[i] <= 253:
			cl := 0
			switch {
			case bs[i]&252 == 252:
				cl = 6
			case bs[i]&248 == 248:
				cl = 5
			case bs[i]&240 == 240:
				cl = 4
			case bs[i]&224 == 224:
				cl = 3
			default:
				cl = 2
			}
			if bl+1 == cl {
				return string(bs[:i+cl])
			}
			return string(bs[:i])
		}
	}
	return ""
}
func Upload(filename string, path string) (fileId string, err error) {
	file, err := os.Open(path + filename)
	if err != nil {
		return "", err
	}
	defer file.Close()
	if fileContent, err := ioutil.ReadAll(file); err == nil {
		if fileId, err := mongo.GridFSUpload(common.MC_File, filename, fileContent); err == nil {
			return fileId, nil
		} else {
			return "", err
		}
	} else {
		return "", err
	}
	return
}

// 删除包
func DeleteLteSystemPackage(packageId int) (string, error) {
	fileId, err := mysql_model.Get(packageId)
	if err != nil {
		return "", err
	}
	err = mongo.GridFSRemoveFile(common.MC_File, bson.ObjectIdHex(fileId))
	if err != nil {
		return "", err
	}
	_, err = new(mysql_model.BeehiveLteSystemPackage).RemoveById(packageId)
	if err != nil {
		return "", err
	}
	return "success", nil
}

// 获取系统状态
func GetSystemState(taskId int) (map[string]interface{}, error) {
	var m map[string]interface{}
	m = make(map[string]interface{})
	lteSystemModel := mysql_model.BeehiveLteSystem{}
	_, err := lteSystemModel.FindByTaskId(taskId)
	if err != nil {
		return m, err
	}
	id := lteSystemModel.Id
	if id == 0 {
		m["task_id"] = lteSystemModel.TaskId
		m["status"] = lteSystemModel.Status
		return m, nil
	}
	m["task_id"] = lteSystemModel.TaskId
	m["status"] = lteSystemModel.Status
	m["apn"] = lteSystemModel.Apn
	m["imsi"] = lteSystemModel.Imsi
	m["ip"] = lteSystemModel.Ip
	return m, nil
}
