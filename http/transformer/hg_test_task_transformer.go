package transformer

import (
	"time"

	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/logic"
	"skygo_detection/mongo_model"
)

type HgTestTaskTransformer struct {
	mongo.BaseTransformer
}

func (h *HgTestTaskTransformer) AdditionalFields(data qmap.QM) qmap.QM {
	countInfo := new(logic.HgTestTaskLogic).GetCountInfo(data)
	data["test_case_count"] = countInfo

	lastConnectTime, _ := data.TryInt("last_connect_time")
	if lastConnectTime == 0 {
		data["connect_status"] = mongo_model.HgTestTaskConnectStatusNever
	} else {
		// 5s内做为已连接
		if (time.Now().UnixNano()/1e6 - int64(lastConnectTime)) < 10000 {
			data["connect_status"] = mongo_model.HgTestTaskConnectStatusYes
		} else {
			data["connect_status"] = mongo_model.HgTestTaskConnectStatusNo
		}
	}
	return data
}
