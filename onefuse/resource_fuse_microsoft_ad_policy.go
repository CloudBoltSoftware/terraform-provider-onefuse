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
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func resourceMicrosoftADPolicy() *schema.Resource {
	return &schema.Resource{
		Create: resourceMicrosoftADPolicyCreate,
		Read:   resourceMicrosoftADPolicyRead,
		Update: resourceMicrosoftADPolicyUpdate,
		Delete: resourceMicrosoftADPolicyDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"microsoft_endpoint_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"computer_name_letter_case": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Either Lowercase or Uppercase",
			},
			"ou": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace_url": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"security_groups": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"create_ou": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
			"remove_ou": {
				Type:     schema.TypeBool,
				Computed: true,
				Optional: true,
			},
		},
	}
}

func bindMicrosoftADPolicyResource(d *schema.ResourceData, policy *MicrosoftADPolicy) error {
	log.Println("onefuse.bindMicrosoftADPolicyResource")

	if err := d.Set("name", policy.Name); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("Cannot set name: '%s'", policy.Name))
	}

	if err := d.Set("description", policy.Description); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("Cannot set description: '%s'", policy.Description))
	}

	if err := d.Set("workspace_url", policy.Links.Workspace.Href); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("Cannot set workspace: '%s'", policy.Links.Workspace.Href))
	}

	if err := d.Set("computer_name_letter_case", policy.ComputerNameLetterCase); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("Cannot set computer_name_letter_case: '%s'", policy.ComputerNameLetterCase))
	}

	if err := d.Set("ou", policy.OU); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("Cannot set OU: '%s'", policy.OU))
	}

	if err := d.Set("create_ou", policy.CreateOU); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("Cannot set Create OU: %t", policy.CreateOU))
	}

	if err := d.Set("remove_ou", policy.RemoveOU); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("Cannot set Remove OU: %t", policy.CreateOU))
	}

	if err := d.Set("security_groups", policy.SecurityGroups); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("Cannot set Security Groups: %#v", policy.SecurityGroups))
	}

	microsoftEndpointURLSplit := strings.Split(policy.Links.MicrosoftEndpoint.Href, "/")
	microsoftEndpointID := microsoftEndpointURLSplit[len(microsoftEndpointURLSplit)-2]
	microsoftEndpointIDInt, err := strconv.Atoi(microsoftEndpointID)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("Expected to convert '%s' to int value.", microsoftEndpointID))
	}
	if err := d.Set("microsoft_endpoint_id", microsoftEndpointIDInt); err != nil {
		return errors.WithMessage(err, "Cannot set microsoft_endpoint_id")
	}

	return nil
}

func resourceMicrosoftADPolicyCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceMicrosoftADPolicyCreate")

	config := m.(Config)

	var securityGroups []string
	for _, group := range d.Get("security_groups").([]interface{}) {
		securityGroups = append(securityGroups, group.(string))
	}

	newPolicy := MicrosoftADPolicy{
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		OU:                     d.Get("ou").(string),
		MicrosoftEndpointID:    d.Get("microsoft_endpoint_id").(int),
		ComputerNameLetterCase: d.Get("computer_name_letter_case").(string),
		WorkspaceURL:           d.Get("workspace_url").(string),
		CreateOU:               d.Get("create_ou").(bool),
		RemoveOU:               d.Get("remove_ou").(bool),
		SecurityGroups:         securityGroups,
	}

	policy, err := config.NewOneFuseApiClient().CreateMicrosoftADPolicy(&newPolicy)
	if err != nil {
		return errors.WithMessage(err, "Failed to create Microsoft AD Policy")
	}
	d.SetId(strconv.Itoa(policy.ID))

	return resourceMicrosoftADPolicyRead(d, m)
}

func resourceMicrosoftADPolicyRead(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceMicrosoftADPolicyRead")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return errors.WithMessage(err, "Failed to convert integer ID to string: "+string(id))
	}

	policy, err := config.NewOneFuseApiClient().GetMicrosoftADPolicy(intID)
	if err != nil {
		return errors.WithMessage(err, "Failed to read Microsoft AD Policy")
	}

	return bindMicrosoftADPolicyResource(d, policy)
}

func resourceMicrosoftADPolicyUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceMicrosoftADPolicyUpdate")

	// Determine if a change is needed
	changed := (d.HasChange("name") ||
		d.HasChange("description") ||
		d.HasChange("microsoft_endpoint_id") ||
		d.HasChange("computer_name_letter_case") ||
		d.HasChange("workspace_url") ||
		d.HasChange("ou") ||
		d.HasChange("create_ou") ||
		d.HasChange("remove_ou") ||
		d.HasChange("security_groups"))

	if !changed {
		return nil
	}

	// Make the API call to update the policy
	config := m.(Config)

	// Create the desired AD Policy object
	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	var securityGroups []string
	for _, group := range d.Get("security_groups").([]interface{}) {
		securityGroups = append(securityGroups, group.(string))
	}

	desiredPolicy := MicrosoftADPolicy{
		Name:                   d.Get("name").(string),
		Description:            d.Get("description").(string),
		MicrosoftEndpointID:    d.Get("microsoft_endpoint_id").(int),
		ComputerNameLetterCase: d.Get("computer_name_letter_case").(string),
		WorkspaceURL:           d.Get("workspace_url").(string),
		OU:                     d.Get("ou").(string),
		CreateOU:               d.Get("create_ou").(bool),
		RemoveOU:               d.Get("remove_ou").(bool),
		SecurityGroups:         securityGroups,
	}

	_, err = config.NewOneFuseApiClient().UpdateMicrosoftADPolicy(intID, &desiredPolicy)
	if err != nil {
		return errors.WithMessage(err, "Failed to updated Microsoft AD Policy")
	}

	return resourceMicrosoftADPolicyRead(d, m)
}

func resourceMicrosoftADPolicyDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceMicrosoftADPolicyDelete")

	config := m.(Config)

	id := d.Id()
	inID, err := strconv.Atoi(id)
	if err != nil {
		return errors.WithMessage(err, "Failed to delete Microsoft AD Policy")
	}

	return config.NewOneFuseApiClient().DeleteMicrosoftADPolicy(inID)
}
