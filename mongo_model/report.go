package mongo_model

import (
	"context"
	"errors"
	"strconv"
	"strings"
	"time"

	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/report"
	"skygo_detection/mysql_model"

	"github.com/globalsign/mgo/bson"
)

type Report struct {
	Id                bson.ObjectId `bson:"_id,omitempty"`
	ProjectId         string        `bson:"project_id" json:"project_id"` // 项目id
	Name              string        `bson:"name"`                         // 报告名称名称
	FileId            string        `bson:"file_id"`                      // 文件id
	FileSize          int           `bson:"file_size"`                    // 文件大小
	ReportType        string        `bson:"report_type"`                  // 报告类型(周报:week,  初测报告:test, 复测报告:retest)
	RelativeId        []int         `bson:"relative_id"`                  // 相关用户
	Status            int           `bson:"status"`                       // 报告状态（0 创建，1审核，-1失败，2成功）
	OperatorId        int           `bson:"operator_id"`                  // 操作人员id
	CurrentOperatorId int           `bson:"current_operator_id"`          // 当前操作人员id
	UploadType        int64         `bson:"upload_type"`                  // 报告类型 1 上传报告 2 导出报告
	UpdateTime        int64         `bson:"update_time"`
	CreateTime        int64         `bson:"create_time"`
}

func (this *Report) Create(opId int, rawInfo *qmap.QM) (*Report, error) {
	this.Id = bson.NewObjectId()
	this.ProjectId = rawInfo.String("project_id")
	this.ReportType = rawInfo.String("report_type")

	if rawInfo.Int64("upload_type") == 1 {
		this.UploadType = rawInfo.Int64("upload_type")
	} else {
		this.UploadType = 2
	}
	this.Status = common.RAS_NEW
	this.OperatorId = opId
	this.CurrentOperatorId = opId
	this.CreateTime = custom_util.GetCurrentMilliSecond()
	fileId := rawInfo.MustString("file_id")
	if fi, err := mongo.GridFSOpenId(common.MC_File, bson.ObjectIdHex(fileId)); err == nil {
		this.FileSize = int(fi.Size())
		this.Name = fi.Name()
		fi.Close()
	} else {
		panic(err)
	}
	if newFileId, err := mongo.GridFSRename(common.MC_File, this.Name, bson.ObjectIdHex(fileId)); err == nil {
		this.FileId = newFileId
	} else {
		panic(err)
	}
	mongoClient := mongo.NewMgoSession(common.MC_REPORT)
	if err := mongoClient.Insert(this); err == nil {
		// 创建报告，添加报告节点记录
		nodeInfo := qmap.QM{
			"project_id": this.ProjectId,
			"report_id":  this.Id.Hex(),
			"auditor_id": opId,
			"name":       "创建",
			"status":     "create",
			"history": History{
				Result:    true,
				Operation: "创建",
				Comment:   "",
				OpId:      opId,
				OpTime:    custom_util.GetCurrentMilliSecond(),
			},
		}
		if _, err := new(ReportNode).Create(nodeInfo); err != nil {
			return nil, errors.New("Create Report Node Error")
		}
		return this, nil
	} else {
		return nil, err
	}
}

func (this *Report) Update(reportId string, update qmap.QM) error {
	selector := bson.M{"_id": bson.ObjectIdHex(reportId)}
	updateItem := bson.M{
		"$set": update,
	}
	return mongo.NewMgoSession(common.MC_REPORT).Update(selector, updateItem)
}

func (this *Report) ExportReport(opId int, projectId, reportType string, evaluateItems []string, ctx context.Context) (*Report, error) {
	todayReportCount, countErr := this.CountTodayReport(projectId, reportType)
	custom_util.CheckErr(countErr)
	var fileId string
	var reportErr error
	if strings.HasPrefix(reportType, common.REPORT_ASSETPRE) {
		fileId, reportErr = report.GetAssetReportDocx(projectId, reportType, todayReportCount+1, evaluateItems)
	} else {
		fileId, reportErr = report.GetItemFirstReportDocx(projectId, reportType, todayReportCount+1, evaluateItems)
	}
	custom_util.CheckErr(reportErr)
	reportQM := &qmap.QM{
		"project_id":  projectId,
		"file_id":     fileId,
		"report_type": reportType,
	}
	return this.Create(opId, reportQM)
}

func (this *Report) CountTodayReport(projectId, reportType string) (int, error) {
	time, _ := time.Parse("2006-01-02", time.Now().Format("2006-01-02"))
	nano := (time.Unix() - 3600*8) * 1000
	params := qmap.QM{
		"e_project_id":    projectId,
		"gte_create_time": nano,
		"e_report_type":   reportType,
		"e_upload_type":   2,
	}
	count, err := mongo.NewMgoSession(common.MC_REPORT).AddCondition(params).SetLimit(1000).Count()
	return count, err
}

/*
*
拼接报告数据
[

	{
		"evaluate_items": [
			{
				"_id": "5e79d30b24b647600e526504",
				"create_time": 1585042187598,
				"evaluate_type": "",
				"has_vulnerability": 0,
				"item_vulnerability": {
					"influence": "",
					"level": 0,
					"name": "",
					"status": 0,
					"suggest": ""
				},
				"level": 2,
				"name": "扽还放扽",
				"objective": "发森阿森",
				"op_id": 1,
				"op_username": "李清",
				"procedure": "风赛风赛风",
				"project_id": "5e73a77c24b64720bdd8c9ea",
				"result": "发布仍行将赛庚庚庚",
				"status": 1,
				"target_id": "5e8163cb24b64764f4342f7a",
				"update_time": 1585042187598
			}
		],
		"target_info": {
			"_id": "5e8163cb24b64764f4342f7a",
			"attributes": {},
			"create_time": 1585537995419,
			"evaluate_type": "NEW APP",
			"name": "恒润TBOX",
			"op_id": 2,
			"project_id": "5e73a77c24b64720bdd8c9ea",
			"update_time": 1585537995419
		}
	}

]
*/
func (this *Report) assembleReportData(projectId string, evaluateItems []bson.ObjectId, ctx context.Context) ([]qmap.QM, error) {
	params := qmap.QM{
		"e_project_id": projectId,
		"in__id":       evaluateItems,
	}
	ormSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, params)
	ormSession.SetTransformFunc(func(qm qmap.QM) qmap.QM {
		if userMap, err := new(mysql_model.SysUser).GetUserInfo(qm.MustInt("op_id")); err == nil {
			qm["op_username"] = userMap.String("realname")
		}
		return qm
	})
	if data, err := ormSession.Get(); err == nil {
		evaluateTargets := qmap.QM{}
		for _, tempItem := range *data {
			var item qmap.QM = tempItem
			targetId := item.String("target_id")
			if tempItems, has := evaluateTargets.TryInterface(targetId); !has {
				evaluateTargets[targetId] = []qmap.QM{item}
			} else {
				var items = tempItems.([]qmap.QM)
				evaluateTargets[targetId] = append(items, item)
			}
		}
		reportContent := []qmap.QM{}
		for targetId, items := range evaluateTargets {
			params := qmap.QM{
				"e__id": bson.ObjectIdHex(targetId),
			}
			if target, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, params).GetOne(); err == nil {
				evaluateTarget := qmap.QM{
					"target_info":    *target,
					"evaluate_items": items,
				}
				reportContent = append(reportContent, evaluateTarget)
			} else {
				return nil, err
			}
		}
		return reportContent, nil
	} else {
		return nil, err
	}
}

func (this *Report) getReportData(projectId string, evaluateItems []string, ctx context.Context) (map[string]qmap.QM, error) {
	// 获取Project信息
	projectInfo, projectErr := mongo.NewMgoSessionWithCond(common.MC_PROJECT, qmap.QM{"e__id": bson.ObjectIdHex(projectId)}).GetOne()
	custom_util.CheckErr(projectErr)
	// 获取测试项列表
	itemList, itemErr := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, qmap.QM{"e_project_id": projectId, "in__id": evaluateItems}).SetLimit(1000).Get()

	custom_util.CheckErr(itemErr)
	itemData := map[string]map[string]qmap.QM{}
	targetIds := []string{}
	vulIds := []bson.ObjectId{}
	vulCount := qmap.QM{}
	for index, item := range *itemList {
		if itemData[item["target_id"].(string)] == nil {
			itemData[item["target_id"].(string)] = map[string]qmap.QM{}
		}
		itemData[item["target_id"].(string)][strconv.Itoa(index)] = item

		targetIds = append(targetIds, item["target_id"].(string))

		if item["has_vulnerability"].(int) >= 0 {
			for _, vulId := range item["item_vulnerability"].([]interface{}) {
				vulIds = append(vulIds, bson.ObjectIdHex(vulId.(string)))
			}
		}

		vulCount[item["target_id"].(string)] = map[int]int{
			0: 0,
			1: 0,
			2: 0,
			3: 0,
			4: 0,
		}
	}
	// 获取 测试对象 列表
	tarCondition := qmap.QM{"e_project_id": projectId, "in__id": targetIds}
	target, targetErr := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, tarCondition).SetLimit(1000).Get()
	custom_util.CheckErr(targetErr)
	targetList := []qmap.QM{}
	for _, tar := range *target {
		target := qmap.QM{
			"name": tar["name"],
			"id":   tar["_id"].(bson.ObjectId).Hex(),
		}
		targetList = append(targetList, target)
	}
	// 获取漏洞列表
	vul, vulErr := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VULNERABILITY, qmap.QM{"in__id": vulIds}).SetLimit(1000).Get()
	custom_util.CheckErr(vulErr)

	vulList := map[string]map[string]qmap.QM{}

	for index, vulItem := range *vul {
		targetId := vulItem["target_id"].(string)
		vulCount[targetId].(map[int]int)[vulItem["level"].(int)]++

		if vulList[targetId] == nil {
			vulList[targetId] = map[string]qmap.QM{}
		}
		vulList[targetId][strconv.Itoa(index)] = vulItem
	}

	reportData := map[string]qmap.QM{
		"project": {
			"id":   projectId,
			"name": projectInfo.String("name"),
		},
		"target": {
			"list": *target,
		},
		"item": {
			"list":   itemData,
			"target": targetList,
		},
		"vul": {
			"list":   vulList,
			"target": targetList,
		},
		"vulCount": qmap.QM{
			"total":  vulCount,
			"target": targetList,
		},
	}

	return reportData, nil
}

func (this *Report) GetList(rawInfo qmap.QM, uid int, ctx context.Context) (*qmap.QM, error) {
	match := bson.M{
		"$or": []bson.M{{"operator_id": uid}, {"relative_id": uid}},
	}

	if name, has := rawInfo.TryString("name"); has {
		match["name"] = bson.M{
			"$regex": name,
		}
	}
	if status, has := rawInfo.TryInt("status"); has {
		match["status"] = status
	}
	if myStage := rawInfo.Int("my_stage"); myStage > 0 {
		match["current_op_id"] = uid
	}
	Operations := []bson.M{
		{"$match": match},
	}

	ormSession := mongo.NewMgoSession(common.MC_REPORT)
	if page, has := rawInfo.TryInt("page"); has {
		ormSession = ormSession.SetPage(page)
	}

	if limit, has := rawInfo.TryInt("limit"); has {
		ormSession = ormSession.SetLimit(limit)
	}

	return ormSession.MATCHGetPage(Operations)
}

/**
 * @Description:统计待操作的项目报告数量
 * @param userId
 * @return int
 */
func (this *Report) CountOperateReport(userId int64) int {
	params := qmap.QM{
		"e_current_operator_id": userId,
	}
	if count, err := mongo.NewMgoSessionWithCond(common.MC_REPORT, params).Count(); err == nil {
		return count
	} else {
		return 0
	}
}

/**
 * @Description:统计最近某一段时间内，发布的报告数量
 * @param userId
 * @return int
 */
func (this *Report) CountPublishedReport(projectIds []string, startTime int64) int {
	params := qmap.QM{
		"in_project_id":   projectIds,
		"e_status":        common.RS_SUCCESS,
		"gte_update_time": startTime,
	}
	if count, err := mongo.NewMgoSessionWithCond(common.MC_REPORT, params).Count(); err == nil {
		return count
	} else {
		return 0
	}
}
