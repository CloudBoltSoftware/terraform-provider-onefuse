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
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

const ApiVersion = "api/v3"
const ApiNamespace = "onefuse"
const AnsibleTowerDeploymentResourceType = "ansibleTowerDeployments"
const AnsibleTowerPolicyResourceType = "ansibleTowerPolicies"
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

type AnsibleTowerDeployment struct {
	Links *struct {
		Self        LinkRef `json:"self,omitempty"`
		Workspace   LinkRef `json:"workspace,omitempty"`
		Policy      LinkRef `json:"policy,omitempty"`
		JobMetadata LinkRef `json:"jobMetadata,omitempty"`
	} `json:"_links,omitempty"`
	ID                     int      `json:"id,omitempty"`
	PolicyID               int      `json:"policyId,omitempty"`
	Policy                 string   `json:"policy,omitempty"`
	WorkspaceURL           string   `json:"workspace,omitempty"`
	Limit                  string   `json:"limit,omitempty"`
	InventoryName          string   `json:"inventoryName,omitempty"`
	Hosts                  []string `json:"hosts,omitempty"`
	Archived               bool     `json:"archived,omitempty"`
	ProvisioningJobResults []struct {
		output          string `json:"output"`
		status          string `json:"status"`
		jobTemplateName string `json:"jobTemplateName"`
	} `json:"provisioningJobResults,omitempty"`
	DeprovisioningJobResults []struct {
		output          string `json:"output"`
		status          string `json:"status"`
		jobTemplateName string `json:"jobTemplateName"`
	} `json:"deprovisioningJobResults,omitempty"`
	TemplateProperties map[string]interface{} `json:"template_properties,omitempty"`
}

func (c *Config) NewOneFuseApiClient() *OneFuseAPIClient {
	return &OneFuseAPIClient{
		config: c,
	}
}

func (apiClient *OneFuseAPIClient) GenerateCustomName(namingPolicyID string, workspaceID string, templateProperties map[string]interface{}) (*CustomName, error) {
	log.Println("onefuse.apiClient: GenerateCustomName")

	config := apiClient.config

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

	var req *http.Request
	var err error
	if req, err = buildPostRequest(config, NamingResourceType, postBody); err != nil {
		return nil, err
	}

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

	endpoints := EndpointsListResponse{}
	entity, err := findEntityByName(config, name, ModuleEndpointResourceType, &endpoints, "Endpoints", ";type:microsoft")
	if err != nil {
		return nil, err
	}
	endpoint := entity.(MicrosoftEndpoint)
	return &endpoint, nil
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

	var err error
	if newPolicy.WorkspaceURL, err = findWorkspaceURLOrDefault(config, newPolicy.WorkspaceURL); err != nil {
		return nil, err
	}

	var req *http.Request
	if req, err = buildPostRequest(config, MicrosoftADPolicyResourceType, newPolicy); err != nil {
		return nil, err
	}

	client := getHttpClient(config)

	res, err := client.Do(req)
	if err != nil {
		body, _ := ioutil.ReadAll(req.Body)
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to do request POST %s %s", req.URL, body))
	}

	if err = checkForErrors(res); err != nil {
		body, _ := ioutil.ReadAll(req.Body)
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Request failed POST %s %s", req.URL, body))
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		body, _ := ioutil.ReadAll(req.Body)
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to read response body from POST %s %s", req.URL, body))
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

	policy := MicrosoftADPolicy{}
	err := doGet(config, url, &policy)
	if err != nil {
		return nil, err
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

	var err error
	if updatedPolicy.WorkspaceURL, err = findWorkspaceURLOrDefault(config, updatedPolicy.WorkspaceURL); err != nil {
		return nil, err
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

	var err error
	if newComputerAccount.WorkspaceURL, err = findWorkspaceURLOrDefault(config, newComputerAccount.WorkspaceURL); err != nil {
		return nil, err
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

	var req *http.Request
	if req, err = buildPostRequest(config, MicrosoftADComputerAccountResourceType, newComputerAccount); err != nil {
		return nil, err
	}

	computerAccount := MicrosoftADComputerAccount{}
	_, err = handleAsyncRequestAndFetchManagdObject(req, config, &computerAccount, "POST")
	if err != nil {
		return nil, err
	}

	return &computerAccount, nil
}

func (apiClient *OneFuseAPIClient) GetMicrosoftADComputerAccount(id int) (*MicrosoftADComputerAccount, error) {
	log.Println("onefuse.apiClient: GetMicrosoftADComputerAccount")

	config := apiClient.config

	url := itemURL(config, MicrosoftADComputerAccountResourceType, id)

	computerAccount := MicrosoftADComputerAccount{}
	err := doGet(config, url, &computerAccount)
	if err != nil {
		return nil, err
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

	_, err = handleAsyncRequest(req, config, "DELETE")
	return err
}

//DNS Functions

//Create DNS Reservation

func (apiClient *OneFuseAPIClient) CreateDNSReservation(newDNSRecord *DNSReservation) (*DNSReservation, error) {
	log.Println("onefuse.apiClient: CreateDNSReservation")

	config := apiClient.config

	var err error
	if newDNSRecord.WorkspaceURL, err = findWorkspaceURLOrDefault(config, newDNSRecord.WorkspaceURL); err != nil {
		return nil, err
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

	var req *http.Request
	if req, err = buildPostRequest(config, DNSReservationResourceType, newDNSRecord); err != nil {
		return nil, err
	}

	dnsRecord := DNSReservation{}

	_, err = handleAsyncRequestAndFetchManagdObject(req, config, &dnsRecord, "POST")
	if err != nil {
		return nil, err
	}

	return &dnsRecord, nil
}

func buildPostRequest(config *Config, resourceType string, requestEntity interface{}) (*http.Request, error) {
	url := collectionURL(config, resourceType)

	jsonBytes, err := json.Marshal(requestEntity)
	if err != nil {
		return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to marshal request body to JSON")
	}

	requestBody := string(jsonBytes)
	payload := strings.NewReader(requestBody)

	req, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Unable to create request POST %s %s", url, requestBody))
	}

	setHeaders(req, config)

	return req, nil
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

	_, err = handleAsyncRequest(req, config, "DELETE")
	return err
}

//Create IPAM Reservation

func (apiClient *OneFuseAPIClient) CreateIPAMReservation(newIPAMRecord *IPAMReservation) (*IPAMReservation, error) {
	log.Println("onefuse.apiClient: CreateIPAMReservation")

	config := apiClient.config

	var err error
	if newIPAMRecord.WorkspaceURL, err = findWorkspaceURLOrDefault(config, newIPAMRecord.WorkspaceURL); err != nil {
		return nil, err
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

	var req *http.Request
	if req, err = buildPostRequest(config, IPAMReservationResourceType, newIPAMRecord); err != nil {
		return nil, err
	}

	ipamRecord := IPAMReservation{}

	_, err = handleAsyncRequestAndFetchManagdObject(req, config, &ipamRecord, "POST")
	if err != nil {
		return nil, err
	}
	return &ipamRecord, nil
}

//Get IPAM Reservation

func (apiClient *OneFuseAPIClient) GetIPAMReservation(id int) (*IPAMReservation, error) {
	log.Println("onefuse.apiClient: GetIPAMReservation")

	config := apiClient.config

	url := itemURL(config, IPAMReservationResourceType, id)

	ipamRecord := IPAMReservation{}
	err := doGet(config, url, &ipamRecord)
	if err != nil {
		return nil, err
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

	_, err = handleAsyncRequest(req, config, "DELETE")
	return err
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

	ipamPolicies := IPAMPolicyResponse{}
	entity, err := findEntityByName(config, name, IPAMPolicyResourceType, &ipamPolicies, "IPAMPolicies", "")
	if err != nil {
		return nil, err
	}
	ipamPolicy := entity.(IPAMPolicy)
	return &ipamPolicy, nil
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

	namingPolicies := NamingPolicyResponse{}
	entity, err := findEntityByName(config, name, NamingPolicyResourceType, &namingPolicies, "NamingPolicies", "")
	if err != nil {
		return nil, err
	}
	namingPolicy := entity.(NamingPolicy)
	return &namingPolicy, nil
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

	adPolicies := ADPolicyResponse{}
	entity, err := findEntityByName(config, name, ADPolicyResourceType, &adPolicies, "ADPolicies", "")
	if err != nil {
		return nil, err
	}
	adPolicy := entity.(ADPolicy)
	return &adPolicy, nil
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

	dnsPolicies := DNSPolicyResponse{}
	entity, err := findEntityByName(config, name, DNSPolicyResourceType, &dnsPolicies, "DNSPolicies", "")
	if err != nil {
		return nil, err
	}
	dnsPolicy := entity.(DNSPolicy)
	return &dnsPolicy, nil
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

	staticPropertySets := StaticPropertySetResponse{}
	entity, err := findEntityByName(config, name, StaticPropertySetResourceType, &staticPropertySets, "PropertySets", ";type:static")
	if err != nil {
		return nil, err
	}
	staticPropertySet := entity.(StaticPropertySet)
	return &staticPropertySet, nil
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
			return nil, errors.New("Timed out while waiting for job to complete.")
		}
	}
	return jobStatus, nil
}

func findWorkspaceURLOrDefault(config *Config, workspaceURL string) (string, error) {
	// Default workspace if it was not provided
	if workspaceURL == "" {
		workspaceID, err := findDefaultWorkspaceID(config)
		if err != nil {
			return "", errors.WithMessage(err, "onefuse.apiClient: Failed to find default workspace")
		}
		workspaceIDInt, err := strconv.Atoi(workspaceID)
		if err != nil {
			return "", errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to convert Workspace ID '%s' to integer", workspaceID))
		}

		workspaceURL = itemURL(config, WorkspaceResourceType, workspaceIDInt)
	}
	return workspaceURL, nil
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

// Finds an entity on OneFuse of type "resourceType" with name "name", using the supplied "collectionResponse" interface
// embeddedStructFieldName as the name of the embedded inner collection name.
// Additional filters to the collection will be appened to the name filter in the URL.
func findEntityByName(config *Config, name string, resourceType string, collectionResponse interface{},
	embeddedStructFieldName string, additionalFilters string) (interface{}, error) {

	url := fmt.Sprintf("%s?filter=name:%s%s", collectionURL(config, resourceType), name, additionalFilters)

	err := doGet(config, url, &collectionResponse)
	if err != nil {
		return nil, err
	}

	// TODO: Reflection logic below could likely use some better safeguards, but we can assume that OneFuse won't return a non-error Response
	// that doesn't follow HAL+JSON's embedded structure.
	embeddedField := reflect.Indirect(reflect.ValueOf(collectionResponse)).FieldByName("Embedded")
	embedded := embeddedField.Interface()

	collectionField := reflect.Indirect(reflect.ValueOf(embedded)).FieldByName(embeddedStructFieldName)

	if collectionField.Len() < 1 {
		return nil, errors.New(fmt.Sprintf("onefuse.apiClient: Could not find %s '%s'!", resourceType, name))
	}

	entity := collectionField.Index(0).Interface()

	return entity, err
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
