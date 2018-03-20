---
layout: "outscale"
page_title: "OUTSCALE: outscale_lin_internet_gateway_link"
sidebar_current: "docs-outscale-resource-lin-internet-gateway-link"
description: |-
  Attaches an Internet gateway to a VPC.
---

# outscale_lin_internet_gateway_link

To enable the connection between the Internet and a VPC, you must attach an Internet gateway to this VPC.

## Example Usage

```hcl
resource "outscale_lin_internet_gateway" "gateway" {}

resource "outscale_lin" "vpc" {
	cidr_block = "10.0.0.0/16"
}

resource "outscale_lin_internet_gateway_link" "link" {
	vpc_id = "${outscale_lin.vpc.id}"
	internet_gateway_id = "${outscale_lin_internet_gateway.gateway.id}"
}
```

## Argument Reference

The following arguments are supported:

* `internet_gateway_id` - The ID of the Internet gateway.
* `vpc_id` - The ID of the VPC.

See detailed information in [Outscale Attach Internet Gateway](http://docs.outscale.com/api_fcu/operations/Action_AttachInternetGateway_get.html#_api_fcu-action_attachinternetgateway_get).


## Attributes Reference

The following attributes are exported:

* `request_id` - The CIDR block of the VPC, in the [16;28] range (for example, 10.84.7.0/24).

See detailed information in [Outscale Attach Internet Gateway](http://docs.outscale.com/api_fcu/operations/Action_AttachInternetGateway_get.html#_api_fcu-action_attachinternetgateway_get).
