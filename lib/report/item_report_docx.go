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

var reportTemplateConfig *service.ReportTemplateConfig

type ItemFirstReport struct {
	ReportCommont
}

func GetItemFirstReportDocx(projectId, reportType string, todayReportCount int, evaluateItems []string) (fileId string, err error) {
	reportTemplateConfig = service.LoadReportTemplateConfig()
	allData, err := getAllData(projectId, reportType, evaluateItems)
	if err != nil {
		panic(err)
	}
	if doc, err := document.OpenTemplate(reportTemplateConfig.ReportTemplate); err == nil {
		weekReport := WeekReportDocx{&ReportParse{doc}}
		itemFirstReport := ItemFirstReport{ReportCommont{weekReport}}
		itemFirstReport.addHeader()
		itemFirstReport.addFooter()
		itemFirstReport.Title(allData, reportType)
		itemFirstReport.Catalog()
		itemFirstReport.Statement(allData)
		itemFirstReport.Abbreviation()
		itemFirstReport.ItemSummarize(allData)
		itemFirstReport.Summarize(allData)
		itemFirstReport.PenetrationTest(allData)
		itemFirstReport.Appendix()
		itemFirstReport.SaveToFile("./tmp.docx")

		buffer := new(bytes.Buffer)
		itemFirstReport.Save(buffer)
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

func getAllData(projectId, reportType string, evaluateItems []string) (qmap.QM, error) {
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(projectId),
	}
	// 获取项目信息
	project, err := mongo.NewMgoSessionWithCond(common.MC_PROJECT, params).GetOne()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("project is %s", err.Error()))
	}
	// 获取测试用例信息
	var items []map[string]interface{}
	paramsItem := qmap.QM{}
	if len(evaluateItems) == 0 {
		paramsItem = qmap.QM{
			"e_project_id": projectId,
		}
	} else {
		paramsItem = qmap.QM{
			"in__id": evaluateItems,
		}
	}
	// 有些测试用例没有这条，所以先不添加
	// if reportType == common.RT_TEST {
	//	paramsItem["e_test_phase"] = 1
	// }
	mgoSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_ITEM, paramsItem)
	mgoSession.SetLimit(1000)
	items = mgoSession.All()
	// 通过测试用例获取资产id
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
	// 获取漏洞信息
	params = qmap.QM{
		"in_asset_id": assetIds,
	}
	ormSession := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VULNERABILITY, params)
	ormSession.SetTransformFunc(func(qm qmap.QM) qmap.QM {
		riskType := qm.Int("risk_type")
		params := qmap.QM{
			"e__id": riskType,
		}
		vulType, _ := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VUL_TYPE, params).GetOne()
		qm["risk_type_name"] = vulType.DefaultString("name", "未知")
		return qm
	})
	vuls := ormSession.All()
	// 获取车企的名字
	params = qmap.QM{
		"e__id": bson.ObjectIdHex(project.String("company")),
	}
	factory, err := mongo.NewMgoSessionWithCond(common.MC_FACTORY, params).GetOne()
	if err != nil {
		factory = &qmap.QM{"name": " 未知 "}
		// return nil, errors.New(fmt.Sprintf("factory is %s", err.Error()))
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

// 首页
func (this *ItemFirstReport) Title(allData qmap.QM, reportType string) {
	// 获取工程名称
	project := allData["project"]
	tmpProject := project.(*qmap.QM)
	name := tmpProject.String("name")
	this.AddMultiBreak(3)

	if reportType == "retest" {
		this.AddTitleWihtStyle().AddText(fmt.Sprintf("%s复测报告", name))
	} else {
		this.AddTitleWihtStyle().AddText(fmt.Sprintf("%s初测报告", name))
	}
	this.AddMultiBreak(20)
	niceNumber := niceNumber(1)
	this.AddTitle5().AddText(fmt.Sprintf("报告编号：SKYGO-%s-%s", time.Now().Format("20060102"), niceNumber))
	this.AddTitle5().AddText("报告提供商：360 SKY-GO智能网联汽车安全实验室")
	this.AddParagraph().AddRun().AddPageBreak()
}

// 目录
func (this *ItemFirstReport) Catalog() {
	this.AddTitleWihtStyle().AddText("目录")
	doc := this.Document
	doc.AddParagraph().AddRun().AddField(document.FieldTOC)
	doc.AddParagraph().AddRun().AddPageBreak()
}

// 1 声明
func (this *ItemFirstReport) Statement(allData qmap.QM) {
	this.AddParagraph()
	content1 := "本报告是针对 %s-%s 的安全评测报告。"
	content2 := "本报告测评结论的有效性建立在被测车联网系统提供相关证据的真实性基础之上。"
	content3 := "本报告中给出的测评结论仅对被测车联网系统当时的安全状态有效。当评测工作完成后，被测零部件因功能迭代升级而涉及到的相关组件发生变化，本报告将不再适用。"
	content4 := "在任何情况下，若需引用本报告中的测评结果或结论都应保持其原有的意义，不得对相关内容擅自进行增加、修改和伪造或掩盖事实。"
	project := allData["project"]
	factory := allData["factory"]
	tmpProject := project.(*qmap.QM)
	tmpFactory := factory.(*qmap.QM)
	// 获取车企名称
	name := tmpFactory.String("name")
	// 获取车型带号
	codeName := tmpProject.String("code_name")
	this.AddHeading1().AddText("声明")
	this.AddMain2().AddText(fmt.Sprintf(content1, name, codeName))
	this.AddMain2().AddText(content2)
	this.AddMain2().AddText(content3)
	this.AddMain2().AddText(content4)
}

// 2 缩写词汇表
func (this *ItemFirstReport) Abbreviation() {
	this.AddParagraph()
	this.AddHeading1().AddText("缩写词汇表")
	table := this.AddTable()
	table.Properties().SetWidth(5.7 * measurement.Inch)
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
}

// 3 项目综述
func (this *ItemFirstReport) ItemSummarize(allData qmap.QM) {
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
	codeName := project.String("code_name")

	// 3 项目综述
	this.AddHeading1().AddText("项目综述")
	// 3.1 项目目的
	this.AddHeading2().AddText("项目目的")
	this.AddMain2().AddText(fmt.Sprintf(content, companyName, codeName, companyName, companyName, companyName, codeName))
	// 3.2 项目范围
	this.AddHeading2().AddText("项目范围")
	this.AddMain().AddText("表1 资产表")
	table := this.AddTable()
	table.Properties().SetWidth(5.7 * measurement.Inch)
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
		if attributes == "" {
			attributes = "--"
		}
		cell.AddParagraph().AddRun().AddText(attributes)
	}
}

// 4总体评价
func (this *ItemFirstReport) Summarize(allData qmap.QM) {
	this.AddParagraph()
	tmpVuls := allData["vuls"]
	vuls := tmpVuls.([]map[string]interface{})
	vulNum := len(vuls)
	var serious int
	var hight int
	var middle int
	var low int
	var info int
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
		case 0:
			info++
		}
		switch vul["status"] {
		case 0:
			unrepair++
		case 1:
			repair++
		}
	}
	// 4 总体评价
	project := allData["project"]
	factory := allData["factory"]
	tmpProject := project.(*qmap.QM)
	tmpFactory := factory.(*qmap.QM)
	// 获取车企名称
	name := tmpFactory.String("name")
	// 获取车型名称
	brand := tmpProject.String("brand")
	this.AddHeading1().AddText("总体评价")
	// 4.1 项目总结论
	this.AddHeading2().AddText("项目总结论")
	content := "本项目%s提供的待测系统为车联网系统。" +
		"360 Sky-Go团队在整个渗透测试过程中测试内容涵盖应用安全、管理安全、通信安全、数据安全等内容，主要手段有动态调试、" +
		"逆向分析、漏洞扫描、渗透测试等手段，在%s %s架构环境下进行了渗透测试，并对其中存在的漏洞进行了综合分析评估。"
	this.AddMain2().AddText(fmt.Sprintf(content, name, brand))
	content = "在首次渗透测试中，" +
		"360 Sky-Go团队发现车联网系统的安全问题共计%v个（严重：%v个；高危：%v个；中危：%v个；低危：%v个；提示：%v个），" +
		"%v个已修复，%v个未修复，其中%v个漏洞，%v方表示可接受。因此我司评估%v联网系统的网络安全水平为XX。"
	this.AddMain2().AddText(fmt.Sprintf(content, vulNum, serious, hight, middle, low, info, repair, unrepair, vulNum, name, name))
	timeLayout := "2006年01月02日"
	start := tmpProject.Int("start_time")
	startTime := time.Unix(int64(start), 0).Format(timeLayout)
	endTime := time.Now().Format(timeLayout)
	content = "本次测试自%s至%s结束，整个测试实施过程共分为测试环境准备阶段、" +
		"初次测试阶段、初测报告编制阶段、复测阶段、复测报告编制阶段，共五个阶段。"
	this.AddMain2().AddText(fmt.Sprintf(content, startTime, endTime))

	// todo 甘特图还没弄
	// this.AddMain2().AddText("甘特图 无 ")

	// 4.2 漏洞统计与分布
	this.AddHeading2().AddText("漏洞统计与分布")
	// 如果不存在漏洞，就返回 本车型未发现漏洞
	if len(vuls) == 0 {
		this.AddMain2().AddText("本车型未发现漏洞")
		return
	}
	// 4.2.1 总统计图
	this.AddHeading3().AddText("总统计图")
	if serious == 0 && hight == 0 && middle == 0 && low == 0 && info == 0 {
		this.AddMain2().AddText("图1 漏洞等级分布 暂无")
	} else {
		pictureParams := qmap.QM{
			"red":    qmap.QM{"value": serious, "label": fmt.Sprintf("严重%d个", serious)},
			"yellow": qmap.QM{"value": hight, "label": fmt.Sprintf("高危%d个", hight)},
			"blue":   qmap.QM{"value": middle, "label": fmt.Sprintf("中危%d个", middle)},
			"orange": qmap.QM{"value": low, "label": fmt.Sprintf("低危%d个", low)},
			"gray":   qmap.QM{"value": info, "label": fmt.Sprintf("提示%d个", info)},
		}
		picture1 := this.NewPieImage(pictureParams, "图1", "")
		this.AddImageFromFile(picture1)
		this.AddMain2().AddText("图1 漏洞等级分布")
	}
	this.AddMain2().AddTab()
	// 资产漏洞
	tmpAssets := allData["assets"]
	assets := tmpAssets.([]map[string]interface{})
	for _, asset := range assets {
		assetId := asset["_id"]
		var serious int
		var hight int
		var middle int
		var low int
		var info int
		for _, vul := range vuls {
			if assetId == vul["asset_id"] {
				switch vul["level"] {
				case 1:
					low++
				case 2:
					middle++
				case 3:
					hight++
				case 4:
					serious++
				case 0:
					info++
				}
			}
		}
		tmpParams := qmap.QM{}
		if serious != 0 {
			tmpParams["red"] = qmap.QM{"value": serious, "label": "严重"}
		}
		if hight != 0 {
			tmpParams["yellow"] = qmap.QM{"value": hight, "label": "高危"}
		}
		if middle != 0 {
			tmpParams["blue"] = qmap.QM{"value": middle, "label": "中危"}
		}
		if low != 0 {
			tmpParams["orange"] = qmap.QM{"value": low, "label": "低危"}
		}
		if info != 0 {
			tmpParams["gray"] = qmap.QM{"value": info, "label": "提示"}
		}
		if serious == 0 && hight == 0 && middle == 0 && low == 0 && info == 0 {
			this.AddMain2().AddText(fmt.Sprintf("图2 %s漏洞等级分布 暂无", asset["name"]))
		} else {
			picture1 := this.NewDonutImage(tmpParams, "图2", "")
			this.AddImageFromFile(picture1)
			this.AddMain2().AddText(fmt.Sprintf("图2 %s漏洞等级分布", asset["name"]))
		}

	}
	this.AddMain2().AddTab()

	// 所有漏洞类型分布
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
		// 根据漏洞类型获取漏洞名称
		if v != 0 {
			isVul = true
		}
		tmp := qmap.QM{"value": float64(v), "label": vulType.String("name")}
		vulTypes = append(vulTypes, tmp)
	}
	if isVul {
		picture1 := this.NewPieImage1(vulTypes, "图3", "")
		this.AddImageFromFile(picture1)
		this.AddMain2().AddText("图3 漏洞类型分类")
	} else {
		this.AddMain2().AddText("图3 漏洞类型分类 暂无")
	}

	// 4.2.2 各资产统计图
	this.AddHeading3().AddText("各资产统计图")
	// todo 各个资产
	for _, asset := range assets {
		assetName := asset["name"]
		this.AddMain2().AddText(assetName.(string))

		assetId := asset["_id"]
		// 漏洞等级分布
		var serious int
		var hight int
		var middle int
		var low int
		var info int
		for _, vul := range vuls {
			if assetId == vul["asset_id"] {
				switch vul["level"] {
				case 1:
					low++
				case 2:
					middle++
				case 3:
					hight++
				case 4:
					serious++
				case 0:
					info++
				}
			}
		}
		if serious == 0 && hight == 0 && middle == 0 && low == 0 && info == 0 {
			this.AddMain2().AddText(fmt.Sprintf("图5 %s漏洞等级分布 暂无", assetName))
		} else {
			tmpParams := qmap.QM{
				"red":    qmap.QM{"value": serious, "label": "严重"},
				"yellow": qmap.QM{"value": hight, "label": "高危"},
				"blue":   qmap.QM{"value": middle, "label": "中危"},
				"orange": qmap.QM{"value": low, "label": "低危"},
				"gray":   qmap.QM{"value": info, "label": "提示"},
			}
			picture1 := this.NewPieImage(tmpParams, "图5", "")
			this.AddImageFromFile(picture1)
			this.AddMain2().AddText(fmt.Sprintf("图5 %s漏洞等级分布", assetName))
		}

		// 漏洞类型分布
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
			// 根据漏洞类型获取漏洞名称
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

	// 4.2.3 漏洞总表
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
}

// 5渗透测试详情
func (this *ItemFirstReport) PenetrationTest(allData qmap.QM) {
	// 5 渗透测试详情
	this.AddHeading1().AddText("渗透测试详情")
	// 获取资产
	tmpAssets := allData["assets"]
	assets := tmpAssets.([]map[string]interface{})
	// 资产
	for _, asset := range assets {
		assetId := asset["_id"]
		// 5.1 资产
		// this.AddMain2().AddText("工作环境架构逻辑图")
		// 测试组件
		moduleTypeIds := asset["module_type_id"]
		for _, moduleTypeId := range moduleTypeIds.([]interface{}) {
			paramsModule := qmap.QM{
				"e__id": bson.ObjectIdHex(moduleTypeId.(string)),
			}
			if modult, err := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_MODULE, paramsModule).GetOne(); err == nil {
				// moduleTitle := fmt.Sprintf("5.%v.%v 测试组件:%v", i+1, moduleIndex+1, modult.String("module_name"))
				moduleTitle := fmt.Sprintf("测试组件:%v", modult.String("module_name"))
				// 测试用例
				tmpItems := allData["items"]
				items := tmpItems.([]map[string]interface{})
				for _, item := range items {
					tmpModuleTypeId := item["module_type_id"]
					if assetId == item["asset_id"] && moduleTypeId.(string) == tmpModuleTypeId.(string) {
						// 5.1 资产
						assetTitle := fmt.Sprintf("资产:%s", asset["name"])
						this.AddHeading2().AddText(assetTitle)
						// 5.1.1 测试组件
						this.AddHeading3().AddText(moduleTitle)
						itemId := item["_id"]
						// this.AddMain2().AddText(fmt.Sprintf("测试用例:%s", item["name"]))
						// 5.1.1.1 测试用例
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
								// vulRiskType := vul["risk_type"]
								vulRiskTypeName := vul["risk_type_name"]
								vulDescribe := vul["describe"]
								vulInfluence := vul["influence"]
								// this.AddMain2().AddText(fmt.Sprintf("5.1.1.1 测试用例名称%s", vul["name"]))
								table := this.AddTable()
								table.Properties().SetWidth(5.7 * measurement.Inch)
								borders := table.Properties().Borders()
								borders.SetAll(wml.ST_BorderSingle, color.Auto, 2*measurement.Point)

								// 漏洞信息
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
								cell.AddParagraph().AddRun().AddText(fmt.Sprintf("%v", vulRiskTypeName))

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
						attachment := map[string]interface{}{}
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
	this.AddParagraph().AddRun().AddPageBreak()
}

// 4版权
func (this *ItemFirstReport) Copyright() {}

// 5附录
func (this *ItemFirstReport) Appendix() {
	this.Appendix1()
	this.Appendix2()
	this.Appendix3()
	this.Appendix4()
	this.Appendix5()
}
func (this *ItemFirstReport) Appendix1() {
	this.AddTitle1().AddText("附录1：漏洞等级评级标准")
	table := this.AddTable()
	table.Properties().SetWidth(5.7 * measurement.Inch)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, 1*measurement.Point)

	row := table.AddRow()
	cell := row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeRestart)
	cell.Properties().SetVerticalAlignment(wml.ST_VerticalJcCenter)
	cell.AddParagraph().AddRun().AddText("漏洞等级")
	row.AddCell().AddParagraph().AddRun().AddText("分类")
	row.AddCell().AddParagraph().AddRun().AddText("描述")

	row = table.AddRow()
	cell = row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
	cell.AddParagraph().AddRun().AddText("漏洞等级")
	row.AddCell().AddParagraph().AddRun().AddText("严重")
	row.AddCell().AddParagraph().AddRun().AddText("危及人身生命，或严重的经济损失，或社会/国家安全")

	row = table.AddRow()
	cell = row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
	cell.AddParagraph().AddRun().AddText("漏洞等级")
	row.AddCell().AddParagraph().AddRun().AddText("高危")
	row.AddCell().AddParagraph().AddRun().AddText("危及危害到系统安全")

	row = table.AddRow()
	cell = row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
	cell.AddParagraph().AddRun().AddText("漏洞等级")
	row.AddCell().AddParagraph().AddRun().AddText("中危")
	row.AddCell().AddParagraph().AddRun().AddText("对系统、服务和设备的性能和状态有一定影响")

	row = table.AddRow()
	cell = row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
	cell.AddParagraph().AddRun().AddText("漏洞等级")
	row.AddCell().AddParagraph().AddRun().AddText("低危")
	row.AddCell().AddParagraph().AddRun().AddText("对敏感信息的访问")

	row = table.AddRow()
	cell = row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
	cell.AddParagraph().AddRun().AddText("漏洞等级")
	row.AddCell().AddParagraph().AddRun().AddText("提示")
	row.AddCell().AddParagraph().AddRun().AddText("不是安全风险，推荐使用符合安全开发规范方式实现。")

}
func (this *ItemFirstReport) Appendix2() {
	this.AddTitle1().AddText("附录2：修复建议参考标准")
	table := this.AddTable()
	// table.Properties().SetWidthPercent(100)
	table.Properties().SetWidth(5.7 * measurement.Inch)
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

	this.AddBreak()

	table = this.AddTable()
	// table.Properties().SetWidthPercent(100)
	table.Properties().SetWidth(5.7 * measurement.Inch)
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
func (this *ItemFirstReport) Appendix3() {
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
	run := row.AddCell().AddParagraph().AddRun()
	run.AddText("良好状态")
	run.AddText("信息系统处于良好运行状态，没有发现或只存在零星的低风险安全问题，此时只要保持现有安全策略就满足了本系统的安全等级要求。")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("2")
	run = row.AddCell().AddParagraph().AddRun()
	run.AddText("预警状态")
	run.AddText("信息系统中存在一些漏洞或安全隐患，此时需根据评估中发现的网络、主机、应用和管理等方面的问题对进行有针对性的加固或改进。")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("3")
	run = row.AddCell().AddParagraph().AddRun()
	run.AddText("严重状态")
	run.AddText("信息系统中发现存在严重漏洞或可能严重威胁到系统正常运行的安全问题，此时需要立刻采取措施，例如安装补丁或重新部署安全系统进行防护等等。")
	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("4")
	run = row.AddCell().AddParagraph().AddRun()
	run.AddText("紧急状态")
	run.AddText("信息系统面临严峻的网络安全态势，对组织的重大经济利益或政治利益可能造成严重损害。此时需要与其他安全部门通力协作采取紧急防御措施。")
}
func (this *ItemFirstReport) Appendix4() {
	this.AddTitle1().AddText("附录4：工具列表")

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
func (this *ItemFirstReport) Appendix5() {
	this.AddTitle1().AddText("附录5：参考文献")

	table := this.AddTable()
	table.Properties().SetWidth(5.7 * measurement.Inch)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Auto, measurement.Zero)
	row := table.AddRow()
	para := row.AddCell().AddParagraph()
	para.Properties().SetAlignment(wml.ST_JcCenter)
	run := para.AddRun()
	run.Properties().SetBold(true)
	run.AddText("序号")
	para = row.AddCell().AddParagraph()
	para.Properties().SetAlignment(wml.ST_JcCenter)
	run = para.AddRun()
	run.Properties().SetBold(true)
	run.AddText("标准名称")
	standardNames := []string{
		"《ISO26262汽车电子电气系统的功能安全标准》", "《T/ITS 0054-2016基于公众电信网的汽车网关检测方法》", "《T/ITS 0057-2016车辆安全通信系统 设备性能测试规程》",
		"《GBT 20984-2007信息安全技术 信息安全风险评估规范》", "《ISO/IEC TR 13335信息安全管理指南》", "《Owasp Testing Guide v4测试指南》",
		"《Owsap ASVS 3.0 WEB应用安全评估标准》", "《移动互联网应用软件安全评估大纲》", "《中国国家信息安全漏洞库CNNVD》",
		"《OWASP Mobile Top 10 2017》", "《OWASP 安全检测指南》", "《GBT 36627-2018 信息安全技术 网络安全等级保护测试评估技术指南》",
		"《ISO-11898 道路车辆 控制器局域网络（CAN）》", "《ISO-14229 统一诊断服务》", "《ISO 26262 道路车辆功能安全》",
		"《GB 26149-2017 乘用车轮胎气压监测系统的性能要求和试验方法》", "《GB/T 12572-2008 无线电发射设备参数通用要求和测量方法》", "《GB/T 18314-2009 全球定位系统（GPS）测量规范》",
		"《GB/T 20008-2005 信息安全技术 操作系统安全评估准则》", "《GB/T 20271 信息安全技术 信息系统通用 安全技术要求》", "《GB/T 20277-2006 信息安全技术 网络和终端设备隔离部件测试评价方法》",
		"《GB/T 22186-2008 信息安全技术 具有中央处理器的集成电路（IC）卡芯片安 全技术要求(评估保证级 4 增强级)》", "《GB/T 25068.3-2010 信息技术 安全技术 IT 网络安全 第 3 部分：使用安 全网关的网间通信安全保护》",
		"《GB/T 26256-2010 2.4GHz 频段无线电通信 设备的相互干扰限制与共存要 求及测试方法》", "《GB/T 30284-2013 信息安全技术 移动通信智能终端操作系统安全技术要求》",
		"《GB/T 30290.3-2013 卫星定位车辆信息服务系统 第 3 部分：信息安全规范》", "《GB/T 30976.1-2014 工业控制系统信息安全 第 1 部分：评估规范》", "《GB/T 30976.2-2014 工业控制系统信息安全 第 2 部分：验收规范》",
		"《GB/T 32415-2015 GSM/CDMA/WCDMA 数字蜂窝移动通信网塔顶放大器 技术指标和测试方法》", "《GB/T 32420-2015 无线局域网测试规范 汽车产品信息安全测试评价规范》",
		"《GB/T 34975-2017 信息安全技术 移动智能终端应用软件安全技术要求和测试评价方法》", "《GB/T 34976-2017 信息安全技术 移动智能终端操作系统安全技术要求和测 试评价方法》",
		"《GB/T 34977-2017 信息安全技术 移动智能终端数据存储安全技术要求与测 试评价方法》", "《GB/T 35291-2017 信息安全技术 智能密码钥匙应用接口规范》",
		"《IEEE 802.15 基于蓝牙的局域网标准》", "《IEEE 802.11 无线局域网标准》", "《YD/T 2585-2016 互联网数据中心安全防护检测要求》",
		"《YD/T 2408-2013 移动智能终端安全能力测试方法》", "《YD/T 883-1996 900MHz TDMA 数字蜂窝移动通信网 基站无线设备技术指标及测试方法》",
	}
	for i := 1; i <= 39; i++ {
		serialNumber := fmt.Sprintf("%v", i)
		standardName := standardNames[i-1]
		row = table.AddRow()
		row.AddCell().AddParagraph().AddRun().AddText(serialNumber)
		row.AddCell().AddParagraph().AddRun().AddText(standardName)
	}
}
