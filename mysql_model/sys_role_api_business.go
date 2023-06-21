package mysql_model

import (
	"xorm.io/xorm"

	"skygo_detection/lib/common_lib/mysql"
)

type SysRoleApiBusiness struct {
	Id     int `xorm:"not null pk autoincr comment('编号') INT(11)"`
	RoleId int `xorm:"not null default 0 comment('角色ID (与sys_role表id关联)') unique(role_access) INT(11)"`
	ApiId  int `xorm:"not null default 0 comment('访问接口ID (与sys_api表的id关联)') unique(role_access) INT(11)"`
}

func ApiRoleBusinessGetWhere(params map[string]interface{}) *xorm.Session {
	session := mysql.GetSession()
	if id, ok := params["id"]; ok {
		session.And("id=?", id)
	}
	if roleId, ok := params["roleId"]; ok {
		session.And("role_id=?", roleId)
	}
	if apiId, ok := params["apiId"]; ok {
		session.And("api_id=?", apiId)
	}

	return session
}

func (this SysRoleApiBusiness) CheckPrivilege(roleId int64, url, method string) bool {
	apiId := ApiFetchIdByUrl(url, method)
	roleApi := new(SysRoleApi)
	params := map[string]interface{}{
		"roleId": roleId,
		"apiId":  apiId,
	}
	if has, err := ApiRoleGetWhere(params).Get(roleApi); err == nil && has {
		return true
	}
	return false
}

func (this SysRoleApiBusiness) GetRoleApi(roleId int) []int {
	roleApi := []int{}
	if err := mysql.GetSession().Where("role_id = ?", roleId).Table(this).Cols("api_id").Find(&roleApi); err == nil {
		return roleApi
	} else {
		panic(err)
	}
}

func (this SysRoleApiBusiness) UpdateRoleApi(roleId int, privileges []int) error {
	{
		mysql.GetSession().Where("role_id = ?", roleId).Delete(this)
	}

	{
		relations := []SysRoleApi{}
		for _, apiId := range privileges {
			temp := SysRoleApi{
				RoleId: roleId,
				ApiId:  apiId,
			}
			relations = append(relations, temp)
		}
		_, err := mysql.GetSession().Insert(relations)
		return err
	}
}
