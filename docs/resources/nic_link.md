---
layout: "outscale"
page_title: "OUTSCALE: outscale_nic_link"
sidebar_current: "outscale-nic-link"
description: |-
  [Manages a NIC link.]
---

# outscale_nic_link Resource

Manages a NIC link.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-FNIs.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-nic).

## Example Usage

### Required resources

```hcl
resource "outscale_net" "net01" {
	ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet01" {
	subregion_name = "${var.region}a"
	ip_range       = "10.0.0.0/16"
	net_id         = outscale_net.net01.net_id
}

resource "outscale_vm" "vm01" {
	image_id     = var.image_id
	vm_type      = var.vm_type
	keypair_name = var.keypair_name
	subnet_id    = outscale_subnet.subnet01.subnet_id
}

resource "outscale_nic" "nic01" {
	subnet_id = outscale_subnet.subnet01.subnet_id
}
```

### Link a NIC to a VM

```hcl
resource "outscale_nic_link" "nic_link01" {
	device_number = "1"
	vm_id         = outscale_vm.vm01.vm_id
	nic_id        = outscale_nic.nic01.nic_id
}
```

## Argument Reference

The following arguments are supported:

* `device_number` - (Required) The index of the VM device for the NIC attachment (between 1 and 7, both included).
* `nic_id` - (Required) The ID of the NIC you want to attach.
* `vm_id` - (Required) The ID of the VM to which you want to attach the NIC.

## Attribute Reference

The following attributes are exported:

* `link_nic_id` - The ID of the NIC attachment.

## Import

A NIC link can be imported using the NIC ID. For example:

```console

$ terraform import outscale_nic_link.ImportedNicLink eni-12345678

```