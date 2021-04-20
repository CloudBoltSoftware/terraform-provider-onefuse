// Comment out for Terraform 0.12
terraform {
  required_providers {
    onefuse = {
      source  = "CloudBoltSoftware/onefuse"
      version = ">= 1.10.1"
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

// OneFuse Static Property Set
data "onefuse_static_property_set" "linux" {
  name = "linux"
}

// Policy Data Source: IPAM
data "onefuse_ipam_policy" "ipam_policy" {
  name = "infoblox851_ipampolicy"
}

// Policy Data Source: Naming
data "onefuse_naming_policy" "machine" {
  name = "machine"
}

// Policy Data Source: Active Directory
data "onefuse_ad_policy" "default" {
  name = "default"
}

// Policy Data Source: DNS
data "onefuse_dns_policy" "my_dns" {
  name = "infoblox851_dnspolicy"
}

// Name Object Resource
resource "onefuse_naming" "machine-name" {
  naming_policy_id  = data.onefuse_naming_policy.machine.id // Refers to onefuse_naming_policy data source to retrieve ID
  workspace_url     = ""                                    // Leave blank for default workspace
  dns_suffix        = ""
  template_properties = {                                   // Your properties and its values to pass into module
    property1 = "value1"
    property2 = "value2"
    property3 = data.onefuse_static_property_set.linux.properties.myPropName // Replace "myPropName" with the key name of the property defined in Static Property Set
  }
}

// AD Computer Object Resource
resource "onefuse_microsoft_ad_computer_account" "dev" {
  name          = onefuse_naming.machine-name.name      // Refers to onefuse_naming_policy for computer name
  policy_id     = data.onefuse_ad_policy.default.id     // Refers to onefuse_ad_policy data source to retrieve ID
  workspace_url = ""                                    // Leave blank for default workspace
  template_properties = {                               // Your properties and its values to pass into module
    property1 = "value1"
    property2 = "value2"
    property3 = "value3"
  }
}

// IP Address Object Resource
resource "onefuse_ipam_record" "my-ipam-record" {
  hostname      = onefuse_naming.machine-name.name  // Refers to onefuse_naming_resource for computer name
  policy_id     = data.onefuse_ipam_policy.dev.id   // Refers to onefuse_ipam_policy data source to retrieve ID
  workspace_url = ""                                // Leave blank for default workspace
  template_properties = {                           // Your properties and its values to pass into module
    property1 = "value1"
    property2        = "value2"
    property3        = "value3"
  }
}

// DNS Record Object Resource
resource "onefuse_dns_record" "my-dns-record" {
  name          = onefuse_naming.machine-name.name              // Refers to onefuse_naming resource for computer name
  policy_id     = data.onefuse_dns_policy.dev.id                // Refers to onefuse_dns_policy data source to retrieve ID
  workspace_url = ""                                            // Leave blank for default workspace
  zones         = [onefuse_naming.machine-name.dns_suffix]      // Comma separated list of zones, example gets zone from Naming Policy
  value         = onefuse_ipam_record.my-ipam-record.ip_address // Refers to onefuse_ipam resource for computer name
  template_properties = {                                       // Your properties and its values to pass into module
    property1 = "value1"
    property2        = "value2"
    property3        = "value3"
  }
}
