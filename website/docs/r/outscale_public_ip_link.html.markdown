---
layout: "outscale"
page_title: "OUTSCALE: outscale_public_ip_link"
sidebar_current: "docs-outscale-resource-public_ip_link"
description: |-
  Associates an External IP address (EIP) with an instance or a network interface.
An EIP address is a static IP address designed for dynamic Cloud computing. It can be used for instances in the public Cloud (standard) or in a Virtual Private Cloud (VPC).
If you want to associate a new EIP to an instance that is already associated with another EIP, this action disassociates the old EIP and associates the new one. If you do not specify any network interface to associate the EIP with, it is associated with the primary network interface.
---

NOTE: You can associate an EIP with a NAT gateway only when creating the NAT gateway. To modify its EIP, you need to delete the NAT gateway and create a new one. For more information, see the CreateNatGateway method.

## Example Usage

```hcl

resource "outscale_public_ip_link" "oip_assoc" {
  instance_id   = "${outscale_vm.web.id}"
  allocation_id = "${outscale_public_ip.example.id}"
}

resource "outscale_vm" "web" {
 image_id = "ami-8a6a0120"
 instance_type = "t2.micro"

}

resource "outscale_public_ip" "example" {}
```

## Argument Reference

The following arguments are supported:

* `allocation_id` - (Optional) The allocation ID, required for instances in a VPC.
* `allow_reassociation` - (VPC only) If set to true, allows an EIP that is already associated with an instance or a network interface to be reassociated with the instance or network interface you specify.
* `instance_id` - (Optional) The ID of the instance or of the network interface.
* `network_interface_id` - (Optional) The ID of the network interface, required if the instance has more than one network interface.
* `private_ip_address` - (Optional) The primary or secondary private IP address to associate with the External IP address.
* `public_ip` - (Optional) The External IP address.

## Attributes Reference

* `association_id` - The ID that represents the association of the EIP with an instance or a network interface.
* `allocation_id` - As above
* `instance_id` - As above
* `network_interface_id` - As above
* `private_ip_address` - As above
* `public_ip` - As above

See detailed information in [Describe Addresses](http://docs.outscale.com/api_fcu/operations/Action_DescribeAddresses_get.html#_api_fcu-action_describeaddresses_get).