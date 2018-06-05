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
resource "outscale_lin" "outscale_lin" {
    count = 1

    cidr_block = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    count = 1

    availability_zone   = "eu-west-2a"
    cidr_block          = "10.0.0.0/16"
    vpc_id              = "${outscale_lin.outscale_lin.vpc_id}"
}

resource "outscale_nic" "outscale_nic" {
    count = 1

    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}
```

## Argument Reference

The following arguments are supported:

* `subnet_id` - The ID of the subnet
* `description` - (Optional) A description of the network interface.
* `private_ip_adress` - (Optional) The private IP address of the network interface.\
\- This IP address must be in the CIDR of the subnet you specify with the SubnetId attribute. \
\- If you do not specify a private IP address, Outscale selects one in the CIDR of the subnet.
* `security_group_id.N` - One or more security group IDs for the network interface.

## Attributes

* `association` - The association information for an External IP associated with the network interface.
* `attachment` - The network interface attachment.
* `availability_zone` - The Availability Zone in which the network interface is located.
* `description` - A description of the network interface.
* `group_set` - One or more security groups for the network interface.
* `mac_address` - The MAC address.
* `network_interface_id` - The ID of the network interface.
* `owner_id` - The account ID of the owner of the network interface.
* `private_dns_name` - The name of the private DNS.
* `private_ip_address` - The private IP addresses assigned to the network interface, in the CIDR of its subnet.
* `private_ip_addresses_set.N` - Information about one or more private IP addresses assigned to the network interface.
* `requester_id` - The ID of the requester that launched the instances on your behalf.
* `requester_managed` - If true, the network interface is being managed by Outscale.
* `source_dest_check` - If true, the traffic to or from the instance is validated.
* `status` - The state of the network interface (available | attaching | in-use | detaching).
* `subnet_id` - The ID of the subnet.
* `tag_set.N` - One or more tags associated with the network interface.
* `vpc_id` - The ID of the VPC.

[See detailed information](http://docs.outscale.com/api_fcu/operations/Action_CreateNetworkInterface_get.html#_api_fcu-action_createnetworkinterface_get).
