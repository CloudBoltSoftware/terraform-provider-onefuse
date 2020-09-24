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
  template_properties     = var.onefuse_template_properties
}

resource "onefuse_microsoft_ad_computer_account" "my_ad_computer_account" {
    
    name = onefuse_naming.my-onefuse-name.name
    policy_id = var.onefuse_ad_policy_id
    workspace_url = var.workspace_url
}
