// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"fmt"
	"log"
	"strconv"

	"github.com/hashicorp/terraform/helper/schema"
)

func dataSourceScriptingPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceScriptingPolicyRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceScriptingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("onefuse.dataSourceScriptingPolicyRead")

	config := meta.(Config)
	apiClient := config.NewOneFuseApiClient()

	scriptingPolicy, err := apiClient.GetScriptingPolicyByName(d.Get("name").(string))

	if err != nil {
		return fmt.Errorf("Error loading Scripting Policy: %s", err)
	}

	d.SetId(strconv.Itoa(scriptingPolicy.ID))
	d.Set("name", scriptingPolicy.Name)
	d.Set("description", scriptingPolicy.Description)

	return nil
}
