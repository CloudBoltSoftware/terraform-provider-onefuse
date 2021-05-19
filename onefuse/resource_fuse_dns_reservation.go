// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func resourceDNSReservation() *schema.Resource {
	return &schema.Resource{
		Create: resourceDNSReservationCreate,
		Read:   resourceDNSReservationRead,
		Update: resourceDNSReservationUpdate,
		Delete: resourceDNSReservationDelete,
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"workspace_url": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"value": {
				Type:     schema.TypeString,
				Required: true,
			},
			"zones": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Required: true,
			},
			"template_properties": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
			Delete: schema.DefaultTimeout(5 * time.Minute),
		},
	}
}

func bindDNSReservationResource(d *schema.ResourceData, dnsRecord *DNSReservation) error {
	log.Println("onefuse.bindDNSReservationResource")

	if err := d.Set("name", dnsRecord.Name); err != nil {
		return errors.WithMessage(err, "Cannot set name: "+dnsRecord.Name)
	}

	if err := d.Set("workspace_url", dnsRecord.Links.Workspace.Href); err != nil {
		return errors.WithMessage(err, "Cannot set workspace: "+dnsRecord.Links.Workspace.Href)
	}

	dnsPolicyURLSplit := strings.Split(dnsRecord.Links.Policy.Href, "/")
	dnsPolicyID := dnsPolicyURLSplit[len(dnsPolicyURLSplit)-2]
	dnsPolicyIDInt, _ := strconv.Atoi(dnsPolicyID)
	if err := d.Set("policy_id", dnsPolicyIDInt); err != nil {
		return errors.WithMessage(err, "Cannot set policy")
	}

	return nil
}

func resourceDNSReservationCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceDNSReservationCreate")

	var dnsZones []string
	for _, group := range d.Get("zones").([]interface{}) {
		dnsZones = append(dnsZones, group.(string))
	}

	config := m.(Config)

	newDNSRecord := DNSReservation{
		Name:               d.Get("name").(string),
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		Value:              d.Get("value").(string),
		Zones:              dnsZones,
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	dnsRecord, err := config.NewOneFuseApiClient().CreateDNSReservation(&newDNSRecord)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(dnsRecord.ID))

	return bindDNSReservationResource(d, dnsRecord)
}

func resourceDNSReservationRead(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceDNSReservationRead")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	dnsRecord, err := config.NewOneFuseApiClient().GetDNSReservation(intID)
	if err != nil {
		return err
	}

	return bindDNSReservationResource(d, dnsRecord)
}

func resourceDNSReservationUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceDNSReservationUpdate")

	// Determine if a change is needed
	changed := (d.HasChange("name") ||
		d.HasChange("policy_id") ||
		d.HasChange("workspace_url")) ||
		d.HasChange("value") ||
		d.HasChange("zones") ||
		d.HasChange("template_properties")

	if !changed {
		return nil
	}

	var dnsZones []string
	for _, group := range d.Get("zones").([]interface{}) {
		dnsZones = append(dnsZones, group.(string))
	}

	// Make the API call to update the computer account
	config := m.(Config)

	// Create the desired AD Computer Account object
	id := d.Id()
	desiredDNSRecord := DNSReservation{
		Name:               d.Get("name").(string),
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		Value:              d.Get("value").(string),
		Zones:              dnsZones,
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	dnsRecord, err := config.NewOneFuseApiClient().UpdateDNSReservation(intID, &desiredDNSRecord)
	if err != nil {
		return err
	}

	return bindDNSReservationResource(d, dnsRecord)
}

func resourceDNSReservationDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceDNSReservationDelete")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return config.NewOneFuseApiClient().DeleteDNSReservation(intID)
}
