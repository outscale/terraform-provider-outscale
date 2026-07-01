---
layout: "outscale"
page_title: "OUTSCALE: outscale_virtual_gateway_route_propagation"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-virtual-gateway-route-propagation"
description: |-
  [Manages a virtual gateway route propagation.]
---

# outscale_virtual_gateway_route_propagation Resource

Manages a virtual gateway route propagation.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Routing-Configuration-for-VPN-Connections.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#updateroutepropagation).

## Example Usage

### Required resources

```hcl
resource "outscale_virtual_gateway" "virtual_gateway01" {
	connection_type = "ipsec.1"
}

resource "outscale_net" "net01" {
	ip_range = "10.0.0.0/16"
}

resource "outscale_route_table" "route_table01" {
	net_id = outscale_net.net01.net_id
}

resource "outscale_virtual_gateway_link" "virtual_gateway_link01" {
	virtual_gateway_id = outscale_virtual_gateway.virtual_gateway01.virtual_gateway_id
	net_id             = outscale_net.net01.net_id
}
```

### Activate the propagation of routes to a route table of a Net by a virtual gateway

```hcl
resource "outscale_virtual_gateway_route_propagation" "virtual_gateway_route_propagation01" {
	enable             = true
	virtual_gateway_id = outscale_virtual_gateway.virtual_gateway01.virtual_gateway_id
	route_table_id     = outscale_route_table.route_table01.route_table_id
	depends_on         = [outscale_virtual_gateway_link.virtual_gateway_link01]
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required) If true, a virtual gateway can propagate routes to a specified route table of a Net. If false, the propagation is disabled.
* `route_table_id` - (Required) The ID of the route table.
* `virtual_gateway_id` - (Required) The ID of the virtual gateway.

## Attribute Reference

The following attributes are exported:

* `enable` - If true, a virtual gateway can propagate routes to a specified route table of a Net. If false, the propagation is disabled.
* `route_table_id` - The ID of the route table.
* `virtual_gateway_id` - The ID of the virtual gateway.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 10 minutes.
* `read` - Defaults to 5 minutes.
* `delete` - Defaults to 5 minutes.
