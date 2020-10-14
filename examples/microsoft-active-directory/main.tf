

data "onefuse_static_property_set" "linux" {
  name = "linux"
}

data "onefuse_ipam_policy" "dev" {
  name = "infoblox851_ipampolicy"
}

data "onefuse_naming_policy" "machine" {
  name = "machineNaming"
}

data "onefuse_naming_policy" "deployment" {
  name = "deploymentNaming"
}

data "onefuse_ad_policy" "default" {
  name = "defaultAdPolicy"
}

data "onefuse_dns_policy" "dev" {
  name = "infoblox851_dnspolicy"
}

resource "onefuse_naming" "machine-name" {
  count                   = "2"
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

resource "onefuse_naming" "deployment-name" {
  naming_policy_id        = data.onefuse_naming_policy.deployment.id
  dns_suffix              = ""
  template_properties = {
      deployNameRequestSource     = "TF"
      deployNameEnv               = "PROD"
      deployNameApp               = "WEB"
  }
}

resource "onefuse_microsoft_ad_computer_account" "dev" {
    
    name = onefuse_naming.machine-name[0].name
    policy_id = data.onefuse_ad_policy.default.id
    workspace_url = var.workspace_url
    template_properties = var.onefuse_template_properties
}

resource "onefuse_ipam_record" "my-ipam-record" {
    
    hostname = format("%s.%s", onefuse_naming.machine-name[0].name, onefuse_naming.machine-name[0].dns_suffix)
    policy_id = data.onefuse_ipam_policy.dev.id
    workspace_url = var.workspace_url
    template_properties = var.onefuse_template_properties
}

resource "onefuse_dns_record" "my-dns-record" {
    
    name = onefuse_naming.machine-name[0].name
    policy_id = data.onefuse_dns_policy.dev.id
    workspace_url = var.workspace_url
    zones = [onefuse_naming.machine-name[0].dns_suffix]
    value = onefuse_ipam_record.my-ipam-record.ip_address
    template_properties = var.onefuse_template_properties
}


output "hostname_machine_1" {
  value = onefuse_naming.machine-name[0].name
}

output "deployment_name" {
  value = onefuse_naming.deployment-name.name
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
  value = format("%s.%s", onefuse_naming.machine-name[0].name, onefuse_naming.machine-name[0].dns_suffix)
}

output "ad_ou" {
  value = onefuse_microsoft_ad_computer_account.dev.final_ou
}

output "hostname_machine_2" {
  value = onefuse_naming.machine-name[1].name
}

