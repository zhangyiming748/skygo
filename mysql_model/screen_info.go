package mysql_model

import "skygo_detection/lib/common_lib/mysql"

// 大屏主要信息，包括车型/零部件/测试任务/测试用例
type ScreenInfo struct {
	Id              int `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id" `
	VechielNumber   int `xorm:"not null comment('车型个数') INT(11)" json:"vechiel_number"`
	TestPieceNumber int `xorm:"not null comment('零部件/测试件个数') INT(11)" json:"test_piece_number"`
	TaskNumber      int `xorm:"not null comment('测试任务个数') INT(11)" json:"task_number"`
	TestCaseNumber  int `xorm:"not null comment('测试用例个数') INT(11)" json:"test_case_number"`
}

func (this *ScreenInfo) Create() (*ScreenInfo, error) {
	_, err := mysql.GetSession().InsertOne(this)
	return this, err
}

func (this *ScreenInfo) Remove(id int) error {
	_, err := mysql.GetSession().ID(id).Delete(this)
	return err
}
