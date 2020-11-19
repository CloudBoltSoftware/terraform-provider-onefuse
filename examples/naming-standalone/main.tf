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
  address    = "onefuse_fqdn"
  port       = "port"
  user       = "admin"
  password   = "admin"
  verify_ssl = "false"
}

// Naming Policy data source
data "onefuse_naming_policy" "machine" {
  name = "machineNaming"
}

resource "onefuse_naming" "my-onefuse-name" {
  naming_policy_id = data.onefuse_naming_policy.machine.id // Refers to onefuse_naming_policy data source to retrieve ID
  workspace_url    = ""                                    // Leave blank for default workspace
  dns_suffix       = ""
  template_properties = {
    property1 = "value1" // Your properties and values to pass into module
    proeprty2 = "value2"
    property3 = "value3"
  }
}

// Outputs
output "name" {
  value = onefuse_naming.my-onefuse-name.name
}

output "dns_suffix" {
  value = onefuse_naming.my-onefuse-name.dns_suffix // Refers to dns_sudffix output by naming is defined in policy
}
