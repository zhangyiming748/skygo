package mongo_model

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/orm_mongo"

	"github.com/globalsign/mgo/bson"
)

type ToolTaskData struct {
	Id               primitive.ObjectID `bson:"_id" json:"_id"`
	ProTaskId        string             `bson:"pro_task_id" json:"pro_task_id"`
	TaskName         string             `bson:"task_name" json:"task_name"`
	TaskNumber       string             `bson:"task_number" json:"task_number"`
	ToolName         string             `bson:"tool_name" json:"tool_name"`
	ToolId           string             `bson:"tool_id" json:"tool_id"`
	ToolCategoryName string             `bson:"tool_category_name" json:"tool_category_name"`
	AssetsId         string             `bson:"assets_id" json:"assets_id"`               // 资源ID
	AssetsName       string             `bson:"assets_name" json:"assets_name"`           // 资源名称
	FileId           string             `bson:"file_id" json:"file_id"`                   // 资源ID
	FileName         string             `bson:"file_name" json:"file_name"`               // 资源ID
	AssetsCateName   string             `bson:"assets_cate_name" json:"assets_cate_name"` // 资源类别名称
	CreateUserID     int                `bson:"create_user_id" json:"create_user_id"`
	CreateUserName   string             `bson:"create_user_name" json:"create_user_name"`
	CreateTime       int                `bson:"create_time" json:"create_time"`
	UpdateTime       int                `bson:"update_time" json:"update_time"`
	Status           int                `bson:"status" json:"status"` // 0 删除 1 已完成 2 运行中 3 暂停 4 运行失败
}

type ToolTaskResultBindTest struct {
	Id             primitive.ObjectID `bson:"_id" json:"_id"`
	ProTaskId      string             `bson:"pro_task_id" json:"pro_task_id"`   // 项目任务ID
	ToolTaskId     string             `bson:"tool_task_id" json:"tool_task_id"` // 工具任务ID
	ResultId       string             `bson:"result_id" json:"result_id"`       // 工具任务解析结果ID
	TestId         string             `bson:"test_id" json:"test_id"`           // 测试用例ID
	RecordId       string             `bson:"record_id" json:"record_id"`       // 测试记录id
	CreateUserID   int                `bson:"create_user_id" json:"create_user_id"`
	CreateUserName string             `bson:"create_user_name" json:"create_user_name"`
	CreateTime     int                `bson:"create_time" json:"create_time"`
	UpdateTime     int                `bson:"update_time" json:"update_time"`
	Status         int                `bson:"status" json:"status"` // 0 删除 1 绑定成功 2 解除绑定
}

type LoopholeData struct {
	Id           primitive.ObjectID `bson:"_id,omitempty" json:"_id"`
	Name         string             `bson:"name" json:"name"`                   // 漏洞名称
	Status       int                `bson:"status" json:"status"`               // 漏洞状态（0:未修复 1:已修复 2:重打开 3:未涉及）
	Level        int                `bson:"level"  json:"level"`                // 漏洞级别（ 0:提示 1:低危 2:中危 3:高危 4:严重 ）
	RiskType     string             `bson:"risk_type" json:"risk_type"`         // 风险根源类型（配置 设计 代码 其他）
	Describe     string             `bson:"describe" json:"describe"`           // 漏洞描述
	DateBulletin string             `bson:"date_bulletin" json:"date_bulletin"` // 漏洞公布时间
	DateExposure string             `bson:"date_exposure" json:"date_exposure"` // 漏洞暴露时间
	Sketch       string             `bson:"sketch" json:"sketch"`               // 漏洞简述
	CveType      string             `bson:"cve_type" json:"cve_type"`           // 漏洞类型
	CreateTime   int                `bson:"create_time"`
	TaskNumber   string             `bson:"task_number"`
}

type ToolTaskLinkDataResponse struct {
	Id         primitive.ObjectID `bson:"_id"`
	CreateTime int                `bson:"create_time"`
	UpdateTime int                `bson:"update_time"`
}

type ToolTaskDataResponse struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`
	ProTaskId  string             `bson:"pro_task_id"` // 项目ID
	TaskNumber string             `bson:"task_number"` // 名称
	TaskName   string             `bson:"task_name"`   // 名称
}

/*
 * 检测项目任务处于什么阶段
 * 1:创建、2:任务审核、3:测试、4:报告审核、5:完成
 */
func (this *ToolTaskData) CheckProTaskStatus(evaluateID string) (status int, err error) {
	_id, err := primitive.ObjectIDFromHex(evaluateID)
	params := qmap.QM{
		"e__id": _id,
	}
	w := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_TASK, params)
	if rts, err := w.Get(); err == nil {
		return rts.MustInt("status"), nil
	} else {
		return -1, errors.New("未发现此项目任务")
	}
}

func (this *ToolTaskData) Create(rawInfo *qmap.QM, UserID int, UserName string) (*ToolTaskDataResponse, error) {

	toolName := rawInfo.MustString("tool_name")              // 工具名称
	taskName := rawInfo.MustString("task_name")              // 工具任务名称
	toolId := rawInfo.MustString("tool_id")                  // 工具ID
	categoryName := rawInfo.MustString("category_name")      // 工具分类名称
	assetsCateName := rawInfo.MustString("assets_cate_name") // 资源分类名称
	assetsId := rawInfo.MustString("assets_id")              // 资源ID
	assetsName := rawInfo.MustString("assets_name")          // 资源名称
	fileId := rawInfo.MustString("file_id")                  // 文件ID
	fileName := rawInfo.String("file_name")                  // 文件名称
	nowTime := custom_util.GetCurrentMilliSecond() / 1000
	// nowTimeStr := util.TimestampToStringNoSpace(nowTime)
	proTaskId := rawInfo.MustString("pro_task_id")          // 项目任务ID
	taskNumber := custom_util.RandStringBytesMaskImprSrc(6) // 任务编号

	// 非测试阶段，不允许添加扫描任务
	if proTaskStatus, PTErr := this.CheckProTaskStatus(proTaskId); PTErr == nil && proTaskStatus == 3 {
		this.Id = primitive.NewObjectID()
		// this.TaskName = "task_" + nowTimeStr
		this.TaskName = taskName
		this.TaskNumber = taskNumber
		this.ToolId = toolId
		this.ToolName = toolName
		this.AssetsCateName = assetsCateName
		this.AssetsId = assetsId
		this.AssetsName = assetsName
		this.FileId = fileId
		this.FileName = fileName
		this.ProTaskId = proTaskId

		this.CreateUserName = UserName
		this.CreateUserID = UserID
		this.ToolCategoryName = categoryName
		this.CreateTime = int(nowTime)
		this.UpdateTime = int(nowTime)
		this.Status = 2
		coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_TOOL_TASK)
		if _, err := coll.InsertOne(context.Background(), this); err == nil {
			Rts := ToolTaskDataResponse{
				Id:         this.Id,
				ProTaskId:  this.ProTaskId,
				TaskName:   this.TaskName,
				TaskNumber: this.TaskNumber,
			}
			return &Rts, nil
		} else {
			return nil, err
		}
	} else {
		return nil, errors.New("项目任务非测试阶段,不允许添加扫描任务")
	}

}

// 检测工具工具任务是否已绑定结果
func (this *ToolTaskData) CheckToolTaskIsBind(toolTaskId, proTaskId string) bool {
	params := qmap.QM{
		"e_pro_task_id":  proTaskId,  // 项目任务ID
		"e_tool_task_id": toolTaskId, // 工具任务ID
		"e_status":       1,
	}
	if rtsNum, _ := orm_mongo.NewWidgetWithParams(common.MC_TOOL_TASK_RESULT_BIND_TEST, params).Count(); rtsNum > 0 {
		return true
	} else {
		return false
	}
}

func (this *ToolTaskData) UpdateStatus(masterID, proTaskId string, status int) (*ToolTaskDataResponse, error) {
	_id, _ := primitive.ObjectIDFromHex(masterID)
	params := qmap.QM{
		"e__id":         _id,
		"e_pro_task_id": proTaskId,
		"ne_status":     0,
	}
	w := orm_mongo.NewWidgetWithParams(common.MC_TOOL_TASK, params)
	if rts, err := w.Get(); err == nil {
		if status == 0 {
			// 如果是删除工具任务,则检测此工具是否已经绑定结果信息
			if isBind := this.CheckToolTaskIsBind(masterID, proTaskId); isBind == true {
				return nil, errors.New("该工具任务已经绑定测试结果，无法删除该工具任务")
			}
		}

		nowTime := custom_util.GetCurrentMilliSecond() / 1000
		update := bson.M{
			"$set": bson.M{
				"update_time": int(nowTime),
				"status":      status,
			},
		}
		_id, _ := primitive.ObjectIDFromHex(masterID)
		coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_TOOL_TASK)
		if _, err := coll.UpdateOne(context.Background(), bson.M{"_id": _id}, update); err != nil {
			return nil, err
		} else {
			_id, _ := primitive.ObjectIDFromHex(masterID)
			Rts := ToolTaskDataResponse{
				Id:         _id,
				ProTaskId:  rts.MustString("pro_task_id"),
				TaskName:   rts.MustString("task_name"),
				TaskNumber: rts.MustString("task_number"),
			}
			return &Rts, nil
		}
	} else {
		return nil, err
	}
}

func (this *ToolTaskResultBindTest) BindTestResult(rawInfo *qmap.QM, UserID int, UserName string) (*ToolTaskLinkDataResponse, error) {
	nowTime := custom_util.GetCurrentMilliSecond() / 1000
	// nowTimeStr := custom_util.TimestampToStringNoSpace(nowTime)
	proTaskId := rawInfo.MustString("pro_task_id")
	toolTaskId := rawInfo.MustString("tool_task_id")
	resultId := rawInfo.MustString("result_id")
	TestId := rawInfo.MustString("test_id")
	recordId := rawInfo.MustString("record_id")

	// 检索此数据是否已存在 已经绑定
	params := qmap.QM{
		"e_pro_task_id":  proTaskId,
		"e_tool_task_id": toolTaskId,
		"e_result_id":    resultId,
		"e_test_id":      TestId,
		"e_status":       1,
	}
	if rtsNum, _ := orm_mongo.NewWidgetWithParams(common.MC_TOOL_TASK_RESULT_BIND_TEST, params).Count(); rtsNum > 0 {
		return nil, errors.New("该数据已绑定,请勿重复操作")
	} else {
		this.Id = primitive.NewObjectID()
		this.ProTaskId = proTaskId
		this.ToolTaskId = toolTaskId
		this.ResultId = resultId
		this.TestId = TestId
		this.RecordId = recordId
		this.CreateUserName = UserName
		this.CreateUserID = UserID
		this.CreateTime = int(nowTime)
		this.UpdateTime = int(nowTime)
		this.Status = 1
		coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_TOOL_TASK_RESULT_BIND_TEST)
		if _, err := coll.InsertOne(context.Background(), this); err == nil {
			Rts := ToolTaskLinkDataResponse{
				Id:         this.Id,
				CreateTime: int(nowTime),
				UpdateTime: int(nowTime),
			}
			return &Rts, nil
		} else {
			return nil, err
		}
	}
}
func (this *ToolTaskResultBindTest) UnBindTestResult(rawInfo *qmap.QM, UserID int, UserName string) (*ToolTaskLinkDataResponse, error) {
	masterID, _ := primitive.ObjectIDFromHex(rawInfo.MustString("id"))

	// 检索此数据是否已绑定过
	params := qmap.QM{
		"e__id":    masterID,
		"e_status": 1,
	}
	w := orm_mongo.NewWidgetWithParams(common.MC_TOOL_TASK_RESULT_BIND_TEST, params)
	if rtsNum, _ := w.Count(); rtsNum > 0 {
		nowTime := custom_util.GetCurrentMilliSecond() / 1000
		update := bson.M{
			"$set": bson.M{
				"update_time": int(nowTime),
				"status":      2,
			},
		}
		coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_TOOL_TASK_RESULT_BIND_TEST)
		if _, err := coll.UpdateOne(context.Background(), bson.M{"_id": masterID}, update); err != nil {
			return nil, err
		} else {
			Rts := ToolTaskLinkDataResponse{
				Id:         masterID,
				CreateTime: int(nowTime),
				UpdateTime: int(nowTime),
			}
			return &Rts, nil
		}
	} else {
		return nil, errors.New("该数据不存在")
	}
}

// 脚本更新 根据任务编号判断任务是否存在
func (this *ToolTaskResultBindTest) ScriptIsSetToolTask(taskNumber string) bool {
	getResultParams := qmap.QM{
		"e_task_number": taskNumber,
		"ne_status":     0,
	}
	if _, err := orm_mongo.NewWidgetWithParams(common.MC_TOOL_TASK, getResultParams).Get(); err == nil {
		return true
	} else {
		return false
	}
}

// 脚本更新 根据任务编号更新任务状态为完成
func (this *ToolTaskResultBindTest) ScriptUpdateToolTaskStatus(taskNumber string) error {
	getResultParams := qmap.QM{
		"e_task_number": taskNumber,
		"ne_status":     0,
	}
	nowTime := custom_util.GetCurrentMilliSecond() / 1000
	if taskRts, err := orm_mongo.NewWidgetWithParams(common.MC_TOOL_TASK, getResultParams).Get(); err == nil {
		// 如果存在此任务 序列化数据到漏洞信息表
		// step 1 获取cve漏洞信息数据
		getCveParams := qmap.QM{
			"e_task_id": taskNumber,
		}
		if cveRts, cveErr := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_VUL_INFO, getCveParams).Find(); cveErr == nil {
			// 存在cve消息
			loopData := LoopholeData{}
			for _, val := range cveRts {
				loopData.TaskNumber = taskNumber
				loopData.CreateTime = int(nowTime)
				loopData.Id = primitive.NewObjectID()
				loopData.Name = val["cve_id"].(string)
				loopData.Describe = val["description"].(string)
				loopData.RiskType = val["involve_module"].(string)
				loopData.CveType = val["cve_type"].(string)
				loopData.DateBulletin = val["date_bulletin"].(string)
				loopData.DateExposure = val["date_exposure"].(string)
				loopData.Sketch = val["sketch"].(string)

				// 漏洞状态（0:未修复 1:已修复 2:重打开 3:未涉及）
				cveStatus := int(val["fix_status"].(float64))
				if cveStatus == 1 {
					loopData.Status = 0
				} else if cveStatus == 2 {
					loopData.Status = 1
				} else {
					loopData.Status = cveStatus // 3 未涉及
				}
				// 漏洞级别（ 0:提示 1:低危 2:中危 3:高危 4:严重 ）
				cveLevel := int(val["google_severity_level"].(float64))
				if cveLevel == 0 {
					loopData.Level = 1
				} else if cveLevel == 1 {
					loopData.Level = 2
				} else if cveLevel == 2 {
					loopData.Level = 3
				} else if cveLevel == 3 {
					loopData.Level = 4
				} else {
					loopData.Level = 0
				}
				coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_TOOl_TASK_LOOPHOLE)
				coll.InsertOne(context.Background(), loopData)
			}
			// 更新tool_task 任务状态
			updateData := bson.M{
				"$set": bson.M{
					"status":      1,
					"update_time": nowTime,
				},
			}
			masterId := taskRts.Interface("_id").(primitive.ObjectID)
			updateParams := bson.M{"task_number": taskNumber, "_id": masterId}
			coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_TOOL_TASK)
			if _, err := coll.UpdateOne(context.Background(), updateParams, updateData); err != nil {
				return nil
			} else {
				return errors.New("update error ")
			}
		} else {
			return errors.New("未发现此任务结果")
		}

	} else {
		return errors.New("not found this task number")
	}
}

func (this *ToolTaskResultBindTest) GetToolScanResult(recordId string) ([]map[string]interface{}, error) {
	getResultParams := qmap.QM{
		"e_record_id": recordId,
		"e_status":    1,
	}
	if resultBinds, err := orm_mongo.NewWidgetWithParams(common.MC_TOOL_TASK_RESULT_BIND_TEST, getResultParams).Find(); err == nil {
		scanIds := []interface{}{}
		for _, item := range resultBinds {
			var itemQM qmap.QM = item
			_id, _ := primitive.ObjectIDFromHex(itemQM.String("result_id"))
			scanIds = append(scanIds, _id)
		}
		query := qmap.QM{
			"in__id": scanIds,
		}
		if scanResults, err := orm_mongo.NewWidgetWithParams(common.MC_TOOl_TASK_LOOPHOLE, query).Find(); err == nil {
			return scanResults, err
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
