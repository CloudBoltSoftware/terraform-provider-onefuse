# Terraform Provider OneFuse

## [Unreleased]

The big features in this commit are the following:

* Added data source "onefuse_microsoft_endpoint"
* Added resource "onefuse_microsoft_ad_policy"
* Added resource "onefuse_microsoft_ad_computer_account"
* Added AD Poicy CRUD methods
* Added AD Computer Object CR_D (no Update) methods
* Added Microsoft Endpoint Read methods
* Added integration tests for new API Client components
* Re-added "scheme" Provider parameter, "http" or "https"
* Verfied the Provider supports Terraform 0.13
* Makefile got an updated to help with Terraform 0.13+ development
* Standaridize error message conventions format (mostly)

## 1.0.0

* API Client for OneFuse APIv3.
* Custom Naming resource type.
* Custom Naming Examples.
