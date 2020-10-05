provider "onefuse" {
  scheme      = var.onefuse_scheme
  address     = var.onefuse_address
  port        = var.onefuse_port
  user        = var.onefuse_user
  password    = var.onefuse_password
  verify_ssl  = var.onefuse_verify_ssl
}

provider "vsphere" {
  user           = var.vsphere_user
  password       = var.vsphere_password
  vsphere_server = var.vsphere_server
  version = "~> 1.20"

  # If you have a self-signed cert
  allow_unverified_ssl = true
}

###  Get OneFuse Resources ###

resource "onefuse_naming" "my-onefuse-name" {
  naming_policy_id        = var.onefuse_naming_policy_id
  dns_suffix              = ""
  template_properties     = var.onefuse_template_properties
}

resource "onefuse_microsoft_ad_computer_account" "my-ad-computer-account" {
    
    name = onefuse_naming.my-onefuse-name.name
    policy_id = var.onefuse_ad_policy_id
    workspace_url = var.workspace_url
    template_properties = var.onefuse_template_properties
}


resource "onefuse_ipam_record" "my-ipam-record" {
    
    hostname = "${format("%s.%s", onefuse_naming.my-onefuse-name.name, onefuse_naming.my-onefuse-name.dns_suffix)}"
    policy_id = var.onefuse_ipam_policy_id
    workspace_url = var.workspace_url
    template_properties = var.onefuse_template_properties
}

resource "onefuse_dns_record" "my-dns-record" {
    
    name = onefuse_naming.my-onefuse-name.name
    policy_id = var.onefuse_dns_policy_id
    workspace_url = var.workspace_url
    zones = var.onefuse_dns_zones
    value = onefuse_ipam_record.my-ipam-record.ip_address
    template_properties = var.onefuse_template_properties
}

###  vSphere Machine Deployment ###

#Data Sources
data "vsphere_datacenter" "dc" {
  name = "SovLabs"
}

data "vsphere_datastore_cluster" "datastore_cluster" {
  name          = "SovLabs_XtremIO"
  datacenter_id = data.vsphere_datacenter.dc.id
}
 
data "vsphere_compute_cluster" "cluster" {
  name          = "Cluster1"
  datacenter_id = data.vsphere_datacenter.dc.id
}
 
data "vsphere_network" "network" {
  name          = "dvs_SovLabs_329_10.30.29.0_24"
  datacenter_id = data.vsphere_datacenter.dc.id
}
 
data "vsphere_virtual_machine" "template" {
  name          = "CentOS7"
  datacenter_id = data.vsphere_datacenter.dc.id
}


#Virtual Machine Resource
resource "vsphere_virtual_machine" "vsphereweb1" {


    // Use OneFuse generated name for VM hostname and domain
    name = onefuse_naming.my-onefuse-name.name

  resource_pool_id = data.vsphere_compute_cluster.cluster.resource_pool_id
  datastore_cluster_id = data.vsphere_datastore_cluster.datastore_cluster.id
  folder = "VRM-BACKUPEXCLUDED/pre-sales-demo/"
 
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
        ipv4_address = onefuse_ipam_record.my-ipam-record.ip_address
        ipv4_netmask = 24
      }
 
      ipv4_gateway = var.onefuse_gateway
    }
  }
}
