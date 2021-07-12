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
const ModulePolicyResourceType = "modulePolicies"
const ModuleDepoloymentResourceType = "moduleManagedObjects"
const DNSReservationResourceType = "dnsReservations"
const IPAMReservationResourceType = "ipamReservations"
const StaticPropertySetResourceType = "propertySets"
const RenderTemplateType = "templateTester"
const IPAMPolicyResourceType = "ipamPolicies"
const NamingPolicyResourceType = "namingPolicies"
const ADPolicyResourceType = "microsoftADPolicies"
const DNSPolicyResourceType = "dnsPolicies"
const ScriptingPolicyResourceType = "scriptingPolicies"
const JobStatusResourceType = "jobStatus"
const ScriptingDepoloymentResourceType = "scriptingDeployments"
const VraDeploymentResourceType = "vraDeployments"
const ServicenowCMDBPolicyResourceType = "servicenowCMDBPolicies"
const ServicenowCMDBDepoloymentResourceType = "servicenowCMDBDeployments"
const VraPolicyResourceType = "vraPolicies"

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
	TemplateProperties map[string]interface{} `json:"templateProperties"`
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
	TemplateProperties map[string]interface{} `json:"templateProperties"`
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
	Subnet             string                 `json:"subnet,omitempty"`
	DNSSuffix          string                 `json:"dnsSuffix,omitempty"`
	Netmask            string                 `json:"netmask,omitempty"`
	NicLabel           string                 `json:"nicLabel,omitempty"`
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
	Name        string                 `json:"name,omitempty"`
	Description string                 `json:"description,omitempty"`
	Properties  map[string]interface{} `json:"properties,omitempty"`
	Raw         string
}

type RenderTemplateResponse struct {
	Value string `json:"value,omitempty"`
}

type RenderTemplateRequest struct {
	Template           string                 `json:"template,omitempty"`
	TemplateProperties map[string]interface{} `json:"template_properties,omitempty"`
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

type ServicenowCMDBPolicyResponse struct {
	Embedded struct {
		ServicenowCMDBPolicies []ServicenowCMDBPolicy `json:"servicenowCMDBPolicies"`
	} `json:"_embedded"`
}

type ServicenowCMDBPolicy struct {
	Links *struct {
		Self      LinkRef `json:"self,omitempty"`
		Workspace LinkRef `json:"workspace,omitempty"`
	} `json:"_links,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type ServicenowCMDBDeployment struct {
	Links *struct {
		Self        LinkRef `json:"self,omitempty"`
		Workspace   LinkRef `json:"workspace,omitempty"`
		Policy      LinkRef `json:"policy,omitempty"`
		JobMetadata LinkRef `json:"jobMetadata,omitempty"`
	} `json:"_links,omitempty"`
	ID                     int                      `json:"id,omitempty"`
	PolicyID               int                      `json:"policyId,omitempty"`
	Policy                 string                   `json:"policy,omitempty"`
	WorkspaceURL           string                   `json:"workspace,omitempty"`
	ConfigurationItemsInfo []map[string]interface{} `json:"configurationItemsInfo,omitempty"`
	ExecutionDetails       map[string]interface{}   `json:"executionDetails,omitempty"`
	Archived               bool                     `json:"archived,omitempty"`
	TemplateProperties     map[string]interface{}   `json:"templateProperties"`
}

type JobStatus struct {
	Links *struct {
		Self          LinkRef `json:"self,omitempty"`
		JobMetadata   LinkRef `json:"jobMetadata,omitempty"`
		ManagedObject LinkRef `json:"managedObject,omitempty"`
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
	} `json:"errorDetails,omitempty"`
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
		Output          string `json:"output"`
		Status          string `json:"status"`
		JobTemplateName string `json:"jobTemplateName"`
	} `json:"provisioningJobResults,omitempty"`
	DeprovisioningJobResults *struct {
		Output          string `json:"output"`
		Status          string `json:"status"`
		JobTemplateName string `json:"jobTemplateName"`
	} `json:"deprovisioningJobResults,omitempty"`
	TemplateProperties map[string]interface{} `json:"templateProperties"`
}

type ScriptingDeployment struct {
	Links *struct {
		Self        LinkRef `json:"self,omitempty"`
		Workspace   LinkRef `json:"workspace,omitempty"`
		Policy      LinkRef `json:"policy,omitempty"`
		JobMetadata LinkRef `json:"jobMetadata,omitempty"`
	} `json:"_links,omitempty"`
	ID                  int    `json:"id,omitempty"`
	PolicyID            int    `json:"policyId,omitempty"`
	Policy              string `json:"policy,omitempty"`
	WorkspaceURL        string `json:"workspace,omitempty"`
	Hostname            string `json:"hostname,omitempty"`
	ProvisioningDetails *struct {
		Status string   `json:"status"`
		Output []string `json:"output"`
	} `json:"provisioningDetails,omitempty"`
	DeprovisioningDetails *struct {
		Status string   `json:"status"`
		Output []string `json:"output"`
	} `json:"deprovisioningDetails,omitempty"`
	Archived           bool                   `json:"archived,omitempty"`
	TemplateProperties map[string]interface{} `json:"templateProperties"`
}

type ScriptingPolicyResponse struct {
	Embedded struct {
		ScriptingPolicies []ScriptingPolicy `json:"scriptingPolicies"`
	} `json:"_embedded"`
}

type ScriptingPolicy struct {
	Links *struct {
		Self      LinkRef `json:"self,omitempty"`
		Workspace LinkRef `json:"workspace,omitempty"`
	} `json:"_links,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type AnsibleTowerPolicyResponse struct {
	Embedded struct {
		AnsibleTowerPolicies []AnsibleTowerPolicy `json:"ansibleTowerPolicies"`
	} `json:"_embedded"`
}

type AnsibleTowerPolicy struct {
	Links *struct {
		Self      LinkRef `json:"self,omitempty"`
		Workspace LinkRef `json:"workspace,omitempty"`
	} `json:"_links,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// add outputs to this struct once deploy is done
// like the provisioningdetails for scripting above
type VraDeployment struct {
	Links *struct {
		Self        LinkRef `json:"self,omitempty"`
		Workspace   LinkRef `json:"workspace,omitempty"`
		Policy      LinkRef `json:"policy,omitempty"`
		JobMetadata LinkRef `json:"jobMetadata,omitempty"`
	} `json:"_links,omitempty"`
	ID                 int                    `json:"id,omitempty"`
	PolicyID           int                    `json:"policyId,omitempty"`
	Policy             string                 `json:"policy,omitempty"`
	WorkspaceURL       string                 `json:"workspace,omitempty"`
	DeploymentName     string                 `json:"deploymentName,omitempty"`
	Name               string                 `json:"name,omitempty"`
	Archived           bool                   `json:"archived,omitempty"`
	TemplateProperties map[string]interface{} `json:"templateProperties"`
	DeploymentInfo     map[string]interface{} `json:"deploymentInfo,omitempty"`
	BlueprintName      string                 `json:"blueprintName,omitempty"`
	ProjectName        string                 `json:"projectName,omitempty"`
}

type VraPolicyResponse struct {
	Embedded struct {
		VraPolicies []VraPolicy `json:"vraPolicies"`
	} `json:"_embedded"`
}

type VraPolicy struct {
	Links *struct {
		Self      LinkRef `json:"self,omitempty"`
		Workspace LinkRef `json:"workspace,omitempty"`
	} `json:"_links,omitempty"`
	ID          int    `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
type ModuleDeployment struct {
	Links *struct {
		Self        LinkRef `json:"self,omitempty"`
		Workspace   LinkRef `json:"workspace,omitempty"`
		Policy      LinkRef `json:"policy,omitempty"`
		JobMetadata LinkRef `json:"jobMetadata,omitempty"`
	} `json:"_links,omitempty"`
	ID                       int                    `json:"id,omitempty"`
	PolicyID                 int                    `json:"policyId,omitempty"`
	Policy                   string                 `json:"policy,omitempty"`
	WorkspaceURL             string                 `json:"workspace,omitempty"`
	Name                     string                 `json:"name,omitempty"`
	Archived                 bool                   `json:"archived,omitempty"`
	TemplateProperties       map[string]interface{} `json:"templateProperties"`
	ProvisioningJobResults   map[string]interface{} `json:"provisioningJobResults,omitempty"`
	DeprovisioningJobResults map[string]interface{} `json:"deprovisioningJobResults,omitempty"`
}

type ModulePolicyResponse struct {
	Embedded struct {
		ModulePolicies []ModulePolicy `json:"modulePolicies"`
	} `json:"_embedded"`
}

type ModulePolicy struct {
	Links *struct {
		Self      LinkRef `json:"self,omitempty"`
		Workspace LinkRef `json:"workspace,omitempty"`
	} `json:"_links,omitempty"`
	ID             int    `json:"id,omitempty"`
	Name           string `json:"name,omitempty"`
	Description    string `json:"description,omitempty"`
	PolicyTemplate string `json:"policyTemplate,omitempty"`
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

func buildPutRequest(config *Config, resourceType string, requestEntity interface{}, id int) (*http.Request, error) {
	url := itemURL(config, resourceType, id)

	jsonBytes, err := json.Marshal(requestEntity)
	if err != nil {
		return nil, errors.WithMessage(err, "onefuse.apiClient: Failed to marshal request body to JSON")
	}

	requestBody := string(jsonBytes)
	payload := strings.NewReader(requestBody)

	req, err := http.NewRequest("PUT", url, payload)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Unable to create request PUT %s %s", url, requestBody))
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

// Start Ansible Tower Deployment

func (apiClient *OneFuseAPIClient) CreateAnsibleTowerDeployment(newAnsibleTowerDeployment *AnsibleTowerDeployment) (*AnsibleTowerDeployment, error) {
	log.Println("onefuse.apiClient: CreateAnsibleTowerDeployment")

	config := apiClient.config

	var err error
	if newAnsibleTowerDeployment.WorkspaceURL, err = findWorkspaceURLOrDefault(config, newAnsibleTowerDeployment.WorkspaceURL); err != nil {
		return nil, err
	}

	if newAnsibleTowerDeployment.Policy == "" {
		if newAnsibleTowerDeployment.PolicyID != 0 {
			newAnsibleTowerDeployment.Policy = itemURL(config, WorkspaceResourceType, newAnsibleTowerDeployment.PolicyID)
		} else {
			return nil, errors.New("onefuse.apiClient: Ansible Tower Deployment Create requires a PolicyID or Policy URL")
		}
	} else {
		return nil, errors.New("onefuse.apiClient: Ansible Tower Deployment Create requires a PolicyID or Policy URL")
	}

	var req *http.Request
	if req, err = buildPostRequest(config, AnsibleTowerDeploymentResourceType, newAnsibleTowerDeployment); err != nil {
		return nil, err
	}

	ansibleTowerDeployment := AnsibleTowerDeployment{}

	_, err = handleAsyncRequestAndFetchManagdObject(req, config, &ansibleTowerDeployment, "POST")
	if err != nil {
		return nil, err
	}

	return &ansibleTowerDeployment, nil
}

func (apiClient *OneFuseAPIClient) GetAnsibleTowerDeployment(id int) (*AnsibleTowerDeployment, error) {
	log.Println("onefuse.apiClient: GetAnsibleTowerDeployment")

	config := apiClient.config

	url := itemURL(config, AnsibleTowerDeploymentResourceType, id)

	ansibleTowerDeployment := AnsibleTowerDeployment{}
	err := doGet(config, url, &ansibleTowerDeployment)
	if err != nil {
		return nil, err
	}
	return &ansibleTowerDeployment, err
}

func (apiClient *OneFuseAPIClient) UpdateAnsibleTowerDeployment(id int, updatedAnsibleTowerDeployment *AnsibleTowerDeployment) (*AnsibleTowerDeployment, error) {
	log.Println("onefuse.apiClient: UpdateAnsibleTowerDeployment")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) DeleteAnsibleTowerDeployment(id int) error {
	log.Println("onefuse.apiClient: DeleteAnsibleTowerDeployment")

	config := apiClient.config

	url := itemURL(config, AnsibleTowerDeploymentResourceType, id)

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

// End Ansible Tower Deployment

// Start vRA Deployment

func (apiClient *OneFuseAPIClient) CreateVraDeployment(newVraDeployment *VraDeployment) (*VraDeployment, error) {
	log.Println("onefuse.apiClient: CreateVraDeployment")

	config := apiClient.config

	var err error
	if newVraDeployment.WorkspaceURL, err = findWorkspaceURLOrDefault(config, newVraDeployment.WorkspaceURL); err != nil {
		return nil, err
	}

	if newVraDeployment.Policy == "" {
		if newVraDeployment.PolicyID != 0 {
			newVraDeployment.Policy = itemURL(config, VraPolicyResourceType, newVraDeployment.PolicyID)
		} else {
			return nil, errors.New("onefuse.apiClient: vRA Deployment Create requires a PolicyID or Policy URL")
		}
	} else {
		return nil, errors.New("onefuse.apiClient: vRA Deployment Create requires a PolicyID or Policy URL")
	}

	var req *http.Request
	if req, err = buildPostRequest(config, VraDeploymentResourceType, newVraDeployment); err != nil {
		return nil, err
	}

	vraDeployment := VraDeployment{}

	_, err = handleAsyncRequestAndFetchManagdObject(req, config, &vraDeployment, "POST")
	if err != nil {
		return nil, err
	}

	return &vraDeployment, nil
}

func (apiClient *OneFuseAPIClient) GetVraDeployment(id int) (*VraDeployment, error) {
	log.Println("onefuse.apiClient: GetVraDeployment")

	config := apiClient.config

	url := itemURL(config, VraDeploymentResourceType, id)

	vraDeployment := VraDeployment{}
	err := doGet(config, url, &vraDeployment)
	if err != nil {
		return nil, err
	}
	return &vraDeployment, err
}

func (apiClient *OneFuseAPIClient) UpdateVraDeployment(id int, updatedVraDeployment *VraDeployment) (*VraDeployment, error) {
	log.Println("onefuse.apiClient: UpdateVraDeployment")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) DeleteVraDeployment(id int) error {
	log.Println("onefuse.apiClient: DeleteVraDeployment")

	config := apiClient.config

	url := itemURL(config, VraDeploymentResourceType, id)

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

// End vRA Deployment

// Start IPAM Policies

func (apiClient *OneFuseAPIClient) GetIPAMPolicy(id int) (*IPAMPolicy, error) {
	log.Println("onefuse.apiClient: GetIPAMPolicy")
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

// Start Scripting

func (apiClient *OneFuseAPIClient) CreateScriptingDeployment(newScriptingDeployment *ScriptingDeployment) (*ScriptingDeployment, error) {
	log.Println("onefuse.apiClient: CreateScriptingDeployment")

	config := apiClient.config

	var err error
	if newScriptingDeployment.WorkspaceURL, err = findWorkspaceURLOrDefault(config, newScriptingDeployment.WorkspaceURL); err != nil {
		return nil, err
	}

	if newScriptingDeployment.Policy == "" {
		if newScriptingDeployment.PolicyID != 0 {
			newScriptingDeployment.Policy = itemURL(config, WorkspaceResourceType, newScriptingDeployment.PolicyID)
		} else {
			return nil, errors.New("onefuse.apiClient: Scripting Deployment Create requires a PolicyID or Policy URL")
		}
	} else {
		return nil, errors.New("onefuse.apiClient: Scripting Deployment Create requires a PolicyID or Policy URL")
	}

	var req *http.Request
	if req, err = buildPostRequest(config, ScriptingDepoloymentResourceType, newScriptingDeployment); err != nil {
		return nil, err
	}

	scriptingDeployment := ScriptingDeployment{}
	_, err = handleAsyncRequestAndFetchManagdObject(req, config, &scriptingDeployment, "POST")
	if err != nil {
		return nil, err
	}

	return &scriptingDeployment, nil
}

func (apiClient *OneFuseAPIClient) GetScriptingDeployment(id int) (*ScriptingDeployment, error) {
	log.Println("onefuse.apiClient: GetScriptingDeployment")

	config := apiClient.config

	url := itemURL(config, ScriptingDepoloymentResourceType, id)
	scriptingDeployment := ScriptingDeployment{}
	err := doGet(config, url, &scriptingDeployment)
	if err != nil {
		return nil, err
	}
	return &scriptingDeployment, err
}

func (apiClient *OneFuseAPIClient) UpdateScriptingDeployment(id int, desiredScriptingDeployment *ScriptingDeployment) (*ScriptingDeployment, error) {
	log.Println("onefuse.apiClient: UpdateScriptingDeployment")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) DeleteScriptingDeployment(id int) error {
	log.Println("onefuse.apiClient: DeleteScriptingDeployment")

	config := apiClient.config

	url := itemURL(config, ScriptingDepoloymentResourceType, id)

	req, err := http.NewRequest("DELETE", url, nil)
	if err != nil {
		return errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to create request DELETE %s", url))
	}

	setHeaders(req, config)

	_, err = handleAsyncRequest(req, config, "DELETE")
	return err
}

// End Scripting

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

// Start Scripting Policies

func (apiClient *OneFuseAPIClient) GetScriptingPolicy(id int) (*ScriptingPolicy, error) {
	log.Println("onefuse.apiClient: ScriptingPolicy")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetScriptingPolicyByName(name string) (*ScriptingPolicy, error) {
	log.Println("onefuse.apiClient: GetScriptingPolicyByName")

	config := apiClient.config

	scriptingPolicies := ScriptingPolicyResponse{}
	entity, err := findEntityByName(config, name, ScriptingPolicyResourceType, &scriptingPolicies, "ScriptingPolicies", "")
	if err != nil {
		return nil, err
	}
	scriptingPolicy := entity.(ScriptingPolicy)
	return &scriptingPolicy, nil
}

// End Scripting Policies

// Start Ansible Tower Policies

func (apiClient *OneFuseAPIClient) GetAnsibleTowerPolicy(id int) (*AnsibleTowerPolicy, error) {
	log.Println("onefuse.apiClient: AnsibleTowerPolicy")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetAnsibleTowerPolicyByName(name string) (*AnsibleTowerPolicy, error) {
	log.Println("onefuse.apiClient: GetAnsibleTowerPolicyByName")

	config := apiClient.config

	ansibleTowerPolicies := AnsibleTowerPolicyResponse{}
	entity, err := findEntityByName(config, name, AnsibleTowerPolicyResourceType, &ansibleTowerPolicies, "AnsibleTowerPolicies", "")
	if err != nil {
		return nil, err
	}
	ansibleTowerPolicy := entity.(AnsibleTowerPolicy)
	return &ansibleTowerPolicy, nil
}

// End Ansible Tower Policies

// Start ServicenowCMDB Policies

func (apiClient *OneFuseAPIClient) GetServicenowCMDBPolicy(id int) (*ServicenowCMDBPolicy, error) {
	log.Println("onefuse.apiClient: ServicenowCMDBPolicy")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetServicenowCMDBPolicyByName(name string) (*ServicenowCMDBPolicy, error) {
	log.Println("onefuse.apiClient: GetServicenowCMDBPolicyByName")

	config := apiClient.config

	servicenowCMDBPolicies := ServicenowCMDBPolicyResponse{}
	entity, err := findEntityByName(config, name, ServicenowCMDBPolicyResourceType, &servicenowCMDBPolicies, "ServicenowCMDBPolicies", "")
	if err != nil {
		return nil, err
	}
	servicenowCMDBPolicy := entity.(ServicenowCMDBPolicy)
	return &servicenowCMDBPolicy, nil
}

// End ServicenowCMDB Policies

// Start vRA Policies

func (apiClient *OneFuseAPIClient) GetVraPolicy(id int) (*VraPolicy, error) {
	log.Println("onefuse.apiClient: VraPolicy")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetVraPolicyByName(name string) (*VraPolicy, error) {
	log.Println("onefuse.apiClient: GetVraPolicyByName")

	config := apiClient.config

	vraPolicies := VraPolicyResponse{}
	entity, err := findEntityByName(config, name, VraPolicyResourceType, &vraPolicies, "VraPolicies", "")
	if err != nil {
		return nil, err
	}
	vraPolicy := entity.(VraPolicy)
	return &vraPolicy, nil
}

// End vRA Policies

// Start Static Property Set

func (apiClient *OneFuseAPIClient) GetStaticPropertySet(id int) (*StaticPropertySet, error) {
	log.Println("onefuse.apiClient: GetStaticPropertySet")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetStaticPropertySetByName(name string) (*StaticPropertySet, error) {
	log.Println("onefuse.apiClient: GetStaticPropertySetByName")

	config := apiClient.config

	staticPropertySets := StaticPropertySetResponse{}
	entity, err := findEntityByName(config, name, StaticPropertySetResourceType, &staticPropertySets, "PropertySets")
	if err != nil {
		return nil, err
	}

	staticPropertySet := entity.(StaticPropertySet)

	raw, err := json.Marshal(staticPropertySet.Properties)
	if err != nil {
		return nil, err
	}
	staticPropertySet.Raw = string(raw)

	return &staticPropertySet, nil
}

// End Static Property Set

// Start Jobs

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

// End Jobs

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
	PollingTimeoutMS := 3600000
	PollingIntervalMS := 5000
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

// Start Render Template

func (apiClient *OneFuseAPIClient) RenderTemplate(template string, templateProperties map[string]interface{}) (*RenderTemplateResponse, error) {
	// this API endpoint is a POST, but only so we can pass in a body to be rendered by the templating engine
	// it behaves mostly like a GET, and doesn't create an object, just returns the rendered value.
	log.Println("onefuse.apiClient: RenderTemplate")

	config := apiClient.config

	requestBody := RenderTemplateRequest{
		Template:           template,
		TemplateProperties: templateProperties,
	}

	var err error

	var req *http.Request
	if req, err = buildPostRequest(config, RenderTemplateType, requestBody); err != nil {
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

	renderTemplateResponse := RenderTemplateResponse{}
	if err = json.Unmarshal(body, &renderTemplateResponse); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}
	renderedTemplate := renderTemplateResponse

	return &renderedTemplate, nil
}

// End Render Template

// Start ServiceNow CMDB Deployment

func (apiClient *OneFuseAPIClient) CreateServicenowCMDBDeployment(newServicenowCMDBDeployment *ServicenowCMDBDeployment) (*ServicenowCMDBDeployment, error) {
	log.Println("onefuse.apiClient: CreateServicenowCMDBDeployment")

	config := apiClient.config

	var err error
	if newServicenowCMDBDeployment.WorkspaceURL, err = findWorkspaceURLOrDefault(config, newServicenowCMDBDeployment.WorkspaceURL); err != nil {
		return nil, err
	}

	if newServicenowCMDBDeployment.Policy == "" {
		if newServicenowCMDBDeployment.PolicyID != 0 {
			newServicenowCMDBDeployment.Policy = itemURL(config, WorkspaceResourceType, newServicenowCMDBDeployment.PolicyID)
		} else {
			return nil, errors.New("onefuse.apiClient: ServiceNow CMDB Deployment Create requires a PolicyID or Policy URL")
		}
	} else {
		return nil, errors.New("onefuse.apiClient: ServiceNow CMDB Deployment Create requires a PolicyID or Policy URL")
	}

	var req *http.Request
	if req, err = buildPostRequest(config, ServicenowCMDBDepoloymentResourceType, newServicenowCMDBDeployment); err != nil {
		return nil, err
	}

	servicenowCMDBDeployment := ServicenowCMDBDeployment{}

	_, err = handleAsyncRequestAndFetchManagdObject(req, config, &servicenowCMDBDeployment, "POST")
	if err != nil {
		return nil, err
	}

	return &servicenowCMDBDeployment, nil
}

func (apiClient *OneFuseAPIClient) GetServicenowCMDBDeployment(id int) (*ServicenowCMDBDeployment, error) {
	log.Println("onefuse.apiClient: GetServicenowCMDBDeployment")

	config := apiClient.config

	url := itemURL(config, ServicenowCMDBDepoloymentResourceType, id)

	servicenowCMDBDeployment := ServicenowCMDBDeployment{}
	err := doGet(config, url, &servicenowCMDBDeployment)
	if err != nil {
		return nil, err
	}
	return &servicenowCMDBDeployment, err
}

func (apiClient *OneFuseAPIClient) UpdateServicenowCMDBDeployment(id int, updatedServicenowCMDBDeployment *ServicenowCMDBDeployment) (*ServicenowCMDBDeployment, error) {
	log.Println("onefuse.apiClient: UpdateServicenowCMDBDeployment")

	config := apiClient.config

	var err error
	if updatedServicenowCMDBDeployment.WorkspaceURL, err = findWorkspaceURLOrDefault(config, updatedServicenowCMDBDeployment.WorkspaceURL); err != nil {
		return nil, err
	}

	if updatedServicenowCMDBDeployment.Policy == "" {
		if updatedServicenowCMDBDeployment.PolicyID != 0 {
			updatedServicenowCMDBDeployment.Policy = itemURL(config, WorkspaceResourceType, updatedServicenowCMDBDeployment.PolicyID)
		} else {
			return nil, errors.New("onefuse.apiClient: ServiceNow CMDB Deployment Update requires a PolicyID or Policy URL")
		}
	} else {
		return nil, errors.New("onefuse.apiClient: ServiceNow CMDB Deployment Update requires a PolicyID or Policy URL")
	}

	var req *http.Request
	if req, err = buildPutRequest(config, ServicenowCMDBDepoloymentResourceType, updatedServicenowCMDBDeployment, id); err != nil {
		return nil, err
	}

	servicenowCMDBDeployment := ServicenowCMDBDeployment{}

	_, err = handleAsyncRequestAndFetchManagdObject(req, config, &servicenowCMDBDeployment, "PUT")
	if err != nil {
		return nil, err
	}

	return &servicenowCMDBDeployment, nil
}

func (apiClient *OneFuseAPIClient) DeleteServicenowCMDBDeployment(id int) error {
	log.Println("onefuse.apiClient: DeleteServicenowCMDBDeployment")

	config := apiClient.config

	url := itemURL(config, ServicenowCMDBDepoloymentResourceType, id)

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

// End ServiceNow CMDB Deployment

// Start Module Policies

func (apiClient *OneFuseAPIClient) GetModulePolicy(id int) (*ModulePolicy, error) {
	log.Println("onefuse.apiClient: GetModulePolicy")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetModulePolicyByName(name string) (*ModulePolicy, error) {
	log.Println("onefuse.apiClient: GetModulePolicyByName")

	config := apiClient.config

	modulePolicies := ModulePolicyResponse{}
	entity, err := findEntityByName(config, name, ModulePolicyResourceType, &modulePolicies, "ModulePolicies", "")
	if err != nil {
		return nil, err
	}
	modulePolicy := entity.(ModulePolicy)
	return &modulePolicy, nil
}

// End Module Policies

// Start Module Deployment

func (apiClient *OneFuseAPIClient) CreateModuleDeployment(newModuleDeployment *ModuleDeployment) (*ModuleDeployment, error) {
	log.Println("onefuse.apiClient: CreateModuleDeployment")

	config := apiClient.config

	var err error
	if newModuleDeployment.WorkspaceURL, err = findWorkspaceURLOrDefault(config, newModuleDeployment.WorkspaceURL); err != nil {
		return nil, err
	}

	if newModuleDeployment.Policy == "" {
		if newModuleDeployment.PolicyID != 0 {
			newModuleDeployment.Policy = itemURL(config, WorkspaceResourceType, newModuleDeployment.PolicyID)
		} else {
			return nil, errors.New("onefuse.apiClient: Module Deployment Create requires a PolicyID or Policy URL")
		}
	} else {
		return nil, errors.New("onefuse.apiClient: Module Deployment Create requires a PolicyID or Policy URL")
	}

	var req *http.Request
	if req, err = buildPostRequest(config, ModuleDepoloymentResourceType, newModuleDeployment); err != nil {
		return nil, err
	}

	ModuleDeployment := ModuleDeployment{}

	_, err = handleAsyncRequestAndFetchManagdObject(req, config, &ModuleDeployment, "POST")
	if err != nil {
		return nil, err
	}

	return &ModuleDeployment, nil
}

func (apiClient *OneFuseAPIClient) GetModuleDeployment(id int) (*ModuleDeployment, error) {
	log.Println("onefuse.apiClient: GetModuleDeployment")

	config := apiClient.config

	url := itemURL(config, ModuleDepoloymentResourceType, id)

	ModuleDeployment := ModuleDeployment{}
	err := doGet(config, url, &ModuleDeployment)
	if err != nil {
		return nil, err
	}
	return &ModuleDeployment, err
}

func (apiClient *OneFuseAPIClient) UpdateModuleDeployment(id int, updatedModuleDeployment *ModuleDeployment) (*ModuleDeployment, error) {
	log.Println("onefuse.apiClient: UpdateModuleDeployment")

	config := apiClient.config

	var err error
	if updatedModuleDeployment.WorkspaceURL, err = findWorkspaceURLOrDefault(config, updatedModuleDeployment.WorkspaceURL); err != nil {
		return nil, err
	}

	if updatedModuleDeployment.Policy == "" {
		if updatedModuleDeployment.PolicyID != 0 {
			updatedModuleDeployment.Policy = itemURL(config, WorkspaceResourceType, updatedModuleDeployment.PolicyID)
		} else {
			return nil, errors.New("onefuse.apiClient: Module Deployment Update requires a PolicyID or Policy URL")
		}
	} else {
		return nil, errors.New("onefuse.apiClient: Module Deployment Update requires a PolicyID or Policy URL")
	}

	var req *http.Request
	if req, err = buildPutRequest(config, ModuleDepoloymentResourceType, updatedModuleDeployment, id); err != nil {
		return nil, err
	}

	ModuleDeployment := ModuleDeployment{}

	_, err = handleAsyncRequestAndFetchManagdObject(req, config, &ModuleDeployment, "PUT")
	if err != nil {
		return nil, err
	}

	return &ModuleDeployment, nil
}

func (apiClient *OneFuseAPIClient) DeleteModuleDeployment(id int) error {
	log.Println("onefuse.apiClient: DeleteModuleDeployment")

	config := apiClient.config

	url := itemURL(config, ModuleDepoloymentResourceType, id)

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

// End Module Deployment

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
