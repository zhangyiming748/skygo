package mysql_model

import (
	"errors"
	"skygo_detection/lib/common_lib/mysql"
)

type BeehiveGsmSystemTty struct {
	Id         int    `xorm:"not null pk autoincr comment('主键id') INT(10)" json:"id"`
	TaskId     int    `xorm:"not null comment('任务id') INT(10)" json:"task_id"`
	Imsi       string `xorm:"not null comment('获取到的imsi') VARCHAR(255)" json:"imsi"`
	Imei       string `xorm:"not null comment('获取到的imei') VARCHAR(255)" json:"imei"`
	Mobile     string `xorm:"comment('获取到的手机号') VARCHAR(255)" json:"mobile"`
	CreateTime string `xorm:"created not null comment('创建时间') DATETIME" json:"create_time"`
}

func (this BeehiveGsmSystemTty) UpdateTty() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}
func (this BeehiveGsmSystemTty) DeleteByTaskId() (int64, error) {
	return mysql.GetSession().Where("task_id = ?", this.TaskId).Delete(this)
}
func (this BeehiveGsmSystemTty) GetTtyList() ([]BeehiveGsmSystemTty, error) {
	models := make([]BeehiveGsmSystemTty, 0)
	err := mysql.GetSession().Where("task_id = ?", this.TaskId).Find(&models)
	if err != nil {
		return []BeehiveGsmSystemTty{}, err
	}
	return models, nil
}
func (this BeehiveGsmSystemTty) FindMobileByImsi() (string, error) {
	model := new(BeehiveGsmSystemTty)
	has, err := mysql.GetSession().Where("mobile = ?", this.Mobile).Get(&model)
	if err != nil {
		return "", err
	}
	if !has {
		return "", errors.New("not found")
	}
	return model.Mobile, nil
}
