package toolbox

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/logic/toolbox"
	"skygo_detection/mysql_model"
	"skygo_detection/service"
)

var (
	ErrParams           = errors.New("参数不正确")
	ErrBufReadString    = errors.New("读取错误")
	ErrEOFBufReadString = errors.New("读取EOF错误")
)

type PrivacyController struct{}

type PrivacyAnalysisRecord struct {
	Uid               string
	Package           string
	Permission        string
	PermissionDefault string
	PermissionState   string
	PermissionMethod  string
	PermissionTime    string
}

func (pc PrivacyController) AnalysisRecord(ctx *gin.Context) {
	logger := service.GetDefaultLogger("privacy_analysis_record")
	defer logger.Sync()
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	taskId, res := req.TryString("task_id")
	if res == false {
		response.RenderFailure(ctx, ErrParams)
		return
	}
	taskUuid := mysql_model.GetIdByUid(taskId)
	if taskUuid == 0 {
		response.RenderFailure(ctx, errors.New("任务ID不存在"))
		return
	}
	var appOpsListCount int64
	var appInfoListCount int64
	// 应用版本记录
	appInfoListChangedFlag, res := req.TryBool("app_info_list_changed")
	if res == false {
		response.RenderFailure(ctx, ErrParams)
		return
	}

	if appInfoListChangedFlag == true {
		var appInfoListSlice []mysql_model.PrivacyAppVersion
		if appInfo, has := req.TrySlice("app_info_list"); has && len(appInfo) > 0 {
			for _, item := range appInfo {
				appInfoList := toolbox.AppInfoList(item.(map[string]interface{}), taskUuid)
				if len(appInfoList) > 0 {
					appInfoListSlice = append(appInfoListSlice, appInfoList...)
				}
			}
		}
		if len(appInfoListSlice) > 0 {
			// 添加应用信息
			count, err := toolbox.AddAppInfoLogic(taskUuid, appInfoListSlice)
			if err != nil {
				response.RenderFailure(ctx, err)
				return
			}
			appInfoListCount += count
		}
	}

	// 应用隐私记录
	appOpsFlag, has := req.TryBool("app_ops_changed")
	if has == false {
		response.RenderFailure(ctx, ErrParams)
		return
	}
	if appOpsFlag == true {
		str, has := req.TryString("app_ops")
		if has == false {
			response.RenderFailure(ctx, ErrParams)
			return
		}
		count, err := toolbox.AnalysisRecordLogic(str, taskUuid)
		if err != nil {
			response.RenderFailure(ctx, err)
			return
		}
		appOpsListCount += count
	}
	response.RenderSuccess(ctx, qmap.QM{
		"app_info_count": appInfoListCount,
		"app_ops_count":  appOpsListCount,
	})
	return
}

// 应用隐私调用统计
func (pc PrivacyController) AppList(ctx *gin.Context) {
	taskId := ctx.Request.FormValue("task_id")
	data, err := toolbox.AppListLogic(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, data)
	return
}

// 应用调用频次统计
func (pc PrivacyController) AppPerList(ctx *gin.Context) {
	taskId := ctx.Request.FormValue("task_id")
	data, err := toolbox.AppPerListLogic(taskId)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, data)
	return
}

// 应用权限请求列表
func (pc PrivacyController) PerCountList(ctx *gin.Context) {
	taskId := ctx.Request.FormValue("task_id")
	appName := ctx.Request.FormValue("app_name")
	data, err := toolbox.PerCountListLogic(taskId, appName)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, data)
	return
}

// 请求记录列表
func (pc PrivacyController) RecordList(ctx *gin.Context) {
	taskId := ctx.Request.FormValue("task_id")
	appName := ctx.Request.FormValue("app_name")
	data, err := toolbox.RecordListLogic(taskId, appName)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, data)
	return
}

func (pc PrivacyController) AppCount(ctx *gin.Context) {
	taskId := ctx.Request.FormValue("task_id")
	appName := ctx.Request.FormValue("app_name")
	fmt.Println(taskId)
	data, err := toolbox.AppCountLogic(taskId, appName)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	response.RenderSuccess(ctx, data)
	return
}

/**
 * apiType http
 * @api {get} /api/privacy/v1/log/:task_id 获取任务的日志
 * @apiVersion 1.0.0
 * @apiName GetLog
 * @apiGroup Privacy
 * @apiDescription 获取任务的日志
 *
 *
 * @apiExample {curl} 请求示例:
 * curl -i -X GET -d  http://localhost/api/privacy/v1/log/23
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *   "code": 0,
 *   "data": [
 *       {
 *           "id": 1,
 *           "task_id": 23,
 *           "create_time": 22,
 *           "content": "fgbfg"
 *       },
 *       {
 *           "id": 2,
 *           "task_id": 23,
 *           "create_time": 345,
 *           "content": "sgb"
 *       }
 *   ],
 *   "msg": ""
 * }
 */
func (pc PrivacyController) GetLog(ctx *gin.Context) {
	tid := request.ParamInt(ctx, "task_id")
	res := mysql_model.GetPrivacyLog(tid)
	response.RenderSuccess(ctx, res)
}
