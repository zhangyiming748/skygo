package mysql_model

import (
	"math/rand"
	"time"

	"xorm.io/xorm"

	"skygo_detection/lib/common_lib/mysql"
)

type SysVehicleFactory struct {
	Id         int    `xorm:"not null pk autoincr comment('主键') INT(10)"`
	Name       string `xorm:"default '' comment('车厂名称') VARCHAR(64)"`
	Type       int    `xorm:"not null default 0 comment('车厂类型(0:车厂, 1:供应商)') TINYINT(4)"`
	ChannelId  string `xorm:"not null comment('渠道号') unique CHAR(6)"`
	Status     int    `xorm:"not null default 1 comment('状态(0:禁用， 1:启用)') TINYINT(4)"`
	UpdateTime int    `xorm:"not null default 0 INT(10)"`
	CreateTime int    `xorm:"created not null default 0 INT(10)"`
}

const (
	FT_VEHICLE  = 0 //车厂
	FT_SUPPLIER = 1 //供应商

	CP_VEHICLE  = "T" //车厂渠道号前缀
	CP_SUPPLIER = "O" //供应商渠道号前缀
)

func (f SysVehicleFactory) GetWhere(params map[string]interface{}) *xorm.Session {
	session := mysql.GetSession()
	if id, ok := params["id"]; ok {
		session.And("id=?", id)
	}
	if sn, ok := params["sn"]; ok && sn != "" {
		session.And("sn=?", sn)
	}
	if channelId, ok := params["channelId"]; ok && channelId != "" {
		session.And("channel_id=?", channelId)
	}
	return session
}

/**
 * 生成一个未被使用过的渠道号
 *
 * @param {int}     type    车厂类型(0:车厂，1:供应商)
 *
 * @return string
 */
func (this *SysVehicleFactory) GenerateChannelId(factoryType int) string {
	newChannelID := ""
	if factoryType == FT_VEHICLE {
		newChannelID = CP_VEHICLE
	} else {
		newChannelID = CP_SUPPLIER
	}

	raw := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < 5; i++ {
		newChannelID += raw[rand.Intn(10)]
	}

	if has, _ := mysql.GetSession().Where("channel_id = ?", newChannelID).Get(this); has {
		return this.GenerateChannelId(factoryType)
	} else {
		return newChannelID
	}
}

func (this SysVehicleFactory) GetChannelName(channelId string) string {
	factory := new(SysVehicleFactory)
	params := map[string]interface{}{"channelId": channelId}
	if has, err := this.GetWhere(params).Get(factory); err == nil && has {
		return factory.Name
	}
	return ""
}
