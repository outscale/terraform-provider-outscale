---
layout: "outscale"
page_title: "OUTSCALE: outscale_vpn_connection"
sidebar_current: "docs-outscale-datasource-vpn-connection"
description: |-
  Describes the VPN connections
---

# outscale_volume

Describes the VPN connection.

## Example Usage

```hcl
data "outscale_vpn_connection" "outscale_vpn_connection" {
    vpn_connection_id = "${outscale_vpn_connection.outscale_vpn_connection.id}"
}


resource "outscale_vpn_connection" "outscale_vpn_connection" {
    customer_gateway_id = "${outscale_client_endpoint.outscale_client_endpoint.id}"
    vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway.id}"
    type = "ipsec.1"
}

resource "outscale_vpn_gateway" "outscale_vpn_gateway" {
    type = "ipsec.1"
}

resource "outscale_client_endpoint" "outscale_client_endpoint" {
    bgp_asn     = "3"
    ip_address  = "171.33.74.122"
    type        = "ipsec.1"
}

resource "outscale_lin" "outscale_lin" {
    cidr_block = "10.0.0.0/16"
}

resource "outscale_vpn_gateway_link" "outscale_vpn_gateway_link" {
    vpc_id         = "${outscale_lin.outscale_lin.vpc_id}"
    #vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway.vpn_gateway_id}"
    vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway.id}"
}
```

## Argument Reference

The following arguments are supported:

* `vpn_connection_id` - (Optional) One or more VPN connections IDs.

See detailed information in [Outscale VPN Connections](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Filters

Use the Filter.N parameter to filter the described VPN Connection on the following properties:

* `bgp-asn`: - The Border Gateway Protocol (BGP) Autonomous System Number (ASN) of the connection.
* `customer-gateway-configuration`: - The XML configuration of the customer gateway connection.
* `customer-gateway-id`: - The ID of the customer gateway.
* `option.static-routes-only`: - Whether the connection has static routes only.
* `route.destination-cidr-block`: - The destination CIDR block.
* `state`: - The state of the connection (pending | available | deleting | deleted).
* `type`: - The type of connection (only ipsec.1 is supported).
* `vpn-gateway-id`: - The ID of the virtual private gateway.
* `vpn-connection-id`: - The ID of the VPN connection.
* `tag`: - The key/value combination of a tag associated with the resource.
* `tag-key`: - The key of a tag associated with the resource.
* `tag-value`: - The value of a tag associated with the resource.

## Attributes Reference

The following attributes are exported:

* `customerGatewayConfiguration` - The configuration to apply to the customer gateway to establish the VPN connection, in XML format.
* `customerGatewayId` - The ID of the customer gateway used on the customer end of the connection.
* `options` - One or more options for the VPN connection.
* `routes` - Information about one or more static routes associated with the VPN connection, if any.
* `state` - The state of the VPN connection.
* `tagSet` - One or more tags associated with the VPN connection.
* `type` - The type of VPN connection (always ipsec.1).
* `vgwTelemetry` - Information about the state of the VPN tunnel (this list contains only one element, as Outscale supports one tunnel per VPN connection).
* `vpnConnectionId` - The ID of the VPN connection.
* `vpnGatewayId` - The ID of the virtual private gateway used on the Outscale end of the connection.
* `requestId` - The ID of the request.

See detailed information in [Describe VPN Connection](http://docs.outscale.com/api_fcu/operations/Action_DescribeVpnConnections_get.html#_api_fcu-action_describevpnconnections_get).
