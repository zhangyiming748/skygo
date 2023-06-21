package mysql_model

import (
	"skygo_detection/lib/common_lib/mysql"
)

type PrivacyAppVersion struct {
	Id          int    `xorm:"not null pk comment('id') INT(11)" json:"id"`
	TaskId      int    `xorm:"not null  comment('任务ID') INT(11)" json:"task_id"`
	AppName     string `xorm:"not null comment('包名') VARCHAR(255)" json:"app_name"`
	AppNameZh   string `xorm:"not null comment('包中文名') VARCHAR(255)" json:"app_name_zh"`
	Permission  int    `xorm:"not null comment('权限名') INT(11)" json:"permission"`
	VersionCode string `xorm:"not null comment('版本代号') VARCHAR(255)" json:"version_code"`
	VersionName string `xorm:"not null comment('版本号') VARCHAR(255)" json:"version_name"`
	Path        string `xorm:"not null comment('应用路径') VARCHAR(255)" json:"path"`
	Uid         int    `xorm:"not null comment('uid') VARCHAR(255)" json:"uid"`
	CreateTime  string `xorm:"not null comment('创建时间') DATETIME" json:"create_time"`
	UpdateTime  string `xorm:"not null comment('更新时间') DATETIME" json:"update_time"`
}
type PToName struct {
	Permission int
	Name       string
}

// SELECT v.permission,m.name FROM `privacy_app_version` AS `v` INNER JOIN `privacy_app_permission_map` as `m` ON v.permission = m.id WHERE (v.permission = 19) AND (v.app_name = 'com.android.defcontainer')
func FindZhWithPermission() {

}

// 批量插入
func BatchesInsert(data []PrivacyAppVersion) (int64, error) {
	//privacyAppVersions := make([]*PrivacyAppVersion, 0)
	return mysql.GetSession().Insert(&data)
}

func (this *PrivacyAppVersion) Create() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}

// 查找与任务相关的app
func (this *PrivacyAppVersion) FindAppListByTaskId() ([]PrivacyAppVersion, error) {
	result := make([]PrivacyAppVersion, 0)
	err := mysql.GetSession().Select("*").Where("task_id = ?", this.TaskId).GroupBy("app_name").Find(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 查找与app相关的版本号
func (this *PrivacyAppVersion) FindAppVersionList(name string) ([]PrivacyAppVersion, error) {
	result := make([]PrivacyAppVersion, 0)
	err := mysql.GetSession().Select("*").
		Where("app_name = ? ", name).
		GroupBy("version_name").Find(&result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// 查找app指定版本号的权限
func (this *PrivacyAppVersion) FindPermissionByVersion() ([]PrivacyAppVersion, error) {
	result := make([]PrivacyAppVersion, 0)
	err := mysql.GetSession().Where("version_name = ?", this.VersionName).Find(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 查找app所有版本涉及到的权限
func (this PrivacyAppVersion) FindAppAllPermission() ([]PrivacyAppVersion, error) {
	result := make([]PrivacyAppVersion, 0)
	err := mysql.GetSession().Where("app_name = ?", this.AppName).Distinct("permission").Find(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// 查找app权限对应的版本号
func (this PrivacyAppVersion) FindVersionWithPermission() ([]PrivacyAppVersion, error) {
	result := make([]PrivacyAppVersion, 0)
	err := mysql.GetSession().Where("permission = ?", this.Permission).And("app_name = ?", this.AppName).Find(&result)
	if err != nil {
		return nil, err
	}
	return result, nil
}
