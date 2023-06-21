package mysql_model

type FirmwareReportRtsCert struct {
	Id         int    `xorm:"not null pk autoincr INT(10)"`
	ScannerId  int    `xorm:"comment('固件扫描任务id') INT(10)"`
	Type       string `xorm:"VARCHAR(255)"`
	FileName   string `xorm:"not null VARCHAR(512)"`
	Info       string `xorm:"VARCHAR(1024)"`
	Path       string `xorm:"VARCHAR(1024)"`
	Content    string `xorm:"VARCHAR(512)"`
	CreateTime int    `xorm:"INT(10)"`
}
