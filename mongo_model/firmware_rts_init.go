package mongo_model

import "github.com/globalsign/mgo/bson"

type FirmWareRtsInit struct {
	ID          bson.ObjectId `bson:"_id" json:"_id"`
	UploadLogId string        `bson:"upload_log_id" json:"upload_log_id"`
	ProjectId   string        `bson:"project_id" json:"project_id"`
	TasksId     string        `bson:"tasks_id" json:"tasks_id"`
	FirmwareMd5 string        `bson:"firmware_md5" json:"firmware_md5"`
	DirNum      int           `bson:"dir_num" json:"dir_num"`
	FileNum     int           `bson:"file_num" json:"file_num"`
	LinkNum     int           `bson:"link_num" json:"link_num"`
	NodeNume    int           `bson:"node_num" json:"node_num"`
}
