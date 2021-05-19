// Commented out for Terraform 0.12

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


// Inititalize OneFuse Provider
provider "onefuse" {
  scheme     = "https"
  address    = "onefuse_fqdn"
  port       = "port"
  user       = "admin"
  password   = "admin"
  verify_ssl = "false"
}

// DNS Policy data source
//data "onefuse_dns_policy" "my_dns" {
//  name = "infoblox851_dnspolicy"
//}

// DNS computer object resource
resource "onefuse_dns_record" "my-dns-record" {
    name = "hostname"
    policy_id = 1 //data.onefuse_dns_policy.my_dns.id // Refers to onefuse_dns_policy data source to retrieve ID
    workspace_url = "" // Leave blank for default workspace
    zones = ["dnszone1,dnszone2"] // Comma seperated dns zones.  At least one zone required
    value = "ipAddress" // IP Address
    template_properties = {
        property1        = "value1" // Your properties and values to pass into module
        proeprty2        = "value2"
        property3        = "value3"
  }
}
