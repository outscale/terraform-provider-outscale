---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_subnet"
sidebar_current: "docs-outscale-resource-subnet"
description: |-
  [Manages a subnet.]
---

# outscale_subnet

Manages a subnet.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+VPCs).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-subnet).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `ip_range` - (Required) The IP range in the Subnet, in CIDR notation (for example, 10.0.0.0/16).
* `net_id` - (Required) The ID of the Net for which you want to create a Subnet.
* `subregion_name` - (Optional) The name of the Subregion in which you want to create the Subnet.

## Attribute Reference

The following attributes are exported:

* `subnet` - Information about the Subnet.
  * `available_ips_count` - The number of available IP addresses in the Subnets.
  * `ip_range` - The IP range in the Subnet, in CIDR notation (for example, 10.0.0.0/16).
  * `map_public_ip_on_launch` - If `true`, a public IP address is assigned to the network interface cards (NICs) created in the specified Subnet.
  * `net_id` - The ID of the Net in which the Subnet is.
  * `state` - The state of the Subnet (`pending` \| `available`).
  * `subnet_id` - The ID of the Subnet.
  * `subregion_name` - The name of the Subregion in which the Subnet is located.
  * `tags` - One or more tags associated with the Subnet.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
