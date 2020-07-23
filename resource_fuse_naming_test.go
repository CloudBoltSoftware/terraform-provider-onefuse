// Copyright 2020 CloudBolt Software
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"log"
	"strconv"
	"testing"
)

func TestGenerateCustomName(t *testing.T) {
	config := GetConfig()
	cn := config.NewFuseApiClient().GenerateCustomName("sovlabs.net", "2", "", nil)
	log.Println(strconv.Itoa(cn.Id) + ": " + cn.Name + "." + cn.DnsSuffix + " version:" + strconv.Itoa(cn.Version))
	if cn.Id <= 0 {
		t.Errorf("customName.Id=%d; want > 0", cn.Id)
	}
	if cn.DnsSuffix != "sovlabs.net" {
		t.Errorf("customName.DnsSuffix=%s; want sovlabs.net", cn.DnsSuffix)
	}
	if cn.Name == "" {
		t.Errorf("customName.Name=%s; want non-empty string", cn.Name)
	}
}

func TestGetCustomName(t *testing.T) {
	config := GetConfig()
	cn1, _ := config.NewFuseApiClient().GenerateCustomName("sovlabs.net", "2", "", nil)
	cn2 := config.NewFuseApiClient().GetCustomName(cn1.Id)
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
		scheme:   "http",
		address:  "localhost",
		port:     "8000",
		user:     "admin2",
		password: "adminpass",
	}
	return config
}
