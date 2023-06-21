package transformer

import (
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/mysql_model"
)

type AssetTestPieceTransformer struct {
	orm.Transformer
}

func (h *AssetTestPieceTransformer) ModifyItem(data qmap.QM) qmap.QM {
	// 获取测试件id
	id := data.MustInt("id")
	assetVehicleId := data.MustInt("asset_vehicle_id")

	// 车型信息
	vehicleModel := mysql_model.AssetVehicle{}
	mysql.FindById(assetVehicleId, &vehicleModel)
	data["brand"] = vehicleModel.Brand // 车型品牌
	data["code"] = vehicleModel.Code   // 车型代号

	// 查询最新版本
	vModel, _ := mysql_model.AssetTestPieceVersionFindLatest(id)
	data["version"] = vModel.Version             // 测试件最新版本
	data["create_user_id"] = vModel.CreateUserId // 创建人
	// data["update_time"] = vModel.UpdateTime // 更新时间
	tmp := mysql_model.AssetTestPieceVersionFindAll(id)
	data["versions_history"] = tmp // 测试件里的所有版本
	return data
}

type AssetTestPieceDetailTransformer struct {
	orm.Transformer
}

func (h *AssetTestPieceDetailTransformer) ModifyItem(data qmap.QM) qmap.QM {
	// 获取测试件id
	id := data.MustInt("id")
	assetVehicleId := data.MustInt("asset_vehicle_id")

	// 车型信息
	vehicleModel := mysql_model.AssetVehicle{}
	mysql.FindById(assetVehicleId, &vehicleModel)
	data["brand"] = vehicleModel.Brand // 车型品牌
	data["code"] = vehicleModel.Code   // 车型代号

	// 查询最新版本
	vModel, _ := mysql_model.AssetTestPieceVersionFindLatest(id)
	data["version"] = vModel.Version             // 测试件最新版本
	data["create_user_id"] = vModel.CreateUserId // 创建人
	// data["update_time"] = vModel.UpdateTime // 更新时间

	// 查询组件的所有版本记录列表，取版本名称
	vModels := make([]mysql_model.AssetTestPieceVersion, 0)
	s := mysql.GetSession().Where("asset_test_piece_id = ?", id).OrderBy("create_time desc").Limit(20)
	s.Find(&vModels)

	versionNames := make([]string, 0)
	versionIds := make([]int, 0)
	for _, v := range vModels {
		versionNames = append(versionNames, v.Version)
		versionIds = append(versionIds, v.Id)
	}

	data["version_names"] = versionNames
	data["version_ids"] = versionIds

	return data
}
