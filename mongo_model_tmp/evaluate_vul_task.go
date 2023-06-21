package mongo_model_tmp

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
)

// 测试项
type EvaluateVulTask struct {
	ID            bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Name          string        `bson:"name" json:"name"`                     // 任务名称
	TaskID        string        `bson:"task_id" json:"task_id"`               // 任务ID
	VulScannerID  string        `bson:"vul_scanner_id" json:"vul_scanner_id"` // 漏洞检测id
	Status        int           `bson:"status" json:"status"`                 // 状态（4未展示，0未开始，1测试中，2测试完成）
	TestTime      int           `bson:"test_time" json:"test_time"`           // 测试时间
	CreateTime    int           `bson:"create_time" json:"create_time"`       // 创建时间
	SearchContent string        `bson:"search_content" json:"search_content"` // 搜索字段
	ParentTaskId  int           `bson:"parent_task_id" json:"parent_task_id"` // 父类任务ID
}

func (this *EvaluateVulTask) GetAll(queryParams string) (*qmap.QM, error) {
	params := qmap.QM{
		"ne_status": common.VUL_UNSHOW,
	}
	mgoSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_TASK, params).AddUrlQueryCondition(queryParams)
	if res, err := mgoSession.GetPage(); err == nil {
		return res, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateVulTask) CheckNameExist(name string) bool {
	params := qmap.QM{
		"e_name": name,
	}
	mgoSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_TASK, params)
	if _, err := mgoSession.GetOne(); err == nil {
		return true
	}
	return false
}

func (this *EvaluateVulTask) Create(taskId string, rawInfo qmap.QM) (*EvaluateVulTask, error) {
	name := rawInfo.MustString("name")
	// 检查name是否存在
	if this.CheckNameExist(name) {
		return nil, errors.New("名称已存在")
	}

	this.ID = bson.NewObjectId()
	this.Name = rawInfo.String("name")
	if taskId != "" {
		this.TaskID = taskId
	} else {
		this.TaskID = this.getTaskId()
	}
	this.Status = common.VUL_UNSTART
	this.TestTime = 0
	this.CreateTime = int(custom_util.GetCurrentMilliSecond())
	this.SearchContent = fmt.Sprintf("%s_%s", this.Name, this.TaskID)
	this.ParentTaskId = rawInfo.MustInt("parent_task_id")
	if err := mongo.NewMgoSession(common.MC_EVALUATE_VUL_TASK).Insert(this); err == nil {
		return this, nil
	} else {
		return nil, err
	}
}

func (this *EvaluateVulTask) BulkDelete(rawIds []string) (*qmap.QM, error) {
	// 删除 测试项
	effectNum := 0
	ids := []bson.ObjectId{}
	for _, id := range rawIds {
		ids = append(ids, bson.ObjectIdHex(id))
	}
	if len(ids) > 0 {
		match := bson.M{
			"_id": bson.M{"$in": ids},
		}
		if changeInfo, err := mongo.NewMgoSession(common.MC_EVALUATE_VUL_TASK).RemoveAll(match); err == nil {
			effectNum = changeInfo.Removed
		} else {
			return nil, err
		}
	}
	return &qmap.QM{"number": effectNum}, nil
}

func (this *EvaluateVulTask) getTaskId() string {
	now := int(time.Now().UnixNano() / 1000)
	str := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K", "L", "M", "N", "O", "P", "Q", "R", "S", "T", "U", "V", "W", "X", "Y", "Z"}
	assetID := ""
	var remainder int
	var remainderStr string
	for now != 0 {
		remainder = now % 36
		if remainder < 36 && remainder > 9 {
			remainderStr = str[remainder]
		} else {
			remainderStr = strconv.Itoa(remainder)
		}
		assetID = remainderStr + assetID
		now = now / 36
	}
	if len(assetID) > 8 {
		rs := []rune(assetID)
		assetID = string(rs[:8])
	}

	return assetID
}

func (this *EvaluateVulTask) GetOneByParentId(parentId int) (*qmap.QM, error) {
	params := qmap.QM{
		"ne_status":        common.VUL_UNSHOW,
		"e_parent_task_id": parentId,
	}
	mgoSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_TASK, params)
	if res, err := mgoSession.GetOne(); err == nil {
		return res, nil
	} else {
		return nil, err
	}
}
