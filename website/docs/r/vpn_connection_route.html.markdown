---
layout: "outscale"
page_title: "OUTSCALE: outscale_vpn_connection_route"
sidebar_current: "outscale-vpn-connection-route"
description: |-
  [Manages a VPN connection route.]
---

# outscale_vpn_connection_route Resource

Manages a VPN connection route.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Routing-Configuration-for-VPN-Connections.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-vpnconnection).

## Example Usage

### Required resources

```hcl
resource "outscale_client_gateway" "client_gateway01" {
	bgp_asn         = 65000
	public_ip       = "111.11.11.111"
	connection_type = "ipsec.1"
}

resource "outscale_virtual_gateway" "virtual_gateway01" {
	connection_type = "ipsec.1"
}

resource "outscale_vpn_connection" "vpn_connection01" {
	client_gateway_id  = outscale_client_gateway.client_gateway01.client_gateway_id
	virtual_gateway_id = outscale_virtual_gateway.virtual_gateway01.virtual_gateway_id
	connection_type    = "ipsec.1"
	static_routes_only = true
}
```

### Create a static route to a VPN connection

```hcl
resource "outscale_vpn_connection_route" "vpn_connection_route01" {
	vpn_connection_id    = outscale_vpn_connection.vpn_connection01.vpn_connection_id
	destination_ip_range = "10.0.0.0/16"
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

```console

$ terraform import outscale_vpn_connection_route.ImportedRoute vpn-12345678_10.0.0.0/0

```