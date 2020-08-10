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

func dataSourceMicrosoftEndpoint() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceMicrosoftEndpointRead,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceMicrosoftEndpointRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("onefuse.dataSourceMicrosoftEndpointRead")

	config := meta.(Config)
	apiClient := config.NewOneFuseApiClient()

	endpoint, err := apiClient.GetMicrosoftEndpointByName(d.Get("name").(string))

	if err != nil {
		return fmt.Errorf("Error loading Microsoft Endpoint: %s", err)
	}

	d.SetId(strconv.Itoa(endpoint.ID))
	d.Set("name", endpoint.Name)

	return nil
}
