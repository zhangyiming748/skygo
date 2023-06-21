package mysql_model

type FirmwareReportRtsLeaks struct {
	Id         int    `xorm:"not null pk autoincr INT(10)"`
	ScannerId  int    `xorm:"comment('固件扫描任务id') INT(10)"`
	Info       string `xorm:"not null VARCHAR(512)"`
	Type       string `xorm:"VARCHAR(255)"`
	Path       string `xorm:"VARCHAR(512)"`
	FullPath   string `xorm:"VARCHAR(1024)"`
	Origin     string `xorm:"VARCHAR(512)"`
	CreateTime int    `xorm:"INT(10)"`
}
