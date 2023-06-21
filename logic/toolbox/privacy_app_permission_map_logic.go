package toolbox

import (
	"skygo_detection/common"
	"skygo_detection/mysql_model"
)

// 权限id更换为对应权限中文名
func Permission2PermissionZh(p []P2V) []P2V {
	var np []P2V
	for _, v := range p {
		var nnp P2V
		pid := v.Permission
		nnp.PermissionZH, _ = mysql_model.FindNameByPermissionId(pid)
		nnp.AppName = v.AppName
		nnp.VersionList = v.VersionList
		np = append(np, nnp)
	}
	return np
}

// 权限id更换为对应权限中文名
func Permission2PermissionZhWithMap(p []P2V) []P2V {
	var np []P2V
	for _, v := range p {
		var nnp P2V
		pid := v.Permission
		if val, ok := common.PermissionIdZH[pid]; ok {
			nnp.PermissionZH = val
		} else {
			nnp.Permission = v.Permission
		}
		nnp.AppName = v.AppName
		nnp.VersionList = v.VersionList
		np = append(np, nnp)
	}
	return np
}
