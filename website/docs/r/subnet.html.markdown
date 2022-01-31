---
layout: "outscale"
page_title: "OUTSCALE: outscale_subnet"
sidebar_current: "outscale-subnet"
description: |-
  [Manages a Subnet.]
---

# outscale_subnet Resource

Manages a Subnet.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-VPCs.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-subnet).

## Example Usage

### Required resource

```hcl
resource "outscale_net" "net01" {
	ip_range = "10.0.0.0/16"
}
```

### Create a subnet

```hcl
resource "outscale_subnet" "subnet01" {
	net_id   = outscale_net.net01.net_id
	ip_range = "10.0.0.0/18"
}
```

## Argument Reference

The following arguments are supported:

* `ip_range` - (Required) The IP range in the Subnet, in CIDR notation (for example, 10.0.0.0/16).
* `net_id` - (Required) The ID of the Net for which you want to create a Subnet.
* `subregion_name` - (Optional) The name of the Subregion in which you want to create the Subnet.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `available_ips_count` - The number of available IP addresses in the Subnets.
* `ip_range` - The IP range in the Subnet, in CIDR notation (for example, 10.0.0.0/16).
* `map_public_ip_on_launch` - If true, a public IP is assigned to the network interface cards (NICs) created in the specified Subnet.
* `net_id` - The ID of the Net in which the Subnet is.
* `state` - The state of the Subnet (`pending` \| `available`).
* `subnet_id` - The ID of the Subnet.
* `subregion_name` - The name of the Subregion in which the Subnet is located.
* `tags` - One or more tags associated with the Subnet.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Import

A subnet can be imported using its ID. For example:

```console

$ terraform import outscale_subnet.ImportedSubnet subnet-12345678

```