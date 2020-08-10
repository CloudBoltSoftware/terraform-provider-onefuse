// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"os"
	"testing"
)

func TestGenerateCustomName(t *testing.T) {
	config := GetConfig()
	cn, err := config.NewOneFuseApiClient().GenerateCustomName(getEnv("CB_ONEFUSE_CFG_NAMING_POLICY_ID", "1"), "", nil)
	if err != nil {
		t.Errorf("generate custom name error '%s'", err)
		return
	}
	if cn.Id <= 0 {
		t.Errorf("customName.Id=%d; want > 0", cn.Id)
	}
	// TODO: this assertion is only true if the NamingPolicy has its dnsSuffix set to "sovlabs.net". The
	//       value being passed into GenerateCustomName is not being used anywhere.
	//
	// if cn.DnsSuffix != "sovlabs.net" {
	// 	t.Errorf("customName.DnsSuffix=%s; want sovlabs.net", cn.DnsSuffix)
	// }
	if cn.Name == "" {
		t.Errorf("customName.Name=%s; want non-empty string", cn.Name)
	}
}

func TestGetCustomName(t *testing.T) {
	config := GetConfig()
	cn1, err := config.NewOneFuseApiClient().GenerateCustomName("sovlabs.net", getEnv("CB_ONEFUSE_CFG_NAMING_POLICY_ID", "1"), "", nil)
	if err != nil {
		t.Errorf("generate custom name error '%s'", err)
		return
	}
	cn2, err := config.NewOneFuseApiClient().GetCustomName(cn1.Id)
	if err != nil {
		t.Errorf("get custom name error '%s'", err)
		return
	}
	if cn1.Id != cn2.Id {
		t.Error("Reserved customName.Id does not match after retrieval")
	}
	if cn1.Name != cn2.Name {
		t.Error("Reserved customName.Name does not match after retrieval")
	}
	if cn1.DnsSuffix != cn2.DnsSuffix {
		t.Error("Reserved customName.DnsSuffix does not match after retrieval")
	}

}

func GetConfig() Config {
	config := Config{
		scheme:   getEnv("CB_ONEFUSE_CFG_SCHEME", "https"),
		address:  getEnv("CB_ONEFUSE_CFG_ADDRESS", "localhost"),
		port:     getEnv("CB_ONEFUSE_CFG_PORT", "443"),
		user:     getEnv("CB_ONEFUSE_CFG_USER", "admin"),
		password: getEnv("CB_ONEFUSE_CFG_PASSWORD", "admin"),
	}
	return config
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}
