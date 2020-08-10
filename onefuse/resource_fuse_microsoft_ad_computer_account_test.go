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

func createADComputerAccount(name string) (*MicrosoftADComputerAccount, error) {
	config := GetConfig()
	adPolicyID, _ := strconv.Atoi(getEnv("CB_ONEFUSE_CFG_MICROSOFT_AD_POLICY_ID", "1"))
	newComputerAccount := MicrosoftADComputerAccount{
		Name:         name,
		PolicyID:     adPolicyID,
		WorkspaceURL: "",
	}
	computerAccount, err := config.NewOneFuseApiClient().CreateMicrosoftADComputerAccount(&newComputerAccount)
	if err != nil {
		return computerAccount, err
	}
	// Verify the create
	verifyComputerAccount, err := config.NewOneFuseApiClient().GetMicrosoftADComputerAccount(computerAccount.ID)
	if verifyComputerAccount.ID == 0 {
		err = errors.New("Error verifying created Microsoft AD Computer Account")
	}
	return computerAccount, err
}

func deleteADComputerAccount(id int) error {
	config := GetConfig()
	err := config.NewOneFuseApiClient().DeleteMicrosoftADComputerAccount(id)
	if err != nil {
		return err
	}

	// Verify the delete
	computerAccount, err := config.NewOneFuseApiClient().GetMicrosoftADComputerAccount(id)
	if computerAccount != nil {
		return errors.New("Error verifying deleted Microsoft AD ComputerAccount")
	}

	return nil
}

func TestResourceMicrosoftADComputerAccountCreate(t *testing.T) {
	computerAccount, err := createADComputerAccount("myComputerAccout")
	if err != nil {
		t.Errorf("Error creating Microsoft AD Computer Account: '%s'", err)
		return
	}
	err = deleteADComputerAccount(computerAccount.ID)
	if err != nil {
		t.Errorf("Error cleaning up Microsoft AD Computer Account: '%s'", err)
		return
	}
}

func TestResourceMicrosoftADComputerAccountGet(t *testing.T) {
	computerAccount, err := createADComputerAccount("myComputerAccout")
	if err != nil {
		t.Errorf("Error creating Microsoft AD Policy for read: '%s'", err)
		return
	}

	config := GetConfig()
	computerAccountVerify, err := config.NewOneFuseApiClient().GetMicrosoftADComputerAccount(computerAccount.ID)
	if err != nil {
		t.Errorf("Error creating Microsoft AD Computer Account for read: '%s'", err)
	}
	if computerAccountVerify.Name != computerAccount.Name {
		t.Errorf("Bad name for AD Computer Account; expected '%s' but got '%s'", computerAccountVerify.Name, computerAccount.Name)
	}
	if computerAccountVerify.PolicyID != computerAccount.PolicyID {
		t.Errorf("Bad policy for AD Computer Account; expected '%d' but got '%d'", computerAccountVerify.PolicyID, computerAccount.PolicyID)
	}
	if computerAccountVerify.WorkspaceURL != computerAccount.WorkspaceURL {
		t.Errorf("Bad workspace for AD Computer Account; expected '%s' but got '%s'", computerAccountVerify.WorkspaceURL, computerAccount.WorkspaceURL)
	}

	err = deleteADComputerAccount(computerAccount.ID)
	if err != nil {
		t.Errorf("Error cleaning up Microsoft AD Computer Account: '%s'", err)
		return
	}
}

func TestResourceMicrosoftADComputerAccountUpdate(t *testing.T) {
	// Updates not yet supported for Microsoft Active Directory Computer Names.
	return
}

func TestResourceMicrosoftADComputerAccountDelete(t *testing.T) {
	computerAccount, err := createADComputerAccount("myComputerAccout")
	if err != nil {
		t.Errorf("Error creating Microsoft AD Computer Account for delete: '%s'", err)
		return
	}
	err = deleteADComputerAccount(computerAccount.ID)
	if err != nil {
		t.Errorf("Error deleting Microsoft AD Computer Account: '%s'", err)
		return
	}
}
