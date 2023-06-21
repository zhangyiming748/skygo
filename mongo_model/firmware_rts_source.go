package mongo_model

import (
	"encoding/base64"

	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
)

type FirmWareSource struct {
	Id         bson.ObjectId `bson:"_id" json:"_id"`
	MasterId   string        `bson:"master_id" json:"master_id"`     //source表Id
	ProjectId  int           `bson:"project_id" json:"project_id"`   //工程的ID
	TaskId     int           `bson:"task_id" json:"task_id"`         //任务ID
	TemplateId int           `bson:"template_id" json:"template_id"` //模板ID
	IsExecuted int           `bson:"is_executed" json:"is_executed"` //是或否已执行（0：未执行，-1：正在执行，1：已完成）
	SourceData string        `bson:"source_data" json:"source_data"` //GridFS文件ID
}

func (this *FirmWareSource) Create(rawInfo qmap.QM) (*FirmWareSource, error) {
	if val, has := rawInfo.TryString("master_id"); has {
		this.MasterId = val
	}
	if val, has := rawInfo.TryInt("project_id"); has {
		this.ProjectId = val
	}
	if val, has := rawInfo.TryInt("task_id"); has {
		this.TaskId = val
	}
	if val, has := rawInfo.TryInt("template_id"); has {
		this.TemplateId = val
	}
	this.Id = bson.NewObjectId()
	this.IsExecuted = 0
	data := rawInfo.MustString("source_data")
	base64Encode := base64.StdEncoding.EncodeToString([]byte(data))
	if fileId, err := mongo.GridFSUpload(common.MC_File, "Source Data", []byte(base64Encode)); err == nil {
		this.SourceData = fileId
	} else {
		this.SourceData = base64Encode
	}

	err := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_SOURCE).Insert(this)
	return this, err
}
