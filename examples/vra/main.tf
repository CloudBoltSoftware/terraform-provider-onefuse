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
  scheme     = "https"
  address    = "onefuse_fqdn"
  port       = "port"
  user       = "admin"
  password   = "admin"
  verify_ssl = "false"
}

// vRealize Automation Deployment Resource
resource "onefuse_vra_deployment" "vra_deployment_01" {
  policy_id       = 1                               // Refers to vRealize Automation Policy ID (integer)
  workspace_url   = ""                              // Leave blank for default workspace
  deployment_name = "tf_vra_deployment"             // vRA Deployment Name
  template_properties = {                           // Your properties and its values to pass into module
    property1 = "value1"
    property2 = "value2"
    property3 = "value3"
  }
  timeouts {
    create = "12m"
    delete = "3m"
  }
}

// Output
output "deployment_info" {
  value = onefuse_vra_deployment.vra_deployment_01.deployment_info
}
