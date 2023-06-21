package mysql_model

import (
	"fmt"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
)

type BeehiveGsmSystemSms struct {
	Id         int    `xorm:"not null pk autoincr comment('主键id') INT(10)" json:"id"`
	TaskId     int    `xorm:"not null comment('任务id') INT(10)" json:"task_id"`
	Time       string `xorm:"comment('收件时间') VARCHAR(255)" json:"time"`
	RecvMobile string `xorm:"comment('收件手机号') VARCHAR(255)" json:"recv_mobile"`
	RecvImsi   string `xorm:"comment('收件imsi') VARCHAR(255)" json:"recv_imsi"`
	SendMobile string `xorm:"comment('发件手机号') VARCHAR(255)" json:"send_mobile"`
	SendImsi   string `xorm:"comment('发件imsi') VARCHAR(255)" json:"send_imsi"`
	SmsContent string `xorm:"comment('短信内容') VARCHAR(255)" json:"sms_content"`
	CreateTime string `xorm:"created not null comment('创建时间') DATETIME" json:"create_time"`
	UpdateTime string `xorm:"created not null comment('更新时间') DATETIME" json:"update_time"`
	DeleteTime string `xorm:"deleted comment('删除时间') DATETIME" json:"delete_time"`
}

// 软删除
func (this BeehiveGsmSystemSms) DeleteSMS(ids []int) int {
	fail := 0
	for _, id := range ids {
		_, err := mysql.GetSession().ID(id).Delete(&this)
		if err != nil {
			fail++
		}
	}
	return fail
}

// 硬删除
func (this BeehiveGsmSystemSms) RealDelete() (int64, error) {
	return mysql.GetSession().Unscoped().Delete(&this)
}
func (this BeehiveGsmSystemSms) Delete(ids []int) (int64, error) {
	return mysql.GetSession().ID(ids).Unscoped().Delete(&this)
}

// 获取全部短信
func GetAllSms(tid int) map[string]interface{} {
	all := mysql.GetSession().
		Table("beehive_gsm_system_sms").
		Where("task_id = ?", tid)

	widget := orm.PWidget{}
	widget.AddSorter(*(orm.NewSorter("time", 1)))
	res := widget.PaginatorFind(all, &[]BeehiveGsmSystemSms{})
	return res
}

// 短信数量角标
func (this BeehiveGsmSystemSms) GsmSystemGetSMSNum() (int, error) {
	count, err := mysql.GetSession().
		Table("beehive_gsm_system_sms").
		Where("task_id = ?", this.TaskId).
		Count()
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return int(count), nil
}

//搜索短信

func GsmSystemSearchSms(tid int, key, q string) map[string]interface{} {
	s := mysql.GetSession().
		Table("beehive_gsm_system_sms").Where("task_id = ?", tid)
	if key != "" {
		s = s.Where("recv_mobile like ? or send_mobile like ? or sms_content like ?", "%"+key+"%", "%"+key+"%", "%"+key+"%")
	}
	//Table("beehive_gsm_system_sms ").
	//Where("task_id = ?", tid).
	//Or("recv_mobile like ? ", "%"+key+"%").
	//Or("send_mobile like ?", "%"+key+"%").
	//Or("sms_content like ?", "%"+key+"%")

	widget := orm.PWidget{}
	widget.SetQueryStr(q)
	widget.AddSorter(*(orm.NewSorter("time", 1)))
	all := widget.PaginatorFind(s, &[]BeehiveGsmSystemSms{})
	return all
}

func (this BeehiveGsmSystemSms) GsmSystemSaveSms() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}

//func (this []BeehiveGsmSystemSms)GsmSystemSaveAll()(int64,error){
//	mysql.GetSession().Insert(this)
//}
