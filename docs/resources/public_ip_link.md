---
layout: "outscale"
page_title: "OUTSCALE: outscale_public_ip_link"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-public-ip-link"
description: |-
  [Manages a public IP link.]
---

# outscale_public_ip_link Resource

Manages a public IP link.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Public-IPs.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-publicip).

## Example Usage

### Required resources

```hcl
resource "outscale_public_ip" "public_ip01" {
}

resource "outscale_vm" "vm01" {
	image_id           = var.image_id
	vm_type            = var.vm_type
	keypair_name       = var.keypair_name
	security_group_ids = [var.security_group_id]
}
```

### Link a public IP address to a VM

```hcl
resource "outscale_public_ip_link" "public_ip_link01" {
	vm_id     = outscale_vm.vm01.vm_id
	public_ip = outscale_public_ip.public_ip01.public_ip
}
```

## Argument Reference

The following arguments are supported:

* `allow_relink` - (Optional) If true, allows the public IP to be associated with the VM or NIC that you specify even if it is already associated with another VM or NIC. If false, prevents the public IP from being associated with the VM or NIC that you specify if it is already associated with another VM or NIC. (By default, true in the public Cloud, false in a Net.)
* `nic_id` - (Optional) (Net only) The ID of the NIC. This parameter is required if the VM has more than one NIC attached. Otherwise, you need to specify the `vm_id` parameter instead. You cannot specify both parameters at the same time.
* `private_ip` - (Optional) (Net only) The primary or secondary private IP of the specified NIC. By default, the primary private IP.
* `public_ip_id` - (Optional) The allocation ID of the public IP. This parameter is required unless you use the `public_ip` parameter.
* `public_ip` - (Optional) The public IP. This parameter is required unless you use the `public_ip_id` parameter.
* `vm_id` - (Optional) The ID of the VM.<br />- In the public Cloud, this parameter is required.<br />- In a Net, this parameter is required if the VM has only one NIC. Otherwise, you need to specify the `nic_id` parameter instead. You cannot specify both parameters at the same time.

## Attribute Reference

The following attributes are exported:

* `link_public_ip_id` - (Net only) The ID representing the association of the public IP with the VM or the NIC.

## Import

A public IP link can be imported using the public IP or the public IP link ID. For example:

```console

$ terraform import outscale_public_ip_link.ImportedPublicIpLink eipassoc-12345678

```