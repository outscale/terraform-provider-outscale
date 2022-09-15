---
layout: "outscale"
page_title: "OUTSCALE: outscale_nets"
sidebar_current: "outscale-nets"
description: |-
  [Provides information about Nets.]
---

# outscale_nets Data Source

Provides information about Nets.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-VPCs.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-net).

## Example Usage

```hcl
data "outscale_nets" "nets01" {
  filter {
    name   = "net_ids"
    values = ["vpc-12345678", "vpc-12345679", "vpc-12345680"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `dhcp_options_set_ids` - (Optional) The IDs of the DHCP options sets.
    * `ip_ranges` - (Optional) The IP ranges for the Nets, in CIDR notation (for example, 10.0.0.0/16).
    * `net_ids` - (Optional) The IDs of the Nets.
    * `states` - (Optional) The states of the Nets (`pending` \| `available`).
    * `tag_keys` - (Optional) The keys of the tags associated with the Nets.
    * `tag_values` - (Optional) The values of the tags associated with the Nets.
    * `tags` - (Optional) The key/value combination of the tags associated with the Nets, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.

## Attribute Reference

The following attributes are exported:

* `nets` - Information about the described Nets.
    * `dhcp_options_set_id` - The ID of the DHCP options set (or `default` if you want to associate the default one).
    * `ip_range` - The IP range for the Net, in CIDR notation (for example, 10.0.0.0/16).
    * `net_id` - The ID of the Net.
    * `state` - The state of the Net (`pending` \| `available`).
    * `tags` - One or more tags associated with the Net.
        * `key` - The key of the tag, with a minimum of 1 character.
        * `value` - The value of the tag, between 0 and 255 characters.
    * `tenancy` - The VM tenancy in a Net.
