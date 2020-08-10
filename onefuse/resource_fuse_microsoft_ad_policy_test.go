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

func createADPolicy(name string) (*MicrosoftADPolicy, error) {
	config := GetConfig()
	endpointID, _ := strconv.Atoi(getEnv("CB_ONEFUSE_CFG_MICROSOFT_ENDPOINT_ID", "1"))
	newPolicy := MicrosoftADPolicy{
		Name:                   name,
		Description:            "Description",
		OU:                     "OU=Foo,DC=Bar",
		CreateOU:               true,
		RemoveOU:               true,
		SecurityGroups:         []string{"CN=Group,OU=Groups,DC=Bar"},
		MicrosoftEndpointID:    endpointID,
		ComputerNameLetterCase: "UPPER",
		WorkspaceURL:           "",
	}
	policy, err := config.NewOneFuseApiClient().CreateMicrosoftADPolicy(&newPolicy)
	if err != nil {
		return policy, err
	}
	// Verify the create
	verifyPolicy, err := config.NewOneFuseApiClient().GetMicrosoftADPolicy(policy.ID)
	if verifyPolicy.ID == 0 {
		err = errors.New("Error verifying created Microsoft AD policy")
	}
	return policy, err
}

func deleteADPolicy(id int) error {
	config := GetConfig()
	err := config.NewOneFuseApiClient().DeleteMicrosoftADPolicy(id)
	if err != nil {
		return err
	}

	// Verify the delete
	// We expect a 404 error
	policy, _ := config.NewOneFuseApiClient().GetMicrosoftADPolicy(id)
	if policy != nil {
		err = errors.New("Error verifying deleted Microsoft AD policy")
	}

	return nil
}

func TestResourceMicrosoftADPolicyCreate(t *testing.T) {
	policy, err := createADPolicy("myMicrosoftADPolicy")
	if err != nil {
		t.Errorf("Error creating Microsoft AD Policy: '%s'", err)
		return
	}
	err = deleteADPolicy(policy.ID)
	if err != nil {
		t.Errorf("Error cleaning up Microsoft AD Policy: '%s'", err)
		return
	}
}

func TestResourceMicrosoftADPolicyGet(t *testing.T) {
	policy, err := createADPolicy("myMicrosoftADPolicy")
	if err != nil {
		t.Errorf("Error creating Microsoft AD Policy for read: '%s'", err)
		return
	}

	config := GetConfig()
	policyVerify, err := config.NewOneFuseApiClient().GetMicrosoftADPolicy(policy.ID)
	if err != nil {
		t.Errorf("Error creating Microsoft AD Policy for read: '%s'", err)
	}
	if policyVerify.Name != policy.Name {
		t.Errorf("Bad name for AD policy; expected '%s' but got '%s'", policyVerify.Name, policy.Name)
	}
	if policyVerify.Description != policy.Description {
		t.Errorf("Bad description for AD policy; expected '%s' but got '%s'", policyVerify.Description, policy.Description)
	}
	if policyVerify.MicrosoftEndpointID != policy.MicrosoftEndpointID {
		t.Errorf("Bad endpoint for AD policy; expected '%d' but got '%d'", policyVerify.MicrosoftEndpointID, policy.MicrosoftEndpointID)
	}
	if policyVerify.ComputerNameLetterCase != policy.ComputerNameLetterCase {
		t.Errorf("Bad ComputerNameLetterCase for AD policy; expected '%s' but got '%s'", policyVerify.ComputerNameLetterCase, policy.ComputerNameLetterCase)
	}
	if policyVerify.WorkspaceURL != policy.WorkspaceURL {
		t.Errorf("Bad workspace for AD policy; expected '%s' but got '%s'", policyVerify.WorkspaceURL, policy.WorkspaceURL)
	}
	if policyVerify.OU != policy.OU {
		t.Errorf("Bad OU for AD policy; expected '%s' but got '%s'", policyVerify.OU, policy.OU)
	}
	if policyVerify.CreateOU != policy.CreateOU {
		t.Errorf("Bad CreateOU for AD policy; expected '%s' but got '%s'", policyVerify.CreateOU, policy.CreateOU)
	}
	if policyVerify.RemoveOU != policy.RemoveOU {
		t.Errorf("Bad RemoveOU for AD policy; expected '%s' but got '%s'", policyVerify.RemoveOU, policy.RemoveOU)
	}
	if !testSliceEq(policyVerify.SecurityGroups, policy.SecurityGroups) {
		t.Errorf("Bad SecurityGroups for AD policy; expected '%s' but got '%s'", policyVerify.SecurityGroups, policy.SecurityGroups)
	}

	err = deleteADPolicy(policy.ID)
	if err != nil {
		t.Errorf("Error cleaning up Microsoft AD Policy: '%s'", err)
		return
	}
}

func TestResourceMicrosoftADPolicyUpdate(t *testing.T) {
	policy, err := createADPolicy("myMicrosoftADPolicy")
	if err != nil {
		t.Errorf("Error creating Microsoft AD Policy for read: '%s'", err)
		return
	}

	config := GetConfig()
	policy, err = config.NewOneFuseApiClient().GetMicrosoftADPolicy(policy.ID)
	newPolicy := MicrosoftADPolicy{
		Name:                   policy.Name,
		Description:            "I am a changed policy, I tell you!",
		OU:                     "OU=Updated,DC=Woohoo",
		ComputerNameLetterCase: "LOWER",
	}

	updatedPolicy, err := config.NewOneFuseApiClient().UpdateMicrosoftADPolicy(policy.ID, &newPolicy)
	if err != nil {
		t.Errorf("Error updating Microsoft AD Policy: '%s'", err)
	}
	if updatedPolicy.Name != newPolicy.Name {
		t.Errorf("Bad Name for AD policy; expected '%s' but got '%s'", updatedPolicy.Name, policy.Name)
	}
	if updatedPolicy.Description != newPolicy.Description {
		t.Errorf("Bad description for AD policy; expected '%s' but got '%s'", updatedPolicy.Description, policy.Description)
	}
	if updatedPolicy.ComputerNameLetterCase != newPolicy.ComputerNameLetterCase {
		t.Errorf("Bad ComputerNameLetterCase for AD policy; expected '%s' but got '%s'", updatedPolicy.ComputerNameLetterCase, policy.ComputerNameLetterCase)
	}
	if updatedPolicy.OU != newPolicy.OU {
		t.Errorf("Bad OU for AD policy; expected '%s' but got '%s'", updatedPolicy.OU, policy.OU)
	}
	if updatedPolicy.CreateOU != newPolicy.CreateOU {
		t.Errorf("Bad CreateOU for AD policy; expected '%s' but got '%s'", updatedPolicy.CreateOU, newPolicy.CreateOU)
	}
	if updatedPolicy.RemoveOU != newPolicy.RemoveOU {
		t.Errorf("Bad RemoveOU for AD policy; expected '%s' but got '%s'", updatedPolicy.RemoveOU, newPolicy.RemoveOU)
	}
	if !testSliceEq(updatedPolicy.SecurityGroups, newPolicy.SecurityGroups) {
		t.Errorf("Bad SecurityGroups for AD policy; expected '%s' but got '%s'", updatedPolicy.SecurityGroups, newPolicy.SecurityGroups)
	}

	err = deleteADPolicy(policy.ID)
	if err != nil {
		t.Errorf("Error cleaning up Microsoft AD Policy: '%s'", err)
		return
	}
}

func TestResourceMicrosoftADPolicyDelete(t *testing.T) {
	policy, err := createADPolicy("myMicrosoftADPolicy")
	if err != nil {
		t.Errorf("Error creating Microsoft AD Policy for delete: '%s'", err)
		return
	}

	err = deleteADPolicy(policy.ID)
	if err != nil {
		t.Errorf("Error deleting Microsoft AD Policy: '%s'", err)
	}
}

func testSliceEq(a []string, b []string) bool {
	// Source: https://stackoverflow.com/a/15312097/6500622
	// If one is nil, the other must also be nil.
	if (a == nil) != (b == nil) {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
