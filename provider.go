package main

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"address": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEFUSE_ADDRESS", nil),
				Description: "OneFuse REST endpoint service host address",
			},
			"port": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEFUSE_PORT", nil),
				Description: "OneFuse REST endpoint service port number",
			},
			"user": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEFUSE_USER", nil),
				Description: "OneFuse REST endpoint user name",
			},
			"password": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				DefaultFunc: schema.EnvDefaultFunc("ONEFUSE_PASSWORD", nil),
				Description: "OneFuse REST endpoint password",
			},
			"verify_ssl": {
				Type:        schema.TypeBool,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEFUSE_VERIFY_SSL", true),
				Description: "Verify SSL certificates for OneFuse endpoints",
			},
		},
		ResourcesMap: map[string]*schema.Resource{
			"onefuse_naming": resourceCustomNaming(),
		},
		ConfigureFunc: configureProvider,
	}
}

type Config struct {
	scheme    string
	address   string
	port      string
	user      string
	password  string
	verifySSL bool
}

func configureProvider(d *schema.ResourceData) (interface{}, error) {
	return Config{
		scheme:    "https",
		address:   d.Get("address").(string),
		port:      d.Get("port").(string),
		user:      d.Get("user").(string),
		password:  d.Get("password").(string),
		verifySSL: d.Get("verify_ssl").(bool),
	}, nil
}
