package mysql_model

type FirmwareReportInit struct {
	Id          int    `xorm:"not null pk autoincr INT(10)"`
	ScannerId   int    `xorm:"comment('固件扫描任务id') INT(10)"`
	DirNum      int    `xorm:"INT(10)"`
	FileNum     int    `xorm:"INT(10)"`
	LinkNum     int    `xorm:"INT(10)"`
	NodeNum     int    `xorm:"INT(10)"`
	Arch        string `xorm:"MEDIUMTEXT"`
	System      string `xorm:"MEDIUMTEXT"`
	Compressed  string `xorm:"MEDIUMTEXT"`
	DiskUsage   string `xorm:"MEDIUMTEXT"`
	Dulplicated string `xorm:"MEDIUMTEXT"`
	Filesystem  string `xorm:"MEDIUMTEXT"`
	Firmware    string `xorm:"MEDIUMTEXT"`
	SystemGuess string `xorm:"MEDIUMTEXT"`
	TotalSize   string `xorm:"VARCHAR(255)"`
	CreateTime  int    `xorm:"INT(10)"`
}
