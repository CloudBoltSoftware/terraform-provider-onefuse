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

// AD Policy data source
data "onefuse_ad_policy" "default" {
  name = "default"
}

// Ad computer object resource
resource "onefuse_microsoft_ad_computer_account" "my_ad_computer" {

  name          = "namehere"
  policy_id     = data.onefuse_ad_policy.default.id // Refers to onefuse_ad_policy data source to retrieve ID
  workspace_url = ""                                // Leave blank for default workspace
  template_properties = {
    property1 = "value1" // Your properties and values to pass into module
    proeprty2 = "value2"
    property3 = "value3"
  }
}

// Output Result for AD OU Placement
output "ad_ou" {
  value = onefuse_microsoft_ad_computer_account.my_ad_computer.final_ou
}
