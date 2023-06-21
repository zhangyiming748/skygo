package transformer

import (
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/mysql_model"
)

type AssetVehicleTransformer struct {
	orm.Transformer
}

func (h *AssetVehicleTransformer) ModifyItem(data qmap.QM) qmap.QM {
	createUserId := data.MustInt("create_user_id")

	userModel, err := mysql_model.SysUserFindById(createUserId)
	if err == nil {
		data["create_user_name"] = userModel.Realname
	} else {
		data["create_user_name"] = "-"
	}

	return data
}
