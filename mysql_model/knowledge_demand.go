package mysql_model

import (
	"errors"
	"time"

	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/mysql"
)

// 检测知识库-安全需求
type KnowledgeDemand struct {
	Id            int    `xorm:"not null pk autoincr comment('自增主键id') INT(11)" json:"id"`
	Name          string `xorm:"not null comment('安全需求名称，非空') VARCHAR(255)" json:"name"`
	Category      int    `xorm:"not null comment('需求类型，1企业内部标准2法规标准3渗透测试4其他') TINYINT(3)" json:"category"`
	Code          string `xorm:"not null comment('标准编号') VARCHAR(255)" json:"code"`
	ImplementTime int    `xorm:"not null comment('实施日期（秒）') INT(11)" json:"implement_time"`
	Detail        string `xorm:"not null default '' comment('描述') VARCHAR(255)" json:"detail"`
	CreateTime    int    `xorm:"not null comment('创建日期（秒）') INT(11)" json:"create_time"`
	UpdateTime    int    `xorm:"not null default 0 comment('更新日期（秒）') INT(11)" json:"update_time"`
	CreateUserId  int    `xorm:"not null comment('创建用户id') INT(11)" json:"create_user_id"`
}

func (this *KnowledgeDemand) Create() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}

func (this *KnowledgeDemand) Update(cols ...string) (int64, error) {
	return mysql.GetSession().Table(this).ID(this.Id).Cols(cols...).Update(this)
}

func (this *KnowledgeDemand) Remove() (int64, error) {
	return mysql.GetSession().ID(this.Id).Delete(this)
}

// 按照id查询某个章节记录
func KnowledgeDemandFindById(id int) (*KnowledgeDemand, bool) {
	model := KnowledgeDemand{}
	has, err := mysql.FindById(id, &model)
	if err != nil {
		panic(err)
	}
	return &model, has
}

// 查询当前Id的记录是否存在
func KnowLedgeDemandExistById(id int) (bool, error) {
	return mysql.GetSession().Where("id = ?", id).Exist(&KnowledgeDemand{})
}

// 查询当前名称的记录是否存在
func KnowLedgeDemandExistByName(name string) (bool, error) {
	return mysql.GetSession().Where("name = ?", name).Exist(&KnowledgeDemand{})
}

// ------------------------------------------------------------------------
// 创建表单
type KnowledgeDemandCreateForm struct {
	Name          string `json:"name"`           // 安全需求名称
	Category      int    `json:"category"`       // 类型
	Code          string `json:"code"`           // 标准编号
	ImplementTime int    `json:"implement_time"` // 实施日期
	Detail        string `json:"detail"`         // 测试件描述
}

// 根据表单创建测试件记录， 返回新建记录id
func KnowledgeDemandCreateFromForm(form *KnowledgeDemandCreateForm, uid int) (*KnowledgeDemand, error) {
	// 判断存在
	if has, _ := KnowLedgeDemandExistByName(form.Name); has {
		return nil, errors.New("名称已经存在")
	}

	// 创建测试件主记录
	model := KnowledgeDemand{}
	model.Name = form.Name
	// 类型
	// 1 企业内部标准
	// 2 法规标准
	// 3 渗透测试
	// 4 其他
	model.Category = form.Category
	model.Code = form.Code
	model.Detail = form.Detail
	model.ImplementTime = form.ImplementTime
	model.CreateTime = int(time.Now().Unix())
	model.UpdateTime = int(time.Now().Unix())
	model.CreateUserId = uid

	_, err := mysql.GetSession().Insert(&model)
	return &model, err
}

// 按照id更新
func KnowledgeDemandUpdateById(id int, data qmap.QM) (*KnowledgeDemand, error) {
	model := KnowledgeDemand{}
	has, err := mysql.GetSession().Where("id=?", id).Get(&model)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("数据不存在")
	}

	if name, has := data.TryString("name"); has {
		model.Name = name
	}

	if category, has := data.TryInt("category"); has {
		model.Category = category
	}

	if code, has := data.TryString("code"); has {
		model.Code = code
	}

	if detail, has := data.TryString("detail"); has {
		model.Detail = detail
	}

	if implementTime, has := data.TryInt("implement_time"); has {
		model.ImplementTime = implementTime
	}
	model.UpdateTime = int(time.Now().Unix())

	if _, err := mysql.GetSession().ID(id).Update(&model); err != nil {
		return nil, err
	} else {
		return &model, nil
	}
}

// 根据需求id，删除单条记录，注意级联删除
func KnowledgeDemandDeleteById(id int) error {
	// 启动事务
	session := mysql.GetSession()
	err := session.Begin()
	if err != nil {
		return err
	}
	defer session.Close()

	// 删除需求记录, 通过返回的int值判断记录是否存在
	if count, err := session.ID(id).Delete(KnowledgeDemand{}); err != nil {
		session.Rollback()
		return err
	} else {
		if count == 0 {
			session.Rollback()
			return errors.New("记录不存在")
		}
	}

	// 删除需求记录, 通过返回的int值判断记录是否存在
	if count, err := session.Where("knowledge_demand_id = ?", id).Delete(KnowledgeDemandChapter{}); err != nil {
		session.Rollback()
		return err
	} else {
		if count == 0 {
			session.Rollback()
			return errors.New("记录不存在")
		}
	}

	if err := session.Commit(); err != nil {
		return err
	}

	return nil
}

func NewKnowledgeDemandDeleteById(id int) error {
	// 启动事务
	session := mysql.GetSession()
	// 删除需求记录, 通过返回的int值判断记录是否存在
	_, err := session.ID(id).Delete(KnowledgeDemand{})
	if err != nil {
		return err
	}
	session = mysql.GetSession()
	_, err = session.Where("knowledge_demand_id = ?", id).Delete(KnowledgeDemandChapter{})
	if err != nil {
		return err
	}
	return nil
}

// 知识库需求查询
func KnowledgeDemandTree(id int, chapterIds []int) []*KnowLedgeDemandTreeNode {
	lists := make([]KnowledgeDemandChapter, 0)
	session := mysql.GetSession().Where("knowledge_demand_id = ?", id).And("parent_id = 0")
	if len(chapterIds) > 0 {
		session.In("id", chapterIds)
	}
	session.Find(&lists)

	result := make([]*KnowLedgeDemandTreeNode, 0)

	if len(lists) > 0 {
		for _, v := range lists {
			// 基于章节记录，构建出一个节点对象
			node := &KnowLedgeDemandTreeNode{
				Id:       v.Id,
				Code:     v.Code,
				Title:    v.Title,
				Children: make([]*KnowLedgeDemandTreeNode, 0),
			}
			// 节点对象开始获取子节点，递归方式
			node.FetchChild(chapterIds)

			result = append(result, node)
		}
	}
	return result
}

// 知识库需求查询
func KnowledgeScenarioTree(demandId int, chapterIds []int) []*KnowLedgeDemandTreeNode {
	chapterIds = GetChapterParentId(chapterIds)
	lists := make([]KnowledgeDemandChapter, 0)
	session := mysql.GetSession().Where("knowledge_demand_id = ?", demandId).And("parent_id = 0")
	if len(chapterIds) > 0 {
		session.In("id", chapterIds)
	}
	session.Find(&lists)

	result := make([]*KnowLedgeDemandTreeNode, 0)

	if len(lists) > 0 {
		for _, v := range lists {
			// 基于章节记录，构建出一个节点对象
			node := &KnowLedgeDemandTreeNode{
				Id:       v.Id,
				Code:     v.Code,
				Title:    v.Title,
				Children: make([]*KnowLedgeDemandTreeNode, 0),
			}
			// 节点对象开始获取子节点，递归方式
			node.FetchChild(chapterIds)

			result = append(result, node)
		}
	}
	return result
}

func GetChapterParentId(chapterIds []int) []int {
	result := chapterIds
	for {
		parentIds := []int{}
		if err := sys_service.NewSession().Session.Table(new(KnowledgeDemandChapter)).Cols("parent_id").In("id", chapterIds).And("parent_id<>?", 0).Find(&parentIds); err == nil {
			if len(parentIds) > 0 {
				result = append(result, parentIds...)
				chapterIds = parentIds
			} else {
				break
			}
		} else {
			panic(err)
		}
	}
	return result
}

// -------------------------------------------------------------------------------
// 下拉列表
type KnowLedgeDemandSelectListType struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func KnowLedgeDemandSelectList() []KnowLedgeDemandSelectListType {
	models := make([]KnowLedgeDemandSelectListType, 0)
	mysql.GetSession().Select("id, name").Table(KnowledgeDemand{}).
		OrderBy("id asc").
		Limit(10000).
		Find(&models)
	return models
}

// 按照name 查询
func (this *KnowledgeDemand) KnowledgeDemandFindByName(name string) (*KnowledgeDemand, bool) {
	model := KnowledgeDemand{}
	has, _ := mysql.GetSession().Where("name=?", name).Get(&model)
	return &model, has
}

// 按照id查询
func (this *KnowledgeDemand) KnowledgeDemandFindById(id int) (*KnowledgeDemand, bool) {
	if has, err := mysql.GetSession().ID(id).Get(this); err != nil {
		return nil, has
	} else {
		return this, has
	}
}
