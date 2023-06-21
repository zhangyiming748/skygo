package mongo_model

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
)

// 合规测试
// 设备对应的测试模板
type HgTestTemplate struct {
	ID            bson.ObjectId `bson:"_id,omitempty" json:"_id"` // 记录_id
	Name          string        `bson:"name" json:"name"`         // 测试模板名称
	HgClientInfo  HgClientInfo  `bson:"hg_client_info" json:"hg_client_info"`
	HgTestCaseIds []string      `bson:"hg_test_case_ids" json:"hg_test_case_ids"` // 测试用例_id列表
	File          HgFile        `bson:"file" json:"file"`                         // 文件
	CreateTime    time.Time     `bson:"create_time" json:"create_time"`           // 创建时间
}

type HgFile struct {
	Name string `bson:"name" json:"name"` // 测试工具名字、测试脚本名称、测试环境搭建示意图名称
	Id   string `bson:"id" json:"id"`     // 测试工具id、测试脚本id、测试环境搭建示意图名称id
}

// 根据设备信息，获取一条对应的测试模板
func (h *HgTestTemplate) FindByClientInfo(info *HgClientInfo) (*HgTestTemplate, error) {
	params := qmap.QM{
		"hg_client_info.os_type":    info.OsType,
		"hg_client_info.os_version": info.OsVersion,
		"hg_client_info.cpu":        info.Cpu,
	}
	model := HgTestTemplate{}

	session := mongo.NewMgoSession(common.McHgTestTemplate)
	err := session.AddCondition(params).One(&model)
	if err != nil {
		return nil, err
	}
	return &model, nil
}
