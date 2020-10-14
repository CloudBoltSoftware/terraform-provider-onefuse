// Start DNS Policies

func (apiClient *OneFuseAPIClient) GetDNSPolicy(id int) (*DNSPolicy, error) {
	log.Println("onefuse.apiClient: DNSPolicy")
	return nil, errors.New("onefuse.apiClient: Not implemented yet")
}

func (apiClient *OneFuseAPIClient) GetDNSPolicyByName(name string) (*DNSPolicy, error) {
	log.Println("onefuse.apiClient: GetDNSPolicyByName")

	config := apiClient.config
	url := fmt.Sprintf("%s?filter=name:%s", collectionURL(config, DNSPolicyResourceType), name)

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

	dnsPolicies := DNSPolicyResponse{}
	err = json.Unmarshal(body, &dnsPolicies)
	if err != nil {
		return nil, errors.WithMessage(err, fmt.Sprintf("onefuse.apiClient: Failed to unmarshal response %s", string(body)))
	}

	if len(dnsPolicies.Embedded.DNSPolicies) < 1 {
		return nil, errors.New(fmt.Sprintf("onefuse.apiClient: Could not find AD Policy '%s'!", name))
	}

	dnsPolicy := dnsPolicies.Embedded.DNSPolicies[0]

	return &dnsPolicy, err
}

// End DNS Policies