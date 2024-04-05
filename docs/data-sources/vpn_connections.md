---
layout: "outscale"
page_title: "OUTSCALE: outscale_vpn_connections"
sidebar_current: "outscale-vpn-connections"
description: |-
  [Provides information about VPN connections.]
---

# outscale_vpn_connections Data Source

Provides information about VPN connections.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-VPN-Connections.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-vpnconnection).

## Example Usage

```hcl
data "outscale_vpn_connections" "vpn_connections01" {
	filter {
		name   = "client_gateway_ids"
		values = ["cgw-12345678"]
	}
	filter {
		name   = "virtual_gateway_ids"
		values = ["vgw-12345678", "vgw-12345678"]
	}
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `bgp_asns` - (Optional) The Border Gateway Protocol (BGP) Autonomous System Numbers (ASNs) of the connections.
    * `client_gateway_ids` - (Optional) The IDs of the client gateways.
    * `connection_types` - (Optional) The types of the VPN connections (only `ipsec.1` is supported).
    * `route_destination_ip_ranges` - (Optional) The destination IP ranges.
    * `states` - (Optional) The states of the VPN connections (`pending` \| `available` \| `deleting` \| `deleted`).
    * `static_routes_only` - (Optional) If false, the VPN connection uses dynamic routing with Border Gateway Protocol (BGP). If true, routing is controlled using static routes. For more information about how to create and delete static routes, see [CreateVpnConnectionRoute](https://docs.outscale.com/api#createvpnconnectionroute) and [DeleteVpnConnectionRoute](https://docs.outscale.com/api#deletevpnconnectionroute).
    * `tag_keys` - (Optional) The keys of the tags associated with the VPN connections.
    * `tag_values` - (Optional) The values of the tags associated with the VPN connections.
    * `tags` - (Optional) The key/value combinations of the tags associated with the VPN connections, in the following format: `TAGKEY=TAGVALUE`.
    * `virtual_gateway_ids` - (Optional) The IDs of the virtual gateways.
    * `vpn_connection_ids` - (Optional) The IDs of the VPN connections.
* `next_page_token` - (Optional) The token to request the next page of results. Each token refers to a specific page.
* `results_per_page` - (Optional) The maximum number of logs returned in a single response (between `1`and `1000`, both included). By default, `100`.

## Attribute Reference

The following attributes are exported:

* `next_page_token` - The token to request the next page of results. Each token refers to a specific page.
* `vpn_connections` - Information about one or more VPN connections.
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
