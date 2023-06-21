package mysql_model

import (
	"errors"
	"github.com/gin-gonic/gin"
	"skygo_detection/common"
	"skygo_detection/guardian/app/sys_service"
	"skygo_detection/guardian/src/net/qmap"
	"skygo_detection/lib/common_lib/http_ctx"
	"skygo_detection/lib/common_lib/mysql"
	"time"
)

type Task struct {
	Id                 int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	TaskUuid           string `xorm:"not null comment('全局唯一id') VARCHAR(255)" json:"task_uuid"`
	Name               string `xorm:"not null comment('任务名称') VARCHAR(255)" json:"name"`
	Category           string `xorm:"not null comment('类型') VARCHAR(255)" json:"category"`
	AssetVehicleId     int    `xorm:"not null comment('车型id') INT(11)" json:"asset_vehicle_id"`
	PieceId            int    `xorm:"not null comment('测试件id') INT(11)" json:"piece_id"`
	PieceVersionId     int    `xorm:"not null comment('测试件id') INT(11)" json:"piece_version_id"`
	NeedConnected      int    `xorm:"not null comment('是否需要连接设备， 1是 2否') INT(11)" json:"need_connected"`
	FirmwareTemplateId int    `xorm:"not null comment('合规测试，测试模板id') INT(11)" json:"firmware_template_id"`
	ScenarioId         int    `xorm:"not null comment('场景id') INT(11)" json:"scenario_id"`
	ToolId             string `xorm:"not null comment('工具id,如果is_tool_task为不是工具任务，那就不关心这个值') VARCHAR(255)" json:"tool_id"`
	Tool               string `xorm:"not null comment('工具名称') VARCHAR(255)" json:"tool"`
	ToolTaskId         string `xorm:"not null comment('工具任务id') VARCHAR(255)" json:"tool_task_id"`
	IsToolTask         int    `xorm:"not null comment('是否是工具任务 1：是，0：不是') INT(11)" json:"is_tool_id"`
	Status             int    `xorm:"not null comment('状态') INT(11)" json:"status"`
	Describe           string `xorm:"not null comment('描述') VARCHAR(255)" json:"describe"`
	CreateUserId       int    `xorm:"not null comment('创建人id') INT(11)" json:"create_user_id"`
	LastOpId           int    `xorm:"not null comment('最近操作用户id') INT(11)" json:"last_op_id"`
	LastConnectTime    int    `xorm:"not null comment('上次连接更新时间，单位秒') INT(11)" json:"last_connect_time"`
	UpdateTime         int    `xorm:"not null comment('更新时间') INT(11)" json:"update_time"`
	HgClientInfo       string `xorm:"not null comment('合规上传的硬件信息，json格式') VARCHAR(255)" json:"hg_client_info"`
	HgFileUuid         string `xorm:"not null comment('合规硬件信息得到后匹配的测试用例压缩包的uuid') VARCHAR(255)" json:"hg_file_uuid"`
	CreateTime         int    `xorm:"created not null comment('创建时间') INT(11)" json:"create_time"`
	ClientInfoTime     int    `xorm:"not null comment('拿到终端上传信息的时间') INT(11)" json:"client_info_time"`
	CompleteTime       int    `xorm:"not null comment('任务完成的时间') INT(11)" json:"complete_time"`
	ToolType           string `xorm:"not null comment('一些新增加的类型') VARCHAR(255)" json:"tool_type"`
}

func (this *Task) Create() (int64, error) {
	return mysql.GetSession().InsertOne(this)
}
func (this Task) Sync(tid int) (int64, error) {
	return mysql.GetSession().ID(tid).Update(this)
}
func (this *Task) Update(cols ...string) (int64, error) {
	return mysql.GetSession().Table(this).ID(this.Id).Cols(cols...).Update(this)
}

func (this *Task) Remove() (int64, error) {
	return mysql.GetSession().ID(this.Id).Delete(this)
}

func UpdateStatus(taskId, status int) (int64, error) {
	data := map[string]int{"status": status}
	return mysql.GetSession().Table(Task{}).ID(taskId).Cols("status").Update(data)
}

// 更新子任务状态
func UpdateStatusBySubInfo(taskId int, data map[string]interface{}) (int64, error) {
	return mysql.GetSession().Table(Task{}).ID(taskId).Cols("status").Update(data)
}

// 根据工具类型和子任务ID ，查询记录
func GetTaskBySubInfo(subTaskId int, toolType string, taskUuid string) (*Task, error) {
	model := Task{}
	session := mysql.GetSession()
	session.Where("sub_task_id = ?", subTaskId).And("tool_type = ?", toolType).
		And("task_uuid = ?", taskUuid)
	_, err := session.Get(&model)
	if err != nil {
		return nil, err
	}
	return &model, nil
}

// 删除任务
func DeleteTaskByTaskId(taskId int) (res int64, err error) {
	res, err = mysql.GetSession().Where("id = ?", taskId).Delete(new(Task))
	if err != nil {
		return
	}
	return
}

// 创建表单
type TaskCreateForm struct {
	Name               string `json:"name"`                 // 任务名称，text框
	ScenarioId         int    `json:"scenario_id"`          // 场景id
	AssetVehicleId     int    `json:"asset_vehicle_id"`     // 通过车型品牌、车型代码两个下拉列表得到记录id
	PieceId            int    `json:"piece_id"`             // 测试件id
	PieceVersionId     int    `json:"piece_version_id"`     // 测试件id
	FirmwareTemplateId int    `json:"firmware_template_id"` // 合规测试，测试模板id， 目前是前端写死的 // [{"id":71,"name":"通用IoT固件检测模板"},{"id":103,"name":"APK扫描模板"},{"id":102,"name":"fs"}]
	NeedConnected      int    `json:"need_connected"`
	Describe           string `json:"describe"` // 描述
	ToolId             string `json:"tool_id"`
	IsToolTask         int    `json:"is_tool_id"`
	Tool               string `json:"tool"`
	Category           string `json:"category"`
}

// 任务创建
func TaskCreate(form *TaskCreateForm, ctx *gin.Context) (*Task, error) {
	if _, has := TaskFindByName(form.Name); has {
		return nil, errors.New("存在相同名称任务")
	}
	model := Task{}
	model.Name = form.Name
	model.AssetVehicleId = form.AssetVehicleId
	model.PieceId = form.PieceId
	model.PieceVersionId = form.PieceVersionId
	model.FirmwareTemplateId = form.FirmwareTemplateId
	model.NeedConnected = form.NeedConnected
	model.ScenarioId = form.ScenarioId
	model.Describe = form.Describe
	model.CreateUserId = int(http_ctx.GetUserId(ctx))
	model.LastOpId = int(http_ctx.GetUserId(ctx))
	model.CreateTime = int(time.Now().Unix())
	model.UpdateTime = 0
	model.CompleteTime = 0
	model.Status = common.TASK_STATUS_RUNNING // 创建后立即进行中
	model.TaskUuid = mysql.GetTaskId()

	if form.IsToolTask == 1 {
		model.ToolId = form.ToolId
		model.IsToolTask = form.IsToolTask
		model.Tool = form.Tool
		model.Category = form.Category // 类型分别为 固件扫描工具，车机漏扫工具 等等，不固定，所以存储字段
	} else {
		if scenario, has, err := KnowledgeScenarioFindById(model.ScenarioId); err == nil && has {
			model.Category = scenario.Name
		} else {
			model.Category = "测试模板"
		}
	}

	if _, err := mysql.GetSession().InsertOne(&model); err == nil {
		new(TaskLog).Insert(int(http_ctx.GetUserId(ctx)), http_ctx.GetUserName(ctx), &model)
		return &model, nil
	} else {
		return nil, err
	}
}

// 根据任务名称，查询记录
func TaskFindByName(name string) (*Task, bool) {
	model := Task{}
	session := mysql.GetSession()
	session.Where("name = ?", name)
	if has, err := session.Get(&model); err != nil {
		panic(err)
	} else {
		return &model, has
	}
}

// ------------------------------------------------
// 合规终端信息，在task表中hg_client_info字段存为json
type HgClientInfo struct {
	Cpu       string `json:"cpu"`        // cpu位数
	OsType    string `json:"os_type"`    // 操作系统类型
	OsVersion string `json:"os_version"` // 操作系统版本
}

func (this *Task) UpdateTaskById(id int, info qmap.QM, userId int, username string) (*Task, error) {
	has, err := sys_service.NewSession().Session.Where("id = ?", id).Get(this)
	if err != nil {
		return nil, err
	}
	if !has {
		return nil, errors.New("Item not found")
	}
	if val, has := info.TryInt("status"); has {
		this.Status = val
		this.CompleteTime = int(time.Now().Unix())
	}
	_, err = sys_service.NewOrm().Table(this).ID(this.Id).AllCols().Update(this)
	if err == nil {
		if _, has := info.TryInt("status"); has {
			// 如果修改了任务状态，则记录日志
			new(TaskLog).Insert(userId, username, this)
		}
	}
	return this, err
}

// 根据任务id ，查询记录
func TaskFindById(id int) (*Task, bool) {
	model := Task{}
	session := mysql.GetSession()
	session.Where("id = ?", id)
	if has, err := session.Get(&model); err != nil {
		panic(err)
	} else {
		return &model, has
	}
}

// 根据任务id ，查询记录
func GetTaskByUuid(uuid string) (*Task, error) {
	model := Task{}
	session := mysql.GetSession()
	session.Where("task_uuid = ?", uuid)
	if _, err := session.Get(&model); err != nil {
		return nil, err
	}
	return &model, nil
}

func GetTaskCategoryList(uid int) []string {
	models := make([]Task, 0)
	// 只能看到自己创建的全部任务属性
	mysql.GetSession().Where("create_user_id=?", uid).Distinct("category").Table(Task{}).
		Limit(10000).
		Find(&models)
	var result = make([]string, 0)
	for _, model := range models {
		result = append(result, model.Category)
	}
	return result
}

// 根据任务id ，查询记录
func TaskGetStatusById(id int) int {
	model := Task{}
	if has, err := sys_service.NewSession().Session.Where("id = ?", id).Get(&model); err != nil || !has {
		return 0
	} else {
		return model.Status
	}
}

// 根据uid ，查询id
func GetIdByUid(uid string) int {
	model := Task{}
	if has, err := sys_service.NewSession().Session.Where("task_uuid = ?", uid).Get(&model); err != nil || !has {
		return 0
	} else {
		return model.Id
	}
}
