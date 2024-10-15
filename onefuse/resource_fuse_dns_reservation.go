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
		Importer: &schema.ResourceImporter{
			State: importDNSReservation,
	    },
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
			"records": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"type": {
							Type:     schema.TypeString,
							Required: true,
						},
						"name": {
							Type:     schema.TypeString,
							Required: true,
						},
						"value": {
							Type:     schema.TypeString,
							Required: true,
						},
					},
				},
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

func importDNSReservation(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
    log.Println("onefuse.importDNSReservation - Starting the import")

    config, ok := meta.(Config)
    if !ok {
        return nil, errors.New("invalid meta type")
    }

    id := d.Id()
    intID, err := strconv.Atoi(id)
    if err != nil {
        log.Printf("Error converting ID to int: %v", err)
        return nil, errors.Wrap(err, "invalid ID format")
    }

    dnsRecord, err := config.NewOneFuseApiClient().GetDNSReservation(intID)
    if err != nil {
        log.Printf("Error fetching IPAM reservation: %v", err)
        return nil, errors.Wrap(err, "error fetching DNS reservation")
    }

    // Bind the DNS reservationn
    if err := bindDNSReservationResource(d, dnsRecord); err != nil {
        log.Printf("Error binding IPAM reservation resource: %v", err)
        return nil, errors.Wrap(err, "failed to bind IPAM reservation data")
    }

	if len(dnsRecord.Records) > 0 {
		if err := d.Set("value", dnsRecord.Records[0]["value"]); err != nil {
			return nil, errors.Wrap(err, "Cannot set the dnsRecord value")
		}
		zone := strings.Join(strings.Split(dnsRecord.Records[0]["name"], ".")[1:], ".")
		zoneSlice := []string{zone}

		if err := d.Set("zones", zoneSlice); err != nil {
			return nil, errors.Wrap(err, "Cannot set the dnsRecord zones")
		}
	} else {
		return nil, errors.New("dnsRecord.Records is empty")
	}

    jobMetaDataRecord, err := fetchDnsJobMetaData(dnsRecord, &config)
    if err != nil {
        log.Printf("Error fetching job metadata: %v", err)
        return nil, errors.Wrap(err, "error fetching job metadata during import")
    }

    if jobMetaDataRecord == nil {
        log.Println("jobMetaDataRecord is nil after fetching job metadata")
        return nil, errors.New("jobMetaDataRecord is nil after fetching job metadata")
    }

    log.Printf("Template Properties are: %+v", jobMetaDataRecord.ResolvedProperties)
    log.Println("onefuse.importDNSReservation - import completed successfully")

    return []*schema.ResourceData{d}, nil
}

func fetchDnsJobMetaData(dnsRecord *DNSReservation, config *Config) (*JobMetaData, error){
	log.Println("Fetching the job metadata - Start")

    jobMetaDataURLSplit := strings.Split(dnsRecord.Links.JobMetadata.Href, "/")
    jobMetaDataId := jobMetaDataURLSplit[len(jobMetaDataURLSplit)-2]
    jobMetaDataIdInt, err := strconv.Atoi(jobMetaDataId)
    if err != nil {
        return nil, errors.Wrap(err, "failed to convert job metadata ID to int")
    }

    jobMetaDataRecord, err := GetJobMetaData(jobMetaDataIdInt, config)
    if err != nil {
        return nil, errors.Wrap(err, "failed to fetch job metadata")
    }

	log.Println("Fetching the job metadata - Completed")

    return jobMetaDataRecord, nil
}
