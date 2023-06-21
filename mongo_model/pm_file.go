package mongo_model

import (
	"context"
	"errors"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"

	"github.com/globalsign/mgo/bson"
)

type PMFile struct {
	Id           primitive.ObjectID `bson:"_id" json:"id"`
	ProjectId    string             `bson:"project_id" json:"project_id"`         // 文件所属项目id
	ParentId     string             `bson:"parent_id" json:"parent_id"`           // 所属文件夹的文件id
	MetaFileId   string             `bson:"meta_file_id" json:"meta_file_id"`     // 源数据文件id
	MetaFileSize int                `bson:"meta_file_size" json:"meta_file_size"` // 源数据文件大小
	FileName     string             `bson:"file_name" json:"file_name"`           // 文件名称
	FileType     string             `bson:"file_type" json:"file_type"`           // 文件类型(dir:文件夹,doc:文件)
	OpId         int                `bson:"op_id" json:"op_id"`                   // 文件操作人id
	CreateTime   int64              `bson:"create_time" json:"create_time"`       // 文件创建时间
}

const (
	FILE_TYPE_DIR = "dir" // 文件类型:目录文件
	FILE_TYPE_DOC = "doc" // 文件类型:文档文件
)

func (this *PMFile) Create(ctx context.Context, projectId, metaFileId, fileName, fileType, parentId string, metaFileSize, opId int) (*PMFile, error) {
	this.Id = primitive.NewObjectID()
	this.ProjectId = projectId
	this.ParentId = parentId
	this.MetaFileId = metaFileId
	this.MetaFileSize = metaFileSize
	this.FileName = fileName
	this.OpId = opId
	this.CreateTime = custom_util.GetCurrentMilliSecond()
	if fileType := strings.ToLower(fileType); fileType == FILE_TYPE_DIR || fileType == FILE_TYPE_DOC {
		this.FileType = fileType
	} else {
		return nil, errors.New("Unknown file type!")
	}

	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_PROJECT_File)
	if _, err := coll.InsertOne(ctx, this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *PMFile) IsFileExist(projectId, fileName, fileType, parentId string) bool {
	params := qmap.QM{
		"e_project_id": projectId,
		"e_file_name":  fileName,
		"e_file_type":  fileType,
		"e_parent_id":  parentId,
	}
	if _, err := mongo.NewMgoSessionWithCond(common.MC_PROJECT_File, params).GetOne(); err == nil {
		return true
	} else {
		return false
	}
}

func (this *PMFile) Rename(projectId, id, newFileName string) (*PMFile, error) {
	params := qmap.QM{
		"e__id":        bson.ObjectIdHex(id),
		"e_project_id": projectId,
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_PROJECT_File, params)
	if err := mongoClient.One(&this); err == nil {
		// 如果重命名的是文档文件，则先修改文档文件的名字
		if this.FileType == "doc" && this.MetaFileId != "" {
			if newMetaFileId, err := mongo.GridFSRename(common.MC_File, newFileName, bson.ObjectIdHex(this.MetaFileId)); err == nil {
				this.MetaFileId = newMetaFileId
			} else {
				return nil, err
			}
		}
		this.FileName = newFileName
		if err := mongoClient.Update(bson.M{"_id": this.Id}, this); err == nil {
			return this, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}
