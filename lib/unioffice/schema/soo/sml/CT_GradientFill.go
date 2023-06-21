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

	"skygo_detection/lib/unioffice"
)

type CT_GradientFill struct {
	// Gradient Fill Type
	TypeAttr ST_GradientType
	// Linear Gradient Degree
	DegreeAttr *float64
	// Left Convergence
	LeftAttr *float64
	// Right Convergence
	RightAttr *float64
	// Top Gradient Convergence
	TopAttr *float64
	// Bottom Convergence
	BottomAttr *float64
	// Gradient Stop
	Stop []*CT_GradientStop
}

func NewCT_GradientFill() *CT_GradientFill {
	ret := &CT_GradientFill{}
	return ret
}

func (m *CT_GradientFill) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if m.TypeAttr != ST_GradientTypeUnset {
		attr, err := m.TypeAttr.MarshalXMLAttr(xml.Name{Local: "type"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.DegreeAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "degree"},
			Value: fmt.Sprintf("%v", *m.DegreeAttr)})
	}
	if m.LeftAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "left"},
			Value: fmt.Sprintf("%v", *m.LeftAttr)})
	}
	if m.RightAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "right"},
			Value: fmt.Sprintf("%v", *m.RightAttr)})
	}
	if m.TopAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "top"},
			Value: fmt.Sprintf("%v", *m.TopAttr)})
	}
	if m.BottomAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "bottom"},
			Value: fmt.Sprintf("%v", *m.BottomAttr)})
	}
	e.EncodeToken(start)
	if m.Stop != nil {
		sestop := xml.StartElement{Name: xml.Name{Local: "ma:stop"}}
		for _, c := range m.Stop {
			e.EncodeElement(c, sestop)
		}
	}
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (m *CT_GradientFill) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initialize to default
	for _, attr := range start.Attr {
		if attr.Name.Local == "type" {
			m.TypeAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "degree" {
			parsed, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return err
			}
			m.DegreeAttr = &parsed
			continue
		}
		if attr.Name.Local == "left" {
			parsed, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return err
			}
			m.LeftAttr = &parsed
			continue
		}
		if attr.Name.Local == "right" {
			parsed, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return err
			}
			m.RightAttr = &parsed
			continue
		}
		if attr.Name.Local == "top" {
			parsed, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return err
			}
			m.TopAttr = &parsed
			continue
		}
		if attr.Name.Local == "bottom" {
			parsed, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return err
			}
			m.BottomAttr = &parsed
			continue
		}
	}
lCT_GradientFill:
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch el := tok.(type) {
		case xml.StartElement:
			switch el.Name {
			case xml.Name{Space: "http://schemas.openxmlformats.org/spreadsheetml/2006/main", Local: "stop"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/spreadsheetml/main", Local: "stop"}:
				tmp := NewCT_GradientStop()
				if err := d.DecodeElement(tmp, &el); err != nil {
					return err
				}
				m.Stop = append(m.Stop, tmp)
			default:
				unioffice.Log("skipping unsupported element on CT_GradientFill %v", el.Name)
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			break lCT_GradientFill
		case xml.CharData:
		}
	}
	return nil
}

// Validate validates the CT_GradientFill and its children
func (m *CT_GradientFill) Validate() error {
	return m.ValidateWithPath("CT_GradientFill")
}

// ValidateWithPath validates the CT_GradientFill and its children, prefixing error messages with path
func (m *CT_GradientFill) ValidateWithPath(path string) error {
	if err := m.TypeAttr.ValidateWithPath(path + "/TypeAttr"); err != nil {
		return err
	}
	for i, v := range m.Stop {
		if err := v.ValidateWithPath(fmt.Sprintf("%s/Stop[%d]", path, i)); err != nil {
			return err
		}
	}
	return nil
}
