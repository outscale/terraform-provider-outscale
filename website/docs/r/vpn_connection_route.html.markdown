---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_vpn_connection_route"
sidebar_current: "outscale-vpn-connection-route"
description: |-
  [Manages a VPN connection route.]
---

# outscale_vpn_connection_route Resource

Manages a VPN connection route.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Routing+Configuration+for+VPN+Connections).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-vpnconnection).

## Example Usage

```hcl
#resource "outscale_vpn_connection" "vpn_connection01" {
#	client_gateway_id  = "cgw-12345678"
#	virtual_gateway_id = "vgw-12345678"
#	connection_type    = "ipsec.1"
#	static_routes_only = false
#}

resource "outscale_vpn_connection_route" "vpn_connection_route01" {
	vpn_connection_id    = outscale_vpn_connection.vpn_connection01.vpn_connection_id
	destination_ip_range = "10.0.0.0/0"
}
```

## Argument Reference

The following arguments are supported:

* `destination_ip_range` - (Required) The network prefix of the route, in CIDR notation (for example, 10.12.0.0/16).
* `vpn_connection_id` - (Required) The ID of the target VPN connection of the static route.

## Attribute Reference

No attribute is exported.

## Import

A VPN connection route can be imported using the VPN connection ID and the route destination IP range. For example:

```

$ terraform import outscale_vpn_connection_route.ImportedRoute vpn-12345678_10.0.0.0/0

```