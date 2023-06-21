package mysql_model

import "skygo_detection/lib/common_lib/mysql"

type BeehiveLteSystemPackage struct {
	Id         int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	TaskId     int    `xorm:"not null comment('任务id') INT(11)" json:"task_id"`
	FileId     string `xorm:"not null comment('mongo file_id') VARCHAR(255)" json:"file_id"`
	Name       string `xorm:"comment('包文件名字') VARCHAR(255)" json:"name"`
	Size       string `xorm:"comment('包文件大小') VARCHAR(20)" json:"size"`
	CreateTime string `xorm:"created not null comment('创建时间') DATETIME" json:"create_time"`
	UpdateTime string `xorm:"created not null comment('修改时间') DATETIME" json:"update_time"`
}

func (this *BeehiveLteSystemPackage) RemoveById(id int) (int64, error) {
	return mysql.GetSession().ID(id).Delete(this)
}

func (this *BeehiveLteSystemPackage) Create() error {
	// 创建场景数据
	session := mysql.GetSession()
	_, err := session.InsertOne(this)
	if err != nil {
		return err
	}
	return nil
}
func GetAll() ([]BeehiveLteSystemPackage, error) {
	model := make([]BeehiveLteSystemPackage, 0)
	err := mysql.GetSession().Find(&model)
	if err != nil {
		return nil, err
	}
	return model, nil
}
func Get(id int) (string, error) {
	model := make([]BeehiveLteSystemPackage, 0)
	err := mysql.GetSession().Where("id = ?", id).Find(&model)
	if err != nil {
		return "", err
	}
	fileId := ""
	for _, v := range model {
		fileId = v.FileId
	}
	return fileId, err
}
