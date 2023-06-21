package mysql_model

type MicroService struct {
	Id          int    `xorm:"not null pk autoincr comment('id') INT(10)"`
	Service     string `xorm:"not null comment('服务') VARCHAR(255)"`
	ServiceName string `xorm:"comment('服务名称') VARCHAR(255)"`
	CreateTime  int    `xorm:"created not null default 0 comment('创建时间') INT(10)"`
}
