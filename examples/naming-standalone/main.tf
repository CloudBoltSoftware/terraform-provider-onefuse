provider "onefuse" {
  address     = var.onefuse_address
  port        = var.onefuse_port
  user        = var.onefuse_user
  password    = var.onefuse_password
  verify_ssl  = var.onefuse_verify_ssl
}

resource "onefuse_naming" "my-onefuse-name" {
  naming_policy_id        = var.onefuse_naming_policy_id
  dns_suffix              = var.onefuse_dns_suffix
  template_properties     = var.onefuse_template_properties
}
