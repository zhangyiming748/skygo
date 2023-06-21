package report

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/unioffice/color"
	"skygo_detection/lib/unioffice/document"
	"skygo_detection/lib/unioffice/measurement"
	"skygo_detection/lib/unioffice/schema/soo/wml"
	"skygo_detection/service"
)

type AssetReport struct {
	ReportCommont
}

func GetAssetReportDocx(projectId, reportType string, todayReportCount int, evaluateItems []string) (fileId string, err error) {
	reportTemplateConfig = service.LoadReportTemplateConfig()
	allData, err := getAssetData(projectId, reportType, evaluateItems)
	if err != nil {
		panic(err)
	}
	if doc, err := document.OpenTemplate(reportTemplateConfig.ReportTemplate); err == nil {
		weekReport := WeekReportDocx{&ReportParse{doc}}
		assetReport := AssetReport{ReportCommont{weekReport}}
		assetReport.addHeader()
		assetReport.addFooter()
		assetReport.Title(allData, reportType)
		assetReport.Catalog()
		assetReport.Statement(allData)
		assetReport.Abbreviation()
		assetReport.ItemSummarize(allData)
		assetReport.Summarize(allData)
		assetReport.PenetrationTest(allData)
		assetReport.Copyright()
		assetReport.Appendix()
		assetReport.SaveToFile("./tmp.docx")

		buffer := new(bytes.Buffer)
		assetReport.Save(buffer)
		if fileContent, err := ioutil.ReadAll(buffer); err == nil {
			fileName := GenerateReportName(projectId, reportType, todayReportCount)
			if fileId, err := mongo.GridFSUpload(common.MC_File, fileName, fileContent); err == nil {
				return fileId, nil
			} else {
				return "", err
			}
		} else {
			return "", err
		}

	} else {
		panic(err)
	}
}

func getAssetData(projectId, reportType string, evaluateItems []string) (qmap.QM, error) {
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(projectId),
	}
	//获取项目信息
	project, err := mongo.NewMgoSessionWithCond(common.MC_PROJECT, params).GetOne()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("project is %s", err.Error()))
	}
	//获取测试用例信息
	var items []map[string]interface{}
	if len(evaluateItems) == 0 {
		paramsItem := qmap.QM{
			"e_project_id": projectId,
		}
		if reportType == common.RT_TEST {
			paramsItem["e_test_phase"] = 1
		}
		items = mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, paramsItem).All()
	} else {
		paramsItem := qmap.QM{
			"in__id": evaluateItems,
		}
		if reportType == common.RT_TEST {
			paramsItem["e_test_phase"] = 1
		}
		items = mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, paramsItem).All()
	}

	//通过测试用例获取资产id
	assetIds := make([]string, 0, len(items))
	tmpAssetIds := map[string]interface{}{}
	for _, item := range items {
		if v, ok := item["asset_id"]; ok {
			tmpAssetIds[v.(string)] = 1
		}
	}
	for k := range tmpAssetIds {
		assetIds = append(assetIds, k)
	}

	params = qmap.QM{
		"in__id": assetIds,
	}
	assets := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, params).All()
	//获取漏洞信息
	params = qmap.QM{
		"in_asset_id": assetIds,
	}
	vuls := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VULNERABILITY, params).All()

	//获取车企的名字
	params = qmap.QM{
		"e__id": bson.ObjectIdHex(project.String("company")),
	}
	factory, err := mongo.NewMgoSessionWithCond(common.MC_FACTORY, params).GetOne()
	if err != nil {
		factory = &qmap.QM{"name": " 未知 "}
		//return nil, errors.New(fmt.Sprintf("factory is %s", err.Error()))
	}
	allData := qmap.QM{
		"project": project,
		"assets":  assets,
		"vuls":    vuls,
		"items":   items,
		"factory": factory,
	}
	return allData, nil
}

func GetReportName(projectId, reportType string, todayReportCount int) string {
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(projectId),
	}
	projectName := ""
	if project, err := mongo.NewMgoSessionWithCond(common.MC_PROJECT, params).GetOne(); err == nil {
		projectName = project.String("name")
	}
	switch reportType {
	case "test":
		return fmt.Sprintf("%s-资产初测报告-%s-%d.docx", projectName, time.Now().Format("20060102"), todayReportCount)
	case "retest":
		return fmt.Sprintf("%s-资产复测报告-%s-%d.docx", projectName, time.Now().Format("20060102"), todayReportCount)
	default:
		panic("unknown report type!")
	}
}

// 首页
func (this *AssetReport) Title(allData qmap.QM, reportType string) {
	//获取工程名称
	project := allData["project"]
	tmpProject := project.(*qmap.QM)
	name := tmpProject.String("name")

	run := this.AddTitleWihtStyle()
	if reportType == common.REPORT_ASSETRETEST {
		run.AddText(fmt.Sprintf("%s资产复测报告", name))
	} else {
		run.AddText(fmt.Sprintf("%s资产初测报告", name))
	}
	this.AddMultiBreak(26)
	niceNumber := niceNumber(1)
	this.AddTitle5().AddText(fmt.Sprintf("报告编号：SKYGO-%s-%s", time.Now().Format("20060102"), niceNumber))
	this.AddTitle5().AddText("报告提供商：360 SKY-GO智能网联汽车安全实验室")
	this.AddParagraph().AddRun().AddPageBreak()
}

// 目录
func (this *AssetReport) Catalog() {
	doc := this.Document
	doc.AddParagraph().AddRun().AddField(document.FieldTOC)
	doc.AddParagraph().AddRun().AddPageBreak()
}

// 1 声明
func (this *AssetReport) Statement(allData qmap.QM) {
	this.AddParagraph()
	content1 := "本报告是针对 %s-%s 的安全评测报告。"
	content2 := "本报告测评结论的有效性建立在被测车联网系统提供相关证据的真实性基础之上。"
	content3 := "本报告中给出的测评结论仅对被测车联网系统当时的安全状态有效。当评测工作完成后，被测零部件因功能迭代升级而涉及到的相关组件发生变化，本报告将不再适用。"
	content4 := "在任何情况下，若需引用本报告中的测评结果或结论都应保持其原有的意义，不得对相关内容擅自进行增加、修改和伪造或掩盖事实。"
	project := allData["project"]
	tmpProject := project.(*qmap.QM)
	//获取车企名称
	name := tmpProject.String("name")
	//获取车型名称
	brand := tmpProject.String("brand")
	this.AddTitleWihtStyle().AddText("声明")
	this.AddMain2().AddText(fmt.Sprintf(content1, name, brand))
	this.AddMain2().AddText(content2)
	this.AddMain2().AddText(content3)
	this.AddMain2().AddText(content4)
	this.AddParagraph().AddRun().AddPageBreak()
}

// 2 缩写词汇表
func (this *AssetReport) Abbreviation() {
	this.AddTitleWihtStyle().AddText("缩写词汇表")
	table := this.AddTable()
	table.Properties().SetWidthPercent(100)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, measurement.Zero)
	row := table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("缩写")
	row.AddCell().AddParagraph().AddRun().AddText("释义")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("ASIL")
	row.AddCell().AddParagraph().AddRun().AddText("Automotive Safety Integrity Level")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("AVES")
	row.AddCell().AddParagraph().AddRun().AddText("Alliance Vehicle Evaluation Standard")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("B")
	row.AddCell().AddParagraph().AddRun().AddText("B sample - See definition in section: V3P LOGIC – ECU TECHNICAL DEFINITION")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("BDV")
	row.AddCell().AddParagraph().AddRun().AddText("BDV sample - See definition in section: V3P LOGIC – ECU TECHNICAL DEFINITION")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("BOM")
	row.AddCell().AddParagraph().AddRun().AddText("Bill Of Materials")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("C")
	row.AddCell().AddParagraph().AddRun().AddText("C sample - See definition in section: V3P LOGIC – ECU TECHNICAL DEFINITION")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("CAN")
	row.AddCell().AddParagraph().AddRun().AddText("Controller Area Network")
	this.AddParagraph().AddRun().AddPageBreak()
}

// 1 项目综述
func (this *AssetReport) ItemSummarize(allData qmap.QM) {
	this.AddParagraph()
	this.AddParagraph()
	content := "本次测试通过对 %s%s 联网系统进行渗透测试，及时发现车联网系统中存在的安全问题；" +
		"然后将渗透测试的结果以及相应安全修复建议反馈给%s，用于指导 %s 组织开展安全修复工作；" +
		"最后通过对修复后的车联网系统进行安全复测，确认车联网系统安全问题修复的有效性，最终达到提升%s%s车联网系统的安全性的目的。"
	tmpFactory := allData["factory"]
	factory := tmpFactory.(*qmap.QM)
	companyName := factory.String("name")
	tmpProject := allData["project"]
	project := tmpProject.(*qmap.QM)
	brand := project.String("brand")

	//1 项目综述
	this.AddHeading1().AddText("项目综述")
	//1.1 项目目的
	this.AddHeading2().AddText("项目目的")
	this.AddMain2().AddText(fmt.Sprintf(content, companyName, brand, companyName, companyName, companyName, brand))
	//1.2 项目范围
	this.AddHeading2().AddText("项目范围")

	table := this.AddTable()
	table.Properties().SetWidthPercent(100)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, measurement.Zero)
	row := table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("资产名称")
	row.AddCell().AddParagraph().AddRun().AddText("资产版本")
	row.AddCell().AddParagraph().AddRun().AddText("其他信息")
	tmpAssets := allData["assets"]
	assets := tmpAssets.([]map[string]interface{})
	for _, asset := range assets {
		name := asset["name"]
		version := asset["version"]
		other := asset["attributes"]
		var attributes string
		for k, v := range other.(map[string]interface{}) {
			attributes += fmt.Sprintf("%v:%v ", k, v)
		}
		row := table.AddRow()
		cell := row.AddCell()
		cell.AddParagraph().AddRun().AddText(name.(string))

		cell = row.AddCell()
		cell.AddParagraph().AddRun().AddText(version.(string))

		cell = row.AddCell()
		cell.AddParagraph().AddRun().AddText(attributes)
	}
	this.AddParagraph().AddRun().AddPageBreak()
}

// 2 总体评价
func (this *AssetReport) Summarize(allData qmap.QM) {
	this.AddParagraph()
	tmpVuls := allData["vuls"]
	vuls := tmpVuls.([]map[string]interface{})
	vulNum := len(vuls)
	var serious int
	var hight int
	var middle int
	var low int
	var repair int
	var unrepair int
	for _, vul := range vuls {
		switch vul["level"] {
		case 1:
			low++
		case 2:
			middle++
		case 3:
			hight++
		case 4:
			serious++
		}
		switch vul["status"] {
		case 0:
			unrepair++
		case 1:
			repair++
		}
	}
	//4 总体评价
	project := allData["project"]
	tmpProject := project.(*qmap.QM)
	//获取车企名称
	name := tmpProject.String("name")
	//获取车型名称
	brand := tmpProject.String("brand")
	this.AddHeading1().AddText("总体评价")
	//4.1 项目总结论
	this.AddHeading2().AddText("项目总结论")
	content := "本项目%s提供的待测系统为车联网系统。" +
		"360 Sky-Go团队在整个渗透测试过程中测试内容涵盖应用安全、管理安全、通信安全、数据安全等内容，主要手段有动态调试、" +
		"逆向分析、漏洞扫描、渗透测试等手段，在%s 架构环境下进行了渗透测试，并对其中存在的漏洞进行了综合分析评估。"
	this.AddMain2().AddText(fmt.Sprintf(content, name, brand))
	content = "在首次渗透测试中，" +
		"360 Sky-Go团队发现车联网系统的安全问题共计%v个（严重：%v个；高危：%v个；中危：%v个；低危：%v个），" +
		"%v个已修复，%v个未修复，其中%v个漏洞，%v方表示可接受。因此我司评估%v联网系统的网络安全水平为XX。"
	this.AddMain2().AddText(fmt.Sprintf(content, vulNum, serious, hight, middle, low, repair, unrepair, vulNum, name, name))
	this.AddMain2().AddText("本次测试自2020年3月27日至2020年12月30日结束，整个测试实施过程共分为测试环境准备阶段、" +
		"初次测试阶段、初测报告编制阶段、复测阶段、复测报告编制阶段，共五个阶段。")
	//todo 甘特图还没弄
	this.AddMain2().AddText("甘特图 无 ")

	//4.2 漏洞统计与分布
	this.AddHeading2().AddText("漏洞统计与分布")
	//如果不存在漏洞，就返回 本车型未发现漏洞
	if len(vuls) == 0 {
		this.AddMain2().AddText("本车型未发现漏洞")
		return
	}
	//4.2.1 总统计图
	this.AddHeading3().AddText("总统计图")
	if serious == 0 && hight == 0 && middle == 0 {
		this.AddMain2().AddText("图2 漏洞等级分布 暂无")
	} else {
		pictureParams := qmap.QM{
			"red":    qmap.QM{"value": float64(serious), "label": fmt.Sprintf("严重%d个", serious)},
			"yellow": qmap.QM{"value": float64(hight), "label": fmt.Sprintf("高危%d个", hight)},
			"blue":   qmap.QM{"value": float64(middle), "label": fmt.Sprintf("中危%d个", middle)},
		}
		picture1 := this.NewPieImage(pictureParams, "图2", "")
		this.AddImageFromFile(picture1)
		this.AddMain2().AddText("图2 漏洞等级分布")
	}
	this.AddMain2().AddTab()
	//资产漏洞
	tmpAssets := allData["assets"]
	assets := tmpAssets.([]map[string]interface{})
	for _, asset := range assets {
		assetId := asset["_id"]
		var l float64
		var m float64
		var h float64
		var s float64
		for _, vul := range vuls {
			if assetId == vul["asset_id"] {
				switch vul["level"] {
				case 1:
					l++
				case 2:
					m++
				case 3:
					h++
				case 4:
					s++
				}
			}
		}
		if s == 0 && h == 0 && m == 0 {
			this.AddMain2().AddText(fmt.Sprintf("图3 %s漏洞等级分布 暂无", asset["name"]))
		} else {
			tmpParams := qmap.QM{
				"red":    qmap.QM{"value": s, "label": "严重"},
				"yellow": qmap.QM{"value": h, "label": "高危"},
				"blue":   qmap.QM{"value": m, "label": "中危"},
			}
			picture1 := this.NewDonutImage(tmpParams, "图3", "")
			this.AddImageFromFile(picture1)
			this.AddMain2().AddText(fmt.Sprintf("图3 %s漏洞等级分布", asset["name"]))
		}
	}
	this.AddMain2().AddTab()

	//所有漏洞类型分布
	var riskTypes = map[int]int{}
	for _, vul := range vuls {
		tmp := vul["risk_type"]
		riskTypes[tmp.(int)] += 1
	}

	var vulTypes = []qmap.QM{}
	var isVul = false
	for k, v := range riskTypes {
		params := qmap.QM{
			"e__id": k,
		}
		vulType, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_TYPE, params).GetOne()
		if err != nil {
			continue
		}
		//根据漏洞类型获取漏洞名称
		if v != 0 {
			isVul = true
		}
		tmp := qmap.QM{"value": float64(v), "label": vulType.String("name")}
		vulTypes = append(vulTypes, tmp)
	}
	if isVul {
		picture1 := this.NewPieImage1(vulTypes, "图4", "")
		this.AddImageFromFile(picture1)
		this.AddMain2().AddText("图4 漏洞类型分类")
	} else {
		this.AddMain2().AddText("图4 漏洞类型分类 暂无")
	}

	//4.2.2 各资产统计图
	this.AddHeading3().AddText("各资产统计图")
	//todo 各个资产
	for _, asset := range assets {
		assetName := asset["name"]
		this.AddMain2().AddText(assetName.(string))

		assetId := asset["_id"]
		//漏洞等级分布
		var l float64
		var m float64
		var h float64
		var s float64
		for _, vul := range vuls {
			if assetId == vul["asset_id"] {
				switch vul["level"] {
				case 1:
					l++
				case 2:
					m++
				case 3:
					h++
				case 4:
					s++
				}
			}
		}
		if s == 0 && h == 0 && m == 0 {
			this.AddMain2().AddText(fmt.Sprintf("图5 %s漏洞等级分布 暂无", assetName))
		} else {
			tmpParams := qmap.QM{
				"red":    qmap.QM{"value": s, "label": "严重"},
				"yellow": qmap.QM{"value": h, "label": "高危"},
				"blue":   qmap.QM{"value": m, "label": "中危"},
			}
			picture1 := this.NewPieImage(tmpParams, "图5", "")
			this.AddImageFromFile(picture1)
			this.AddMain2().AddText(fmt.Sprintf("图5 %s漏洞等级分布", assetName))
		}

		//漏洞类型分布
		var assetRiskTypes = map[int]int{}
		var assetVulTypes = []qmap.QM{}
		for _, vul := range vuls {
			if assetId == vul["asset_id"] {
				tmp := vul["risk_type"]
				assetRiskTypes[tmp.(int)] += 1
			}
		}
		var isVul = false
		for k, v := range assetRiskTypes {
			params := qmap.QM{
				"e__id": k,
			}
			vulType, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_TYPE, params).GetOne()
			if err != nil {
				continue
			}
			if v != 0 {
				isVul = true
			}
			//根据漏洞类型获取漏洞名称
			tmp := qmap.QM{"value": float64(v), "label": vulType.String("name")}
			assetVulTypes = append(assetVulTypes, tmp)
		}
		if isVul {
			picture1 := this.NewDonutImage1(assetVulTypes, "图6", "")
			this.AddImageFromFile(picture1)
			this.AddMain2().AddText(fmt.Sprintf("图6 %s漏洞类型分类", assetName))
		} else {
			this.AddMain2().AddText(fmt.Sprintf("图6 %s漏洞类型分类 暂无", assetName))
		}
	}

	//4.2.3 漏洞总表
	this.AddHeading3().AddText("漏洞总表")
	this.AddMain2().AddText("表2：漏洞列表")

	if len(vuls) != 0 {
		table := this.AddTable()
		table.Properties().SetWidthPercent(100)
		borders := table.Properties().Borders()
		borders.SetAll(wml.ST_BorderSingle, color.Auto, measurement.Zero)
		row := table.AddRow()
		row.AddCell().AddParagraph().AddRun().AddText("漏洞ID")
		row.AddCell().AddParagraph().AddRun().AddText("漏洞名称")
		row.AddCell().AddParagraph().AddRun().AddText("漏洞等级")
		row.AddCell().AddParagraph().AddRun().AddText("资产")
		row.AddCell().AddParagraph().AddRun().AddText("修复状态/接受情况")

		for _, vul := range vuls {
			id := vul["_id"]
			name := vul["name"]
			level := vul["level"]
			var levelName string
			switch level {
			case 0:
				levelName = "提示"
			case 1:
				levelName = "低危"
			case 2:
				levelName = "中危"
			case 3:
				levelName = "高危"
			case 4:
				levelName = "严重"
			}
			assetId := vul["asset_id"]
			var assetName string
			assetParams := qmap.QM{
				"e__id": assetId,
			}
			if asset, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ASSET, assetParams).GetOne(); err == nil {
				assetName = asset.String("name")
			}
			status := vul["status"]
			var fixStatus string
			switch status {
			case 0:
				fixStatus = "未修复"
			case 1:
				fixStatus = "已修复"
			case 2:
				fixStatus = "重打开"
			}
			row := table.AddRow()
			row.AddCell().AddParagraph().AddRun().AddText(id.(bson.ObjectId).Hex())
			row.AddCell().AddParagraph().AddRun().AddText(name.(string))
			row.AddCell().AddParagraph().AddRun().AddText(levelName)
			row.AddCell().AddParagraph().AddRun().AddText(assetName)
			row.AddCell().AddParagraph().AddRun().AddText(fixStatus)
		}
	}
	this.AddParagraph().AddRun().AddPageBreak()
}

// 5渗透测试详情
func (this *AssetReport) PenetrationTest(allData qmap.QM) {
	//5 渗透测试详情
	this.AddHeading1().AddText("渗透测试详情")
	//获取资产
	tmpAssets := allData["assets"]
	assets := tmpAssets.([]map[string]interface{})
	//资产
	for _, asset := range assets {
		assetId := asset["_id"]
		//5.1 资产
		//this.AddMain2().AddText("工作环境架构逻辑图")
		//测试组件
		moduleTypeIds := asset["module_type_id"]
		for _, moduleTypeId := range moduleTypeIds.([]interface{}) {
			paramsModule := qmap.QM{
				"e__id": bson.ObjectIdHex(moduleTypeId.(string)),
			}
			if modult, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, paramsModule).GetOne(); err == nil {
				//moduleTitle := fmt.Sprintf("5.%v.%v 测试组件:%v", i+1, moduleIndex+1, modult.String("module_name"))
				moduleTitle := fmt.Sprintf("测试组件:%v", modult.String("module_name"))
				//测试用例
				tmpItems := allData["items"]
				items := tmpItems.([]map[string]interface{})
				for _, item := range items {
					tmpModuleTypeId := item["module_type_id"]
					if assetId == item["asset_id"] && moduleTypeId.(string) == tmpModuleTypeId.(string) {
						//5.1 资产
						assetTitle := fmt.Sprintf("资产:%s", asset["name"])
						this.AddHeading2().AddText(assetTitle)
						//5.1.1 测试组件
						this.AddHeading3().AddText(moduleTitle)
						itemId := item["_id"]
						//this.AddMain2().AddText(fmt.Sprintf("测试用例:%s", item["name"]))
						//5.1.1.1 测试用例
						this.AddHeading4().AddText(fmt.Sprintf("测试用例:%s", item["name"]))
						tmpVuls := allData["vuls"]
						vuls := tmpVuls.([]map[string]interface{})
						for _, vul := range vuls {
							if itemId == vul["item_id"] {
								vulName := vul["name"]
								vulId := vul["_id"]
								vulLevel := vul["level"]
								var vulLevelName string
								switch vulLevel {
								case 0:
									vulLevelName = "提示"
								case 1:
									vulLevelName = "低危"
								case 2:
									vulLevelName = "中危"
								case 3:
									vulLevelName = "高危"
								case 4:
									vulLevelName = "严重"
								}
								vulStatus := vul["status"]
								var vulStatusName string
								switch vulStatus {
								case 0:
									vulStatusName = "未修复"
								case 1:
									vulStatusName = "已修复"
								case 2:
									vulStatusName = "重打开"
								}
								vulRiskType := vul["risk_type"]
								vulDescribe := vul["describe"]
								vulInfluence := vul["influence"]
								this.AddMain2().AddText(fmt.Sprintf("5.1.1.1 测试用例名称%s", vul["name"]))
								table := this.AddTable()
								table.Properties().SetWidth(5.7 * measurement.Inch)
								borders := table.Properties().Borders()
								borders.SetAll(wml.ST_BorderSingle, color.Auto, 2*measurement.Point)

								//漏洞信息
								row := table.AddRow()
								cell := row.AddCell()
								cell.AddParagraph().AddRun().AddText("漏洞名称")

								cell = row.AddCell()
								cell.AddParagraph().AddRun().AddText(vulName.(string))

								cell = row.AddCell()
								cell.AddParagraph().AddRun().AddText("漏洞ID")

								cell = row.AddCell()
								cell.AddParagraph().AddRun().AddText(vulId.(bson.ObjectId).Hex())

								row = table.AddRow()
								row.AddCell().AddParagraph().AddRun().AddText("漏洞等级")
								row.AddCell().AddParagraph().AddRun().AddText(vulLevelName)
								row.AddCell().AddParagraph().AddRun().AddText("修复状态")
								row.AddCell().AddParagraph().AddRun().AddText(vulStatusName)

								row = table.AddRow()
								cell = row.AddCell()
								cell.Properties().SetColumnSpan(2)
								cell.AddParagraph().AddRun().AddText("漏洞类型")
								cell = row.AddCell()
								cell.Properties().SetColumnSpan(2)
								cell.AddParagraph().AddRun().AddText(fmt.Sprintf("%v", vulRiskType))

								row = table.AddRow()
								cell = row.AddCell()
								cell.Properties().SetColumnSpan(2)
								cell.AddParagraph().AddRun().AddText("漏洞描述")
								cell = row.AddCell()
								cell.Properties().SetColumnSpan(2)
								cell.AddParagraph().AddRun().AddText(vulDescribe.(string))

								row = table.AddRow()
								cell = row.AddCell()
								cell.Properties().SetColumnSpan(2)
								cell.AddParagraph().AddRun().AddText("影响范围")
								cell = row.AddCell()
								cell.Properties().SetColumnSpan(2)
								cell.AddParagraph().AddRun().AddText(vulInfluence.(string))
							}
						}
						this.AddMain2().AddText("攻击链示意图")
						paramsRecord := qmap.QM{
							"e_item_id": itemId,
						}
						records := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_RECORD, paramsRecord).All()
						var testProcedure string
						var toolTestResult string
						var attachment map[string]interface{}
						for _, record := range records {
							if tmpTestProcedure, exit := record["test_procedure"]; exit {
								testProcedure += tmpTestProcedure.(string)
							}
							if tmpToolTestResult, exit := record["tool_test_result"]; exit {
								toolTestResult += tmpToolTestResult.(string)
							}
							if tmpAttachment, exit := record["attachment"]; exit {
								attachments := tmpAttachment.([]interface{})
								for _, tmpAttachment := range attachments {
									tmp := tmpAttachment.(map[string]interface{})
									for k, v := range tmp {
										attachment[k] = v
									}
								}
							}
						}
						this.AddMain2().AddText("测试过程")
						this.AddMain2().AddText(testProcedure)
						this.AddMain2().AddText("工具测试结果")
						this.AddMain2().AddText(toolTestResult)
						this.AddMain2().AddText("附件")
						this.AddMain2().AddText("")
						for k, v := range attachment {
							this.AddMain2().AddText(k)
							this.AddMain2().AddText(v.(string))
						}
					}
				}
			}
		}
	}
}

// 4版权
func (this *AssetReport) Copyright() {
	this.AddHeading1().AddText("版权")
	content := "本测试报告版权归北京奇虎科技有限公司智能网联汽车信息安全实验室所有，未经北京奇虎科技有限公司书面授权，" +
		"任何人不得以任何形式或者任何手段（电子，机械，显微复印，照相复印或其他手段）对其中内容的全部或者部分进行使用、复制、传输、公开传播，" +
		"泄露于其他任何第三方，或者将其存放在检索文件中。北京奇虎科技有限公司将保留对本文全部或者部分内容的盗用或者泄密责任追究权。"
	this.AddParagraph().AddRun().AddText(content)
}

// 5附录
func (this *AssetReport) Appendix() {
	this.Appendix1()
	this.Appendix2()
	this.Appendix3()
	this.Appendix4()
	this.Appendix5()
}
func (this *AssetReport) Appendix1() {
	this.AddTitle1().AddText("附录：漏洞等级评级标准")
	this.AddTable()
}
func (this *AssetReport) Appendix2() {
	this.AddTitle1().AddText("附录2：修复建议参考标准")
	table := this.AddTable()
	table.Properties().SetWidthPercent(100)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, measurement.Zero)

	row := table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("重要程度")
	row.AddCell().AddParagraph().AddRun().AddText("说明")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("必要")
	row.AddCell().AddParagraph().AddRun().AddText("在没有其他解决方案的情况下，优先考虑此项")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("推荐")
	row.AddCell().AddParagraph().AddRun().AddText("在其他方案无法满足的情况下，或者风险无法完全避免的情况下，可以考虑此项")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("可选")
	row.AddCell().AddParagraph().AddRun().AddText("在资金和时间都有尚余力的情况下，可以考虑此项")

	table = this.AddTable()
	table.Properties().SetWidthPercent(100)
	borders = table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, measurement.Zero)
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("修复成本")
	row.AddCell().AddParagraph().AddRun().AddText("说明")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("极高")
	row.AddCell().AddParagraph().AddRun().AddText("涉及多个供应商，现存案例少，生产工艺或系统架构需要重新设计，反复评审，消耗大量资源")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("高")
	row.AddCell().AddParagraph().AddRun().AddText("涉及多个供应商，现存案例多，部分设计需要修改，反复评审，消耗大量资源")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("中")
	row.AddCell().AddParagraph().AddRun().AddText("涉及1个供应商，需要经过评审，增加或修改部分设计，需要多人完成")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("低")
	row.AddCell().AddParagraph().AddRun().AddText("需要经过评审，讨论需求，进行修复，能在一周左右解决")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("极低")
	row.AddCell().AddParagraph().AddRun().AddText("无需评审，能在数日内解决")

}
func (this *AssetReport) Appendix3() {
	this.AddTitle1().AddText("附录3：安全等级评级标准")
	this.AddMain2().AddText("安全风险状况等级说明")

	table := this.AddTable()
	table.Properties().SetWidth(5.7 * measurement.Inch)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, measurement.Zero)
	row := table.AddRow()
	cell := row.AddCell()
	cell.Properties().SetColumnSpan(2)
	cell.AddParagraph().AddRun().AddText("安全风险状况说明")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("1")
	row.AddCell().AddParagraph().AddRun().AddText("良好状态\n信息系统处于良好运行状态，没有发现或只存在零星的低风险安全问题，此时只要保持现有安全策略就满足了本系统的安全等级要求。")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("2")
	row.AddCell().AddParagraph().AddRun().AddText("预警状态\n信息系统中存在一些漏洞或安全隐患，此时需根据评估中发现的网络、主机、应用和管理等方面的问题对进行有针对性的加固或改进。")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("3")
	row.AddCell().AddParagraph().AddRun().AddText("严重状态\n信息系统中发现存在严重漏洞或可能严重威胁到系统正常运行的安全问题，此时需要立刻采取措施，例如安装补丁或重新部署安全系统进行防护等等。")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("4")
	row.AddCell().AddParagraph().AddRun().AddText("紧急状态\n信息系统面临严峻的网络安全态势，对组织的重大经济利益或政治利益可能造成严重损害。此时需要与其他安全部门通力协作采取紧急防御措施。")
}
func (this *AssetReport) Appendix4() {
	this.AddTitle1().AddText("附录4：工具列表")
	this.AddMain2().AddText("安全风险状况等级说明")

	table := this.AddTable()
	table.Properties().SetWidth(5.7 * measurement.Inch)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, measurement.Zero)
	row := table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("工具名称")
	row.AddCell().AddParagraph().AddRun().AddText("工具简介")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("xxx")
	row.AddCell().AddParagraph().AddRun().AddText("xxxx")
}
func (this *AssetReport) Appendix5() {
	this.AddTitle1().AddText("附录5：参考文献")
	this.AddMain2().AddText("[1]刘国钧，陈绍业，王凤翥. 图书馆目录[M]. 北京：高等教育出版社，1957.15-18.[1]刘国钧，陈绍业，王凤翥. 图书馆目录[M]. 北京：高等教育出版社，1957.15-18.")
	this.AddMain2().AddText("[2]辛希孟. 信息技术和信息服务国际研讨会论文集：A集[C]. 北京：中国社会科学出版社，1994")
}
