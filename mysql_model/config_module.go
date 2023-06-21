package mysql_model

type ConfigModule struct {
	Id             int    `xorm:"not null pk comment('自增主键id') INT(11)"`
	ModuleType     string `xorm:"not null comment('组件分类，如蓝牙钥匙等') VARCHAR(255)"`
	ModuleTypeCode int    `xorm:"not null comment('组件分类编码') INT(11)"`
	ModeleName     string `xorm:"not null comment('组件名称') VARCHAR(255)"`
	ModuleNameCode int    `xorm:"not null comment('组件名称编码') INT(11)"`
}
