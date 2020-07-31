# Terraform Provider for OneFuse

Terraform Provider for integrating with [SovLabs OneFuse](https://www.sovlabs.com/products/onefuse).

## Requirements

* [Terraform](https://www.terraform.io/downloads.html) 0.12.x
* [Go](https://golang.org/doc/install) >= 1.14 (to build the provider plugin)

## Building the provider

Clone repository to: `$GOPATH/src/github.com/cloudboltsoftware/terraform-provider-onefuse`

```sh
$ mkdir -p $GOPATH/src/github.com/CloudBoltSoftware
$ cd $GOPATH/src/github.com/CloudBoltSoftware
$ git clone https://github.com/CloudBoltSoftware/terraform-provider-onefuse.git
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/CloudBoltSoftware/terraform-provider-onefuse.git
$ make build
```

To install the provider, create the `$HOME/.terraform.d/plugins/` directory and run the `make install` command:

```sh
$ mkdir $HOME/.terraform.d/plugins/
$ make install
```

Then copy the binary to your terraform plugins directory

## Using the provider

### Sample Terraform Configuration

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
