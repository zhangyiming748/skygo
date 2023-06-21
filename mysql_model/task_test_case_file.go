package mysql_model

// 任务下的测试用例的文件内容
type TaskTestCaseFile struct {
	Id             int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	TaskTestCaseId int    `xorm:"not null comment('测试案例id') INT(11)" json:"task_test_case_id"`
	TestFileType   int    `xorm:"not null comment('文件类型，1测试附件 2设备上传文件') TINYINT(3)" json:"storage_type"`
	FileName       string `xorm:"not null comment('文件名称') VARCHAR(255)" json:"file_name"`
	FileSize       int64  `xorm:"not null default 0 comment('文件大小（kb）') BIGINT(20)" json:"file_size"`
	StorageType    int    `xorm:"not null comment('存储类型，1mongodb') TINYINT(3)" json:"storage_type"`
	FileUuid       string `xorm:"not null default '' comment('文件存储唯一标识uuid') VARCHAR(255)" json:"file_uuid"`
	CreateTime     int    `xorm:"not null comment('创建时间，即文件上传时间') INT(11)" json:"create_time"`
	IsDelete       int    `xorm:"not null default 1 comment('是否删除， 1否 2是') TINYINT(3)" json:"is_delete"`
	DeleteUserId   int    `xorm:"not null comment('文件删除操作用户id') INT(11)" json:"delete_user_id"`
	DeleteTime     int    `xorm:"not null comment('删除时间') INT(11)" json:"delete_time"`
}
