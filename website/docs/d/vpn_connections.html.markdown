---
layout: "outscale"
page_title: "OUTSCALE: outscale_vpn_connections"
sidebar_current: "docs-outscale-datasource-vpn-connections"
description: |-
  Describes the VPN connections.
---

# outscale_volume

  Describes the VPN connections.

## Example Usage

```hcl
data "outscale_vpn_connections" "outscale_vpn_connections" {
    vpn_connection_id = ["${outscale_vpn_connection.outscale_vpn_connection.id}","${outscale_vpn_connection.outscale_vpn_connection2.id}"]
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

resource "outscale_vpn_gateway_link" "outscale_vpn_gateway_link" {
    vpc_id         = "${outscale_lin.outscale_lin.vpc_id}"
    #vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway.vpn_gateway_id}"
    vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway.id}"
}

resource "outscale_lin" "outscale_lin" {
    cidr_block = "10.0.0.0/16"
}

resource "outscale_vpn_connection" "outscale_vpn_connection2" {
    customer_gateway_id = "${outscale_client_endpoint.outscale_client_endpoint2.id}"
    vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway2.id}"
    type = "ipsec.1" 
}

resource "outscale_vpn_gateway" "outscale_vpn_gateway2" {
    type = "ipsec.1" 
}

resource "outscale_client_endpoint" "outscale_client_endpoint2" {
    bgp_asn     = "3"
    ip_address  = "171.33.74.123"
    type        = "ipsec.1"
}

resource "outscale_vpn_gateway_link" "outscale_vpn_gateway_link2" {
    vpc_id         = "${outscale_lin.outscale_lin2.vpc_id}"
    #vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway2.vpn_gateway_id}"
    vpn_gateway_id = "${outscale_vpn_gateway.outscale_vpn_gateway2.id}"
}

resource "outscale_lin" "outscale_lin2" {
    cidr_block = "10.0.0.0/16"
}
```

## Argument Reference

The following arguments are supported:

* `vpn_connection_id`	- (Optional) One or more VPN connections IDs.

See detailed information in [Outscale VPN Connections](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Filters

Use the Filter.N parameter to filter the described VPN Connections on the following properties:

* `bgp-asn`: -	The Border Gateway Protocol (BGP) Autonomous System Number (ASN) of the connection.	
* `customer-gateway-configuration`: -	The XML configuration of the customer gateway connection.	
* `customer-gateway-id`: -	The ID of the customer gateway.	
* `option.static-routes-only`: -	Whether the connection has static routes only.	
* `route.destination-cidr-block`: -	The destination CIDR block.	
* `state`: -	The state of the connection (pending | available | deleting | deleted).	
* `type`: -	The type of connection (only ipsec.1 is supported).	
* `vpn-gateway-id`: -	The ID of the virtual private gateway.
* `vpn-connection-id`: -	The ID of the VPN connection.	
* `tag`: -	The key/value combination of a tag associated with the resource.
* `tag-key`: -	The key of a tag associated with the resource.
* `tag-value`: -	The value of a tag associated with the resource.


## Attributes Reference

The following attributes are exported:

* `vpn_connection_set` -	Information about one or more VPN connections.	false	VpnConnection
* `request_id`-	The ID of the request	false	string

See detailed information in [Describe VPN Connections](http://docs.outscale.com/api_fcu/operations/Action_DescribeVpnConnections_get.html#_api_fcu-action_describevpnconnections_get).