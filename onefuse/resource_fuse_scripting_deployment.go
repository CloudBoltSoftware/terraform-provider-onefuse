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

func resourceScriptingDeployment() *schema.Resource {
	return &schema.Resource{
		Create: resourceScriptingDeploymentCreate,
		Read:   resourceScriptingDeploymentRead,
		Update: resourceScriptingDeploymentUpdate,
		Delete: resourceScriptingDeploymentDelete,
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
			"template_properties": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func bindScriptingDeploymentResource(d *schema.ResourceData, scriptingDeployment *ScriptingDeployment) error {
	log.Println("onefuse.bindScriptingDeploymentResource")

	if err := d.Set("workspace_url", scriptingDeployment.Links.Workspace.Href); err != nil {
		return errors.WithMessage(err, "Cannot set workspace: "+scriptingDeployment.Links.Workspace.Href)
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
