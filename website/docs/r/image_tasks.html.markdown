---
layout: "outscale"
page_title: "OUTSCALE: outscale_image_tasks"
sidebar_current: "docs-outscale-resource-image-tasks"
description: |-
  Exports an Outscale machine image (OMI) to an Object Storage Unit (OSU) bucket.
---

# outscale_image_tasks

Exports an Outscale machine image (OMI) to an Object Storage Unit (OSU) bucket.. [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
resource "outscale_vm" "outscale_vm" {
    count = 1

    image_id                    = "ami-880caa66"
    instance_type               = "c4.large"

}

resource "outscale_image" "outscale_image" {
    name            = "image_${outscale_vm.outscale_vm.id}"
    instance_id     = "${outscale_vm.outscale_vm.id}"
}

resource "outscale_image_tasks" "outscale_image_tasks" {
    count = 1

		export_to_osu {
			disk_image_format = "raw"
			osu_bucket = "test"
		}
    image_id = "${outscale_image.outscale_image.image_id}"
}

```

## Argument Reference

The following arguments are supported:

* `export_to_osu` - (optional)	Information about the export task (you must specify the OsuBucket and OsuAkSk parameters at least).	
* `ImageId` - (required)	The ID of the OMI to export.	

## Attributes Reference

* `imageExportTask.N` -	Information about one or more image export tasks.
* `requestId` -	The ID of the request.

See detailed information in [Outscale Image Tasks](http://docs.outscale.com/api_fcu/operations/Action_DescribeImageExportTasks_get.html#_api_fcu-action_describeimageexporttasks_get).

