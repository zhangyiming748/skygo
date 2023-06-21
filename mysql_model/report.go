package mysql_model

import (
	"encoding/json"
	"skygo_detection/lib/common_lib/mysql"
)

type Report struct {
	Demands []DemandReport `json:"demands"`
}

type DemandReport struct {
	Id      int          `json:"id"`
	Title   string       `json:"title"`
	Content string       `json:"content"`
	Cases   []CaseReport `json:"cases"`
}

type CaseReport struct {
	Id                   int      `json:"id"`
	TestProcedure        string   `json:"test_procedure"`         // 执行过程中的记录，程序会改
	TestProcedureComment string   `json:"test_procedure_comment"` // 拷贝过来的测试步骤
	TestResultStatus     int      `json:"test_result_status"`     // task_case中字段
	ModuleId             string   `json:"module_id"`              // test_case中字段
	ModuleTypeName       string   `json:"module_type_name"`
	ModuleName           string   `json:"module_name"`
	FixSuggest           []string `json:"fix_suggest"`      // vulnerablity中字段
	CaseResultFile       string   `json:"case_result_file"` // task_case中字段
	TestCaseName         string   `json:"test_case_name"`   // 测试用例名称
	TestStandard         string   `json:"test_standard"`    // 验证标准
}

type ReportTreeNode struct {
	Id       int               `json:"id"`      // 主键id
	Code     string            `json:"code"`    // 章节编号
	Title    string            `json:"title"`   // 章节标题
	Content  string            `json:"content"` // 章节下边的介绍
	Cases    []*CaseReport     `json:"cases"`
	Children []*ReportTreeNode `json:"children"` // 子节点
}

func (this *ReportTreeNode) GetResultDetail(taskId int) []*ReportTreeNode {
	// 取任务
	task := new(Task)
	mysql.GetSession().Where("id = ?", taskId).Get(task)

	// 取senario
	scenario := new(KnowledgeScenario)
	mysql.GetSession().Where("id = ?", task.ScenarioId).Get(scenario)

	// 取一级条目
	parentNodes := FetchParent(scenario.DemandId)

	// 取所有的子条目
	for k, node := range parentNodes {
		parentNodes[k].FetchChild(node)
		parentNodes[k].RangeAllChild(scenario.Id, scenario.DemandId, taskId)
	}

	// 取叶子条目的测试用例
	// for k,_ := range parentNodes {
	// parentNodes[k].RangeAllChild(scenario.Id ,scenario.DemandId,taskId)
	// }
	return parentNodes
}

func FetchParent(demanId int) []*ReportTreeNode {
	var parent = make([]*ReportTreeNode, 0)
	var kdc = make([]KnowledgeDemandChapter, 0)
	_ = mysql.GetSession().Where("knowledge_demand_id=?", demanId).Where("parent_id=?", 0).Find(&kdc)
	for _, v := range kdc {
		var tmp = new(ReportTreeNode)
		tmp.Id = v.Id
		tmp.Code = v.Code
		tmp.Title = v.Title
		tmp.Content = v.Content
		tmp.Cases = make([]*CaseReport, 0)
		tmp.Children = make([]*ReportTreeNode, 0)
		parent = append(parent, tmp)
	}
	return parent
}

func GetCases(demand_chapter_id int, senario_id int, demand_id int, task_id int) []*CaseReport {
	type FS struct {
		Caution      string `json:"caution"`
		RepairCost   string `json:"repair_cost"`
		RepairEffect string `json:"repair_effect"`
		Importance   string `json:"importance"`
	}
	ktcc := make([]KnowledgeTestCaseChapter, 0)
	_ = mysql.GetSession().Where("senario_id=?", senario_id).Where("demand_id=?", demand_id).Where("demand_chapter_id=?", demand_chapter_id).Find(&ktcc)
	var cases = make([]*CaseReport, 0)
	for _, t := range ktcc {
		testStandard, testProcedure := "", ""
		testCaseId := t.TestCaseId
		// 取测试用例
		ktc := new(KnowledgeTestCase)
		if has, _ := mysql.GetSession().Where("id=?", testCaseId).Get(ktc); has {
			testStandard = ktc.TestStandard
			testProcedure = ktc.TestProcedure
		}

		// 取任务里的测试用例
		ttc := new(TaskTestCase)
		_, _ = mysql.GetSession().Where("task_id=?", task_id).Where("test_case_id=?", testCaseId).Get(ttc)
		var tmpCase = &CaseReport{
			Id:                   ttc.Id,
			TestProcedure:        ttc.TestProcedure, // 程序运行中产生的
			TestProcedureComment: testProcedure,     // 拷贝过来的
			TestStandard:         testStandard,
			TestResultStatus:     ttc.TestResultStatus,
			ModuleId:             ttc.ModuleId,
			CaseResultFile:       ttc.CaseResultFile,
			TestCaseName:         ttc.TestCaseName,
			FixSuggest:           []string{},
		}

		// Vulnerability
		fix_suggest := []string{}
		vb := make([]Vulnerability, 0)
		if err := mysql.GetSession().Where("task_id=?", task_id).Where("task_case_id=?", ttc.Id).Find(&vb); err == nil {
			if len(vb) > 0 {
				for _, v := range vb {
					if v.FixSuggest != "" {
						fs := []FS{}
						err := json.Unmarshal([]byte(v.FixSuggest), &fs)
						if err == nil {
							for _, v := range fs {
								fix_suggest = append(fix_suggest, v.Caution)
							}
						}
					}
				}
				tmpCase.FixSuggest = fix_suggest
			}
		}
		cases = append(cases, tmpCase)
	}
	return cases
}

// 节点递归方式获取子节点
func (this *ReportTreeNode) FetchChild(node *ReportTreeNode) {
	var children = make([]*ReportTreeNode, 0)
	var childChpter = make([]KnowledgeDemandChapter, 0)
	_ = mysql.GetSession().Where("parent_id = ?", node.Id).Find(&childChpter)
	for _, v := range childChpter {
		child := ReportTreeNode{
			Id:      v.Id,
			Code:    v.Code,
			Title:   v.Title,
			Content: v.Content,
			Cases:   make([]*CaseReport, 0),
		}
		children = append(children, &child)
	}
	this.Children = children
	if len(this.Children) > 0 {
		for k, v := range this.Children {
			this.Children[k].FetchChild(v)
		}
	}
	return
}

func (this *ReportTreeNode) RangeAllChild(senario_id int, demand_id int, task_id int) {
	if len(this.Children) > 0 {
		for _, v := range this.Children {
			v.RangeAllChild(senario_id, demand_id, task_id)
		}
	} else {
		this.Cases = GetCases(this.Id, senario_id, demand_id, task_id)
	}
}
