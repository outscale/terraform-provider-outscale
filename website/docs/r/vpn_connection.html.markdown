---
layout: "outscale"
page_title: "OUTSCALE: outscale_vpn_connection"
sidebar_current: "docs-outscale-vpn-connection"
description: |-
  Creates a VPN connection between a specified virtual private gateway and a specified customer gateway.
---

# outscale_vpn_connection

Creates a VPN connection between a specified virtual private gateway and a specified customer gateway. Instances also support [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
data "outscale_volume" "outscale_volume" {
  most_recent = true

  filter {
    name   = "volume-type"
    values = ["gp2"]
  }

  filter {
    name   = "tag:Name"
    values = ["Example"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `CustomerGatewayId` - (Required)	The ID of the customer gateway.
* `Options` - (Optional)	Options for the VPN connection.	
* `Type` - (Required)	The type of VPN connection (only ipsec.1 is supported).	
* `VpnGatewayId` - (Required)	The ID of the virtual private gateway.

See detailed information in [Outscale VPN Connection](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Attributes Reference

The following attributes are exported:

* `customer_gateway_configuration` - The configuration to apply to the customer gateway to establish the VPN connection, in XML format.
* `customer_gateway_id` - The ID of the customer gateway used on the customer end of the connection.
* `options` - One or more options for the VPN connection.
* `routes` - Information about one or more static routes associated with the VPN connection, if any.
* `state` - The state of the VPN connection.
* `tag_set` -  One or more tags associated with the VPN connection.
* `type` - The type of VPN connection (always ipsec.1).
* `vgw_telemetry` - Information about the state of the VPN tunnel (this list contains only one element, as Outscale supports one tunnel per VPN connection).
* `vpn_connection_id` - The ID of the VPN connection.
* `vpn_gateway_id` - The ID of the virtual private gateway used on the Outscale end of the connection.
* `request_id` - The ID of the request

See detailed information in [Describe VPN Connection](http://docs.outscale.com/api_fcu/definitions/VpnConnection.html#_api_fcu-vpnconnection).
