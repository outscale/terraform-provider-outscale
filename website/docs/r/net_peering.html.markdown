---
layout: "outscale"
page_title: "OUTSCALE: outscale_net_peering"
sidebar_current: "outscale-net-peering"
description: |-
  [Manages a Net peering.]
---

# outscale_net_peering Resource

Manages a Net peering.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-VPC-Peering-Connections.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-netpeering).

## Example Usage

### Required resources

```hcl
resource "outscale_net" "net01" {
  ip_range = "10.10.0.0/24"
}

resource "outscale_net" "net02" {
  ip_range = "10.31.0.0/16"
}
```

### Peer Nets

```hcl
resource "outscale_net_peering" "net_peering01" {
  accepter_net_id = outscale_net.net01.net_id
  source_net_id   = outscale_net.net02.net_id
}
```

## Argument Reference

The following arguments are supported:

* `accepter_net_id` - (Required) The ID of the Net you want to connect with.
* `source_net_id` - (Required) The ID of the Net you send the peering request from.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

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

## Import

A Net peering can be imported using its ID. For example:

```console

$ terraform import outscale_net_peering.ImportedNetPeering pcx-12345678

```