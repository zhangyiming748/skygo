package mysql_model

import (
	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/mysql"
)

type KnowledgeScenarioChapter struct {
	Id              int `xorm:"not null pk comment('自增主键id') INT(11)"`
	ScenarioId      int `xorm:"not null comment('安全检测场景id') INT(11)"`
	DemandId        int `xorm:"not null comment('关联知识库需求id') INT(11)"`
	DemandChapterId int `xorm:"not null comment('关联知识库需求章节id') INT(11)"`
}

func (this *KnowledgeScenarioChapter) Create() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}

func (this *KnowledgeScenarioChapter) Update(cols ...string) (int64, error) {
	return mysql.GetSession().Table(this).ID(this.Id).Cols(cols...).Update(this)
}

func (this *KnowledgeScenarioChapter) Remove() (int64, error) {
	return mysql.GetSession().Delete(this)
}

func (this *KnowledgeScenarioChapter) GetDemandChapter(scenarioId int) (int, []int) {
	result := make([]int, 0)
	demandChapters := []KnowledgeScenarioChapter{}
	session := mysql.GetSession()
	session.Where("scenario_id=?", scenarioId)
	session.Find(&demandChapters)
	demandId := 0
	for _, demandChapter := range demandChapters {
		demandId = demandChapter.DemandId
		result = append(result, demandChapter.DemandChapterId)
	}
	return demandId, result
}

func (this *KnowledgeScenarioChapter) GetDemandChapters(scenarioId, demandId int, demandNames []string) []int {
	result := make([]int, 0)
	demandChapters := []KnowledgeScenarioChapter{}
	session := mysql.GetSession()
	session.Where("scenario_id=?", scenarioId)
	session.Where("demand_id=?", scenarioId)
	session.Find(&demandChapters)
	for _, demandChapter := range demandChapters {
		result = append(result, demandChapter.DemandChapterId)
	}
	return result
}

func (this *KnowledgeScenarioChapter) GetDemandChapterIds(scenarioId int) []int {
	chapterIds := make([]int, 0)
	if err := sys_service.NewSessionWithCond(qmap.QM{"e_scenario_id": scenarioId}).Cols("demand_chapter_id").Table(this).Find(&chapterIds); err == nil {
		return chapterIds
	} else {
		panic(err)
	}
}
