package mongo_model

import (
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
)

type FirmWareUploadLog struct {
	Id              bson.ObjectId `bson:"_id" json:"_id"`
	ProjectName     string        `bson:"project_name" json:"project_name"`
	DeviceName      int           `bson:"device_name" json:"device_name"`
	DeviceModel     int           `bson:"device_model" json:"device_model"`
	DirmwareVersion int           `bson:"firmware_version" json:"firmware_version"`
	DeviceType      int           `bson:"device_type" json:"device_type"`
	FirmwareName    string        `bson:"firmware_name" json:"firmware_name"`
	TmpHdGilePath   string        `bson:"tmp_hd_file_path" json:"tmp_hd_file_path"`
	CreateTime      int64         `bson:"create_time" json:"create_time"`
	ProjectId       int           `bson:"project_id" json:"project_id"`
	FirmwareSize    int           `bson:"firmware_size" json:"firmware_size"`
	FirmwareMd5     string        `bson:"firmware_md5" json:"firmware_md5"`
	UploadUserId    int           `bson:"upload_user_id" json:"upload_user_id"`
	UploadUser      string        `bson:"upload_user" json:"upload_user"`
	UploadTime      int64         `bson:"upload_time" json:"upload_time"`
	Status          int           `bson:"status" json:"status"` //状态0：未处置，1处置中，2处置结束）
	TaskId          int           `bson:"task_id" json:"task_id"`
	TaskName        string        `bson:"task_name" json:"task_name"`
	TemplateId      int64         `bson:"template_id" json:"template_id"`
	TemplateName    string        `bson:"template_name" json:"template_name"`
	ResponseTime    int           `bson:"response_time" json:"response_time"`
	RealFileName    string        `bson:"real_file_name"` //真实文件名
	ResponseFile    string        `bson:"response_file"`  //接口返回地址
}

func (this *FirmWareUploadLog) Get(rawInfo qmap.QM) (*[]map[string]interface{}, error) {
	params := qmap.QM{}
	if val, has := rawInfo.TryInt("status"); has {
		params["e_status"] = val
	}

	if val, has := rawInfo.TryInt("ne_task_id"); has {
		params["ne_task_id"] = val
	}
	return mongo.NewMgoSession(common.MC_FIRMWARE_UPLOAD_LOG).AddCondition(params).Get()
}

func (this *FirmWareUploadLog) Update(rawInfo qmap.QM) (*FirmWareUploadLog, error) {
	id := rawInfo.MustString("id")
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(id),
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_FIRMWARE_UPLOAD_LOG, params)
	if err := mongoClient.One(&this); err == nil {
		this.Status = 1
		if err := mongoClient.Update(bson.M{"_id": this.Id}, this); err == nil {
			return this, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
