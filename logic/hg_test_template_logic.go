package logic

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/mongo_model"
)

// 逻辑模块 -- 合规检测工具
// 检测模板
type HgTestTemplateLogic struct {
}

// 创建合规测试模板
// ids 对应集合 evaluate_test_case中的_id（string类型）
func (s *HgTestTemplateLogic) Create(name, osType, osVersion, cpu string, ids []string) (*mongo_model.HgTestTemplate, error) {
	m := mongo_model.HgTestTemplate{}
	m.Name = name
	m.CreateTime = time.Now()
	m.HgClientInfo = mongo_model.HgClientInfo{
		OsVersion: osVersion,
		OsType:    osType,
		Cpu:       cpu,
	}
	m.HgTestCaseIds = ids

	// 测试用例id得到测试用例记录
	testCaseModels := new(mongo_model.EvaluateTestCase).FindModelsByIds(ids)

	// 写入zip文件中的内容
	var b bytes.Buffer
	writer := zip.NewWriter(&b)

	// 从测试用例中查询所有要打包的文件
	for _, model := range testCaseModels {
		// 只取第一个测试脚本文件
		for _, fileInfo := range model.TestScript {
			// 根据_id得到evaluate_test_case集合中的文件信息，打包得到一个压缩包
			fileIdStr := fileInfo.Value // 文件id
			file, err := mongo.GridFSOpenId(common.MC_File, bson.ObjectIdHex(fileIdStr))
			if err != nil {
				panic(err)
			}

			// 文件名要统一修改为测试用_id做为名称的文件，保持后缀不变
			index := strings.LastIndex(file.Name(), ".")
			tail := file.Name()[index:]
			fileName := model.Id + tail

			// 打包
			fileWriter, err := writer.Create(fileName)
			if err != nil {
				// todo 打印日志
				panic(err)
			}

			bs, err := ioutil.ReadAll(file)
			if err != nil {
				panic(err)
			}

			_, err = fileWriter.Write(bs)
			if err != nil {
				panic(err)
			}

			break
		}
	}

	if err := writer.Close(); err != nil {
		panic(err)
	}

	// 测试用例的脚本文件打包zip
	zipFileName := fmt.Sprintf("%s.zip", name)
	fileId, err := mongo.GridFSUpload(common.MC_File, zipFileName, b.Bytes())
	if err != nil {
		panic(err)
	}

	m.File.Name = zipFileName
	m.File.Id = fileId
	if err := mongo.NewMgoSession(common.McHgTestTemplate).Insert(&m); err != nil {
		return nil, err
	}

	return &m, nil
}

// func (s *HgTestTemplateLogic) ZipDeCompress(id string) error {
// 	file, err := mongo.GridFSOpenId(common.MC_File, bson.ObjectIdHex(id))
// 	if err != nil {
// 		fmt.Println(err, "8000")
// 	}
//
// 	dataBytes, err := ioutil.ReadAll(file)
//
// 	reader := bytes.NewReader(dataBytes)
//
// 	r, _ := zip.NewReader(reader, reader.Size())
//
//
// 	for _, a := range r.File {
//
// 	}
// 	return nil
// }
