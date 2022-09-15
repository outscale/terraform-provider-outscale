---
layout: "outscale"
page_title: "OUTSCALE: outscale_route_table"
sidebar_current: "outscale-route-table"
description: |-
  [Manages a route table.]
---

# outscale_route_table Resource

Manages a route table.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Route-Tables.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-routetable).

## Example Usage

### Required resource

```hcl
resource "outscale_net" "net01" {
	ip_range = "10.0.0.0/16"
}
```

### Create a route table

```hcl
resource "outscale_route_table" "route_table01" {
	net_id = outscale_net.net01.net_id
}
```

## Argument Reference

The following arguments are supported:

* `net_id` - (Required) The ID of the Net for which you want to create a route table.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `link_route_tables` - One or more associations between the route table and Subnets.
    * `link_route_table_id` - The ID of the association between the route table and the Subnet.
    * `main` - If true, the route table is the main one.
    * `route_table_id` - The ID of the route table.
    * `subnet_id` - The ID of the Subnet.
* `net_id` - The ID of the Net for the route table.
* `route_propagating_virtual_gateways` - Information about virtual gateways propagating routes.
    * `virtual_gateway_id` - The ID of the virtual gateway.
* `route_table_id` - The ID of the route table.
* `routes` - One or more routes in the route table.
    * `creation_method` - The method used to create the route.
    * `destination_ip_range` - The IP range used for the destination match, in CIDR notation (for example, 10.0.0.0/24).
    * `destination_service_id` - The ID of the OUTSCALE service.
    * `gateway_id` - The ID of the Internet service or virtual gateway attached to the Net.
    * `nat_service_id` - The ID of a NAT service attached to the Net.
    * `net_access_point_id` - The ID of the Net access point.
    * `net_peering_id` - The ID of the Net peering connection.
    * `nic_id` - The ID of the NIC.
    * `state` - The state of a route in the route table (`active` \| `blackhole`). The `blackhole` state indicates that the target of the route is not available.
    * `vm_account_id` - The account ID of the owner of the VM.
    * `vm_id` - The ID of a VM specified in a route in the table.
* `tags` - One or more tags associated with the route table.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Import

A route table can be imported using its ID. For example:

```console

$ terraform import outscale_route_table.ImportedRouteTable rtb-12345678

```