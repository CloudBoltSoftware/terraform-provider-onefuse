// Commented out for Terraform 0.12

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


// Inititalize OneFuse Provider
provider "onefuse" {
  scheme     = "https"
  address    = "docker02.sovlabs.net"
  port       = "8722"
  user       = "admin"
  password   = "admin"
  verify_ssl = "false"
}

// ServiceNow CMDB Policy data source
data "onefuse_servicenow_cmdb_policy" "servicenow_cmdb_policy" {
  name = "servicenow_cmdb_policy"
}

resource "onefuse_servicenow_cmdb_deployment" "my-servicenow-cmdb-deployment" {
  policy_id        = data.onefuse_servicenow_cmdb_policy.servicenow_cmdb_policy.id // Refers to onefuse_servicenow_cmdb_policy data source to retrieve ID
  workspace_url    = ""                                    // Leave blank for default workspace
  template_properties = {
    property1 = "value1" // Your properties and values to pass into module
    proeprty2 = "value2"
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