---
layout: "outscale"
page_title: "OUTSCALE: outscale_subnet"
sidebar_current: "docs-outscale-resource-subnet"
description: |-
    Creates a Subnet in an existing Net.
---

# outscale_subnet

To create a Subnet in a Net, you have to provide the ID of the Net and the IP range for the Subnet (its network range). Once the Subnet is created, you cannot modify its IP range.

The IP range of the Subnet can be either the same as the Net one if you create only a single Subnet in this Net, or a subset of the Net one. In case of several Subnets in a Net, their IP ranges must not overlap. The smallest Subnet you can create uses a /30 netmask (four IP addresses).

## Example Usage

```hcl
resource "outscale_net" "net" {
  ip_range = "10.0.0.0/16"
}

resource "outscale_subnet" "subnet" {
  ip_range = "10.0.0.0/16"
  subregion_name = "in-west-2a"
  net_id = "${outscale_net.net.id}"
}
```

## Argument Reference

The following arguments are supported:

* `subregion_name` - (Optional) The name of the Subregion in which you want to create the Subnet.
* `ip_range` - (Required) The IP range in the Subnet, in CIDR notation (for example, 10.0.0.0/16).
* `net_id` - (Required) The ID of the Net for which you want to create a Subnet.

.

## Attributes Reference

* `available_ips_count` - The number of available IP addresses in the Subnets.
* `ip_range` - The IP range in the Subnet, in CIDR notation (for example, 10.0.0.0/16).
* `net_id` - The ID of the Net in which the Subnet is.
* `state` - The state of the Subnet (pending | available).
* `subnet_id` - The ID of the Subnet.
* `subregion_name` - The name of the Subregion in which the Subnet is located.
* `tags` - One or more tags associated with the Subnet.
* `request_id`- The ID of the request.

More info here [Subnet](https://docs-beta.outscale.com/oapi#outscale-api-subnet).
More info here [Subnet](https://docs-beta.outscale.com/oapi#tocssubnet).
See detailed information in [OAPI ReadSubnets](http://docs.outscale.com/api_fcu/operations/Action_DescribeSubnets_get.html#_api_fcu-action_describesubnets_get).
See detailed information in [OAPI CreateSubnet](http://docs.outscale.com/api_fcu/operations/Action_CreateSubnet_get.html#_api_fcu-action_createsubnet_get).