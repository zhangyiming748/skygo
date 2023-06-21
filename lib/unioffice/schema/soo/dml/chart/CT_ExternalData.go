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

type CT_ExternalData struct {
	IdAttr     string
	AutoUpdate *CT_Boolean
}

func NewCT_ExternalData() *CT_ExternalData {
	ret := &CT_ExternalData{}
	return ret
}

func (m *CT_ExternalData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "r:id"},
		Value: fmt.Sprintf("%v", m.IdAttr)})
	e.EncodeToken(start)
	if m.AutoUpdate != nil {
		seautoUpdate := xml.StartElement{Name: xml.Name{Local: "c:autoUpdate"}}
		e.EncodeElement(m.AutoUpdate, seautoUpdate)
	}
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (m *CT_ExternalData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initialize to default
	for _, attr := range start.Attr {
		if attr.Name.Space == "http://schemas.openxmlformats.org/officeDocument/2006/relationships" && attr.Name.Local == "id" ||
			attr.Name.Space == "http://purl.oclc.org/ooxml/officeDocument/relationships" && attr.Name.Local == "id" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.IdAttr = parsed
			continue
		}
	}
lCT_ExternalData:
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch el := tok.(type) {
		case xml.StartElement:
			switch el.Name {
			case xml.Name{Space: "http://schemas.openxmlformats.org/drawingml/2006/chart", Local: "autoUpdate"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/drawingml/chart", Local: "autoUpdate"}:
				m.AutoUpdate = NewCT_Boolean()
				if err := d.DecodeElement(m.AutoUpdate, &el); err != nil {
					return err
				}
			default:
				unioffice.Log("skipping unsupported element on CT_ExternalData %v", el.Name)
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			break lCT_ExternalData
		case xml.CharData:
		}
	}
	return nil
}

// Validate validates the CT_ExternalData and its children
func (m *CT_ExternalData) Validate() error {
	return m.ValidateWithPath("CT_ExternalData")
}

// ValidateWithPath validates the CT_ExternalData and its children, prefixing error messages with path
func (m *CT_ExternalData) ValidateWithPath(path string) error {
	if m.AutoUpdate != nil {
		if err := m.AutoUpdate.ValidateWithPath(path + "/AutoUpdate"); err != nil {
			return err
		}
	}
	return nil
}
