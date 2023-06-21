package mysql_model

import (
	"skygo_detection/lib/common_lib/mysql"
)

type TestCaseModule struct {
	Id             string `xorm:"not null pk autoincr comment('id') VARCHAR(255)"`
	ModuleName     string `xorm:"comment('测试组件') VARCHAR(255)"`       // 测试组件
	ModuleNameCode string `xorm:"not null comment('测试组件编码') INT(11)"` // 测试组件编码
	ModuleType     string `xorm:"comment('测试分类') VARCHAR(255)"`       // 测试分类
	ModuleTypeCode string `xorm:"not null comment('测试组件编码') INT(11)"` // 测试组件编码
}

func (this *TestCaseModule) Create() (*TestCaseModule, error) {
	// 创建场景数据
	session := mysql.GetSession()
	_, err := session.InsertOne(this)
	if err != nil {
		return nil, err
	}
	return this, err
}

func (this *TestCaseModule) Update(id interface{}, cols ...string) (*TestCaseModule, error) {
	_, err := mysql.GetSession().Table(this).ID(id).Cols(cols...).Update(this)
	return this, err
}

func (this *TestCaseModule) Remove() (*TestCaseModule, error) {
	// 删除场景库
	_, err := mysql.GetSession().Delete(this)
	if err != nil {
		return nil, err
	}
	return this, nil
}

func (this *TestCaseModule) RemoveById(id int) (int64, error) {
	return mysql.GetSession().ID(id).Delete(this)
}

func (this *TestCaseModule) GetOne(id int) (*TestCaseModule, error) {
	session := mysql.GetSession()
	session.Where("id=?", id)
	_, err := session.Get(this)
	if err != nil {
		return nil, err
	}
	return this, err
}

func (this *TestCaseModule) GetBydName(moduleName, moduleType string) (*TestCaseModule, error) {
	session := mysql.GetSession()
	session.Where("module_name", moduleName)
	session.Where("module_type", moduleType)
	_, err := session.Get(this)
	if err != nil {
		return nil, err
	}
	return this, err
}
