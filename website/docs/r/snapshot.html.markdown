---
layout: "outscale"
page_title: "OUTSCALE: outscale_snapshot"
sidebar_current: "docs-outscale-resource-snapshot"
description: |-
	Creates a snapshot of a BSU volume.
---

# outscale_snapshot

Creates a snapshot of a BSU volume.
Snapshots are point-in-time images of a volume you can use to back up your data or to create replicas of this volume at the time the snapshot was created.

## Example Usage

```hcl
resource "outscale_volume" "test" {
	availability_zone = "eu-west-2a"
	size = 1
}

resource "outscale_snapshot" "test" {
	volume_id = "${outscale_volume.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `volume_id` - The ID of the BSU volume you want to create a snapshot of.
* `description` - A description for the new snapshot.

## Attributes

* `snapshot_id` - The ID of the newly created snapshot.
* `status` - The state of the snapshot (in-queue| pending | completed).
* `owner_alias` - The account alias of the owner of the snapshot.
* `description` - The description of the snapshot.
* `tag_set.N` - One or more tags associated with the snapshot.
* `volume_id` - The ID of the volume used to create the snapshot.
* `status_message` - The error message in case of snapshot copy operation failure.
* `owner_id` - The account ID of owner of the snapshot.
* `volume_size` - The size of the volume used to create the snapshot, in Gibibytes (GiB).
* `request_id` - The ID of the request.

[See detailed information.](http://docs.outscale.com/api_fcu/operations/Action_CreateSnapshot_get.html#_api_fcu-action_createsnapshot_get)
