// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"strconv"
	"testing"

	"github.com/pkg/errors"
)

func createAnsibleTowerDeployment(name string) (*AnsibleTowerDeployment, error) {
	config := GetConfig()
	ansibleTowerPolicyID, _ := strconv.Atoi(getEnv("CB_ONEFUSE_CFG_ANSIBLE_TOWER_POLICY_ID", "1"))
	newAnsibleTowerDeployment := AnsibleTowerDeployment{
		PolicyID:           ansibleTowerPolicyID,
	}
	ansibleTowerDeployment, err := config.NewOneFuseApiClient().CreateAnsibleTowerDeployment(&newAnsibleTowerDeployment)
	if err != nil {
		return ansibleTowerDeployment, err
	}
	// Verify the create
	verifyAnsibleTowerDeployment, err := config.NewOneFuseApiClient().GetAnsibleTowerDeployment(ansibleTowerDeployment.ID)
	if verifyAnsibleTowerDeployment.ID == 0 {
		err = errors.New("Error verifying created Ansible Tower Deployment")
	}
	return ansibleTowerDeployment, err
}

func deleteAnsibleTowerDeployment(id int) error {
	config := GetConfig()
	err := config.NewOneFuseApiClient().DeleteAnsibleTowerDeployment(id)
	if err != nil {
		return err
	}

	// Verify the delete
	ansibleTowerDeployment, err := config.NewOneFuseApiClient().GetAnsibleTowerDeployment(id)
	if ansibleTowerDeployment != nil {
		return errors.New("Error verifying deleted Ansible Tower Deployment")
	}

	return nil
}

func TestResourceAnsibleTowerDeploymentCreate(t *testing.T) {
	ansibleTowerDeployment, err := createAnsibleTowerDeployment("myAnsibleTowerDeployment")
	if err != nil {
		t.Errorf("Error creating Ansible Tower Deployment: '%s'", err)
		return
	}
	err = deleteAnsibleTowerDeployment(ansibleTowerDeployment.ID)
	if err != nil {
		t.Errorf("Error cleaning up Ansible Tower Deployment: '%s'", err)
		return
	}
}

func TestResourceAnsibleTowerDeploymentGet(t *testing.T) {
	ansibleTowerDeployment, err := createAnsibleTowerDeployment("myAnsibleTowerDeployment")
	if err != nil {
		t.Errorf("Error creating Ansible Tower Deployment for read: '%s'", err)
		return
	}

	config := GetConfig()
	ansibleTowerDeploymentVerify, err := config.NewOneFuseApiClient().GetAnsibleTowerDeployment(ansibleTowerDeployment.ID)
	if err != nil {
		t.Errorf("Error creating Ansible Tower Deployment for read: '%s'", err)
	}
	if ansibleTowerDeploymentVerify.PolicyID != ansibleTowerDeployment.PolicyID {
		t.Errorf("Bad policy for Ansible Tower Deployment; expected '%d' but got '%d'", ansibleTowerDeploymentVerify.PolicyID, ansibleTowerDeployment.PolicyID)
	}
	if ansibleTowerDeploymentVerify.WorkspaceURL != ansibleTowerDeployment.WorkspaceURL {
		t.Errorf("Bad workspace for Ansible Tower Deployment; expected '%s' but got '%s'", ansibleTowerDeploymentVerify.WorkspaceURL, ansibleTowerDeployment.WorkspaceURL)
	}

	err = deleteAnsibleTowerDeployment(ansibleTowerDeployment.ID)
	if err != nil {
		t.Errorf("Error cleaning up Ansible Tower Deployment: '%s'", err)
		return
	}
}

func TestResourceAnsibleTowerDeploymentUpdate(t *testing.T) {
	// Updates not yet supported for Ansible Tower Deployments.
	return
}

func TestResourceAnsibleTowerDeploymentDelete(t *testing.T) {
	ansibleTowerDeployment, err := createAnsibleTowerDeployment("myAnsibleTowerDeployment")
	if err != nil {
		t.Errorf("Error creating Ansible Tower Deployment for delete: '%s'", err)
		return
	}
	err = deleteAnsibleTowerDeployment(ansibleTowerDeployment.ID)
	if err != nil {
		t.Errorf("Error deleting Ansible Tower Deployment: '%s'", err)
		return
	}
}

