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

func dataSourceNamingPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceNamingPolicyRead,
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

func dataSourceNamingPolicyRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("onefuse.dataSourceNamingPolicyRead")

	config := meta.(Config)
	apiClient := config.NewOneFuseApiClient()

	namingPolicy, err := apiClient.GetNamingPolicyByName(d.Get("name").(string))

	if err != nil {
		return fmt.Errorf("Error loading Naming Policy: %s", err)
	}

	d.SetId(strconv.Itoa(namingPolicy.ID))
	d.Set("name", namingPolicy.Name)
	d.Set("description", namingPolicy.Description)

	return nil
}
