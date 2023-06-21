package mysql_model

import (
	"skygo_detection/lib/common_lib/mysql"
)

type ReportTask struct {
	Id         int    `xorm:"not null pk autoincr INT(11)"`
	Name       string `xorm:"VARCHAR(255)"`
	TaskId     int    `xorm:"comment('父类任务id') INT(11)"`
	Status     int    `xorm:"comment('任务状态，0：未执行,1：进行中，2：完成') INT(1)"`
	ReportType int    `xorm:"INT(1)"`
	CreateTime string `xorm:"VARCHAR(11)"`
	EndTime    string `xorm:"VARCHAR(11)"`
	FileId     string `xorm:"VARCHAR(11)"`
	ReportName string `xorm:"VARCHAR(11)"`
	ExcelId    string `xorm:"VARCHAR(11)"`
	ExcelName  string `xorm:"VARCHAR(11)"`
	PdfId      string `xorm:"VARCHAR(11)"`
	PdfName    string `xorm:"VARCHAR(11)"`
	Count      int    `xorm:"INT(10)"`
}

const (
	// 任务状态
	ReportTaskStatusDefault = 1 // 未开始
	ReportTaskStatusRunning = 2 // 进行中
	ReportTaskStatusSuccess = 3 // 成功
	// 任务类型
	ReportTypeVul  = 1 // 车机漏扫任务
	ReportTypeFirm = 2 // 固件扫描任务
	ReportTypeHG   = 3 // 合规扫描任务
)

func (this *ReportTask) Create() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}
func FindByFileID(file_id string) (*ReportTask, error) {
	model := ReportTask{}
	_, err := mysql.GetSession().Where("file_id = ?", file_id).Or("excel_id = ?", file_id).Or("pdf_id = ?", file_id).Get(&model)
	return &model, err
}

func (this *ReportTask) AddDownloadCount(n int) error {
	_, err := mysql.GetSession().ID(this.Id).Incr("count", n).Update(this)
	if err != nil {
		return err
	}
	return nil
}

func (this *ReportTask) UpdateByID() error {
	_, err := mysql.GetSession().ID(this.Id).Update(this)
	if err != nil {
		return err
	}
	return nil
}
