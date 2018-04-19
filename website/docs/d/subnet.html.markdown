---
layout: "outscale"
page_title: "OUTSCALE: outcale_subnet"
sidebar_current: "docs-outscale-datasource-subnet"
description: |-
    Describes one or more subnet

---

# outscale_subnet

Describes one or more of your subnets.
If you do not specify any subnet ID, this action describes all of your subnets.
You can use the Filter.N parameter to filter the subnets on the following properties:

## Example Usage

```hcl
data "outscale_subnet" "by_id" {
  subnet_id = "${outscale_subnet.test.id}"
}

data "outscale_subnet" "by_filter" {
  filter {
    name = "subnet-id"
    values = ["${outscale_subnet.test.id}"]
  }
```

## Argument Reference

The following arguments are supported:

* `subnet_id.N` - (Optional) (Only is provided here).
* `filter.N` - (Optional) One or more filters.





See detailed information in [Outscale subnet](http://docs.outscale.com/api_fcu/operations/Action_DescribeSubnets_get.html#_api_fcu-action_describesubnets_get).

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `subnet-id` - (Optional) The ID of the Subnet.
* `vpc-id` - (Optional) The ID of the VPC in which the subnet is.
* `state` - (Optional) The state of the subnet (pending | available).
* `cidr` - (Optional) The exact CIDR block of the subnet.
* `cidr-block` - (Optional) The exact CIDR block of the subnet (similar to cidr and cidrBlock).
* `cidrBlock` - (Optional) The exact CIDR block of the VPC (similar to cidr and cidr-block).
* `available-ip-address-count` - (Optional) The number of available IP adresses in the subnet.
* `availability-zone` - (Optional) The Availability Zone in which the subnets are located.
* `availabilityZone` - (Optional) Alias for availability-zone.



## Attributes Reference

The following attributes are exported:

* `subnet_id` - (Optional) The ID of the Subnet.
* `vpc_id` - (Optional) The ID of the VPC in which the subnet is.
* `state` - (Optional) The state of the subnet (pending | available).
* `cidr_block` - (Optional) The exact CIDR block of the subnet (similar to cidr and cidrBlock).
* `available_ip_address_count` - (Optional) The number of available IP adresses in the subnet.
* `availability_zone` - (Optional) The Availability Zone in which the subnets are located.
* `tag_set.N` - (Optional) One or more tags associated with the VPC.
* `requester_id` - (Optional) The iD of the request.




See detailed information in [Describe subnet](http://docs.outscale.com/api_fcu/operations/Action_DescribeSubnets_get.html#_api_fcu-action_describesubnets_get).
