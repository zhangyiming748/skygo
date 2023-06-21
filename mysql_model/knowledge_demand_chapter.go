package mysql_model

import (
	"errors"

	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/mysql"
)

// 知识库需求的某个章节
type KnowledgeDemandChapter struct {
	Id                int    `xorm:"not null pk autoincr comment('自增主键id') INT(11)" json:"id"`
	KnowledgeDemandId int    `xorm:"not null comment('表knowledge_demand主键') INT(11)" json:"knowledge_demand_id"`
	Code              string `xorm:"not null comment('章节编号') VARCHAR(50)" json:"code"`
	Title             string `xorm:"not null comment('章节标题') VARCHAR(255)" json:"title"`
	ParentId          int    `xorm:"not null default -1 comment('父章节id，-1标识无父章节') INT(11)" json:"parent_id"`
	ParentCode        string `xorm:"not null comment('父章节编号') VARCHAR(255)" json:"parent_code"`
	Content           string `xorm:"not null default '' comment('内容') VARCHAR(255)" json:"content"`
}

// 查询id查询当前章节是否有子章节
func (this *KnowledgeDemandChapter) ExistChildChapter() (bool, error) {
	id := this.Id
	return mysql.GetSession().Where("parent_id = ?", id).Exist(&KnowledgeDemandChapter{})
}

// 创建表单
type KnowledgeDemandChapterCreateForm struct {
	KnowledgeDemandId int    `json:"knowledge_demand_id"` // 表knowledge_demand主键
	Code              string `json:"code"`                // 章节编号
	Title             string `json:"tile"`                // 章节标题
	ParentId          int    `json:"parent_id"`           // 父章节编号
	Content           string `json:"content"`             // 内容
}

// 根据表单创建测试件记录， 返回新建记录id
func KnowledgeDemandChapterCreate(form *KnowledgeDemandChapterCreateForm) (*KnowledgeDemandChapter, error) {
	// 判断存在
	if has, _ := KnowLedgeDemandExistById(form.KnowledgeDemandId); !has {
		return nil, errors.New("需求不存在")
	}

	// parentId为0，代表没有父节点
	var parentModel *KnowledgeDemandChapter
	var has bool
	if form.ParentId != 0 {
		parentModel, has = KnowledgeDemandChapterFindById(form.ParentId)
		if !has {
			return nil, errors.New("父章节不存在")
		}
	}

	// 创建测试件主记录
	model := KnowledgeDemandChapter{}
	model.KnowledgeDemandId = form.KnowledgeDemandId
	model.Code = form.Code
	model.Title = form.Title
	if form.ParentId != 0 {
		model.ParentId = form.ParentId
		model.ParentCode = parentModel.ParentCode
	} else {
		model.ParentId = 0
		model.ParentCode = ""
	}
	model.Content = form.Content

	_, err := mysql.GetSession().Insert(&model)
	return &model, err
}

// 按照id查询某个章节记录
func KnowledgeDemandChapterFindById(id int) (*KnowledgeDemandChapter, bool) {
	model := KnowledgeDemandChapter{}
	has, err := mysql.FindById(id, &model)
	if err != nil {
		panic(err)
	}
	return &model, has
}

// 按照id删除某个章节记录
func KnowledgeDemandChapterDeleteById(id int) (int64, error) {
	chapterModel, has := KnowledgeDemandChapterFindById(id)
	if !has {
		return 0, errors.New("记录不存在")
	}

	if has, _ := chapterModel.ExistChildChapter(); has {
		return 0, errors.New("存在子章节")
	}

	return mysql.GetSession().ID(id).Delete(chapterModel)
}

// 按照id更新
func KnowledgeDemandChapterUpdateById(id int, data qmap.QM) (*KnowledgeDemandChapter, error) {
	modelPtr, has := KnowledgeDemandChapterFindById(id)
	if !has {
		return nil, errors.New("数据不存在")
	}

	// 章节编号
	if code, has := data.TryString("code"); has {
		modelPtr.Code = code
	}
	// 章节标题
	if title, has := data.TryString("title"); has {
		modelPtr.Title = title
	}
	// 父章节编号
	if parentId, has := data.TryInt("parent_id"); has {
		// todo 判断合法性
		modelPtr.ParentId = parentId
	}
	// 内容
	if content, has := data.TryString("content"); has {
		modelPtr.Content = content
	}

	if _, err := mysql.GetSession().ID(id).Update(modelPtr); err != nil {
		return nil, err
	} else {
		return modelPtr, nil
	}
}

// ------------------------------------------------------------------------
// 页面级联列表树状展示查询
type KnowLedgeDemandTreeNode struct {
	Id       int                        `json:"id"`       // 主键id
	Code     string                     `json:"code"`     // 章节编号
	Title    string                     `json:"title"`    // 章节标题
	Children []*KnowLedgeDemandTreeNode `json:"children"` // 子节点
}

// 节点递归方式获取子节点
func (this *KnowLedgeDemandTreeNode) FetchChild(chapterIds []int) {
	lists := make([]KnowledgeDemandChapter, 0)
	session := mysql.GetSession().Where("parent_id = ?", this.Id)
	if len(chapterIds) > 0 {
		session.In("id", chapterIds)
	}
	session.Find(&lists)
	if len(lists) == 0 {
		return
	} else {
		children := make([]*KnowLedgeDemandTreeNode, 0)
		for _, v := range lists {
			children = append(children, &KnowLedgeDemandTreeNode{
				Id:       v.Id,
				Code:     v.Code,
				Title:    v.Title,
				Children: make([]*KnowLedgeDemandTreeNode, 0),
			})
		}
		this.Children = children
	}

	if len(this.Children) > 0 {
		for _, v := range this.Children {
			v.FetchChild(chapterIds)
		}
	}
	return
}

// -------------------------------------------------------------------------------
// 下拉列表
type KnowLedgeDemandChapterSelectListType struct {
	Id    int    `json:"id"`
	Code  string `json:"code"`
	Title string `json:"title"`
}

// 法规列表
type KnowLedgeDemandCodeList struct {
	Id   int    `json:"id"`
	Code string `json:"code"`
}

func KnowLedgeDemandChapterSelectList(demandId int) []KnowLedgeDemandChapterSelectListType {
	models := make([]KnowLedgeDemandChapterSelectListType, 0)
	mysql.GetSession().Select("id, code, title").Table(KnowledgeDemandChapter{}).
		Where("knowledge_demand_id = ?", demandId).
		OrderBy("id asc").
		Limit(10000).
		Find(&models)
	return models
}

// 按照name 查询
func (this *KnowledgeDemandChapter) KnowledgeDemandChapterFindByName(demandId int, code, title string) (*KnowledgeDemandChapter, bool) {
	model := KnowledgeDemandChapter{}
	session := mysql.GetSession()
	session.Where("knowledge_demand_id=?", demandId)
	session.Where("code=?", code)
	session.Where("title=?", title)
	has, _ := session.Get(&model)
	return &model, has
}

func GetKnowLedgeDemandCodeList() []string {
	models := make([]KnowLedgeDemandCodeList, 0)
	mysql.GetSession().Distinct("code").Table(KnowledgeDemand{}).
		Limit(10000).
		Find(&models)
	var result = make([]string, 0)
	for _, model := range models {
		result = append(result, model.Code)
	}
	return result
}

func GetAllChild(kdtn *KnowLedgeDemandTreeNode) (allId []int) {
	for _, v := range kdtn.Children {
		allId = append(allId, v.Id)
		if len(v.Children) > 0 {
			tmp := GetAllChild(v)
			allId = append(allId, tmp...)
		}
	}
	return
}
