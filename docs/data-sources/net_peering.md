---
layout: "outscale"
page_title: "OUTSCALE: outscale_net_peering"
sidebar_current: "outscale-net-peering"
description: |-
  [Provides information about a specific Net peering.]
---

# outscale_net_peering Data Source

Provides information about a specific Net peering.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-VPC-Peering-Connections.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-netpeering).

## Example Usage

```hcl
data "outscale_net_peering" "net_peering01" {
  filter {
    name   = "net_peering_ids"
    values = ["pcx-12345678"]
  }    
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `accepter_net_account_ids` - (Optional) The account IDs of the owners of the peer Nets.
    * `accepter_net_ip_ranges` - (Optional) The IP ranges of the peer Nets, in CIDR notation (for example, 10.0.0.0/24).
    * `accepter_net_net_ids` - (Optional) The IDs of the peer Nets.
    * `net_peering_ids` - (Optional) The IDs of the Net peering connections.
    * `source_net_account_ids` - (Optional) The account IDs of the owners of the peer Nets.
    * `source_net_ip_ranges` - (Optional) The IP ranges of the peer Nets.
    * `source_net_net_ids` - (Optional) The IDs of the peer Nets.
    * `state_messages` - (Optional) Additional information about the states of the Net peering connections.
    * `state_names` - (Optional) The states of the Net peering connections (`pending-acceptance` \| `active` \| `rejected` \| `failed` \| `expired` \| `deleted`).
    * `tag_keys` - (Optional) The keys of the tags associated with the Net peering connections.
    * `tag_values` - (Optional) The values of the tags associated with the Net peering connections.
    * `tags` - (Optional) The key/value combination of the tags associated with the Net peering connections, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.

## Attribute Reference

The following attributes are exported:

* `accepter_net` - Information about the accepter Net.
    * `account_id` - The account ID of the owner of the accepter Net.
    * `ip_range` - The IP range for the accepter Net, in CIDR notation (for example, 10.0.0.0/16).
    * `net_id` - The ID of the accepter Net.
* `net_peering_id` - The ID of the Net peering connection.
* `source_net` - Information about the source Net.
    * `account_id` - The account ID of the owner of the source Net.
    * `ip_range` - The IP range for the source Net, in CIDR notation (for example, 10.0.0.0/16).
    * `net_id` - The ID of the source Net.
* `state` - Information about the state of the Net peering connection.
    * `message` - Additional information about the state of the Net peering connection.
    * `name` - The state of the Net peering connection (`pending-acceptance` \| `active` \| `rejected` \| `failed` \| `expired` \| `deleted`).
* `tags` - One or more tags associated with the Net peering connection.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
