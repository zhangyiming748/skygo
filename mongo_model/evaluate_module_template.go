package mongo_model

import (
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
)

type EvaluateModuleTemplate struct {
	Id         bson.ObjectId `bson:"_id,omitempty" json:"id"`        //测试组件id
	ModuleName string        `bson:"module_name" json:"module_name"` //测试组件名称
	ModuleType string        `bson:"module_type" json:"module_type"` //测试组件分类
	ItemName   string        `bson:"item_name" json:"item_name"`     //测试组件里的项
	Objective  string        `bson:"objective" json:"objective"`     //测试目的
	Level      int           `bson:"level" json:"level"`             //测试难度（1低、2中、3高）
}

func (this *EvaluateModuleTemplate) Create(rawInfo qmap.QM) (*EvaluateModuleTemplate, error) {
	this.Id = bson.NewObjectId()
	if val, has := rawInfo.TryString("module_name"); has {
		this.ModuleName = val
	}
	if val, has := rawInfo.TryString("module_type"); has {
		this.ModuleType = val
	}
	if val, has := rawInfo.TryString("item_name"); has {
		this.ItemName = val
	}
	if val, has := rawInfo.TryString("objective"); has {
		this.Objective = val
	}
	if val, has := rawInfo.TryInt("level"); has {
		this.Level = val
	}

	if err := mongo.NewMgoSession(common.MC_EvaluateModuleTemplate).Insert(this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateModuleTemplate) BulkDelete(rawInfo qmap.QM) (*qmap.QM, error) {
	effectNum := 0
	if rawIds, has := rawInfo.TrySlice("ids"); has {
		ids := []bson.ObjectId{}
		for _, id := range rawIds {
			ids = append(ids, bson.ObjectIdHex(id.(string)))
		}
		if len(ids) > 0 {
			match := bson.M{
				"_id": bson.M{"$in": ids},
			}
			if changeInfo, err := mongo.NewMgoSession(common.MC_EvaluateModuleTemplate).RemoveAll(match); err == nil {
				effectNum = changeInfo.Removed
			} else {
				return nil, err
			}
		}
	}
	return &qmap.QM{"number": effectNum}, nil
}
