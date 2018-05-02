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

resource "outscale_snapshot_export_tasks" "basic" {
    Snapshot_id = ""
}


```

## Argument Reference

The following arguments are supported:

* `Export_to_osu` - (Optional)	Information about the export (you must at least specify DiskImageFormat and OsuBucket parameters).
* `Snapshot_id` - (Required)	The ID of the snapshot to export.

## Attributes Reference

* `completion` -	The percentage of the task completed.
* `export_to_osu` -	Information to export the snapshot to OSU.
* `snapshot_export` -	Information about the snapshot you want to export.
* `snapshot_export_task_id` -	The ID of the snapshot export task.
* `snapshot_id` -	The ID of the snapshot to be exported.
* `state` -	The state of the snapshot export task (pending | active | completed | failed).
* `status_message` -	If the snapshot export task fails, an error message.
* `request_id` -	The ID of the request.


See detailed information in [Outscale Snapshot Export Tasks](http://docs.outscale.com/api_fcu/operations/Action_DescribeSubnets_get.html#_api_fcu-action_describesubnets_get).