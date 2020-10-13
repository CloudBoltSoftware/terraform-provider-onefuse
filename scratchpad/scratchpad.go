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

	if len(staticPropertySets.Embedded.staticPropertySets) < 1 {
		return nil, errors.New(fmt.Sprintf("onefuse.apiClient: Could not find staticPropertySet '%s'!", name))
	}

	staticPropertySet := staticPropertySets.Embedded.staticPropertySets[0]

	return &staticPropertySet, err
}

// End Static Property Set