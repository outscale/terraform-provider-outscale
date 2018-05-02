---
layout: "outscale"
page_title: "OUTSCALE: outscale_vpn_connection-route"
sidebar_current: "docs-outscale-vpn-connection-route"
description: |-
  Creates a static route to a VPN connection.
---

# outscale_vpn_connection_route

  Creates a static route to a VPN connection. Instances also support [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
	resource "outscale_vpn_gateway" "vpn_gateway" {
		tag {
			Name = "vpn_gateway"
		}
	}

	resource "outscale_client_endpoint" "customer_gateway" {
		bgp_asn = %d
		ip_address = "182.0.0.1"
		type = "ipsec.1"
	}

	resource "outscale_vpn_connection" "vpn_connection" {
		vpn_gateway_id = "${outscale_vpn_gateway.vpn_gateway.id}"
		customer_gateway_id = "${outscale_client_endpoint.customer_gateway.id}"
		type = "ipsec.1"
		options {
					static_routes_only = true
		}
	}

	resource "outscale_vpn_connection_route" "foo" {
	    destination_cidr_block = "172.168.10.0/24"
	    vpn_connection_id = "${outscale_vpn_connection.vpn_connection.id}"
	}
```

## Argument Reference

The following arguments are supported:

* `DestinationCidrBlock` - (Required)	The network prefix of the route, in CIDR notation (for example, 10.12.0.0/16).
* `VpnConnectionId` - (Required)	The ID of target VPN connection of the static route.

See detailed information in [Outscale VPN Connection](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Attributes Reference

The following attributes are exported:

* `destination_cidr_block`	The network prefix of the route, in CIDR notation (for example, 10.12.0.0/16)
* `Vpn_connection_id`	The ID of target VPN connection of the static route
* `request_id`	The ID of the request	

See detailed information in [Describe VPN Connection Route](http://docs.outscale.com/api_fcu/operations/Action_CreateVpnConnectionRoute_get.html#_api_fcu-action_createvpnconnectionroute_get).
