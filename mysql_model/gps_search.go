package mysql_model

import (
	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mysql"
	"strconv"
)

const (
	TASK_TYPE_REALTIME = 1
	QUERY_NUMBER       = 5
)

type GpsSearch struct {
	Id         int     `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	TaskId     int     `xorm:"not null comment('任务id') INT(11)" json:"task_id"`
	TemplateId int     `xorm:"not null comment('模板id') INT(11)" json:"template_id"`
	Type       int     `xorm:"not null comment('类型') INT(10)" json:"type"`
	Start      string  `xorm:"not null comment('起点') VARCHAR(20)" json:"start"`
	Middle     string  `xorm:"not null comment('中间') VARCHAR(20)" json:"middle"`
	End        string  `xorm:"not null comment('终点') VARCHAR(20)" json:"end"`
	Req        string  `xorm:"not null comment('参数') VARCHAR(20)" json:"req"`
	Lng        float32 `xorm:"not null comment('参数') VARCHAR(20)" json:"lng"`
	Lat        float32 `xorm:"not null comment('参数') VARCHAR(20)" json:"lat"`
	CreateTime string  `xorm:"not null comment('创建时间') varchar(20)" json:"create_time"`
}

func (g *GpsSearch) Create() (int64, error) {
	return mysql.GetSession().InsertOne(g)
}

func (g *GpsSearch) Update(cols ...string) (int64, error) {
	return mysql.GetSession().Table(g).ID(g.Id).Cols(cols...).Update(g)
}

func (g *GpsSearch) GetFive(tid string) ([]GpsSearch, error) {
	model := make([]GpsSearch, 0)
	taskId, _ := strconv.Atoi(tid)
	err := mysql.GetSession().Where("task_id = ?", taskId).
		And("type = ?", TASK_TYPE_REALTIME).
		Limit(QUERY_NUMBER).
		Desc("create_time").
		Find(&model)
	if err != nil {
		return model, err
	}
	return model, err
}

func (g *GpsSearch) FindMotion(taskId, templateId, limit int) ([]GpsSearch, error) {
	data := []GpsSearch{}
	err := mysql.GetSession().Where("task_id=?", taskId).Where("template_id=?", templateId).Where("type=?", common.GPS_TYPE_MOTION).OrderBy("id desc").Limit(limit).Find(&data)
	return data, err
}

func (g *GpsSearch) GetOne(searchId int) (bool, error) {
	bool, err := mysql.GetSession().Where("id=?", searchId).Where("type=?", common.GPS_TYPE_MOTION).Get(g)
	return bool, err
}
