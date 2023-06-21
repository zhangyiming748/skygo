package mongo_model

import (
	"errors"
	"fmt"

	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
)

type Factory struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Name       string        `bson:"name"` //车厂名称
	CreateTime int64         `bson:"create_time"`
}

func (this *Factory) Create(rawInfo *qmap.QM) (*Factory, error) {
	name := rawInfo.String("name")
	params := qmap.QM{
		"e_name": name,
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_FACTORY, params)
	if err := mongoClient.One(&this); err == nil {
		return nil, errors.New(fmt.Sprintf("已存在名称为:%s的厂商，请重新输入新的厂商名称", name))
	} else {
		this.Id = bson.NewObjectId()
		this.Name = name
		this.CreateTime = custom_util.GetCurrentMilliSecond()
		if err := mongoClient.Insert(this); err == nil {
			return this, nil
		} else {
			return nil, err
		}
	}
}

func (this *Factory) Update(factoryId string, rawInfo qmap.QM) (*Factory, error) {
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(factoryId),
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_FACTORY, params)
	if err := mongoClient.One(&this); err == nil {
		if val, has := rawInfo.TryString("name"); has {
			this.Name = val
		}
		if err := mongoClient.Update(bson.M{"_id": this.Id}, this); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Factory not found")
	}
	return this, nil
}

func (this *Factory) checkDuplicate(name string) bool {
	params := qmap.QM{
		"e_name": name,
	}
	if result, err := mongo.NewMgoSessionWithCond(common.MC_FACTORY, params).GetOne(); err == nil {
		if len(*result) != 0 {
			//存在重复的数据
			return true
		}
	}
	//不存在重复的数据
	return false
}
