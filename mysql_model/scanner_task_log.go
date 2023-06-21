package mysql_model

import (
	"time"

	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"
)

type ScannerTaskLog struct {
	Id            int    `xorm:"not null pk autoincr INT(10)"`
	ScannerTaskId int    `xorm:"comment('关联的扫描任务id') INT(10)"`
	ScannerId     int    `xorm:"comment('扫描任务id') INT(10)"`
	Name          string `xorm:"comment('任务名称') VARCHAR(255)"`
	Level         string `xorm:"comment('日志类型(error、info)') VARCHAR(255)"`
	ScannerType   string `xorm:"comment('任务类型(固件扫描：scanner_firmware)') VARCHAR(255)"`
	Msg           string `xorm:"comment('日志信息') VARCHAR(255)"`
	ErrMsg        string `xorm:"comment('错误信息') VARCHAR(255)"`
	CreateTime    int    `xorm:"comment('创建时间') INT(11)"`
}

func (t *ScannerTaskLog) Insert(taskInfo qmap.QM, msg, errMsg string) {
	t.ScannerTaskId = taskInfo.MustInt("id")
	t.ScannerId = taskInfo.MustInt("scanner_id")
	t.ScannerType = taskInfo.String("scanner_type")
	t.Name = taskInfo.MustString("name")
	t.ErrMsg = errMsg
	t.Msg = msg
	if errMsg == "" {
		t.Level = "info"
	} else {
		t.Level = "error"
	}
	t.CreateTime = int(time.Now().Unix())
	if _, err := sys_service.NewOrm().InsertOne(t); err != nil {
		panic(err)
	}
}
