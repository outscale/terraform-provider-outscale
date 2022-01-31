---
layout: "outscale"
page_title: "OUTSCALE: outscale_vms"
sidebar_current: "outscale-vms"
description: |-
  [Provides information about virtual machines (VMs).]
---

# outscale_vms Data Source

Provides information about virtual machines (VMs).
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Instances.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-vm).

## Example Usage

```hcl
data "outscale_vms" "vms_01" {
  filter {
    name   = "vm_ids"
    values = ["i-12345678", "i-12345679"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `tag_keys` - (Optional) The keys of the tags associated with the VMs.
    * `tag_values` - (Optional) The values of the tags associated with the VMs.
    * `tags` - (Optional) The key/value combination of the tags associated with the VMs, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.
    * `vm_ids` - (Optional) One or more IDs of VMs.

## Attribute Reference

The following attributes are exported:

* `vms` - Information about one or more VMs.
    * `architecture` - The architecture of the VM (`i386` \| `x86_64`).
    * `block_device_mappings_created` - The block device mapping of the VM.
        * `bsu` - Information about the created BSU volume.
            * `delete_on_vm_deletion` - If true, the volume is deleted when terminating the VM. If false, the volume is not deleted when terminating the VM.
            * `link_date` - The time and date of attachment of the volume to the VM.
            * `state` - The state of the volume.
            * `volume_id` - The ID of the volume.
        * `device_name` - The name of the device.
    * `bsu_optimized` - If true, the VM is optimized for BSU I/O.
    * `client_token` - The idempotency token provided when launching the VM.
    * `deletion_protection` - If true, you cannot terminate the VM using Cockpit, the CLI or the API. If false, you can.
    * `hypervisor` - The hypervisor type of the VMs (`ovm` \| `xen`).
    * `image_id` - The ID of the OMI used to create the VM.
    * `is_source_dest_checked` - (Net only) If true, the source/destination check is enabled. If false, it is disabled. This value must be false for a NAT VM to perform network address translation (NAT) in a Net.
    * `keypair_name` - The name of the keypair used when launching the VM.
    * `launch_number` - The number for the VM when launching a group of several VMs (for example, 0, 1, 2, and so on).
    * `net_id` - The ID of the Net in which the VM is running.
    * `nics` - (Net only) The network interface cards (NICs) the VMs are attached to.
        * `account_id` - The account ID of the owner of the NIC.
        * `description` - The description of the NIC.
        * `is_source_dest_checked` - (Net only) If true, the source/destination check is enabled. If false, it is disabled. This value must be false for a NAT VM to perform network address translation (NAT) in a Net.
        * `link_nic` - Information about the network interface card (NIC).
            * `delete_on_vm_deletion` - If true, the NIC is deleted when the VM is terminated.
            * `device_number` - The device index for the NIC attachment (between 1 and 7, both included).
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
        * `private_ips` - The private IP address or addresses of the NIC.
            * `is_primary` - If true, the IP address is the primary private IP address of the NIC.
            * `link_public_ip` - Information about the public IP associated with the NIC.
                * `public_dns_name` - The name of the public DNS.
                * `public_ip` - The public IP associated with the NIC.
                * `public_ip_account_id` - The account ID of the owner of the public IP.
            * `private_dns_name` - The name of the private DNS.
            * `private_ip` - The private IP address.
        * `security_groups` - One or more IDs of security groups for the NIC.
            * `security_group_id` - The ID of the security group.
            * `security_group_name` - The name of the security group.
        * `state` - The state of the NIC (`available` \| `attaching` \| `in-use` \| `detaching`).
        * `subnet_id` - The ID of the Subnet for the NIC.
    * `os_family` - Indicates the operating system (OS) of the VM.
    * `performance` - The performance of the VM (`medium` \| `high` \|  `highest`).
    * `placement` - Information about the placement of the VM.
        * `subregion_name` - The name of the Subregion.
        * `tenancy` - The tenancy of the VM (`default` \| `dedicated`).
    * `private_dns_name` - The name of the private DNS.
    * `private_ip` - The primary private IP address of the VM.
    * `product_codes` - The product code associated with the OMI used to create the VM (`0001` Linux/Unix \| `0002` Windows \| `0004` Linux/Oracle \| `0005` Windows 10).
    * `public_dns_name` - The name of the public DNS.
    * `public_ip` - The public IP of the VM.
    * `reservation_id` - The reservation ID of the VM.
    * `root_device_name` - The name of the root device for the VM (for example, /dev/vda1).
    * `root_device_type` - The type of root device used by the VM (always `bsu`).
    * `security_groups` - One or more security groups associated with the VM.
        * `security_group_id` - The ID of the security group.
        * `security_group_name` - The name of the security group.
    * `state` - The state of the VM (`pending` \| `running` \| `stopping` \| `stopped` \| `shutting-down` \| `terminated` \| `quarantine`).
    * `state_reason` - The reason explaining the current state of the VM.
    * `subnet_id` - The ID of the Subnet for the VM.
    * `tags` - One or more tags associated with the VM.
        * `key` - The key of the tag, with a minimum of 1 character.
        * `value` - The value of the tag, between 0 and 255 characters.
    * `user_data` - The Base64-encoded MIME user data.
    * `vm_id` - The ID of the VM.
    * `vm_initiated_shutdown_behavior` - The VM behavior when you stop it. By default or if set to `stop`, the VM stops. If set to `restart`, the VM stops then automatically restarts. If set to `terminate`, the VM stops and is deleted.
    * `vm_type` - The type of VM. For more information, see [Instance Types](https://wiki.outscale.net/display/EN/Instance+Types).
