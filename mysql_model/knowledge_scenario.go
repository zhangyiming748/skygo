package mysql_model

import (
	"fmt"

	"skygo_detection/lib/common_lib/log"
	"skygo_detection/lib/common_lib/mysql"
)

type KnowledgeScenario struct {
	Id           int    `xorm:"not null pk autoincr comment('自增主键id') INT(11)" json:"id"`
	Name         string `xorm:"not null comment('安全检测场景名称') VARCHAR(255)" json:"name"`
	DemandId     int    `xorm:"not null comment('关联知识库需求id') INT(11)" json:"demand_id"`
	Detail       string `xorm:"not null default '' comment('场景描述') VARCHAR(255)" json:"detail"`
	CreateTime   int    `xorm:"not null comment('创建时间') INT(11)" json:"create_time"`
	UpdateTime   int    `xorm:"not null comment('更新时间') INT(11)" json:"update_time"`
	CreateUserId int    `xorm:"not null comment('创建人') INT(11)" json:"create_user_id"`
	LastOpId     int    `xorm:"not null comment('最后更新人') INT(11)" json:"last_op_id"`
	Tag          string `xorm:"not null comment('标签') VARCHAR(255)" json:"tag"`
	Describe     string `xorm:"not null default '' comment('介绍') VARCHAR(255)" json:"describe"`
}

// 按照id查询
func KnowledgeScenarioFindById(id int) (*KnowledgeScenario, bool, error) {
	model := KnowledgeScenario{}
	if has, err := mysql.GetSession().ID(id).Get(&model); err != nil {
		return nil, has, err
	} else {
		return &model, has, nil
	}
}

// 按照name 查询
func (this *KnowledgeScenario) KnowledgeScenarioFindByName(name string) (*KnowledgeScenario, bool) {
	model := KnowledgeScenario{}
	has, _ := mysql.GetSession().Where("name=?", name).Get(&model)
	return &model, has
}

func (this *KnowledgeScenario) Create(demandChapterIds []int) (*KnowledgeScenario, error) {
	//创建场景数据
	session := mysql.GetSession()
	_, err := session.InsertOne(this)
	if err != nil {
		return nil, err
	}
	//创建场景和关联安全需求条目数据
	for _, demandChapterId := range demandChapterIds {
		scenarioChapter := new(KnowledgeScenarioChapter)
		scenarioChapter.DemandId = this.DemandId
		scenarioChapter.DemandChapterId = demandChapterId
		scenarioChapter.ScenarioId = this.Id
		_, err = scenarioChapter.Create()
		if err != nil {
			return nil, err
		}
	}
	return this, err
}

func (this *KnowledgeScenario) Update(demandChapterIds []interface{}, cols ...string) (*KnowledgeScenario, error) {
	//更新场景库
	_, err := mysql.GetSession().Table(this).ID(this.Id).Cols(cols...).Update(this)
	if err != nil {
		log.GetHttpLogLogger().Error(fmt.Sprintf("%v", err))
		return nil, err
	}
	//更新场景 和 需求的关系库
	scenarioChapter := new(KnowledgeScenarioChapter)
	scenarioChapter.ScenarioId = this.Id
	scenarioChapter.DemandId = this.DemandId
	//先删除关系库的关系
	_, err = mysql.GetSession().Delete(scenarioChapter)
	if err != nil {
		log.GetHttpLogLogger().Error(fmt.Sprintf("%v", err))
		return nil, err
	}
	//增加关系
	for _, demandChapterId := range demandChapterIds {
		tmp := demandChapterId.(float64)
		scenarioChapter.DemandChapterId = int(tmp)
		_, err = scenarioChapter.Create()
		if err != nil {
			return nil, err
		}
	}
	return this, err
}

func (this *KnowledgeScenario) Remove() (*KnowledgeScenario, error) {
	//删除关系表的库
	scenarioChapter := new(KnowledgeScenarioChapter)
	scenarioChapter.ScenarioId = this.Id
	_, err := scenarioChapter.Remove()
	if err != nil {
		return nil, err
	}
	//删除场景库
	_, err = mysql.GetSession().Delete(this)
	if err != nil {
		return nil, err
	}
	return this, nil
}
func (this *KnowledgeScenario) UpdateTag(id int, tag string) (int64, error) {
	this.Id = id
	this.Tag = tag
	//has , err := mysql.GetSession().Get(this)
	return mysql.GetSession().ID(this.Id).Cols("tag").Update(this)
}
