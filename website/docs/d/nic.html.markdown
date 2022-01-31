---
layout: "outscale"
page_title: "OUTSCALE: outscale_nic"
sidebar_current: "outscale-nic"
description: |-
  [Provides information about a specific network interface card (NIC).]
---

# outscale_nic Data Source

Provides information about a specific network interface card (NIC).
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-FNIs.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-nic).

## Example Usage

```hcl
data "outscale_nic" "nic01" {
  filter {
    name   = "nic_ids"
    values = ["eni-12345678"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `descriptions` - (Optional) The descriptions of the NICs.
    * `is_source_dest_check` - (Optional) Whether the source/destination checking is enabled (true) or disabled (false).
    * `link_nic_delete_on_vm_deletion` - (Optional) Whether the NICs are deleted when the VMs they are attached to are terminated.
    * `link_nic_device_numbers` - (Optional) The device numbers the NICs are attached to.
    * `link_nic_link_nic_ids` - (Optional) The attachment IDs of the NICs.
    * `link_nic_states` - (Optional) The states of the attachments.
    * `link_nic_vm_account_ids` - (Optional) The account IDs of the owners of the VMs the NICs are attached to.
    * `link_nic_vm_ids` - (Optional) The IDs of the VMs the NICs are attached to.
    * `link_public_ip_account_ids` - (Optional) The account IDs of the owners of the public IPs associated with the NICs.
    * `link_public_ip_link_public_ip_ids` - (Optional) The association IDs returned when the public IPs were associated with the NICs.
    * `link_public_ip_public_ip_ids` - (Optional) The allocation IDs returned when the public IPs were allocated to their accounts.
    * `link_public_ip_public_ips` - (Optional) The public IPs associated with the NICs.
    * `mac_addresses` - (Optional) The Media Access Control (MAC) addresses of the NICs.
    * `net_ids` - (Optional) The IDs of the Nets where the NICs are located.
    * `nic_ids` - (Optional) The IDs of the NICs.
    * `private_dns_names` - (Optional) The private DNS names associated with the primary private IP addresses.
    * `private_ips_link_public_ip_account_ids` - (Optional) The account IDs of the owner of the public IPs associated with the private IP addresses.
    * `private_ips_link_public_ip_public_ips` - (Optional) The public IPs associated with the private IP addresses.
    * `private_ips_primary_ip` - (Optional) Whether the private IP address is the primary IP address associated with the NIC.
    * `private_ips_private_ips` - (Optional) The private IP addresses of the NICs.
    * `security_group_ids` - (Optional) The IDs of the security groups associated with the NICs.
    * `security_group_names` - (Optional) The names of the security groups associated with the NICs.
    * `states` - (Optional) The states of the NICs.
    * `subnet_ids` - (Optional) The IDs of the Subnets for the NICs.
    * `subregion_names` - (Optional) The Subregions where the NICs are located.
    * `tag_keys` - (Optional) The keys of the tags associated with the NICs.
    * `tag_values` - (Optional) The values of the tags associated with the NICs.
    * `tags` - (Optional) The key/value combination of the tags associated with the NICs, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.

## Attribute Reference

The following attributes are exported:

* `account_id` - The account ID of the owner of the NIC.
* `description` - The description of the NIC.
* `is_source_dest_checked` - (Net only) If true, the source/destination check is enabled. If false, it is disabled. This value must be false for a NAT VM to perform network address translation (NAT) in a Net.
* `link_nic` - Information about the NIC attachment.
    * `delete_on_vm_deletion` - If true, the NIC is deleted when the VM is terminated.
    * `device_number` - The device index for the NIC attachment (between 1 and 7, both included).
    * `link_nic_id` - The ID of the NIC to attach.
    * `state` - The state of the attachment (`attaching` \| `attached` \| `detaching` \| `detached`).
    * `vm_account_id` - The account ID of the owner of the VM.
    * `vm_id` - The ID of the VM.
* `link_public_ip` - Information about the public IP association.
    * `link_public_ip_id` - (Required in a Net) The ID representing the association of the public IP with the VM or the NIC.
    * `public_dns_name` - The name of the public DNS.
    * `public_ip` - The public IP associated with the NIC.
    * `public_ip_account_id` - The account ID of the owner of the public IP.
    * `public_ip_id` - The allocation ID of the public IP.
* `mac_address` - The Media Access Control (MAC) address of the NIC.
* `net_id` - The ID of the Net for the NIC.
* `nic_id` - The ID of the NIC.
* `private_dns_name` - The name of the private DNS.
* `private_ips` - The private IP addresses of the NIC.
    * `is_primary` - If true, the IP address is the primary private IP address of the NIC.
    * `link_public_ip` - Information about the public IP association.
        * `link_public_ip_id` - (Required in a Net) The ID representing the association of the public IP with the VM or the NIC.
        * `public_dns_name` - The name of the public DNS.
        * `public_ip` - The public IP associated with the NIC.
        * `public_ip_account_id` - The account ID of the owner of the public IP.
        * `public_ip_id` - The allocation ID of the public IP.
    * `private_dns_name` - The name of the private DNS.
    * `private_ip` - The private IP address of the NIC.
* `security_groups` - One or more IDs of security groups for the NIC.
    * `security_group_id` - The ID of the security group.
    * `security_group_name` - The name of the security group.
* `state` - The state of the NIC (`available` \| `attaching` \| `in-use` \| `detaching`).
* `subnet_id` - The ID of the Subnet.
* `subregion_name` - The Subregion in which the NIC is located.
* `tags` - One or more tags associated with the NIC.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
