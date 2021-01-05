// Commented out for Terraform 0.12

terraform {
  required_providers {
    onefuse = {
      source  = "CloudBoltSoftware/onefuse"
      version = ">= 1.1.0"
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

// Static Property Set
// We expect a Static Property Set called SPS01 with the following properties:
// {
//   "someKey": "{{templatedValue}}",
//   "parent": {
//     "someNestedKey": "someNestedValue"
//   }
// }
//
data "onefuse_static_property_set" "sps01" {
    name = "SPS01"
}

// Only top-level key:value (string:string) pairs
output "some_value" {
  value = data.onefuse_static_property_set.sps01.properties
}

// All key:value pairs, arbitrary nesting allowed
output "raw_value" {
  value = data.onefuse_static_property_set.sps01.raw
}

locals  {
  some_nested_value = jsondecode(data.onefuse_static_property_set.sps01.raw).parent.someNestedKey
}

output "some_nested_value" {
  value = local.some_nested_value
}

// Render the template
data "onefuse_rendered_template" "rendered_template" {
    template = data.onefuse_static_property_set.sps01.raw
    template_properties = {"templatedValue": "rendered_template"}
}

output "rendered_template" {
  value = data.onefuse_rendered_template.rendered_template.value
}
