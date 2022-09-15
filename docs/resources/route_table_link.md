---
layout: "outscale"
page_title: "OUTSCALE: outscale_route_table_link"
sidebar_current: "outscale-route-table-link"
description: |-
  [Manages a route table link.]
---

# outscale_route_table_link Resource

Manages a route table link.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Route-Tables.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-routetable).

## Example Usage

### Required resources

```hcl
resource "outscale_net" "net01" {
	ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet01" {
	net_id   = outscale_net.net01.net_id
	ip_range = "10.0.0.0/18"
}

resource "outscale_route_table" "route_table01" {
	net_id = outscale_net.net01.net_id
}
```

### Link a route table to a subnet

```hcl
resource "outscale_route_table_link" "route_table_link01" {
	subnet_id      = outscale_subnet.subnet01.subnet_id
	route_table_id = outscale_route_table.route_table01.route_table_id
}
```

## Argument Reference

The following arguments are supported:

* `route_table_id` - (Required) The ID of the route table.
* `subnet_id` - (Required) The ID of the Subnet.

## Attribute Reference

The following attributes are exported:

* `link_route_table_id` - The ID of the association between the route table and the Subnet.
* `main` - If true, the route table is the main one.
* `route_table_id` - The ID of the route table.
* `subnet_id` - The ID of the Subnet.

## Import

A route table link can be imported using the route table ID and the route table link ID. For example:

```console

$ terraform import outscale_route_table_link.ImportedRouteTableLink rtb-12345678_rtbassoc-87654321

```