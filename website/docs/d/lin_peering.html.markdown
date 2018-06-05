---
layout: "outscale"
page_title: "OUTSCALE: outscale_lin_peering"
sidebar_current: "docs-outscale-datasource-lin-peering"
description: |-
  Describes one or more peering connections between two Virtual Private Clouds (VPCs).
---

# outscale_lin_peering

Describes one or more peering connections between two Virtual Private Clouds (VPCs).

## Example Usage

```hcl
resource "outscale_lin" "foo" {
  cidr_block = "10.1.0.0/16"

  tag {
      Name = "terraform-testacc-vpc-peering-connection-data-source-foo"
  }
}

resource "outscale_lin" "bar" {
  cidr_block = "10.2.0.0/16"

  tag {
      Name = "terraform-testacc-vpc-peering-connection-data-source-bar"
  }
}

resource "outscale_lin_peering" "test" {
    vpc_id = "${outscale_lin.foo.id}"
    peer_vpc_id = "${outscale_lin.bar.id}"

    tag {
      Name = "terraform-testacc-vpc-peering-connection-data-source-foo-to-bar"
    }
}

data "outscale_lin_peerings" "test_by_id" {
    vpc_peering_connection_id = ["${outscale_lin_peering.test.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `filter.N` (Optional). One or more filters.
* `vpc_peering_connection_id.N` (Optional). One or more VPC peering connection IDs.

## Filters

You can use the Filter.N parameter to filter the described VPC peering connections on the following properties:

* `accepter-vpc-info.cidr-block` - The CIDR block of the peer VPC.
* `accepter-vpc-info.owner-id` - The account ID of the owner of the peer VPC.
* `accepter-vpc-info.vpc-id` - The ID of the peer VPC.
* `expiration-time` - The date after which the connection expires.
* `requester-vpc-info.cidr-block` - The CIDR block of the requester VPC.
* `requester-vpc-info.owner-id` - The ID of the owner of the requester VPC.
* `requester-vpc-info.vpc-id` - The ID of the requester VPC.
* `status-code` - The state of the VPC peering connection (pending-acceptance | rejected | expired | active | deleted)
* `status-message` - Information about the VPC peering connection status-code.
* `vpc-peering-connection-id` - The ID of the VPC peering connection.
* `tag` - The key/value combination of a tag associated with the resource.
* `tag-key` - The key of a tag associated with the resource.
* `tag-value` - The value of a tag associated with the resource.

## Attributes Reference

The following attributes are exported:

* `status` - The state of the VPC peering connection.
* `accepter_vpc_info` - Information about the peer VPC of the VPC peering connection.
* `requester_vpc_info` - Information about the source VPC of the VPC peering connection.
* `tag_set.N` - One or more tags associated with the VPC peering connection.
* `vpc_peering_connection_id` - The ID of the VPC peering connection.

[See detailed description](http://docs.outscale.com/api_fcu/operations/Action_DescribeVpcPeeringConnections_get.html#_api_fcu-action_describevpcpeeringconnections_get)
