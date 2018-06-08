---
layout: "outscale"
page_title: "OUTSCALE: outscale_snapshot_copy"
sidebar_current: "docs-outscale-resource-snapshot-copy"
description: |-
  Copies a snapshot to your account, from an account in the same Region.
---

# outscale_snapshot_copy

Copies a snapshot to your account, from an account in the same Region.
To do so, the owner of the source snapshot must share it with your account. For more information about how to share a snapshot with another account, see [ModifySnapshotAttribute](http://docs.outscale.com/api_fcu/operations/Action_ModifySnapshotAttribute_post.html#_api_fcu_modifysnapshotattribute_post).
The copy of the source snapshot is independent and belongs to you.
To copy a snapshot between accounts in different Regions, the owner of the source snapshot must export it to an Object Storage Unit (OSU) bucket using the CreateSnapshotExportTask method, and then you need to import it using the ImportSnapshot method. For more information, see [CreateSnapshotExportTask](http://docs.outscale.com/api_fcu/operations/Action_CreateSnapshotExportTask_post.html#_api_fcu-action_createsnapshotexporttask_post) and [ImportSnapshot](http://docs.outscale.com/api_fcu/operations/Action_ImportSnapshot_post.html#_api_fcu_importsnapshot_post).  

## Example Usage

```hcl
resource "outscale_volume" "test" {
    availability_zone = "eu-west-2a"
    size = 1
}

resource "outscale_snapshot" "test" {
    volume_id = "${outscale_volume.test.id}"
    description = "Snapshot Acceptance Test"
}

resource "outscale_snapshot_copy" "test" {
    source_region =  "eu-west-2"
    source_snapshot_id = "${outscale_snapshot.test.id}"
}
```

## Argument Reference

The following arguments are supported:

* `source_snapshot_id` - The ID of the snapshot you want to copy.
* `source_region` - The name of the destination Region.
* `description` - A description for the new snapshot (if different from the source snapshot one).
* `destination_region` - The name of the destination Region.

## Attributes Reference

* `snapshot_id` - The ID of the new snapshot.

[See detailed description](http://docs.outscale.com/api_fcu/operations/Action_CopySnapshot_get.html#_api_fcu-action_copysnapshot_get)
