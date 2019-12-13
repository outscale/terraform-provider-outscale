---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_net_peering_acceptation"
sidebar_current: "outscale-net-peering-acceptation"
description: |-
  [Manages a Net peering acceptation.]
---

# outscale_net_peering_acceptation Resource

Manages a Net peering acceptation.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+VPC+Peering+Connections).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#acceptnetpeering).

## Example Usage

```hcl

resource "outscale_net_peering_acceptation" "outscale_net_peering_acceptation01" {
  net_peering_id = outscale_net_peering.outscale_net_peering01.net_peering_id
}


```

## Argument Reference

The following arguments are supported:

* `net_peering_id` - (Required) The ID of the Net peering connection you want to accept.

## Attribute Reference

The following attributes are exported:

* `net_peering` - Information about the Net peering connection.
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
