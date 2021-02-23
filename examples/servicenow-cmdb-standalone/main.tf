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

// Outputs
output "policy-id" {
  value = data.onefuse_servicenow_cmdb_policy.servicenow_cmdb_policy.id
}