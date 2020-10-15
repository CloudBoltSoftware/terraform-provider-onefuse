// Uncomment the following declaration for Terraform 0.13, leave commented for Terraform 0.12
//
// terraform {
//  required_providers {
//    onefuse = {
//      source  = "cloudbolt.io/cloudbolt/onefuse"
//      version = ">= 1.1.0"
//    }
//  }
//  required_version = ">= 0.13"
//}

provider "onefuse" {

  scheme     = var.onefuse_scheme
  address    = var.onefuse_address
  port       = var.onefuse_port
  user       = var.onefuse_user
  password   = var.onefuse_password
  verify_ssl = var.onefuse_verify_ssl
}


data "onefuse_static_property_set" "linux" {
  name = "linux"
}

data "onefuse_ipam_policy" "dev" {
  name = "infoblox851_ipampolicy"
}

data "onefuse_naming_policy" "machine" {
  name = "machineNaming"
}

data "onefuse_ad_policy" "default" {
  name = "default"
}

data "onefuse_dns_policy" "dev" {
  name = "infoblox851_dnspolicy"
}

resource "onefuse_naming" "machine-name" {
  naming_policy_id        = data.onefuse_naming_policy.machine.id
  dns_suffix              = ""
  template_properties = {
      nameEnv               = "dev"
      nameOs                = data.onefuse_static_property_set.linux.properties.nameOs
      nameDatacenter        = "por"
      nameApp               = "web"
      nameDomain            = "sovlabs.net"
      nameLocation          = "atl"
      testOU	              = "sidtest"
  }
}

resource "onefuse_microsoft_ad_computer_account" "dev" {
    
    name = onefuse_naming.machine-name.name
    policy_id = data.onefuse_ad_policy.default.id
    workspace_url = var.workspace_url
    template_properties = var.onefuse_template_properties
}

resource "onefuse_ipam_record" "my-ipam-record" {
    
    hostname = onefuse_naming.machine-name.name
    policy_id = data.onefuse_ipam_policy.dev.id
    workspace_url = var.workspace_url
    template_properties = var.onefuse_template_properties
}

resource "onefuse_dns_record" "my-dns-record" {
    
    name = onefuse_naming.machine-name.name
    policy_id = data.onefuse_dns_policy.dev.id
    workspace_url = var.workspace_url
    zones = [onefuse_naming.machine-name.dns_suffix]
    value = onefuse_ipam_record.my-ipam-record.ip_address
    template_properties = var.onefuse_template_properties
}


output "hostname" {
  value = onefuse_naming.machine-name.name
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

output "fqdn" {
  value = format("%s.%s", onefuse_naming.machine-name.name, onefuse_naming.machine-name.dns_suffix)
}

output "ad_ou" {
  value = onefuse_microsoft_ad_computer_account.dev.final_ou
}

