package mysql_model

type ToolHistory struct {
	Id            int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	Vsersion      int    `xorm:"not null comment('版本') INT(11)" json:"version"`
	Date          string `xorm:"not null default '' comment('时间戳') VARCHAR(255)" json:"data"`
	Updateser     int    `xorm:"not null pk autoincr comment('修改人id') INT(11)" json:"updateuser"`
	Updatecontent string `xorm:"not null default '' comment('修改内容') VARCHAR(255)" json:"updatecontent"`
}
