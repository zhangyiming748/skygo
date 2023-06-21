// Copyright 2017 FoxyUtils ehf. All rights reserved.
//
// DO NOT EDIT: generated by gooxml ECMA-376 generator
//
// Use of this source code is governed by the terms of the Affero GNU General
// Public License version 3.0 as published by the Free Software Foundation and
// appearing in the file LICENSE included in the packaging of this file. A
// commercial license can be purchased via https://unidoc.io website.

package wml

import (
	"encoding/xml"
	"fmt"

	"skygo_detection/lib/unioffice"
)

type CT_AltChunk struct {
	IdAttr *string
	// External Content Import Properties
	AltChunkPr *CT_AltChunkPr
}

func NewCT_AltChunk() *CT_AltChunk {
	ret := &CT_AltChunk{}
	return ret
}

func (m *CT_AltChunk) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if m.IdAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "r:id"},
			Value: fmt.Sprintf("%v", *m.IdAttr)})
	}
	e.EncodeToken(start)
	if m.AltChunkPr != nil {
		sealtChunkPr := xml.StartElement{Name: xml.Name{Local: "w:altChunkPr"}}
		e.EncodeElement(m.AltChunkPr, sealtChunkPr)
	}
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (m *CT_AltChunk) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
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
	}
lCT_AltChunk:
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch el := tok.(type) {
		case xml.StartElement:
			switch el.Name {
			case xml.Name{Space: "http://schemas.openxmlformats.org/wordprocessingml/2006/main", Local: "altChunkPr"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/wordprocessingml/main", Local: "altChunkPr"}:
				m.AltChunkPr = NewCT_AltChunkPr()
				if err := d.DecodeElement(m.AltChunkPr, &el); err != nil {
					return err
				}
			default:
				unioffice.Log("skipping unsupported element on CT_AltChunk %v", el.Name)
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			break lCT_AltChunk
		case xml.CharData:
		}
	}
	return nil
}

// Validate validates the CT_AltChunk and its children
func (m *CT_AltChunk) Validate() error {
	return m.ValidateWithPath("CT_AltChunk")
}

// ValidateWithPath validates the CT_AltChunk and its children, prefixing error messages with path
func (m *CT_AltChunk) ValidateWithPath(path string) error {
	if m.AltChunkPr != nil {
		if err := m.AltChunkPr.ValidateWithPath(path + "/AltChunkPr"); err != nil {
			return err
		}
	}
	return nil
}
