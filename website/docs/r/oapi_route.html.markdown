---
layout: "outscale"
page_title: "OUTSCALE: outscale_route"
sidebar_current: "docs-outscale-resource-route"
description: |-
	Creates a route in a specified route table within a specified Net.
---

# outscale_route_table

Creates a route in a specified route table within a specified Net.

## Example Usage

```hcl
resource "outscale_net" "test" {
  ip_range = "10.10.0.0/16"
}

resource "outscale_route_table" "test" {
  net_id = "${outscale_net.test.id}"
}

resource "outscale_subnet" "test" {
  net_id = "${outscale_net.test.id}"
  subregion_name = "in-west-2a"
  ip_range = "10.10.10.0/24"
}

resource "outscale_route" "test" {
  route_table_id = "${outscale_route_table.test.id}"
  destination_ip_range = "0.0.0.0/0"
  vm_id = "${outscale_vm.nat.id}"
}

resource "outscale_vm" "nat" {
  image_id = "ami-8a6a0120"
  instance_type = "t2.micro"
  subnet_id = "${outscale_subnet.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `destination_ip_range` - (Required)   The IP range used for the destination match, in CIDR notation (for example, 10.0.0.0/24).
* `gateway_id` -	(Optional)	The ID of the Internet service or virtual gateway attached to the Net.
* `nat_service_id` -	(Optional)	The ID of a NAT service attached to the Net.
* `vm_id` -	(Optional)	The ID of a VM specified in a route in the table.
* `nic_id` - (Optional) The ID of the NIC.
* `net_peering_id` -	(Optional)	The ID of the Net peering connection.
* `route_table_id` -	(Required)	The ID of the route table.

## Attributes

* `destination_ip_range` -	The IP range used for the destination match, in CIDR notation (for example, 10.0.0.0/24).
* `gateway_id` -	The ID of the Internet service or virtual gateway attached to the Net.
* `vm_id` - The ID of a VM specified in a route in the table.
* `nat_service_id` -	The ID of a NAT service attached to the Net.
* `nic_id` -	The ID of the NIC.
* `route_table_id` -	The ID of the route table.
* `net_peering_id` -	The ID of the Net peering connection.
* `destination_prefix_list_id` -    The prefix ID(s) of the service(s) specified in routes in the tables.
* `vm_account_id` -	The account ID of the owner of the VM.
* `creation_method` -	The method used to create the route.
* `state` - The state of a route in the route table (active | blackhole). The blackhole state indicates that the target of the route is not available.
* `request_id` -	The ID of the request


See detailed information in [Create Route](https://docs-beta.outscale.com/#createroute).
