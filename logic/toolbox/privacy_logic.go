package toolbox

import (
	"bufio"
	"regexp"
	"skygo_detection/common"
	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/mysql_model"
	"strconv"
	"strings"
	"time"
)

// 批量添加数据
func AddAppInfoLogic(taskId int, privacyAppVersionSlice []mysql_model.PrivacyAppVersion) (count int64, err error) {
	if len(privacyAppVersionSlice) > 0 {
		_ = new(mysql_model.PrivacyAnalysisLog).SetPrivacyLog(taskId, "开始应用信息分析")
		count, err = mysql_model.BatchesInsert(privacyAppVersionSlice)
		if err != nil {
			return 0, err
		}
		_ = new(mysql_model.PrivacyAnalysisLog).SetPrivacyLog(taskId, "结束应用信息分析")
	}
	return count, nil
}

// 应用数据处理
func AppInfoList(appInfo qmap.QM, taskId int) []mysql_model.PrivacyAppVersion {
	nowTime := time.Now()
	var privacyAppVersionSlice []mysql_model.PrivacyAppVersion
	permissions := appInfo.String("permissions")
	if permissions != "" {
		// 权限分割
		perSlice := strings.Split(permissions, `|`)
		for _, permission := range perSlice {
			per, _ := strconv.Atoi(permission)
			privacyAppVersion := mysql_model.PrivacyAppVersion{
				TaskId:      taskId,
				AppName:     appInfo.String("package_name"),
				AppNameZh:   appInfo.String("name"),
				Permission:  per,
				VersionCode: appInfo.String("version_code"),
				VersionName: appInfo.String("version_name"),
				Path:        appInfo.String("path"),
				Uid:         appInfo.Int("uid"),
				CreateTime:  nowTime.Format("2006-01-02 15:04:05"),
				UpdateTime:  nowTime.Format("2006-01-02 15:04:05"),
			}
			privacyAppVersionSlice = append(privacyAppVersionSlice, privacyAppVersion)
		}
	}
	return privacyAppVersionSlice
}

func AnalysisRecordLogic(str string, taskId int) (data int64, err error) {
	str = str + "END_FLAG"
	str = strings.Replace(str, "\n", "#", -1)
	sr := strings.NewReader(str)
	buf := bufio.NewReader(sr)
	privacyAnalysisRecord := new(mysql_model.PrivacyAnalysisRecord)
	privacyAnalysisRecord.TaskId = taskId
	var count int64
	err = new(mysql_model.PrivacyAnalysisLog).SetPrivacyLog(taskId, "开始隐私日志分析")
	if err != nil {
		return 0, err
	}
	// 查询该任务下最近的一条时间
	perTime, err := new(mysql_model.PrivacyAnalysisRecord).LastTimeByTaskId(taskId)
	if err != nil {
		return 0, err
	}
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation("2006-01-02 15:04:05", perTime, loc)
	unixLastTime := theTime.Unix()
	for {
		line, _ := buf.ReadString('#')
		line = strings.TrimSpace(line)
		//if err != nil {
		//	if err == io.EOF {
		//		response.RenderFailure(ctx, ErrEOFBufReadString)
		//		return
		//	}
		//	response.RenderFailure(ctx, ErrBufReadString)
		//	return
		//}
		// 匹配 Uid
		r1 := regexp.MustCompile(`Uid\s*(.*?)\s*:`)
		uid := r1.FindAllStringSubmatch(line, -1)
		if uid != nil {
			privacyAnalysisRecord.Uid = uid[0][1]
			continue
		}
		// 匹配 Package
		r2 := regexp.MustCompile(`Package\s*(.*?)\s*:`)
		packages := r2.FindAllStringSubmatch(line, -1)
		if packages != nil {
			privacyAnalysisRecord.AppName = packages[0][1]
			continue
		}
		// 匹配权限行
		r3 := regexp.MustCompile(`[A-Z_]+\s*(.*?)\s*\):`)
		allow := r3.FindAllStringSubmatch(line, -1)
		if allow != nil {
			// 匹配权限
			r31 := regexp.MustCompile(`\s*(.*?)\s*\(`)
			permission := r31.FindAllStringSubmatch(allow[0][0], -1)
			privacyAnalysisRecord.Permission = permission[0][1]
			// 匹配默认权限
			r32 := regexp.MustCompile(`\(\s*(.*?)\s*\):`)
			permissionDefault := r32.FindAllStringSubmatch(allow[0][0], -1)
			privacyAnalysisRecord.PermissionDefault = permissionDefault[0][1]
			continue
		}
		// 匹配 PermissionMethod  相关数据
		r5 := regexp.MustCompile(`(pers|top|fgsvc|fg|bg|cch)`)
		permissionMethod := r5.FindAllStringSubmatch(line, -1)
		if permissionMethod != nil {
			// 匹配 PermissionState  相关数据
			r4 := regexp.MustCompile(`(Access|Reject)`)
			permissionState := r4.FindAllStringSubmatch(line, -1)
			if permissionState != nil {
				privacyAnalysisRecord.PermissionState = permissionState[0][1]
			}
			privacyAnalysisRecord.PermissionMethod = permissionMethod[0][1]
			// 匹配 PermissionTime 相关数据
			r6 := regexp.MustCompile(`\d{4}-\d{1,2}-\d{1,2} \d{1,2}:\d{1,2}:\d{1,2}.\d{1,3}`)
			permissionTime := r6.FindAllStringSubmatch(line, -1)
			if permissionTime != nil {
				currentPermissionTime, _ := time.ParseInLocation("2006-01-02 15:04:05", permissionTime[0][0], loc)
				currentUnixTime := currentPermissionTime.Unix()
				// 当前权限的时间大于等于最近一条数据的时间
				if currentUnixTime >= unixLastTime {
					privacyAnalysisRecord.PermissionTime = permissionTime[0][0]
					nowTime := time.Now()
					privacyAnalysisRecord.UpdateTime = nowTime.Format("2006-01-02 15:04:05")
					privacyAnalysisRecord.CreateTime = nowTime.Format("2006-01-02 15:04:05")
					// 添加数据
					_, err := privacyAnalysisRecord.Create()
					if err != nil {
						return 0, err
					}
					count++
				}
			}
			continue
		}
		// 匹配结尾
		r7 := regexp.MustCompile(`END_FLAG\s*(.*?)\s*`)
		endFlag := r7.FindAllStringSubmatch(line, -1)
		if endFlag != nil {
			break
		}
	}
	err = new(mysql_model.PrivacyAnalysisLog).SetPrivacyLog(taskId, "结束隐私日志分析")
	if err != nil {
		return 0, err
	}
	return count, nil
}

func AppCountLogic(taskId string, appName string) (appInfo map[string]interface{}, err error) {
	// 请求权限数量
	permissionTotal, err := new(mysql_model.PrivacyAnalysisRecord).PermissionTotal(appName, taskId)
	if err != nil {
		return appInfo, err
	}
	// 请求权限次数
	permissionTimes, err := new(mysql_model.PrivacyAnalysisRecord).PermissionTimes(appName, taskId)
	if err != nil {
		return appInfo, err
	}
	// 权限请求列表
	permissionList, err := new(mysql_model.PrivacyAnalysisRecord).PermissionList(appName, taskId)
	if err != nil {
		return appInfo, err
	}
	// 权限请求次数
	for k, v := range permissionList {
		permissionCount, err := new(mysql_model.PrivacyAnalysisRecord).PermissionCount(appName, taskId, v.Permission)
		if err != nil {
			return appInfo, err
		}
		permissionList[k].PermissionTimes = permissionCount
		// 权限中文名称
		permissionList[k].PermissionZH = permissionList[k].Permission
		if zh, ok := common.Permission[permissionList[k].Permission]; ok {
			permissionList[k].PermissionZH = zh
		}
	}
	// 权限请求记录
	appInfo = make(map[string]interface{}, 0)
	appInfo = map[string]interface{}{
		"count": permissionTotal,
		"times": permissionTimes,
		"list":  permissionList,
	}
	return appInfo, nil
}

func AppListLogic(taskId string) (appInfo []map[string]interface{}, err error) {
	appInfoSlice := make([]map[string]interface{}, 0)
	appPerSlice := make([]mysql_model.PrivacyAnalysisRecord, 0)
	list, err := new(mysql_model.PrivacyAnalysisRecord).AppList(taskId)
	if err != nil {
		return
	}
	for _, v := range list {
		// 应用调用的总次数
		transferCount, err := new(mysql_model.PrivacyAnalysisRecord).TransferCount(v.AppName, taskId)
		if err != nil {
			return appInfo, err
		}
		// 应用调用权限次数
		permissionTotal, err := new(mysql_model.PrivacyAnalysisRecord).PermissionTotal(v.AppName, taskId)
		if err != nil {
			return appInfo, err
		}
		var zh_cn string
		if value, ok := common.AppsName[v.AppName]; ok {
			zh_cn = value
		} else {
			zh_cn = v.AppName
		}
		temp := map[string]interface{}{
			"app_name":               v.AppName,
			"app_name_zh":            zh_cn,
			"app_transfer_count":     transferCount,   // 调用总次数
			"app_transfer_per_count": permissionTotal, // 调用权限次数
		}
		appInfoSlice = append(appInfoSlice, temp)
		// 每个应用的权限
		appPerList, err := new(mysql_model.PrivacyAnalysisRecord).AppPerList(v.AppName, taskId)
		if err != nil {
			return appInfo, err
		}
		appPerSlice = append(appPerSlice, appPerList...)
	}
	return appInfoSlice, nil
}

func PerCountListLogic(taskId string, appName string) (list []mysql_model.PerTransfer, err error) {
	if appName != "" {
		list, err = new(mysql_model.PrivacyAnalysisRecord).WithAppNamePerList(taskId, appName)
	} else {
		list, err = new(mysql_model.PrivacyAnalysisRecord).PerTransferList(taskId)
	}
	if err != nil {
		return
	}
	for k := range list {
		if zh_cn, ok := common.Permission[list[k].Permission]; ok {
			list[k].PermissionZH = zh_cn
		} else {
			list[k].PermissionZH = list[k].Permission
		}
	}
	return list, nil
}

func RecordListLogic(taskId string, appName string) (list []map[string]interface{}, err error) {
	slice := make([]map[string]interface{}, 0)
	var recordList []mysql_model.PrivacyAnalysisRecord
	if appName != "" {
		recordList, err = new(mysql_model.PrivacyAnalysisRecord).WithAppNamePerRecordList(taskId, appName)
	} else {
		recordList, err = new(mysql_model.PrivacyAnalysisRecord).PerRecordList(taskId)
	}
	if err != nil {
		return nil, err
	}
	for _, v := range recordList {
		temp := map[string]interface{}{
			"app_name":              v.AppName,
			"app_name_zh":           v.AppName,
			"permission":            v.Permission,
			"permission_zh":         v.Permission,
			"permission_time":       v.PermissionTime,
			"permission_default":    v.PermissionDefault,
			"permission_default_zh": v.PermissionDefault,
			"permission_method":     v.PermissionMethod,
			"permission_method_zh":  "前台请求",
		}
		ancn, ok := common.AppsName[v.AppName]
		if ok {
			temp["app_name_zh"] = ancn
		}
		pcn, ok := common.Permission[v.Permission]
		if ok {
			temp["permission_zh"] = pcn
		}
		pdcn, ok := common.Permission_default[v.PermissionDefault]
		if ok {
			temp["permission_default_zh"] = pdcn
		}
		slice = append(slice, temp)
	}
	return slice, nil
}

func AppPerListLogic(taskId string) (appInfo []map[string]interface{}, err error) {
	// 包名列表
	appInfoSlice := make([]map[string]interface{}, 0)
	list, err := new(mysql_model.PrivacyAnalysisRecord).AppList(taskId)
	if err != nil {
		return
	}
	for _, v := range list {
		// 权限请求列表
		permissionList, err := new(mysql_model.PrivacyAnalysisRecord).PermissionList(v.AppName, taskId)
		if err != nil {
			return appInfo, err
		}
		// 处理本包下所有权限调用频次统计
		for k, item := range permissionList {
			permissionCount, err := new(mysql_model.PrivacyAnalysisRecord).PermissionCount(v.AppName, taskId, item.Permission)
			if err != nil {
				return appInfo, err
			}
			// item.PermissionTimes = permissionCount
			pcn, ok := common.Permission[item.Permission]
			if ok {
				permissionList[k].PermissionZH = pcn
			} else {
				permissionList[k].PermissionZH = item.Permission
			}
			permissionList[k].PermissionTimes = permissionCount
		}
		// 应用调用的总次数
		temp := map[string]interface{}{
			"app_name":    v.AppName,
			"app_name_zh": v.AppName,
			"per_list":    permissionList,
		}
		ancn, ok := common.AppsName[v.AppName]
		if ok {
			temp["app_name_zh"] = ancn
		}
		appInfoSlice = append(appInfoSlice, temp)
	}
	return appInfoSlice, nil
}
