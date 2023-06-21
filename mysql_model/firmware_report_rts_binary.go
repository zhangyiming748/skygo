package mysql_model

type FirmwareReportRtsBinary struct {
	Id         int    `xorm:"not null pk autoincr INT(10)"`
	ScannerId  int    `xorm:"comment('固件扫描任务id') INT(10)"`
	IsDoubt    int    `xorm:"not null TINYINT(11)"`
	RelaPath   string `xorm:"not null VARBINARY(255)"`
	Hardenable int    `xorm:"not null TINYINT(11)"`
	MagicInfo  string `xorm:"TEXT"`
	Type       string `xorm:"VARCHAR(255)"`
	IsElf      int    `xorm:"TINYINT(255)"`
	Result     string `xorm:"TEXT"`
	FileName   string `xorm:"VARCHAR(255)"`
	FullPath   string `xorm:"VARCHAR(255)"`
	CreateTime int    `xorm:"INT(10)"`
}
