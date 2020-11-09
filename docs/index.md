# <provider> Provider

Summary of what the provider is for, including use cases and links to
app/service documentation.

Only uncomment the following declaration if using Terraform v0.13
Leave commented for Terraform v0.12

```hcl
terraform {
  required_providers {
    onefuse = {
    source = "CloudBoltSoftware/onefuse"
    version = ">= 1.10.0
    }
  }
required_version = ">= 0.13"
}

```

## Example Usage

```hcl
provider "onefuse" {
  address = "my-onefuse.example.com" //OneFuse Host
  port = "443" //OneFuse Port
  user = "admin" //OneFuse User with Workspace Admin or Member role
  password = "my-password" //OneFuse User's password
  scheme = "https" //OneFuse Protocol
  verify_ssl = false //Verify OneFuse SSL - true || false
}
```

## Argument Reference

- List any arguments for the provider block.
