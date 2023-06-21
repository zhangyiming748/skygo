package controller

import (
	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/common_lib/session"
	"skygo_detection/mongo_model"
)

type EvaluateTestCaseController struct{}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_test_cases 分页查询测试用例列表
 * @apiVersion 0.1.0
 * @apiName GetAll
 * @apiGroup EvaluateTestCase
 *
 * @apiDescription 分页查询测试用例列表
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "list": [
 *             {
 *                 "_id": "5fcf691d5d502646ee8c5725",
 *                 "content": "内容变更2",
 *                 "content_version": "v1.1",
 *                 "create_time": 1607428381189,
 *                 "diff_content": [
 *                     {
 *                         "_id": "5fcf69485d502646ee8c5727",
 *                         "content": "内容变更2",
 *                         "op_id": 23,
 *                         "test_case_id": "5fcf691d5d502646ee8c5725",
 *                         "timestamp": 1607428424410,
 *                         "user_name": "用户名",
 *                         "version": "v1.1"
 *                     },
 *                     {
 *                         "_id": "5fcf691d5d502646ee8c5726",
 *                         "content": "内容变更",
 *                         "op_id": 23,
 *                         "test_case_id": "5fcf691d5d502646ee8c5725",
 *                         "timestamp": 1607428381189,
 *                         "user_name": "用户名",
 *                         "version": "v1.0"
 *                     }
 *                 ],
 *                 "last_update_op_id": 23,
 *                 "module_name": "测试组件",
 *                 "module_type": "测试分类",
 *                 "op_id": 23,
 *                 "status": 1,
 *                 "test_auto_test_degree": "自动化程度2",
 *                 "test_case_level": "测试用例级别2",
 *                 "test_case_number": "Test1000000000001",
 *                 "test_env_diagram": "测试环境搭建示意图id号码abcd1234562",
 *                 "test_level": 1,
 *                 "test_name": "测试用例名称",
 *                 "test_objective": "测试目的2",
 *                 "test_procedure": "测试步骤2",
 *                 "test_scripts": [
 *                     {
 *                         "name": "脚本1",
 *                         "value": "脚本id"
 *                     },
 *                     {
 *                         "name": "脚本2",
 *                         "value": "脚本id"
 *                     }
 *                 ],
 *                 "test_standard": "测试标准2",
 *                 "test_tools": [
 *                     {
 *                         "name": "工具1",
 *                         "value": "工具id"
 *                     },
 *                     {
 *                         "name": "工具2",
 *                         "value": "工具id"
 *                     }
 *                 ],
 *                 "update_time": 1607428424410,
 *                 "user_name": "用户名"
 *             },
 *             ...
 *         ],
 *         "pagination": {
 *             "count": 3,
 *             "current_page": 1,
 *             "per_page": 20,
 *             "total": 3,
 *             "total_pages": 1
 *         }
 *     }
 * }
 */
func (this EvaluateTestCaseController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	data, err := new(mongo_model.EvaluateTestCase).GetAll(queryParams)
	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, data)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_test_cases/:id 查询某一测试用例
 * @apiVersion 0.1.0
 * @apiName GetOne
 * @apiGroup EvaluateTestCase
 *
 * @apiDescription 根据id查询某一测试用例
 *
 * @apiParam {string}   id  		测试用例id
 *
 * curl 10.16.133.118:3001/api/v1/evaluate_test_cases/5fcf691d5d502646ee8c5725
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "_id": "5fcf691d5d502646ee8c5725",
 *         "content": "内容变更2",
 *         "content_version": "v1.1",
 *         "create_time": 1607428381189,
 *         "diff_content": [
 *             {
 *                 "_id": "5fcf69485d502646ee8c5727",
 *                 "content": "内容变更2",
 *                 "op_id": 23,
 *                 "test_case_id": "5fcf691d5d502646ee8c5725",
 *                 "timestamp": 1607428424410,
 *                 "user_name": "用户名",
 *                 "version": "v1.1"
 *             },
 *             {
 *                 "_id": "5fcf691d5d502646ee8c5726",
 *                 "content": "内容变更",
 *                 "op_id": 23,
 *                 "test_case_id": "5fcf691d5d502646ee8c5725",
 *                 "timestamp": 1607428381189,
 *                 "user_name": "用户名",
 *                 "version": "v1.0"
 *             }
 *         ],
 *         "last_update_op_id": 23,
 *         "module_name": "测试组件",
 *         "module_type": "测试分类",
 *         "op_id": 23,
 *         "status": 1,
 *         "test_auto_test_degree": "自动化程度2",
 *         "test_case_level": "测试用例级别2",
 *         "test_case_number": "Test1000000000001",
 *         "test_env_diagram": "测试环境搭建示意图id号码abcd1234562",
 *         "test_level": 1,
 *         "test_name": "测试用例名称",
 *         "test_objective": "测试目的2",
 *         "test_procedure": "测试步骤2",
 *         "test_scripts": [
 *             {
 *                 "name": "脚本1",
 *                 "value": "脚本id"
 *             },
 *             {
 *                 "name": "脚本2",
 *                 "value": "脚本id"
 *             }
 *         ],
 *         "test_standard": "测试标准2",
 *         "test_tools": [
 *             {
 *                 "name": "工具1",
 *                 "value": "工具id"
 *             },
 *             {
 *                 "name": "工具2",
 *                 "value": "工具id"
 *             }
 *         ],
 *         "update_time": 1607428424410,
 *         "user_name": "用户名"
 *     }
 * }
 */
func (this EvaluateTestCaseController) GetOne(ctx *gin.Context) {
	id := ctx.Param("id")

	if testCase, err := new(mongo_model.EvaluateTestCase).GetOne(id); err == nil {
		response.RenderSuccess(ctx, testCase)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_test_cases 添加测试用例
 * @apiVersion 0.1.0
 * @apiName Create
 * @apiGroup EvaluateTestCase
 *
 * @apiDescription 添加测试用例
 *
 * @apiParam {string}      test_objective       测试目的
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *     "test_objective": "测试目的2",
 *     "test_procedure": "测试步骤2",
 *     "test_standard": "测试标准2",
 *     "test_level": 1,
 *     "test_case_level": "测试用例级别2",
 *     "test_auto_test_degree": "自动化程度2",
 *     "test_scripts": [
 *         {
 *             "name": "脚本1",
 *             "value": "脚本id"
 *         },
 *         {
 *             "name": "脚本2",
 *             "value": "脚本id"
 *         }
 *     ],
 *     "test_tools": [
 *         {
 *             "name": "工具1",
 *             "value": "工具id"
 *         },
 *         {
 *             "name": "工具2",
 *             "value": "工具id"
 *         }
 *     ],
 *     "test_env_diagram": "测试环境搭建示意图id号码abcd1234562",
 *     "content": "内容变更2"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *     {
 *         "code": 0,
 *         "data": {
 *             "_id": "5fcf691d5d502646ee8c5725",
 *             "content": "内容变更",
 *             "content_version": "v1.0",
 *             "create_time": 1607428381189,
 *             "last_update_op_id": 23,
 *             "module_name": "测试组件",
 *             "module_type": "测试分类",
 *             "op_id": 23,
 *             "status": 1,
 *             "test_auto_test_degree": "自动化程度",
 *             "test_case_level": "测试用例级别",
 *             "test_case_number": "Test1000000000001",
 *             "test_env_diagram": "测试环境搭建示意图id号码abcd123456",
 *             "test_level": 1,
 *             "test_name": "测试用例名称",
 *             "test_objective": "测试目的",
 *             "test_procedure": "测试步骤",
 *             "test_scripts": [
 *                 {
 *                     "name": "脚本1",
 *                     "value": "脚本id"
 *                 },
 *                 {
 *                     "name": "脚本2",
 *                     "value": "脚本id"
 *                 }
 *             ],
 *             "test_standard": "测试标准",
 *             "test_tools": [
 *                 {
 *                     "name": "工具1",
 *                     "value": "工具id"
 *                 },
 *                 {
 *                     "name": "工具2",
 *                     "value": "工具id"
 *                 }
 *             ],
 *             "update_time": 1607428381189,
 *             "user_name": "用户名"
 *         }
 *     }
 */
func (this EvaluateTestCaseController) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if testCase, err := new(mongo_model.EvaluateTestCase).Create(*req, int(session.GetUserId(http_ctx.GetHttpCtx(ctx)))); err == nil {
		response.RenderSuccess(ctx, testCase)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {put} /api/v1/evaluate_test_cases/:id 更新测试用例
 * @apiVersion 0.1.0
 * @apiName Update
 * @apiGroup EvaluateTestCase
 *
 * @apiDescription 根据测试用例ID,更新测试用例
 *
 * @apiParam {string}           id                      测试用例id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *     "test_objective": "测试目的2",
 *     "test_procedure": "测试步骤2",
 *     "test_standard": "测试标准2",
 *     "test_level": 1,
 *     "test_case_level": "测试用例级别2",
 *     "test_auto_test_degree": "自动化程度2",
 *     "test_scripts": [
 *         {
 *             "name": "脚本1",
 *             "value": "脚本id"
 *         },
 *         {
 *             "name": "脚本2",
 *             "value": "脚本id"
 *         }
 *     ],
 *     "test_tools": [
 *         {
 *             "name": "工具1",
 *             "value": "工具id"
 *         },
 *         {
 *             "name": "工具2",
 *             "value": "工具id"
 *         }
 *     ],
 *     "test_env_diagram": "测试环境搭建示意图id号码abcd1234562",
 *     "content": "内容变更2"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "_id": "5fcf691d5d502646ee8c5725",
 *         "content": "内容变更2",
 *         "content_version": "v1.1",
 *         "create_time": 1607428381189,
 *         "last_update_op_id": 23,
 *         "module_name": "测试组件",
 *         "module_type": "测试分类",
 *         "op_id": 23,
 *         "status": 1,
 *         "test_auto_test_degree": "自动化程度2",
 *         "test_case_level": "测试用例级别2",
 *         "test_case_number": "Test1000000000001",
 *         "test_env_diagram": "测试环境搭建示意图id号码abcd1234562",
 *         "test_level": 1,
 *         "test_name": "测试用例名称",
 *         "test_objective": "测试目的2",
 *         "test_procedure": "测试步骤2",
 *         "test_scripts": [
 *             {
 *                 "name": "脚本1",
 *                 "value": "脚本id"
 *             },
 *             {
 *                 "name": "脚本2",
 *                 "value": "脚本id"
 *             }
 *         ],
 *         "test_standard": "测试标准2",
 *         "test_tools": [
 *             {
 *                 "name": "工具1",
 *                 "value": "工具id"
 *             },
 *             {
 *                 "name": "工具2",
 *                 "value": "工具id"
 *             }
 *         ],
 *         "update_time": 1607428424410,
 *         "user_name": "用户名"
 *     }
 * }
 */
func (this EvaluateTestCaseController) Update(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	(*req)["id"] = ctx.Param("id")
	// (*req)["query_params"] = ctx.Request.URL.RawQuery

	id := req.MustString("id")
	if testCase, err := new(mongo_model.EvaluateTestCase).Update(id, *req, int(session.GetUserId(http_ctx.GetHttpCtx(ctx)))); err == nil {
		response.RenderSuccess(ctx, testCase)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/evaluate_test_cases 批量删除测试用例
 * @apiVersion 0.1.0
 * @apiName BulkDelete
 * @apiGroup EvaluateTestCase
 *
 * @apiDescription 批量删除测试用例
 *
 * @apiParam {[]string}   ids  测试用例id
 *
 * @apiParamExample {json}  请求参数示例:
 *     {
 *         "ids": [
 *             "5fca0de55d502636d84ffd9b"
 *         ]
 *     }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *      {
 *           "code": 0
 *			 "data":{
 *				"number":1
 *			}
 *      }
 */
func (this EvaluateTestCaseController) BulkDelete(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	// (*req)["id"] = ctx.Param("id")
	// (*req)["query_params"] = ctx.Request.URL.RawQuery

	rawIds := []string{}
	if ids, has := req.TrySlice("ids"); has {
		for _, id := range ids {
			rawIds = append(rawIds, id.(string))
		}
	}
	data, err := new(mongo_model.EvaluateTestCase).BulkDelete(rawIds)
	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, data)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_test_case/upload 批量导入测试用例
 * @apiVersion 1.0.0
 * @apiName Upload
 * @apiGroup EvaluateTestCase
 *
 * @apiDescription 测试用例导入
 *
 * @apiUse authHeader
 *
 * @apiParam {string} 	[file_name]       	文件名称
 * @apiParam {file}		file 				文件
 *
 * @apiExample {curl} 请求示例:
 * curl -i http://localhost/api/v1/evaluate_test_case/upload
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
func (this EvaluateTestCaseController) Upload(ctx *gin.Context) {
	//ctx.Redirect(http.StatusMovedPermanently,"http://www.baidu.com")
	new(KnowledgeTestCaseController).Upload(ctx)
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_test_case/download 测试用例导出
 * @apiVersion 1.0.0
 * @apiName Download
 * @apiGroup EvaluateTestCase
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
 * curl -i http://localhost/api/v1/evaluate_test_case/download
 */
func (this EvaluateTestCaseController) Download(ctx *gin.Context) {
	// params := qmap.QM{}
	// params = params.Merge(*request.GetRequestBody(ctx))
	// ids, _ :=json.Marshal(params)
	// req := project_manage.TestCaseDownloadRequest{
	// 	Ids: string(ids),
	// }
	// rpcClient := client.NewGRpcClient(common.PM_SERVICE, http_ctx.NewOutputContext(ctx))
	// defer rpcClient.Close()
	// if resp, err := project_manage.NewEvaluateTestCaseClient(rpcClient.Client).Download(rpcClient.Ctx, &req); err == nil {
	// 	fileContent := []byte{}
	// 	var fileName string
	// 	for {
	// 		if fileBlock, recvErr := resp.Recv(); recvErr == io.EOF {
	// 			break
	// 		} else if recvErr != nil {
	// 			panic(recvErr)
	// 		} else {
	// 			fileName = fileBlock.FileName
	// 			fileContent = append(fileContent, fileBlock.FileContent...)
	// 		}
	// 	}
	// 	ctx.Writer.WriteHeader(http.StatusOK)
	// 	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	// 	ctx.Header("Content-Type", "*")
	// 	ctx.Header("Accept-Length", fmt.Sprintf("%d", len(fileContent)))
	// 	ctx.Writer.Write(fileContent)
	//
	// } else {
	// 	panic(err)
	// }
}
