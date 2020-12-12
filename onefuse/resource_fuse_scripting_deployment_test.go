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

func createScriptingDeployment() (*ScriptingDeployment, error) {
	config := GetConfig()

	// Get raw user input from the environment
	scriptingPolicyID, _ := strconv.Atoi(getEnv("CB_ONEFUSE_CFG_SCRIPTING_POLICY_ID", "1"))
	scriptingDeploymentTemplatePropertiesStr := getEnv("CB_ONEFUSE_CFG_SCRIPTING_DEPLOYMENT_TEMPLATE_PROPERTIES", "{}")

	// Parse string input into structures
	var scriptingDeploymentTemplateProperties map[string]interface{}
	json.Unmarshal([]byte(scriptingDeploymentTemplatePropertiesStr), &scriptingDeploymentTemplateProperties)

	newScriptingDeployment := ScriptingDeployment{
		PolicyID:           scriptingPolicyID,
		TemplateProperties: scriptingDeploymentTemplateProperties,
		// TODO: Support [Hosts] list too
	}

	scriptingDeployment, err := config.NewOneFuseApiClient().CreateScriptingDeployment(&newScriptingDeployment)
	if err != nil {
		return scriptingDeployment, err
	}

	// Verify the create
	verifyScriptingDeployment, err := config.NewOneFuseApiClient().GetScriptingDeployment(scriptingDeployment.ID)
	if verifyScriptingDeployment.ID == 0 {
		err = errors.New("Error verifying created Scripting Deployment")
	}

	return scriptingDeployment, err
}

func deleteScriptingDeployment(id int) error {
	config := GetConfig()
	err := config.NewOneFuseApiClient().DeleteScriptingDeployment(id)
	if err != nil {
		return err
	}

	// Verify the delete
	scriptingDeployment, err := config.NewOneFuseApiClient().GetScriptingDeployment(id)

	if !scriptingDeployment.Archived {
		return errors.New("Error verifying deleted Scripting Deployment")
	}

	return nil
}

func TestResourceScriptingDeploymentCreate(t *testing.T) {
	scriptingDeployment, err := createScriptingDeployment()
	if err != nil {
		t.Errorf("Error creating Scripting Deployment: '%s'", err)
		return
	}
	err = deleteScriptingDeployment(scriptingDeployment.ID)
	if err != nil {
		t.Errorf("Error cleaning up Scripting Deployment: '%s'", err)
		return
	}
}

func TestResourceScriptingDeploymentGet(t *testing.T) {
	scriptingDeployment, err := createScriptingDeployment()
	if err != nil {
		t.Errorf("Error creating Scripting Deployment for read: '%s'", err)
		return
	}

	config := GetConfig()
	scriptingDeploymentVerify, err := config.NewOneFuseApiClient().GetScriptingDeployment(scriptingDeployment.ID)
	if err != nil {
		t.Errorf("Error creating Scripting Deployment for read: '%s'", err)
	}
	if scriptingDeploymentVerify.PolicyID != scriptingDeployment.PolicyID {
		t.Errorf("Bad policy for Scripting Deployment; expected '%d' but got '%d'", scriptingDeploymentVerify.PolicyID, scriptingDeployment.PolicyID)
	}
	if scriptingDeploymentVerify.WorkspaceURL != scriptingDeployment.WorkspaceURL {
		t.Errorf("Bad workspace for Scripting Deployment; expected '%s' but got '%s'", scriptingDeploymentVerify.WorkspaceURL, scriptingDeployment.WorkspaceURL)
	}

	err = deleteScriptingDeployment(scriptingDeployment.ID)
	if err != nil {
		t.Errorf("Error cleaning up Scripting Deployment: '%s'", err)
		return
	}
}

func TestResourceScriptingDeploymentUpdate(t *testing.T) {
	// Updates not yet supported for Scripting Deployments.
	return
}

func TestResourceScriptingDeploymentDelete(t *testing.T) {
	scriptingDeployment, err := createScriptingDeployment()
	if err != nil {
		t.Errorf("Error creating Scripting Deployment for delete: '%s'", err)
		return
	}
	err = deleteScriptingDeployment(scriptingDeployment.ID)
	if err != nil {
		t.Errorf("Error deleting Scripting Deployment: '%s'", err)
		return
	}
}
