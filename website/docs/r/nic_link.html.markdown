---
layout: "outscale"
page_title: "OUTSCALE: outscale_nic_link"
sidebar_current: "docs-outscale-resource-nic-link"
description: |-
  Attaches a network interface to an instance.
---

# outscale_nic_link

Attaches a network interface to an instance.
The interface and the instance must be in the same Availability Zone (AZ). The instance can be either running or stopped. The network interface must be in the available state.

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

resource "outscale_nic_link" "outscale_nic_link" {
    device_index            = "1"
    instance_id             = "${outscale_vm.outscale_instance.id}"
    network_interface_id    = "${outscale_nic.outscale_nic.id}"
}
```

## Argument Reference

The following arguments are supported:

* `device_index` - The index of the instance device for the network interface (between 1 and 7, both included).
* `instance_id` - The ID of the instance to which you want to attach the network interface.
* `network_interface_id` - (Optional) The ID of the network interface you want to attach.

## Attributes

* `nic_sort_number` - The instance device index of the network interface attachment.
* `vm_id` - The ID of the instance.
* `nic_id` - The ID of the network interface.
* `request_id` - The ID of tue request.

[See detailed information](http://docs.outscale.com/api_fcu/operations/Action_AttachNetworkInterface_get.html#_api_fcu-action_attachnetworkinterface_get).
