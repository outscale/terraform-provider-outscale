---
layout: "outscale"
page_title: "OUTSCALE: outscale_subnet"
sidebar_current: "docs-outscale-resource-subnet"
description: |-
To create a subnet in a VPC, you have to provide the ID of the VPC and the CIDR block for the subnet (its network range). Once the subnet is created, you cannot modify its CIDR block.

---

# outscale_subnet

NOTE: The CIDR block of the subnet can be either the same as the VPC one if you create only a single subnet in this VPC, or a subset of the VPC one. In case of several subnets in a VPC, their CIDR blocks must no overlap. The smallest subnet you can create uses a /30 netmask (four IP addresses).

## Example Usage

```hcl

resource "outscale_subnet" "basic" {
    cidr_block = "10.0.0.0/16"
    vpc_id = "vpc-2f09a348"
    availability_zone = "eu-west-1"
}


```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional) The name of the Availability Zone in which you want to create the subnet.
* `cidr_block` - (Required) The CIDR block for the subnet (for example, 10.0.0.0/24).
* `vpc_id` - (Required) The ID of the VPC..

.

## Attributes Reference

* `availability_zone` - (Optional) The name of the Availability Zone in which the subnet is located.	.
* `available_ip_address_count` - (Optional) The number of unused IP addresses in the subnet	.
* `cidr_block` - (Optional) The CIDR block of the subnet (for example, 10.84.7.0/24).
* `state` - (Optional) The state of the subnet (pending | available).
* `subnet_id` - (Optional) The ID of the subnet.
* `tag_set.N` - (Optional) One or more tags associated with the VPC.
* `vpc_id` - (Optional) The ID of the VPC where the subnet is.
* `request_id`- (Optional) The ID of the request.

See detailed information in [FCU DescribeSubNet](http://docs.outscale.com/api_fcu/operations/Action_DescribeSubnets_get.html#_api_fcu-action_describesubnets_get).
See detailed information in [FCU CreateSubNet](http://docs.outscale.com/api_fcu/operations/Action_CreateSubnet_get.html#_api_fcu-action_createsubnet_get).