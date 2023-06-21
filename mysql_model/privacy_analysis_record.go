package mysql_model

import (
	"skygo_detection/lib/common_lib/mysql"
)

type PrivacyAnalysisRecord struct {
	Id                int    `xorm:"not null pk comment('id') INT(11)" json:"id"`
	TaskId            int    `xorm:"not null  comment('任务ID') INT(11)" json:"task_id"`
	Uid               string `xorm:"not null comment('uid') VARCHAR(255)" json:"uid"`
	AppName           string `xorm:"not null comment('包名') VARCHAR(255)" json:"app_name"`
	Permission        string `xorm:"not null comment('权限') VARCHAR(255)" json:"permission"`
	PermissionDefault string `xorm:"not null comment('权限默认值') VARCHAR(255)" json:"permission_default"`
	PermissionState   string `xorm:"not null comment('权限状态') VARCHAR(255)"`
	PermissionMethod  string `xorm:"not null comment('权限请求方式') VARCHAR(255)"`
	PermissionTime    string `xorm:"not null comment('权限请求时间') VARCHAR(255)"`
	CreateTime        string `xorm:"not null comment('创建时间') DATETIME"`
	UpdateTime        string `xorm:"not null comment('更新时间') DATETIME"`
}

const PAGE_SIZE = 25

type PermissionList struct {
	Permission      string `json:"name"`
	PermissionZH    string `json:"name_zh"`
	PermissionTimes int64  `json:"times"`
}

func (this *PrivacyAnalysisRecord) Create() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}

// 应用调用包名列表
func (this *PrivacyAnalysisRecord) AppList(taskId string) (list []PrivacyAnalysisRecord, err error) {
	list = make([]PrivacyAnalysisRecord, 0)
	err = mysql.GetSession().
		Select("*").
		Where("task_id = ? ", taskId).
		GroupBy("app_name,task_id").Find(&list)
	return list, err
}

type PerTransfer struct {
	Count        int64  `json:"count"`
	PermissionZH string `json:"permission_zh"`
	Permission   string `json:"permission"`
}

// 权限请求（应用全部权限-不带包名）
func (this *PrivacyAnalysisRecord) PerTransferList(taskId string) (list []PerTransfer, err error) {
	list = make([]PerTransfer, 0)
	err = mysql.GetSession().Table("privacy_analysis_record").
		Select("count(*) count,permission").
		Where("task_id = ? ", taskId).
		GroupBy("permission").Find(&list)
	return list, err
}

// 权限请求（应用权限-带包名）
func (this *PrivacyAnalysisRecord) WithAppNamePerList(taskId string, appName string) (list []PerTransfer, err error) {
	list = make([]PerTransfer, 0)
	err = mysql.GetSession().Table("privacy_analysis_record").
		Select("count(*) count,permission").
		Where("task_id = ? AND app_name = ?", taskId, appName).
		GroupBy("permission").Limit(PAGE_SIZE).Find(&list)
	return list, err
}

// 请求记录（全部）
func (this *PrivacyAnalysisRecord) PerRecordList(taskId string) (list []PrivacyAnalysisRecord, err error) {
	list = make([]PrivacyAnalysisRecord, 0)
	err = mysql.GetSession().Table("privacy_analysis_record").
		Where("task_id = ? ", taskId).Limit(PAGE_SIZE).OrderBy("permission_time desc").Find(&list)
	return list, err
}

// 请求记录（带包名）
func (this *PrivacyAnalysisRecord) WithAppNamePerRecordList(taskId string, appName string) (list []PrivacyAnalysisRecord, err error) {
	list = make([]PrivacyAnalysisRecord, 0)
	err = mysql.GetSession().Table("privacy_analysis_record").
		Where("task_id = ? AND app_name = ?", taskId, appName).OrderBy("permission_time desc").Find(&list)
	return list, err
}

// 应用调用权限列表
func (this *PrivacyAnalysisRecord) PerList(taskId string) (list []PrivacyAnalysisRecord, err error) {
	list = make([]PrivacyAnalysisRecord, 0)
	err = mysql.GetSession().
		Select("*").
		Where("task_id = ? ", taskId).
		GroupBy("app_name,task_id,permission").Find(&list)
	return list, err
}

// 应用调用的次数
func (this *PrivacyAnalysisRecord) TransferCount(appName string, taskId string) (count int64, err error) {
	count, _ = mysql.GetSession().Table(PrivacyAnalysisRecord{}).
		Where("task_id = ? AND app_name = ?", taskId, appName).
		Count()
	return count, err
}

// 权限请求次数统计（不带权限名称）
func (this *PrivacyAnalysisRecord) PermissionTotal(appName string, taskId string) (count int64, err error) {
	count, _ = mysql.GetSession().Table(PrivacyAnalysisRecord{}).Distinct("permission").
		Where("task_id = ? AND app_name = ?", taskId, appName).
		Count()
	return count, err
}

// 权限请求次数统计（带权限）
func (this *PrivacyAnalysisRecord) PermissionCount(appName string, taskId string, permission string) (count int64, err error) {
	count, _ = mysql.GetSession().Table(PrivacyAnalysisRecord{}).
		Where("task_id = ? AND app_name = ? AND permission = ?", taskId, appName, permission).
		Count()
	return count, err
}

// 应用的权限
func (this *PrivacyAnalysisRecord) AppPerList(appName string, taskId string) (list []PrivacyAnalysisRecord, err error) {
	mysql.GetSession().Table(PrivacyAnalysisRecord{}).
		Where("task_id = ? AND app_name = ? ", taskId, appName).
		Find(&list)
	return list, err
}

// 应用权限请求列表
func (this *PrivacyAnalysisRecord) PermissionList(appName string, taskId string) (list []PermissionList, err error) {
	list = make([]PermissionList, 0)
	mysql.GetSession().Table(PrivacyAnalysisRecord{}).
		Where("task_id = ? AND app_name = ?", taskId, appName).Select("permission").
		Distinct("permission").
		Find(&list)
	return list, err
}

func (this *PrivacyAnalysisRecord) PermissionTimes(appName string, taskId string) (count int64, err error) {
	count, _ = mysql.GetSession().Table(PrivacyAnalysisRecord{}).
		Where("task_id = ? AND app_name = ?", taskId, appName).
		Count()
	return count, err
}

func (this *PrivacyAnalysisRecord) LastTimeByTaskId(taskId int) (perTime string, err error) {
	model := new(PrivacyAnalysisRecord)
	has, err := mysql.GetSession().Select("max(permission_time) permission_time").Where("task_id = ? ", taskId).Get(model)
	if !has {
		return "not found", err
	}
	if err != nil {
		return "something wrong", err
	}
	return model.PermissionTime, nil
}
