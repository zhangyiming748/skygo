package hg_scanner_logic

import (
	"archive/zip"
	"bytes"
	"errors"
	"github.com/globalsign/mgo/bson"
	"go.uber.org/zap"
	"io/ioutil"
	"skygo_detection/common"
	"skygo_detection/custom_util/clog"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/mysql_model"
	"strings"
	"sync"
)

var (
	ErrorGetTestCaseId      = errors.New("获取测试用例错误")
	ErrorIllegalTestCaseIds = errors.New("测试用例不存在")
	ErrorGetTestScriptById  = errors.New("获取测试脚本错误")
	ErrorTestScriptSlice    = errors.New("测试脚本不存在")
)

type HandleFunc func()

// 获取测试脚本
func GetTestScript(taskUuid string) (testScriptSlice []string, err error) {
	// 获取测试用例ID
	testCaseIds, err := mysql_model.GetTestCaseIdByTaskUuid(taskUuid)
	if err != nil {
		return testScriptSlice, ErrorGetTestCaseId
	}
	if len(testCaseIds) <= 0 {
		return testScriptSlice, ErrorIllegalTestCaseIds
	}
	// 获取测试脚本
	testScriptSlice, err = mysql_model.GetTestScriptById(testCaseIds)
	if err != nil {
		return testScriptSlice, ErrorGetTestScriptById
	}
	if len(testScriptSlice) <= 0 {
		return nil, ErrorTestScriptSlice
	}
	return testScriptSlice, nil
}

// 写文件
func writeFile(zipWriter *zip.Writer, firstTestScript string, wg *sync.WaitGroup) {
	defer wg.Done()
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			clog.Error("writeFile recover Err", zap.Any("error", recoverErr))
		}
	}()

	file, err := mongo.GridFSOpenId(common.MC_File, bson.ObjectIdHex(firstTestScript))
	if err != nil {
		clog.Error("writeFile mongo.GridFSOpenId Err", zap.Any("error", err))
		panic(err)
	}

	fileWriter, err := zipWriter.Create(file.Name())
	if err != nil {
		clog.Error("writeFile zipWriter.Create Err", zap.Any("error", err))
		panic(err)
	}

	bs, err := ioutil.ReadAll(file)
	if err != nil {
		clog.Error("writeFile ioUtil.ReadAll Err", zap.Any("error", err))

		panic(err)
	}

	_, err = fileWriter.Write(bs)
	if err != nil {
		clog.Error("writeFile fileWriter.Write  Err", zap.Any("error", err))
		panic(err)
	}
}

// 打包文件
func ArchiveFile(testScriptSlice []string) (b bytes.Buffer, err error) {
	zipWriter := zip.NewWriter(&b)
	defer zipWriter.Close()
	var wg sync.WaitGroup
	for _, v := range testScriptSlice {
		// 当测试脚本为多个时取第一个
		firstTestScript := strings.Split(v, "|")[0]
		if len(firstTestScript) > 0 {
			wg.Add(1)
			go writeFile(zipWriter, firstTestScript, &wg)
		}
	}
	wg.Wait()
	return b, nil
}
