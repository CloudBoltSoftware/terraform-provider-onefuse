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

// IPAM Policy data source
data "onefuse_ipam_policy" "ipam_policy" {
  name = "infoblox851_ipampolicy"
}

// IPAM reservation resource
resource "onefuse_ipam_record" "my-ipam-record" {
    hostname = "testname"
    policy_id = data.onefuse_ipam_policy.ipam_policy.id // Refers to onefuse_ipam_policy data source to retrieve ID
    workspace_url = "" // Leave blank for default workspace
    template_properties = {
        property1        = "value1" // Your properties and values to pass into module
        property2        = "value2"
        property3        = "value3"
  }
}

// Outputs
output "ip_address" {
  value = onefuse_ipam_record.my-ipam-record.ip_address
}

output "netmask" {
  value = onefuse_ipam_record.my-ipam-record.netmask
}

output "gateway" {
  value = onefuse_ipam_record.my-ipam-record.gateway
}

output "network" {
  value = onefuse_ipam_record.my-ipam-record.network
}

output "primary_dns" {
  value = onefuse_ipam_record.my-ipam-record.primary_dns
}

output "secondary_dns" {
  value = onefuse_ipam_record.my-ipam-record.secondary_dns
}

output "nic_label" {
  value = onefuse_ipam_record.my-ipam-record.nic_label
}