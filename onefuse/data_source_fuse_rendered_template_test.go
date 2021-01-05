// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"testing"
)

func TestRenderedTemplate(t *testing.T) {
	config := GetConfig()

	template := "template {{templatedValue}}"
	templateProperties := make(map[string]interface{})
	templateProperties["templatedValue"] = "this is the value"

	renderedTemplate, err := config.NewOneFuseApiClient().RenderTemplate(template, templateProperties)

	if err != nil {
		t.Errorf("Error rendering template: '%s'", err)
		return
	}

	t.Logf("Rendered template: %v", renderedTemplate)

	if renderedTemplate.Value != "template this is the value" {
		t.Errorf("Error: Rendered template does not match expected. '%s'", err)
		return
	}
}
