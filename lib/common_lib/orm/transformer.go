package orm

import (
	"skygo_detection/guardian/src/net/qmap"
)

// 对map进行转换
type TransformerFunc func(qm qmap.QM) qmap.QM

// 接口定义
type TransformerIf interface {
	ModifyItem(qm qmap.QM) qmap.QM // 修改一个map，如添加其字段、删除字段等
	ExcludeFields() []string       // 删除map中指定的字段
}

// 用户根据业务自定义的struct包含此结构体即可实现TransformerIf接口
// 在通过重写的ModifyTtem、ExcludeFields函数来实现具体逻辑
type Transformer struct{}

func (t Transformer) ModifyItem(qm qmap.QM) qmap.QM {
	return qm
}

func (t Transformer) ExcludeFields() []string {
	return nil
}
