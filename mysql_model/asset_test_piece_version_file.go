package mysql_model

import (
	"time"

	"skygo_detection/lib/common_lib/mysql"
)

type AssetTestPieceVersionFile struct {
	Id           int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	VersionId    int    `xorm:"not null comment('固件某版本记录的id') INT(11)" json:"version_id"`
	FileName     string `xorm:"not null comment('文件名称') VARCHAR(255)" json:"file_name"`
	FileSize     int64  `xorm:"not null default 0 comment('文件大小（kb）') BIGINT(20)" json:"file_size"`
	StorageType  int    `xorm:"not null comment('存储类型，1mongodb') TINYINT(3)" json:"storage_type"`
	FileUuid     string `xorm:"not null default '' comment('文件存储唯一标识uuid') VARCHAR(255)" json:"file_uuid"`
	CreateTime   int    `xorm:"not null comment('创建时间，即文件上传时间') INT(11)" json:"create_time"`
	IsDelete     int    `xorm:"not null default 1 comment('是否删除， 1否 2是') TINYINT(3)" json:"is_delete"`
	DeleteUserId int    `xorm:"not null comment('文件删除操作用户id') INT(11)" json:"delete_user_id"`
	DeleteTime   int    `xorm:"not null comment('删除时间') INT(11)" json:"delete_time"`
}

const (
	IsDeleteNo  = 1
	IsDeleteYes = 2
)

// 获取测试件某个版本的 “文件记录列表”
// versionId 测试件某个版本的记录ID
func AssetTestPieceVersionFileFindDetail(versionId int) []AssetTestPieceVersionFile {
	models := make([]AssetTestPieceVersionFile, 0)

	session := mysql.GetSession()
	session.Where("is_delete = ?", 1)
	session.And("version_id = ?", versionId)
	session.Find(&models)

	return models
}

// 新增一个文件记录
func AssetTestPieceVersionFileCreate(versionId int, fileName string, fileSize int64, fileUuid string) error {
	model := AssetTestPieceVersionFile{}
	model.VersionId = versionId
	model.FileName = fileName
	model.FileSize = fileSize
	model.FileUuid = fileUuid
	model.StorageType = StorageTypeMongo
	model.CreateTime = int(time.Now().Unix())
	model.IsDelete = IsDeleteNo
	model.DeleteTime = 0
	model.DeleteUserId = 0

	_, err := mysql.GetSession().InsertOne(&model)
	return err
}

// 删除测试件版本
func (this *AssetTestPieceVersionFile) Delete(id int) (int64, error) {
	return mysql.GetSession().ID(id).Delete(this)
}
