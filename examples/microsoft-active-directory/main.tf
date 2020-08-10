provider "onefuse" {
  scheme      = var.onefuse_scheme
  address     = var.onefuse_address
  port        = var.onefuse_port
  user        = var.onefuse_user
  password    = var.onefuse_password
  verify_ssl  = var.onefuse_verify_ssl
}

data "onefuse_microsoft_endpoint" "my_microsoft_endpoint" {
    name = var.onefuse_microsoft_endpoint
}

resource "onefuse_microsoft_ad_policy" "my_ad_policy" {
    name = var.ad_policy_name
    description = var.ad_policy_description

    ou = var.ad_ou
    create_ou = var.ad_create_ou
    remove_ou = var.ad_remove_ou

    security_groups = var.ad_security_groups

    computer_name_letter_case = var.ad_computer_name_letter_case

    microsoft_endpoint_id = data.onefuse_microsoft_endpoint.my_microsoft_endpoint.id
    workspace_url = var.ad_workspace_url
}

resource "onefuse_microsoft_ad_computer_account" "my_ad_computer_account" {
    name = var.ad_computer_account_name

    policy_id = onefuse_microsoft_ad_policy.my_ad_policy.id
    workspace_url = var.ad_workspace_url
}
