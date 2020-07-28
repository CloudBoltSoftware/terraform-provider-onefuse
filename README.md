# terraform-provider-onefuse
Terraform Provider for integrating with SovLabs OneFuse.

## Sample Terraform Configuration
To get started with the Terraform Provider for SovLabs OneFuse, put the following into a file called `main.tf`.

Fill in the `provider "onefuse"` section with details about your SovLabs OneFuse instance.

```hcl
provider "onefuse" {
  address     = "localhost"
  port        = "8000"
  user        = "admin"
  password    = "my-password"
  scheme      = "http"
  verify_ssl  = false
}

resource "onefuse_naming" "my-onefuse-name" {
  naming_policy_id        = "2"
  dns_suffix              = "sovlabs.net"
  workspace_id            = "6"
  template_properties     = {
      "ownerName"               = "jsmith@company.com"
      "Environment"             = "dev"
      "OS"                      = "Linux"
      "Application"             = "Web Servers"
      "suffix"                  = "sovlabs.net"
      "tenant"                  =  "mytenant"
  }
}
```