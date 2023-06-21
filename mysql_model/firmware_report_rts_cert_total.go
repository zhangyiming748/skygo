package mysql_model

type FirmwareReportRtsCertTotal struct {
	Id                       int    `xorm:"not null pk autoincr INT(10)"`
	ScannerId                int    `xorm:"comment('固件扫描任务id') index INT(10)"`
	CertificateCount         int    `xorm:"not null default 0 INT(10)"`
	CertificateOverdateCount int    `xorm:"not null default 0 INT(10)"`
	PrivateKeyCount          int    `xorm:"not null default 0 INT(10)"`
	Level                    int    `xorm:"INT(10)"`
	Part                     string `xorm:"VARCHAR(255)"`
	CreateTime               int    `xorm:"INT(10)"`
}
