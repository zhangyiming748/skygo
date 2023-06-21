package mongo_model

import (
	"context"
	"errors"
	"time"

	"github.com/globalsign/mgo/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/orm_mongo"
)

type ScannerTask struct {
	Id           primitive.ObjectID `bson:"_id"`
	TaskId       string             `bson:"task_id"`
	Name         string             `bson:"name"`
	TaskType     string             `bson:"task_type"`
	Status       int                `bson:"status"`
	RetryTimes   int                `bson:"retry_times"`
	NextExecTime int64              `bson:"next_exec_time"`
}

func (this *ScannerTask) TestInsert() {
	this.Id = primitive.NewObjectID()
	this.TaskId = ""
	this.Name = "测试任务"
	this.TaskType = "scanner_test"
	this.Status = common.SCANNER_STATUS_READY
	this.RetryTimes = 0
	this.NextExecTime = time.Now().Unix()

	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_SCANNER_TASK)
	if _, err := coll.InsertOne(context.Background(), this); err == nil {
		panic(err)
	}
}

func (this *ScannerTask) TaskInsert(taskId, taskName, taskType string) error {
	this.Id = primitive.NewObjectID()
	this.TaskId = taskId
	this.Name = taskName
	this.TaskType = taskType
	this.Status = common.SCANNER_STATUS_READY
	this.RetryTimes = 0
	this.NextExecTime = time.Now().Unix()

	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_SCANNER_TASK)
	if _, err := coll.InsertOne(context.Background(), this); err == nil {
		return nil
	} else {
		return err
	}
}

func (this *ScannerTask) Update(id string, rawInfo qmap.QM) error {
	_id, _ := primitive.ObjectIDFromHex(id)
	params := qmap.QM{
		"e__id": _id,
	}
	w := orm_mongo.NewWidgetWithParams(common.MC_SCANNER_TASK, params)
	if err := w.One(&this); err == nil {
		if val, has := rawInfo.TryString("name"); has {
			this.Name = val
		}
		if val, has := rawInfo.TryString("task_type"); has {
			this.TaskType = val
		}
		if val, has := rawInfo.TryInt("status"); has {
			this.Status = val
		}
		if val, has := rawInfo.TryInt("retry_times"); has {
			this.RetryTimes = val
		}
		if val, has := rawInfo.TryInt("status"); has {
			this.Status = val
		}
		if val, has := rawInfo.TryInt("next_exec_time"); has {
			this.NextExecTime = int64(val)
		}
		coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_SCANNER_TASK)
		if _, err := coll.UpdateOne(context.Background(), bson.M{"_id": this.Id}, bson.M{"$set": this}); err != nil {
			return err
		}
	} else {
		return errors.New("Item not found")
	}
	return nil
}
