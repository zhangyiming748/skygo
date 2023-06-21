// Copyright 2017 FoxyUtils ehf. All rights reserved.
//
// DO NOT EDIT: generated by gooxml ECMA-376 generator
//
// Use of this source code is governed by the terms of the Affero GNU General
// Public License version 3.0 as published by the Free Software Foundation and
// appearing in the file LICENSE included in the packaging of this file. A
// commercial license can be purchased via https://unidoc.io website.

package sml_test

import (
	"encoding/xml"
	"testing"

	"skygo_detection/lib/unioffice/schema/soo/sml"
)

func TestCT_FunctionGroupsConstructor(t *testing.T) {
	v := sml.NewCT_FunctionGroups()
	if v == nil {
		t.Errorf("sml.NewCT_FunctionGroups must return a non-nil value")
	}
	if err := v.Validate(); err != nil {
		t.Errorf("newly constructed sml.CT_FunctionGroups should validate: %s", err)
	}
}

func TestCT_FunctionGroupsMarshalUnmarshal(t *testing.T) {
	v := sml.NewCT_FunctionGroups()
	buf, _ := xml.Marshal(v)
	v2 := sml.NewCT_FunctionGroups()
	xml.Unmarshal(buf, v2)
}
