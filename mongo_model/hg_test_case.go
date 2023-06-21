package mongo_model

import (
	"time"

	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
)

type HgTestCase struct {
	ID         bson.ObjectId `bson:"_id,omitempty" json:"_id"`       // 记录_id
	Code       string        `bson:"code" json:"code"`               // 用例编号
	Name       string        `bson:"name" json:"name"`               // 用例名称
	TestType   int           `bson:"test_type" json:"test_type"`     // 用例类型
	Version    string        `bson:"version" json:"version"`         // 用例版本
	Detail     string        `bson:"detail" json:"detail"`           // 用例说明
	CreateTime time.Time     `bson:"create_time" json:"create_time"` // 创建时间
}

const HgTestCaseTestTypeAuto = 1  // 合规测试用例 -- 自动化测试
const HgTestCaseTestTypeInter = 2 // 合规测试用例 -- 交互测试

// 根据一组ID，查询出对应的一组测试案例
func (h *HgTestCase) FindListByIds(ids []bson.ObjectId) []HgTestCase {
	session := mongo.NewMgoSession(common.McHgTestCase)
	params := qmap.QM{
		"in_uuid": ids,
	}
	session.AddCondition(params)

	models := make([]HgTestCase, 0)
	session.Session.Find(nil).All(&models)

	return models
}
