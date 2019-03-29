---
layout: "outscale"
page_title: "OUTSCALE: outscale_route_table"
sidebar_current: "docs-outscale-resource-route-table"
description: |-
  Creates a route table for a specified VPC. You can then add routes and associate this route table with a subnet.
---

# outscale_route_table

Creates a route table for a specified VPC. You can then add routes and associate this route table with a subnet.

## Example Usage

```hcl
resource "outscale_lin" "foo" {
    cidr_block = "10.1.0.0/16"
}

resource "outscale_lin_internet_gateway" "foo" {
    #vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route_table" "foo" {
    vpc_id = "${outscale_lin.foo.id}"
}
```

## Argument Reference

The following arguments are supported:

* `vpc_id` -	(Required)	The ID of the VPC.	false	string

## Attributes

* `association_set` -	One or more associations between the route table and the subnets.
* `propagating_vgw_set` -	Information about virtual private gateways propagating routes.
* `route_set` -	One or more routes in the route table.
* `route_table_id` -	The ID of the route table.
* `tag_set` -	One or more tags associated with the route table.
* `vpc_id` -	The ID of the VPC.
* `request_id` -	The ID of the request

[See detailed information](http://docs.outscale.com/api_fcu/operations/Action_CreateRouteTable_get.html#_api_fcu-action_createroutetable_get).
