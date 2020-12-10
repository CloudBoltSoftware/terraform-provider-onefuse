// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const ApiVersion = "api/v3"
const ApiNamespace = "onefuse"
const NamingResourceType = "customNames"
const WorkspaceResourceType = "workspaces"
const MicrosoftADPolicyResourceType = "microsoftADPolicies"
const MicrosoftADComputerAccountResourceType = "microsoftADComputerAccounts"
const ModuleEndpointResourceType = "endpoints"
const DNSReservationResourceType = "dnsReservations"
const IPAMReservationResourceType = "ipamReservations"
const StaticPropertySetResourceType = "propertySets"
const IPAMPolicyResourceType = "ipamPolicies"
const NamingPolicyResourceType = "namingPolicies"
const ADPolicyResourceType = "microsoftADPolicies"
const DNSPolicyResourceType = "dnsPolicies"
const JobStatusResourceType = "jobStatus"

const JobSuccess = "Successful"
const JobFailed = "Failed"

type OneFuseAPIClient struct {
	config *Config
}

type CustomName struct {
	Id        int
	Name      string
	DnsSuffix string
}

type LinkRef struct {
	Href  string `json:"href,omitempty"`
	Title string `json:"title,omitempty"`
}

type Workspace struct {
	Links *struct {
		Self LinkRef `json:"self,omitempty"`
	} `json:"_links,omitempty"`
	ID   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type WorkspacesListResponse struct {
	Embedded struct {
		Workspaces []Workspace `json:"workspaces"`
	} `json:"_embedded"`
}

type EndpointsListResponse struct {
	Embedded struct {
		Endpoints []MicrosoftEndpoint `json:"endpoints"` // TODO: Generalize to Endpoints
	} `json:"_embedded"`
}

type MicrosoftEndpoint struct {
	Links *struct {
		Self       LinkRef `json:"self,omitempty"`
		Workspace  LinkRef `json:"workspace,omitempty"`
		Credential LinkRef `json:"credential,omitempty"`
	} `json:"_links,omitempty"`
	ID               int    `json:"id,omitempty"`
	Type             string `json:"type,omitempty"`
	Name             string `json:"name,omitempty"`
	Description      string `json:"description,omitempty"`
	Host             string `json:"host,omitempty"`
	Port             int    `json:"port,omitempty"`
	SSL              bool   `json:"ssl,omitempty"`
	MicrosoftVersion string `json:"microsoftVersion,omitempty"`
}

type MicrosoftADPolicy struct {
	Links *struct {
		Self              LinkRef `json:"self,omitempty"`
		Workspace         LinkRef `json:"workspace,omitempty"`
		MicrosoftEndpoint LinkRef `json:"microsoftEndpoint,omitempty"`
	} `json:"_links,omitempty"`
	Name                   string   `json:"name,omitempty"`
	ID                     int      `json:"id,omitempty"`
	Description            string   `json:"description,omitempty"`
	MicrosoftEndpointID    int      `json:"microsoftEndpointId,omitempty"`
	MicrosoftEndpoint      string   `json:"microsoftEndpoint,omitempty"`
	ComputerNameLetterCase string   `json:"computerNameLetterCase,omitempty"`
	WorkspaceURL           string   `json:"workspace,omitempty"`
	OU                     string   `json:"ou,omitempty"`
	CreateOU               bool     `json:"createOrganizationalUnit,omitempty"`
	RemoveOU               bool     `json:"removeOrganizationalUnit,omitempty"`
	SecurityGroups         []string `json:"securityGroups,omitempty"`
}

type MicrosoftADComputerAccount struct {
	Links *struct {
		Self        LinkRef `json:"self,omitempty"`
		Workspace   LinkRef `json:"workspace,omitempty"`
		Policy      LinkRef `json:"policy,omitempty"`
		JobMetadata LinkRef `json:"jobMetadata,omitempty"`
	} `json:"_links,omitempty"`
	ID                 int                    `json:"id,omitempty"`
	Name               string                 `json:"name,omitempty"`
	FinalOU            string                 `json:"finalOu"`
	PolicyID           int                    `json:"policyId,omitempty"`
	Policy             string                 `json:"policy,omitempty"`
	WorkspaceURL       string                 `json:"workspace,omitempty"`
	TemplateProperties map[string]interface{} `json:"templateProperties,omitempty"`
}

type DNSReservation struct {
	Links *struct {
		Self        LinkRef `json:"self,omitempty"`
		Workspace   LinkRef `json:"workspace,omitempty"`
		Policy      LinkRef `json:"policy,omitempty"`
		JobMetadata LinkRef `json:"jobMetadata,omitempty"`
	} `json:"_links,omitempty"`
	ID                 int                    `json:"id,omitempty"`
	Name               string                 `json:"name,omitempty"`
	PolicyID           int                    `json:"policyId,omitempty"`
	Policy             string                 `json:"policy,omitempty"`
	WorkspaceURL       string                 `json:"workspace,omitempty"`
	Value              string                 `json:"value,omitempty"`
	Zones              []string               `json:"zones,omitempty"`
	TemplateProperties map[string]interface{} `json:"templateProperties,omitempty"`
}

type IPAMReservation struct {
	Links *struct {
		Self        LinkRef `json:"self,omitempty"`
		Workspace   LinkRef `json:"workspace,omitempty"`
		Policy      LinkRef `json:"policy,omitempty"`
		JobMetadata LinkRef `json:"jobMetadata,omitempty"`
	} `json:"_links,omitempty"`
	ID                 int                    `json:"id,omitempty"`
	Hostname           string                 `json:"hostname,omitempty"`
	PolicyID           int                    `json:"policyId,omitempty"`
	Policy             string                 `json:"policy,omitempty"`
	WorkspaceURL       string                 `json:"workspace,omitempty"`
	IPaddress          string                 `json:"ipAddress,omitempty"`
	Gateway            string                 `json:"gateway,omitempty"`
	PrimaryDNS         string                 `json:"primaryDns"`
	SecondaryDNS       string                 `json:"secondaryDns"`
	Network            string                 `json:"network,omitempty"`
	DNSSuffix          string                 `json:"dnsSuffix,omitempty"`
	Netmask            string                 `json:"netmask,omitempty"`
	TemplateProperties map[string]interface{} `json:"template_properties,omitempty"`
}

type StaticPropertySetResponse struct {
	Embedded struct {
		PropertySets []StaticPropertySet `json:"propertySets"`
	} `json:"_embedded"`
}

type StaticPropertySet struct {
	Links *struct {
		Self      LinkRef `json:"self,omitempty"`
		Workspace LinkRef `json:"workspace,omitempty"`
	} `json:"_links,omitempty"`
	ID          int                    `json:"id,omitempty"`
	Type        string                 `json:"type,omitempty"`
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
}

type IPAMPolicyResponse struct {
	Embedded struct {
		IPAMPolicies []IPAMPolicy `json:"ipamPolicies"`
	} `json:"_embedded"`
}

type IPAMPolicy struct {
	Links *struct {
		Self      LinkRef `json:"self,omitempty"`
		Workspace LinkRef `json:"workspace,omitempty"`
	} `json:"_links,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type NamingPolicyResponse struct {
	Embedded struct {
		NamingPolicies []NamingPolicy `json:"namingPolicies"`
	} `json:"_embedded"`
}

type NamingPolicy struct {
	Links *struct {
		Self      LinkRef `json:"self,omitempty"`
		Workspace LinkRef `json:"workspace,omitempty"`
	} `json:"_links,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ADPolicyResponse struct {
	Embedded struct {
		ADPolicies []ADPolicy `json:"microsoftADPolicies"`
	} `json:"_embedded"`
}

type ADPolicy struct {
	Links *struct {
		Self      LinkRef `json:"self,omitempty"`
		Workspace LinkRef `json:"workspace,omitempty"`
	} `json:"_links,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type DNSPolicyResponse struct {
	Embedded struct {
		DNSPolicies []DNSPolicy `json:"dnsPolicies"`
	} `json:"_embedded"`
}

type DNSPolicy struct {
	Links *struct {
		Self      LinkRef `json:"self,omitempty"`
		Workspace LinkRef `json:"workspace,omitempty"`
	} `json:"_links,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type JobStatus struct {
	Links *struct {
		Self          LinkRef `json:"self,omitempty"`
		JobMetadata   LinkRef `json:"jobMetadata,omitempty"`
		ManagedObject LinkRef `json:managedObject,omitempty"`
		Policy        LinkRef `json:"policy,omitempty"`
		Workspace     LinkRef `json:"workspace,omitempty"`
	} `json:"_links,omitempty"`
	ID                  int    `json:"id,omitempty"`
	JobStateDescription string `json:"jobStateDescription,omitempty"`
	JobState            string `json:"jobState,omitempty"`
	JobTrackingID       string `json:"jobTrackingId,omitempty"`
	JobType             string `json:"jobType,omitempty"`
	ErrorDetails        *struct {
		Code   int `json:"code,omitempty"`
		Errors *[]struct {
			Message string `json:"message,omitempty"`
		} `json:"errors,omitempty"`
	} `json:errorDetails,omitempty"`
}

func (c *Config) NewOneFuseApiClient() *OneFuseAPIClient {
	return &OneFuseAPIClient{
		config: c,
	}
}

func (apiClient *OneFuseAPIClient) GenerateCustomName(namingPolicyID string, workspaceID string, templateProperties map[string]interface{}) (*CustomName, error) {
	log.Println("onefuse.apiClient: GenerateCustomName")

	config := apiClient.config
	url := collectionURL(config, NamingResourceType)

	if templateProperties == nil {
		templateProperties = make(map[string]interface{})
	}

	if workspaceID == "" {
		defaultWorkspaceID, err := findDefaultWorkspaceID(config)
		if err != nil {
			return nil, err
		}
		workspaceID = defaultWorkspaceID
	}

	postBody := map[string]interface{}{
		"policy":             fmt.Sprintf("/%s/%s/namingPolicies/%s/", ApiVersion, ApiNamespace, namingPolicyID),
		"templateProperties": templateProperties,
		"workspace":          fmt.Sprintf("/%s/%s/workspaces/%s/", ApiVersion, ApiNamespace, workspaceID),
	}

	jsonBytes, err := json.Marshal(postBody)
	if err != nil {
		return nil, errors.WithMessage(err, "onefuse.apiClient: Unable to marshal request body to JSON")
	}

	requestBody := string(jsonBytes)

	payload := strings.NewReader(requestBody)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Unable to create request POST %s %s", url, requestBody))
	}

	setHeaders(req, config)

	customName := CustomName{}
	var jobStatus *JobStatus
	jobStatus, err = handleAsyncRequestAndFetchManagdObject(req, config, &customName, "POST")
	if err != nil {
		return nil, err
	}

	if err = checkForJobErrors(jobStatus); err != nil {
		return nil, err
	}

	return &customName, nil
}

func (apiClient *OneFuseAPIClient) GetCustomName(id int) (*CustomName, error) {
	log.Println("onefuse.apiClient: GetCustomName")

	config := apiClient.config

	url := itemURL(config, NamingResourceType, id)
	customName := CustomName{}

	err := doGet(config, url, &customName)
	if err != nil {
		return nil, err
	}

	return &customName, nil
}

func (apiClient *OneFuseAPIClient) DeleteCustomName(id int) error {
	log.Println("onefuse.apiClient: DeleteCustomName")

	config := apiClient.config

	url := itemURL(config, NamingResourceType, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request DELETE %s", url))
	}

	setHeaders(req, config)

	if _, err = handleAsyncRequest(req, config, "DELETE"); err != nil {
		return err
	}

	return nil
}

func (apiClient *OneFuseAPIClient) CreateMicrosoftEndpoint(newEndpoint MicrosoftEndpoint) (*MicrosoftEndpoint, error) {
	log.Println("onefuse.apiClient: CreateMicrosoftEndpoint")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetMicrosoftEndpoint(id int) (*MicrosoftEndpoint, error) {
	log.Println("onefuse.apiClient: GetMicrosoftEndpoint")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetMicrosoftEndpointByName(name string) (*MicrosoftEndpoint, error) {
	log.Println("onefuse.apiClient: GetMicrosoftEndpointByName")

	config := apiClient.config
	url := fmt.Sprintf("%s?filter=name:%s;type:microsoft", collectionURL(config, ModuleEndpointResourceType), name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request GET %s", url))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request GET %s", url))
	}

	if err = checkForErrors(res); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Request failed GET %s", url))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from GET %s", url))
	}
	defer res.Body.Close()

	endpoints := EndpointsListResponse{}
	err = json.Unmarshal(body, &endpoints)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	if len(endpoints.Embedded.Endpoints) < 1 {
		return nil, errors.New(fmt.Sprintf("onefuse.apiClient: Could not find Microsoft Endpoint '%s'!", name))
	}

	endpoint := endpoints.Embedded.Endpoints[0]

	return &endpoint, err
}

func (apiClient *OneFuseAPIClient) UpdateMicrosoftEndpoint(id int, updatedEndpoint MicrosoftEndpoint) (*MicrosoftEndpoint, error) {
	log.Println("onefuse.apiClient: UpdateMicrosoftEndpoint")

	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) DeleteMicrosoftEndpoint(id int) error {
	log.Println("onefuse.apiClient: DeleteMicrosoftEndpoint")

	return errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) CreateMicrosoftADPolicy(newPolicy *MicrosoftADPolicy) (*MicrosoftADPolicy, error) {
	log.Println("onefuse.apiClient: CreateMicrosoftADPolicy")

	config := apiClient.config

	// Default workspace if it was not provided
	if newPolicy.WorkspaceURL == "" {
		workspaceID, err := findDefaultWorkspaceID(config)
		if err != nil {
			return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to find default workspace")
		}
		workspaceIDInt, err := strconv.Atoi(workspaceID)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to convert Workspace ID '%s' to integer", workspaceID))
		}

		newPolicy.WorkspaceURL = itemURL(config, WorkspaceResourceType, workspaceIDInt)
	}

	// Construct a URL we are going to POST to
	// /api/v3/onefuse/microsoftADPolicies/
	url := collectionURL(config, MicrosoftADPolicyResourceType)

	jsonBytes, err := json.Marshal(newPolicy)
	if err != nil {
		return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to marshal request body to JSON")
	}

	requestBody := string(jsonBytes)
	payload := strings.NewReader(requestBody)

	// Create the create request
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Unable to create request POST %s %s", url, requestBody))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request POST %s %s", url, requestBody))
	}

	if err = checkForErrors(res); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Request failed POST %s %s", url, requestBody))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from POST %s %s", url, requestBody))
	}
	defer res.Body.Close()

	policy := MicrosoftADPolicy{}
	if err = json.Unmarshal(body, &policy); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	return &policy, nil
}

func (apiClient *OneFuseAPIClient) GetMicrosoftADPolicy(id int) (*MicrosoftADPolicy, error) {
	log.Println("onefuse.apiClient: GetMicrosoftADPolicy")

	config := apiClient.config

	url := itemURL(config, MicrosoftADPolicyResourceType, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request GET %s", url))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request GET %s", url))
	}

	err = checkForErrors(res)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Error from request GET %s %s", url, err))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from GET %s", url))
	}
	defer res.Body.Close()

	policy := MicrosoftADPolicy{}
	if err = json.Unmarshal(body, &policy); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	return &policy, err
}

func (apiClient *OneFuseAPIClient) UpdateMicrosoftADPolicy(id int, updatedPolicy *MicrosoftADPolicy) (*MicrosoftADPolicy, error) {
	log.Println("onefuse.apiClient: UpdateMicrosoftADPolicy")

	config := apiClient.config

	url := itemURL(config, MicrosoftADPolicyResourceType, id)

	if updatedPolicy.Name == "" {
		return nil, errors.New("onefuse.apiClient: Microsoft AD Policy Updates Require a Name")
	}

	if updatedPolicy.WorkspaceURL == "" {
		workspaceID, err := findDefaultWorkspaceID(config)
		if err != nil {
			return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to find default workspace")
		}
		workspaceIDInt, err := strconv.Atoi(workspaceID)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to convert Workspace ID '%s' to integer", workspaceID))
		}

		updatedPolicy.WorkspaceURL = itemURL(config, WorkspaceResourceType, workspaceIDInt)
	}

	jsonBytes, err := json.Marshal(updatedPolicy)
	if err != nil {
		return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to marshal request body to JSON")
	}

	requestBody := string(jsonBytes)
	payload := strings.NewReader(requestBody)

	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request PUT %s %s", url, requestBody))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request PUT %s %s", url, requestBody))
	}

	err = checkForErrors(res)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Error from request PUT %s %s", url, err))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from PUT %s %s", url, requestBody))
	}
	defer res.Body.Close()

	policy := MicrosoftADPolicy{}
	if err = json.Unmarshal(body, &policy); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarhsal response %s", string(body)))
	}

	return &policy, nil
}

func (apiClient *OneFuseAPIClient) DeleteMicrosoftADPolicy(id int) error {
	log.Println("onefuse.apiClient: DeleteMicrosoftADPolicy")

	config := apiClient.config

	url := itemURL(config, MicrosoftADPolicyResourceType, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request DELETE %s", url))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request DELETE %s", url))
	}

	return checkForErrors(res)
}

func (apiClient *OneFuseAPIClient) CreateMicrosoftADComputerAccount(newComputerAccount *MicrosoftADComputerAccount) (*MicrosoftADComputerAccount, error) {
	log.Println("onefuse.apiClient: CreateMicrosoftADComputerAccount")

	config := apiClient.config

	// Default workspace if it was not provided
	if newComputerAccount.WorkspaceURL == "" {
		workspaceID, err := findDefaultWorkspaceID(config)
		if err != nil {
			return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to find default workspace")
		}
		workspaceIDInt, err := strconv.Atoi(workspaceID)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to convert Workspace ID '%s' to integer", workspaceID))
		}

		newComputerAccount.WorkspaceURL = itemURL(config, WorkspaceResourceType, workspaceIDInt)
	}

	if newComputerAccount.Policy == "" {
		if newComputerAccount.PolicyID != 0 {
			newComputerAccount.Policy = itemURL(config, WorkspaceResourceType, newComputerAccount.PolicyID)
		} else {
			return nil, errors.New("onefuse.apiClient: Microsoft AD Computer Account Create requires a PolicyID or Policy URL")
		}
	} else {
		return nil, errors.New("onefuse.apiClient: Microsoft AD Computer Account Create requires a PolicyID or Policy URL")
	}

	// Construct a URL we are going to POST to
	// /api/v3/onefuse/microsoftADComputerAccounts/
	url := collectionURL(config, MicrosoftADComputerAccountResourceType)

	jsonBytes, err := json.Marshal(newComputerAccount)
	if err != nil {
		return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to marshal request body to JSON")
	}

	requestBody := string(jsonBytes)
	payload := strings.NewReader(requestBody)

	// Create the create request
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Unable to create request POST %s %s", url, requestBody))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	// Make the create request
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request POST %s %s", url, requestBody))
	}

	if err = checkForErrors(res); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Request failed POST %s %s", url, requestBody))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from POST %s %s", url, requestBody))
	}
	defer res.Body.Close()

	computerAccount := MicrosoftADComputerAccount{}
	if err = json.Unmarshal(body, &computerAccount); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	return &computerAccount, nil
}

func (apiClient *OneFuseAPIClient) GetMicrosoftADComputerAccount(id int) (*MicrosoftADComputerAccount, error) {
	log.Println("onefuse.apiClient: GetMicrosoftADComputerAccount")

	config := apiClient.config

	url := itemURL(config, MicrosoftADComputerAccountResourceType, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request GET %s %s", url, err))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request GET %s %s", url, err))
	}

	if err = checkForErrors(res); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Error from request GET %s %s", url, err))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from GET %s %s", url, err))
	}
	defer res.Body.Close()

	computerAccount := MicrosoftADComputerAccount{}
	if err = json.Unmarshal(body, &computerAccount); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	return &computerAccount, err
}

func (apiClient *OneFuseAPIClient) UpdateMicrosoftADComputerAccount(id int, updatedComputerAccount *MicrosoftADComputerAccount) (*MicrosoftADComputerAccount, error) {
	log.Println("onefuse.apiClient: UpdateMicrosoftADComputerAccount")

	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) DeleteMicrosoftADComputerAccount(id int) error {
	log.Println("onefuse.apiClient: DeleteMicrosoftADComputerAccount")

	config := apiClient.config

	url := itemURL(config, MicrosoftADComputerAccountResourceType, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request DELETE %s", url))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request DELETE %s", url))
	}

	return checkForErrors(res)
}

//DNS Functions

//Create DNS Reservation

func (apiClient *OneFuseAPIClient) CreateDNSReservation(newDNSRecord *DNSReservation) (*DNSReservation, error) {
	log.Println("onefuse.apiClient: CreateDNSReservation")

	config := apiClient.config

	// Default workspace if it was not provided
	if newDNSRecord.WorkspaceURL == "" {
		workspaceID, err := findDefaultWorkspaceID(config)
		if err != nil {
			return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to find default workspace")
		}
		workspaceIDInt, err := strconv.Atoi(workspaceID)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to convert Workspace ID '%s' to integer", workspaceID))
		}

		newDNSRecord.WorkspaceURL = itemURL(config, WorkspaceResourceType, workspaceIDInt)
	}

	if newDNSRecord.Policy == "" {
		if newDNSRecord.PolicyID != 0 {
			newDNSRecord.Policy = itemURL(config, WorkspaceResourceType, newDNSRecord.PolicyID)
		} else {
			return nil, errors.New("onefuse.apiClient: DNS Record Create requires a PolicyID or Policy URL")
		}
	} else {
		return nil, errors.New("onefuse.apiClient: DNS Record Create requires a PolicyID or Policy URL")
	}

	// Construct a URL we are going to POST to
	// /api/v3/onefuse/dnsReservations/
	url := collectionURL(config, DNSReservationResourceType)

	jsonBytes, err := json.Marshal(newDNSRecord)
	if err != nil {
		return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to marshal request body to JSON")
	}

	requestBody := string(jsonBytes)
	payload := strings.NewReader(requestBody)

	// Create the create request
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Unable to create request POST %s %s", url, requestBody))
	}

	setHeaders(req, config)

	dnsRecord := DNSReservation{}
	var jobStatus *JobStatus

	jobStatus, err = handleAsyncRequestAndFetchManagdObject(req, config, &dnsRecord, "POST")
	if err != nil {
		return nil, err
	}

	return &dnsRecord, nil
}

//Get DNS Reservation

func (apiClient *OneFuseAPIClient) GetDNSReservation(id int) (*DNSReservation, error) {
	log.Println("onefuse.apiClient: GetDNSReservation")

	config := apiClient.config

	url := itemURL(config, DNSReservationResourceType, id)

	dnsRecord := DNSReservation{}

	err := doGet(config, url, &dnsRecord)
	if err != nil {
		return nil, err
	}

	return &dnsRecord, err
}

//Update DNS Record

func (apiClient *OneFuseAPIClient) UpdateDNSReservation(id int, updatedDNSReservation *DNSReservation) (*DNSReservation, error) {
	log.Println("onefuse.apiClient: UpdateDNSReservation")

	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) DeleteDNSReservation(id int) error {
	log.Println("onefuse.apiClient: DeleteDNSReservation")

	config := apiClient.config

	url := itemURL(config, DNSReservationResourceType, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request DELETE %s", url))
	}

	setHeaders(req, config)

	if _, err = handleAsyncRequest(req, config, "DELETE"); err != nil {
		return err
	}
	return nil
}

//Create IPAM Reservation

func (apiClient *OneFuseAPIClient) CreateIPAMReservation(newIPAMRecord *IPAMReservation) (*IPAMReservation, error) {
	log.Println("onefuse.apiClient: CreateIPAMReservation")

	config := apiClient.config

	// Default workspace if it was not provided
	if newIPAMRecord.WorkspaceURL == "" {
		workspaceID, err := findDefaultWorkspaceID(config)
		if err != nil {
			return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to find default workspace")
		}
		workspaceIDInt, err := strconv.Atoi(workspaceID)
		if err != nil {
			return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to convert Workspace ID '%s' to integer", workspaceID))
		}

		newIPAMRecord.WorkspaceURL = itemURL(config, WorkspaceResourceType, workspaceIDInt)
	}

	if newIPAMRecord.Policy == "" {
		if newIPAMRecord.PolicyID != 0 {
			newIPAMRecord.Policy = itemURL(config, WorkspaceResourceType, newIPAMRecord.PolicyID)
		} else {
			return nil, errors.New("onefuse.apiClient: IPAM Record Create requires a PolicyID or Policy URL")
		}
	} else {
		return nil, errors.New("onefuse.apiClient: IPAM Record Create requires a PolicyID or Policy URL")
	}

	// Construct a URL we are going to POST to
	// /api/v3/onefuse/ipamReservations/
	url := collectionURL(config, IPAMReservationResourceType)

	jsonBytes, err := json.Marshal(newIPAMRecord)
	if err != nil {
		return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to marshal request body to JSON")
	}

	requestBody := string(jsonBytes)
	payload := strings.NewReader(requestBody)

	// Create the create request
	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Unable to create request POST %s %s", url, requestBody))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	// Make the create request
	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request POST %s %s", url, requestBody))
	}

	if err = checkForErrors(res); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Request failed POST %s %s", url, requestBody))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from POST %s %s", url, requestBody))
	}
	defer res.Body.Close()

	ipamRecord := IPAMReservation{}
	if err = json.Unmarshal(body, &ipamRecord); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	return &ipamRecord, nil
}

//Get IPAM Reservation

func (apiClient *OneFuseAPIClient) GetIPAMReservation(id int) (*IPAMReservation, error) {
	log.Println("onefuse.apiClient: GetIPAMReservation")

	config := apiClient.config

	url := itemURL(config, IPAMReservationResourceType, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request GET %s %s", url, err))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request GET %s %s", url, err))
	}

	if err = checkForErrors(res); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Error from request GET %s %s", url, err))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from GET %s %s", url, err))
	}
	defer res.Body.Close()

	ipamRecord := IPAMReservation{}
	if err = json.Unmarshal(body, &ipamRecord); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	return &ipamRecord, err
}

//Update IPAM Record

func (apiClient *OneFuseAPIClient) UpdateIPAMReservation(id int, updatedIPAMReservation *IPAMReservation) (*IPAMReservation, error) {
	log.Println("onefuse.apiClient: UpdateIPAMReservation")

	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) DeleteIPAMReservation(id int) error {
	log.Println("onefuse.apiClient: DeleteIPAMReservation")

	config := apiClient.config

	url := itemURL(config, IPAMReservationResourceType, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request DELETE %s", url))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request DELETE %s", url))
	}

	return checkForErrors(res)
}

// End IPAM

// Start IPAM Policies

func (apiClient *OneFuseAPIClient) GetIPAMPolicy(id int) (*IPAMPolicy, error) {
	log.Println("onefuse.apiClient: IPAMPolicy")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetIPAMPolicyByName(name string) (*IPAMPolicy, error) {
	log.Println("onefuse.apiClient: GetIPAMPolicyByName")

	config := apiClient.config
	url := fmt.Sprintf("%s?filter=name:%s", collectionURL(config, IPAMPolicyResourceType), name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request GET %s", url))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request GET %s", url))
	}

	if err = checkForErrors(res); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Request failed GET %s", url))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from GET %s", url))
	}
	defer res.Body.Close()

	ipamPolicies := IPAMPolicyResponse{}
	err = json.Unmarshal(body, &ipamPolicies)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	if len(ipamPolicies.Embedded.IPAMPolicies) < 1 {
		return nil, errors.New(fmt.Sprintf("onefuse.apiClient: Could not find IPAM Policy '%s'!", name))
	}

	ipamPolicy := ipamPolicies.Embedded.IPAMPolicies[0]

	return &ipamPolicy, err
}

// End IPAM Policies

// Start Naming Policies

func (apiClient *OneFuseAPIClient) GetNamingPolicy(id int) (*NamingPolicy, error) {
	log.Println("onefuse.apiClient: NamingPolicy")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetNamingPolicyByName(name string) (*NamingPolicy, error) {
	log.Println("onefuse.apiClient: GetNamingPolicyByName")

	config := apiClient.config
	url := fmt.Sprintf("%s?filter=name:%s", collectionURL(config, NamingPolicyResourceType), name)

	namingPolicies := NamingPolicyResponse{}
	err := doGet(config, url, &namingPolicies)
	if err != nil {
		return nil, err
	}

	if len(namingPolicies.Embedded.NamingPolicies) < 1 {
		return nil, errors.New(fmt.Sprintf("onefuse.apiClient: Could not find Naming Policy '%s'!", name))
	}

	namingPolicy := namingPolicies.Embedded.NamingPolicies[0]

	return &namingPolicy, err
}

// End Naming Policies

// Start AD Policies

func (apiClient *OneFuseAPIClient) GetADPolicy(id int) (*ADPolicy, error) {
	log.Println("onefuse.apiClient: ADPolicy")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetADPolicyByName(name string) (*ADPolicy, error) {
	log.Println("onefuse.apiClient: GetADPolicyByName")

	config := apiClient.config
	url := fmt.Sprintf("%s?filter=name:%s", collectionURL(config, ADPolicyResourceType), name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request GET %s", url))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request GET %s", url))
	}

	if err = checkForErrors(res); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Request failed GET %s", url))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from GET %s", url))
	}
	defer res.Body.Close()

	adPolicies := ADPolicyResponse{}
	err = json.Unmarshal(body, &adPolicies)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	if len(adPolicies.Embedded.ADPolicies) < 1 {
		return nil, errors.New(fmt.Sprintf("onefuse.apiClient: Could not find AD Policy '%s'!", name))
	}

	adPolicy := adPolicies.Embedded.ADPolicies[0]

	return &adPolicy, err
}

// End AD Policies

// Start DNS Policies

func (apiClient *OneFuseAPIClient) GetDNSPolicy(id int) (*DNSPolicy, error) {
	log.Println("onefuse.apiClient: DNSPolicy")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetDNSPolicyByName(name string) (*DNSPolicy, error) {
	log.Println("onefuse.apiClient: GetDNSPolicyByName")

	config := apiClient.config
	url := fmt.Sprintf("%s?filter=name:%s", collectionURL(config, DNSPolicyResourceType), name)

	dnsPolicies := DNSPolicyResponse{}
	err := doGet(config, url, &dnsPolicies)
	if err != nil {
		return nil, err
	}

	if len(dnsPolicies.Embedded.DNSPolicies) < 1 {
		return nil, errors.New(fmt.Sprintf("onefuse.apiClient: Could not find AD Policy '%s'!", name))
	}

	dnsPolicy := dnsPolicies.Embedded.DNSPolicies[0]

	return &dnsPolicy, err
}

// End DNS Policies

// Start Static Property Set

func (apiClient *OneFuseAPIClient) GetStaticPropertySet(id int) (*StaticPropertySet, error) {
	log.Println("onefuse.apiClient: GetStaticPropertySet")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetStaticPropertySetByName(name string) (*StaticPropertySet, error) {
	log.Println("onefuse.apiClient: GetStaticPropertySetByName")

	config := apiClient.config
	url := fmt.Sprintf("%s?filter=name:%s;type:static", collectionURL(config, StaticPropertySetResourceType), name)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request GET %s", url))
	}

	setHeaders(req, config)

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request GET %s", url))
	}

	if err = checkForErrors(res); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Request failed GET %s", url))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from GET %s", url))
	}
	defer res.Body.Close()

	staticPropertySets := StaticPropertySetResponse{}
	err = json.Unmarshal(body, &staticPropertySets)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	if len(staticPropertySets.Embedded.PropertySets) < 1 {
		return nil, errors.New(fmt.Sprintf("onefuse.apiClient: Could not find staticPropertySet '%s'!", name))
	}

	staticPropertySet := staticPropertySets.Embedded.PropertySets[0]

	return &staticPropertySet, err
}

// End Static Property Set

func GetJobStatus(id int, config *Config) (*JobStatus, error) {
	log.Println("onefuse.apiClient: GetJobStatus")

	url := itemURL(config, JobStatusResourceType, id)
	result := JobStatus{}

	err := doGet(config, url, &result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func handleAsyncRequestAndFetchManagdObject(req *http.Request, config *Config, responseObject interface{}, httpVerb string) (jobStatus *JobStatus, err error) {

	if jobStatus, err = handleAsyncRequest(req, config, httpVerb); err != nil {
		return
	}

	url := urlFromHref(config, jobStatus.Links.ManagedObject.Href)
	err = doGet(config, url, &responseObject)
	if err != nil {
		return nil, err
	}

	return jobStatus, nil
}

func handleAsyncRequest(req *http.Request, config *Config, httpVerb string) (jobStatus *JobStatus, err error) {

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		body, _ := ioutil.ReadAll(req.Body)
		return jobStatus, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request %s %s %s", httpVerb, req.URL, body))
	}

	body, err := readResponse(res)
	if err != nil {
		body, _ := ioutil.ReadAll(req.Body)
		return jobStatus, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from %s %s %s", httpVerb, req.URL, body))
	}
	defer res.Body.Close()

	if err = json.Unmarshal(body, &jobStatus); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	jobStatus, err = waitForJob(jobStatus.ID, config)
	if err != nil {
		return
	}

	if err = checkForJobErrors(jobStatus); err != nil {
		return nil, err
	}

	return jobStatus, nil
}

func doGet(config *Config, url string, v interface{}) (err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request GET %s", url))
	}

	setHeaders(req, config)

	client := getHttpClient(config)
	res, err := client.Do(req)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request GET %s", url))
	}

	if err = checkForErrors(res); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Request failed GET %s", url))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from GET %s", url))
	}
	defer res.Body.Close()

	if err = json.Unmarshal(body, &v); err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	return nil
}

func waitForJob(jobID int, config *Config) (jobStatus *JobStatus, err error) {
	jobStatusDescription := ""
	PollingTimeoutMS := 60000
	PollingIntervalMS := 2000
	startTime := time.Now()
	for jobStatusDescription != JobSuccess && jobStatusDescription != JobFailed {
		jobStatus, err = GetJobStatus(jobID, config)
		if err != nil {
			return nil, err
		}

		jobStatusDescription = jobStatus.JobState
		log.Println(jobStatus)

		time.Sleep(time.Duration(PollingIntervalMS) * time.Millisecond)
		if time.Since(startTime) > (time.Duration(PollingTimeoutMS) * time.Millisecond) {
			panic("Timed out while waiting for job to complete.")
		}
	}
	return jobStatus, nil
}

func findDefaultWorkspaceID(config *Config) (workspaceID string, err error) {
	fmt.Println("onefuse.findDefaultWorkspaceID")

	filter := "filter=name.exact:Default"
	url := fmt.Sprintf("%s?%s", collectionURL(config, WorkspaceResourceType), filter)

	req, clientErr := http.NewRequest("GET", url, nil)
	if clientErr != nil {
		err = errors.WithMessage(clientErr, fmt.Sprintf("onefuse.findDefaultWorkspaceID: Failed to make request GET %s", url))
		return
	}

	setHeaders(req, config)

	client := getHttpClient(config)
	res, clientErr := client.Do(req)
	if clientErr != nil {
		err = errors.WithMessage(clientErr, fmt.Sprintf("onefuse.findDefaultWorkspaceID: Failed to do request GET %s", url))
		return
	}

	body, err := readResponse(res)
	if err != nil {
		return
	}
	defer res.Body.Close()

	var data WorkspacesListResponse
	json.Unmarshal(body, &data)

	workspaces := data.Embedded.Workspaces
	if len(workspaces) == 0 {
		err = errors.WithMessage(clientErr, "onefuse.findDefaultWorkspaceID: Failed to find default workspace!")
		return
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

func readResponse(res *http.Response) (bytes []byte, err error) {
	err = checkForErrors(res)
	if err != nil {
		return
	}

	bytes, err = ioutil.ReadAll(res.Body)
	return
}

func checkForErrors(res *http.Response) error {
	if res.StatusCode >= 500 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		return errors.New(string(b))
	} else if res.StatusCode >= 400 {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		defer res.Body.Close()
		return errors.New(string(b))
	}
	return nil
}

func checkForJobErrors(jobStatus *JobStatus) error {
	if jobStatus.JobState != JobSuccess {
		return errors.New(fmt.Sprintf("Job %s (%d) failed with message %v", jobStatus.JobType, jobStatus.ID, *jobStatus.ErrorDetails.Errors))
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
	req.Header.Add("Host", fmt.Sprintf("%s:%s", config.address, config.port))
	req.Header.Add("SOURCE", "Terraform")
	req.SetBasicAuth(config.user, config.password)
}

func collectionURL(config *Config, resourceType string) string {
	baseURL := fmt.Sprintf("%s://%s:%s", config.scheme, config.address, config.port)
	endpoint := path.Join(ApiVersion, ApiNamespace, resourceType)
	return fmt.Sprintf("%s/%s/", baseURL, endpoint)
}

func urlFromHref(config *Config, href string) string {
	return fmt.Sprintf("%s://%s:%s%s", config.scheme, config.address, config.port, href)
}

func itemURL(config *Config, resourceType string, id int) string {
	idString := strconv.Itoa(id)
	baseURL := collectionURL(config, resourceType)
	return fmt.Sprintf("%s%s/", baseURL, idString)
}
