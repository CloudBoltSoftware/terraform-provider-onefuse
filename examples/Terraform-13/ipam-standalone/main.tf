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

data "onefuse_ipam_policy" "ipam_policy" {
  name = "infoblox851_ipampolicy"
}

resource "onefuse_ipam_record" "my-ipam-record" {
    
    hostname = "testname"
    policy_id = data.onefuse_ipam_policy.ipam_policy.id
    workspace_url = var.workspace_url
    template_properties = var.onefuse_template_properties
}
