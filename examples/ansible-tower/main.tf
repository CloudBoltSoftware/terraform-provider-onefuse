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

// Ansible Tower Deployment Resource
resource "onefuse_ansible_tower_deployment" "bar" {
  policy_id     = 1                                 // Refers to Ansible Tower Policy ID (integer)
  workspace_url = ""                                // Leave blank for default workspace
  limit         = "..."                             // Ansible Tower Policy Limit
  hosts         = [ "host1", "host2" ]              // Hosts to run the Policy against
  template_properties = {                           // Your properties and its values to pass into module
    property1 = "value1"
    property2 = "value2"
    property3 = "value3"
  }
  timeouts {
    create = "12m"
    delete = "3m"
  }
}
