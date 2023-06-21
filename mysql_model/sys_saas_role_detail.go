package mysql_model

import (
	"skygo_detection/lib/common_lib/mysql"
)

type SysSaasRoleDetail struct {
	Id              int    `xorm:"not null pk autoincr INT(10)" json:"id"`
	SaasRoleId      int    `xorm:"not null comment('系统角色id') INT(10)" json:"saas_role_id"`
	ServiceRoleId   int    `xorm:"comment('微服务角色id') INT(10)" json:"service_role_id"`
	ServiceRoleName string `xorm:"comment('微服务角色名称') VARCHAR(255)" json:"service_role_name"`
	Service         string `xorm:"comment('微服务') VARCHAR(255)" json:"service"`
	ServiceName     string `xorm:"comment('微服务名称') VARCHAR(255)" json:"service_name"`
	OpId            int    `xorm:"comment('操作用户id') INT(10)" json:"op_id"`
	UpdateTime      int    `xorm:"updated comment('更新时间') INT(10)" json:"update_time"`
	CreateTime      int    `xorm:"created comment('创建时间') INT(10)" json:"create_time"`
}

func (this *SysSaasRoleDetail) GetMServiceRoleId(saasRoleId int64, service string) int64 {
	has, _ := mysql.GetSession().Where("saas_role_id = ?", saasRoleId).And("service = ?", service).Get(this)
	if has {
		return int64(this.ServiceRoleId)
	}
	return 0
}

// 查询有某个服务、某个角色的所有系统角色id
func (this *SysSaasRoleDetail) GetServiceSaasRoleIds(service string, serviceRoleId int) []int {
	saasRoleIds := []int{}

	session := mysql.GetSession().Where("service = ?", service)
	if serviceRoleId > 0 {
		session.And("service_role_id = ?", serviceRoleId)
	}

	session.Table("sys_saas_role_detail").Select("distinct(saas_role_id)").Find(&saasRoleIds)
	return saasRoleIds
}
