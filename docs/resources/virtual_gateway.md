---
layout: "outscale"
page_title: "OUTSCALE: outscale_virtual_gateway"
sidebar_current: "outscale-virtual-gateway"
description: |-
  [Manages a virtual gateway.]
---

# outscale_virtual_gateway Resource

Manages a virtual gateway.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Virtual-Private-Gateways.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-virtualgateway).

## Example Usage

```hcl
resource "outscale_virtual_gateway" "virtual_gateway01" {
	connection_type = "ipsec.1"
	tags {
		key   = "name"
		value = "terraform-virtual-gateway"
	}
}
```

## Argument Reference

The following arguments are supported:

* `connection_type` - (Required) The type of VPN connection supported by the virtual gateway (only `ipsec.1` is supported).
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `connection_type` - The type of VPN connection supported by the virtual gateway (only `ipsec.1` is supported).
* `net_to_virtual_gateway_links` - The Net to which the virtual gateway is attached.
    * `net_id` - The ID of the Net to which the virtual gateway is attached.
    * `state` - The state of the attachment (`attaching` \| `attached` \| `detaching` \| `detached`).
* `state` - The state of the virtual gateway (`pending` \| `available` \| `deleting` \| `deleted`).
* `tags` - One or more tags associated with the virtual gateway.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
* `virtual_gateway_id` - The ID of the virtual gateway.

## Import

A virtual gateway can be imported using its ID. For example:

```console

$ terraform import outscale_virtual_gateway.ImportedVirtualGateway vgw-12345678

```
