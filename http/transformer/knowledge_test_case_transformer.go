package transformer

import (
	"strings"

	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/orm"
	"skygo_detection/mongo_model"
	"skygo_detection/mysql_model"
)

type KnowledgeTestCaseTransformer struct {
	orm.Transformer
}

func (h *KnowledgeTestCaseTransformer) ModifyItem(data qmap.QM) qmap.QM {
	//查询测试组件分类名字
	moduleId := data.MustString("module_id")
	evaluateModule, err := new(mongo_model.EvaluateModule).FindById(moduleId)
	if err == nil {
		data["module_type_name"] = evaluateModule.ModuleType
		data["module_name"] = evaluateModule.ModuleName
	} else {
		data["module_type_name"] = ""
		data["module_name"] = ""
	}
	//查询需求名字
	//查询章节名字
	//查询章节编号
	testCaseId := data.MustInt("id")
	dcid, dcname := new(mysql_model.KnowledgeTestCaseChapter).GetByTestCaseId(testCaseId)
	data["demand_chapter_id"] = dcid
	data["demand_chapter_name"] = dcname
	codes := []string{}
	for _, v := range dcid {
		model, has := mysql_model.KnowledgeDemandChapterFindById(v)
		if has {
			codes = append(codes, model.Code)
		}
	}
	data["demand_chapter_code"] = codes

	// 获取变更内容
	history := new(mysql_model.KnowledgeTestCaseHistory).GetAllByCaseId(testCaseId)
	content := []qmap.QM{}
	var opids []int
	if len(history) > 0 {
		for _, v := range history {
			opids = append(opids, v.OPId)
		}
		users, _ := new(mysql_model.SysUser).FindByIds(opids)
		for _, v := range history {
			tmp := qmap.QM{
				"version":    v.Version,
				"content":    v.Content,
				"time_stamp": v.TimeStamp,
				"user_name":  "",
			}
			if len(users) > 0 {
				for _, u := range users {
					if v.OPId == u.Id {
						tmp["user_name"] = u.Realname
						break
					}
				}
			}
			content = append(content, tmp)
		}

	}
	data["content"] = content

	//查询场景名称
	scenarioId := data.Int("scenario_id")
	if scenario, has, _ := mysql_model.KnowledgeScenarioFindById(scenarioId); has {
		data["scenario_name"] = scenario.Name
	}
	demandId := data.Int("demand_id")
	if demand, has := new(mysql_model.KnowledgeDemand).KnowledgeDemandFindById(demandId); has {
		data["demand_name"] = demand.Name
	}

	//获取测试脚本数组
	var testScriptsList = make([]qmap.QM, 0)
	testScriptsTmp := data.String("test_script")
	if len(testScriptsTmp) != 0 {
		testScripts := strings.Split(testScriptsTmp, "|")
		for _, testScript := range testScripts {
			if testScript != "" {
				f, err := mongo.GridFSOpenId(common.MC_File, bson.ObjectIdHex(testScript))
				if err == nil {
					tmp := qmap.QM{
						"name":  f.Name(),
						"value": testScript,
					}
					testScriptsList = append(testScriptsList, tmp)
				}
			}
		}
	}
	data["test_scripts"] = testScriptsList

	//获取环境搭建示意图
	var testSketchMapsList = make([]qmap.QM, 0)
	testSketchMapsTmp := data.String("test_sketch_map")
	if len(testSketchMapsTmp) != 0 {
		testSketchMaps := strings.Split(testSketchMapsTmp, "|")
		for _, testSketchMap := range testSketchMaps {
			if testSketchMap != "" {
				f, err := mongo.GridFSOpenId(common.MC_File, bson.ObjectIdHex(testSketchMap))
				if err == nil {
					tmp := qmap.QM{
						"name":  f.Name(),
						"value": testSketchMap,
					}
					testSketchMapsList = append(testSketchMapsList, tmp)
				}

			}
		}
	}
	data["test_sketch_maps"] = testSketchMapsList
	return data
}
