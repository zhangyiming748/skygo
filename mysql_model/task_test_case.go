package mysql_model

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"skygo_detection/custom_util/clog"
	"time"

	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/mongo_model_tmp"
)

// 任务下的测试用例
type TaskTestCase struct {
	Id               int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	TaskId           int    `xorm:"not null comment('任务表id') INT(11)" json:"task_id"`
	TaskName         string `xorm:"not null comment('任务名字') VARCHAR(255)" json:"task_name"`
	TaskUuid         string `xorm:"not null comment('全局唯一id') VARCHAR(255)" json:"task_uuid"`
	TestCaseId       int    `xorm:"not null comment('测试用例表id') INT(11)" json:"test_case_id"`
	CaseUuid         string `xorm:"not null comment('测试用例表id') INT(11)" json:"case_uuid"`
	TestCaseName     string `xorm:"not null comment('测试用例名称') VARCHAR(255)" json:"test_case_name"`
	ToolTaskId       string `xorm:"not null comment('扫描任务id') VARCHAR(255)" json:"tool_task_id"`
	TestToolName     string `xorm:"not null comment('测试工具名称') VARCHAR(255)" json:"test_tool_name"`
	TestTool         string `xorm:"not null comment('测试工具类型,firmware_scanner/vul_scanner/hg_scanner') int(11)" json:"test_tool"`
	FileId           string `xorm:"not null comment('固件文件id') VARCHAR(255)" json:"file_id"`
	TemplateId       int    `xorm:"not null comment('固件模板id') int(11)" json:"template_id"`
	AutoTestLevel    int    `xorm:"not null comment('自动化测试程度 1人工 2半自动化 3自动化') TINYINT(3)"`
	TestResultStatus int    `xorm:"not null comment('测试结果(1:通过 2:未通过)') TINYINT(3)" json:"test_result_status"`
	ActionStatus     int    `xorm:"not null comment('执行状态、测试用例状态') TINYINT(3)" json:"action_status"`
	LogResult        int    `xorm:"not null comment('测试过程记录') TINYINT(3)" json:"log_result"` // 【自动化】测试用例，测试完毕后自动填充测试过程记录字段：xxx时间使用xxx工具进行测试。（填充工具日志）
	CreateTime       int    `xorm:"not null comment('创建时间') INT(11)" json:"create_time"`
	UpdateTime       int    `xorm:"not null comment('更新时间') INT(11)" json:"update_time"`
	CompleteTime     int    `xorm:"not null comment('完成时间') INT(11)" json:"complete_time"`
	CasePriority     int    `xorm:"not null comment('用例执行级别') INT(1)" json:"case_priority"`
	CaseResult       string `xorm:"not null comment('任务测试用例结果json格式存储') VARCHAR(1024)" json:"case_result"`
	CaseResultFile   string `xorm:"not null comment('测试用例结果的文件或者图片id') VARCHAR(255)" json:"case_result_file"`
	DemandId         int    `xorm:"not null comment('安全需求id') INT(11)" json:"demand_id"`
	TestProcedure    string `xorm:"comment('测试步骤') TEXT" json:"test_procedure"`
	TestAttachment   string `xorm:"comment('测试附件') VARCHAR(1024)" json:"test_attachment"`
	TestParam        string `xorm:"not null comment('测试参数之前的block_list') VARCHAR(255)"`
	ModuleId         string `xorm:"comment('测试组件/测试分类') VARCHAR(255)"`
	TaskParam        string `xorm:"not null comment('任务参数') VARCHAR(255)"`
}

func (this *TaskTestCase) Create() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}

func (this *TaskTestCase) Update(cols ...string) (int64, error) {
	return mysql.GetSession().Table(this).ID(this.Id).Cols(cols...).Update(this)
}

func (this *TaskTestCase) Remove() (int64, error) {
	return mysql.GetSession().ID(this.Id).Delete(this)
}

// 任务的测试用例列表，查询时要连场景表、用例表等，因此要使用视图
type TaskTestCaseView struct {
	Id               int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	TaskId           int    `xorm:"not null comment('任务表id') INT(11)" json:"task_id"`
	TaskUuid         string `xorm:"not null comment('任务uuid') VARCHAR(255)" json:"task_uuid"`
	TestCaseId       int    `xorm:"not null comment('测试用例表id') INT(11)" json:"test_case_id"`
	TestCaseName     string `xorm:"not null comment('测试用例名称') VARCHAR(255)" json:"test_case_name"`
	ScenarioName     string `xorm:"not null comment('测试用例名称') VARCHAR(255)" json:"scenario_name"`
	ModuleId         int    `xorm:"not null comment('完成时间') INT(11)" json:"module_id"`
	AutoTestLevel    int    `xorm:"not null comment('测试结果') TINYINT(3)" json:"auto_test_level"`
	TestTools        int    `xorm:"not null comment('测试结果') TINYINT(3)" json:"test_tools"`
	TestResultStatus int    `xorm:"not null comment('测试结果') TINYINT(3)" json:"test_result_status"`
	ActionStatus     int    `xorm:"not null comment('完成时间') INT(11)" json:"action_status"`
	TaskParam        string `xorm:"not null comment('任务参数') VARCHAR(255)" json:"task_param"`
}

// 任务测试用例表中，根据任务task_uuid、测试用例id查询出一条记录
func TaskTestCaseFindOne(taskUuid string, testCaseId int) (*TaskTestCase, error) {
	model := TaskTestCase{}
	has, err := mysql.GetSession().Where("task_uuid = ?", taskUuid).
		And("test_case_id = ?", testCaseId).Get(&model)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("not found")
	}
	return &model, err
}

// 任务测试用例表中，根据task_uuid查询其下所有测试用例
func TaskTestCaseFindAll(taskUuid string) []TaskTestCase {
	models := make([]TaskTestCase, 0)
	mysql.GetSession().Where("task_uuid = ?", taskUuid).Find(&models)
	return models
}

// 根据测试用例拷贝出一份任务下的测试用例
func CopyTestCases(testCases []KnowledgeTestCase, parentTask *Task) []TaskTestCase {
	now := int(time.Now().Unix())
	taskTestCases := make([]TaskTestCase, 0)
	var hgScannerTask, vulScannerTask, firmwareScannerTask100, firmwareScannerTask101, firmwareScannerTask102, firmwareScannerTask103, firmwareScannerTask104 string
	for _, testCase := range testCases {
		taskTestCase := new(TaskTestCase)
		taskTestCase.TaskId = parentTask.Id
		taskTestCase.TaskName = parentTask.Name
		taskTestCase.ActionStatus = common.CASE_STATUS_READY
		taskTestCase.TestCaseId = testCase.Id
		taskTestCase.TaskUuid = parentTask.TaskUuid
		taskTestCase.CaseUuid = testCase.CaseUuid
		taskTestCase.TestCaseName = testCase.Name
		taskTestCase.AutoTestLevel = testCase.AutoTestLevel
		taskTestCase.TestToolName = testCase.TestToolName
		taskTestCase.TestTool = testCase.TestTool
		taskTestCase.DemandId = testCase.DemandId
		taskTestCase.TestParam = testCase.TestParam
		taskTestCase.TaskParam = testCase.TaskParam
		taskTestCase.ModuleId = testCase.ModuleId
		taskTestCase.CasePriority = common.CASE_PRIORITY_DEFAULT
		// 获取测试工具任务参数
		paramsStr := testCase.TaskParam
		params := orm.QueryStrToMap(paramsStr)
		if templateId, ok := params["template_id"]; ok {
			taskTestCase.TemplateId = templateId.(int)
		}
		taskTestCase.CreateTime = now
		taskTestCase.UpdateTime = now
		switch taskTestCase.TestTool {
		case common.TOOL_HG_ANDROID_SCANNER:
			taskTestCase.TestTool = common.TOOL_HG_ANDROID_SCANNER
			if parentTask.NeedConnected == common.TASK_CONNECT_STATUS_UNCONNECTED {
				// 如何任务没有设备连接，则安卓合规任务相关的测试用例无效
				taskTestCase.ActionStatus = common.CASE_STATUS_INVALID
			} else {
				if hgScannerTask == "" {
					if taskId, err := CreateSubTask(common.TOOL_HG_ANDROID_SCANNER, taskTestCase, parentTask); err == nil {
						hgScannerTask = taskId
					} else {
						panic(err)
					}
				}
				taskTestCase.ToolTaskId = hgScannerTask
			}
		case common.TOOL_FIRMWARE_SCANNER:
			taskTestCase.TestTool = common.TOOL_FIRMWARE_SCANNER
			templateId := 0
			if taskParam, err := qmap.NewWithString(taskTestCase.TaskParam); err == nil {
				if val, has := taskParam.TryInt("template_id"); has {
					templateId = val
				}
			}
			if templateId == 0 {
				taskTestCase.ActionStatus = common.CASE_STATUS_FAIL
				taskTestCase.TestResultStatus = common.CASE_TEST_STATUS_UNPASS
				taskTestCase.CaseResult = "测试用例参数中未发现固件模板(template_id),固件扫描任务创建失败"
			} else {
				switch templateId {
				case common.FIRMWARE_TEMPLATE_ID_100:
					taskTestCase.TemplateId = 100
					if firmwareScannerTask100 == "" {
						if taskId, err := CreateSubTask(common.TOOL_FIRMWARE_SCANNER, taskTestCase, parentTask); err == nil {
							firmwareScannerTask100 = taskId
						} else {
							panic(err)
						}
					}
					taskTestCase.ToolTaskId = firmwareScannerTask100
				case common.FIRMWARE_TEMPLATE_ID_101:
					taskTestCase.TemplateId = 101
					if firmwareScannerTask101 == "" {
						if taskId, err := CreateSubTask(common.TOOL_FIRMWARE_SCANNER, taskTestCase, parentTask); err == nil {
							firmwareScannerTask101 = taskId
						} else {
							panic(err)
						}
					}
					taskTestCase.ToolTaskId = firmwareScannerTask101
				case common.FIRMWARE_TEMPLATE_ID_102:
					taskTestCase.TemplateId = 102
					if firmwareScannerTask102 == "" {
						if taskId, err := CreateSubTask(common.TOOL_FIRMWARE_SCANNER, taskTestCase, parentTask); err == nil {
							firmwareScannerTask102 = taskId
						} else {
							panic(err)
						}
					}
					taskTestCase.ToolTaskId = firmwareScannerTask102
				case common.FIRMWARE_TEMPLATE_ID_103:
					taskTestCase.TemplateId = 103
					if firmwareScannerTask103 == "" {
						if taskId, err := CreateSubTask(common.TOOL_FIRMWARE_SCANNER, taskTestCase, parentTask); err == nil {
							firmwareScannerTask103 = taskId
						} else {
							panic(err)
						}
					}
					taskTestCase.ToolTaskId = firmwareScannerTask103
				case common.FIRMWARE_TEMPLATE_ID_104:
					taskTestCase.TemplateId = 104
					if firmwareScannerTask104 == "" {
						if taskId, err := CreateSubTask(common.TOOL_FIRMWARE_SCANNER, taskTestCase, parentTask); err == nil {
							firmwareScannerTask104 = taskId
						} else {
							panic(err)
						}
					}
					taskTestCase.ToolTaskId = firmwareScannerTask104
				default:
					taskTestCase.ActionStatus = common.CASE_STATUS_COMPLETED
					taskTestCase.TestResultStatus = common.CASE_TEST_STATUS_UNPASS
					taskTestCase.CaseResult = fmt.Sprintf("未知固件模板id(template_id):%d,固件扫描任务创建失败", templateId)
				}
			}
		case common.TOOL_VUL_SCANNER:
			taskTestCase.TestTool = common.TOOL_VUL_SCANNER
			if parentTask.NeedConnected == common.TASK_CONNECT_STATUS_UNCONNECTED {
				// 如何任务没有设备连接，则车机漏扫相关的测试用例无效
				taskTestCase.ActionStatus = common.CASE_STATUS_INVALID
			} else {
				if vulScannerTask == "" {
					// 创建车机漏扫任务
					if taskId, err := CreateSubTask(common.TOOL_VUL_SCANNER, taskTestCase, parentTask); err == nil {
						vulScannerTask = taskId
					} else {
						panic(err)
					}
				}
				taskTestCase.ToolTaskId = vulScannerTask
			}
		}
		_, err := taskTestCase.Create()
		if err != nil {
			panic(err)
		}
		taskTestCases = append(taskTestCases, *taskTestCase)
	}
	return taskTestCases
}

func CreateSubTask(taskType string, taskTestCase *TaskTestCase, parentTask *Task) (string, error) {
	switch taskType {
	case common.TOOL_FIRMWARE_SCANNER:
		pieceVersion, err := new(AssetTestPieceVersion).FindById(parentTask.PieceVersionId)
		if err != nil {
			return "", err
		}
		// 创建固件扫描任务
		firmwareTask := new(FirmwareTask)
		firmwareTask.TaskId = taskTestCase.TaskId
		firmwareTask.Name = fmt.Sprintf("固件子任务_%s_%d", taskTestCase.TaskName, custom_util.GetCurrentMilliSecond())
		firmwareTask.FileId = pieceVersion.FirmwareFileUuid
		switch pieceVersion.FirmwareDeviceType {
		case common.DEVICE_TYPE_GW:
			firmwareTask.DeviceType = common.DEVICE_TYPE_GW_NAME
		case common.DEVICE_TYPE_ECU:
			firmwareTask.DeviceType = common.DEVICE_TYPE_ECU_NAME
		case common.DEVICE_TYPE_IVI:
			firmwareTask.DeviceType = common.DEVICE_TYPE_IVI_NAME
		}
		firmwareTask.FirmwareVersion = pieceVersion.Version
		firmwareTask.FirmwareName = pieceVersion.FirmwareName
		firmwareTask.TemplateId = taskTestCase.TemplateId
		firmwareTask.CreateTime = int(time.Now().Unix())
		firmwareTask.UpdateTime = int(time.Now().Unix())
		firmwareTask.Status = common.FIRMWARE_STATUS_PROJECT_CREATE
		if _, err := firmwareTask.Create(); err == nil {
			new(ScannerTask).TaskInsert(firmwareTask.Id, firmwareTask.Name, common.TOOL_FIRMWARE_SCANNER)
			return fmt.Sprintf("%d", firmwareTask.Id), nil
		} else {
			return "", err
		}
	case common.TOOL_HG_ANDROID_SCANNER:
		// 创建合规扫描任务
		hgTask := new(HgTestTask)
		hgTask.Name = fmt.Sprintf("合规子任务_%s_%d", taskTestCase.TaskName, custom_util.GetCurrentMilliSecond())
		hgTask.TaskUuid = taskTestCase.TaskUuid
		hgTask.Status = common.HG_TEST_TASK_STATUS_CREATE
		hgTask.CreateTime = int(time.Now().Unix())
		if err := hgTask.Create(); err == nil {
			return fmt.Sprintf("%d", hgTask.Id), nil
		} else {
			return "", err
		}
	case common.TOOL_VUL_SCANNER:
		// 创建车机漏扫任务
		req := qmap.QM{
			"name":           fmt.Sprintf("合规子任务_%s_%d", taskTestCase.TaskName, custom_util.GetCurrentMilliSecond()),
			"parent_task_id": parentTask.Id,
		}
		if result, err := new(mongo_model_tmp.EvaluateVulTask).Create(parentTask.TaskUuid, req); err == nil {
			// todo mysql里也创建任务，目前这个任务还没有串联起来，因为客户端的上传结果还是存mongo库里边呢
			_, err = new(VulTask).Create(result.Name, taskTestCase.TaskId, result.TaskID)
			if err != nil {
				panic(err)
			}
			return result.TaskID, nil
		} else {
			return "", err
		}
	}
	return "", nil
}

// 获取下一个待执行的合规测试用例
func (this *TaskTestCase) GetNextHgExecTestCase(taskUuid string) (qmap.QM, bool) {
	{
		// 首先判断有没有正在测试的合规测试用例，如果有，则把该测试用例推送给终端
		params := qmap.QM{
			"e_task_uuid":     taskUuid,
			"e_test_tool":     common.TOOL_HG_ANDROID_SCANNER,
			"e_action_status": common.CASE_STATUS_TESTING,
		}
		if has, testCase := sys_service.NewSessionWithCond(params).GetOne(this); has {
			return *testCase, true
		}
	}

	{
		// 优先查询处于队列中的合规测试用例
		params := qmap.QM{
			"e_task_uuid":     taskUuid,
			"e_test_tool":     common.TOOL_HG_ANDROID_SCANNER,
			"e_action_status": common.CASE_STATUS_QUEUING,
		}
		if has, testCase := sys_service.NewSessionWithCond(params).GetOne(this); has {
			return *testCase, true
		}
	}
	{
		// 优先查询处于待测试的合规测试用例
		params := qmap.QM{
			"e_task_uuid":     taskUuid,
			"e_test_tool":     common.TOOL_HG_ANDROID_SCANNER,
			"e_action_status": common.CASE_STATUS_READY,
		}
		if has, testCase := sys_service.NewSessionWithCond(params).GetOne(this); has {
			return *testCase, true
		} else {
			return nil, false
		}
	}
}

func (this *TaskTestCase) UpdateCase(taskUuid, caseUuid string, info qmap.QM) (*TaskTestCase, error) {
	has, err := sys_service.NewSession().Session.Where("task_uuid = ? and case_uuid = ?", taskUuid, caseUuid).Get(this)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("Item not found")
	}
	if val, has := info.TryInt("action_status"); has {
		this.ActionStatus = val
	}
	if val, has := info.TryInt("test_result_status"); has {
		this.TestResultStatus = val
		if val == common.CASE_TEST_STATUS_UNPASS || val == common.CASE_TEST_STATUS_PASS {
			this.CompleteTime = int(time.Now().Unix())
		}
	}
	if val, has := info.TryString("case_result"); has {
		this.CaseResult = val
	}
	if val, has := info.TryMap("case_result_file"); has {
		if this.CaseResultFile != "" {
			if temp, err := custom_util.StringToSlice(this.CaseResultFile); err == nil {
				temp = append(temp, val)
				this.CaseResultFile = custom_util.SliceToString(temp)
			} else {
				this.CaseResultFile = custom_util.SliceToString([]interface{}{val})
			}
		} else {
			this.CaseResultFile = custom_util.SliceToString([]interface{}{val})
		}
	}
	if val, has := info.TryString("test_procedure"); has {
		this.TestProcedure = val
	}
	_, err = sys_service.NewOrm().Table(this).ID(this.Id).AllCols().Update(this)
	return this, err
}

func (this *TaskTestCase) UpdateCaseById(id int, info qmap.QM) (*TaskTestCase, error) {
	has, err := sys_service.NewSession().Session.Where("id = ?", id).Get(this)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("Item not found")
	}
	if val, has := info.TryInt("action_status"); has {
		this.ActionStatus = val
		if this.ActionStatus == common.CASE_STATUS_QUEUING {
			// 如果测试状态重置为队列中，则重置测试用例的测试结果
			this.TestResultStatus = 0
		} else if this.ActionStatus == common.CASE_STATUS_CLOSED || this.ActionStatus == common.CASE_STATUS_CANCELED {
			// 如果测试用例为 "已忽略" 或者 "已取消"状态，则重置所有测试结果
			this.TestResultStatus = 0  // 重置测试结果状态
			this.TestProcedure = ""    // 重置测试过程
			this.CaseResultFile = "[]" // 重置测试结果文件
			this.CaseResult = ""       // 重置测试结果
			this.CompleteTime = 0      // 重置测试完成时间
			if err := new(Vulnerability).RemoveCaseVul(this.TaskId, id); err != nil {
				panic(err)
			}
		} else if this.ActionStatus == common.CASE_STATUS_TESTING {
			this.TestResultStatus = 0 // 重置测试结果状态
			this.CompleteTime = 0     // 重置测试完成时间
		}
	}
	if val, has := info.TryInt("test_result_status"); has {
		this.TestResultStatus = val
		if val == common.CASE_TEST_STATUS_UNPASS || val == common.CASE_TEST_STATUS_PASS {
			this.CompleteTime = int(time.Now().Unix())
		}
	}
	if val, has := info.TryString("case_result_file"); has {
		this.CaseResultFile = val
	}
	if val, has := info.TryString("test_procedure"); has {
		this.TestProcedure = val
	}
	if val, has := info.TryString("test_attachment"); has {
		this.TestAttachment = val
	}
	_, err = sys_service.NewOrm().Table(this).ID(this.Id).AllCols().Update(this)
	return this, err
}

func (this *TaskTestCase) FindOne(taskUuid, caseUuid string) (*TaskTestCase, error) {
	has, err := mysql.GetSession().Where("task_uuid = ?", taskUuid).
		And("case_uuid = ?", caseUuid).Get(this)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("not found")
	}
	return this, err
}

// 判断用例是否已经完成，是否属于重新测试
func IsRestartTest(caseId int) bool {
	taskTestCase := new(TaskTestCase)
	if has, _ := mysql.GetSession().Where("id = ?", caseId).Get(taskTestCase); has {
		if taskTestCase.ActionStatus == common.CASE_STATUS_COMPLETED {
			return true
		}
	}
	return false
}

// 重新测试的用例，添加到任务中
func AddTaskWithCase(taskId, caseId int) {
	// 1.判断这条任务是否已经完成，
}

func (this *TaskTestCase) SetCaseStatus(taskId int) error {
	cases := make([]TaskTestCase, 0)
	err := sys_service.NewSession().Session.Where("task_id = ?", taskId).Find(&cases)
	if err != nil {
		return err
	}
	for _, taskTestCase := range cases {
		if taskTestCase.ActionStatus == common.CASE_STATUS_COMPLETED {
			continue
		}
		taskTestCase.ActionStatus = common.CASE_STATUS_FAIL
		_, err := sys_service.NewOrm().Table(this).ID(taskTestCase.Id).Cols("action_status").Update(taskTestCase)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *TaskTestCase) GetCaseParams(taskId, caseUuid string) (qmap.QM, error) {
	testCase, err := this.FindOne(taskId, caseUuid)
	if err != nil {
		return nil, err
	}
	if taskParam, err := qmap.NewWithString(testCase.TaskParam); err == nil {
		return taskParam, nil
	} else {
		return nil, err
	}
}

// 当固件任务执行失败，将固件扫描任务关联的所有测试用例状态置为"扫描失败"
func (this *TaskTestCase) SetFirmwareCaseStatusToFailure(firmwareScannerId int) error {
	cases := make([]TaskTestCase, 0)
	err := sys_service.NewSession().Session.Where("tool_task_id = ?", firmwareScannerId).In("action_status = ?", []interface{}{common.CASE_STATUS_READY, common.CASE_STATUS_QUEUING}).Find(&cases)
	if err != nil {
		return err
	}
	for _, taskTestCase := range cases {
		taskTestCase.ActionStatus = common.CASE_STATUS_FAIL
		_, err := sys_service.NewOrm().Table(this).ID(taskTestCase.Id).Cols("action_status").Update(taskTestCase)
		if err != nil {
			return err
		}
	}
	return nil
}

// 通过任务的 uuid 获取测试脚本
func GetTestScriptByTaskUuid(taskUuid string) (testScriptSlice []string, err error) {
	testScriptSlice = make([]string, 0)
	err = sys_service.NewOrm().Table("task_test_case").Alias("ttc").
		Cols("ktc.test_script").
		Join("INNER", []string{"knowledge_test_case", "ktc"}, "ktc.id = ttc.test_case_id").
		Where("ttc.task_uuid = ?", taskUuid).
		And("ttc.test_tool = ?", common.TOOL_HG_ANDROID_SCANNER).
		Find(&testScriptSlice)
	if err != nil {
		clog.Error("GetTestScriptByTaskUuid Find Err", zap.Any("error", err))
		return testScriptSlice, err
	}
	clog.Info("sliceTestScript Info", zap.Any("Info", testScriptSlice))
	return testScriptSlice, nil
}

// 通过任务的 uuid 查询测试用例
func GetTestCaseIdByTaskUuid(taskUuid string) (testCaseIdSlice []int, err error) {
	testCaseIdSlice = make([]int, 0)
	err = sys_service.NewOrm().Table("task_test_case").
		Cols("test_case_id").
		Where("task_uuid = ?", taskUuid).
		And("test_tool = ?", common.TOOL_HG_ANDROID_SCANNER).
		Find(&testCaseIdSlice)
	if err != nil {
		clog.Error("GetTestCaseIdByTaskUuid Find Err", zap.Any("error", err))
		return testCaseIdSlice, err
	}
	clog.Debug("testCaseIdSlice Info", zap.Any("info", testCaseIdSlice))
	return testCaseIdSlice, nil
}
