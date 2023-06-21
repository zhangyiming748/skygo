package mysql_model

import "skygo_detection/lib/common_lib/mysql"

type GpsCheat struct {
	Id         int     `xorm:"not null pk autoincr comment('主键id') INT(10)" json:"id"`
	TaskId     int     `xorm:"not null comment('任务id') INT(10)" json:"task_id"`
	SearchId   int     `xorm:"comment('模板id') INT(10)" json:"Search_id"`
	Type       int     `xorm:"comment('类型') INT(10)" json:"type"`
	Start      string  `xorm:"comment('起点') VARCHAR(256)" json:"start"`
	Middle     string  `xorm:"comment('中间') VARCHAR(1024)" json:"middle"`
	End        string  `xorm:"comment('终点') VARCHAR(256)" json:"end"`
	Status     int     `xorm:"comment('类型') TEXT(3)" json:"status"`
	Req        string  `xorm:"comment('参数') VARCHAR(20)" json:"req"`
	Resp       string  `xorm:"comment('参数') VARCHAR(128)" json:"resp"`
	Lng        float32 `xorm:"comment('参数') VARCHAR(12)" json:"lng"`
	Lat        float32 `xorm:"comment('参数') VARCHAR(12)" json:"lat"`
	CreateTime string  `xorm:"comment('创建时间') varchar(20)" json:"create_time"`
}

func (g *GpsCheat) Create() (int64, error) {
	return mysql.GetSession().InsertOne(g)
}

func (g *GpsCheat) GetLatestByTaskId(taskId, limit int) (bool, error) {
	b, err := mysql.GetSession().Where("task_id=?", taskId).OrderBy("id desc").Limit(limit).Get(g)
	return b, err
}

func (this *GpsCheat) Update(tid int, cols ...string) error {
	_, err := mysql.GetSession().Where("task_id = ?", tid).
		Cols(cols...).
		Update(this)
	return err
}
