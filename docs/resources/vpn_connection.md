---
layout: "outscale"
page_title: "OUTSCALE: outscale_vpn_connection"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-vpn-connection"
description: |-
  [Manages a VPN connection.]
---

# outscale_vpn_connection Resource

Manages a VPN connection.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-VPN-Connections.html).  
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
```

### Create a VPN connection

```hcl
resource "outscale_vpn_connection" "vpn_connection01" {
	client_gateway_id  = outscale_client_gateway.client_gateway01.client_gateway_id
	virtual_gateway_id = outscale_virtual_gateway.virtual_gateway01.virtual_gateway_id
	connection_type    = "ipsec.1"
	static_routes_only = true
	tags {
		key   = "Name"
		value = "vpn01"
	}
}
```

## Argument Reference

The following arguments are supported:

* `client_gateway_id` - (Required) The ID of the client gateway.
* `connection_type` - (Required) The type of VPN connection (always `ipsec.1`).
* `static_routes_only` - (Optional) By default or if false, the VPN connection uses dynamic routing with Border Gateway Protocol (BGP). If true, routing is controlled using static routes. For more information about how to create and delete static routes, see [CreateVpnConnectionRoute](https://docs.outscale.com/api#createvpnconnectionroute) and [DeleteVpnConnectionRoute](https://docs.outscale.com/api#deletevpnconnectionroute).
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.
* `virtual_gateway_id` - (Required) The ID of the virtual gateway.

## Attribute Reference

The following attributes are exported:

* `client_gateway_configuration` - Example configuration for the client gateway.
* `client_gateway_id` - The ID of the client gateway used on the client end of the connection.
* `connection_type` - The type of VPN connection (always `ipsec.1`).
* `routes` - Information about one or more static routes associated with the VPN connection, if any.
    * `destination_ip_range` - The IP range used for the destination match, in CIDR notation (for example, `10.0.0.0/24`).
    * `route_type` - The type of route (always `static`).
    * `state` - The current state of the static route (`pending` \| `available` \| `deleting` \| `deleted`).
* `state` - The state of the VPN connection (`pending` \| `available` \| `deleting` \| `deleted`).
* `static_routes_only` - If false, the VPN connection uses dynamic routing with Border Gateway Protocol (BGP). If true, routing is controlled using static routes. For more information about how to create and delete static routes, see [CreateVpnConnectionRoute](https://docs.outscale.com/api#createvpnconnectionroute) and [DeleteVpnConnectionRoute](https://docs.outscale.com/api#deletevpnconnectionroute).
* `tags` - One or more tags associated with the VPN connection.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
* `vgw_telemetries` - Information about the current state of one or more of the VPN tunnels.
    * `accepted_route_count` - The number of routes accepted through BGP (Border Gateway Protocol) route exchanges.
    * `last_state_change_date` - The date and time (UTC) of the latest state update.
    * `outside_ip_address` - The IP on the OUTSCALE side of the tunnel.
    * `state` - The state of the IPSEC tunnel (`UP` \| `DOWN`).
    * `state_description` - A description of the current state of the tunnel.
* `virtual_gateway_id` - The ID of the virtual gateway used on the OUTSCALE end of the connection.
* `vpn_connection_id` - The ID of the VPN connection.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 10 minutes.
* `read` - Defaults to 5 minutes.
* `update` - Defaults to 10 minutes.
* `delete` - Defaults to 5 minutes.

## Import

A VPN connection can be imported using its ID. For example:

```console

$ terraform import outscale_vpn_connection.ImportedVPN vpn-12345678

```