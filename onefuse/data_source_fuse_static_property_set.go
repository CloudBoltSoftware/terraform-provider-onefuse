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

func dataSourceStaticPropertySet() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceStaticPropertySetRead,
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
			"properties": {
				Type:     schema.TypeMap,
				Computed: true,
			},
			"raw": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceStaticPropertySetRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("onefuse.dataSourceStaticPropertySetRead")

	config := meta.(Config)
	apiClient := config.NewOneFuseApiClient()

	staticPropertySet, err := apiClient.GetStaticPropertySetByName(d.Get("name").(string))

	if err != nil {
		return fmt.Errorf("Error loading Static Property Set: %s", err)
	}

	d.SetId(strconv.Itoa(staticPropertySet.ID))
	d.Set("name", staticPropertySet.Name)
	d.Set("properties", staticPropertySet.Properties)
	d.Set("raw", staticPropertySet.Raw)

	return nil
}
