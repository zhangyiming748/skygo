package mysql_model

import (
	"skygo_detection/lib/common_lib/mysql"
)

type PrivacyAppPermissionMap struct {
	Id            int    `xorm:"not null pk comment('权限id') INT(10)" json:"id"`
	GroupId       int    `xorm:"not null comment('组权限id') INT(10)" json:"group_id"`
	FullName      string `xorm:"not null comment('权限全称') VARCHAR(255)" json:"full_name"`
	PartName      string `xorm:"not null comment('权限部分全称') VARCHAR(255)" json:"part_name"`
	Name          string `xorm:"not null comment('权限中文名称') VARCHAR(255)" json:"name"`
	GroupFullName string `xorm:"not null comment('组权限全称') VARCHAR(255)" json:"group_full_name"`
	GroupName     string `xorm:"not null comment('组权限部分全称') VARCHAR(255)" json:"group_name"`
}

func FindNameByPermissionId(pid int) (string, error) {
	model := new(PrivacyAppPermissionMap)
	has, err := mysql.GetSession().ID(pid).Get(model)
	if !has {
		return "not found", err
	}
	if err != nil {
		return "something wrong", err
	}
	return model.Name, nil
}
