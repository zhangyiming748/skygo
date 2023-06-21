package mongo_model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"
	"skygo_detection/mysql_model"

	"github.com/globalsign/mgo/bson"
)

type EvaluateTestCase struct {
	Id             string          `bson:"_id,omitempty"  json:"_id"`                  // 测试用例编码
	Name           string          `bson:"name" json:"name"`                           // 测试用例名称
	ModuleTypeId   string          `bson:"module_type_id" json:"module_type_id"`       // 测试分类ID
	Objective      string          `bson:"objective" json:"objective"`                 // 测试目的
	ExternalInput  string          `bson:"external_input" json:"external_input"`       // 外部输入
	TestProcedure  string          `bson:"test_procedure" json:"test_procedure"`       // 测试步骤
	TestStandard   string          `bson:"test_standard" json:"test_standard"`         // 测试标准
	Level          int             `bson:"level" json:"level"`                         // 测试难度, 低、中、高
	TestCaseLevel  string          `bson:"test_case_level" json:"test_case_level"`     // 测试用例级别
	TestMethod     string          `bson:"test_method" json:"test_method"`             // 测试方法， 黑盒、灰盒、白盒
	AutoTestLevel  string          `bson:"auto_test_level" json:"auto_test_level"`     // 自动化测试程度, 取值范围（"自动化"、"人工"）
	TestScript     []TestValue     `bson:"test_script" json:"test_script"`             // 测试脚本
	TestTools      []TestValue     `bson:"test_tools" json:"test_tools"`               // 测试工具
	TestSketchMap  []TestValue     `bson:"test_sketch_map" json:"test_env_diagrams"`   // 测试环境搭建示意图
	Content        string          `bson:"content" json:"content"`                     // 内容变更
	ContentVersion string          `bson:"content_version" json:"content_version"`     // 内容变更版本
	OpId           int             `bson:"op_id" json:"op_id"`                         // 操作人id
	LastUpdateOpId int             `bson:"last_update_op_id" json:"last_update_op_id"` // 最后更新人id
	Status         int             `bson:"status" json:"status"`                       // 测试用例状态
	UpdateTime     int             `bson:"update_time" json:"update_time"`             // 更新时间
	CreateTime     int             `bson:"create_time" json:"create_time"`             // 创建时间
	SearchContent  string          `bson:"search_content" json:"search_content"`       // 查询条件
	BlockList      []TestCaseBlock `bson:"block_list" json:"block_list"`               // 测试用例block列表
}

const EvaluateTestCaseAutoTestLevelAuto = "自动化"
const EvaluateTestCaseAutoTestLevelInter = "人工"

type TestValue struct {
	Name  string `bson:"name" json:"name"`   // 测试工具名字、测试脚本名称、测试环境搭建示意图名称
	Value string `bson:"value" json:"value"` // 测试工具id、测试脚本id、测试环境搭建示意图名称id
}

func (this *EvaluateTestCase) Create(rawInfo qmap.QM, opId int) (*EvaluateTestCase, error) {
	// 判断库中是否有重名的测试用例
	if name := rawInfo.String("name"); name != "" {
		params := qmap.QM{
			"e_name": name,
		}
		n, err := orm_mongo.NewWidgetWithParams(common.MC_TEST_CASE, params).Count()
		if err != nil {
			return nil, err
		}
		if n != 0 {
			return nil, errors.New(fmt.Sprintf("已存在该测试用例：%s", name))
		}
	}
	this.Name = rawInfo.String("name")
	this.ModuleTypeId = rawInfo.String("module_type_id")
	this.Objective = rawInfo.String("objective")
	this.ExternalInput = rawInfo.String("external_input")
	this.TestProcedure = rawInfo.String("test_procedure")
	this.TestStandard = rawInfo.String("test_standard")
	this.Level = rawInfo.DefaultInt("level", 0)
	this.TestCaseLevel = rawInfo.String("test_case_level")
	this.TestMethod = rawInfo.String("test_method")
	this.AutoTestLevel = rawInfo.String("auto_test_level")
	if content, has := rawInfo.TryString("content"); has {
		this.Content = content
	} else {
		this.Content = "创建"
	}
	this.ContentVersion = changeDiffContentVersion("")
	this.OpId = opId
	this.Status = common.STATUS_EXIST
	this.LastUpdateOpId = opId
	this.UpdateTime = int(custom_util.GetCurrentMilliSecond())
	this.CreateTime = int(custom_util.GetCurrentMilliSecond())
	// 获取测试组件名称计算测试用例编号
	moduleName, moduleType := this.GetModuleName(this.ModuleTypeId)
	this.Id = this.GetTestCaseId(moduleName, moduleType, "M")
	// 组件查询条件
	this.SearchContent = fmt.Sprintf("%s_%s_%s_%s_%s", this.Id, this.Name, moduleName, moduleType, this.Objective)
	// 处理测试用例
	{
		scripts := []TestValue{}
		testScripts := rawInfo.Slice("test_script")
		for _, testScript := range testScripts {
			script := new(TestValue)
			tmp := testScript.(map[string]interface{})
			ttmp := qmap.QM(tmp)
			script.Name = ttmp.String("name")
			script.Value = ttmp.String("value")
			scripts = append(scripts, *script)
		}
		this.TestScript = scripts
	}
	// 处理测试工具
	{
		tools := []TestValue{}
		testTools := rawInfo.Slice("test_tools")
		for _, testTool := range testTools {
			tool := new(TestValue)
			tmp := testTool.(map[string]interface{})
			ttmp := qmap.QM(tmp)
			tool.Name = ttmp.String("name")
			tool.Value = ttmp.String("value")
			tools = append(tools, *tool)
		}
		this.TestTools = tools
	}
	// 处理测试环境示意图
	{
		sketchMaps := []TestValue{}
		envDiagrams := rawInfo.Slice("test_sketch_map")
		for _, envDiagram := range envDiagrams {
			sketchMap := new(TestValue)
			tmp := envDiagram.(map[string]interface{})
			ttmp := qmap.QM(tmp)
			sketchMap.Name = ttmp.String("name")
			sketchMap.Value = ttmp.String("value")
			sketchMaps = append(sketchMaps, *sketchMap)
		}
		this.TestSketchMap = sketchMaps
	}

	if err := mongo.NewMgoSession(common.MC_TEST_CASE).Insert(this); err == nil {
		// 更新历史记录
		{
			rawInfo := qmap.QM{
				"test_case_id":    this.Id,
				"content":         this.Content,
				"content_version": this.ContentVersion,
				"op_id":           this.OpId,
				"timestamp":       this.UpdateTime,
			}
			if _, err := new(TestCaseHistoryContent).Create(rawInfo); err != nil {
				return nil, err
			}
		}
		return this, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateTestCase) Update(id string, rawInfo qmap.QM, opId int) (*EvaluateTestCase, error) {
	params := qmap.QM{
		"e__id": id,
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_TEST_CASE, params)
	if err := mongoClient.One(&this); err == nil {
		if val, has := rawInfo.TryString("objective"); has {
			this.Objective = val
		}
		if val, has := rawInfo.TryString("external_input"); has {
			this.ExternalInput = val
		}
		if val, has := rawInfo.TryString("test_procedure"); has {
			this.TestProcedure = val
		}
		if val, has := rawInfo.TryString("test_standard"); has {
			this.TestStandard = val
		}
		if val, has := rawInfo.TryInt("level"); has {
			this.Level = val
		}
		if val, has := rawInfo.TryString("test_case_level"); has {
			this.TestCaseLevel = val
		}
		if val, has := rawInfo.TryString("test_method"); has {
			this.TestMethod = val
		}
		if val, has := rawInfo.TryString("auto_test_level"); has {
			this.AutoTestLevel = val
		}
		if val, has := rawInfo.TrySlice("test_script"); has {
			scripts := []TestValue{}
			testScripts := val
			for _, testScript := range testScripts {
				script := new(TestValue)
				tmp := testScript.(map[string]interface{})
				ttmp := qmap.QM(tmp)
				script.Name = ttmp.String("name")
				script.Value = ttmp.String("value")
				scripts = append(scripts, *script)
			}
			this.TestScript = scripts
		}
		if val, has := rawInfo.TrySlice("test_tools"); has {
			tools := []TestValue{}
			testTools := val
			for _, testTool := range testTools {
				tool := new(TestValue)
				tmp := testTool.(map[string]interface{})
				ttmp := qmap.QM(tmp)
				tool.Name = ttmp.String("name")
				tool.Value = ttmp.String("value")
				tools = append(tools, *tool)
			}
			this.TestTools = tools
		}
		if val, has := rawInfo.TrySlice("test_sketch_map"); has {
			// 处理测试环境示意图
			diagrams := []TestValue{}
			envDiagrams := val
			for _, envDiagram := range envDiagrams {
				diagram := new(TestValue)
				tmp := envDiagram.(map[string]interface{})
				ttmp := qmap.QM(tmp)
				diagram.Name = ttmp.String("name")
				diagram.Value = ttmp.String("value")
				diagrams = append(diagrams, *diagram)
			}
			this.TestSketchMap = diagrams
		}
		if val, has := rawInfo.TryString("content"); has {
			this.Content = val
		}
		contentVersion := this.ContentVersion
		this.ContentVersion = changeDiffContentVersion(contentVersion)
		// 最后更新人
		this.LastUpdateOpId = opId
		this.UpdateTime = int(custom_util.GetCurrentMilliSecond())
		if err := mongoClient.Update(bson.M{"_id": this.Id}, this); err != nil {
			return nil, err
		} else {
			// 更改测试用例记录
			{
				rawInfo := qmap.QM{
					"test_case_id":    this.Id,
					"content":         this.Content,
					"content_version": this.ContentVersion,
					"op_id":           opId,
					"timestamp":       this.UpdateTime,
				}
				if _, err := new(TestCaseHistoryContent).Create(rawInfo); err != nil {
					return nil, err
				}
			}
			return this, nil
		}
	}
	return nil, nil
}

// 获取所有的分页数据
func (this *EvaluateTestCase) GetAll(queryParams string) (*qmap.QM, error) {
	mgoSession := mongo.NewMgoSession(common.MC_TEST_CASE).AddUrlQueryCondition(queryParams)
	mgoSession.AddSorter("create_time", 1)
	mgoSession.SetTransformFunc(this.testCaseTransform)
	if res, err := mgoSession.GetPage(); err == nil {
		return res, nil
	} else {
		return nil, err
	}
}

// 获取所有的数据
func (this *EvaluateTestCase) Get(queryParams string) (*[]map[string]interface{}, error) {
	mgoSession := mongo.NewMgoSession(common.MC_TEST_CASE).AddUrlQueryCondition(queryParams)
	mgoSession.SetTransformFunc(this.testCaseTransform)
	return mgoSession.SetLimit(10000).Get()
}

func (this *EvaluateTestCase) GetOne(id string) (*qmap.QM, error) {
	params := qmap.QM{
		"e__id": id,
	}
	mgoSession := mongo.NewMgoSessionWithCond(common.MC_TEST_CASE, params)
	mgoSession.SetTransformFunc(this.testCaseTransform)
	return mgoSession.GetOne()
}

// 真实删除测试用例
func (this *EvaluateTestCase) BulkDelete(rawIds []string) (*qmap.QM, error) {
	// 删除 测试用例
	effectNum := 0
	if len(rawIds) > 0 {
		match := bson.M{
			"_id": bson.M{"$in": rawIds},
		}
		if changeInfo, err := mongo.NewMgoSession(common.MC_TEST_CASE).RemoveAll(match); err == nil {
			new(TestCaseHistoryContent).BulkDelete(rawIds)
			effectNum = changeInfo.Removed
		} else {
			return nil, err
		}
	}
	// 删除测试用例的变更历史

	return &qmap.QM{"number": effectNum}, nil
}

// 假删除测试用例
func (this *EvaluateTestCase) BulkDeleteWithStatus(rawIds []string) (*qmap.QM, error) {
	// 删除 测试用例
	mongoClient := mongo.NewMgoSession(common.MC_TEST_CASE)
	for _, id := range rawIds {
		selector := qmap.QM{
			"_id": bson.ObjectIdHex(id),
		}
		update := bson.M{
			"$set": bson.M{
				"status": common.STATUS_DELETE,
			},
		}
		if err := mongoClient.Update(selector, update); err != nil {
			return nil, err
		}
	}
	return nil, nil
}

func changeDiffContentVersion(version string) string {
	if version == "" {
		return "v1.0"
	}
	if strings.Contains(version, "v") {
		ver := strings.TrimLeft(version, "v")
		if vver, err := strconv.ParseFloat(ver, 64); err == nil {
			vver += 0.1
			return fmt.Sprintf("v%0.1f", vver)
		}
	}
	return version
}

func (this *EvaluateTestCase) testCaseTransform(data qmap.QM) qmap.QM {
	id := data["_id"]
	rawInfo := qmap.QM{
		"test_case_id": id,
	}
	if diffContents, err := new(TestCaseHistoryContent).Get(rawInfo); err != nil {
		return data
	} else {
		data["diff_content"] = diffContents
	}
	if module, err := new(EvaluateModule).FindById(data.MustString("module_type_id")); err == nil {
		data["module_name"] = module.ModuleName
		data["module_type"] = module.ModuleType
	}
	// todo 临时脚本，刷search_content内容
	if searchContent := data.String("search_content"); searchContent == "" {
		searchContent = fmt.Sprintf("%s_%s_%s_%s_%s", data.String("_id"), data.String("name"), data.String("module_name"), data.String("module_type"), data.String("objective"))
		selector := bson.M{"_id": data.String("_id")}
		update := bson.M{
			"$set": bson.M{"search_content": searchContent},
		}
		mongo.NewMgoSession(common.MC_TEST_CASE).Update(selector, update)
		data["search_content"] = searchContent
	}

	// 查询创建这名字
	if opId := data.Int("op_id"); opId > 0 {
		// 查询操作人员信息
		if rsp, err := new(mysql_model.SysUser).GetUserInfo(opId); err == nil {
			if realname := rsp.String("realname"); realname != "" {
				data["op_name"] = realname
			} else {
				data["op_name"] = rsp.String("username")
			}

		} else {
			data["op_name"] = ""
		}
	} else {
		data["op_name"] = ""
	}
	return data
}

func (this *EvaluateTestCase) Upload(rawList []qmap.QM, opId int64) (int64, int64, string) {
	successNumber := int64(0)
	failureNumber := int64(0)
	errorInfo := ""
	errorLists := []qmap.QM{}

	if len(rawList) == 0 {
		return successNumber, failureNumber, errorInfo
	}
	for _, raw := range rawList {
		if _, err := this.CheckItemRules(raw); err != nil {
			errString := err.Error()
			errorList := qmap.QM{
				"row":         raw.Int("row_number"),
				"id":          raw.String("_id"),
				"name":        raw.String("name"),
				"module_name": raw.String("module_name"),
				"module_type": raw.String("module_type"),
				"err_code":    errString,
			}

			if errString != "" {
				errorLists = append(errorLists, errorList)
				failureNumber++
				continue
			}
		}

		if err := this.Upsert(raw, opId); err == nil {
			successNumber++
		} else {
			failureNumber++
			errorList := qmap.QM{
				"row":         raw.Int("row_number"),
				"id":          raw.String("_id"),
				"name":        raw.String("name"),
				"module_name": raw.String("module_name"),
				"module_type": raw.String("module_type"),
				"err_code":    err.Error(),
			}
			errorLists = append(errorLists, errorList)
		}
	}
	if errInfo, err := json.Marshal(errorLists); err == nil {
		errorInfo = string(errInfo)
	}
	return successNumber, failureNumber, errorInfo
}

func (this *EvaluateTestCase) Upsert(rawInfo qmap.QM, opId int64) (err error) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			err = errors.New(fmt.Sprintf("%v", recoverErr))
		}
	}()
	moduleName := rawInfo.MustString("module_name")
	moduleType := rawInfo.MustString("module_type")
	moduleTypeId := ""
	if module, err := new(EvaluateModule).Find("", moduleName, moduleType); err == nil {
		moduleTypeId = module.Id.Hex()
	} else {
		return errors.New("测试组件和测试分类不存在")
	}
	id := rawInfo.String("_id")
	name := rawInfo.MustString("name")
	objective := rawInfo.MustString("objective")
	level := rawInfo.MustInt("level")
	searchContent := fmt.Sprintf("%s_%s_%s_%s_%s", id, name, moduleName, moduleType, objective)
	item := qmap.QM{
		"_id":             id,
		"name":            name,
		"module_type_id":  moduleTypeId,
		"objective":       objective,
		"external_input":  rawInfo.String("external_input"),
		"test_procedure":  rawInfo.String("test_procedure"),
		"test_standard":   rawInfo.String("test_standard"),
		"level":           level,
		"test_case_level": rawInfo.MustString("test_case_level"),
		"test_method":     rawInfo.MustString("test_method"),
		"auto_test_level": rawInfo.MustString("auto_test_level"),
		"op_id":           opId,
		"search_content":  searchContent,
	}
	selector := bson.M{
		"_id": bson.M{"$eq": id},
	}
	upsertItem := bson.M{
		"$setOnInsert": bson.M{
			"create_time": custom_util.GetCurrentMilliSecond(),
		},
		"$set": item,
	}
	_, err = mongo.NewMgoSession(common.MC_TEST_CASE).Upsert(selector, upsertItem)
	// 增加测试用例记录
	{
		nowtime := int(custom_util.GetCurrentMilliSecond())
		rawInfo := qmap.QM{
			"test_case_id":    id,
			"content":         "excel导入创建",
			"content_version": "--",
			"op_id":           opId,
			"timestamp":       nowtime,
		}
		new(TestCaseHistoryContent).Create(rawInfo)
	}
	return err
}

func (this *EvaluateTestCase) GetTestCaseId(moduleName, moduleType, autoStatus string) string {
	if autoStatus == "a" || autoStatus == "A" {
		autoStatus = "A"
	} else if autoStatus == "m" || autoStatus == "M" {
		autoStatus = "M"
	} else {
		return ""
	}
	moudleNumber := common.Module.DefaultString(moduleName, "999")
	moudleTypeNumber := common.ModuleType.DefaultString(moduleType, "999")
	// id 格式 TC001001A001
	id := fmt.Sprintf("%s%s%s%s", common.FirstTestCase, moudleNumber, moudleTypeNumber, autoStatus)
	params := qmap.QM{
		"l__id": id,
	}
	// 从数据库里查询最大的测试用例编号
	widget := orm_mongo.NewWidgetWithParams(common.MC_TEST_CASE, params)
	widget.AddSorter("_id", 1)
	// 给测试用例加1
	if result, err := widget.Find(); err != nil || len(result) == 0 {
		return id + "001"
	} else {
		lastTestCase := result[0]
		lastTestCaseId := lastTestCase["_id"].(string)
		lastTestCaseIdSlise := strings.Split(lastTestCaseId, autoStatus)
		tmp := lastTestCaseIdSlise[1]
		num, _ := strconv.Atoi(tmp)
		return fmt.Sprintf("%s%03d", id, num+1)
	}
}

func (this *EvaluateTestCase) GetModuleName(moduleId string) (moduleName, moduleType string) {
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(moduleId),
	}
	result := qmap.QM{}
	orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_MODULE, params).One(&result)
	moduleName = result.MustString("module_name")
	moduleType = result.MustString("module_type")
	return
}

func (this *EvaluateTestCase) CheckItemRules(raw qmap.QM) (bool, error) {
	// 判断一下关键的字段存在是否存在
	var errString string
	// 测试组件必须存在
	moduleName := raw.String("module_name")
	if moduleName == "" {
		errString = errString + "不存在测试组件;\n"
	}
	// 测试分类必须存在
	moduleType := raw.String("module_type")
	if moduleType == "" {
		errString = errString + "不存在测试分类;\n "
	}
	// 测试用例名称必须存在
	if name := raw.String("name"); name == "" {
		errString = errString + "不存在测试用例名称;\n"
	}
	// 测试用例级别必须存在
	if testCaseLevel := raw.String("test_case_level"); testCaseLevel == "" {
		errString = errString + "不存在测试用例级别;\n"
	} else {
		switch testCaseLevel {
		case "基础测试", "全面测试", "提高测试", "专家模式":
		default:
			errString = errString + "测试用例级别应该为(基础测试，全面测试，提高测试，专家模式);\n"
		}
	}
	// 测试方法必须存在
	if testMethod := raw.String("test_method"); testMethod == "" {
		errString = errString + "不存在测试方式;\n"
	} else {
		switch testMethod {
		case "黑盒", "灰盒", "白盒":
		default:
			errString = errString + "测试方式应该为(黑盒，灰盒，白盒);\n"
		}
	}
	// 测试目的必须存在
	if objective := raw.String("objective"); objective == "" {
		errString = errString + "不存在测试目的;\n"
	}
	// 自动化测试程度必须存在
	if autoTestLevel := raw.String("auto_test_level"); autoTestLevel == "" {
		errString = errString + "不存在自动化测试程度;\n"
	} else {
		switch autoTestLevel {
		case "人工", "自动化":
		default:
			errString = errString + "自动化测试程度应该为(人工，自动化);\n"
		}
	}
	// 测试步骤必须存在
	if testProcedure := raw.String("test_procedure"); testProcedure == "" {
		errString = errString + "不存在测试步骤;\n"
	}
	// 测试标准必须存在
	if testStandard := raw.String("test_standard"); testStandard == "" {
		errString = errString + "不存在测试标准;\n"
	}
	var id string = raw.String("_id")

	if id == "" {
		id = this.GetTestCaseId(moduleName, moduleType, common.TEST_AUTO)

		// 当测试用例无id的时候 需要查一下是否重名
		evaluateModule, _ := new(EvaluateModule).Find("", moduleName, moduleType)
		evaluateModuleId := evaluateModule.Id.Hex()
		name := raw.String("name")
		params := qmap.QM{
			"e_module_type_id": evaluateModuleId,
			"e_name":           name,
		}
		if number, _ := orm_mongo.NewWidgetWithParams(common.MC_TEST_CASE, params).Count(); number != 0 {
			errString = errString + "测试用例导入失败，测试名称重名，请修改后重试;\n"
		}

	} else {
		// 库中查询id是否存在，存在就更新，不存在抛弃掉，不新建
		params := qmap.QM{
			"e__id": id,
		}
		if number, _ := orm_mongo.NewWidgetWithParams(common.MC_TEST_CASE, params).Count(); number == 0 {
			errString = errString + "不存在该测试用例ID;\n"
		}
	}
	raw["_id"] = id

	if errString != "" {
		return false, errors.New(errString)
	}
	return true, nil
}

// 根据一组ID，查询出对应的一组测试案例
func (this *EvaluateTestCase) FindModelsByIds(ids []string) []EvaluateTestCase {
	session := mongo.NewMgoSession(common.MC_TEST_CASE)

	params := bson.M{
		"_id": bson.M{
			"$in": ids,
		},
	}
	models := make([]EvaluateTestCase, 0)

	session.Session.Find(params).All(&models)
	return models
}
