package transformer

import (
	"fmt"

	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/common_lib/mysql"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/mysql_model"
)

type ReportTransformer struct {
	orm.Transformer
}

func (h *ReportTransformer) ModifyItem(data qmap.QM) qmap.QM {
	//查找这个用例的 章节
	var chapters = make([]mysql_model.KnowledgeTestCaseChapter, 0)
	testCaseId := data.MustInt("test_case_id")
	demandId := data.MustInt("damand_id")
	s := mysql.GetSession()
	s.Where("test_case_id=?", testCaseId)
	s.Where("damand_id=?", demandId)
	if err := s.Find(&chapters); err == nil {
		for _, chapter := range chapters {
			//查找章节名字
			demandChapter := new(mysql_model.KnowledgeDemandChapter)
			sChapter := mysql.GetSession()
			if has, _ := sChapter.ID(chapter.DemandChapterId).Get(demandChapter); has {

			}
		}
	} else {
		fmt.Println(err)
	}
	return data
}

type ReportDetailTransformer struct {
	orm.Transformer
}

func (h *ReportDetailTransformer) ModifyItem(data qmap.QM) qmap.QM {
	taskId := data.MustInt("task_id")
	data["task_name"] = data["name"]
	if task, has := mysql_model.TaskFindById(taskId); has {
		data["task_name"] = task.Name
	}
	return data
}
