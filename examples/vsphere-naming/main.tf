provider "onefuse" {
  address     = var.onefuse_address
  port        = var.onefuse_port
  user        = var.onefuse_user
  password    = var.onefuse_password
  verify_ssl  = var.onefuse_verify_ssl
}

resource "onefuse_naming" "my-onefuse-name" {
  naming_policy_id        = var.onefuse_naming_policy_id
  dns_suffix              = var.onefuse_dns_suffix
  template_properties     = var.onefuse_template_properties
}

provider "vsphere" {
  user           = var.vsphere_user
  password       = var.vsphere_password
  vsphere_server = var.vsphere_server
  version = "~> 1.20"

  # If you have a self-signed cert
  allow_unverified_ssl = true
}

#Data Sources
data "vsphere_datacenter" "dc" {
  name = "replace_with_dc_name"
}

data "vsphere_datastore_cluster" "datastore_cluster" {
  name          = "replace_with_datastor_name"
  datacenter_id = data.vsphere_datacenter.dc.id
}
 
data "vsphere_compute_cluster" "cluster" {
  name          = "replace_with_cluster_name"
  datacenter_id = data.vsphere_datacenter.dc.id
}
 
data "vsphere_network" "network" {
  name          = "replace_with_portgroup_name"
  datacenter_id = data.vsphere_datacenter.dc.id
}
 
data "vsphere_virtual_machine" "template" {
  name          = "replace_with_template_name"
  datacenter_id = data.vsphere_datacenter.dc.id
}


#Virtual Machine Resource
resource "vsphere_virtual_machine" "vsphereweb1" {


  // Use OneFuse generated name for VM hostname and domain
  name = onefuse_naming.my-onefuse-name.name
	
  resource_pool_id = data.vsphere_compute_cluster.cluster.resource_pool_id
  datastore_cluster_id = data.vsphere_datastore_cluster.datastore_cluster.id
  folder = "replace_with_folder_path"
 
  num_cpus = 1
  memory   = 512
  guest_id = data.vsphere_virtual_machine.template.guest_id
 
  scsi_type = data.vsphere_virtual_machine.template.scsi_type
 
  network_interface {
    network_id   = data.vsphere_network.network.id
    adapter_type = "vmxnet3"
  }
 
  disk {
    label            = "disk0"
    size             = data.vsphere_virtual_machine.template.disks.0.size
    eagerly_scrub    = data.vsphere_virtual_machine.template.disks.0.eagerly_scrub
    thin_provisioned = data.vsphere_virtual_machine.template.disks.0.thin_provisioned
  }
 
  clone {
    template_uuid = data.vsphere_virtual_machine.template.id
 
    customize {
      linux_options {
        host_name  = onefuse_naming.my-onefuse-name.name
        domain = onefuse_naming.my-onefuse-name.dns_suffix
      }
 
      network_interface {
        ipv4_address = "replace_with_ip"
        ipv4_netmask = 24
      }
 
      ipv4_gateway = "replace_with_gateway_address"
    }
  }
}
