package mongo_model

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
)

type ReportNode struct {
	Id         bson.ObjectId `bson:"_id,omitempty"`
	ProjectId  string        `bson:"project_id"`             //项目ID
	ReportId   string        `bson:"report_id"`              //报告ID
	Name       string        `bson:"name"`                   //阶段名称
	Result     int           `bson:"result"`                 //审核结果
	AuditorId  int           `bson:"auditor_id"`             //审核人ID
	AuditTime  int           `bson:"audit_time"`             //审核时间
	History    []History     `bson:"history" json:"history"` //阶段历史
	CreateTime int64         `bson:"create_time"`
}

func (this *ReportNode) Create(rawInfo qmap.QM) (*ReportNode, error) {
	this.Id = bson.NewObjectId()
	this.ProjectId = rawInfo.MustString("project_id")
	this.ReportId = rawInfo.MustString("report_id")
	this.AuditorId = rawInfo.MustInt("auditor_id")
	status, hasStatus := rawInfo.TryString("status")
	if hasStatus && status == "create" {
		this.Result = common.RAS_SUCCESS
	} else {
		this.Result = common.RAS_NEW
	}
	this.Name = rawInfo.MustString("name")
	this.History = append(this.History, rawInfo["history"].(History))
	mongoClient := mongo.NewMgoSession(common.MC_REPORT_NODE)
	this.CreateTime = custom_util.GetCurrentMilliSecond()
	if err := mongoClient.Insert(this); err == nil {
		//创建节点后，报告状态改为审核中
		if !hasStatus || status != "create" {
			selector := bson.M{"_id": bson.ObjectIdHex(this.ReportId)}
			updateItem := bson.M{
				"$set": qmap.QM{
					"status": common.RS_AUDIT,
				},
			}
			if err := mongo.NewMgoSession(common.MC_REPORT).Update(selector, updateItem); err != nil {
				fmt.Println(err)
				return nil, err
			}
		}
		return this, nil
	} else {
		return nil, err
	}
}
