---
layout: "outscale"
page_title: "OUTSCALE: outscale_nat_service"
sidebar_current: "docs-outscale-datasource-nat-service"
description: |-
    Describes one or more network address translation (NAT) gateways.


---

# outscale_nat_service

Describes one or more network address translation (NAT) gateways.


## Example Usage

```hcl
data "outscale_nat_services" "nat" {
	nat_gateway_id = ["nat-08f41400"]
}
```

## Argument Reference

The following arguments are supported:

* `nat_gateway_id` - (Optional) One or more IDs of NAT gateways.

See detailed information in [Outscale Instances](http://docs.outscale.com/api_fcu/operations/Action_CreateNatGateway_get.html#_api_fcu-action_createnatgateway_get).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `nat-gateway-id` - The ID of the NAT gateway.
* `state` The state of the NAT gateway (pending | available | deleting | deleted).
* `subnet-id` The ID of the subnet in which the NAT gateway is.
* `vpc-id` The ID of the VPC in which the NAT gateway is.


## Attributes Reference

The following attributes are exported:

* `nat_gateway_address ` - Information about the External IP address (EIP) associated with the NAT gateway.
* `nat_gateway_id` - The ID of the NAT gateway.
* `state` - The state of the NAT gateway (pending | available| deleting | deleted).
* `subnet-id` - The ID of the subnet in which the NAT gateway is.
* `vpc_id` - The ID of the VPC in which the NAT gateway is.
* `request_id` - The ID of the request.

See detailed information in [Describe Instances](http://docs.outscale.com/api_fcu/operations/Action_DescribeNatGateways_get.html#_api_fcu-action_describenatgateways_get).
