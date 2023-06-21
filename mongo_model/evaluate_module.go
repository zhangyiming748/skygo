package mongo_model

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"

	"github.com/globalsign/mgo/bson"
)

type EvaluateModule struct {
	Id             bson.ObjectId `bson:"_id,omitempty" json:"id"`
	ModuleName     string        `bson:"module_name" json:"module_name"`           // 测试组件
	ModuleNameCode string        `bson:"module_name_code" json:"module_name_code"` // 测试组件编码
	ModuleType     string        `bson:"module_type" json:"module_type"`           // 测试分类
	ModuleTypeCode string        `bson:"module_type_code" json:"module_type_code"` // 测试组件编码
}

func (this *EvaluateModule) Create(rawInfo qmap.QM) (*EvaluateModule, error) {
	this.ModuleName = rawInfo.MustString("module_name")
	this.ModuleType = rawInfo.MustString("module_type")
	this.ModuleTypeCode = rawInfo.MustString("module_type_code")
	this.Id = bson.NewObjectId()

	moduleNameCode, has := rawInfo.TryString("module_name_code")
	if has == false {
		return nil, errors.New("module_name_code 不能为空")
	}

	this.ModuleNameCode = moduleNameCode

	if err := this.CheckModuleNameCode(this.ModuleName, this.ModuleNameCode, this.Id.Hex()); err != nil {
		return nil, err
	}

	if err := this.CheckModuleTypeCode(this.ModuleType, this.ModuleTypeCode, this.Id.Hex()); err != nil {
		return nil, err
	}

	if _, err := new(EvaluateModule).Find("", this.ModuleName, this.ModuleType); err == nil {
		return nil, errors.New(fmt.Sprintf("测试组件:%s 测试分类:%s 已经存在", this.ModuleName, this.ModuleType))
	}
	if err := mongo.NewMgoSession(common.MC_EVALUATE_MODULE).Insert(this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateModule) Update(id string, rawInfo qmap.QM) (*EvaluateModule, error) {
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(id),
	}
	err := sys_service.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, params).One(this)
	if err != nil {
		return nil, err
	}
	if moduleName, has := rawInfo.TryString("module_name"); has {
		this.ModuleName = moduleName
	}
	if moduleNameCode, has := rawInfo.TryString("module_name_code"); has {
		this.ModuleNameCode = moduleNameCode
	}
	if moduleType, has := rawInfo.TryString("module_type"); has {
		this.ModuleType = moduleType
	}
	if moduleTypeCode, has := rawInfo.TryString("module_type_code"); has {
		this.ModuleTypeCode = moduleTypeCode
	}
	if err := this.CheckModuleNameCode(this.ModuleName, this.ModuleNameCode, this.Id.Hex()); err != nil {
		return nil, err
	}

	if err := this.CheckModuleTypeCode(this.ModuleType, this.ModuleTypeCode, this.Id.Hex()); err != nil {
		return nil, err
	}
	if _, err := new(EvaluateModule).Find(this.Id.Hex(), this.ModuleName, this.ModuleType); err == nil {
		return nil, errors.New(fmt.Sprintf("更新失败，测试组件:%s 测试分类:%s 已经存在", this.ModuleName, this.ModuleType))
	}

	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_MODULE)
	if _, err := coll.UpdateOne(context.Background(), bson.M{"_id": this.Id}, bson.M{"$set": this}); err != nil {
		return nil, err
	} else {
		return this, nil
	}
}

func (this *EvaluateModule) Find(excludeId, moduleName, moduleType string) (*EvaluateModule, error) {
	params := qmap.QM{
		"e_module_name": moduleName,
		"e_module_type": moduleType,
	}
	if excludeId != "" {
		params["ne__id"] = bson.ObjectIdHex(excludeId)
	}
	err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_MODULE, params).One(this)
	return this, err
}

func (this *EvaluateModule) GetModuleByTypeIds(typeIds []interface{}) ([]map[string]interface{}, error) {
	idObject := []primitive.ObjectID{}
	for _, id := range typeIds {
		_id, _ := primitive.ObjectIDFromHex(id.(string))
		idObject = append(idObject, _id)
	}
	queryParams := qmap.QM{
		"in__id": idObject,
	}
	all, err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_MODULE, queryParams).Find()
	return all, err
}

func (this *EvaluateModule) GetModuleSlice() (*qmap.QM, error) {
	moduleList, err := mongo.NewMgoSession(common.MC_EVALUATE_MODULE).SetLimit(5000).Get()
	if err != nil {
		return nil, err
	}
	moduleMap := qmap.QM{}
	for _, item := range *moduleList {
		id := item["_id"]
		itType := reflect.TypeOf(id)
		switch itType.Name() {
		case "bson.ObjectId":
			moduleMap[item["_id"].(bson.ObjectId).Hex()] = item["module_name"]
		default:
			continue
		}
	}
	return &moduleMap, nil
}

func (this *EvaluateModule) FindById(id string) (*EvaluateModule, error) {
	_id, _ := primitive.ObjectIDFromHex(id)
	params := qmap.QM{
		"e__id": _id,
	}
	if err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_MODULE, params).One(this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateModule) Update1() {
	mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_MODULE)
	mgoSession.SetLimit(100)
	res, _ := mgoSession.Get()
	for _, module := range *res {
		moduleName := module["module_name"]
		moduleType := module["module_type"]
		update := bson.M{
			"$set": bson.M{
				"module_name_code": common.Module.DefaultString(moduleName.(string), "999"),
				"module_type_code": common.ModuleType.DefaultString(moduleType.(string), "999"),
			},
		}
		mongoClient := mongo.NewMgoSession(common.MC_EVALUATE_MODULE)
		if err := mongoClient.Update(bson.M{"_id": module["_id"]}, &update); err == nil {
		} else {
			panic(err)
		}
	}
}

// 校验测试组件码是否符合条件
// 如果测试组件名称已经存在，则添加的测试组件编码必须与库中已有测试组件编码一致
// 如果测试组件编号已经存在，则添加的测试组件名称必须与库中已有测试组件名称一致
func (this *EvaluateModule) CheckModuleNameCode(moduleName, moduleNameCode, id string) error {
	if existCode := new(EvaluateModule).GetModuleNameCodeByName(moduleName, id); existCode != "" && existCode != moduleNameCode {
		return errors.New(fmt.Sprintf("测试组件:%s 编码与现有测试组件编码不一致", moduleName))
	}
	if existName := new(EvaluateModule).GetModuleNameByCode(moduleNameCode); existName != "" && existName != moduleName {
		return errors.New(fmt.Sprintf("测试组件编码:%s 已经存在", moduleNameCode))
	}
	return nil
}

func (this *EvaluateModule) GetModuleNameCodeByName(moduleName, excludeId string) string {
	params := qmap.QM{
		"e_module_name": moduleName,
	}
	if excludeId != "" {
		params["ne__id"] = bson.ObjectIdHex(excludeId)
	}
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, params).One(this); err == nil {
		return this.ModuleNameCode
	}
	return ""
}

func (this *EvaluateModule) GetModuleNameByCode(moduleNameCode string) string {
	params := qmap.QM{
		"e_module_name_code": moduleNameCode,
	}
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, params).One(this); err == nil {
		return this.ModuleName
	}
	return ""
}

// 校验测试分类编码是否符合条件
// 如果测试分类已经存在，则添加的测试分类编码必须与库中已有测试分类编码一致
// 如果测试分类编号已经存在，则添加的测试分类名称必须与库中已有测试分类名称一致
func (this *EvaluateModule) CheckModuleTypeCode(moduleType, moduleTypeCode, id string) error {
	existCode := new(EvaluateModule).GetModuleTypeCode(moduleType, id)
	if existCode != "" && existCode != moduleTypeCode {
		return errors.New(fmt.Sprintf("测试分类:%s 编码与现有测试分类编码不一致", moduleType))
	}
	if existName := new(EvaluateModule).GetModuleTypeName(moduleTypeCode, id); existName != "" && existName != moduleType {
		return errors.New(fmt.Sprintf("测试分类编码:%s 已经存在", moduleTypeCode))
	}
	return nil
}

func (this *EvaluateModule) GetModuleTypeCode(moduleType, id string) string {
	params := qmap.QM{
		"e_module_type": moduleType,
	}
	if id != "" {
		params["ne__id"] = bson.ObjectIdHex(id)
	}
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, params).One(this); err == nil {
		return this.ModuleTypeCode
	}
	return ""
}

func (this *EvaluateModule) GetModuleTypeName(moduleTypeCode, id string) string {
	params := qmap.QM{
		"e_module_type_code": moduleTypeCode,
	}
	if id != "" {
		params["ne__id"] = bson.ObjectIdHex(id)
	}
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, params).One(this); err == nil {
		return this.ModuleType
	}
	return ""
}

func (this *EvaluateModule) GetRecommendModuleNameCode() string {
	maxModuleNameCode := new(EvaluateModule).GetMaxModuleNameCode()
	maxCode := 1
	if temp, err := strconv.Atoi(maxModuleNameCode); err == nil {
		maxCode = temp
	}
	maxCode += 10
	for ; maxCode < 999; maxCode += 10 {
		maxCodeStr := ""
		if maxCode < 100 {
			maxCodeStr = fmt.Sprintf("0%d", maxCode)
		} else {
			maxCodeStr = fmt.Sprintf("%d", maxCode)
		}
		if new(EvaluateModule).IsModuleNameCodeAvailable(maxCodeStr) {
			return maxCodeStr
		}
	}
	for maxCode = 998; maxCode > 0; maxCode-- {
		maxCodeStr := ""
		if maxCode < 100 {
			maxCodeStr = fmt.Sprintf("0%d", maxCode)
		} else {
			maxCodeStr = fmt.Sprintf("%d", maxCode)
		}
		if new(EvaluateModule).IsModuleNameCodeAvailable(maxCodeStr) {
			return maxCodeStr
		}
	}
	return "000"
}

// 查询到最大的module_type_code
// 如果其小于999，就判断是否存在，不存在就累加10，正在判断是否存在，找到一个不存在的code返回
func (this *EvaluateModule) GetRecommendModuleTypeCode() string {
	maxModuleTypeCode := new(EvaluateModule).GetMaxModuleTypeCode()
	maxCode := 1
	if temp, err := strconv.Atoi(maxModuleTypeCode); err == nil {
		maxCode = temp
	}
	maxCode += 10
	for ; maxCode < 999; maxCode += 10 {
		maxCodeStr := ""
		if maxCode < 100 {
			maxCodeStr = fmt.Sprintf("0%d", maxCode)
		} else {
			maxCodeStr = fmt.Sprintf("%d", maxCode)
		}
		if new(EvaluateModule).IsModuleTypeCodeAvailable(maxCodeStr) {
			return maxCodeStr
		}
	}
	for maxCode = 998; maxCode > 0; maxCode-- {
		maxCodeStr := ""
		if maxCode < 100 {
			maxCodeStr = fmt.Sprintf("0%d", maxCode)
		} else {
			maxCodeStr = fmt.Sprintf("%d", maxCode)
		}
		if new(EvaluateModule).IsModuleNameCodeAvailable(maxCodeStr) {
			return maxCodeStr
		}
	}
	return "000"
}

// 从集合 evaluate_module 中按照module_name_code倒序查询，查出最大值
// 要求小于999，如果查询不到或者出问题，就默认998
func (this *EvaluateModule) GetMaxModuleNameCode() string {
	param := qmap.QM{
		"lt_module_name_code": "999",
	}
	mgoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, param).SetLimit(1)
	mgoClient.AddSorter("module_name_code", 1)
	if list, err := mgoClient.Get(); err == nil && len(*list) > 0 {
		var itemQM qmap.QM = (*list)[0]
		return itemQM.String("module_name_code")
	} else {
		return "998"
	}
}

// 从集合 evaluate_module 中按照module_type_code倒序查询，查出最大值
// 要求小于999，如果查询不到或者出问题，就默认998
func (this *EvaluateModule) GetMaxModuleTypeCode() string {
	param := qmap.QM{
		"lt_module_type_code": "999",
	}
	mgoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, param).SetLimit(1)
	mgoClient.AddSorter("module_type_code", 1)
	if list, err := mgoClient.Get(); err == nil && len(*list) > 0 {
		var itemQM qmap.QM = (*list)[0]
		return itemQM.String("module_type_code")
	} else {
		return "998"
	}
}

// 没查到一样的值就是false
func (this *EvaluateModule) IsModuleNameCodeAvailable(moduleNameCode string) bool {
	params := qmap.QM{
		"e_module_name_code": moduleNameCode,
	}
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, params).One(this); err != nil {
		return true
	} else {
		return false
	}
}

// 没查到一样的值就是false
func (this *EvaluateModule) IsModuleTypeCodeAvailable(moduleTypeCode string) bool {
	params := qmap.QM{
		"e_module_type_code": moduleTypeCode,
	}
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, params).One(this); err != nil {
		return true
	} else {
		return false
	}
}

func (this *EvaluateModule) GetBydName(moduleName, moduleType string) (*EvaluateModule, error) {
	params := qmap.QM{
		"e_module_name": moduleName,
		"e_module_type": moduleType,
	}
	err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, params).One(this)
	return this, err
}
