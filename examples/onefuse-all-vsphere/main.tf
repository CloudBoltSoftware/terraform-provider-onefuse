// Commented out for Terraform 0.12

 terraform {
  required_providers {
    onefuse = {
      source = "CloudBoltSoftware/onefuse"
      version = ">= 1.10.1"
    }
  }
  required_version = ">= 0.13"
}

// Comment out above for Terraform 0.12


// Inititalize OneFuse Provider
provider "onefuse" {

  scheme     = "https"
  address    = "onefuse_fqdn"
  port       = "port"
  user       = "admin"
  password   = "admin"
  verify_ssl = "false"
}

// OneFuse Static Property Set
data "onefuse_static_property_set" "linux" {
  name = "linux"
}

// IPAM Policy data source
data "onefuse_ipam_policy" "ipam_policy" {
  name = "infoblox851_ipampolicy"
}

// Naming Policy data source
data "onefuse_naming_policy" "machine" {
  name = "machineNaming"
}

// AD Policy data source
data "onefuse_ad_policy" "default" {
  name = "default"
}

// DNS Policy data source
data "onefuse_dns_policy" "my_dns" {
  name = "infoblox851_dnspolicy"
}

resource "onefuse_naming" "machine-name" {
  naming_policy_id        = data.onefuse_naming_policy.machine.id // Refers to onefuse_naming_policy data source to retrieve ID
  workspace_url = "" // Leave blank for default workspace
  dns_suffix              = ""
  template_properties = {
        property1        = "value1" // Your properties and values to pass into module
        proeprty2        = "value2"
        property3        = data.onefuse_static_property_set.linux.properties.{{propName}} // Reference value defined in Static Property Set
  }
}

resource "onefuse_microsoft_ad_computer_account" "dev" {
    
    name = onefuse_naming.machine-name.name // Refers to onefuse_naming_policy for computer name
    policy_id = data.onefuse_ad_policy.default.id // Refers to onefuse_ad_policy data source to retrieve ID
    workspace_url = "" // Leave blank for default workspace
    template_properties = {
        property1        = "value1" // Your properties and values to pass into module
        proeprty2        = "value2"
        property3        = "value3"
  }
}

resource "onefuse_ipam_record" "my-ipam-record" {
    
    hostname = onefuse_naming.machine-name.name  // Refers to onefuse_naming_resource for computer name
    policy_id = data.onefuse_ipam_policy.dev.id // Refers to onefuse_ipam_policy data source to retrieve ID
    workspace_url = "" // Leave blank for default workspace
    template_properties = {
        property1        = "value1" // Your properties and values to pass into module
        proeprty2        = "value2"
        property3        = "value3"
  }
}

resource "onefuse_dns_record" "my-dns-record" {
    
    name = onefuse_naming.machine-name.name // Refers to onefuse_naming resource for computer name
    policy_id = data.onefuse_dns_policy.dev.id // Refers to onefuse_dns_policy data source to retrieve ID
    workspace_url = "" // Leave blank for default workspace
    zones = [onefuse_naming.machine-name.dns_suffix] // Comma seperated list of zones, example grabbing zone from naming policy
    value = onefuse_ipam_record.my-ipam-record.ip_address // Refers to onefuse_ipam resource for computer name
    template_properties = {
        property1        = "value1" // Your properties and values to pass into module
        proeprty2        = "value2"
        property3        = "value3"
  }
}
