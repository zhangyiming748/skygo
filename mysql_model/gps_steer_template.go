package mysql_model

import (
	"skygo_detection/lib/common_lib/mysql"
)

type GpsSteerTemplate struct {
	Id                      int     `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	Name                    string  `xorm:"not null comment('行驶模板名称') VARCHAR(50)" json:"name"`
	MaxLatAcc               float32 `xorm:"not null comment('最大横向加速度') FLOAT(255)" json:"max_lat_acc"`
	MaxLongAcc              float32 `xorm:"not null comment('最大纵向加速度') FLOAT(255)" json:"max_long_acc"`
	MaxJerk                 int     `xorm:"not null comment('急动度') int(11)" json:"max_jerk"`
	MaxSpeed                float32 `xorm:"not null comment('最大速度') FLOAT(255)" json:"max_speed"`
	StationaryPeriod        float32 `xorm:"not null comment('起步时间') FLOAT(255)" json:"stationary_period"`
	StationaryPeriodEnd     float32 `xorm:"not null comment('停车时间') FLOAT(255)" json:"stationary_period_end"`
	PositionSmoothingFactor int     `xorm:"not null comment('移动平滑系数') INT(11)" json:"position_smoothing_factor"`
	SpeedSmoothingFactor    int     `xorm:"not null comment('速度抖动系数') INT(11)" json:"speed_smoothing_factor"`
	Creator                 string  `xorm:"not null comment('创建人') VARCHAR(20)" json:"creator"`
	CreateTime              string  `xorm:"created not null comment('创建时间') DATETIME" json:"create_time"`
	UpdateTime              string  `xorm:"created not null comment('修改时间') DATETIME" json:"update_time"`
	Type                    int     `xorm:"comment('模板类型') INT(11)" json:"type"`
	FileId                  string  `xorm:"comment('mongo 文件id') VARCHAR(100)" json:"file_id"`
	PictureName             string  `xorm:"comment('文件name') VARCHAR(255)" json:"file_name"`
	TemplateState           int     `xorm:"comment('1为系统模板不能删除，2为新增模板可以删除') INT(11)" json:"template_state"`
	Sn                      int     `xorm:"comment('排序') INT(3)" json:"sn"`
}

func (this *GpsSteerTemplate) Create() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}

func (this *GpsSteerTemplate) Remove(id int) (int64, error) {
	return mysql.GetSession().ID(id).Delete(this)
}

func (this *GpsSteerTemplate) Update(id interface{}, cols ...string) (int64, error) {
	return mysql.GetSession().ID(id).Cols(cols...).Update(this)
}

func (this *GpsSteerTemplate) FindById(id int) (GpsSteerTemplate, bool) {
	model := GpsSteerTemplate{}
	has, err := mysql.FindById(id, &model)
	if err != nil {
		panic(err)
	}
	return model, has
}

func GetAllTemplate() []GpsSteerTemplate {
	arr := make([]GpsSteerTemplate, 0)
	mysql.GetSession().Find(&arr)
	return arr
}

func (this *GpsSteerTemplate) GetOne(id int) (bool, error) {
	bool, err := mysql.GetSession().Where("id=?", id).Get(this)
	return bool, err
}
