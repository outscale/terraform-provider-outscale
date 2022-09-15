---
layout: "outscale"
page_title: "OUTSCALE: outscale_snapshot_export_task"
sidebar_current: "outscale-snapshot-export-task"
description: |-
  [Provides information about a specific snapshot export task.]
---

# outscale_snapshot_export_task Data Source

Provides information about a specific snapshot export task.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Snapshots.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-snapshot).

## Example Usage

```hcl
data "outscale_snapshot_export_task" "snapshot_export_task01" {
  filter {
    name   = "task_ids"
    values = ["snap-export-12345678"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `task_ids` - (Optional) The IDs of the export tasks.

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
