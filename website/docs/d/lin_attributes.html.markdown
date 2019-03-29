---
layout: "outscale"
page_title: "OUTSCALE: outscale_lin_attributes"
sidebar_current: "docs-outscale-datasource-lin-attributes"
description: |-
Describes a specified attribute of a VPC.


---

# outscale_lin_attributes

Describes a specified attribute of a VPC.


## Example Usage

```hcl
data "outscale_lin_attributes" "test" {
	vpc_id = "${outscale_lin.vpc.id}"
	attribute = "enableDnsSupport"
}
```

## Argument Reference

The following arguments are supported:

* `attribute` (Required) The Attribute name (enableDnsSupport or enableDnsHostnames)
* `vpc_id` - (Required) The ID of the VPC.
* `filter.N` - (Optional) One or more filters.


See detailed information in [Outscale linAttributes](http://docs.outscale.com/api_fcu/operations/Action_DescribeVpcAttribute_get.html#_api_fcu-action_describevpcattribute_get).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `EnableDnsSupport`  Whether a DNS resolution is supported for the VPC.
* `EnableDnsHostNames`  Whether the instances launched in the VPC get DNS hostnames.


## Attributes Reference

The following attributes are exported:

* `enable_dns_hostnames` - The status of the enableDnsHostnames attribute.
* `enable_dns_support` - The status of the enableDnsSupport attribute.
* `vpc_id` - The ID of the VPC.
* `request_id` - The ID of the request.

See detailed information in [Describe linAttributes](http://docs.outscale.com/api_fcu/operations/Action_DescribeVpcAttribute_get.html#_api_fcu-action_describevpcattribute_get).
