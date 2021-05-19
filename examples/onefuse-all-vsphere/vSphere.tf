provider "vsphere" {
  user           = "username"
  password       = "password"
  vsphere_server = "vCenter_Address"
  version = "~> 1.20"

  # If you have a self-signed cert
  allow_unverified_ssl = true
}

###  vSphere Machine Deployment ###

#Data Sources
data "vsphere_datacenter" "datacenter {
  name = "datacenter01"
}

data "vsphere_datastore_cluster" "datastore_cluster" {
  name          = "datastore"
  datacenter_id = data.vsphere_datacenter.dc.id
}
 
data "vsphere_compute_cluster" "cluster" {
  name          = "cluster"
  datacenter_id = data.vsphere_datacenter.dc.id
}
 
data "vsphere_network" "network" {
  name          = onefuse_ipam_record.my-ipam-record.network // Assign network from OneFuse ipam resource
  datacenter_id = data.vsphere_datacenter.dc.id
}
 
data "vsphere_virtual_machine" "template" {
  name          = data.onefuse_static_property_set.linux.properties.{{template}} //Assign template name from definition within static property set
  datacenter_id = data.vsphere_datacenter.dc.id
}


#Virtual Machine Resource
resource "vsphere_virtual_machine" "vsphereweb1" {

  // Use OneFuse generated name for VM hostname and domain
  name = onefuse_naming.machine-name.name

  resource_pool_id = data.vsphere_compute_cluster.cluster.resource_pool_id
  datastore_cluster_id = data.vsphere_datastore_cluster.datastore_cluster.id
  folder = "vCenter/folder/path"
 
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
        host_name  = onefuse_naming.machine-name.name                 //  Assign name from OneFuse naming resource
        domain = onefuse_naming.machine-name.dns_suffix               // Assign DNS Suffix from OneFuse naming resource
      }
 
      network_interface {
        ipv4_address = onefuse_ipam_record.my-ipam-record.ip_address // Assign IP Address from OneFuse ipam resource
        ipv4_netmask = 24
      }
 
      ipv4_gateway = onefuse_ipam_record.my-ipam-record.gateway     // Assign gateway from OneFuse ipam resource
    }
  }
}
