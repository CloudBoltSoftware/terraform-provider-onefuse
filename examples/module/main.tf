// Comment out for Terraform 0.12
terraform {
  required_providers {
    onefuse = {
      source  = "CloudBoltSoftware/onefuse"
      version = ">= 1.10.3"
    }
  }
  required_version = ">= 0.13"
}
// Comment out above for Terraform 0.12


// Initialize OneFuse Provider
provider "onefuse" {
  scheme     = "http"
  address    = "onefuse_fqdn"
  port       = "port"
  user       = "admin"
  password   = "admin"
  verify_ssl = "false"
}

data "onefuse_module_policy" "my_policy" {
  name = "my_module_policy"
}

resource "onefuse_module_deployment" "bar" {
  policy_id = data.onefuse_module_policy.my_policy.id
  workspace_url = "" // Leave blank for default workspace
  template_properties = {
        property1        = "value1" // Your properties and values to pass into module
        property2        = "value2"
  }
  timeouts {
    create = "60m"
    delete = "30m"
  }
}
