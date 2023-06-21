package mysql_model

import (
	"errors"
	"fmt"

	"skygo_detection/guardian/src/net/qmap"
	"xorm.io/xorm"

	"skygo_detection/lib/common_lib/mysql"
)

type SysModule struct {
	Id          int    `xorm:"not null pk autoincr comment('模块序号') SMALLINT(11)" json:"id"`
	Name        string `xorm:"not null comment('模块名称') VARCHAR(64)" json:"name"`
	Rank        int    `xorm:"not null comment('排序') index SMALLINT(5)" json:"rank"`
	IconName    string `xorm:"comment('图标名') VARCHAR(32)" json:"icon_name"`
	Enable      int    `xorm:"not null default 0 comment('禁用状态(0:正常 1:禁用)') TINYINT(1)" json:"enable"`
	ParentId    int    `xorm:"not null default 0 comment('父ID') SMALLINT(11)" json:"parent_id"`
	ForeignLink string `xorm:"comment('外部链接') VARCHAR(128)" json:"foreign_link"`
	IsMenu      int    `xorm:"not null default 0 comment('是否为菜单(0:权限 1:菜单)') TINYINT(1)" json:"is_menu"`
}

func SysModuleGetWhere(params map[string]interface{}) *xorm.Session {
	session := mysql.GetSession()
	if id, ok := params["id"]; ok {
		session.And("id=?", id)
	}

	return session
}

func (this *SysModule) Delete(id int) (int64, error) {
	params := qmap.QM{
		"id": id,
	}
	return SysModuleGetWhere(params).Delete(this)
}

func (this *SysModule) Create() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}

func (this *SysModule) FindById(id int) (*SysModule, bool) {
	if has, err := mysql.GetSession().ID(id).Get(this); err != nil {
		panic(err)
	} else {
		return this, has
	}
}

func (this *SysModule) UpdateById(id int, data qmap.QM) (*SysModule, error) {
	if _, has := this.FindById(id); has {
		if _, err := mysql.GetSession().Table(this).ID(this.Id).Update(data); err != nil {
			return nil, err
		} else {
			newModule := new(SysModule)
			newModule.FindById(this.Id)
			return newModule, nil
		}
	} else {
		return nil, errors.New(fmt.Sprintf("module id:%v is not found", id))
	}
}

func (this *SysModule) GetModuleTree() []qmap.QM {
	modules := []SysModule{}
	if err := mysql.GetSession().Find(&modules); err == nil {
		return this.GetMenusTree(modules, 3)
	} else {
		panic(err)
	}
}

func (this *SysModule) GetAllModules() []qmap.QM {
	res := []SysModule{}
	if err := mysql.GetSession().Find(&res); err == nil {
		modules := []qmap.QM{}
		for _, module := range res {
			temp := qmap.QM{
				"id":           module.Id,
				"name":         module.Name,
				"rank":         module.Rank,
				"icon_name":    module.IconName,
				"foreign_link": module.ForeignLink,
			}
			modules = append(modules, temp)
		}
		return modules
	} else {
		panic(err)
	}
}

func (this *SysModule) GetMenusByRoleID(roleId int) []qmap.QM {
	modules := this.GetModulesByRoleID(roleId)
	return this.GetMenusTree(modules, 3)
}

func (this *SysModule) GetModulesByRoleID(roleId int) []SysModule {
	moduleIds := new(SysRoleModule).GetModuleIds(roleId)
	modules := []SysModule{}
	if err := mysql.GetSession().Table(this).In("id", moduleIds).Find(&modules); err == nil {
		return modules
	} else {
		panic(err)
	}
}

/*
*
获取菜单树
*/
func (this *SysModule) GetMenusTree(modules []SysModule, level int) []qmap.QM {
	menusTree := []qmap.QM{}
	for _, module := range modules {
		if module.ParentId <= 0 {
			menu := qmap.QM{
				"id":           module.Id,
				"name":         module.Name,
				"rank":         module.Rank,
				"icon_name":    module.IconName,
				"foreign_link": module.ForeignLink,
				"children":     this.GetModuleChildren(module.Id, modules, level-1),
			}
			menusTree = append(menusTree, menu)
		}
	}
	return menusTree
}

func (this *SysModule) GetModuleChildren(parentId int, modules []SysModule, level int) []qmap.QM {
	if level <= 0 {
		return []qmap.QM{}
	}
	children := []qmap.QM{}
	for _, module := range modules {
		if module.ParentId == parentId {
			menu := qmap.QM{
				"id":           module.Id,
				"name":         module.Name,
				"rank":         module.Rank,
				"icon_name":    module.IconName,
				"foreign_link": module.ForeignLink,
				"children":     this.GetModuleChildren(module.Id, modules, level-1),
			}
			children = append(children, menu)
		}
	}
	return children
}
