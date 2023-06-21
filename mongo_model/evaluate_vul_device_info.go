package mongo_model

import (
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
)

type EvaluateVulDeviceInfo struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	TaskId     string        `bson:"task_id" json:"task_id"`         //任务id
	Company    string        `bson:"company" json:"company"`         //车机厂商
	SysVersion string        `bson:"sys_version" json:"sys_version"` //系统版本
	CPUMode    string        `bson:"cpu_mode" json:"cpu_mode"`       //芯片型号
	CPUVersion string        `bson:"cpu_version" json:"cpu_version"` //芯片版本
	Platform   string        `bson:"platform" json:"platform"`       //平台
	SysSdkVer  string        `bson:"sys_sdk_ver" json:"sys_sdk_ver"` //sdk版本
	CarMode    string        `bson:"car_mode" json:"car_mode"`       //车机型号
	Brand      string        `bson:"brand" json:"brand"`             //车机品牌
}

func (this *EvaluateVulDeviceInfo) Create(rawInfo qmap.QM) (*EvaluateVulDeviceInfo, error) {
	if err := mongo.NewMgoSession(common.MC_EVALUATE_VUL_DEVICE_INFO).Insert(rawInfo); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateVulDeviceInfo) BulkDelete(rawIds []string) (*qmap.QM, error) {
	// 删除 测试项
	effectNum := 0
	ids := []string{}
	for _, id := range rawIds {
		ids = append(ids, id)
	}
	if len(ids) > 0 {
		match := bson.M{
			"task_id": bson.M{"$in": ids},
		}
		if changeInfo, err := mongo.NewMgoSession(common.MC_EVALUATE_VUL_DEVICE_INFO).RemoveAll(match); err == nil {
			effectNum = changeInfo.Removed
			// 根据item_id删除 测试项里的漏洞
			new(EvaluateVulnerability).BulkDeleteByItemIds(rawIds)
		} else {
			return nil, err
		}
	}
	return &qmap.QM{"number": effectNum}, nil
}

func (this *EvaluateVulDeviceInfo) GetOne(taskId string) (*qmap.QM, error) {
	params := qmap.QM{
		"e_task_id": taskId,
	}
	return mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_DEVICE_INFO, params).GetOne()
}
