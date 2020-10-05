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
  default = "2"
}

variable "onefuse_dns_policy_id" {
  type = string
  default = "1"
}

variable "onefuse_dns_zones" {
  type = list
  default = ["w2k19ad.sovlabs.net"]
}

variable "onefuse_dns_ip" {
  type = string
  default = "10.30.6.221"
}

variable "onefuse_gateway" {
  type = string
  default = "10.30.29.1"
}

variable "workspace_url" {
  type = string
  default = "" // Default
}

variable "onefuse_ipam_policy_id" {
  type = string
  default = "2"
}

variable "onefuse_template_properties" {
  type = map
  default = {
      "nameEnv"               = "dev"
      "nameOs"                = "w"
      "nameDatacenter"        = "por"
      "nameApp"               = "web"
      "nameDomain"            = "sovlabs.net"
      "nameLocation"          = "atl"
      "testOU"	              = "sidtest"
  }
}

variable "vsphere_user" {
        type = string
        default = "vrasvc@sovlabs.net"
}

variable "vsphere_password" {
        type = string
        default = "VmwareS0v"
}

variable "vsphere_server" {
        type  = string
        default = "vcenter01.sovlabs.net"
}
