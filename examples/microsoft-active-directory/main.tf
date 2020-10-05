provider "onefuse" {
  scheme      = var.onefuse_scheme
  address     = var.onefuse_address
  port        = var.onefuse_port
  user        = var.onefuse_user
  password    = var.onefuse_password
  verify_ssl  = var.onefuse_verify_ssl
}

resource "onefuse_naming" "my-onefuse-name" {
  naming_policy_id        = var.onefuse_naming_policy_id
  dns_suffix              = ""
  template_properties     = var.onefuse_template_properties
}

resource "onefuse_microsoft_ad_computer_account" "my-ad-computer-account" {
    
    name = onefuse_naming.my-onefuse-name.name
    policy_id = var.onefuse_ad_policy_id
    workspace_url = var.workspace_url
    template_properties = var.onefuse_template_properties
}

resource "onefuse_ipam_record" "my-ipam-record" {
    
    hostname = "${format("%s.%s", onefuse_naming.my-onefuse-name.name, onefuse_naming.my-onefuse-name.dns_suffix)}"
    policy_id = var.onefuse_ipam_policy_id
    workspace_url = var.workspace_url
    template_properties = var.onefuse_template_properties
}

resource "onefuse_dns_record" "my-dns-record" {
    
    name = onefuse_naming.my-onefuse-name.name
    policy_id = var.onefuse_dns_policy_id
    workspace_url = var.workspace_url
    zones = var.onefuse_dns_zones
    value = onefuse_ipam_record.my-ipam-record.ip_address
    template_properties = var.onefuse_template_properties
}