---
layout: "outscale"
page_title: "OUTSCALE: outscale_route_tables"
sidebar_current: "docs-outscale-datasource-route-tables"
description: |-
  Describes route tables
---

# outscale_route_tables

Describes one or more of your route tables.
In your Virtual Private Cloud (VPC), each subnet must be associated with a route table. If a subnet is not explicitly associated with a route table, it is implicitly associated with the main route table of the VPC.

## Example Usage

```hcl
resource "outscale_lin" "test" {
  cidr_block = "172.16.0.0/16"

  tag {
    Name = "terraform-testacc-data-source"
  }
}

resource "outscale_subnet" "test" {
  cidr_block = "172.16.0.0/24"
  vpc_id     = "${outscale_lin.test.id}"
  tag {
    Name = "terraform-testacc-data-source"
  }
}

resource "outscale_route_table" "test" {
  vpc_id = "${outscale_lin.test.id}"
  tag {
    Name = "terraform-testacc-routetable-data-source"
  }
}

data "outscale_route_tables" "by_filter" {
  filter {
    name = "route-table-id"
    values = ["${outscale_route_table.test.id}"]
  }
}

data "outscale_route_tables" "by_id" {
  route_table_id = ["${outscale_route_table.test.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `Filter.N` (Optional). One or more filters.
* `RouteTableId.N` (Optional). One or more route table IDs.

## Filters

You can use the Filter.N parameter to filter the route tables on the following properties:

* `association.route-table-association-id` The ID of an association ID for the route table.
* `association.route-table-id` The ID of the route table involved in the association.
* `association.subnet-id` The ID of the subnet involved in the association.
* `association.main` Indicates whether the route table is the main route table for the VPC (true | false).
* `route-table-id` The ID of the route table.
* `route.destination-cidr-block` The CIDR range specified in a route in the table.
* `route.destination-prefix-list-id` The prefix ID of the service specified in a route in the table.
* `route.gateway-id` The ID of a gateway specified in a route in the table.
* `route.instance-id` The ID of an instance specified in a route in the table.
* `route.nat-gateway-id` The ID of a NAT gateway specified in a route in the table.
* `route.origin` How the route was created.
* `route.state` The state of a route in the route table (active | blackhole). The blackhole state indicates that the target of the route is not available.
* `route.vpc-peering-connection-id` The ID of a VPC peering connection specified in a route in the table.
* `tag` The key/value combination of a tag associated with the resource, in the following format: key=value.
* `tag-key` The key of a tag associated with the resource.
* `tag-value` The value of a tag associated with the resource.
* `vpc-id` The ID of the VPC for the route table.

## Attributes Reference

The following attributes are exported:

* `route_table_set.N` - Information about one or more route tables, each containing the following attributes:
  - `association_set.N` - One or more associations between the route table and the subnets.	false	RouteTableAssociation
  - `propagating_vgw_set.N` - Information about virtual private gateways propagating routes.	false	PropagatingVgw
  - `route_set.N` - One or more routes in the route table.
  - `route_table_id` - The ID of the route table.
  - `tag_set.N` - One or more tags associated with the route table.
  - `vpc_id` - The ID of the VPC.
  - `request_id` - The ID of the request

[See detailed description](http://docs.outscale.com/api_fcu/operations/Action_DescribeRouteTables_get.html#_api_fcu-action_describeroutetables_get)
