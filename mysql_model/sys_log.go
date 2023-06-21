package mysql_model

type SysLog struct {
	Id        int    `xorm:"not null pk autoincr comment('主键,自增长') INT(10)"`
	ChannelId string `xorm:"not null default '' comment('渠道号') CHAR(6)"`
	Uid       string `xorm:"comment('用户id（与sys_user用户表的ID相关联）') CHAR(32)"`
	Time      int    `xorm:"not null comment('时间') INT(10)"`
	Type      string `xorm:"default 'info' comment('日志类型('debug', 'info','warn', 'error')') VARCHAR(16)"`
	Url       string `xorm:"comment('访问URL') VARCHAR(128)"`
	Content   string `xorm:"comment('日志内容') LONGTEXT"`
}
