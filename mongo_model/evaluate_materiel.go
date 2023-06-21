package mongo_model

import (
	"errors"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"
)

type EvaluateMateriel struct {
	Id         bson.ObjectId `bson:"_id,omitempty" json:"_id,omitempty"`
	ProjectId  string        `bson:"project_id" json:"project_id"`   //所属项目id
	Name       string        `bson:"name" json:"name"`               //物料名称
	AssetId    string        `bson:"asset_id" json:"asset_id"`       //资产ID
	Number     int           `bson:"number" json:"number"`           //设备数量
	Image      string        `bson:"image" json:"image"`             //设备照片
	Comment    string        `bson:"comment" json:"comment"`         //备注
	CreateTime int64         `bson:"create_time" json:"create_time"` //登记时间
}

func (this *EvaluateMateriel) Create(rawInfo *qmap.QM) (*EvaluateMateriel, error) {

	this.Id = bson.NewObjectId()
	this.ProjectId = rawInfo.MustString("project_id")
	this.Name = rawInfo.String("name")
	this.Number = rawInfo.Int("number")
	this.Image = rawInfo.MustString("image")
	this.AssetId = rawInfo.String("asset_id")
	this.Comment = rawInfo.String("comment")
	this.CreateTime = custom_util.GetCurrentMilliSecond()
	if this.checkNameExist(this.ProjectId, this.Name) {
		return this, errors.New("该物料名称已被使用")
	}
	mongoClient := mongo.NewMgoSession(common.MC_EVALUATE_MATERIEL)
	if err := mongoClient.Insert(this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateMateriel) checkNameExist(projectId, name string) bool {
	params := qmap.QM{
		"e_project_id": projectId,
		"e_name":       name,
	}
	if result, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MATERIEL, params).GetOne(); err == nil {
		if len(*result) != 0 {
			//存在重复的数据
			return true
		}
	}

	//不存在重复的数据
	return false
}

func (this *EvaluateMateriel) Update(projectId string, rawInfo qmap.QM) (*EvaluateMateriel, error) {
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(projectId),
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MATERIEL, params)
	if err := mongoClient.One(&this); err == nil {
		//字段project_id不可修改

		if val, has := rawInfo.TryString("name"); has {
			this.Name = val
		}
		if val, has := rawInfo.TryInt("number"); has {
			this.Number = val
		}
		if val, has := rawInfo.TryString("image"); has {
			this.Image = val
		}
		if val, has := rawInfo.TryString("comment"); has {
			this.Comment = val
		}
		if val, has := rawInfo.TryString("asset_id"); has {
			this.AssetId = val
		}

		if err := mongoClient.Update(bson.M{"_id": this.Id}, this); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Materiel not found")
	}

	return this, nil
}

/*
func (this *EvaluateMateriel) GetAllWithTaskId1(rawInfo *qmap.QM) (*qmap.QM, error) {
	//1. 从任务和测试用例的关系表中，查询出项目id和资产id
	var projectId string //通过任务id查询资产，最终所有的资产里项目id都是一个，所以是固定的
	assetIds := make([]string, 0) //通过任务id查询资产id，可以是多个
	params := qmap.QM{
		"e_evaluate_task_id": rawInfo.MustString("id"),
	}
	if taskItems, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK_ITEM, params).Get(); err == nil {
		for _, taskItem := range *taskItems {
			projectId = taskItem["project_id"].(string)
			assetId := taskItem["asset_id"].(string)
			assetIds = append(assetIds, assetId)
		}
		paramsMaterial := qmap.QM{
			"e_project_id": projectId,
			"in_asset_id": assetIds,
		}
		//2. 根据项目id和资产id，查询出符合条件的物料
		mgoSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MATERIEL, paramsMaterial)
		mgoSession.SetTransformFunc(this.EvaluateMaterialTransformer)
		return mgoSession.GetPage()
	} else {
		return nil, err
	}
}
*/

func (this *EvaluateMateriel) GetAllWithTaskId(rawInfo *qmap.QM) (qmap.QM, error) {
	//1. 通过task查询 出asset
	taskId := rawInfo.MustString("id")
	_id, _ := primitive.ObjectIDFromHex(taskId)
	params := qmap.QM{
		"e__id": _id,
	}
	task := new(EvaluateTask)
	if err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_TASK, params).One(task); err != nil {
		return nil, err
	}
	//2. 查询materiel
	assets := []string{}
	for asset, _ := range task.AssetVersions {
		assets = append(assets, asset)
	}
	params = qmap.QM{
		"in_asset_id": assets,
	}
	widget := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_MATERIEL, params)
	widget.SetTransformerFunc(this.EvaluateMaterialTransformer)
	return widget.PaginatorFind()
}

func (this EvaluateMateriel) EvaluateMaterialTransformer(data qmap.QM) qmap.QM {
	assetId := data.MustString("asset_id")
	asset, _ := new(EvaluateAsset).GetOne(assetId)
	assetName := asset.String("name")
	evaluateType := asset.String("evaluate_type")
	data["asset_name"] = assetName
	data["evaluate_type"] = evaluateType
	return data
}
