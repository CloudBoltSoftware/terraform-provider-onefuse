// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"testing"
)

// Requires a Microsoft Endpoint named "myMicrosoftEndpoint"
func TestGetMicrosoftEndpointByName(t *testing.T) {
	tables := []struct {
		name   string
		result bool
	}{
		{getEnv("CB_ONEFUSE_CFG_MICROSOFT_ENDPOINT_NAME", "myMicrosoftEndpoint"), true},
		{"idontexits", false},
	}
	config := GetConfig()

	for _, table := range tables {
		endpoint, err := config.NewOneFuseApiClient().GetMicrosoftEndpointByName(table.name)
		t.Logf("Endpoint: %v", endpoint)
		if table.result == true && err != nil {
			t.Errorf("Error getting expected endpoint by name '%s'", table.name)
		} else if table.result == false && err == nil {
			t.Errorf("Missing error getting nonexistent endpoint by name '%s'", table.name)
		}
	}

}
