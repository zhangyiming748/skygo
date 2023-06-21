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
	"strconv"

	"skygo_detection/lib/unioffice/schema/soo/ofc/sharedTypes"
)

type CT_ImageData struct {
	EmbosscolorAttr      *string
	RecolortargetAttr    *string
	HrefAttr             *string
	AlthrefAttr          *string
	TitleAttr            *string
	OleidAttr            *float32
	DetectmouseclickAttr sharedTypes.ST_TrueFalse
	MovieAttr            *float32
	RelidAttr            *string
	IdAttr               *string
	PictAttr             *string
	RHrefAttr            *string
	SIdAttr              *string
	SrcAttr              *string
	CropleftAttr         *string
	CroptopAttr          *string
	CroprightAttr        *string
	CropbottomAttr       *string
	GainAttr             *string
	BlacklevelAttr       *string
	GammaAttr            *string
	GrayscaleAttr        sharedTypes.ST_TrueFalse
	BilevelAttr          sharedTypes.ST_TrueFalse
	ChromakeyAttr        *string
}

func NewCT_ImageData() *CT_ImageData {
	ret := &CT_ImageData{}
	return ret
}

func (m *CT_ImageData) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	if m.EmbosscolorAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "embosscolor"},
			Value: fmt.Sprintf("%v", *m.EmbosscolorAttr)})
	}
	if m.RecolortargetAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "recolortarget"},
			Value: fmt.Sprintf("%v", *m.RecolortargetAttr)})
	}
	if m.HrefAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "o:href"},
			Value: fmt.Sprintf("%v", *m.HrefAttr)})
	}
	if m.AlthrefAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "o:althref"},
			Value: fmt.Sprintf("%v", *m.AlthrefAttr)})
	}
	if m.TitleAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "o:title"},
			Value: fmt.Sprintf("%v", *m.TitleAttr)})
	}
	if m.OleidAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "o:oleid"},
			Value: fmt.Sprintf("%v", *m.OleidAttr)})
	}
	if m.DetectmouseclickAttr != sharedTypes.ST_TrueFalseUnset {
		attr, err := m.DetectmouseclickAttr.MarshalXMLAttr(xml.Name{Local: "detectmouseclick"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.MovieAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "o:movie"},
			Value: fmt.Sprintf("%v", *m.MovieAttr)})
	}
	if m.RelidAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "o:relid"},
			Value: fmt.Sprintf("%v", *m.RelidAttr)})
	}
	if m.IdAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "r:id"},
			Value: fmt.Sprintf("%v", *m.IdAttr)})
	}
	if m.PictAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "r:pict"},
			Value: fmt.Sprintf("%v", *m.PictAttr)})
	}
	if m.RHrefAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "r:href"},
			Value: fmt.Sprintf("%v", *m.RHrefAttr)})
	}
	if m.SIdAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "id"},
			Value: fmt.Sprintf("%v", *m.SIdAttr)})
	}
	if m.SrcAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "src"},
			Value: fmt.Sprintf("%v", *m.SrcAttr)})
	}
	if m.CropleftAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "cropleft"},
			Value: fmt.Sprintf("%v", *m.CropleftAttr)})
	}
	if m.CroptopAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "croptop"},
			Value: fmt.Sprintf("%v", *m.CroptopAttr)})
	}
	if m.CroprightAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "cropright"},
			Value: fmt.Sprintf("%v", *m.CroprightAttr)})
	}
	if m.CropbottomAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "cropbottom"},
			Value: fmt.Sprintf("%v", *m.CropbottomAttr)})
	}
	if m.GainAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "gain"},
			Value: fmt.Sprintf("%v", *m.GainAttr)})
	}
	if m.BlacklevelAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "blacklevel"},
			Value: fmt.Sprintf("%v", *m.BlacklevelAttr)})
	}
	if m.GammaAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "gamma"},
			Value: fmt.Sprintf("%v", *m.GammaAttr)})
	}
	if m.GrayscaleAttr != sharedTypes.ST_TrueFalseUnset {
		attr, err := m.GrayscaleAttr.MarshalXMLAttr(xml.Name{Local: "grayscale"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.BilevelAttr != sharedTypes.ST_TrueFalseUnset {
		attr, err := m.BilevelAttr.MarshalXMLAttr(xml.Name{Local: "bilevel"})
		if err != nil {
			return err
		}
		start.Attr = append(start.Attr, attr)
	}
	if m.ChromakeyAttr != nil {
		start.Attr = append(start.Attr, xml.Attr{Name: xml.Name{Local: "chromakey"},
			Value: fmt.Sprintf("%v", *m.ChromakeyAttr)})
	}
	e.EncodeToken(start)
	e.EncodeToken(xml.EndElement{Name: start.Name})
	return nil
}

func (m *CT_ImageData) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	// initialize to default
	for _, attr := range start.Attr {
		if attr.Name.Space == "http://schemas.openxmlformats.org/officeDocument/2006/relationships" && attr.Name.Local == "pict" ||
			attr.Name.Space == "http://purl.oclc.org/ooxml/officeDocument/relationships" && attr.Name.Local == "pict" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.PictAttr = &parsed
			continue
		}
		if attr.Name.Space == "http://schemas.openxmlformats.org/officeDocument/2006/relationships" && attr.Name.Local == "href" ||
			attr.Name.Space == "http://purl.oclc.org/ooxml/officeDocument/relationships" && attr.Name.Local == "href" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.RHrefAttr = &parsed
			continue
		}
		if attr.Name.Space == "urn:schemas-microsoft-com:office:office" && attr.Name.Local == "href" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.HrefAttr = &parsed
			continue
		}
		if attr.Name.Space == "urn:schemas-microsoft-com:office:office" && attr.Name.Local == "althref" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.AlthrefAttr = &parsed
			continue
		}
		if attr.Name.Space == "urn:schemas-microsoft-com:office:office" && attr.Name.Local == "title" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.TitleAttr = &parsed
			continue
		}
		if attr.Name.Space == "urn:schemas-microsoft-com:office:office" && attr.Name.Local == "oleid" {
			parsed, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return err
			}
			pt := float32(parsed)
			m.OleidAttr = &pt
			continue
		}
		if attr.Name.Space == "urn:schemas-microsoft-com:office:office" && attr.Name.Local == "detectmouseclick" {
			m.DetectmouseclickAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Space == "urn:schemas-microsoft-com:office:office" && attr.Name.Local == "movie" {
			parsed, err := strconv.ParseFloat(attr.Value, 64)
			if err != nil {
				return err
			}
			pt := float32(parsed)
			m.MovieAttr = &pt
			continue
		}
		if attr.Name.Space == "urn:schemas-microsoft-com:office:office" && attr.Name.Local == "relid" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.RelidAttr = &parsed
			continue
		}
		if attr.Name.Space == "http://schemas.openxmlformats.org/officeDocument/2006/relationships" && attr.Name.Local == "id" ||
			attr.Name.Space == "http://purl.oclc.org/ooxml/officeDocument/relationships" && attr.Name.Local == "id" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.IdAttr = &parsed
			continue
		}
		if attr.Name.Local == "id" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.SIdAttr = &parsed
			continue
		}
		if attr.Name.Local == "cropbottom" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.CropbottomAttr = &parsed
			continue
		}
		if attr.Name.Local == "embosscolor" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.EmbosscolorAttr = &parsed
			continue
		}
		if attr.Name.Local == "src" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.SrcAttr = &parsed
			continue
		}
		if attr.Name.Local == "cropleft" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.CropleftAttr = &parsed
			continue
		}
		if attr.Name.Local == "croptop" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.CroptopAttr = &parsed
			continue
		}
		if attr.Name.Local == "cropright" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.CroprightAttr = &parsed
			continue
		}
		if attr.Name.Local == "recolortarget" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.RecolortargetAttr = &parsed
			continue
		}
		if attr.Name.Local == "gain" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.GainAttr = &parsed
			continue
		}
		if attr.Name.Local == "blacklevel" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.BlacklevelAttr = &parsed
			continue
		}
		if attr.Name.Local == "gamma" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.GammaAttr = &parsed
			continue
		}
		if attr.Name.Local == "grayscale" {
			m.GrayscaleAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "bilevel" {
			m.BilevelAttr.UnmarshalXMLAttr(attr)
			continue
		}
		if attr.Name.Local == "chromakey" {
			parsed, err := attr.Value, error(nil)
			if err != nil {
				return err
			}
			m.ChromakeyAttr = &parsed
			continue
		}
	}
	// skip any extensions we may find, but don't support
	for {
		tok, err := d.Token()
		if err != nil {
			return fmt.Errorf("parsing CT_ImageData: %s", err)
		}
		if el, ok := tok.(xml.EndElement); ok && el.Name == start.Name {
			break
		}
	}
	return nil
}

// Validate validates the CT_ImageData and its children
func (m *CT_ImageData) Validate() error {
	return m.ValidateWithPath("CT_ImageData")
}

// ValidateWithPath validates the CT_ImageData and its children, prefixing error messages with path
func (m *CT_ImageData) ValidateWithPath(path string) error {
	if err := m.DetectmouseclickAttr.ValidateWithPath(path + "/DetectmouseclickAttr"); err != nil {
		return err
	}
	if err := m.GrayscaleAttr.ValidateWithPath(path + "/GrayscaleAttr"); err != nil {
		return err
	}
	if err := m.BilevelAttr.ValidateWithPath(path + "/BilevelAttr"); err != nil {
		return err
	}
	return nil
}