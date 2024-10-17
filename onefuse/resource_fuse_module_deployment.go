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

func resourceModuleDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceModuleDeploymentCreate,
		Read:   resourceModuleDeploymentRead,
		Update: resourceModuleDeploymentUpdate,
		Delete: resourceModuleDeploymentDelete,
		Importer: &schema.ResourceImporter{
			State: importModuleDeployment,
		},
		Schema: map[string]*schema.Schema{
			"name": {
				Type:     schema.TypeString,
				Required: false,
				Computed: true,
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
			"template_properties": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"provisioning_job_results": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"deprovisioning_job_results": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func bindModuleDeploymentResource(d *schema.ResourceData, ModuleDeployment *ModuleDeployment) error {
	log.Println("onefuse.bindModuleDeploymentResource")

	if err := d.Set("workspace_url", ModuleDeployment.Links.Workspace.Href); err != nil {
		return errors.WithMessage(err, "Cannot set workspace: "+ModuleDeployment.Links.Workspace.Href)
	}

	if err := d.Set("name", ModuleDeployment.Name); err != nil {
		return errors.WithMessage(err, "Cannot set name: "+ModuleDeployment.Name)
	}

	provisioningJobResultsJson, err := json.Marshal(ModuleDeployment.ProvisioningJobResults)
	if err != nil {
		return errors.WithMessage(err, "Unable to Marshal provisioning_job_results into string")
	}

	provisioningJobResultsString := string(provisioningJobResultsJson)
	if err := d.Set("provisioning_job_results", provisioningJobResultsString); err != nil {
		return errors.WithMessage(err, "Cannot set provisioning_job_results: "+provisioningJobResultsString)
	}

	deprovisioningJobResultsJson, err := json.Marshal(ModuleDeployment.DeprovisioningJobResults)
	if err != nil {
		return errors.WithMessage(err, "Unable to Marshal deprovisioning_job_results into string")
	}

	deprovisioningJobResultsString := string(deprovisioningJobResultsJson)
	if err := d.Set("deprovisioning_job_results", deprovisioningJobResultsString); err != nil {
		return errors.WithMessage(err, "Cannot set deprovisioning_job_results: "+deprovisioningJobResultsString)
	}

	ModulePolicyURLSplit := strings.Split(ModuleDeployment.Links.Policy.Href, "/")
	ModulePolicyID := ModulePolicyURLSplit[len(ModulePolicyURLSplit)-2]
	ModulePolicyIDInt, _ := strconv.Atoi(ModulePolicyID)
	if err := d.Set("policy_id", ModulePolicyIDInt); err != nil {
		return errors.WithMessage(err, "Cannot set policy")
	}

	return nil
}

func resourceModuleDeploymentCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceModuleDeploymentCreate")

	config := m.(Config)

	newModuleDeployment := ModuleDeployment{
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	ModuleDeployment, err := config.NewOneFuseApiClient().CreateModuleDeployment(&newModuleDeployment)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(ModuleDeployment.ID))

	return bindModuleDeploymentResource(d, ModuleDeployment)
}

func resourceModuleDeploymentRead(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceModuleDeploymentRead")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	ModuleDeployment, err := config.NewOneFuseApiClient().GetModuleDeployment(intID)
	if err != nil {
		return err
	}

	return bindModuleDeploymentResource(d, ModuleDeployment)
}

func resourceModuleDeploymentUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceModuleDeploymentUpdate")

	// Determine if a change is needed
	changed := (d.HasChange("policy_id") ||
		d.HasChange("workspace_url") ||
		d.HasChange("template_properties"))

	if !changed {
		return nil
	}

	// Make the API call to update the computer account
	config := m.(Config)

	// Create the desired Module Deployment
	id := d.Id()
	desiredModuleDeployment := ModuleDeployment{
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	ModuleDeployment, err := config.NewOneFuseApiClient().UpdateModuleDeployment(intID, &desiredModuleDeployment)
	if err != nil {
		return err
	}

	return bindModuleDeploymentResource(d, ModuleDeployment)
}

func resourceModuleDeploymentDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceModuleDeploymentDelete")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return config.NewOneFuseApiClient().DeleteModuleDeployment(intID)
}

func importModuleDeployment(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
    log.Println("onefuse.importModuleDeployment - Starting the import")

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

    moduleRecord, err := config.NewOneFuseApiClient().GetModuleDeployment(intID)
    if err != nil {
        log.Printf("Error fetching module reservation: %v", err)
        return nil, errors.Wrap(err, "error fetching module reservation")
    }

    // Bind the Pluggable Module reservation
    if err := bindModuleDeploymentResource(d, moduleRecord); err != nil {
        log.Printf("Error binding module reservation resource: %v", err)
        return nil, errors.Wrap(err, "failed to bind module reservation data")
    }

    jobMetaDataRecord, err := fetchModuleJobMetaData(moduleRecord, &config)
    if err != nil {
        log.Printf("Error fetching job metadata: %v", err)
        return nil, errors.Wrap(err, "error fetching job metadata during import")
    }

    if jobMetaDataRecord == nil {
        log.Println("jobMetaDataRecord is nil after fetching job metadata")
        return nil, errors.New("jobMetaDataRecord is nil after fetching job metadata")
    }

	log.Printf("Template Properties are: %+v", jobMetaDataRecord.ResolvedProperties)
    log.Println("onefuse.importModuleDeployment - import completed successfully")

    return []*schema.ResourceData{d}, nil
}

func fetchModuleJobMetaData(moduleRecord *ModuleDeployment, config *Config) (*JobMetaData, error){
	log.Println("Fetching the job metadata - Start")

    jobMetaDataURLSplit := strings.Split(moduleRecord.Links.JobMetadata.Href, "/")
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
