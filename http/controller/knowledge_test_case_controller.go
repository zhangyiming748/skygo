package controller

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/http/transformer"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/log"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mongo_model"
	"skygo_detection/mysql_model"
)

type KnowledgeTestCaseController struct{}

/**
 * apiType http
 * @api {get} /api/v1/knowledge_test_cases 查询knowledge测试用例
 * @apiVersion 0.1.0
 * @apiName GetAll
 * @apiGroup KnowledgeTestCase
 *
 * @apiDescription 根据id,查询knowledge测试用例
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *   "code": 0,
 *   "data": {
 *       "list": [
 *           {
 *               "auto_test_level": 2,
 *               "case_uuid": "CASE_31_1",
 *               "create_time": 0,
 *               "create_user_id": 0,
 *               "demand_chapter_id": [],
 *               "demand_chapter_name": [],
 *               "demand_id": 1,
 *                "demand_name": "车载信息交互系统信息安全技术要求",
 *                "id": 1,
 *                "input": "",
 *                "last_op_id": 0,
 *                "level": 0,
 *                "module_id": "",
 *                "module_name": "",
 *                "module_type_name": "",
 *                "name": "拨打电话安全试验方法",
 *                "objective": "在应用软件内调用拨打电话操作，检查应用软件调用执行拨打电话操作时，是否在用户明示同 意的情况下，才能执行拨打操作； ",
 *                "scenario_id": 1,
 *                "scenario_name": "车载合规测试",
 *                "tag": "123",
 *                "task_param": "",
 *                "test_case_level": 2,
 *                "test_method": 2,
 *                "test_param": "",
 *                "test_procedure": "调用安卓系统拨打电话的接口，检测Activity变化，截图",
 *                "test_result": "",
 *                "test_script": "617bb2b0e1382311f6245bfb",
 *                "test_scripts": [
 *                    {
 *                        "name": "CASE_31_1.zip",
 *                        "value": "617bb2b0e1382311f6245bfb"
 *                    }
 *                ],
 *                "test_sketch_map": "",
 *                "test_standard": "如果未弹出拨打电话授权窗口则不通过",
 *                "test_tool": "hg_scanner",
 *                "test_tool_name": "车机检测工具",
 *                "test_tools_id": ""
 *           }
 *       ],
 *       "pagination": {
 *           "current_page": 1,
 *           "per_page": 20,
 *           "total": 118,
 *           "total_pages": 6
 *       }
 *   },
 *   "msg": ""
 * }
 */
func (this KnowledgeTestCaseController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	s := mysql.GetSession()

	// module_id 单独处理
	widget := orm.PWidget{}
	if moduleIdStr := request.QueryString(ctx, "module_id"); moduleIdStr != "" {
		if moduleId, err := custom_util.StringToSlice(moduleIdStr); err == nil {
			s.In("module_id", moduleId)
		}
	}

	// demand_chapter_id 单独处理
	demand_chapter_id := request.QueryInt(ctx, "demand_chapter_id")
	if demand_chapter_id > 0 {
		kdtn := new(mysql_model.KnowLedgeDemandTreeNode)
		kdtn.Id = demand_chapter_id
		kdtn.FetchChild([]int{})

		ids := []int{}
		ids = mysql_model.GetAllChild(kdtn)
		ids = append(ids, demand_chapter_id)
		ktcc := new(mysql_model.KnowledgeTestCaseChapter)
		data := ktcc.GetByDemandChapterIds(ids)
		test_case_ids := []int{}
		if len(data) > 0 {
			for _, v := range data {
				test_case_ids = append(test_case_ids, v.TestCaseId)
			}
			s.In("id", test_case_ids)
		} else {
			s.Where("1!=1")
		}
	}

	widget.SetQueryStr(queryParams)
	widget.AddSorter(*(orm.NewSorter("id", 1)))
	widget.SetTransformer(&transformer.KnowledgeTestCaseTransformer{})
	all := widget.PaginatorFind(s, &[]mysql_model.KnowledgeTestCase{})
	response.RenderSuccess(ctx, all)
}

/**
 * apiType http
 * @api {get} /api/v1/knowledge_test_cases/:id 查询knowledge测试用例
 * @apiVersion 0.1.0
 * @apiName GetOne
 * @apiGroup KnowledgeTestCase
 *
 * @apiDescription 根据id,查询knowledge测试用例
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *   "code": 0,
 *   "data": {
 *        "auto_test_level": 3,
 *        "case_uuid": "CASE_51",
 *        "create_time": 0,
 *        "create_user_id": 0,
 *        "demand_chapter_id": [],
 *        "demand_chapter_name": [],
 *        "demand_id": 1,
 *        "demand_name": "车载信息交互系统信息安全技术要求",
 *        "id": 113,
 *        "input": "",
 *        "last_op_id": 0,
 *        "level": 0,
 *        "module_id": "",
 *        "module_name": "",
 *        "module_type_name": "",
 *        "name": "定位功能试验方法",
 *        "objective": "检查当应用软件调用定位功能时，车载信息交互系统是否在用户界面上有相应的状态提示。",
 *        "scenario_id": 1,
 *        "scenario_name": "车载合规测试",
 *        "tag": "",
 *        "task_param": "",
 *        "test_case_level": 2,
 *        "test_method": 2,
 *        "test_param": "",
 *        "test_procedure": "使用签名的应用调用定位功能，检查是否在用户界面上有相应的状态",
 *        "test_result": "",
 *        "test_script": "61820226e830c64451f07287",
 *        "test_scripts": [
 *            {
 *                 "name": "CASE_51.zip",
 *                 "value": "61820226e830c64451f07287"
 *            }
 *        ],
 *        "test_sketch_map": "",
 *        "test_sketch_maps":[
 *			{
 *				"name": "106f5.png",
 *				"value": "6189f6dae830c613130eb1c4"
 *			},
 *			{
 *				"name": "106f5.png",
 *				"value": "6189f6dfe830c613130eb1c6"
 *			}
 *			],
 *        "test_standard": "",
 *        "test_tool": "",
 *        "test_tool_name": "车机检测工具",
 *        "test_tools_id": ""
 *   },
 *   "msg": ""
 * }
 */
func (this KnowledgeTestCaseController) GetOne(ctx *gin.Context) {
	id := request.ParamString(ctx, "id")
	s := mysql.GetSession()
	s.Where("id=?", id)

	w := orm.PWidget{}
	w.SetTransformer(&transformer.KnowledgeTestCaseTransformer{})
	result, err := w.One(s, &mysql_model.KnowledgeTestCase{})

	if err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/knowledge_test_cases 创建knowledge测试用例
 * @apiVersion 0.1.0
 * @apiName Create
 * @apiGroup KnowledgeTestCase
 *
 * @apiDescription 创建一条knowledge测试用例
 *
 * @apiUse authHeader
 *
 * @apiParam {string} 			name  			测试用例名称(该字段必须传)
 * @apiParam {array} 			demand_chapter_id  			章节id组成的数组[3,4,5]
 * @apiParam {string} 			module_id  	测试组件/测试分类
 * @apiParam {int}       scenario_id    	安全检测场景
 * @apiParam {int}		demand_id			关联安全需求id，需求条目存另一个表
 * @apiParam {string} 		objective			测试目的
 * @apiParam {string} 		input			外部输入
 * @apiParam {string} 		test_procedure			测试步骤
 * @apiParam {string} 		test_standard			验证标准
 * @apiParam {int} 		level			测试难度
 * @apiParam {int} 		test_case_level			测试用例级别
 * @apiParam {int} 		test_method			测试方式，1黑盒 2灰盒 3白盒
 * @apiParam {int} 		auto_test_level			自动化测试程度 1人工 2半自动化 3自动化
 * @apiParam {string} 		test_tool			测试工具
 * @apiParam {string} 		test_tool_name			测试工具名称
 * @apiParam {string} 		test_tools_id			测试工具id
 * @apiParam {string} 		task_param			任务参数
 * @apiParam {int} 		create_user_id			创建用户id
 * @apiParam {int} 		last_op_id			最近操作用户id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *	"name":"测试用例名称(该字段必须传)",
 *  "demand_chapter_id":[3,4,5]
 *	"module_id":"12",
 *	"scenario_id":2,
 *	"demand_id":3,
 *	"objective":"abc",
 *	"input":"",
 *	"test_procedure":"",
 *	"test_standard":"",
 *	"level":1,
 *	"test_case_level":1,
 *	"test_method":1,
 *	"auto_test_level":1,
 *	"test_tool":"",
 *	"test_tool_name":"",
 *	"test_tools_id":"",
 *	"test_script":"6189f5a8e830c613130eb1bc|6189f5a8e830c613130eb1ba",
 *	"test_sketch_map":"6189f6dae830c613130eb1c4|6189f6dfe830c613130eb1c6",
 *	"task_param":"",
 *	"create_user_id":0,
 *	"last_op_id":0
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *	"code": 0,
 *	"data": {
 *		 "auto_test_level": 1,
 *		 "case_uuid": "",
 *		 "create_time": 1635930902,
 *		 "create_user_id": 0,
 *		 "demand_id": 3,
 *		 "id": 118,
 *		 "input": "",
 *		 "last_op_id": 0,
 *		 "level": 1,
 *		 "module_id": "12",
 *		 "name": "测试用例名称",
 *		 "objective": "abc",
 *		 "scenario_id": 2,
 *		 "tag": "",
 *		 "task_param": "",
 *		 "test_case_level": 1,
 *		 "test_method": 1,
 *		 "test_param": "[{"case_type":"jar","block_name":"captureScreen","case_id":"hg_CASE_6","ret_of_crash":"","test_level":1,"name":"截图","time_out":10,"upload_crash_log":false}]",
 *		 "test_procedure": "",
 *		 "test_result": "",
 *		 "test_script": "",
 *		 "test_sketch_map": "",
 *		 "test_standard": "",
 *		 "test_tool": "",
 *		 "test_tool_name": "",
 *		 "test_tools_id": ""
 *	},
 *	"msg": ""
 * }
 */
func (this KnowledgeTestCaseController) Create(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	demandChapterIds := req.SliceInt("demand_chapter_id")

	testCase := new(mysql_model.KnowledgeTestCase)
	testCase.Name = req.MustString("name")
	testCase.ModuleId = req.String("module_id")
	testCase.ScenarioId = req.DefaultInt("scenario_id", 0)
	testCase.DemandId = req.DefaultInt("demand_id", 0)
	testCase.Objective = req.DefaultString("objective", "")
	testCase.Input = req.DefaultString("input", "")
	testCase.TestProcedure = req.DefaultString("test_procedure", "")
	testCase.TestStandard = req.DefaultString("test_standard", "")
	testCase.Level = req.DefaultInt("level", 1)
	testCase.TestCaseLevel = req.DefaultInt("test_case_level", 1)
	testCase.TestMethod = req.DefaultInt("test_method", 1)

	testCase.TestTool = req.String("test_tool")
	testCase.TestToolName = req.String("test_tool_name")
	testCase.TestToolsId = req.String("test_tools_id")
	testCase.AutoTestLevel = req.DefaultInt("auto_test_level", 1)

	// if testCase.TestToolsId != "" {
	// 	params := qmap.QM{
	// 			"e__id": bson.ObjectIdHex(testCase.TestToolsId),
	// 		}
	// 	ormSession := mongo.NewMgoSessionWithCond(common.MC_TOOL, params)
	// 	data, _ := ormSession.GetOne()
	// 	if len(*data) > 0 {
	// 		rs := *data
	// 		testCase.TestToolName = rs["tool"].(string)
	// 	}
	// }
	testCase.TaskParam = req.DefaultString("task_param", "")
	testCase.CreateUserId = int(http_ctx.GetUserId(ctx))
	testCase.LastOpId = req.DefaultInt("last_op_id", 0)
	testCase.TestScript = req.String("test_script")
	testCase.TestSketchMap = req.String("test_sketch_map")
	// testCase.TestToolParams = req.String("test_tool_params")
	testCase.CreateTime = int(time.Now().Unix())
	id, err := mysql_model.GetBydName(testCase.Name, testCase.ScenarioId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	if id > 0 {
		response.RenderFailure(ctx, errors.New("测试用例名重复了"))
		return
	}
	if _, err := testCase.Create(demandChapterIds); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, orm.StructToMap(*testCase))
	}
}

/**
 * apiType http
 * @api {put} /api/v1/knowledge_test_cases/:id 修改knowledge测试用例
 * @apiVersion 0.1.0
 * @apiName Update
 * @apiGroup KnowledgeTestCase
 *
 * @apiDescription 修改一条knowledge测试用例
 *
 * @apiUse authHeader
 *
 * @apiParam {string} 			name  			测试用例名称(该字段必须传)
 * @apiParam {array} 			demand_chapter_id  			章节id组成的数组[3,4,5]
 * @apiParam {string} 			module_id  	测试组件/测试分类
 * @apiParam {int}       scenario_id    	安全检测场景
 * @apiParam {int}		demand_id			关联安全需求id，需求条目存另一个表
 * @apiParam {string} 		objective			测试目的
 * @apiParam {string} 		input			外部输入
 * @apiParam {string} 		test_procedure			测试步骤
 * @apiParam {string} 		test_standard			验证标准
 * @apiParam {int} 		level			测试难度
 * @apiParam {int} 		test_case_level			测试用例级别
 * @apiParam {int} 		test_method			测试方式，1黑盒 2灰盒 3白盒
 * @apiParam {int} 		auto_test_level			自动化测试程度 1人工 2半自动化 3自动化
 * @apiParam {string} 		test_tool			测试工具
 * @apiParam {string} 		test_tool_name			测试工具名称
 * @apiParam {string} 		test_tools_id			测试工具id
 * @apiParam {string} 		task_param			任务参数
 * @apiParam {int} 		create_user_id			创建用户id
 * @apiParam {int} 		last_op_id			最近操作用户id
 * @apiParam {string}   content			变更内容
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *	"name":"测试用例名称(该字段必须传)",
 *  "demand_chapter_id":[3,4,5]
 *	"module_id":"12",
 *	"scenario_id":2,
 *	"demand_id":3,
 *	"objective":"abc",
 *	"input":"",
 *	"test_procedure":"",
 *	"test_standard":"",
 *	"level":1,
 *	"test_case_level":1,
 *	"test_method":1,
 *	"auto_test_level":1,
 *	"test_tool":"",
 *	"test_tool_name":"",
 *	"test_tools_id":"",
 *	"test_script":"6189f5a8e830c613130eb1bc|6189f5a8e830c613130eb1ba",
 *	"test_sketch_map":"6189f6dae830c613130eb1c4|6189f6dfe830c613130eb1c6",
 *	"task_param":"",
 *	"create_user_id":0,
 *	"last_op_id":0
 *	"content":"修改备注"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *	"code": 0,
 *	"data": {
 *		 "auto_test_level": 1,
 *		 "case_uuid": "",
 *		 "create_time": 1635930902,
 *		 "create_user_id": 0,
 *		 "demand_id": 3,
 *		 "id": 118,
 *		 "input": "",
 *		 "last_op_id": 0,
 *		 "level": 1,
 *		 "module_id": "12",
 *		 "name": "测试用例名称",
 *		 "objective": "abc",
 *		 "scenario_id": 2,
 *		 "tag": "",
 *		 "task_param": "",
 *		 "test_case_level": 1,
 *		 "test_method": 1,
 *		 "test_param": "[{"case_type":"jar","block_name":"captureScreen","case_id":"hg_CASE_6","ret_of_crash":"","test_level":1,"name":"截图","time_out":10,"upload_crash_log":false}]",
 *		 "test_procedure": "",
 *		 "test_result": "",
 *		 "test_script": "",
 *		 "test_sketch_map": "",
 *		 "test_standard": "",
 *		 "test_tool": "",
 *		 "test_tools_id": ""
 *	},
 *	"msg": ""
 * }
 */
func (this KnowledgeTestCaseController) Update(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	id := request.ParamInt(ctx, "id")
	demandChapterIds := req.SliceInt("demand_chapter_id")
	testCase := new(mysql_model.KnowledgeTestCase)

	cols := []string{}
	if request.IsExist(ctx, "name") {
		testCase.Name = req.MustString("name")
		cols = append(cols, "name")
	}
	if request.IsExist(ctx, "module_id") {
		testCase.ModuleId = req.String("module_id")
		cols = append(cols, "module_id")
	}
	if request.IsExist(ctx, "scenario_id") {
		testCase.ScenarioId = req.MustInt("scenario_id")
		cols = append(cols, "scenario_id")
	}
	if request.IsExist(ctx, "demand_id") {
		testCase.DemandId = req.MustInt("demand_id")
		cols = append(cols, "demand_id")
	}
	if request.IsExist(ctx, "objective") {
		testCase.Objective = req.MustString("objective")
		cols = append(cols, "objective")
	}
	if request.IsExist(ctx, "input") {
		testCase.Input = req.MustString("input")
		cols = append(cols, "input")
	}
	if request.IsExist(ctx, "test_procedure") {
		testCase.TestProcedure = req.MustString("test_procedure")
		cols = append(cols, "test_procedure")
	}
	if request.IsExist(ctx, "test_standard") {
		testCase.TestStandard = req.MustString("test_standard")
		cols = append(cols, "test_procedure")
	}
	if request.IsExist(ctx, "level") {
		testCase.Level = req.MustInt("level")
		cols = append(cols, "level")
	}
	if request.IsExist(ctx, "test_case_level") {
		testCase.TestCaseLevel = req.MustInt("test_case_level")
		cols = append(cols, "test_case_level")
	}
	if request.IsExist(ctx, "test_method") {
		testCase.TestMethod = req.MustInt("test_method")
		cols = append(cols, "test_method")
	}

	if request.IsExist(ctx, "test_tool") {
		testCase.TestTool = req.String("test_tool")
		cols = append(cols, "test_tool")
	}
	if request.IsExist(ctx, "test_tool_name") {
		testCase.TestToolName = req.String("test_tool_name")
		cols = append(cols, "test_tool_name")
	}
	if request.IsExist(ctx, "test_tools_id") {
		testCase.TestToolsId = req.String("test_tools_id")
		cols = append(cols, "test_tools_id")
	}
	// if testCase.TestToolsId != "" {
	// 	params := qmap.QM{
	// 			"e__id": bson.ObjectIdHex(testCase.TestToolsId),
	// 		}
	// 	ormSession := mongo.NewMgoSessionWithCond(common.MC_TOOL, params)
	// 	data, _ := ormSession.GetOne()
	// 	if len(*data) > 0 {
	// 		rs := *data
	// 		testCase.TestToolName = rs["tool"].(string)
	// 	}
	// }
	if request.IsExist(ctx, "test_script") {
		testCase.TestScript = req.String("test_script")
		cols = append(cols, "test_script")
		if result := strings.Split(testCase.TestScript, "|"); len(result) > 0 && result[0] != "" {
			if fi, err := mongo.GridFSOpenId(common.MC_File, bson.ObjectIdHex(result[0])); err == nil {
				defer fi.Close()
				fileName := fi.Name()
				filesExt := path.Ext(fileName)
				testCase.CaseUuid = fileName[0 : len(fileName)-len(filesExt)]
				cols = append(cols, "case_uuid")
			}
		}
	}
	if request.IsExist(ctx, "test_sketch_map") {
		testCase.TestSketchMap = req.String("test_sketch_map")
		cols = append(cols, "test_sketch_map")
	}
	// if request.IsExist(ctx,"test_tool_params") {
	//	testCase.TestToolParams = req.String("test_tool_params")
	//  cols=append(cols,"test_tool_params")
	// }
	if request.IsExist(ctx, "task_param") {
		testCase.TaskParam = req.MustString("task_param")
		cols = append(cols, "task_param")
	}
	if request.IsExist(ctx, "create_user_id") {
		testCase.CreateUserId = req.MustInt("create_user_id")
		cols = append(cols, "create_user_id")
	}
	if request.IsExist(ctx, "last_op_id") {
		testCase.LastOpId = req.MustInt("last_op_id")
		cols = append(cols, "last_op_id")
	}
	// 如果为人工，不会显示测试工具
	if request.IsExist(ctx, "auto_test_level") {
		level := req.Int("auto_test_level")
		if level == 1 {
			testCase.TestToolName = ""
			testCase.TestToolsId = ""
			testCase.TestTool = ""
			cols = append(cols, "test_tool_name")
			cols = append(cols, "test_tools_id")
			cols = append(cols, "test_tool")
			testCase.AutoTestLevel = level
			cols = append(cols, "auto_test_level")
		} else {
			testCase.AutoTestLevel = level
			cols = append(cols, "auto_test_level")
		}
	}
	testCase.Id = id
	// cols 是"仅更新"而不是"更新的时候带上这些空值"
	if _, err := testCase.Update(demandChapterIds, cols...); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		// todo  更新测试用例记录
		history := new(mysql_model.KnowledgeTestCaseHistory)
		version, _ := history.FindLastVersion(testCase.Id)
		history = new(mysql_model.KnowledgeTestCaseHistory)
		history.TestCaseId = testCase.Id
		history.TimeStamp = int(time.Now().Unix())
		history.Content = req.String("content")
		history.OPId = int(http_ctx.GetUserId(ctx))
		history.Version = version
		history.Create()
		response.RenderSuccess(ctx, orm.StructToMap(*testCase))
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/knowledge_test_cases 批量删除knowledge测试用例
 * @apiVersion 0.1.0
 * @apiName BulkDelete
 * @apiGroup KnowledgeTestCase
 *
 * @apiDescription 批量删除knowledge测试用例
 *
 * @apiUse authHeader
 *
 * @apiParam {array} 			ids  			测试用例id组成的数组[3,4,5]
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  "ids":[3,4,5]
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *	"code": 0,
 *	"data": {
 *		 "number": 2
 *	},
 *	"msg": ""
 * }
 */
func (this KnowledgeTestCaseController) BulkDelete(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	successNum := 0
	if _, has := req.TrySlice("ids"); has {
		ids := req.SliceInt("ids")
		for _, id := range ids {
			_, err := new(mysql_model.KnowledgeTestCase).RemoveById(id)
			if err != nil {
				log.GetHttpLogLogger().Error(fmt.Sprintf("%v", err))
				break
			} else {
				successNum++
			}
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"number": successNum})
}

// 批量导入测试用例
func (this KnowledgeTestCaseController) Upload(ctx *gin.Context) {
	// 获取上传的文件
	file, _ := ctx.FormFile("file")
	// 上传文件到指定的路径
	dst := strings.Join([]string{"/tmp", file.Filename}, string(os.PathSeparator))
	if runtime.GOOS == "windows" {
		fmt.Println("Windows下调试")
		dst = strings.Join([]string{"tmp", file.Filename}, string(os.PathSeparator))
	}
	os.Remove(dst)
	if err := ctx.SaveUploadedFile(file, dst); err != nil {
		fmt.Println(err.Error())
	}

	outFilePath := strings.TrimSuffix(dst, ".zip")

	os.RemoveAll(outFilePath)
	if err := unzip(dst, outFilePath); err != nil {
		fmt.Println(err.Error())
	}

	xlsxList := getXlsxFilelist(outFilePath)
	total := 0
	failure := []string{}
	for _, xxx := range xlsxList {
		// fmt.Println("xxx:", xxx)
		f, err := excelize.OpenFile(xxx)
		if err != nil {
			fmt.Println(err.Error())
		}
		rows := f.GetRows("Sheet1")
		if err != nil {
			fmt.Println(err) // 打开文件失败
			return
		}
		for index, row := range rows {
			if index == 0 { // 跳过表头
				continue
			}
			testCase := new(mysql_model.KnowledgeTestCase)
			testCase.Name = row[0]   // A 测试用例名称
			if testCase.Name == "" { // 跳过空行
				continue
			}
			total++

			moduleName := row[1] // B 测试组件
			moduleType := row[2] // C 测试分类
			// 获取组件和分类id
			module, err := new(mongo_model.EvaluateModule).GetBydName(moduleName, moduleType)
			if err != nil {
				testCase.ModuleId = ""
				failure = append(failure, fmt.Sprintf("第%d条数据导入错误,组件和分类在平台中不存在", index))
				continue
			} else {
				testCase.ModuleId = module.Id.Hex()
			}
			// 获取场景id
			scenarioName := row[3] // D 安全检测场景
			if scenario, has := new(mysql_model.KnowledgeScenario).KnowledgeScenarioFindByName(scenarioName); has {
				testCase.ScenarioId = scenario.Id
			} else {
				failure = append(failure, fmt.Sprintf("第%d条数据导入错误,安全检测场景字段内容在平台中不存在", index))
				continue
			}

			// 测试用例名+场景id 重复
			id, _ := mysql_model.GetBydName(testCase.Name, testCase.ScenarioId)
			if id > 0 {
				failure = append(failure, fmt.Sprintf("第%d条数据导入错误,测试用例名称重复", index))
				continue
			}

			// todo 获取安全需求id
			demandName := row[4] // E 关联安全需求
			demand, has := new(mysql_model.KnowledgeDemand).KnowledgeDemandFindByName(demandName)
			if has {
				testCase.DemandId = demand.Id
			}
			// todo 安全条目
			demandChapterName := row[5] // F 安全需求条目
			demandChapterList := strings.Split(demandChapterName, ",")
			var demandChapterIds = make([]int, 0)
			warnLine := 0 // 错误的条数
			for _, demandChapterContent := range demandChapterList {
				tmp := strings.Split(demandChapterContent, " ")
				code := tmp[0]
				title := tmp[1]
				ttmp, has := new(mysql_model.KnowledgeDemandChapter).KnowledgeDemandChapterFindByName(demand.Id, code, title)
				if has {
					demandChapterIds = append(demandChapterIds, ttmp.Id)
				} else {
					warnLine++
				}
			}
			if warnLine > 0 {
				failure = append(failure, fmt.Sprintf("第%d条数据导入错误,不存在一条或多条对应的安全条目章节", index))
				continue
			}
			testCase.Objective = row[6] // G 测试目的
			testCase.Input = row[7]     // H 外部输入
			// if strings.Contains(row[7], " ") {
			//	failure = append(failure, fmt.Sprintf("第%d条数据导入错误,外部输入不能包含空格", index))
			//	continue
			// }
			if strings.HasPrefix(row[7], " ") || strings.HasPrefix(row[7], "\u3000") { // 中文全角空格
				failure = append(failure, fmt.Sprintf("第%d条数据导入错误,外部输入不能以空格作为开头", index))
				continue
			}
			if strings.HasSuffix(row[7], " ") || strings.HasSuffix(row[7], "\u3000") {
				failure = append(failure, fmt.Sprintf("第%d条数据导入错误,外部输入不能以空格作为结尾", index))
				continue
			}
			testCase.TestProcedure = row[8] // I 测试步骤
			testCase.TestStandard = row[9]  // J 验证标准
			level := row[10]                // K 测试难度
			switch level {
			case "低":
				testCase.Level = common.KNOWLEDGE_TEST_CASE_LEVEL_LOW
			case "中":
				testCase.Level = common.KNOWLEDGE_TEST_CASE_LEVEL_MIDDLE
			case "高":
				testCase.Level = common.KNOWLEDGE_TEST_CASE_LEVEL_HIGH
			default:
				failure = append(failure, fmt.Sprintf("第%d条数据导入错误,测试难度是错误的值", index))
				continue
			}
			testCaseLevel := row[11] // L 测试用例级别
			switch testCaseLevel {
			case "基础测试":
				testCase.TestCaseLevel = common.KNOWLEDGE_TEST_CASE_LEVEL_BASIC
			case "全面测试":
				testCase.TestCaseLevel = common.KNOWLEDGE_TEST_CASE_LEVEL_COMPLETE
			case "提高测试":
				testCase.TestCaseLevel = common.KNOWLEDGE_TEST_CASE_LEVEL_IMPROVE
			case "专家模式":
				testCase.TestCaseLevel = common.KNOWLEDGE_TEST_CASE_LEVEL_EXPERT
			default:
				failure = append(failure, fmt.Sprintf("第%d条数据导入错误,测试用例级别是错误的值", index))
				continue
			}
			method := row[12] // M 测试方式
			if method == "" {
				failure = append(failure, fmt.Sprintf("第%d条数据导入错误,测试方式为必填项", index))
				continue
			}
			switch method {
			case "黑盒":
				testCase.TestMethod = common.KNOWLEDGE_TEST_CASE_METHOD_BLACK
			case "灰盒":
				testCase.TestMethod = common.KNOWLEDGE_TEST_CASE_METHOD_GRAY
			case "白盒":
				testCase.TestMethod = common.KNOWLEDGE_TEST_CASE_METHOD_WHITE
			default:
				failure = append(failure, fmt.Sprintf("第%d条数据导入错误,测试方式只能填写黑盒、灰盒、白盒中的一项", index))
				continue
			}
			autoTestLevel := row[13] // N 自动化测试程度
			switch autoTestLevel {
			case "人工":
				testCase.AutoTestLevel = common.IS_TASK_CASE_MAN
			case "半自动化":
				testCase.AutoTestLevel = common.IS_TASK_CASE_SEMI
			case "自动化":
				testCase.AutoTestLevel = common.IS_TASK_CASE_AUTO
			default:
				failure = append(failure, fmt.Sprintf("第%d条数据导入错误,测试方式只能填写人工、半自动化、自动化中的一项", index))
				continue
			}
			testCase.TestToolName = row[14] // O 测试工具
			if testCase.TestToolName == "" && autoTestLevel != "人工" {
				failure = append(failure, fmt.Sprintf("第%d条数据导入错误,测试方式非人工时测试工具不能为空", index))
				continue
			}
			if testCase.TestToolName != "" && autoTestLevel == "人工" {
				failure = append(failure, fmt.Sprintf("第%d条数据导入错误,测试方式为人工时不应该存在测试工具", index))
				continue
			}
			switch testCase.TestToolName {
			case common.TOOL_HG_ANDROID_SCANNER_NAME:
				testCase.TestTool = common.TOOL_HG_ANDROID_SCANNER
			case common.TOOL_FIRMWARE_SCANNER_NAME:
				testCase.TestTool = common.TOOL_FIRMWARE_SCANNER
			case common.TOOL_VUL_SCANNER_NAME:
				testCase.TestTool = common.TOOL_VUL_SCANNER
			}
			testCase.TaskParam = row[15] // P 测试工具任务参数
			// 需要的路径是 zip包解压的位置 + attachment + 文件名称
			testScriptName := row[16] // Q 测试脚本
			// testScriptPath := path.Join(outFilePath, "attachment", testScriptName)
			testScriptPath := strings.Join([]string{outFilePath, "attachment", testScriptName}, string(os.PathSeparator))
			testScriptId := createFiletoMongo(testScriptPath)
			testCase.TestScript = testScriptId
			testCase.TestResult = row[17] // R 测试结果抽取
			var sb = new(strings.Builder)
			testSketchMaps := row[18] // S 测试环境搭建示意图
			fail := 0
			if testSketchMaps != "" {
				testSketchMapsList := strings.Split(testSketchMaps, ",")
				comma := []byte("|")
				for index, testSketchMapName := range testSketchMapsList { // 遍历每个图片文件名
					testSketchMapPath := strings.Join([]string{outFilePath, "attachment", testSketchMapName}, string(os.PathSeparator)) // 拼接图片绝对路径
					if illegal(testSketchMapPath) {
						continue
					}
					testSketchMapId := createFiletoMongo(testSketchMapPath)
					// tmpByte, _ := hex.DecodeString(testSketchMapId)
					if testSketchMapId == "" {
						fail++
					}
					if index == len(testSketchMapsList)-1 {
						sb.Write([]byte(testSketchMapId))
					} else {
						sb.Write([]byte(testSketchMapId))
						sb.Write(comma)
					}
				}
			}
			if fail > 0 {
				failure = append(failure, fmt.Sprintf("第%d条数据导入错误,一张或多张搭建示意图上传失败", index))
				continue
			}
			testCase.TestSketchMap = sb.String()
			var caseUuid = strings.TrimSuffix(testScriptName, ".zip")
			// if testScriptName != "" {
			//	switch testCase.TestToolName {
			//	case common.TOOL_HG_ANDROID_SCANNER_NAME:
			//		caseUuid = fmt.Sprintf("%s%s", common.CASE_HG_PRE, strings.TrimSuffix(testScriptName, ".zip"))
			//	}
			// }
			testCase.CaseUuid = caseUuid
			// _, err = testCase.CreateOnly()
			// if err != nil {
			//	panic(err)
			// }
			// 创建测试用例，场景，安全需求和章节的关系表
			testCase.CreateUserId = int(http_ctx.GetUserId(ctx))
			testCase.Create(demandChapterIds)
		}

	}
	if err := os.Remove(dst); err != nil {
		fmt.Println(err.Error())
	}
	if err := os.RemoveAll(outFilePath); err != nil {
		fmt.Println(err.Error())
	}
	success := total - len(failure)
	fail := len(failure)
	result := qmap.QM{
		"success": success,
		"fail":    fail,
		"reasons": failure,
	}
	response.RenderSuccess(ctx, result)
}

// 解压文件
func unzip(zipFile string, destDir string) error {
	zipReader, err := zip.OpenReader(zipFile)
	if err != nil {
		return err
	}
	defer zipReader.Close()
	var decodeName string
	for _, f := range zipReader.File {
		// fpath := filepath.Join(destDir, f.Name)
		if f.Flags == 0 {
			// 如果标志位是0  则是默认的本地编码   默认为gbk
			i := bytes.NewReader([]byte(f.Name))
			decoder := transform.NewReader(i, simplifiedchinese.GB18030.NewDecoder())
			content, _ := ioutil.ReadAll(decoder)
			decodeName = string(content)
		} else {
			// 如果标志位是 1 << 11也就是 2048  则是utf-8编码
			decodeName = f.Name
		}
		fpath := strings.Join([]string{destDir, decodeName}, string(os.PathSeparator))
		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
		} else {
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				return err
			}
			inFile, err := f.Open()
			if err != nil {
				return err
			}
			defer inFile.Close()

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer outFile.Close()

			_, err = io.Copy(outFile, inFile)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func getXlsxFilelist(dirPath string) []string {
	var xlsxList = make([]string, 0)
	fileInfos, err := ioutil.ReadDir(dirPath)
	if err != nil {
		fmt.Println(err.Error())
	}
	for _, fileInfo := range fileInfos {
		if index := strings.IndexAny(fileInfo.Name(), ".xlsx"); index != -1 {
			// xlsxList = append(xlsxList, path.Join(dirPath, fileInfo.Name()))
			xlsxList = append(xlsxList, strings.Join([]string{dirPath, fileInfo.Name()}, string(os.PathSeparator)))
		}
	}
	return xlsxList
}

func createFiletoMongo(filePath string) string {
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err.Error())
		return ""
	}
	fileName := path.Base(filePath)
	// 上传文件到mongo
	if fileId, err := mongo.GridFSUpload(common.MC_PROJECT, fileName, fileContent); err != nil {
		fmt.Println(err.Error())
		return ""
	} else {
		return fileId
	}
}

/**
 * apiType http
 * @api {post} /api/v1/knowledge_test_cases/:id/tag 给knowledge测试用例添加tag
 * @apiVersion 0.1.0
 * @apiName AddTag
 * @apiGroup KnowledgeTestCase
 *
 * @apiDescription 给knowledge测试用例添加tag
 *
 * @apiUse authHeader
 *
 * @apiParam {string} 			tag  			标签
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  "tag":"我是标签"
 * }
 *
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *   "code": 0,
 *   "data": {
 *        "auto_test_level": 3,
 *        "case_uuid": "CASE_51",
 *        "create_time": 0,
 *        "create_user_id": 0,
 *        "demand_chapter_id": [],
 *        "demand_chapter_name": [],
 *        "demand_id": 1,
 *        "demand_name": "车载信息交互系统信息安全技术要求",
 *        "id": 113,
 *        "input": "",
 *        "last_op_id": 0,
 *        "level": 0,
 *        "module_id": "",
 *        "module_name": "",
 *        "module_type_name": "",
 *        "name": "定位功能试验方法",
 *        "objective": "检查当应用软件调用定位功能时，车载信息交互系统是否在用户界面上有相应的状态提示。",
 *        "scenario_id": 1,
 *        "scenario_name": "车载合规测试",
 *        "tag": "",
 *        "task_param": "",
 *        "test_case_level": 2,
 *        "test_method": 2,
 *        "test_param": "",
 *        "test_procedure": "使用签名的应用调用定位功能，检查是否在用户界面上有相应的状态",
 *        "test_result": "",
 *        "test_script": "61820226e830c64451f07287",
 *        "test_scripts": [
 *            {
 *         *        "name": "CASE_51.zip",
 *         *        "value": "61820226e830c64451f07287"
 *            }
 *        ],
 *        "test_sketch_map": "",
 *        "test_standard": "",
 *        "test_tool": "",
 *        "test_tool_name": "车机检测工具",
 *        "test_tools_id": ""
 *   },
 *   "msg": ""
 * }
 */
func (this KnowledgeTestCaseController) AddTag(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	id := request.ParamInt(ctx, "id")
	tag := req.MustString("tag")
	testCase := new(mysql_model.KnowledgeTestCase)
	if _, err := testCase.UpdateTag(id, tag); err != nil {
		response.RenderFailure(ctx, err)
	}
	response.RenderSuccess(ctx, orm.StructToMap(*testCase))
}

func (this KnowledgeTestCaseController) Copy(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	// 拷贝测试用例
	scenarioId := req.MustInt("scenario_id")
	demandId := req.MustInt("demand_id")
	chapters := req.MustSlice("chapters")
	testCases := req.MustSlice("test_cases")
	num, _ := new(mysql_model.KnowledgeTestCase).Copy(scenarioId, demandId, chapters, testCases)
	result := qmap.QM{
		"success": num,
	}
	response.RenderSuccess(ctx, result)
}

/**
 * apiType http
 * @api {get} /api/v1/knowledge_test_this_tools/:id/this_tool_list 场景中涉及到的测试工具
 * @apiVersion 0.1.0
 * @apiName ScenarioCaseTool
 * @apiGroup KnowledgeTestCase
 *
 * @apiDescription 安全检测场景主页中测试工具的下拉列表
 *
 * @apiUse authHeader
 *
 * @apiParam {string} 			id  			场景id(该字段必须传)
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *	curl http://127.0.0.1:82/api/v1/knowledge_test_this_tools/2/this_tool_list
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "id": 0,
 *             "name": "",
 *             "case_uuid": "",
 *             "module_id": "",
 *             "scenario_id": 0,
 *             "demand_id": 0,
 *             "objective": "",
 *             "input": "",
 *             "test_procedure": "",
 *             "test_standard": "",
 *             "level": 0,
 *             "test_case_level": 0,
 *             "test_method": 0,
 *             "auto_test_level": 0,
 *             "test_tool": "",
 *             "test_tool_name": "",
 *             "test_tools_id": "",
 *             "task_param": "",
 *             "test_script": "",
 *             "test_result": "",
 *             "test_sketch_map": "",
 *             "create_user_id": 0,
 *             "last_op_id": 0,
 *             "create_time": 0,
 *             "test_param": "",
 *             "tag": ""
 *         }
 *     ],
 *     "msg": ""
 * }
 */
func (this KnowledgeTestCaseController) ScenarioToolList(ctx *gin.Context) {
	scenarioId := request.ParamInt(ctx, "id")
	UserId := http_ctx.GetUserId(ctx)
	KnowledgeTestCases := new(mysql_model.KnowledgeTestCase).GetToolListByScenarioId(scenarioId, int(UserId))
	response.RenderSuccess(ctx, KnowledgeTestCases)
}
func illegal(s string) bool {
	if strings.HasSuffix(s, string(os.PathSeparator)) { // 目录文件
		return true
	}
	if strings.Contains(s, strings.Join([]string{string(os.PathSeparator), "."}, "")) { // 挂载点文件
		return true
	}
	return false
}
