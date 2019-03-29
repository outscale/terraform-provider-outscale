---
layout: "outscale"
page_title: "OUTSCALE: outscale_lin_internet_gateways"
sidebar_current: "docs-outscale-datasource-lin-internet-gateways"
description: |-
Describes one or more of your Internet gateways.
An Internet gateway enables your instances launched in a VPC to connect to the Internet. By default, a VPC includes an Internet gateway, and each subnet is public. Every instance launched within a default subnet has a private and a public IP addresses.


---

# outscale_lin_internet_gateways

Describes one or more of your Internet gateways.
An Internet gateway enables your instances launched in a VPC to connect to the Internet.


## Example Usage

```hcl
data "outscale_lin_internet_gateways" "outscale_lin_internet_gateways" {
  filter {
		name = "internet-gateway-id"
		values = ["${outscale_lin_internet_gateway.gateway.id}"]
	}
}
```

## Argument Reference

The following arguments are supported:

* `internet_gateway_id.N` One or more internet gateways IDS - (Optional)
* `filter.N` - (Optional) One or more filters.

See detailed information in [Outscale InternetGateways](http://docs.outscale.com/api_fcu/operations/Action_DescribeInternetGateways_get.html#_api_fcu-action_describeinternetgateways_get).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `attachment.state`  The ID of the NAT gateway.
* `attachment.vpc-id` The ID of an attachment VPC
* `internet-gateway-id` The ID of the internet Gateway
* `tag` The key/value combination of a tag that is assigned to the resource, in the following format: key = value
* `tag-key` The key of a tag that is assigned to the resource. You can use this filter alongside the tag-value filter. In that case, you filter the resources corresponding to each tag, regardless of the filter.
* `tag-value` The value of a tag that is assigned to the resource. You can use this filter alongside the tag-key filter. In the case, you filter the resource corresponding to each tag, regardlessof the other filter.



## Attributes Reference

The following attributes are exported:

* `attachment_set.N ` - One or more VPCs attached to the internet gateway associated with the NAT gateway.
* `internet_gateway_id` - The ID of the Internet gateway
* `tag_set.N` - One or more tags associated with the internet gateway.
* `request_id` - The ID of the request.

See detailed information in [Describe InternetGateways](http://docs.outscale.com/api_fcu/operations/Action_DescribeNatGateways_get.html#_api_fcu-action_describenatgateways_get).
