package mysql_model

import (
	"fmt"
	"time"

	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/mysql"
)

type FirmwareTask struct {
	Id                int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	Name              string `xorm:"not null comment('固件任务名称') INT(11)" json:"name"`
	TaskId            int    `xorm:"not null comment('任务id') INT(11)" json:"task_id"`
	TaskName          string `xorm:"not null comment('任务名称') VARCHAR(255)" json:"task_name"`
	YafafProjectId    int    `xorm:"not null comment('扫描项目id') INT(11)" json:"yafaf_project_id"`
	YafafId           int    `xorm:"not null comment('扫描任务id') INT(11)" json:"yafaf_id"`
	FileId            string `xorm:"not null comment('文件存储id') VARCHAR(255)" json:"file_id"`
	TemplateId        int    `xorm:"not null comment('模板id') INT(11)" json:"template_id"`
	TemplateName      string `xorm:"not null comment('模板名称') VARCHAR(255)" json:"template_name"`
	YafafDownloadPath string `xorm:"not null comment('yafaf平台固件下载路径') VARCHAR(255)" json:"yafaf_download_path"`
	DeviceName        string `xorm:"not null comment('硬件名称') VARCHAR(255)" json:"device_name"`
	DeviceModel       string `xorm:"not null comment('硬件模型') VARCHAR(255)" json:"device_model"`
	FirmwareVersion   string `xorm:"not null comment('固件版本') VARCHAR(255)" json:"firmware_version"`
	DeviceType        string `xorm:"not null comment('硬件类型') VARCHAR(255)" json:"device_type"`
	FirmwareName      string `xorm:"not null comment('固件名称') VARCHAR(255)" json:"firmware_name"`
	FirmwareSize      int    `xorm:"not null comment('固件大小') INT(11)" json:"firmware_size"`
	FirmwareMd5       string `xorm:"not null comment('固件md5') VARCHAR(255)" json:"firmware_md5"`
	Progress          int    `xorm:"not null comment('') INT(11)" json:"progress"`
	Status            int    `xorm:"not null comment('状态') INT(11)" json:"status"`
	SourceReport      string `xorm:"MEDIUMTEXT" json:"source_report"`
	UpdateTime        int    `xorm:"not null comment('') INT(11)" json:"update_time"`
	CreateTime        int    `xorm:"not null comment('创建时间') INT(11)" json:"create_time"`
}

func (this *FirmwareTask) Create() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}

func (t *FirmwareTask) Update(id int, rawInfo qmap.QM) error {
	params := qmap.QM{
		"e_id": id,
	}
	has, _ := sys_service.NewSessionWithCond(params).GetOne(t)
	if has {
		if val, has := rawInfo.TryString("name"); has {
			t.Name = val
		}
		if val, has := rawInfo.TryInt("yafaf_project_id"); has {
			t.YafafProjectId = val
		}
		if val, has := rawInfo.TryString("yafaf_download_path"); has {
			t.YafafDownloadPath = val
		}
		if val, has := rawInfo.TryInt("yafaf_id"); has {
			t.YafafId = val
		}
		if val, has := rawInfo.TryInt("status"); has {
			t.Status = val
		}
		if val, has := rawInfo.TryString("source_report"); has {
			t.SourceReport = val
		}
		_, err := sys_service.NewOrm().Table(t).ID(t.Id).Update(t)
		return err
	} else {
		return nil
	}
}

func (this *FirmwareTask) CreateByParams(rawInfo qmap.QM) (int64, error) {
	now := int(time.Now().Unix())
	if val, has := rawInfo.TryInt("task_id"); has {
		this.TaskId = val
	}
	if val, has := rawInfo.TryString("name"); has {
		this.Name = fmt.Sprintf("固件子任务_%s", val)
	}
	if val, has := rawInfo.TryString("file_id"); has {
		this.FileId = val
	}
	if val, has := rawInfo.TryInt("template_id"); has {
		this.TemplateId = val
	}
	this.CreateTime = now
	this.UpdateTime = now
	return mysql.GetSession().InsertOne(this)
}
