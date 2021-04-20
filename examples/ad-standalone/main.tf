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

// Policy Data Source: Active Directory
data "onefuse_ad_policy" "default" {
  name = "default"
}

// AD Computer Object Resource
resource "onefuse_microsoft_ad_computer_account" "my_ad_computer" {
  name          = "myhostname"
  policy_id     = data.onefuse_ad_policy.default.id // Refers to onefuse_ad_policy data source to retrieve ID
  workspace_url = ""                                // Leave blank for default workspace
  template_properties = {                           // Your properties and its values to pass into module
    property1 = "value1"
    property2 = "value2"
    property3 = "value3"
  }
}

// Output Result for AD OU Placement
output "ad_ou" {
  value = onefuse_microsoft_ad_computer_account.my_ad_computer.final_ou
}
