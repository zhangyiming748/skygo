package mongo_model

import (
	"errors"
	"time"

	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
)

type HgTestTask struct {
	ID              bson.ObjectId `bson:"_id,omitempty" json:"_id"`                   // 记录_id
	Name            string        `bson:"name" json:"name"`                           // 任务名称
	TaskId          string        `bson:"task_id" json:"task_id"`                     // 任务全局id,6位string
	Status          string        `bson:"status" json:"status"`                       // 任务状态
	LastConnectTime int           `bson:"last_connect_time" json:"last_connect_time"` // 最近连接日期
	TestCase        []TestCase    `bson:"test_case" json:"test_case"`                 // 测试用例
	HgClientInfo    *HgClientInfo `bson:"hg_client_info" json:"hg_client_info"`       // 终端信息（内嵌文档） -- 设备上传
	StatusFlow      *StatusFlow   `bson:"status_flow" json:"status_flow"`             // 状态流程
	UserId          int           `bson:"user_id" json:"user_id"`                     // 用户ID
	TemplateFile    *HgFile       `bson:"template_file" json:"template_file"`         // 文件
	TemplateId      string        `bson:"template_id" json:"template_id"`             // 模板ID
}

const HgTestTaskStatusCreate = "create"                   // 任务状态 - 创建
const HgTestTaskStatusClientInfo = "client_info"          // 任务状态 - 获取信息
const HgTestTaskStatusChooseTestCase = "choose_test_case" // 任务状态 - 适配用例
const HgTestTaskStatusAutoTest = "auto_test"              // 任务状态 - 测试
const HgTestTaskStatusComplete = "complete"               // 任务状态 - 完成

// 作废
// const HgTestTaskStatusDefault = 1 // 任务状态 - 未开始
// const HgTestTaskStatusRunning = 2 // 任务状态 - 进行中
// const HgTestTaskStatusWarning = 3 // 任务状态 - 警告
// const HgTestTaskStatusFailed = 4  // 任务状态 - 失败
// const HgTestTaskStatusSuccess = 5 // 任务状态 - 成功

const HgTestTaskConnectStatusNever = 1 // 连接状态 - 未连接
const HgTestTaskConnectStatusYes = 2   // 连接状态 - 已连接
const HgTestTaskConnectStatusNo = 3    // 连接状态 - 连接断开

const HgTestTaskFlowNo = 1  // 状态流程节点状态 - 未完成
const HgTestTaskFlowYes = 2 // 状态流程节点状态 - 完成

const HgTestTaskCaseStatusPrepare = 1 // 任务中测试案例状态 - 待测试
const HgTestTaskCaseStatusRunning = 2 // 任务中测试案例状态 - 测试中
const HgTestTaskCaseStatusAnalyse = 3 // 任务中测试案例状态 - 分析中
const HgTestTaskCaseStatusPass = 4    // 任务中测试案例状态 - 通过
const HgTestTaskCaseStatusUnPass = 5  // 任务中测试案例状态 - 未通过
const HgTestTaskCaseStatusFail = 6    // 任务中测试案例状态 - 失败

const HgTestTaskCaseBlockStatusSuccess = 1 // 任务中测试案例Block的状态 - 失败
const HgTestTaskCaseBlockStatusFail = 2    // 任务中测试案例Block的状态 - 成功
const HgTestTaskCaseBlockStatusRunning = 3 // 任务中测试案例Block的状态 - 运行中

var HgTestTaskCaseStatusList = map[int]string{
	1: "待测试",
	2: "测试中",
	3: "分析中",
	4: "通过",
	5: "未通过",
	6: "失败",
}

type HgClientInfo struct {
	Cpu       string `bson:"cpu" json:"cpu"`               // cpu位数
	OsType    string `bson:"os_type" json:"os_type"`       // 操作系统类型
	OsVersion string `bson:"os_version" json:"os_version"` // 操作系统版本
}

type StatusFlow struct {
	CreateStage         int `bson:"create_stage" json:"create_stage"`                     // 创建阶段, 1进行中 2 完成
	CreateTime          int `bson:"create_time" json:"create_time"`                       // 创建时间
	ClientInfoStage     int `bson:"client_info_stage" json:"client_info_stage"`           // 获取设备信息阶段
	ClientInfoTime      int `bson:"client_info_time" json:"client_info_time"`             // 获取设备信息时间
	ChooseTestCaseStage int `bson:"choose_test_case_stage" json:"choose_test_case_stage"` // 选择测试用例阶段
	ChooseTestCaseTime  int `bson:"choose_test_case_time" json:"choose_test_case_time"`   // 选择测试用例时间
	TestingStage        int `bson:"testing_stage" json:"auto_test_stage"`                 // 测试阶段
	TestingTime         int `bson:"testing_time" json:"auto_test_time"`                   // 测试时间
	CompleteStage       int `bson:"complete_stage" json:"complete_stage"`                 // 任务完成阶段
	CompleteTime        int `bson:"complete_time" json:"complete_time"`                   // 任务完成时间
}

type TestCase struct {
	AutoTestLevel     string          `bson:"auto_test_level" json:"auto_test_level"`           // 类型， "自动化" "人工"
	TestCaseId        string          `bson:"test_case_id" json:"test_case_id"`                 // 检测用例_id，对应集合evaluate_test_case集合， 也代表 检测用例编码
	Name              string          `bson:"name" json:"name"`                                 // 检测用例名称
	TestProcedure     string          `bson:"test_procedure" json:"test_procedure"`             // 测试步骤
	Status            int             `bson:"status" json:"status"`                             // 检测用例状态  1待测试 2通过 3失败
	BlockList         []TestCaseBlock `bson:"block_list" json:"block_list"`                     // 测试用例block列表
	ApiResultDesc     string          `bson:"api_result_desc" json:"api_result_desc"`           // 接口返回的扫描结果内容
	ManMadeResultDesc string          `bson:"man_made_result_desc" json:"man_made_result_desc"` // 人工录入扫描结果内容
	FileList          []HgFile        `bson:"file_list" json:"file_list"`                       // 文件列表, (测试时block结果中会带上，截图时会带上)
}

type TestCaseBlock struct {
	CaseType  string `bson:"case_type" json:"case_type"`   // 测试用例类型 jar/so/apk
	BlockName string `bson:"block_name" json:"block_name"` // block英文名， 比如 runBlock
	Name      string `bson:"name" json:"name"`             // block英文名， 比如 runBlock
	TimeOut   int    `bson:"time_out" json:"time_out"`     // block执行超时（单位:秒）
}

func (h *HgTestTask) Create(rawInfo qmap.QM) (*HgTestTask, error) {
	if err := mongo.NewMgoSession(common.MC_EVALUATE_VUL_DEVICE_INFO).Insert(rawInfo); err == nil {
		return h, nil
	} else {
		return nil, err
	}
}

func (h *HgTestTask) Update(taskId string, rawInfo qmap.QM) (*HgTestTask, error) {
	params := qmap.QM{
		"e_task_id": taskId,
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.McHgTestTask, params)

	taskModel := HgTestTask{}

	if err := mongoClient.One(&taskModel); err == nil {
		if status, has := rawInfo.TryString("status"); has {
			taskModel.Status = status
			switch status {
			case HgTestTaskStatusChooseTestCase:
				taskModel.StatusFlow.ChooseTestCaseTime = int(time.Now().UnixNano() / 1e6)
				taskModel.StatusFlow.ChooseTestCaseStage = HgTestTaskFlowYes
			case HgTestTaskStatusComplete:
				taskModel.StatusFlow.TestingStage = HgTestTaskFlowYes
				taskModel.StatusFlow.TestingTime = int(time.Now().UnixNano() / 1e6)

				taskModel.StatusFlow.CompleteStage = HgTestTaskFlowYes
				taskModel.StatusFlow.CompleteTime = int(time.Now().UnixNano() / 1e6)

				// 把所有进行中的测试案例值为失败
				statusList := []int{
					HgTestTaskCaseStatusPrepare,
					HgTestTaskCaseStatusRunning,
					HgTestTaskCaseStatusAnalyse,
					// HgTestTaskCaseStatusPass,
					// HgTestTaskCaseStatusUnPass,
					// HgTestTaskCaseStatusFail,
				}
				for k, tc := range taskModel.TestCase {
					if custom_util.InIntSlice(tc.Status, statusList) {
						taskModel.TestCase[k].Status = HgTestTaskCaseStatusFail
					}
				}
			}
		}

		// 最后更新人
		if err := mongoClient.Update(bson.M{"task_id": taskId}, taskModel); err != nil {
			return nil, err
		} else {
			return &taskModel, nil
		}
	}
	return nil, errors.New("未查到任务记录")
}

func (h *HgTestTask) BulkDelete(rawIds []string) (*qmap.QM, error) {
	// 删除 测试项
	effectNum := 0
	ids := []string{}
	for _, id := range rawIds {
		ids = append(ids, id)
	}
	if len(ids) > 0 {
		match := bson.M{
			"task_id": bson.M{"$in": ids},
		}
		if changeInfo, err := mongo.NewMgoSession(common.MC_EVALUATE_VUL_DEVICE_INFO).RemoveAll(match); err == nil {
			effectNum = changeInfo.Removed
			// 根据item_id删除 测试项里的漏洞
			new(EvaluateVulnerability).BulkDeleteByItemIds(rawIds)
		} else {
			return nil, err
		}
	}
	return &qmap.QM{"number": effectNum}, nil
}

func (h *HgTestTask) FindByTaskId(taskId string) (*HgTestTask, error) {
	session := mongo.NewMgoSession(common.McHgTestTask)
	param := qmap.QM{
		"e_task_id": taskId,
	}

	// 1、检测任务必须存在
	taskModel := HgTestTask{}
	err := session.AddCondition(param).One(&taskModel)
	return &taskModel, err
}

// 获取任务中的单个测试案例
func (h *HgTestTask) FindTaskTestCase(taskId string, testCaseId string) (*TestCase, error) {
	taskPtr, err := h.FindByTaskId(taskId)
	if err != nil {
		return nil, err
	}
	for _, t := range taskPtr.TestCase {
		if t.TestCaseId == testCaseId {
			return &t, nil
		}
	}
	return nil, errors.New("记录不存在")
}

// 更新任务中的测试用例列表，包含“自动” 和 “交互”测试用例
func (h *HgTestTask) UpdateTestCase(taskId string, tc []TestCase, template *HgTestTemplate) error {
	update := bson.M{
		"$set": bson.M{
			"test_case":     tc,
			"template_file": template.File,
			"template_id":   template.ID,
		},
	}

	return mongo.NewMgoSession(common.McHgTestTask).Update(bson.M{"task_id": taskId}, update)
}

// 更新任务中的设备最近在线时间
func (h *HgTestTask) UpdateConnected(taskId string) error {
	update := bson.M{
		"$set": bson.M{
			"last_connect_time": time.Now().UnixNano() / 1e6,
		},
	}
	return mongo.NewMgoSession(common.McHgTestTask).Update(bson.M{"task_id": taskId}, update)
}

// 更新任务中设备信息， 包含状态流的更新
func (h *HgTestTask) UpdateClientInfo(taskId string, clientInfo *HgClientInfo) error {
	update := bson.M{
		"$set": bson.M{
			"hg_client_info":                clientInfo,                  // 设备信息
			"status_flow.client_info_stage": HgTestTaskFlowYes,           // 状态流中的具体字段
			"status_flow.client_info_time":  time.Now().UnixNano() / 1e6, // 状态流中的具体字段
			"status":                        HgTestTaskStatusClientInfo,
		},
	}
	return mongo.NewMgoSession(common.McHgTestTask).Update(bson.M{"task_id": taskId}, update)
}

func (h *HgTestTask) UpdateStatusFlow(taskId string, flag HgTestTaskStatusFlowFlag) error {
	stageKey := flag.PrintFlag()
	update := bson.M{
		"$set": bson.M{
			"status_flow." + stageKey + "_stage": HgTestTaskFlowYes,
			"status_flow." + stageKey + "_time":  time.Now(),
		},
	}
	return mongo.NewMgoSession(common.McHgTestTask).Update(bson.M{"task_id": taskId}, update)
}

// ---------------------------------------------
// 任务状态流转，每个状态的识别符号
type HgTestTaskStatusFlowFlag interface {
	PrintFlag() (stageKey string)
}

type HgFlagCreate struct{}

func (h *HgFlagCreate) PrintFlag() string {
	return "create"
}
func (h *HgFlagCreate) PrintChinese() string {
	return "创建成功"
}

// 更新设备信息时就顺带更新了状态流，这个不使用
type HgFlagClientInfo struct{}

func (h *HgFlagClientInfo) PrintFlag() string {
	return "client_info"
}
func (h *HgFlagClientInfo) PrintChinese() string {
	return "获取信息"
}

type HgFlagChooseTestCase struct{}

func (h *HgFlagChooseTestCase) PrintFlag() string {
	return "choose_test_case"
}
func (h *HgFlagChooseTestCase) PrintChinese() string {
	return "匹配用例"
}

type HgFlagTesting struct{}

func (h *HgFlagTesting) PrintFlag() string {
	return "testing"
}
func (h *HgFlagTesting) PrintChinese() string {
	return "测试阶段"
}

type HgFlagComplete struct{}

func (h *HgFlagComplete) PrintFlag() string {
	return "complete"
}
func (h *HgFlagComplete) PrintChinese() string {
	return "完成"
}
