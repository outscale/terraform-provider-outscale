---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_route_table_link"
sidebar_current: "docs-outscale-resource-route-table-link"
description: |-
  [Manages a route table link.]
---

# outscale_route_table_link Resource

Manages a route table link.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Route+Tables).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#linkroutetable).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `route_table_id` - (Required) The ID of the route table.
* `subnet_id` - (Required) The ID of the Subnet.

## Attribute Reference

The following attributes are exported:

* `link_route_table_id` - The ID of the route table association.
* `main` - If true, the route table is the main one.
* `route_table_id` - The ID of the route table.
* `subnet_id` - The ID of the subnet.
