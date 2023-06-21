package mysql_model

import (
	"skygo_detection/lib/common_lib/mysql"
)

type EquipmentRegisterLog struct {
	Id            int    `xorm:"not null pk autoincr comment('id') INT(11)"`
	ChannelId     string `xorm:"not null default '' comment('渠道号') CHAR(6)"`
	Sn            string `xorm:"default '' comment('sn号') VARCHAR(256)"`
	EquipmentType string `xorm:"ENUM('tbox','vehicle')"`
	Udid          string `xorm:"not null default '' comment('设备唯一识别号') CHAR(6)"`
	CryptKey      string `xorm:"not null VARCHAR(256)"`
	CreateTime    int    `xorm:"created not null default '0' comment('创建时间') int"`
}

func EquipmentRegisterLogInsert(register *EquipmentRegister) {
	registerLog := new(EquipmentRegisterLog)
	registerLog.Sn = register.Sn
	registerLog.ChannelId = register.ChannelId
	registerLog.EquipmentType = register.EquipmentType
	registerLog.Udid = register.Udid
	registerLog.CryptKey = register.CryptKey
	if _, err := mysql.GetSession().InsertOne(registerLog); err != nil {
		panic(err)
	}
}
