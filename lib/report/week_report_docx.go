package report

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/globalsign/mgo/bson"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/common"
	"skygo_detection/custom_util"
	"skygo_detection/lib/common_lib/mongo"
	"skygo_detection/lib/unioffice/color"
	unioffice_common "skygo_detection/lib/unioffice/common"
	"skygo_detection/lib/unioffice/document"
	"skygo_detection/lib/unioffice/measurement"
	"skygo_detection/lib/unioffice/schema/soo/wml"
	"skygo_detection/service"
)

type WeekReportDocx struct {
	*ReportParse
}

type Report struct {
	ProjectId string
	VulItem   []*WeekReportTarget
	OtherItem []*WeekReportDocx
}

type WeekReportTarget struct{}

//const WEEK_REPORT_DOCX_TEMPLATE = "/report_template/week_report_template.docx"
//const TMP_TEMPLATE = "/home/q/system/services/8100/project_manage/report_template/template_item_report.docx"
//const TMP_TEMPLATE = "./report_template/template_item_report.docx"
//const HEARDER = "./report_template/header.png"
//const FOOTER = "./report_template/footer.png"
//const HEARDER = "/home/q/system/services/8100/project_manage/report_template/header.png"
//const FOOTER = "/home/q/system/services/8100/project_manage/report_template/footer.png"
//const fontPath = "/home/q/system/services/8100/project_manage/report_template/Microsoft Yahei.ttf"
//const fontPath = "./report_template/Microsoft Yahei.ttf"
//const outputImage = "/home/q/system/services/8100/project_manage/report/image/%s.png"
//const outputImage = "./report/image/%s.png"

/**
{
    "item": {
        "list": {
            "5f86c60bf98f924de4d98456": {
                "0": {
                    "_id": "5f87af5d24b64741e994b113",
                    "create_time": 1602727773471,
                    "evaluate_type": "test",
                    "has_vulnerability": 1,
                    "item_vulnerability": [
                        "5f87af5d24b64741e994b114",
                        "5f87af5d24b64741e994b115"
                    ],
                    "last_update_op_id": 1,
                    "level": 2,
                    "matters_needing_attention": "",
                    "module_name": "测试组件",
                    "module_type": "测试分类1",
                    "name": "测试项1",
                    "objective": "11",
                    "op_id": 1,
                    "procedure": "##### 11111\nljadlsjfaljdfjadlfjl\n![image.png](/api/v1/project_file/image?file_id=5f86e4af24b6473f3747d290)\n1. s11111\n2. 222222222\n3. 333333\n\n![image.png](/api/v1/project_file/image?file_id=5f86e4d824b6473f3747d292)\n- 42352435243\n- 55555\n- sfgsg",
                    "project_id": "5f85808324b6472dd945f9da",
                    "repair_effect": "",
                    "result": "测试结果",
                    "status": 0,
                    "tag": [],
                    "target_id": "5f86c60bf98f924de4d98456",
                    "update_time": 1602727773471
                }
            }
        },
        "target": [
            {
                "id": "5f86c60bf98f924de4d98456",
                "name": "哈哈哈"
            }
        ]
    },
    "project": {
        "id": "5f85808324b6472dd945f9da",
        "name": "test-lq"
    },
    "target": {
        "list": [
            {
                "_id": "5f86c60bf98f924de4d98456",
                "attributes": {
                    "测试": "给"
                },
                "create_time": 1602668043865,
                "evaluate_type": "test",
                "name": "哈哈哈",
                "op_id": 23,
                "project_id": "5f85808324b6472dd945f9da",
                "tags": [],
                "update_time": 1602668043865
            }
        ]
    },
    "vul": {
        "list": {
            "5f86c60bf98f924de4d98456": {
                "0": {
                    "_id": "5f87af5d24b64741e994b115",
                    "describe": "带带 大扥",
                    "importance": 1,
                    "influence": "安赛飞",
                    "last_update_op_id": 1,
                    "level": 2,
                    "matters_needing_attention": "111",
                    "name": "漏洞2",
                    "op_id": 1,
                    "project_id": "5f85808324b6472dd945f9da",
                    "repair_cost": 0,
                    "repair_effect": 2,
                    "retest_record": [],
                    "risk_type": 1,
                    "status": 0,
                    "suggest": "撒扥as啊",
                    "target_id": "5f86c60bf98f924de4d98456"
                },
                "1": {
                    "_id": "5f87af5d24b64741e994b114",
                    "describe": "阿塞阀赛风",
                    "importance": 1,
                    "influence": "112123",
                    "last_update_op_id": 1,
                    "level": 3,
                    "matters_needing_attention": "扥安抚",
                    "name": "漏洞1",
                    "op_id": 1,
                    "project_id": "5f85808324b6472dd945f9da",
                    "repair_cost": 1,
                    "repair_effect": 0,
                    "retest_record": [],
                    "risk_type": 3,
                    "status": 0,
                    "suggest": "ad发",
                    "target_id": "5f86c60bf98f924de4d98456"
                }
            }
        },
        "target": [
            {
                "id": "5f86c60bf98f924de4d98456",
                "name": "哈哈哈"
            }
        ]
    },
    "vulCount": {
        "target": [
            {
                "id": "5f86c60bf98f924de4d98456",
                "name": "哈哈哈"
            }
        ],
        "total": {
            "5f86c60bf98f924de4d98456": {
                "0": 0,
                "1": 0,
                "2": 1,
                "3": 1,
                "4": 0
            }
        }
    }
}
*/

var ROOT_DIR = "."

func GetWeekReportDocx(reportInfo map[string]qmap.QM, todayReportCount int) (fileId string, err error) {
	if doc, err := document.OpenTemplate(fmt.Sprintf("%s%s", ROOT_DIR, reportTemplateConfig.ReportTemplate)); err == nil {
		weekReport := WeekReportDocx{&ReportParse{doc}}
		//weekReport.addDef()
		//weekReport.addHeader()
		//weekReport.addFooter()
		weekReport.addHomePage(todayReportCount)
		//weekReport.addStatement(reportInfo["project"])
		//weekReport.addTotalEvaluate(reportInfo)
		//weekReport.addToc()
		//weekReport.addOverView(reportInfo["target"])
		//weekReport.addTotalSecure(reportInfo["vulCount"])
		//weekReport.addVulDetail(reportInfo["vul"])
		//weekReport.addItem(reportInfo["item"])
		weekReport.addNoteA()
		weekReport.addNoteB()
		weekReport.addNoteC()
		weekReport.addNoteD()
		buffer := new(bytes.Buffer)
		weekReport.Save(buffer)
		if fileContent, err := ioutil.ReadAll(buffer); err == nil {
			fileName := GenerateReportName(reportInfo["project"]["id"].(string), "test", todayReportCount)
			if fileId, err := mongo.GridFSUpload(common.MC_File, fileName, fileContent); err == nil {
				return fileId, nil
			} else {
				return "", err
			}
		} else {
			return "", err
		}
	} else {
		return "", err
	}
}

func (this *WeekReportDocx) addHomePage(index int) {
	this.AddTitle().AddText("《XXX项目 -XXX信息安全测试报告》")
	this.AddBreak()
	niceNumber := niceNumber(index)
	this.AddTitle5().AddText(fmt.Sprintf("报告编号：SKYGO-%s-%s", time.Now().Format("20060102"), niceNumber))
	this.AddBreak()

	para := this.AddParagraph()
	para.SetStyle("R-title5")
	para.AddRun().AddText("实施方案说明   版本号")
	run := para.AddRun()
	run.AddText(fmt.Sprintf("V1.0.0 %s", time.Now().Format("2006")))
	run.Properties().SetColor(color.OrangeRed)

	this.AddTitle5().AddText("Proposal of Project Implementation")
	this.AddMultiBreak(4)
	run = this.AddIndent()
	run.Properties().SetBold(true)
	run.AddText("本测试报告版权归北京奇虎科技有限公司智能网联汽车信息安全实验室所有，未经北京奇虎科技有限公司书面授权，任何人不得以任何形式或者任何手段（电子，机械，显微复印，照相复印或其他手段）对其中内容的全部或者部分进行使用、复制、传输、公开传播，泄露于其他任何第三方，或者将其存放在检索文件中。北京奇虎科技有限公司将保留对本文全部或者部分内容的盗用或者泄密责任追究权。")
	this.AddMultiBreak(6)

	para = this.AddParagraph()
	para.SetStyle("R-main3")
	run = para.AddRun()
	run.AddText("报告提供商：")
	run.Properties().SetBold(true)
	run.Properties().SetColor(color.OrangeRed)
	run = para.AddRun()
	run.Properties().SetBold(true)
	run.AddText("北京奇虎科技有限公司")

	para = this.AddParagraph()
	para.SetStyle("R-main3")
	run = para.AddRun()
	run.AddText("报告时间：")
	run.Properties().SetBold(true)
	run.Properties().SetColor(color.OrangeRed)
	run = para.AddRun()
	run.Properties().SetBold(true)
	run.AddText(time.Now().Format("2006 年 1 月"))
	this.AddIndent().AddPageBreak()
}

func (this *WeekReportDocx) addStatement(projectInfo qmap.QM) {
	this.AddMiddle(17).AddText("声 明")
	this.AddBreak()
	this.AddIndent().AddText(fmt.Sprintf("本报告是针对 %s 的安全评测报告。", projectInfo["name"]))
	this.AddIndent().AddText(fmt.Sprintf("本报告测评结论的有效性建立在被测 %s 提供相关证据的真实性基础之上。", projectInfo["name"]))
	this.AddIndent().AddText(fmt.Sprintf("本报告中给出的测评结论仅对被测 %s 当时的安全状态有效。当评测工作完成后，被测零部件因功能迭代升级而涉及到的相关组件发生变化，本报告将不再适用。", projectInfo["name"]))
	this.AddIndent().AddText("在任何情况下，若需引用本报告中的测评结果或结论都应保持其原有的意义，不得对相关内容擅自进行增加、修改和伪造或掩盖事实。")
	this.AddMultiBreak(10)
	this.AddLeft().AddText("北京奇虎科技有限公司")
	this.AddLeft().AddText(time.Now().Format("2006 年 1 月 2 日"))
	this.AddIndent().AddPageBreak()
}

func (this *WeekReportDocx) addTotalEvaluate(reportInfo map[string]qmap.QM) {
	targetData := reportInfo["target"]["list"]
	targetLink := ""
	for _, target := range targetData.([]map[string]interface{}) {
		targetLink += target["name"].(string) + " "
	}

	vulCount := reportInfo["vulCount"]["total"]
	total, heavy, high, mid, low, tips := 0, 0, 0, 0, 0, 0
	for _, countItem := range vulCount.(qmap.QM) {
		heavy += countItem.(map[int]int)[4]
		high += countItem.(map[int]int)[3]
		mid += countItem.(map[int]int)[2]
		low += countItem.(map[int]int)[1]
		tips += countItem.(map[int]int)[0]
	}
	total = heavy + high + mid + low

	this.AddMiddle(15).AddText("总体评价")
	this.AddBreak()
	this.AddIndent().AddText(fmt.Sprintf("本项目XX提供的待测零部件为 %s。360 Sky-Go团队在整个安全测试过程中测试内容涵盖{{此处为项目实施人员编写}}，主要手段有动态调试、逆向分析、漏洞扫描、渗透测试等手段，在 XX 架构环境下进行了安全测试，并对其中存在的漏洞进行了综合分析评估。", targetLink))
	this.AddIndent().AddText(fmt.Sprintf("在首次安全测试中，360 Sky-Go团队发现 %s 的安全问题共计 %d 个（严重： %d 个；高危： %d 个；中危： %d 个；低危： %d 个；提示： %d 个），包括通讯数据唯一性、安全秘钥硬编码、APP未验证用户凭证、访问用户主页泄露注册用户手机号等安全问题。", targetLink, total, heavy, high, mid, low, tips))
	para := this.AddParagraph()
	para.Properties().SetFirstLineIndent(2 * measurement.ChChar)
	para.AddRun().AddText("综上所述，首测结束后，360 Sky-Go团队发现多个中危安全问题，因此我们评估分析XX安全状态为 ")
	run := para.AddRun()
	run.AddText("预警状态")
	run.Properties().SetColor(color.DarkBlue)
	para.AddRun().AddText("。 ")
	this.AddMultiBreak(2)
	this.AddIndent().AddText("{{此处为项目实施人员编写}}")
	this.AddIndent().AddPageBreak()

}

func niceNumber(number int) string {
	prefix := "0"
	if number < 10 {
		prefix = "00"
	} else if number >= 100 {
		prefix = ""
	}
	return fmt.Sprintf("%s%d", prefix, number)
}

func (this *WeekReportDocx) addOverView(targetData qmap.QM) {
	this.AddHeading1().AddText("项目综述")
	this.AddHeading2().AddText("项目目的")
	this.AddMain().AddText("本次测试通过对XXXXXX进行安全测试，及时发现XXXX中存在的安全问题；然后将安全测试的结果以及相应安全修复建议反馈给XX，用于指导XX组织开展安全修复工作；最后通过对修复后的XX（测试对象）进行安全复测，确认XX（测试对象）安全问题修复的有效性，最终达到提升XX（测试对象或车联网）的安全。")

	this.AddHeading2().AddText("测试对象")
	var target qmap.QM
	for _, target = range (targetData)["list"].([]map[string]interface{}) {
		this.AddMain().AddText(target.String("name"))
	}
	this.AddHeading2().AddText("测试过程")
	table := this.AddTable()
	table.Properties().SetWidthPercent(100)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Gray, 1*measurement.Point)
	row := table.AddRow()
	cell := row.AddCell()
	cell.AddParagraph().AddRun().AddText("测试阶段")
	cell.Properties().SetWidthPercent(18)
	row.AddCell().AddParagraph().AddRun().AddText("测试内容")
	cell = row.AddCell()
	cell.AddParagraph().AddRun().AddText("开始时间")
	cell.Properties().SetWidthPercent(13)
	cell = row.AddCell()
	cell.AddParagraph().AddRun().AddText("结束时间")
	cell.Properties().SetWidthPercent(13)

	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("测试环境准备")
	row.AddCell().AddParagraph().AddRun().AddText("  本阶段360 Sky-Go团队对被测零部件进行了详细调研，根据掌握的被测零部件的基本情况搭建了相应的台架测试环境。")
	row.AddCell().AddParagraph().AddRun().AddText("")
	row.AddCell().AddParagraph().AddRun().AddText("")

	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("零部件测试")
	row.AddCell().AddParagraph().AddRun().AddText("  本阶段360 Sky-Go团队通过对被测零部件进行了全面的安全评估，验证被测零部件是否存在安全风险点。")
	row.AddCell().AddParagraph().AddRun().AddText("")
	row.AddCell().AddParagraph().AddRun().AddText("")

	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("实车测试")
	row.AddCell().AddParagraph().AddRun().AddText("  本阶段360 Sky-Go 团队通过对被测零部件进行了 实车验证测试 ，验证被测零部件 对整车的实际影响程度 。")
	row.AddCell().AddParagraph().AddRun().AddText("")
	row.AddCell().AddParagraph().AddRun().AddText("")

	row = table.AddRow()
	row.AddCell().AddParagraph().AddRun().AddText("编写初测报告")
	row.AddCell().AddParagraph().AddRun().AddText("  360 Sky-Go 团队通过现场测试阶段 和实车验证测试阶段 获得的测试证据，经整体测试分析和风险分析，最终形成《安全测试报告》。")
	row.AddCell().AddParagraph().AddRun().AddText("")
	row.AddCell().AddParagraph().AddRun().AddText("")

	this.AddMultiBreak(2)
}

func (this *WeekReportDocx) addTotalSecure(vulCount qmap.QM) {
	this.AddHeading1().AddText("整体安全测试情况")
	this.AddHeading2().AddText("总体安全情况统计")
	vulTotal := vulCount["total"].(qmap.QM)
	for _, target := range vulCount["target"].([]qmap.QM) {
		//添加表格
		table := this.AddTable()
		table.Properties().SetWidthPercent(100)
		borders := table.Properties().Borders()
		borders.SetAll(wml.ST_BorderSingle, color.Gray, 1*measurement.Point)
		targetName := target.String("name")
		targetId := target.String("id")

		row := table.AddRow()
		cell := row.AddCell()
		cell.AddParagraph().AddRun().AddText(targetName)
		cell.Properties().SetVerticalMerge(wml.ST_MergeRestart)
		cell.Properties().SetVerticalAlignment(wml.ST_VerticalJcCenter)
		cell.Properties().SetWidthPercent(15)

		row.AddCell().AddParagraph().AddRun().AddText("漏洞等级")
		row.AddCell().AddParagraph().AddRun().AddText("初测漏洞数量")
		row.AddCell().AddParagraph().AddRun().AddText("安全修复漏洞数量")
		row.AddCell().AddParagraph().AddRun().AddText("安全复测漏洞数量")

		tmpVulTotal := vulTotal[targetId].(map[int]int)
		for i := 0; i < len(tmpVulTotal); i++ {
			row := table.AddRow()
			cell := row.AddCell()
			cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
			cell.AddParagraph().AddRun().AddText("")
			row.AddCell().AddParagraph().AddRun().AddText(this.convertLevel(i))
			row.AddCell().AddParagraph().AddRun().AddText(fmt.Sprintf("%d", tmpVulTotal[i]))
			row.AddCell().AddParagraph().AddRun().AddText("0")
			row.AddCell().AddParagraph().AddRun().AddText("0")
		}
		//for level, count := range vulTotal[targetId].(map[int]int) {
		//	row := table.AddRow()
		//	cell := row.AddCell()
		//	cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
		//	cell.AddParagraph().AddRun().AddText("")
		//	row.AddCell().AddParagraph().AddRun().AddText(this.convertLevel(level))
		//	row.AddCell().AddParagraph().AddRun().AddText(fmt.Sprintf("%d", count))
		//	row.AddCell().AddParagraph().AddRun().AddText("0")
		//	row.AddCell().AddParagraph().AddRun().AddText("0")
		//}
		this.AddBreak()
	}

}

func (this *WeekReportDocx) addVulDetail(vul qmap.QM) {
	vulList := vul["list"].(map[string]map[string]qmap.QM)
	for _, target := range vul["target"].([]qmap.QM) {
		targetName := target.String("name")
		targetId := target.String("id")
		vulItem := vulList[targetId]
		this.AddHeading2().AddText(targetName)
		this.AddMain2().AddText("测试对象为小节")
		if len(vulItem) == 0 {
			this.AddMain2().AddText("未发现漏洞")
			continue
		}
		{
			//添加表格
			table := this.AddTable()
			table.Properties().SetWidthPercent(100)
			borders := table.Properties().Borders()
			borders.SetAll(wml.ST_BorderSingle, color.Gray, 1*measurement.Point)
			row := table.AddRow()
			cell := row.AddCell()
			cell.Properties().SetWidthPercent(10)
			cell.AddParagraph().AddRun().AddText("序号")
			cell = row.AddCell()
			cell.Properties().SetWidthPercent(40)
			cell.AddParagraph().AddRun().AddText("初次测试发现的主要问题")
			cell = row.AddCell()
			cell.Properties().SetWidthPercent(15)
			cell.AddParagraph().AddRun().AddText("危害程度")
			cell = row.AddCell()
			cell.Properties().SetWidthPercent(35)
			cell.AddParagraph().AddRun().AddText("备注")

			line := 1
			for _, item := range vulItem {
				row := table.AddRow()
				row.AddCell().AddParagraph().AddRun().AddText(strconv.Itoa(line))
				row.AddCell().AddParagraph().AddRun().AddText(item["name"].(string))
				row.AddCell().AddParagraph().AddRun().AddText(this.convertLevel(item["level"].(int)))
				row.AddCell().AddParagraph().AddRun().AddText("")
				line++
			}
		}

	}
	this.AddBreak()
}

func (this *WeekReportDocx) addItem(itemData qmap.QM) {
	itemList := itemData["list"].(map[string]map[string]qmap.QM)
	for _, target := range itemData["target"].([]qmap.QM) {
		targetName := target.String("name")
		targetId := target.String("id")
		items := itemList[targetId]
		this.AddHeading1().AddText(fmt.Sprintf("%s测试详情", targetName))

		itemType := map[string]map[string]qmap.QM{}
		for typeIndex, line := range items {
			if itemType[line["module_type"].(string)] == nil {
				itemType[line["module_type"].(string)] = map[string]qmap.QM{}
			}
			itemType[line["module_type"].(string)][typeIndex] = line
		}

		for typeName, lineDataMap := range itemType {
			this.AddHeading2().AddText(typeName)

			for _, lineData := range lineDataMap {
				this.AddHeading3().AddText(lineData["name"].(string))

				run := this.AddMain2()
				run.AddText("测试目的")
				run.Properties().SetBold(true)
				if lineData["objective"].(string) == "" {
					this.AddMain2().AddText("无")
				} else {
					this.AddMain2().AddText(lineData["objective"].(string))
				}

				run = this.AddMain2()
				run.AddText("测试步骤")
				run.Properties().SetBold(true)
				this.ParseMarkdown(lineData["procedure"].(string))
				this.AddBreak()

				run = this.AddMain2()
				run.AddText("测试结果")
				run.Properties().SetBold(true)

				if lineData["has_vulnerability"] == 0 {
					this.AddMain2().AddText("通过以上测试，未发现问题")
				} else {
					this.AddMain2().AddText("通过以上测试，发现如下漏洞：")

					//添加表格
					table := this.AddTable()
					table.Properties().SetWidthPercent(100)
					borders := table.Properties().Borders()
					borders.SetAll(wml.ST_BorderSingle, color.Gray, 1*measurement.Point)

					row := table.AddRow()
					row.AddCell().AddParagraph().AddRun().AddText("序号")
					row.AddCell().AddParagraph().AddRun().AddText("漏洞名称")
					row.AddCell().AddParagraph().AddRun().AddText("漏洞等级")

					vulList := map[string]qmap.QM{}
					for index, vulId := range lineData["item_vulnerability"].([]interface{}) {
						vul, vulErr := mongo.NewMgoSessionWithCond(common.MC_EVALUATE_VULNERABILITY, qmap.QM{"e__id": bson.ObjectIdHex(vulId.(string))}).GetOne()
						vulList[vulId.(string)] = *vul
						custom_util.CheckErr(vulErr)
						row := table.AddRow()
						row.AddCell().AddParagraph().AddRun().AddText(strconv.Itoa(index + 1))
						row.AddCell().AddParagraph().AddRun().AddText(vul.String("name"))
						row.AddCell().AddParagraph().AddRun().AddText(this.convertLevel(vul.Int("level")))
					}

					for _, vulItem := range vulList {
						this.AddHeading4().AddText(vulItem["name"].(string))
						{
							//添加表格
							table := this.AddTable()
							table.Properties().SetStyle("skygo")
							table.Properties().SetWidthPercent(100)

							row := table.AddRow()
							cell := row.AddCell()
							cell.AddParagraph().AddRun().AddText("漏洞属性")
							cell.Properties().SetColumnSpan(3)

							row = table.AddRow()
							cell = row.AddCell()
							cell.AddParagraph().AddRun().AddText("基本信息")
							cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
							cell.Properties().SetVerticalAlignment(wml.ST_VerticalJcCenter)
							cell = row.AddCell()
							cell.AddParagraph().AddRun().AddText("测试对象")
							cell = row.AddCell()
							cell.AddParagraph().AddRun().AddText(targetName)

							row = table.AddRow()
							cell = row.AddCell()
							cell.AddParagraph().AddRun().AddText(" ")
							cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
							cell = row.AddCell()
							cell.AddParagraph().AddRun().AddText("测试负责人")
							cell = row.AddCell()
							cell.AddParagraph().AddRun().AddText("")

							row = table.AddRow()
							cell = row.AddCell()
							cell.AddParagraph().AddRun().AddText(" ")
							cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
							cell = row.AddCell()
							cell.AddParagraph().AddRun().AddText("风险根源类型")

							riskType := "配置"
							if vulItem["risk_type"].(int) == 1 {
								riskType = "设计"
							} else if vulItem["risk_type"].(int) == 2 {
								riskType = "代码"
							} else if vulItem["risk_type"].(int) == 3 {
								riskType = "其他"
							}

							cell = row.AddCell()
							cell.AddParagraph().AddRun().AddText(riskType)

							row = table.AddRow()
							cell = row.AddCell()
							cell.AddParagraph().AddRun().AddText("漏洞评级")
							cell = row.AddCell()
							cell.AddParagraph().AddRun().AddText(this.convertLevel(vulItem["level"].(int)))
							cell.Properties().SetColumnSpan(2)
						}

						this.AddBreak()
						run = this.AddMain2()
						run.AddText("问题/漏洞描述")
						run.Properties().SetBold(true)
						this.ParseMarkdown(vulItem["describe"].(string))

						run = this.AddMain2()
						run.AddText("修复建议")
						run.Properties().SetBold(true)

						//添加表格
						{
							table = this.AddTable()
							table.Properties().SetStyle("skygo")
							table.Properties().SetWidthPercent(100)
							row := table.AddRow()
							row.AddCell().AddParagraph().AddRun().AddText("序号")
							row.AddCell().AddParagraph().AddRun().AddText("注意事项")
							row.AddCell().AddParagraph().AddRun().AddText("修复效果")
							row.AddCell().AddParagraph().AddRun().AddText("修复成本")
							row.AddCell().AddParagraph().AddRun().AddText("重要程度")
							row = table.AddRow()
							row.AddCell().AddParagraph().AddRun().AddText("1)")
							MattersNeedingAttention := ""
							if vulItem["matters_needing_attention"] != nil {
								MattersNeedingAttention = vulItem["matters_needing_attention"].(string)
							}

							row.AddCell().AddParagraph().AddRun().AddText(MattersNeedingAttention)
							RepairEffect := "彻底"
							if vulItem["repair_effect"].(int) == 1 {
								RepairEffect = "显著"
							} else if vulItem["repair_effect"].(int) == 2 {
								RepairEffect = "基础"
							}
							row.AddCell().AddParagraph().AddRun().AddText(RepairEffect)
							RepairCost := "极低"
							if vulItem["repair_cost"].(int) == 1 {
								RepairCost = "低"
							} else if vulItem["repair_cost"].(int) == 2 {
								RepairCost = "中"
							} else if vulItem["repair_cost"].(int) == 3 {
								RepairCost = "高"
							} else if vulItem["repair_cost"].(int) == 4 {
								RepairCost = "极高"
							}
							row.AddCell().AddParagraph().AddRun().AddText(RepairCost)
							inportance := "必要"
							if vulItem["importance"].(int) == 1 {
								inportance = "推荐"
							} else if vulItem["importance"].(int) == 2 {
								inportance = "可选"
							}
							row.AddCell().AddParagraph().AddRun().AddText(inportance)
						}

						this.AddBreak()
						this.ParseMarkdown(vulItem["suggest"].(string))
					}
				}
			}

		}

	}
	this.AddMultiBreak(2)
}

func (this *WeekReportDocx) convertLevel(level int) string {
	levelMap := map[int]string{
		0: "提示",
		1: "低危",
		2: "中危",
		3: "高危",
		4: "严重",
	}
	return levelMap[level]
}

// 添加附录A
func (this *WeekReportDocx) addNoteA() {
	this.AddAppendix().AddText("附录A：漏洞风险评级")

	//添加表格
	table := this.AddTable()
	table.Properties().SetWidthPercent(100)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Gray, 1*measurement.Point)

	row := table.AddRow()
	cell := row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
	cell.Properties().SetVerticalAlignment(wml.ST_VerticalJcCenter)
	cell.AddParagraph().AddRun().AddText("漏洞类型")
	cell.Properties().SetWidthPercent(15)
	cell = row.AddCell()
	cell.AddParagraph().AddRun().AddText("问题分类")
	cell.Properties().SetWidthPercent(15)
	row.AddCell().AddParagraph().AddRun().AddText("问题描述")

	row = table.AddRow()
	cell = row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
	cell.AddParagraph().AddRun().AddText("")
	row.AddCell().AddParagraph().AddRun().AddText("配置")
	row.AddCell().AddParagraph().AddRun().AddText("此攻击风险的根源来源于系统层、应用层或网络层等")

	row = table.AddRow()
	cell = row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
	cell.AddParagraph().AddRun().AddText("")
	row.AddCell().AddParagraph().AddRun().AddText("设计")
	row.AddCell().AddParagraph().AddRun().AddText("此攻击风险的根源来源于产品设计或者逻辑错误等")

	row = table.AddRow()
	cell = row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
	cell.AddParagraph().AddRun().AddText("")
	row.AddCell().AddParagraph().AddRun().AddText("代码")
	row.AddCell().AddParagraph().AddRun().AddText("此攻击风险的根源来源于错误的开发代码本身")

	row = table.AddRow()
	cell = row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
	cell.AddParagraph().AddRun().AddText("")
	row.AddCell().AddParagraph().AddRun().AddText("其他")
	row.AddCell().AddParagraph().AddRun().AddText("此攻击风险不是产品漏洞本身造成，由于产品第三方所造成的或未知的原因")

	row = table.AddRow()
	cell = row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeRestart)
	cell.Properties().SetVerticalAlignment(wml.ST_VerticalJcCenter)
	cell.AddParagraph().AddRun().AddText("漏洞危害")
	row.AddCell().AddParagraph().AddRun().AddText("分类")
	row.AddCell().AddParagraph().AddRun().AddText("描述")

	row = table.AddRow()
	cell = row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
	cell.AddParagraph().AddRun().AddText("")
	row.AddCell().AddParagraph().AddRun().AddText("高")
	row.AddCell().AddParagraph().AddRun().AddText("危及危害到系统安全")

	row = table.AddRow()
	cell = row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
	cell.AddParagraph().AddRun().AddText("")
	row.AddCell().AddParagraph().AddRun().AddText("中")
	row.AddCell().AddParagraph().AddRun().AddText("对系统、服务和设备的性能和状态有一定影响")

	row = table.AddRow()
	cell = row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
	cell.AddParagraph().AddRun().AddText("")
	row.AddCell().AddParagraph().AddRun().AddText("低")
	row.AddCell().AddParagraph().AddRun().AddText("对敏感信息的访问")

	row = table.AddRow()
	cell = row.AddCell()
	cell.Properties().SetVerticalMerge(wml.ST_MergeContinue)
	cell.AddParagraph().AddRun().AddText("")
	row.AddCell().AddParagraph().AddRun().AddText("提示")
	row.AddCell().AddParagraph().AddRun().AddText("不是安全风险，推荐使用符合安全开发规范方式实现")
	this.AddBreak()
	this.AddBreak()
}

// 添加附录B
func (this *WeekReportDocx) addNoteB() {
	this.AddAppendix().AddText("附录B：安全等级评级标准")
	this.AddMain().AddText("安全风险状况等级说明")

	//添加表格
	table := this.AddTable()
	table.Properties().SetWidthPercent(100)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Gray, 1*measurement.Point)

	row := table.AddRow()
	cell := row.AddCell()
	cell.Properties().SetColumnSpan(2)
	run := cell.AddParagraph().AddRun()
	run.AddText("安全风险状况说明")
	run.Properties().SetColor(color.White)
	cell.Properties().SetShading(wml.ST_ShdSolid, color.MyDeepBlue, color.Auto)

	row = table.AddRow()
	cell = row.AddCell()
	run = cell.AddParagraph().AddRun()
	run.AddText("1")
	run.Properties().SetColor(color.Green)
	cell.Properties().SetWidthPercent(15)
	cell = row.AddCell()
	run = cell.AddParagraph().AddRun()
	run.AddText("良好状态")
	run.Properties().SetColor(color.Green)
	cell.AddParagraph().AddRun().AddText("信息系统处于良好运行状态，没有发现或只存在零星的低风险安全问题，此时只要保持现有安全策略就满足了本系统的安全等级要求。")

	row = table.AddRow()
	run = row.AddCell().AddParagraph().AddRun()
	run.AddText("2")
	run.Properties().SetColor(color.DarkBlue)
	cell = row.AddCell()
	run = cell.AddParagraph().AddRun()
	run.AddText("预警状态")
	run.Properties().SetColor(color.DarkBlue)
	cell.AddParagraph().AddRun().AddText("信息系统中存在一些漏洞或安全隐患，此时需根据评估中发现的网络、主机、应用和管理等方面的问题对进行有针对性的加固或改进。")

	row = table.AddRow()
	run = row.AddCell().AddParagraph().AddRun()
	run.AddText("3")
	run.Properties().SetColor(color.OrangeRed)
	cell = row.AddCell()
	run = cell.AddParagraph().AddRun()
	run.AddText("严重状态")
	run.Properties().SetColor(color.OrangeRed)
	cell.AddParagraph().AddRun().AddText("信息系统中发现存在严重漏洞或可能严重威胁到系统正常运行的安全问题，此时需要立刻采取措施，例如安装补丁或重新部署安全系统进行防护等等。")

	row = table.AddRow()
	run = row.AddCell().AddParagraph().AddRun()
	run.AddText("4")
	run.Properties().SetColor(color.Red)
	cell = row.AddCell()
	run = cell.AddParagraph().AddRun()
	run.AddText("紧急状态")
	run.Properties().SetColor(color.Red)
	cell.AddParagraph().AddRun().AddText("信息系统面临严峻的网络安全态势，对组织的重大经济利益或政治利益可能造成严重损害。此时需要与其他安全部门通力协作采取紧急防御措施。")
	this.AddBreak()
	this.AddBreak()
}

// 添加附录C
func (this *WeekReportDocx) addNoteC() {
	this.AddAppendix().AddText("附录C：修复建议参考标准")

	//添加表格
	{
		table := this.AddTable()
		table.Properties().SetWidthPercent(100)
		borders := table.Properties().Borders()
		borders.SetAll(wml.ST_BorderSingle, color.MyDeepBlue, 1*measurement.Point)

		row := table.AddRow()
		cell := row.AddCell()
		run := cell.AddParagraph().AddRun()
		run.AddText("重要程度")
		run.Properties().SetColor(color.White)
		cell.Properties().SetShading(wml.ST_ShdSolid, color.MyDeepBlue, color.Auto)
		cell = row.AddCell()
		run = cell.AddParagraph().AddRun()
		run.AddText("说明")
		run.Properties().SetColor(color.White)
		cell.Properties().SetShading(wml.ST_ShdSolid, color.MyDeepBlue, color.Auto)

		row = table.AddRow()
		cell = row.AddCell()
		run = cell.AddParagraph().AddRun()
		run.AddText("必要")
		run.Properties().SetBold(true)
		cell.Properties().SetWidthPercent(15)
		cell.Properties().SetShading(wml.ST_ShdSolid, color.MyLightBlue, color.Auto)
		row.AddCell().AddParagraph().AddRun().AddText("在没有其他解决方案的情况下，优先考虑此项")

		row = table.AddRow()
		cell = row.AddCell()
		run = cell.AddParagraph().AddRun()
		run.AddText("推荐")
		run.Properties().SetBold(true)
		cell.Properties().SetWidthPercent(15)
		row.AddCell().AddParagraph().AddRun().AddText("在其他方案无法满足的情况下，或者风险无法完全避免的情况下，可以考虑此项")

		row = table.AddRow()
		cell = row.AddCell()
		run = cell.AddParagraph().AddRun()
		run.AddText("可选")
		run.Properties().SetBold(true)
		cell.Properties().SetWidthPercent(15)
		cell.Properties().SetShading(wml.ST_ShdSolid, color.MyLightBlue, color.Auto)
		row.AddCell().AddParagraph().AddRun().AddText("在资金和时间都有尚余力的情况下，可以考虑此项")
	}

	this.AddBreak()

	{
		//添加表格
		table := this.AddTable()
		table.Properties().SetWidthPercent(100)
		borders := table.Properties().Borders()
		borders.SetAll(wml.ST_BorderSingle, color.MyDeepBlue, 1*measurement.Point)

		row := table.AddRow()
		cell := row.AddCell()
		run := cell.AddParagraph().AddRun()
		run.AddText("修复成本")
		run.Properties().SetColor(color.White)
		cell.Properties().SetShading(wml.ST_ShdSolid, color.MyDeepBlue, color.Auto)
		cell = row.AddCell()
		run = cell.AddParagraph().AddRun()
		run.AddText("说明")
		run.Properties().SetColor(color.White)
		cell.Properties().SetShading(wml.ST_ShdSolid, color.MyDeepBlue, color.Auto)

		row = table.AddRow()
		cell = row.AddCell()
		run = cell.AddParagraph().AddRun()
		run.AddText("极高")
		run.Properties().SetBold(true)
		cell.Properties().SetWidthPercent(15)
		cell.Properties().SetShading(wml.ST_ShdSolid, color.MyLightBlue, color.Auto)
		row.AddCell().AddParagraph().AddRun().AddText("涉及多个供应商，现存案例少，生产工艺或系统架构需要重新设计，反复评审，消耗大量资源")

		row = table.AddRow()
		cell = row.AddCell()
		run = cell.AddParagraph().AddRun()
		run.AddText("高")
		run.Properties().SetBold(true)
		cell.Properties().SetWidthPercent(15)
		row.AddCell().AddParagraph().AddRun().AddText("涉及多个供应商，现存案例多，部分设计需要修改，反复评审，消耗大量资源")

		row = table.AddRow()
		cell = row.AddCell()
		run = cell.AddParagraph().AddRun()
		run.AddText("中")
		run.Properties().SetBold(true)
		cell.Properties().SetWidthPercent(15)
		cell.Properties().SetShading(wml.ST_ShdSolid, color.MyLightBlue, color.Auto)
		row.AddCell().AddParagraph().AddRun().AddText("涉及1个供应商，需要经过评审，增加或修改部分设计，需要多人完成")

		row = table.AddRow()
		cell = row.AddCell()
		run = cell.AddParagraph().AddRun()
		run.AddText("低")
		run.Properties().SetBold(true)
		cell.Properties().SetWidthPercent(15)
		row.AddCell().AddParagraph().AddRun().AddText("需要经过评审，讨论需求，进行修复，能在一周左右解决")

		row = table.AddRow()
		cell = row.AddCell()
		run = cell.AddParagraph().AddRun()
		run.AddText("极地")
		run.Properties().SetBold(true)
		cell.Properties().SetWidthPercent(15)
		cell.Properties().SetShading(wml.ST_ShdSolid, color.MyLightBlue, color.Auto)
		row.AddCell().AddParagraph().AddRun().AddText("无需评审，能在数日内解决")
	}

	this.AddBreak()
	this.AddBreak()
}

// 添加附录D
func (this *WeekReportDocx) addNoteD() {
	//添加附录D
	this.AddAppendix().AddText("附录D：测试依据")

	table := this.AddTable()
	table.Properties().SetWidthPercent(100)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Gray, 1*measurement.Point)
	this.AddDefaultRow(table, "序号", "标准名称")

	rowList := []string{
		"《ISO26262汽车电子电气系统的功能安全标准》",
		"《T/ITS 0054-2016基于公众电信网的汽车网关检测方法》",
		"《T/ITS 0057-2016车辆安全通信系统 设备性能测试规程》",
		"《GBT 20984-2007信息安全技术 信息安全风险评估规范》",
		"《ISO/IEC TR 13335信息安全管理指南》",
		"《Owasp Testing Guide v4测试指南》",
		"《Owsap ASVS 3.0 WEB应用安全评估标准》",
		"《移动互联网应用软件安全评估大纲》",
		"《中国国家信息安全漏洞库CNNVD》",
		"《OWASP Mobile Top 10 2017》",
		"《OWASP 安全检测指南》",
		"《GBT 36627-2018 信息安全技术 网络安全等级保护测试评估技术指南》",
		"《ISO-11898 道路车辆 控制器局域网络（CAN）》",
		"《ISO-14229 统一诊断服务》",
		"《ISO 26262 道路车辆功能安全》",
		"《GB 26149-2017 乘用车轮胎气压监测系统的性能要求和试验方法》",
		"《GB/T 12572-2008 无线电发射设备参数通用要求和测量方法》",
		"《GB/T 18314-2009 全球定位系统（GPS）测量规范》",
		"《GB/T 20008-2005 信息安全技术 操作系统安全评估准则》",
		"《GB/T 20271 信息安全技术 信息系统通用 安全技术要求》",
		"《GB/T 20277-2006 信息安全技术 网络和终端设备隔离部件测试评价方法》",
		"《GB/T 22186-2008 信息安全技术 具有中央处理器的集成电路（IC）卡芯片安 全技术要求(评估保证级 4 增强级)》",
		"《GB/T 25068.3-2010 信息技术 安全技术 IT 网络安全 第 3 部分：使用安 全网关的网间通信安全保护》",
		"《GB/T 26256-2010 2.4GHz 频段无线电通信 设备的相互干扰限制与共存要 求及测试方法》",
		"《GB/T 30284-2013 信息安全技术 移动通信智能终端操作系统安全技术要求》",
		"《GB/T 30290.3-2013 卫星定位车辆信息服务系统 第 3 部分：信息安全规范》",
		"《GB/T 30976.1-2014 工业控制系统信息安全 第 1 部分：评估规范》",
		"《GB/T 30976.2-2014 工业控制系统信息安全 第 2 部分：验收规范》",
		"《GB/T 32415-2015 GSM/CDMA/WCDMA 数字蜂窝移动通信网塔顶放大器 技术指标和测试方法》",
		"《GB/T 32420-2015 无线局域网测试规范 汽车产品信息安全测试评价规范》",
		"《GB/T 34975-2017 信息安全技术 移动智能终端应用软件安全技术要求和测试评价方法》",
		"《GB/T 34976-2017 信息安全技术 移动智能终端操作系统安全技术要求和测 试评价方法》",
		"《GB/T 34977-2017 信息安全技术 移动智能终端数据存储安全技术要求与测 试评价方法》",
		"《GB/T 35291-2017 信息安全技术 智能密码钥匙应用接口规范》",
		"《IEEE 802.15 基于蓝牙的局域网标准》",
		"《IEEE 802.11 无线局域网标准》",
		"《YD/T 2585-2016 互联网数据中心安全防护检测要求》",
		"《YD/T 2408-2013 移动智能终端安全能力测试方法》",
		"《YD/T 883-1996 900MHz TDMA 数字蜂窝移动通信网 基站无线设备技术指标及测试方法》",
	}

	for index, title := range rowList {
		this.AddDefaultRow(table, strconv.Itoa(index+1), title)
	}

}

type ReportParse struct {
	*document.Document
}

func (this *ReportParse) AddDefaultRow(table document.Table, left, right string) document.Row {
	row := table.AddRow()
	cell := row.AddCell()
	cell.AddParagraph().AddRun().AddText(left)
	cell.Properties().SetWidthPercent(15)
	row.AddCell().AddParagraph().AddRun().AddText(right)

	return row
}

func (this *ReportParse) AddTitle() document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-title")
	run := para.AddRun()
	run.Properties().SetBold(true)
	run.Properties().SetSize(20)
	return run
}

func (this *ReportParse) AddTitleWihtStyle() document.Run {
	para := this.AddParagraph()
	para.Properties().SetAlignment(wml.ST_JcCenter)
	run := para.AddRun()
	run.Properties().SetBold(true)
	run.Properties().SetSize(20)
	return run
}

func (this *ReportParse) AddTitle1() document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-title1")
	run := para.AddRun()
	run.Properties().SetBold(true)
	return run
}

func (this *ReportParse) AddTitle5() document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-title5")
	run := para.AddRun()
	run.Properties().SetBold(true)
	return run
}

func (this *ReportParse) AddMiddle(indent int) document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-title1")
	//para.Properties().SetFirstLineIndent(measurement.Distance(indent) * measurement.ChChar)
	para.Properties().AddTabStop(2.5*measurement.Inch, wml.ST_TabJcCenter, wml.ST_TabTlcNone)
	run := para.AddRun()
	run.AddTab()
	run.Properties().SetBold(true)
	run.Properties().SetSize(16)
	return run
}

func (this *ReportParse) AddLeft() document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-main")
	para.Properties().AddTabStop(4*measurement.Inch, wml.ST_TabJcLeft, wml.ST_TabTlcNone)
	run := para.AddRun()
	run.AddTab()
	return run
}

func (this *ReportParse) AddRight() document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-main")
	para.Properties().AddTabStop(0.1*measurement.Inch, wml.ST_TabJcRight, wml.ST_TabTlcNone)
	//para.Properties().AddTabStop(4.5*measurement.Inch, wml.ST_TabJcLeft, wml.ST_TabTlcNone)
	run := para.AddRun()
	run.AddTab()
	return run
}

func (this *ReportParse) AddBreak() document.Run {
	para := this.AddParagraph()
	run := para.AddRun()
	return run
}

func (this *ReportParse) AddMultiBreak(num int) {
	for i := 0; i < num; i++ {
		this.AddBreak()
	}
}

func (this *ReportParse) AddHeading1() document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-heading1")
	return para.AddRun()
}

func (this *ReportParse) AddHeading2() document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-heading2")
	return para.AddRun()
}

func (this *ReportParse) AddHeading3() document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-heading3")
	return para.AddRun()
}

func (this *ReportParse) AddHeading4() document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-heading4")
	return para.AddRun()
}

func (this *ReportParse) AddHeading5() document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-heading5")
	return para.AddRun()
}

func (this *ReportParse) AddAppendix() document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-appendix")
	return para.AddRun()
}

func (this *ReportParse) AddMain() document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-main")
	return para.AddRun()
}

func (this *ReportParse) AddMain2() document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-main2")
	return para.AddRun()
}

func (this *ReportParse) AddIndent() document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-main")
	para.Properties().SetFirstLineIndent(2 * measurement.ChChar)
	return para.AddRun()
}

func (this *ReportParse) AddMargin(indent int) document.Run {
	para := this.AddParagraph()
	para.SetStyle("R-main")
	para.Properties().SetFirstLineIndent(measurement.Distance(indent) * measurement.ChChar)
	return para.AddRun()
}

type HtmlElem struct {
	Text string
	Html string
}

func (this *ReportParse) addParagraph(segment *service.MdSegment) {
	para := this.AddParagraph()
	run := para.AddRun()
	for _, elem := range segment.Elements {
		switch elem.Type {
		case service.MD_ELEM_H1, service.MD_ELEM_H2, service.MD_ELEM_H3, service.MD_ELEM_H4, service.MD_ELEM_H5, service.MD_ELEM_H6:
			para.SetStyle("R-heading5")
			this.AddText(run, elem.Content)
		case service.MD_ELEM_IMG:
			this.InsertImage(run, elem.Content)
		case service.MD_ELEM_ORDER_LIST:
			para.SetStyle("R-main2")
			//para.SetStyle("R-order")
			//判断一下 列表里是内容 还是图片
			if elemImage, isImage := service.IsImage(elem.Content); isImage {
				this.InsertImage(run, elemImage)
			} else {
				this.AddText(run, elem.Content)
			}
		case service.MD_ELEM_UNORDER_LIST:
			//para.SetStyle("R-main2")
			para.SetStyle("R-unorder")
			//判断一下 列表里是内容 还是图片
			if elemImage, isImage := service.IsImage(elem.Content); isImage {
				this.InsertImage(run, elemImage)
			} else {
				this.AddText(run, elem.Content)
			}
		default:
			this.AddText(run, elem.Content)
		}
	}
}
func (this *ReportParse) AddText(docRun document.Run, text string) {
	docRun.AddText(text)
}

func (this *ReportParse) AddDotText(docRun document.Run, text string) {
	docRun.Properties().SetSize(28)
	docRun.AddText("· ")
	//docRun.Properties().SetSize(12)
	docRun.AddText(text)
}

func (this *ReportParse) InsertImage(docRun document.Run, url string) {
	imgContent, err := this.DownloadImageFromMongoDB(url)
	if err != nil {
		return
	}
	if newImg, err := unioffice_common.ImageFromBytes(imgContent); err == nil {
		if img1ref, err := this.AddImage(newImg); err == nil {
			if inl, err := docRun.AddDrawingInline(img1ref); err == nil {
				ratio := measurement.Distance(newImg.Size.Y) / measurement.Distance(newImg.Size.X)
				inl.SetSize(5.5*measurement.Inch, 5.5*measurement.Inch*ratio)
			}
		}
	}
}

func (this *ReportParse) ParseMarkdown(mdContent string) {
	markdownSegments := service.ParseMarkdown(mdContent)
	for _, segment := range markdownSegments {
		this.addParagraph(segment)
	}
}

func (this *ReportParse) DownloadImageFromMongoDB(imgUrl string) ([]byte, error) {
	urlReg := regexp.MustCompile(`^[^:]*file_id=([^\s]*)`)
	if params := urlReg.FindStringSubmatch(imgUrl); len(params) == 2 {
		if fi, err := mongo.GridFSOpenId(common.MC_File, bson.ObjectIdHex(params[1])); err == nil {
			return ioutil.ReadAll(fi)
		} else {
			return nil, err
		}
	} else {
		return nil, errors.New(fmt.Sprintf("Image url parse error:%s", imgUrl))
	}
}

func (this *ReportParse) ParseTable(mdContent string) string {
	re := regexp.MustCompile(`(\|[^\n]+\|)`)
	matched := re.FindAllStringSubmatch(mdContent, -1)
	mdContent = re.ReplaceAllString(mdContent, "pppppp\n")
	re2 := regexp.MustCompile(`(pppppp\n)`)
	mdContent = re2.ReplaceAllString(mdContent, "")
	//添加表格
	table := this.AddTable()
	table.Properties().SetWidthPercent(100)
	borders := table.Properties().Borders()
	borders.SetAll(wml.ST_BorderSingle, color.Gray, 1*measurement.Point)

	for itemIndex, match := range matched {
		if itemIndex == 1 {
			continue
		}

		matchSlice := strings.Split(match[0], "|")
		row := table.AddRow()
		for colIndex, item := range matchSlice {
			if colIndex == 0 || colIndex == len(matchSlice)-1 {
				continue
			}
			row.AddCell().AddParagraph().AddRun().AddText(item)
		}
	}
	return mdContent
}

func (this *ReportParse) ParseBreak(mdContent string) {
	textSlice := strings.Split(mdContent, "\n")
	for _, text := range textSlice {
		para := this.AddParagraph()
		run := para.AddRun()
		run.AddText(text)
	}
}

func GenerateReportName(projectId, reportType string, todayReportCount int) string {
	params := qmap.QM{
		"e__id": bson.ObjectIdHex(projectId),
	}
	projectName := ""
	if project, err := mongo.NewMgoSessionWithCond(common.MC_PROJECT, params).GetOne(); err == nil {
		projectName = project.String("name")
	}
	switch reportType {
	case "week":
		return fmt.Sprintf("%s-周报告-%s-%d.docx", projectName, time.Now().Format("20060102"), todayReportCount)
	case common.REPORT_TEST:
		return fmt.Sprintf("%s-初测报告-%s-%d.docx", projectName, time.Now().Format("20060102"), todayReportCount)
	case common.REPORT_RETEST:
		return fmt.Sprintf("%s-复测报告-%s-%d.docx", projectName, time.Now().Format("20060102"), todayReportCount)
	case common.REPORT_ASSETTEST:
		return fmt.Sprintf("%s-资产初测报告-%s-%d.docx", projectName, time.Now().Format("20060102"), todayReportCount)
	case common.REPORT_ASSETRETEST:
		return fmt.Sprintf("%s-资产复测报告-%s-%d.docx", projectName, time.Now().Format("20060102"), todayReportCount)
	default:
		panic("unknown report type!")
	}
}

func (this *WeekReportDocx) addHeader() {
	doc := this.Document
	img, err := unioffice_common.ImageFromFile(reportTemplateConfig.ReportHearder)
	if err != nil {
		log.Fatalf("unable to create image: %s", err)
	}

	hdr := doc.AddHeader()
	// We need to add a reference of the image to the header instead of to the
	// document
	iref, err := hdr.AddImage(img)
	if err != nil {
		log.Fatalf("unable to to add image to document: %s", err)
	}

	para := hdr.AddParagraph()
	imgInl, _ := para.AddRun().AddDrawingInline(iref)
	imgInl.SetSize(10*measurement.Inch, 0.6*measurement.Inch)
	doc.BodySection().SetHeader(hdr, wml.ST_HdrFtrDefault)
}

func (this *WeekReportDocx) addFooter() {
	doc := this.Document
	img, err := unioffice_common.ImageFromFile(reportTemplateConfig.ReportFooter)
	if err != nil {
		log.Fatalf("unable to create image: %s", err)
	}

	ftr := doc.AddFooter()
	para := ftr.AddParagraph()

	para.Properties().AddTabStop(6*measurement.Inch, wml.ST_TabJcRight, wml.ST_TabTlcNone)
	run := para.AddRun()
	run.Properties().SetBold(true)
	run.Properties().SetSize(10)
	run.AddText("360 SKY-GO智能网联汽车安全实验室")
	run.AddTab()

	iref, err := ftr.AddImage(img)
	if err != nil {
		log.Fatalf("unable to to add image to document: %s", err)
	}
	para = ftr.AddParagraph()
	para.Properties().AddTabStop(0.1*measurement.Inch, wml.ST_TabJcLeft, wml.ST_TabTlcNone)
	imgInl, _ := para.AddRun().AddDrawingInline(iref)
	imgInl.SetSize(6*measurement.Inch, 0.3*measurement.Inch)

	para = ftr.AddParagraph()
	para.Properties().AddTabStop(2.5*measurement.Inch, wml.ST_TabJcLeft, wml.ST_TabTlcNone)
	run = para.AddRun()
	run.AddTab()
	run.AddTab()
	run.Properties().SetBold(true)
	run.Properties().SetSize(10)
	run.AddText("Pg ")
	run.AddField(document.FieldCurrentPage)
	run.AddText(" of ")
	run.AddField(document.FieldNumberOfPages)
	doc.BodySection().SetFooter(ftr, wml.ST_HdrFtrDefault)

}

// 添加目录
func (this *WeekReportDocx) addToc() {
	doc := this.Document
	//doc.Settings.SetUpdateFieldsOnOpen(true)
	// Add a TOC
	doc.AddParagraph().AddRun().AddField(document.FieldTOC)
	// followed by a page break
	//doc.AddParagraph().Properties().AddSection(wml.ST_SectionMarkNextPage)
	//nd := doc.Numbering.AddDefinition()
	//for i := 0; i < 9; i++ {
	//	lvl := nd.AddLevel()
	//	lvl.SetFormat(wml.ST_NumberFormatNone)
	//	lvl.SetAlignment(wml.ST_JcLeft)
	//	//if i%2 == 0 {
	//	//	lvl.SetFormat(wml.ST_NumberFormatBullet)
	//	//	lvl.RunProperties().SetFontFamily("Symbol")
	//	//	lvl.SetText("")
	//	//}
	//	lvl.Properties().SetLeftIndent(0.5 * measurement.Distance(i) * measurement.Inch)
	//}

	doc.AddParagraph().AddRun().AddPageBreak()
}

func (this *WeekReportDocx) addDef() {
	doc := this.Document
	nd := doc.Numbering.AddDefinition()
	for i := 0; i < 9; i++ {
		lvl := nd.AddLevel()
		lvl.SetFormat(wml.ST_NumberFormatNone)
		lvl.SetAlignment(wml.ST_JcLeft)
		if i%2 == 0 {
			lvl.SetFormat(wml.ST_NumberFormatBullet)
			lvl.RunProperties().SetFontFamily("Symbol")
			lvl.SetText("")
		}
		lvl.Properties().SetLeftIndent(0.5 * measurement.Distance(i) * measurement.Inch)
	}
}
