// Comment out for Terraform 0.12
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


// Initialize OneFuse Provider
provider "onefuse" {
  scheme     = "https"
  address    = "onefuse_fqdn"
  port       = "port"
  user       = "admin"
  password   = "admin"
  verify_ssl = "false"
}

// Scripting Deployment Object Resource
resource "onefuse_scripting_deployment" "my-scripting-deployment" {
  policy_id     = 1                       // Refers to Scripting Policy ID (integer)
  workspace_url = ""                      // Leave blank for default workspace
  template_properties = {                 // Your properties and its values to pass into module
    property1 = "value1"
    property2 = "value2"
    property3 = "value3"
  }
}
