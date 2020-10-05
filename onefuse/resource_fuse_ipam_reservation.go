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

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func resourceIPAMReservation() *schema.Resource {
	return &schema.Resource{
		Create: resourceIPAMReservationCreate,
		Read:   resourceIPAMReservationRead,
		Update: resourceIPAMReservationUpdate,
		Delete: resourceIPAMReservationDelete,
		Schema: map[string]*schema.Schema{
			"hostname": {
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
			"ip_address": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"gateway": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"primary_dns": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"secondary_dns": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"dns_suffix": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"dns_search_suffix": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				Optional: true,
			},
			"template_properties": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func bindIPAMReservationResource(d *schema.ResourceData, ipamRecord *IPAMReservation) error {
	log.Println("onefuse.bindIPAMReservationResource")

	if err := d.Set("hostname", ipamRecord.Hostname); err != nil {
		return errors.WithMessage(err, "Cannot set name: "+ipamRecord.Hostname)
	}

	if err := d.Set("workspace_url", ipamRecord.Links.Workspace.Href); err != nil {
		return errors.WithMessage(err, "Cannot set workspace: "+ipamRecord.Links.Workspace.Href)
	}

	if err := d.Set("ip_address", ipamRecord.IPaddress); err != nil {
		return errors.WithMessage(err, "Cannot set IP address: "+ipamRecord.IPaddress)
	}

	if err := d.Set("primary_dns", ipamRecord.PrimaryDNS); err != nil {
		return errors.WithMessage(err, "Cannot set Primmary DNS: "+ipamRecord.PrimaryDNS)
	}

	if err := d.Set("secondary_dns", ipamRecord.SecondaryDNS); err != nil {
		return errors.WithMessage(err, "Cannot set Secondary DNS: "+ipamRecord.SecondaryDNS)
	}

	if err := d.Set("gateway", ipamRecord.Gateway); err != nil {
		return errors.WithMessage(err, "Cannot set Gateway: "+ipamRecord.Gateway)
	}

	ipamPolicyURLSplit := strings.Split(ipamRecord.Links.Policy.Href, "/")
	ipamPolicyID := ipamPolicyURLSplit[len(ipamPolicyURLSplit)-2]
	ipamPolicyIDInt, _ := strconv.Atoi(ipamPolicyID)
	if err := d.Set("policy_id", ipamPolicyIDInt); err != nil {
		return errors.WithMessage(err, "Cannot set policy")
	}

	return nil
}

func resourceIPAMReservationCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceIPAMReservationCreate")

	var ipam_Suffixes []string
	for _, group := range d.Get("dns_search_suffix").([]interface{}) {
		ipam_Suffixes = append(ipam_Suffixes, group.(string))
	}

	config := m.(Config)

	newIPAMRecord := IPAMReservation{
		Hostname:           d.Get("hostname").(string),
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		IPaddress:          d.Get("ip_address").(string),
		Gateway:            d.Get("gateway").(string),
		PrimaryDNS:         d.Get("primary_dns").(string),
		SecondaryDNS:       d.Get("secondary_dns").(string),
		DNSSuffix:          d.Get("dns_suffix").(string),
		DNSSearchSuffixes:  ipam_Suffixes,
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	ipamRecord, err := config.NewOneFuseApiClient().CreateIPAMReservation(&newIPAMRecord)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(ipamRecord.ID))

	return bindIPAMReservationResource(d, ipamRecord)
}

func resourceIPAMReservationRead(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceIPAMReservationRead")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	ipamRecord, err := config.NewOneFuseApiClient().GetIPAMReservation(intID)
	if err != nil {
		return err
	}

	return bindIPAMReservationResource(d, ipamRecord)
}

func resourceIPAMReservationUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceIPAMReservationUpdate")

	// Determine if a change is needed
	changed := (d.HasChange("hostname") ||
		d.HasChange("policy_id") ||
		d.HasChange("workspace_url")) ||
		d.HasChange("ip_address") ||
		d.HasChange("primary_dns") ||
		d.HasChange("secondary_dns") ||
		d.HasChange("dns_suffix") ||
		d.HasChange("template_properties")

	if !changed {
		return nil
	}

	var ipam_Suffixes []string
	for _, group := range d.Get("dns_search_suffix").([]interface{}) {
		ipam_Suffixes = append(ipam_Suffixes, group.(string))
	}

	// Make the API call to update the computer account
	config := m.(Config)

	// Create the desired IPAM Reservation
	id := d.Id()
	desiredIPAMRecord := IPAMReservation{
		Hostname:           d.Get("hostname").(string),
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		IPaddress:          d.Get("ip_address").(string),
		PrimaryDNS:         d.Get("primary_dns").(string),
		SecondaryDNS:       d.Get("secondary_dns").(string),
		DNSSuffix:          d.Get("dns_suffix").(string),
		DNSSearchSuffixes:  ipam_Suffixes,
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	ipamRecord, err := config.NewOneFuseApiClient().UpdateIPAMReservation(intID, &desiredIPAMRecord)
	if err != nil {
		return err
	}

	return bindIPAMReservationResource(d, ipamRecord)
}

func resourceIPAMReservationDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceIPAMReservationDelete")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return config.NewOneFuseApiClient().DeleteIPAMReservation(intID)
}
