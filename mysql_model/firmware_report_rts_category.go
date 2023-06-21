package mysql_model

type FirmwareReportRtsCategory struct {
	Id                  int `xorm:"not null pk autoincr INT(10)"`
	ScannerId           int `xorm:"comment('固件扫描任务id') INT(10)"`
	TemplateId          int `xorm:"not null INT(10)"`
	InitFiles           int `xorm:"not null default 0 INT(10)"`
	ElfScanner          int `xorm:"not null default 0 INT(10)"`
	BinaryHardening     int `xorm:"not null default 0 INT(10)"`
	SymbolsXrefs        int `xorm:"not null default 0 INT(10)"`
	VersionScanner      int `xorm:"not null default 0 INT(10)"`
	CertificatesScanner int `xorm:"not null default 0 INT(10)"`
	LeaksScanner        int `xorm:"not null default 0 INT(10)"`
	PasswordScanner     int `xorm:"not null default 0 INT(10)"`
	ApkInfo             int `xorm:"not null default 0 INT(10)"`
	ApkCommonVul        int `xorm:"not null default 0 INT(10)"`
	ApkSensitiveInfo    int `xorm:"not null default 0 INT(10)"`
	LinuxBasicAudit     int `xorm:"not null default 0 INT(10)"`
	IsElf               int `xorm:"not null default 0 INT(10)"`
}
