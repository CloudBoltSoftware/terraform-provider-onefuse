# Data Source: onefuse_dns_policy

Use this resource to lookup Policy ID by Policy Name.

## Example Usage

```hcl
data "onefuse_dns_policy" "dns_policy" {
  name = "my_dnspolicy_name"                       // Replace with Policy Name
}
```

## Argument Reference

* `name` - (Required) The name of the DNS policy

## Attribute Reference

* `ID` - ID of the DNS policy

* `description` - The description of the DNS policy
