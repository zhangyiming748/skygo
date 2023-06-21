package mongo_model

import (
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
)

type EvaluateAssetVersion struct {
	Id           bson.ObjectId `bson:"_id"`            //版本id
	AssetId      string        `bson:"asset_id"`       //资产id
	ProjectId    string        `bson:"project_id"`     //项目id
	Name         string        `bson:"name"`           //资产名称
	Version      string        `bson:"version"`        //资产版本
	EvaluateType string        `bson:"evaluate_type"`  //资产类型
	ModuleTypeId []interface{} `bson:"module_type_id"` //资产分类
	Attributes   interface{}   `bson:"attributes"`     //测试对象属性
	OpId         int64         `bson:"op_id"`          //操作人id
	Note         string        `bson:"note"`           //备注
	CreateTime   int64         `bson:"create_time"`    //创建时间
	UpdateTime   int64         `bson:"update_time"`    //更新时间
}

func (this *EvaluateAssetVersion) Create(assetStruct *EvaluateAsset) (*EvaluateAssetVersion, error) {
	this.Id = bson.NewObjectId()
	this.AssetId = assetStruct.Id
	this.ProjectId = assetStruct.ProjectId
	this.Version = assetStruct.Version
	this.ModuleTypeId = assetStruct.ModuleTypeId
	this.Name = assetStruct.Name
	this.EvaluateType = assetStruct.EvaluateType
	this.OpId = assetStruct.OpId
	this.Note = assetStruct.Note
	this.Attributes = assetStruct.Attributes
	microSecond := custom_util.GetCurrentMilliSecond()
	this.UpdateTime = microSecond
	this.CreateTime = microSecond
	if err := mongo.NewMgoSession(common.MC_EVALUATE_ASSET_VERSION).Insert(this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateAssetVersion) GetVersion(assetId string) (*qmap.QM, error) {
	params := qmap.QM{
		"e_asset_id": assetId,
	}
	info, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET_VERSION, params).GetOne()
	if err == nil {
		idObject := []bson.ObjectId{}
		for _, id := range (*info)["module_type_id"].([]interface{}) {
			idObject = append(idObject, bson.ObjectIdHex(id.(string)))
		}
		moduleType, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, qmap.QM{"in__id": idObject}).Get()
		(*info)["module_type"] = moduleType
		return info, err
	}
	return info, err
}
