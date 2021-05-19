// Comment out for Terraform 0.12
terraform {
  required_providers {
    onefuse = {
      source  = "CloudBoltSoftware/onefuse"
      version = ">= 1.10.3"
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

data "onefuse_ansible_tower_policy" "my_policy" {
  name = "my_ansible_tower_policy"
}

resource "onefuse_ansible_tower_deployment" "bar" {
  policy_id = data.onefuse_ansible_tower_policy.my_policy.id
  workspace_url = "" // Leave blank for default workspace
  limit = "*" // Ansible Tower Policy Limit
  hosts = [ "host1", "host2", ] // Hosts to run the policy against
  template_properties = {
        property1        = "value1" // Your properties and values to pass into module
        property2        = "value2"
  }
  timeouts {
    create = "12m"
    delete = "3m"
  }
}
