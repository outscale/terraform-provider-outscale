---
layout: "outscale"
page_title: "OUTSCALE: outscale_nic"
sidebar_current: "docs-outscale-resource-nic"
description: |-
  Creates a network interface in the specified subnet.
---

# outscale_nic

Creates a network interface in the specified subnet.

## Example Usage

```hcl
resource "outscale_net" "outscale_net" {
    ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    subregion_name   = "eu-west-2a"
    ip_range          = "10.0.0.0/16"
    net_id              = "${outscale_net.outscale_net.net_id}"
}

resource "outscale_nic" "outscale_nic" {
    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}
```

## Argument Reference

The following arguments are supported:

* `subnet_id` - The ID of the subnet
* `description` - (Optional) A description of the network interface.
* `private_ip` - (Optional) The private IP address of the network interface.
    This IP address must be in the CIDR of the subnet you specify with the SubnetId attribute.
    If you do not specify a private IP address, Outscale selects one in the CIDR of the subnet.
* `security_group_id` - One or more security group IDs for the network interface.

## Attributes

* `link_public_ip` - The association information for an External IP associated with the network interface
  * `public_ip_id` ­- The ID of the allocation.
  * `link_public_ip_id` ­- The ID of the association.
  * `public_ip_account_id`­ - The ID of owner of the External IP address.
  * `public_dns_name`­ - The name of the public DNS.
  * `public_ip`­ - The External IP address (EIP) associated with the network interface.
* `link_nic` - The network interface attachment.
  * `link_nic_id` ­- The ID of the network interface attachment.
  * `delete_on_vm_deletion` ­- If true, the volume is deleted when the instance is terminated.
  * `device_number` ­- The instance device index of the network interface attachment.
  * `vm_id`­- The ID of the instance.
  * `vm_account_id` ­- The account ID of the owner of the instance.
  * `state`­- The attachment state (attaching | attached | detaching | detached).
* `subregion_name` -The Availability Zone in which the network interface is located.
* `description` - A description of the network interface.
* `security_groups` - One or more security groups for the network interface.
  * `security_group_id` - The ID of the security group.
  * `security_group_name` - The name of the security group.
* `mac_address` - The MAC address.
* `nic_id` - The ID of the network interface.
* `account_id` - The account ID of the owner of the network interface.
* `private_dns_name` - The name of the private DNS.
* `private_ips` - Information about one or more private IP addresses assigned to the network interface.
  * `link_public_ip` - The association information for an External IP address associated with the network interface.
    * `public_ip_id` - The ID of the allocation.
    * `link_public_ip_id` - The ID of the association.
    * `public_ip_account_id` - The ID of owner of the External IP address.
    * `public_dns_name` - The name of the public DNS.
    * `public_ip` - The External IP address (EIP) associated with the network interface. 
  * `is_primary` - If true, the IP address is the primary private IP address of the network interface.
  * `private_dns_name` - The name of the private DNS.
  * `private_ip` - The private IP address.
* `requester_managed` - If true, the network interface is being managed by Outscale.
* `is_source_dest_checked` - If true, the traffic to or from the instance is validated.
* `state` - The state of the network interface (available | attaching | in-use | detaching).
* `tags` - One or more tags associated with the network interface.
  * `key` - The key of the tag.
  * `value` - The value of the tag.
* `net_id` - The ID of the VPC.
* `request_id` - The ID of the request.


Se detailed information: [CreateNic](https://docs-beta.outscale.com/#createnic)