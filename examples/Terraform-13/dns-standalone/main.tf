terraform {
  required_providers {
    onefuse = {
      source  = "cloudbolt.io/cloudbolt/onefuse"
      version = ">= 1.1.0"
    }
  }
  required_version = ">= 0.13"
}

provider "onefuse" {

  scheme     = var.onefuse_scheme
  address    = var.onefuse_address
  port       = var.onefuse_port
  user       = var.onefuse_user
  password   = var.onefuse_password
  verify_ssl = var.onefuse_verify_ssl
}

data "onefuse_dns_policy" "my_dns" {
  name = "infoblox851_dnspolicy"
}

resource "onefuse_dns_record" "my-dns-record" {
    
    name = "hostname"
    policy_id = data.onefuse_dns_policy.my_dns.id
    workspace_url = var.workspace_url
    zones = ["dnszone1,dnszone2"]
    value = "ipAddress"
    template_properties = var.onefuse_template_properties
}
