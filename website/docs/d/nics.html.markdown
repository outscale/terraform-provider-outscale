---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_nic"
sidebar_current: "outscale-nic"
description: |-
  [Provides information about NICs.]
---

# outscale_nic Data Source

Provides information about NICs.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+FNIs).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-nic).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `filters` - One or more filters.
  * `link_nic_sort_numbers` - (Optional) The device numbers the NICs are attached to.
  * `link_nic_vm_ids` - (Optional) The IDs of the VMs the NICs are attached to.
  * `nic_ids` - (Optional) The IDs of the NICs.
  * `private_ips_private_ips` - (Optional) The private IP addresses of the NICs.
  * `subnet_ids` - (Optional) The IDs of the Subnets for the NICs.

## Attribute Reference

The following attributes are exported:

* `nics` - Information about one or more NICs.
  * `account_id` - The account ID of the owner of the NIC.
  * `description` - The description of the NIC.
  * `is_source_dest_checked` - (Net only) If `true`, the source/destination check is enabled. If `false`, it is disabled. This value must be `false` for a NAT VM to perform network address translation (NAT) in a Net.
  * `link_nic` - Information about the NIC attachment.
    * `delete_on_vm_deletion` - If `true`, the volume is deleted when the VM is terminated.
    * `device_number` - The device index for the NIC attachment (between 1 and 7, both included).
    * `link_nic_id` - The ID of the NIC to attach.
    * `state` - The state of the attachment (`attaching` \| `attached` \| `detaching` \| `detached`).
    * `vm_account_id` - The account ID of the owner of the VM.
    * `vm_id` - The ID of the VM.
  * `link_public_ip` - Information about the EIP association.
    * `link_public_ip_id` - (Required in a Net) The ID representing the association of the EIP with the VM or the NIC.
    * `public_dns_name` - The name of the public DNS.
    * `public_ip` - The External IP address (EIP) associated with the NIC.
    * `public_ip_account_id` - The account ID of the owner of the EIP.
    * `public_ip_id` - The allocation ID of the EIP.
  * `mac_address` - The Media Access Control (MAC) address of the NIC.
  * `net_id` - The ID of the Net for the NIC.
  * `nic_id` - The ID of the NIC.
  * `private_dns_name` - The name of the private DNS.
  * `private_ips` - The private IP addresses of the NIC.
    * `is_primary` - If `true`, the IP address is the primary private IP address of the NIC.
    * `link_public_ip` - Information about the EIP association.
      * `link_public_ip_id` - (Required in a Net) The ID representing the association of the EIP with the VM or the NIC.
      * `public_dns_name` - The name of the public DNS.
      * `public_ip` - The External IP address (EIP) associated with the NIC.
      * `public_ip_account_id` - The account ID of the owner of the EIP.
      * `public_ip_id` - The allocation ID of the EIP.
    * `private_dns_name` - The name of the private DNS.
    * `private_ip` - The private IP address of the NIC.
  * `security_groups` - One or more IDs of security groups for the NIC.
    * `security_group_id` - The ID of the security group.
    * `security_group_name` - (Public Cloud only) The name of the security group.
  * `state` - The state of the NIC (`available` \| `attaching` \| `in-use` \| `detaching`).
  * `subnet_id` - The ID of the Subnet.
  * `subregion_name` - The Subregion in which the NIC is located.
  * `tags` - One or more tags associated with the NIC.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
