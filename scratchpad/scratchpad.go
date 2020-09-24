//DNS Functions

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
	url := collectionURL(config, DNSResourceType)

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

	dnsRecord := DNSReservation{}
	if err = json.Unmarshal(body, &dnsRecord); err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	return &dnsRecord, nil
}