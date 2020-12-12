// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func resourceCustomNaming() *schema.Resource {
	return &schema.Resource{
		Create: resourceCustomNameCreate,
		Read:   resourceCustomNameRead,
		Update: resourceCustomNameUpdate,
		Delete: resourceCustomNameDelete,
		Schema: map[string]*schema.Schema{
			"custom_name_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"naming_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Fuse Template Properties",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func bindCustomNamingResource(d *schema.ResourceData, cn *CustomName) error {
	log.Println("onefuse.bindCustomNamingResource")

	// setting the ID is REALLY necessary here
	// we use the FQDN instead of the numeric ID as it is more likely to remain consistent as a composite key in TF
	d.SetId(cn.Name + "." + cn.DnsSuffix)

	if err := d.Set("custom_name_id", cn.Id); err != nil {
		return errors.WithMessage(err, "cannot set custom_name_id")
	}
	if err := d.Set("name", cn.Name); err != nil {
		return errors.WithMessage(err, "cannot set name")
	}
	if err := d.Set("dns_suffix", cn.DnsSuffix); err != nil {
		return errors.WithMessage(err, "cannot set dns_suffix")
	}
	return nil
}

func resourceCustomNameCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceCustomNameCreate")

	config := m.(Config)

	namingPolicyID := d.Get("naming_policy_id").(string)
	workspaceID := d.Get("workspace_id").(string)
	templateProperties := d.Get("template_properties").(map[string]interface{})

	cn, err := config.NewOneFuseApiClient().GenerateCustomName(namingPolicyID, workspaceID, templateProperties)
	if err != nil {
		return err
	}

	return bindCustomNamingResource(d, cn)
}

func resourceCustomNameRead(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceCustomNameRead")

	config := m.(Config)

	id := d.Get("custom_name_id").(int)

	customName, err := config.NewOneFuseApiClient().GetCustomName(id)
	if err != nil {
		return err
	}

	return bindCustomNamingResource(d, customName)
}

func resourceCustomNameUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceCustomNameUpdate")
	return nil
}

func resourceCustomNameDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceCustomNameDelete")

	config := m.(Config)

	id := d.Get("custom_name_id").(int)

	return config.NewOneFuseApiClient().DeleteCustomName(id)
}
