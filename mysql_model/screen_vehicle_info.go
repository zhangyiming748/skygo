package mysql_model

import "skygo_detection/lib/common_lib/mysql"

// 车辆上的信息
type ScreenVehicleInfo struct {
	Id       int     `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	Name     string  `xorm:"not null comment('名称') VARCHAR(255)" json:"name"`
	Standard string  `xorm:"not null comment('标准') VARCHAR(255)" json:"standard"`
	TaskRate float64 `xorm:"not null comment('任务进度') FLOAT(11)" json:"task_rate"`
	PassRate float64 `xorm:"not null comment('符合度') FLOAT(11)" json:"pass_rate"`
}

func (this *ScreenVehicleInfo) Create() (*ScreenVehicleInfo, error) {
	_, err := mysql.GetSession().InsertOne(this)
	return this, err
}

func (this *ScreenVehicleInfo) Remove(id int) error {
	_, err := mysql.GetSession().ID(id).Delete(this)
	return err
}
