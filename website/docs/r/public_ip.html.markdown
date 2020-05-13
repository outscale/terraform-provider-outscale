---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_public_ip"
sidebar_current: "outscale-public-ip"
description: |-
  [Manages a public IP.]
---

# outscale_public_ip Resource

Manages a public IP.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+EIPs).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-publicip).

## Example Usage

```hcl

resource "outscale_public_ip" "public_ip01" {
}


```

## Argument Reference

The following arguments are supported:

* `tags` - One or more tags to add to this resource.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `public_ip` - Information about the public IP.
  * `link_public_ip_id` - (Required in a Net) The ID representing the association of the EIP with the VM or the NIC.
  * `nic_account_id` - The account ID of the owner of the NIC.
  * `nic_id` - The ID of the NIC the EIP is associated with (if any).
  * `private_ip` - The private IP address associated with the EIP.
  * `public_ip` - The External IP address (EIP) associated with the NAT service.
  * `public_ip_id` - The allocation ID of the EIP associated with the NAT service.
  * `tags` - One or more tags associated with the EIP.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
  * `vm_id` - The ID of the VM the External IP (EIP) is associated with (if any).

## Import

A public IP can be imported using its ID. For example:

```

$ terraform import outscale_public_ip.ImportedPublicIp 111.11.111.11

```