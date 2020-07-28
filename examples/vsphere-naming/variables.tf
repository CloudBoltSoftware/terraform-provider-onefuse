variable "onefuse_user" {
  type = string
}

variable "onefuse_password" {
  type = string
}

variable "onefuse_address" {
  type = string
}

variable "onefuse_port" {
  type = string
  default = "443"
}

variable "onefuse_naming_policy_id" {
  type = string
  default = 2
}

variable "onefuse_dns_suffix" {
  type = string
}

variable "vsphere_user" {
	type = string
}

variable "vsphere_password" {
	type = string
}

variable "vsphere_server" {
	type  = string
}

variable "onefuse_template_properties" {
  type = map
  default = {
      "nameEnv"               = "p"
      "nameOs"         	      = "w"
      "nameDatacenter"        = "por"
      "nameApp"               = "ap"
      "nameLocation"          = "atl"
  }
}

variable "onefuse_verify_ssl" {
  type = bool
  default = false
}
