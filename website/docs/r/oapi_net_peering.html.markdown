---
layout: "outscale"
page_title: "OUTSCALE: outscale_net_peering"
sidebar_current: "docs-outscale-resource-net-peering"
description: |-
  Requests a VPC peering connection between a VPC you own and a peer VPC that can belong to another Outscale account.
---

# outscale_net_peering

Requests a VPC peering connection between a VPC you own and a peer VPC that can belong to another Outscale account.
The two VPCs must not have overlapping CIDR blocks, otherwise the VPC peering connection is created with a failed status.
The created VPC peering connection remains in the pending-acceptance state until it is accepted by the owner of the peer VPC. The request expires after seven days.

## Example Usage

```hcl
resource "outscale_net" "foo" {
	cidr_block = "10.0.0.0/16"
	tag {
		Name = "TestAccOutscaleOAPInetPeeringConnection_basic"
	}
}

resource "outscale_net" "bar" {
	cidr_block = "10.1.0.0/16"
}

resource "outscale_net_peering" "foo" {
	vpc_id = "${outscale_net.foo.id}"
	peer_vpc_id = "${outscale_net.bar.id}"
}
```

## Argument Reference

The following arguments are supported:

* `source_net_account_id` - (Optional) The ID of the owner of the peer VPC.
* `accepter_net_id` - (Optional) The ID of the peer VPC with which you want to connect the source VPC.
* `source_net_id` - (Optional) The ID of the requester VPC.

## Attributes

* `accepter_net` - Information about the peer VPC of the VPC peering connection.
  * `ip_range` - The CIDR block of the VPC.
  * `account_id` - The account ID of the owner of the VPC.
  * `net_id` - The ID of the VPC.
* `source_net` - Information about the source VPC of the VPC peering connection.
  * `ip_range` - The CIDR block of the VPC.
  * `account_id` - The account ID of the owner of the VPC.
  * `net_id` - The ID of the VPC.
* `state` - The state of the VPC peering connection.
  * `name` - The state of the VPC peering connection (pending-acceptance | active| deleted | rejected | failed | expired | deleted).
  * `message` - Additional information about the state of the VPC peering connection.
* `tags` - One or more tags associated with the VPC peering connection.
  * `key` - The key of the tag.	
  * `value` - The value of the tag
* `net_peering_id` - The ID of the VPC peering connection.
* `request_id` - The ID of the request
