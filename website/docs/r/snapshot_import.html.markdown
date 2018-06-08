---
layout: "outscale"
page_title: "OUTSCALE: outscale_snapshot_import"
sidebar_current: "docs-outscale-resource-snapshot-import"
description: |-
  Imports a snapshot from an Object Storage Unit (OSU) bucket to create a copy of this snapshot in your account.
---

# outscale_snapshot_import

Imports a snapshot from an Object Storage Unit (OSU) bucket to create a copy of this snapshot in your account.
This method enables you to copy a snapshot from another account that is either within the same Region as the OSU bucket, or in a different one. To copy a snapshot within the same Region, you can also use the CopySnapshot direct method. For more information, see [CopySnapshot](http://docs.outscale.com/api_fcu/operations/Action_CopySnapshot_post.html#_api_fcu-action_copysnapshot_post).
The copy of the source snapshot is independent and belongs to you.
You can import a snapshot using a pre-signed URL. You do not need any permission for this snapshot, or the bucket in which it is contained. The pre-signed URL is valid for seven days (you can re-generate a new one if needed). For more information about how to export a snapshot an OSU bucket, see [CreateSnapshotExportTask](http://docs.outscale.com/api_fcu/operations/Action_CreateSnapshotExportTask_post.html#_api_fcu-action_createsnapshotexporttask_post).

## Example Usage

```hcl
resource "outscale_snapshot_import" "test" {
    snapshot_location = ""
    snapshot_size = ""
}
```

## Argument Reference

The following arguments are supported:

* `snapshot_location` - The pre-signed URL of the snapshot you want to import from the OSU bucket.
* `snapshot_size` - The size of the snapshot created in your account, in Gibibytes (GiB). This size must be exactly the same as the source snapshot one.
* `description` - The description for the snapshot created in your account.

## Attributes Reference

* `description` - The description of the snapshot created in your account.
* `encrypted` - Indicates whether the snapshot is encrypted or not (always false).
* `owner_alias` - The alias of the owner of the snapshot created in your account.
* `progress` - The percentage of the task completed.
* `status` - The state of the snapshot created in your account (error | completed).
* `volume_size` - The ID of the new snapshot.
* `id` - The size of the snapshot created in your account, in Gibibytes (GiB).

[See detailed description](http://docs.outscale.com/api_fcu/operations/Action_ImportSnapshot_get.html#_api_fcu-action_importsnapshot_get)
