---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_vm"
sidebar_current: "outscale-vm"
description: |-
  [Manages a virtual machine (VM).]
---

# outscale_vm Resource

Manages a virtual machine (VM).
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Instances).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-vm).

## Example Usage

```hcl

# Create a VM in the Public Cloud

resource "outscale_vm" "vm01" {
  image_id                 = var.image_id
  vm_type                  = var.vm_type
  keypair_name             = var.keypair_name
  security_group_ids       = [var.security_group_id]
  placement_subregion_name = "eu-west-2a"
  placement_tenancy        = "default"
  tags {
    key   = "name"
    value = "terraform-public-vm"
  }
}

# Create a VM with block device mappings

resource "outscale_vm" "vm02" {
  image_id                = var.image_id
  vm_type                 = var.vm_type
  keypair_name            = var.keypair_name
  block_device_mappings {
    device_name = "/dev/sdb"
    bsu  {
      volume_size = 15
      volume_type = "gp2"
      snapshot_id = var.snapshot_id
    }
  }  
  block_device_mappings {
    device_name = "/dev/sdc"
    bsu  {
      volume_size           = 22
      volume_type           = "io1"
      iops                  = 150
      delete_on_vm_deletion = true
    }
  }
}

# Create a VM in the Private Cloud

resource "outscale_net" "net01" {
  ip_range = "10.0.0.0/16"
  tags  {
    key   = "name"
    value = "terraform-net-for-vm"
  }
}

resource "outscale_subnet" "subnet01" {
    net_id         = outscale_net.net01.net_id
    ip_range       = "10.0.0.0/24"
    subregion_name = "eu-west-2b"
    tags {
        key   = "name"
        value = "terraform-subnet-for-vm"
      }
}

resource "outscale_security_group" "security_01" {
  description         = "Terraform security group for VM"
  security_group_name = "terraform-security-group-for-vm"
  net_id              = outscale_net.net01.net_id
}

resource "outscale_internet_service" "internet_service01" {
}

resource "outscale_route_table" "route_table01" {
  net_id = outscale_net.net01.net_id
  tags {
    key   = "name"
    value = "terraform-route-table-for-vm"
  }
}

resource "outscale_route_table_link" "route_table_link01" {
  route_table_id = outscale_route_table.route_table01.route_table_id
  subnet_id      = outscale_subnet.subnet01.subnet_id
}

resource "outscale_internet_service_link" "internet_service_link01" {
  internet_service_id = outscale_internet_service.internet_service01.internet_service_id
  net_id              = outscale_net.net01.net_id
}

resource "outscale_route" "route01" {
  gateway_id           = outscale_internet_service.internet_service01.internet_service_id
  destination_ip_range = "0.0.0.0/0"
  route_table_id       = outscale_route_table.route_table01.route_table_id
}

resource "outscale_vm" "vm03" {
  image_id           = var.image_id
  vm_type            = var.vm_type
  keypair_name       = var.keypair_name
  security_group_ids = [outscale_security_group.security_group01.security_group_id]
  subnet_id          = outscale_subnet.subnet01.subnet_id
}

# Create a VM with a NIC

resource "outscale_net" "net02" {
   ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet02" {
  net_id         = outscale_net.net02.net_id
  ip_range       = "10.0.0.0/24"
  subregion_name = "eu-west-2a"
}

resource "outscale_nic" "nic01" {
  subnet_id = outscale_subnet.subnet02.subnet_id
}

resource "outscale_vm" "vm04" {
  image_id     = var.image_id
  vm_type      = "c4.large"
  keypair_name = var.keypair_name
  nics {
    nic_id        = outscale_nic.nic01.nic_id
    device_number = "0"
  }
}


```

## Argument Reference

The following arguments are supported:

* `block_device_mappings` - (Optional) One or more block device mappings.
  * `bsu` - Information about the BSU volume to create.
    * `delete_on_vm_deletion` - (Optional) Set to `true` by default, which means that the volume is deleted when the VM is terminated. If set to `false`, the volume is not deleted when the VM is terminated.
    * `iops` - (Optional) The number of I/O operations per second (IOPS). This parameter must be specified only if you create an `io1` volume. The maximum number of IOPS allowed for `io1` volumes is `13000`.
    * `snapshot_id` - (Optional) The ID of the snapshot used to create the volume.
    * `volume_size` - (Optional) The size of the volume, in gibibytes (GiB).<br />
If you specify a snapshot ID, the volume size must be at least equal to the snapshot size.<br />
If you specify a snapshot ID but no volume size, the volume is created with a size similar to the snapshot one.
    * `volume_type` - (Optional) The type of the volume (`standard` \| `io1` \| `gp2`). If not specified in the request, a `standard` volume is created.<br />
For more information about volume types, see [Volume Types and IOPS](https://wiki.outscale.net/display/EN/About+Volumes#AboutVolumes-VolumeTypesVolumeTypesandIOPS).
  * `device_name` - (Optional) The name of the device.
  * `no_device` - (Optional) Removes the device which is included in the block device mapping of the OMI.
  * `virtual_device_name` - (Optional) The name of the virtual device (ephemeralN).
* `boot_on_creation` - (Optional) By default or if `true`, the VM is started on creation. If `false`, the VM is stopped on creation.
* `bsu_optimized` - (Optional) If `true`, the VM is created with optimized BSU I/O.
* `client_token` - (Optional) A unique identifier which enables you to manage the idempotency.
* `deletion_protection` - (Optional) If `true`, you cannot terminate the VM using Cockpit, the CLI or the API. If `false`, you can.
* `image_id` - (Required) The ID of the OMI used to create the VM. You can find the list of OMIs by calling the [ReadImages](https://docs.outscale.com/api#readimages) method.
* `keypair_name` - (Optional) The name of the keypair.
* `max_vms_count` - (Optional) The maximum number of VMs you want to create. If all the VMs cannot be created, the largest possible number of VMs above MinVmsCount is created.
* `min_vms_count` - (Optional) The minimum number of VMs you want to create. If this number of VMs cannot be created, no VMs are created.
* `nics` - (Optional) One or more NICs. If you specify this parameter, you must define one NIC as the primary network interface of the VM with `0` as its device number.
  * `delete_on_vm_deletion` - (Optional) If `true`, the NIC is deleted when the VM is terminated. You can specify `true` only if you create a NIC when creating a VM.
  * `description` - (Optional) The description of the NIC, if you are creating a NIC when creating the VM.
  * `device_number` - (Optional) The index of the VM device for the NIC attachment (between 0 and 7, both included). This parameter is required if you create a NIC when creating the VM.
  * `nic_id` - (Optional) The ID of the NIC, if you are attaching an existing NIC when creating a VM.
  * `private_ips` - (Optional) One or more private IP addresses to assign to the NIC, if you create a NIC when creating a VM. Only one private IP address can be the primary private IP address.
    * `is_primary` - (Optional) If `true`, the IP address is the primary private IP address of the NIC.
    * `private_ip` - (Optional) The private IP address of the NIC.
  * `secondary_private_ip_count` - (Optional) The number of secondary private IP addresses, if you create a NIC when creating a VM. This parameter cannot be specified if you specified more than one private IP address in the `PrivateIps` parameter.
  * `security_group_ids` - (Optional) One or more IDs of security groups for the NIC, if you acreate a NIC when creating a VM.
  * `subnet_id` - (Optional) The ID of the Subnet for the NIC, if you create a NIC when creating a VM.
* `performance` - (Optional) The performance of the VM (`standard` \| `high` \|  `highest`).
* `placement_subregion_name` - (Optional) The name of the Subregion where the VM is placed.
* `placement_tenancy` - (Optional) The tenancy of the VM (`default` | `dedicated`).      
* `private_ips` - (Optional) One or more private IP addresses of the VM.
* `security_group_ids` - (Optional) One or more IDs of security group for the VMs.
* `security_group_names` - (Optional) One or more names of security groups for the VMs.
* `subnet_id` - (Optional) The ID of the Subnet in which you want to create the VM.
* `user_data` - (Optional) Data or script used to add a specific configuration to the VM. It must be base64-encoded.
* `vm_initiated_shutdown_behavior` - (Optional) The VM behavior when you stop it. By default or if set to `stop`, the VM stops. If set to `restart`, the VM stops then automatically restarts. If set to `terminate`, the VM stops and is terminated.
* `vm_type` - (Optional) The type of VM (`tinav2.c1r2` by default).<br />
For more information, see [Instance Types](https://wiki.outscale.net/display/EN/Instance+Types).
* `tags` - One or more tags to add to this resource.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
    
## Attribute Reference

The following attributes are exported:

* `vms` - Information about one or more created VMs.
  * `architecture` - The architecture of the VM (`i386` \| `x86_64`).
  * `block_device_mappings_created` - The block device mapping of the VM.
    * `bsu` - Information about the created BSU volume.
      * `delete_on_vm_deletion` - Set to `true` by default, which means that the volume is deleted when the VM is terminated. If set to `false`, the volume is not deleted when the VM is terminated.
      * `link_date` - The time and date of attachment of the volume to the VM.
      * `state` - The state of the volume.
      * `volume_id` - The ID of the volume.
    * `device_name` - The name of the device.
  * `bsu_optimized` - If `true`, the VM is optimized for BSU I/O.
  * `client_token` - The idempotency token provided when launching the VM.
  * `deletion_protection` - If `true`, you cannot terminate the VM using Cockpit, the CLI or the API. If `false`, you can.
  * `hypervisor` - The hypervisor type of the VMs (`ovm` \| `xen`).
  * `image_id` - The ID of the OMI used to create the VM.
  * `is_source_dest_checked` - (Net only) If `true`, the source/destination check is enabled. If `false`, it is disabled. This value must be `false` for a NAT VM to perform network address translation (NAT) in a Net.
  * `keypair_name` - The name of the keypair used when launching the VM.
  * `launch_number` - The number for the VM when launching a group of several VMs (for example, 0, 1, 2, and so on).
  * `net_id` - The ID of the Net in which the VM is running.
  * `nics` - The network interface cards (NICs) the VMs are attached to.
    * `account_id` - The account ID of the owner of the NIC.
    * `description` - The description of the NIC.
    * `is_source_dest_checked` - (Net only) If `true`, the source/destination check is enabled. If `false`, it is disabled. This value must be `false` for a NAT VM to perform network address translation (NAT) in a Net.
    * `link_nic` - Information about the network interface card (NIC).
      * `delete_on_vm_deletion` - If `true`, the volume is deleted when the VM is terminated.
      * `device_number` - The device index for the NIC attachment (between 1 and 7, both included).
      * `link_nic_id` - The ID of the NIC to attach.
      * `state` - The state of the attachment (`attaching` \| `attached` \| `detaching` \| `detached`).
    * `link_public_ip` - Information about the EIP associated with the NIC.
      * `public_dns_name` - The name of the public DNS.
      * `public_ip` - The External IP address (EIP) associated with the NIC.
      * `public_ip_account_id` - The account ID of the owner of the EIP.
    * `mac_address` - The Media Access Control (MAC) address of the NIC.
    * `net_id` - The ID of the Net for the NIC.
    * `nic_id` - The ID of the NIC.
    * `private_dns_name` - The name of the private DNS.
    * `private_ips` - The private IP address or addresses of the NIC.
      * `is_primary` - If `true`, the IP address is the primary private IP address of the NIC.
      * `link_public_ip` - Information about the EIP associated with the NIC.
        * `public_dns_name` - The name of the public DNS.
        * `public_ip` - The External IP address (EIP) associated with the NIC.
        * `public_ip_account_id` - The account ID of the owner of the EIP.
      * `private_dns_name` - The name of the private DNS.
      * `private_ip` - The private IP address.
    * `security_groups` - One or more IDs of security groups for the NIC.
      * `security_group_id` - The ID of the security group.
      * `security_group_name` - (Public Cloud only) The name of the security group.
    * `state` - The state of the NIC (`available` \| `attaching` \| `in-use` \| `detaching`).
    * `subnet_id` - The ID of the Subnet for the NIC.
  * `os_family` - Indicates the operating system (OS) of the VM.
  * `performance` - The performance of the VM (`standard` \| `high` \|  `highest`).
  * `placement` - Information about the placement of the VM.
    * `subregion_name` - The name of the Subregion.
    * `tenancy` - The tenancy of the VM (`default` \| `dedicated`).
  * `private_dns_name` - The name of the private DNS.
  * `private_ip` - The primary private IP address of the VM.
  * `product_codes` - The product code associated with the OMI used to create the VM (`0001` Linux/Unix \| `0002` Windows \| `0004` Linux/Oracle \| `0005` Windows 10).
  * `public_dns_name` - The name of the public DNS.
  * `public_ip` - The public IP address of the VM.
  * `reservation_id` - The reservation ID of the VM.
  * `root_device_name` - The name of the root device for the VM (for example, /dev/vda1).
  * `root_device_type` - The type of root device used by the VM (always `bsu`).
  * `security_groups` - One or more security groups associated with the VM.
    * `security_group_id` - The ID of the security group.
    * `security_group_name` - (Public Cloud only) The name of the security group.
  * `state` - The state of the VM (`pending` \| `running` \| `shutting-down` \| `terminated` \| `stopping` \| `stopped`).
  * `state_reason` - The reason explaining the current state of the VM.
  * `subnet_id` - The ID of the Subnet for the VM.
  * `tags` - One or more tags associated with the VM.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
  * `user_data` - The Base64-encoded MIME user data.
  * `vm_id` - The ID of the VM.
  * `vm_initiated_shutdown_behavior` - The VM behavior when you stop it. By default or if set to `stop`, the VM stops. If set to `restart`, the VM stops then automatically restarts. If set to `delete`, the VM stops and is deleted.
  * `vm_type` - The type of VM. For more information, see [Instance Types](https://wiki.outscale.net/display/EN/Instance+Types).
