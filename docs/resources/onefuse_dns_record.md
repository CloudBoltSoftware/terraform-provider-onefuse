# Resource: onefuse_dns_policy

Use this resource to create a DNS record.

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

* `name` - (Required) The short name of the computer object.

* `policy_id` - (Required) The id of the policy object in OneFuse (add example of format) 

* `workspace_url` - (Optional) The URL of the workspace being used in OneFuse

* `zones` - (Required) An array of DNS zones

* `value` - (Required) The value of the DNS record (i,e. IP address)

* `template_properties` - (Optional) Additional properties that can be pushed to OneFuse and referenced within the policy

## Attribute Reference

* `workspace_url` - Value of default Workspace URL, if no URL is provided
