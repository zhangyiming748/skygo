package mysql_model

import (
	"reflect"
	"testing"
)

func TestGetTestScriptById(t *testing.T) {
	TestingInit()
	testCaseId := []int{2357, 2354, 1802}
	taskUuid := "G5BUEFMB"
	testScriptSlice, err := GetTestScriptById(testCaseId)
	if err != nil {
		t.Errorf("GetTestScriptById Err %s", err)
	}
	testScriptSliceTemp, err := GetTestScriptByTaskUuid(taskUuid)
	if err != nil {
		t.Errorf("GetTestScriptByTaskUuid Err %s", err)
	}
	if reflect.DeepEqual(testScriptSlice, testScriptSliceTemp) {
		t.Errorf("expect %s got %s", testScriptSlice, testScriptSliceTemp)
	}
}
