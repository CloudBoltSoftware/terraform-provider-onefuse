// Provider setup

variable "onefuse_scheme" {
  type = string
  default = "https"
}

variable "onefuse_address" {
  type = string
  default = "se-onefuse-dev.sovlabs.net"
}

variable "onefuse_port" {
  type = string
  default = "8080"
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

variable "onefuse_naming_policy_id" {
  type = string
  default = "1"
}

variable "onefuse_ad_policy_id" {
  type = string
  default = "1"
}

variable "workspace_url" {
  type = string
  default = "" // Default
}


variable "onefuse_template_properties" {
  type = map
  default = {
      "nameEnv"               = "prod"
      "nameOs"                = "w"
      "nameDatacenter"        = "por"
      "nameApp"               = "web"
      "nameDomain"            = "sovlabs.net"
      "nameLocation"          = "atl"
      "testOU"	              = "build"
  }
}