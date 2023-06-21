package mysql_model

import (
	"errors"
	"time"

	"skygo_detection/lib/common_lib/mysql"

	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"
)

type HgTestTask struct {
	Id              int    `xorm:"not null pk autoincr INT(10)"`
	Name            string `xorm:"VARCHAR(255)"`
	TaskUuid        string `xorm:"comment('场景任务id') VARCHAR(255)"`
	Status          string `xorm:"VARCHAR(64)"`
	Cpu             string `xorm:"VARCHAR(255)"`
	OsType          string `xorm:"VARCHAR(255)"`
	OsVersion       string `xorm:"VARCHAR(255)"`
	LastConnectTime int64  `xorm:"default 0 BIGINT(16)"`
	CreateTime      int    `xorm:"INT(10)"`
}

var HgTestTaskCaseStatusList = map[int]string{
	1: "待测试",
	2: "测试中",
	3: "分析中",
	4: "通过",
	5: "未通过",
	6: "失败",
}

func (this *HgTestTask) FindOne(taskUuid string) (*HgTestTask, error) {
	has, err := sys_service.NewSession().Session.Where("task_uuid = ?", taskUuid).Get(this)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("Item not found")
	}
	return this, nil
}

func (this *HgTestTask) Update(taskUuid string, info qmap.QM) (*HgTestTask, error) {
	has, err := sys_service.NewSession().Session.Where("task_uuid = ?", taskUuid).Get(this)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("Item not found")
	}
	if val, has := info.TryString("cpu"); has {
		this.Cpu = val
	}
	if val, has := info.TryString("os_type"); has {
		this.OsType = val
	}
	if val, has := info.TryString("os_version"); has {
		this.OsVersion = val
	}
	if val, has := info.TryString("status"); has {
		this.Status = val
	}
	_, err = sys_service.NewSession().Session.ID(this.Id).Update(this)
	return this, err
}

func (this *HgTestTask) Create() error {
	_, err := mysql.GetSession().InsertOne(this)
	return err
}

func (this *HgTestTask) GetTaskInfo(taskUuid string) (*HgTestTask, error) {
	has, err := sys_service.NewSession().Session.Where("task_uuid = ?", taskUuid).Get(this)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("Item not found")
	}
	return this, nil
}

func (this *HgTestTask) UpdateTerminalConnectionTime(taskUuid string) error {
	has, err := sys_service.NewSession().Session.Where("task_uuid = ?", taskUuid).Get(this)
	if err != nil {
		return err
	}
	if !has {
		return errors.New("Item not found")
	}
	this.LastConnectTime = time.Now().Unix()
	_, err = sys_service.NewSession().Session.ID(this.Id).Update(this)
	return err
}
