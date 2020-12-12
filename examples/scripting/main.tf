terraform {
  required_providers {
    onefuse = {
      source = "CloudBoltSoftware/onefuse"
      version = ">= 1.1.0"
    }
  }
  required_version = ">= 0.13"
}

provider "onefuse" {

  scheme     = "https"
  address    = "onefuse_fqdn"
  port       = "port"
  user       = "admin"
  password   = "admin"
  verify_ssl = "false"
}

resource "onefuse_scripting_deployment" "my-scripting-deployment" {
    policy_id = 1 //data.onefuse_dns_policy.my_dns.id // Refers to onefuse_dns_policy data source to retrieve ID
    workspace_url = "" // Leave blank for default workspace
    template_properties = {
        property1        = "value1" // Your properties and values to pass into module
        proeprty2        = "value2"
        property3        = "value3"
  }
}
