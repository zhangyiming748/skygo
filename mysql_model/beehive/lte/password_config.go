package lteModel

import "skygo_detection/lib/common_lib/mysql"

type BeehiveLteSystemPasswordConfig struct {
	Id           int    `xorm:"not null pk comment('自增主键id') INT(11)" json:"id"`
	Status       int    `xorm:"not null comment('密码配置状态：1:默认密码 2:上传密码') INT(11)" json:"status"`
	OriginFileId string `xorm:"not null comment('默认的密码文件ID') varchar(20)" json:"origin_file_id"`
	UploadFileId string `xorm:"not null comment('上传的密码文件ID') varchar(20)" json:"upload_file_id"`
	Name         string `xorm:"not null comment('文件名字') varchar(50)" json:"name"`
	CreateTime   string `xorm:"not null comment('创建时间') varchar(20)" json:"create_time"`
	UpdateTime   string `xorm:"not null comment('修改时间') varchar(20)" json:"update_time"`
}

const (
	StatusDefault = 1
	StatusUpload  = 2
)

func (b *BeehiveLteSystemPasswordConfig) Create() (int64, error) {
	return mysql.GetSession().InsertOne(b)
}

func (b *BeehiveLteSystemPasswordConfig) Update(cols ...string) (int64, error) {
	return mysql.GetSession().Table(b).ID(b.Id).Cols(cols...).Update(b)
}

func (b *BeehiveLteSystemPasswordConfig) Find() (bool, error) {
	return mysql.GetSession().Table(b).Get(b)
}

func (b BeehiveLteSystemPasswordConfig) GetOne() (BeehiveLteSystemPasswordConfig, error) {
	_, err := mysql.GetSession().Get(&b)
	if err != nil {
		return b, err
	}
	return b, err
}
