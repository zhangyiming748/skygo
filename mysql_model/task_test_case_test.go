package mysql_model

import (
	"go.uber.org/zap"
	"reflect"
	"skygo_detection/common"
	"skygo_detection/custom_util/clog"
	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/lib/common_lib/mysql"
	"testing"
)

func TestingInit() {
	// TODO
	sys_service.InitConfigWatcher("qa", "../config/qa/config.tml")
	mysql.InitMysqlEngine()
}

func TestGetTestCaseIdByTaskUuid(t *testing.T) {
	TestingInit()
	taskUuid := "G5BUEFMB"
	testCaseIds, err := GetTestCaseIdByTaskUuid(taskUuid)
	if err != nil {
		t.Errorf("GetTestCaseIdByTaskUuid Err %s", err)
	}
	for _, v := range testCaseIds {
		taskTestCase, err := TaskTestCaseFindOne(taskUuid, v)
		if err != nil {
			t.Errorf("TaskTestCaseFindOne Err %s", err)
		}
		if common.TOOL_HG_ANDROID_SCANNER != taskTestCase.TestTool {
			t.Errorf("TaskTestCaseFindOne TestTool Err %s", err)
		}
	}
	clog.Info("got %s", zap.Any("value", testCaseIds))
}

func TestGetTestScriptByTaskUuid(t *testing.T) {
	TestingInit()
	taskUuid := "G5BUEFMB"
	_ = [3]int{2357, 2354, 1802}
	expect := []string{
		"61b32421e830c630780178ee",
		"61b1a96ee830c632c3b50547",
		"61af1ffab1b91923d0225cd6|61b1ba23e830c667aaa6b8bf"}
	testScriptSlice, err := GetTestScriptByTaskUuid(taskUuid)
	if err != nil {
		t.Errorf("GetTestScriptByTaskUuid Err %s", err)
	}
	if reflect.DeepEqual(expect, testScriptSlice) {
		t.Errorf("expect %s got %s", expect, testScriptSlice)
	}
	clog.Info("got %s", zap.Any("value", testScriptSlice))
}
