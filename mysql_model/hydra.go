package mysql_model

import (
	"errors"
	"skygo_detection/lib/common_lib/mysql"
)

type HydraTask struct {
	Id            int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	TaskId        int    `xorm:"comment('任务id') INT(11)" json:"task_id"`
	TaskName      string `xorm:"comment('任务名称') VARCHAR(255)" json:"task_name"`
	Address       string `xorm:"comment('服务器地址') VARCHAR(255)" json:"address"`
	Port          string `xorm:"comment('手动输入端口号') VARCHAR(255)" json:"port"`
	Protocol      string `xorm:"comment('协议类型') VARCHAR(255)" json:"protocol"`
	Path          string `xorm:"comment('登录接口地址') VARCHAR(255)" json:"path"`
	Form          string `xorm:"comment('登录表单') VARCHAR(255)" json:"form"`
	UserName      string `xorm:"comment('用户手动输入用户名字典') VARCHAR(255)" json:"user_name"`
	UserNameFile  string `xorm:"comment('用户上传的用户名字典文件') VARCHAR(512)" json:"user_name_file"`
	UserNameType  int    `xorm:"comment('用户名类型 0:默认字典 1:手动录入 2:自定义上传) INT(10)" json:"user_name_type"`
	Passwd        string `xorm:"comment('用户手动输入密码字典') VARCHAR(255)" json:"passwd"`
	PasswdFile    string `xorm:"comment('用户上传的密码字典文件') VARCHAR(512)" json:"passwd_file"`
	PasswdType    int    `xorm:"comment('密码类型 0:默认字典 1:手动录入 2:自定义上传') INT(10)" json:"passwd_type"`
	Sid           string `xorm:"comment('Oracle Sid') VARCHAR(255)" json:"sid"`
	RequestHost   string `xorm:"comment('发出查询请求的host') VARCHAR(255)" json:"request_host"`
	CreateTime    string `xorm:"created comment('创建时间') DATETIME" json:"create_time"`
	UserId        int    `xorm:"comment('创建人') INT(11)" json:"user_id"`
	UpdateTime    string `xorm:"updated comment('更新时间') DATETIME" json:"update_time"`
	DeleteTime    string `xorm:"deleted comment('删除时间') DATETIME" json:"delete_time"`
	Status        int    `xorm:"comment('任务状态 1运行中 2成功 3失败 4已取消') TINYINT(3)" json:"status"`
	OriginResults string `xorm:"comment('原始结果') TEXT" json:"origin_results"`
	Success       int    `xorm:"comment('hydra状态 1成功 2失败') TINYINT(3)" json:"success"`
	Results       string `xorm:"comment('破解结果报告') TEXT" json:"results"`
}

func (this HydraTask) Create() error {
	_, err := mysql.GetSession().InsertOne(this)
	if err != nil {
		return err
	}
	return nil
}
func (this HydraTask) UpdateByTaskId(tid int) error {
	_, err := mysql.GetSession().Where("task_id = ?", tid).Update(this)
	if err != nil {
		return err
	}
	return nil
}
func (this HydraTask) DeleteByTaskId(tid int) (int64, error) {
	return mysql.GetSession().Where("task_id = ?", tid).Unscoped().Delete(&this)
}
func (this HydraTask) DeleteByTaskIds(ids []int) (int64, error) {
	return mysql.GetSession().ID(ids).Delete(&this)
}
func (this HydraTask) FindByTaskId(tid int) (HydraTask, error) {
	has, err := mysql.GetSession().Where("task_id = ?", tid).Get(&this)
	if !has {
		err = errors.New("not found")
	}
	if err != nil {
		return HydraTask{}, err
	}
	return this, nil
}
