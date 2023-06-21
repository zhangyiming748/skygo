package controller

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"skygo_detection/guardian/src/net/qmap"

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

type KnowledgeScenarioController struct{}

/*
*
  - apiType http
  - @api {get} /api/xxxx xxx
  - @apiVersion 0.1.0
  - @apiName GetAll
  - @apiGroup KnowledgeScenario
    *
  - @apiDescription xxx
    *
  - @apiUse authHeader
    *
  - @apiSuccessExample {json} 请求成功示例:
*/
func (this KnowledgeScenarioController) GetAll(ctx *gin.Context) {
	queryParams := ctx.Request.URL.RawQuery
	s := mysql.GetSession()

	// 查询组键
	widget := orm.PWidget{}
	widget.SetQueryStr(queryParams)
	widget.SetTransformer(&transformer.KnowledgeScenarioTransformer{})
	all := widget.PaginatorFind(s, &[]mysql_model.KnowledgeScenario{})
	response.RenderSuccess(ctx, all)
}

func (this KnowledgeScenarioController) GetOne(ctx *gin.Context) {
	id := request.ParamString(ctx, "id")
	s := mysql.GetSession()
	s.Where("id=?", id)

	w := orm.PWidget{}
	w.SetTransformer(&transformer.KnowledgeScenarioTransformer{})
	result, err := w.One(s, &mysql_model.KnowledgeScenario{})

	if err == nil {
		response.RenderSuccess(ctx, result)
	} else {
		response.RenderFailure(ctx, err)
	}
}

func (this KnowledgeScenarioController) Create(ctx *gin.Context) {
	uid := session.GetUserId(http_ctx.GetHttpCtx(ctx))
	req := request.GetRequestBody(ctx)
	demandChapterId := req.SliceInt("demand_chapter_id")
	// 创建场景数据
	scenario := new(mysql_model.KnowledgeScenario)
	scenario.Name = req.MustString("name")
	scenario.DemandId = req.Int("demand_id")
	scenario.Detail = req.MustString("detail")
	scenario.Describe = scenario.Detail
	scenario.CreateTime = int(time.Now().Unix())
	scenario.UpdateTime = scenario.CreateTime
	scenario.CreateUserId = int(uid)
	scenario.LastOpId = scenario.CreateUserId
	_, err := scenario.Create(demandChapterId)
	if err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, orm.StructToMap(*scenario))
	}
}

/**
 * apiType http
 * @api {put} /api/v1/knowledge_scenarios/:id 更新安全检测场景
 * @apiVersion 1.0.0
 * @apiName update
 * @apiGroup KnowledgeScenario
 *
 * @apiDescription 更新安全检测场景
 *
 * @apiUse authHeader
 *
 * @apiParam {string}           tag 		场景标签
 *
 * @apiParamExample {json} 请求参数示例:
 * {
 *     "tag":"安全"
 * }
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": {
 *         "create_time": 0,
 *         "create_user_id": 0,
 *         "demand_id": 0,
 *         "detail": "",
 *         "id": 2,
 *         "last_op_id": 0,
 *         "name": "",
 *         "tag": "安全",
 *         "update_time": 1636509940
 *     },
 *     "msg": ""
 * }
 */
func (this KnowledgeScenarioController) Update(ctx *gin.Context) {
	id := request.ParamInt(ctx, "id")
	scenario := new(mysql_model.KnowledgeScenario)
	scenario.Id = id

	lastOpId := session.GetUserId(http_ctx.GetHttpCtx(ctx))

	if request.IsExist(ctx, "name") {
		scenario.Name = request.MustString(ctx, "name")
	}
	if request.IsExist(ctx, "demand_id") {
		scenario.DemandId = request.MustInt(ctx, "demand_id")
	}
	if request.IsExist(ctx, "detail") {
		scenario.Detail = request.MustString(ctx, "detail")
		// detail和describe保持一致
		scenario.Describe = scenario.Detail
	}
	if request.IsExist(ctx, "tag") {
		scenario.Tag = request.String(ctx, "tag")
	}

	updateTime := time.Now().Unix()
	scenario.UpdateTime = int(updateTime)

	scenario.LastOpId = int(lastOpId)

	demandChapterIds := request.Slice(ctx, "demand_chapter_id")
	{ // 更新场景关联的章节的时候，如果删除的章节已经关联了场景中的测试用例，则删除失败
		preChapterIds := new(mysql_model.KnowledgeTestCaseChapter).GetDemandChapterIdsBySenarioId(id)
		if len(preChapterIds) > 0 {
			for _, val := range preChapterIds {
				flag := false
				for _, chapterId := range demandChapterIds {
					fChapterId := chapterId.(float64)
					iChapterId := int(fChapterId)
					if iChapterId == val.DemandChapterId {
						flag = true
						break
					}
				}
				if flag == false {
					panic(errors.New(`"请勿删除已经关联了测试用例的"安全需求条目"`))
				}
			}
		}
	}
	if _, err := scenario.Update(demandChapterIds); err != nil {
		response.RenderFailure(ctx, err)
	} else {
		response.RenderSuccess(ctx, orm.StructToMap(*scenario))
	}

}

func (this KnowledgeScenarioController) BulkDelete(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)

	successNum := 0
	if _, has := req.TrySlice("ids"); has {
		ids := req.SliceInt("ids")

		s := mysql.GetSession()
		s.Begin()
		for _, id := range ids {
			// 1先删除场景关系表，然后删除场景
			// 1.1删除场景关系表
			scenarioChapter := new(mysql_model.KnowledgeScenarioChapter)
			scenarioChapter.ScenarioId = id
			_, err := s.Delete(scenarioChapter)
			if err != nil {
				s.Rollback()
				log.GetHttpLogLogger().Error(fmt.Sprintf("%v", err))
				response.RenderFailure(ctx, err)
			}
			// 1.2删除场景表
			_, err = s.ID(id).Delete(&mysql_model.KnowledgeScenario{})
			if err != nil {
				s.Rollback()
				log.GetHttpLogLogger().Error(fmt.Sprintf("%v", err))
				response.RenderFailure(ctx, err)
			} else {
				successNum++
			}
			err = s.Commit()
			if err != nil {
				s.Rollback()
				log.GetHttpLogLogger().Error(fmt.Sprintf("%v", err))
				response.RenderFailure(ctx, err)
			}
		}
	}
	response.RenderSuccess(ctx, qmap.QM{"number": successNum})
}

/*
*
  - apiType http
  - @api {post} /api/v1/ks_tag 更新标签
  - @apiVersion 1.0.0
  - @apiName updateTag
  - @apiGroup KnowledgeScenario
    *
  - @apiDescription 单独修改安全检测场景标签
    *
  - @apiUse authHeader
    *
  - @apiParam {int}              id          场景id
  - @apiParam {string}           tag 		场景标签
    *
  - @apiParamExample {json} 请求参数示例:
  - {
  - "id":30
  - "tag":"安全"
  - }
    *
  - @apiSuccessExample {json} 请求成功示例:
  - {
  - "code": 0,
  - "data": {
  - "id": 30,
  - "name": "",
  - "demand_id": 0,
  - "detail": "",
  - "create_time": 0,
  - "update_time": 0,
  - "create_user_id": 0,
  - "last_op_id": 0,
  - "tag": "安全"
  - },
  - "msg": ""
    }
*/
func (this KnowledgeScenarioController) UpdateTag(ctx *gin.Context) {
	req := request.GetRequestBody(ctx)
	id := request.ParamInt(ctx, "id")

	tag := req.MustString("tag")
	scenario := new(mysql_model.KnowledgeScenario)
	if _, err := scenario.UpdateTag(id, tag); err != nil {
		response.RenderFailure(ctx, err)
	}
	response.RenderSuccess(ctx, scenario)

}

/**
 * apiType http
 * @api {get} /api/v1/knowledge_scenario/chapter_tree/:id 场景安全需求-章节级联列表
 * @apiVersion 0.1.0
 * @apiName ChapterTree
 * @apiGroup KnowledgeScenario
 *
 * @apiDescription 知识库安全需求-章节级联列表
 *
 * @apiParam {int}      	   knowledge_scenario_id    	安全场景id
 *
 * @apiUse authHeader
 *
 * @apiSuccessExample {json} 请求成功示例:
 * {
 *     "code": 0,
 *     "data": [
 *         {
 *             "id": 70,
 *             "code": "1.1",
 *             "title": "1.1test",
 *             "children": []
 *         }
 *     ],
 *     "msg": ""
 * }
 */
func (this KnowledgeScenarioController) ChapterTree(ctx *gin.Context) {
	scenarioId := request.ParamInt(ctx, "id")
	demandId, _ := new(mysql_model.KnowledgeScenarioChapter).GetDemandChapter(scenarioId)
	data := mysql_model.KnowledgeScenarioTree(demandId, []int{})
	response.RenderSuccess(ctx, data)
}
