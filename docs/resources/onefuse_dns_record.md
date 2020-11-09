# <onefuse_dns_policy> Resource

Description of what this resource does, with links to official
app/service documentation.

## Example Usage

```hcl
// Code block with an example of how to use this resource.
// OneFuse Resource for DNS Record
resource "onefuse_dns_record" "my_dns_record" {
  name = "computer_name"                           // Required
  policy_id = data.onefuse_dns_policy.my_dns.id    // Required
  workspace_url = ""                               // Optional - Set to "" to use default
  zones = ["example.com,example1.com"]             // Required
  value = "10.1.1.1"                               // Required
  template_properties = {                          // Optional
    "Environment" = "development"
    "OS"          = "Linux"
    "Application" = "Web Servers"
    "suffix"      = "example.com"
  }
}
```

## Argument Reference

* `name` - (Required) The name of the dns record.
* 'policy_id' - (Required) 

## Attribute Reference

* `attribute_name` - List attributes that this resource exports.
