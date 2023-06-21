package transformer

import (
	"encoding/json"
	"fmt"

	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/mongo_model"
	"skygo_detection/mongo_model_tmp"
	"skygo_detection/mysql_model"
)

type TaskTransformer struct {
	orm.Transformer
}

func (h *TaskTransformer) ModifyItem(data qmap.QM) qmap.QM {
	// 用户名
	createUserId := data.MustInt("create_user_id")
	userModel, err := mysql_model.SysUserFindById(createUserId)
	if err == nil {
		data["create_user_name"] = userModel.Realname
	} else {
		data["create_user_name"] = "-"
	}

	// 前端需要展示车型信息，用 “brand/code”
	assetVehicleId := data.MustInt("asset_vehicle_id")
	assetVehicleModel, has, _ := mysql_model.AssetVehicleFindById(assetVehicleId)
	if has {
		data["asset_vehicle_code"] = fmt.Sprintf("%s/%s", assetVehicleModel.Brand, assetVehicleModel.Code)
		data["asset_vehicle_brand"] = assetVehicleModel.Brand
		data["asset_vehicle_brand_code"] = assetVehicleModel.Code
	} else {
		data["asset_vehicle_code"] = "未查到车型"
		data["asset_vehicle_brand"] = ""
		data["asset_vehicle_brand_code"] = ""
	}

	// 测试件名称
	pieceId := data.MustInt("piece_id")
	pieceModel, has, _ := mysql_model.AssetTestPieceFindById(pieceId)
	if has {
		data["piece_name"] = pieceModel.Name
	} else {
		data["piece_name"] = ""
	}
	// 测试件版本
	if pieceVerion, err := new(mysql_model.AssetTestPieceVersion).FindById(data.Int("piece_version_id")); err == nil {
		data["piece_version"] = pieceVerion.Version
	} else {
		data["piece_version"] = ""
	}

	// 场景名称
	scenarioId := data.MustInt("scenario_id")
	scenarioModel, has, _ := mysql_model.KnowledgeScenarioFindById(scenarioId)
	if has {
		data["scenario_name"] = scenarioModel.Name
	} else {
		data["scenario_name"] = ""
	}

	// 如果这个父任务是车机漏洞扫描的话，需要添加漏扫的任务信息
	if data.String("category") == common.TOOL_VUL_SCANNER_NAME {
		parentTaskId := data.MustInt("id")
		vulTask, _ := new(mongo_model_tmp.EvaluateVulTask).GetOneByParentId(parentTaskId)
		data["vul_task"] = vulTask
	}
	return data
}

// 任务详情页的转换逻辑
type TaskOneTransformer struct {
	orm.Transformer
}

func (h *TaskOneTransformer) ModifyItem(data qmap.QM) qmap.QM {
	// 用户名
	createUserId := data.MustInt("create_user_id")
	userModel, err := mysql_model.SysUserFindById(createUserId)
	if err == nil {
		data["create_user_name"] = userModel.Realname
	} else {
		data["create_user_name"] = "-"
	}

	// 前端需要展示车型信息，用 “brand/code”
	assetVehicleId := data.MustInt("asset_vehicle_id")
	assetVehicleModel, has, _ := mysql_model.AssetVehicleFindById(assetVehicleId)
	if has {
		data["asset_vehicle_brand"] = assetVehicleModel.Brand
	} else {
		data["asset_vehicle_code"] = assetVehicleModel.Code
	}

	// 测试件名称
	pieceId := data.MustInt("piece_id")
	pieceModel, has, _ := mysql_model.AssetTestPieceFindById(pieceId)
	if has {
		data["piece_name"] = pieceModel.Name
	} else {
		data["piece_name"] = ""
	}

	// 场景名称
	scenarioId := data.MustInt("scenario_id")
	scenarioModel, has, _ := mysql_model.KnowledgeScenarioFindById(scenarioId)
	if has {
		data["scenario_name"] = scenarioModel.Name
	} else {
		data["scenario_name"] = ""
	}
	return data
}

// 任务用例里工具的名称
type TaskToolTransformer struct {
	orm.Transformer
}

func (h *TaskToolTransformer) ModifyItem(data qmap.QM) qmap.QM {
	// 场景名称
	fileId := data.MustString("file_id")
	if fileId != "" {
		rawInfo := qmap.QM{
			"id": fileId,
		}
		tool, err := new(mongo_model.ToolData).GetOne(&rawInfo)
		if err != nil {
			data["tool_name"] = "不存在这个工具"
		} else {
			if val, has := tool.TryString("tool_name"); has {
				data["tool_name"] = val
			} else {
				// todo 测试工具调试中
				data["tool_name"] = "测试工具调试中"
			}
		}
		return data
	}
	blockList := data["test_param"]
	if blockList != "" {
		ttmp := make([]qmap.QM, 0)
		blockListString := blockList.(string)
		err := json.Unmarshal([]byte(blockListString), &ttmp)
		if err != nil {
			fmt.Println(err)
		}
		data["test_param"] = ttmp
	}
	data["tool_name"] = "无测试工具"
	return data
}
