---
layout: "outscale"
page_title: "OUTSCALE: outscale_lin_internet_gateway"
sidebar_current: "docs-outscale-resource-lin-internet-gateway"
description: |-
  Creates an Internet gateway you can use with a VPC.
---

# outscale_lin_internet_gateway

An Internet gateway enables your instances launched in a VPC to connect to the Internet. By default, a VPC includes an Internet gateway, and each subnet is public. Every instance launched within a default subnet has a private and a public IP addresses.

## Example Usage

```hcl
resource "outscale_lin_internet_gateway" "gateway" {}
```
## Attributes Reference

The following attributes are exported:

* `attachement_set` - One or more VPCs attached to the Internet gateway.
* `internet_gateway_id` - The ID of the Internet gateway.
* `tag_set` - One or more tags associated with the Internet gateway.
* `request_id` - The ID of the request.

See detailed information in [Create Internet Gateway](http://docs.outscale.com/api_fcu/operations/Action_CreateInternetGateway_get.html#_api_fcu-action_createinternetgateway_get).
