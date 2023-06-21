package mysql_model

import "skygo_detection/lib/common_lib/mysql"

type KnowledgeTestCaseChapter struct {
	Id              int `xorm:"not null pk comment('自增主键id') INT(11)"`
	TestCaseId      int `xorm:"not null comment('测试用例id') INT(11)"`
	SenarioId       int `xorm:"not null comment('安全检测场景id') INT(11)"`
	DemandId        int `xorm:"not null comment('关联知识库需求id') INT(11)"`
	DemandChapterId int `xorm:"not null comment('关联知识库需求章节id') INT(11)"`
}

func (this *KnowledgeTestCaseChapter) Create() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}

func (this *KnowledgeTestCaseChapter) Update(cols ...string) (int64, error) {
	return mysql.GetSession().Table(this).ID(this.Id).Cols(cols...).Update(this)
}

func (this *KnowledgeTestCaseChapter) Remove() (int64, error) {
	return mysql.GetSession().Delete(this)
}

func (this *KnowledgeTestCaseChapter) GetByTestCaseId(testCaseId int) ([]int, []string) {
	resultInd := make([]int, 0)
	resultString := make([]string, 0)
	testCaseChapters := []KnowledgeTestCaseChapter{}
	session := mysql.GetSession()
	session.Where("test_case_id=?", testCaseId)
	session.Find(&testCaseChapters)
	for _, testCaseChapter := range testCaseChapters {
		demandChapterId := testCaseChapter.DemandChapterId
		resultInd = append(resultInd, demandChapterId)
		tmp, _ := KnowledgeDemandChapterFindById(demandChapterId)
		resultString = append(resultString, tmp.Title)
	}
	return resultInd, resultString
}

// 查询测试用例
func (this *KnowledgeTestCaseChapter) GetTestCases(demandId int, chapterIds []interface{}) []KnowledgeTestCaseChapter {
	model := make([]KnowledgeTestCaseChapter, 0)
	session := mysql.GetSession()
	session.Where("demand_id=?", demandId)
	session.In("demand_chapter_id", chapterIds)
	session.Find(&model)
	return model
}

func (this *KnowledgeTestCaseChapter) GetByDemandChapterIds(demandChapterIds []int) []KnowledgeTestCaseChapter {
	model := make([]KnowledgeTestCaseChapter, 0)
	session := mysql.GetSession()
	session.In("demand_chapter_id", demandChapterIds)
	session.Find(&model)
	return model
}

func (this *KnowledgeTestCaseChapter) GetDemandChapterIdsBySenarioId(senarioId int) []KnowledgeTestCaseChapter {
	model := make([]KnowledgeTestCaseChapter, 0)
	session := mysql.GetSession()
	session.Where("senario_id=?", senarioId)
	session.Find(&model)
	return model
}
