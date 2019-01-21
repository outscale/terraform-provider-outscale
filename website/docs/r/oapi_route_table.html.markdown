---
layout: "outscale"
page_title: "OUTSCALE: outscale_route_table"
sidebar_current: "docs-outscale-resource-route-table"
description: |-
  Creates a route table for a specified Net.
---

# outscale_route_table

Creates a route table for a specified Net. You can then add routes and associate this route table with a subnet.

## Example Usage

```hcl
resource "outscale_net" "foo" {
    ip_range = "10.1.0.0/16"
}

resource "outscale_internet_service" "foo" {}

resource "outscale_route_table" "foo" {
    net_id = "${outscale_net.foo.id}"
}
```

## Argument Reference

The following arguments are supported:

* `net_id` -	(Required)	The ID of the Net for which you want to create a route table.

## Attributes

* `net_id` -	The ID of the Net for which you want to create a route table.
* `route_table_id` -	The ID of the route table.
* `tags` -	One or more tags associated with the route table.
* `link_route_tables` -	One or more associations between the route table and the subnets.
* `route_propagating_virtual_gateways` -    Information about virtual private gateways propagating routes.
* `routes` -	One or more routes in the route table.
* `request_id` -	The ID of the request

[See detailed information](http://docs.outscale.com/api_fcu/operations/Action_CreateRouteTable_get.html#_api_fcu-action_createroutetable_get).
