// Copyright 2017 FoxyUtils ehf. All rights reserved.
//
// DO NOT EDIT: generated by gooxml ECMA-376 generator
//
// Use of this source code is governed by the terms of the Affero GNU General
// Public License version 3.0 as published by the Free Software Foundation and
// appearing in the file LICENSE included in the packaging of this file. A
// commercial license can be purchased via https://unidoc.io website.

package vml

import (
	"encoding/xml"

	"skygo_detection/lib/unioffice"
)

type OfcCT_ShapeLayout struct {
	Idmap        *OfcCT_IdMap
	Regrouptable *OfcCT_RegroupTable
	Rules        *OfcCT_Rules
	ExtAttr      ST_Ext
}

func NewOfcCT_ShapeLayout() *OfcCT_ShapeLayout {
	ret := &OfcCT_ShapeLayout{}
	return ret
}

func (m *OfcCT_ShapeLayout) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if m.ExtAttr != ST_ExtUnset {
		attr, err := m.ExtAttr.MarshalXMLAttr(xml.Name{Local: "ext"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	e.EncodeToken(start)
	if m.Idmap != nil {
		seidmap := xml.StartElement{Name: xml.Name{Local: "o:idmap"}}
		e.EncodeElement(m.Idmap, seidmap)
	}
	if m.Regrouptable != nil {
		seregrouptable := xml.StartElement{Name: xml.Name{Local: "o:regrouptable"}}
		e.EncodeElement(m.Regrouptable, seregrouptable)
	}
	if m.Rules != nil {
		serules := xml.StartElement{Name: xml.Name{Local: "o:rules"}}
		e.EncodeElement(m.Rules, serules)
	}
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (m *OfcCT_ShapeLayout) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initialize to default
	for _, attr := range start.Attr {
		if attr.Name.Local == "ext" {
			m.ExtAttr.UnmarshalXMLAttr(attr)
			continue
		}
	}
lOfcCT_ShapeLayout:
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch el := tok.(type) {
		case xml.StartElement:
			switch el.Name {
			case xml.Name{Space: "urn:schemas-microsoft-com:office:office", Local: "idmap"}:
				m.Idmap = NewOfcCT_IdMap()
				if err := d.DecodeElement(m.Idmap, &el); err != nil {
					return err
				}
			case xml.Name{Space: "urn:schemas-microsoft-com:office:office", Local: "regrouptable"}:
				m.Regrouptable = NewOfcCT_RegroupTable()
				if err := d.DecodeElement(m.Regrouptable, &el); err != nil {
					return err
				}
			case xml.Name{Space: "urn:schemas-microsoft-com:office:office", Local: "rules"}:
				m.Rules = NewOfcCT_Rules()
				if err := d.DecodeElement(m.Rules, &el); err != nil {
					return err
				}
			default:
				unioffice.Log("skipping unsupported element on OfcCT_ShapeLayout %v", el.Name)
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			break lOfcCT_ShapeLayout
		case xml.CharData:
		}
	}
	return nil
}

// Validate validates the OfcCT_ShapeLayout and its children
func (m *OfcCT_ShapeLayout) Validate() error {
	return m.ValidateWithPath("OfcCT_ShapeLayout")
}

// ValidateWithPath validates the OfcCT_ShapeLayout and its children, prefixing error messages with path
func (m *OfcCT_ShapeLayout) ValidateWithPath(path string) error {
	if m.Idmap != nil {
		if err := m.Idmap.ValidateWithPath(path + "/Idmap"); err != nil {
			return err
		}
	}
	if m.Regrouptable != nil {
		if err := m.Regrouptable.ValidateWithPath(path + "/Regrouptable"); err != nil {
			return err
		}
	}
	if m.Rules != nil {
		if err := m.Rules.ValidateWithPath(path + "/Rules"); err != nil {
			return err
		}
	}
	if err := m.ExtAttr.ValidateWithPath(path + "/ExtAttr"); err != nil {
		return err
	}
	return nil
}
