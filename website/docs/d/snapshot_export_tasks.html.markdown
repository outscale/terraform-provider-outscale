---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_snapshot_export_tasks"
sidebar_current: "outscale-snapshot-export-tasks"
description: |-
  [Provides information about snapshot export tasks.]
---

# outscale_snapshot_export_tasks Data Source

Provides information about snapshot export tasks.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Snapshots).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#createsnapshotexporttask).

## Example Usage

```hcl

data "outscale_snapshot_export_tasks" "snapshot_export_tasks01" {
  filter {
    name   = "task_ids"
    values = ["snap-export-12345678", "snap-export-12345679"]
  }
}


```

## Argument Reference

The following arguments are supported:

* `filter` - One or more filters.
    * `task_ids` - (Optional) The IDs of the export tasks.

## Attribute Reference

The following attributes are exported:

* `snapshot_export_tasks` - Information about one or more snapshot export tasks.
    * `comment` - If the snapshot export task fails, an error message appears.
    * `osu_export` - Information about the OSU export.
        * `disk_image_format` - The format of the export disk (`qcow2` \| `raw`).
        * `osu_api_key` - Information about the OSU API key.
            * `api_key_id` - The API key of the OSU account that enables you to access the bucket.
            * `secret_key` - The secret key of the OSU account that enables you to access the bucket.
        * `osu_bucket` - The name of the OSU bucket where you want to export the object.
        * `osu_manifest_url` - The URL of the manifest file.
        * `osu_prefix` - The prefix for the key of the OSU object.
    * `progress` - The progress of the snapshot export task, as a percentage.
    * `snapshot_id` - The ID of the snapshot to be exported.
    * `state` - The state of the snapshot export task (`pending` \| `active` \| `completed` \| `failed`).
    * `tags` - One or more tags associated with the snapshot export task.
        * `key` - The key of the tag, with a minimum of 1 character.
        * `value` - The value of the tag, between 0 and 255 characters.
    * `task_id` - The ID of the snapshot export task.
