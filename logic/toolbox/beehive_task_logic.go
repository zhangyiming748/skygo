package toolbox

import (
	"errors"
	"fmt"
	"skygo_detection/common"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/redis"
	"skygo_detection/mysql_model"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	TASK_DEAFAULT = 0 // 待测试
	TASK_TESTING  = 1 // 测试中
	TASK_COMPLETE = 2 // 完成
	TASK_FAIL     = 3 // 失败
	TASK_PAUSE    = 4 // 暂停---测过，但还没手动点完成
)

type Task struct{}

func (t Task) Create(ctx *gin.Context, name, tool_type, describe string) (int64, error) {
	task := new(mysql_model.Task)
	task.Name = name
	task.ToolType = tool_type
	task.Describe = describe
	task.Category = "隐私安全"
	task.Tool = "privacy"
	task.TaskUuid = mysql.GetTaskId()
	task.CreateUserId = int(http_ctx.GetUserId(ctx))
	task.CreateTime = int(time.Now().Unix())
	if _, err := task.Create(); err != nil {
		return 0, err
	}
	return int64(task.Id), nil
}

// 编辑任务基础信息
func (t Task) Update(ctx *gin.Context, id int, name, tool_type, describe string) (int64, error) {
	task, bool := mysql_model.TaskFindById(id)
	if !bool {
		return 0, errors.New("没找到记录")
	}
	if task.Status > 0 {
		return 0, errors.New("该任务已执行过，不能修改了")
	}
	task.Name = name
	task.ToolType = tool_type
	task.Describe = describe
	return task.Update()
}

// 获取任务基础信息
func (t Task) GetOne(ctx *gin.Context, id int) map[string]interface{} {
	info := make(map[string]interface{})
	task, bool := mysql_model.TaskFindById(id)
	if !bool {
		return nil
	}
	info["name"] = task.Name
	info["tool_type"] = task.ToolType
	info["test_type"] = t.ToolTypeCh(task.ToolType)

	info["describe"] = task.Describe
	info["user"] = ""
	if userModel, err := mysql_model.SysUserFindById(task.CreateUserId); err == nil {
		info["user"] = userModel.Realname
	}

	info["create_time"] = fmt.Sprint(time.Unix(int64(task.CreateTime), 0).Format("2006-01-02 15:04:05"))
	info["task_uuid"] = task.TaskUuid
	info["status"] = task.Status
	info["complete_time"] = fmt.Sprint(time.Unix(int64(task.CompleteTime), 0).Format("2006-01-02 15:04:05"))
	return info
}

// 完成任务
func (t Task) Complete(ctx *gin.Context, id int) (int64, error) {
	task, bool := mysql_model.TaskFindById(id)
	if !bool {
		return 0, errors.New("没找到记录")
	}
	task.Status = TASK_COMPLETE
	task.CompleteTime = int(time.Now().Unix())
	return task.Update()
}

// 判断字符串是否是预定的某个
func (t Task) CheckToolType(toolType string) bool {
	toolTypes := [3]string{"gsm-sniffer", "gsm-system", "lte-system"}
	for _, v := range toolTypes {
		if v == toolType {
			return true
		}
	}
	return false
}

// ToolType翻译
func (t Task) ToolTypeCh(toolType string) string {
	toolTypes := map[string]string{"gsm-sniffer": "GSM嗅探", "gsm-system": "GSM模拟", "lte-system": "LTE模拟"}
	v, ok := toolTypes[toolType]
	if ok {
		return v
	}
	return ""
}

func Start(uuid string) error {
	r := new(redis.Redis_service)
	if isExist(common.PRIVACY_TASK_START, uuid) {
		return nil
	}
	err := r.RPush(common.PRIVACY_TASK_START, uuid).Err()
	if err != nil {
		return err
	}
	return nil
}
func Stop(uuid string) error {
	r := new(redis.Redis_service)
	if isExist(common.PRIVACY_TASK_STOP, uuid) {
		return nil
	}
	err := r.RPush(common.PRIVACY_TASK_STOP, uuid).Err()
	if err != nil {
		return err
	}
	return nil
}
func isExist(key, uuid string) bool {
	r := new(redis.Redis_service)
	list := r.LRange(key, 0, -1).Val()
	for _, v := range list {
		if v == uuid {
			return true
		}
	}
	return false
}
