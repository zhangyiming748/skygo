package mongo_model

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"
)

// 测试记录
type EvaluateRecord struct {
	Id             bson.ObjectId `bson:"_id" json:"id"`                            // 测试记录id
	ProjectId      string        `bson:"project_id" json:"project_id"`             // 项目id
	EvaluateTaskId string        `bson:"evaluate_task_id" json:"evaluate_task_id"` // 项目任务id
	AssetId        string        `bson:"asset_id" json:"asset_id"`                 // 资产id
	AssetVersion   string        `bson:"asset_version" json:"asset_version"`       // 资产版本
	ItemId         string        `bson:"item_id" json:"item_id"`                   // 测试用例id
	TestPhase      int           `bson:"test_phase" json:"test_phase"`             // 测试阶段 （1:初测、2:复测1、3:复测2、4:复测3 ...）
	TestTime       int           `bson:"test_time" json:"test_time"`               // 测试次数
	TestProcedure  string        `bson:"test_procedure" json:"test_procedure"`     // 测试过程
	ToolId         string        `bson:"tool_id" json:"tool_id"`                   // 测试工具id
	ToolTestResult string        `bson:"tool_test_result" json:"tool_test_result"` // 工具测试结果
	Attachment     interface{}   `bson:"attachment" json:"attachment"`             // 测试附件
	OpId           int64         `bson:"op_id" json:"op_id"`                       // 操作人id
	UpdateTime     int64         `bson:"update_time" json:"update_time"`           // 更新时间
	CreateTime     int64         `bson:"create_time" json:"create_time,omitempty"` // 创建时间
}

// 添加或者修改测试记录
func (this *EvaluateRecord) Upsert(rawInfo qmap.QM, opId int64) error {
	id := rawInfo.MustString("id")
	// RecordID 和 testId 值一样
	taskItemId := id
	taskItem, err := new(EvaluateTaskItem).GetOne(taskItemId)
	if err != nil {
		return err
	}

	// 判断是否有操作测试用例的权限
	itemId := taskItem.ItemId
	evaluateItem, err := new(EvaluateItem).GetOne(itemId)
	if err != nil {
		return err
	}
	item := qmap.QM{
		"project_id":       evaluateItem.ProjectId,
		"task_item_id":     taskItemId,
		"evaluate_task_id": taskItem.EvaluateTaskId,
		"asset_id":         evaluateItem.AssetId,
		"asset_version":    new(EvaluateTask).GetAssetVersion(taskItem.EvaluateTaskId, evaluateItem.AssetId),
		"test_phase":       taskItem.TestPhase,
		"item_id":          itemId,
		"test_procedure":   rawInfo.MustString("test_procedure"),
		"attachment":       rawInfo.Interface("attachment"),
		"op_id":            opId,
	}
	selector := bson.M{
		"_id": bson.M{"$eq": bson.ObjectIdHex(id)},
	}
	upsertItem := bson.M{
		"$setOnInsert": bson.M{
			"create_time": custom_util.GetCurrentMilliSecond(),
		},
		"$set": item,
	}
	// 更新测试用例的测试时间
	if err := this.UpdateTestTime(id, itemId); err != nil {
		fmt.Println(err)
	}

	if _, err = mongo.NewMgoSession(common.MC_EVALUATE_RECORD).Upsert(selector, upsertItem); err == nil {
	} else {
		return err
	}
	return nil
}

// 如果测试记录是第一次添加，则更新测试用例关系表中的测试时间
func (this *EvaluateRecord) UpdateTestTime(recordId, itemId string) error {
	fmt.Println(recordId)
	if _, err := this.GetOne(recordId); err != nil {
		// 如果测试记录不存在,则更新测试用例的测试时间
		testTime := custom_util.GetCurrentMilliSecond()
		rawInfo := qmap.QM{
			"test_time": testTime,
		}
		if _, err := new(EvaluateTaskItem).Update(recordId, rawInfo); err != nil {
			return err
		}
		fmt.Println(itemId, rawInfo)
		//同步更新item主表测试时间
		if _, err := new(EvaluateItem).Update(itemId, rawInfo); err != nil {
			return err
		}
	}
	return nil
}

// 通过查询测试记录
func (this *EvaluateRecord) GetOne(id string) (interface{}, error) {
	_id, _ := primitive.ObjectIDFromHex(id)
	params := qmap.QM{
		"e__id": _id,
	}
	return orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_RECORD, params).Get()
}
