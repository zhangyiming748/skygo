package mysql_model

import "skygo_detection/lib/common_lib/mysql"

// 整车测试进展
type ScreenVehicleTestProgress struct {
	Id        int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	Company   string `xorm:"not null comment('车企名称') VARCHAR(255)" json:"company"`
	Brand     string `xorm:"not null comment('车型名称') VARCHAR(255)" json:"brand"`
	Status    int    `xorm:"not null comment('状态') INT(1)" json:"status"`
	StartTime int    `xorm:"not null comment('开始时间') INT(11)" json:"start_time"`
	EndTime   int    `xorm:"not null comment('结束时间') INT(11)" json:"end_time"`
}

func (this *ScreenVehicleTestProgress) Create() (*ScreenVehicleTestProgress, error) {
	_, err := mysql.GetSession().InsertOne(this)
	return this, err
}

func (this *ScreenVehicleTestProgress) Remove(id int) error {
	_, err := mysql.GetSession().ID(id).Delete(this)
	return err
}
