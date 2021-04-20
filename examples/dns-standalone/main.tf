// Comment out for Terraform 0.12
terraform {
  required_providers {
    onefuse = {
      source  = "CloudBoltSoftware/onefuse"
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

// Policy Data Source: DNS
data "onefuse_dns_policy" "default" {
  name = "default"
}

// DNS Record Object Resource
resource "onefuse_dns_record" "my-dns-record" {
  name          = "hostname"
  policy_id     = data.onefuse_dns_policy.default.id // Refers to onefuse_dns_policy data source to retrieve ID
  workspace_url = ""                                 // Leave blank for default workspace
  zones         = ["dnszone1,dnszone2"]              // Comma separated DNS Zones.  At least one zone required
  value         = "ipAddress"                        // IP Address
  template_properties = {                            // Your properties and its values to pass into module
    property1 = "value1"
    property2 = "value2"
    property3 = "value3"
  }
}
