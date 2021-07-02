---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_snapshot_export_task"
sidebar_current: "outscale-snapshot-export-task"
description: |-
  [Manages a snapshot export task.]
---

# outscale_snapshot_export_task Resource

Manages a snapshot export task.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Snapshots).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-snapshot).

## Example Usage

```hcl
resource "outscale_snapshot_export_task" "snapshot_export_task01" {
    snapshot_id = "snap-12345678"
    osu_export {
        disk_image_format = "qcow2"
        osu_bucket        = "terraform-bucket"
        osu_prefix        = "new-export"
        osu_api_key {
            api_key_id = var.access_key_id
            secret_key = var.secret_key_id
        }
    }
    tags {
        key   = "Name"
        value = "terraform-snapshot-export-task"
    }
}
```

## Argument Reference

The following arguments are supported:

* `osu_export` - Information about the OSU export.
    * `disk_image_format` - (Optional) The format of the export disk (`qcow2` \| `raw`).
    * `osu_api_key` - Information about the OSU API key.
        * `api_key_id` - (Optional) The API key of the OSU account that enables you to access the bucket.
        * `secret_key` - (Optional) The secret key of the OSU account that enables you to access the bucket.
    * `osu_bucket` - (Optional) The name of the OSU bucket where you want to export the object.
    * `osu_manifest_url` - (Optional) The URL of the manifest file.
    * `osu_prefix` - (Optional) The prefix for the key of the OSU object.
* `snapshot_id` - (Required) The ID of the snapshot to export.
* `tags` - A tag to add to this resource. You can specify this argument several times.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `snapshot_export_task` - Information about the snapshot export task.
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

