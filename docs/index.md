# OneFuse Provider

OneFuse is the first codeless solution for automating, integrating, and extending private and hybrid cloud infrastructure. Through OneFuse’s dynamic templating technology, enterprises can build integrations (e.g., for IT technologies like IPAM, DNS, networking and security, etc.) into API-consumable policies for sharing across various IT teams and cloud environments. [Cloudbolt Documentation](https://docs.cloudbolt.io/)


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

* `address` - (Optional) OneFuse REST endpoint service port number

* `port` - (Required) OneFuse REST endpoint service port number

* `user` - (Required) OneFuse REST endpoint user name

* `password` - (Required) OneFuse REST endpoint password

* `scheme` - (Required) OneFuse REST endpoint service host address

* `verify_ss1` - (Required) Verify SSL certificates for OneFuse endpoints
