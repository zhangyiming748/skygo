package toolbox

import (
	"skygo_detection/mysql_model"
)

// 获取任务相关的所有app
func GetAllTaskApp(tid int) ([]mysql_model.PrivacyAppVersion, error) {
	task := new(mysql_model.PrivacyAppVersion)
	task.TaskId = tid
	return task.FindAppListByTaskId()
}

// 获取选中app的所有版本
func GetAppAllVersion(app string) ([]mysql_model.PrivacyAppVersion, error) {
	this := new(mysql_model.PrivacyAppVersion)
	this.AppName = app
	this.AppNameZh = app
	return this.FindAppVersionList(app)
}

// 获取选中app版本的权限
func GetPermission(name, version string) ([]mysql_model.PrivacyAppVersion, error) {
	this := new(mysql_model.PrivacyAppVersion)
	this.AppName = name
	this.VersionName = version
	return this.FindPermissionByVersion()
}

// 获取app所有版本涉及到的权限
func GetAppAllPermission(app string) []int {
	this := new(mysql_model.PrivacyAppVersion)
	this.AppName = app
	permission, err := this.FindAppAllPermission()
	if err != nil {
		return nil
	}
	plist := make([]int, 0)
	for _, v := range permission {
		plist = append(plist, v.Permission)
	}
	return plist
}

type P2V struct {
	AppName      string   `json:"app_name"`      //应用名
	Permission   int      `json:"permission"`    //权限id
	PermissionZH string   `json:"permission_zh"` //权限id
	VersionList  []string `json:"version_list"`  //对应的版本号数组
}

func UsePermission(app string) ([]P2V, error) {
	plist := GetAppAllPermission(app)
	list := make([]P2V, 0)
	for _, v := range plist {
		var p2v P2V
		s := new(mysql_model.PrivacyAppVersion)
		s.AppName = app
		s.Permission = v
		p2v.AppName = app
		p2v.Permission = v
		permission, err := s.FindVersionWithPermission()
		if err != nil {
			return nil, err
		}
		for _, v := range permission {
			p2v.VersionList = append(p2v.VersionList, v.VersionName)
		}
		list = append(list, p2v)
	}
	return list, nil
}

//type Compare struct {
//	Name           string `json:"name"`
//	Permission     string `json:"permission"`
//	PermissionIdZH string `json:"permission_zh"`
//	Version        string `json:"version"`
//}
//
////获取全部选中版本的权限
//func GetAllPermission(name string, versions []string) ([]Compare, error) {
//	c := make([]Compare, 0)
//	for _, version := range versions {
//		permissions, err := GetPermission(name, version)
//		if err != nil {
//			return nil, err
//		}
//		for _, permission := range permissions {
//			var one Compare
//			one.Name = permission.AppNameZh
//			one.Permission = permission.Permission
//			pid, err := strconv.Atoi(permission.Permission)
//			if err != nil {
//				return nil, err
//			}
//			if v, ok := common.PermissionIdZH[pid]; ok {
//				one.PermissionIdZH = v
//			}
//			one.Version = permission.VersionName
//			c = append(c, one)
//		}
//	}
//	return c, nil
//}
