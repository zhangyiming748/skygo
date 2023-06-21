package mysql_model

import (
	"errors"

	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mysql"
)

const CRYPTTYPE_REGISTER = "1" //注册密钥加密
const CRYPTTYPE_SESSION = "2"  //会话密钥加密

type EquipmentSession struct {
	Id         int    `xorm:"not null pk autoincr comment('id') INT(11)"`
	Sn         string `xorm:"default '' comment('sn号') index(idx_channelid_sn) VARCHAR(256)"`
	SessionKey string `xorm:"default '' comment('session key') VARCHAR(256)"`
	UpdateTime int    `xorm:"updated not null default 0 comment('更新时间') INT(11)"`
	CreateTime int    `xorm:"created not null default 0 comment('创建时间') INT(11)"`
}

func (this *EquipmentSession) Update(sn, sessionKey string) string {
	has, _ := mysql.GetSession().Where("sn = ?", sn).Get(this)
	this.SessionKey = sessionKey
	if !has {
		this.Sn = sn
		if _, err := mysql.GetSession().InsertOne(this); err != nil {
			panic(err)
		}
	} else {
		//两个坑:
		// ①必须要用map的参数形式传递，调用Update()方法才能更新字段中的0值，否则会被忽略；
		// ②结构体转化为map时，参数不能传递指针，需要先解构一下
		if _, err := mysql.GetSession().Table(new(EquipmentSession)).ID(this.Id).Update(custom_util.StructToMap(*this)); err != nil {
			panic(err)
		}
	}
	return this.SessionKey
}

func (this *EquipmentSession) GetKey(sn string) string {
	if has, _ := mysql.GetSession().Where("sn = ?", sn).Get(this); has && len(this.SessionKey) == 32 {
		return this.SessionKey
	} else {
		panic(errors.New("SessionKeyError"))
	}
}
