package mysql_model

import (
	"xorm.io/builder"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mysql"
)

type SysSaasRole struct {
	Id        int    `xorm:"not null pk autoincr INT(10)"`
	ChannelId string `xorm:"not null default '' comment('渠道号') CHAR(6)"`
	Name      string `xorm:"default '' comment('角色名称') VARCHAR(255)"`
	ParentId  int    `xorm:"not null default 0 comment('父角色id') INT(10)"`
}

// 查询子角色id列表
func GetSubRoleIds(roleId int) []int {
	if roleId <= 0 {
		return []int{}
	}
	//getSubRoles := func(arg ...interface{}) interface{} {
	//	maxRoleDeep := 15 //最大用户角色查询深度
	//	queryRoles := []int{roleId}
	//	subRoles := []int{roleId}
	//	for maxRoleDeep > 0 && len(queryRoles) > 0 {
	//		params := qmap.QM{
	//			"in_parent_id": queryRoles,
	//		}
	//		roles := []int{}
	//		if err := sys_service.NewSessionWithCond(params).Table("sys_saas_role").Cols("id").Find(&roles); err == nil {
	//			queryRoles = roles
	//		} else {
	//			panic(err)
	//		}
	//		subRoles = append(subRoles, roles...)
	//		maxRoleDeep--
	//	}
	//
	//	return qmap.QM{"sub_roles": subRoles}
	//}
	//data := service.CacheQMDefault(fmt.Sprintf("mservice_auth:sub_roles:%d", roleId), getSubRoles)

	maxRoleDeep := 15 //最大用户角色查询深度
	queryRoles := []int{roleId}
	subRoles := []int{roleId}
	if roleId == common.SUPER_ADMINISTRATE_ROLE_ID {
		subRoles = append(subRoles, 0)
	}
	for maxRoleDeep > 0 && len(queryRoles) > 0 {
		var roles []int
		if err := mysql.GetSession().Where(builder.In("parent_id", queryRoles)).Table("sys_saas_role").Cols("id").Find(&roles); err == nil {
			queryRoles = roles
		} else {
			panic(err)
		}
		subRoles = append(subRoles, roles...)
		maxRoleDeep--
	}
	return subRoles
}

/*
判断某个角色是否有管理另一个角色的权限
*/
func HasRoleManagePrivilege(ownerRoleId, targetRoleId int) bool {
	// 如果是超级管理员或者角色相同直接返回true
	if ownerRoleId == common.SUPER_ADMINISTRATE_ROLE_ID || ownerRoleId == targetRoleId {
		return true
	} else {
		ownerSubRoleIds := GetSubRoleIds(ownerRoleId)
		for _, subRoleId := range ownerSubRoleIds {
			if subRoleId == targetRoleId {
				return true
			}
		}
	}
	return false
}

func (this *SysSaasRole) GetSaasRoleName(id int) string {
	if has, err := mysql.GetSession().ID(id).Get(this); err != nil || !has {
		return ""
	} else {
		return this.Name
	}
}
