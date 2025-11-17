---
layout: "outscale"
page_title: "OUTSCALE: outscale_vm"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-vm"
description: |-
  [Provides information about a virtual machine (VM).]
---

# outscale_vm Data Source

Provides information about a virtual machine (VM).

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-VMs.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-vm).

## Example Usage

```hcl
data "outscale_vm" "vm01" {
    filter {
        name   = "vm_ids"
        values = ["i-12345678"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `architectures` - (Optional) The architectures of the VMs (`i386` \| `x86_64`).
    * `block_device_mapping_delete_on_vm_deletion` - (Optional) Whether the BSU volumes are deleted when terminating the VMs.
    * `block_device_mapping_device_names` - (Optional) The device names for the BSU volumes (in the format `/dev/sdX`, `/dev/sdXX`, `/dev/xvdX`, or `/dev/xvdXX`).
    * `block_device_mapping_link_dates` - (Optional) The link dates for the BSU volumes mapped to the VMs (for example, `2016-01-23T18:45:30.000Z`).
    * `block_device_mapping_states` - (Optional) The states for the BSU volumes (`attaching` \| `attached` \| `detaching` \| `detached`).
    * `block_device_mapping_volume_ids` - (Optional) The volume IDs of the BSU volumes.
    * `boot_modes` - (Optional) The boot modes of the VMs. Possible values: `uefi` | `legacy`.
    * `client_tokens` - (Optional) The idempotency tokens provided when launching the VMs.
    * `creation_dates` - (Optional) The dates when the VMs were launched.
    * `image_ids` - (Optional) The IDs of the OMIs used to launch the VMs.
    * `is_source_dest_checked` - (Optional) Whether the source/destination checking is enabled (true) or disabled (false).
    * `keypair_names` - (Optional) The names of the keypairs used when launching the VMs.
    * `launch_numbers` - (Optional) The numbers for the VMs when launching a group of several VMs (for example, `0`, `1`, `2`, and so on).
    * `lifecycles` - (Optional) Whether the VMs are Spot Instances (spot).
    * `net_ids` - (Optional) The IDs of the Nets in which the VMs are running.
    * `nic_account_ids` - (Optional) The IDs of the NICs.
    * `nic_descriptions` - (Optional) The descriptions of the NICs.
    * `nic_is_source_dest_checked` - (Optional) Whether the source/destination checking is enabled (true) or disabled (false).
    * `nic_link_nic_delete_on_vm_deletion` - (Optional) Whether the NICs are deleted when the VMs they are attached to are deleted.
    * `nic_link_nic_device_numbers` - (Optional) The device numbers the NICs are attached to.
    * `nic_link_nic_link_nic_dates` - (Optional) The dates and times (UTC) when the NICs were attached to the VMs.
    * `nic_link_nic_link_nic_ids` - (Optional) The IDs of the NIC attachments.
    * `nic_link_nic_states` - (Optional) The states of the attachments.
    * `nic_link_nic_vm_account_ids` - (Optional) The account IDs of the owners of the VMs the NICs are attached to.
    * `nic_link_nic_vm_ids` - (Optional) The IDs of the VMs the NICs are attached to.
    * `nic_link_public_ip_account_ids` - (Optional) The account IDs of the owners of the public IPs associated with the NICs.
    * `nic_link_public_ip_link_public_ip_ids` - (Optional) The association IDs returned when the public IPs were associated with the NICs.
    * `nic_link_public_ip_public_ip_ids` - (Optional) The allocation IDs returned when the public IPs were allocated to their accounts.
    * `nic_link_public_ip_public_ips` - (Optional) The public IPs associated with the NICs.
    * `nic_mac_addresses` - (Optional) The Media Access Control (MAC) addresses of the NICs.
    * `nic_net_ids` - (Optional) The IDs of the Nets where the NICs are located.
    * `nic_nic_ids` - (Optional) The IDs of the NICs.
    * `nic_private_ips_link_public_ip_account_ids` - (Optional) The account IDs of the owner of the public IPs associated with the private IPs.
    * `nic_private_ips_link_public_ip_ids` - (Optional) The public IPs associated with the private IPs.
    * `nic_private_ips_primary_ip` - (Optional) Whether the private IPs are the primary IPs associated with the NICs.
    * `nic_private_ips_private_ips` - (Optional) The private IPs of the NICs.
    * `nic_security_group_ids` - (Optional) The IDs of the security groups associated with the NICs.
    * `nic_security_group_names` - (Optional) The names of the security groups associated with the NICs.
    * `nic_states` - (Optional) The states of the NICs (`available` \| `in-use`).
    * `nic_subnet_ids` - (Optional) The IDs of the Subnets for the NICs.
    * `nic_subregion_names` - (Optional) The Subregions where the NICs are located.
    * `platforms` - (Optional) The platforms. Use windows if you have Windows VMs. Otherwise, leave this filter blank.
    * `private_ips` - (Optional) The private IPs of the VMs.
    * `product_codes` - (Optional) The product codes associated with the OMI used to create the VMs.
    * `public_ips` - (Optional) The public IPs of the VMs.
    * `reservation_ids` - (Optional) The IDs of the reservation of the VMs, created every time you launch VMs. These reservation IDs can be associated with several VMs when you launch a group of VMs using the same launch request.
    * `root_device_names` - (Optional) The names of the root devices for the VMs (for example, `/dev/sda1`)
    * `root_device_types` - (Optional) The root devices types used by the VMs (always `ebs`)
    * `security_group_ids` - (Optional) The IDs of the security groups for the VMs (only in the public Cloud).
    * `security_group_names` - (Optional) The names of the security groups for the VMs (only in the public Cloud).
    * `state_reason_codes` - (Optional) The reason codes for the state changes.
    * `state_reason_messages` - (Optional) The messages describing the state changes.
    * `state_reasons` - (Optional) The reasons explaining the current states of the VMs. This filter is like the `state_reason_codes` one.
    * `subnet_ids` - (Optional) The IDs of the Subnets for the VMs.
    * `subregion_names` - (Optional) The names of the Subregions of the VMs.
    * `tag_keys` - (Optional) The keys of the tags associated with the VMs.
    * `tag_values` - (Optional) The values of the tags associated with the VMs.
    * `tags` - (Optional) The key/value combinations of the tags associated with the VMs, in the following format: `TAGKEY=TAGVALUE`.
    * `tenancies` - (Optional) The tenancies of the VMs (`dedicated` \| `default` \| `host`).
    * `vm_ids` - (Optional) One or more IDs of VMs.
    * `vm_security_group_ids` - (Optional) The IDs of the security groups for the VMs.
    * `vm_security_group_names` - (Optional) The names of the security group for the VMs.
    * `vm_state_codes` - (Optional) The state codes of the VMs: `-1` (quarantine), `0` (pending), `16` (running), `32` (shutting-down), `48` (terminated), `64` (stopping), and `80` (stopped).
    * `vm_state_names` - (Optional) The state names of the VMs (`pending` \| `running` \| `stopping` \| `stopped` \| `shutting-down` \| `terminated` \| `quarantine`).
    * `vm_types` - (Optional) The VM types (for example, t2.micro). For more information, see [VM Types](https://docs.outscale.com/en/userguide/VM-Types.html).

## Attribute Reference

The following attributes are exported:

* `actions_on_next_boot` - The action to perform on the next boot of the VM.
    * `secure_boot` - One action to perform on the next boot of the VM. For more information, see [About Secure Boot](https://docs.outscale.com/en/userguide/About-Secure-Boot.html#_secure_boot_actions).
* `architecture` - The architecture of the VM (`i386` \| `x86_64`).
* `block_device_mappings_created` - The block device mapping of the VM.
    * `bsu` - Information about the created BSU volume.
        * `delete_on_vm_deletion` - If true, the volume is deleted when terminating the VM. If false, the volume is not deleted when terminating the VM.
        * `link_date` - The date and time (UTC) at which the volume was attached to the VM, in ISO 8601 date-time format.
        * `state` - The state of the volume.
        * `volume_id` - The ID of the volume.
    * `device_name` - The name of the device.
* `boot_mode` - The boot mode of the VM. Possible values: `uefi` | `legacy`.
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
* `placement` - Information about the placement of the VM.
    * `subregion_name` - The name of the Subregion. If you specify this parameter, you must not specify the `nics` parameter.
    * `tenancy` - The tenancy of the VM (`default`, `dedicated`, or a dedicated group ID).
* `private_dns_name` - The name of the private DNS.
* `private_ip` - The primary private IP of the VM.
* `product_codes` - The product codes associated with the OMI used to create the VM.
* `public_dns_name` - The name of the public DNS.
* `public_ip` - The public IP of the VM.
* `reservation_id` - The reservation ID of the VM.
* `root_device_name` - The name of the root device for the VM (for example, `/dev/sda1`).
* `root_device_type` - The type of root device used by the VM (always `bsu`).
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
