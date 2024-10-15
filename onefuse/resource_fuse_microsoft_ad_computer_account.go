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

func resourceMicrosoftADComputerAccount() *schema.Resource {
	return &schema.Resource{
		Create: resourceMicrosoftADComputerAccountCreate,
		Read:   resourceMicrosoftADComputerAccountRead,
		Update: resourceMicrosoftADComputerAccountUpdate,
		Delete: resourceMicrosoftADComputerAccountDelete,
		Importer: &schema.ResourceImporter{
			State: importADReservation,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				// Updates not yet supported for Microsoft Active Directory Computer Names.
				ForceNew: true,
				// Suppress diff if both names are the same in Lowercase or Uppercase
				DiffSuppressFunc: func(k string, oldName string, newName string, d *schema.ResourceData) bool {
					if strings.ToLower(oldName) == strings.ToLower(newName) {
						return true
					} else if strings.ToUpper(oldName) == strings.ToUpper(newName) {
						return true
					} else {
						return false
					}
				},
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
			"final_ou": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"template_properties": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func bindMicrosoftADComputerAccountResource(d *schema.ResourceData, computerAccount *MicrosoftADComputerAccount) error {
	log.Println("onefuse.bindMicrosoftADComputerAccountResource")

	if err := d.Set("name", computerAccount.Name); err != nil {
		return errors.WithMessage(err, "Cannot set name: "+computerAccount.Name)
	}
	if err := d.Set("final_ou", computerAccount.FinalOU); err != nil {
		return errors.WithMessage(err, "Cannot set final OU: "+computerAccount.FinalOU)
	}

	if err := d.Set("workspace_url", computerAccount.Links.Workspace.Href); err != nil {
		return errors.WithMessage(err, "Cannot set workspace: "+computerAccount.Links.Workspace.Href)
	}

	microsoftADPolicyURLSplit := strings.Split(computerAccount.Links.Policy.Href, "/")
	microsoftADPolicyID := microsoftADPolicyURLSplit[len(microsoftADPolicyURLSplit)-2]
	microsoftADPolicyIDInt, _ := strconv.Atoi(microsoftADPolicyID)
	if err := d.Set("policy_id", microsoftADPolicyIDInt); err != nil {
		return errors.WithMessage(err, "Cannot set policy")
	}

	return nil
}

func resourceMicrosoftADComputerAccountCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceMicrosoftADComputerAccountCreate")

	config := m.(Config)

	newComputerAccount := MicrosoftADComputerAccount{
		Name:               d.Get("name").(string),
		FinalOU:            d.Get("final_ou").(string),
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	computerAccount, err := config.NewOneFuseApiClient().CreateMicrosoftADComputerAccount(&newComputerAccount)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(computerAccount.ID))

	return bindMicrosoftADComputerAccountResource(d, computerAccount)
}

func resourceMicrosoftADComputerAccountRead(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceMicrosoftADComputerAccountRead")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	computerAccount, err := config.NewOneFuseApiClient().GetMicrosoftADComputerAccount(intID)
	if err != nil {
		return err
	}

	return bindMicrosoftADComputerAccountResource(d, computerAccount)
}

func resourceMicrosoftADComputerAccountUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceMicrosoftADComputerAccountUpdate")

	// Determine if a change is needed
	changed := (d.HasChange("name") ||
		d.HasChange("policy_id") ||
		d.HasChange("final_ou") ||
		d.HasChange("workspace_url"))

	if !changed {
		return nil
	}

	// Make the API call to update the computer account
	config := m.(Config)

	// Create the desired AD Computer Account object
	id := d.Id()
	desiredComputerAccount := MicrosoftADComputerAccount{
		Name:               d.Get("name").(string),
		FinalOU:            d.Get("final_ou").(string),
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	computerAccount, err := config.NewOneFuseApiClient().UpdateMicrosoftADComputerAccount(intID, &desiredComputerAccount)
	if err != nil {
		return err
	}

	return bindMicrosoftADComputerAccountResource(d, computerAccount)
}

func resourceMicrosoftADComputerAccountDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceMicrosoftADComputerAccountDelete")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return config.NewOneFuseApiClient().DeleteMicrosoftADComputerAccount(intID)
}

func importADReservation(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
    log.Println("onefuse.importADReservation - Starting the import")

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

    adRecord, err := config.NewOneFuseApiClient().GetMicrosoftADComputerAccount(intID)
    if err != nil {
        log.Printf("Error fetching AD reservation: %v", err)
        return nil, errors.Wrap(err, "error fetching AD reservation")
    }

    // Bind the AD reservation
    if err := bindMicrosoftADComputerAccountResource(d, adRecord); err != nil {
        log.Printf("Error binding AD reservation resource: %v", err)
        return nil, errors.Wrap(err, "failed to bind AD reservation data")
    }

    jobMetaDataRecord, err := fetchAdJobMetaData(adRecord, &config)
    if err != nil {
        log.Printf("Error fetching job metadata: %v", err)
        return nil, errors.Wrap(err, "error fetching job metadata during import")
    }

    if jobMetaDataRecord == nil {
        log.Println("jobMetaDataRecord is nil after fetching job metadata")
        return nil, errors.New("jobMetaDataRecord is nil after fetching job metadata")
    }

    log.Printf("Template Properties are: %+v", jobMetaDataRecord.ResolvedProperties)
    log.Println("onefuse.importADReservation - import completed successfully")

    return []*schema.ResourceData{d}, nil
}

func fetchAdJobMetaData(adRecord *MicrosoftADComputerAccount, config *Config) (*JobMetaData, error){
	log.Println("Fetching the job metadata - Start")

    jobMetaDataURLSplit := strings.Split(adRecord.Links.JobMetadata.Href, "/")
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
