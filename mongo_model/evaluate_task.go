package mongo_model

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/log"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm_mongo"
	"skygo_detection/mysql_model"

	"github.com/globalsign/mgo/bson"
)

// 任务
type EvaluateTask struct {
	Id                primitive.ObjectID `bson:"_id,omitempty"`                                  // 项目任务id
	Name              string             `bson:"name" json:"name"`                               // 任务名称
	ProjectId         string             `bson:"project_id" json:"project_id"`                   // 项目id
	AssetVersions     qmap.QM            `bson:"asset_versions" json:"asset_versions"`           // 任务资产
	Status            int                `bson:"status" json:"status"`                           // 项目任务状态  1:创建、2:任务审核、3:测试、4:报告审核、5:完成
	IsDeleted         int                `bson:"is_deleted" json:"is_deleted"`                   // 项目任务删除 0未删除，1已删除
	TestPhase         int                `bson:"test_phase" json:"test_phase"`                   // 测试阶段 （1:初测、2:复测1、3:复测2、4:复测3 ...
	EvaluateItemIds   []string           `bson:"evaluate_item_ids" json:"evaluate_item_ids"`     // 项目任务包含的测试用例id
	AuditStatus       int                `bson:"audit_status" json:"audit_status"`               // 审核状态  0未审核，1已审核
	AuditRecord       []interface{}      `bson:"audit_record" json:"audit_record"`               // 审核记录
	ReportAuditStatus int                `bson:"report_audit_status" json:"report_audit_status"` // 报告审核状态  0未审核，1已审核
	ReportAuditRecord []interface{}      `bson:"report_audit_record" json:"report_audit_record"` // 报告审核记录
	OpId              int                `bson:"op_id" json:"op_id"`                             // 操作人id
	TaskAuditorId     int                `bson:"task_auditor_id" json:"task_auditor_id"`         // 任务审核人id
	TaskAuditorTime   int64              `bson:"task_auditor_time" json:"task_auditor_time"`     // 任务审核时间
	ReportAuditorId   int                `bson:"report_auditor_id" json:"report_auditor_id"`     // 报告审核人id
	ReportAuditorTime int64              `bson:"report_auditor_time" json:"report_auditor_time"` // 报告审核时间
	TesterId          int                `bson:"tester_id" json:"tester_id"`                     // 测试团队id
	CurrentOpId       int                `bson:"current_op_id" json:"current_op_id"`             // 当前操作人员id
	UpdateTime        int64              `bson:"update_time" json:"update_time"`                 // 更新时间
	CreateTime        int64              `bson:"create_time" json:"create_time,omitempty"`       // 创建时间
	SubmitTime        int64              `bson:"submit_time" json:"submit_time"`                 // 提交时间
}

// func (this *EvaluateTask) Create(rawInfo qmap.QM, opId int64) (*EvaluateTask, error) {
// 	// 测试任务所属项目必须存在
// 	projectId := rawInfo.MustString("project_id")
// 	name := rawInfo.MustString("name")
// 	params := qmap.QM{
// 		"e__id": bson.ObjectIdHex(projectId),
// 	}
//
// 	if err := CheckIsProjectManager(projectId, opId); err != nil {
// 		return nil, err
// 	}
//
// 	evaluateItemIds := rawInfo.SliceString("evaluate_item_ids")
// 	if len(evaluateItemIds) == 0 {
// 		return nil, errors.New("测试用例不能为空")
// 	}
//
// 	// 任务名称不能重复
// 	params = qmap.QM{
// 		"e_project_id": projectId,
// 		"e_name":       name,
// 	}
// 	if info, _ := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params).GetOne(); (*info)["_id"] != nil {
// 		return nil, errors.New("该任务名称已被使用")
// 	}
//
// 	// 新建项目任务中关联的测试用例必须还没有与其他项目任务关联: 判断status=0可绑定，status=1不可绑定
// 	for _, item := range evaluateItemIds {
// 		params = qmap.QM{
// 			"e__id":        item,
// 			"e_project_id": projectId,
// 		}
// 		if temp, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params).GetOne(); err != nil || temp.Int("status") != 0 {
// 			return nil, errors.New(fmt.Sprintf("项目任务创建失败:测试用例%s无法添加到该项目任务中", item))
// 		}
// 	}
// 	this.Id = bson.NewObjectId()
// 	this.ProjectId = projectId
// 	this.Name = name
// 	this.AssetVersions = *rawInfo.Map("asset_versions")
// 	this.EvaluateItemIds = evaluateItemIds
//
// 	this.Status = common.PTS_CREATE
// 	this.IsDeleted = common.PSD_DEFAULT
// 	this.AuditStatus = common.PTS_AUDIT_STATUS_NEW
// 	this.ReportAuditStatus = common.PTS_AUDIT_STATUS_NEW
// 	this.TestPhase = rawInfo.Int("test_phase")
// 	this.UpdateTime = custom_util.GetCurrentMilliSecond()
// 	this.CreateTime = custom_util.GetCurrentMilliSecond()
// 	this.SubmitTime = 0
// 	this.TaskAuditorTime = 0
// 	this.ReportAuditorTime = 0
// 	this.OpId = int(opId)
// 	this.CurrentOpId = int(opId)
//
// 	if err := mongo.NewMgoSession(common.MC_EVALUATE_TASK).Insert(this); err == nil {
// 		// 项目任务创建成功后将测试用例与项目任务绑定(绑定后，测试用例将不能与其他项目任务绑定)
// 		new(EvaluateItem).BindEvaluateTask(evaluateItemIds, this.Id.Hex(), common.NOT_PREBIND)
// 		// 创建阶段记录
// 		_ = this.createTaskNode()
// 		return this, nil
// 	} else {
// 		return nil, err
// 	}
// }

// func (this *EvaluateTask) createTaskNode() error {
// 	nodeInfo := qmap.QM{
// 		"project_id": this.ProjectId,
// 		"task_id":    this.Id.Hex(),
// 		"status":     this.Status,
// 		"name":       common.PTS_SIGN_CREATE,
// 		"history": History{
// 			Result:    true,
// 			Operation: "任务创建",
// 			Comment:   "任务创建",
// 			OpId:      this.OpId,
// 			OpTime:    this.CreateTime,
// 		},
// 	}
// 	_, err := new(EvaluateTaskNode).Create(nodeInfo, this.OpId)
// 	return err
// }

// func (this *EvaluateTask) Update(id string, rawInfo qmap.QM) (*EvaluateTask, error) {
// 	params := qmap.QM{
// 		"e__id": bson.ObjectIdHex(id),
// 	}
// 	mongoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params)
// 	if err := mongoClient.One(&this); err == nil {
// 		if val, has := rawInfo.TryString("name"); has {
// 			this.Name = val
// 		}
// 		itemIds, has := rawInfo.TrySlice("evaluate_item_ids")
//
// 		delItemIds := []string{}
// 		newItemIds := rawInfo.SliceString("evaluate_item_ids")
//
// 		if has && len(itemIds) > 0 {
// 			delItemIds = custom_util.DifferenceDel(this.EvaluateItemIds, newItemIds)
// 			this.EvaluateItemIds = newItemIds
// 		} else {
// 			return nil, errors.New("项目任务更新失败，测试用例数量不能为0个")
// 		}
//
// 		if assetVersions, has := rawInfo.TryMap("asset_versions"); has {
// 			this.AssetVersions = *assetVersions
// 		}
// 		if err := mongoClient.Update(bson.M{"_id": this.Id}, this); err != nil {
// 			return nil, err
// 		}
// 		// 记录更新任务节点记录
// 		// 指派记录
// 		nodeInfo := qmap.QM{
// 			"project_id": rawInfo.MustString("project_id"),
// 			"task_id":    id,
// 			"status":     this.Status,
// 			"name":       common.PTS_SIGN_CREATE,
// 			"history": History{
// 				Result:    true,
// 				Operation: "任务修改",
// 				Comment:   "任务修改",
// 				OpId:      this.OpId,
// 				OpTime:    custom_util.GetCurrentMilliSecond(),
// 			},
// 		}
//
// 		// 检查当前阶段节点是否存在
// 		if _, err := new(EvaluateTaskNode).Update(nodeInfo); err != nil {
// 			return nil, errors.New("Update Task Node Error")
// 		}
//
// 		// 更新任务后，清除测试用例审核标记
// 		selector := bson.M{"evaluate_task_id": id}
// 		updateItem := bson.M{
// 			"$set": qmap.QM{
// 				"audit_status": common.EIAS_DEFAULT,
// 			},
// 		}
// 		if _, err := mongo.NewMgoSession(common.MC_EVALUATE_TASK_ITEM).UpdateAll(selector, updateItem); err != nil {
// 			return nil, err
// 		}
//
// 		// 将预绑定的正式绑定
// 		new(EvaluateItem).bindEvaluateTask(newItemIds)
//
// 		// 将删除的测试用例解绑
// 		if len(delItemIds) > 0 {
// 			new(EvaluateItem).UnbindEvaluateTask(delItemIds)
// 		}
//
// 	} else {
// 		return nil, errors.New("Item not found")
// 	}
// 	return this, nil
// }

func (this *EvaluateTask) GetAll(rawInfo qmap.QM, ctx context.Context) (*qmap.QM, error) {
	queryParams := rawInfo.String("query_params")
	mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_TASK).AddUrlQueryCondition(queryParams)
	projectSlice := GetProjectSlice()
	if res, err := mgoSession.GetPage(); err == nil {
		for index, item := range (*res)["list"].([]map[string]interface{}) {
			projectName := item["project_id"].(string)
			if projectSlice[item["project_id"].(string)] != nil {
				projectName = projectSlice[item["project_id"].(string)].(string)
			}
			(*res)["list"].([]map[string]interface{})[index]["project_name"] = projectName
			(*res)["list"].([]map[string]interface{})[index]["op_name"] = GetUserName(item["op_id"].(int), ctx)
			(*res)["list"].([]map[string]interface{})[index]["task_auditor_name"] = GetUserName(item["task_auditor_id"].(int), ctx)
			(*res)["list"].([]map[string]interface{})[index]["report_auditor_name"] = GetUserName(item["report_auditor_id"].(int), ctx)
			(*res)["list"].([]map[string]interface{})[index]["tester_name"] = GetUserName(item["tester_id"].(int), ctx)
		}
		return res, nil
	} else {
		return nil, err
	}

}

func GetUserName(uid int, ctx context.Context) string {
	// 通过id获取realname
	if rsp, err := new(mysql_model.SysUser).GetUserInfo(uid); err == nil {
		if str := rsp.String("realname"); str != "" {
			return str
		}
	}
	return "未知"
}

func (this *EvaluateTask) GetList(rawInfo qmap.QM, uid int, ctx context.Context) (*qmap.QM, error) {
	match := bson.M{
		"$or": []bson.M{{"op_id": uid}, {"task_auditor_id": uid}, {"report_auditor_id": uid}, {"tester_id": uid}},
	}
	if projectId, has := rawInfo.TryString("project_id"); has {
		match["project_id"] = projectId
	}
	if name, has := rawInfo.TryString("name"); has {
		match["name"] = bson.M{
			"$regex": name,
		}
	}
	if myStage := rawInfo.Int("my_stage"); myStage > 0 {
		match["current_operator_id"] = uid
	}
	if status, has := rawInfo.TryInt("status"); has {
		match["status"] = status
	}
	if auditStatus, has := rawInfo.TryInt("audit_status"); has {
		match["audit_status"] = auditStatus
	}
	if phase, has := rawInfo.TryInt("test_phase"); has {
		match["test_phase"] = phase
	}
	if isDeleted, has := rawInfo.TryInt("is_deleted"); has {
		match["is_deleted"] = isDeleted
	}
	submitTime := bson.M{}
	if startTime, has := rawInfo.TryInt("start_time"); has {
		submitTime["$gt"] = startTime
	}
	if endTime, has := rawInfo.TryInt("end_time"); has {
		submitTime["$lt"] = endTime
	}
	if len(submitTime) > 0 {
		match["submit_time"] = submitTime
	}
	Operations := []bson.M{
		{"$match": match},
	}

	fmt.Println(match)
	ormSession := mongo.NewMgoSession(common.MC_EVALUATE_TASK)
	if page, has := rawInfo.TryInt("page"); has {
		ormSession = ormSession.SetPage(page)
	}

	if limit, has := rawInfo.TryInt("limit"); has {
		ormSession = ormSession.SetLimit(limit)
	}

	list, err := ormSession.MATCHGetPage(Operations)

	if err == nil {
		projectSlice := GetProjectSlice()
		for index, item := range (*list)["list"].([]map[string]interface{}) {
			projectId := item["project_id"].(string)
			projectName := projectId
			if projectSlice[projectId] != nil {
				projectName = projectSlice[projectId].(string)
			}
			// 查找项目经理
			projectInfo, _ := mongo.NewMgoSession(common.MC_PROJECT).AddCondition(qmap.QM{"e__id": bson.ObjectIdHex(projectId)}).GetOne()
			(*list)["list"].([]map[string]interface{})[index]["manager_name"] = GetUserName((*projectInfo)["manager_id"].(int), ctx)
			(*list)["list"].([]map[string]interface{})[index]["project_name"] = projectName
			(*list)["list"].([]map[string]interface{})[index]["op_name"] = GetUserName(item["op_id"].(int), ctx)
			(*list)["list"].([]map[string]interface{})[index]["task_auditor_name"] = GetUserName(item["task_auditor_id"].(int), ctx)
			(*list)["list"].([]map[string]interface{})[index]["report_auditor_name"] = GetUserName(item["report_auditor_id"].(int), ctx)
			(*list)["list"].([]map[string]interface{})[index]["tester_name"] = GetUserName(item["tester_id"].(int), ctx)
		}
		return list, nil
	}
	return nil, err
}

func (this *EvaluateTask) GetTaskItemList(rawInfo qmap.QM, uid int, ctx context.Context) (*qmap.QM, error) {
	taskId := rawInfo.MustString("task_id")
	taskInfo, err := mongo.NewMgoSession(common.MC_EVALUATE_TASK).AddCondition(qmap.QM{"e__id": bson.ObjectIdHex(taskId)}).GetOne()
	if err != nil {
		return nil, errors.New("Task Not Found")
	}

	match := bson.M{}

	if name, has := rawInfo.TryString("name"); has {
		match["name"] = bson.M{
			"$regex": name,
		}
	}

	if testMethod, has := rawInfo.TryString("test_method"); has {
		match["test_method"] = testMethod
	}
	if autoTestLevel, has := rawInfo.TryString("auto_test_level"); has {
		match["auto_test_level"] = autoTestLevel
	}

	// 是否正式绑定
	isRelatived := true
	if (*taskInfo)["status"] == common.PTS_CREATE || (*taskInfo)["status"] == common.PTS_TASK_AUDIT {
		isRelatived = false
	}

	var ormSession *mongo.MgoOrmSession
	if isRelatived == false { // 创建阶段任务测试用例，从Item表中直接查询
		if ids, has := rawInfo.TrySlice("ids"); has {
			match["_id"] = bson.M{"$in": ids}
		}
		if itemId, has := rawInfo.TryString("item_id"); has {
			match["_id"] = itemId
		}
		if status, has := rawInfo.TryInt("test_status"); has {
			match["test_status"] = status
		}
		match["pre_bind"] = taskId
		ormSession = mongo.NewMgoSession(common.MC_EVALUATE_ITEM)
	} else { // 其他阶段，测试用例从task_item关系表中查询
		match["evaluate_task_id"] = taskId
		if ids, has := rawInfo.TrySlice("ids"); has {
			match["item_id"] = bson.M{"$in": ids}
		}
		if itemId, has := rawInfo.TryString("item_id"); has {
			match["item_id"] = itemId
		}
		if auditStatus, has := rawInfo.TryInt("audit_status"); has {
			match["audit_status"] = auditStatus
		}
		if status, has := rawInfo.TryInt("test_status"); has {
			match["status"] = status
		}
		ormSession = mongo.NewMgoSession(common.MC_EVALUATE_TASK_ITEM)
	}
	Operations := []bson.M{
		{"$match": match},
	}

	fmt.Println(Operations)
	if page, has := rawInfo.TryInt("page"); has {
		ormSession = ormSession.SetPage(page)
	}

	if limit, has := rawInfo.TryInt("limit"); has {
		ormSession = ormSession.SetLimit(limit)
	}
	list, err := ormSession.MATCHGetPage(Operations)

	for index, info := range (*list)["list"].([]map[string]interface{}) {
		// 需要追加测试人员
		(*list)["list"].([]map[string]interface{})[index]["tester_id"] = (*taskInfo)["tester_id"]
		(*list)["list"].([]map[string]interface{})[index]["tester_name"] = GetNameById((*taskInfo)["tester_id"].(int))

		if (*taskInfo)["status"] != common.PTS_CREATE && (*taskInfo)["status"] != common.PTS_TASK_AUDIT {
			testId := (*list)["list"].([]map[string]interface{})[index]["_id"]
			(*list)["list"].([]map[string]interface{})[index]["_id"] = (*list)["list"].([]map[string]interface{})[index]["item_id"]
			(*list)["list"].([]map[string]interface{})[index]["test_id"] = testId
		}

		if isRelatived == true {
			// 已经正式绑定的，接口追加相关的数据
			(*list)["list"].([]map[string]interface{})[index]["pre_bind"] = info["evaluate_task_id"]
			params := qmap.QM{
				"e__id": info["item_id"],
			}
			if itemInfo, err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).AddCondition(params).GetOne(); err == nil {
				(*list)["list"].([]map[string]interface{})[index]["vul_number"] = (*itemInfo)["vul_number"]
			}
			(*list)["list"].([]map[string]interface{})[index]["test_status"] = info["status"]
		} else {
			(*list)["list"].([]map[string]interface{})[index]["test_time"] = 0
			(*list)["list"].([]map[string]interface{})[index]["test_status"] = 0
		}
	}

	return list, err
}

func GetNameById(id int) string {
	// 通过id获取realname等用户信息
	if rsp, err := new(mysql_model.SysUser).GetUserInfo(id); err == nil {
		if realname := rsp.String("realname"); realname != "" {
			return realname
		} else {
			return rsp.String("username")
		}
	}
	return ""
}

func (this *EvaluateTask) GetOne(taskId string) (*qmap.QM, error) {
	param := qmap.QM{
		"e__id": bson.ObjectIdHex(taskId),
	}
	return mongo.NewMgoSession(common.MC_EVALUATE_TASK).AddCondition(param).GetOne()
}

func (this *EvaluateTask) GetPhase(projectId string) (*[]int, error) {
	result := []int{}
	groupOperations := []bson.M{
		{"$match": bson.M{"project_id": projectId}},
		{"$project": bson.M{"test_phase": 1}},
		{"$group": bson.M{"_id": bson.M{"test_phase": "$test_phase"}, "count": bson.M{"$sum": 1}}},
	}
	list, _ := mongo.NewMgoSession(common.MC_EVALUATE_TASK).QueryGet(groupOperations)
	for _, item := range *list {
		result = append(result, item["_id"].(map[string]interface{})["test_phase"].(int))
	}

	return &result, nil
}

func (this *EvaluateTask) GetTester(projectId string) (*[]interface{}, error) {
	result := []interface{}{}
	groupOperations := []bson.M{
		{"$match": bson.M{"project_id": projectId}},
		{"$project": bson.M{"tester_id": 1}},
		{"$group": bson.M{"_id": bson.M{"tester_id": "$tester_id"}, "count": bson.M{"$sum": 1}}},
	}
	list, _ := mongo.NewMgoSession(common.MC_EVALUATE_TASK).QueryGet(groupOperations)

	for _, item := range *list {
		uid := item["_id"].(map[string]interface{})["tester_id"]
		if uid != 0 && uid != "" {
			result = append(result, uid)
		}
	}
	return &result, nil
}

// func (this *EvaluateTask) Assign(req qmap.QM, opId int, userList map[int]string) (bool, error) {
// 	assignType := req.MustString("type")
// 	userId := req.MustInt("user_id")
// 	id := req.MustString("id")
// 	params := qmap.QM{
// 		"e__id": bson.ObjectIdHex(id),
// 	}
// 	mongoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params)
// 	if err := mongoClient.One(&this); err == nil {
// 		switch assignType {
// 		case "task_auditor_id": // 指派审核员
// 			// 检查任务阶段和权限
// 			if this.Status != common.PTS_CREATE || this.OpId != opId {
// 				return false, errors.New("Phase Error Or Not OpId")
// 			}
//
// 			statusData := qmap.QM{
// 				"project_id":  this.ProjectId,
// 				"id":          this.Id.Hex(),
// 				"status":      this.Status,
// 				"comment_tag": userList[userId],
// 				"result":      true,
// 				"op_id":       opId,
// 			}
// 			if _, err := this.ChangeStatusNode(statusData, "指派审核人"); err != nil {
// 				return false, err
// 			}
// 			this.Status = common.PTS_TASK_AUDIT
// 			this.CurrentOpId = userId
// 			this.AuditStatus = common.PTS_AUDIT_STATUS_NEW
// 			this.TaskAuditorId = userId
// 			this.TaskAuditorTime = 0
// 			// 指派任务审核人时，任务下所有的测试用例驳回状态的，改为默认
// 			mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).Session
// 			selector := bson.M{
// 				"pre_bind":     id,
// 				"audit_status": common.EIAS_REJECT,
// 			}
// 			data := qmap.QM{
// 				"audit_status": common.EIAS_DEFAULT,
// 			}
// 			if _, err := mgoSession.UpdateAll(selector, qmap.QM{"$set": data}); err != nil {
// 				return false, err
// 			}
//
// 		case "report_auditor_id":
// 			statusData := qmap.QM{
// 				"project_id":  this.ProjectId,
// 				"id":          this.Id.Hex(),
// 				"status":      this.Status,
// 				"comment_tag": userList[userId],
// 				"result":      true,
// 				"op_id":       opId,
// 			}
// 			if _, err := this.ChangeStatusNode(statusData, "指派报告审核人"); err != nil {
// 				return false, err
// 			}
// 			this.ReportAuditorId = userId
// 			this.CurrentOpId = userId
// 		case "tester_id": // 指派测试团队
// 			// 如果当前阶段是创建阶段，只有创建人可以指派测试团队，并且状态直接变为测试
// 			if this.Status == common.PTS_CREATE {
// 				if this.OpId != opId {
// 					return false, errors.New("创建阶段只有创建者可以指派测试人员！")
// 				}
// 			}
//
// 			statusData := qmap.QM{
// 				"project_id":  this.ProjectId,
// 				"id":          this.Id.Hex(),
// 				"status":      this.Status,
// 				"comment_tag": userList[userId],
// 				"result":      true,
// 				"op_id":       opId,
// 			}
// 			if _, err := this.ChangeStatusNode(statusData, "派发任务"); err != nil {
// 				return false, err
// 			}
// 			this.Status = common.PTS_TEST
// 			this.TesterId = userId
// 			this.ReportAuditorTime = 0
// 			this.CurrentOpId = userId
//
// 			// 指派测试团队时，将预绑定测试用例，建立正式任务关系
// 			new(EvaluateItem).WriteItemToTaskItem(id, opId)
// 		}
//
// 		if err := mongoClient.Update(bson.M{"_id": this.Id}, this); err != nil {
// 			return false, err
// 		}
// 		if this.Status == common.PTS_TEST {
// 			if _, err := new(Project).ChangeStatus(this.ProjectId, common.PS_TEST); err != nil {
// 				return false, err
// 			}
// 		}
//
// 		return true, nil
// 	}
// 	return false, errors.New("Item not found")
// }

func (this *EvaluateTask) ChangeStatusNode(rawInfo qmap.QM, comment string) (*EvaluateTask, error) {
	_id, _ := primitive.ObjectIDFromHex(rawInfo.MustString("id"))
	params := qmap.QM{
		"e__id": _id,
	}
	mongoClient := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_TASK, params)
	if err := mongoClient.One(&this); err == nil {
		status := rawInfo.MustInt("status")
		taskId := this.Id.Hex()
		name := TranslateStatusKey(status)
		commentTag := rawInfo.MustString("comment_tag")
		if commentTag != "" {
			comment += ": " + commentTag
		}

		// 指派记录
		nodeInfo := qmap.QM{
			"project_id": rawInfo.MustString("project_id"),
			"task_id":    taskId,
			"status":     status,
			"name":       name,
			"history": History{
				Result:    rawInfo["result"].(bool),
				Operation: comment,
				Comment:   comment,
				OpId:      rawInfo.MustInt("op_id"),
				OpTime:    custom_util.GetCurrentMilliSecond(),
			},
		}

		// 检查当前阶段节点是否存在
		if err := CheckNodeExist(taskId, name); err == nil {
			if _, err := new(EvaluateTaskNode).Create(nodeInfo, rawInfo.MustInt("op_id")); err != nil {
				return nil, errors.New("Create Task Node Error")
			}
		} else {
			if _, err := new(EvaluateTaskNode).Update(nodeInfo); err != nil {
				log.GetHttpLogLogger().Error(fmt.Sprintf("%s", err))
				return nil, errors.New("Update Task Node Error")
			}
		}
		return nil, err
	}
	return nil, errors.New("Item not found")
}

/**
 * 提交报告后，调用此方法，更新任务状态并记录节点日志
 */
func (this *EvaluateTask) changeStatusReportAudit(id string, opId int) (bool, error) {
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(id),
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params)
	if err := mongoClient.One(&this); err != nil {
		return false, errors.New("Item not found")
	}
	statusData := qmap.QM{
		"project_id":  this.ProjectId,
		"id":          this.Id.Hex(),
		"status":      this.Status,
		"comment_tag": "",
		"result":      true,
		"op_id":       opId,
	}
	if _, err := this.ChangeStatusNode(statusData, "提交报告审核"); err != nil {
		return false, err
	}
	this.Status = common.PTS_REPORT_AUDIT
	this.SubmitTime = custom_util.GetCurrentMilliSecond()
	if err := mongoClient.Update(bson.M{"_id": this.Id}, this); err != nil {
		return false, err
	}
	return true, nil
}

func (this *EvaluateTask) Submit(req qmap.QM, opId int) (bool, error) {
	taskId := req.MustString("id")
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(taskId),
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params)
	if err := mongoClient.One(&this); err == nil {
		// 检查当前任务状态，测试中
		if this.Status != common.PTS_TEST {
			return false, errors.New("只有测试状态下的任务，才可以提交")
		}

		// 检查任务下所有测试用例状态是否含有未开始和待补充的
		params := qmap.QM{
			"e_evaluate_task_id": taskId,
			"in_status": []int{
				common.TIS_PART_TEST_COMPLETE,
				common.TIS_READY,
			},
		}

		if _, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK_ITEM, params).GetOne(); err == nil {
			return false, errors.New("任务提交失败，有测试用例仍未测试完成")
		}

		// 将测试用例测试记录审核状态 驳回的 改为 默认
		selector := bson.M{
			"evaluate_task_id":    taskId,
			"record_audit_status": common.IRAS_REJECT,
		}
		updateItem := bson.M{
			"$set": qmap.QM{
				"record_audit_status": common.IRAS_DEFAULT,
			},
		}
		if _, err := mongo.NewMgoSession(common.MC_EVALUATE_TASK_ITEM).UpdateAll(selector, updateItem); err != nil {
			return false, err
		}

		// 更改任务状态
		return this.changeStatusReportAudit(taskId, opId)
	}
	return false, errors.New("Item not found")
}

func (this *EvaluateTask) Audit(req qmap.QM, opId int) (bool, error) {
	id := req.MustString("id")
	auditType := req.MustString("type")
	result := req.MustInt("result")
	comment := req.MustString("comment")

	params := qmap.QM{
		"e__id": bson.ObjectIdHex(id),
	}
	mongoClient := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params)
	if err := mongoClient.One(&this); err == nil {
		commentTag := "通过"
		resultTag := true
		if result == -1 {
			commentTag = "驳回"
			resultTag = false
		}
		switch auditType {
		case "task_audit": // 任务审核
			// 检查任务阶段和权限
			if this.Status != common.PTS_TASK_AUDIT {
				return false, errors.New("只有任务审核阶段可以审核任务")
			} else if this.OpId != opId && this.TaskAuditorId != opId { // 只有在创建者和指派的审核人可以审核
				return false, errors.New("有在创建者和指派的审核人可以审核!")
			}

			statusData := qmap.QM{
				"project_id":  this.ProjectId,
				"id":          this.Id.Hex(),
				"status":      this.Status,
				"comment_tag": commentTag,
				"result":      resultTag,
				"op_id":       opId,
			}
			if _, err := this.ChangeStatusNode(statusData, "任务审核"); err != nil {
				return false, err
			}

			if result == -1 {
				this.Status = common.PTS_CREATE
				// 任务审核驳回，重置任务审核人
				this.TaskAuditorId = 0
				this.TaskAuditorTime = 0
				// 任务驳回时，将未驳回的测试用例 改为通过
				mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).Session
				selector := bson.M{
					"pre_bind":     id,
					"audit_status": common.EIAS_DEFAULT,
				}
				data := qmap.QM{
					"audit_status": common.EIAS_ACCEPT,
				}
				if _, err := mgoSession.UpdateAll(selector, qmap.QM{"$set": data}); err != nil {
					return false, err
				}
			} else {
				// 任务审核通过，将测试用例状态改为通过
				mgoSession := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).Session
				data := qmap.QM{
					"audit_status": common.EIAS_ACCEPT,
				}
				if _, err := mgoSession.UpdateAll(bson.M{"pre_bind": id}, qmap.QM{"$set": data}); err != nil {
					return false, err
				}
			}
			this.AuditStatus = result
			this.TaskAuditorTime = custom_util.GetCurrentMilliSecond()
			this.AuditRecord = append(this.AuditRecord, qmap.QM{
				"result":  result,
				"comment": comment,
			})

		case "report_audit": // 报告审核
			// 检查任务阶段和权限
			if this.Status != common.PTS_REPORT_AUDIT {
				return false, errors.New("只有报告审核阶段可以审核报告")
			} else if this.OpId != opId && this.ReportAuditorId != opId { // 只有在创建者和指派的审核人可以审核
				return false, errors.New("只有创建者和报告审核员可以审核报告")
			}

			statusData := qmap.QM{
				"project_id":  this.ProjectId,
				"id":          this.Id.Hex(),
				"status":      this.Status,
				"comment_tag": commentTag,
				"result":      resultTag,
				"op_id":       opId,
			}
			if _, err := this.ChangeStatusNode(statusData, "报告审核"); err != nil {
				return false, err
			}

			if result == -1 {
				this.Status = common.PTS_TEST
				this.ReportAuditorId = 0
				this.ReportAuditorTime = 0
				this.SubmitTime = 0

				// 将所有用例记录状态置为：通过
				selector := bson.M{
					"evaluate_task_id": id,
				}
				updateItem := bson.M{
					"$set": qmap.QM{
						"record_audit_status": common.IRAS_ACCEPT,
					},
				}
				if _, err := mongo.NewMgoSession(common.MC_EVALUATE_TASK_ITEM).UpdateAll(selector, updateItem); err != nil {
					return false, err
				}

				// 将驳回的测试用例状态置为 待补充
				itemIds := req.MustSlice("item_ids")
				selector = bson.M{
					"evaluate_task_id": bson.M{"$eq": id},
					"item_id":          bson.M{"$in": itemIds},
				}
				updateItem = bson.M{
					"$set": qmap.QM{
						"status":              common.TIS_PART_TEST_COMPLETE,
						"record_audit_status": common.IRAS_REJECT,
					},
				}
				if _, err := mongo.NewMgoSession(common.MC_EVALUATE_TASK_ITEM).UpdateAll(selector, updateItem); err != nil {
					return false, err
				}

				// 更新Item主表中 test_status 待补充
				if err := new(EvaluateItem).ChangeTestStatus(itemIds, common.TIS_PART_TEST_COMPLETE); err != nil {
					return false, errors.New("更新Item主表信息失败")
				}

			} else {
				// 报告审核通过，任务完成
				statusData := qmap.QM{
					"project_id":  this.ProjectId,
					"id":          this.Id.Hex(),
					"status":      common.PTS_FINISH,
					"comment_tag": commentTag,
					"result":      resultTag,
					"op_id":       opId,
				}
				if _, err := this.ChangeStatusNode(statusData, "报告审核"); err != nil {
					return false, err
				}
				// 将测试用例 测试状态 改为完成
				selector := bson.M{
					"evaluate_task_id": id,
				}
				updateItem := bson.M{
					"$set": qmap.QM{
						"status":              common.TIS_COMPLETE,
						"record_audit_status": common.IRAS_ACCEPT,
					},
				}
				if _, err := mongo.NewMgoSession(common.MC_EVALUATE_TASK_ITEM).UpdateAll(selector, updateItem); err != nil {
					return false, err
				}
				// 将原始测试用例表中，测试用例改为可绑定任务
				itemIds := new(EvaluateTaskItem).GetItemIdsByTaskId(id)
				selector = bson.M{
					"_id": bson.M{"$in": itemIds},
				}
				updateItem = bson.M{
					"$set": qmap.QM{
						"status":       common.EIS_FREE,
						"pre_bind":     "",
						"last_task_id": id,
						"test_status":  common.TIS_COMPLETE,
					},
				}
				if _, err := mongo.NewMgoSession(common.MC_EVALUATE_ITEM).UpdateAll(selector, updateItem); err != nil {
					return false, err
				}

				//将任务的漏洞，更新到ITEM的漏洞主表中
				// new(EvaluateItem).UpdateVuls(id, itemIds)
				// 这块采用新的逻辑
				new(EvaluateItem).UpdateVuls2(id, itemIds, opId)

				this.Status = common.PTS_FINISH
				this.CurrentOpId = 0
				this.UpdateTime = custom_util.GetCurrentMilliSecond()
				this.ReportAuditorTime = custom_util.GetCurrentMilliSecond()
				// 如果未分配报告审核人，创建者直接审核的时候，将报告审核人标记为当前审核人ID
				if this.ReportAuditorId == 0 {
					this.ReportAuditorId = opId
				}

			}

			this.ReportAuditStatus = result
			this.ReportAuditRecord = append(this.ReportAuditRecord, qmap.QM{
				"result":  result,
				"comment": comment,
			})
		}

		if err := mongoClient.Update(bson.M{"_id": this.Id}, this); err != nil {
			return false, err
		}

		return true, nil
	}
	return false, errors.New("Item not found")
}

// 判断任务是否存在
func checkTaskExist(taskId string) error {
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(taskId),
	}
	if _, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params).GetOne(); err == nil {
		return errors.New(fmt.Sprintf("任务: %s 不存在！", taskId))
	}
	return nil
}

// 判断任务是否存在
func checkProjectTaskExist(projectId, taskId string) error {
	params := qmap.QM{
		"e__id":        bson.ObjectIdHex(taskId),
		"e_project_id": projectId,
	}
	if _, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params).GetOne(); err == nil {
		return errors.New(fmt.Sprintf("项目任务: %s 不存在！或与 %s 不匹配", taskId, projectId))
	}
	return nil
}

func (this *EvaluateTask) GetAssetVersion(id, assetId string) string {
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(id),
	}
	if err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params).One(this); err == nil {
		if version, has := this.AssetVersions.TryString(assetId); has {
			return version
		}
	}
	return ""
}

func TranslateStatusKey(status int) string {
	phase := map[int]string{
		1: common.PTS_SIGN_CREATE,
		2: common.PTS_SIGN_TASK_AUDIT,
		3: common.PTS_SIGN_TEST,
		4: common.PTS_SIGN_REPORT_AUDIT,
		5: common.PTS_SIGN_FINISH,
	}
	if phase[status] != "" {
		return phase[status]
	}
	return ""
}

/**
 * @Description: 统计项目中的任务数量
 * @param projectIds
 * @param status
 * @return int
 */
func (this *EvaluateTask) CountProjectTask(projectIds, status []string) int {
	params := qmap.QM{
		"in_project_id": projectIds,
		"in_status":     status,
	}
	if count, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params).Count(); err == nil {
		return count
	} else {
		return 0
	}
}

/**
 * @Description:统计项目包含的测试人员数量
 * @param projectIds
 * @return int
 */
func (this *EvaluateTask) CountProjectsTaskTester(projectIds []string) int {
	match := bson.M{
		"project_id": bson.M{"$in": projectIds},
		"tester_id":  bson.M{"$gt": 0},
	}
	group := bson.M{"_id": bson.M{"utester_id": "$member_ids"}}

	operations := []bson.M{
		{"$match": match},
		{"$group": group},
		{"$group": bson.M{"_id": nil, "count": bson.M{"$sum": 1}}},
	}
	incident := mongo.NewMgoSession(common.MC_EVALUATE_TASK).Session
	resp := []bson.M{}
	if err := incident.Pipe(operations).All(&resp); err == nil {
		if len(resp) > 0 {
			data := resp[0]
			return data["count"].(int)
		}
	} else {
		panic(err)
	}
	return 0
}

/**
 * @Description:统计某个用户当前待操作的项目任务数量
 * @param userId
 * @return int
 */
func (this *EvaluateTask) CountBacklogTask(userId int64) int {
	params := qmap.QM{
		"e_current_op_id": userId,
	}
	if count, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params).Count(); err == nil {
		return count
	} else {
		return 0
	}
}

/**
 * @Description:统计最近某一段时间内，完成的项目任务数量
 * @param userId
 * @return int
 */
func (this *EvaluateTask) CountCompleteTask(projectIds []string, startTime int64) int {
	params := qmap.QM{
		"in_project_id":   projectIds,
		"e_status":        common.PTS_FINISH,
		"gte_update_time": startTime,
	}
	if count, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_TASK, params).Count(); err == nil {
		return count
	} else {
		return 0
	}
}

// -------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------
// -------------------------------------------------------------------------------------

// 指派审核人或测试团队
// req参数：
//
//	前端传的json体内容 {"id":"60f92b1e898f98590df9233b","type":"task_auditor_id","user_id":24}
//	id为任务表记录的_id， type是操作类型， user_id是设置的用户id
func LogicEvaluateTaskAssign(req qmap.QM, opId int, userList map[int]string) (bool, error) {
	assignType := req.MustString("type")
	userId := req.MustInt("user_id")
	id := req.MustString("id")
	_id, _ := primitive.ObjectIDFromHex(id)
	params := qmap.QM{
		"e__id": _id,
	}

	this := EvaluateTask{}
	widget := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_TASK, params)
	if err := widget.One(&this); err == nil {
		switch assignType {
		case "task_auditor_id": // 指派审核员
			// 检查任务阶段和权限
			if this.Status != common.PTS_CREATE || this.OpId != opId {
				return false, errors.New("Phase Error Or Not OpId")
			}

			statusData := qmap.QM{
				"project_id":  this.ProjectId,
				"id":          this.Id.Hex(),
				"status":      this.Status,
				"comment_tag": userList[userId],
				"result":      true,
				"op_id":       opId,
			}
			if _, err := this.ChangeStatusNode(statusData, "指派审核人"); err != nil {
				return false, err
			}

			this.Status = common.PTS_TASK_AUDIT
			this.CurrentOpId = userId
			this.AuditStatus = common.PTS_AUDIT_STATUS_NEW
			this.TaskAuditorId = userId
			this.TaskAuditorTime = 0
			// 指派任务审核人时，任务下所有的测试用例驳回状态的，改为默认
			coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_ITEM)
			selector := bson.M{
				"pre_bind":     id,
				"audit_status": common.EIAS_REJECT,
			}
			data := qmap.QM{
				"audit_status": common.EIAS_DEFAULT,
			}
			if _, err := coll.UpdateMany(context.Background(), selector, qmap.QM{"$set": data}); err != nil {
				return false, err
			}

		case "report_auditor_id":
			statusData := qmap.QM{
				"project_id":  this.ProjectId,
				"id":          this.Id.Hex(),
				"status":      this.Status,
				"comment_tag": userList[userId],
				"result":      true,
				"op_id":       opId,
			}
			if _, err := this.ChangeStatusNode(statusData, "指派报告审核人"); err != nil {
				return false, err
			}
			this.ReportAuditorId = userId
			this.CurrentOpId = userId
		case "tester_id": // 指派测试团队
			// 如果当前阶段是创建阶段，只有创建人可以指派测试团队，并且状态直接变为测试
			if this.Status == common.PTS_CREATE {
				if this.OpId != opId {
					return false, errors.New("创建阶段只有创建者可以指派测试人员！")
				}
			}

			statusData := qmap.QM{
				"project_id":  this.ProjectId,
				"id":          this.Id.Hex(),
				"status":      this.Status,
				"comment_tag": userList[userId],
				"result":      true,
				"op_id":       opId,
			}
			if _, err := this.ChangeStatusNode(statusData, "派发任务"); err != nil {
				return false, err
			}
			this.Status = common.PTS_TEST
			this.TesterId = userId
			this.ReportAuditorTime = 0
			this.CurrentOpId = userId

			// 指派测试团队时，将预绑定测试用例，建立正式任务关系
			new(EvaluateItem).WriteItemToTaskItem(id, opId)
		}

		coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_TASK)
		if _, err := coll.UpdateOne(context.Background(), bson.M{"_id": this.Id}, bson.M{"$set": this}); err != nil {
			return false, err
		}

		if this.Status == common.PTS_TEST {
			if _, err := new(Project).ChangeStatus(this.ProjectId, common.PS_TEST); err != nil {
				return false, err
			}
		}

		return true, nil
	}
	return false, errors.New("Item not found")
}

// 新增任务的节点记录（创建任务时的节点日志）
func EvaluateTaskCreateTaskNodeCreate(ctx context.Context, task *EvaluateTask) (*EvaluateTaskNode, error) {
	history := &History{
		Result:    true,
		Operation: "任务创建",
		Comment:   "任务创建",
		OpId:      task.OpId,
		OpTime:    task.CreateTime,
	}
	taskNodePtr, err := EvaluateTaskNodeCreate(ctx, task, history, common.PTS_SIGN_CREATE)
	return taskNodePtr, err
}

// 新增任务的节点记录（创建任务时的节点日志）
func EvaluateTaskCreateTaskNodeUpdate(ctx context.Context, task *EvaluateTask) (*EvaluateTaskNode, error) {
	history := &History{
		Result:    true,
		Operation: "任务修改",
		Comment:   "任务修改",
		OpId:      task.OpId,
		OpTime:    custom_util.GetCurrentMilliSecond(),
	}
	taskNodePtr, err := EvaluateTaskNodeUpdate(ctx, task, history, common.PTS_SIGN_TASK_AUDIT)
	return taskNodePtr, err
}

// 新增任务记录
func EvaluateTaskInsert(ctx context.Context, modelPrt *EvaluateTask) (*primitive.ObjectID, error) {
	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_TASK)
	result, err := coll.InsertOne(ctx, modelPrt)
	if err != nil {
		return nil, err
	}
	a := result.InsertedID.(primitive.ObjectID)
	return &a, nil
}

// 更新任务记录
func EvaluateTaskUpdate(ctx context.Context, id interface{}, task *EvaluateTask) error {
	coll := orm_mongo.GetDefaultMongoDatabase().Collection(common.MC_EVALUATE_TASK)
	_, err := coll.UpdateByID(ctx, id, task)
	return err
}
