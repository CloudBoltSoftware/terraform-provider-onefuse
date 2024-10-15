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

func resourceScriptingDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceScriptingDeploymentCreate,
		Read:   resourceScriptingDeploymentRead,
		Update: resourceScriptingDeploymentUpdate,
		Delete: resourceScriptingDeploymentDelete,
		Importer: &schema.ResourceImporter{
			State: importScriptingReservation,
		},
		Schema: map[string]*schema.Schema{
			"hostname": {
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
			"provisioning_details": {
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

func bindScriptingDeploymentResource(d *schema.ResourceData, scriptingDeployment *ScriptingDeployment) error {
	log.Println("onefuse.bindScriptingDeploymentResource")

	if err := d.Set("workspace_url", scriptingDeployment.Links.Workspace.Href); err != nil {
		return errors.WithMessage(err, "Cannot set workspace: "+scriptingDeployment.Links.Workspace.Href)
	}

	if err := d.Set("hostname", scriptingDeployment.Hostname); err != nil {
		return errors.WithMessage(err, "Cannot set hostname: "+scriptingDeployment.Hostname)
	}

	provisioningDetailsJson, err := json.Marshal(scriptingDeployment.ProvisioningDetails)
	if err != nil {
		return errors.WithMessage(err, "Unable to Marshal provisioning_details into string")
	}

	provisioningDetailsString := string(provisioningDetailsJson)
	if err := d.Set("provisioning_details", provisioningDetailsString); err != nil {
		return errors.WithMessage(err, "Cannot set provisioning_details: "+provisioningDetailsString)
	}

	scriptingPolicyURLSplit := strings.Split(scriptingDeployment.Links.Policy.Href, "/")
	scriptingPolicyID := scriptingPolicyURLSplit[len(scriptingPolicyURLSplit)-2]
	scriptingPolicyIDInt, _ := strconv.Atoi(scriptingPolicyID)
	if err := d.Set("policy_id", scriptingPolicyIDInt); err != nil {
		return errors.WithMessage(err, "Cannot set policy")
	}

	return nil
}

func resourceScriptingDeploymentCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceScriptingDeploymentCreate")

	config := m.(Config)

	newScriptingDeployment := ScriptingDeployment{
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	scriptingDeployment, err := config.NewOneFuseApiClient().CreateScriptingDeployment(&newScriptingDeployment)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(scriptingDeployment.ID))

	return bindScriptingDeploymentResource(d, scriptingDeployment)
}

func resourceScriptingDeploymentRead(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceScriptingDeploymentRead")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	scriptingDeployment, err := config.NewOneFuseApiClient().GetScriptingDeployment(intID)
	if err != nil {
		return err
	}

	return bindScriptingDeploymentResource(d, scriptingDeployment)
}

func resourceScriptingDeploymentUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceScriptingDeploymentUpdate")

	// Determine if a change is needed
	changed := (d.HasChange("policy_id") ||
		d.HasChange("workspace_url") ||
		d.HasChange("template_properties"))

	if !changed {
		return nil
	}

	// Make the API call to update the computer account
	config := m.(Config)

	// Create the desired Scripting Deployment
	id := d.Id()
	desiredScriptingDeployment := ScriptingDeployment{
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	scriptingDeployment, err := config.NewOneFuseApiClient().UpdateScriptingDeployment(intID, &desiredScriptingDeployment)
	if err != nil {
		return err
	}

	return bindScriptingDeploymentResource(d, scriptingDeployment)
}

func resourceScriptingDeploymentDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceScriptingDeploymentDelete")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return config.NewOneFuseApiClient().DeleteScriptingDeployment(intID)
}

func importScriptingReservation(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Println("onefuse.importScriptingReservation - Starting the import")

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

	scriptRecord, err := config.NewOneFuseApiClient().GetScriptingDeployment(intID)
	if err != nil {
		log.Printf("Error fetching script reservation: %v", err)
		return nil, errors.Wrap(err, "error fetching script reservation")
	}

	// Bind the scripting reservation
	if err := bindScriptingDeploymentResource(d, scriptRecord); err != nil {
		log.Printf("Error binding script reservation resource: %v", err)
		return nil, errors.Wrap(err, "failed to bind script reservation data")
	}

	jobMetaDataRecord, err := fetchScriptJobMetaData(scriptRecord, &config)
	if err != nil {
		log.Printf("Error fetching job metadata: %v", err)
		return nil, errors.Wrap(err, "error fetching job metadata during import")
	}

	if jobMetaDataRecord == nil {
		log.Println("jobMetaDataRecord is nil after fetching job metadata")
		return nil, errors.New("jobMetaDataRecord is nil after fetching job metadata")
	}

	log.Printf("Template Properties are: %+v", jobMetaDataRecord.ResolvedProperties)
	log.Println("onefuse.importScriptingReservation - import completed successfully")
	return []*schema.ResourceData{d}, nil
}

func fetchScriptJobMetaData(scriptRecord *ScriptingDeployment, config *Config) (*JobMetaData, error){
	log.Println("Fetching the job metadata - Start")

	jobMetaDataURLSplit := strings.Split(scriptRecord.Links.JobMetadata.Href, "/")
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
