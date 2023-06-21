package mysql_model

import (
	"encoding/json"
	"errors"
	"time"

	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
)

type AssetTestPieceVersion struct {
	Id                 int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	AssetTestPieceId   int    `xorm:"not null comment('测试件记录id') INT(11)" json:"asset_test_piece_id"`
	Version            string `xorm:"comment('测试件版本') VARCHAR(255)" json:"version"`
	StorageType        int    `xorm:"not null default 1 comment('存储系统分类，1mongodb') TINYINT(3)" json:"storage_type"`
	CreateUserId       int    `xorm:"not null default 0 comment('创建人id') INT(11)" json:"create_user_id"`
	UpdateTime         int    `xorm:"not null default 0 comment('版本记录修改时间（秒）') INT(11)" json:"update_time"`
	FirmwareFileName   string `xorm:"not null default '' comment('固件文件名称') VARCHAR(255)" json:"firmware_file_name"`
	FirmwareFileUuid   string `xorm:"not null default '' comment('固件文件的唯一标识') VARCHAR(255)" json:"firmware_file_uuid"`
	FirmwareName       string `xorm:"not null default '' comment('固件名称') VARCHAR(255)" json:"firmware_name"`
	FirmwareSize       int64  `xorm:"not null default 0 comment('固件大小（kb）') BIGINT(20)" json:"firmware_size"`
	FirmwareDeviceType int    `xorm:"not null comment('1汽车网关(GW) 2远程通信单元(ECU) 3信息娱乐单元(IVI)') TINYINT(3)" json:"firmware_device_type"`
	IsDelete           int    `xorm:"not null default 1 comment('是否删除， 1否 2是') TINYINT(3)" json:"is_delete"`
	CreateTime         int    `xorm:"not null default 0 comment('版本记录创建时间（秒）') INT(11)" json:"create_time"`
	FirmwareVersion    string `xorm:"comment('测试件版本') VARCHAR(255)" json:"firmware_version"`
}

// 根据测试件记录id，查询出最新的一个版本信息
func AssetTestPieceVersionFindLatest(assetTestPieceId int) (*AssetTestPieceVersion, bool) {
	model := AssetTestPieceVersion{}
	session := mysql.GetSession()
	session.Where("asset_test_piece_id = ?", assetTestPieceId)
	session.OrderBy("create_time desc")
	if has, err := session.Get(&model); err != nil {
		panic(err)
	} else {
		return &model, has
	}
}

// 根据测试件记录id，上传一个固件
func AssetTestPieceVersionUploadFirmware(versionId int, firmwareName, fileName string, fileSize int64, fileUuid string, firmwareDeviceType int, version string) error {
	model := AssetTestPieceVersion{}
	has, err := mysql.GetSession().ID(versionId).Get(&model)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("记录不存在")
	}

	model.UpdateTime = int(time.Now().Unix())
	model.FirmwareFileUuid = fileUuid
	model.FirmwareSize = fileSize
	model.FirmwareName = firmwareName
	model.FirmwareFileName = fileName
	model.FirmwareVersion = version
	model.FirmwareDeviceType = firmwareDeviceType

	_, err = mysql.GetSession().ID(versionId).Update(&model)
	return err
}

// 根据测试件记录id，更新固件信息
func AssetTestPieceVersionUploadInfo(versionId int, firmwareName, version string, firmwareDeviceType int) error {
	model := AssetTestPieceVersion{}
	has, err := mysql.GetSession().ID(versionId).Get(&model)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("记录不存在")
	}

	model.UpdateTime = int(time.Now().Unix())
	model.FirmwareName = firmwareName
	model.FirmwareVersion = version
	model.FirmwareDeviceType = firmwareDeviceType

	_, err = mysql.GetSession().ID(versionId).Update(&model)
	return err
}

// 根据测试件记录id，查询所有版本
func AssetTestPieceVersionFindAll(assetTestPieceId int) []map[string]interface{} {
	models := make([]AssetTestPieceVersion, 0)
	session := mysql.GetSession()
	session.Where("asset_test_piece_id = ?", assetTestPieceId)
	session.OrderBy("create_time desc")
	if err := session.Find(&models); err != nil {
		return nil
	}
	// all := orm.AllResult{}
	// fmt.Println("models:", models)
	// for i := 0; i < reflect.ValueOf(models).Elem().Len(); i++ {
	//	// reflect.ValueOf(modelsPtr).Elem().Index(i).Type().Kind() is reflect.Struct
	//	one := reflect.ValueOf(models).Elem().Index(i).Interface()
	//	all = append(all, orm.StructToMap(one))
	// }
	// for key, one := range all {
	//	all[key] = one
	// }
	tmp, _ := json.Marshal(models)
	all := orm.AllResult{}
	json.Unmarshal(tmp, &all)
	return all

}

// 新增测试件版本
func (this *AssetTestPieceVersion) Create() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}

// 删除测试件版本
func (this *AssetTestPieceVersion) Delete(id int) (int64, error) {
	return mysql.GetSession().ID(id).Delete(this)
}

// 按照id更新
func (this *AssetTestPieceVersion) UpdateById(id int, data qmap.QM) (*AssetTestPieceVersion, error) {
	has, err := mysql.GetSession().ID(id).Get(this)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("数据不存在")
	}
	if version, has := data.TryString("version"); has {
		this.Version = version
	}
	if _, err := mysql.GetSession().ID(id).Update(this); err != nil {
		return nil, err
	} else {
		return this, nil
	}
}

// 按照id更新
func (this *AssetTestPieceVersion) FindById(id int) (*AssetTestPieceVersion, error) {
	if has, err := mysql.GetSession().ID(id).Get(this); err == nil {
		if has {
			return this, nil
		} else {
			return nil, errors.New("测试件不存在")
		}
	} else {
		return nil, err
	}
}
