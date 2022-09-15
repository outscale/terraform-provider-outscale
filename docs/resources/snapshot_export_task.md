---
layout: "outscale"
page_title: "OUTSCALE: outscale_snapshot_export_task"
sidebar_current: "outscale-snapshot-export-task"
description: |-
  [Manages a snapshot export task.]
---

# outscale_snapshot_export_task Resource

Manages a snapshot export task.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Snapshots.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-snapshot).

## Example Usage

### Required resources

```hcl
resource "outscale_volume" "volume01" {
	subregion_name = "${var.region}a"
	size           = 40
}

resource "outscale_snapshot" "snapshot01" {
	volume_id = outscale_volume.volume01.volume_id
}
```

### Create a snapshot export task

```hcl
resource "outscale_snapshot_export_task" "snapshot_export_task01" {
	snapshot_id = outscale_snapshot.snapshot01.snapshot_id
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

* `osu_export` - Information about the OOS export task to create.
    * `disk_image_format` - (Optional) The format of the export disk (`qcow2` \| `raw`).
    * `osu_api_key` - Information about the OOS API key.
        * `api_key_id` - (Optional) The API key of the OOS account that enables you to access the bucket.
        * `secret_key` - (Optional) The secret key of the OOS account that enables you to access the bucket.
    * `osu_bucket` - (Optional) The name of the OOS bucket where you want to export the object.
    * `osu_manifest_url` - (Optional) The URL of the manifest file.
    * `osu_prefix` - (Optional) The prefix for the key of the OOS object.
* `snapshot_id` - (Required) The ID of the snapshot to export.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `comment` - If the snapshot export task fails, an error message appears.
* `osu_export` - Information about the snapshot export task.
    * `disk_image_format` - The format of the export disk (`qcow2` \| `raw`).
    * `osu_bucket` - The name of the OOS bucket the snapshot is exported to.
    * `osu_prefix` - The prefix for the key of the OOS object corresponding to the snapshot.
* `progress` - The progress of the snapshot export task, as a percentage.
* `snapshot_id` - The ID of the snapshot to be exported.
* `state` - The state of the snapshot export task (`pending` \| `active` \| `completed` \| `failed`).
* `tags` - One or more tags associated with the snapshot export task.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
* `task_id` - The ID of the snapshot export task.

