package mysql_model

import (
	"errors"
	"fmt"
	"time"

	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/log"
	"skygo_detection/lib/common_lib/mysql"
)

type AssetVehicle struct {
	Id           int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	SerialNumber string `xorm:"not null comment('车型编号，程序自动生成') VARCHAR(255)" json:"serial_number"`
	Brand        string `xorm:"not null default '' comment('车型品牌') VARCHAR(255)" json:"brand"`
	Code         string `xorm:"not null default '' comment('车型代号') VARCHAR(255)" json:"code"`
	Detail       string `xorm:"not null default '' comment('车型描述') VARCHAR(255)" json:"detail"`
	CreateUserId int    `xorm:"not null comment('创建用户id') INT(11)" json:"create_user_id"`
	UpdateTime   int    `xorm:"not null default 0 updated comment('创建时间（秒）') INT(11)" json:"update_time"`
	CreateTime   int    `xorm:"not null default 0 comment('更新时间') INT(11)" json:"create_time"`
}

// 创建表单
type AssetVehicleCreateForm struct {
	Brand  string `json:"brand"`
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

// 基于表单创建
func AssetVehicleCreateFromForm(form *AssetVehicleCreateForm, uid int) (*AssetVehicle, error) {
	model := AssetVehicle{}
	// todo 车型编号
	model.SerialNumber = ""
	model.Brand = form.Brand
	model.Code = form.Code
	model.Detail = form.Detail
	model.CreateUserId = uid
	model.UpdateTime = 0
	model.CreateTime = int(time.Now().Unix())

	// 查询code是否存在
	_, has := AssetVehicleGetByCode(model.Code)
	if has {
		return nil, errors.New("车型代号已经存在")
	}

	_, err := mysql.GetSession().InsertOne(&model)
	return &model, err
}

func generateAssetVehicleNumber() {

}

// 按照id查询
func AssetVehicleFindById(userId int) (*AssetVehicle, bool, error) {
	model := AssetVehicle{}
	if has, err := mysql.GetSession().ID(userId).Get(&model); err != nil {
		return nil, has, err
	} else {
		return &model, has, nil
	}
}

func AssetVehicleGetByCode(code string) (*AssetVehicle, bool) {
	model := AssetVehicle{}
	session := mysql.GetSession()
	has, err := session.Where("code =?", code).Get(&model)
	if err != nil {
		log.GetHttpLogLogger().Error(fmt.Sprintf("%v", err))
		return nil, false
	}

	return &model, has
}

// 按照id更新
func AssetVehicleUpdateById(id int, data qmap.QM) (*AssetVehicle, error) {
	model := AssetVehicle{}
	has, err := mysql.GetSession().Get(&model)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("数据不存在")
	}

	if brand, has := data.TryString("brand"); has {
		model.Brand = brand
	}

	if code, has := data.TryString("code"); has {
		model.Code = code
	}

	if detail, has := data.TryString("detail"); has {
		model.Detail = detail
	}

	if _, err := mysql.GetSession().ID(id).Update(&model); err != nil {
		return nil, err
	} else {
		return &model, nil
	}
}

// ---------------- 查询分组列表 -------------
type AssetVehicleSelectListItem struct {
	Brand string                    `json:"brand"`
	Codes []AssetVehicleSelectCodes `json:"codes"`
}

type AssetVehicleSelectCodes struct {
	Id   int    `json:"id"`
	Code string `json:"code"`
}

func AssetVehicleSelectList() []AssetVehicleSelectListItem {
	lists := make([]AssetVehicle, 0)
	err := mysql.GetSession().OrderBy("brand asc").Find(&lists)
	if err != nil {
		panic(err)
	}

	// 构建数据slice， 因为从数据库中查询的数据已经按照brand正序排列，遍历lists，相同brand的放到同一个AssetVehicleSelectListItem对象中
	data := make([]AssetVehicleSelectListItem, 0)
	index := 0
	brand := ""
	for k, v := range lists {
		// 首次直接创建元素
		if k == 0 {
			brand = v.Brand
			data = append(data, AssetVehicleSelectListItem{
				Brand: v.Brand,
				Codes: []AssetVehicleSelectCodes{{Id: v.Id, Code: v.Code}},
			})
			continue
		}

		// 当brand不变时，追加到当前AssetVehicleSelectListItem对象
		if brand == v.Brand {
			data[index].Codes = append(data[index].Codes, AssetVehicleSelectCodes{Id: v.Id, Code: v.Code})
		} else {
			// 当遍历数据发现brand跟上一个不同时，要另存到新的AssetVehicleSelectListItem对象
			index++ // 方便接下来相同brand时直接拿到slice中具体位置
			brand = v.Brand
			data = append(data, AssetVehicleSelectListItem{
				Brand: v.Brand,
				Codes: []AssetVehicleSelectCodes{{Id: v.Id, Code: v.Code}},
			})
		}
	}

	return data
}
