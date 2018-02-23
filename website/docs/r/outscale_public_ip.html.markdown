---
layout: "outscale"
page_title: "OUTSCALE: outscale_public_ip"
sidebar_current: "docs-outscale-resource-public_ip"
description: |-
  Provides an Outscale Public IP Association as a top level resource, to associate and disassociate Public IPs from Outscale VMs and Network Interfaces.
---

# public_ip

NOTE: outscale_public_ip is useful in scenarios where Public IPs are either pre-existing or distributed to customers or users and therefore cannot be changed.

## Example Usage

```hcl

resource "outscale_public_ip" "oip_assoc" {
  instance_id   = "${outscale_vm.web.id}"
  allocation_id = "${outscale_public_ip.example.id}"
}

resource "outscale_vm" "web" {
 image_id = "ami-8a6a0120"
instance_type = "t2.micro"

  tags {
    Name = "HelloWorld"
  }
}

resource "outscale_public_ip" "example" {
  vpc = true
}
```

## Argument Reference

The following arguments are supported:

* `domain` - (Optional) The type of platform in which you want to use the EIP (standard | vpc).

## Attributes Reference

* `allocation_id` - (Optional) The ID that represents the allocation of the EIP for use with instances in a VPC.
* `association_id` - (Optional) The association ID for the EIP.
* `domain` - (Optional) The type of platform in which you can use the EIP.
* `instance_id` - (Optional) The ID of the instance. This is required for Outscale VM-Classic. For Outscale VM-VPC, you can specify either the instance ID or the network interface ID, but not both. The operation fails if you specify an instance ID unless exactly one network interface is attached.
* `network_interface_id` - (Optional) The ID of the network interface. If the instance has more than one network interface, you must specify a network interface ID.
* `network_interface_owner_id` - (Optional) The account ID of the owner.
* `private_ip_address` - (Optional) The primary or secondary private IP address to associate with the Public IP address. If no private IP address is specified, the Public IP address is associated with the primary private IP address.
* `public_ip` - (Optional) The External IP address.