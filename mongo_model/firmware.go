package mongo_model

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"
	"skygo_detection/service"

	"github.com/globalsign/mgo/bson"
)

type FirmWareData struct {
	Id              bson.ObjectId `bson:"_id"`
	ProjectName     string        `bson:"project_name"`     // 项目名称 工程名称
	DeviceName      string        `bson:"device_name"`      // 设备名称
	DeviceModel     string        `bson:"device_model"`     // 设备型号
	FirmwareVersion string        `bson:"firmware_version"` // 固件版本
	DeviceType      string        `bson:"device_type"`      // 设备类型
	FirmwareName    string        `bson:"firmware_name"`    // 文件名 固件名称
	TmpHdFilePath   string        `bson:"tmp_hd_file_path"` // 固件本地存储地址
	CreateTime      int64         `bson:"create_time"`      // 文件创建时间
	ProjectId       int           `bson:"project_id"`       // 项目 工程ID 接口返回值
	FirmwareSize    int           `bson:"firmware_size"`    // 固件大小
	FirmwareMd5     string        `bson:"firmware_md5"`     // 固件md5
	UploadUser      string        `bson:"upload_user"`      // 上传者用户名
	UploadUserId    int           `bson:"upload_user_id"`   // 上传者ID
	UploadTime      int64         `bson:"upload_time"`      // 上传时间
	Progress        int           `bson:"progress"`         // 上传进度
	Status          int           `bson:"status"`           // 状态 1 待上传 2 上传完成 3 上传失败 4 取消上传 5 (下载完成) 扫描中 6 (创建任务) 扫描中 7 取消扫描 8 扫描完成 9 扫描失败 10 已解析 0 已删除
	TaskId          int           `bson:"task_id"`          // 任务ID 接口返回值
	TaskName        string        `bson:"task_name"`        // 任务名称
	TemplateId      int64         `bson:"template_id"`      // 模板ID
	TemplateName    string        `bson:"template_name"`    // 模板ID
	ResponseTime    int64         `bson:"response_time"`    // 接口返回时间
	RealFileName    string        `bson:"real_file_name"`   // 真实文件名
	ResponseFile    string        `bson:"response_file"`    // 接口返回地址
}

/*
 * 上传固件 文件流上传 已废弃
 */
func (this *FirmWareData) Create(FileName, FilePath, UploadReturnFilePath, DeviceName, DeviceModel, FirmwareVersion, DeviceType, FirmwareMd5 string, FirmwareSize int, TemplateId int64) (*FirmWareData, error) {
	nowTime := custom_util.GetCurrentMilliSecond() / 1000
	tempStr := custom_util.TimestampToString(nowTime)

	this.Id = bson.NewObjectId()
	this.ProjectName = "工程名" + tempStr
	this.DeviceName = DeviceName
	this.RealFileName = UploadReturnFilePath
	this.TmpHdFilePath = FilePath
	this.FirmwareName = FileName
	this.DeviceModel = DeviceModel
	this.FirmwareVersion = FirmwareVersion
	this.DeviceType = DeviceType
	this.FirmwareMd5 = FirmwareMd5
	this.FirmwareSize = FirmwareSize
	this.CreateTime = nowTime
	this.TemplateId = TemplateId
	this.Status = 0
	return this, errors.New("api not found")
}

/*
 * 上传固件文本信息活的上传授权
 */
func (this *FirmWareData) CreateFirmWareMsg(FirmwareName, DeviceName, DeviceModel, FirmwareVersion, DeviceType, UserName string, TemplateId, UserID int) (*FirmWareData, error) {
	nowTime := custom_util.GetCurrentMilliSecond() / 1000
	tempStr := custom_util.TimestampToStringNoSpace(nowTime)

	this.Id = bson.NewObjectId()
	this.ProjectName = "ProName_" + tempStr
	this.DeviceName = DeviceName
	this.DeviceModel = DeviceModel
	this.FirmwareVersion = FirmwareVersion
	this.DeviceType = DeviceType
	this.CreateTime = nowTime
	this.UploadUserId = UserID
	this.UploadUser = UserName
	this.FirmwareName = FirmwareName
	this.TemplateId = int64(TemplateId)
	this.Status = 1
	this.Progress = 0
	this.TaskId = 0

	this.ProjectId = 0
	mongoClient := mongo.NewMgoSession(common.MC_FIRMWARE_UPLOAD_LOG)
	if err := mongoClient.Insert(this); err == nil {
		return this, nil

	} else {
		return nil, errors.New("mongo db save message error")
	}

	// ProjectId := CreateProject(this.ProjectName, DeviceName, DeviceModel, FirmwareVersion, DeviceType)
	// if ProjectId > 0 {
	//	this.ProjectId = ProjectId
	//	mongoClient := mongo.NewMgoSession(common.MC_FIRMWARE_UPLOAD_LOG)
	//	if err := mongoClient.Insert(this); err == nil {
	//		return this, nil
	//	} else {
	//		return nil, errors.New("mongo db save message error")
	//	}
	// } else {
	//	return nil, errors.New("api project create error")
	// }
}

type FirmwareCreateProResp struct {
	Code int            `bson:"code"`
	Data map[string]int `bson:"data"`
}

// /*
// * 创建上传固件工程 调用安全团队api
// */
// func CreateProject(projectName, deviceName, deviceModel, firmwareVersion, deviceType string) int {
//
//	url := common.FIRM_WARE_API + "/api/projects"
//	var jsonStr = []byte(`
// {
//	"project_name": "` + projectName + `",
//    "device_name": "` + deviceName + `",
//    "device_model": "` + deviceModel + `",
//    "firmware_version": "` + firmwareVersion + `",
//    "device_type": "` + deviceType + `"
// }`)
//	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
//	//req.Header.Set("X-Custom-Header", "myvalue")
//	req.Header.Set("Content-Type", "application/json")
//
//	client := &http.Client{}
//	resp, err := client.Do(req)
//	if err != nil {
//		panic(err)
//	}
//	defer resp.Body.Close()
//	respBodyJson, _ := ioutil.ReadAll(resp.Body)
//
//	respBodyMap := []byte(respBodyJson)
//	resps := FirmwareCreateProResp{}
//	json.Unmarshal(respBodyMap, &resps)
//
//	var proId int
//	for _, i := range resps.Data {
//		proId = i
//
//	}
//	return proId
// }

/*
 * 上传固件url信息
 */
func (this *FirmWareData) UploadFirmWareUrl(ProjectId, fileSize, UserID int, filePath, RealFileName, firmwareMd5, masterId, UserName string) error {

	id, _ := primitive.ObjectIDFromHex(masterId)
	getFirmwareParams := qmap.QM{
		"e_project_id": ProjectId,
		"e__id":        id,
		"in_status":    []int{1, 3, 4},
	}
	// 判断当前上传的固件记录是否存在
	if firmWare, err := orm_mongo.NewWidgetWithParams(common.MC_FIRMWARE_UPLOAD_LOG, getFirmwareParams).Get(); err == nil {
		update := bson.M{
			"$set": bson.M{
				"real_file_name":   RealFileName,
				"tmp_hd_file_path": filePath,
				"firmware_size":    fileSize,
				"firmware_md5":     firmwareMd5,
				"status":           2, // 更新固件信息为上传成功
			},
		}

		coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_FIRMWARE_UPLOAD_LOG)
		if _, err := coll.UpdateOne(context.Background(), bson.M{"_id": id}, update); err == nil {

			// 插入队列任务信息
			taskName := firmWare.MustString("project_name") + " 扫描任务 " + masterId
			if err := new(ScannerTask).TaskInsert(masterId, taskName, common.TOOL_FIRMWARE_SCANNER); err == nil {
				return nil
			} else {
				return errors.New("mongo db save queue scanner msg error")
			}
		} else {
			return errors.New("scanning num > 4 and update status = 2 error")
		}
	} else {
		return errors.New("not found this firmware project")
	}

	// //查看此固件当前多少个处于扫描中状态
	// scanningParams := qmap.QM{
	//	"in_status":        []int{5, 6, 8},
	//	"e_upload_user_id": UserID,
	//	"e_upload_user":    UserName,
	// }
	// nowScanningCount, _ := mongo.NewMgoSessionWithCond(common.MC_FIRMWARE_UPLOAD_LOG, scanningParams).Count()
	// //如果当前用户扫描中用户小于4个
	// if nowScanningCount < 4 {
	//	getFirmwareParams := qmap.QM{
	//		"e_project_id": ProjectId,
	//		"e__id":        bson.ObjectIdHex(masterId),
	//		"in_status":    []int{1, 3, 4},
	//	}
	//	//判断当前上传的固件记录是否存在
	//	if _, err := mongo.NewMgoSessionWithCond(common.MC_FIRMWARE_UPLOAD_LOG, getFirmwareParams).GetOne(); err == nil {
	//		//step 1 首先更新固件为上传成功
	//		update := bson.M{
	//			"$set": bson.M{
	//				"real_file_name":   RealFileName,
	//				"tmp_hd_file_path": filePath,
	//				"firmware_size":    fileSize,
	//				"firmware_md5":     firmwareMd5,
	//				"status":           2, //更新固件信息为上传成功
	//			},
	//		}
	//
	//		if err := mongo.NewMgoSession(common.MC_FIRMWARE_UPLOAD_LOG).Update(bson.M{"_id": bson.ObjectIdHex(masterId)}, update); err == nil {
	//			//post url 并且开始任务
	//			postDownloadUrl := common.FIRM_WARE_DOWNLOAD_MOMAIN+"/api/v1/firmware/download?name="+RealFileName
	//			if _, urlErr := RemoteDownload(ProjectId, postDownloadUrl); urlErr == nil {
	//				updateSM := bson.M{
	//					"$set": bson.M{
	//						"status": 5, //更新固件信息 post url成功 并且等待下载
	//					},
	//				}
	//				if err := mongo.NewMgoSession(common.MC_FIRMWARE_UPLOAD_LOG).Update(bson.M{"_id": bson.ObjectIdHex(masterId), "status": 2}, updateSM); err == nil {
	//					return nil
	//				} else {
	//					return errors.New("update status 5 msg error")
	//
	//				}
	//			} else {
	//				return errors.New("post url api error")
	//			}
	//		} else {
	//			return errors.New("firmware update msg error")
	//		}
	//	} else {
	//		return errors.New("not found this firmware project")
	//	}
	// } else {
	//	//当前任务大于4个 固件状态为2 上传成功 等待其它任务扫描完成后 后台轮询脚本自动执行
	//	update := bson.M{
	//		"$set": bson.M{
	//			"real_file_name":   RealFileName,
	//			"tmp_hd_file_path": filePath,
	//			"firmware_size":    fileSize,
	//			"firmware_md5":     firmwareMd5,
	//			"status":           2, //更新固件信息为上传成功
	//		},
	//	}
	//
	//	if err := mongo.NewMgoSession(common.MC_FIRMWARE_UPLOAD_LOG).Update(bson.M{"_id": bson.ObjectIdHex(masterId)}, update); err == nil {
	//		return nil
	//	} else {
	//		return errors.New("scanning num > 4 and update status = 2 error")
	//	}
	//
	// }
}

/*
* 上传固件url post到安全团队api
 */
func RemoteDownload(ProjectId int, FileUrl string) (int, error) {

	projectId := strconv.Itoa(ProjectId)

	var jsonStr = []byte(`
{
	"pid":` + projectId + `,
	"url":"` + FileUrl + `"
}`)
	firmwareConfig := service.LoadFirmwareConfig()
	url := firmwareConfig.ScanHost + "/api/remotedownload"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBodyJson, _ := ioutil.ReadAll(resp.Body)

	respBodyMap := []byte(respBodyJson)
	TaskResps := FirmwareCreateProResp{}
	json.Unmarshal(respBodyMap, &TaskResps)

	if TaskResps.Code == 1000 {
		return 1, nil
	} else {
		return -1, errors.New("start task error")
	}
}

type FirmwareGetDownloadResp struct {
	Code int                         `bson:"code"`
	Data FirmwareGetDownloadRespData `bson:"data"`
}

type FirmwareGetDownloadRespData struct {
	Pid    int    `bson:"pid"`
	Status string `bson:"status"`
	Path   string `bson:"path"`
}

/*
 * 获取已post上传的Url下载状态
 * 如果已下载则创建任务 开始任务扫描
 */
func GetRemoteDownload(ProjectId int, MasterId string) string {

	nowTime := custom_util.GetCurrentMilliSecond() / 1000
	projectId := strconv.Itoa(ProjectId)
	firmwareConfig := service.LoadFirmwareConfig()
	url := firmwareConfig.ScanHost + "/api/remotedownload?pid=" + projectId
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBodyJson, _ := ioutil.ReadAll(resp.Body)

	respBodyMap := []byte(respBodyJson)
	DownloadResps := FirmwareGetDownloadResp{}
	json.Unmarshal(respBodyMap, &DownloadResps)

	if DownloadResps.Data.Status == "finished" {
		getFirmwareParams := qmap.QM{
			"e_project_id": ProjectId,
			"e__id":        bson.ObjectIdHex(MasterId),
			// "e_status":    int(4),
		}
		ResponseFile := DownloadResps.Data.Path
		// 判断是否存在此固件记录
		if msg, err := mongo.NewMgoSessionWithCond(common.MC_FIRMWARE_UPLOAD_LOG, getFirmwareParams).GetOne(); err == nil {
			templateId := msg.MustInt("template_id")
			ProjectName := msg.MustString("project_name")
			// FilePath := msg.MustString("tmp_hd_file_path")
			// 创建开始任务
			TaskId, TaskName := CreateTask(ProjectName, ResponseFile, ProjectId, int64(templateId))
			if TaskId > 0 {
				if taskRts, _ := ChangeTaskStatus(TaskId, "start"); taskRts == 1 {
					updateStatus := bson.M{
						"$set": bson.M{
							"task_id":       TaskId,
							"task_name":     TaskName,
							"upload_time":   nowTime,
							"response_file": ResponseFile,
							"status":        6,
						},
					}
					mongoClient := mongo.NewMgoSession(common.MC_FIRMWARE_UPLOAD_LOG)
					if err := mongoClient.Update(bson.M{"_id": bson.ObjectIdHex(MasterId)}, updateStatus); err != nil {
						return "firmware update msg error"
					} else {
						return "success"
					}
				} else {
					return "start task error"
				}
			} else {
				updateStatus := bson.M{
					"$set": bson.M{
						"status": 9,
					},
				}
				mongo.NewMgoSession(common.MC_FIRMWARE_UPLOAD_LOG).Update(bson.M{"_id": bson.ObjectIdHex(MasterId)}, updateStatus)
				return "task create error set status = 9"
			}
		} else {
			return "not found " + MasterId + " this upload msg"
		}
	} else if DownloadResps.Data.Status == "none" {
		updateStatus := bson.M{
			"$set": bson.M{
				"status": 9,
			},
		}
		mongo.NewMgoSession(common.MC_FIRMWARE_UPLOAD_LOG).Update(bson.M{"_id": bson.ObjectIdHex(MasterId)}, updateStatus)
		return "not found push url for firmware to api"
	} else {
		return DownloadResps.Data.Status
	}
}

/*
 * 创建扫描任务
 */
func CreateTask(ProjectName, ResponseFile string, ProjectId int, TemplateId int64) (int, string) {

	proStr := strings.Split(ProjectName, "_")
	TaskName := "TaskName" + proStr[len(proStr)-1]

	projectId := strconv.Itoa(ProjectId)
	templateId := strconv.FormatInt(TemplateId, 10)

	var jsonStr = []byte(`
{
	"task_name":"` + TaskName + `",
	"project_name":"` + ProjectName + `",
	"project_id":` + projectId + `,
	"status":"ready",
	"file_type":"multifile",
	"template_id":` + templateId + `,
	"template_name":"HERMES全扫描",
	"tempfile":"` + ResponseFile + `",
	"workpath":"" 
}`)
	firmwareConfig := service.LoadFirmwareConfig()
	url := firmwareConfig.ScanHost + "/api/tasks"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	// fmt.Println(resp)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBodyJson, _ := ioutil.ReadAll(resp.Body)

	respBodyMap := []byte(respBodyJson)
	TaskResps := FirmwareCreateProResp{}
	json.Unmarshal(respBodyMap, &TaskResps)
	if TaskResps.Code == 1000 {
		var taskId int
		for _, i := range TaskResps.Data {
			taskId = i

		}
		return taskId, TaskName
	} else {
		return -1, "crate task code is " + string(TaskResps.Code)
	}
}

/*
 * 更改固件扫描任务状态 start stop
 */
func ChangeTaskStatus(TaskId int, Status string) (int, error) {
	TaskIdStr := strconv.Itoa(TaskId)
	var jsonStr = []byte(`
{
	"id":` + TaskIdStr + `,
	"operation":"` + Status + `" 
}`)
	firmwareConfig := service.LoadFirmwareConfig()
	url := firmwareConfig.ScanHost + "/api/tasks"
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	respBodyJson, _ := ioutil.ReadAll(resp.Body)

	respBodyMap := []byte(respBodyJson)
	TaskResps := FirmwareCreateProResp{}
	json.Unmarshal(respBodyMap, &TaskResps)
	if TaskResps.Code == 1000 {
		return 1, nil
	} else {
		return -1, errors.New("start task error")
	}
}

func (this *FirmWareData) IsSetFirmware(firmwareMd5 string, projectId int) bool {
	params := qmap.QM{
		"e_project_id":   projectId,
		"e_firmware_md5": firmwareMd5,
	}
	if _, err := mongo.NewMgoSessionWithCond(common.MC_FIRMWARE_UPLOAD_LOG, params).GetOne(); err == nil {
		return true
	} else {
		return false
	}
}

func (this *FirmWareData) IsSetFirmwareForProID(masterId string, taskId, projectId int) bool {
	getFirmwareParams := qmap.QM{
		"e_project_id": projectId,
		"e__id":        bson.ObjectIdHex(masterId),
		"ne_status":    int(0),
	}
	if _, err := mongo.NewMgoSessionWithCond(common.MC_FIRMWARE_UPLOAD_LOG, getFirmwareParams).GetOne(); err == nil {
		codeResult, _ := ChangeTaskStatus(taskId, "stop")
		if int(codeResult) == 1 {
			// log.Println("this result is one")
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

/*
 * 停止扫描 只有状态为扫描中的状态 status=6 or status=8 才可以暂停扫描
 */
func (this *FirmWareData) StopScanning(masterId, userName string, taskId, projectId, userID int) (bool, error) {
	getFirmwareParams := qmap.QM{
		"e_project_id":     projectId,
		"e__id":            bson.ObjectIdHex(masterId),
		"in_status":        []int{5, 6, 8},
		"e_upload_user":    userName,
		"e_upload_user_id": userID,
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_FIRMWARE_UPLOAD_LOG, getFirmwareParams)
	if upRts, err := mongoClient.GetOne(); err == nil {

		// 如果不存在此任务 则 直接更改状态为7 取消扫描
		rtsTaskID := upRts.Int("task_id")
		// rtsTaskName := upRts.String("task_name")
		if rtsTaskID == 0 {
			if err := mongoClient.Update(bson.M{"_id": bson.ObjectIdHex(masterId)}, bson.M{"$set": bson.M{"status": 7}}); err == nil {
				return true, nil
			} else {
				return false, errors.New("task id = 0 and update status = 7 error")
			}
		} else {
			codeResult, _ := ChangeTaskStatus(taskId, "stop")
			if int(codeResult) == 1 {
				// 更新数据状态status=7 取消\暂停扫描状态
				if err := mongoClient.Update(bson.M{"_id": bson.ObjectIdHex(masterId)}, bson.M{"$set": bson.M{"status": 7}}); err == nil {
					return true, nil
				} else {
					return false, errors.New("stop scanning success but update status error")
				}
			} else {
				return false, errors.New("api stop scanning error")
			}
		}

	} else {
		return false, errors.New("not found this message")
	}
}

/*
 * 开始扫描 只有状态为扫描中的状态 status=7 才可以暂停扫描
 */
func (this *FirmWareData) StartScanning(masterId, userName string, taskId, projectId, userID int) (bool, error) {
	getFirmwareParams := qmap.QM{
		"e_project_id":     projectId,
		"e__id":            bson.ObjectIdHex(masterId),
		"e_status":         7,
		"e_upload_user":    userName,
		"e_upload_user_id": userID,
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_FIRMWARE_UPLOAD_LOG, getFirmwareParams)
	if _, err := mongoClient.GetOne(); err == nil {
		// 直接将状态改为status=5 待后台脚本轮询解决
		if err := mongoClient.Update(bson.M{"_id": bson.ObjectIdHex(masterId)}, bson.M{"$set": bson.M{"status": 5}}); err == nil {
			return true, nil
		} else {
			return false, errors.New("update status = 5 error")
		}
	} else {
		return false, errors.New("not found this message")
	}
}

func (this *FirmWareData) GetAll(rawInfo *qmap.QM) (*FirmWareData, error) {
	return this, nil
}

func (this *FirmWareData) GetOne(rawInfo *qmap.QM) (*FirmWareData, error) {
	return this, nil
}
