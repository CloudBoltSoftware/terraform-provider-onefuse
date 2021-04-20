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

func dataSourceServicenowCMDBPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceServicenowCMDBPolicyRead,
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

func dataSourceServicenowCMDBPolicyRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("onefuse.dataSourceServicenowCMDBPolicyRead")

	config := meta.(Config)
	apiClient := config.NewOneFuseApiClient()

	servicenowCMDBPolicy, err := apiClient.GetServicenowCMDBPolicyByName(d.Get("name").(string))

	if err != nil {
		return fmt.Errorf("Error loading ServicenowCMDB Policy: %s", err)
	}

	d.SetId(strconv.Itoa(servicenowCMDBPolicy.ID))
	d.Set("name", servicenowCMDBPolicy.Name)
	d.Set("description", servicenowCMDBPolicy.Description)

	return nil
}
