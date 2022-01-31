---
layout: "outscale"
page_title: "OUTSCALE: outscale_public_ip"
sidebar_current: "outscale-public-ip"
description: |-
  [Manages a public IP.]
---

# outscale_public_ip Resource

Manages a public IP.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-EIPs.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-publicip).

## Example Usage

```hcl
resource "outscale_public_ip" "public_ip01" {
}
```

## Argument Reference

The following arguments are supported:

* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `link_public_ip_id` - (Required in a Net) The ID representing the association of the public IP with the VM or the NIC.
* `nic_account_id` - The account ID of the owner of the NIC.
* `nic_id` - The ID of the NIC the public IP is associated with (if any).
* `private_ip` - The private IP address associated with the public IP.
* `public_ip_id` - The allocation ID of the public IP.
* `public_ip` - The public IP.
* `tags` - One or more tags associated with the public IP.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
* `vm_id` - The ID of the VM the public IP is associated with (if any).

## Import

A public IP can be imported using its ID. For example:

```console

$ terraform import outscale_public_ip.ImportedPublicIp eipalloc-12345678

```