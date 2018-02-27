---
layout: "outscale"
page_title: "OUTSCALE: outscale_public_ip_link"
sidebar_current: "docs-outscale-resource-public_ip_link"
description: |-
  Provides an Outscale Public IP Association as a top level resource, to associate and disassociate Public IPs from Outscale VMs and Network Interfaces.
---

# public_ip_link

NOTE: outscale_public_ip_link is useful in scenarios where Public IPs are either pre-existing or distributed to customers or users and therefore cannot be changed.

## Example Usage

```hcl

resource "outscale_public_ip_link" "oip_assoc" {
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

* `allocation_id` - (Optional) The allocation ID. This is required for Outscale VM-VPC.
* `allow_reassociation` - (Optional, Boolean) Whether to allow an Public IP to be re-associated. Defaults to true in VPC.
* `instance_id` - (Optional) The ID of the instance. This is required for Outscale VM-Classic. For Outscale VM-VPC, you can specify either the instance ID or the network interface ID, but not both. The operation fails if you specify an instance ID unless exactly one network interface is attached.
* `network_interface_id` - (Optional) The ID of the network interface. If the instance has more than one network interface, you must specify a network interface ID.
* `private_ip_address` - (Optional) The primary or secondary private IP address to associate with the Public IP address. If no private IP address is specified, the Public IP address is associated with the primary private IP address.
* `public_ip` - (Optional) The Public IP address. This is required for Outscale VM-Classic.

## Attributes Reference

* `association_id` - The ID that represents the association of the Public IP address with an instance.
* `allocation_id` - As above
* `instance_id` - As above
* `network_interface_id` - As above
* `private_ip_address` - As above
* `public_ip` - As above