package logic

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/orm_mongo"
	"skygo_detection/mongo_model"
)

// 新增任务
// 1、检查，
//
//	任务所属的项目必须存在，且当前操作用户必须是管理员
//	参数中必须要有evaluate_item_ids，即测试用例不能为空
//	任务名称不能重复，检查集合evaluate_task
//	检测项目的测试用例数据，即集合evaluate_item, project_id对应所有的测试用例，检查这次选中的测试用例，其status必须是0，才能创建任务
//
// 2、插入一条任务evaluate_task记录
// 3、批量修改测试用例数据，即集合evaluate_item,把这次涉及到的测试用例都改为status =1，表示以及被任务占用了。
// 4. 创建任务阶段数据，即集合evaluate_task_node。
func EvaluateTaskCreate(rawInfo qmap.QM, opId int64) (*mongo_model.EvaluateTask, error) {
	// 测试任务所属项目必须存在
	projectId := rawInfo.MustString("project_id")
	name := rawInfo.MustString("name")
	_id, _ := primitive.ObjectIDFromHex(projectId)
	params := qmap.QM{
		"e__id": _id,
	}
	if err := mongo_model.CheckIsProjectManager(projectId, opId); err != nil {
		return nil, err
	}

	evaluateItemIds := rawInfo.SliceString("evaluate_item_ids")
	if len(evaluateItemIds) == 0 {
		return nil, errors.New("测试用例不能为空")
	}

	// 任务名称不能重复
	params = qmap.QM{
		"e_project_id": projectId,
		"e_name":       name,
	}
	if info, _ := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_TASK, params).Get(); info["_id"] != nil {
		return nil, errors.New("该任务名称已被使用")
	}

	// 新建项目任务中关联的测试用例必须还没有与其他项目任务关联: 判断status=0可绑定，status=1不可绑定
	for _, item := range evaluateItemIds {
		params = qmap.QM{
			"e__id":        item,
			"e_project_id": projectId,
		}
		if temp, err := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_ITEM, params).Get(); err != nil || temp.Int("status") != 0 {
			return nil, errors.New(fmt.Sprintf("项目任务创建失败:测试用例%s无法添加到该项目任务中", item))
		}
	}

	model := mongo_model.EvaluateTask{}
	model.Id = primitive.NewObjectID()
	model.ProjectId = projectId
	model.Name = name
	model.AssetVersions = rawInfo.Map("asset_versions")
	model.EvaluateItemIds = evaluateItemIds
	model.Status = common.PTS_CREATE
	model.IsDeleted = common.PSD_DEFAULT
	model.AuditStatus = common.PTS_AUDIT_STATUS_NEW
	model.ReportAuditStatus = common.PTS_AUDIT_STATUS_NEW
	model.TestPhase = rawInfo.Int("test_phase")
	model.UpdateTime = custom_util.GetCurrentMilliSecond()
	model.CreateTime = custom_util.GetCurrentMilliSecond()
	model.SubmitTime = 0
	model.TaskAuditorTime = 0
	model.ReportAuditorTime = 0
	model.OpId = int(opId)
	model.CurrentOpId = int(opId)

	// todo 加事务

	ctx := context.Background()

	// 插入任务记录
	if _, err := mongo_model.EvaluateTaskInsert(ctx, &model); err != nil {
		return nil, err
	}

	// 项目任务创建成功后将测试用例与项目任务绑定(绑定后，测试用例将不能与其他项目任务绑定)
	if err := mongo_model.EvaluateItemBindEvaluateTask(ctx, evaluateItemIds, model.Id.Hex(), common.NOT_PREBIND); err != nil {
		return nil, err
	}

	// 创建阶段记录
	if _, err := mongo_model.EvaluateTaskCreateTaskNodeCreate(ctx, &model); err != nil {
		return nil, err
	}
	return &model, nil
}

// 修改任务
func EvaluateTaskUpdate(idHex string, rawInfo qmap.QM) (*mongo_model.EvaluateTask, error) {
	_id, _ := primitive.ObjectIDFromHex(idHex)
	params := qmap.QM{
		"e__id": _id,
	}

	// 查询任务记录
	this := mongo_model.EvaluateTask{}
	coll := orm_mongo.NewWidgetWithParams(common.MC_EVALUATE_TASK, params)
	if err := coll.One(&this); err != nil {
		return nil, errors.New("item not found")
	}

	// 解析要更新的数据
	if val, has := rawInfo.TryString("name"); has {
		this.Name = val
	}
	itemIds, has := rawInfo.TrySlice("evaluate_item_ids")

	delItemIds := []string{}
	newItemIds := rawInfo.SliceString("evaluate_item_ids")

	if has && len(itemIds) > 0 {
		delItemIds = custom_util.DifferenceDel(this.EvaluateItemIds, newItemIds)
		this.EvaluateItemIds = newItemIds
	} else {
		return nil, errors.New("项目任务更新失败，测试用例数量不能为0个")
	}

	if assetVersions, has := rawInfo.TryMap("asset_versions"); has {
		this.AssetVersions = *assetVersions
	}

	// todo 事务
	ctx := context.Background()

	err := mongo_model.EvaluateTaskUpdate(ctx, this.Id, &this)
	if err != nil {
		return nil, err
	}

	// 记录更新任务节点记录
	// 指派记录
	// 检查当前阶段节点是否存在
	_, err = mongo_model.EvaluateTaskCreateTaskNodeUpdate(ctx, &this)
	if err != nil {
		return nil, errors.New("Update Task Node Error")
	}

	// 更新任务后，清除测试用例审核标记
	err = mongo_model.EvaluateTaskItemUpdateAuditStatus(ctx, idHex, common.EIAS_DEFAULT)
	if err != nil {
		return nil, err
	}

	// 将预绑定的正式绑定
	err = mongo_model.EvaluateItem_bindEvaluateTask(ctx, newItemIds)
	if err != nil {
		return nil, err
	}

	// 将删除的测试用例解绑
	if len(delItemIds) > 0 {
		err := mongo_model.EvaluateItem_UnbindEvaluateTask(ctx, delItemIds)
		if err != nil {
			return nil, err
		}
	}

	return &this, nil
}
