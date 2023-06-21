package mongo_model

import (
	"errors"
	"time"

	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"

	"github.com/globalsign/mgo/bson"
)

type EvaluateVulType struct {
	Id         int    `bson:"_id,omitempty"`
	Name       string `bson:"name"`        // 漏洞类型名称
	UpdateTime int64  `bson:"update_time"` // 更新时间
	CreateTime int64  `bson:"create_time"` // 创建时间
}

func (this *EvaluateVulType) GetAll() (*[]map[string]interface{}, error) {
	mongoClient := mongo.NewMgoSession(common.MC_EVALUATE_VUL_TYPE)
	mongoClient.SetLimit(1000000)
	return mongoClient.Get()
}

func (this *EvaluateVulType) Upsert(rawInfo qmap.QM) (int, error) {
	name := rawInfo.MustString("name")
	id := rawInfo.Int("id")
	if val, has := rawInfo.TryInt("id"); has {
		id = val
	} else {
		id = this.GetMaxId() + 1
	}

	params := qmap.QM{
		"e_name": name,
	}
	vulType := new(EvaluateVulType)
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_TYPE, params).One(vulType); err == nil {
		// 不允许漏洞类型名重复
		if id != vulType.Id {
			return 0, errors.New("该漏洞类型名称已经存在！")
		}
	}

	item := qmap.QM{
		"name":        rawInfo.MustString("name"),
		"update_time": time.Now().Unix(),
	}
	selector := bson.M{
		"_id": bson.M{"$eq": id},
	}
	upsertItem := bson.M{
		"$setOnInsert": bson.M{
			"create_time": time.Now().Unix(),
		},
		"$set": item,
	}
	_, err := mongo.NewMgoSession(common.MC_EVALUATE_VUL_TYPE).Upsert(selector, upsertItem)
	return id, err
}

func (this *EvaluateVulType) GetMaxId() int {
	if err := mongo.NewMgoSession(common.MC_EVALUATE_VUL_TYPE).Session.Find(nil).Sort("-_id").One(this); err == nil {
		return this.Id
	} else {
		return 0
	}
}

func (this *EvaluateVulType) GetVulIdTypeByName(name string) int {
	params := qmap.QM{
		"e_name": name,
	}
	vulType := new(EvaluateVulType)
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_TYPE, params).One(vulType); err == nil {
		return vulType.Id
	} else {
		return 0
	}
}
