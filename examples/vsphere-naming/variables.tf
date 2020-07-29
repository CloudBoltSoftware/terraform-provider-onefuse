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

variable "vsphere_datacenter" {
    type = string
    default = "replace_with_dc_name"
}

variable "vsphere_datastore_cluster" {
    type = string
    default = "replace_with_datastore_name"
}

variable "vsphere_compute_cluster" {
    type = string
    default = "replace_with_cluster_name"
}

variable "vsphere_network" {
    type = string
    default = "replace_with_portgroup_name"
}

variable "vsphere_virtual_machine" {
    type = string
    default = "replace_with_template_name"
}

variable "vsphereweb1_folder" {
    type = string
    default = "replace_with_folder_path"
}

variable "vsphereweb1_ip" {
    type = string
    default = "replace_with_op"
}

variable "vsphereweb1_gateway" {
    type = string
    default = "replace_with_gatway_addresst"
}

