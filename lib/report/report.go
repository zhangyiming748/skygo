package report

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/golang/freetype/truetype"
	"github.com/wcharczuk/go-chart/v2"
	"github.com/wcharczuk/go-chart/v2/drawing"
	"skygo_detection/guardian/src/net/qmap"

	"skygo_detection/lib/unioffice/common"
	"skygo_detection/lib/unioffice/measurement"
)

var (
	Blue   = drawing.Color{R: 0, G: 0, B: 128, A: 255}
	Red    = drawing.Color{R: 163, G: 0, B: 20, A: 255}
	Orange = drawing.Color{R: 251, G: 212, B: 161, A: 255}
	Yellow = drawing.Color{R: 255, G: 204, B: 0, A: 255}
	Gray   = drawing.Color{R: 187, G: 187, B: 187, A: 255}
)

type ReportCommont struct {
	WeekReportDocx
}

// 从文件中添加图片
func (this *ReportCommont) AddImageFromFile(path string) {
	if path == "" {
		return
	}
	doc := this.Document
	img, err := common.ImageFromFile(path)
	if err != nil {
		log.Fatalf("unable to create image: %s", err)
	}
	imgref, err := doc.AddImage(img)
	if err != nil {
		log.Fatalf("unable to add image to document: %s", err)
	}
	para := doc.AddParagraph()
	run := para.AddRun()
	anchored, err := run.AddDrawingInline(imgref)
	if err != nil {
		log.Fatalf("unable to add anchored image: %s", err)
	}
	anchored.SetSize(2*measurement.Inch, 2*measurement.Inch)
}

// 数据中图片
func (this *ReportCommont) AddImageFromBytes(data []byte) {
	doc := this.Document
	img, err := common.ImageFromBytes(data)
	if err != nil {
		log.Fatalf("unable to create image: %s", err)
	}
	imgref, err := doc.AddImage(img)
	if err != nil {
		log.Fatalf("unable to add image to document: %s", err)
	}
	para := doc.AddParagraph()
	run := para.AddRun()
	anchored, err := run.AddDrawingInline(imgref)
	if err != nil {
		log.Fatalf("unable to add anchored image: %s", err)
	}
	anchored.SetSize(2*measurement.Inch, 2*measurement.Inch)
}

// 生成饼图
func (this *ReportCommont) NewPieImage(params qmap.QM, name string, title string) string {
	values := []chart.Value{}
	for key, value := range params {
		tmpValue := value.(qmap.QM)
		t := tmpValue.Int("value")
		if t == 0 {
			continue
		}
		switch key {
		case "red":
			tmpStyle := chart.Style{FillColor: Red}
			tmpChart := chart.Value{Value: float64(t), Label: tmpValue.String("label"), Style: tmpStyle}
			values = append(values, tmpChart)
		case "yellow":
			tmpStyle := chart.Style{FillColor: Yellow}
			tmpChart := chart.Value{Value: float64(t), Label: tmpValue.String("label"), Style: tmpStyle}
			values = append(values, tmpChart)
		case "blue":
			tmpStyle := chart.Style{FillColor: Blue}
			tmpChart := chart.Value{Value: float64(t), Label: tmpValue.String("label"), Style: tmpStyle}
			values = append(values, tmpChart)
		case "orange":
			tmpStyle := chart.Style{FillColor: Orange}
			tmpChart := chart.Value{Value: float64(t), Label: tmpValue.String("label"), Style: tmpStyle}
			values = append(values, tmpChart)
		case "gray":
			tmpStyle := chart.Style{FillColor: Gray}
			tmpChart := chart.Value{Value: float64(t), Label: tmpValue.String("label"), Style: tmpStyle}
			values = append(values, tmpChart)
		}
	}
	f1 := this.getZWFont(reportTemplateConfig.FontPath)
	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: values,
		Font:   f1,
		Title:  title,
	}
	f, _ := os.Create(fmt.Sprintf(reportTemplateConfig.OutputImage, name))
	defer f.Close()
	pie.Render(chart.PNG, f)
	return f.Name()
}

func (this *ReportCommont) NewPieImage1(params []qmap.QM, name string, title string) string {
	values := []chart.Value{}
	for _, param := range params {
		value := param.Float64("value")
		label := param.String("label")
		values = append(values, chart.Value{Value: value, Label: label})
	}
	f1 := this.getZWFont(reportTemplateConfig.FontPath)
	pie := chart.PieChart{
		Width:  512,
		Height: 512,
		Values: values,
		Font:   f1,
		Title:  title,
	}
	f, _ := os.Create(fmt.Sprintf(reportTemplateConfig.OutputImage, name))
	defer f.Close()
	pie.Render(chart.PNG, f)
	return f.Name()
}

// 生成环状图
func (this *ReportCommont) NewDonutImage(params qmap.QM, name string, title string) string {
	values := []chart.Value{}
	for key, value := range params {
		tmpValue := value.(qmap.QM)
		t := tmpValue.Int("value")
		if t == 0 {
			continue
		}
		switch key {
		case "red":
			tmpStyle := chart.Style{FillColor: Red}
			tmpChart := chart.Value{Value: float64(t), Label: tmpValue.String("label"), Style: tmpStyle}
			values = append(values, tmpChart)
		case "yellow":
			tmpStyle := chart.Style{FillColor: Yellow}
			tmpChart := chart.Value{Value: float64(t), Label: tmpValue.String("label"), Style: tmpStyle}
			values = append(values, tmpChart)
		case "blue":
			tmpStyle := chart.Style{FillColor: Blue}
			tmpChart := chart.Value{Value: float64(t), Label: tmpValue.String("label"), Style: tmpStyle}
			values = append(values, tmpChart)
		case "orange":
			tmpStyle := chart.Style{FillColor: Orange}
			tmpChart := chart.Value{Value: float64(t), Label: tmpValue.String("label"), Style: tmpStyle}
			values = append(values, tmpChart)
		case "gray":
			tmpStyle := chart.Style{FillColor: Gray}
			tmpChart := chart.Value{Value: float64(t), Label: tmpValue.String("label"), Style: tmpStyle}
			values = append(values, tmpChart)
		}
	}
	f1 := this.getZWFont(reportTemplateConfig.FontPath)
	pie := chart.DonutChart{
		Width:  512,
		Height: 512,
		Values: values,
		Font:   f1,
		Title:  title,
	}

	f, _ := os.Create(fmt.Sprintf(reportTemplateConfig.OutputImage, name))
	defer f.Close()
	pie.Render(chart.PNG, f)
	return f.Name()
}

func (this *ReportCommont) NewDonutImage1(params []qmap.QM, name string, title string) string {
	values := []chart.Value{}
	for _, param := range params {
		value := param.Float64("value")
		label := param.String("label")
		values = append(values, chart.Value{Value: value, Label: label})
	}
	f1 := this.getZWFont(reportTemplateConfig.FontPath)
	pie := chart.DonutChart{
		Width:  512,
		Height: 512,
		Values: values,
		Font:   f1,
		Title:  title,
	}
	f, _ := os.Create(fmt.Sprintf(reportTemplateConfig.OutputImage, name))
	defer f.Close()
	pie.Render(chart.PNG, f)
	return f.Name()
}

// 获取中文编码
func (this *ReportCommont) getZWFont(path string) *truetype.Font {
	fontFile := path
	// 读字体数据
	fontBytes, err := ioutil.ReadFile(fontFile)
	if err != nil {
		log.Println(err)
		return nil
	}
	font, err := truetype.Parse(fontBytes)
	if err != nil {
		log.Println(err)
		return nil
	}
	return font
}
