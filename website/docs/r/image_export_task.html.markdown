---
layout: "outscale"
page_title: "OUTSCALE: outscale_image_export_task"
sidebar_current: "outscale-image-export-task"
description: |-
  [Manages an image export task.]
---

# outscale_image_export_task Resource

Manages an image export task.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-OMIs.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-image).

## Example Usage

### Required resource

```hcl
resource "outscale_image" "image01" {
  image_name = "terraform-image-to-export"
  vm_id      = "i-12345678"
}
```

### Create an image export task

```hcl
resource "outscale_image_export_task" "image_export_task01" {
	image_id = outscale_image.image01.image_id
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

* `image_id` - (Required) The ID of the OMI to export.
* `osu_export` - Information about the OOS export task to create.
    * `disk_image_format` - (Optional) The format of the export disk (`qcow2` \| `raw`).
    * `osu_api_key` - Information about the OOS API key.
        * `api_key_id` - (Optional) The API key of the OOS account that enables you to access the bucket.
        * `secret_key` - (Optional) The secret key of the OOS account that enables you to access the bucket.
    * `osu_bucket` - (Optional) The name of the OOS bucket where you want to export the object.
    * `osu_manifest_url` - (Optional) The URL of the manifest file.
    * `osu_prefix` - (Optional) The prefix for the key of the OOS object.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.

## Attribute Reference

The following attributes are exported:

* `comment` - If the OMI export task fails, an error message appears.
* `image_id` - The ID of the OMI to be exported.
* `osu_export` - Information about the OMI export task.
    * `disk_image_format` - The format of the export disk (`qcow2` \| `raw`).
    * `osu_bucket` - The name of the OOS bucket the OMI is exported to.
    * `osu_manifest_url` - The URL of the manifest file.
    * `osu_prefix` - The prefix for the key of the OOS object corresponding to the image.
* `progress` - The progress of the OMI export task, as a percentage.
* `state` - The state of the OMI export task (`pending/queued` \| `pending` \| `completed` \| `failed` \| `cancelled`).
* `tags` - One or more tags associated with the image export task.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
* `task_id` - The ID of the OMI export task.

