---
layout: "outscale"
page_title: "OUTSCALE: outscale_nat_service"
sidebar_current: "docs-outscale-resource-nat-service"
description: |-
  Provides an Outscale Nat Gateway resource. This allows instances to be created, described, and deleted. Nat Gateway also support provisioning.
---

# outscale_nat_service

  Provides an Outscale Nat Gateway resource. This allows instances to be created, described, and deleted. Nat Gateway also support [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
resource "outscale_nat_service" "gateway" {
    reservation_id = "eipalloc-32e506e8"
    subnet_id = "subnet-861fbecc"
}
```

## Argument Reference

The following arguments are supported:

* `allocation_id` - (Requiered) The allocation ID of the EIP to associate with the NAT gateway. 
* `client_token` - (Optional) A unique identifier which enables you to manage the idempotency.
* `subnet_id` - (Requiered) The public subnet where you want to create the NAT gateway..


See detailed information in [Outscale Instances](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).


## Attributes Reference

The following attributes are exported:

* `nat_gateway_address` - A unique identifier which enables you to manage the idempotency.
* `nat_gateway_id` - Information about the newly created NAT gateway.
* `state` - The ID of the request.
* `subnet_id` - Information about the newly created NAT gateway.
* `vpc_id` - The ID of the request.

See detailed information in [Tags](http://docs.outscale.com/api_fcu/definitions/NatGateway.html#_api_fcu-natgateway.
