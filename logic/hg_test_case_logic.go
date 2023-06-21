package logic

import (
	"time"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/mongo_model"
)

// 逻辑模块 -- 合规检测工具
// 测试用例
type HgTestCaseLogic struct {
}

func (s *HgTestCaseLogic) Create(code, name string, testType int, version, detail string) (*mongo_model.HgTestCase, error) {
	m := mongo_model.HgTestCase{}
	m.Code = code
	m.Name = name
	m.TestType = testType
	m.Version = version
	m.Detail = detail
	m.CreateTime = time.Now()

	if err := mongo.NewMgoSession(common.McHgTestCase).Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}
