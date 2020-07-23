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
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
)

const ApiVersion = "/api/v3/"
const ApiNamespace = "onefuse"
const NamingResourceType = "customNames"
const WorkspaceResourceType = "workspaces"

type OneFuseAPIClient struct {
	config *Config
}

type CustomName struct {
	Id        int
	Version   int
	Name      string
	DnsSuffix string
}

type WorkspacesListResponse struct {
	Embedded struct {
		Workspaces []struct {
			Name string `json:"name"`
			ID   int    `json:"id"`
		} `json:"workspaces`
	} `json:"_embedded"`
}

func (c *Config) NewOneFuseApiClient() *OneFuseAPIClient {
	return &OneFuseAPIClient{
		config: c,
	}
}

func (apiClient *OneFuseAPIClient) GenerateCustomName(dnsSuffix string, namingPolicyID string, workspaceID string,
	templateProperties map[string]interface{}) (result *CustomName, err error) {

	config := apiClient.config
	url := collectionURL(config, NamingResourceType)
	log.Println("reserving custom name from " + url + "  dnsSuffix=" + dnsSuffix)

	if templateProperties == nil {
		templateProperties = make(map[string]interface{})
	}
	if workspaceID == "" {
		workspaceID, err = findDefaultWorkspaceID(config)
		if err != nil {
			return
		}
	}

	postBody := map[string]interface{}{
		"namingPolicy":       fmt.Sprintf("%s%s/namingPolicies/%s/", ApiVersion, ApiNamespace, namingPolicyID),
		"templateProperties": templateProperties,
		"workspace":          fmt.Sprintf("%s%s/workspaces/%s/", ApiVersion, ApiNamespace, workspaceID),
	}
	var jsonBytes []byte
	jsonBytes, err = json.Marshal(postBody)
	requestBody := string(jsonBytes)
	if err != nil {
		err = errors.New("unable to marshal request body to JSON")
		return
	}
	payload := strings.NewReader(requestBody)

	log.Println("CONFIG:")
	log.Println(config)
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return
	}
	log.Println("HTTP PAYLOAD to " + url + ":")
	log.Println(postBody)

	setHeaders(req, config)

	client := getHttpClient(config)
	var res *http.Response
	res, err = client.Do(req)
	if err != nil {
		return
	}

	checkForErrors(res)

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	log.Println("HTTP POST RESULTS:")
	log.Println(string(body))
	json.Unmarshal(body, &result)
	res.Body.Close()

	if result == nil {
		err = errors.New("invalid response " + strconv.Itoa(res.StatusCode) + " while generating a custom name: " + string(body))
		return
	}

	log.Println("custom name reserved: " +
		"custom_name_id=" + strconv.Itoa(result.Id) +
		" name=" + result.Name +
		" dnsSuffix=" + result.DnsSuffix)
	return
}

func (apiClient *OneFuseAPIClient) GetCustomName(id int) (result CustomName, err error) {
	config := apiClient.config
	url := itemURL(config, NamingResourceType, id)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	setHeaders(req, config)

	log.Println("REQUEST:")
	log.Println(req)
	client := getHttpClient(config)
	res, _ := client.Do(req)
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return
	}
	log.Println("HTTP GET RESULTS:")
	log.Println(string(body))

	json.Unmarshal(body, &result)
	res.Body.Close()
	return
}

func (apiClient *OneFuseAPIClient) DeleteCustomName(id int) error {
	config := apiClient.config
	url := itemURL(config, NamingResourceType, id)
	req, _ := http.NewRequest("DELETE", url, nil)
	setHeaders(req, config)
	client := getHttpClient(config)
	res, err := client.Do(req)
	if err != nil {
		return err
	}
	return checkForErrors(res)
}

func findDefaultWorkspaceID(config *Config) (workspaceID string, err error) {
	filter := "filter=name.exact:Default"
	url := fmt.Sprintf("%s?%s", collectionURL(config, WorkspaceResourceType), filter)
	req, clientErr := http.NewRequest("GET", url, nil)
	if clientErr != nil {
		err = clientErr
		return
	}

	setHeaders(req, config)

	client := getHttpClient(config)
	res, clientErr := client.Do(req)
	if clientErr != nil {
		err = clientErr
		return
	}

	checkForErrors(res)

	body, clientErr := ioutil.ReadAll(res.Body)
	if clientErr != nil {
		err = clientErr
		return
	}

	var data WorkspacesListResponse
	json.Unmarshal(body, &data)
	res.Body.Close()

	workspaces := data.Embedded.Workspaces
	if len(workspaces) == 0 {
		panic("Unable to find default workspace.")
	}
	workspaceID = strconv.Itoa(workspaces[0].ID)
	return
}

func getHttpClient(config *Config) *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.verifySSL},
	}
	return &http.Client{Transport: tr}
}

func checkForErrors(res *http.Response) error {
	if res.StatusCode >= 500 {
		b, _ := ioutil.ReadAll(res.Body)
		return errors.New(string(b))
	}
	return nil
}

func setStandardHeaders(req *http.Request) {
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "*/*")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("accept-encoding", "gzip, deflate")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("cache-control", "no-cache")
}

func setHeaders(req *http.Request, config *Config) {
	setStandardHeaders(req)
	req.Header.Add("Host", config.address+":"+config.port)
	req.Header.Add("SOURCE", "Terraform")
	req.SetBasicAuth(config.user, config.password)
}

func collectionURL(config *Config, resourceType string) string {
	address := config.address
	port := config.port
	return config.scheme + "://" + address + ":" + port + ApiVersion + ApiNamespace + "/" + resourceType + "/"
}

func itemURL(config *Config, resourceType string, id int) string {
	address := config.address
	port := config.port
	idString := strconv.Itoa(id)
	return config.scheme + "://" + address + ":" + port + ApiVersion + ApiNamespace + "/" + resourceType + "/" + idString + "/"
}
