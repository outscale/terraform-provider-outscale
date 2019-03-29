---
layout: "outscale"
page_title: "OUTSCALE: outscale_vpn_gateway_route_propagation"
sidebar_current: "docs-outscale-resource-vpn-gateway_route_propagation"
description: |-
  Provides an Outscale resource to enable a virtual private gateway to propagate routes to a specified route table of a VPC.
---

# outscale_vpn_gateway_route_propagation

Provides an Outscale resource to enable a virtual private gateway to propagate routes to a specified route table of a VPC. [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
resource "outscale_vpn_gateway" "outscale_vpn_gateway" {
    count = 1

    type = "ipsec.1" 
}

resource "outscale_lin" "outscale_lin" {
    count = 1

    cidr_block = "10.0.0.0/16"
}

resource "outscale_route_table" "outscale_route_table" {
    count = 1

    vpc_id = "${outscale_lin.outscale_lin.vpc_id}"
}

resource "outscale_vpn_gateway_route_propagation" "outscale_vpn_gateway_route_propagation" {
    gateway_id      = "${outscale_vpn_gateway.outscale_vpn_gateway.vpn_gateway_id}"
    route_table_id  = "${outscale_route_table.outscale_route_table.route_table_id}"
}
```

## Argument Reference

The following arguments are supported:

* `gateway_id` - (Required)	The ID of the virtual private gateway.
* `route_table_id` - (Required)	The ID of the route table.

See detailed information in [Outscale VPN Gateway Route Propagation](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Attributes Reference

The following attributes are exported:

* `gateway_id`	The ID of the virtual private gateway.
* `route_table_id`	The ID of the route table.
* `request_id`	The ID of the request.

See detailed information in [Describe VPN Gateway Route Propagation](http://docs.outscale.com/api_fcu/operations/Action_EnableVgwRoutePropagation_get.html#_api_fcu-action_enablevgwroutepropagation_get).
