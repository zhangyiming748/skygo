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
)

type CT_MetadataType struct {
	// Metadata Type Name
	NameAttr string
	// Minimum Supported Version
	MinSupportedVersionAttr uint32
	// Metadata Ghost Row
	GhostRowAttr *bool
	// Metadata Ghost Column
	GhostColAttr *bool
	// Metadata Edit
	EditAttr *bool
	// Metadata Cell Value Delete
	DeleteAttr *bool
	// Metadata Copy
	CopyAttr *bool
	// Metadata Paste All
	PasteAllAttr *bool
	// Metadata Paste Formulas
	PasteFormulasAttr *bool
	// Metadata Paste Special Values
	PasteValuesAttr *bool
	// Metadata Paste Formats
	PasteFormatsAttr *bool
	// Metadata Paste Comments
	PasteCommentsAttr *bool
	// Metadata Paste Data Validation
	PasteDataValidationAttr *bool
	// Metadata Paste Borders
	PasteBordersAttr *bool
	// Metadata Paste Column Widths
	PasteColWidthsAttr *bool
	// Metadata Paste Number Formats
	PasteNumberFormatsAttr *bool
	// Metadata Merge
	MergeAttr *bool
	// Meatadata Split First
	SplitFirstAttr *bool
	// Metadata Split All
	SplitAllAttr *bool
	// Metadata Insert Delete
	RowColShiftAttr *bool
	// Metadata Clear All
	ClearAllAttr *bool
	// Metadata Clear Formats
	ClearFormatsAttr *bool
	// Metadata Clear Contents
	ClearContentsAttr *bool
	// Metadata Clear Comments
	ClearCommentsAttr *bool
	// Metadata Formula Assignment
	AssignAttr *bool
	// Metadata Coercion
	CoerceAttr *bool
	// Adjust Metadata
	AdjustAttr *bool
	// Cell Metadata
	CellMetaAttr *bool
}

func NewCT_MetadataType() *CT_MetadataType {
	ret := &CT_MetadataType{}
	return ret
}

func (m *CT_MetadataType) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "name"},
		Value: fmt.Sprintf("%v", m.NameAttr)})
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "minSupportedVersion"},
		Value: fmt.Sprintf("%v", m.MinSupportedVersionAttr)})
	if m.GhostRowAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "ghostRow"},
			Value: fmt.Sprintf("%d", b2i(*m.GhostRowAttr))})
	}
	if m.GhostColAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "ghostCol"},
			Value: fmt.Sprintf("%d", b2i(*m.GhostColAttr))})
	}
	if m.EditAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "edit"},
			Value: fmt.Sprintf("%d", b2i(*m.EditAttr))})
	}
	if m.DeleteAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "delete"},
			Value: fmt.Sprintf("%d", b2i(*m.DeleteAttr))})
	}
	if m.CopyAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "copy"},
			Value: fmt.Sprintf("%d", b2i(*m.CopyAttr))})
	}
	if m.PasteAllAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "pasteAll"},
			Value: fmt.Sprintf("%d", b2i(*m.PasteAllAttr))})
	}
	if m.PasteFormulasAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "pasteFormulas"},
			Value: fmt.Sprintf("%d", b2i(*m.PasteFormulasAttr))})
	}
	if m.PasteValuesAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "pasteValues"},
			Value: fmt.Sprintf("%d", b2i(*m.PasteValuesAttr))})
	}
	if m.PasteFormatsAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "pasteFormats"},
			Value: fmt.Sprintf("%d", b2i(*m.PasteFormatsAttr))})
	}
	if m.PasteCommentsAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "pasteComments"},
			Value: fmt.Sprintf("%d", b2i(*m.PasteCommentsAttr))})
	}
	if m.PasteDataValidationAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "pasteDataValidation"},
			Value: fmt.Sprintf("%d", b2i(*m.PasteDataValidationAttr))})
	}
	if m.PasteBordersAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "pasteBorders"},
			Value: fmt.Sprintf("%d", b2i(*m.PasteBordersAttr))})
	}
	if m.PasteColWidthsAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "pasteColWidths"},
			Value: fmt.Sprintf("%d", b2i(*m.PasteColWidthsAttr))})
	}
	if m.PasteNumberFormatsAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "pasteNumberFormats"},
			Value: fmt.Sprintf("%d", b2i(*m.PasteNumberFormatsAttr))})
	}
	if m.MergeAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "merge"},
			Value: fmt.Sprintf("%d", b2i(*m.MergeAttr))})
	}
	if m.SplitFirstAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "splitFirst"},
			Value: fmt.Sprintf("%d", b2i(*m.SplitFirstAttr))})
	}
	if m.SplitAllAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "splitAll"},
			Value: fmt.Sprintf("%d", b2i(*m.SplitAllAttr))})
	}
	if m.RowColShiftAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "rowColShift"},
			Value: fmt.Sprintf("%d", b2i(*m.RowColShiftAttr))})
	}
	if m.ClearAllAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "clearAll"},
			Value: fmt.Sprintf("%d", b2i(*m.ClearAllAttr))})
	}
	if m.ClearFormatsAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "clearFormats"},
			Value: fmt.Sprintf("%d", b2i(*m.ClearFormatsAttr))})
	}
	if m.ClearContentsAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "clearContents"},
			Value: fmt.Sprintf("%d", b2i(*m.ClearContentsAttr))})
	}
	if m.ClearCommentsAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "clearComments"},
			Value: fmt.Sprintf("%d", b2i(*m.ClearCommentsAttr))})
	}
	if m.AssignAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "assign"},
			Value: fmt.Sprintf("%d", b2i(*m.AssignAttr))})
	}
	if m.CoerceAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "coerce"},
			Value: fmt.Sprintf("%d", b2i(*m.CoerceAttr))})
	}
	if m.AdjustAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "adjust"},
			Value: fmt.Sprintf("%d", b2i(*m.AdjustAttr))})
	}
	if m.CellMetaAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "cellMeta"},
			Value: fmt.Sprintf("%d", b2i(*m.CellMetaAttr))})
	}
	e.EncodeToken(start)
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (m *CT_MetadataType) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initialize to default
	for _, attr := range start.Attr {
		if attr.Name.Local == "pasteColWidths" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.PasteColWidthsAttr = &parsed
			continue
		}
		if attr.Name.Local == "name" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.NameAttr = parsed
			continue
		}
		if attr.Name.Local == "pasteNumberFormats" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.PasteNumberFormatsAttr = &parsed
			continue
		}
		if attr.Name.Local == "ghostRow" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.GhostRowAttr = &parsed
			continue
		}
		if attr.Name.Local == "merge" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.MergeAttr = &parsed
			continue
		}
		if attr.Name.Local == "edit" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.EditAttr = &parsed
			continue
		}
		if attr.Name.Local == "splitFirst" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.SplitFirstAttr = &parsed
			continue
		}
		if attr.Name.Local == "copy" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.CopyAttr = &parsed
			continue
		}
		if attr.Name.Local == "splitAll" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.SplitAllAttr = &parsed
			continue
		}
		if attr.Name.Local == "pasteFormulas" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.PasteFormulasAttr = &parsed
			continue
		}
		if attr.Name.Local == "cellMeta" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.CellMetaAttr = &parsed
			continue
		}
		if attr.Name.Local == "clearAll" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.ClearAllAttr = &parsed
			continue
		}
		if attr.Name.Local == "minSupportedVersion" {
			parsed, err := strconv.ParseUint(attr.Value, 10, 32)
			if err != nil {
				return err
			}
			m.MinSupportedVersionAttr = uint32(parsed)
			continue
		}
		if attr.Name.Local == "adjust" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.AdjustAttr = &parsed
			continue
		}
		if attr.Name.Local == "clearContents" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.ClearContentsAttr = &parsed
			continue
		}
		if attr.Name.Local == "pasteValues" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.PasteValuesAttr = &parsed
			continue
		}
		if attr.Name.Local == "rowColShift" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.RowColShiftAttr = &parsed
			continue
		}
		if attr.Name.Local == "pasteComments" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.PasteCommentsAttr = &parsed
			continue
		}
		if attr.Name.Local == "clearFormats" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.ClearFormatsAttr = &parsed
			continue
		}
		if attr.Name.Local == "ghostCol" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.GhostColAttr = &parsed
			continue
		}
		if attr.Name.Local == "coerce" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.CoerceAttr = &parsed
			continue
		}
		if attr.Name.Local == "clearComments" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.ClearCommentsAttr = &parsed
			continue
		}
		if attr.Name.Local == "pasteAll" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.PasteAllAttr = &parsed
			continue
		}
		if attr.Name.Local == "pasteBorders" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.PasteBordersAttr = &parsed
			continue
		}
		if attr.Name.Local == "pasteFormats" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.PasteFormatsAttr = &parsed
			continue
		}
		if attr.Name.Local == "pasteDataValidation" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.PasteDataValidationAttr = &parsed
			continue
		}
		if attr.Name.Local == "delete" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.DeleteAttr = &parsed
			continue
		}
		if attr.Name.Local == "assign" {
			parsed, err := strconv.ParseBool(attr.Value)
			if err != nil {
				return err
			}
			m.AssignAttr = &parsed
			continue
		}
	}
	// skip any extensions we may find, but don't support
	for {
		tok, err := d.Token()
		if err != nil {
			return fmt.Errorf("parsing CT_MetadataType: %s", err)
		}
		if el, ok := tok.(xml.EndElement); ok && el.Name == start.Name {
			break
		}
	}
	return nil
}

// Validate validates the CT_MetadataType and its children
func (m *CT_MetadataType) Validate() error {
	return m.ValidateWithPath("CT_MetadataType")
}

// ValidateWithPath validates the CT_MetadataType and its children, prefixing error messages with path
func (m *CT_MetadataType) ValidateWithPath(path string) error {
	return nil
}
