---
layout: "outscale"
page_title: "OUTSCALE: outscale_snapshots"
sidebar_current: "docs-outscale-datasource-outscale-snapshots"
description: |-
  Describes one or more snapshots that are available to you.
---

# outscale_snapshots

Describes one or more snapshots that are available to you.

## Example Usage

```hcl
resource "outscale_volume" "example" {
    availability_zone = "eu-west-2a"
    volume_type = "gp2"
    size = 40
    tag {
        Name = "External Volume"
    }
}

resource "outscale_snapshot" "snapshot" {
    volume_id = "${outscale_volume.example.id}"
}

data "outscale_snapshots" "outscale_snapshots" {
    snapshot_id = ["${outscale_snapshot.snapshot.id}"]
}
```

## Argument Reference

The following arguments are supported:

* `snapshot_id.N` - One or more snapshot IDs.
* `filter.N` - One or more filters
* `owner.N` - One or more owners of the snapshots.
* `restorable_by.N` - One or more accounts IDs that have the permissions to create volumes from the snapshot.

## Filters

You can filter the described snapshots using the snapshot_id.N, the owner.N and the restorable_by.N parameters.
You can also use the Filter.N parameter to filter the snapshots on the following properties:

* `description` - The description of the snapshot.
* `owner-alias` - The account alias of the owner of the snapshot
* `owner-id` - The account ID of the owner of the snapshot.
* `progress` - The progress of the snapshot, as a percentage.
* `snapshot-id` - The ID of the snapshot.
* `start-time` - The time at which the snapshot was initiated.
* `status` - The state of the snapshot (in-queue | pending | completed).
* `volume-id` - The ID of the volume used to create the snapshot.
* `volume-size` - The size of the volume used to create the snapshot, in Gibibytes (GiB).
* `tag` - The key/value combination of a tag associated with the resource.
* `tag-key` - The key of a tag associated with the resource.
* `tag-value` - The value of a tag associated with the resource.

## Attributes Reference

The following attributes are exported:

* `snapshot_set.N` - Information about one or more snapshots; each described with the following attributes:
  * `progress` - The progress of the snapshot, as a percentage.
  * `status` - The state of the snapshot (in-queue| pending | completed).
  * `owner_alias` - The account alias of the owner of the snapshot.
  * `description` - The description of the snapshot.
  * `tag_set.N` - One or more tags associated with the snapshot.
  * `volume_id` - The ID of the volume used to create the snapshot.
  * `status_message` - The error message in case of snapshot copy operation failure.
   `owner_id` - The account ID of owner of the snapshot.
  * `volume_size` - The size of the volume used to create the snapshot, in Gibibytes (GiB).
  * `start_time` - The time at which the snapshot was initiated.

[See detailed description](http://docs.outscale.com/api_fcu/operations/Action_DescribeSnapshots_get.html#_api_fcu-action_describesnapshots_get)
