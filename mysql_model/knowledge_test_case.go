package mysql_model

import (
	"archive/zip"
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"skygo_detection/custom_util"
	"skygo_detection/custom_util/clog"
	"strings"

	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/app/sys_service"
	"xorm.io/builder"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/log"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/common_lib/mysql"
)

type KnowledgeTestCase struct {
	Id            int    `xorm:"not null pk autoincr comment('主键id') INT(11)" json:"id"`
	Name          string `xorm:"not null comment('测试用例名称') VARCHAR(255)" json:"name"`
	CaseUuid      string `xorm:"not null comment('测试用例真实id') VARCHAR(255)" json:"case_uuid"`
	ModuleId      string `xorm:"comment('测试组件/测试分类') VARCHAR(255)" json:"module_id"`
	ScenarioId    int    `xorm:"not null comment('安全检测场景') INT(11)" json:"scenario_id"`
	DemandId      int    `xorm:"not null comment('关联安全需求id，需求条目存另一个表') INT(11)" json:"demand_id"`
	Objective     string `xorm:"not null comment('测试目的') VARCHAR(255)" json:"objective"`
	Input         string `xorm:"not null comment('外部输入') VARCHAR(255)" json:"input"`
	TestProcedure string `xorm:"not null comment('测试步骤') VARCHAR(255)" json:"test_procedure"`
	TestStandard  string `xorm:"not null comment('验证标准') VARCHAR(255)" json:"test_standard"`
	Level         int    `xorm:"not null comment('测试难度') TINYINT(3)" json:"level"`
	TestCaseLevel int    `xorm:"not null comment('测试用例级别 1基础测试，2全面测试，3提高测试，4专家模式') TINYINT(3)" json:"test_case_level"`
	TestMethod    int    `xorm:"not null comment('测试方式，1黑盒 2灰盒 3白盒') TINYINT(3)" json:"test_method"`
	AutoTestLevel int    `xorm:"not null comment('自动化测试程度 1人工 2半自动化 3自动化') TINYINT(3)" json:"auto_test_level"`
	TestTool      string `xorm:"not null comment('测试工具，todo是否多选') VARCHAR(255)" json:"test_tool"`
	TestToolName  string `xorm:"not null comment('测试工具，todo是否多选') VARCHAR(255)" json:"test_tool_name"`
	TestToolsId   string `xorm:"not null comment('测试工具id，todo是否多选') TINYINT(3)" json:"test_tools_id"`
	TaskParam     string `xorm:"not null comment('任务参数') VARCHAR(255)" json:"task_param"`
	TestScript    string `xorm:"not null comment('测试脚本') VARCHAR(255)" json:"test_script"`
	TestResult    string `xorm:"not null comment('测试结果抽取') VARCHAR(255)" json:"test_result"`
	TestSketchMap string `xorm:"not null comment('测试环境搭建示意图') VARCHAR(255)" json:"test_sketch_map"`
	CreateUserId  int    `xorm:"not null comment('创建用户id') INT(11)" json:"create_user_id"`
	LastOpId      int    `xorm:"not null comment('最近操作用户id') INT(11)" json:"last_op_id"`
	CreateTime    int    `xorm:"not null comment('创建时间（秒）') INT(11)" json:"create_time"`
	TestParam     string `xorm:"not null comment('测试参数之前的block_list') VARCHAR(255)" json:"test_param"`
	Tag           string `xorm:"not null comment('标签，以，分割开') VARCHAR(255)" json:"tag"`
}

const (
	TestToolsHgClient = 1 // 测试工具， 1代表合规客户端工具
)

const KnowledgeTestCaseAutoTestLevelInter = 1
const KnowledgeTestCaseAutoTestLevelHalf = 2 // 半自动化 -- 新版有的，
const KnowledgeTestCaseAutoTestLevelAuto = 3 // 自动化测试程度

func (this *KnowledgeTestCase) Create(demandChapterIds []int) (*KnowledgeTestCase, error) {
	// 创建测试用例
	// todo 这部分先写死，后续人工填写这个参数值
	this.TestParam = common.BlockList
	_, err := mysql.GetSession().InsertOne(this)
	if err != nil {
		return nil, err
	}
	// 创建测试用例，关系需求，章节等关系库
	chapter := new(KnowledgeTestCaseChapter)
	for _, demandChapterId := range demandChapterIds {
		chapter.TestCaseId = this.Id
		chapter.SenarioId = this.ScenarioId
		chapter.DemandId = this.DemandId
		chapter.DemandChapterId = demandChapterId
		_, err = chapter.Create()
		if err != nil {
			return nil, err
		}
	}
	// 创建测试用例的版本信息
	history := new(KnowledgeTestCaseHistory)
	history.Content = "创建"
	history.OPId = this.CreateUserId
	history.TestCaseId = this.Id
	history.Create()
	return this, err
}

func (this *KnowledgeTestCase) CreateOnly() (int64, error) {
	// 创建测试用例
	return mysql.GetSession().InsertOne(this)
}

func (this *KnowledgeTestCase) Update(demandChapterIds []int, cols ...string) (int64, error) {
	// 更新关系表的库
	if len(demandChapterIds) != 0 {
		chapter := new(KnowledgeTestCaseChapter)
		// 先删除关系表，在新创建关系表
		chapter.TestCaseId = this.Id
		chapter.Remove()
		for _, demandChapterId := range demandChapterIds {
			chapter.SenarioId = this.ScenarioId
			chapter.DemandId = this.DemandId
			chapter.DemandChapterId = demandChapterId
			chapter.Create()
		}
	}
	// 更新测试用例库
	return mysql.GetSession().Table(this).ID(this.Id).Cols(cols...).Update(this)
}

func (this *KnowledgeTestCase) Remove() (int64, error) {
	return mysql.GetSession().Delete(this)
}

func (this *KnowledgeTestCase) RemoveById(id int) (int64, error) {
	// 删除关系库
	chapter := new(KnowledgeTestCaseChapter)
	chapter.TestCaseId = id
	_, err := chapter.Remove()
	if err != nil {
		return 0, err
	}
	// 删除用例库
	return mysql.GetSession().ID(id).Delete(this)
}

// 根据一组ID，查询出对应的一组测试案例
func KnowledgeTestCaseFindByIds(ids []int) []KnowledgeTestCase {
	models := make([]KnowledgeTestCase, 0)
	mysql.GetSession().Where(builder.In("id", ids)).Find(&models)
	return models
}

// 基于一组测试用例id，它们都支持合规测试，把它们进行整合，得到可用测试工具压缩包
// uuid 即通过它可用从mongodb中获取文件
func KnowledgeTestCaseMergeFiles(ids []int, fileName string) (uuid string, err error) {
	// 测试用例id得到测试用例记录
	tcModels := KnowledgeTestCaseFindByIds(ids)

	// 写入zip文件中的内容
	var b bytes.Buffer
	writer := zip.NewWriter(&b)

	// 从测试用例中查询所有要打包的文件
	for _, model := range tcModels {
		// 只取第一个测试脚本文件
		tcfModel := KnowledgeTestCaseFile{}
		if has, err := mysql.GetSession().Where("test_case_id = ?", model.Id).
			And("category = ?", KnowledgeTestCaseFileCategoryScript).Get(&tcfModel); err != nil {
			log.GetHttpLogLogger().Error(err.Error())
		} else {
			if !has {
				log.GetHttpLogLogger().Error("test_case_id do not has script")
				continue
			}
		}

		// 根据_id得到evaluate_test_case集合中的文件信息，打包得到一个压缩包
		fileIdStr := tcfModel.FileUuid // 文件id
		file, err := mongo.GridFSOpenId(common.MC_File, bson.ObjectIdHex(fileIdStr))
		if err != nil {
			panic(err)
		}

		// 文件名要统一修改为测试用_id做为名称的文件，保持后缀不变
		index := strings.LastIndex(file.Name(), ".")
		tail := file.Name()[index:]
		fileName := fmt.Sprintf("%d%s", model.Id, tail)

		// 打包
		fileWriter, err := writer.Create(fileName)
		if err != nil {
			// todo 打印日志
			panic(err)
		}

		bs, err := ioutil.ReadAll(file)
		if err != nil {
			panic(err)
		}

		_, err = fileWriter.Write(bs)
		if err != nil {
			panic(err)
		}
	}

	if err := writer.Close(); err != nil {
		panic(err)
	}

	// 测试用例的脚本文件打包zip
	zipFileName := fmt.Sprintf("%s.zip", fileName)
	fileId, err := mongo.GridFSUpload(common.MC_File, zipFileName, b.Bytes())
	if err != nil {
		panic(err)
	}

	return fileId, nil
}

// 通过检测场景id，查询里面测试用例是否需要固件检测工具，返回bool
func (this *KnowledgeTestCase) HasFirmwareToolsByScenarioId(id int) bool {
	// 遍历场景里每个用例
	// 判断用例是否包含固件检测工具
	return true
}

// 通过检测场景id，检测场景中存在需要连接设备的测试用例，返回bool
func (this *KnowledgeTestCase) HasWebSocketByScenarioId(id int) bool {
	// 遍历场景里每条用例
	// 判断用例是否需要长连接
	return true
}

func (this *KnowledgeTestCase) GetCountByScenarioId(scenarioId int) int64 {
	// 创建测试用例
	session := mysql.GetSession()
	session.Where("scenario_id=?", scenarioId)
	n, _ := session.FindAndCount(&[]KnowledgeTestCase{})
	return n
}

// 根据场景一组ID，查询出对应的一组测试案例
func (this *KnowledgeTestCase) KnowledgeTestCaseFindByScenarioIds(scenarioId int) []KnowledgeTestCase {
	models := make([]KnowledgeTestCase, 0)
	mysql.GetSession().Where("scenario_id=?", scenarioId).Find(&models)
	return models
}

// 根据场景ID查出涉及到的测试用例
func (this *KnowledgeTestCase) GetToolListByScenarioId(scenarioId, uid int) interface{} {
	type Item struct {
		TestTool     string `json:"test_tool"`
		TestToolName string `json:"test_tool_name"`
	}
	models := make([]Item, 0)
	if err := sys_service.NewSession().Session.Table(this).Cols("test_tool", "test_tool_name").Where("scenario_id=?", scenarioId).Distinct("test_tool").And("test_tool <> ?", "").Find(&models); err != nil {
		panic(err)
	}
	return models
}
func (this *KnowledgeTestCase) UpdateTag(id int, tag string) (int64, error) {
	this.Id = id
	this.Tag = tag
	return mysql.GetSession().Table(this).ID(this.Id).Cols("tag").Update(this)
}

func (this *KnowledgeTestCase) Copy(scenarioId, demandId int, chapters, testCases []interface{}) (int64, error) {
	var successNum int64
	for _, testCaseId := range testCases {
		var newCase = new(KnowledgeTestCase)
		if has, _ := mysql.GetSession().Where("id=?", testCaseId).Get(newCase); has {
			newCase.Id = 0
			newCase.ScenarioId = scenarioId
			newCase.DemandId = demandId
			if _, err := mysql.GetSession().InsertOne(newCase); err == nil {
				// 添加测试用例对应的章节
				for _, chapter := range chapters {
					tmpChapter := chapter.(float64)
					var newChapter = new(KnowledgeTestCaseChapter)
					newChapter.TestCaseId = newCase.Id
					newChapter.SenarioId = scenarioId
					newChapter.DemandId = demandId
					newChapter.DemandChapterId = int(tmpChapter)
					mysql.GetSession().InsertOne(newChapter)
				}
				successNum++
			}
		}
	}
	return successNum, nil
}

func GetBydName(caseName string, scenario_id int) (int, error) {
	models := KnowledgeTestCase{}
	session := mysql.GetSession()
	session.Where("name=?", caseName).Where("scenario_id=?", scenario_id)
	_, err := session.Get(&models)
	if err != nil {
		return 0, err
	}
	return models.Id, err
}

// 通过测试案例 id 获取测试脚本
func GetTestScriptById(testCaseIds []int) (testScriptSlice []string, err error) {
	testScriptSlice = make([]string, 0)
	err = sys_service.NewOrm().Table("knowledge_test_case").
		Cols("test_script").
		Where(custom_util.SliceToSqlInString(testCaseIds, "id")).
		Find(&testScriptSlice)
	if err != nil {
		clog.Error("GetTestScriptById Find Err", zap.Any("error", err))
		return testScriptSlice, err
	}
	clog.Debug("GetTestScriptById Info", zap.Any("Info", testScriptSlice))
	return testScriptSlice, nil
}
