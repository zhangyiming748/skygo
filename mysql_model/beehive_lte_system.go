package mysql_model

import (
	"skygo_detection/lib/common_lib/mysql"
)

type BeehiveLteSystem struct {
	Id         int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	TaskId     int    `xorm:"not null comment('任务id') INT(11)" json:"task_id"`
	Status     int    `xorm:"not null default '0' comment('系统状态：1:运行中 2:已关闭') INT(3)" json:"status"`
	Apn        string `xorm:"comment('终端设备apn') VARCHAR(20)" json:"apn"`
	Imsi       string `xorm:"default '' comment('终端设备imsi') VARCHAR(20)" json:"imsi"`
	Ip         string `xorm:"not null default 0 comment('终端设备ip') VARCHAR(50)" json:"ip"`
	UserName   string `xorm:"comment('密码破解用户名') VARCHAR(50)" json:"user_name"`
	Password   string `xorm:"comment('密码') VARCHAR(50)" json:"password"`
	CreateTime string `xorm:"created not null comment('创建时间') DATETIME" json:"create_time"`
	UpdateTime string `xorm:"created not null comment('修改时间') DATETIME" json:"update_time"`
}

func (this *BeehiveLteSystem) Create() error {
	// 创建场景数据
	session := mysql.GetSession()
	_, err := session.InsertOne(this)
	if err != nil {
		return err
	}
	return nil
}
func (this *BeehiveLteSystem) Get(taskId int, imsi string) int {
	models := make([]BeehiveLteSystem, 0)
	session := mysql.GetSession()
	session.Where("task_id = ?", taskId)
	session.And("imsi = ?", imsi)
	session.Find(&models)
	var id = 0
	for _, v := range models {
		id = v.Id
	}
	return id
}

func (this *BeehiveLteSystem) Update(cols ...string) (int64, error) {
	return mysql.GetSession().Table(this).ID(this.Id).Cols(cols...).Update(this)
}

func (b *BeehiveLteSystem) FindByTaskId(taskId int) (bool, error) {
	return mysql.GetSession().Table(b).Where("task_id=?", taskId).Get(b)
}

func (this BeehiveLteSystem) ForceUpdateSystemStatus() (int64, error) {
	s := mysql.GetSession()
	i, err := s.Where("task_id = ?", this.TaskId).
		Cols("status").Update(this)
	if err != nil {
		return 0, err
	}
	return i, nil
}
