package mysql_model

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mysql"
)

type VulTask struct {
	Id            int    `xorm:"not null pk autoincr INT(11)"`
	Name          string `xorm:"comment('任务名称') VARCHAR(255)"`
	TaskId        string `xorm:"comment('下发到终端的ID') VARCHAR(255)"`
	ParentId      int    `xorm:"comment('父任务ID') VARCHAR(255)"`
	VulScannerId  string `xorm:"comment('漏洞检测id') VARCHAR(255)"`
	Status        int    `xorm:"comment('任务状态(4未展示，0未开始，1测试中，2测试完成)') TINYINT(255)"`
	TestTime      int    `xorm:"comment('测试时间') INT(11)"`
	CreateTime    int    `xorm:"updated comment('创建时间') INT(11)"`
	SearchContent string `xorm:"comment('搜索字段') VARCHAR(255)"`
}

func (this *VulTask) Create(name string, parentId int, taskId string) (*VulTask, error) {
	// 检查name是否存在
	if this.CheckNameExist(name) {
		return nil, errors.New("名称已存在")
	}
	this.Name = name
	if taskId != "" {
		this.TaskId = taskId
	} else {
		this.TaskId = this.getTaskId()
	}
	this.ParentId = parentId
	this.Status = common.VUL_UNSTART
	this.TestTime = 0
	this.CreateTime = int(custom_util.GetCurrentMilliSecond())
	this.SearchContent = fmt.Sprintf("%s_%s", this.Name, this.TaskId)
	if _, err := mysql.GetSession().InsertOne(this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *VulTask) CheckNameExist(name string) bool {
	session := mysql.GetSession()
	session.Where("name=?", name)
	if num, _ := session.Count(this); num != 0 {
		return true
	}
	return false
}

func (this *VulTask) getTaskId() string {
	now := int(time.Now().UnixNano() / 1000)
	str := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	assetID := ""
	var remainder int
	var remainderStr string
	for now != 0 {
		remainder = now % 36
		if remainder < 36 && remainder > 9 {
			remainderStr = str[remainder]
		} else {
			remainderStr = strconv.Itoa(remainder)
		}
		assetID = remainderStr + assetID
		now = now / 36
	}
	if len(assetID) > 8 {
		rs := []rune(assetID)
		assetID = string(rs[:8])
	}

	return assetID
}

func (this *VulTask) UpdateByTaskId(taskId string, rawInfo qmap.QM) error {
	params := qmap.QM{
		"e_task_id": taskId,
	}
	has, _ := sys_service.NewSessionWithCond(params).GetOne(this)
	if has {
		if val, has := rawInfo.TryInt("status"); has {
			this.Status = val
		}
		_, err := sys_service.NewOrm().Table(this).ID(this.Id).Update(this)
		return err
	} else {
		return nil
	}
}
