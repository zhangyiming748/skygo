package mysql_model

type SysRole struct {
	Id   int    `xorm:"not null pk autoincr comment('编号') INT(11)"`
	Name string `xorm:"not null comment('角色') index VARCHAR(64)"`
}
