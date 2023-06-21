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
	"fmt"

	"skygo_detection/lib/unioffice/schema/soo/ofc/sharedTypes"
)

type CT_TextPath struct {
	OnAttr       sharedTypes.ST_TrueFalse
	FitshapeAttr sharedTypes.ST_TrueFalse
	FitpathAttr  sharedTypes.ST_TrueFalse
	TrimAttr     sharedTypes.ST_TrueFalse
	XscaleAttr   sharedTypes.ST_TrueFalse
	StringAttr   *string
	IdAttr       *string
	StyleAttr    *string
}

func NewCT_TextPath() *CT_TextPath {
	ret := &CT_TextPath{}
	return ret
}

func (m *CT_TextPath) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if m.OnAttr != sharedTypes.ST_TrueFalseUnset {
		attr, err := m.OnAttr.MarshalXMLAttr(xml.Name{Local: "on"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.FitshapeAttr != sharedTypes.ST_TrueFalseUnset {
		attr, err := m.FitshapeAttr.MarshalXMLAttr(xml.Name{Local: "fitshape"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.FitpathAttr != sharedTypes.ST_TrueFalseUnset {
		attr, err := m.FitpathAttr.MarshalXMLAttr(xml.Name{Local: "fitpath"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.TrimAttr != sharedTypes.ST_TrueFalseUnset {
		attr, err := m.TrimAttr.MarshalXMLAttr(xml.Name{Local: "trim"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.XscaleAttr != sharedTypes.ST_TrueFalseUnset {
		attr, err := m.XscaleAttr.MarshalXMLAttr(xml.Name{Local: "xscale"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.StringAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "string"},
			Value: fmt.Sprintf("%v", *m.StringAttr)})
	}
	if m.IdAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "id"},
			Value: fmt.Sprintf("%v", *m.IdAttr)})
	}
	if m.StyleAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "style"},
			Value: fmt.Sprintf("%v", *m.StyleAttr)})
	}
	e.EncodeToken(start)
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (m *CT_TextPath) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initialize to default
	for _, attr := range start.Attr {
		if attr.Name.Local == "on" {
			m.OnAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "fitshape" {
			m.FitshapeAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "fitpath" {
			m.FitpathAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "trim" {
			m.TrimAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "xscale" {
			m.XscaleAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "string" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.StringAttr = &parsed
			continue
		}
		if attr.Name.Local == "id" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.IdAttr = &parsed
			continue
		}
		if attr.Name.Local == "style" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.StyleAttr = &parsed
			continue
		}
	}
	// skip any extensions we may find, but don't support
	for {
		tok, err := d.Token()
		if err != nil {
			return fmt.Errorf("parsing CT_TextPath: %s", err)
		}
		if el, ok := tok.(xml.EndElement); ok && el.Name == start.Name {
			break
		}
	}
	return nil
}

// Validate validates the CT_TextPath and its children
func (m *CT_TextPath) Validate() error {
	return m.ValidateWithPath("CT_TextPath")
}

// ValidateWithPath validates the CT_TextPath and its children, prefixing error messages with path
func (m *CT_TextPath) ValidateWithPath(path string) error {
	if err := m.OnAttr.ValidateWithPath(path + "/OnAttr"); err != nil {
		return err
	}
	if err := m.FitshapeAttr.ValidateWithPath(path + "/FitshapeAttr"); err != nil {
		return err
	}
	if err := m.FitpathAttr.ValidateWithPath(path + "/FitpathAttr"); err != nil {
		return err
	}
	if err := m.TrimAttr.ValidateWithPath(path + "/TrimAttr"); err != nil {
		return err
	}
	if err := m.XscaleAttr.ValidateWithPath(path + "/XscaleAttr"); err != nil {
		return err
	}
	return nil
}