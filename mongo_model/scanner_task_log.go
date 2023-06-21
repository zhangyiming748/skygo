package mongo_model

import (
	"time"

	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"

	"github.com/globalsign/mgo/bson"
)

type ScannerTaskLog struct {
	Id       bson.ObjectId `bson:"_id"`
	TaskId   string        `bson:"task_id"`
	Name     string        `bson:"name"`
	Level    string        `bson:"status"`
	TaskType string        `bson:"task_type"`
	//TaskDetail string        `bson:"task_detail"`
	ErrMsg     string `bson:"err_msg"`
	CreateTime int64  `bson:"create_time"`
}

func (this *ScannerTaskLog) Insert(taskInfo qmap.QM, errMsg string) error {
	this.Id = bson.NewObjectId()
	this.Name = taskInfo.MustString("name")
	this.TaskId = taskInfo.String("id")
	this.ErrMsg = errMsg
	if errMsg == "" {
		this.Level = "info"
	} else {
		this.Level = "error"
	}
	this.TaskType = taskInfo.String("task_type")
	//this.TaskDetail = taskInfo.Map("task_detail").ToString()
	this.CreateTime = time.Now().Unix()
	if err := mongo.NewMgoSession(common.MC_SCANNER_TASK_LOG).Insert(this); err != nil {
		return err
	} else {
		return nil
	}
}
