package mongo_model

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"

	"github.com/globalsign/mgo/bson"
)

// 项目任务绑定测试用例测试表
type EvaluateTaskItem struct {
	Id                bson.ObjectId `bson:"_id,omitempty"`                                  // 测试id
	ItemId            string        `bson:"item_id" json:"item_id"`                         // 测试用例ID
	Name              string        `bson:"name" json:"name"`                               // 测试用例名称
	ProjectId         string        `bson:"project_id" json:"project_id"`                   // 项目id
	EvaluateTaskId    string        `bson:"evaluate_task_id" json:"evaluate_task_id"`       // 关联的项目任务id
	AssetId           string        `bson:"asset_id" json:"asset_id"`                       // 资产id
	ModuleTypeId      string        `bson:"module_type_id" json:"module_type_id"`           // 组件分类id
	TestMethod        string        `bson:"test_method" json:"test_method"`                 // 测试方法（黑盒、白盒）
	AutoTestLevel     string        `bson:"auto_test_level" json:"auto_test_level"`         // 自动化测试程度（自动化、人工）
	TestTime          int           `bson:"test_time" json:"test_time"`                     // 测试时间
	TestCount         int           `bson:"test_count" json:"test_count"`                   // 测试次数
	Status            int           `bson:"status" json:"status"`                           // 测试状态 （0:待测试 1:测试完成,2待补充，3审核通过）
	RecordAuditStatus int           `bson:"record_audit_status" json:"record_audit_status"` // 测试记录审核状态 （1:通过 0:待审核 -1:驳回）
	TestPhase         int           `bson:"test_phase" json:"test_phase"`                   // 测试阶段 （1:初测、2:复测1、3:复测2、4:复测3 ...）
	RecordId          string        `bson:"record_id" json:"record_id"`                     // 当前测试记录id(测试用例每进入一个新的项目任务，就会重新预创建一个当前测试记录id)
	OpId              int           `bson:"op_id" json:"op_id"`                             // 操作人id
	UpdateTime        int64         `bson:"update_time" json:"update_time"`                 // 更新时间
	CreateTime        int64         `bson:"create_time" json:"create_time,omitempty"`       // 创建时间
}

func (this *EvaluateTaskItem) Create(rawInfo qmap.QM, opId int) (*EvaluateTaskItem, error) {
	projectId := rawInfo.MustString("project_id")
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(projectId),
	}
	if _, err := mongo.NewMgoSessionWithCond(common.MC_PROJECT, params).GetOne(); err != nil {
		panic(fmt.Sprintf("项目: %s 不存在！", projectId))
	}

	taskId := rawInfo.MustString("task_id")
	itemId := rawInfo.MustString("item_id")
	params = qmap.QM{
		"e__id": itemId,
	}
	if _, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params).GetOne(); err != nil {
		panic(fmt.Sprintf("测试用例: %s 不存在！", itemId))
	}

	if id, has := rawInfo.TryInterface("id"); has {
		this.Id = id.(bson.ObjectId)
	} else {
		this.Id = bson.NewObjectId()
	}

	this.ProjectId = projectId
	this.EvaluateTaskId = taskId
	this.ItemId = itemId
	this.Name = rawInfo.MustString("name")
	this.AssetId = rawInfo.MustString("asset_id")
	this.ModuleTypeId = rawInfo.MustString("module_type_id")
	this.TestMethod = rawInfo.MustString("test_method")
	this.TestPhase = rawInfo.MustInt("test_phase")
	this.AutoTestLevel = rawInfo.MustString("auto_test_level")
	this.Status = common.TIS_READY
	this.RecordAuditStatus = common.IRAS_DEFAULT
	this.RecordId = this.Id.Hex()
	this.OpId = int(opId)
	this.CreateTime = custom_util.GetCurrentMilliSecond()
	this.UpdateTime = custom_util.GetCurrentMilliSecond()

	if err := mongo.NewMgoSession(common.MC_EVALUATE_TASK_ITEM).Insert(this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

// 获取测试信息： 包括测试用例和各种测试状态
func (this *EvaluateTaskItem) GetInfo(id string) (*qmap.QM, error) {
	param := qmap.QM{
		"e__id": bson.ObjectIdHex(id),
	}
	err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK_ITEM, param).One(this)
	if err == nil {
		param = qmap.QM{
			"e__id": this.ItemId,
		}
		if itemInfo, err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).GetOne(); err == nil {
			(*itemInfo)["evaluate_task_id"] = this.EvaluateTaskId
			(*itemInfo)["test_status"] = this.Status
			(*itemInfo)["record_id"] = this.RecordId
			(*itemInfo)["op_id"] = this.OpId
			(*itemInfo)["record_audit_status"] = this.RecordAuditStatus
			return itemInfo, nil
		}
	}
	return nil, err
}

// 获取测试信息： 包括测试用例和各种测试状态
func (this *EvaluateTaskItem) GetFullListByTaskId(taskId string) (*[]qmap.QM, error) {
	param := qmap.QM{
		"e_task_id": taskId,
	}
	result := []qmap.QM{}
	if list, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK_ITEM, param).SetLimit(10000).Get(); err == nil {
		for _, taskItem := range *list {
			param = qmap.QM{
				"e__id": taskItem["item_id"],
			}
			if itemInfo, err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).GetOne(); err == nil {
				(*itemInfo)["evaluate_task_id"] = this.EvaluateTaskId
				(*itemInfo)["test_status"] = this.Status
				(*itemInfo)["record_id"] = this.RecordId
				(*itemInfo)["op_id"] = this.OpId
				(*itemInfo)["record_audit_status"] = this.RecordAuditStatus
				result = append(result, *itemInfo)
			}
		}
	}
	return &result, nil
}

func (this *EvaluateTaskItem) GetOne(id string) (*EvaluateTaskItem, error) {
	param := qmap.QM{
		"e__id": bson.ObjectIdHex(id),
	}
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK_ITEM, param).One(this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateTaskItem) GetListByTaskId(taskId string) (*[]map[string]interface{}, error) {
	param := qmap.QM{
		"e_evaluate_task_id": taskId,
	}
	return mongo.NewMgoSession(common.MC_EVALUATE_TASK_ITEM).AddCondition(param).Get()
}

func (this *EvaluateTaskItem) GetItemIdsByTaskId(taskId string) []string {
	itemIds := []string{}
	if list, err := this.GetListByTaskId(taskId); err == nil {
		for _, info := range *list {
			itemIds = append(itemIds, info["item_id"].(string))
		}
	}
	return itemIds
}

func (this *EvaluateTaskItem) Update(id string, rawInfo qmap.QM) (*EvaluateTaskItem, error) {
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(id),
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK_ITEM, params)
	if err := mongoClient.One(&this); err == nil {
		if val, has := rawInfo.TryInt("test_time"); has {
			this.TestTime = val
		}
		if val, has := rawInfo.TryInt("test_count"); has {
			this.TestCount = val
		}
		if val, has := rawInfo.TryInt("status"); has {
			this.Status = val
		}
		if val, has := rawInfo.TryInt("record_audit_status"); has {
			this.RecordAuditStatus = val
		}
		if val, has := rawInfo.TryInt("test_phase"); has {
			this.TestPhase = val
		}

		this.UpdateTime = custom_util.GetCurrentMilliSecond()
		if err := mongoClient.Update(bson.M{"_id": this.Id}, this); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Item not found")
	}
	return this, nil
}

func (this *EvaluateTaskItem) GetEvaluateTaskReport(evaluateTaskId string, rawInfo qmap.QM) (interface{}, error) {
	assetRecords := []*qmap.QM{}
	// 查询项目任务关联的所有资产
	assetIds := this.GetEvaluateTaskRelatedAssets(evaluateTaskId)
	for _, assetId := range assetIds {
		// 查询每一个资产关联的测试用例（用 key（测试组件|测试分类|测试分类id`进行分组）
		itemRecords := map[string]*[]qmap.QM{}
		if items, err := this.GetRelatedItems(evaluateTaskId, assetId, rawInfo); err == nil {
			for _, item := range items {
				var itemQM qmap.QM = item
				// 配合前端需要加的两个字段
				itemQM["label"] = itemQM.String("name")
				itemQM["children"] = []qmap.QM{}
				// 查询测试用例关联的测试记录
				if record, err := new(EvaluateRecord).GetOne(itemQM.String("record_id")); err == nil {
					itemQM["records"] = []interface{}{record}
				}
				// 查询测试用例关联的所有漏洞信息
				if vuls, err := new(EvaluateTaskVulnerability).GetItemTaskVuls(evaluateTaskId, itemQM.String("item_id")); err == nil {
					itemQM["vulnerabilities"] = vuls
				}
				if scanResults, err := new(ToolTaskResultBindTest).GetToolScanResult(itemQM.String("record_id")); err == nil {
					itemQM["scan_results"] = scanResults
				}
				// 查询测试用例关联的[测试组件] [测试分类] [测试分类id] 信息
				if module, err := new(EvaluateModule).FindById(itemQM.String("module_type_id")); err == nil {
					itemQM["module_name"] = module.ModuleName
					itemQM["module_type"] = module.ModuleType
					moduleKey := fmt.Sprintf("%s|%s|%s", module.ModuleName, module.ModuleType, module.Id.Hex())
					if items, has := itemRecords[moduleKey]; has {
						*items = append(*items, itemQM)
					} else {
						itemRecords[moduleKey] = &[]qmap.QM{itemQM}
					}
				}
			}
		}
		// 查询资产关联的信息
		if asset, err := new(EvaluateAsset).One(assetId); err == nil {
			assetChildren := []*qmap.QM{}
			for key, items := range itemRecords {
				// 从key中拆分出 [测试组件] [测试分类] [测试分类id] 信息
				info := strings.Split(key, "|")
				// 查询该测试用例所属的测试组件
				currentModule := &qmap.QM{}
				for _, module := range assetChildren {
					if moduleName := module.String("label"); moduleName == info[0] {
						currentModule = module
						break
					}
				}
				if len(*currentModule) == 0 {
					(*currentModule)["label"] = info[0]
					(*currentModule)["children"] = &([]qmap.QM{})
					assetChildren = append(assetChildren, currentModule)
				}
				// 创建测试分类
				currentModuleType := qmap.QM{
					"id":       info[2],
					"label":    info[1],
					"children": items,
				}
				moduleChildren := currentModule.Interface("children").(*[]qmap.QM)
				*moduleChildren = append(*moduleChildren, currentModuleType)
			}
			if len(assetChildren) == 0 {
				continue
			}
			assetRecord := &qmap.QM{
				"id":       asset.Id,
				"label":    asset.Name,
				"children": assetChildren,
			}
			assetRecords = append(assetRecords, assetRecord)
		}
	}
	return assetRecords, nil
}

// 查询项目任务关联的所有测试资产列表
func (this *EvaluateTaskItem) GetEvaluateTaskRelatedAssets(evaluateTaskId string) []string {
	query := bson.M{
		"evaluate_task_id": evaluateTaskId,
	}
	assetIds := make([]interface{}, 0)
	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_TASK_ITEM)
	assetIds, err := coll.Distinct(context.Background(), "asset_id", query)
	if err != nil {
		return []string{}
	}

	assertIdsStr := make([]string, 0)

	for _, v := range assetIds {
		assertIdsStr = append(assertIdsStr, v.(string))
	}
	return assertIdsStr
}

// 查询项目任务和资产关联的所有测试用例
func (this *EvaluateTaskItem) GetRelatedItems(evaluateTaskId, assetId string, rawInfo qmap.QM) ([]map[string]interface{}, error) {
	params := qmap.QM{
		"e_evaluate_task_id": evaluateTaskId,
		"e_asset_id":         assetId,
	}
	if itemId, has := rawInfo.TryString("item_id"); has && itemId != "" {
		params["e_item_id"] = itemId
	}
	if itemName, has := rawInfo.TryString("item_name"); has && itemName != "" {
		params["l_name"] = itemName
	}
	if recordAuditStatus, has := rawInfo.TryInt("record_audit_status"); has {
		params["e_record_audit_status"] = recordAuditStatus
	}
	return orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_TASK_ITEM, params).SetLimit(1000000).Find()
}

// ------------------------------------------------------
// ------------------------------------------------------
// ------------------------------------------------------

func EvaluateTaskItemUpdateAuditStatus(ctx context.Context, taskIdHex string, status int) error {
	selector := bson.M{"evaluate_task_id": taskIdHex}
	updateItem := bson.M{
		"$set": qmap.QM{
			"audit_status": status,
		},
	}

	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_TASK_ITEM)
	if _, err := coll.UpdateMany(ctx, selector, updateItem); err != nil {
		return err
	}
	return nil
}
