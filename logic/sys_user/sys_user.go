package sys_user

import (
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/logic/ump"
	"skygo_detection/mysql_model"
	"skygo_detection/service"
)

// 统一登录创建用户
func CreateUserByUmp(userInfo ump.UserInfo, password string) (res int64, err error) {
	newUser := mysql_model.SysUser{
		Username:    userInfo.UserName,
		ChannelId:   "Q00001",
		RoleId:      1,
		Realname:    userInfo.UserName,
		Password:    service.HashPassword(password),
		AccountType: "user",
		Nickname:    userInfo.UserName,
		Email:       userInfo.Email,
		Status:      mysql_model.USER_STATUS_NORMAL,
	}
	s := mysql.GetSession().Table(mysql_model.SysUser{})
	res, err = s.InsertOne(newUser)
	if err != nil {
		return
	}
	return
}
