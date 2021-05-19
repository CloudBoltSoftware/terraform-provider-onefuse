// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"crypto/sha256"
	"fmt"
	"log"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceRenderedTemplate() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceRenderedTemplateRead,
		Schema: map[string]*schema.Schema{
			"template": {
				Type:     schema.TypeString,
				Required: true,
			},
			"template_properties": {
				Type:     schema.TypeMap,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceRenderedTemplateRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("onefuse.dataSourceRenderedTemplateRead")

	config := meta.(Config)
	apiClient := config.NewOneFuseApiClient()

	renderedTemplate, err := apiClient.RenderTemplate(d.Get("template").(string), d.Get("template_properties").(map[string]interface{}))

	if err != nil {
		return fmt.Errorf("Error loading Rendered Template: %s", err)
	}

	// a resource needs an ID, otherwise it will be destroyed, so here is a fun hack to make up an ID bc we dont have one
	inputVars := fmt.Sprint(d.Get("template_properties").(map[string]interface{}))
	concatVars := inputVars + d.Get("template").(string)
	id := sha256.Sum256([]byte(concatVars))

	d.SetId(fmt.Sprintf("%x", id))
	d.Set("value", renderedTemplate.Value)

	log.Println("onefuse.dataSourceRenderedTemplate " + renderedTemplate.Value)

	return nil
}
