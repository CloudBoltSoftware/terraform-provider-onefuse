# Terraform Provider for OneFuse

Terraform Provider for integrating with [OneFuse](https://www.sovlabs.com/products/onefuse).

## Requirements

* [Terraform](https://www.terraform.io/downloads.html) 0.13.x
* [Go](https://golang.org/doc/install) >= 1.14 (to build the provider plugin)

*Note*: _Onefuse will drop support for Terraform 0.12.x after the release of Terraform Provider for OneFuse 1.0.1._

## Building the provider

Clone repository to: `$GOPATH/src/github.com/cloudboltsoftware/terraform-provider-onefuse`

```sh
$ mkdir -p $GOPATH/src/github.com/CloudBoltSoftware
$ cd $GOPATH/src/github.com/CloudBoltSoftware
$ git clone https://github.com/CloudBoltSoftware/terraform-provider-onefuse.git
```

Enter the provider directory and install the provider's dependencies

```sh
$ cd $GOPATH/src/github.com/CloudBoltSoftware/terraform-provider-onefuse.git
$ make install
```

To build the provider binary, create the `$HOME/.terraform.d/plugins/` directory and run the `make install` command:

```sh
$ make build
$ mkdir $HOME/.terraform.d/plugins/
$ mv terraform-provider-onefuse_v* $HOME/.terraform.d/plugins/
```

Then copy the binary to your terraform plugins directory

_You may want to use Make 4.3+ to ensure all make features work._

## Using the provider

### Sample Terraform Configuration

To get started with the Terraform Provider for OneFuse, put the following into a file called `main.tf`.

Fill in the `provider "onefuse"` section with details about your OneFuse instance.

```hcl
provider "onefuse" {
  address     = "localhost"
  port        = "8000"
  user        = "admin"
  password    = "my-password"
  scheme      = "https"
  verify_ssl  = false
}

resource "onefuse_naming" "my-onefuse-name" {
  naming_policy_id        = "2"
  dns_suffix              = "company.com"
  workspace_id            = "6"
  template_properties     = {
      "ownerName"               = "jsmith@company.com"
      "Environment"             = "dev"
      "OS"                      = "Linux"
      "Application"             = "Web Servers"
      "suffix"                  = "company.com"
      "tenant"                  =  "mytenant"
  }
}
```
## Releases
> To learn more, please visit our [docs](https://docs.cloudbolt.io/articles/onefuse-upstream-platforms-latest/hashicorp-terraform)
### v1.1
###### October 14, 2020
- OneFuse DNS module support
- OneFuse IPAM module support
- OneFuse Microsoft Active Directory module support

### v1.0
###### July 15, 2020
- OneFuse Naming module support

## Contributing

Interested in contributing? Wonderful!

* If you spot a problem, or room for improvement, please [create an issue][issue_url].
* If you are interested in fixing an issue, please [make a pull request][pr_url].
* A CloudBolt Developer will review your submission within a few days.

For more information about contributing to Terraform Provider Onefuse, see [CONTRIBUTING.md](CONTRIBUTING.md).

## License

Terraform Provider Onefuse is licensed under Mozilla Public License 2.0.
For more information see [LICENSE](LICENSE)

[issue_url]: https://github.com/CloudBoltSoftware/terraform-provider-onefuse/issues
[pr_url]: https://github.com/CloudBoltSoftware/terraform-provider-onefuse/pulls
