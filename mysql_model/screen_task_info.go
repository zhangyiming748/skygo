package mysql_model

import "skygo_detection/lib/common_lib/mysql"

// 任务信息
type ScreenTaskInfo struct {
	Id        int     `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	Name      string  `xorm:"not null comment('名称') VARCHAR(255)" json:"name"`
	Supplier  string  `xorm:"not null comment('供应商') VARCHAR(255)" json:"supplier"`
	TestPiece string  `xorm:"not null comment('零部件') VARCHAR(255)" json:"test_piece"`
	Category  string  `xorm:"not null comment('类型') VARCHAR(255)" json:"category"`
	Rate      float64 `xorm:"not null comment('类型') FLOAT(255)" json:"rate"`
}

func (this *ScreenTaskInfo) Create() (*ScreenTaskInfo, error) {
	_, err := mysql.GetSession().InsertOne(this)
	return this, err
}

func (this *ScreenTaskInfo) Remove(id int) error {
	_, err := mysql.GetSession().ID(id).Delete(this)
	return err
}
