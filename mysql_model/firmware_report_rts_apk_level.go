package mysql_model

type FirmwareReportRtsApkLevel struct {
	Id              int    `xorm:"not null pk autoincr INT(10)"`
	ScannerId       int    `xorm:"comment('固件扫描任务id') INT(10)"`
	Type            string `xorm:"VARCHAR(255)"`
	No              int    `xorm:"not null INT(10)"`
	Low             int    `xorm:"not null default 0 INT(10)"`
	Mid             int    `xorm:"not null default 0 INT(10)"`
	High            int    `xorm:"not null default 0 INT(10)"`
	Heavy           int    `xorm:"not null default 0 INT(10)"`
	Total           int    `xorm:"not null default 0 INT(10)"`
	So              string `xorm:"TEXT"`
	SensitiveCount  string `xorm:"TEXT"`
	Level           int    `xorm:"not null default 0 INT(10)"`
	OriginalContent string `xorm:"not null TEXT"`
	CreateTime      int    `xorm:"not null INT(10)"`
}
