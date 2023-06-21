package controller

import (
	"encoding/base64"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/service"

	"github.com/gin-gonic/gin"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mongo_model"
	"skygo_detection/mysql_model"
)

type ToolController struct{}

// @auto_generated_api_begin
/**
 * apiType http
 * @api {get} /api/v1/tool/list 工具查询列表
 * @apiVersion 0.1.0
 * @apiName List
 * @apiGroup Tool
 *
 * @apiDescription 工具管理-工具查询列表
 *
 * @apiParam {string}      			search      	工具搜索内容
 * @apiParam {string}      			category_id   	工具分类ID
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "list": [
 *             {
 *                 "_id": "5fce00df89e90f9b41b639d4",
 *                 "brand": "brand",
 *                 "category_id": 101,
 *                 "category_name": "主机安全检测",
 *                 "create_time": 1607336159,
 *                 "create_user_id": 0,
 *                 "create_user_name": "",
 *                 "hardware_version": "hardware_version",
 *                 "link_pic": [
 *                     {
 *                         "name": "linkpic1.png",
 *                         "value": "000001"
 *                     },
 *                     {
 *                         "name": "linkpic2,jpg",
 *                         "value": "0000002"
 *                     }
 *                 ],
 *                 "params_json": "cGFyYW1zX2pzb25wYXJhbXNfanNvbg==",
 *                 "params_remarks": "remarksremarksremarksremarksremarks",
 *                 "script": [
 *                     {
 *                         "name": "script.sh",
 *                         "value": "000001"
 *                     },
 *                     {
 *                         "name": "script02.sh",
 *                         "value": "0000002"
 *                     }
 *                 ],
 *                 "search": [
 *                     "zhuph20208888888",
 *                     "主机安全检测",
 *                     "101",
 *                     "10001"
 *                 ],
 *                 "software_version": "software_version",
 *                 "sp_version": "sp_version",
 *                 "status": 1,
 *                 "system_version": "system_version",
 *                 "test_pic": [
 *                     {
 *                         "name": "name1.png",
 *                         "value": "000001"
 *                     },
 *                     {
 *                         "name": "name2,jpg",
 *                         "value": "0000002"
 *                     }
 *                 ],
 *                 "tool_detail": "tool_detailtool_detailtool_detailtool_detail",
 *                 "tool_name": "zhuph20208888888",
 *                 "tool_number": 20201207181559,
 *                 "update_time": 1607336159,
 *                 "use_detail": "use_detailuse_detailuse_detailuse_detail",
 *                 "use_manual": [
 *                     {
 *                         "name": "um1.png",
 *                         "value": "000001"
 *                     },
 *                     {
 *                         "name": "um2",
 *                         "value": "0000002"
 *                     }
 *                 ],
 *                 "use_manual_link": [
 *                     "htttp:wwww.baicu.com.gpd",
 *                     "http://wwww.biadu.com/asdlfjsdf_pnd"
 *                 ]
 *             }
 *         ],
 *         "pagination": {
 *             "count": 1,
 *             "current_page": 1,
 *             "per_page": 20,
 *             "total": 1,
 *             "total_pages": 1
 *
 *         }
 *     }
 * }
 */
func (this ToolController) List(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	search := req.String("search")
	categoryID := req.Int("category_id")
	queryParams := qmap.QM{
		"e_status": 1,
	}
	if search != "" {
		queryParams["l_search"] = search
	}
	if categoryID > 0 {
		queryParams["e_category_id"] = categoryID
	}
	mgoSession := sys_service.NewMgoSessionWithCond(common.MC_TOOL, queryParams).AddUrlQueryCondition(ctx.Request.URL.RawQuery)
	if res, err := mgoSession.GetPage(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/tool/category 获取工具类别信息
 * @apiVersion 0.1.0
 * @apiName Category
 * @apiGroup Tool
 *
 * @apiDescription 获取工具类别信息
 *
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":{
 *				"number":2
 *			}
 *      }
 */
func (this ToolController) Category(ctx *gin.Context) {
	response.RenderSuccess(ctx, map[string]interface{}{
		"101": "硬件安全",
		"102": "系统安全",
		"103": "应用安全",
		"104": "无线电安全",
		"105": "车载网络安全",
		"106": "代码安全",
		"107": "固件安全",
		"108": "云端安全",
		"109": "隐私安全",
		"110": "蜂窝网络安全",
	})
}

/**
 * apiType http
 * @api {get} /api/v1/tool/category_tool 获取工具类别联动信息
 * @apiVersion 0.1.0
 * @apiName Category
 * @apiGroup Tool
 *
 * @apiDescription 获取工具类别联动信息
 *
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":{
 *				"number":2
 *			}
 *      }
 */
func (this ToolController) CategoryTool(ctx *gin.Context) {
	if ret, err := new(mongo_model.ToolData).GetToolCate(); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, ret)
	}

}

/**
 * apiType http
 * @api {get} /api/v1/tool/detail 获取工具详情信息
 * @apiVersion 0.1.0
 * @apiName Detail
 * @apiGroup Tool
 *
 * @apiDescription 获取工具详情信息
 *
 * @apiParam {string}      			id      	工具ID
 *
 *
 * @apiSuccessExample {json} 请求成功示例:
 *
 *  {
 *     "code": 0,
 *     "data": {
 *         "brand": "brand",
 *         "category_id": 101,
 *         "category_name": "主机安全检测",
 *         "create_time": 1607336159,
 *         "create_user_id": 0,
 *         "create_user_name": "",
 *         "hardware_version": "hardware_version",
 *         "id": "5fce00df89e90f9b41b639d4",
 *         "link_pic": [
 *             {
 *                 "name": "linkpic1.png",
 *                 "value": "000001"
 *             },
 *             {
 *                 "name": "linkpic2,jpg",
 *                 "value": "0000002"
 *             }
 *         ],
 *         "params_json": "cGFyYW1zX2pzb25wYXJhbXNfanNvbg==",
 *         "params_remarks": "remarksremarksremarksremarksremarks",
 *         "script": [
 *             {
 *                 "name": "script.sh",
 *                 "value": "000001"
 *             },
 *             {
 *                 "name": "script02.sh",
 *                 "value": "0000002"
 *             }
 *         ],
 *         "search": [
 *             "zhuph20208888888",
 *             "主机安全检测",
 *             "101",
 *             "10001"
 *         ],
 *         "software_version": "software_version",
 *         "sp_version": "sp_version",
 *         "system_version": "system_version",
 *         "test_pic": [
 *             {
 *                 "name": "name1.png",
 *                 "value": "000001"
 *             },
 *             {
 *                 "name": "name2,jpg",
 *                 "value": "0000002"
 *             }
 *         ],
 *         "tool_detail": "tool_detailtool_detailtool_detailtool_detail",
 *         "tool_name": "zhuph20208888888",
 *         "tool_number": 20201207181559,
 *         "update_time": 1607336159,
 *         "use_detail": "use_detailuse_detailuse_detailuse_detail",
 *         "use_manual": [
 *             {
 *                 "name": "um1.png",
 *                 "value": "000001"
 *             },
 *             {
 *                 "name": "um2",
 *                 "value": "0000002"
 *             }
 *         ],
 *         "use_manual_link": [
 *             "htttp:wwww.baicu.com.gpd",
 *             "http://wwww.biadu.com/asdlfjsdf_pnd"
 *         ]
 *     }
 * }
 */
func (this ToolController) Detail(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if tempRts, err := new(mongo_model.ToolData).GetOne(req); err != nil {
		response.RenderFailure(ctx, err)
		return
	} else {
		ret := *tempRts
		ret["id"] = ret["_id"]
		paramsJsonStr, _ := base64.StdEncoding.DecodeString(ret["params_json"].(string))
		ret["params_json"] = string(paramsJsonStr)
		delete(ret, "status")
		delete(ret, "_id")
		response.RenderSuccess(ctx, ret)
	}
}

/*
*
  - apiType http
  - @api {post} /api/v1/tool/add 添加工具
  - @apiVersion 0.1.0
  - @apiName Add
  - @apiGroup Tool
    *
  - @apiDescription 添加工具
    *
  - @apiParam {string}      			name      	工具名称
  - @apiParam {string}      			test_pic   			测试工具图片 [{"name":"name1.png","value":"000001"},{"name":"name2,jpg","value":"0000002"}] name 为文件名称 value为file_id
  - @apiParam {string}      			category_name   			分类名称
  - @apiParam {int}      			category_id   				分类ID
  - @apiParam {string}      			tool_detail   				工具介绍
  - @apiParam {string}      			use_detail   				使用方法详情
  - @apiParam {string}      			use_manual   				使用手册 [{"name":"name1.png","value":"000001"},{"name":"name2,jpg","value":"0000002"}] name 为文件名称 value为file_id
  - @apiParam {string}      			link_pic   					工具连接示意图 [{"name":"name1.png","value":"000001"},{"name":"name2,jpg","value":"0000002"}] name 为文件名称 value为file_id
  - @apiParam {string}      			script   					脚本 [{"name":"name1.png","value":"000001"},{"name":"name2,jpg","value":"0000002"}] name 为文件名称 value为file_id
  - @apiParam {string}      			params_json   				工具参数配置
  - @apiParam {string}      			brand   					工具品牌
  - @apiParam {string}      			remarks   					工具参数配置备注
  - @apiParam {string}      			software_version   			软件版本
  - @apiParam {string}      			hardware_version   			硬件版本
  - @apiParam {string}      			sp_version   				设备型号
  - @apiParam {string}      			system_version   			系统版本
  - @apiParam {string}      			use_manual_link   			使用手册Link 多个URL使用逗号分隔

// * @apiParam {string}               history_version             变更历史迭代版本号 //To Do 添加工具时默认为"v1.0"

	*
	*
	* @apiParamExample {json}  请求参数示例:
	* {
	*    "name": "zhuph20201208015",
	*    "test_pic": "[{\"name\":\"test_pic01.png\",\"value\":\"000001\"},{\"name\":\"test_pic02,jpg\",\"value\":\"0000002\"}]",
	*    "category_name": "主机安全检测",
	*    "category_id": 101,
	*    "use_detail": "use_detailuse_detailuse_detailuse_detail",
	*    "use_manual": "[{\"name\":\"um01.pdf\",\"value\":\"000001\"},{\"name\":\"um02.pdf\",\"value\":\"0000002\"}]",
	*    "link_pic": "[{\"name\":\"lk01.png\",\"value\":\"000001\"},{\"name\":\"lk02,jpg\",\"value\":\"0000002\"}]",
	*    "script": "[{\"name\":\"script01.sh\",\"value\":\"000001\"},{\"name\":\"script02.c\",\"value\":\"0000002\"}]",
	*    "params_json": "",
	*    "brand": "brand",
	*    "remarks": "",
	*    "software_version": "software_version",
	*    "hardware_version": "hardware_version",
	*    "sp_version": "sp_version",
	*    "system_version": "system_version",
	*    "use_manual_link": "htttp:wwww.baicu.com.gpd,http://wwww.biadu.com/asdlfjsdf_pnd",
	*    "tool_detail": "tool_detailtool_detailtool_detailtool_detail"
	* }

	*
	* @apiSuccessExample {json} 请求成功示例:
	* {
	*     "code": 0,
	*     "data": {
	*         "Id": "5fc75bef89e90f61c041a2d3",
	*         "Name": "zhuph20201202",
	*         "ToolNumber": 20201202171839
	*     }
	* }
*/
func (this ToolController) Add(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	UserID := int(http_ctx.GetUserId(ctx))
	UserName := http_ctx.GetUserName(ctx)
	if ret, err := new(mongo_model.ToolData).Create(req, UserID, UserName); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, ret)
	}
}

/*
*

  - apiType http

  - @api {post} /api/v1/tool/edit 编辑工具

  - @apiVersion 0.1.0

  - @apiName Edit

  - @apiGroup Tool
    *

  - @apiDescription 编辑工具
    *

  - @apiParam {string}      			id      					工具id

  - @apiParam {string}      			name      					工具名称

  - @apiParam {string}      			test_pic   					测试工具图片 [{"name":"name1.png","value":"000001"},{"name":"name2,jpg","value":"0000002"}] name 为文件名称 value为file_id

  - @apiParam {string}      			category_name   			分类名称

  - @apiParam {int}      			category_id   				分类ID

  - @apiParam {string}      			tool_detail   				工具介绍

  - @apiParam {string}      			use_detail   				使用方法详情

  - @apiParam {string}      			use_manual   				使用手册 [{"name":"name1.png","value":"000001"},{"name":"name2,jpg","value":"0000002"}] name 为文件名称 value为file_id

  - @apiParam {string}      			link_pic   					工具连接示意图 [{"name":"name1.png","value":"000001"},{"name":"name2,jpg","value":"0000002"}] name 为文件名称 value为file_id

  - @apiParam {string}      			script   					脚本 [{"name":"name1.png","value":"000001"},{"name":"name2,jpg","value":"0000002"}] name 为文件名称 value为file_id

  - @apiParam {string}      			params_json   				工具参数配置

  - @apiParam {string}      			brand   					工具品牌

  - @apiParam {string}      			remarks   					工具参数配置备注

  - @apiParam {string}      			software_version   			软件版本

  - @apiParam {string}      			hardware_version   			硬件版本

  - @apiParam {string}      			sp_version   				设备型号

  - @apiParam {string}      			system_version   			系统版本

  - @apiParam {string}      			use_manual_link   			使用手册Link 多个URL使用逗号分隔

  - @apiParam {string}      			history_content   			版本变更log记录
    *

  - @apiParamExample {json}  请求参数示例:

  - {

  - "id":"adshfkjasdgkhsg",

  - "name": "zhuph20201208015",

  - "test_pic": "[{\"name\":\"test_pic01.png\",\"value\":\"000001\"},{\"name\":\"test_pic02,jpg\",\"value\":\"0000002\"}]",

  - "category_name": "主机安全检测",

  - "category_id": 101,

  - "use_detail": "use_detailuse_detailuse_detailuse_detail",

  - "use_manual": "[{\"name\":\"um01.pdf\",\"value\":\"000001\"},{\"name\":\"um02.pdf\",\"value\":\"0000002\"}]",

  - "link_pic": "[{\"name\":\"lk01.png\",\"value\":\"000001\"},{\"name\":\"lk02,jpg\",\"value\":\"0000002\"}]",

  - "script": "[{\"name\":\"script01.sh\",\"value\":\"000001\"},{\"name\":\"script02.c\",\"value\":\"0000002\"}]",

  - "params_json": "",

  - "brand": "brand",

  - "remarks": "",

  - "software_version": "software_version",

  - "hardware_version": "hardware_version",

  - "sp_version": "sp_version",

  - "system_version": "system_version",

  - "use_manual_link": "htttp:wwww.baicu.com.gpd,http://wwww.biadu.com/asdlfjsdf_pnd",

  - "tool_detail": "tool_detailtool_detailtool_detailtool_detail",

  - "history_content":"版本变更历史记录"

  - }

    *

  - @apiSuccessExample {json} 请求成功示例:

  - {

  - "code": 0,

  - "data": {

  - "Id": "5fc75bef89e90f61c041a2d3",

  - "Name": "zhuph20201202",

  - "ToolNumber": 20201202171839

  - }

  - }
*/
func (this ToolController) Edit(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	UserID := int(http_ctx.GetUserId(ctx))
	UserName := http_ctx.GetUserName(ctx)
	if ret, err := new(mongo_model.ToolData).Edit(req, UserID, UserName); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, ret)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/tool/del 删除工具
 * @apiVersion 0.1.0
 * @apiName Delete
 * @apiGroup Tool
 *
 * @apiDescription 删除工具
 *
 * @apiParam {string}      			id      					工具id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *    "id":"adshfkjasdgkhsg",
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "Id": "5fc75bef89e90f61c041a2d3",
 *         "Name": "zhuph20201202",
 *         "ToolNumber": 20201202171839
 *     }
 * }
 */
func (this ToolController) Delete(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	// todo 如果测试任务中有这个tool则删除失败
	if isBusy(req.String("id")) {
		response.RenderFailure(ctx, errors.New("当前测试工具正在被使用"))
	} else {
		if ret, err := new(mongo_model.ToolData).Del(req); err != nil {
			response.RenderFailure(ctx, err)
		} else {
			response.RenderSuccess(ctx, ret)
		}
	}
}

// @auto_generated_api_end
/**
 * apiType http
 * @api {post} /api/v1/tool/upload 上传工具
 * @apiVersion 0.1.0
 * @apiName Upload
 * @apiGroup Tool
 *
 * @apiDescription 上传工具
 *
 * @apiParam {form-data}      			file      					工具包文件
 * @apiParam {string}      				file_name      				工具名称
 * @apiParam {string}      				file_type      				工具类型 doc pdf
 *
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "file_id": "5fc75ce489e90f61c041a2d4"
 *     }
 * }
 */
func (this ToolController) Upload(ctx *gin.Context) {
	fileName := ctx.Request.FormValue("file_name")
	file, header, _ := ctx.Request.FormFile("file")
	if fileName == "" && header != nil {
		fileName = header.Filename
	}
	fileContent := make([]byte, 0)
	file.Read(fileContent)
	if fileId, err := mongo.GridFSUpload(common.MC_GSF_TOOL, fileName, fileContent); err == nil {
		res := &qmap.QM{"file_id": fileId}
		response.RenderSuccess(ctx, res)
	} else {
		panic(err)
	}
}

/*
*

  - apiType http

  - @api {post} /api/v1/tool/edit_Tag 编辑工具标签

  - @apiVersion 0.1.0

  - @apiName EditTag

  - @apiGroup Tool
    *

  - @apiDescription 编辑工具标签
    *

  - @apiParam {string}      			id      					工具id

  - @apiParam {string}      			name      					工具名称

  - @apiParam {string}               tag                         工具标签

  - @apiParamExample {json}  请求参数示例:

  - {

  - "id": "618cca93b1b91940accb8d74",

  - "name": "skill",

  - "tag": "手动"

  - }

    *

  - @apiSuccessExample {json} 请求成功示例:

  - {

  - "code": 0,

  - "data": {

  - "Id": "618cca93b1b91940accb8d74",

  - "Name": "skill",

  - "ToolNumber": 20211111154731

  - },

  - "msg": ""

  - }
*/
func (this ToolController) EditTag(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	UserID := int(http_ctx.GetUserId(ctx))
	UserName := http_ctx.GetUserName(ctx)
	if ret, err := new(mongo_model.ToolData).UpdateTag(req, UserID, UserName); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, ret)
	}
}

// 查询是否正在被任务使用的测试工具
func isBusy(tid string) bool {
	models := make([]mysql_model.Task, 0)
	mysql.GetSession().Table(mysql_model.Task{}).Find(&models)
	var result = make([]string, 0)
	for _, model := range models {
		result = append(result, model.ToolId)
	}
	for _, v := range result {
		if v == tid {
			return true
		}
	}
	return false
}

/**
 * apiType http
 * @api {get} /api/v1/tool/download_app 下载检测工具包
 * @apiVersion 0.1.0
 * @apiName DownloadApp
 * @apiGroup Tool
 *
 * @apiDescription 下载检测工具包
 *
 * @apiParam {string}      			tool_type      					工具类型(hg_scanner、vul_scanner)
 *
 */
func (this ToolController) DownloadApp(ctx *gin.Context) {
	toolType := request.QueryString(ctx, "tool_type")
	fileName := ""
	relatePath := ""
	switch toolType {
	case common.TOOL_VUL_SCANNER:
		fileName = "vul_scan.apk"
		relatePath = "vul_scanner/app/"
	case common.TOOL_HG_ANDROID_SCANNER:
		fileName = "hg_scan.apk"
		relatePath = "hg_scanner/app/"
	default:
		panic(errors.New("位置工具类型:" + toolType))
	}
	packagePath := fmt.Sprintf("%s/%s%s", service.LoadConfig().Download.DownloadPath, relatePath, fileName)
	if content, err := ioutil.ReadFile(packagePath); err == nil {
		ctx.Writer.WriteHeader(http.StatusOK)
		ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
		ctx.Header("Content-Type", "*")
		ctx.Header("Accept-Length", fmt.Sprintf("%d", len(content)))
		ctx.Writer.Write(content)
	} else {
		response.RenderFailure(ctx, err)
	}
}
