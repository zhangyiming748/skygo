package mysql_model

type KnowledgeTestTools struct {
	Id           int    `xorm:"not null pk autoincr comment('主键id') INT(11)"`
	Name         string `xorm:"not null comment('工具名称') VARCHAR(255)"`
	Category     string `xorm:"not null comment('工具分类') VARCHAR(255)"`
	Introduce    string `xorm:"not null comment('工具介绍') VARCHAR(255)"`
	Brand        string `xorm:"not null comment('品牌') VARCHAR(255)"`
	Version      string `xorm:"not null comment('规格型号/版本') VARCHAR(255)"`
	CreateTime   int    `xorm:"not null comment('创建时间') INT(11)"`
	UpdateTime   int    `xorm:"not null comment('更新时间') INT(11)"`
	LastOpId     int    `xorm:"not null comment('最近操作用户id') INT(11)"`
	CreateUserId int    `xorm:"not null comment('创建用户id') INT(11)"`
}
