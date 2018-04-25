---
layout: "outscale"
page_title: "OUTSCALE: outscale_nat_service"
sidebar_current: "docs-outscale-nat-service"
description: |-
  Creates a network address translation (NAT) gateway in the specified public subnet of a VPC.
---

# outscale_image

Creates a network address translation (NAT) gateway in the specified public subnet of a VPC.

## Example Usage

```hcl
resource "outscale_lin" "vpc" {
	cidr_block = "10.0.0.0/16"
}
resource "outscale_subnet" "subnet" {
	cidr_block = "10.0.0.0/24"
	vpc_id = "${outscale_lin.vpc.id}"
}

resource "outscale_public_ip" "bar" {
	domain = "standard"
}

resource "outscale_nat_service" "gateway" {
    allocation_id = "${outscale_public_ip.bar.allocation_id}"
    subnet_id = "${outscale_subnet.subnet.id}"
}
```

## Argument Reference

The following arguments are supported:

* `allocation_id` - (Required) The allocation ID of the EIP to associate with the NAT gateway.
If the EIP is already associated with another resource, you must first disassociate it.
* `client_token` - (Optional) A unique identifier to manage the idempotency.
* `subnet_id` - (Required) The public subnet where you want to create the NAT gateway.

# Attributes

* `natGatewayAddress`	- Information about the External IP address (EIP) associated with the NAT gateway.
* `natGatewayId` -	The ID of the NAT gateway.
* `state`	- The state of the NAT gateway (pending | available| deleting | deleted).	
* `subnetId` -	The ID of the public subnet in which the NAT gateway is.
* `vpcId`	- The ID of the VPC in which the NAT gateway is.
* `requestId` -	The ID of the request.


See detailed information in [Create Nat Service](http://docs.outscale.com/api_fcu/operations/Action_CreateNatGateway_get.html#_api_fcu-action_createnatgateway_get).

