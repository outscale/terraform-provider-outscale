---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_vpn_connection"
sidebar_current: "docs-outscale-resource-vpn-connection"
description: |-
  [Manages a VPN connection.]
---

# outscale_vpn_connection Resource

Manages a VPN connection.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+VPN+Connections).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-vpnconnection).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `client_gateway_id` - (Required) The ID of the client gateway.
* `connection_type` - (Required) The type of VPN connection (only `ipsec.1` is supported).
* `static_routes_only` - (Optional) If `false`, the VPN connection uses dynamic routing with Border Gateway Protocol (BGP). If `true`, routing is controlled using static routes. For more information about how to create and delete static routes, see [CreateVpnConnectionRoute](https://docs-beta.outscale.com/#createvpnconnectionroute) and [DeleteVpnConnectionRoute](https://docs-beta.outscale.com/#deletevpnconnectionroute).
* `virtual_gateway_id` - (Required) The ID of the virtual gateway.

## Attribute Reference

The following attributes are exported:

* `vpn_connection` - Information about a VPN connection.
  * `client_gateway_configuration` - The configuration to apply to the client gateway to establish the VPN connection, in XML format.
  * `client_gateway_id` - The ID of the client gateway used on the client end of the connection.
  * `connection_type` - The type of VPN connection (always `ipsec.1`).
  * `routes` - Information about one or more static routes associated with the VPN connection, if any.
    * `destination_ip_range` - The IP range used for the destination match, in CIDR notation (for example, 10.0.0.0/24).
    * `route_type` - The type of route (always `static`).
    * `state` - The current state of the static route (`pending` \| `available` \| `deleting` \| `deleted`).
  * `state` - The state of the VPN connection (`pending` \| `available` \| `deleting` \| `deleted`).
  * `static_routes_only` - If `false`, the VPN connection uses dynamic routing with Border Gateway Protocol (BGP). If `true`, routing is controlled using static routes. For more information about how to create and delete static routes, see [CreateVpnConnectionRoute](https://docs-beta.outscale.com/#createvpnconnectionroute) and [DeleteVpnConnectionRoute](https://docs-beta.outscale.com/#deletevpnconnectionroute).
  * `tags` - One or more tags associated with the VPN connection.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
  * `virtual_gateway_id` - The ID of the virtual gateway used on the 3DS OUTSCALE end of the connection.
  * `vpn_connection_id` - The ID of the VPN connection.
