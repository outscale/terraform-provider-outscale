---
layout: "outscale"
page_title: "OUTSCALE: outscale_route_tables"
sidebar_current: "docs-outscale-datasource-route-tables"
description: |-
  Lists one or more of your route tables.
---

# outscale_route_tables

In your Net, each Subnet must be associated with a route table. If a Subnet is not explicitly associated with a route table, it is implicitly associated with the main route table of the Net.

## Example Usage

```hcl
resource "outscale_net" "test" {
    ip_range = "172.16.0.0/16"

    tags {
        key = "Name"
        value = "terraform-testacc-data-source"
    }
}

resource "outscale_subnet" "test" {
    ip_range = "172.16.0.0/24"
    net_id     = "${outscale_net.test.id}"
}

resource "outscale_route_table" "test" {
    net_id = "${outscale_net.test.id}"
    tags {
        key = "Name"
        value = "terraform-testacc-routetable-data-source"
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

* `Filters` (Optional). One or more filters.

## Filters

You can use the Filter.N parameter to filter the route tables on the following properties:

* `link_route_table.route-table-ids` The ID(s) of the associations between the route table and the Subnets.
* `link_route_table.route-table-ids` The ID(s) of the route table(s).
* `link_route_table.link-subnet-ids` The ID(s) of the Subnet(s) involved in the associations.
* `routes.destination-ip-ranges` The IP range(s) specified in a route in the table.
* `routes.destination-prefix-list-ids` The prefix ID(s) of the service(s) specified in routes in the tables.
* `routes.gateway-ids` The ID(s) of the gateway(s) specified in routes in the tables.
* `routes.vm-ids` The ID(s) of the VM(s) specified in routes in the tables.
* `routes.nat-service-ids` The ID(s) of the NAT service(s) specified in routes in the tables.
* `routes.states` The state(s) of routes in the route table (active | blackhole). The blackhole state indicates that the target of the route is not available.
* `routes.net-peering-ids` The ID(s) of the Net peering connection(s) specified in routes in the tables.
* `tags` The key/value combination of the tag(s) associated with your resources, in the following format: "Filters":{"Tags":["TAGKEY=TAGVALUE"]}.
* `tag-keys` The key(s) of the tag(s) associated with your resources.
* `tag-value` The value(s) of the tag(s) associated with your resources.
* `net-ids` The ID(s) of the Net(s) for the route tables.

## Attributes Reference

The following attributes are exported:

* `link_route_tables` - One or more associations between the route table and the subnets.
* `route_propagating_virtual_gateways` - Information about virtual private gateways propagating routes.
* `routes` - One or more routes in the route table.
* `route_table_id` - The ID of the route table.
* `tags` - One or more tags associated with the route table.
* `net_id` - The ID of the VPC.
* `request_id` - The ID of the request

[See detailed description](http://docs.outscale.com/api_fcu/operations/Action_DescribeRouteTables_get.html#_api_fcu-action_describeroutetables_get)
