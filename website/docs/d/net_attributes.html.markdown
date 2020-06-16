---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_net_attributes"
sidebar_current: "outscale-net-attributes"
description: |-
  [Provides information about Net attributes.]
---

# outscale_net_attributes Data Source

Provides information about Net attributes.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+VPCs).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#updatenet).

## Example Usage

```hcl

data "outscale_net_attributes" "net_attributes01" {
  net_id = "vpc-12345678"
}


```

## Argument Reference

The following arguments are supported:

* `filters` - One or more filters.
  * `dhcp_options_set_id` - (Optional) The ID of the DHCP options set.
  * `ip_range` - (Optional) The IP range for the Net, in CIDR notation (for example, 10.0.0.0/16).
  * `net_id` - (Optional) The ID of the Net.
  * `state` - (Optional) The state of the Net (`pending` | `available`).
  * `tags` - (Optional) The key/value combination of the tags associated with the security groups, in the following format: "Filters":{"Tags":["TAGKEY=TAGVALUE"]}. 

## Attribute Reference

The following attributes are exported:

* `net_attributes` - Information about one or more Net attributes.
  * `dhcp_options_set_id` - The ID of the DHCP options set (or `default` if you want to associate the default one).
  * `ip_range` - The IP range for the Net, in CIDR notation (for example 10.0.0.0/16).
  * `net_id` - The ID of the Net.
  * `state` - The state of the Net (`pending` | `available`).
  * `tags` - One or more tags associated with the Net.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
