// Copyright 2017 FoxyUtils ehf. All rights reserved.
//
// DO NOT EDIT: generated by gooxml ECMA-376 generator
//
// Use of this source code is governed by the terms of the Affero GNU General
// Public License version 3.0 as published by the Free Software Foundation and
// appearing in the file LICENSE included in the packaging of this file. A
// commercial license can be purchased via https://unidoc.io website.

package math_test

import (
	"encoding/xml"
	"testing"

	"skygo_detection/lib/unioffice/schema/soo/ofc/math"
)

func TestCT_RChoiceConstructor(t *testing.T) {
	v := math.NewCT_RChoice()
	if v == nil {
		t.Errorf("math.NewCT_RChoice must return a non-nil value")
	}
	if err := v.Validate(); err != nil {
		t.Errorf("newly constructed math.CT_RChoice should validate: %s", err)
	}
}

func TestCT_RChoiceMarshalUnmarshal(t *testing.T) {
	v := math.NewCT_RChoice()
	buf, _ := xml.Marshal(v)
	v2 := math.NewCT_RChoice()
	xml.Unmarshal(buf, v2)
}
