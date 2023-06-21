package controller

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/http/transformer"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/log"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/orm_mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/common_lib/session"
	"skygo_detection/mysql_model"
)

type AssetTestPieceController struct{}

/**
 * apiType http
 * @api {get} /api/v1/asset_test_pieces 测试件列表查询
 * @apiVersion 0.1.0
 * @apiName GetAll
 * @apiGroup AssetTestPiece
 *
 * @apiDescription 测试件列表查询
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "list": [
 *             {
 *                 "brand": "b1111112",
 *                 "code": "c111",
 *                 "create_time": 1628674339,
 *                 "create_user_id": 0,
 *                 "detail": "d111",
 *                 "id": 4,
 *                 "serial_number": "",
 *                 "update_time": 0
 *             }
 *         ],
 *         "pagination": {
 *             "current_page": 1,
 *             "per_page": 20,
 *             "total": 5,
 *             "total_pages": 1
 *         }
 *     },
 *     "msg": ""
 * }
 */
func (this AssetTestPieceController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	s := mysql.GetSession()

	// 查询组键
	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	widget.SetTransformer(&transformer.AssetTestPieceTransformer{})
	all := widget.PaginatorFind(s, &[]mysql_model.AssetTestPiece{})
	response.RenderSuccess(ctx, all)
}

func (this AssetTestPieceController) GetAllWithNotPage(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	s := mysql.GetSession()

	// 查询组键
	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	widget.SetTransformer(&transformer.AssetTestPieceTransformer{})
	all, _ := widget.All(s, &[]mysql_model.AssetTestPiece{})
	response.RenderSuccess(ctx, all)
}

/**
 * apiType http
 * @api {post} /api/v1/asset_test_pieces 创建测试件记录
 * @apiVersion 0.1.0
 * @apiName Create
 * @apiGroup AssetTestPiece
 *
 * @apiDescription 创建测试件记录
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 0
 *     },
 *     "msg": ""
 * }
 */
func (this AssetTestPieceController) Create(ctx *gin.Context) {
	// 表单
	form := &mysql_model.AssetTestPieceCreateForm{}
	form.Name = request.MustString(ctx, "name")
	form.Version = request.MustString(ctx, "version")
	form.AssetVehicleId = request.MustInt(ctx, "asset_vehicle_id")
	form.Detail = request.String(ctx, "detail")

	uid := session.GetUserId(http_ctx.GetHttpCtx(ctx))
	id, err := mysql_model.AssetTestPieceCreateFromForm(form, int(uid))

	if err == nil {
		response.RenderSuccess(ctx, gin.H{"id": id})
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/asset_test_pieces/:id 测试件记录详情
 * @apiVersion 0.1.0
 * @apiName GetOne
 * @apiGroup AssetTestPiece
 *
 * @apiDescription 测试件记录详情
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "files": [
 *             {
 *                 "id": 1,
 *                 "version_id": 8,
 *                 "file_name": "检测平台.postman_collection.json",
 *                 "file_size": 35351,
 *                 "storage_type": 1,
 *                 "file_uuid": "6115163cd4f5de9efa985b2d",
 *                 "create_time": 1628771900,
 *                 "is_delete": 1,
 *                 "delete_user_id": 0,
 *                 "delete_time": 0
 *             }
 *         ],
 *         "model": {
 *             "id": 8,
 *             "asset_test_piece_id": 8,
 *             "version": "c111",
 *             "storage_type": 1,
 *             "create_user_id": 0,
 *             "update_time": 1628763862,
 *             "firmware_file_uuid": "",
 *             "firmware_name": "",
 *             "firmware_size": 0,
 *             "firmware_device_type": 0,
 *             "is_delete": 2,
 *             "create_time": 1628763862
 *         }
 *     },
 *     "msg": ""
 * }
 */
func (this AssetTestPieceController) GetOne(ctx *gin.Context) {
	id := request.ParamInt(ctx, "id")
	s := mysql.GetSession()
	s.Where("id=?", id)

	widget := orm.PWidget{}
	widget.SetTransformer(&transformer.AssetTestPieceDetailTransformer{})
	result, err := widget.One(s, &mysql_model.AssetTestPiece{})

	if err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {put} /api/v1/asset_test_pieces/:id 测试件记录更新
 * @apiVersion 0.1.0
 * @apiName Update
 * @apiGroup AssetTestPiece
 *
 * @apiDescription 测试件记录更新
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 4,
 *         "serial_number": "",
 *         "brand": "b1111112",
 *         "code": "c111",
 *         "detail": "d111",
 *         "create_user_id": 0,
 *         "update_time": 0,
 *         "create_time": 1628674339
 *     },
 *     "msg": ""
 * }
 */
func (this AssetTestPieceController) Update(ctx *gin.Context) {
	data := request.GetRequestBody(ctx)
	id := request.ParamInt(ctx, "id")

	if model, err := mysql_model.AssetTestPieceUpdateById(id, *data); err == nil {
		response.RenderSuccess(ctx, model)
		return
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/asset_test_pieces 测试件批量删除
 * @apiVersion 0.1.0
 * @apiName BulkDelete
 * @apiGroup AssetVehicle
 *
 * @apiDescription 测试件批量删除
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "number": 1
 *     },
 *     "msg": ""
 * }
 */
func (this AssetTestPieceController) BulkDelete(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)

	successNum := 0
	if _, has := req.TrySlice("ids"); has {
		ids := req.SliceInt("ids")

		for _, id := range ids {
			// todo 检查是否可以删除

			err := mysql_model.AssetTestPieceDeleteById(id)
			if err != nil {
				log.GetHttpLogLogger().Error(fmt.Sprintf("%v", err))
				break
			} else {
				successNum++
			}
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"number": successNum})
}

/**
 * apiType http
 * @api {get} /api/v1/asset_test_piece/get_by_version_id/:id 测试组件获取指定版本的详细信息
 * @apiVersion 0.1.0
 * @apiName GetByVersionId
 * @apiGroup AssetTestPiece
 *
 * @apiDescription 测试组件获取指定版本的详细信息
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 0
 *     },
 *     "msg": ""
 * }
 */
func (this AssetTestPieceController) GetByVersionId(ctx *gin.Context) {
	versionId := request.ParamInt(ctx, "id")

	result := gin.H{}

	// 查询测试组件特定版本记录
	versionModel := mysql_model.AssetTestPieceVersion{}
	if has, err := mysql.FindById(versionId, &versionModel); err != nil {
		panic(err)
	} else {
		if !has {
			response.RenderFailure(ctx, err)
			return
		} else {
			result["model"] = versionModel
			fileModels := mysql_model.AssetTestPieceVersionFileFindDetail(versionId)
			result["files"] = fileModels
		}
	}

	response.RenderSuccess(ctx, result)
	return
}

//func (this AssetTestPieceController) UpdateByVersionId(ctx *gin.Context) {
//	versionId := request.ParamInt(ctx, "id")
//	version:= request.ParamString(ctx,"version")
//	println("fVersion",version)
//	result := gin.H{}
//
//	// 查询测试组件特定版本记录
//	versionModel := mysql_model.AssetTestPieceVersion{}
//	if has, err := mysql.FindById(versionId, &versionModel); err != nil {
//		panic(err)
//	} else {
//		if !has {
//			response.RenderFailure(ctx, err)
//			return
//		} else {
//			result["model"] = versionModel
//			fileModels := mysql_model.AssetTestPieceVersionFileFindDetail(versionId)
//			result["files"] = fileModels
//		}
//	}
//
//	response.RenderSuccess(ctx, result)
//	return
//}

/**
 * apiType http
 * @api {get} /api/v1/asset_test_piece/upload_file/ 测试组件指定版本上传固件
 * @apiVersion 0.1.0
 * @apiName UpdateFirmware
 * @apiGroup AssetTestPiece
 *
 * @apiDescription 测试组件指定版本上传固件
 *
 * @apiUse authHeader
 *
 * @apiParam {int}      version_id	 测试件版本id
 * @apiParam {binary}   file		文件内容
 *
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 0
 *     },
 *     "msg": ""
 * }
 */
func (this AssetTestPieceController) UpdateFirmware(ctx *gin.Context) {
	// version_id 对应
	versionIdStr := ctx.Request.FormValue("version_id")
	firmware_version := ctx.Request.FormValue("firmware_version")
	versionId, err := strconv.Atoi(versionIdStr)
	if err != nil {
		panic(err)
	}

	// 固件名称
	firmwareName := ctx.Request.FormValue("firmware_name")
	// 固件设备类型 todo
	// 1 远程通信单元(ECU)
	// 2 信息娱乐单元(IVI)
	// 3 汽车网关(GW)
	firmwareDeviceType := ctx.Request.FormValue("firmware_device_type")
	firmwareDeviceTypeInt, _ := strconv.Atoi(firmwareDeviceType)

	// file字段内容是binary，即文件的具体内容
	if file, header, _ := ctx.Request.FormFile("file"); file != nil {
		fileName := header.Filename // 文件名

		// 文件内容
		fileContent := make([]byte, header.Size)
		num, err := file.Read(fileContent)
		if err != nil {
			response.RenderFailure(ctx, errors.New("文件处理失败"))
			return
		}
		if num <= 0 {
			response.RenderFailure(ctx, errors.New("您上传的文件为空文件！"))
			return
		}

		maxSize := 2147483648 // byte,最大2G
		if num > maxSize {
			response.RenderFailure(ctx, errors.New("您上传的文件太大了！"))
			return
		}

		// 使用mongodb存储文件
		fileUuid, err2 := orm_mongo.GridFSUpload(common.MC_PROJECT, fileName, fileContent)
		if err2 != nil {
			response.RenderFailure(ctx, errors.New("文件上传失败"))
			return
		}
		// 记录到数据库表中
		if err := mysql_model.AssetTestPieceVersionUploadFirmware(versionId, firmwareName, fileName, header.Size, fileUuid, firmwareDeviceTypeInt, firmware_version); err != nil {
			response.RenderFailure(ctx, err)
			return
		} else {
			response.RenderSuccess(ctx, gin.H{})
			return
		}
	} else {
		if err := mysql_model.AssetTestPieceVersionUploadInfo(versionId, firmwareName, firmware_version, firmwareDeviceTypeInt); err != nil {
			response.RenderFailure(ctx, err)
			return
		} else {
			response.RenderSuccess(ctx, gin.H{})
			return
		}
	}

}

/**
 * apiType http
 * @api {get} /api/v1/asset_test_piece/upload_file/ 测试组件指定版本上传文件
 * @apiVersion 0.1.0
 * @apiName UpdateFile
 * @apiGroup AssetTestPiece
 *
 * @apiDescription 测试组件指定版本上传文件
 *
 * @apiUse authHeader
 *
 * @apiParam {int}      version_id	 测试件版本id
 * @apiParam {binary}   file		文件内容
 *
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 0
 *     },
 *     "msg": ""
 * }
 */
func (this AssetTestPieceController) UpdateFile(ctx *gin.Context) {
	// version_id 对应
	versionIdStr := ctx.Request.FormValue("version_id")
	versionId, err := strconv.Atoi(versionIdStr)
	if err != nil {
		panic(err)
	}

	// file字段内容是binary，即文件的具体内容
	file, header, _ := ctx.Request.FormFile("file")

	fileName := header.Filename // 文件名

	// 文件内容
	fileContent := make([]byte, header.Size)
	num, err := file.Read(fileContent)
	if err != nil {
		response.RenderFailure(ctx, errors.New("文件处理失败"))
		return
	}
	if num <= 0 {
		response.RenderFailure(ctx, errors.New("您上传的文件为空文件！"))
		return
	}
	maxSize := 2147483648 // byte,最大2G
	if num > maxSize {
		response.RenderFailure(ctx, errors.New("您上传的文件太大了！"))
		return
	}

	// 使用mongodb存储文件
	fileUuid, err2 := orm_mongo.GridFSUpload(common.MC_PROJECT, fileName, fileContent)
	if err2 != nil {
		response.RenderFailure(ctx, errors.New("文件上传失败"))
		return
	}

	// 记录到数据库表中
	if err := mysql_model.AssetTestPieceVersionFileCreate(versionId, fileName, header.Size, fileUuid); err != nil {
		response.RenderFailure(ctx, err)
		return
	} else {
		response.RenderSuccess(ctx, gin.H{})
		return
	}
}

// 添加测试件版本
func (this AssetTestPieceController) CreateVersion(ctx *gin.Context) {
	version := request.MustString(ctx, "version")
	id := request.ParamInt(ctx, "id")
	uid := session.GetUserId(http_ctx.GetHttpCtx(ctx))

	testPieceVersion := new(mysql_model.AssetTestPieceVersion)
	testPieceVersion.AssetTestPieceId = id
	testPieceVersion.Version = version
	testPieceVersion.StorageType = mysql_model.StorageTypeMongo
	testPieceVersion.CreateUserId = int(uid)
	now := int(time.Now().Unix())
	testPieceVersion.CreateTime = now
	testPieceVersion.UpdateTime = now
	_, err := testPieceVersion.Create()
	if err == nil {
		response.RenderSuccess(ctx, gin.H{"id": id})
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this AssetTestPieceController) UpdateVersion(ctx *gin.Context) {
	data := request.GetRequestBody(ctx)
	id := request.Int(ctx, "version_id")
	if model, err := new(mysql_model.AssetTestPieceVersion).UpdateById(id, *data); err == nil {
		response.RenderSuccess(ctx, model)
		return
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}

// 添加测试件版本
func (this AssetTestPieceController) BulkDeleteVersion(ctx *gin.Context) {
	// assetTestPieceId := request.MustString(ctx, "version")
	req := request.GetRequestBody(ctx)

	successNum := 0
	if _, has := req.TrySlice("version_id"); has {
		ids := req.SliceInt("version_id")

		for _, id := range ids {
			// todo 检查是否可以删除

			_, err := new(mysql_model.AssetTestPieceVersion).Delete(id)
			if err != nil {
				log.GetHttpLogLogger().Error(fmt.Sprintf("%v", err))
				break
			} else {
				successNum++
			}
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"number": successNum})
}

func (this AssetTestPieceController) DeleteVersionFile(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	successNum := 0
	if _, has := req.TrySlice("version_file_id"); has {
		ids := req.SliceInt("version_file_id")
		for _, id := range ids {
			_, err := new(mysql_model.AssetTestPieceVersionFile).Delete(id)
			if err != nil {
				log.GetHttpLogLogger().Error(fmt.Sprintf("%v", err))
				break
			} else {
				successNum++
			}
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"number": successNum})
}

func (this AssetTestPieceController) DownloadFirmware(ctx *gin.Context) {
	fileId := ctx.Query("file_id")
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
