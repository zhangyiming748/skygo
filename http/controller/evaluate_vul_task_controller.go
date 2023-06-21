package controller

import (
	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/mongo_model_tmp"
)

type EvaluateVulTaskController struct{}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_vul_tasks 查询所有漏洞检测任务
 * @apiVersion 0.1.0
 * @apiName GetAll
 * @apiGroup EvaluateVulTask
 *
 * @apiDescription 查询所有漏洞检测任务
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiDescription 分页查询测试项列表，
 * {
 *    "code": 0,
 *   "data": {
 *        "list": [
 *           {
 *                "_id": "5fdc27195d502660456242b3",
 *                "create_time": 1608263449101,
 *                "name": "漏洞检测",
 *                "status": 0,
 *                "task_id": "1234ABCD",
 *                "test_time": 1608263449101,
 *                "vul_scanner_id": "漏洞详情ID"
 *            }
 *        ],
 *        "pagination": {
 *            "count": 1,
 *            "current_page": 1,
 *            "per_page": 20,
 *            "total": 1,
 *            "total_pages": 1
 *        }
 *    }
 * }
 */
func (this EvaluateVulTaskController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	data, _ := new(mongo_model_tmp.EvaluateVulTask).GetAll(queryParams)
	response.RenderSuccess(ctx, data)
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_vul_tasks  添加漏洞检测任务
 * @apiVersion 0.1.0
 * @apiName Create
 * @apiGroup EvaluateVulTask
 *
 * @apiDescription 添加漏洞检测任务
 *
 * @apiUse authHeader
 *
 * @apiParam {string}  name    漏洞检测任务名称
 *
 * @apiParamExample {json}  请求参数示例:
 *  {"name":"漏洞检测任务名称"}
 *
 * @apiSuccessExample {json} 请求成功示例:
 *{
 *    "code": 0,
 *    "data": {
 *        "create_time": 1608263449101,
 *        "id": "5fdc27195d502660456242b3",
 *        "name": "漏洞检测",
 *        "status": 0,
 *        "task_id": "1234ABCD",
 *        "test_time": 1608263449101,
 *        "vul_scanner_id": ""
 *    }
 *}
 */
func (this EvaluateVulTaskController) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if result, err := new(mongo_model_tmp.EvaluateVulTask).Create("", *req); err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/evaluate_vul_tasks 批量删除漏洞检测任务
 * @apiVersion 0.1.0
 * @apiName BulkDelete
 * @apiGroup EvaluateVulTask
 *
 * @apiDescription 批量删除漏洞检测任务
 *
 * @apiUse authHeader
 *
 * @apiParam {[]string}   ids  测试项id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  "ids":[
 *		"5e688f7a24b6476b74bb3548"
 * 	]
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":{
 *				"number":1
 *			}
 *      }
 */
func (this EvaluateVulTaskController) BulkDelete(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	rawIds := []string{}
	if ids, has := req.TrySlice("ids"); has {
		for _, id := range ids {
			rawIds = append(rawIds, id.(string))
		}
	}
	data, err := new(mongo_model_tmp.EvaluateVulTask).BulkDelete(rawIds)
	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, data)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_vul_task/download 下载漏洞工具
 * @apiVersion 0.1.0
 * @apiName DownloadTool
 * @apiGroup EvaluateVulTask
 *
 * @apiDescription 下载漏洞工具
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 *{
 *    "code": 0,
 *    "data": {
 *        "url": "http://pub-zzdt.s3.360.cn/sk3b8c6"
 *    }
 *}
 */
func (this EvaluateVulTaskController) DownloadTool(ctx *gin.Context) {
	// path := "http://10.4.7.250/vulscan.zip"
	res := &qmap.QM{
		"code": 0,
		"data": qmap.QM{"url": "path"},
	}
	response.RenderSuccess(ctx, res)

	// s3 := service.NewS3Client()
	// if path, err := s3.GetSignedUrl("vulscan.zip"); err != nil {
	//	response.RenderFailure(ctx, err)
	// } else {
	//	res := &qmap.QM{
	//		"code": 0,
	//		"data": qmap.QM{"url":path},
	//	}
	//	response.RenderSuccess(ctx, res)
	// }
}
