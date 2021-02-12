// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"encoding/json"
	"strconv"
	"testing"

	"github.com/pkg/errors"
)

func createVraDeployment() (*VraDeployment, error) {
	config := GetConfig()

	// Get raw user input from the environment
	vraPolicyID, _ := strconv.Atoi(getEnv("CB_ONEFUSE_CFG_VRA_POLICY_ID", "2348"))
	vraDeploymentTemplatePropertiesStr := getEnv("CB_ONEFUSE_CFG_VRA_DEPLOYMENT_TEMPLATE_PROPERTIES", "{\"property1\": \"test\"}")
	vraDeploymentName := getEnv("CB_ONEFUSE_CFG_VRA_DEPLOYMENT_NAME", "")

	// Parse string input into structures
	var vraDeploymentTemplateProperties map[string]interface{}
	json.Unmarshal([]byte(vraDeploymentTemplatePropertiesStr), &vraDeploymentTemplateProperties)

	newVraDeployment := VraDeployment{
		PolicyID:           vraPolicyID,
		TemplateProperties: vraDeploymentTemplateProperties,
		DeploymentName:     vraDeploymentName,
	}

	vraDeployment, err := config.NewOneFuseApiClient().CreateVraDeployment(&newVraDeployment)
	if err != nil {
		return vraDeployment, err
	}

	// Verify the create
	verifyVraDeployment, err := config.NewOneFuseApiClient().GetVraDeployment(vraDeployment.ID)
	if verifyVraDeployment.ID == 0 {
		err = errors.New("Error verifying created vRealize Automation Deployment")
	}

	return vraDeployment, err
}

func deleteVraDeployment(id int) error {
	config := GetConfig()
	err := config.NewOneFuseApiClient().DeleteVraDeployment(id)
	if err != nil {
		return err
	}

	// Verify the delete
	vraDeployment, err := config.NewOneFuseApiClient().GetVraDeployment(id)
	if vraDeployment != nil {
		return errors.New("Error verifying deleted vRealize Automation Deployment")
	}

	return nil
}

func TestResourceVraDeploymentCreate(t *testing.T) {
	vraDeployment, err := createVraDeployment()
	if err != nil {
		t.Errorf("Error creating vRealize Automation Deployment: '%s'", err)
		return
	}
	err = deleteVraDeployment(vraDeployment.ID)
	if err != nil {
		t.Errorf("Error cleaning up vRealize Automation Deployment: '%s'", err)
		return
	}
}

func TestResourceVraDeploymentGet(t *testing.T) {
	vraDeployment, err := createVraDeployment()
	if err != nil {
		t.Errorf("Error creating vRealize Automation Deployment for read: '%s'", err)
		return
	}

	config := GetConfig()
	vraDeploymentVerify, err := config.NewOneFuseApiClient().GetVraDeployment(vraDeployment.ID)
	if err != nil {
		t.Errorf("Error creating vRealize Automation Deployment for read: '%s'", err)
	}
	if vraDeploymentVerify.PolicyID != vraDeployment.PolicyID {
		t.Errorf("Bad policy for vRealize Automation Deployment; expected '%d' but got '%d'", vraDeploymentVerify.PolicyID, vraDeployment.PolicyID)
	}
	if vraDeploymentVerify.WorkspaceURL != vraDeployment.WorkspaceURL {
		t.Errorf("Bad workspace for vRealize Automation Deployment; expected '%s' but got '%s'", vraDeploymentVerify.WorkspaceURL, vraDeployment.WorkspaceURL)
	}

	err = deleteVraDeployment(vraDeployment.ID)
	if err != nil {
		t.Errorf("Error cleaning up vRealize Automation Deployment: '%s'", err)
		return
	}
}

func TestResourceVraDeploymentUpdate(t *testing.T) {
	// Updates not yet supported for vRealize Automation Deployments.
	return
}

func TestResourceVraDeploymentDelete(t *testing.T) {
	vraDeployment, err := createVraDeployment()
	if err != nil {
		t.Errorf("Error creating vRealize Automation Deployment for delete: '%s'", err)
		return
	}
	err = deleteVraDeployment(vraDeployment.ID)
	if err != nil {
		t.Errorf("Error deleting vRealize Automation Deployment: '%s'", err)
		return
	}
}
