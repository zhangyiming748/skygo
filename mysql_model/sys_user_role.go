package mysql_model

import (
	"skygo_detection/guardian/src/net/qmap"
	"xorm.io/xorm"

	"skygo_detection/lib/common_lib/mysql"
)

type SysUserRole struct {
	Id          int    `xorm:"not null pk autoincr comment('主键') INT(11)"`
	UserId      int    `xorm:"not null comment('用户id') INT(11)"`
	RoleId      int    `xorm:"not null default 0 comment('所属服务角色id') INT(11)"`
	RoleName    string `xorm:"comment('所属服务角色名称') VARCHAR(255)"`
	Service     string `xorm:"comment('所属服务') VARCHAR(255)"`
	ServiceName string `xorm:"comment('所属服务名称') VARCHAR(255)"`
	OpId        int    `xorm:"not null comment('操作用户id') INT(11)"`
	UpdateTime  int    `xorm:"updated not null comment('更新时间') INT(11)"`
	CreateTime  int    `xorm:"created not null comment('创建时间') INT(11)"`
}

func UserRoleGetWhere(params map[string]interface{}) *xorm.Session {
	session := mysql.GetSession()
	if id, ok := params["id"]; ok {
		session.And("id=?", id)
	}
	if roleId, ok := params["roleId"]; ok {
		session.And("role_id=?", roleId)
	}
	if service, ok := params["service"]; ok {
		session.And("service=?", service)
	}
	if userId, ok := params["userId"]; ok {
		session.And("user_id=?", userId)
	}

	return session
}

//func UserRoleGetRoleId(userId int64, service string) int64 {
//	userRole := new(SysUserRole)
//	params := map[string]interface{}{"userId": userId, "service": service}
//	if has, err := UserRoleGetWhere(params).Get(userRole); err == nil {
//		if has {
//			return int64(userRole.RoleId)
//		}
//	}
//	return 0
//}

func (this *SysUserRole) GetSpecifiedServiceUserIds(service string, roleId int) []int {
	params := qmap.QM{
		"e_service": service,
	}
	if roleId > 0 {
		params["e_role_id"] = roleId
	}
	userIds := []int{}
	if err := mysql.GetSession().Table(this).Cols("user_id").Find(&userIds); err == nil {
		return userIds
	} else {
		panic(err)
	}
}
