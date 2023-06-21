package service

import (
	"bytes"
	"fmt"

	"github.com/360EntSecGroup-Skylar/excelize"
)

// 1. 创建一个写excel对象
type ExcelObj struct {
	ExcelFile    *excelize.File
	SheetIndex   int
	Sheet        string
	ContentTitle []string
	ContentKey   []string
	Content      []interface{}
	RowIndex     int
	Row          string
}

// 2. 创建excel
func (this *ExcelObj) NewExcel() *ExcelObj {
	this.ExcelFile = excelize.NewFile()
	return this
}

// 3. 创建sheet
func (this *ExcelObj) NewSheet() int {
	this.SheetIndex = this.SheetIndex + 1
	this.Sheet = fmt.Sprintf("Sheet%v", this.SheetIndex)
	this.InitRow()
	return this.ExcelFile.NewSheet(this.Sheet)
}

// 新增row
func (this *ExcelObj) AddRow() {
	this.RowIndex = this.RowIndex + 1
	this.Row = fmt.Sprintf("A%v", this.RowIndex)
}

func (this *ExcelObj) InitRow() {
	this.RowIndex = 1
	this.Row = fmt.Sprintf("A%v", this.RowIndex)
}

// 写title
func (this *ExcelObj) WriteTitle() {
	this.ExcelFile.SetSheetRow(this.Sheet, this.Row, &this.ContentTitle)
}

// 4. 写excel内容
func (this *ExcelObj) WriteContent() {
	this.AddRow()
	// this.ExcelFile.SetDefaultFont("微软雅黑")
	this.ExcelFile.SetSheetRow(this.Sheet, this.Row, &this.Content)
}

// 5. 返回bytes
func (this *ExcelObj) Output() *bytes.Buffer {
	buffer := new(bytes.Buffer)
	this.ExcelFile.Write(buffer)
	return buffer
}
