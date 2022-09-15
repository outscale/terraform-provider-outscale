---
layout: "outscale"
page_title: "OUTSCALE: outscale_virtual_gateway_route_propagation"
sidebar_current: "outscale-virtual-gateway-route-propagation"
description: |-
  [Manages a virtual gateway route propagation.]
---

# outscale_virtual_gateway_route_propagation Resource

Manages a virtual gateway route propagation.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Routing-Configuration-for-VPN-Connections.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#updateroutepropagation).

## Example Usage

### Required resources

```hcl
resource "outscale_virtual_gateway" "virtual_gateway01" {
	connection_type = "ipsec.1"
}

resource "outscale_net" "net01" {
	ip_range = "10.0.0.0/16"
}

resource "outscale_route_table" "route_table01" {
	net_id = outscale_net.net01.net_id
}

resource "outscale_virtual_gateway_link" "virtual_gateway_link01" {
	virtual_gateway_id = outscale_virtual_gateway.virtual_gateway01.virtual_gateway_id
	net_id             = outscale_net.net01.net_id
}
```

### Activate the propagation of routes to a route table of a Net by a virtual gateway

```hcl
resource "outscale_virtual_gateway_route_propagation" "virtual_gateway_route_propagation01" {
	enable             = true
	virtual_gateway_id = outscale_virtual_gateway.virtual_gateway01.virtual_gateway_id
	route_table_id     = outscale_route_table.route_table01.route_table_id
	depends_on         = [outscale_virtual_gateway_link.virtual_gateway_link01]
}
```

## Argument Reference

The following arguments are supported:

* `enable` - (Required) If true, a virtual gateway can propagate routes to a specified route table of a Net. If false, the propagation is disabled.
* `route_table_id` - (Required) The ID of the route table.
* `virtual_gateway_id` - (Required) The ID of the virtual gateway.

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

