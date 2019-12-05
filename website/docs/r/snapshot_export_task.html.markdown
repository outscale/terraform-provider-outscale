---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_snapshot_export_task"
sidebar_current: "docs-outscale-resource-snapshot-export-task"
description: |-
  [Manages a snapshot export task.]
---

# outscale_snapshot_export_task Resource

Manages a snapshot export task.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Snapshots).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-snapshot).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `osu_export` - Information about the OSU export.
  * `disk_image_format` - (Optional) The format of the export disk (`qcow2` \| `vdi` \| `vmdk`).
  * `osu_api_key` - Information about the OSU API key.
    * `api_key_id` - (Optional) The API key of the OSU account that enables you to access the bucket.
    * `secret_key` - (Optional) The secret key of the OSU account that enables you to access the bucket.
  * `osu_bucket` - (Optional) The name of the OSU bucket you want to export the object to.
  * `osu_manifest_url` - (Optional) The URL of the manifest file.
  * `osu_prefix` - (Optional) The prefix for the key of the OSU object. This key follows this format: `prefix + object_export_task_id + '.' + disk_image_format`.
* `snapshot_id` - (Required) The ID of the snapshot to export.

## Attribute Reference

The following attributes are exported:

* `snapshot_export_task` - Information about the snapshot export task.
  * `comment` - If the snapshot export task fails, an error message appears.
  * `osu_export` - Information about the OSU export.
    * `disk_image_format` - The format of the export disk (`qcow2` \| `vdi` \| `vmdk`).
    * `osu_api_key` - Information about the OSU API key.
      * `api_key_id` - The API key of the OSU account that enables you to access the bucket.
      * `secret_key` - The secret key of the OSU account that enables you to access the bucket.
    * `osu_bucket` - The name of the OSU bucket you want to export the object to.
    * `osu_manifest_url` - The URL of the manifest file.
    * `osu_prefix` - The prefix for the key of the OSU object. This key follows this format: `prefix + object_export_task_id + '.' + disk_image_format`.
  * `progress` - The progress of the snapshot export task, as a percentage.
  * `snapshot_id` - The ID of the snapshot to be exported.
  * `state` - The state of the snapshot export task (`pending` \| `active` \| `completed` \| `failed`).
  * `tags` - One or more tags associated with the snapshot export task.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
  * `task_id` - The ID of the snapshot export task.
