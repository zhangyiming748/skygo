// Copyright 2017 FoxyUtils ehf. All rights reserved.
//
// DO NOT EDIT: generated by gooxml ECMA-376 generator
//
// Use of this source code is governed by the terms of the Affero GNU General
// Public License version 3.0 as published by the Free Software Foundation and
// appearing in the file LICENSE included in the packaging of this file. A
// commercial license can be purchased via https://unidoc.io website.

package pml

import (
	"encoding/xml"

	"skygo_detection/lib/unioffice"
	"skygo_detection/lib/unioffice/schema/soo/dml"
)

type NotesMaster struct {
	CT_NotesMaster
}

func NewNotesMaster() *NotesMaster {
	ret := &NotesMaster{}
	ret.CT_NotesMaster = *NewCT_NotesMaster()
	return ret
}

func (m *NotesMaster) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xmlns"}, Value: "http://schemas.openxmlformats.org/presentationml/2006/main"})
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xmlns:a"}, Value: "http://schemas.openxmlformats.org/drawingml/2006/main"})
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xmlns:p"}, Value: "http://schemas.openxmlformats.org/presentationml/2006/main"})
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xmlns:r"}, Value: "http://schemas.openxmlformats.org/officeDocument/2006/relationships"})
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xmlns:sh"}, Value: "http://schemas.openxmlformats.org/officeDocument/2006/sharedTypes"})
	start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "xmlns:xml"}, Value: "http://www.w3.org/XML/1998/namespace"})
	start.Name.Local = "p:notesMaster"
	return m.CT_NotesMaster.MarshalXML(e, start)
}

func (m *NotesMaster) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initialize to default
	m.CT_NotesMaster = *NewCT_NotesMaster()
lNotesMaster:
	for {
		tok, err := d.Token()
		if err != nil {
			return err
		}
		switch el := tok.(type) {
		case xml.StartElement:
			switch el.Name {
			case xml.Name{Space: "http://schemas.openxmlformats.org/presentationml/2006/main", Local: "cSld"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/presentationml/main", Local: "cSld"}:
				if err := d.DecodeElement(m.CSld, &el); err != nil {
					return err
				}
			case xml.Name{Space: "http://schemas.openxmlformats.org/presentationml/2006/main", Local: "clrMap"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/presentationml/main", Local: "clrMap"}:
				if err := d.DecodeElement(m.ClrMap, &el); err != nil {
					return err
				}
			case xml.Name{Space: "http://schemas.openxmlformats.org/presentationml/2006/main", Local: "hf"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/presentationml/main", Local: "hf"}:
				m.Hf = NewCT_HeaderFooter()
				if err := d.DecodeElement(m.Hf, &el); err != nil {
					return err
				}
			case xml.Name{Space: "http://schemas.openxmlformats.org/presentationml/2006/main", Local: "notesStyle"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/presentationml/main", Local: "notesStyle"}:
				m.NotesStyle = dml.NewCT_TextListStyle()
				if err := d.DecodeElement(m.NotesStyle, &el); err != nil {
					return err
				}
			case xml.Name{Space: "http://schemas.openxmlformats.org/presentationml/2006/main", Local: "extLst"},
				xml.Name{Space: "http://purl.oclc.org/ooxml/presentationml/main", Local: "extLst"}:
				m.ExtLst = NewCT_ExtensionListModify()
				if err := d.DecodeElement(m.ExtLst, &el); err != nil {
					return err
				}
			default:
				unioffice.Log("skipping unsupported element on NotesMaster %v", el.Name)
				if err := d.Skip(); err != nil {
					return err
				}
			}
		case xml.EndElement:
			break lNotesMaster
		case xml.CharData:
		}
	}
	return nil
}

// Validate validates the NotesMaster and its children
func (m *NotesMaster) Validate() error {
	return m.ValidateWithPath("NotesMaster")
}

// ValidateWithPath validates the NotesMaster and its children, prefixing error messages with path
func (m *NotesMaster) ValidateWithPath(path string) error {
	if err := m.CT_NotesMaster.ValidateWithPath(path); err != nil {
		return err
	}
	return nil
}
