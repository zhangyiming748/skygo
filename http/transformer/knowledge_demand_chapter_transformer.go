package transformer

import (
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/mysql_model"
)

type KnowledgeDemandChapterTransformer struct {
	orm.Transformer
}

func (h *KnowledgeDemandChapterTransformer) ModifyItem(data qmap.QM) qmap.QM {
	parentId := data.MustInt("parent_id")

	model, has := mysql_model.KnowledgeDemandChapterFindById(parentId)
	if has {
		data["parent_code"] = model.Code
	} else {
		data["parent_code"] = ""
	}

	return data
}
