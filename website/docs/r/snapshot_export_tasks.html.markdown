---
layout: "outscale"
page_title: "OUTSCALE: outscale_snapshot_export_tasks"
sidebar_current: "docs-outscale-snapshot-export-tasks"
description: |-
    Exports a snapshot to an Object Storage Unit (OSU) bucket.

---

# outscale_subnet

Exports a snapshot to an Object Storage Unit (OSU) bucket.

## Example Usage

```hcl

resource "outscale_volume" "test" {
  availability_zone = "eu-west-2a"
  size = 1
}

resource "outscale_snapshot" "test" {
  volume_id = "${outscale_volume.test.id}"
}

resource "outscale_snapshot_export_tasks" "outscale_snapshot_export_tasks" {
  snapshot_id = "${outscale_snapshot.test.id}"
  disk_image_format = "raw"
  osu_bucket = "customer_tooling_%d"
}


```

## Argument Reference

The following arguments are supported:

* `disk_image_format` (Required) - The format of the export disk (qcow2 | vdi | vmdk).
* `osu_bucket` (Required) - The name of the OSU bucket you want to export the snapshot to.
* `osu_key` (Optional) - The key of the OSU object corresponding to the snapshot. This element only appears in results.
* `osu_prefix` (Optional) - The prefix for the key of the OSU object corresponding to the snapshot. This key follows the prefix + snapshot_export_task_id + '.' + disk_image_format format.
* `export_to_osu_aksk` (Optional) - The access key and secret key of the OSU account used to access the bucket.
  * `access_key` - The access key of the OSU account that enables you to access the bucket.
  * `secret_key` - The secret key of the OSU account that enables you to access the bucket.
* `snapshot_id` (Required) - The ID of the snapshot to export.

## Attributes Reference

* `completion` - The percentage of the task completed.
* * `osu_key` - The key of the OSU object corresponding to the snapshot. This element only appears in results.
* `osu_prefix` -  The prefix for the key of the OSU object corresponding to the snapshot. This key follows the prefix + snapshot_export_task_id + '.' + disk_image_format format.
* `snapshot_export` - Information about the snapshot you want to export.
* `snapshot_export_task_id` - The ID of the snapshot export task.
* `snapshot_id` - The ID of the snapshot to be exported.
* `state` - The state of the snapshot export task (pending | active | completed | failed).
* `status_message` - If the snapshot export task fails, an error message.
* `request_id` - The ID of the request.

See detailed information in [Outscale Snapshot Export Tasks](http://docs.outscale.com/api_fcu/operations/Action_DescribeSubnets_get.html#_api_fcu-action_describesubnets_get).