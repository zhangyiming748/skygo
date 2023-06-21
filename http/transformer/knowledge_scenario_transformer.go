package transformer

import (
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/mysql_model"
)

type KnowledgeScenarioTransformer struct {
	orm.Transformer
}

func (h *KnowledgeScenarioTransformer) ModifyItem(data qmap.QM) qmap.QM {
	createUserId := data.MustInt("create_user_id")
	userModel, err := mysql_model.SysUserFindById(createUserId)
	if err == nil {
		data["create_user_name"] = userModel.Realname
	} else {
		data["create_user_name"] = "-"
	}

	demandId := data.MustInt("demand_id")
	model, has := mysql_model.KnowledgeDemandFindById(demandId)
	if has {
		data["demand_name"] = model.Name
	} else {
		data["demand_name"] = ""
	}

	scenarioId := data.MustInt("id")
	n := new(mysql_model.KnowledgeTestCase).GetCountByScenarioId(scenarioId)
	data["test_case_count"] = n
	// 查询关联安全需求
	// 查询关联安全需求条目
	// 查询场景对应的章节名称
	_, chapterId := new(mysql_model.KnowledgeScenarioChapter).GetDemandChapter(scenarioId)
	data["demand_chapter_id"] = chapterId
	return data
}
