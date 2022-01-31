---
layout: "outscale"
page_title: "OUTSCALE: outscale_virtual_gateway_link"
sidebar_current: "outscale-virtual-gateway-link"
description: |-
  [Manages a virtual gateway link.]
---

# outscale_virtual_gateway_link Resource

Manages a virtual gateway link.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Virtual-Private-Gateways.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-virtualgateway).

## Example Usage

### Required resources

```hcl
resource "outscale_virtual_gateway" "virtual_gateway01" {
	connection_type = "ipsec.1"
}

resource "outscale_net" "net01" {
	ip_range = "10.0.0.0/16"
}
```

### Link a virtual gateway to a Net

```hcl
resource "outscale_virtual_gateway_link" "virtual_gateway_link01" {
	virtual_gateway_id = outscale_virtual_gateway.virtual_gateway01.virtual_gateway_id
	net_id             = outscale_net.net01.net_id
}
```

## Argument Reference

The following arguments are supported:

* `net_id` - (Required) The ID of the Net to which you want to attach the virtual gateway.
* `virtual_gateway_id` - (Required) The ID of the virtual gateway.

## Attribute Reference

The following attributes are exported:

* `net_id` - The ID of the Net to which the virtual gateway is attached.
* `state` - The state of the attachment (`attaching` \| `attached` \| `detaching` \| `detached`).

## Import

A virtual gateway link can be imported using its virtual gateway ID. For example:

```console

$ terraform import outscale_virtual_gateway_link.ImportedVirtualGatewayLink vgw-12345678

```
