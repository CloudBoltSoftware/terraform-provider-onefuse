// Provider setup

variable "onefuse_scheme" {
  type = string
  default = "https"
}

variable "onefuse_address" {
  type = string
  default = "localhost"
}

variable "onefuse_port" {
  type = string
  default = "443"
}

variable "onefuse_user" {
  type = string
  default = "admin"
}

variable "onefuse_password" {
  type = string
  default = "admin"
}

variable "onefuse_verify_ssl" {
  type = bool
  default = false
}

// Microsoft AD Endpoint

variable "onefuse_microsoft_endpoint" {
  type = string
  default = "myMSEndpoint"
}

// Microsoft AD Policy

variable "ad_policy_name" {
  type = string
  default = "myADPolicy01"
}

variable "ad_policy_description" {
  type = string
  default = "Created with Terraform"
}

variable "ad_computer_name_letter_case" {
  type = string
  default = "Lowercase"
}

variable "ad_ou" {
  type = string
  default = "OU=Accounting,DC=yourOrg,DC=com"
}

variable "ad_workspace_url" {
  type = string
  default = "" // Default
}

variable "ad_security_groups" {
  type = list(string)
  default = ["CN=SomeSecurityGroup,OU=Accounting,DC=yourOrg,DC=com"]
  
}

variable "ad_create_ou" {
  type = bool
  default = true
}

variable "ad_remove_ou" {
  type = bool
  default = true
}

// Microsoft AD Object
variable "ad_computer_account_name" {
  type = string
  default = "someComputerName"
}
