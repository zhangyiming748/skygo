package mysql_model

import (
	"errors"
	"time"

	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"
	"xorm.io/builder"

	"skygo_detection/lib/common_lib/mysql"
)

type AssetTestPiece struct {
	Id             int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	Name           string `xorm:"not null default '' comment('测试件名称') VARCHAR(255)" json:"name"`
	AssetVehicleId int    `xorm:"not null comment('车型记录id') INT(11)" json:"asset_vehicle_id"`
	PieceType      string `xorm:"not null comment('测试件型号') VARCHAR(255)" json:"piece_type"`
	Detail         string `xorm:"not null default '' comment('测试件描述') VARCHAR(255)" json:"detail"`
	CreateTime     int    `xorm:"not null default 0 comment('创建时间（秒）') INT(11)" json:"create_time"`
	UpdateTime     int    `xorm:"default 0 comment('更新时间（秒）') INT(11)" json:"update_time"`
	CreateUserId   int    `xorm:"not null default 0 comment('创建人id') INT(11)" json:"create_user_id"`
}

const (
	StorageTypeMongo = 1 // 存储类型，即存储文件使用的存储方式， 1. mongo
)

// 创建表单
type AssetTestPieceCreateForm struct {
	Name           string `json:"name"`             // 测试件名称
	Version        string `json:"code"`             // 测试件版本
	AssetVehicleId int    `json:"asset_vehicle_id"` // 车型记录id，来自车型代号下拉列表
	Detail         string `json:"detail"`           // 测试件描述
}

// 根据表单创建测试件记录， 返回新建记录id
func AssetTestPieceCreateFromForm(form *AssetTestPieceCreateForm, uid int) (int, error) {
	// 判断存在
	model := AssetVehicle{}
	if has, err := mysql.FindById(form.AssetVehicleId, &model); err != nil {
		panic(err)
	} else {
		if !has {
			return 0, errors.New("车型不存在")
		}
	}

	nowTime := int(time.Now().Unix())

	// 开启事务
	session := mysql.GetSession()
	defer session.Close()
	if err := session.Begin(); err != nil {
		return 0, err
	}

	// 创建测试件主记录
	pModel := AssetTestPiece{}
	pModel.Name = form.Name
	pModel.AssetVehicleId = form.AssetVehicleId
	pModel.Detail = form.Detail
	pModel.CreateTime = nowTime
	pModel.CreateUserId = uid
	if _, err := session.Insert(&pModel); err != nil {
		session.Rollback()
		return 0, err
	}

	// 创建测试件版本记录
	vModel := AssetTestPieceVersion{}
	vModel.AssetTestPieceId = pModel.Id
	vModel.Version = form.Version
	vModel.StorageType = StorageTypeMongo
	vModel.CreateUserId = uid
	vModel.UpdateTime = 0
	vModel.CreateTime = nowTime
	vModel.FirmwareFileUuid = "" // 还没上传固件文件
	vModel.FirmwareName = ""
	vModel.FirmwareSize = 0
	vModel.FirmwareDeviceType = 0
	if _, err := session.InsertOne(&vModel); err != nil {
		session.Rollback()
		return 0, err
	}

	// 提交事务
	if err := session.Commit(); err != nil {
		return 0, err
	}
	return pModel.Id, nil
}

// 按照id查询
func AssetTestPieceFindById(userId int) (*AssetTestPiece, bool, error) {
	model := AssetTestPiece{}
	if has, err := mysql.GetSession().ID(userId).Get(&model); err != nil {
		return nil, has, err
	} else {
		return &model, has, nil
	}
}

// 按照id更新
func AssetTestPieceUpdateById(id int, data qmap.QM) (*AssetTestPiece, error) {
	model := AssetTestPiece{}
	has, err := mysql.GetSession().ID(id).Get(&model)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("数据不存在")
	}

	// 测试件名称
	if name, has := data.TryString("name"); has {
		model.Name = name
	}

	// 车型记录id
	if assetVehicleId, has := data.TryInt("asset_vehicle_id"); has {
		model.AssetVehicleId = assetVehicleId
	}

	if detail, has := data.TryString("detail"); has {
		model.Detail = detail
	}
	//获取更新时间
	nowTime := int(time.Now().Unix())
	model.UpdateTime = nowTime

	if _, err := mysql.GetSession().ID(id).Update(&model); err != nil {
		return nil, err
	} else {
		return &model, nil
	}
}

// 根据测试件id，删除单条记录，注意级联删除
func AssetTestPieceDeleteById(id int) error {
	// 启动事务
	session := mysql.GetSession()
	err := session.Begin()
	if err != nil {
		return err
	}
	defer session.Close()

	// 删除测试件记录, 通过返回的int值判断记录是否存在
	if count, err := session.ID(id).Delete(AssetTestPiece{}); err != nil {
		session.Rollback()
		return err
	} else {
		if count == 0 {
			session.Rollback()
			return errors.New("记录不存在")
		}
	}

	// 删除测试件的所有版本记录（软删除）
	update := gin.H{
		"is_delete": 2,
	}
	if _, err := session.Table(AssetTestPieceVersion{}).
		Where("asset_test_piece_id = ?", id).Update(update); err != nil {
		session.Rollback()
		return err
	}

	// 删除测试件版本记录中的文件记录（软删除）
	versionIds := make([]int, 0)
	session.Table(AssetTestPieceVersion{}).Select("id").Where("asset_test_piece_id").Find(&versionIds)
	if _, err := session.Table(AssetTestPieceVersionFile{}).
		Where(builder.In("version_id", versionIds)).Update(update); err != nil {
		session.Rollback()
		return err
	}

	if err := session.Commit(); err != nil {
		return err
	}

	return nil
}
