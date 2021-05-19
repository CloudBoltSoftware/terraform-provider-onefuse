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

func resourceVraDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceVraDeploymentCreate,
		Read:   resourceVraDeploymentRead,
		Update: resourceVraDeploymentUpdate,
		Delete: resourceVraDeploymentDelete,
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
			"deployment_name": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"template_properties": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"deployment_info": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"blueprint_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
			"project_name": {
				Type:     schema.TypeString,
				Computed: true,
				Optional: true,
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func bindVraDeploymentResource(d *schema.ResourceData, vraDeployment *VraDeployment) error {
	log.Println("onefuse.bindVraDeploymentResource")

	if err := d.Set("workspace_url", vraDeployment.Links.Workspace.Href); err != nil {
		return errors.WithMessage(err, "Cannot set workspace: "+vraDeployment.Links.Workspace.Href)
	}

	if err := d.Set("deployment_name", vraDeployment.Name); err != nil {
		return errors.WithMessage(err, "Cannot set deployment name: "+vraDeployment.Name)
	}

	deploymentInfoJSON, err := json.Marshal(vraDeployment.DeploymentInfo)
	if err != nil {
		return errors.WithMessage(err, "Unable to Marshal deployment_info into string")
	}
	deploymentInfoString := string(deploymentInfoJSON)
	if err := d.Set("deployment_info", deploymentInfoString); err != nil {
		return errors.WithMessage(err, "Cannot set deployment_info: "+deploymentInfoString)
	}

	if err := d.Set("blueprint_name", vraDeployment.BlueprintName); err != nil {
		return errors.WithMessage(err, "Cannot set blueprint name: "+vraDeployment.BlueprintName)
	}

	if err := d.Set("project_name", vraDeployment.ProjectName); err != nil {
		return errors.WithMessage(err, "Cannot set project name: "+vraDeployment.ProjectName)
	}

	vraPolicyURLSplit := strings.Split(vraDeployment.Links.Policy.Href, "/")
	vraPolicyID := vraPolicyURLSplit[len(vraPolicyURLSplit)-2]
	vraPolicyIDInt, _ := strconv.Atoi(vraPolicyID)
	if err := d.Set("policy_id", vraPolicyIDInt); err != nil {
		return errors.WithMessage(err, "Cannot set policy")
	}

	return nil
}

func resourceVraDeploymentCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceVraDeploymentCreate")

	config := m.(Config)

	newVraDeployment := VraDeployment{
		PolicyID:           d.Get("policy_id").(int),
		WorkspaceURL:       d.Get("workspace_url").(string),
		DeploymentName:     d.Get("deployment_name").(string),
		TemplateProperties: d.Get("template_properties").(map[string]interface{}),
	}

	vraDeployment, err := config.NewOneFuseApiClient().CreateVraDeployment(&newVraDeployment)
	if err != nil {
		return err
	}
	d.SetId(strconv.Itoa(vraDeployment.ID))

	return bindVraDeploymentResource(d, vraDeployment)
}

func resourceVraDeploymentRead(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceVraDeploymentRead")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	vraDeployment, err := config.NewOneFuseApiClient().GetVraDeployment(intID)
	if err != nil {
		return err
	}

	return bindVraDeploymentResource(d, vraDeployment)
}

func resourceVraDeploymentUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceVraDeploymentUpdate")
	log.Println("No Op!")
	return nil
}

func resourceVraDeploymentDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceVraDeploymentDelete")

	config := m.(Config)

	id := d.Id()
	intID, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	return config.NewOneFuseApiClient().DeleteVraDeployment(intID)
}
