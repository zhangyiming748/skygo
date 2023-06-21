package mongo_model

import (
	"github.com/globalsign/mgo/bson"
)

type FirmWareBinary struct {
	Id         bson.ObjectId `bson:"_id" json:"_id"`
	MasterId   string        `bson:"master_id" json:"master_id"`     //source表Id
	ProjectId  int           `bson:"project_id" json:"project_id"`   //工程的ID
	TaskId     int           `bson:"task_id" json:"task_id"`         //任务ID
	TemplateId int           `bson:"template_id" json:"template_id"` //模板ID
	IsExecuted int           `bson:"is_executed" json:"is_executed"` //是或否已执行（0：未执行，1：正在执行，2：已完成）
	SourceData string        `bson:"source_data" json:"source_data"` //base64加密后的原始数据
}
