---
layout: "outscale"
page_title: "OUTSCALE: outscale_route_table_link"
sidebar_current: "docs-outscale-resource-route-table-link"
description: |-
	Associates a Subnet with a route table.
---

# outscale_route_table_link

The Subnet and the route table must be in the same Net. The traffic is routed according to the route table defined within this Net. You can associate a route table with several Subnets.

## Example Usage

```hcl
resource "outscale_net" "foo" {
	ip_range = "10.1.0.0/16"
}

resource "outscale_subnet" "foo" {
	net_id = "${outscale_net.foo.id}"
	ip_range = "10.1.1.0/24"
}

resource "outscale_route_table" "foo" {
	net_id = "${outscale_net.foo.id}"
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

* `route_table_id` - The ID of the route table.
* `subnet_id` -	The ID of the subnet.
* `link_id` -	The ID of the route table association.
* `request_id` -	The ID of the request.

See detailed information in [Associates Route Table Link](http://docs.outscale.com/api_fcu/operations/Action_AssociateRouteTable_get.html#_api_fcu-action_associateroutetable_get).
