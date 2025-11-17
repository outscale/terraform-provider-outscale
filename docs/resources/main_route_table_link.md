---
layout: "outscale"
page_title: "OUTSCALE: outscale_main_route_table_link"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-main-route-table-link"
description: |-
  [Manages a main route table link.]
---

# outscale_main_route_table_link Resource

Manages a main route table link.


~> **Note:** On Net creation, the OUTSCALE API always creates an initial main route table. The `main_route_table_link`resource records the ID of the inital route table under the `default_route_table_id` attribute. The "Destroy" action for a `main_route_table_link` consists of resetting the original route table as the main route table for the Net. The additional route table must remain intact in order for the `main_route_table_link` destroy to work properly.

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

### Link a main route table

```hcl
resource "outscale_main_route_table_link" "main" {
  net_id   = outscale_net.net01.net_id
  route_table_id = outscale_route_table.route_table01.route_table_id
}
```

## Argument Reference

The following arguments are supported:

* `net_id` - (Required) The ID of the Net.
* `route_table_id` - (Required) The ID of the route table.

## Attribute Reference

The following attributes are exported:

* `default_route_table_id` - The ID of the default route table.
* `link_route_table_id` - The ID of the association between the route table and the Subnet.
* `main` - If true, the route table is the main one.
* `net_id` - The ID of the Net.
* `route_table_id` - The ID of the route table.

