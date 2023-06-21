// Copyright 2017 FoxyUtils ehf. All rights reserved.
//
// DO NOT EDIT: generated by gooxml ECMA-376 generator
//
// Use of this source code is governed by the terms of the Affero GNU General
// Public License version 3.0 as published by the Free Software Foundation and
// appearing in the file LICENSE included in the packaging of this file. A
// commercial license can be purchased via https://unidoc.io website.

package terms_test

import (
	"encoding/xml"
	"testing"

	"skygo_detection/lib/unioffice/schema/purl.org/dc/terms"
)

func TestElementsAndRefinementsGroupChoiceConstructor(t *testing.T) {
	v := terms.NewElementsAndRefinementsGroupChoice()
	if v == nil {
		t.Errorf("terms.NewElementsAndRefinementsGroupChoice must return a non-nil value")
	}
	if err := v.Validate(); err != nil {
		t.Errorf("newly constructed terms.ElementsAndRefinementsGroupChoice should validate: %s", err)
	}
}

func TestElementsAndRefinementsGroupChoiceMarshalUnmarshal(t *testing.T) {
	v := terms.NewElementsAndRefinementsGroupChoice()
	buf, _ := xml.Marshal(v)
	v2 := terms.NewElementsAndRefinementsGroupChoice()
	xml.Unmarshal(buf, v2)
}