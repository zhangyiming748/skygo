package controller

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/custom_util"
	"skygo_detection/http/transformer"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/log"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/lib/common_lib/request"
	"skygo_detection/lib/common_lib/response"
	"skygo_detection/lib/common_lib/session"
	"skygo_detection/mysql_model"
)

type KnowledgeDemandController struct{}

/**
 * apiType http
 * @api {get} /api/v1/knowledge_demands 知识库安全需求列表
 * @apiVersion 0.1.0
 * @apiName GetAll
 * @apiGroup KnowledgeDemand
 *
 * @apiDescription 知识库安全需求列表
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
func (this KnowledgeDemandController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	s := mysql.GetSession()

	// 查询组键
	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	widget.AddSorter(*(orm.NewSorter("id", 1)))
	all := widget.PaginatorFind(s, &[]mysql_model.KnowledgeDemand{})
	response.RenderSuccess(ctx, all)
}

/**
 * apiType http
 * @api {post} /api/v1/knowledge_demands 知识库安全需求创建
 * @apiVersion 0.1.0
 * @apiName Create
 * @apiGroup KnowledgeDemand
 *
 * @apiDescription 知识库安全需求创建
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
func (this KnowledgeDemandController) Create(ctx *gin.Context) {
	// 表单
	form := &mysql_model.KnowledgeDemandCreateForm{}
	form.Name = request.MustString(ctx, "name")
	form.Category = request.MustInt(ctx, "category")
	form.Code = request.MustString(ctx, "code")
	form.ImplementTime = request.MustInt(ctx, "implement_time")
	form.Detail = request.String(ctx, "detail")

	uid := session.GetUserId(http_ctx.GetHttpCtx(ctx))
	id, err := mysql_model.KnowledgeDemandCreateFromForm(form, int(uid))

	if err == nil {
		response.RenderSuccess(ctx, gin.H{"id": id})
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {get} /api/v1/knowledge_demands/:id 知识库安全需求详情
 * @apiVersion 0.1.0
 * @apiName GetOne
 * @apiGroup KnowledgeDemand
 *
 * @apiDescription 知识库安全需求详情
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
func (this KnowledgeDemandController) GetOne(ctx *gin.Context) {
	id := request.ParamString(ctx, "id")
	s := mysql.GetSession()
	s.Where("id=?", id)

	w := orm.PWidget{}
	result, err := w.One(s, &mysql_model.KnowledgeDemand{})

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
func (this KnowledgeDemandController) Update(ctx *gin.Context) {
	data := request.GetRequestBody(ctx)
	id := request.ParamInt(ctx, "id")

	if model, err := mysql_model.KnowledgeDemandUpdateById(id, *data); err == nil {
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
func (this KnowledgeDemandController) BulkDelete(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)

	successNum := 0
	if _, has := req.TrySlice("ids"); has {
		ids := req.SliceInt("ids")

		for _, id := range ids {
			// todo 检查是否可以删除

			err := mysql_model.NewKnowledgeDemandDeleteById(id)
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
 * @api {get} /api/v1/knowledge_demand/select_list 知识库安全需求下拉列表
 * @apiVersion 0.1.0
 * @apiName SelectList
 * @apiGroup KnowledgeDemand
 *
 * @apiDescription 知识库安全需求下拉列表
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
func (this KnowledgeDemandController) SelectList(ctx *gin.Context) {
	lists := mysql_model.KnowLedgeDemandSelectList()
	response.RenderSuccess(ctx, lists)
}

/**
 * apiType http
 * @api {get} /api/v1/knowledge_demand/chapter_tree/:id 知识库安全需求-章节级联列表
 * @apiVersion 0.1.0
 * @apiName ChapterTree
 * @apiGroup KnowledgeDemand
 *
 * @apiDescription 知识库安全需求-章节级联列表
 *
 * @apiParam {string}      	   knowledge_demand_id    	安全需求id
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 1
 *     },
 *     "msg": ""
 * }
 */
func (this KnowledgeDemandController) ChapterTree(ctx *gin.Context) {
	id := request.ParamInt(ctx, "id")
	data := mysql_model.KnowledgeDemandTree(id, []int{})
	response.RenderSuccess(ctx, data)
}

/**
 * apiType http
 * @api {post} /api/v1/knowledge_demand/chapter_all?knowledge_demand_id=1 知识库安全需求-章节列表页
 * @apiVersion 0.1.0
 * @apiName ChapterAll
 * @apiGroup KnowledgeDemand
 *
 * @apiDescription 知识库安全需求-章节列表页
 *
 * @apiParam {string}      	   knowledge_demand_id    	安全需求id
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 1
 *     },
 *     "msg": ""
 * }
 */
func (this KnowledgeDemandController) ChapterAll(ctx *gin.Context) {
	// 安全需求id, 是通过?knowledge_demand_id=1方式传的
	// data := request.GetRequestBody(ctx)
	// knowledgeDemandId := data.MustInt("knowledge_demand_id")

	knowledgeDemandId := request.QueryInt(ctx, "knowledge_demand_id")
	id := request.QueryInt(ctx, "id")

	queryParams := ctx.Request.URL.RawQuery // ?后的参数
	s := mysql.GetSession()
	if id > 0 {
		kdtn := new(mysql_model.KnowLedgeDemandTreeNode)
		kdtn.Id = id
		kdtn.FetchChild([]int{})

		ids := []int{}
		ids = mysql_model.GetAllChild(kdtn)
		ids = append(ids, id)
		s = s.In("id", ids).OrderBy("code desc")
	} else if knowledgeDemandId > 0 {
		s = s.Where("knowledge_demand_id = ?", knowledgeDemandId).OrderBy("code desc")
	}

	// 查询组件
	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	widget.SetTransformer(&transformer.KnowledgeDemandChapterTransformer{})
	all := widget.PaginatorFind(s, &[]mysql_model.KnowledgeDemandChapter{})
	response.RenderSuccess(ctx, all)
}

/**
 * apiType http
 * @api {post} /api/v1/knowledge_demand/chapter_one 知识库安全需求-章节详情
 * @apiVersion 0.1.0
 * @apiName ChapterOne
 * @apiGroup KnowledgeDemand
 *
 * @apiDescription 知识库安全需求-章节详情
 *
 * @apiParam {string}      	   id    	章节id
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 1
 *     },
 *     "msg": ""
 * }
 */
func (this KnowledgeDemandController) ChapterOne(ctx *gin.Context) {
	data := request.GetRequestBody(ctx)
	id := data.MustInt("id") // 章节id

	s := mysql.GetSession().Where("id=?", id)

	// 查询组件
	w := orm.PWidget{}
	result, err := w.One(s, &mysql_model.KnowledgeDemandChapter{})

	if err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

/**
 * apiType http
 * @api {post} /api/v1/knowledge_demand/chapter_create 知识库安全需求-添加章节
 * @apiVersion 0.1.0
 * @apiName ChapterCreate
 * @apiGroup KnowledgeDemand
 *
 * @apiDescription 知识库安全需求-添加章节
 *
 * @apiParam {string}      	   knowledge_demand_id    	安全需求id
 * @apiParam {string}      	   code    			   		章节编号
 * @apiParam {string}      	   title    		    	章节标题
 * @apiParam {int} 	     	   parent_id    			父章节标题id
 * @apiParam {string} 	           content			    	内容
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 1
 *     },
 *     "msg": ""
 * }
 */
func (this KnowledgeDemandController) ChapterCreate(ctx *gin.Context) {
	// 表单
	form := &mysql_model.KnowledgeDemandChapterCreateForm{}
	form.Code = request.MustString(ctx, "code")
	form.Content = request.String(ctx, "content")
	form.Title = request.MustString(ctx, "title")
	form.KnowledgeDemandId = request.MustInt(ctx, "knowledge_demand_id")
	form.ParentId = request.DefaultInt(ctx, "parent_id", 0)
	model, err := mysql_model.KnowledgeDemandChapterCreate(form)
	if err != nil {
		response.RenderFailure(ctx, err)
		return
	} else {
		// if model.ParentCode == "" {
		// 	model.ParentCode = "--" // 前端展示用
		// }
		response.RenderSuccess(ctx, model)
		return
	}
}

/**
 * apiType http
 * @api {post} /api/v1/knowledge_demand/chapter_delete 知识库安全需求-章节删除
 * @apiVersion 0.1.0
 * @apiName ChapterDelete
 * @apiGroup KnowledgeDemand
 *
 * @apiDescription  知识库安全需求-章节删除
 *
 * @apiParam {string}      	   id    			    	安全需求的章节id
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
func (this KnowledgeDemandController) ChapterDelete(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	ids := req.SliceInt("ids")
	var num int
	for _, id := range ids {
		if tmpNum, err := mysql_model.KnowledgeDemandChapterDeleteById(id); err != nil {
			response.RenderFailure(ctx, err)
			return
		} else {
			num += int(tmpNum)
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"number": num})
}

/**
 * apiType http
 * @api {post} /api/v1/knowledge_demand/chapter_update 知识库安全需求-章节更新
 * @apiVersion 0.1.0
 * @apiName ChapterUpdate
 * @apiGroup KnowledgeDemand
 *
 * @apiDescription  知识库安全需求-章节更新
 *
 * @apiParam {string}      	                        id    			    	安全需求的章节id
 * @apiParam {string}      	                        code    			    章节编号
 * @apiParam {string}      	                       	title    		    	章节标题
 * @apiParam {int} 	     	                       	parent_id    			父章节标题id
 * @apiParam {string} 				     	            content			    	内容
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "id": 1
 *     },
 *     "msg": ""
 * }
 */
func (this KnowledgeDemandController) ChapterUpdate(ctx *gin.Context) {
	data := request.GetRequestBody(ctx)
	id := data.MustInt("id")

	if model, err := mysql_model.KnowledgeDemandChapterUpdateById(id, *data); err == nil {
		response.RenderSuccess(ctx, model)
		return
	} else {
		response.RenderFailure(ctx, err)
		return
	}
}

/**
 * apiType http
 * @api {get} /api/v1/knowledge_demand/chapter_select_list?knowledge_demand_id=1 知识库安全需求条目的下拉列表
 * @apiVersion 0.1.0
 * @apiName ChapterSelectList
 * @apiGroup KnowledgeDemand
 *
 * @apiDescription 知识库安全需求条目的下拉列表
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
func (this KnowledgeDemandController) ChapterSelectList(ctx *gin.Context) {
	knowledgeDemandId := request.QueryInt(ctx, "knowledge_demand_id")
	lists := mysql_model.KnowLedgeDemandChapterSelectList(knowledgeDemandId)
	response.RenderSuccess(ctx, lists)
}

func (this KnowledgeDemandController) CodetList(ctx *gin.Context) {
	lists := mysql_model.GetKnowLedgeDemandCodeList()
	response.RenderSuccess(ctx, lists)
}

// 获取需求/章节下的测试用例
func (this KnowledgeDemandController) GetTestCases(ctx *gin.Context) {
	data := request.GetRequestBody(ctx)
	demandId := data.MustInt("demand_id")
	chapterIds := data.MustSlice("demand_chapter_ids")
	lists := new(mysql_model.KnowledgeTestCaseChapter).GetTestCases(demandId, chapterIds)
	var testCasesId = make([]int, 0)
	for _, list := range lists {
		testCasesId = append(testCasesId, list.TestCaseId)
	}
	cases := mysql_model.KnowledgeTestCaseFindByIds(testCasesId)
	if result, err := custom_util.StructToList(cases); err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}
