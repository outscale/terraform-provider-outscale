---
layout: "outscale"
page_title: "OUTSCALE: outscale_lin_peering_acceptation"
sidebar_current: "docs-outscale-resource-lin-peering-acceptation"
description: |-
	Accepts a VPC peering connection request.
---

# outscale_lin_peering_acceptation

Accepts a VPC peering connection request.
To accept this request, you must be the owner of the peer VPC.

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

// Accepter's side of the connection.
resource "outscale_lin_peering_acceptation" "peer" {
    vpc_peering_connection_id = "${outscale_lin_peering.foo.id}"

    tag {
       Side = "Accepter"
    }
}
```

## Argument Reference

The following arguments are supported:

* `vpc_peering_connection_id` - The ID of the VPC peering connection you want to accept.

## Attributes

* `status` - The state of the VPC peering connection.
* `accepter_vpc_info` - Information about the peer VPC of the VPC peering connection.
* `requester_vpc_info` - Information about the source VPC of the VPC peering connection.
* `tag_set.N` - One or more tags associated with the VPC peering connection.
* `vpc_peering_connection_id` - The ID of the VPC peering connection.

[See detailed information.](http://docs.outscale.com/api_fcu/definitions/VpcPeeringConnection.html#_api_fcu-vpcpeeringconnection)
