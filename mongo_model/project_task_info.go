package mongo_model

import (
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
)

type ProjectTaskInfo struct {
	Id                 string `bson:"_id,omitempty"`
	ProjectId          string `bson:"project_id"`           // 项目Id
	CreateTaskNumber   int    `bson:"create_task_number"`   // 创建任务数量
	TestingTaskNumber  int    `bson:"testing_task_number"`  // 测试中任务数量
	CompleteTaskNumber int    `bson:"complete_task_number"` // 结束任务数量
	CreateFormatTime   string `bson:"create_format_time"`
	CreateTime         int64  `bson:"create_time"`
}

func (this *ProjectTaskInfo) GetProjectTaskSeries(projectId string) qmap.QM {
	params := qmap.QM{
		"e_project_id": projectId,
	}
	mgoClient := mongo.NewMgoSessionWithCond(common.MC_PROJECT_TASK_INFO, params)
	mgoClient.AddSorter("_id", 0)
	xAxis := []string{}
	createArr := []int{}
	testingArr := []int{}
	completeArr := []int{}
	if list, err := mgoClient.SetLimit(10000).Get(); err == nil {
		for _, item := range *list {
			var itemQM qmap.QM = item
			xAxis = append(xAxis, itemQM.String("create_format_time"))
			createArr = append(createArr, itemQM.Int("create_task_number"))
			testingArr = append(testingArr, itemQM.Int("testing_task_number"))
			completeArr = append(completeArr, itemQM.Int("complete_task_number"))
		}
	}
	result := qmap.QM{
		"x_axis":   xAxis,
		"create":   createArr,
		"testing":  testingArr,
		"complete": completeArr,
	}
	return result
}
