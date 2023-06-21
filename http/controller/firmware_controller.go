package controller

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/app/sys_service"

	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/common_lib/session"
	"skygo_detection/mongo_model"
	"skygo_detection/mysql_model"
)

type FirmwareController struct{}

type FirmwareUploadResp struct {
	Code int         `bson:"code"`
	Data interface{} `bson:"data"`
}

/*
 * apiType http
 * @api {post} /api/v1/firmware/upload_firmware_msg 上传固件
 * @apiVersion 1.0.1
 * @apiName UploadFirmWare
 * @apiGroup Firmware
 *
 * @apiDescription 上传固件信息获取上传授权
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {string}      			device_model      	设备型号
 * @apiParam {string}      			device_name       	设备名称
 * @apiParam {string}      			firmware_version       	固件版本
 * @apiParam {string}      			device_type       	固件类型
 * @apiParam {string}      			template_id       	模板名称
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/firmware/upload_firmware_msg
 *
 * {
 *    "code": 0,
 *    "data": {
 *        "id": "5f8fe0b089e90f2127ffe4c5",//授权ID
 *        "project_id": 121//项目ID
 *    }
 * }
 *
 */
func (this FirmwareController) UploadFirmWareMsg(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	DeviceModel := params.MustString("device_model")         // 设备型号
	DeviceName := params.MustString("device_name")           // 设备名称
	FirmwareVersion := params.MustString("firmware_version") // 固件版本
	DeviceType := params.MustString("device_type")           // 设备类型
	TemplateId := params.MustInt("template_id")              // 模板ID
	FirmwareName := params.MustString("file_name")           // 此固件名称为用户上传的固件名称

	UserID := int(session.GetUserId(ctx))
	UserName := session.GetUserName(ctx)
	if res, err := new(mongo_model.FirmWareData).CreateFirmWareMsg(FirmwareName, DeviceName, DeviceModel, FirmwareVersion, DeviceType, UserName, TemplateId, UserID); err == nil {
		data := map[string]interface{}{
			"id":         res.Id,
			"project_id": res.ProjectId,
		}
		response.RenderSuccess(ctx, data)
		// qm := qmap.QM{
		//	"data": data,
		// }
		// return &qm, nil
	} else {
		response.RenderSuccess(ctx, err)
	}

}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

/*
 * apiType http
 * @api {get} /api/v1/firmware/download 下载固件
 * @apiVersion 1.0.1
 * @apiName DownloadFirmWare
 * @apiGroup Firmware
 *
 * @apiDescription 下载已上传固件
 *
 * @apiUse authHeader urlQueryParams
 *
 * @apiParam {string}      			name       	上传反回值固件名称
 * @apiParam {bearer token}      		token       	授权token
 *
 */
// func (this FirmwareController) DownloadFirmWare(ctx *gin.Context) {
//	fileName := ctx.Request.FormValue("name")
//	//fmt.Println(fileName)
//	FilePath := "/data/firmware_store/" + fileName
//	//FilePath := "/Users/zhupenghui/data/firmware_store/"+fileName
//	_, err := os.Stat(FilePath) //os.Stat获取文件信息
//	//fmt.Println(err)
//	if err == nil || os.IsExist(err) {
//		file, opErr := os.Open(FilePath)
//		if opErr != nil {
//			//fmt.Println(opErr)
//			ctx.JSON(400, gin.H{
//				"code": -1,
//				"msg":  "os open error",
//			})
//		}
//		defer file.Close()
//		content, conErr := ioutil.ReadAll(file)
//		if conErr != nil {
//			ctx.JSON(400, gin.H{
//				"code": -1,
//				"msg":  "file read error",
//			})
//		}
//
//		//存在
//		//ctx.Stream()
//		//fileContent := this.DownloadExcel_tem()
//		ctx.Writer.WriteHeader(200)
//		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
//		ctx.Header("Content-Type", "*")
//		ctx.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
//		ctx.Writer.Write(content)
//	} else {
//		//response.RenderSuccess(ctx, &qmap.QM{"id": resp.FileId})
//		ctx.JSON(400, gin.H{
//			"code": -1,
//			"msg":  "file not found",
//		})
//	}
// }

/*
 * apiType http
 * @api {post} /api/v1/firmware/upload_firmware_file 上传固件File信息
 * @apiVersion 1.0.1
 * @apiName UploadFirmWareFile
 * @apiGroup Firmware
 *
 * @apiDescription 上传固件FIle信息
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			project_id      	项目ID
 * @apiParam {string}      			id       	上传固件授权ID
 * @apiParam {form-data}      		file       	固件文件
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/firmware/upload_firmware_file
 *
 * {
 *    "code": 0,
 *    "data": {
 *        "id": "5f8fe0b089e90f2127ffe4c5",
 *        "project_id": 121
 *    }
 * }
 *
 */
func (this FirmwareController) FirmWareUpload(ctx *gin.Context) {
	projectId := ctx.Request.FormValue("project_id")
	masterId := ctx.Request.FormValue("id")
	fileName := ctx.Request.FormValue("file_name")
	userId := http_ctx.GetUserId(ctx)
	userName := http_ctx.GetUserName(ctx)
	if file, header, err := ctx.Request.FormFile("file"); err == nil {
		if fileName == "" {
			fileName = header.Filename
		}
		fileContent, _ := ioutil.ReadAll(file)
		fileSize := len(fileContent)
		fileMd5Byte16 := md5.Sum(fileContent)
		fileMd5 := fmt.Sprintf("%x", fileMd5Byte16)
		if fileId, err := mongo.GridFSUpload(common.MC_File, fileName, fileContent); err == nil {
			proId, _ := strconv.Atoi(projectId)
			pUserIDStr := strconv.FormatInt(userId, 10)
			pUserID, _ := strconv.Atoi(pUserIDStr)
			if err := new(mongo_model.FirmWareData).UploadFirmWareUrl(proId, fileSize, pUserID, fileId, fileId, fileMd5, masterId, userName); err == nil {
				resule := qmap.QM{
					"id":         masterId,
					"project_id": projectId,
				}
				response.RenderSuccess(ctx, resule)
			} else {
				response.RenderFailure(ctx, err)
			}
		} else {
			response.RenderFailure(ctx, err)
		}
	}
}

func (this FirmwareController) DownloadFirmWare(ctx *gin.Context) {
	fileId := ctx.Query("name")
	if fi, err := mongo.GridFSOpenId(common.MC_File, bson.ObjectIdHex(fileId)); err == nil {
		fileContent, _ := ioutil.ReadAll(fi)
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fi.Name()))
		ctx.Header("Content-Type", "*")
		ctx.Header("Accept-Length", fmt.Sprintf("%d", len(fileContent)))
		ctx.Writer.Write(fileContent)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/*
 * apiType http
 * @api {post} /api/v1/firmware/update_upload_status 更改固件上传状态
 * @apiVersion 1.0.1
 * @apiName UpdateUploadStatus
 * @apiGroup Firmware
 *
 * @apiDescription 更改固件上传中固件显示状态
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			project_id      	项目ID
 * @apiParam {int}      			status       	固件状态
 * @apiParam {string}      			master_id       	主键ID
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/firmware/update_upload_status
 *
 * {
 *    "code": 0,
 *    "data": {
 *        "id": "5f8fe0b089e90f2127ffe4c5",
 *        "project_id": 121
 *    }
 * }
 *
 */
func (this FirmwareController) UpdateUploadStatus(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	projectId := request.MustInt(ctx, "project_id")
	masterId := request.MustString(ctx, "master_id")
	status := request.MustInt(ctx, "status")
	Progress := request.MustInt(ctx, "progress")
	UserID := int(session.GetUserId(ctx))
	UserName := session.GetUserName(ctx)
	if status != 3 && status != 4 && status != 1 {
		response.RenderFailure(ctx, errors.New("status error"))
	}
	// 更新status
	params1 := qmap.QM{
		"e_project_id":     projectId,
		"e__id":            bson.ObjectIdHex(masterId),
		"e_upload_user":    UserName,
		"e_upload_user_id": UserID,
		"ne_status":        0,
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_FIRMWARE_UPLOAD_LOG, params1)
	update := bson.M{}
	if status == 1 {
		update = bson.M{"$set": bson.M{"progress": Progress}}
	} else {
		update = bson.M{"$set": bson.M{"status": status}}
	}
	if err := mongoClient.Update(bson.M{"_id": bson.ObjectIdHex(masterId)}, update); err != nil {
		response.RenderFailure(ctx, errors.New("update status error"))
	} else {
		data := map[string]interface{}{
			"id":         masterId,
			"project_id": projectId,
		}
		response.RenderSuccess(ctx, data)
	}
}

/*
 * apiType http
 * @api {get} /api/v1/firmware/list 固件列表
 * @apiVersion 1.0.1
 * @apiName UploadFirmWare
 * @apiGroup List
 *
 * @apiDescription 固件检测列表
 *
 * @apiUse authHeader
 *
 * @apiParam {string}      			current_page      	当前页
 * @apiParam {int}      			limit       			单页数量
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/firmware/list
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *     	"list":[
 * 	    	{
 * 	    		"id": "5f681e3589e90ff95d7bf66b",	//主键ID
 * 		        "project_id":105,
 * 		        "task_id":111,
 * 		        "template_id":77,
 * 		        "firmware_name":"固件名称xxxx",
 * 		        "firmware_size":1024,
 * 		        "status": 1,		//状态 0 未解析 1 解析完成 2 已删除 3 已取消任务
 * 		        "create_time":1600422799,
 * 	    	},
 * 	    	{
 * 	    		"id": "5f681e3589e90ff95d7bf66c",	//主键ID
 * 		        "project_id":106,
 * 		        "task_id":112,
 * 		        "template_id":98,
 * 		        "firmware_name":"固件名称xxxx",
 * 		        "firmware_size":1024,
 * 		        "status": 1,		//状态 0 未解析 1 解析完成 2 已删除 3
 * 		        "create_time":1600422799,
 * 	    	}
 *
 * 	    ]
 *
 * 	},
 * 	"pagination": {
 *         "count": 20,
 *         "current_page": 1,
 *         "total": 202,
 *         "total_page": 11
 *     },
 * }
 *
 */
/*
 * 固件列表
 * 状态 	1 待上传 2 上传完成 3 上传失败 4 取消上传 5 (下载完成) 扫描中 6 (创建任务) 扫描中 7 取消扫描 8 扫描完成 9 扫描失败 10 已解析 0 已删除
 */
func (this FirmwareController) List(ctx *gin.Context) {
	UserID := int(http_ctx.GetUserId(ctx))
	UserName := http_ctx.GetUserName(ctx)
	queryParams := qmap.QM{
		// "in_status": []int{1, 4},
		"ne_status":        0,
		"e_upload_user":    UserName,
		"e_upload_user_id": UserID,
	}

	mgoSession := mongo.NewMgoSession(common.MC_FIRMWARE_UPLOAD_LOG).AddCondition(queryParams).AddUrlQueryCondition(ctx.Request.URL.RawQuery)

	mgoSession.SetTransformFunc(this.FirmWareListFormat)
	if res, err := mgoSession.GetPage(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this FirmwareController) FirmWareListFormat(data qmap.QM) qmap.QM {

	data["id"] = data["_id"]
	delete(data, "_id")
	delete(data, "device_model")
	delete(data, "device_name")
	delete(data, "device_type")
	delete(data, "firmware_md5")
	// delete(data, "response_time")
	delete(data, "task_name")
	delete(data, "temp_file")
	delete(data, "tmp_file_path")
	delete(data, "tmp_hd_file_path")
	delete(data, "upload_time")
	delete(data, "template_name")

	return data
}

/*
 * apiType http
 * @api {post} /api/v1/firmware/start_scanning 开始扫描 重新扫描
 * @apiVersion 1.0.1
 * @apiName StartScanning
 * @apiGroup Firmware
 *
 * @apiDescription 上传固件 开始扫描 重新扫描
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			project_id      	项目ID
 * @apiParam {int}      			status       	固件状态
 * @apiParam {string}      			master_id       	主键ID
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/firmware/start_scanning
 *
 * {
 *    "code": 0,
 *    "data": {
 *        "id": "5f8fe0b089e90f2127ffe4c5",
 *        "project_id": 121,
 * 		  "task_id":102
 *    }
 * }
 *
 */
func (this FirmwareController) StartScanning(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	projectId := params.MustInt("project_id")
	taskId := params.MustInt("task_id")
	masterId := params.MustString("master_id")

	UserID := int(session.GetUserId(ctx))
	UserName := session.GetUserName(ctx)
	setStatus, err := new(mongo_model.FirmWareData).StartScanning(masterId, UserName, taskId, projectId, UserID)
	if setStatus {
		data := map[string]interface{}{
			"id":         masterId,
			"project_id": projectId,
			"task_id":    taskId,
		}
		response.RenderSuccess(ctx, data)

	} else {
		response.RenderFailure(ctx, err)
	}
}

/*
 * apiType http
 * @api {post} /api/v1/firmware/start_scanning 取消扫描 暂停扫描
 * @apiVersion 1.0.1
 * @apiName CancelScanning
 * @apiGroup Firmware
 *
 * @apiDescription 固件信息取消扫描 暂停扫描
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			project_id      	项目ID
 * @apiParam {int}      			status       	固件状态
 * @apiParam {string}      			master_id       	主键ID
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/firmware/start_scanning
 *
 * {
 *    "code": 0,
 *    "data": {
 *        "id": "5f8fe0b089e90f2127ffe4c5",
 *        "project_id": 121,
 * 		  "task_id":102
 *    }
 * }
 *
 */
/*
 * 取消扫描  暂停扫描 status 7
 * 只有扫描中的status=6的固件才可以取消扫描
 */
func (this FirmwareController) CancelScanning(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	projectId := params.MustInt("project_id")
	taskId := params.MustInt("task_id")
	masterId := params.MustString("master_id")

	UserID := int(http_ctx.GetUserId(ctx))
	UserName := http_ctx.GetUserName(ctx)
	setStatus, err := new(mongo_model.FirmWareData).StopScanning(masterId, UserName, taskId, projectId, UserID)

	if setStatus {
		data := map[string]interface{}{
			"id":         masterId,
			"project_id": projectId,
			"task_id":    taskId,
		}
		response.RenderSuccess(ctx, data)

	} else {
		response.RenderFailure(ctx, err)
	}
}

/*
 * apiType http
 * @api {post} /api/v1/firmware/del_task 删除任务
 * @apiVersion 1.0.1
 * @apiName DelTask
 * @apiGroup Firmware
 *
 * @apiDescription 删除任务
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			project_id      	项目ID
 * @apiParam {int}      			status       	固件状态
 * @apiParam {string}      			master_id       	主键ID
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/firmware/del_task
 *
 * {
 *    "code": 0,
 *    "data": {
 *        "id": "5f8fe0b089e90f2127ffe4c5",
 *        "project_id": 121,
 * 		  "task_id":102
 *    }
 * }
 *
 */
func (this FirmwareController) DelTask(ctx *gin.Context) {
	UserID := int(http_ctx.GetUserId(ctx))
	UserName := http_ctx.GetUserName(ctx)
	masterId := request.MustString(ctx, "master_id")
	projectId := request.MustInt(ctx, "project_id")
	taskId := request.MustInt(ctx, "task_id")
	params := qmap.QM{
		"e_project_id":     projectId,
		"e_task_id":        taskId,
		"e__id":            bson.ObjectIdHex(masterId),
		"e_upload_user":    UserName,
		"e_upload_user_id": UserID,
		"ne_status":        int(0),
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_FIRMWARE_UPLOAD_LOG, params)

	if err := mongoClient.Update(bson.M{"_id": bson.ObjectIdHex(masterId)}, bson.M{"$set": bson.M{"status": 0}}); err != nil {
		response.RenderFailure(ctx, errors.New("msg update error"))
	} else {
		data := map[string]interface{}{
			"id":         bson.ObjectIdHex(masterId),
			"project_id": projectId,
			"task_id":    taskId,
		}
		response.RenderSuccess(ctx, data)
	}
}

/*
 * apiType http
 * @api {get} /api/v1/firmware/basic 获取固件基本信息
 * @apiVersion 1.0.1
 * @apiName Basic
 * @apiGroup Firmware
 *
 * @apiDescription 获取固件基本信息
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			project_id      	项目ID
 * @apiParam {int}      			task_id       	    任务ID
 * @apiParam {string}      			master_id       	主键ID
 *
 */
func (this FirmwareController) Basic(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	queryParams := qmap.QM{
		"e_project_id": params.MustInt("project_id"),
		"e_task_id":    params.MustInt("task_id"),
		"e__id":        bson.ObjectIdHex(params.MustString("master_id")),
		// "e_upload_user":    UserName,
		// "e_upload_user_id": UserID,
		"ne_status": int(0),
	}
	mgoSession := mongo.NewMgoSession(common.MC_FIRMWARE_UPLOAD_LOG).AddCondition(queryParams)
	mgoSession.SetTransformFunc(this.FirmWareBasicFormat)
	if res, err := mgoSession.GetOne(); err == nil {
		// 获取elf统计数据
		queryElfParams := qmap.QM{
			"e_project_id": params.MustInt("project_id"),
			"e_task_id":    params.MustInt("task_id"),
			"e_master_id":  params.MustString("master_id"),
		}
		elfRes := map[string]interface{}{}
		if tempElfRes, elfErr := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_TOTAL_ELF).AddCondition(queryElfParams).GetOne(); elfErr == nil {
			elfRes = *tempElfRes
			delete(elfRes, "_id")
			delete(elfRes, "master_id")
			delete(elfRes, "project_id")
			delete(elfRes, "task_id")
		}
		initRes := map[string]interface{}{}
		if tempInitRes, initErr := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_TOTAL_INIT).AddCondition(queryElfParams).GetOne(); initErr == nil {
			initRes = *tempInitRes
			delete(initRes, "_id")
			delete(initRes, "master_id")
			delete(initRes, "project_id")
			delete(initRes, "task_id")
		}
		result := qmap.QM{
			"basic": res,
			"elf":   elfRes,
			"init":  initRes,
		}
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/*
 * apiType http
 * @api {get} /api/v1/firmware/category_detail 扫描检测结果类别详情
 * @apiVersion 1.0.1
 * @apiName CategoryDetail
 * @apiGroup Firmware
 *
 * @apiDescription 扫描检测结果类别详情
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			project_id      	项目ID
 * @apiParam {int}      			task_id       	    任务ID
 * @apiParam {string}      			master_id       	主键ID
 *
 */
func (this FirmwareController) CategoryDetail(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	queryParams := qmap.QM{
		"e_project_id": params.MustInt("project_id"),
		"e_task_id":    params.MustInt("task_id"),
	}

	mgoSession := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_CATEGORY).AddCondition(queryParams)
	mgoSession.SetTransformFunc(this.FirmWareIDFormat)
	if res, err := mgoSession.GetOne(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this FirmwareController) FirmWareBasicFormat(data qmap.QM) qmap.QM {

	data["id"] = data["_id"]
	delete(data, "_id")
	delete(data, "status")
	delete(data, "tmp_hd_file_path")
	delete(data, "temp_file")

	return data
}

/*
 * apiType http
 * @api {get} /api/v1/firmware/analysis_detail 扫描检测结果详情
 * @apiVersion 1.0.1
 * @apiName AnalysisDetail
 * @apiGroup Firmware
 *
 * @apiDescription 扫描检测结果类别详情
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			project_id      	项目ID
 * @apiParam {int}      			task_id       	    任务ID
 * @apiParam {string}      			master_id       	主键ID
 *
 */
/**
 * 扫描检测结果详情 暂时无用处
 */
func (this FirmwareController) AnalysisDetail(ctx *gin.Context) {

	data := map[string]interface{}{
		"id":                   "5f681e3589e90ff95d7bf66b",
		"project_id":           105,
		"task_id":              102,
		"template_id":          77,
		"init_files":           map[string]interface{}{},
		"elf_scanner":          map[string]interface{}{},
		"binary_hardening":     map[string]interface{}{},
		"symbols_xrefs":        map[string]interface{}{},
		"version_scanner":      map[string]interface{}{},
		"certificates_scanner": map[string]interface{}{},
		"leaks_scanner":        map[string]interface{}{},
		"password_scanner":     map[string]interface{}{},
	}
	response.RenderSuccess(ctx, data)
}

/*
 * apiType http
 * @api {get} /api/v1/firmware/analysis_download 扫描检测结果下载
 * @apiVersion 1.0.1
 * @apiName AnalysisDetail
 * @apiGroup Firmware
 *
 * @apiDescription 扫描检测结果下载
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			project_id      	项目ID
 * @apiParam {int}      			task_id       	    任务ID
 * @apiParam {string}      			master_id       	主键ID
 *
 */
func (this FirmwareController) AnalysisDownload(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	queryParams := qmap.QM{
		"e_project_id": params.MustInt("project_id"),
		"e_task_id":    params.MustInt("task_id"),
		// "e__id":  bson.ObjectIdHex(req.MustString("master_id")),
	}

	mgoSession := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_SOURCE).AddCondition(queryParams)
	mgoSession.SetTransformFunc(this.FirmWareDownloadFormat)
	if res, err := mgoSession.GetOne(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this FirmwareController) FirmWareDownloadFormat(data qmap.QM) qmap.QM {

	data["id"] = data["_id"]
	content, _ := base64.StdEncoding.DecodeString(data["source_data"].(string))
	tmpJson := map[string]interface{}{}
	json.Unmarshal(content, &tmpJson)
	data["content"] = tmpJson
	delete(data, "_id")
	delete(data, "source_data")

	return data
}

/*
 * apiType http
 * @api {get} /api/v1/firmware/analysis_detail_page 扫描检测结果分页数据
 * @apiVersion 1.0.1
 * @apiName AnalysisDetailPage
 * @apiGroup Firmware
 *
 * @apiDescription 扫描检测结果分页数据
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			project_id      	项目ID
 * @apiParam {int}      			task_id       	    任务ID
 * @apiParam {string}      			master_id       	主键ID
 * @apiParam {int}      			current_page       	当前页
 * @apiParam {int}      			count       	    单页数量
 * @apiParam {string}      			type       			主类型
 * @apiParam {string}      			type_child       	子类型
 *
 */
/**
 * 扫描检测结果分页数据
 */
func (this FirmwareController) AnalysisDetailPage(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	var tableName string
	var childFieldVal string
	var searchOperationsRegex []bson.M
	var searchOperationsMatch bson.M
	var groupOperationsMatch bson.M
	var groupOperationsFiled bson.M
	var operations []bson.M
	var groupOperations []bson.M
	isSearch := false
	parentType := params.MustString("type")
	childType := params.MustString("type_child")
	projectId := params.MustInt("project_id")
	taskId := params.MustInt("task_id")
	search := params.MustString("search")
	cerCount := params.Int("cer_count")
	// 二进制数据统计
	binaryType := params.Int("binary_type")
	matchCountAll := params.Int("match_count_all")
	cveType := params.Int("cve_type")
	if binaryType == 0 {
		binaryType = 3 // 全部
	} else if binaryType == 1 {
		binaryType = 1 // 有问题 数据库真是字段值为1
	} else if binaryType == 2 {
		binaryType = 0 // 无问题 数据库真是字段值为0
	} else {
		response.RenderFailure(ctx, errors.New("binary_type error"))
	}

	if parentType == "leaks_scanner" {
		tableName = common.MC_FIRMWARE_RTS_LEAKS
		if childType == "url" {
			childFieldVal = "url"
		} else if childType == "email" {
			childFieldVal = "email"
		} else if childType == "ipv4_public" {
			childFieldVal = "ipv4_public"
		} else if childType == "ipv4_private" {
			childFieldVal = "ipv4_private"
		} else {
			response.RenderFailure(ctx, errors.New("child type error 1"))
		}
		if search != "" {
			isSearch = true
			searchOperationsRegex = []bson.M{
				{
					"info": bson.M{
						"$regex": search,
					},
				},
			}
		}
	} else if parentType == "certificates_scanner" {
		tableName = common.MC_FIRMWARE_RTS_CERT
		if childType == "unknown" {
			childFieldVal = "unknown"
		} else if childType == "public_key" {
			childFieldVal = "public_key"
		} else if childType == "certificate" {
			childFieldVal = "certificate"
			if search != "" {
				isSearch = true
				searchOperationsRegex = []bson.M{
					{
						"info.json.Issuer": bson.M{
							"$regex": search,
						},
					},
					{
						"info.json.Subject": bson.M{
							"$regex": search,
						},
					},
				}
			}

		} else if childType == "private_key" {
			childFieldVal = "private_key"
		} else {
			response.RenderFailure(ctx, errors.New("child type error 2"))
		}

	} else if parentType == "binary_hardening" {
		tableName = common.MC_FIRMWARE_RTS_BINARY
		if childType == "binary_hardening" {
			childFieldVal = "binary_hardening"
			// 二进制数据筛选
			if search != "" {
				isSearch = true
				searchOperationsRegex = []bson.M{
					{
						"file_name": bson.M{
							"$regex": search,
						},
					},
				}
			}
		} else {
			response.RenderFailure(ctx, errors.New("child type error 3"))
		}

	} else if parentType == "version_scanner" {
		tableName = common.MC_FIRMWARE_RTS_CVE
		if childType == "version_scanner" {
			childFieldVal = "version_scanner"
			if search != "" {
				isSearch = true
				searchOperationsRegex = []bson.M{
					{
						"cve": bson.M{
							"$regex": search,
						},
					},
					{
						"vector": bson.M{
							"$regex": search,
						},
					},
					{
						"vendor": bson.M{
							"$regex": search,
						},
					},
				}
			}
		} else {
			response.RenderFailure(ctx, errors.New("child type error 4"))
		}

	} else if parentType == "password_scanner" {
		tableName = common.MC_FIRMWARE_RTS_PWD
		if childType == "printable" {
			childFieldVal = "printable"
		} else {
			response.RenderFailure(ctx, errors.New("child type error 5"))
		}

	} else if parentType == "apk_sensitive" {
		tableName = common.MC_FIRMWARE_RTS_APK_SENITIVE
		if childType == "urls" {
			childFieldVal = "urls"
		} else if childType == "ip" {
			childFieldVal = "ip"
		} else if childType == "email" {
			childFieldVal = "email"
		} else if childType == "token" {
			childFieldVal = "token"
		} else if childType == "access_key" {
			childFieldVal = "access_key"
		} else if childType == "cert" {
			childFieldVal = "cert"
		} else {
			response.RenderFailure(ctx, errors.New("child type error 5"))
		}
		if search != "" {
			isSearch = true
			searchOperationsRegex = []bson.M{
				{
					"content": bson.M{
						"$regex": search,
					},
				},
			}
		}

	} else if parentType == "linux_basic_audit" {
		tableName = common.MC_FIRMWARE_RTS_LINUX
		if childType == "linux_basic_audit" {
			childFieldVal = "linux_basic_audit"
			if search != "" {
				isSearch = true
				searchOperationsRegex = []bson.M{
					{
						"fullpath": bson.M{
							"$regex": search,
						},
					},
				}
			}
		} else {
			response.RenderFailure(ctx, errors.New("child type error 6"))
		}

	} else {
		response.RenderFailure(ctx, errors.New("type error 7"))
	}

	searchOperationsMatch = bson.M{
		"project_id": projectId,
		"task_id":    taskId,
	}

	// version_scanner
	if parentType == "linux_basic_audit" {
		// linux基线检测搜索
		if isSearch {
			searchOperationsMatch["$or"] = searchOperationsRegex
		}
	} else if parentType == "binary_hardening" {
		if isSearch {
			if binaryType == 3 && parentType == "binary_hardening" {
				searchOperationsMatch["$or"] = searchOperationsRegex
				searchOperationsMatch["type"] = childFieldVal

			} else {
				searchOperationsMatch["$or"] = searchOperationsRegex
				searchOperationsMatch["type"] = childFieldVal
				searchOperationsMatch["is_doubt"] = binaryType

			}
		} else {
			// 搜索结果条件  无search 不搜索条件下
			if binaryType == 3 && parentType == "binary_hardening" {
				searchOperationsMatch["type"] = childFieldVal

			} else {
				searchOperationsMatch["type"] = childFieldVal
				searchOperationsMatch["is_doubt"] = binaryType
			}
		}
	} else if parentType == "version_scanner" {
		if isSearch {
			if cveType == 0 {
				searchOperationsMatch["$or"] = searchOperationsRegex
				searchOperationsMatch["type"] = childFieldVal
			} else {
				searchOperationsMatch["$or"] = searchOperationsRegex
				searchOperationsMatch["type"] = childFieldVal
				searchOperationsMatch["level"] = cveType
			}
		} else {
			// 搜索结果条件  无search 不搜索条件下
			if cveType == 0 {
				searchOperationsMatch["type"] = childFieldVal

			} else {
				searchOperationsMatch["type"] = childFieldVal
				searchOperationsMatch["level"] = cveType
			}

		}
	} else {
		// 其它类型数据搜索结果条件  具备搜索条件下
		if isSearch {
			searchOperationsMatch["$or"] = searchOperationsRegex
			searchOperationsMatch["type"] = childFieldVal
		} else {
			// 搜索结果条件  无search 不搜索条件下
			searchOperationsMatch["type"] = childFieldVal

		}
	}
	// 合并数据搜索条件
	operations = []bson.M{
		{
			"$match": searchOperationsMatch,
		},
	}

	if res, err := mongo.NewMgoSession(tableName).AddUrlQueryCondition(params.String("query_params")).MATCHGetPage(operations); err == nil {
		groupOperationsMatch = bson.M{
			"project_id": projectId,
			"task_id":    taskId,
		}
		// 获取证书检测列表聚合统计数据
		if cerCount == 1 {
			if parentType == "certificates_scanner" {
				groupOperationsFiled = bson.M{"_id": "$type", "count": bson.M{"$sum": 1}}
				if isSearch && childFieldVal == "certificate" {
					// 拼装搜索条件统计
					groupOperationsMatch["$or"] = searchOperationsRegex
				}
				groupOperations = []bson.M{
					{"$match": groupOperationsMatch},
					{"$group": groupOperationsFiled},
				}
				tempGroupRts, _ := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_CERT).MATCHALL(groupOperations)

				res.Merge(*tempGroupRts)
			} else if parentType == "leaks_scanner" {
				groupOperationsFiled = bson.M{"_id": "$type", "count": bson.M{"$sum": 1}}
				if isSearch {
					// 拼装搜索条件统计
					groupOperationsMatch["$or"] = searchOperationsRegex
					groupOperationsMatch["type"] = childFieldVal

				}
				groupOperations = []bson.M{
					{"$match": groupOperationsMatch},
					{"$group": groupOperationsFiled},
				}
				tempGroupRts, _ := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_LEAKS).MATCHALL(groupOperations)

				res.Merge(*tempGroupRts)
			} else if parentType == "apk_sensitive" {
				groupOperationsFiled = bson.M{"_id": "$type", "count": bson.M{"$sum": 1}}
				if isSearch {
					// 拼装搜索条件统计
					groupOperationsMatch["$or"] = searchOperationsRegex
					groupOperationsMatch["type"] = childFieldVal

				}
				groupOperations = []bson.M{
					{"$match": groupOperationsMatch},
					{"$group": groupOperationsFiled},
				}
				tempGroupRts, _ := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_APK_SENITIVE).MATCHALL(groupOperations)
				res.Merge(*tempGroupRts)
			} else if parentType == "binary_hardening" {
				groupOperationsFiled = bson.M{"_id": "$is_doubt", "count": bson.M{"$sum": 1}}
				if isSearch {
					// 拼装包含搜索条件统计
					if binaryType == 3 {
						groupOperationsMatch["$or"] = searchOperationsRegex
						groupOperationsMatch["type"] = childFieldVal

					} else {
						groupOperationsMatch["$or"] = searchOperationsRegex
						groupOperationsMatch["type"] = childFieldVal
						groupOperationsMatch["is_doubt"] = binaryType
						if matchCountAll == 1 {
							// 强制统计全部
							delete(groupOperationsMatch, "is_doubt")
						}
					}
				} else {
					// 拼装不包含搜索条件统计
					if binaryType == 3 {
						groupOperationsMatch["type"] = childFieldVal
					} else {
						groupOperationsMatch["type"] = childFieldVal
						groupOperationsMatch["is_doubt"] = binaryType
						if matchCountAll == 1 {
							// 强制统计全部
							delete(groupOperationsMatch, "is_doubt")
						}
					}

				}
				groupOperations = []bson.M{
					{"$match": groupOperationsMatch},
					{"$group": groupOperationsFiled},
				}
				tempGroupRts, _ := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_BINARY).MATCHALL(groupOperations)

				secondGroupRts := *tempGroupRts
				resultTemplate := []map[string]interface{}{
					{"_id": 1, "count": 0},
					{"_id": 0, "count": 0},
				}

				matchAllResult := resetResultMatchAll(resultTemplate, secondGroupRts)
				res.Merge(matchAllResult)

			} else if parentType == "version_scanner" {
				groupOperationsFiled = bson.M{"_id": "$level", "count": bson.M{"$sum": 1}}
				if isSearch {
					// 拼装包含搜索条件统计
					if cveType == 0 {
						groupOperationsMatch["$or"] = searchOperationsRegex
						groupOperationsMatch["type"] = childFieldVal
					} else {
						groupOperationsMatch["$or"] = searchOperationsRegex
						groupOperationsMatch["type"] = childFieldVal
						groupOperationsMatch["level"] = cveType
						if matchCountAll == 1 {
							// 强制统计全部
							delete(groupOperationsMatch, "level")
						}

					}
				} else {
					// 拼装不包含搜索条件统计
					if cveType == 0 {
						groupOperationsMatch["type"] = childFieldVal
					} else {
						groupOperationsMatch["type"] = childFieldVal
						groupOperationsMatch["level"] = cveType
						if matchCountAll == 1 {
							// 强制统计全部
							delete(groupOperationsMatch, "level")
						}
					}

				}
				groupOperations = []bson.M{
					{"$match": groupOperationsMatch},
					{"$group": groupOperationsFiled},
				}
				tempGroupRts, _ := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_CVE).MATCHALL(groupOperations)
				// 特殊处理CVE数据统计结果数据
				secondGroupRts := *tempGroupRts
				resultTemplate := []map[string]interface{}{
					{"_id": 1, "count": 0},
					{"_id": 2, "count": 0},
					{"_id": 3, "count": 0},
					{"_id": 4, "count": 0},
				}
				matchAllResult := resetResultMatchAll(resultTemplate, secondGroupRts)
				res.Merge(matchAllResult)
			}
		}
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func resetResultMatchAll(resultTemplate []map[string]interface{}, secondGroupRts qmap.QM) qmap.QM {
	resultMatchAll := []map[string]interface{}{}
	for _, val := range resultTemplate {
		tempID := val["_id"].(int)
		tempSetStatus := false
		for _, v := range secondGroupRts["match_all"].([]map[string]interface{}) {
			tempDbID := v["_id"].(int)
			tempDbCount := v["count"]
			if tempDbID == tempID {
				resultMatchAll = append(resultMatchAll, map[string]interface{}{"_id": tempDbID, "count": tempDbCount})
				tempSetStatus = true
			}
		}
		if tempSetStatus == false {
			resultMatchAll = append(resultMatchAll, map[string]interface{}{"_id": tempID, "count": 0})
		}
	}
	return map[string]interface{}{"match_all": resultMatchAll}
}

/*
 * apiType http
 * @api {get} /api/v1/firmware/analysis_tracker 扫描检测结果统计数据
 * @apiVersion 1.0.1
 * @apiName AnalysisTracker
 * @apiGroup Firmware
 *
 * @apiDescription 扫描检测结果统计数据
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			project_id      	项目ID
 * @apiParam {int}      			task_id       	    任务ID
 * @apiParam {string}      			master_id       	主键ID
 * @apiParam {string}      			type       			主类型
 *
 */
func (this FirmwareController) AnalysisTracker(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	queryParams := qmap.QM{
		"e_project_id": params.MustInt("project_id"),
		"e_task_id":    params.MustInt("task_id"),
	}
	cveRes := map[string]interface{}{}
	if tempCveRes, cveErr := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_TOTAL_CVE).AddCondition(queryParams).SetTransformFunc(this.FirmWareIDFormat).GetOne(); cveErr == nil {
		cveRes = *tempCveRes
		delete(cveRes, "master_id")
		delete(cveRes, "project_id")
		delete(cveRes, "task_id")
		delete(cveRes, "_id")
		delete(cveRes, "create_time")
		delete(cveRes, "update_time")
	}

	binRes := map[string]interface{}{}
	if tempBinRes, binErr := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_TOTAL_BINARY).AddCondition(queryParams).SetTransformFunc(this.FirmWareIDFormat).GetOne(); binErr == nil {
		binRes = *tempBinRes
		delete(binRes, "master_id")
		delete(binRes, "project_id")
		delete(binRes, "task_id")
		delete(binRes, "_id")
		delete(binRes, "create_time")
		delete(binRes, "update_time")
	}

	certRes := map[string]interface{}{}
	if tempCertRes, certErr := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_TOTAL_CERT).AddCondition(queryParams).SetTransformFunc(this.FirmWareIDFormat).GetOne(); certErr == nil {
		certRes = *tempCertRes
		delete(certRes, "master_id")
		delete(certRes, "project_id")
		delete(certRes, "task_id")
		delete(certRes, "_id")
		delete(certRes, "create_time")
		delete(certRes, "update_time")
	}

	riskRes := map[string]interface{}{}
	if tempRiskRes, riskErr := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_TOTAL_RISK).AddCondition(queryParams).SetTransformFunc(this.FirmWareIDFormat).GetOne(); riskErr == nil {

		riskRes = *tempRiskRes
		delete(riskRes, "list")
		delete(riskRes, "master_id")
		delete(riskRes, "project_id")
		delete(riskRes, "task_id")
		delete(riskRes, "_id")
		delete(riskRes, "create_time")
		delete(riskRes, "update_time")
	}

	linuxRes := map[string]interface{}{}
	if templinuxRes, linuxErr := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_TOTAL_LINUX).AddCondition(queryParams).SetTransformFunc(this.FirmWareIDFormat).GetOne(); linuxErr == nil {
		linuxRes = *templinuxRes
		delete(linuxRes, "master_id")
		delete(linuxRes, "project_id")
		delete(linuxRes, "task_id")
		delete(linuxRes, "_id")
		delete(linuxRes, "id")
		delete(linuxRes, "create_time")
		delete(linuxRes, "update_time")
	}

	levelRes := map[string]interface{}{}
	levelQueryParams := qmap.QM{
		"e_project_id": params.MustInt("project_id"),
		"e_task_id":    params.MustInt("task_id"),
	}
	mongoSession := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_TOTAL_LEVEL).AddCondition(levelQueryParams)
	mongoSession.AddSorter("level", 1)
	if tempLevelRes, levelErr := mongoSession.Get(); levelErr == nil {

		if len(*tempLevelRes) > 0 {
			levelRes = (*tempLevelRes)[0]
			delete(levelRes, "master_id")
			delete(levelRes, "project_id")
			delete(levelRes, "task_id")
			delete(levelRes, "_id")
			delete(levelRes, "id")
			delete(levelRes, "create_time")
			delete(levelRes, "update_time")
		}

	} else {
		levelRes = map[string]interface{}{
			"level": 0,
		}
	}
	result := &qmap.QM{
		"cve":               cveRes,
		"binary":            binRes,
		"cert":              certRes,
		"risk":              riskRes,
		"linux_basic_audit": linuxRes,
		"firmware":          levelRes,
	}
	response.RenderSuccess(ctx, result)
}

func (this FirmwareController) FirmWareIDFormat(data qmap.QM) qmap.QM {

	data["id"] = data["_id"]
	delete(data, "_id")
	return data
}

/*
 * apiType http
 * @api {get} /api/v1/firmware/apk_basic_info 扫描检测结果统计数据
 * @apiVersion 1.0.1
 * @apiName AnalysisTracker
 * @apiGroup Firmware
 *
 * @apiDescription 扫描检测结果统计数据
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			project_id      	项目ID
 * @apiParam {int}      			task_id       	    任务ID
 * @apiParam {string}      			master_id       	主键ID
 * @apiParam {string}      			type       			主类型
 *
 */
func (this FirmwareController) ApkBasicInfo(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	queryParams := qmap.QM{
		"e_project_id": params.MustInt("project_id"),
		"e_task_id":    params.MustInt("task_id"),
	}

	resultRes := map[string]interface{}{}
	cveRes := map[string]interface{}{}
	if tempCveRes, cveErr := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_APK).AddCondition(queryParams).SetTransformFunc(this.FirmWareIDFormat).GetOne(); cveErr == nil {

		cveRes = *tempCveRes
		resultRes["app_name"], _ = cveRes["app_name"].(string)
		resultRes["apkfile_name"], _ = cveRes["apkfile_name"].(string)
		resultRes["version_name"], _ = cveRes["version_name"].(string)
		resultRes["version_code"], _ = cveRes["version_code"].(string)
		resultRes["debuggable"], _ = cveRes["debuggable"].(string)
		resultRes["is_protect"], _ = cveRes["is_protect"].(string)
		resultRes["allowBackup"], _ = cveRes["allowBackup"].(string)
		resultRes["dex_protect"], _ = cveRes["dex_protect"].(string)
		resultRes["apkfile_md5"], _ = cveRes["apkfile_md5"]
		resultRes["apkfile_sha1"], _ = cveRes["apkfile_sha1"]
		resultRes["apkfile_sha256"], _ = cveRes["apkfile_sha256"]
		resultRes["min_sdk_version"], _ = cveRes["min_sdk_version"]
		resultRes["target_sdk_version"], _ = cveRes["target_sdk_version"]
		resultRes["created_at"], _ = cveRes["created_at"]
		resultRes["updated_at"], _ = cveRes["updated_at"]
		resultRes["sign_serial_Number"] = cveRes["sign"].(map[string]interface{})["serial_Number"].(string)
		resultRes["sign_valid_not_after"] = cveRes["sign"].(map[string]interface{})["valid_not_after"].(string)
		resultRes["sign_valid_not_before"] = cveRes["sign"].(map[string]interface{})["valid_not_before"].(string)
		resultRes["receiver"] = cveRes["receiver"].(map[string]interface{})
		resultRes["activity"] = cveRes["activity"].(map[string]interface{})
		resultRes["provider"] = cveRes["provider"].(map[string]interface{})
		resultRes["service"] = cveRes["service"].(map[string]interface{})
		for i := 1; i <= 10; i++ {
			key := fmt.Sprintf("sign_v%d", i)
			if signVal, isSetKey := cveRes["sign"].(map[string]interface{})[key].(string); isSetKey == true {
				resultRes[key] = signVal
			} else {
				break
			}
		}
	}
	response.RenderSuccess(ctx, resultRes)
}

/*
 * apiType http
 * @api {get} /api/v1/firmware/apk_common_vue 扫描检测结果统计数据
 * @apiVersion 1.0.1
 * @apiName AnalysisTracker
 * @apiGroup Firmware
 *
 * @apiDescription 扫描检测结果统计数据
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			project_id      	项目ID
 * @apiParam {int}      			task_id       	    任务ID
 * @apiParam {string}      			master_id       	主键ID
 * @apiParam {string}      			type       			主类型
 *
 */
func (this FirmwareController) ApkCommonVue(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	operations := []bson.M{
		{
			"$match": bson.M{
				"project_id":  params.MustInt("project_id"),
				"task_id":     params.MustInt("task_id"),
				"parent_type": params.MustString("type"),
			},
		},
		{
			"$project": bson.M{"detail": 0},
		},
	}

	if res, err := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_APK_VUL).AddUrlQueryCondition(params.String("query_params")).QueryGet(operations); err == nil {
		result := *res
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/*
 * apiType http
 * @api {get} /api/v1/firmware/apk_common_vue_detail 扫描检测结果统计数据
 * @apiVersion 1.0.1
 * @apiName ApkCommonVueDetail
 * @apiGroup Firmware
 *
 * @apiDescription 扫描检测结果统计数据
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			project_id      	项目ID
 * @apiParam {int}      			task_id       	    任务ID
 * @apiParam {string}      			master_id       	主键ID
 *
 */
func (this FirmwareController) ApkCommonVueDetail(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))

	operations := []bson.M{
		{
			"$match": bson.M{
				"project_id": params.MustInt("project_id"),
				"task_id":    params.MustInt("task_id"),
				"_id":        bson.ObjectIdHex(params.String("master_id")),
			},
		},
		{
			"$project": bson.M{"detail": 1},
		},
	}

	if res, err := mongo.NewMgoSession(common.MC_FIRMWARE_RTS_APK_VUL).QueryGet(operations); err == nil {
		result := *res
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

// ################# 固件检测V2 #########################
/*
 * apiType http
 * @api {get} /api/v2/firmware/basic 获取固件基本信息
 * @apiVersion 1.0.1
 * @apiName Basic
 * @apiGroup Firmware
 *
 * @apiDescription 获取固件基本信息
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			task_id       	    任务ID
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "basic": {
 *             "create_time": 0,
 *             "device_model": "NTG7",
 *             "device_name": "奔驰⻋机1",
 *             "device_type": "IVI",
 *             "file_id": "61499cd224b647313abcba4b",
 *             "firmware_md5": "",
 *             "firmware_name": "",
 *             "firmware_size": 0,
 *             "firmware_version": "20180001",
 *             "id": 1,
 *             "name": "固件扫描任务TEST16",
 *             "progress": 0,
 *             "source_report": "",
 *             "status": 8,
 *             "task_id": 1,
 *             "task_name": "",
 *             "template_id": 71,
 *             "template_name": "通用IoT固件检测模板",
 *             "update_time": 0,
 *             "yafaf_download_path": "remotedown/1632304364_rootfs.zip",
 *             "yafaf_id": 6,
 *             "yafaf_project_id": 28
 *         },
 *         "elf": {
 *             "create_time": 1632448855,
 *             "executable": "/workspace/固件扫描任务TEST16_2021-09-221632304368/extracted/rootfs/app/lib/libwebsdk_postevent.so",
 *             "id": 1,
 *             "scanner_id": 1,
 *             "type": "shared_lib"
 *         },
 *         "init": {
 *             "dir_num": 239,
 *             "file_num": 729,
 *             "link_num": 75,
 *             "node_num": 0
 *         }
 *     },
 *     "msg": ""
 * }
 */
func (this FirmwareController) BasicV2(ctx *gin.Context) {
	params := qmap.QM{
		"e_id": request.QueryInt(ctx, "task_id"),
	}
	if has, item := sys_service.NewSession().AddCondition(params).SetTransformFunc(this.BasicFormatV2Transform).GetOne(new(mysql_model.FirmwareTask)); has {
		response.RenderSuccess(ctx, item)
	} else {
		panic(errors.New("Item not found"))
	}
}

func (this FirmwareController) BasicFormatV2Transform(data qmap.QM) qmap.QM {
	data["source_report"] = ""
	result := qmap.QM{
		"basic": data,
	}
	id := data.Int("id")
	params := qmap.QM{
		"e_scanner_id": id,
	}
	if has, elf := sys_service.NewSession().AddCondition(params).GetOne(new(mysql_model.FirmwareReportRtsElf)); has {
		result["elf"] = elf
	} else {
		result["elf"] = qmap.QM{}
	}
	_, init := sys_service.NewSession().AddCondition(params).GetOne(new(mysql_model.FirmwareReportInit))
	result["init"] = qmap.QM{
		"dir_num":  init.Int("dir_num"),
		"file_num": init.Int("file_num"),
		"link_num": init.Int("link_num"),
		"node_num": init.Int("node_num"),
	}
	return result
}

/*
 * apiType http
 * @api {get} /api/v2/firmware/category_detail 扫描检测结果类别详情
 * @apiVersion 1.0.1
 * @apiName CategoryDetailV2
 * @apiGroup Firmware
 *
 * @apiDescription 扫描检测结果类别详情
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}			task_id       	    任务ID
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "apk_common_vul": 0,
 *         "apk_info": 0,
 *         "apk_sensitive_info": 0,
 *         "binary_hardening": 1,
 *         "certificates_scanner": 1,
 *         "elf_scanner": 1,
 *         "id": 1,
 *         "init_files": 1,
 *         "is_elf": 0,
 *         "leaks_scanner": 1,
 *         "linux_basic_audit": 1,
 *         "password_scanner": 1,
 *         "scanner_id": 0,
 *         "symbols_xrefs": 0,
 *         "template_id": 0,
 *         "version_scanner": 1
 *     },
 *     "msg": ""
 * }
 */
func (this FirmwareController) CategoryDetailV2(ctx *gin.Context) {
	queryParams := qmap.QM{
		"e_scanner_id": request.QueryInt(ctx, "task_id"),
	}
	if has, item := sys_service.NewSessionWithCond(queryParams).GetOne(new(mysql_model.FirmwareReportRtsCategory)); has {
		response.RenderSuccess(ctx, item)
	} else {
		response.RenderSuccess(ctx, nil)
	}
}

/*
 * apiType http
 * @api {get} /api/v2/firmware/analysis_tracker 扫描检测结果统计数据
 * @apiVersion 1.0.1
 * @apiName AnalysisTrackerV2
 * @apiGroup Firmware
 *
 * @apiDescription 扫描检测结果统计数据
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			task_id       	    任务ID
 * @apiParam {string}      			type       			主类型
 *
 * @apiSuccessExample {json} 请求成功示例:
 *{
 *    "code": 0,
 *    "data": {
 *        "binary": {
 *            "create_time": 1632448856,
 *            "file_name": "bt_daemon",
 *            "full_path": "/workspace/固件扫描任务TEST16_2021-09-221632304368/extracted/rootfs/app/bin/bt_daemon",
 *            "hardenable": 0,
 *            "id": 1,
 *            "is_doubt": 1,
 *            "is_elf": 1,
 *            "magic_info": "ELF 32-bit LSB shared object, ARM, EABI5 version 1 (SYSV), dynamically linked, interpreter /system/bin/linker, stripped",
 *            "rela_path": "/rootfs/app/bin",
 *            "result": "{\"canary\":\"yes\",\"debugsym\":\"no\",\"fortified\":\"1\",\"fortify_able\":\"5\",\"fortify_source\":\"yes\",\"nx\":\"yes\",\"pie\":\"yes\",\"relro\":\"full\",\"rpath\":\"no\",\"runpath\":\"no\",\"stripped\":\"yes\"}",
 *            "scanner_id": 1,
 *            "type": "binary_hardening"
 *        },
 *        "cert": {
 *            "content": "-----BEGIN CERTIFICATE-----\nMIIEtTCCA52gAwIBAgIJAKd3ZJUaqxk4MA0GCSqGSIb3DQEBBQUAMIGYMQswCQYD\nVQQGEwJDTjESMBAGA1UECBMJR3Vhbmdkb25nMRIwEAYDVQQHEwlTaGVuZ3poZW4x\nDzANBgNVBAoTBkh1YXdlaTEYMBYGA1UECxMPVGVybWluYWxDb21wYW55MRQwEgYD\nVQQDEwtBbmRyb2lkVGVhbTEgMB4GCSqGSIb3DQEJARYRbW9iaWxlQGh1YXdlaS5j\nb20wHhcNMTEwNTI1MDYxMDQ4WhcNMzYwNTE4MDYxMDQ4WjCBmDELMAkGA1UEBhMC\nQ04xEjAQBgNVBAgTCUd1YW5nZG9uZzESMBAGA1UEBxMJU2hlbmd6aGVuMQ8wDQYD\nVQQKEwZIdWF3ZWkxGDAWBgNVBAsTD1Rlcm1pbmFsQ29tcGFueTEUMBIGA1UEAxML\nQW5kcm9pZFRlYW0xIDAeBgkqhkiG9",
 *            "create_time": 1632448859,
 *            "file_name": "testkey.x509.pem",
 *            "id": 1,
 *            "info": "",
 *            "path": "/rootfs/system/etc/security/otacerts.zip_unpacked/testkey.x509.pem",
 *            "scanner_id": 1,
 *            "type": "certificate"
 *        },
 *        "cve": {
 *            "create_time": 1632448859,
 *            "cve": "",
 *            "cvssv2": "{\"cvssV2\":{\"accessComplexity\":\"MEDIUM\",\"accessVector\":\"NETWORK\",\"authentication\":\"NONE\",\"availabilityImpact\":\"PARTIAL\",\"baseScore\":6.8,\"confidentialityImpact\":\"PARTIAL\",\"integrityImpact\":\"PARTIAL\",\"vectorString\":\"AV:N/AC:M/Au:N/C:P/I:P/A:P\",\"version\":\"2.0\"},\"exploitabilityScore\":8.6,\"impactScore\":6.4,\"obtainAllPrivilege\":false,\"obtainOtherPrivilege\":false,\"obtainUserPrivilege\":false,\"severity\":\"MEDIUM\",\"userInteractionRequired\":false}",
 *            "cvssv2score": 8,
 *            "cvssv3": "{\"cvssV3\":{\"attackComplexity\":\"HIGH\",\"attackVector\":\"NETWORK\",\"availabilityImpact\":\"HIGH\",\"baseScore\":8.1,\"baseSeverity\":\"HIGH\",\"confidentialityImpact\":\"HIGH\",\"integrityImpact\":\"HIGH\",\"privilegesRequired\":\"NONE\",\"scope\":\"UNCHANGED\",\"userInteraction\":\"NONE\",\"vectorString\":\"CVSS:3.1/AV:N/AC:H/PR:N/UI:N/S:U/C:H/I:H/A:H\",\"version\":\"3.1\"},\"exploitabilityScore\":2.2,\"impactScore\":5.9}",
 *            "description": "[{\"lang\":\"en\",\"value\":\"Busybox contains a Missing SSL certificate validation vulnerability in The \\\"busybox wget\\\" applet that can result in arbitrary code execution. This attack appear to be exploitable via Simply download any file over HTTPS using \\\"busybox wget https://compromised-domain.com/important-file\\\".\"}]",
 *            "file_name": "busybox",
 *            "id": 1,
 *            "level": 3,
 *            "path": "/workspace/固件扫描任务TEST16_2021-09-221632304368/extracted/rootfs/system/bin/busybox",
 *            "scanner_id": 1,
 *            "type": "version_scanner",
 *            "vector": "NETWORK",
 *            "vendor": "",
 *            "version": "1.21.1",
 *            "version_end_excluding": "1.32.0",
 *            "version_end_including": "",
 *            "version_start_excluding": "",
 *            "version_start_including": ""
 *        },
 *        "firmware": {
 *            "level": 0
 *        },
 *        "linux_basic_audit": {
 *            "create_time": 1632448860,
 *            "detail": "{\"checkname\":\"检查syslog远程日志功能配置1\",\"conf\":\"Not Found\",\"details\":\"远程日志可将系统日志发送到另一台服务器上便于统一管理,当前系统使用syslog记录日志，但在syslog.conf配置文件中未发现远程日志IP主机配置项\",\"hardenable\":true,\"repair\":\"建议根据生产需求开启远程日志功能，syslog.conf示例配置规则1:*.*   @192.168.1.10:514。(IP/Domain配置2选一即可)\"}",
 *            "full_path": "/rootfs/cust/app/etc/syslog-ng.conf",
 *            "id": 1,
 *            "scanner_id": 1,
 *            "type": "check_logging_service"
 *        },
 *        "risk": {
 *            "binary_count": 199,
 *            "id": 1,
 *            "linux_count": 6,
 *            "over_cert_count": 2,
 *            "pass_risk_count": 29,
 *            "risk_count": 435,
 *            "risk_suspect_count": 281,
 *            "scanner_id": 1
 *        }
 *    },
 *    "msg": ""
 *}
 */
func (this FirmwareController) AnalysisTrackerV2(ctx *gin.Context) {
	queryParams := qmap.QM{
		"e_scanner_id": request.QueryInt(ctx, "task_id"),
	}
	result := qmap.QM{}
	if has, item := sys_service.NewSessionWithCond(queryParams).GetOne(new(mysql_model.FirmwareReportRtsCveTotal)); has {
		result["cve"] = item
	} else {
		result["cve"] = qmap.QM{}
	}
	if has, item := sys_service.NewSessionWithCond(queryParams).GetOne(new(mysql_model.FirmwareReportRtsBinaryTotal)); has {
		delete(*item, "scanner_id")
		result["binary"] = item
	} else {
		result["binary"] = qmap.QM{}
	}
	if has, item := sys_service.NewSessionWithCond(queryParams).GetOne(new(mysql_model.FirmwareReportRtsCertTotal)); has {
		result["cert"] = item
	} else {
		result["cert"] = qmap.QM{}
	}
	if has, item := sys_service.NewSessionWithCond(queryParams).GetOne(new(mysql_model.FirmwareReportRtsRisk)); has {
		result["risk"] = item
	} else {
		result["risk"] = qmap.QM{}
	}
	if has, item := sys_service.NewSessionWithCond(queryParams).GetOne(new(mysql_model.FirmwareReportRtsLinuxTotal)); has {
		result["linux_basic_audit"] = item
	} else {
		result["linux_basic_audit"] = qmap.QM{}
	}
	if has, item := sys_service.NewSessionWithCond(queryParams).GetOne(new(mysql_model.FirmwareReportRtsApkLevel)); has {
		result["firmware"] = item
	} else {
		result["firmware"] = qmap.QM{}
	}
	response.RenderSuccess(ctx, result)
}

/*
 * apiType http
 * @api {get} /api/v2/firmware/analysis_detail_page 扫描检测结果分页数据
 * @apiVersion 1.0.1
 * @apiName AnalysisDetailPageV2
 * @apiGroup Firmware
 *
 * @apiDescription 扫描检测结果分页数据
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			task_id       	    任务ID
 * @apiParam {int}      			current_page       	当前页
 * @apiParam {int}      			count       	    单页数量
 * @apiParam {string}      			type       			主类型
 * @apiParam {string}      			type_child       	子类型
 *
 */
func (this FirmwareController) AnalysisDetailPageV2(ctx *gin.Context) {
	switch request.QueryString(ctx, "type") {
	case "leaks_scanner":
		this.leaksScannerAnalysisPageV2(ctx)
	case "certificates_scanner":
		this.certificatesScannerAnalysisPageV2(ctx)
	case "binary_hardening":
		this.binaryHardeningScannerAnalysisPageV2(ctx)
	case "version_scanner":
		this.versionScannerAnalysisPageV2(ctx)
	case "password_scanner":
		this.passwordScannerAnalysisPageV2(ctx)
	case "linux_basic_audit":
		this.linuxBasicAuditScannerAnalysisPageV2(ctx)
	case "apk_sensitive":
		this.apkSensitiveScannerAnalysisPageV2(ctx)
	default:
		response.RenderFailure(ctx, errors.New("unknown type "+request.QueryString(ctx, "type")))
	}
	return
}

func (this FirmwareController) leaksScannerAnalysisPageV2(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))
	(*params)["e_scanner_id"] = request.QueryInt(ctx, "task_id")
	if typeChild := request.DefaultQueryString(ctx, "type_child", ""); typeChild != "" {
		(*params)["e_type"] = typeChild
	}
	if search := request.DefaultQueryString(ctx, "search", ""); search != "" {
		(*params)["lk_origin"] = search
	}
	if result, err := sys_service.NewSessionWithCond(*params).GetPage(new(mysql_model.FirmwareReportRtsLeaks), &[]mysql_model.FirmwareReportRtsLeaks{}); err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this FirmwareController) certificatesScannerAnalysisPageV2(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))
	(*params)["e_scanner_id"] = request.QueryInt(ctx, "task_id")
	if typeChild := request.DefaultQueryString(ctx, "type_child", ""); typeChild != "" {
		(*params)["e_type"] = typeChild
	}
	if search := request.DefaultQueryString(ctx, "search", ""); search != "" {
		(*params)["lk_info"] = search
	}
	if result, err := sys_service.NewSessionWithCond(*params).SetTransformFunc(this.certificatesScannerAnalysisPageV2TransformerFunc).GetPage(new(mysql_model.FirmwareReportRtsCert), &[]mysql_model.FirmwareReportRtsCert{}); err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this FirmwareController) certificatesScannerAnalysisPageV2TransformerFunc(item qmap.QM) qmap.QM {
	if info, has := item.TryString("info"); has {
		if infoMap, err := qmap.NewWithString(info); err == nil {
			item["info"] = infoMap
		} else {
			item["info"] = qmap.QM{}
		}
	} else {
		item["info"] = qmap.QM{}
	}
	return item
}

func (this FirmwareController) binaryHardeningScannerAnalysisPageV2(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))
	(*params)["e_scanner_id"] = request.QueryInt(ctx, "task_id")

	if binaryType := request.DefaultQueryInt(ctx, "binary_type", 0); binaryType == 1 {
		(*params)["e_is_doubt"] = binaryType
	}
	if search := request.DefaultQueryString(ctx, "search", ""); search != "" {
		(*params)["lk_file_name"] = search
	}
	if result, err := sys_service.NewSessionWithCond(*params).GetPage(new(mysql_model.FirmwareReportRtsBinary), &[]mysql_model.FirmwareReportRtsBinary{}); err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this FirmwareController) versionScannerAnalysisPageV2(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))
	(*params)["e_scanner_id"] = request.QueryInt(ctx, "task_id")
	(*params)["e_type"] = "version_scanner"
	cveType := request.DefaultQueryInt(ctx, "cve_type", 0)
	if cveType > 0 {
		(*params)["e_level"] = cveType
	}
	if search := request.DefaultQueryString(ctx, "search", ""); search != "" {
		(*params)["lk_cve"] = search
		(*params)["orlk_vendor"] = search
		(*params)["orlk_vector"] = search
	}
	if result, err := sys_service.NewSessionWithCond(*params).GetPage(new(mysql_model.FirmwareReportRtsCve), &[]mysql_model.FirmwareReportRtsCve{}); err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this FirmwareController) passwordScannerAnalysisPageV2(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))
	(*params)["e_scanner_id"] = request.QueryInt(ctx, "task_id")
	if result, err := sys_service.NewSessionWithCond(*params).GetPage(new(mysql_model.FirmwareReportRtsPwd), &[]mysql_model.FirmwareReportRtsPwd{}); err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this FirmwareController) linuxBasicAuditScannerAnalysisPageV2(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))
	(*params)["e_scanner_id"] = request.QueryInt(ctx, "task_id")
	if search := request.DefaultQueryString(ctx, "search", ""); search != "" {
		(*params)["lk_full_path"] = search
	}
	if result, err := sys_service.NewSessionWithCond(*params).GetPage(new(mysql_model.FirmwareReportRtsLinux), &[]mysql_model.FirmwareReportRtsLinux{}); err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this FirmwareController) apkSensitiveScannerAnalysisPageV2(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))
	(*params)["e_scanner_id"] = request.QueryInt(ctx, "task_id")
	if typeChild := request.DefaultQueryString(ctx, "type_child", ""); typeChild != "" {
		(*params)["e_type"] = typeChild
	}
	if search := request.DefaultQueryString(ctx, "search", ""); search != "" {
		(*params)["lk_content"] = search
	}
	if result, err := sys_service.NewSessionWithCond(*params).GetPage(new(mysql_model.FirmwareReportRtsApkSensitive), &[]mysql_model.FirmwareReportRtsApkSensitive{}); err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/*
 * apiType http
 * @api {get} /api/v2/firmware/apk_common_vue 扫描检测结果统计数据
 * @apiVersion 1.0.1
 * @apiName ApkCommonVueV2
 * @apiGroup Firmware
 *
 * @apiDescription 扫描检测结果统计数据
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			task_id       	    任务ID
 * @apiParam {string}      			type       			主类型
 *
 */
func (this FirmwareController) ApkCommonVueV2(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))
	(*params)["e_scanner_id"] = params.MustInt("task_id")
	(*params)["e_parent_type"] = params.MustString("type")

	if result, err := sys_service.NewSessionWithCond(*params).GetPage(new(mysql_model.FirmwareReportRtsApkVul), &[]mysql_model.FirmwareReportRtsApkVul{}); err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/*
 * apiType http
 * @api {get} /api/v2/firmware/apk_common_vue_detail 扫描检测结果统计数据
 * @apiVersion 1.0.1
 * @apiName ApkCommonVueDetailV2
 * @apiGroup Firmware
 *
 * @apiDescription 扫描检测结果统计数据
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			task_id       	    任务ID
 *
 */
func (this FirmwareController) ApkCommonVueDetailV2(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))
	(*params)["e_scanner_id"] = params.MustInt("task_id")

	if result, err := sys_service.NewSessionWithCond(*params).GetPage(new(mysql_model.FirmwareReportRtsApkVul), &[]mysql_model.FirmwareReportRtsApkVul{}); err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/*
 * apiType http
 * @api {get} /api/v2/firmware/apk_basic_info 扫描检测结果统计数据
 * @apiVersion 1.0.1
 * @apiName ApkBasicInfoV2
 * @apiGroup Firmware
 *
 * @apiDescription 扫描检测结果统计数据
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {int}      			task_id       	    任务ID
 * @apiParam {string}      			type       			主类型
 *
 */
func (this FirmwareController) ApkBasicInfoV2(ctx *gin.Context) {
	params := &qmap.QM{
		"query_params": ctx.Request.URL.RawQuery,
	}
	*params = params.Merge(*request.GetRequestQueryParams(ctx))
	*params = params.Merge(*request.GetRequestBody(ctx))
	(*params)["e_scanner_id"] = params.MustInt("task_id")

	if has, result := sys_service.NewSessionWithCond(*params).GetOne(new(mysql_model.FirmwareReportRtsApkLevel)); has {
		if origin, err := qmap.NewWithString(result.String("original_content")); err == nil {
			(*result)["original_content"] = origin
		} else {
			(*result)["original_content"] = qmap.QM{}
		}
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, errors.New("Item not found"))
	}
}

func (this FirmwareController) Template(ctx *gin.Context) {
	item := qmap.QM{
		"id":   71,
		"name": "通用IoT固件检测模板",
	}
	items := []qmap.QM{item}
	response.RenderSuccess(ctx, items)
}
