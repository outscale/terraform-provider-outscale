---
layout: "outscale"
page_title: "OUTSCALE: outscale_nic_private_ip"
sidebar_current: "docs-outscale-resource-nic-private-ip"
description: |-
  Assigns one or more secondary private IP addresses to a specified network interface.
---

# outscale_nic_private_ip

Assigns one or more secondary private IP addresses to a specified network interface.
This action is only available in a VPC.
The private IP addresses to be assigned can be added individually using the PrivateIpAddress.N parameter, or you can specify the number of private IP addresses to be automatically chosen within the subnet range using the SecondaryPrivateIpAddressCount parameter. You can specify only one of these two parameters. If none of these parameters are specified, a private IP address is chosen within the subnet range.

## Example Usage

```hcl
resource "outscale_vm" "outscale_instance" {
    image_id                    = "ami-880caa66"
    instance_type               = "c4.large"
    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}

resource "outscale_lin" "outscale_lin" {
    cidr_block          = "10.0.0.0/16"
}

resource "outscale_subnet" "outscale_subnet" {
    availability_zone   = "eu-west-2a"
    cidr_block          = "10.0.0.0/16"
    vpc_id              = "${outscale_lin.outscale_lin.id}"
}

resource "outscale_nic" "outscale_nic" {
    subnet_id = "${outscale_subnet.outscale_subnet.subnet_id}"
}

resource "outscale_nic_private_ip" "outscale_nic_private_ip" {
    network_interface_id    = "${outscale_nic.outscale_nic.id}"
}
```

## Argument Reference

The following arguments are supported:

* `network_interface_id` - The ID of the network interface you want to attach.
* `allow_reassignment` - (Optional) If set to true, allows an IP address that is already assigned   to another network interface in the same subnet to be assigned to the network interface you     specified.
* `secondary_private_ip_address_count` - (Optional) The number of secondary private IP addresses to be automatically assigned to the network interface, chosen within the subnet range.
* `private_ip_address.N` - The secondary private IP address you want to assign to the network interface.

## Attributes

* `network_interface_id` - The ID of the network interface you want to attach.
* `allow_reassignment` - (Optional) If set to true, allows an IP address that is already assigned   to another network interface in the same subnet to be assigned to the network interface you     specified.
* `secondary_private_ip_address_count` - (Optional) The number of secondary private IP addresses to be automatically assigned to the network interface, chosen within the subnet range.
* `private_ip_address.N` - The secondary private IP address you want to assign to the network interface.

[See detailed information.](http://docs.outscale.com/api_fcu/operations/Action_AssignPrivateIpAddresses_get.html#_api_fcu-action_assignprivateipaddresses_get)
