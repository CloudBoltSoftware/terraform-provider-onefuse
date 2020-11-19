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

func dataSourceIPAMPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceIPAMPolicyRead,
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

func dataSourceIPAMPolicyRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("onefuse.dataSourceIPAMPolicyRead")

	config := meta.(Config)
	apiClient := config.NewOneFuseApiClient()

	ipamPolicy, err := apiClient.GetIPAMPolicyByName(d.Get("name").(string))

	if err != nil {
		return fmt.Errorf("Error loading IPAM Policy: %s", err)
	}

	d.SetId(strconv.Itoa(ipamPolicy.ID))
	d.Set("name", ipamPolicy.Name)
	d.Set("description", ipamPolicy.Description)

	return nil
}
