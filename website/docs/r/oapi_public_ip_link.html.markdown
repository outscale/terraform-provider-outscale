---
layout: "outscale"
page_title: "OUTSCALE: outscale_public_ip_link"
sidebar_current: "docs-outscale-resource-public_ip_link"
description: |-
  Associates an External IP address (EIP) with an instance or a network interface.
---

# outscale_public_ip_link

An EIP address is a static IP address designed for dynamic Cloud computing. It can be used for instances in the public Cloud (standard) or in a Virtual Private Cloud (VPC).
If you want to associate a new EIP to an instance that is already associated with another EIP, this action disassociates the old EIP and associates the new one. If you do not specify any network interface to associate the EIP with, it is associated with the primary network interface.

NOTE: You can associate an EIP with a NAT gateway only when creating the NAT gateway. To modify its EIP, you need to delete the NAT gateway and create a new one. For more information, see the CreateNatGateway method.

## Example Usage

```hcl
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
	subnet_id = "subnet-861fbecc"
}

resource "outscale_public_ip" "bar" {}

resource "outscale_public_ip_link" "by_public_ip" {
	public_ip = "${outscale_public_ip.bar.public_ip}"
	vm_id = "${outscale_vm.basic.id}"
  depends_on = ["outscale_vm.basic", "outscale_public_ip.bar"]
}
```

## Argument Reference

The following arguments are supported:

* `reservation_id` - (Optional) The allocation ID, required for instances in a VPC.
* `allow_relink` - (VPC only) If set to true, allows an EIP that is already associated with an instance or a network interface to be reassociated with the instance or network interface you specify.
* `vm_id` - (Optional) The ID of the instance or of the network interface.
* `nic_id` - (Optional) The ID of the network interface, required if the instance has more than one network interface.
* `private_ip` - (Optional) The primary or secondary private IP address to associate with the External IP address.
* `public_ip` - (Optional) The External IP address.

## Attributes Reference

* `reservation_id` - The ID of the address allocation for use with a VPC. The ID of the allocation.
* `link_id` - The ID of the address association with an instance in a VPC. The ID of the association.
* `vm_id` - The ID of the instance the address is associated with. The ID of the instance.
* `nic_id` - The ID of the instance or of the network interface.
* `private_ip` - The primary or secondary private IP address to associate with the External IP address.
* `public_ip` - The External IP address.
* `request_id` - The ID of the request.
