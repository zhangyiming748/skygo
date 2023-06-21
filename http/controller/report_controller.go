package controller

import (
	"errors"
	"fmt"
	"sort"
	"time"

	"skygo_detection/guardian/src/net/qmap"

	"github.com/gin-gonic/gin"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/log"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mongo_model"
	"skygo_detection/mysql_model"
)

type ReportController struct{}

func (this ReportController) GetResultPASSRate(ctx *gin.Context) {
	taskId := ctx.Query("task_id")
	s := mysql.GetSession()
	s.Where("task_id=?", taskId)
	s.Where("test_result_status=?", common.CASE_TEST_STATUS_PASS)
	pass, _ := s.FindAndCount(&[]mysql_model.TaskTestCase{})
	s = mysql.GetSession()
	s.Where("task_id=?", taskId)
	s.Where("test_result_status=?", common.CASE_TEST_STATUS_UNPASS)
	fail, _ := s.FindAndCount(&[]mysql_model.TaskTestCase{})
	result := qmap.QM{
		"pass": pass,
		"fail": fail,
	}
	response.RenderSuccess(ctx, result)
}

func (this ReportController) GetResultDistribution(ctx *gin.Context) {
	var taskTestCases = make([]mysql_model.TaskTestCase, 0)
	// 查询所有的任务测试用例
	taskId := ctx.Query("task_id")
	s := mysql.GetSession()
	s.Where("task_id=?", taskId)
	err := s.Find(&taskTestCases)
	if err != nil {
		response.RenderFailure(ctx, err)
	}
	failCalc := map[string]int{}
	passCalc := map[string]int{}
	moduleNameMap := map[string]string{}

	for _, taskTestCase := range taskTestCases {
		moduleName := ""
		if val, has := moduleNameMap[taskTestCase.ModuleId]; has {
			moduleName = val
		} else {
			if module, err := new(mongo_model.EvaluateModule).FindById(taskTestCase.ModuleId); err == nil {
				// moduleName = module.ModuleName + "/" + module.ModuleType
				moduleName = module.ModuleName
			} else {
				moduleName = "不存在这个测试组件，请检查"
			}
			moduleNameMap[taskTestCase.ModuleId] = moduleName
		}
		if taskTestCase.TestResultStatus == common.CASE_TEST_STATUS_PASS {
			if val, has := passCalc[taskTestCase.ModuleId]; has {
				passCalc[taskTestCase.ModuleId] = val + 1
			} else {
				passCalc[taskTestCase.ModuleId] = 1
				failCalc[taskTestCase.ModuleId] = 0
			}
		} else if taskTestCase.TestResultStatus == common.CASE_TEST_STATUS_UNPASS {
			if val, has := failCalc[taskTestCase.ModuleId]; has {
				failCalc[taskTestCase.ModuleId] = val + 1
			} else {
				passCalc[taskTestCase.ModuleId] = 0
				failCalc[taskTestCase.ModuleId] = 1
			}
		}
	}
	var result = make([]qmap.QM, 0)
	// 然后计算pass,fail个数
	for k, v := range moduleNameMap {
		if passCalc[k] > 0 || failCalc[k] > 0 {
			tmpResult := qmap.QM{
				"name": v,
				"pass": passCalc[k],
				"fail": failCalc[k],
			}
			result = append(result, tmpResult)
		}
	}
	response.RenderSuccess(ctx, result)
}

func (this ReportController) GetResultView(ctx *gin.Context) {
	// 查询所有的任务测试用例
	var taskTestCases = make([]mysql_model.TaskTestCase, 0)
	taskId := ctx.Query("task_id")
	s := mysql.GetSession()
	s.Where("task_id=?", taskId)
	s.Find(&taskTestCases)
	var casesResult = map[int]int{}
	for _, taskTestCase := range taskTestCases {
		casesResult[taskTestCase.TestCaseId] = taskTestCase.TestResultStatus
	}

	// 取任务
	task := new(mysql_model.Task)
	mysql.GetSession().Where("id = ?", taskId).Get(task)
	scenario_id := task.ScenarioId

	// 章节
	var chapterRsult = make(map[int]int, 0)
	for _, taskTestCase := range taskTestCases {
		var chapters = make([]mysql_model.KnowledgeTestCaseChapter, 0)
		s = mysql.GetSession()
		s.Where("test_case_id=?", taskTestCase.TestCaseId)
		s.Where("senario_id=?", scenario_id)
		s.Where("demand_id=?", taskTestCase.DemandId)
		s.Find(&chapters)
		// 单条章节数据处理
		for _, chapter := range chapters {
			chapterId := chapter.DemandChapterId
			caseId := chapter.TestCaseId
			// status 数据库的状态
			if status, ok := casesResult[caseId]; ok {
				// tmpStatus 已存的状态
				if tmpStatus, ok := chapterRsult[chapterId]; ok {
					// 如果已存的状态不是失败，那就采用数据的状态。如果章节已经失败，就不管
					if tmpStatus != common.CASE_TEST_STATUS_UNPASS {
						if tmpStatus == 1 {
							chapterRsult[chapterId] = 1
						}
					}
				} else {
					chapterRsult[chapterId] = status
				}
			}
		}
	}
	// 处理章节id，章节的数据结构是 章节id:1
	var result = make([]qmap.QM, 0)
	for k, v := range chapterRsult {
		tmp := qmap.QM{}
		var tt = new(mysql_model.KnowledgeDemandChapter)
		s = mysql.GetSession()
		s.Where("id=?", k)
		if has, _ := s.Get(tt); has {
			tmp["name"] = tt.Title
			tmp["code"] = tt.Code
			tmp["result"] = v
			result = append(result, tmp)
		}
	}
	if len(result) > 0 {
		sort.Slice(result, func(i, j int) bool {
			iv := result[i]["code"].(string)
			jv := result[j]["code"].(string)
			return iv < jv
		})
	}
	response.RenderSuccess(ctx, result)

}

// func (this ReportController) GetResultDetail1(ctx *gin.Context) {
// 	// 查询所有的任务测试用例
// 	var taskTestCases = make([]mysql_model.TaskTestCase, 0)
// 	taskId := ctx.Query("task_id")
// 	s := mysql.GetSession()
// 	s.Where("task_id=?", taskId)
// 	s.Find(&taskTestCases)
// 	var demandList = map[int][]mysql_model.CaseReport{}
// 	for _, taskTestCase := range taskTestCases {
// 		// 从reportResult中获取结果，如果存在
// 		var testCase = new(mysql_model.KnowledgeTestCase)
// 		testProcedure := "不存在"
// 		s := mysql.GetSession()
// 		s.Where("id=?", taskTestCase.TestCaseId)
// 		if has, _ := s.Get(testCase); has {
// 			testProcedure = testCase.TestProcedure
// 		}
// 		caseReport := mysql_model.CaseReport{
// 			Id:               taskTestCase.Id,
// 			TestProcedure:    testProcedure,
// 			TestResultStatus: taskTestCase.TestResultStatus,
// 			ModuleId:         taskTestCase.ModuleId,
// 			FixSuggest:       "未提供",
// 			CaseResultFile:   taskTestCase.CaseResultFile,
// 		}
// 		// 先获取已存在的case，然后增加
// 		if tmpCases, ok := demandList[taskTestCase.DemandId]; ok {
// 			tmpCases = append(tmpCases, caseReport)
// 			demandList[taskTestCase.DemandId] = tmpCases
// 		} else {
// 			tmpCases = []mysql_model.CaseReport{caseReport}
// 			demandList[taskTestCase.DemandId] = tmpCases
// 		}
// 	}
// 	var reportResult = mysql_model.Report{}
// 	// demandId列表
// 	for demandId, caseReports := range demandList {
// 		var demand = new(mysql_model.KnowledgeDemand)
// 		s := mysql.GetSession()
// 		s.Where("id=?", demandId)
// 		if has, _ := s.Get(demand); has {
// 			var demandReport = mysql_model.DemandReport{}
// 			demandReport.Id = demand.Id
// 			demandReport.Title = demand.Name
// 			demandReport.Content = demand.Detail
// 			demandReport.Cases = caseReports
// 			tmpDemands := reportResult.Demands
// 			tmpDemands = append(tmpDemands, demandReport)
// 			reportResult.Demands = tmpDemands
// 		} else {
// 			fmt.Println("has not demand id:", demandId)
// 		}
// 	}
// 	if result, err := json.Marshal(reportResult); err != nil {
// 		response.RenderFailure(ctx, err)
// 	} else {
// 		var tmp = map[string]interface{}{}
// 		json.Unmarshal(result, &tmp)
// 		response.RenderSuccess(ctx, tmp)
// 	}
// }

func (this ReportController) GetResultDetail(ctx *gin.Context) {
	taskId := request.QueryInt(ctx, "task_id")
	data := new(mysql_model.ReportTreeNode).GetResultDetail(taskId)

	// 查询 module_name module_type_name
	getModuleName(data)
	response.RenderSuccess(ctx, data)
}

func getModuleName(data []*mysql_model.ReportTreeNode) {
	for k, d := range data {
		if len(d.Children) > 0 {
			getModuleName(data[k].Children)
		}
		for ck, cv := range data[k].Cases {
			evaluateModule, err := new(mongo_model.EvaluateModule).FindById(cv.ModuleId)
			if err == nil {
				data[k].Cases[ck].ModuleTypeName = evaluateModule.ModuleType
				data[k].Cases[ck].ModuleName = evaluateModule.ModuleName
			}
		}
	}
}

func (this ReportController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	s := mysql.GetSession()
	// s.Where("status=?", common.REPORT_SUCCESS)
	// 查询组键
	content := request.QueryString(ctx, "content")
	if content != "" {
		s.Where("name like ? or report_name like ?", "%"+content+"%", "%"+content+"%")
	}
	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	widget.AddSorter(*(orm.NewSorter("id", 1)))
	// widget.SetTransformer(&transformer.ReportDetailTransformer{})
	all := widget.PaginatorFind(s, &[]mysql_model.ReportTask{})
	response.RenderSuccess(ctx, all)
}
func (this ReportController) GetRange(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	start := params.String("start_time") //任务开始时间
	fmt.Printf("获取到的任务开始时间是%v,类型是%T\n", start, start)
	end := params.String("end_time") //任务结束时间
	fmt.Printf("获取到的任务结束时间是%v,类型是%T\n", end, end)

	starttime, _ := time.Parse("2006-01-02 15:04:05", start)
	endtime, _ := time.Parse("2006-01-02 15:04:05", end)
	fmt.Printf("转换成时间类型之后的\nstart:%v\nend:%v\n", starttime, endtime)
	sub := int(endtime.Sub(starttime).Hours())
	fmt.Printf("时间差:%d\n", sub)
	if sub < 0 {
		response.RenderFailure(ctx, errors.New("不合法的时间范围"))
	}
	s := mysql.GetSession()
	s.Where("status=?", common.REPORT_SUCCESS)
	s.Where("create_time >= ? and create_time <= ?", start, end)
	widget := orm.PWidget{}
	all := widget.PaginatorFind(s, &[]mysql_model.ReportTask{})
	response.RenderSuccess(ctx, all)
}
func (this ReportController) IfGenerate(ctx *gin.Context) {
	id := request.ParamInt(ctx, "id")
	s := mysql.GetSession()
	s.Where("task_id=?", id)

	widget := orm.PWidget{}
	var result = make(map[string]interface{}, 0)
	rs, err := widget.One(s, &mysql_model.ReportTask{})
	if err != nil {
		log.GetHttpLogLogger().Error(fmt.Sprintf("report_task is null,err:%v", err))
	}
	result["file_id"] = rs["file_id"]         // 报告文件id
	result["report_name"] = rs["report_name"] // 报告名称

	response.RenderSuccess(ctx, result)
}

func (this ReportController) Create(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	rt := new(mysql_model.ReportTask)
	taskId := req.MustInt("task_id")
	task := new(mysql_model.Task)
	if _, err := mysql.GetSession().Where("id = ?", taskId).Get(task); err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	if task.Id < 1 {
		response.RenderFailure(ctx, errors.New("任务不存在"))
		return
	}
	rt.TaskId = taskId
	rt.Name = task.Name
	rt.ReportType = req.DefaultInt("report_type", 4) // 默认正常的报告
	rt.CreateTime = time.Now().Format("2006-01-02 15:04:05")
	if _, err := rt.Create(); err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, nil)

}
