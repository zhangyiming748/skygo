// Copyright 2017 FoxyUtils ehf. All rights reserved.
//
// DO NOT EDIT: generated by gooxml ECMA-376 generator
//
// Use of this source code is governed by the terms of the Affero GNU General
// Public License version 3.0 as published by the Free Software Foundation and
// appearing in the file LICENSE included in the packaging of this file. A
// commercial license can be purchased via https://unidoc.io website.

package sml

import (
	"encoding/xml"
	"fmt"
	"strconv"

	"skygo_detection/lib/unioffice/schema/soo/ofc/sharedTypes"
)

type CT_PageSetup struct {
	// Paper Size
	PaperSizeAttr *uint32
	// Paper Height
	PaperHeightAttr *string
	// Paper Width
	PaperWidthAttr *string
	// Print Scale
	ScaleAttr *uint32
	// First Page Number
	FirstPageNumberAttr *uint32
	// Fit To Width
	FitToWidthAttr *uint32
	// Fit To Height
	FitToHeightAttr *uint32
	// Page Order
	PageOrderAttr ST_PageOrder
	// Orientation
	OrientationAttr ST_Orientation
	// Use Printer Defaults
	UsePrinterDefaultsAttr *bool
	// Black And White
	BlackAndWhiteAttr *bool
	// Draft
	DraftAttr *bool
	// Print Cell Comments
	CellCommentsAttr ST_CellComments
	// Use First Page Number
	UseFirstPageNumberAttr *bool
	// Print Error Handling
	ErrorsAttr ST_PrintError
	// Horizontal DPI
	HorizontalDpiAttr *uint32
	// Vertical DPI
	VerticalDpiAttr *uint32
	// Number Of Copies
	CopiesAttr *uint32
	IdAttr     *string
}

func NewCT_PageSetup() *CT_PageSetup {
	ret := &CT_PageSetup{}
	return ret
}

func (m *CT_PageSetup) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if m.PaperSizeAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "paperSize"},
			Value: fmt.Sprintf("%v", *m.PaperSizeAttr)})
	}
	if m.PaperHeightAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "paperHeight"},
			Value: fmt.Sprintf("%v", *m.PaperHeightAttr)})
	}
	if m.PaperWidthAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "paperWidth"},
			Value: fmt.Sprintf("%v", *m.PaperWidthAttr)})
	}
	if m.ScaleAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "scale"},
			Value: fmt.Sprintf("%v", *m.ScaleAttr)})
	}
	if m.FirstPageNumberAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "firstPageNumber"},
			Value: fmt.Sprintf("%v", *m.FirstPageNumberAttr)})
	}
	if m.FitToWidthAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "fitToWidth"},
			Value: fmt.Sprintf("%v", *m.FitToWidthAttr)})
	}
	if m.FitToHeightAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "fitToHeight"},
			Value: fmt.Sprintf("%v", *m.FitToHeightAttr)})
	}
	if m.PageOrderAttr != ST_PageOrderUnset {
		attr, err := m.PageOrderAttr.MarshalXMLAttr(xml.Name{Local: "pageOrder"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.OrientationAttr != ST_OrientationUnset {
		attr, err := m.OrientationAttr.MarshalXMLAttr(xml.Name{Local: "orientation"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.UsePrinterDefaultsAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "usePrinterDefaults"},
			Value: fmt.Sprintf("%d", b2i(*m.UsePrinterDefaultsAttr))})
	}
	if m.BlackAndWhiteAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "blackAndWhite"},
			Value: fmt.Sprintf("%d", b2i(*m.BlackAndWhiteAttr))})
	}
	if m.DraftAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "draft"},
			Value: fmt.Sprintf("%d", b2i(*m.DraftAttr))})
	}
	if m.CellCommentsAttr != ST_CellCommentsUnset {
		attr, err := m.CellCommentsAttr.MarshalXMLAttr(xml.Name{Local: "cellComments"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.UseFirstPageNumberAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "useFirstPageNumber"},
			Value: fmt.Sprintf("%d", b2i(*m.UseFirstPageNumberAttr))})
	}
	if m.ErrorsAttr != ST_PrintErrorUnset {
		attr, err := m.ErrorsAttr.MarshalXMLAttr(xml.Name{Local: "errors"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.HorizontalDpiAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "horizontalDpi"},
			Value: fmt.Sprintf("%v", *m.HorizontalDpiAttr)})
	}
	if m.VerticalDpiAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "verticalDpi"},
			Value: fmt.Sprintf("%v", *m.VerticalDpiAttr)})
	}
	if m.CopiesAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "copies"},
			Value: fmt.Sprintf("%v", *m.CopiesAttr)})
	}
	if m.IdAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "r:id"},
			Value: fmt.Sprintf("%v", *m.IdAttr)})
	}
	e.EncodeToken(start)
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (m *CT_PageSetup) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initialize to default
	for _, attr := range start.Attr {
		if attr.Name.Space == "http://schemas.openxmlformats.org/officeDocument/2006/relationships" && attr.Name.Local == "id" ||
			attr.Name.Space == "http://purl.oclc.org/ooxml/officeDocument/relationships" && attr.Name.Local == "id" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.IdAttr = &parsed
			continue
		}
		if attr.Name.Local == "paperSize" {
			parsed, err := strconv.ParseUint(attr.Value, 10, 32)
			if err != nil {
				return err
			}
			pt := uint32(parsed)
			m.PaperSizeAttr = &pt
			continue
		}
		if attr.Name.Local == "blackAndWhite" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.BlackAndWhiteAttr = &parsed
			continue
		}
		if attr.Name.Local == "draft" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.DraftAttr = &parsed
			continue
		}
		if attr.Name.Local == "scale" {
			parsed, err := strconv.ParseUint(attr.Value, 10, 32)
			if err != nil {
				return err
			}
			pt := uint32(parsed)
			m.ScaleAttr = &pt
			continue
		}
		if attr.Name.Local == "cellComments" {
			m.CellCommentsAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "fitToWidth" {
			parsed, err := strconv.ParseUint(attr.Value, 10, 32)
			if err != nil {
				return err
			}
			pt := uint32(parsed)
			m.FitToWidthAttr = &pt
			continue
		}
		if attr.Name.Local == "pageOrder" {
			m.PageOrderAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "orientation" {
			m.OrientationAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "paperHeight" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.PaperHeightAttr = &parsed
			continue
		}
		if attr.Name.Local == "paperWidth" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.PaperWidthAttr = &parsed
			continue
		}
		if attr.Name.Local == "firstPageNumber" {
			parsed, err := strconv.ParseUint(attr.Value, 10, 32)
			if err != nil {
				return err
			}
			pt := uint32(parsed)
			m.FirstPageNumberAttr = &pt
			continue
		}
		if attr.Name.Local == "fitToHeight" {
			parsed, err := strconv.ParseUint(attr.Value, 10, 32)
			if err != nil {
				return err
			}
			pt := uint32(parsed)
			m.FitToHeightAttr = &pt
			continue
		}
		if attr.Name.Local == "useFirstPageNumber" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.UseFirstPageNumberAttr = &parsed
			continue
		}
		if attr.Name.Local == "errors" {
			m.ErrorsAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "horizontalDpi" {
			parsed, err := strconv.ParseUint(attr.Value, 10, 32)
			if err != nil {
				return err
			}
			pt := uint32(parsed)
			m.HorizontalDpiAttr = &pt
			continue
		}
		if attr.Name.Local == "verticalDpi" {
			parsed, err := strconv.ParseUint(attr.Value, 10, 32)
			if err != nil {
				return err
			}
			pt := uint32(parsed)
			m.VerticalDpiAttr = &pt
			continue
		}
		if attr.Name.Local == "copies" {
			parsed, err := strconv.ParseUint(attr.Value, 10, 32)
			if err != nil {
				return err
			}
			pt := uint32(parsed)
			m.CopiesAttr = &pt
			continue
		}
		if attr.Name.Local == "usePrinterDefaults" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.UsePrinterDefaultsAttr = &parsed
			continue
		}
	}
	// skip any extensions we may find, but don't support
	for {
		tok, err := d.Token()
		if err != nil {
			return fmt.Errorf("parsing CT_PageSetup: %s", err)
		}
		if el, ok := tok.(xml.EndElement); ok && el.Name == start.Name {
			break
		}
	}
	return nil
}

// Validate validates the CT_PageSetup and its children
func (m *CT_PageSetup) Validate() error {
	return m.ValidateWithPath("CT_PageSetup")
}

// ValidateWithPath validates the CT_PageSetup and its children, prefixing error messages with path
func (m *CT_PageSetup) ValidateWithPath(path string) error {
	if m.PaperHeightAttr != nil {
		if !sharedTypes.ST_PositiveUniversalMeasurePatternRe.MatchString(*m.PaperHeightAttr) {
			return fmt.Errorf(`%s/m.PaperHeightAttr must match '%s' (have %v)`, path, sharedTypes.ST_PositiveUniversalMeasurePatternRe, *m.PaperHeightAttr)
		}
	}
	if m.PaperHeightAttr != nil {
		if !sharedTypes.ST_UniversalMeasurePatternRe.MatchString(*m.PaperHeightAttr) {
			return fmt.Errorf(`%s/m.PaperHeightAttr must match '%s' (have %v)`, path, sharedTypes.ST_UniversalMeasurePatternRe, *m.PaperHeightAttr)
		}
	}
	if m.PaperWidthAttr != nil {
		if !sharedTypes.ST_PositiveUniversalMeasurePatternRe.MatchString(*m.PaperWidthAttr) {
			return fmt.Errorf(`%s/m.PaperWidthAttr must match '%s' (have %v)`, path, sharedTypes.ST_PositiveUniversalMeasurePatternRe, *m.PaperWidthAttr)
		}
	}
	if m.PaperWidthAttr != nil {
		if !sharedTypes.ST_UniversalMeasurePatternRe.MatchString(*m.PaperWidthAttr) {
			return fmt.Errorf(`%s/m.PaperWidthAttr must match '%s' (have %v)`, path, sharedTypes.ST_UniversalMeasurePatternRe, *m.PaperWidthAttr)
		}
	}
	if err := m.PageOrderAttr.ValidateWithPath(path + "/PageOrderAttr"); err != nil {
		return err
	}
	if err := m.OrientationAttr.ValidateWithPath(path + "/OrientationAttr"); err != nil {
		return err
	}
	if err := m.CellCommentsAttr.ValidateWithPath(path + "/CellCommentsAttr"); err != nil {
		return err
	}
	if err := m.ErrorsAttr.ValidateWithPath(path + "/ErrorsAttr"); err != nil {
		return err
	}
	return nil
}