---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_public_ip_link"
sidebar_current: "outscale-public-ip-link"
description: |-
  [Manages a public IP link.]
---

# outscale_public_ip_link Resource

Manages a public IP link.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+EIPs#AboutEIPs-EipAssocationEIPAssociation).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#linkpublicip).

## Example Usage

```hcl

#resource "outscale_public_ip" "public_ip01" {
#}

#resource "outscale_vm" "vm01" {
#	image_id           = var.image_id
#	vm_type            = var.vm_type
#	keypair_name       = var.keypair_name
#	security_group_ids = [var.security_group_id]
#}

resource "outscale_public_ip_link" "public_ip_link01" {
	vm_id     = outscale_vm.vm01.vm_id
	public_ip = outscale_public_ip.public_ip01.public_ip
}


```

## Argument Reference

The following arguments are supported:

* `allow_relink` - (Optional) - If true, allows the EIP to be associated with the VM or NIC that you specify even if it is already associated with another VM or NIC.<br />
- If false, prevents the EIP from being associated with the VM or NIC that you specify if it is already associated with another VM or NIC.<br /><br />

(By default, true in the public Cloud, false in a Net.)
* `nic_id` - (Optional) (Net only) The ID of the NIC. This parameter is required if the VM has more than one NIC attached. Otherwise, you need to specify the `vm_id` parameter instead. You cannot specify both parameters at the same time.
* `private_ip` - (Optional) (Net only) The primary or secondary private IP address of the specified NIC. By default, the primary private IP address.
* `public_ip_id` - (Optional) The allocation ID of the EIP. In a Net, this parameter is required.
* `public_ip` - (Optional) The EIP. In the public Cloud, this parameter is required.
* `vm_id` - (Optional) The ID of the VM.<br />
- In the public Cloud, this parameter is required.<br />
- In a Net, this parameter is required if the VM has only one NIC. Otherwise, you need to specify the `nic_id` parameter instead. You cannot specify both parameters at the same time.

## Attribute Reference

The following attributes are exported:

* `link_public_ip_id` - (Net only) The ID representing the association of the EIP with the VM or the NIC.

## Import

A public IP link can be imported using the public IP or the public IP link ID. For example:

```

$ terraform import outscale_public_ip_link.ImportedPublicIpLink eipassoc-12345678

```