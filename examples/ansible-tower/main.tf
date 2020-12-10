terraform {
  required_providers {
    onefuse = {
      source  = "CloudBoltSoftware/onefuse"
      version = ">= 1.10.3"
    }
  }
  required_version = ">= 0.13"
}

// Inititalize OneFuse Provider
provider "onefuse" {
  scheme     = "https"
  address    = "onefuse_fqdn"
  port       = "port"
  user       = "admin"
  password   = "admin"
  verify_ssl = "false"
}

/*
// OneFuse Ansible Tower Policy
data "onefuse_ansible_tower_policy" "foo" {
  name = "my-ansible-tower-policy-name"
}
*/

resource "onefuse_ansible_tower_deployment" "bar" {
  policy_id = 1 // Refers to onefuse_ansible_tower_deployment data source to retrieve ID
  workspace_url = "" // Leave blank for default workspace
  limit = "..." // Ansible Tower Policy Limit
  hosts = [ "host1", "host2", ] // Hosts to run the policy against
  template_properties = {
        property1        = "value1" // Your properties and values to pass into module
        proeprty2        = "value2"
  }
}
