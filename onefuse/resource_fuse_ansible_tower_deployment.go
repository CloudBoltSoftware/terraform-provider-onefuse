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

func resourceAnsibleTowerDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceAnsibleTowerDeploymentCreate,
		Read:   resourceAnsibleTowerDeploymentRead,
		Update: resourceAnsibleTowerDeploymentUpdate,
		Delete: resourceAnsibleTowerDeploymentDelete,
		Importer: &schema.ResourceImporter{
			State: importAnsibleReservation,
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
			"limit": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"hosts": {
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
			"inventory_name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"provisioning_job_results": {
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

func bindAnsibleTowerDeploymentResource(d *schema.ResourceData, ansibleDeployment *AnsibleTowerDeployment) error {
	log.Println("onefuse.bindAnsibleTowerDeploymentResource")

	if err := d.Set("workspace_url", ansibleDeployment.Links.Workspace.Href); err != nil {
		return errors.WithMessage(err, "Cannot set workspace: "+ansibleDeployment.Links.Workspace.Href)
	}

	if err := d.Set("hosts", ansibleDeployment.Hosts); err != nil {
		hosts := strings.Join(ansibleDeployment.Hosts[:], ",")
		return errors.WithMessage(err, "Cannot set hosts: "+hosts)
	}

	if err := d.Set("limit", ansibleDeployment.Limit); err != nil {
		return errors.WithMessage(err, "Cannot set limit: "+ansibleDeployment.Limit)
	}

	if err := d.Set("inventory_name", ansibleDeployment.InventoryName); err != nil {
		return errors.WithMessage(err, "Cannot set inventory name: "+ansibleDeployment.InventoryName)
	}

	provisioningJobResultsJson, err := json.Marshal(ansibleDeployment.ProvisioningJobResults)
	if err != nil {
		return errors.WithMessage(err, "Unable to Marshal provisioning_job_results into string")
	}
	provisioningJobResultsString := string(provisioningJobResultsJson)
	if err := d.Set("provisioning_job_results", provisioningJobResultsString); err != nil {
		return errors.WithMessage(err, "Cannot set provisioning_job_results: "+provisioningJobResultsString)
	}

	ansibleTowerPolicyURLSplit := strings.Split(ansibleDeployment.Links.Policy.Href, "/")
	ansibleTowerPolicyID := ansibleTowerPolicyURLSplit[len(ansibleTowerPolicyURLSplit)-2]
	ansibleTowerPolicyIDInt, _ := strconv.Atoi(ansibleTowerPolicyID)
	if err := d.Set("policy_id", ansibleTowerPolicyIDInt); err != nil {
		return errors.WithMessage(err, "Cannot set policy")
	}

	return nil
}

func resourceAnsibleTowerDeploymentCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceAnsibleTowerDeploymentCreate")

	var hosts []string
	for _, group := range d.Get("hosts").([]interface{}) {
		hosts = append(hosts, group.(string))
	}

	config := m.(Config)

	newAnsibleTowerDeployment := AnsibleTowerDeployment{
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		Hosts:              hosts,
		Limit:              d.Get("limit").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	ansibleDeployment, err := config.NewOneFuseApiClient().CreateAnsibleTowerDeployment(&newAnsibleTowerDeployment)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(ansibleDeployment.ID))

	return bindAnsibleTowerDeploymentResource(d, ansibleDeployment)
}

func resourceAnsibleTowerDeploymentRead(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceAnsibleTowerDeploymentRead")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	ansibleDeployment, err := config.NewOneFuseApiClient().GetAnsibleTowerDeployment(intID)
	if err != nil {
		return err
	}

	return bindAnsibleTowerDeploymentResource(d, ansibleDeployment)
}

func resourceAnsibleTowerDeploymentUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceAnsibleTowerDeploymentUpdate")
	log.Println("No Op!")
	return nil
}

func resourceAnsibleTowerDeploymentDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceAnsibleTowerDeploymentDelete")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return config.NewOneFuseApiClient().DeleteAnsibleTowerDeployment(intID)
}

func importAnsibleReservation(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
    log.Println("onefuse.importAnsibleReservation - Starting the import")

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

    ansibleRecord, err := config.NewOneFuseApiClient().GetAnsibleTowerDeployment(intID)
    if err != nil {
        log.Printf("Error fetching Ansible reservation: %v", err)
        return nil, errors.Wrap(err, "error fetching Ansible reservation")
    }

    // Bind the Ansible reservatin
    if err := bindAnsibleTowerDeploymentResource(d, ansibleRecord); err != nil {
        log.Printf("Error binding Ansible reservation resource: %v", err)
        return nil, errors.Wrap(err, "failed to bind Ansible reservation data")
    }

    jobMetaDataRecord, err := fetchAnsibleJobMetaData(ansibleRecord, &config)
    if err != nil {
        log.Printf("Error fetching job metadata: %v", err)
        return nil, errors.Wrap(err, "error fetching job metadata during import")
    }

    if jobMetaDataRecord == nil {
        log.Println("jobMetaDataRecord is nil after fetching job metadata")
        return nil, errors.New("jobMetaDataRecord is nil after fetching job metadata")
    }

    log.Printf("Template Properties are: %+v", jobMetaDataRecord.ResolvedProperties)
    log.Println("onefuse.importAnsibleReservation - import completed successfully")

    return []*schema.ResourceData{d}, nil
}

func fetchAnsibleJobMetaData(ansibleRecord *AnsibleTowerDeployment, config *Config) (*JobMetaData, error){
	log.Println("Fetching the job metadata - Started")

    jobMetaDataURLSplit := strings.Split(ansibleRecord.Links.JobMetadata.Href, "/")
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
