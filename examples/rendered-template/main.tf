// Commented out for Terraform 0.12

terraform {
  required_providers {
    onefuse = {
      source  = "CloudBoltSoftware/onefuse"
      version = ">= 1.1.0"
    }
  }
  required_version = ">= 0.13"
}

// Comment out above for Terraform 0.12


// Inititalize OneFuse Provider
provider "onefuse" {
  scheme     = "https"
  address    = "onefuse_fqdn"
  port       = "port"
  user       = "admin"
  password   = "admin"
  verify_ssl = "false"
}

// Render Template
// Nested template properties (e.g., {{key.subKey.subSubkey}} ) are not yet supported
// See https://github.com/CloudBoltSoftware/terraform-provider-onefuse/issues/15
data "onefuse_rendered_template" "template01" {
    template = "template {{env}}"
    template_properties = {"env": "prod"}
}

data "onefuse_rendered_template" "template02" {
    template = "template {{env}}"
    template_properties = {"env": "dev"}
}

output "rendered_template01" {
  value = data.onefuse_rendered_template.template01.value
}

output "rendered_template02" {
  value = data.onefuse_rendered_template.template02.value
}
