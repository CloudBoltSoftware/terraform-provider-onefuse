// Copyright 2020 CloudBolt Software
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at https://mozilla.org/MPL/2.0/.

package onefuse

import (
	"github.com/hashicorp/terraform/helper/schema"
)

func Provider() *schema.Provider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"scheme": {
				Type:        schema.TypeString,
				Optional:    true,
				DefaultFunc: schema.EnvDefaultFunc("ONEFUSE_SCHEME", "https"),
				Description: "OneFuse REST endpoint service http(s) scheme",
			},
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
			"onefuse_naming":                        resourceCustomNaming(),
			"onefuse_microsoft_ad_policy":           resourceMicrosoftADPolicy(),
			"onefuse_microsoft_ad_computer_account": resourceMicrosoftADComputerAccount(),
			"onefuse_dns_record":                    resourceDNSReservation(),
			"onefuse_ipam_record":                   resourceIPAMReservation(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"onefuse_microsoft_endpoint": dataSourceMicrosoftEndpoint(),
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
	return NewConfig(
		d.Get("scheme").(string),
		d.Get("address").(string),
		d.Get("port").(string),
		d.Get("user").(string),
		d.Get("password").(string),
		d.Get("verify_ssl").(bool),
	), nil
}

func NewConfig(scheme string, address string, port string, user string, password string, verifySSL bool) Config {
	return Config{
		scheme:    scheme,
		address:   address,
		port:      port,
		user:      user,
		password:  password,
		verifySSL: verifySSL,
	}
}
