package logic_test

import (
	"fmt"
	"skygo_detection/guardian/src/config/watcher"
	"skygo_detection/logic"
	"skygo_detection/mongo_model"
	"testing"
)

/*
*
模拟任务创建
go test -v -count=1 logic_test.go -test.run Create
*/
func TestCreate(t *testing.T) {
	watcher.InitWatch("../config/dev_sww")
	name := "xxx任务"
	uuid := 12

	a, b := new(logic.HgTestTaskLogic).CreateTask(uuid, name)
	fmt.Println(a, b)
}

/*
*
go test -v -count=1 logic_test.go -test.run UpdateStatusFlow
*/
func TestUpdateStatusFlow(t *testing.T) {
	watcher.InitWatch("../config/op_qa")
	uuid := "U921SA"
	flowName := "client_info"

	new(logic.HgTestTaskLogic).UpdateStatusFlow(uuid, flowName)
}

/*
*
模式车机进行车机的测试用例匹配
go test -v -count=1 logic_test.go -test.run HgTaskChoose
*/
func TestHgTaskChoose(t *testing.T) {
	watcher.InitWatch("../config/dev_sww")
	taskId := "TT24L0"

	e := new(logic.HgTestTaskLogic).ChooseTestCase(taskId)
	fmt.Println(e, "1213")
}

/*
*
go test -v -count=1 logic_test.go -test.run HgTaskUpdateConnected
*/
func TestHgTaskUpdateConnected(t *testing.T) {
	watcher.InitWatch("../config/op_qa")
	uuid := "U921SA"

	e := new(mongo_model.HgTestTask).UpdateConnected(uuid)
	fmt.Println(e, "1213")
}

/*
*
go test -v -count=1 logic_test.go -test.run HgTestCaseCreate
*/
func TestHgTestCaseCreate(t *testing.T) {
	watcher.InitWatch("../config/op_qa")

	code := "TC0001"
	name := "测试用例"
	testType := 1
	version := "1.0.0"
	detail := "测试说明xxxxxxxxxxx测试说明xxxxxxxxxxx测试说明xxxxxxxxxxx测试说明xxxxxxxxxxx测试说明xxxxxxxxxxx"

	new(logic.HgTestCaseLogic).Create(code, name, testType, version, detail)
}

/*
*
基于测试用例列表，创建一个模板
go test -v -count=1 logic_test.go -test.run HgTestTemplateCreate
*/
func TestHgTestTemplateCreate(t *testing.T) {
	watcher.InitWatch("../config/op_qa")

	name := "0001模板"
	osType := "android"
	osVersion := "9.0"
	cpu := "64"
	ids := []string{
		"TC999030M001",
		"TC999030M108",
		"TC999030M109",
		"TC999030M110",
	}
	new(logic.HgTestTemplateLogic).Create(name, osType, osVersion, cpu, ids)
}

/*
*
模拟上传硬件信息
go test -v -count=1 logic_test.go -test.run HgTestTaskClientInfo
*/
func TestHgTestTaskClientInfo(t *testing.T) {
	watcher.InitWatch("../config/dev_sww")

	uuid := "TT24L0"
	clientInfo := &mongo_model.HgClientInfo{
		Cpu:       "64",
		OsType:    "android",
		OsVersion: "9.0",
	}
	new(logic.HgTestTaskLogic).UpdateClientInfo(uuid, clientInfo)
}

/**
测试解压缩zip
go test -v -count=1 logic_test.go -test.run TestHgTestUnzip
*/
// func TestHgTestUnzip(t *testing.T) {
// 	watcher.InitWatch("../config/dev_sww")
//
// 	new(logic.HgTestTemplateLogic).ZipDeCompress("60a0a706e1382389e99f69a6")
// }
