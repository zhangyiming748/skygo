package mysql_model

import "skygo_detection/lib/common_lib/mysql"

// 零部件测试进展
type ScreenPieceTestProgress struct {
	Id        int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	Supplier  string `xorm:"not null comment('供应商') VARCHAR(255)" json:"supplier"`
	TestPiece string `xorm:"not null comment('零部件') VARCHAR(255)" json:"test_piece"`
	Status    int    `xorm:"not null comment('状态') INT(1)" json:"status"`
	StartTime int    `xorm:"not null comment('开始时间') INT(11)" json:"start_time"`
	EndTime   int    `xorm:"not null comment('结束时间') INT(11)" json:"end_time"`
}

func (this *ScreenPieceTestProgress) Create() (*ScreenPieceTestProgress, error) {
	_, err := mysql.GetSession().InsertOne(this)
	return this, err
}

func (this *ScreenPieceTestProgress) Remove(id int) error {
	_, err := mysql.GetSession().ID(id).Delete(this)
	return err
}
