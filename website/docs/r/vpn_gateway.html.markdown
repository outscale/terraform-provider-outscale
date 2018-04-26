---
layout: "outscale"
page_title: "OUTSCALE: outscale_vpn_gateway"
sidebar_current: "docs-outscale-resource-vpn-gateway"
description: |-
  Provides an Outscale Volume resource to Create a virtual private gateway. A virtual private gateway is the endpoint on the VPC side of a VPN connection.
---

# outscale_vpn_gateway

  Provides an Outscale Volume resource to Create a virtual private gateway. A virtual private gateway is the endpoint on the VPC side of a VPN connection. [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
resource "outscale_lin" "test" {
	cidr_block = "10.0.0.0/16"
}

resource "outscale_vpn_gateway" "test" { 
	type = "ipsec.1" 
}

resource "outscale_vpn_gateway_link" "test" {
	vpc_id = "${outscale_lin.test.id}"
	vpn_gateway_id = "${outscale_vpn_gateway.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `Type` - (Required)	The type of VPN connection supported by the virtual private gateway (only ipsec.1 is supported).

See detailed information in [Outscale Instances](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Attributes Reference

The following attributes are exported:

* `attachment_set` - Information about your volume attachment.

* `Attachments`	The VPC to which the virtual private gateway is attached.
* `State`	The state of the virtual private gateway (pending | available | deleting | deleted).
* `TagSet`	One or more tags associated with the virtual private gateway.
* `Type`	The type of VPN connection supported by the virtual private gateway (only ipsec.1 is supported)
* `VpnGatewayId`	The ID of the virtual private gateway.
* `requestId`	The ID of the request.

See detailed information in [Describe Instances](http://docs.outscale.com/api_fcu/definitions/VpnGateway.html#_api_fcu-vpngateway).
