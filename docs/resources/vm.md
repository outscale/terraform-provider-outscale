---
layout: "outscale"
page_title: "OUTSCALE: outscale_vm"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-vm"
description: |-
  [Manages a virtual machine (VM).]
---

# outscale_vm Resource

Manages a virtual machine (VM).

~> **Important** Consider using the `primary_nic` argument if you plan to use the `outscale_nic_link`resource.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-VMs.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-vm).

## Example Usage

### Optional resource

```hcl
resource "outscale_keypair" "keypair01" {
	keypair_name = "terraform-keypair-for-vm"
}
```

### Create a VM in the public Cloud

```hcl
resource "outscale_security_group" "security_group01" {
  description         = "vm security group"
  security_group_name = "vm_security_group1"
}

resource "outscale_vm" "vm01" {
	image_id                 = var.image_id
	vm_type                  = var.vm_type
	keypair_name             = var.keypair_name
	security_group_ids       = [outscale_security_group.security_group01.security_group_id]
	placement_subregion_name = "eu-west-2a"
	placement_tenancy        = "default"
	tags {
		key   = "name"
		value = "terraform-public-vm"
	}
	user_data                = base64encode(<<EOF
	<CONFIGURATION>
	EOF
	)
}
```

### Create a VM with block device mappings

```hcl
resource "outscale_security_group" "security_group01" {
  description         = "vm security group"
  security_group_name = "vm_security_group1"
}

resource "outscale_vm" "vm02" {
	image_id                = var.image_id
	vm_type                 = var.vm_type
	keypair_name            = var.keypair_name
	security_group_ids  	= [outscale_security_group.security_group01.security_group_id]
	block_device_mappings {
		device_name = "/dev/sdb" # /dev/sdb
		bsu {
			volume_size = 15
			volume_type = "gp2"
			snapshot_id = var.snapshot_id
		}
	}
	block_device_mappings {
		device_name = "/dev/sdc"
		bsu {
			volume_size           = 22
			volume_type           = "io1"
			iops                  = 150
			delete_on_vm_deletion = true
		}
	}
}


resource "outscale_security_group" "security_group01" {
  description         = "vm security group"
  security_group_name = "vm_security_group1"
}

resource "outscale_vm" "vm02" {
	image_id 				= var.image_id
 	vm_type 				= var.vm_type
 	keypair_name 			= var.keypair_name
	security_group_ids 		= [outscale_security_group.security_group01.security_group_id]
 	block_device_mappings {
		device_name = "/dev/sdb"
		bsu {
			volume_size           = 30
			volume_type           = "gp2"
			snapshot_id           = outscale_snapshot.snapshot.id
			delete_on_vm_deletion = false
			tags {
				key                   = "Name"
				value                 = "bsu-tags-gp2"
			}
		}
	}
}
```

### Create a VM in a Net with a network

```hcl
resource "outscale_net" "net01" {
	ip_range = "10.0.0.0/16"
	tags {
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

resource "outscale_security_group" "security_group01" {
	description          = "Terraform security group for VM"
	security_group_name = "terraform-security-group-for-vm"
	net_id               = outscale_net.net01.net_id
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
```

### Create a VM with a primary NIC

~> **Note:** If you plan to use the `outscale_nic_link`resource, it is recommended to specify the `primary_nic` argument to define the primary network interface of a VM.

```hcl
resource "outscale_security_group" "security_group01" {
  description         = "vm security group"
  security_group_name = "vm_security_group1"
}

resource "outscale_net" "net02" {
	ip_range = "10.0.0.0/16"
	tags {
		key   = "name"
		value = "terraform-net-for-vm-with-nic"
	}
}

resource "outscale_subnet" "subnet02" {
	net_id         = outscale_net.net02.net_id
	ip_range       = "10.0.0.0/24"
	subregion_name = "eu-west-2a"
	tags {
		key   = "name"
		value = "terraform-subnet-for-vm-with-nic"
	}
}
resource "outscale_nic" "nic01" {
	subnet_id = outscale_subnet.subnet02.subnet_id
}

resource "outscale_vm" "vm04" {
	image_id     		= var.image_id
	vm_type      		= "tinav5.c1r1p2"
	keypair_name 		= var.keypair_name
	security_group_ids  = [outscale_security_group.security_group01.security_group_id]
	primary_nic {
		nic_id        = outscale_nic.nic01.nic_id
		device_number = "0"
	}
}
```

### Create a VM with secondary NICs

```hcl
resource "outscale_security_group" "security_group01" {
  description         = "vm security group"
  security_group_name = "vm_security_group1"
}

resource "outscale_net" "net02" {
    ip_range = "10.0.0.0/16"
    tags {
        key   = "name"
        value = "terraform-net-for-vm-with-nic"
    }
}

resource "outscale_subnet" "subnet02" {
    net_id         = outscale_net.net02.net_id
    ip_range       = "10.0.0.0/24"
    subregion_name = "eu-west-2a"
    tags {
        key   = "name"
        value = "terraform-subnet-for-vm-with-nic"
    }
}
resource "outscale_nic" "nic01" {
    subnet_id = outscale_subnet.subnet02.subnet_id
}

resource "outscale_vm" "vm04" {
    image_id     		= var.image_id
    vm_type      		= "tinav5.c1r1p2"
    keypair_name 		= var.keypair_name
	security_group_ids  = [outscale_security_group.security_group01.security_group_id]
    nics {
        nic_id        = outscale_nic.nic01.nic_id
        device_number = "0"
	}
	nics {	
		nic_id        = outscale_nic.nic02.nic_id
        device_number = "1"
    }
}
```

### Create a VM with Secure Boot

~> **Important** Secure Boot is only available with VMs booting in Unified Extensible Firmware Interface (UEFI).

```hcl
resource "outscale_security_group" "security_group01" {
  description         = "vm security group"
  security_group_name = "vm_security_group1"
}

resource "outscale_vm" "outscale_vm_TF206" {
  image_id            = var.image_id
  vm_type             = "tinav5.c3r3"
  keypair_name 		= var.keypair_name
  security_group_ids       = [outscale_security_group.security_group01.security_group_id]
  deletion_protection = false
  state               = "stopped"
  boot_mode           = "uefi"
  secure_boot_action  = "enable"
}
```

### Create a VM with an ephemeral keypair

~> **Important** Ephemeral keypairs are only available in Terraform versions 1.10 and up.

```hcl
ephemeral "outscale_keypair" "ephemeral_keypair" {
  keypair_name = "ephemeral-keypair"
  }

resource "outscale_security_group" "security_group01" {
  description         = "vm security group"
  security_group_name = "vm_security_group12"
}

resource "outscale_vm" "outscale_vm_TF206" {
  image_id            = var.image_id
  vm_type             = var.vm_type
  keypair_name_wo     = ephemeral.outscale_keypair.ephemeral_keypair.keypair_name
  security_group_ids  = [outscale_security_group.security_group01.security_group_id] 
}
```

## Argument Reference

The following arguments are supported:

* `block_device_mappings` - (Optional) One or more block device mappings.
    * `bsu` - Information about the BSU volume to create.
      * `delete_on_vm_deletion` - (Optional) By default or if set to true, the volume is deleted when terminating the VM. If false, the volume is not deleted when terminating the VM.
      * `iops` - (Optional) The number of I/O operations per second (IOPS). This parameter must be specified only if you create an `io1` volume. The maximum number of IOPS allowed for `io1` volumes is `13000` with a maximum performance ratio of 300 IOPS per gibibyte.
      * `snapshot_id` - (Optional) The ID of the snapshot used to create the volume.
      * `volume_size` - (Optional) The size of the volume, in gibibytes (GiB).
      * `tags`- One or more tags associated with the VM.
          * `key`- The key of the tag with a minimum of 1 character.
          * `value` - The value of the tag, between 0 and 255 characters.

* `boot_mode` - (Optional) The boot mode of the VM. Possible values: `uefi` | `legacy`.
* `client_token` - (Optional) A unique identifier which enables you to manage the idempotency.
* `deletion_protection` - (Optional) If true, you cannot delete the VM unless you change this parameter back to false.
* `get_admin_password` - (Optional) (Windows VM only) If true, waits for the administrator password of the VM to become available in order to retrieve the VM. The password is exported to the `admin_password` attribute.
* `image_id` - (Required) The ID of the OMI used to create the VM. You can find the list of OMIs by calling the [ReadImages](https://docs.outscale.com/api#readimages) method.
* `keypair_name_wo` - The name of the keypair. This write-only parameter is required to use the ephemeral keypair resource.
* `keypair_name` - (Optional) The name of the keypair.
* `nested_virtualization` - (Optional) (dedicated tenancy only) If true, nested virtualization is enabled. If false, it is disabled.
* `nics` - (Optional) One or more NICs. If you specify this parameter, you must not specify the `subnet_id` and `subregion_name` parameters. To define a NIC as the primary network interface of the VM, use the `primary_nic` argument.
    * `delete_on_vm_deletion` - (Optional) If true, the NIC is deleted when the VM is terminated. You can specify this parameter only for a new NIC. To modify this value for an existing NIC, see [UpdateNic](https://docs.outscale.com/api#updatenic).
    * `description` - (Optional) The description of the NIC, if you are creating a NIC when creating the VM.
    * `device_number` - (Optional) The index of the VM device for the NIC attachment (between `1` and `7`, both included). This parameter is required if you create a NIC when creating the VM.
    * `nic_id` - (Optional) The ID of the NIC, if you are attaching an existing NIC when creating a VM.
    * `private_ips` - (Optional) One or more private IPs to assign to the NIC, if you create a NIC when creating a VM. Only one private IP can be the primary private IP.
        * `is_primary` - (Optional) If true, the IP is the primary private IP of the NIC.
        * `private_ip` - (Optional) The private IP of the NIC.
    * `secondary_private_ip_count` - (Optional) The number of secondary private IPs, if you create a NIC when creating a VM. This parameter cannot be specified if you specified more than one private IP in the `private_ips` parameter.
    * `security_group_ids` - (Optional) One or more IDs of security groups for the NIC, if you create a NIC when creating a VM.
    * `subnet_id` - (Optional) The ID of the Subnet for the NIC, if you create a NIC when creating a VM. This parameter is required if you create a NIC when creating the VM.

* `performance` - (Optional) The performance of the VM (`medium` | `high` | `highest`). Updating this parameter will trigger a stop/start of the VM.
* `placement_subregion_name` - (Optional) The name of the Subregion where the VM is placed.
* `placement_tenancy` - (Optional) The tenancy of the VM (`default` | `dedicated`).
* `primary_nic` - (Optional) The primary network interface of the VM.
    * `delete_on_vm_deletion` - (Optional) If true, the NIC is deleted when the VM is terminated. You can specify this parameter only for a new NIC. To modify this value for an existing NIC, see [UpdateNic](https://docs.outscale.com/api#updatenic).
    * `description` - (Optional) The description of the NIC, if you are creating a NIC when creating the VM.
    * `device_number` - (Optional) The index of the VM device for the NIC attachment (must be `0`). This parameter is required if you create a NIC when creating the VM.
    * `nic_id` - (Optional) The ID of the NIC, if you are attaching an existing NIC when creating a VM.
    * `private_ips` - (Optional) One or more private IPs to assign to the NIC, if you create a NIC when creating a VM. Only one private IP can be the primary private IP.
        * `is_primary` - (Optional) If true, the IP is the primary private IP of the NIC.
        * `private_ip` - (Optional) The private IP of the NIC.
    * `secondary_private_ip_count` - (Optional) The number of secondary private IPs, if you create a NIC when creating a VM. This parameter cannot be specified if you specified more than one private IP in the `private_ips` parameter.
    * `security_group_ids` - (Optional) One or more IDs of security groups for the NIC, if you create a NIC when creating a VM.
    * `subnet_id` - (Optional) The ID of the Subnet for the NIC, if you create a NIC when creating a VM. This parameter is required if you create a NIC when creating the VM.

* `private_ips` - (Optional) One or more private IPs of the VM.
* `secure_boot_action` - (Optional) One action to perform on the next boot of the VM (`enable` | `disable` | `setup-mode` |`none`).<br /> For more information, see [About Secure Boot](https://docs.outscale.com/en/userguide/About-Secure-Boot.html#_secure_boot_actions).
* `security_group_ids` - (Optional) One or more IDs of security group for the VMs. You must specify at least one of the following parameters: `security_group_ids` or `security_group_names`.
* `security_group_names` - (Optional) One or more names of security groups for the VMs. You must specify at least one of the following parameters: `security_group_ids` or `security_group_names`.
* `state` - The state of the VM (`running` | `stopped`). If set to `stopped`, the VM is stopped regardless of the value of the `vm_initiated_shutdown_behavior` argument.
* `subnet_id` - (Optional) The ID of the Subnet in which you want to create the VM. If you specify this parameter, you must not specify the `nics` parameter.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.
* `user_data` - (Optional) Data or script used to add a specific configuration to the VM. It must be Base64-encoded, either directly or using the [base64encode](https://www.terraform.io/docs/configuration/functions/base64encode.html) Terraform function. For multiline strings, use [heredoc syntax](https://www.terraform.io/docs/configuration/expressions.html#string-literals). Updating this parameter will trigger a stop/start of the VM.
* `vm_initiated_shutdown_behavior` - (Optional) The VM behavior when you stop it. By default or if set to `stop`, the VM stops. If set to `restart`, the VM stops then automatically restarts. If set to `terminate`, the VM stops and is terminated.
* `vm_type` - (Optional) The type of VM (`t2.small` by default). Updating this parameter will trigger a stop/start of the VM.<br /> For more information, see [VM Types](https://docs.outscale.com/en/userguide/VM-Types.html).

## Attribute Reference

The following attributes are exported:

* `actions_on_next_boot` - The action to perform on the next boot of the VM.
    * `secure_boot` - One action to perform on the next boot of the VM. For more information, see [About Secure Boot](https://docs.outscale.com/en/userguide/About-Secure-Boot.html#_secure_boot_actions).
* `admin_password` - (Windows VM only) The administrator password of the VM. This password is encrypted with the keypair you specified when launching the VM and encoded in Base64. You need to wait about 10 minutes after launching the VM to be able to retrieve this password.<br />If `get_admin_password` is false or not specified, the VM resource is created without the `admin_password` attribute. Once `admin_password` is available, it will appear in the Terraform state after the next **refresh** or **apply** command.<br />If `get_admin_password` is true, the VM resource itself is not considered created until the `admin_password` attribute is available.<br />Note also that after the first reboot of the VM, this attribute can no longer be retrieved. For more information on how to use this password to connect to the VM, see [Accessing a Windows VM](https://docs.outscale.com/en/userguide/Accessing-a-Windows-VM.html).
* `architecture` - The architecture of the VM (`i386` \| `x86_64`).
* `block_device_mappings_created` - The block device mapping of the VM.
    * `bsu` - Information about the created BSU volume.
        * `delete_on_vm_deletion` - If true, the volume is deleted when terminating the VM. If false, the volume is not deleted when terminating the VM.
        * `link_date` - The date and time (UTC) at which the volume was attached to the VM, in ISO 8601 date-time format.
        * `state` - The state of the volume.
        * `volume_id` - The ID of the volume.
    * `device_name` - The name of the device.
* `boot_mode` - The boot mode of the VM (`legacy` | `uefi`).
* `client_token` - The idempotency token provided when launching the VM.
* `creation_date` - The date and time (UTC) at which the VM was created.
* `deletion_protection` - If true, you cannot delete the VM unless you change this parameter back to false.
* `hypervisor` - The hypervisor type of the VMs (`ovm` \| `xen`).
* `image_id` - The ID of the OMI used to create the VM.
* `is_source_dest_checked` - (Net only) If true, the source/destination check is enabled. If false, it is disabled.
* `keypair_name` - The name of the keypair used when launching the VM.
* `launch_number` - The number for the VM when launching a group of several VMs (for example, `0`, `1`, `2`, and so on).
* `nested_virtualization` - If true, nested virtualization is enabled. If false, it is disabled.
* `net_id` - The ID of the Net in which the VM is running.
* `nics` - (Net only) The network interface cards (NICs) the VMs are attached to.
    * `account_id` - The account ID of the owner of the NIC.
    * `description` - The description of the NIC.
    * `is_source_dest_checked` - (Net only) If true, the source/destination check is enabled. If false, it is disabled.
    * `link_nic` - Information about the network interface card (NIC).
        * `delete_on_vm_deletion` - If true, the NIC is deleted when the VM is terminated.
        * `device_number` - The device index for the NIC attachment (between `1` and `7`, both included).
        * `link_nic_id` - The ID of the NIC to attach.
        * `state` - The state of the attachment (`attaching` \| `attached` \| `detaching` \| `detached`).
    * `link_public_ip` - Information about the public IP associated with the NIC.
        * `public_dns_name` - The name of the public DNS.
        * `public_ip` - The public IP associated with the NIC.
        * `public_ip_account_id` - The account ID of the owner of the public IP.
    * `mac_address` - The Media Access Control (MAC) address of the NIC.
    * `net_id` - The ID of the Net for the NIC.
    * `nic_id` - The ID of the NIC.
    * `private_dns_name` - The name of the private DNS.
    * `private_ips` - The private IP or IPs of the NIC.
        * `is_primary` - If true, the IP is the primary private IP of the NIC.
        * `link_public_ip` - Information about the public IP associated with the NIC.
            * `public_dns_name` - The name of the public DNS.
            * `public_ip` - The public IP associated with the NIC.
            * `public_ip_account_id` - The account ID of the owner of the public IP.
        * `private_dns_name` - The name of the private DNS.
        * `private_ip` - The private IP.
    * `security_groups` - One or more IDs of security groups for the NIC.
        * `security_group_id` - The ID of the security group.
        * `security_group_name` - The name of the security group.
    * `state` - The state of the NIC (`available` \| `attaching` \| `in-use` \| `detaching`).
    * `subnet_id` - The ID of the Subnet for the NIC.
* `os_family` - Indicates the operating system (OS) of the VM.
* `performance` - The performance of the VM (`medium` \| `high` \|  `highest`).
* `placement_subregion_name` - The name of the Subregion where the VM is placed.
* `placement_tenancy` - The tenancy of the VM (`default` | `dedicated`).
* `private_dns_name` - The name of the private DNS.
* `private_ip` - The primary private IP of the VM.
* `product_codes` - The product codes associated with the OMI used to create the VM.
* `public_dns_name` - The name of the public DNS.
* `public_ip` - The public IP of the VM.
* `reservation_id` - The reservation ID of the VM.
* `root_device_name` - The name of the root device for the VM (for example, `/dev/sda1`).
* `root_device_type` - The type of root device used by the VM (always `bsu`).
* `secure_boot_action` - The action to perform on the next boot of the VM (`enable` | `disable` | `setup-mode` |`none`).
* `security_groups` - One or more security groups associated with the VM.
    * `security_group_id` - The ID of the security group.
    * `security_group_name` - The name of the security group.
* `state_reason` - The reason explaining the current state of the VM.
* `state` - The state of the VM (`pending` \| `running` \| `stopping` \| `stopped` \| `shutting-down` \| `terminated` \| `quarantine`).
* `subnet_id` - The ID of the Subnet for the VM.
* `tags` - One or more tags associated with the VM.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
* `user_data` - The Base64-encoded MIME user data.
* `vm_id` - The ID of the VM.
* `vm_initiated_shutdown_behavior` - The VM behavior when you stop it. If set to `stop`, the VM stops. If set to `restart`, the VM stops then automatically restarts. If set to `terminate`, the VM stops and is deleted.
* `vm_type` - The type of VM. For more information, see [VM Types](https://docs.outscale.com/en/userguide/VM-Types.html).

## Import

A VM can be imported using its ID. For example:

```console

$ terraform import outscale_vm.ImportedVm i-12345678

```