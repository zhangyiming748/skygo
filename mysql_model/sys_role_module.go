package mysql_model

import (
	"skygo_detection/lib/common_lib/mysql"
)

type SysRoleModule struct {
	Id       int `xorm:"not null pk autoincr comment('编号') INT(11)"`
	RoleId   int `xorm:"not null default 0 comment('角色ID (与sys_role表id关联)') INT(11)"`
	ModuleId int `xorm:"not null default 0 comment('模块ID (与sys_module表的id关联)') INT(11)"`
}

func (this *SysRoleModule) GetModuleIds(roleId int) []int {
	var moduleIds []int
	if err := mysql.GetSession().And("role_id=?", roleId).Table("sys_role_module").Cols("module_id").Find(&moduleIds); err == nil {
		return moduleIds
	} else {
		panic(err)
	}
}

func (this *SysRoleModule) GetRoleModule(roleId int) []int {
	roleModule := []int{}
	if err := mysql.GetSession().Where("role_id = ?", roleId).Table(this).Cols("module_id").Find(&roleModule); err == nil {
		return roleModule
	} else {
		panic(err)
	}
}

func (this *SysRoleModule) UpdateRoleModule(roleId int, modules []int) error {

	mysql.GetSession().Where("role_id = ?", roleId).Delete(this)

	{
		relations := []SysRoleModule{}
		for _, moduleId := range modules {
			temp := SysRoleModule{
				RoleId:   roleId,
				ModuleId: moduleId,
			}
			relations = append(relations, temp)
		}
		_, err := mysql.GetSession().Insert(relations)
		return err
	}
}
