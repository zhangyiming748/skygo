package mongo_model

import (
	"fmt"

	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
)

type EvaluateVulInfo struct {
	ID                  bson.ObjectId `bson:"_id,omitempty" json:"_id"`
	TaskId              string        `bson:"task_id" json:"task_id"`                             //任务id
	CveId               string        `bson:"cve_id" json:"cve_id"`                               //cve ID编号
	CveType             string        `bson:"cve_type" json:"cve_type"`                           //漏洞类型
	FixStatus           string        `bson:"fix_status" json:"fix_status"`                       //漏洞状态
	InvolveModule       string        `bson:"involve_module" json:"involve_module"`               //漏洞模块
	GoogleSeverityLevel string        `bson:"google_severity_level" json:"google_severity_level"` //漏洞级别
	DateExposure        string        `bson:"date_exposure" json:"date_exposure"`                 //披露时间
	DateBulletin        string        `bson:"date_bulletin" json:"date_bulletin"`                 //发布时间
	Sketch              string        `bson:"sketch" json:"sketch"`                               //漏洞简述
	Description         string        `bson:"description" json:"description"`                     //漏洞详细描述
	SearchContent       string        `bson:"search_content" json:"search_content"`               //搜索字段
}

func (this *EvaluateVulInfo) Create(rawInfo qmap.QM) (*EvaluateVulInfo, error) {
	searchContent := fmt.Sprintf("%s_%s", rawInfo.String("cve_id"), rawInfo.String("sketch"))
	rawInfo["search_content"] = searchContent
	if err := mongo.NewMgoSession(common.MC_EVALUATE_VUL_INFO).Insert(rawInfo); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateVulInfo) BulkDelete(rawIds []string) (*qmap.QM, error) {
	// 删除 测试项
	effectNum := 0
	ids := []string{}
	for _, id := range rawIds {
		ids = append(ids, id)
	}
	if len(ids) > 0 {
		match := bson.M{
			"task_id": bson.M{"$in": ids},
		}
		if changeInfo, err := mongo.NewMgoSession(common.MC_EVALUATE_VUL_INFO).RemoveAll(match); err == nil {
			effectNum = changeInfo.Removed
			// 根据item_id删除 测试项里的漏洞
			new(EvaluateVulnerability).BulkDeleteByItemIds(rawIds)
		} else {
			return nil, err
		}
	}
	return &qmap.QM{"number": effectNum}, nil
}

func (this *EvaluateVulInfo) GetOne(taskId string) (*qmap.QM, error) {
	params := qmap.QM{
		"e_task_id": taskId,
	}
	return mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_INFO, params).GetOne()
}
