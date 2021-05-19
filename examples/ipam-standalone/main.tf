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

// Policy Data Source: IPAM
data "onefuse_ipam_policy" "default" {
  name = "default"
}

// IP Address Object Resource
resource "onefuse_ipam_record" "my-ipam-record" {
  hostname      = "hostname"
  policy_id     = data.onefuse_ipam_policy.default.id // Refers to onefuse_ipam_policy data source to retrieve ID
  workspace_url = ""                                  // Leave blank for default workspace
  template_properties = {                             // Your properties and its values to pass into module
    property1 = "value1"
    property2 = "value2"
    property3 = "value3"
  }
}

// Outputs
output "hostname" {
  value = onefuse_ipam_record.my-ipam-record.computed_hostname
}
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