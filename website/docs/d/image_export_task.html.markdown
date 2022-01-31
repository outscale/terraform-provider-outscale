---
layout: "outscale"
page_title: "OUTSCALE: outscale_image_export_task"
sidebar_current: "outscale-image-export-task"
description: |-
  [Provides information about a specific image export task.]
---

# outscale_image_export_task Data Source

Provides information about a specific image export task.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-OMIs.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-image).

## Example Usage

```hcl
data "outscale_image_export_task" "image_export_task01" {
  filter {
    name   = "task_ids"
    values = ["image-export-12345678"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `task_ids` - (Optional) The IDs of the export tasks.

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
