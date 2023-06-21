package mysql_model

type FirmwareReportRtsCve struct {
	Id                    int     `xorm:"not null pk autoincr INT(10)"`
	ScannerId             int     `xorm:"comment('固件扫描任务id') INT(10)"`
	Level                 int     `xorm:"TINYINT(255)"`
	Description           string  `xorm:"not null TEXT"`
	VersionEndIncluding   string  `xorm:"VARCHAR(1024)"`
	Cvssv3                string  `xorm:"TEXT"`
	Vendor                string  `xorm:"VARCHAR(255)"`
	Cvssv2score           float32 `xorm:"FLOAT(255)"`
	Version               string  `xorm:"VARCHAR(255)"`
	Type                  string  `xorm:"VARCHAR(255)"`
	Cvssv2                string  `xorm:"TEXT"`
	FileName              string  `xorm:"VARCHAR(255)"`
	Vector                string  `xorm:"VARCHAR(255)"`
	VersionStartExcluding string  `xorm:"VARCHAR(255)"`
	Cve                   string  `xorm:"VARCHAR(255)"`
	VersionEndExcluding   string  `xorm:"VARCHAR(255)"`
	VersionStartIncluding string  `xorm:"VARCHAR(255)"`
	Path                  string  `xorm:"VARCHAR(1024)"`
	CreateTime            int     `xorm:"INT(10)"`
}
