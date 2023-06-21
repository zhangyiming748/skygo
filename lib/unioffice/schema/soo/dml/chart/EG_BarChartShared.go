// Copyright 2017 FoxyUtils ehf. All rights reserved.
//
// DO NOT EDIT: generated by gooxml ECMA-376 generator
//
// Use of this source code is governed by the terms of the Affero GNU General
// Public License version 3.0 as published by the Free Software Foundation and
// appearing in the file LICENSE included in the packaging of this file. A
// commercial license can be purchased via https://unidoc.io website.

package chart

import (
	"encoding/xml"
	"fmt"

	"skygo_detection/lib/unioffice"
)

type EG_BarChartShared struct {
	BarDir     *CT_BarDir
	Grouping   *CT_BarGrouping
	VaryColors *CT_Boolean
	Ser        []*CT_BarSer
	DLbls      *CT_DLbls
}

func NewEG_BarChartShared() *EG_BarChartShared {
	ret := &EG_BarChartShared{}
	ret.BarDir = NewCT_BarDir()
	return ret
}

func (m *EG_BarChartShared) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	sebarDir := xml.StartElement{Name: xml.Name{Local: "c:barDir"}}
	e.EncodeElement(m.BarDir, sebarDir)
	if m.Grouping != nil {
		segrouping := xml.StartElement{Name: xml.Name{Local: "c:grouping"}}
		e.EncodeElement(m.Grouping, segrouping)
	}
	if m.VaryColors != nil {
		sevaryColors := xml.StartElement{Name: xml.Name{Local: "c:varyColors"}}
		e.EncodeElement(m.VaryColors, sevaryColors)
	}
	if m.Ser != nil {
		seser := xml.StartElement{Name: xml.Name{Local: "c:ser"}}
		for _, c := range m.Ser {
			e.EncodeElement(c, seser)
		}
	}
	if m.DLbls != nil {
		sedLbls := xml.StartElement{Name: xml.Name{Local: "c:dLbls"}}
		e.EncodeElement(m.DLbls, sedLbls)
	}
	return nil
}

func (m *EG_BarChartShared) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initialize to default
	m.BarDir = NewCT_BarDir()
lEG_BarChartShared:
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch el := tok.(type) {
		case xml.StartElement:
			switch el.Name {
			case xml.Name{Space: "http://schemas.openxmlformats.org/drawingml/2006/chart", Local: "barDir"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/drawingml/chart", Local: "barDir"}:
				if err := d.DecodeElement(m.BarDir, &el); err != nil {
					return err
				}
			case xml.Name{Space: "http://schemas.openxmlformats.org/drawingml/2006/chart", Local: "grouping"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/drawingml/chart", Local: "grouping"}:
				m.Grouping = NewCT_BarGrouping()
				if err := d.DecodeElement(m.Grouping, &el); err != nil {
					return err
				}
			case xml.Name{Space: "http://schemas.openxmlformats.org/drawingml/2006/chart", Local: "varyColors"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/drawingml/chart", Local: "varyColors"}:
				m.VaryColors = NewCT_Boolean()
				if err := d.DecodeElement(m.VaryColors, &el); err != nil {
					return err
				}
			case xml.Name{Space: "http://schemas.openxmlformats.org/drawingml/2006/chart", Local: "ser"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/drawingml/chart", Local: "ser"}:
				tmp := NewCT_BarSer()
				if err := d.DecodeElement(tmp, &el); err != nil {
					return err
				}
				m.Ser = append(m.Ser, tmp)
			case xml.Name{Space: "http://schemas.openxmlformats.org/drawingml/2006/chart", Local: "dLbls"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/drawingml/chart", Local: "dLbls"}:
				m.DLbls = NewCT_DLbls()
				if err := d.DecodeElement(m.DLbls, &el); err != nil {
					return err
				}
			default:
				unioffice.Log("skipping unsupported element on EG_BarChartShared %v", el.Name)
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			break lEG_BarChartShared
		case xml.CharData:
		}
	}
	return nil
}

// Validate validates the EG_BarChartShared and its children
func (m *EG_BarChartShared) Validate() error {
	return m.ValidateWithPath("EG_BarChartShared")
}

// ValidateWithPath validates the EG_BarChartShared and its children, prefixing error messages with path
func (m *EG_BarChartShared) ValidateWithPath(path string) error {
	if err := m.BarDir.ValidateWithPath(path + "/BarDir"); err != nil {
		return err
	}
	if m.Grouping != nil {
		if err := m.Grouping.ValidateWithPath(path + "/Grouping"); err != nil {
			return err
		}
	}
	if m.VaryColors != nil {
		if err := m.VaryColors.ValidateWithPath(path + "/VaryColors"); err != nil {
			return err
		}
	}
	for i, v := range m.Ser {
		if err := v.ValidateWithPath(fmt.Sprintf("%s/Ser[%d]", path, i)); err != nil {
			return err
		}
	}
	if m.DLbls != nil {
		if err := m.DLbls.ValidateWithPath(path + "/DLbls"); err != nil {
			return err
		}
	}
	return nil
}
