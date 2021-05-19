// Commented out for Terraform 0.12

terraform {
  required_providers {
    onefuse = {
      source = "CloudBoltSoftware/onefuse"
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
  port       = "443"
  user       = "admin"
  password   = "admin"
  verify_ssl = "false"
}

data "onefuse_scripting_policy" "my_policy" {
  // name = "My Scripting Policy"
  name = "myScriptingPolicy"
}

// Onefuse Scripting Deployment
resource "onefuse_scripting_deployment" "my-scripting-deployment" {
    policy_id = data.onefuse_scripting_policy.my_policy.id
    workspace_url = ""
    template_properties = {
        property1        = "value1" // Your properties and values to pass into module
        proeprty2        = "value2"
        property3        = "value3"
  }
}
