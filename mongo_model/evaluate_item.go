package mongo_model

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"
)

// 测试用例基础表
type EvaluateItem struct {
	Id            string      `bson:"_id" json:"id"`                            // 测试用例id
	Name          string      `bson:"name" json:"name"`                         // 测试用例名称
	ProjectId     string      `bson:"project_id" json:"project_id"`             // 项目id
	AssetId       string      `bson:"asset_id" json:"asset_id"`                 // 所属资产id
	ModuleTypeId  string      `bson:"module_type_id" json:"module_type_id"`     // 测试分类ID
	Objective     string      `bson:"objective" json:"objective"`               // 测试目的
	ExternalInput string      `bson:"external_input" json:"external_input"`     // 外部输入
	TestProcedure string      `bson:"test_procedure" json:"test_procedure"`     // 测试步骤
	TestStandard  string      `bson:"test_standard" json:"test_standard"`       // 测试标准
	Level         int         `bson:"level" json:"level"`                       // 测试难度（0:默认、1:低、2:中、3:高）
	TestCaseLevel string      `bson:"test_case_level" json:"test_case_level"`   // 测试用例级别（基础测试、全面测试、提高测试、专家测试）
	TestMethod    string      `bson:"test_method" json:"test_method"`           // 测试方法（黑盒、白盒）
	AutoTestLevel string      `bson:"auto_test_level" json:"auto_test_level"`   // 自动化测试程度（自动化、人工）
	TestScript    interface{} `bson:"test_script" json:"test_script"`           // 测试脚本
	TestSketchMap interface{} `bson:"test_sketch_map" json:"test_sketch_map"`   // 测试环境示意图
	TestTime      int         `bson:"test_time" json:"test_time"`               // 测试时间
	TestCount     int         `bson:"test_count" json:"test_count"`             // 测试次数
	VulNumber     int         `bson:"vul_number" json:"vul_number"`             // 漏洞数量
	Status        int         `bson:"status" json:"status"`                     // 测试状态 （0:可创建任务 1:不可创建任务）
	TestStatus    int         `bson:"test_status" json:"test_status"`           // 测试状态 （0:待测试 1:测试完成,2待补充，3审核通过）
	AuditStatus   int         `bson:"audit_status" json:"audit_status"`         // 审核状态 （1:通过,0默认, -1:驳回）
	PreBind       string      `bson:"pre_bind" json:"pre_bind"`                 // 当前绑定任务ID
	LastTaskId    string      `bson:"last_task_id" json:"last_task_id"`         // 上一次任务ID
	IsPreBind     int         `bson:"is_pre_bind" json:"is_pre_bind"`           // 是否为预绑定，0正常，1预绑定
	IsPreDelete   int         `bson:"is_pre_delete" json:"is_pre_delete"`       // 是否为预解绑，0正常，1预解绑
	Tag           []string    `bson:"tag" json:"tag"`                           // 测试用例标签
	OpId          int64       `bson:"op_id" json:"op_id"`                       // 操作人id
	UpdateTime    int64       `bson:"update_time" json:"update_time"`           // 更新时间
	CreateTime    int64       `bson:"create_time" json:"create_time,omitempty"` // 创建时间
}

func (this *EvaluateItem) Create(ctx context.Context, rawInfo qmap.QM, opId int64) (*EvaluateItem, error) {
	// 测试项所属项目必须存在
	projectId := rawInfo.MustString("project_id")
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(projectId),
	}
	if _, err := mongo.NewMgoSessionWithCond(common.MC_PROJECT, params).GetOne(); err != nil {
		panic(fmt.Sprintf("项目: %s 不存在！", projectId))
	}

	// 测试项所属资产必须存在
	assetId := rawInfo.MustString("asset_id")
	params = qmap.QM{
		"e__id": assetId,
	}
	if _, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, params).GetOne(); err != nil {
		panic(fmt.Sprintf("资产: %s 不存在！", assetId))
	}

	name := rawInfo.MustString("name")
	params = qmap.QM{
		"e_project_id": projectId,
		"e_asset_id":   assetId,
		"e_name":       name,
	}
	if n, _ := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params).Count(); n != 0 {
		return nil, errors.New("该测试用例名称已被占用")
	}

	this.Name = name
	this.AssetId = assetId
	this.ProjectId = projectId
	this.Level = rawInfo.DefaultInt("level", 0)
	this.ModuleTypeId = rawInfo.MustString("module_type_id")
	this.Objective = rawInfo.MustString("objective")
	this.ExternalInput = rawInfo.String("external_input")
	this.TestCaseLevel = rawInfo.MustString("test_case_level")
	this.TestMethod = rawInfo.MustString("test_method")
	this.AutoTestLevel = rawInfo.MustString("auto_test_level")
	this.TestProcedure = rawInfo.MustString("test_procedure")
	this.TestStandard = rawInfo.MustString("test_standard")
	this.AuditStatus = common.EIAS_DEFAULT
	this.TestStatus = common.TIS_READY
	this.UpdateTime = custom_util.GetCurrentMilliSecond()
	this.CreateTime = custom_util.GetCurrentMilliSecond()
	this.OpId = opId
	// this.Id = this.GetItemId(assetId, this.ModuleTypeId, "M")
	if id, has := rawInfo.TryString("id"); has {
		this.Id = id
	} else {
		if suffixId, has := rawInfo.TryString("suffix_id"); has {
			this.Id = this.GetItemIdWithSuffixId(assetId, this.ModuleTypeId, suffixId)
		} else {
			this.Id = this.GetItemId(assetId, this.ModuleTypeId, "M")
		}
	}
	if val, has := rawInfo.TryInterface("test_script"); has {
		this.TestScript = val
	}
	if val, has := rawInfo.TryInterface("test_sketch_map"); has {
		this.TestSketchMap = val
	}

	if _, err := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_ITEM).InsertOne(ctx, this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateItem) Update(id string, rawInfo qmap.QM) (*EvaluateItem, error) {
	params := qmap.QM{
		"e__id": id,
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params)
	if err := mongoClient.One(&this); err == nil {
		if val, has := rawInfo.TryString("name"); has {
			this.Name = val
		}
		if val, has := rawInfo.TryString("asset_id"); has {
			this.AssetId = val
		}
		if val, has := rawInfo.TryString("module_type_id"); has {
			this.ModuleTypeId = val
		}
		if val, has := rawInfo.TryInt("level"); has {
			this.Level = val
		}
		if val, has := rawInfo.TryInt("status"); has {
			this.Status = val
		}
		if auditStatus, has := rawInfo.TryInt("audit_status"); has {
			this.AuditStatus = auditStatus
		}
		if val, has := rawInfo.TryString("objective"); has {
			this.Objective = val
		}
		if val, has := rawInfo.TryString("test_procedure"); has {
			this.TestProcedure = val
		}
		if val, has := rawInfo.TryString("test_standard"); has {
			this.TestStandard = val
		}
		if val, has := rawInfo.TryString("test_method"); has {
			this.TestMethod = val
		}
		if val, has := rawInfo.TryInt("vul_number"); has {
			this.VulNumber = val
		}
		if testTime, has := rawInfo.TryInt("test_time"); has {
			this.TestTime = testTime
		}
		if val, has := rawInfo.TryString("auto_test_level"); has {
			this.AutoTestLevel = val
		}
		if val, has := rawInfo.TryString("test_case_level"); has {
			this.TestCaseLevel = val
		}
		if val, has := rawInfo.TryString("external_input"); has {
			this.ExternalInput = val
		}
		if val, has := rawInfo.TryInterface("test_script"); has {
			this.TestScript = val
		}
		if val, has := rawInfo.TryInterface("test_sketch_map"); has {
			this.TestSketchMap = val
		}
		if val, has := rawInfo.TryInt("op_id"); has {
			this.OpId = int64(val)
		}
		if err := mongoClient.Update(bson.M{"_id": this.Id}, this); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Item not found")
	}
	return this, nil
}

func (this *EvaluateItem) Import(rawList []qmap.QM, opId int64) (int64, int64, string) {
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

func (this *EvaluateItem) Upsert(rawInfo qmap.QM, opId int64) (err error) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			err = errors.New(fmt.Sprintf("%v", recoverErr))
		}
	}()
	id := rawInfo.MustString("id")
	// 测试项所属项目必须存在
	projectId := rawInfo.MustString("project_id")
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(projectId),
	}
	if _, err := mongo.NewMgoSessionWithCond(common.MC_PROJECT, params).GetOne(); err != nil {
		return err
		// panic(fmt.Sprintf("项目: %s 不存在！", projectId))
	}

	// 测试项所属资产必须存在
	assetId := rawInfo.MustString("asset_id")
	params = qmap.QM{
		"e__id": assetId,
	}
	if _, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, params).GetOne(); err != nil {
		return err
		// panic(fmt.Sprintf("资产: %s 不存在！", assetId))
	}
	moduleTypeId := ""
	if module, err := new(EvaluateModule).Find("", rawInfo.MustString("module_name"), rawInfo.MustString("module_type")); err == nil {
		moduleTypeId = module.Id.Hex()
	} else {
		return errors.New("未找到测试组件和测试分类")
	}

	levelString := rawInfo.MustString("level")
	level := common.TEST_LEVEL_DEFAULT
	switch levelString {
	case "低":
		level = common.TEST_LEVEL_LOW
	case "中":
		level = common.TEST_LEVEL_MIDDLE
	case "高":
		level = common.TEST_LEVEL_HIGH
	default:
		level = common.TEST_LEVEL_DEFAULT
	}

	item := qmap.QM{
		"_id":              id,
		"name":             rawInfo.MustString("name"),
		"project_id":       rawInfo.MustString("project_id"),
		"asset_id":         rawInfo.MustString("asset_id"),
		"module_type_id":   moduleTypeId,
		"objective":        rawInfo.MustString("objective"),
		"external_input":   rawInfo.String("external_input"),
		"test_procedure":   rawInfo.String("test_procedure"),
		"test_standard":    rawInfo.String("test_standard"),
		"level":            level,
		"test_case_level":  rawInfo.MustString("test_case_level"),
		"test_method":      rawInfo.MustString("test_method"),
		"auto_test_level":  rawInfo.MustString("auto_test_level"),
		"op_id":            opId,
		"test_time":        0,
		"vul_number":       0,
		"evaluate_task_id": "",
		"pre_bind":         "",
		"status":           common.EIS_FREE,
		"test_status":      common.TIS_READY,
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
	_, err = mongo.NewMgoSession(common.MC_EVALUATE_ITEM).Upsert(selector, upsertItem)
	return err
}

// ===============================================//

// 将某一批测试用例与项目任务绑定
func (this *EvaluateItem) BindEvaluateTask(itemIds []string, evaluateTaskId string, isPreBind int) {
	// 第一步，将测试用例status设置为1，表示测试用例已被占用
	mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).Session
	selector := bson.M{"_id": bson.M{"$in": itemIds}}
	data := qmap.QM{
		"status":      common.EIS_INUSE,
		"pre_bind":    evaluateTaskId,
		"is_pre_bind": isPreBind,
		"test_status": common.TIS_READY,
		"test_time":   0,
	}
	if _, err := mgoSession.UpdateAll(selector, qmap.QM{"$set": data}); err != nil {
		fmt.Println(err)
		panic(err)
	}
}

// 将测试用例和任务信息，写入任务用例关系表
func (this *EvaluateItem) WriteItemToTaskItem(taskId string, opId int) {
	// 第一步，查出Item表中正式需要绑定的测试用例
	param := qmap.QM{
		"e_pre_bind":      taskId,
		"e_is_pre_bind":   common.NOT_PREBIND,
		"e_is_pre_delete": common.NOT_PREDEL,
	}
	taskInfo, err := new(EvaluateTask).GetOne(taskId)
	if err != nil {
		panic(err)
	}

	if itemList, err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).AddCondition(param).SetLimit(10000).Get(); err == nil {
		for _, item := range *itemList {
			testId := bson.NewObjectId()
			taskItemData := qmap.QM{
				"id":              testId,
				"project_id":      item["project_id"],
				"task_id":         taskId,
				"test_phase":      (*taskInfo)["test_phase"],
				"item_id":         item["_id"],
				"name":            item["name"],
				"asset_id":        item["asset_id"],
				"module_type_id":  item["module_type_id"],
				"test_method":     item["test_method"],
				"auto_test_level": item["auto_test_level"],
			}
			_, err := new(EvaluateTaskItem).Create(taskItemData, opId)
			if err != nil {
				panic(err)
			}

			// 创建正式关系后，复制任务下测试用例之前已产生的漏洞，到任务漏洞表中
			// 加注释： 旧测试记录id、 新测试记录id、操作人id
			if err := new(EvaluateTaskVulnerability).CopyItemVulsToNewTask(item["_id"].(string), testId.Hex(), int64(opId)); err != nil {
				panic(err)
			}
		}

	}

}

// 将测试用例与任务预解绑
func (this *EvaluateItem) PreDeleteTask(itemIds []string) {
	mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).Session
	selector := bson.M{"_id": bson.M{"$in": itemIds}}
	set := qmap.QM{
		"is_pre_delete": common.IS_PREDEL,
	}
	update := qmap.QM{"$set": set}
	if _, err := mgoSession.UpdateAll(selector, update); err != nil {
		println(err)
	}
}

// 将预绑定的正式绑定
// func (this *EvaluateItem) bindEvaluateTask(ids []string) {
// 	mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).Session
// 	selector := bson.M{"_id": bson.M{"$in": ids}}
// 	update := qmap.QM{"$set": qmap.QM{"is_pre_bind": common.NOT_PREBIND}}
// 	fmt.Println(selector)
// 	fmt.Println(update)
// 	if _, err := mgoSession.UpdateAll(selector, update); err != nil {
// 		println(err)
// 	}
//
// }

// 将某一批测试用例与项目任务解绑
// func (this *EvaluateItem) UnbindEvaluateTask(itemIds []string) {
// 	mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).Session
// 	selector := bson.M{"_id": bson.M{"$in": itemIds}}
// 	set := qmap.QM{
// 		"status":        common.EIS_FREE, // 将测试用例状态，变更为可绑定
// 		"is_pre_delete": common.NOT_PREDEL,
// 		"pre_bind":      "",
// 	}
// 	update := qmap.QM{"$set": set}
// 	fmt.Println(selector)
// 	fmt.Println(update)
// 	if _, err := mgoSession.UpdateAll(selector, update); err != nil {
// 		println(err)
// 	}
// }

func (this *EvaluateItem) GetProjectRelatedAssets(projectId, ids string) []string {
	query := bson.M{
		"project_id": projectId,
		"_id": bson.M{
			"$in": strings.Split(ids, "|"),
		},
	}
	assetIds := []string{}
	if err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).Session.Find(query).Distinct("asset_id", &assetIds); err == nil {
		return assetIds
	} else {
		return []string{}
	}
}

func (this *EvaluateItem) GetTaskRelatedAssets(taskId string) []string {
	query := bson.M{
		"evaluate_task_id": taskId,
	}
	assetIds := []string{}
	if err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).Session.Find(query).Distinct("asset_id", &assetIds); err == nil {
		return assetIds
	} else {
		return []string{}
	}
}

func (this *EvaluateItem) GetOne(id string) (*EvaluateItem, error) {
	param := qmap.QM{
		"e__id": id,
	}
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, param).One(this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateItem) GetItemId(assetId, moduleId, autoCode string) string {
	if autoCode == "a" || autoCode == "A" {
		autoCode = "A"
	} else if autoCode == "m" || autoCode == "M" {
		autoCode = "M"
	} else {
		return ""
	}
	// 获取资产id assetId
	// 获取模块编码
	// 获取模块类型编码
	var moduleNameCode interface{}
	var moduleTypeCode interface{}
	{
		params := qmap.QM{
			"e__id": bson.ObjectIdHex(moduleId),
		}
		if module, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, params).GetOne(); err == nil {
			moduleNameCode = (*module)["module_name_code"]
			moduleTypeCode = (*module)["module_type_code"]
		}
	}
	// 获取 A 或 M  autoCode
	// 获取最后三位编码 id 格式 TC001001A001
	id := fmt.Sprintf("%s%s%s%s%s", assetId, common.FirstTestCase, moduleNameCode, moduleTypeCode, autoCode)
	params := qmap.QM{
		"l__id": id,
	}
	// 从数据库里查询最大的测试用例编号
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params)
	mongoClient.AddSorter("_id", 1)
	// 给测试用例加id递增
	if result, err := mongoClient.Get(); err != nil || len(*result) == 0 {
		return id + "001"
	} else {
		lastTestCase := (*result)[0]
		lastTestCaseId := lastTestCase["_id"].(string)
		lastTestCaseIdSlise := strings.Split(lastTestCaseId, autoCode)
		tmp := lastTestCaseIdSlise[len(lastTestCaseIdSlise)-1]
		num, _ := strconv.Atoi(tmp)
		return fmt.Sprintf("%s%03d", id, num+1)
	}
}

func (this *EvaluateItem) GetItemIdWithModule(assetId, moduleName, moduleType, autoCode string) string {
	if autoCode == "a" || autoCode == "A" {
		autoCode = "A"
	} else if autoCode == "m" || autoCode == "M" {
		autoCode = "M"
	} else {
		return ""
	}
	// 获取资产id assetId
	// 获取模块编码
	// 获取模块类型编码
	var moduleNameCode interface{}
	var moduleTypeCode interface{}
	{
		params := qmap.QM{
			"e_module_name": moduleName,
			"e_module_type": moduleType,
		}
		if module, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, params).GetOne(); err == nil {
			moduleNameCode = (*module)["module_name_code"]
			moduleTypeCode = (*module)["module_type_code"]
		}
	}
	// 获取 A 或 M  autoCode
	// 获取最后三位编码 id 格式 TC001001A001
	id := fmt.Sprintf("%s%s%s%s%s", assetId, common.FirstTestCase, moduleNameCode, moduleTypeCode, autoCode)
	params := qmap.QM{
		"l__id": id,
	}
	// 从数据库里查询最大的测试用例编号
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params)
	mongoClient.AddSorter("_id", 1)
	// 给测试用例加id递增
	if result, err := mongoClient.Get(); err != nil || len(*result) == 0 {
		return id + "001"
	} else {
		lastTestCase := (*result)[0]
		lastTestCaseId := lastTestCase["_id"].(string)
		lastTestCaseIdSlise := strings.Split(lastTestCaseId, autoCode)
		tmp := lastTestCaseIdSlise[1]
		num, _ := strconv.Atoi(tmp)
		return fmt.Sprintf("%s%03d", id, num+1)
	}
}

func (this *EvaluateItem) GetItemIdWithSuffixId(assetId, moduleId, suffixId string) string {
	var moduleNameCode interface{}
	var moduleTypeCode interface{}
	{
		params := qmap.QM{
			"e__id": bson.ObjectIdHex(moduleId),
		}
		if module, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, params).GetOne(); err == nil {
			moduleNameCode = (*module)["module_name_code"]
			moduleTypeCode = (*module)["module_type_code"]
		}
	}
	// 获取 A 或 M  autoCode
	// 获取最后三位编码 id 格式 TC001001A001
	id := fmt.Sprintf("%s%s%s%s%s", assetId, common.FirstTestCase, moduleNameCode, moduleTypeCode, suffixId)
	return id
}

func (this *EvaluateItem) CheckItemRules(raw qmap.QM) (bool, error) {
	// 判断一下关键的字段存在是否存在
	var errString string
	// 测试用例名称必须存在
	name := raw.String("name")
	if name == "" {
		return false, errors.New("不存在测试用例名称;")
	}
	// 资产ID必须存在
	assetId := raw.String("asset_id")
	if assetId == "" {
		errString = errString + "excle下，不存在资产ID;"
	}
	projectId := raw.String("project_id")
	params := qmap.QM{
		"e__id":        assetId,
		"e_project_id": projectId,
	}
	if number, _ := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, params).Count(); number == 0 {
		errString = errString + "该项目下，不存在资产ID;\n"
	}
	// 测试组件必须存在
	moduleName := raw.String("module_name")
	if moduleName == "" {
		errString = errString + "不存在测试组件;\n"
	}
	// 测试分类必须存在
	moduleType := raw.String("module_type")
	if moduleType == "" {
		errString = errString + "不存在测试分类;\n"
	}
	// 测试用例级别必须存在
	if testCaseLevel := raw.String("test_case_level"); testCaseLevel == "" {
		errString = errString + "不存在测试用例级别;\n"
	} else {
		switch testCaseLevel {
		case "基础测试", "全面测试", "提高测试", "专家模式":
		default:
			errString = errString + "测试用例级别不属于(基础测试，全面测试，提高测试，专家模式);\n"
		}
	}

	// 测试方法必须存在
	if testMethod := raw.String("test_method"); testMethod == "" {
		errString = errString + "不存在测试方式;\n"
	} else {
		switch testMethod {
		case "黑盒", "灰盒", "白盒":
		default:
			errString = errString + "测试方式不属于(黑盒，灰盒，白盒);\n"
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
			errString = errString + "自动化测试程度输入有误，应该为(人工，自动化);\n"
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
	// 获取测试用例的id
	// id 是 TCxxxxxxxxx
	var id string = raw.String("_id")
	// itemId 是 assetIdTCxxxxxxxxx
	var itemId string
	if id == "" {
		itemId = this.GetItemIdWithModule(assetId, moduleName, moduleType, common.TEST_AUTO)
		// 如果不存在用例ID，判断一下 该资产下是否有重名的用例
		params = qmap.QM{
			"e_name":     name,
			"e_asset_id": assetId,
		}
		if number, _ := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params).Count(); number != 0 {
			errString = errString + "该资产下，测试用例重名;\n"
		}
	} else {
		itemId = fmt.Sprintf("%s%s", assetId, id)
		// 库中查询itemId是否存在，存在就更新，不存在再接着判断
		params := qmap.QM{
			"e__id": itemId,
		}
		if number, _ := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params).Count(); number == 0 {
			// 查询测试用例库里是否存在这条用例，如果存在就插入
			paramsByTestCase := qmap.QM{
				"e__id": id,
			}
			if n, _ := mongo.NewMgoSessionWithCond(common.MC_TEST_CASE, paramsByTestCase).Count(); n == 0 {
				errString = errString + "不存在该测试用例ID;\n"
			}
		}
	}
	raw["id"] = itemId
	if errString != "" {
		return false, errors.New(errString)
	}
	return true, nil
}

// 检测对应资产的测试组件分类下，是否还存在用例，如果不存在则删除该资产对应的组件和分类
func (this *EvaluateItem) CheckAssetItemRelation(itemId string) (bool, error) {
	itemInfo, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, qmap.QM{"e__id": itemId}).GetOne()
	if err == nil {
		params := qmap.QM{
			"e_project_id":     (*itemInfo)["project_id"],
			"e_module_type_id": (*itemInfo)["module_type_id"],
		}
		records, _ := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params).Count()
		if records <= 1 { // 判断同一个组件分类下是否还有测试用例，如果没有的话，删除对应资产的组件分裂
			assetParams := qmap.QM{
				"e__id": (*itemInfo)["asset_id"],
			}
			if assetInfo, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, assetParams).GetOne(); err == nil {
				moduleTypeId := []string{}
				selectStruct := []interface{}{}
				if len((*assetInfo)["module_type_id"].([]interface{})) > 0 {
					for _, id := range (*assetInfo)["module_type_id"].([]interface{}) {
						if id.(string) != (*itemInfo)["module_type_id"] {
							moduleTypeId = append(moduleTypeId, id.(string))
						}
					}
				}

				if len((*assetInfo)["select_struct"].([]interface{})) > 0 {
					for _, structItem := range (*assetInfo)["select_struct"].([]interface{}) {
						if structItem.([]interface{})[1].(string) != (*itemInfo)["module_type_id"] {
							selectStruct = append(selectStruct, structItem.([]interface{}))
						}
					}
				}

				selector := bson.M{
					"_id": (*assetInfo)["_id"],
				}
				updateItem := bson.M{
					"$set": qmap.QM{
						"module_type_id": moduleTypeId,
						"select_struct":  selectStruct,
					},
				}

				if err := mongo.NewMgoSession(common.MC_EVALUATE_ASSET).Update(selector, updateItem); err != nil {
					return false, err
				}
			}
		}
	}
	return true, nil
}

func (this *EvaluateItem) ChangeTestStatus(itemIds []interface{}, status int) error {
	selector := qmap.QM{
		"_id": bson.M{"$in": itemIds},
	}
	update := bson.M{
		"$set": bson.M{
			"test_status": status,
		},
	}
	if _, err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).UpdateAll(selector, update); err != nil {
		return err
	}
	return nil
}

// 删除项目漏洞，然后又创建新的项目漏洞。这操作，对于漏洞管理模块的功能不太好支持，需要进行重写，此代码保留，新建一个函数去写逻辑 UpdateVuls2
func (this *EvaluateItem) UpdateVuls(taskId string, itemIds []string) error {
	// 删除ITEM表中已存在的漏洞
	match := bson.M{
		"item_id": bson.M{"$in": itemIds},
	}
	if _, err := mongo.NewMgoSession(common.MC_EVALUATE_VULNERABILITY).RemoveAll(match); err == nil {
		for _, itemId := range itemIds {
			// 查询测试用例副本的漏洞
			vuls, err := new(EvaluateTaskVulnerability).GetItemTaskVuls(taskId, itemId)
			if err == nil {
				for _, vul := range vuls {
					data := vul
					data["test_id"] = vul["record_id"]
					data["_id"] = vul["vul_id"]
					if _, err := new(EvaluateVulnerability).Create(data, int64(vul["op_id"].(int))); err != nil {
						return err
					}
				}
			}
			// 更新ITEM表中，漏洞统计
			params := qmap.QM{
				"e_item_id": itemId,
			}
			vulList, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VULNERABILITY, params).Get()

			selector := bson.M{"_id": itemId}
			updateItem := bson.M{
				"$set": qmap.QM{
					"vul_number": len(*vulList),
				},
			}

			if err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).Update(selector, updateItem); err != nil {
				return err
			}
		}
	}

	return nil
}

// 将任务的漏洞，更新到ITEM的漏洞主表中
func (this *EvaluateItem) UpdateVuls2(taskId string, itemIds []string, opId int) error {
	for _, itemId := range itemIds {
		// 查询测试用例副本的漏洞
		vuls, err := new(EvaluateTaskVulnerability).GetItemTaskVuls(taskId, itemId)
		if err == nil {
			for _, vul := range vuls {
				// 拿到一条任务漏洞，根据其evaluate_vulnerability_id字段，看项目漏洞中有没有这条记录，有就更新
				taskVulId := vul["_id"].(bson.ObjectId).Hex()
				evaluateVulnerabilityId := vul["evaluate_vulnerability_id"].(string)
				if ev, err := new(EvaluateVulnerability).One(evaluateVulnerabilityId); err != nil {
					if err == mgo.ErrNotFound {
						// 数据没有查询到，新增主漏洞表记录
						evaluateVulnerabilityId := vul["evaluate_vulnerability_id"].(string)
						data := vul
						data["test_id"] = vul["record_id"]
						data["_id"] = evaluateVulnerabilityId
						if _, err := new(EvaluateVulnerability).Create(data, int64(vul["op_id"].(int))); err != nil {
							return err
						} else {
							// 创建日志记录 EvaluateVulnerabilityLog
							new(EvaluateVulnerabilityLog).AddVulTimeLine(taskVulId, opId)
						}
					} else {
						return err
					}
				} else {
					ev.Update(evaluateVulnerabilityId, vul, opId)
					// 创建日志记录 EvaluateVulnerabilityLog
					new(EvaluateVulnerabilityLog).AddVulTimeLine(taskVulId, opId)
				}
			}
		}
		// 更新ITEM表中，漏洞统计
		params := qmap.QM{
			"e_item_id": itemId,
		}
		vulList, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VULNERABILITY, params).Get()

		selector := bson.M{"_id": itemId}
		updateItem := bson.M{
			"$set": qmap.QM{
				"vul_number": len(*vulList),
			},
		}

		if err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).Update(selector, updateItem); err != nil {
			return err
		}
	}
	return nil
}

// 统计项目测试用例数量
func (this *EvaluateItem) CountItem(projectIds []string) int {
	params := qmap.QM{
		"in_project_id": projectIds,
	}
	if count, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params).Count(); err == nil {
		return count
	}
	return 0
}

// 统计项目测试用例测试情况
func (this *EvaluateItem) StatisticItemInfo(projectId string) []qmap.QM {
	match := bson.M{
		"project_id": bson.M{"$eq": projectId},
	}
	group := bson.M{
		"_id":         bson.M{"t_status": "$test_status"},
		"test_status": bson.M{"$first": "$test_status"},
		"number":      bson.M{"$sum": 1},
	}
	operations := []bson.M{
		{"$match": match},
		{"$group": group},
	}
	ready := 0
	testing := 0
	complete := 0
	resp := []bson.M{}
	if err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).Session.Pipe(operations).All(&resp); err == nil {
		for _, item := range resp {
			var itemQM qmap.QM = map[string]interface{}(item)
			switch itemQM.Int("test_status") {
			case common.TIS_READY:
				ready += itemQM.Int("number")
			case common.TIS_TEST_COMPLETE, common.TIS_PART_TEST_COMPLETE:
				testing += itemQM.Int("number")
			case common.TIS_COMPLETE:
				complete += itemQM.Int("number")
			}
		}
	} else {
		panic(err)
	}

	result := []qmap.QM{
		{
			"name":   "已测试",
			"number": complete,
		},
		{
			"name":   "测试中",
			"number": testing,
		},
		{
			"name":   "未测试",
			"number": ready,
		},
	}
	return result
}

// ----------------------------------------
// ----------------------------------------
// ----------------------------------------

// 将预绑定的正式绑定
func EvaluateItem_bindEvaluateTask(ctx context.Context, ids []string) error {
	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_ITEM)
	selector := bson.M{"_id": bson.M{"$in": ids}}
	update := qmap.QM{"$set": qmap.QM{"is_pre_bind": common.NOT_PREBIND}}
	_, err := coll.UpdateMany(ctx, selector, update)
	return err
}

// 将某一批测试用例与项目任务解绑
func EvaluateItem_UnbindEvaluateTask(ctx context.Context, itemIds []string) error {
	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_ITEM)
	selector := bson.M{"_id": bson.M{"$in": itemIds}}
	set := qmap.QM{
		"status":        common.EIS_FREE, // 将测试用例状态，变更为可绑定
		"is_pre_delete": common.NOT_PREDEL,
		"pre_bind":      "",
	}
	update := qmap.QM{"$set": set}
	_, err := coll.UpdateMany(ctx, selector, update)
	return err
}

// 将某一批测试用例与项目任务绑定
// status设置为1, 表示使用中
// pre_bind设置为任务的id
func EvaluateItemBindEvaluateTask(ctx context.Context, itemIds []string, evaluateTaskId string, isPreBind int) error {
	// 第一步，将测试用例status设置为1，表示测试用例已被占用
	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_ITEM)
	selector := bson.M{"_id": bson.M{"$in": itemIds}} // MC_EVALUATE_ITEM的_id 是Hex格式
	data := qmap.QM{
		"status":      common.EIS_INUSE,
		"pre_bind":    evaluateTaskId,
		"is_pre_bind": isPreBind,
		"test_status": common.TIS_READY,
		"test_time":   0,
	}

	_, err := coll.UpdateMany(ctx, selector, qmap.QM{"$set": data})
	return err
}
