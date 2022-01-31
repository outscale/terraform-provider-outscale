---
layout: "outscale"
page_title: "OUTSCALE: outscale_snapshot"
sidebar_current: "outscale-snapshot"
description: |-
  [Provides information about a specific snapshot.]
---

# outscale_snapshot Data Source

Provides information about a specific snapshot.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Snapshots.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-snapshot).

## Example Usage

```hcl
data "outscale_snapshot" "snapshot01" {
  filter {
    name   = "snapshot_ids"
    values = ["snap-12345678"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `account_aliases` - (Optional) The account aliases of the owners of the snapshots.
    * `account_ids` - (Optional) The account IDs of the owners of the snapshots.
    * `descriptions` - (Optional) The descriptions of the snapshots.
    * `permissions_to_create_volume_account_ids` - (Optional) The account IDs of one or more users who have permissions to create volumes.
    * `permissions_to_create_volume_global_permission` - (Optional) If true, lists all public volumes. If false, lists all private volumes.
    * `progresses` - (Optional) The progresses of the snapshots, as a percentage.
    * `snapshot_ids` - (Optional) The IDs of the snapshots.
    * `states` - (Optional) The states of the snapshots (`in-queue` \| `completed` \| `error`).
    * `tag_keys` - (Optional) The keys of the tags associated with the snapshots.
    * `tag_values` - (Optional) The values of the tags associated with the snapshots.
    * `tags` - (Optional) The key/value combination of the tags associated with the snapshots, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.
    * `volume_ids` - (Optional) The IDs of the volumes used to create the snapshots.
    * `volume_sizes` - (Optional) The sizes of the volumes used to create the snapshots, in gibibytes (GiB).

## Attribute Reference

The following attributes are exported:

* `account_alias` - The account alias of the owner of the snapshot.
* `account_id` - The account ID of the owner of the snapshot.
* `description` - The description of the snapshot.
* `permissions_to_create_volume` - Information about the users who have permissions for the resource.
    * `account_ids` - The account ID of one or more users who have permissions for the resource.
    * `global_permission` - If true, the resource is public. If false, the resource is private.
* `progress` - The progress of the snapshot, as a percentage.
* `snapshot_id` - The ID of the snapshot.
* `state` - The state of the snapshot (`in-queue` \| `completed` \| `error`).
* `tags` - One or more tags associated with the snapshot.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
* `volume_id` - The ID of the volume used to create the snapshot.
* `volume_size` - The size of the volume used to create the snapshot, in gibibytes (GiB).
