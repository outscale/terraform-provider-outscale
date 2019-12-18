---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_route"
sidebar_current: "outscale-route"
description: |-
  [Manages a route.]
---

# outscale_route Resource

Manages a route.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Route+Tables#AboutRouteTables-RoutesRoutesandRoutingOptions).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-route).

## Example Usage

```hcl

#resource "outscale_internet_service" "internet_service01" {
#}

#resource "outscale_internet_service_link" "internet_service_link01" {
#  internet_service_id = outscale_internet_service.internet_service01.internet_service_id
#  net_id              = outscale_net.net01.net_id
#}

resource "outscale_route" "route01" {
  gateway_id           = outscale_internet_service.internet_service01.id
  destination_ip_range = "20.0.0.0/16"
  route_table_id       = outscale_route_table.route_table01.route_table_id
}


```

## Argument Reference

The following arguments are supported:

* `destination_ip_range` - (Required) The IP range used for the destination match, in CIDR notation (for example, 10.0.0.0/24).
* `gateway_id` - (Optional) The ID of an Internet service or virtual gateway attached to your Net.
* `nat_service_id` - (Optional) The ID of a NAT service.
* `net_peering_id` - (Optional) The ID of a Net peering connection.
* `nic_id` - (Optional) The ID of a NIC.
* `route_table_id` - (Required) The ID of the route table for which you want to create a route.
* `vm_id` - (Optional) The ID of a NAT VM in your Net (attached to exactly one NIC).

## Attribute Reference

The following attributes are exported:

* `route_table` - Information about the route table.
  * `link_route_tables` - One or more associations between the route table and Subnets.
    * `link_route_table_id` - The ID of the association between the route table and the Subnet.
    * `main` - If `true`, the route table is the main one.
    * `route_table_id` - The ID of the route table.
    * `subnet_id` - The ID of the Subnet.
  * `net_id` - The ID of the Net for the route table.
  * `route_propagating_virtual_gateways` - Information about virtual gateways propagating routes.
    * `virtual_gateway_id` - The ID of the virtual gateway.
  * `route_table_id` - The ID of the route table.
  * `routes` - One or more routes in the route table.
    * `creation_method` - The method used to create the route.
    * `destination_ip_range` - The IP range used for the destination match, in CIDR notation (for example, 10.0.0.0/24).
    * `destination_service_id` - The ID of the 3DS OUTSCALE service.
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
