---
layout: "outscale"
page_title: "OUTSCALE: outscale_route_table_link"
sidebar_current: "docs-outscale-resource-route-table-link"
description: |-
	Associates a subnet with a route table
---

# outscale_route_table_link

Associates a subnet with a route table

## Example Usage

```hcl
resource "outscale_lin" "foo" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_subnet" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
	cidr_block = "10.1.1.0/24"
}

resource "outscale_route_table" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
}

resource "outscale_route_table_link" "foo" {
	route_table_id = "${outscale_route_table.foo.id}"
	subnet_id = "${outscale_subnet.foo.id}"
}
```

## Argument Reference

The following arguments are supported:

* `route_table_id` -	(Required)	The ID of the route table.
* `subnet_id` -	(Required)	The ID of the subnet.

## Attributes

* `route_table_id` - The ID of the route table.	true	string
* `subnet_id` -	The ID of the subnet.	true	string
* `association_id` -	The ID of the route table association.	false	string
* `request_id` -	The ID of the request	false	string



See detailed information in [Associates Route Table Link](http://docs.outscale.com/api_fcu/operations/Action_AssociateRouteTable_get.html#_api_fcu-action_associateroutetable_get).
