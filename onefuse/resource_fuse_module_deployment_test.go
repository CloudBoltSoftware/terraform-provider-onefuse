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

func createModuleDeployment() (*ModuleDeployment, error) {
	config := GetConfig()

	// Get raw user input from the environment
	ModulePolicyID, _ := strconv.Atoi(getEnv("CB_ONEFUSE_CFG_MODULE_POLICY_ID", "1"))
	ModuleDeploymentTemplatePropertiesStr := getEnv("CB_ONEFUSE_CFG_MODULE_DEPLOYMENT_TEMPLATE_PROPERTIES", "{}")

	// Parse string input into structures
	var ModuleDeploymentTemplateProperties map[string]interface{}
	json.Unmarshal([]byte(ModuleDeploymentTemplatePropertiesStr), &ModuleDeploymentTemplateProperties)

	newModuleDeployment := ModuleDeployment{
		PolicyID:           ModulePolicyID,
		TemplateProperties: ModuleDeploymentTemplateProperties,
	}

	ModuleDeployment, err := config.NewOneFuseApiClient().CreateModuleDeployment(&newModuleDeployment)
	if err != nil {
		return ModuleDeployment, err
	}

	// Verify the create
	verifyModuleDeployment, err := config.NewOneFuseApiClient().GetModuleDeployment(ModuleDeployment.ID)
	if verifyModuleDeployment.ID == 0 {
		err = errors.New("Error verifying created Module Deployment")
	}

	return ModuleDeployment, err
}

func deleteModuleDeployment(id int) error {
	config := GetConfig()
	err := config.NewOneFuseApiClient().DeleteModuleDeployment(id)
	if err != nil {
		return err
	}

	// Verify the delete
	ModuleDeployment, err := config.NewOneFuseApiClient().GetModuleDeployment(id)
	if ModuleDeployment != nil {
		return errors.New("Error verifying deleted Module Deployment")
	}

	return nil
}

func TestResourceModuleDeploymentCreate(t *testing.T) {
	ModuleDeployment, err := createModuleDeployment()
	if err != nil {
		t.Errorf("Error creating Module Deployment: '%s'", err)
		return
	}
	err = deleteModuleDeployment(ModuleDeployment.ID)
	if err != nil {
		t.Errorf("Error cleaning up Module Deployment: '%s'", err)
		return
	}
}

func TestResourceModuleDeploymentGet(t *testing.T) {
	ModuleDeployment, err := createModuleDeployment()
	if err != nil {
		t.Errorf("Error creating Module Deployment for read: '%s'", err)
		return
	}

	config := GetConfig()
	ModuleDeploymentVerify, err := config.NewOneFuseApiClient().GetModuleDeployment(ModuleDeployment.ID)
	if err != nil {
		t.Errorf("Error creating Module Deployment for read: '%s'", err)
	}
	if ModuleDeploymentVerify.PolicyID != ModuleDeployment.PolicyID {
		t.Errorf("Bad policy for Module Deployment; expected '%d' but got '%d'", ModuleDeploymentVerify.PolicyID, ModuleDeployment.PolicyID)
	}
	if ModuleDeploymentVerify.WorkspaceURL != ModuleDeployment.WorkspaceURL {
		t.Errorf("Bad workspace for Module Deployment; expected '%s' but got '%s'", ModuleDeploymentVerify.WorkspaceURL, ModuleDeployment.WorkspaceURL)
	}

	err = deleteModuleDeployment(ModuleDeployment.ID)
	if err != nil {
		t.Errorf("Error cleaning up Module Deployment: '%s'", err)
		return
	}
}

func TestResourceModuleDeploymentUpdate(t *testing.T) {
	// Updates not yet supported for Module Deployments.
	return
}

func TestResourceModuleDeploymentDelete(t *testing.T) {
	ModuleDeployment, err := createModuleDeployment()
	if err != nil {
		t.Errorf("Error creating Module Deployment for delete: '%s'", err)
		return
	}
	err = deleteModuleDeployment(ModuleDeployment.ID)
	if err != nil {
		t.Errorf("Error deleting Module Deployment: '%s'", err)
		return
	}
}
