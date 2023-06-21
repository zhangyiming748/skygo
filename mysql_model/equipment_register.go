package mysql_model

import (
	"errors"
	"fmt"
	"math/rand"
	"strconv"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/service"
)

type EquipmentRegister struct {
	Id            int    `xorm:"not null pk autoincr comment('id') INT(11)"`
	ChannelId     string `xorm:"not null default '' comment('渠道号') CHAR(6)"`
	EquipmentType string `xorm:"ENUM('tbox','vehicle')"`
	Udid          string `xorm:"not null default '' comment('设备唯一识别号') CHAR(6)"`
	Sn            string `xorm:"default '' comment('sn号') VARCHAR(256)"`
	CryptKey      string `xorm:"default '' comment('对称加密密钥') VARCHAR(256)"`
	RegisterNum   int    `xorm:"not null default 1 comment('注册次数') INT(11)"`
	CreateTime    int    `xorm:"created not null default 0 comment('创建时间') INT(11)"`
}

func (this *EquipmentRegister) GetOne(sn string) (*EquipmentRegister, bool) {
	if has, _ := mysql.GetSession().Where("sn = ?", sn).Get(this); has {
		return this, has
	} else {
		return nil, has
	}
}

func EquipmentRegisterInsert(channelId, equipmentType, udid string) (sn, registerKey string) {
	register := new(EquipmentRegister)
	session := mysql.GetSession().Where("channel_id = ?", channelId)
	session = session.And("equipment_type = ?", equipmentType)
	session = session.And("udid = ?", udid)

	has, _ := session.Get(register)
	register.CryptKey = new(service.Crypt_service).GenerateRegisterKey()
	register.RegisterNum = register.RegisterNum + 1
	if has {
		//TBox设备不允许重复注册
		//if equipmentType != common.EQUIPMENT_TBOX {
		if _, err := mysql.GetSession().Table(new(EquipmentRegister)).ID(register.Id).Update(custom_util.StructToMap(*register)); err == nil {
			EquipmentRegisterLogInsert(register)
		} else {
			panic(err)
		}
		//} else {
		//	panic(custom_error.ReRegisterError)
		//}
	} else {
		register.ChannelId = channelId
		register.EquipmentType = equipmentType
		register.Udid = udid
		register.Sn = GenerateSeriesNumber(channelId, equipmentType, udid)
		if _, err := mysql.GetSession().InsertOne(register); err == nil {
			EquipmentRegisterLogInsert(register)
		} else {
			panic(err)
		}
	}
	return register.Sn, register.CryptKey
}

func (this EquipmentRegister) GetKey(sn string) string {
	register := new(EquipmentRegister)
	if has, _ := mysql.GetSession().Where("sn = ?", sn).Get(register); has {
		return register.CryptKey
	} else {
		panic(errors.New("UnRegisterError")) // TODO
	}
}

/*
 * 根据渠道号，设备类型、udid生成唯一序列号
 * 序列号规则:
 * XXXXX XX XXXXXXXXXXXXXXXXXXXX
 * 1~5:渠道号
 * 6~7:设备类型(v:车机，t:tbox)
 * 8~22:唯一id
 */
func GenerateSeriesNumber(channelId, equipmentType, udid string) string {
	channel := getFactorySN(channelId)
	equpment := getEquipmentSN(equipmentType)
	uid := getUniqueID(udid)
	return fmt.Sprintf("%s%s%s", channel, equpment, uid)
}

func getEquipmentSN(equipmentType string) string {
	switch equipmentType {
	case common.ACCOUNT_TYPE_TBOX:
		return "TB"
	case common.ACCOUNT_TYPE_VEHICLE:
		return "VE"
	default:
		panic("unknown equipment type")
	}
}

func getFactorySN(channelId string) string {
	factory := new(SysVehicleFactory)
	if has, _ := mysql.GetSession().Where("channel_id", channelId).Get(factory); has {
		channelType := channelId[0:1]
		channelNumber := channelId[1:]
		if i, err := strconv.Atoi(channelNumber); err == nil {
			return channelType + decimalToThirtySix(i)
		} else {
			panic(err)
		}
	} else {
		panic("unknown channel id")
	}
}

// 十进制转36进制,返回四个字符长度的字符串(转换的最大整形为一百万)
func decimalToThirtySix(raw int) string {
	encoding := []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	encoded := []rune("0000")
	i := 0
	for {
		if raw <= 0 {
			break
		}
		if i > 3 {
			panic("raw id is too large")
		}

		index := raw % 36
		encoded[i] = encoding[index]
		raw = raw / 36
		i++
	}
	return string(encoded)
}

func getUniqueID(udid string) string {
	rand.Seed(stringToInt64(udid))
	return fmt.Sprintf("%v%v", custom_util.GetCurrentMilliSecond(), rand.Intn(100))
}

func stringToInt64(udid string) int64 {
	var res int64 = 0
	var b = 0
	for _, char := range []rune(udid) {
		if b > 18 {
			break
		} else {
			b++
		}
		res = res*10 + int64(char)/10
	}
	return res
}
