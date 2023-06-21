package logic

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"go.uber.org/zap"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/log"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/mongo_model"
)

// 逻辑模块 -- 合规检测工具
type HgTestTaskLogic struct {
}

// 创建合规检测任务
func (s *HgTestTaskLogic) CreateTask(createUserId int, name string) (*mongo_model.HgTestTask, error) {
	m := mongo_model.HgTestTask{}
	m.Name = name
	m.TaskId = custom_util.GetRandomStringUppercase(6)
	m.Status = mongo_model.HgTestTaskStatusCreate
	m.LastConnectTime = 0
	m.TestCase = make([]mongo_model.TestCase, 0)
	m.HgClientInfo = nil
	m.StatusFlow = &mongo_model.StatusFlow{
		CreateStage:         mongo_model.HgTestTaskFlowYes,
		CreateTime:          int(time.Now().UnixNano() / 1e6),
		ClientInfoStage:     mongo_model.HgTestTaskFlowNo,
		ClientInfoTime:      0,
		ChooseTestCaseStage: mongo_model.HgTestTaskFlowNo,
		ChooseTestCaseTime:  0,
		TestingStage:        mongo_model.HgTestTaskFlowNo,
		TestingTime:         0,
		CompleteStage:       mongo_model.HgTestTaskFlowNo,
		CompleteTime:        0,
	}
	m.TemplateId = ""
	m.TemplateFile = nil

	if err := mongo.NewMgoSession(common.McHgTestTask).Insert(&m); err != nil {
		return nil, err
	}
	return &m, nil
}

// 单个任务详情获取
// taskId: 合规任务的唯一id 6位字符串
func (s *HgTestTaskLogic) GetOne(taskId string) (map[string]interface{}, error) {
	model, err := new(mongo_model.HgTestTask).FindByTaskId(taskId)
	if err != nil {
		return nil, err
	}

	// 把测试用例分开返回
	testCaseAuto := make([]mongo_model.TestCase, 0)
	testCaseInter := make([]mongo_model.TestCase, 0)
	for _, tc := range model.TestCase {
		switch tc.AutoTestLevel {
		case mongo_model.EvaluateTestCaseAutoTestLevelAuto:
			testCaseAuto = append(testCaseAuto, tc)
		case mongo_model.EvaluateTestCaseAutoTestLevelInter:
			testCaseInter = append(testCaseInter, tc)
		}
	}

	// 转map
	bytes, _ := json.Marshal(model)
	data := map[string]interface{}{}
	json.Unmarshal(bytes, &data)

	// 获取测试用例数量信息
	countInfo := s.GetCountInfo(data)

	data["test_case_count"] = countInfo

	// 测试用例分开返回，原来的字段test_case不要了
	data["test_case_inter"] = testCaseInter
	data["test_case_auto"] = testCaseAuto
	delete(data, "test_case")

	return data, nil
}

// 删除文档， taskId唯一定位一个文档
func (s *HgTestTaskLogic) DeleteAll(taskIds ...string) error {
	param := bson.M{
		"task_id": bson.M{
			"$in": taskIds,
		},
	}
	a, e := mongo.NewMgoSession(common.McHgTestTask).RemoveAll(param)
	fmt.Println(a, e)
	return e
}

// 单个任务状态流获取
// taskId: 合规任务的唯一id 6位字符串
func (s *HgTestTaskLogic) GetStatusFlow(taskId string) (map[string]interface{}, error) {
	model, err := new(mongo_model.HgTestTask).FindByTaskId(taskId)
	if err != nil {
		return nil, err
	}

	// 从任务中获取流程信息
	statusFlows := []map[string]interface{}{
		// 1.创建
		{
			"name":       new(mongo_model.HgFlagCreate).PrintChinese(),
			"stage_name": new(mongo_model.HgFlagCreate).PrintFlag(),
			"status":     model.StatusFlow.CreateStage,
			"time":       model.StatusFlow.CreateTime,
		},
		// 2.获取信息
		{
			"name":       new(mongo_model.HgFlagClientInfo).PrintChinese(),
			"stage_name": new(mongo_model.HgFlagClientInfo).PrintFlag(),
			"status":     model.StatusFlow.ClientInfoStage,
			"time":       model.StatusFlow.ClientInfoTime,
		},
		// 3.匹配用例
		{
			"name":       new(mongo_model.HgFlagChooseTestCase).PrintChinese(),
			"stage_name": new(mongo_model.HgFlagChooseTestCase).PrintFlag(),
			"status":     model.StatusFlow.ChooseTestCaseStage,
			"time":       model.StatusFlow.ChooseTestCaseTime,
		},
		// 4.测试
		{
			"name":       new(mongo_model.HgFlagTesting).PrintChinese(),
			"stage_name": new(mongo_model.HgFlagTesting).PrintFlag(),
			"status":     model.StatusFlow.TestingStage,
			"time":       model.StatusFlow.TestingTime,
		},
		// 5.完成
		{
			"name":       new(mongo_model.HgFlagComplete).PrintChinese(),
			"stage_name": new(mongo_model.HgFlagComplete).PrintFlag(),
			"status":     model.StatusFlow.CompleteStage,
			"time":       model.StatusFlow.CompleteTime,
		},
	}

	// 转map
	bytes, _ := json.Marshal(model)
	data := map[string]interface{}{}
	json.Unmarshal(bytes, &data)

	// 获取测试用例数量信息
	countInfo := s.GetCountInfo(data)

	result := map[string]interface{}{}
	result["test_case_count"] = countInfo
	result["status_flow"] = statusFlows
	return result, nil
}

// 单个任务状态流获取
// taskId: 合规任务的唯一id 6位字符串
func (s *HgTestTaskLogic) GetTestCase(taskId string) (map[string]interface{}, error) {
	model, err := new(mongo_model.HgTestTask).FindByTaskId(taskId)
	if err != nil {
		return nil, err
	}

	// 把测试用例分开返回
	testCaseAuto := make([]mongo_model.TestCase, 0)
	testCaseInter := make([]mongo_model.TestCase, 0)
	for _, tc := range model.TestCase {
		switch tc.AutoTestLevel {
		case mongo_model.EvaluateTestCaseAutoTestLevelAuto:
			testCaseAuto = append(testCaseAuto, tc)
		case mongo_model.EvaluateTestCaseAutoTestLevelInter:
			testCaseInter = append(testCaseInter, tc)
		}
	}

	result := map[string]interface{}{}
	result["test_case_inter"] = testCaseInter
	result["test_case_auto"] = testCaseAuto

	return result, nil
}

// 单个任务下测试案例的修改， 提供前端使用，用于交互测试手动更新
// taskId: 合规任务的唯一id 6位字符串
// test_case_id: 合规任务的测试案例的id
func (s *HgTestTaskLogic) UpdateTestCase(taskId string, testCaseId string, request map[string]interface{}) error {
	taskPrt, err := new(mongo_model.HgTestTask).FindByTaskId(taskId)
	if err != nil {
		return err
	}

	status, hasStatus := custom_util.MapHasInt(request, "status")
	detail, hasDetail := custom_util.MapHasString(request, "man_made_result_desc")
	deleteFileId, hasDeleteFileId := custom_util.MapHasString(request, "delete_file_id")

	setList := bson.M{}
	pullList := bson.M{}

	for k, tc := range taskPrt.TestCase {
		if testCaseId == tc.TestCaseId {
			if hasStatus {
				setList[fmt.Sprintf("test_case.%d.status", k)] = status
			}
			if hasDetail {
				setList[fmt.Sprintf("test_case.%d.man_made_result_desc", k)] = detail
			}

			// 删除文件id
			if hasDeleteFileId {
				pullList[fmt.Sprintf("test_case.%d.file_list", k)] = bson.M{
					"id": deleteFileId,
				}

				// 主动删除gridFs中内容
				gf := mongo.GetDefaultMongodbDatabase().GridFS(common.MC_File)
				err = gf.RemoveId(bson.ObjectIdHex(deleteFileId))
			}
		}
	}

	// 更新测试用例状态
	// db.hg_test_task.update({
	//     "task_id": "TH57WS"
	// }, {
	//     "$set": {
	//         "test_case.0.status": 3
	//     }
	// })

	action := bson.M{}
	if len(setList) > 0 {
		action["$set"] = setList
	}
	if len(pullList) > 0 {
		action["$pull"] = pullList
	}

	err = mongo.NewMgoSession(common.McHgTestTask).Update(bson.M{"task_id": taskPrt.TaskId}, action)
	if err != nil {
		return err
	}
	// 更新任务测试状态，如果全部测试用例都完成了，就更新任务的测试结果为完成
	if hasStatus {
		return updateTaskStatusByTestCase(taskPrt.TaskId)
	}
	return nil
}

// 更新状态流
// FlowName参数必须符合switch case中的值
func (s *HgTestTaskLogic) UpdateStatusFlow(taskId string, actionFlag string) error {
	session := mongo.NewMgoSession(common.McHgTestTask)
	param := qmap.QM{
		"e_task_id": taskId,
	}

	// 查询到记录
	model := mongo_model.HgTestTask{}
	err := session.AddCondition(param).One(&model)
	if err != nil {
		return errors.New("记录不存在")
	}

	mp := new(mongo_model.HgTestTask)
	switch actionFlag {
	case new(mongo_model.HgFlagChooseTestCase).PrintFlag():
		return mp.UpdateStatusFlow(taskId, &mongo_model.HgFlagChooseTestCase{})
	case new(mongo_model.HgFlagTesting).PrintFlag():
		return mp.UpdateStatusFlow(taskId, &mongo_model.HgFlagTesting{})
	case new(mongo_model.HgFlagComplete).PrintFlag():
		return mp.UpdateStatusFlow(taskId, &mongo_model.HgFlagComplete{})
	default:
		return errors.New("wrong action flag")
	}
}

// 更新连接设备信息
func (s *HgTestTaskLogic) UpdateClientInfo(taskId string, clientInfo *mongo_model.HgClientInfo) error {
	return new(mongo_model.HgTestTask).UpdateClientInfo(taskId, clientInfo)
}

// 测试用例选取
// 根据taskId获取任务信息，根据任务中的硬件信息，选取匹配它的测试模板，从而选取对应的测试案例
func (s *HgTestTaskLogic) ChooseTestCase(taskId string) error {
	// 1、检测任务必须存在
	modelPtr, err := new(mongo_model.HgTestTask).FindByTaskId(taskId)
	if err == mgo.ErrNotFound {
		return errors.New("任务记录不存在")
	} else if err != nil {
		panic(err)
	}

	// 2、检测任务中“设备”信息必须存在
	if modelPtr.HgClientInfo == nil {
		return errors.New("设备信息未获取")
	}

	// 3、设备信息要匹配到具体的“测试模板”
	tModel, err2 := new(mongo_model.HgTestTemplate).FindByClientInfo(modelPtr.HgClientInfo)
	if err2 != nil {
		return errors.New("测试模板不存在")
	}

	// 4、查询测试用例列表
	caseModels := new(mongo_model.EvaluateTestCase).FindModelsByIds(tModel.HgTestCaseIds)

	// 5、更新（添加）任务的测试案例信息
	testCases := make([]mongo_model.TestCase, 0)
	for _, c := range caseModels {
		if c.AutoTestLevel == mongo_model.EvaluateTestCaseAutoTestLevelAuto || c.AutoTestLevel == mongo_model.EvaluateTestCaseAutoTestLevelInter {
			t := mongo_model.TestCase{
				AutoTestLevel:     c.AutoTestLevel,
				TestCaseId:        c.Id,
				Name:              c.Name,
				TestProcedure:     c.TestProcedure,
				Status:            mongo_model.HgTestTaskCaseStatusPrepare,
				ApiResultDesc:     "",
				ManMadeResultDesc: "",
				FileList:          make([]mongo_model.HgFile, 0),
				BlockList:         c.BlockList, // 测试用例的block列表
			}
			testCases = append(testCases, t)
		}
	}

	// 6.更新
	return new(mongo_model.HgTestTask).UpdateTestCase(taskId, testCases, tModel)
}

// 根据一条任务记录，从里面解析出各类测试用例数量信息
func (s *HgTestTaskLogic) GetCountInfo(data map[string]interface{}) map[string]interface{} {
	// 从任务记录（map格式）中解析出测试案例，
	// 计算任务中的"测试用例"列表
	tcModels := custom_util.FetchMapSlice(data, "test_case")
	// 各个状态阶段得测试用例数量
	var autoAll, autoPrepare, autoRunning, autoAnalyse, autoPass, autoUnpass, autoFail int
	var interAll, interPrepare, interRunning, interAnalyse, interPass, interUnpass, interFail int
	for _, tc := range tcModels {
		if testCaseMap, ok := tc.(map[string]interface{}); ok {
			level := custom_util.FetchMapString(testCaseMap, "auto_test_level")
			switch level {
			case mongo_model.EvaluateTestCaseAutoTestLevelAuto:
				if status, has := custom_util.MapHasInt(testCaseMap, "status"); has {
					switch status {
					case mongo_model.HgTestTaskCaseStatusPrepare:
						autoPrepare++
					case mongo_model.HgTestTaskCaseStatusRunning:
						autoRunning++
					case mongo_model.HgTestTaskCaseStatusAnalyse:
						autoAnalyse++
					case mongo_model.HgTestTaskCaseStatusPass:
						autoPass++
					case mongo_model.HgTestTaskCaseStatusUnPass:
						autoUnpass++
					case mongo_model.HgTestTaskCaseStatusFail:
						autoFail++
					}
				}
				autoAll++
			case mongo_model.EvaluateTestCaseAutoTestLevelInter:
				if status, has := custom_util.MapHasInt(testCaseMap, "status"); has {
					switch status {
					case mongo_model.HgTestTaskCaseStatusPrepare:
						interPrepare++
					case mongo_model.HgTestTaskCaseStatusRunning:
						interRunning++
					case mongo_model.HgTestTaskCaseStatusAnalyse:
						interAnalyse++
					case mongo_model.HgTestTaskCaseStatusPass:
						interPass++
					case mongo_model.HgTestTaskCaseStatusUnPass:
						interUnpass++
					case mongo_model.HgTestTaskCaseStatusFail:
						interFail++
					}
				}
				interAll++
			}
		}
	}

	// test_case_count
	result := map[string]interface{}{
		// “自动化”测试
		"auto_all":     autoAll,
		"auto_prepare": autoPrepare,
		"auto_running": autoRunning,
		"auto_analyse": autoAnalyse,
		"auto_pass":    autoPass,
		"auto_unpass":  autoUnpass,
		"auto_fail":    autoFail,
		// “交互”测试
		"inter_all":     interAll,
		"inter_prepare": interPrepare,
		"inter_running": interRunning,
		"inter_analyse": interAnalyse,
		"inter_pass":    interPass,
		"inter_unpass":  interUnpass,
		"inter_fail":    interFail,
	}
	return result
}

// 更新任务中得测试用例内容
// 自动化测试运行时，工具的长连接会调用此rpc接口，实时更新测试案例状态
func (s *HgTestTaskLogic) AddTestCaseResult(data map[string]interface{}) (interface{}, error) {
	log.GetHttpLogLogger().Info(fmt.Sprintf("[AddTestCaseResult]params:%v", data))

	taskId := custom_util.FetchMapString(data, "task_id")
	if taskId == "" {
		return nil, errors.New("task_id参数不正确")
	}

	testCaseId := custom_util.FetchMapString(data, "test_case_id")
	if testCaseId == "" {
		return nil, errors.New("test_case_id参数不正确")
	}

	taskModelPtr, err := new(mongo_model.HgTestTask).FindByTaskId(taskId)
	if err != nil {
		return nil, errors.New("任务不存在")
	}

	taskCaseModelPtr, err2 := new(mongo_model.HgTestTask).FindTaskTestCase(taskId, testCaseId)
	if err2 != nil {
		return nil, errors.New("任务中测试案例不存在")
	}

	var errBack error
	action := custom_util.FetchMapString(data, "action")

	switch action {
	case "update_status":
		errBack = actionUpdateStatus(data, taskModelPtr)
	case "add_block_result":
		errBack = actionAddBlockResult(data, taskModelPtr, taskCaseModelPtr)
	case "add_block_analysis_result":
		errBack = actionAddBlockAnalysisResult(data, taskModelPtr, taskCaseModelPtr)
	case "end_case":
		errBack = actionEndCase(data, taskModelPtr)
	}

	if errBack != nil {
		log.GetHttpLogLogger().Error(fmt.Sprintf("[AddTestCaseResult]params:%v error:%s", data, errBack.Error()))
		return nil, errBack
	}

	if casePtr, err := new(mongo_model.HgTestTask).FindTaskTestCase(taskId, testCaseId); err != nil {
		return nil, err
	} else {
		return casePtr, nil
	}
}

// 修改测试用例的状态
// 只要有一个测试用例开始测试，任务的进入到测试中状态
// 参数data：
//
//	{
//	    "action": "update_status",
//	    "status" : 1,
//	    "task_id": "xxxxxx",
//	    "task_case_id": "xxxxxx"
//	}
func actionUpdateStatus(data map[string]interface{}, taskPrt *mongo_model.HgTestTask) error {
	status := custom_util.FetchMapInt(data, "status")
	if _, ok := mongo_model.HgTestTaskCaseStatusList[status]; !ok {
		return errors.New("status参数不对")
	}

	testCaseId := custom_util.FetchMapString(data, "test_case_id")
	if testCaseId == "" {
		return errors.New("test_case_id参数不正确")
	}

	// 如果任务的状态流中，测试环境还是no，设置为yes
	setList := bson.M{}
	if taskPrt.StatusFlow.TestingStage == mongo_model.HgTestTaskFlowNo {
		setList["status_flow.testing_stage"] = mongo_model.HgTestTaskFlowYes
		setList["status_flow.testing_time"] = time.Now().Nanosecond() / 1e6
	}

	for k, tc := range taskPrt.TestCase {
		if testCaseId == tc.TestCaseId {
			setList[fmt.Sprintf("test_case.%d.status", k)] = status
		}
	}

	// setList长度为0，说明没用测试案例被匹配到，测试用例不存在
	if len(setList) > 0 {
		action := bson.M{
			"$set": setList,
		}
		return mongo.NewMgoSession(common.McHgTestTask).Update(bson.M{"task_id": taskPrt.TaskId}, action)
	} else {
		return errors.New("任务中测试用例不存在")
	}
}

// 修改测试案例block执行结果
// 参数data：
//
//	{
//	    "action": "add_block_result ",
//	    "block_name": "run",
//	    "task_id": "xxxxxx",
//	    "task_case_id": "xxxxxx",
//	    "data": {
//	        "status": "success/fail",
//	        "reason": "",
//	        "case_id": "1212313",
//	        "block_name": "runBlock",
//	        "result": "{\"result_detail\":\"true\",\"attachment\":{\"name\":\"123.text\",\"file_id\":\"us19d1xasg23\"}}"
//	        }
//	    }
//	}
func actionAddBlockResult(ret map[string]interface{}, taskPrt *mongo_model.HgTestTask, taskCasePrt *mongo_model.TestCase) error {
	testCaseId := custom_util.FetchMapString(ret, "test_case_id")
	if testCaseId == "" {
		return errors.New("测试用例参数不正确")
	}

	// defaultBlockStatus := 0 // todo 要不要block状态
	var fileObj *mongo_model.HgFile = nil

	data := custom_util.FetchMapMap(ret, "data")
	dataStatus := custom_util.FetchMapString(data, "status")
	dataResultString := custom_util.FetchMapString(data, "result")
	if dataResultString == "" {
		return errors.New("返回结果result为空字符串")
	}
	var dataResultMap map[string]interface{}
	if err := json.Unmarshal([]byte(dataResultString), &dataResultMap); err != nil {
		return errors.New("返回结果result转为json异常")
	}

	dataResultDetail := custom_util.FetchMapString(dataResultMap, "result_detail")
	dataResultAttachment := custom_util.FetchMapMap(dataResultMap, "attachment")
	if dataResultAttachment != nil {
		fileObj = new(mongo_model.HgFile)
		fileObj.Name = custom_util.FetchMapString(dataResultAttachment, "name")
		fileObj.Id = custom_util.FetchMapString(dataResultAttachment, "file_id")
	}

	if taskCasePrt.AutoTestLevel == mongo_model.EvaluateTestCaseAutoTestLevelAuto {
		if dataStatus == "success" {
			// defaultBlockStatus = mongo_model.HgTestTaskCaseStatusAnalyse
			taskCasePrt.Status = mongo_model.HgTestTaskCaseStatusAnalyse
		} else {
			// defaultBlockStatus = mongo_model.HgTestTaskCase
			taskCasePrt.Status = mongo_model.HgTestTaskCaseStatusFail
		}
	}

	for k, tc := range taskPrt.TestCase {
		if testCaseId == tc.TestCaseId {
			// 更新测试用例状态
			// db.hg_test_task.update({
			//     "task_id": "TH57WS"
			// }, {
			//     "$set": {
			//         "test_case.0.status": 3
			//     }
			// })
			action := bson.M{
				"$set": bson.M{
					fmt.Sprintf("test_case.%d.status", k):          taskCasePrt.Status,
					fmt.Sprintf("test_case.%d.api_result_desc", k): dataResultDetail,
				},
			}
			// 新增终端返回的截图或文件
			// db.hg_test_task.update({
			//     "task_id": "TH57WS"
			// }, {
			//     "$push": {
			//         "test_case.0.file_list": {
			//             "name": "asd",
			//             "id": "12312323"
			//         }
			//     }
			// })
			if fileObj != nil {
				if _, err := mongo.GridFSOpenId(common.MC_File, bson.ObjectIdHex(fileObj.Id)); err != nil {
					return errors.New("文件ID未查到")
				}

				action["$push"] = bson.M{
					fmt.Sprintf("test_case.%d.file_list", k): bson.M{
						"name": fileObj.Name,
						"id":   fileObj.Id,
					},
				}
			}
			err1 := mongo.NewMgoSession(common.McHgTestTask).Update(bson.M{"task_id": taskPrt.TaskId}, action)
			if err1 != nil {
				log.GetHttpLogLogger().Error(fmt.Sprintf("[AddTestCaseResult]params:%v status:%d", ret, taskCasePrt.Status), zap.String(
					"func", "actionAddBlockResult"))
				return err1
			} else {
				log.GetHttpLogLogger().Info(fmt.Sprintf("[AddTestCaseResult]params:%v status:%d", ret, taskCasePrt.Status), zap.String(
					"func", "actionAddBlockResult"))
			}

			return updateTaskStatusByTestCase(taskPrt.TaskId)
		}
	}

	return errors.New("任务中测试用例不存在")
}

// 修改测试案例block执行结果
// 参数data：
//
//	{
//	    "action": "add_block_analysis_result ",
//	    "block_name": "run",
//	    "task_id": "xxxxxx",
//	    "test_case_id": "xxxxxx",
//		   "analyze_result":{
//	                "case_result":"success",
//	                "logs":"",
//	                "remark":""
//	     }
//	    }
func actionAddBlockAnalysisResult(ret map[string]interface{}, taskPrt *mongo_model.HgTestTask, taskCasePrt *mongo_model.TestCase) error {
	// defaultBlockStatus := 0 // todo 要不要block状态

	dataAnalyzeResult := custom_util.FetchMapMap(ret, "analyze_result")
	caseResult := custom_util.FetchMapString(dataAnalyzeResult, "case_result")

	if taskCasePrt.AutoTestLevel == mongo_model.EvaluateTestCaseAutoTestLevelAuto {
		if caseResult == "success" {
			// defaultBlockStatus = mongo_model.HgTestTaskCaseStatusAnalyse
			taskCasePrt.Status = mongo_model.HgTestTaskCaseStatusPass
		} else {
			// defaultBlockStatus = mongo_model.HgTestTaskCase
			taskCasePrt.Status = mongo_model.HgTestTaskCaseStatusFail
		}
	}

	testCaseId := custom_util.FetchMapString(ret, "test_case_id")
	if testCaseId == "" {
		return errors.New("测试用例参数不正确")
	}

	for k, tc := range taskPrt.TestCase {
		if testCaseId == tc.TestCaseId {
			action := bson.M{
				"$set": bson.M{
					fmt.Sprintf("test_case.%d.status", k): taskCasePrt.Status,
				},
			}
			err1 := mongo.NewMgoSession(common.McHgTestTask).Update(bson.M{"task_id": taskPrt.TaskId}, action)
			if err1 != nil {
				log.GetHttpLogLogger().Error(fmt.Sprintf("[AddTestCaseResult]params:%v status:%d", ret, taskCasePrt.Status), zap.String(
					"func", "actionAddBlockAnalysisResult"))
				return err1
			} else {
				log.GetHttpLogLogger().Info(fmt.Sprintf("[AddTestCaseResult]params:%v status:%d", ret, taskCasePrt.Status), zap.String(
					"func", "actionAddBlockAnalysisResult"))
			}

			return updateTaskStatusByTestCase(taskPrt.TaskId)
		}
	}

	return errors.New("任务中测试用例不存在")
}

// 修改测试案例block执行结果
// 参数data：
//
//	{
//	    "action": "end_case ",
//	    "block_name": "run",
//	    "task_id": "xxxxxx",
//	    "task_case_id": "xxxxxx",
//		"analyze_result":{
//	                "case_result":"success",
//	                "logs":"",
//	                "remark":""
//	            }
//	}
func actionEndCase(ret map[string]interface{}, taskPrt *mongo_model.HgTestTask) error {
	// todo
	return nil
}

// 根据任务中所有测试用例的状态，修改任务状态
// 当测试用例全部完成的时候，更新任务状态
func updateTaskStatusByTestCase(taskId string) error {
	taskPtr, err := new(mongo_model.HgTestTask).FindByTaskId(taskId)
	if err != nil {
		return errors.New("任务不存在")
	}

	i := 0
	for _, tc := range taskPtr.TestCase {
		if tc.Status == mongo_model.HgTestTaskCaseStatusPass || tc.Status == mongo_model.HgTestTaskCaseStatusUnPass {
			i++
		}
	}
	if i == len(taskPtr.TestCase) {
		action := bson.M{
			"$set": bson.M{
				"status": mongo_model.HgTestTaskStatusAutoTest,
			},
		}
		return mongo.NewMgoSession(common.McHgTestTask).Update(bson.M{"task_id": taskPtr.TaskId}, action)
	}
	return nil
}
