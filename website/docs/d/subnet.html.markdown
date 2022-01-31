---
layout: "outscale"
page_title: "OUTSCALE: outscale_subnet"
sidebar_current: "outscale-subnet"
description: |-
  [Provides information about a specific Subnet.]
---

# outscale_subnet Data Source

Provides information about a specific Subnet.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-VPCs.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-subnet).

## Example Usage

```hcl
data "outscale_subnet" "subnet01" {
  filter {
    name   = "net_ids"
    values = ["vpc-12345678"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `available_ips_counts` - (Optional) The number of available IPs.
    * `ip_ranges` - (Optional) The IP ranges in the Subnets, in CIDR notation (for example, 10.0.0.0/16).
    * `net_ids` - (Optional) The IDs of the Nets in which the Subnets are.
    * `states` - (Optional) The states of the Subnets (`pending` \| `available`).
    * `subnet_ids` - (Optional) The IDs of the Subnets.
    * `subregion_names` - (Optional) The names of the Subregions in which the Subnets are located.
    * `tag_keys` - (Optional) The keys of the tags associated with the Subnets.
    * `tag_values` - (Optional) The values of the tags associated with the Subnets.
    * `tags` - (Optional) The key/value combination of the tags associated with the Subnets, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.

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
