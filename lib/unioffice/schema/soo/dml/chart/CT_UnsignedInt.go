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
	"strconv"
)

type CT_UnsignedInt struct {
	ValAttr uint32
}

func NewCT_UnsignedInt() *CT_UnsignedInt {
	ret := &CT_UnsignedInt{}
	return ret
}

func (m *CT_UnsignedInt) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "val"},
		Value: fmt.Sprintf("%v", m.ValAttr)})
	e.EncodeToken(start)
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (m *CT_UnsignedInt) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initialize to default
	for _, attr := range start.Attr {
		if attr.Name.Local == "val" {
			parsed, err := strconv.ParseUint(attr.Value, 10, 32)
			if err != nil {
				return err
			}
			m.ValAttr = uint32(parsed)
			continue
		}
	}
	// skip any extensions we may find, but don't support
	for {
		tok, err := d.Token()
		if err != nil {
			return fmt.Errorf("parsing CT_UnsignedInt: %s", err)
		}
		if el, ok := tok.(xml.EndElement); ok && el.Name == start.Name {
			break
		}
	}
	return nil
}

// Validate validates the CT_UnsignedInt and its children
func (m *CT_UnsignedInt) Validate() error {
	return m.ValidateWithPath("CT_UnsignedInt")
}

// ValidateWithPath validates the CT_UnsignedInt and its children, prefixing error messages with path
func (m *CT_UnsignedInt) ValidateWithPath(path string) error {
	return nil
}
