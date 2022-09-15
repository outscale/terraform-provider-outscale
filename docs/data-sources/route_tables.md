---
layout: "outscale"
page_title: "OUTSCALE: outscale_route_tables"
sidebar_current: "outscale-route-tables"
description: |-
  [Provides information about route tables.]
---

# outscale_route_tables Data Source

Provides information about route tables.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Route-Tables.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-routetable).

## Example Usage

```hcl
data "outscale_route_tables" "route_tables01" {
  filter {
    name   = "route_table_ids"
    values = ["rtb-12345678", "rtb-12345679"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `link_route_table_ids` - (Optional) The IDs of the route tables involved in the associations.
    * `link_route_table_link_route_table_ids` - (Optional) The IDs of the associations between the route tables and the Subnets.
    * `link_route_table_main` - (Optional) If true, the route tables are the main ones for their Nets.
    * `link_subnet_ids` - (Optional) The IDs of the Subnets involved in the associations.
    * `net_ids` - (Optional) The IDs of the Nets for the route tables.
    * `route_creation_methods` - (Optional) The methods used to create a route.
    * `route_destination_ip_ranges` - (Optional) The IP ranges specified in routes in the tables.
    * `route_destination_service_ids` - (Optional) The service IDs specified in routes in the tables.
    * `route_gateway_ids` - (Optional) The IDs of the gateways specified in routes in the tables.
    * `route_nat_service_ids` - (Optional) The IDs of the NAT services specified in routes in the tables.
    * `route_net_peering_ids` - (Optional) The IDs of the Net peering connections specified in routes in the tables.
    * `route_states` - (Optional) The states of routes in the route tables (`active` \| `blackhole`). The `blackhole` state indicates that the target of the route is not available.
    * `route_table_ids` - (Optional) The IDs of the route tables.
    * `route_vm_ids` - (Optional) The IDs of the VMs specified in routes in the tables.
    * `tag_keys` - (Optional) The keys of the tags associated with the route tables.
    * `tag_values` - (Optional) The values of the tags associated with the route tables.
    * `tags` - (Optional) The key/value combination of the tags associated with the route tables, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.

## Attribute Reference

The following attributes are exported:

* `route_tables` - Information about one or more route tables.
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
