package ump_model

import "skygo_detection/lib/common_lib/mysql"

type UmpUser struct {
	Id          int    `xorm:"not null pk autoincr INT(11)" json:"id"`
	UserId      int    `xorm:"not null default 0 comment('用户id') INT(11)" json:"user_id"`
	UmpUserId   int    `xorm:"not null default 0 comment('UMP用户id') INT(11)" json:"ump_user_id"`
	UserName    string `xorm:"comment('用户名') unique VARCHAR(64)" json:"user_name"`
	Email       string `xorm:"comment('用户邮箱') unique VARCHAR(200)" json:"email"`
	FirstName   string `xorm:"comment('用户姓') unique VARCHAR(64)" json:"first_name"`
	LastName    string `xorm:"comment('用户名') unique VARCHAR(64)" json:"last_name"`
	IsStaff     int    `xorm:"not null default 0 comment('是否员工(1:否 2:是)') TINYINT(1)" json:"is_staff"`
	IsSuperuser int    `xorm:"not null default 0 comment('是否超级用户(1:否 2:是)') TINYINT(1)" json:"is_superuser"`
	IsActive    int    `xorm:"not null default 0 comment('是否激活(1:否 2:是)') TINYINT(1)" json:"is_active"`
	CreateTime  string `xorm:"not null comment('创建时间') DATETIME" json:"create_time"`
	UpdateTime  string `xorm:"not null comment('更新时间') DATETIME" json:"update_time"`
}

const (
	IsNotStaff     = 1
	IsStaff        = 2
	IsNotActive    = 1
	IsActive       = 2
	IsNotSuperUser = 1
	IsSuperUser    = 2
)

func (this *UmpUser) Create() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}

func (this *UmpUser) Get(userId int) (bool, error) {
	return mysql.GetSession().Where("user_id = ?", userId).Get(this)
}

func (this UmpUser) GetUmpInfo(umpUserId int) (UmpUser, error) {
	_, err := mysql.GetSession().Where("ump_user_id = ?", umpUserId).Get(&this)
	if err != nil {
		return UmpUser{}, nil
	}
	return this, nil
}
