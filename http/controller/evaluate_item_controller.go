package controller

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/tealeg/xlsx/v3"

	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/common_lib/session"
	"skygo_detection/mongo_model"
	"skygo_detection/mysql_model"
	"skygo_detection/service"
)

type EvaluateItemController struct{}

//@auto_generated_api_begin
/**
 * apiType http
 * @api {get} /api/v1/evaluate_items 分页查询测试用例列表
 * @apiVersion 0.1.0
 * @apiName GetAll
 * @apiGroup EvaluateItem
 *
 * @apiDescription 分页查询测试用例列表
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": {
 *                 "list": [
 *                         {
 *                                 "_id": "1",
 *                                 "asset_id": "US4VX1KR",
 *                                 "auto_test_level": "人工",
 *                                 "create_time": 1608365498089,
 *                                 "evaluate_task_id": "",
 *                                 "external_input": "外部输入",
 *                                 "level": "高",
 *                                 "module_name": "无线电",
 *                                 "module_type": "NFC钥匙",
 *                                 "module_type_id": "5fd87545aee3d1849a56ef8f",
 *                                 "name": "测试用例名称",
 *                                 "objective": "测试目的",
 *                                 "op_id": 0,
 *                                 "project_id": "5fd7218624b64712a27f47e8",
 *                                 "record_id": "",
 *                                 "tag": [],
 *                                 "test_case_level": "基础测试",
 *                                 "test_method": "黑盒",
 *                                 "test_procedure": "测试过程",
 *                                 "test_script": [
 *         		                        "name":"脚本",
 *         		                        "value":"dadfa123"
 *                                  ],
 *                                  "test_sketch_map": [
 *         		                        "name":"图片",
 *         		                        "value":"dadfa124"
 *                                  ],
 *                                 "test_standard": "测试标准",
 *                                 "test_time": 0,
 *                                 "update_time": 1608365498089,
 *                                 "vul_number": 0
 *                         }
 *                 ],
 *                 "pagination": {
 *                         "count": 1,
 *                         "current_page": 1,
 *                         "per_page": 20,
 *                         "total": 1,
 *                         "total_pages": 1
 *                 }
 *         }
 * }
 */
func (this EvaluateItemController) GetAll(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	// queryParams := ctx.Request.URL.RawQuery

	mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_ITEM)

	// 针对Post的参数Http_Query_Build化
	if postParams, has := req.TryInterface("post_params"); has {
		q := url.Values{}
		for key, val := range postParams.(map[string]interface{}) {
			if key == "ids" {
				continue
			}
			str := ""
			switch val.(type) {
			case int:
				str = strconv.Itoa(val.(int))
			case string:
				str = val.(string)
			case interface{}:
				data, _ := json.Marshal(val)
				str = string(data)
			}

			q.Set(key, str)
		}
		mgoSession.AddUrlQueryCondition(q.Encode())

		raw := qmap.New(postParams.(map[string]interface{}))
		if ids, has := raw.TryString("ids"); has {
			if idTemp, err := qmap.NewWithString(`{ "ids":` + ids + "}"); err == nil {
				params := qmap.QM{
					"in__id": idTemp.SliceString("ids"),
				}
				mgoSession.AddCondition(params)
			}
		}
	}

	mgoSession.SetTransformFunc(EvaluateItemTransformer)
	if res, err := mgoSession.GetPage(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func EvaluateItemTransformer(data qmap.QM) qmap.QM {
	if module, err := new(mongo_model.EvaluateModule).FindById(data.MustString("module_type_id")); err == nil {
		data["module_name"] = module.ModuleName
		data["module_type"] = module.ModuleType
	}
	data["evaluate_task_name"] = ""
	data["op_name"] = ""
	data["asset_name"] = ""

	testerId := 0
	testCount := 0
	taskId := ""
	// 先通过关系表查出task_id,
	if data.Int("status") == common.EIS_INUSE {
		taskId = data.MustString("pre_bind")
	} else if data.String("last_task_id") != "" { // 如果测试用例未被使用，则取上次绑定的任务ID
		taskId = data.MustString("last_task_id")
	}

	if taskId != "" {
		params := qmap.QM{
			"e__id": bson.ObjectIdHex(taskId),
		}
		mongoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params)
		if evaluateTask, err := mongoClient.GetOne(); err == nil {
			data["evaluate_task_name"] = evaluateTask.String("name")
			testerId = evaluateTask.Int("tester_id")
		}

		params = qmap.QM{
			"e_evaluate_task_id": taskId,
			"e_item_id":          data.MustString("_id"),
		}
		mongoClient = mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK_ITEM, params)
		if evaluateTaskItem, err := mongoClient.GetOne(); err == nil {
			data["test_status"] = evaluateTaskItem.Int("status")
			data["test_time"] = evaluateTaskItem.Int64("test_time")
			testCount = evaluateTaskItem.Int("test_count")
		}
	}

	if assetId, has := data.TryString("asset_id"); has {
		data["asset_name"] = new(mongo_model.EvaluateAsset).GetAssetName(assetId)
	}

	if opId := data.Int("op_id"); opId > 0 {
		data["op_name"] = mongo_model.GetNameById(opId)
	}

	if testerId > 0 {
		data["tester_name"] = mongo_model.GetNameById(testerId)
	}

	data["tester_id"] = testerId
	data["test_count"] = testCount
	return data
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_items/:id 查询某一测试用例
 * @apiVersion 0.1.0
 * @apiName GetOne
 * @apiGroup EvaluateItem
 *
 * @apiDescription 根据id查询某一测试用例
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   id  		测试用例id
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": {
 *                 "_id": "1",
 *                 "asset_id": "US4VX1KR",
 *                 "auto_test_level": "人工",
 *                 "create_time": 1608365498089,
 *                 "evaluate_task_id": "",
 *                 "external_input": "外部输入",
 *                 "level": "高",
 *                 "module_name": "无线电",
 *                 "module_type": "NFC钥匙",
 *                 "module_type_id": "5fd87545aee3d1849a56ef8f",
 *                 "name": "测试用例名称",
 *                 "objective": "测试目的",
 *                 "op_id": 0,
 *                 "project_id": "5fd7218624b64712a27f47e8",
 *                 "record_id": "",
 *                 "tag": [],
 *                 "test_case_level": "基础测试",
 *                 "test_method": "黑盒",
 *                 "test_procedure": "测试过程",
 *                 "test_script": [
 *         		        "name":"脚本",
 *         		         "value":"dadfa123"
 *                 ],
 *                 "test_sketch_map": [
 *         		         "name":"图片",
 *         		         "value":"dadfa124"
 *                 ],
 *                 "test_standard": "测试标准",
 *                 "test_time": 0,
 *                 "update_time": 1608365498089,
 *                 "vul_number": 0
 *         }
 * }
 */
func (this EvaluateItemController) GetOne(ctx *gin.Context) {
	params := qmap.QM{
		"e__id": ctx.Param("id"),
	}
	mgoSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params)
	mgoSession.SetTransformFunc(EvaluateItemTransformer)
	data, _ := mgoSession.GetOne()
	response.RenderSuccess(ctx, data)
}

/**
* apiType http
* @api {post} /api/v1/evaluate_items 添加测试用例
* @apiVersion 0.1.0
* @apiName Create
* @apiGroup EvaluateItem
*
* @apiDescription 添加测试用例
*
* @apiUse authHeader
*
* @apiParam {string}           suffix_id	    		测试用例后缀id
* @apiParam {string}           project_id    		    项目ID
* @apiParam {string}           name    				    测试用例名称
* @apiParam {string}           asset_id    		        测试资产ID
* @apiParam {string}           module_type_id    		测试组件id
* @apiParam {int}              level        			测试难度（高、中、低）
* @apiParam {string}           objective        	 	测试目的
* @apiParam {string}           test_procedure           测试步骤
* @apiParam {string}           test_standard          	测试标准
* @apiParam {string}           test_method          	测试方法
* @apiParam {string}           auto_test_level          自动化测试程度（自动化、人工）
* @apiParam {string}           test_case_level       	测试用例级别（基础测试、全面测试、提高测试、专家测试）
* @apiParam {string}           [external_input]         外部输入
* @apiParam {json}             [test_script]         	测试脚本
* @apiParam {json}             [test_sketch_map]        测试环境示意图
*
* @apiParamExample {json}  请求参数示例:
* {
*         "suffix_id":"1",
*         "project_id": "5fd7218624b64712a27f47e8",
*         "name": "测试用例名称",
*         "asset_id": "US4VX1KR",
*         "module_type_id":"5fd87545aee3d1849a56ef8f",
*         "level": "高",
*         "objective": "测试目的",
*         "external_input": "外部输入",
*         "test_procedure": "测试过程",
*         "test_standard": "测试标准",
*         "test_case_level":"基础测试",
*         "test_method":"黑盒",
*         "auto_test_level":"人工",
*         "test_script": [
*         		"name":"脚本",
*         		"value":"dadfa123"
*        ],
*        "test_sketch_map": [
*         		"name":"图片",
*         		"value":"dadfa124"
*        ],
* }
*
* @apiSuccessExample {json} 请求成功示例:
* {
*         "code": 0,
*         "data": {
*                 "asset_id": "US4VX1KR",
*                 "auto_test_level": "人工",
*                 "create_time": 1608365498089,
*                 "evaluate_task_id": "",
*                 "external_input": "外部输入",
*                 "id": "1",
*                 "level": "高",
*                 "module_type_id": "5fd87545aee3d1849a56ef8f",
*                 "name": "测试用例名称",
*                 "objective": "测试目的",
*                 "op_id": 0,
*                 "project_id": "5fd7218624b64712a27f47e8",
*                 "record_id": "",
*                 "tag": null,
*                 "test_case_level": "基础测试",
*                 "test_method": "黑盒",
*                 "test_procedure": "测试过程",
*                 "test_script": [
*         		        "name":"脚本",
*         		         "value":"dadfa123"
*                 ],
*                 "test_sketch_map": [
*         		         "name":"图片",
*         		         "value":"dadfa124"
*                 ],
*                 "test_standard": "测试标准",
*                 "test_time": 0,
*                 "update_time": 1608365498089,
*                 "vul_number": 0
*         }
* }
 */
func (this EvaluateItemController) Create(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if item, err := new(mongo_model.EvaluateItem).Create(context.TODO(), *req, session.GetUserId(http_ctx.GetHttpCtx(ctx))); err == nil {
		response.RenderSuccess(ctx, item)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {put} /api/v1/evaluate_items/:id 更新测试用例信息
 * @apiVersion 0.1.0
 * @apiName Update
 * @apiGroup EvaluateItem
 *
 * @apiDescription 根据测试用例ID,更新测试用例信息
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           id	    		    	测试用例id
 * @apiParam {string}           [name]    				测试用例名称
 * @apiParam {string}           [asset_id]    		    测试资产ID
 * @apiParam {string}           [module_type_id]    	测试组件id
 * @apiParam {int}           	[level]        			测试难度（高、中、低）
 * @apiParam {string}           [objective]        	 	测试目的
 * @apiParam {string}           [test_procedure]        测试步骤
 * @apiParam {string}           [test_standard]         测试标准
 * @apiParam {string}           [test_method]          	测试方法
 * @apiParam {string}           [auto_test_level]       自动化测试程度（自动化、人工）
 * @apiParam {string}           [test_case_level]       测试用例级别（基础测试、全面测试、提高测试、专家测试）
 * @apiParam {string}           [external_input]        外部输入
 * @apiParam {json}             [test_script]         	测试脚本
 * @apiParam {json}             [test_sketch_map]       测试环境示意图
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *         "id":"1",
 *         "name": "测试用例名称",
 *         "asset_id": "US4VX1KR",
 *         "module_type_id":"5fd87545aee3d1849a56ef8f",
 *         "level": "高",
 *         "objective": "测试目的",
 *         "external_input": "外部输入",
 *         "test_procedure": "测试过程",
 *         "test_standard": "测试标准",
 *         "test_case_level":"基础测试",
 *         "test_method":"黑盒",
 *         "auto_test_level":"人工",
 *          "test_script": [
 *              "name":"脚本",
 *         	    "value":"dadfa123"
 *          ],
 *          "test_sketch_map": [
 *         		"name":"图片",
 *         		"value":"dadfa124"
 *          ]
 * }
 */
func (this EvaluateItemController) Update(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	id := ctx.Param("id")
	opId := session.GetUserId(http_ctx.GetHttpCtx(ctx))
	(*req)["op_id"] = opId
	if item, err := new(mongo_model.EvaluateItem).Update(id, *req); err == nil {
		response.RenderSuccess(ctx, item)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/evaluate_items 批量删除测试用例
 * @apiVersion 0.1.0
 * @apiName BulkDelete
 * @apiGroup EvaluateItem
 *
 * @apiDescription 批量删除测试用例
 *
 * @apiUse authHeader
 *
 * @apiParam {[]string}   ids  测试用例id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  "ids":[
 *		"1",
 *		"2"
 * 	]
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": {
 *                 "failure_number": 0,
 *                 "success_number": 2
 *         }
 * }
 */
func (this EvaluateItemController) BulkDelete(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	// 删除测试用例需要:
	// 1.如果测试用例关联了项目任务或者已经被测试至少一次则不能被删除
	// 2.删除测试用例关联的测试记录
	// 3.删除测试用例关联的漏洞
	total := 0
	successNum := 0
	if ids := req.SliceString("ids"); len(ids) > 0 {
		total = len(ids)
		passIds := []string{}
		for _, id := range ids {
			param := qmap.QM{
				"e__id": id,
			}
			if item, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, param).GetOne(); err == nil {
				if item.String("pre_bind") == "" && item.String("last_task_id") == "" && item.Int("test_count") == 0 {
					passIds = append(passIds, id)
				} else {
					response.RenderFailure(ctx, errors.New("测试用例删除失败，该用例已分配项目任务"))
				}

			}
		}
		if len(passIds) > 0 {
			// 检测对应资产的测试组件分类下，是否还存在用例，如果不存在则删除该资产对应的组件和分类
			for _, itemId := range passIds {
				new(mongo_model.EvaluateItem).CheckAssetItemRelation(itemId)
			}

			// 删除测试用例
			match := bson.M{
				"_id": bson.M{"$in": passIds},
			}
			if changeInfo, err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).RemoveAll(match); err == nil {
				successNum = changeInfo.Removed
				// 删除测试用例关联的测试记录
				match = bson.M{
					"item_id": bson.M{"$in": passIds},
				}
				mongo.NewMgoSession(common.MC_EVALUATE_RECORD).RemoveAll(match)
				// 删除测试用例关联的漏洞
				mongo.NewMgoSession(common.MC_EVALUATE_VULNERABILITY).RemoveAll(match)

			} else {
				response.RenderFailure(ctx, err)
			}
		}
	}
	result := &qmap.QM{
		"success_number": successNum,
		"failure_number": total - successNum,
	}
	response.RenderSuccess(ctx, result)
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_item/tag 获取测试用例标签
 * @apiVersion 0.1.0
 * @apiName GetTag
 * @apiGroup EvaluateItem
 *
 * @apiDescription 获取测试用例标签
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   project_id  		    项目id
 * @apiParam {string}   [task_id]  			项目任务id
 *
 * @apiSuccessExample {json} 请求成功示例:
 *
 *	{
 *		"code": 0,
 *		"data": {
 *			"tag": [
 *				"w1",
 *				"w2"
 *			]
 *		}
 *	}
 */
func (this EvaluateItemController) GetTag(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	query := bson.M{
		"project_id": req.MustString("project_id"),
	}
	if taskId, has := req.TryString("task_id"); has {
		query["evaluate_task_id"] = taskId
	}
	tags := []interface{}{}
	if err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).Session.Find(query).Distinct("tag", &tags); err == nil {
		response.RenderSuccess(ctx, qmap.QM{"tag": tags})
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_item/tag 编辑测试用例标签
 * @apiVersion 0.1.0
 * @apiName EditTag
 * @apiGroup EvaluateItem
 *
 * @apiDescription 编辑测试用例标签
 *
 * @apiUse authHeader
 *
 * @apiParam 	{string}   		item_id  		测试用例id
 * @apiParam 	{string}   		project_id  	项目id
 * @apiParam 	{[]string}   	tag  			测试用例标签
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  	"item_id":"5fd7218624b64712a27f47e8",
 *  	"project_id":"5fd7218624b64712a27f47e8",
 *  	"tag":["a","b"]
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *
 *	{
 *		"code": 0
 *	}
 */
func (this EvaluateItemController) EditTag(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	params := bson.M{
		"_id":        req.MustString("item_id"),
		"project_id": req.MustString("project_id"),
	}
	update := bson.M{
		"$set": bson.M{
			"tag": req.SliceString("tag"),
		},
	}
	if _, err := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_ITEM).UpdateOne(context.Background(), params, update); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, gin.H{})
	}
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_item/navigation 获取测试项导航栏
 * @apiVersion 0.1.0
 * @apiName GetNavigation
 * @apiGroup EvaluateItem
 *
 * @apiDescription 获取测试项导航栏
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   project_id  			项目id
 * @apiParam {string}   [evaluate_task_id]  	项目任务id
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *  "code": 0,
 *  "data": [
 *     {
 *         "children": [
 *                 {
 *                     "children": [
 *                             {
 *                                 "id": "WiFi_组件名称_组件分类1",
 *                                 "label": "组件分类1"
 *                             }
 *                     ],
 *                     "id": "WiFi_组件名称",
 *                     "label": "组件名称"
 *                 }
 *         ],
 *         "id": "WiFi",
 *         "label": "WiFi"
 *     },
 *     {
 *         "children": [
 *                 {
 *                     "children": [
 *                             {
 *                                 "id": "行车记录仪_组件名称_组件分类1",
 *                                 "label": "组件分类1"
 *                             }
 *                     ],
 *                     "id": "行车记录仪_组件名称",
 *                     "label": "组件名称"
 *                 }
 *         ],
 *         "id": "行车记录仪",
 *         "label": "行车记录仪"
 *     }
 *  ]
 * }
 */
func (this EvaluateItemController) GetNavigation(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	params := qmap.QM{
		"e_project_id": req.MustString("project_id"),
	}

	if evaluateTaskId, has := req.TryString("evaluate_task_id"); has {
		taskInfo, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, qmap.QM{"e__id": bson.ObjectIdHex(evaluateTaskId)}).GetOne()
		if err == nil && (*taskInfo)["status"] == common.PTS_FINISH {
			params["e_last_task_id"] = evaluateTaskId
		} else {
			params["e_pre_bind"] = evaluateTaskId
		}
	}
	items, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params).SetLimit(50000).Get()
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	resultMap := qmap.QM{}
	moduleMap, err := new(mongo_model.EvaluateModule).GetModuleSlice()
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}
	assetMap, err := new(mongo_model.EvaluateAsset).GetAssetSlice(req.MustString("project_id"))
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}

	idList := map[string][]string{}
	for _, item := range *items {
		if assetMap[item["asset_id"].(string)] == nil {
			continue
		}
		targetName := assetMap[item["asset_id"].(string)].(qmap.QM)["name"].(string)
		if resultMap[targetName] == nil {
			resultMap[targetName] = qmap.QM{
				"label":       item["target_name"],
				"children":    qmap.QM{},
				"create_time": assetMap[item["asset_id"].(string)].(qmap.QM)["create_time"],
			}
		}
		idListKey := item["module_type_id"].(string) + targetName
		idList[idListKey] = append(idList[idListKey], item["_id"].(string))

		if (*moduleMap)[item["module_type_id"].(string)] != nil {
			moduleTypeItem := (*moduleMap)[item["module_type_id"].(string)].(map[string]interface{})
			childrenModule := resultMap[targetName].(qmap.QM)["children"].(qmap.QM)
			// 添加module name
			if childrenModule[moduleTypeItem["module_name"].(string)] == nil {
				resultMap[targetName].(qmap.QM)["children"].(qmap.QM)[moduleTypeItem["module_name"].(string)] = qmap.QM{
					"label":    moduleTypeItem["module_name"].(string),
					"children": qmap.QM{},
				}
			}

			childrenType := childrenModule[moduleTypeItem["module_name"].(string)].(qmap.QM)["children"].(qmap.QM)
			if childrenType[moduleTypeItem["module_type"].(string)] == nil {
				resultMap[targetName].(qmap.QM)["children"].(qmap.QM)[moduleTypeItem["module_name"].(string)].(qmap.QM)["children"].(qmap.QM)[moduleTypeItem["module_type"].(string)] = qmap.QM{
					"label": moduleTypeItem["module_type"].(string),
					"id":    moduleTypeItem["_id"].(bson.ObjectId).Hex(),
				}
			}
		}
	}

	result := map[interface{}]interface{}{}
	for tagName, item := range resultMap {
		itemSlice := []interface{}{}
		for _, moduleItem := range item.(qmap.QM)["children"].(qmap.QM) {
			moduleSlice := []interface{}{}
			for _, typeItem := range moduleItem.(qmap.QM)["children"].(qmap.QM) {
				idListKey := typeItem.(qmap.QM)["id"].(string) + tagName
				if idList[idListKey] != nil {
					typeItem.(qmap.QM)["item_ids"] = idList[idListKey]
				} else {
					typeItem.(qmap.QM)["item_ids"] = []string{}
				}
				moduleSlice = append(moduleSlice, typeItem)
			}
			moduleItem.(qmap.QM)["children"] = moduleSlice
			itemSlice = append(itemSlice, moduleItem)
		}
		item.(qmap.QM)["children"] = itemSlice
		item.(qmap.QM)["label"] = tagName
		result[item.(qmap.QM)["create_time"]] = item
	}

	response.RenderSuccess(ctx, custom_util.KSort(result))
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_item/asset_versions 查询测试用例资产版本信息
 * @apiVersion 0.1.0
 * @apiName GetAssetVersions
 * @apiGroup EvaluateItem
 *
 * @apiDescription 查询测试用例关联的资产版本信息
 *
 * @apiUse authHeader
 *
 * @apiParam 	{string}   	item_ids  			测试用例id（多个测试用例用"\|"连接）
 * @apiParam 	{string}     project_id  		项目id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  	"item_ids": "2|3",
 *  	"project_id":"5fd7218624b64712a27f47e8"
 * }
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": [
 *                 {
 *                         "id": "US4VX1KR",
 *                         "name": "test",
 *                         "versions": [
 *                                 "1.0"
 *                         ]
 *                 }
 *         ]
 * }
 */
func (this EvaluateItemController) GetAssetVersions(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if assetIds := new(mongo_model.EvaluateItem).GetProjectRelatedAssets(req.MustString("project_id"), req.MustString("item_ids")); len(assetIds) >= 0 {
		assetList := []qmap.QM{}
		for _, assetId := range assetIds {
			if assetInfo, err := new(mongo_model.EvaluateAsset).GetAssetAllVersions(assetId); err == nil {
				assetList = append(assetList, assetInfo)
			}
		}

		response.RenderSuccess(ctx, assetList)
	} else {
		response.RenderSuccess(ctx, gin.H{})
	}
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_item/upsert_record 添加/更新测试记录
 * @apiVersion 0.1.0
 * @apiName UpsertRecord
 * @apiGroup EvaluateItem
 *
 * @apiDescription 根据测试用例id,新增/更新测试记录
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           id                      测试记录id
 * @apiParam {string}           item_id                 测试用例id
 * @apiParam {string}           test_procedure          测试过程
 * @apiParam {json}             attachment         	   测试附件
 *
 * @apiParamExample {json}  请求参数示例:
 *     {
 *		"id":               "5fe464b1f98f923e40e8dd5f",
 *		"item_id":          "12314123",
 *		"test_procedure":   "测试过程",
 *     }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0
 * }
 */
func (this EvaluateItemController) UpsertRecord(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if err := new(mongo_model.EvaluateRecord).Upsert(*req, session.GetUserId(http_ctx.GetHttpCtx(ctx))); err == nil {
		response.RenderSuccess(ctx, gin.H{})
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {delete} /api/v1/evaluate_item/bulk_delete_record 批量删除测试记录
 * @apiVersion 0.1.0
 * @apiName BulkDeleteRecord
 * @apiGroup EvaluateItem
 *
 * @apiDescription 根据测试记录ID, 批量删除测试记录
 *
 * @apiUse authHeader
 *
 * @apiParam {[]string}		record_ids		测试记录id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *         "record_ids" : ["5fe464b1f98f923e40e8dd5f"]
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *			"failure_number":0,
 *			"success_number":1
 *     }
 * }
 */
func (this EvaluateItemController) BulkDeleteRecord(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	total := 0
	successNum := 0
	if _, has := req.TrySlice("record_ids"); has {
		ids := req.SliceString("record_ids")
		total = len(ids)
		idsObject := []bson.ObjectId{}
		for _, id := range ids {
			idsObject = append(idsObject, bson.ObjectIdHex(id))
		}
		// 删除测试记录
		match := bson.M{
			"_id": bson.M{"$in": idsObject},
		}
		if changeInfo, err := mongo.NewMgoSession(common.MC_EVALUATE_RECORD).RemoveAll(match); err == nil {
			successNum = changeInfo.Removed
		} else {
			response.RenderFailure(ctx, err)
			return
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"failure_number": total - successNum, "success_number": successNum})
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_item/audit_task_item 审核测试用例
 * @apiVersion 0.1.0
 * @apiName AuditTaskItem
 * @apiGroup EvaluateItem
 *
 * @apiDescription 审核项目任务中的测试用例
 *
 * @apiUse authHeader
 *
 * @apiParam {string}		evaluate_task_id					项目任务id
 * @apiParam {[]string}		ids									测试用例id
 * @apiParam {int}			audit_status						审核状态 （1:通过 -1:驳回）
 *
 * @apiParamExample {json}  请求参数示例:
 *     {
 *         "evaluate_task_id" : "5fe464b1f98f923e40e8dd51",
 *         "ids" : ["5fe464b1f98f923e40e8dd5f"],
 *         "audit_status" : 0
 *     }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0
 * }
 */
func (this EvaluateItemController) AuditTaskItem(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	if _, has := req.TrySlice("ids"); has {
		ids := req.SliceString("ids")
		auditStatus := req.Int("audit_status")
		for _, id := range ids {
			data := qmap.QM{
				"audit_status": auditStatus,
			}
			if _, err := new(mongo_model.EvaluateItem).Update(id, data); err != nil {
				response.RenderFailure(ctx, err)
				return
			}
		}
	}
	response.RenderSuccess(ctx, gin.H{})
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_item/audited_items 查询审核过的的测试用例
 * @apiVersion 0.1.0
 * @apiName GetAuditedItems
 * @apiGroup EvaluateItem
 *
 * @apiDescription 查询审核过的的测试用例
 *
 * @apiParam {string}		evaluate_task_id					项目任务id
 * @apiParam {int}			audit_status						审核状态 （1:通过 -1:驳回）
 *
 * @apiParamExample {json}  请求参数示例:
 *     {
 *         "evaluate_task_id" : "5fe464b1f98f923e40e8dd5f",
 *         "audit_status" : 0
 *     }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": [
 *                 {
 *                         "_id": "M003",
 *                         "asset_id": "US4VY61L",
 *                         "audit_status": 0,
 *                         "auto_test_level": "自动化",
 *                         "create_time": 1609233504081,
 *                         "evaluate_task_id": "5feb0032f98f9240eeb25cec",
 *                         "external_input": "fasfa",
 *                         "level": 1,
 *                         "module_type_id": "5fea11c35d50267995bc51d5",
 *                         "name": "fffaf",
 *                         "objective": "sfasfa",
 *                         "op_id": 64,
 *                         "project_id": "5fd7218624b64712a27f47e8",
 *                         "record_id": "5feb0032f98f9240eeb25ced",
 *                         "status": 0,
 *                         "tag": [],
 *                         "test_case_level": "基础测试",
 *                         "test_count": 0,
 *                         "test_method": "黑盒",
 *                         "test_phase": 1,
 *                         "test_procedure": "sfasfasfasf",
 *                         "test_script": [
 *                                 {
 *                                         "name": "mailchimp-sq.jpg",
 *                                         "value": "5feaf45ef98f9240b88a05c6"
 *                                 },
 *                                 {
 *                                         "name": "Snipaste_2020-03-10-114334.png",
 *                                         "value": "5feaf504f98f9240b88a05ca"
 *                                 }
 *                         ],
 *                         "test_sketch_map": [
 *                                 {
 *                                         "name": "mailchimp-sq.jpg",
 *                                         "value": "5feaf45ef98f9240b88a05c6"
 *                                 }
 *                         ],
 *                         "test_standard": "fasfasf",
 *                         "test_time": 0,
 *                         "update_time": 1609233504081,
 *                         "vul_number": 0
 *                 }
 *         ]
 * }
 */
func (this EvaluateItemController) GetAuditedItems(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	params := qmap.QM{
		"e_pre_bind":     req.MustString("evaluate_task_id"),
		"e_audit_status": req.MustInt("audit_status"),
	}
	mgoSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params)
	mgoSession.SetLimit(500000)
	lists, err := mgoSession.Get()
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	}

	response.RenderSuccess(ctx, lists)
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_item/records 查询测试用例的所有测试记录
 * @apiVersion 0.1.0
 * @apiName GetItemRecords
 * @apiGroup EvaluateItem
 *
 * @apiDescription 根据测试用例id,查询所有测试记录
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           item_id                 测试用例id
 *
 * @apiParamExample {json}  请求参数示例:
 *     {
 *			"item_id":          "12314123"
 *     }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": [
 *                 {
 *                         "_id": "5fea975524b6476b531c2d44",
 *                         "asset_id": "KJ8CCYXE",
 *                         "asset_version": "",
 *                         "attachment": null,
 *                         "conclude": "测试结论",
 *                         "create_time": 1609209873121,
 *                         "evaluate_task_id": "5fea975424b6476b531c2d41",
 *                         "item_id": "1111114",
 *                         "op_id": 0,
 *                         "project_id": "5fe995fa4806351585aa7568",
 *                         "risk_type": "设计",
 *                         "test_phase": 1,
 *                         "test_procedure": "测试过程"
 *                 }
 *         ]
 * }
 */
func (this EvaluateItemController) GetItemRecords(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	params := qmap.QM{
		"e_item_id":          req.MustString("item_id"),
		"e_evaluate_task_id": req.MustString("evaluate_task_id"),
	}
	mgoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_RECORD, params).SetLimit(1000000)
	mgoClient.SetTransformFunc(ItemRecordsTransformer)
	if res, err := mgoClient.Get(); err == nil {
		response.RenderSuccess(ctx, res)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func ItemRecordsTransformer(data qmap.QM) qmap.QM {
	// if opId := data.Int("op_id"); opId > 0 {
	// 	rpcClient := client.NewGRpcClient(common.AUTH_SERVICE, context.Background())
	// 	defer rpcClient.Close()
	// 	// 查询操作人员信息
	// 	userParam := qmap.QM{
	// 		"id": opId,
	// 	}
	// 	if rsp, err := auth.NewUserClient(rpcClient.Client).GetUserInfo(rpcClient.Ctx, &userParam); err == nil {
	// 		userTempt := rsp.Map("data")
	// 		if realname := userTempt.String("realname"); realname != "" {
	// 			data["op_name"] = realname
	// 		} else {
	// 			data["op_name"] = userTempt.String("username")
	// 		}
	//
	// 	} else {
	// 		data["op_name"] = ""
	// 	}
	// } else {
	// 	data["op_name"] = ""
	// }

	// 把上面的微服务调用改为调用本地auth
	// 根据数据记录中的op_id，做为用户标的id，查询用户信息，如果用户记录的realname字段有值，做为op_name, 否则用用户记录的username字段的值做为op_name
	if opId := data.Int("op_id"); opId > 0 {
		if userTempt, err := new(mysql_model.SysUser).GetUserInfo(opId); err == nil {
			if realname := userTempt.String("realname"); realname != "" {
				data["op_name"] = realname
			} else {
				data["op_name"] = userTempt.String("username")
			}
		} else {
			data["op_name"] = ""
		}
	} else {
		data["op_name"] = ""
	}

	if scanResults, err := new(mongo_model.ToolTaskResultBindTest).GetToolScanResult(data.Interface("_id").(bson.ObjectId).Hex()); err == nil {
		data["scan_results"] = scanResults
	} else {
		data["scan_results"] = []qmap.QM{}
	}
	return data
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_item/audit_record_item 审核测试记录
 * @apiVersion 0.1.0
 * @apiName AuditRecordItem
 * @apiGroup EvaluateItem
 *
 * @apiDescription 审核项目任务中的测试记录
 *
 * @apiUse authHeader
 *
 * @apiParam {string}		evaluate_task_id					项目任务id
 * @apiParam {[]string}	ids									测试用例id
 * @apiParam {int}			record_audit_status					审核状态(1:通过 -1:驳回)
 *
 * @apiParamExample {json}  请求参数示例:
 *     {
 *         "evaluate_task_id" : "5fe464b1f98f923e40e8dd51",
 *         "ids" : ["5fe464b1f98f923e40e8dd5f"],
 *         "record_audit_status" : 1
 *     }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0
 * }
 */
func (this EvaluateItemController) AuditRecordItem(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_TASK_ITEM).Session
	if _, has := req.TrySlice("ids"); has {
		ids := req.SliceString("ids")
		recordAuditStatus := req.MustInt("record_audit_status")
		for _, id := range ids {
			data := qmap.QM{
				"record_audit_status": recordAuditStatus,
			}
			if err := mgoSession.Update(bson.M{"_id": bson.ObjectIdHex(id)}, qmap.QM{"$set": data}); err != nil {
				response.RenderFailure(ctx, err)
				return
			}
		}
	}
	response.RenderSuccess(ctx, gin.H{})
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_item/audited_record_items 查询审核的测试记录
 * @apiVersion 0.1.0
 * @apiName GetAuditedRecordItems
 * @apiGroup EvaluateItem
 *
 * @apiDescription 查询项目任务中被审核的测试记录
 *
 * @apiUse authHeader
 *
 * @apiParam {string}		evaluate_task_id						项目任务id
 * @apiParam {int}			record_audit_status						审核状态(1:通过 -1:驳回)
 *
 * @apiParamExample {json}  请求参数示例:
 *     {
 *         "evaluate_task_id" : "5fe464b1f98f923e40e8dd5f",
 *         "record_audit_status" : -1
 *     }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": [
 *                 {
 *                         "_id": "M003",
 *                         "asset_id": "US4VY61L",
 *                         "audit_status": 0,
 *                         "auto_test_level": "自动化",
 *                         "create_time": 1609233504081,
 *                         "evaluate_task_id": "5feb0032f98f9240eeb25cec",
 *                         "external_input": "fasfa",
 *                         "level": 1,
 *                         "module_type_id": "5fea11c35d50267995bc51d5",
 *                         "name": "fffaf",
 *                         "objective": "sfasfa",
 *                         "op_id": 64,
 *                         "project_id": "5fd7218624b64712a27f47e8",
 *                         "record_id": "5feb0032f98f9240eeb25ced",
 *                         "status": 0,
 *                         "tag": [],
 *                         "test_case_level": "基础测试",
 *                         "test_count": 0,
 *                         "test_method": "黑盒",
 *                         "test_phase": 1,
 *                         "test_procedure": "sfasfasfasf",
 *                         "test_script": [
 *                                 {
 *                                         "name": "mailchimp-sq.jpg",
 *                                         "value": "5feaf45ef98f9240b88a05c6"
 *                                 },
 *                                 {
 *                                         "name": "Snipaste_2020-03-10-114334.png",
 *                                         "value": "5feaf504f98f9240b88a05ca"
 *                                 }
 *                         ],
 *                         "test_sketch_map": [
 *                                 {
 *                                         "name": "mailchimp-sq.jpg",
 *                                         "value": "5feaf45ef98f9240b88a05c6"
 *                                 }
 *                         ],
 *                         "test_standard": "fasfasf",
 *                         "test_time": 0,
 *                         "update_time": 1609233504081,
 *                         "vul_number": 0
 *                 }
 *         ]
 * }
 */
func (this EvaluateItemController) GetAuditedRecordItems(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	params := qmap.QM{
		"e_evaluate_task_id":    req.MustString("evaluate_task_id"),
		"e_record_audit_status": req.MustInt("record_audit_status"),
	}
	mgoSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK_ITEM, params)
	mgoSession.SetLimit(500000)
	if lists, err := mgoSession.Get(); err == nil {
		response.RenderSuccess(ctx, lists)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_item/complete_test 结束测试用例测试
 * @apiVersion 0.1.0
 * @apiName CompleteItemTest
 * @apiGroup EvaluateItem
 *
 * @apiDescription 结束测试用例测试
 *
 * @apiUse authHeader
 *
 * @apiParam {string}		id						测试用例id
 *
 * @apiParamExample {json}  请求参数示例:
 *     {
 *         "id" : "5fe464b1f98f923e40e8dd5f",
 *     }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0
 * }
 */
func (this EvaluateItemController) CompleteItemTest(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	id := req.MustString("id")
	Taskitem, err := new(mongo_model.EvaluateTaskItem).GetOne(id)
	if err == nil {
		if _, err := new(mongo_model.EvaluateRecord).GetOne(Taskitem.RecordId); err != nil {
			response.RenderFailure(ctx, errors.New("请在添加完成测试记录后再进行提交操作！"))
			return
		}
	} else {
		response.RenderFailure(ctx, errors.New("未找到该测试用例"))
		return
	}

	data := qmap.QM{
		"status": common.TIS_TEST_COMPLETE,
	}
	if err := mongo.NewMgoSession(common.MC_EVALUATE_TASK_ITEM).Session.Update(bson.M{"_id": bson.ObjectIdHex(id)}, qmap.QM{"$set": data}); err == nil {
		// 更新Item主表中 test_status 完成
		itemIds := []interface{}{Taskitem.ItemId}
		if err := new(mongo_model.EvaluateItem).ChangeTestStatus(itemIds, common.TIS_TEST_COMPLETE); err == nil {
			response.RenderSuccess(ctx, gin.H{})
		}
	}
	response.RenderFailure(ctx, err)
}

/**
 * apiType http
 * @api {get} /api/v1/evaluate_item/all_ids 查询所有测试用例id
 * @apiVersion 0.1.0
 * @apiName GetAllIds
 * @apiGroup EvaluateItem
 *
 * @apiDescription 查询所有测试用例id
 *
 * @apiUse authHeader
 *
 * @apiUse urlQueryParams
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": [
 *				"US5242TKTC050340M889",
 *              "KKDSHPWVTC090520A001",
 *               "KKDSHPWVTC090510A003"
 *         ]
 * }
 */
func (this EvaluateItemController) GetAllIds(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_ITEM)

	if postParams, has := req.TryInterface("post_params"); has {
		q := url.Values{}
		for key, val := range postParams.(map[string]interface{}) {
			if key == "ids" {
				continue
			}
			str := ""
			switch val.(type) {
			case int:
				str = strconv.Itoa(val.(int))
			case string:
				str = val.(string)
			case interface{}:
				data, _ := json.Marshal(val)
				str = string(data)
			}

			q.Set(key, str)
		}
		mgoSession.AddUrlQueryCondition(q.Encode())

		raw := qmap.New(postParams.(map[string]interface{}))
		if ids, has := raw.TryString("ids"); has {
			if idTemp, err := qmap.NewWithString(`{ "ids":` + ids + "}"); err == nil {
				params := qmap.QM{
					"in__id": idTemp.SliceString("ids"),
				}
				mgoSession.AddCondition(params)
			}
		}
	}

	ids := []string{}
	if result, err := mgoSession.SetLimit(1000000000).Get(); err == nil {
		for _, item := range *result {
			var itemQM qmap.QM = item
			ids = append(ids, itemQM.String("_id"))
		}
	}

	response.RenderSuccess(ctx, ids)
}

//@auto_generated_api_end
/**
 * apiType http
 * @api {post} /api/v1/evaluate_item/export/:id 导出测试用例
 * @apiVersion 0.1.0
 * @apiName Export
 * @apiGroup EvaluateItem
 *
 * @apiDescription 导出测试用例
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   id  项目id
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "success_number":10,
 *         "failure_number":0,
 *         "failure_info":""
 * }
 */
func (this EvaluateItemController) Export(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))
	params := qmap.QM{
		"e_project_id": request.ParamString(ctx, "id"),
	}
	items, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params).SetLimit(10000).Get()
	if err != nil {
		response.RenderFailure(ctx, err)
	}
	result := qmap.QM{}
	for index, item := range *items {
		itemQM := qmap.QM(item)
		// 查询测试分类和组件
		if module, err := new(mongo_model.EvaluateModule).FindById(itemQM.MustString("module_type_id")); err == nil {
			itemQM["module_name"] = module.ModuleName
			itemQM["module_type"] = module.ModuleType
		}
		// 查询测试难度
		levelInt := itemQM.Int("level")
		switch levelInt {
		case common.TEST_LEVEL_LOW:
			itemQM["level"] = "低"
		case common.TEST_LEVEL_MIDDLE:
			itemQM["level"] = "中"
		case common.TEST_LEVEL_HIGH:
			itemQM["level"] = "高"
		default:
			itemQM["level"] = "未知"
		}
		// 导出的用例 不需要带资产编号
		itemId := itemQM.String("_id")
		n := strings.Index(itemId, common.FirstTestCase)
		if n < 0 {
			itemQM["_id"] = itemId
		} else {
			itemQM["_id"] = itemId[n:]
		}
		result[strconv.Itoa(index)] = itemQM
	}

	excelObj := new(service.ExcelObj)
	excelObj.NewExcel()
	index := excelObj.NewSheet()
	excelObj.ContentTitle = []string{"测试用例ID", "资产ID", "测试组件", "测试分类", "测试用例名称", "测试难度", "测试用例级别", "测试方式", "测试目的", "自动化测试程度", "外部输入", "测试步骤", "测试标准"}
	excelObj.WriteTitle()

	//按顺序输出，与页面顺序保持一致
	for i := 0; i < len(result); i++ {
		if data, exit := result[strconv.Itoa(i)]; exit {
			item := data.(qmap.QM)
			excelObj.Content = []interface{}{item.String("_id"), item.String("asset_id"), item.String("module_name"),
				item.String("module_type"), item.String("name"), item.String("level"),
				item.String("test_case_level"), item.String("test_method"), item.String("objective"),
				item.String("auto_test_level"), item.String("external_input"), item.String("test_procedure"), item.String("test_standard")}
			excelObj.WriteContent()
			excelObj.ExcelFile.SetActiveSheet(index)
		}
	}

	data := excelObj.Output()
	ctx.Writer.WriteHeader(http.StatusOK)
	ctx.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", "report.xlsx"))
	ctx.Header("Content-Type", "*")
	ctx.Header("Accept-Length", fmt.Sprintf("%d", len(data.Bytes())))
	ctx.Writer.Write(data.Bytes())

	response.RenderSuccess(ctx, result)
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_item/import/:id 导入测试用例
 * @apiVersion 0.1.0
 * @apiName Import
 * @apiGroup EvaluateItem
 *
 * @apiDescription 导入测试用例
 *
 * @apiUse authHeader
 *
 * @apiParam {string}   project_id  项目id
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "success_number":10,
 *         "failure_number":0,
 *         "failure_info":""
 * }
 */
func (this EvaluateItemController) Import(ctx *gin.Context) {
	fmt.Println("import")
	// todo 用到了rpc流处理，待拆分
	// //获取文件
	fileName := ctx.Request.FormValue("file")
	file, header, _ := ctx.Request.FormFile("file")
	if fileName == "" && header != nil {
		fileName = header.Filename
	}
	projectId := ctx.Request.FormValue("project_id")

	if file != nil {
		fmt.Println("file")
		if fileContent, err := ioutil.ReadAll(file); err == nil {
			xlFile, err := xlsx.OpenBinary(fileContent)
			if err != nil {
				response.RenderFailure(ctx, errors.New("上传的xlsx格式不正确"))
			}
			sheet := xlFile.Sheets[0]
			rowNum := sheet.MaxRow
			data := []qmap.QM{}

			//简单校验一下excel里内容是否是测试用例模板
			checkRow, _ := sheet.Row(0)
			if checkRow.GetCell(0).Value != "测试用例ID" {
				response.RenderFailure(ctx, errors.New("ecxel不是测试用例模板"))
			}
			for i := 1; i < rowNum; i++ {
				row, _ := sheet.Row(i)
				item := qmap.QM{}
				item["row_number"] = i + 1
				item["project_id"] = projectId
				item["_id"] = row.GetCell(0).Value
				item["asset_id"] = row.GetCell(1).Value
				item["module_name"] = row.GetCell(2).Value
				item["module_type"] = row.GetCell(3).Value
				item["name"] = row.GetCell(4).Value
				level := row.GetCell(5).Value
				switch level {
				case "低":
					item["level"] = common.TEST_LEVEL_LOW
				case "中":
					item["level"] = common.TEST_LEVEL_MIDDLE
				case "高":
					item["level"] = common.TEST_LEVEL_HIGH
				default:
					item["level"] = common.TEST_LEVEL_DEFAULT
				}
				item["test_case_level"] = row.GetCell(6).Value
				item["test_method"] = row.GetCell(7).Value
				item["objective"] = row.GetCell(8).Value
				item["auto_test_level"] = row.GetCell(9).Value
				item["external_input"] = row.GetCell(10).Value
				item["test_procedure"] = row.GetCell(11).Value
				item["test_standard"] = row.GetCell(12).Value
				if item["asset_id"] == "" && item["module_name"] == "" && item["module_type"] == "" && item["name"] == "" {
					continue
				}
				data = append(data, item)
				fmt.Println("item:", item)
			}
			successNumber, failureNumber, errInfo := new(mongo_model.EvaluateItem).Import(data, 0)
			errList := []qmap.QM{}
			json.Unmarshal([]byte(errInfo), &errList)
			result := qmap.QM{
				"success_number": successNumber,
				"failure_number": failureNumber,
				"failure_info":   errList,
			}
			if len(errList) != 0 {
				response.Render(ctx, response.RC_FAILURE, "导入用例存在错误", &result)
				return
			}
			response.RenderSuccess(ctx, &result)

		} else {
			response.RenderFailure(ctx, err)
		}
	}
	// var userId int64
	// rpcClient := client.NewGRpcClient(common.PM_SERVICE, http_ctx.NewOutputContext(ctx))
	// defer rpcClient.Close()
	// if uploadStream, err := project_manage.NewItemClient(rpcClient.Client).Import(rpcClient.Ctx); err == nil {
	// 	//如果文件不为空，则开始传输文件内容
	// 	if file != nil {
	// 		//每次最大传输3M
	// 		var fileContent [3145728]byte
	// 		for {
	// 			if len, err := file.Read(fileContent[:]); len == 0 {
	// 				if err == io.EOF {
	// 					break
	// 				} else {
	// 					panic(err)
	// 				}
	// 			} else {
	// 				userId = session.GetUserId(rpcClient.Ctx)
	// 				pushErr := uploadStream.Send(&project_manage.ImportRequest{
	// 					FileContent: fileContent[:len],
	// 				})
	// 				if pushErr != nil {
	// 					panic(pushErr)
	// 				}
	// 			}
	// 		}
	// 	}
	// 	if pushErr := uploadStream.Send(&project_manage.ImportRequest{
	// 		UserId:    userId,
	// 		ProjectId: ctx.Request.FormValue("project_id"),
	// 	}); pushErr != nil {
	// 		ctx.JSON(400, gin.H{
	// 			"code": -1,
	// 			"msg":  pushErr.Error(),
	// 		})
	// 	}
	// 	if resp, closeErr := uploadStream.CloseAndRecv(); closeErr == nil {
	// 		errInfo := resp.ErrorInfo
	// 		errList := []qmap.QM{}
	// 		json.Unmarshal([]byte(errInfo), &errList)
	// 		result := qmap.QM{
	// 			"success_number": resp.SuccessNumber,
	// 			"failure_number": resp.FailureNumber,
	// 			"failure_info":     errList,
	// 		}
	// 		if len(errList) != 0 {
	// 			response.Render(ctx, response.RC_FAILURE, "导入用例存在错误", &result)
	// 			return
	// 		}
	// 		response.RenderSuccess(ctx, &result)
	// 	} else {
	// 		response.RenderFailure(ctx, closeErr)
	//
	// 	}
	// } else {
	// 	response.RenderFailure(ctx, err)
	// }
}

/**
 * apiType http
 * @api {post} /api/v1/evaluate_item/templates 从模板中添加测试项
 * @apiVersion 0.1.0
 * @apiName CreateItemFromTemplate
 * @apiGroup EvaluateItem
 *
 * @apiDescription 从模板中添加测试项
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           project_id    		    项目ID
 * @apiParam {string}           item_id    		        测试项ID
 *
 *curl http://10.16.133.118:3001/api/v1/evaluate_item/templates
 *
 * @apiParamExample {json}  请求参数示例:
 *
 * {
 *    "project_id":"5f4cbd796e655c20df5a7dd9",
 *    "target_id":"5f4fb17f6e655c23aa64f9fd",
 *    "evaluate_type":"APP",
 *    "module_name":"测试组件",
 *    "module_type":"测试分类"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 *	{
 *		"code": 0
 *	}
 */
func (this EvaluateItemController) CreateItemFromTemplate(ctx *gin.Context) {
	response.RenderSuccess(ctx, nil)
}

/**
 * apiType http
 * @api {delete} /api/v1/pre_bind 预绑定任务或解绑任务
 * @apiVersion 0.1.0
 * @apiName PreBind
 * @apiGroup EvaluateItem
 *
 * @apiDescription 批量删除测试用例
 *
 * @apiUse authHeader
 *
 * @apiParam {[]string}   ids  测试用例id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  "ids":[
 *		"1",
 *		"2"
 * 	],
 *  "bind_type" : "bind"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": {
 *                 "result": true,
 *         }
 * }
 */
func (this EvaluateItemController) PreBind(ctx *gin.Context) {
	req := &qmap.QM{}
	*req = req.Merge(*request.GetRequestBody(ctx))
	*req = req.Merge(*request.GetRequestQueryParams(ctx))

	bindType := req.MustString("bind_type")
	taskId := req.MustString("task_id")
	if _, has := req.TrySlice("ids"); has {
		itemIds := req.SliceString("ids")
		if bindType == "bind" { // 预绑定
			new(mongo_model.EvaluateItem).BindEvaluateTask(itemIds, taskId, common.IS_PREBIND)
		} else { // 预解绑
			new(mongo_model.EvaluateItem).PreDeleteTask(itemIds)
		}
		response.RenderSuccess(ctx, qmap.QM{"result": true})
		return
	}
	response.RenderSuccess(ctx, qmap.QM{"result": false})
	return
}

/**
 * apiType http
 * @api {delete} /api/v1/clearBind 预绑定任务或解绑任务
 * @apiVersion 0.1.0
 * @apiName ClearBind
 * @apiGroup EvaluateItem
 *
 * @apiDescription 清除预绑定
 *
 * @apiUse authHeader
 *
 * @apiParam {[]string}   ids  测试用例id
 *
 * @apiParamExample {json}  请求参数示例:
 * {
 *  "task_id" : "123123123"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *         "code": 0,
 *         "data": {
 *                 "result": true,
 *         }
 * }
 */
func (this EvaluateItemController) ClearBind(ctx *gin.Context) {
	taskId := ctx.Param("task_id")

	// 第一步，撤回移除的预绑定item
	selector := bson.M{
		"pre_bind": taskId,
	}
	updateItem := bson.M{
		"$set": qmap.QM{
			"is_pre_delete": common.NOT_PREDEL,
		},
	}
	if _, err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).UpdateAll(selector, updateItem); err != nil {
		fmt.Println(err)
	}

	// 第二步，将item状态重置
	selector = bson.M{
		"pre_bind":    taskId,
		"is_pre_bind": common.IS_PREBIND,
	}
	updateItem = bson.M{
		"$set": qmap.QM{
			"pre_bind":    "",
			"status":      common.EIS_FREE,
			"is_pre_bind": common.NOT_PREBIND,
		},
	}
	if _, err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).UpdateAll(selector, updateItem); err != nil {
		fmt.Println(err)
	}

	response.RenderSuccess(ctx, &qmap.QM{"result": true})
}
