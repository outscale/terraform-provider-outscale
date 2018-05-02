---
layout: "outscale"
page_title: "OUTSCALE: outscale_vpn_gateway_link"
sidebar_current: "docs-outscale-resource-vpn-gateway-link"
description: |-
  Provides an Outscale resource to Attaches a virtual private gateway to a VPC.
---

# outscale_vpn_gateway_link

Provides an Outscale resource to Attaches a virtual private gateway to a VPC. [provisioning](/docs/provisioners/index.html).

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

resource "outscale_vpn_gateway_link" "outscale_vpn_gateway_link" {
    count = 1

    vpc_id          = ""
    vpn_gateway_id  = "${outscale_vpn_gateway.outscale_vpn_gateway.vpn_gateway_id}"
}

```

## Argument Reference

The following arguments are supported:

* `VpcId` - (Required)	The ID of the VPC.
* `VpnGatewayId` - (Required)	The ID of the virtual private gateway.

See detailed information in [Outscale VPN Gateway Link](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Attributes Reference

The following attributes are exported:

* `state`	- The current state of the attachment (attaching | attached | detaching| detached).
* `vpc_id`	- The ID of the Virtual Private Cloud (VPC).
* `vpn_gateway_id`	- The ID of the virtual private gateway.
* `request_id`	- The ID of the request.

See detailed information in [Describe VPN Gateway Link](http://docs.outscale.com/api_fcu/operations/Action_AttachVpnGateway_get.html#_api_fcu-action_attachvpngateway_get
).
