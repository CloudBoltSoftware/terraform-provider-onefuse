// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func resourceServicenowCMDBDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceServicenowCMDBDeploymentCreate,
		Read:   resourceServicenowCMDBDeploymentRead,
		Update: resourceServicenowCMDBDeploymentUpdate,
		Delete: resourceServicenowCMDBDeploymentDelete,
		Importer: &schema.ResourceImporter{
			State: importServiceNowCmdbDeployment,
		},
		Schema: map[string]*schema.Schema{
			"policy_id": {
				Type:     schema.TypeInt,
				Required: true,
			},
			"workspace_url": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"configuration_items_info": {
				Type: schema.TypeList,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
				Computed: true,
				Optional: true,
			},
			"execution_details": {
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

func bindServicenowCMDBDeploymentResource(d *schema.ResourceData, servicenowCMDBDeployment *ServicenowCMDBDeployment) error {
	log.Println("onefuse.bindServicenowCMDBDeploymentResource")

	if err := d.Set("workspace_url", servicenowCMDBDeployment.Links.Workspace.Href); err != nil {
		return errors.WithMessage(err, "Cannot set workspace: "+servicenowCMDBDeployment.Links.Workspace.Href)
	}

	if err := d.Set("configuration_items_info", servicenowCMDBDeployment.ConfigurationItemsInfo); err != nil {
		return errors.WithMessage(err, "Cannot set configuration_items_info")
	}
	
	executionDetailsJSON, err := json.Marshal(servicenowCMDBDeployment.ExecutionDetails)
	if err != nil {
		return errors.WithMessage(err, "Unable to Marshal execution_details into string")
	}
	executionDetailsString := string(executionDetailsJSON)

	if err := d.Set("execution_details", executionDetailsString); err != nil {
		return errors.WithMessage(err, "Cannot set execution_details")
	}

	servicenowCMDBPolicyURLSplit := strings.Split(servicenowCMDBDeployment.Links.Policy.Href, "/")
	servicenowCMDBPolicyID := servicenowCMDBPolicyURLSplit[len(servicenowCMDBPolicyURLSplit)-2]
	servicenowCMDBPolicyIDInt, _ := strconv.Atoi(servicenowCMDBPolicyID)
	if err := d.Set("policy_id", servicenowCMDBPolicyIDInt); err != nil {
		return errors.WithMessage(err, "Cannot set policy")
	}

	return nil
}

func resourceServicenowCMDBDeploymentCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceServicenowCMDBDeploymentCreate")

	config := m.(Config)

	newServicenowCMDBDeployment := ServicenowCMDBDeployment{
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	servicenowCMDBDeployment, err := config.NewOneFuseApiClient().CreateServicenowCMDBDeployment(&newServicenowCMDBDeployment)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(servicenowCMDBDeployment.ID))

	return bindServicenowCMDBDeploymentResource(d, servicenowCMDBDeployment)
}

func resourceServicenowCMDBDeploymentRead(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceServicenowCMDBDeploymentRead")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	servicenowCMDBDeployment, err := config.NewOneFuseApiClient().GetServicenowCMDBDeployment(intID)
	if err != nil {
		return err
	}

	return bindServicenowCMDBDeploymentResource(d, servicenowCMDBDeployment)
}

func resourceServicenowCMDBDeploymentUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceServicenowCMDBDeploymentUpdate")

	// Determine if a change is needed
	changed := (d.HasChange("policy_id") ||
		d.HasChange("workspace_url") ||
		d.HasChange("template_properties"))

	if !changed {
		return nil
	}

	// Get the config
	config := m.(Config)

	// Create the desired ServiceNow CMDB Deployment
	id := d.Id()
	desiredServicenowCMDBDeployment := ServicenowCMDBDeployment{
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	servicenowCMDBDeployment, err := config.NewOneFuseApiClient().UpdateServicenowCMDBDeployment(intID, &desiredServicenowCMDBDeployment)
	if err != nil {
		return err
	}

	return bindServicenowCMDBDeploymentResource(d, servicenowCMDBDeployment)
}

func resourceServicenowCMDBDeploymentDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceServicenowCMDBDeploymentDelete")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return config.NewOneFuseApiClient().DeleteServicenowCMDBDeployment(intID)
}

func importServiceNowCmdbDeployment(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
    log.Println("onefuse.importServiceNowCmdbDeployment - Starting the import")

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

    snowRecord, err := config.NewOneFuseApiClient().GetServicenowCMDBDeployment(intID)
    if err != nil {
        log.Printf("Error fetching ServiceNow reservation: %v", err)
        return nil, errors.Wrap(err, "error fetching ServiceNow reservation")
    }

    // Bind the ServiceNow Cmdb reservation record
    if err := bindServicenowCMDBDeploymentResource(d, snowRecord); err != nil {
        log.Printf("Error binding ServiceNow reservation resource: %v", err)
        return nil, errors.Wrap(err, "failed to bind ServiceNow reservation data")
    }

    jobMetaDataRecord, err := fetchSnowCmdbJobMetaData(snowRecord, &config)
    if err != nil {
        log.Printf("Error fetching job metadata: %v", err)
        return nil, errors.Wrap(err, "error fetching job metadata during import")
    }

    if jobMetaDataRecord == nil {
        log.Println("jobMetaDataRecord is nil after fetching job metadata")
        return nil, errors.New("jobMetaDataRecord is nil after fetching job metadata")
    }

	log.Printf("Template Properties are: %+v", jobMetaDataRecord.ResolvedProperties)
    log.Println("onefuse.importServiceNowCmdbDeployment - import completed successfully")

    return []*schema.ResourceData{d}, nil
}

func fetchSnowCmdbJobMetaData(snowRecord *ServicenowCMDBDeployment, config *Config) (*JobMetaData, error){
	log.Println("Fetching the job metadata - Start")

    jobMetaDataURLSplit := strings.Split(snowRecord.Links.JobMetadata.Href, "/")
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
