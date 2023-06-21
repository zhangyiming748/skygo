package mongo_model

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"
)

// 任务
type EvaluateTaskNode struct {
	Id         primitive.ObjectID `bson:"_id,omitempty"`                            //项目任务id
	Name       string             `bson:"name" json:"name"`                         //阶段名称
	Status     int                `bson:"status" json:"status"`                     //状态
	ProjectId  string             `bson:"project_id" json:"project_id"`             //项目id
	TaskId     string             `bson:"task_id" json:"task_id"`                   //任务id
	History    []History          `bson:"history" json:"history"`                   //阶段日志
	UpdateTime int64              `bson:"update_time" json:"update_time"`           //更新时间
	CreateTime int64              `bson:"create_time" json:"create_time,omitempty"` //创建时间
}

type History struct {
	Result    bool   `bson:"result" json:"result"`
	Comment   string `bson:"comment" json:"comment"`
	Operation string `bson:"operation" json:"operation"`
	OpId      int    `bson:"op_id" json:"op_id"`
	OpTime    int64  `bson:"op_time" json:"op_time"`
}

func (this *EvaluateTaskNode) Create(rawInfo qmap.QM, opId int) (*EvaluateTaskNode, error) {
	//测试任务所属项目和任务必须存在
	projectId := rawInfo.MustString("project_id")
	taskId := rawInfo.MustString("task_id")
	name := rawInfo.MustString("name")
	status := rawInfo.MustInt("status")

	//检查任务是否存在
	if err := checkProjectTaskExist(projectId, taskId); err == nil {
		return nil, err
	}
	//检查当前阶段是否存在
	if err := CheckNodeExist(taskId, name); err != nil {
		return nil, err
	}
	this.Id = primitive.NewObjectID()
	this.ProjectId = projectId
	this.TaskId = taskId
	this.Status = status
	this.Name = name
	this.History = append(this.History, rawInfo["history"].(History))
	this.UpdateTime = custom_util.GetCurrentMilliSecond()
	this.CreateTime = custom_util.GetCurrentMilliSecond()
	if _, err := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_TASK_NODE).InsertOne(context.Background(), this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateTaskNode) Update(rawInfo qmap.QM) (*EvaluateTaskNode, error) {
	params := qmap.QM{
		"e_task_id": rawInfo.MustString("task_id"),
		"e_name":    rawInfo.MustString("name"),
	}
	widget := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_TASK_NODE, params)
	if err := widget.One(&this); err == nil {
		this.History = append(this.History, rawInfo["history"].(History))
		coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_TASK_NODE)
		if _, err := coll.UpdateByID(context.Background(), this.Id, bson.M{"$set": this}); err != nil {
			return nil, err
		}
	} else {
		return nil, errors.New("Item not found")
	}
	return this, nil
}

func (this *EvaluateTaskNode) GetTaskNode(taskId string, status int) (qmap.QM, error) {
	if status == 0 {
		return this.GetTaskNodeDefault()
	}
	return this.GetTaskNodeHistory(taskId, status)
}

func (this *EvaluateTaskNode) GetTaskNodeDefault() (qmap.QM, error) {
	return qmap.QM{
		"index": 0,
		"list":  this.appendNode(0),
	}, nil
}

func (this *EvaluateTaskNode) GetTaskNodeHistory(taskId string, status int) (qmap.QM, error) {
	params := qmap.QM{
		"e_task_id": taskId,
	}

	mongoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK_NODE, params)
	mongoClient.AddSorter("level", 1)
	res, err := mongoClient.SetLimit(5000).Get()
	if err != nil {
		return nil, errors.New("Item not found")
	}
	//检查当前状态，是否有node记录
	hasNode := false
	lastNodeStatus := 0
	appendIndex := status
	index := 1

	for _, item := range *res {
		if item["status"] == status {
			hasNode = true
		}
		if item["status"].(int) > lastNodeStatus {
			lastNodeStatus = item["status"].(int)
		}
	}
	if status <= lastNodeStatus {
		appendIndex = lastNodeStatus
	}

	if hasNode == false {
		appendIndex--
	}

	list := *res
	list = append(list, this.appendNode(appendIndex)...)

	for i, item := range list {
		if item["status"] == status {
			index = i + 1
		}
	}

	return qmap.QM{
		"index": index,
		"list":  list,
	}, nil
}

func (this *EvaluateTaskNode) appendNode(index int) []map[string]interface{} {
	defaultNode := []map[string]interface{}{
		{"status": common.PTS_CREATE, "name": "create", "history": []interface{}{}},
		{"status": common.PTS_TASK_AUDIT, "name": "task_audit", "history": []interface{}{}},
		{"status": common.PTS_TEST, "name": "test", "history": []interface{}{}},
		{"status": common.PTS_REPORT_AUDIT, "name": "report_audit", "history": []interface{}{}},
		{"status": common.PTS_FINISH, "name": "finish", "history": []interface{}{}},
	}
	return defaultNode[index:]
}

// 判断任务是否存在
func CheckNodeExist(taskId, name string) error {
	params := qmap.QM{
		"e_name":    name,
		"e_task_id": taskId,
	}
	if _, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK_NODE, params).GetOne(); err == nil {
		return errors.New(fmt.Sprintf("项目任务: %s 阶段： %s 已存在", taskId, name))
	}
	return nil
}

// -------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------

// 根据任务，创建任务节点记录
// stageName string 任务阶段名称
func EvaluateTaskNodeCreate(ctx context.Context, task *EvaluateTask, history *History, stageName string) (*EvaluateTaskNode, error) {
	// 测试任务所属项目和任务必须存在
	// 检查任务是否存在
	if err := checkProjectTaskExist(task.ProjectId, task.Id.Hex()); err == nil {
		return nil, err
	}

	// 检查当前阶段是否存在
	if err := CheckNodeExist(task.Id.Hex(), stageName); err != nil {
		return nil, err
	}

	model := EvaluateTaskNode{}
	model.Id = primitive.NewObjectID()
	model.TaskId = task.Id.Hex()     // 任务ID，Hex形式
	model.ProjectId = task.ProjectId // 项目ID，Hex形式
	model.Status = task.Status
	model.Name = stageName // 注意这里存的是阶段名称
	model.History = []History{
		*history,
	}
	model.UpdateTime = custom_util.GetCurrentMilliSecond()
	model.CreateTime = custom_util.GetCurrentMilliSecond()

	_, err := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_TASK_NODE).InsertOne(ctx, &model)
	if err != nil {
		return nil, err
	}
	return &model, err
}

// 根据任务的修改，修改任务节点记录
func EvaluateTaskNodeUpdate(ctx context.Context, task *EvaluateTask, history *History, stageName string) (*EvaluateTaskNode, error) {
	params := qmap.QM{
		"e_task_id": task.Id.Hex(),
		"e_name":    stageName,
	}

	model := EvaluateTaskNode{}
	err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_TASK_NODE, params).One(&model)
	if err != nil {
		return nil, errors.New("item not found")
	}

	model.History = append(model.History, *history)

	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_TASK_NODE)
	_, err = coll.UpdateByID(ctx, model.Id, model)
	if err != nil {
		return nil, err
	}
	return &model, nil
}
