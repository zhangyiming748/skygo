package unioffice_test

import (
	"skygo_detection/lib/unioffice/document"
	"skygo_detection/lib/unioffice/spreadsheet"
)

func Example_document() {
	// see the skygo_detection/lib/unioffice/document documentation or _examples/document
	// for more examples
	doc := document.New()
	doc.AddParagraph().AddRun().AddText("Hello World!")
	doc.SaveToFile("document.docx")
}

func Example_spreadsheeet() {
	// see the skygo_detection/lib/unioffice/spreadsheet documentation or _examples/spreadsheet
	// for more examples
	ss := spreadsheet.New()
	sheet := ss.AddSheet()
	sheet.AddRow().AddCell().SetString("Hello")
	sheet.Cell("B1").SetString("World!")
	ss.SaveToFile("workbook.xlsx")
}
