---
layout: "outscale"
page_title: "OUTSCALE: outscale_lin_peering"
sidebar_current: "docs-outscale-resource-lin-peering"
description: |-
	Requests a VPC peering connection between a VPC you own and a peer VPC that can belong to another Outscale account.
---

# outscale_lin_peering

Requests a VPC peering connection between a VPC you own and a peer VPC that can belong to another Outscale account.
The two VPCs must not have overlapping CIDR blocks, otherwise the VPC peering connection is created with a failed status.
The created VPC peering connection remains in the pending-acceptance state until it is accepted by the owner of the peer VPC. The request expires after seven days.

## Example Usage

```hcl
resource "outscale_lin" "foo" {
	cidr_block = "10.0.0.0/16"
	tag {
		Name = "TestAccOutscaleLinPeeringConnection_basic"
	}
}

resource "outscale_lin" "bar" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_lin_peering" "foo" {
	vpc_id = "${outscale_lin.foo.id}"
	peer_vpc_id = "${outscale_lin.bar.id}"
}
```

## Argument Reference

The following arguments are supported:

* `peer_owner_id` - (Optional) The ID of the owner of the peer VPC.
* `peer_vpc_id` - (Optional) The ID of the peer VPC with which you want to connect the source VPC.
* `vpc_id` - (Optional) The ID of the requester VPC.

## Attributes

* `status` - The state of the VPC peering connection.
* `accepter_vpc_info` - Information about the peer VPC of the VPC peering connection.
* `requester_vpc_info` - Information about the source VPC of the VPC peering connection.
* `tag_set.N` - One or more tags associated with the VPC peering connection.
* `vpc_peering_connection_id` - The ID of the VPC peering connection.

[See detailed information.](http://docs.outscale.com/api_fcu/operations/Action_CreateVpcPeeringConnection_get.html#_api_fcu-action_createvpcpeeringconnection_get)
