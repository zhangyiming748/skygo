package mysql_model

import "skygo_detection/lib/common_lib/mysql"

// 合规测试用例分布
type ScreenTestCase struct {
	Id     int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	Name   string `xorm:"not null comment('任务名称') VARCHAR(255)" json:"name"`
	Number int    `xorm:"not null comment('个数') INT(11)" json:"number"`
}

func (this *ScreenTestCase) Create() (*ScreenTestCase, error) {
	_, err := mysql.GetSession().InsertOne(this)
	return this, err
}

func (this *ScreenTestCase) Remove(id int) error {
	_, err := mysql.GetSession().ID(id).Delete(this)
	return err
}
