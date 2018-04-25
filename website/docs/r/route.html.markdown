---
layout: "outscale"
page_title: "OUTSCALE: outscale_route"
sidebar_current: "docs-outscale-resource-route"
description: |-
	Creates a route in a specified route table within a specified VPC.
---

# outscale_route_table

Creates a route in a specified route table within a specified VPC.

## Example Usage

```hcl
resource "outscale_lin" "test" {
  cidr_block = "10.10.0.0/16"
}

resource "outscale_route_table" "test" {
  vpc_id = "${outscale_lin.test.id}"
}

resource "outscale_subnet" "test" {
  vpc_id = "${outscale_lin.test.id}"
  cidr_block = "10.10.10.0/24"
}

resource "outscale_route" "test" {
  route_table_id = "${outscale_route_table.test.id}"
  destination_cidr_block = "0.0.0.0/0"
  instance_id = "${outscale_vm.nat.id}"
}

resource "outscale_vm" "nat" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
  subnet_id = "${outscale_subnet.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `destination_cidr_block` - (Required)	The CIDR block used for the destination match.
* `gateway_id` -	(Optional)	The ID of an Internet gateway or virtual private gateway attached to your VPC.
* `instance_id` -	(Optional)	The ID of a NAT instance in your VPC (attached to exactly one network interface).
* `nat_gateway_id` -	(Optional)	The ID of a NAT gateway.
* `network_interface_id` -	(Optional)	The ID of a network interface.
* `route_table_id` -	(Required)	The ID of the route table.
* `vpc_peering_connection_id` -	(Optional)	The ID of a VPC peering connection.

## Attributes

* `destination_cidr_block` -	The CIDR block used for the destination match.
* `gateway_id` -	The ID of an Internet gateway or virtual private gateway attached to your VPC.
* `instance_id` -	The ID of a NAT instance in your VPC (attached to exactly one network interface).
* `natGateway_id` -	The ID of a NAT gateway.
* `network_interface_id` -	The ID of a network interface.
* `route_table_id` -	The ID of the route table.
* `vpc_peering_connection_id` -	The ID of a VPC peering connection.
* `destination_prefix_list_id` -	The prefix of the Outscale service.
* `instance_owner_id` -	The account ID of the owner of the instance.
* `origin` -	The method used to create the route.
* `state` -	The state of the route.
* `request_id` -	The ID of the request


See detailed information in [Create Route](http://docs.outscale.com/api_fcu/operations/Action_CreateRoute_get.html#_api_fcu-action_createroute_get).
