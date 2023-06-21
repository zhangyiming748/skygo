package mysql_model

import (
	"time"

	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
)

type ScannerTask struct {
	Id           int    `xorm:"not null pk autoincr INT(10)"`
	ScannerId    int    `xorm:"not null comment('扫描任务id') INT(11)"`
	Name         string `xorm:"comment('任务名称') VARCHAR(255)"`
	ScannerType  string `xorm:"comment('任务类型(固件扫描:firmware_scanner)') VARCHAR(255)"`
	Status       int    `xorm:"comment('任务状态(0:待执行，1:执行中，2：执行成功，3：执行失败)') TINYINT(255)"`
	RetryTimes   int    `xorm:"comment('重试次数') INT(11)"`
	NextExecTime int    `xorm:"comment('下一次触发时间') INT(11)"`
}

func (t *ScannerTask) TaskInsert(taskId int, name, scannerType string) error {
	t.ScannerId = taskId
	t.Name = name
	t.ScannerType = scannerType
	t.Status = common.SCANNER_STATUS_READY
	t.RetryTimes = 0
	t.NextExecTime = int(time.Now().Unix())
	_, err := sys_service.NewOrm().InsertOne(t)
	return err
}

func (t *ScannerTask) Update(id int, rawInfo qmap.QM) error {
	params := qmap.QM{
		"e_id": id,
	}
	has, _ := sys_service.NewSessionWithCond(params).GetOne(t)
	if has {
		if status, has := rawInfo.TryInt("status"); has {
			if status == common.SCANNER_STATUS_SUCCESS || status == common.SCANNER_STATUS_FAILURE {
				// 如果即将要更新的任务状态是"执行成功"或者"执行失败"，则直接删除该任务，让扫描任务表始终只保留当前待执行或者正在执行的任务
				t.RemoveById(t.Id)
				return nil
			} else {
				t.Status = status
			}
		}
		if val, has := rawInfo.TryString("name"); has {
			t.Name = val
		}
		if val, has := rawInfo.TryString("task_type"); has {
			t.ScannerType = val
		}
		if val, has := rawInfo.TryInt("retry_times"); has {
			t.RetryTimes = val
		}
		if val, has := rawInfo.TryInt("next_exec_time"); has {
			t.NextExecTime = val
		}
		_, err := sys_service.NewOrm().Table(t).ID(t.Id).AllCols().Update(t)
		return err
	} else {
		return nil
	}
}

func (t *ScannerTask) RemoveById(id int) {
	sys_service.NewSession().DeleteByIds(t, []int{id})
}
