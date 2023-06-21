package mysql_model

type KnowledgeTestCaseFile struct {
	Id           int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	TestCaseId   int    `xorm:"not null comment('测试用例id') INT(11)" json:"test_case_id"`
	Category     int    `xorm:"not null default 1 comment('业务类型，1测试脚本 2测试环境搭建示意图') TINYINT(3)" json:"category"`
	FileName     string `xorm:"not null comment('文件名称') VARCHAR(255)" json:"file_name"`
	FileSize     int64  `xorm:"not null default 0 comment('文件大小（kb）') BIGINT(20)" json:"file_size"`
	StorageType  int    `xorm:"not null comment('存储类型，1mongodb') TINYINT(3)" json:"storage_type"`
	FileUuid     string `xorm:"not null default '' comment('文件存储唯一标识uuid') VARCHAR(255)" json:"file_uuid"`
	CreateTime   int    `xorm:"not null comment('创建时间，即文件上传时间') INT(11)" json:"create_time"`
	IsDelete     int    `xorm:"not null default 1 comment('是否删除， 1否 2是') TINYINT(3)" json:"is_delete"`
	DeleteUserId int    `xorm:"not null comment('文件删除操作用户id') INT(11)" json:"delete_user_id"`
	DeleteTime   int    `xorm:"not null comment('删除时间') INT(11)" json:"delete_time"`
}

const (
	KnowledgeTestCaseFileCategoryScript  = 1 // 测试脚本
	KnowledgeTestCaseFileCategoryPicture = 2 // 示意图
)
