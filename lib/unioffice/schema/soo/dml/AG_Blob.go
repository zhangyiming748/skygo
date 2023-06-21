// Copyright 2017 FoxyUtils ehf. All rights reserved.
//
// DO NOT EDIT: generated by gooxml ECMA-376 generator
//
// Use of this source code is governed by the terms of the Affero GNU General
// Public License version 3.0 as published by the Free Software Foundation and
// appearing in the file LICENSE included in the packaging of this file. A
// commercial license can be purchased via https://unidoc.io website.

package dml

import (
	"encoding/xml"
	"fmt"
)

type AG_Blob struct {
	EmbedAttr *string
	LinkAttr  *string
}

func NewAG_Blob() *AG_Blob {
	ret := &AG_Blob{}
	return ret
}

func (m *AG_Blob) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if m.EmbedAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "r:embed"},
			Value: fmt.Sprintf("%v", *m.EmbedAttr)})
	}
	if m.LinkAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "r:link"},
			Value: fmt.Sprintf("%v", *m.LinkAttr)})
	}
	return nil
}

func (m *AG_Blob) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initialize to default
	for _, attr := range start.Attr {
		if attr.Name.Space == "http://schemas.openxmlformats.org/officeDocument/2006/relationships" && attr.Name.Local == "embed" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.EmbedAttr = &parsed
			continue
		}
		if attr.Name.Space == "http://schemas.openxmlformats.org/officeDocument/2006/relationships" && attr.Name.Local == "link" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.LinkAttr = &parsed
			continue
		}
	}
	// skip any extensions we may find, but don't support
	for {
		tok, err := d.Token()
		if err != nil {
			return fmt.Errorf("parsing AG_Blob: %s", err)
		}
		if el, ok := tok.(xml.EndElement); ok && el.Name == start.Name {
			break
		}
	}
	return nil
}

// Validate validates the AG_Blob and its children
func (m *AG_Blob) Validate() error {
	return m.ValidateWithPath("AG_Blob")
}

// ValidateWithPath validates the AG_Blob and its children, prefixing error messages with path
func (m *AG_Blob) ValidateWithPath(path string) error {
	return nil
}
