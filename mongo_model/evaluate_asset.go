package mongo_model

import (
	"context"
	"errors"
	"fmt"
	"strconv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"
)

type EvaluateAsset struct {
	Id           string                 `bson:"_id"`            // 资产id
	ProjectId    string                 `bson:"project_id"`     // 项目id
	Name         string                 `bson:"name"`           // 资产名称
	Version      string                 `bson:"version"`        // 资产版本
	Image        string                 `bson:"image"`          // 资产示意图
	EvaluateType string                 `bson:"evaluate_type"`  // 资产类型
	ModuleTypeId []interface{}          `bson:"module_type_id"` // 资产类型
	Attributes   map[string]interface{} `bson:"attributes"`     // 测试对象属性
	OpId         int64                  `bson:"op_id"`          // 操作人id
	Note         string                 `bson:"note"`           // 备注
	SelectStruct interface{}            `bson:"select_struct"`  // 前端筛选结构保存
	CreateTime   int64                  `bson:"create_time"`    // 创建时间
	UpdateTime   int64                  `bson:"update_time"`    // 更新时间
}

func (this *EvaluateAsset) Create(ctx context.Context, projectId string, opId int64, rawInfo qmap.QM) (*EvaluateAsset, error) {
	// 判断是否有重复的资产
	this.Name = rawInfo.MustString("name")
	if this.checkNameExist(projectId, this.Name) {
		return this, errors.New("该资产名称已被使用")
	}

	if module, err := new(EvaluateModule).GetModuleByTypeIds(rawInfo["module_type_id"].([]interface{})); err != nil || len(module) <= 0 {
		return this, errors.New("该组件分类不存在")
	}

	typeName := rawInfo.MustString("evaluate_type")
	evaluateType := new(EvaluateType)
	if _, err := evaluateType.GetOne("", typeName); err != nil {
		return nil, errors.New("Unknown evaluate type!")
	}
	microSecond := custom_util.GetCurrentMilliSecond()
	this.Id = this.GenerateAssetID(int(microSecond))
	this.ProjectId = projectId
	this.ModuleTypeId = rawInfo["module_type_id"].([]interface{})
	this.Version = rawInfo.MustString("version")
	this.EvaluateType = typeName
	this.OpId = opId
	if note, has := rawInfo.TryString("note"); has {
		this.Note = note
	}
	if image, has := rawInfo.TryString("image"); has {
		this.Image = image
	}
	if selectStruct, has := rawInfo.TryInterface("select_struct"); has {
		this.SelectStruct = selectStruct
	}
	this.UpdateTime = microSecond
	this.CreateTime = microSecond
	if attributes, err := evaluateType.ExtraAttributeMap(rawInfo); err == nil {
		this.Attributes = attributes
	} else {
		return nil, err
	}
	// if err := mongo.NewMgoSession(common.MC_EVALUATE_ASSET).Insert(this); err == nil {
	if _, err := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_ASSET).InsertOne(ctx, this); err == nil {
		// 创建资产后，将测试用例库导入项目用例
		this.CopyTestCaseToItem(ctx, this.ProjectId, this.Id, qmap.QM{"module_type_id": this.ModuleTypeId}.SliceString("module_type_id"), opId)
		return this, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateAsset) CopyTestCaseToItem(ctx context.Context, projectId, assetId string, moduleTypeIds []string, opId int64) error {
	params := qmap.QM{
		"in_module_type_id": moduleTypeIds,
	}

	widget := orm_mongo.NewWidgetWithCollectionName(common.MC_TEST_CASE).SetParams(params)
	widget.SetLimit(5000)
	testCaseList, _ := widget.Find()
	for _, item := range testCaseList {
		// 检查测试用例是否存在
		checkParams := qmap.QM{
			"e__id":        assetId + item["_id"].(string),
			"e_project_id": projectId,
		}
		// 检查当前测试用例是否存在，存在则跳过
		widget := orm_mongo.NewWidgetWithCollectionName(common.MC_EVALUATE_ITEM).SetParams(checkParams)
		result, err := widget.Get()
		if err == nil && len(result) > 0 {
			continue
		}

		level := 1
		itemQM := qmap.QM(item)
		if itemQM["level"] != "" {
			if lv, has := itemQM.TryInt("level"); has {
				level = lv
			}
		}

		itemData := qmap.QM{
			"id":              assetId + item["_id"].(string),
			"project_id":      projectId,
			"name":            item["name"],
			"asset_id":        assetId,
			"module_type_id":  item["module_type_id"],
			"level":           level,
			"objective":       item["objective"],
			"external_input":  item["external_input"],
			"test_procedure":  item["test_procedure"],
			"test_standard":   item["test_standard"],
			"test_case_level": item["test_case_level"],
			"test_method":     item["test_method"],
			"auto_test_level": item["auto_test_level"],
			"test_script":     "",
			"test_sketch_map": "",
		}
		if _, err := new(EvaluateItem).Create(ctx, itemData, opId); err != nil {
			return err
		}
	}
	return nil
}

// 更新 测试对象的同时，需要更新测试项里的测试对象名称
func (this *EvaluateAsset) Update(ctx context.Context, id string, opId int64, rawInfo qmap.QM) (*EvaluateAsset, error) {
	params := qmap.QM{
		"e__id": id,
	}
	w := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_ASSET, params)
	if err := w.One(&this); err == nil {
		if version, has := rawInfo.TryString("version"); has && version != this.Version {
			// 检查是否变更版本号，如果变更，将历史数据存入版本表
			new(EvaluateAssetVersion).Create(this)
			this.Version = version
		}

		if name, has := rawInfo.TryString("name"); has {
			this.Name = name
		}

		if image, has := rawInfo.TryString("image"); has {
			this.Image = image
		}
		if note, has := rawInfo.TryString("note"); has {
			this.Note = note
		}
		if evaluateType, has := rawInfo.TryString("evaluate_type"); has {
			this.EvaluateType = evaluateType
		}
		if selectStruct, has := rawInfo.TryInterface("select_struct"); has {
			this.SelectStruct = selectStruct
		}

		delSlice := []string{}
		addSlice := []string{}
		if ModuleTypeId, has := rawInfo.TrySlice("module_type_id"); has {
			// 判断新的组件分类和原来组件分类区别， 删除掉的判断测试用例是否已被使用。 新增的复制测试用例
			tmp := qmap.QM{
				"module_type_id": this.ModuleTypeId,
			}
			IdSlice := rawInfo.SliceString("module_type_id")
			delSlice = custom_util.DifferenceDel(tmp.SliceString("module_type_id"), IdSlice)
			addSlice = custom_util.DifferenceDel(IdSlice, tmp.SliceString("module_type_id"))

			// 判断删除掉的判断测试用例是否已被使用
			if _, err := this.CheckModuleTypeIdCanBeDelete(this.Id, delSlice); err != nil {
				return nil, err
			}
			this.ModuleTypeId = ModuleTypeId

		}

		evaluateType := new(EvaluateType)
		if _, err := evaluateType.GetOne("", this.EvaluateType); err != nil {
			return nil, errors.New("Unkown evaluate type!")
		}
		var attributes qmap.QM = this.Attributes
		if attributes, err := evaluateType.ExtraAttributeMap(attributes.Merge(rawInfo)); err == nil {
			this.Attributes = attributes

			coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_ASSET)
			if _, err := coll.UpdateOne(ctx, bson.M{"_id": this.Id}, qmap.QM{"$set": this}); err == nil {
				// 因组件分类删除，删除对应的测试用例
				if len(delSlice) > 0 {
					if _, err := this.DeleteItemByModuleTypeIds(this.Id, delSlice); err != nil {
						return nil, err
					}
				}

				// 因组件分类增加，复制对应的测试用例
				this.CopyTestCaseToItem(ctx, this.ProjectId, this.Id, addSlice, opId)

				return this, nil
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (this *EvaluateAsset) CheckModuleTypeIdCanBeDelete(assetId string, moduleTypeIds []string) (bool, error) {
	// 判断相关的测试用例是否已经绑定任务，如果已经绑定则不能删除
	for _, moduleTypeId := range moduleTypeIds {
		params := qmap.QM{
			"e_asset_id":       assetId,
			"e_module_type_id": moduleTypeId,
			"ne_pre_bind":      "",
		}

		if _, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params).GetOne(); err == nil {
			return false, errors.New("当前选择的资产的测试用例已被分配任务，不支持更新")
		}
	}
	return true, nil
}

func (this *EvaluateAsset) GetOne(id string) (qmap.QM, error) {
	params := qmap.QM{
		"e__id": id,
	}
	info, err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_ASSET, params).Get()
	if err == nil {
		// 根据实际测试用例，更新资产的组件和分类
		items, err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_ITEM, qmap.QM{"e_asset_id": id}).SetLimit(100000).Find()
		moduleTypeIds := []string{}
		for _, item := range items {
			if !custom_util.IndexOfSlice(item["module_type_id"].(string), moduleTypeIds) {
				moduleTypeIds = append(moduleTypeIds, item["module_type_id"].(string))
			}
		}

		idObject := []primitive.ObjectID{}
		for _, id := range moduleTypeIds {
			_id, _ := primitive.ObjectIDFromHex(id)
			idObject = append(idObject, _id)
		}
		moduleType, err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_MODULE, qmap.QM{"in__id": idObject}).SetLimit(100000).Find()

		selectStruct := [][]interface{}{}
		for _, selectItem := range (info)["select_struct"].(primitive.A) {
			item := selectItem.(primitive.A)
			if custom_util.IndexOfSlice(item[1].(string), moduleTypeIds) {
				selectStruct = append(selectStruct, selectItem.(primitive.A))
			}
		}
		(info)["module_type_id"] = moduleTypeIds
		(info)["select_struct"] = selectStruct
		(info)["module_type"] = moduleType

		selector := bson.M{
			"_id": id,
		}
		updateItem := bson.M{
			"$set": qmap.QM{
				"module_type_id": moduleTypeIds,
				"select_struct":  selectStruct,
			},
		}
		coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_ASSET)
		if _, err := coll.UpdateMany(context.Background(), selector, updateItem); err != nil {
			fmt.Println(err)
			return nil, err
		}

		// 追加历史版本
		// 		version, err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_ASSET_VERSION, qmap.QM{"e_asset_id": id}).Get()
		// 		(info)["history"] = version
		// fmt.Println(version,err,  "0000000000000000000")
		version, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET_VERSION, qmap.QM{"e_asset_id": id}).Get()
		(info)["history"] = version

		return info, err
	}
	return info, err
}

func (this *EvaluateAsset) DeleteItemByModuleTypeIds(assetId string, moduleTypeIds []string) (bool, error) {
	// 删除对应的测试用例
	for _, moduleTypeId := range moduleTypeIds {
		match := bson.M{
			"asset_id":         assetId,
			"module_type_id":   moduleTypeId,
			"evaluate_task_id": "",
		}
		if _, err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).RemoveAll(match); err != nil {
			return false, err
		}
	}
	return true, nil
}

func (this *EvaluateAsset) GetTaskAssetInfo(AssetId, taskId string) (qmap.QM, error) {
	params := qmap.QM{
		"e__id": AssetId,
	}
	info, err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_ASSET, params).Get()
	if err == nil {
		_taskId, _ := primitive.ObjectIDFromHex(taskId)
		if taskInfo, err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_TASK, qmap.QM{"e__id": _taskId}).Get(); err == nil {
			// 查出任务的资产版本
			for k, v := range (taskInfo)["asset_versions"].(map[string]interface{}) {
				if k == AssetId {
					(info)["version"] = v
				}
			}

			// 查出任务对应的测试用例，根据测试用例查出组件和分类
			items, err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_ITEM, qmap.QM{"in__id": (taskInfo)["evaluate_item_ids"]}).Find()
			moduleTypeIds := []string{}
			if err == nil {
				for _, item := range items {
					if !custom_util.IndexOfSlice(item["module_type_id"].(string), moduleTypeIds) {
						moduleTypeIds = append(moduleTypeIds, item["module_type_id"].(string))
					}
				}
			}
			idObject := []primitive.ObjectID{}
			(info)["module_type_id"] = moduleTypeIds
			for _, id := range moduleTypeIds {
				_id, _ := primitive.ObjectIDFromHex(id)
				idObject = append(idObject, _id)
			}
			moduleType, err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_MODULE, qmap.QM{"in__id": idObject}).Find()
			(info)["module_type"] = moduleType

		}
		// 追加历史版本
		version, err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_ASSET_VERSION, qmap.QM{"e_asset_id": AssetId}).Find()
		(info)["history"] = version

		return info, err
	}
	return info, err
}

func (this *EvaluateAsset) One(id string) (*EvaluateAsset, error) {
	params := qmap.QM{
		"e__id": id,
	}
	err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_ASSET, params).One(this)
	return this, err
}

func (this *EvaluateAsset) TypeAsset(projectId string) (*qmap.QM, error) {
	params := qmap.QM{
		"e_project_id": projectId,
	}
	list, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, params).SetLimit(5000).Get()
	result := []qmap.QM{}
	typeQM := qmap.QM{}
	if err == nil {
		for _, item := range *list {
			if typeQM[item["evaluate_type"].(string)] == nil {
				typeQM[item["evaluate_type"].(string)] = []qmap.QM{}
			}
			assetStruct := qmap.QM{
				"id":   item["_id"],
				"name": item["name"].(string),
			}
			typeQM[item["evaluate_type"].(string)] = append(typeQM[item["evaluate_type"].(string)].([]qmap.QM), assetStruct)
		}
		for typeName, item := range typeQM {
			itemSlice := qmap.QM{
				"type": typeName,
				"list": item,
			}
			result = append(result, itemSlice)
		}

	}
	return &qmap.QM{"data": result}, err
}

func (this *EvaluateAsset) DeleteModuleType(id string, rawInfo qmap.QM) (*EvaluateAsset, error) {
	params := qmap.QM{
		"e__id": id,
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, params)
	if err := mongoClient.One(&this); err == nil {
		if moduleTypeId, has := rawInfo.TrySlice("module_type_id"); has {
			this.ModuleTypeId = moduleTypeId
		}
		if err := mongoClient.Update(bson.M{"_id": this.Id}, this); err == nil {
			// 更新测试项中的测试对象名称target_name

			return this, nil
		} else {
			return nil, err
		}

	} else {
		return nil, err
	}
}

func (this *EvaluateAsset) checkNameExist(projectId, name string) bool {
	params := qmap.QM{
		"e_project_id": projectId,
		"e_name":       name,
	}
	if result, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, params).GetOne(); err == nil {
		if len(*result) != 0 {
			// 存在重复的数据
			return true
		}
	}
	// 不存在重复的数据
	return false
}

func (this *EvaluateAsset) GetAssetSlice(projectId string) (qmap.QM, error) {
	params := qmap.QM{
		"e_project_id": projectId,
	}
	result := qmap.QM{}
	client := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, params)
	client.AddSorter("create_time", 1)

	assetList, err := client.SetLimit(5000).Get()
	if err == nil {
		for _, item := range *assetList {
			result[item["_id"].(string)] = qmap.QM{
				"name":        item["name"],
				"create_time": item["create_time"],
			}
		}
		return result, err
	}
	return nil, err
}

func (this *EvaluateAsset) GenerateAssetID(time int) string {
	str := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	assetID := ""
	var remainder int
	var remainderStr string
	for time != 0 {
		remainder = time % 36
		if remainder < 36 && remainder > 9 {
			remainderStr = str[remainder]
		} else {
			remainderStr = strconv.Itoa(remainder)
		}
		assetID = remainderStr + assetID
		time = time / 36
	}
	if len(assetID) > 8 {
		rs := []rune(assetID)
		assetID = string(rs[:8])
	}

	return assetID
}

func (this *EvaluateAsset) GetAssetName(id string) (name string) {
	params := qmap.QM{
		"e__id": id,
	}
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, params).One(this); err == nil {
		return this.Name
	}
	return
}

/*
	查询资产关联的所有版本信息
	{
			"id": "US4VX1KR",
			"name": "test",
			"versions": [
					"1.0"
			]
	}
*/

func (this *EvaluateAsset) GetAssetAllVersions(id string) (qmap.QM, error) {
	params := qmap.QM{
		"e__id": id,
	}
	// 查询当前资产版本
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, params).One(this); err == nil {
		assetInfo := qmap.QM{
			"id":   id,
			"name": this.Name,
		}
		versions := []string{this.Version}
		// 查询历史资产版本
		params := qmap.QM{
			"e_asset_id": id,
		}
		if list, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET_VERSION, params).Get(); err == nil {
			for _, item := range *list {
				var itemQM qmap.QM = item
				versions = append(versions, itemQM.String("version"))
			}
		}
		assetInfo["versions"] = versions
		return assetInfo, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateAsset) FindById(id string) (*EvaluateAsset, error) {
	params := qmap.QM{
		"e__id": id,
	}
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, params).One(this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}
