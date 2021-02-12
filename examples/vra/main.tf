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
  scheme     = "http"
  address    = "localhost"
  port       = "8000"
  user       = "admin"
  password   = "admin"
  verify_ssl = "false"
}

resource "onefuse_vra_deployment" "vra_deployment_01" {
  policy_id = 1 // assumes you have a VRA policy with this ID
  workspace_url = "" // Leave blank for default workspace
  deployment_name = "tf_vra_deployment" // VRA Deployment Name
  template_properties = {
        property1        = "value1" // Your properties and values to pass into module
        property2        = "value2"
  }
  timeouts {
    create = "12m"
    delete = "3m"
  }
}

output "deployment_info" {
  value = onefuse_vra_deployment.vra_deployment_01.deployment_info
}
