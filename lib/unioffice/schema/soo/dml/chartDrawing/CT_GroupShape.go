// Copyright 2017 FoxyUtils ehf. All rights reserved.
//
// DO NOT EDIT: generated by gooxml ECMA-376 generator
//
// Use of this source code is governed by the terms of the Affero GNU General
// Public License version 3.0 as published by the Free Software Foundation and
// appearing in the file LICENSE included in the packaging of this file. A
// commercial license can be purchased via https://unidoc.io website.

package chartDrawing

import (
	"encoding/xml"
	"fmt"

	"skygo_detection/lib/unioffice"
	"skygo_detection/lib/unioffice/schema/soo/dml"
)

type CT_GroupShape struct {
	NvGrpSpPr *CT_GroupShapeNonVisual
	GrpSpPr   *dml.CT_GroupShapeProperties
	Choice    []*CT_GroupShapeChoice
}

func NewCT_GroupShape() *CT_GroupShape {
	ret := &CT_GroupShape{}
	ret.NvGrpSpPr = NewCT_GroupShapeNonVisual()
	ret.GrpSpPr = dml.NewCT_GroupShapeProperties()
	return ret
}

func (m *CT_GroupShape) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	e.EncodeToken(start)
	senvGrpSpPr := xml.StartElement{Name: xml.Name{Local: "nvGrpSpPr"}}
	e.EncodeElement(m.NvGrpSpPr, senvGrpSpPr)
	segrpSpPr := xml.StartElement{Name: xml.Name{Local: "grpSpPr"}}
	e.EncodeElement(m.GrpSpPr, segrpSpPr)
	if m.Choice != nil {
		for _, c := range m.Choice {
			c.MarshalXML(e, xml.StartElement{})
		}
	}
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (m *CT_GroupShape) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initialize to default
	m.NvGrpSpPr = NewCT_GroupShapeNonVisual()
	m.GrpSpPr = dml.NewCT_GroupShapeProperties()
lCT_GroupShape:
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch el := tok.(type) {
		case xml.StartElement:
			switch el.Name {
			case xml.Name{Space: "http://schemas.openxmlformats.org/drawingml/2006/chartDrawing", Local: "nvGrpSpPr"}:
				if err := d.DecodeElement(m.NvGrpSpPr, &el); err != nil {
					return err
				}
			case xml.Name{Space: "http://schemas.openxmlformats.org/drawingml/2006/chartDrawing", Local: "grpSpPr"}:
				if err := d.DecodeElement(m.GrpSpPr, &el); err != nil {
					return err
				}
			case xml.Name{Space: "http://schemas.openxmlformats.org/drawingml/2006/chartDrawing", Local: "sp"}:
				tmp := NewCT_GroupShapeChoice()
				if err := d.DecodeElement(&tmp.Sp, &el); err != nil {
					return err
				}
				m.Choice = append(m.Choice, tmp)
			case xml.Name{Space: "http://schemas.openxmlformats.org/drawingml/2006/chartDrawing", Local: "grpSp"}:
				tmp := NewCT_GroupShapeChoice()
				if err := d.DecodeElement(&tmp.GrpSp, &el); err != nil {
					return err
				}
				m.Choice = append(m.Choice, tmp)
			case xml.Name{Space: "http://schemas.openxmlformats.org/drawingml/2006/chartDrawing", Local: "graphicFrame"}:
				tmp := NewCT_GroupShapeChoice()
				if err := d.DecodeElement(&tmp.GraphicFrame, &el); err != nil {
					return err
				}
				m.Choice = append(m.Choice, tmp)
			case xml.Name{Space: "http://schemas.openxmlformats.org/drawingml/2006/chartDrawing", Local: "cxnSp"}:
				tmp := NewCT_GroupShapeChoice()
				if err := d.DecodeElement(&tmp.CxnSp, &el); err != nil {
					return err
				}
				m.Choice = append(m.Choice, tmp)
			case xml.Name{Space: "http://schemas.openxmlformats.org/drawingml/2006/chartDrawing", Local: "pic"}:
				tmp := NewCT_GroupShapeChoice()
				if err := d.DecodeElement(&tmp.Pic, &el); err != nil {
					return err
				}
				m.Choice = append(m.Choice, tmp)
			default:
				unioffice.Log("skipping unsupported element on CT_GroupShape %v", el.Name)
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			break lCT_GroupShape
		case xml.CharData:
		}
	}
	return nil
}

// Validate validates the CT_GroupShape and its children
func (m *CT_GroupShape) Validate() error {
	return m.ValidateWithPath("CT_GroupShape")
}

// ValidateWithPath validates the CT_GroupShape and its children, prefixing error messages with path
func (m *CT_GroupShape) ValidateWithPath(path string) error {
	if err := m.NvGrpSpPr.ValidateWithPath(path + "/NvGrpSpPr"); err != nil {
		return err
	}
	if err := m.GrpSpPr.ValidateWithPath(path + "/GrpSpPr"); err != nil {
		return err
	}
	for i, v := range m.Choice {
		if err := v.ValidateWithPath(fmt.Sprintf("%s/Choice[%d]", path, i)); err != nil {
			return err
		}
	}
	return nil
}