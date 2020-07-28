package main

import (
	"log"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/pkg/errors"
)

func resourceCustomNaming() *schema.Resource {
	return &schema.Resource{
		Create: resourceCustomNameCreate,
		Read:   resourceCustomNameRead,
		Update: resourceCustomNameUpdate,
		Delete: resourceCustomNameDelete,
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
				Required: true,
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
	}
}

func resourceCustomNameCreate(d *schema.ResourceData, m interface{}) error {
	log.Println("calling resourceCustomNameCreate")

	config := m.(Config)
	dnsSuffix := d.Get("dns_suffix").(string)
	namingPolicyID := d.Get("naming_policy_id").(string)
	workspaceID := d.Get("workspace_id").(string)
	templateProperties := d.Get("template_properties").(map[string]interface{})
	cn, err := config.NewOneFuseApiClient().GenerateCustomName(dnsSuffix, namingPolicyID, workspaceID, templateProperties)
	if err != nil {
		return err
	}
	err = bindResource(d, *cn)
	return err
}

func bindResource(d *schema.ResourceData, cn CustomName) error {

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

func resourceCustomNameRead(d *schema.ResourceData, m interface{}) error {
	config := m.(Config)
	id := d.Get("custom_name_id").(int)
	customName, err := config.NewOneFuseApiClient().GetCustomName(id)
	bindResource(d, customName)
	return err
}

func resourceCustomNameUpdate(d *schema.ResourceData, m interface{}) error {
	return nil
}

func resourceCustomNameDelete(d *schema.ResourceData, m interface{}) error {
	config := m.(Config)
	id := d.Get("custom_name_id").(int)
	config.NewOneFuseApiClient().DeleteCustomName(id)
	return nil
}
