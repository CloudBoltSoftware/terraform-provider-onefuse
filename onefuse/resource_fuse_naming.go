// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"log"
	"time"
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func resourceCustomNaming() *schema.Resource {
	return &schema.Resource{
		Create: resourceCustomNameCreate,
		Read:   resourceCustomNameRead,
		Update: resourceCustomNameUpdate,
		Delete: resourceCustomNameDelete,
		Importer: &schema.ResourceImporter{
			State: importNaming,
		},
		Schema: map[string]*schema.Schema{
			"custom_name_id": {
				Type:     schema.TypeInt,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"name": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
				ForceNew: true,
			},
			"naming_policy_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"dns_suffix": {
				Type:     schema.TypeString,
				Optional: true,
				Computed: true,
			},
			"workspace_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"template_properties": {
				Type:        schema.TypeMap,
				Optional:    true,
				Description: "Fuse Template Properties",
			},
		},
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(20 * time.Minute),
			Delete: schema.DefaultTimeout(20 * time.Minute),
		},
	}
}

func bindCustomNamingResource(d *schema.ResourceData, cn *CustomName) error {
	log.Println("onefuse.bindCustomNamingResource")

	// setting the ID is REALLY necessary here
	// we use the FQDN instead of the numeric ID as it is more likely to remain consistent as a composite key in TF
	d.SetId(cn.Name + "." + cn.DnsSuffix)

	if err := d.Set("custom_name_id", cn.Id); err != nil {
		return errors.WithMessage(err, "cannot set custom_name_id")
	}
	if err := d.Set("name", cn.Name); err != nil {
		return errors.WithMessage(err, "cannot set name")
	}
	if err := d.Set("dns_suffix", cn.DnsSuffix); err != nil {
		return errors.WithMessage(err, "cannot set dns_suffix")
	}
	return nil
}

func resourceCustomNameCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceCustomNameCreate")

	config := m.(Config)

	namingPolicyID := d.Get("naming_policy_id").(string)
	workspaceID := d.Get("workspace_id").(string)
	templateProperties := d.Get("template_properties").(map[string]interface{})

	cn, err := config.NewOneFuseApiClient().GenerateCustomName(namingPolicyID, workspaceID, templateProperties)
	if err != nil {
		return err
	}

	return bindCustomNamingResource(d, cn)
}

func resourceCustomNameRead(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceCustomNameRead")

	config := m.(Config)

	id := d.Get("custom_name_id").(int)

	customName, err := config.NewOneFuseApiClient().GetCustomName(id)
	if err != nil {
		return err
	}

	return bindCustomNamingResource(d, customName)
}

func resourceCustomNameUpdate(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceCustomNameUpdate")
	return nil
}

func resourceCustomNameDelete(d *schema.ResourceData, m interface{}) error {
	log.Println("onefuse.resourceCustomNameDelete")

	config := m.(Config)

	id := d.Get("custom_name_id").(int)

	return config.NewOneFuseApiClient().DeleteCustomName(id)
}

func importNaming(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	log.Println("onefuse.importNaming - Starting the import")

	d.SetId(d.Id())
	customNameID, err := strconv.Atoi(d.Id())
	if err != nil {
		return nil, fmt.Errorf("invalid ID format: %s", d.Id())
	}

	config := meta.(Config)
	apiClient := config.NewOneFuseApiClient()

	customName, err := apiClient.GetCustomName(customNameID)
	if err != nil {
		return nil, err
	}

	// Bind the custom name
	if err := bindCustomNamingResource(d, customName); err != nil {
		return nil, err
	}

	jobMetaDataRecord, policyId, err := fetchNameJobMetaData(customName, &config)
	if err != nil {
		log.Printf("Error fetching job metadata: %v", err)
		return nil, errors.Wrap(err, "error fetching job metadata during import")
	}

	if jobMetaDataRecord == nil {
		log.Println("jobMetaDataRecord is nil after fetching job metadata")
		return nil, errors.New("jobMetaDataRecord is nil after fetching job metadata")
	}

	if policyId == "" {
		log.Println("Naming policy id is nil after fetching job metadata")
		return nil, errors.New("Naming policy id is nil after fetching job metadata")
	}

	log.Printf("Template Properties are: %+v", jobMetaDataRecord.ResolvedProperties)
	if err := d.Set("naming_policy_id", policyId); err != nil {
		log.Printf("Error setting policy id: %v", err)
		return nil, errors.Wrap(err, "Cannot set policyId")
	}

	return []*schema.ResourceData{d}, nil
}

func fetchNameJobMetaData(customName *CustomName, config *Config) (jobMetaDataRecord *JobMetaData, policyIdStr string, err error) {
    log.Println("Fetching the job metadata - Start")

    jobMetaDataURLSplit := strings.Split(customName.Links.JobMetadata.Href, "/")
    policyURLSplit := strings.Split(customName.Links.Policy.Href, "/")
    jobMetaDataId := jobMetaDataURLSplit[len(jobMetaDataURLSplit)-2]
    policyIdStr = policyURLSplit[len(policyURLSplit)-2]

    jobMetaDataIdInt, err := strconv.Atoi(jobMetaDataId)
    if err != nil {
        return nil, "", errors.Wrap(err, "failed to convert job metadata ID to int")
    }

    jobMetaDataRecord, err = GetJobMetaData(jobMetaDataIdInt, config)
    if err != nil {
        return nil, "", errors.Wrap(err, "failed to fetch job metadata")
    }

	log.Println("Fetching the job metadata - Completed")

    return jobMetaDataRecord, policyIdStr, nil
}
