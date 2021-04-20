// Comment out for Terraform 0.12
terraform {
  required_providers {
    onefuse = {
      source  = "CloudBoltSoftware/onefuse"
      version = ">= 1.10.1"
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

// Policy Data Source: ServiceNow CMDB
data "onefuse_servicenow_cmdb_policy" "servicenow_cmdb_policy" {
  name = "servicenow_cmdb_policy"
}

resource "onefuse_servicenow_cmdb_deployment" "my-servicenow-cmdb-deployment" {
  policy_id           = data.onefuse_servicenow_cmdb_policy.servicenow_cmdb_policy.id // Refers to onefuse_servicenow_cmdb_policy data source to retrieve ID
  workspace_url       = ""                                                            // Leave blank for default workspace
  template_properties = {                                                             // Your properties and its values to pass into module
    property1 = "value1"
    property2 = "value2"
    property3 = "value3"
  }
}

// Outputs
output "policy-id" {
  value = data.onefuse_servicenow_cmdb_policy.servicenow_cmdb_policy.id
}

// Outputs
output "servicenow-cmdb-response" {
  value = onefuse_servicenow_cmdb_deployment.my-servicenow-cmdb-deployment
}