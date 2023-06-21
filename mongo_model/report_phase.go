package mongo_model

import (
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
)

type ReportPhase struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	Name       string        `bson:"name"`        //审核阶段名称
	Status     int           `bson:"status"`      //审核阶段状态（0 禁用，1激活）
	OperatorId int           `bson:"operator_id"` //操作人员id
	CreateTime int64         `bson:"create_time"`
}

func (this *ReportPhase) Create(opId int, rawInfo *qmap.QM) (*ReportPhase, error) {
	this.Id = bson.NewObjectId()
	this.Status = common.ENABLED
	this.OperatorId = opId
	this.Name = rawInfo.MustString("name")
	this.CreateTime = custom_util.GetCurrentMilliSecond()
	mongoClient := mongo.NewMgoSession(common.MC_REPORT_PHASE)
	if err := mongoClient.Insert(this); err == nil {
		//添加审核节点后，报告状态改为审核中

		return this, nil
	} else {
		return nil, err
	}
}
