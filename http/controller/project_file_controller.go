package controller

import (
	"context"
	"errors"
	"fmt"
	"skygo_detection/logic/project_file_logic"

	"net/http"
	"path/filepath"

	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mongo_model"
	"skygo_detection/mysql_model"

	"github.com/gin-gonic/gin"
)

const TIMEOUT = 200 //毫秒
type ProjectFileController struct{}

/**
 * apiType http
 * @api {post} /api/v1/project_file/upload_project_file 项目文档上传
 * @apiVersion 1.0.0
 * @apiName CreateProjectFile
 * @apiGroup ProjectFile
 *
 * @apiDescription 项目文档上传
 *
 * @apiUse authHeader
 *
 * @apiParam {string}      			project_id      	所属项目id
 * @apiParam {string}      			[parent_id]   		所属文件夹id
 * @apiParam {string}      			file_name       	文件名称
 * @apiParam {string=dir,doc}      	file_type       	文件类型(dir:文件夹,doc:文件)
 * @apiParam {file}      			file       			文件
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/project_file/upload_project_file
 *
 * @apiSuccessExample {json} 请求成功示例:
 *  {
 *		"code":0,
 *		"msg":"",
 *		"data":{
 *			"file_id":"a834qafmcxvadfq1123"
 *		}
 *  }
 */
func (this ProjectFileController) CreateProjectFile(ctx *gin.Context) {
	fileName := ctx.Request.FormValue("file_name")
	fileType := ctx.Request.FormValue("file_type")
	file, header, _ := ctx.Request.FormFile("file")
	if fileName == "" && header != nil {
		fileName = header.Filename
	}

	var projectId, parentId string
	var metaFileSize, opId int

	projectId = ctx.Request.FormValue("project_id")
	parentId = ctx.Request.FormValue("parent_id")
	fileType = ctx.Request.FormValue("file_type")
	opId = int(http_ctx.GetUserId(ctx))

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

	// 文件存储在mongo中的文件id
	metaFileId := ""
	metaFileSize = num

	// 如果上传的文件是文档文件，则先把文件内容存储到mongo中
	if fileType == mongo_model.FILE_TYPE_DOC {
		if fileId, err := mongo.GridFSUpload(common.MC_PROJECT, fileName, fileContent); err != nil {
			panic(err)
		} else {
			metaFileId = fileId
		}
	}

	// 最后一步创建项目文件记录
	if projectFile, err := new(mongo_model.PMFile).Create(context.TODO(), projectId, metaFileId, fileName, fileType, parentId, metaFileSize, opId); err != nil {
		response.RenderFailure(ctx, err)
		return
	} else {
		response.RenderSuccess(ctx, qmap.QM{"id": projectFile.MetaFileId})
		return
	}
}

/**
 * apiType http
 * @api {post} /api/v1/project_file/upload 项目文件上传
 * @apiVersion 1.0.0
 * @apiName Upload
 * @apiGroup ProjectFile
 *
 * @apiDescription 项目文件上传
 *
 * @apiUse authHeader
 *
 * @apiParam {string} 	[file_name]       	文件名称
 * @apiParam {file}		file 				文件
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/project_file/upload
 *
 * @apiSuccessExample {json} 请求成功示例:
 *  {
 *		"code":0,
 *		"msg":"",
 *		"data":{
 *			"file_id":"a834qafmcxvadfq1123"
 *		}
 *  }
 */
func (this ProjectFileController) Upload(ctx *gin.Context) {
	fileName := ctx.Request.FormValue("file_name")
	if file, header, err := ctx.Request.FormFile("file"); err == nil {
		if fileName == "" {
			fileName = header.Filename
		}

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

		if fileId, err := mongo.GridFSUpload(common.MC_File, fileName, fileContent); err != nil {
			panic(err)
		} else {
			m := gin.H{
				"code": 0,
				"data": gin.H{"file_id": fileId},
			}
			ctx.AbortWithStatusJSON(200, m)
			return
		}
	} else {
		panic(err)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/project_file/download 项目管理文件下载
 * @apiVersion 1.0.0
 * @apiName Download
 * @apiGroup ProjectFile
 *
 * @apiDescription 项目管理文件下载
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {string}      file_id       文件id
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/project_file/download
 */
func (this ProjectFileController) Download(ctx *gin.Context) {
	fileId := ctx.Query("file_id")
	fi, fileContent, err := project_file_logic.FindFileByFileID(fileId)
	if err != nil {
		response.RenderFailure(ctx, err)
	}

	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fi.Name()))
	ctx.Header("Content-Type", "*")
	ctx.Header("Accept-Length", fmt.Sprintf("%d", len(fileContent)))
	ctx.Writer.Write(fileContent)

	go project_file_logic.DownloadTimes(fileId)
}

/**
 * apiType http
 * @api {get} /api/v1/project_file/image 项目管理图片查看
 * @apiVersion 1.0.0
 * @apiName ViewImage
 * @apiGroup ProjectFile
 *
 * @apiDescription 项目管理图片查看
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiParam {string}      file_id       文件id
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/project_file/image
 */
func (this ProjectFileController) ViewImage(ctx *gin.Context) {
	fileId := ctx.Query("file_id")
	if fi, err := mongo.GridFSOpenId(common.MC_File, bson.ObjectIdHex(fileId)); err == nil {
		defer fi.Close()

		fileContent := make([]byte, fi.Size())
		if _, readErr := fi.Read(fileContent); readErr != nil {
			panic(readErr)
		}

		fileName := fi.Name()

		contentType := ""
		switch filepath.Base(fileName) {
		case "gif":
			contentType = "image/gif"
		case "jpeg", "jpe":
			contentType = "image/jpeg"
		default:
			contentType = "image/png"
		}
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", fileName))
		ctx.Header("Content-Type", contentType)
		ctx.Header("Accept-Length", fmt.Sprintf("%d", len(fileContent)))
		ctx.Writer.Write(fileContent)
	} else {
		panic(err)
	}
}

// @auto_generated_api_begin
/**
 * apiType http
 * @api {get} /api/v1/project_file/all 查询项目文件列表
 * @apiVersion 0.1.0
 * @apiName GetProjectList
 * @apiGroup ProjectFile
 *
 * @apiDescription 查询项目文件列表
 *
 * @apiParam {string}      			project_id      	所属项目id
 * @apiParam {string}      			dir_id   			文件夹id
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "_id": "5e6b34aa24b647048dd597ad",
 *             "create_time": 1584084138268,
 *             "file_name": "a",
 *             "file_type": "doc",
 *             "meta_file_id": "5e6b34aa24b647048dd597ab",
 *             "parent_id": "5e6b346824b647048dd597aa",
 *             "project_id": "5e66189224b6474fd32bdefc"
 *         }
 *     ]
 * }
 */
func (this ProjectFileController) GetProjectList(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	// (*req)["id"] = ctx.Param("id")
	(*req)["query_params"] = ctx.Request.URL.RawQuery

	params := qmap.QM{
		"e_project_id": req.MustString("project_id"),
		"e_parent_id":  req.String("dir_id"),
	}
	if res, err := mongo.NewMgoSessionWithCond(common.MC_PROJECT_File, params).SetTransformFunc(this.fileListTransform).SetLimit(10000).Get(); err == nil {
		response.RenderSuccess(ctx, res)
		return
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}

func (this ProjectFileController) fileListTransform(data qmap.QM) qmap.QM {
	// 如果meta_file_size为0  代表，是文件件，需要去查找文件夹下边所有文件的meta_file_size值
	if metaFileSize := data.Int("meta_file_size"); metaFileSize == 0 {
		params := bson.M{
			"$match": bson.M{
				"project_id": data.MustString("project_id"),
				"parent_id":  data["_id"].(bson.ObjectId).Hex(),
			},
		}
		agg := bson.M{
			"$group": bson.M{
				"_id": nil,
				"total": bson.M{
					"$sum": "$meta_file_size",
				},
			},
		}
		pipeline := []bson.M{params, agg}
		result := []map[string]interface{}{}
		if err := mongo.NewMgoSession(common.MC_PROJECT_File).Session.Pipe(pipeline).All(&result); err != nil {
			panic(err)
		}
		if len(result) != 0 {
			data["meta_file_size"] = result[0]["total"]
		}
	}

	// 查询操作人员信息
	if rsp, err := new(mysql_model.SysUser).GetUserInfo(data.Int("op_id")); err == nil {
		userTempt := rsp.Map("data")
		if realname := userTempt.String("realname"); realname != "" {
			data["op_name"] = realname
		} else {
			data["op_name"] = userTempt.String("username")
		}
	} else {
		data["op_name"] = ""
	}
	return data
}

/**
 * apiType http
 * @api {delete} /api/v1/project_files 批量删除项目文件
 * @apiVersion 0.1.0
 * @apiName BulkDeleteProjectFile
 * @apiGroup ProjectFile
 *
 * @apiDescription 批量删除项目文件
 *
 * @apiParam {string}      			project_id      	所属项目id
 * @apiParam {[]string}   			ids  				项目文件id
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":{
 *				"number":2
 *			}
 *      }
 */
func (this ProjectFileController) BulkDeleteProjectFile(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	// (*req)["id"] = ctx.Param("id")
	(*req)["query_params"] = ctx.Request.URL.RawQuery

	effectNum := 0
	if rawIds, has := req.TrySlice("ids"); has {
		ids := []bson.ObjectId{}
		for _, id := range rawIds {
			ids = append(ids, bson.ObjectIdHex(id.(string)))
		}
		if len(ids) > 0 {
			match := bson.M{
				"project_id": bson.M{"$eq": req.MustString("project_id")},
				"_id":        bson.M{"$in": ids},
			}
			if changeInfo, err := mongo.NewMgoSession(common.MC_PROJECT_File).RemoveAll(match); err == nil {
				effectNum = changeInfo.Removed
			} else {
				response.RenderFailure(ctx, err)
				return
			}
		}
	}

	response.RenderSuccess(ctx, qmap.QM{"number": effectNum})
	return
}

/**
 * apiType http
 * @api {post} /api/v1/project_file/rename 重命名项目文件
 * @apiVersion 0.1.0
 * @apiName RenameProjectFile
 * @apiGroup ProjectFile
 *
 * @apiDescription 重命名项目文件
 *
 * @apiParam {string}      			project_id      	所属项目id
 * @apiParam {string}      			file_id   			文件id
 * @apiParam {string}      			new_name   			新文件名称
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 * 		"project_id": "5e66189224b6474fd32bdefc",
 * 		"file_id": "5e6b50b624b6470a32fd38d8",
 * 		"new_name": "新文件名",
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "_id": "5e6b34aa24b647048dd597ad",
 *             "create_time": 1584084138268,
 *             "file_name": "a",
 *             "file_type": "doc",
 *             "meta_file_id": "5e6b34aa24b647048dd597ab",
 *             "parent_id": "5e6b346824b647048dd597aa",
 *             "project_id": "5e66189224b6474fd32bdefc"
 *         }
 *     ]
 * }
 */
func (this ProjectFileController) RenameProjectFile(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	// (*req)["id"] = ctx.Param("id")
	(*req)["query_params"] = ctx.Request.URL.RawQuery

	if pmFile, err := new(mongo_model.PMFile).Rename(req.MustString("project_id"), req.MustString("file_id"), req.MustString("new_name")); err == nil {
		response.RenderSuccess(ctx, custom_util.StructToMap2(*pmFile))
		return
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}
