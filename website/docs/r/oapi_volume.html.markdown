---
layout: "outscale"
page_title: "OUTSCALE: outscale_volume"
sidebar_current: "docs-outscale-resource-volume"
description: |-
  Provides an Outscale Volume resource. This allows volumes to be created, deleted, described and imported.
---

# outscale_volume

  Provides an Outscale Volume resource. You can create a new empty volume or restore a volume from an existing snapshot. You can create the following volume types: Enterprise (io1) for provisioned IOPS SSD volumes, Performance (gp2) for general purpose SSD volumes, or Magnetic (standard) volumes

## Example Usage

```hcl
resource "outscale_volume" "test" {
  subregion_name = "eu-west-2a"
  volume_type = "gp2"
  size = 1
  tags {
    key = "Name"
    value = "tf-acc-test-ebs-volume-test"
  }
}
```

## Argument Reference

The following arguments are supported:

* `iops` - (Optional) The number of I/O operations per second (IOPS). This parameter must be specified only if you create an io1 volume. The maximum number of IOPS allowed for io1 volumes is 13 000.
* `size` - (Optional) The size of the volume, in Gibibytes (GiB). The maximum allowed size for a volume is 14,901 GiB.
* `snapshot_id` - (Optional) The ID of the snapshot from which you want to create the volume.
* `subregion_name` - (Required) The Subregion in which you want to create the volume.
* `volume_type` - (Optional) The type of volume you want to create (`io1` | `gp2` | `standard`). If not specified, a standard volume is created.

See detailed information in [Outscale OAPI Volume](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Attributes Reference

The following attributes are exported:

* `subregion_name` - The Subregion in which the volume was created.
* `iops` - The number of I/O operations per second (IOPS):
  * For io1 volumes, the number of provisioned IOPS
  * For gp2 volumes, the baseline performance of the volume
* `size` - The size of the volume, in Gibibytes (GiB).
* `snapshot_id` - The ID of the snapshot from which you want to create the volume.
* `volume_type` - The type of the volume (`standard` | `gp2` | `io1`).
* `linked_volumes` - Information about your volume attachment.
* `state` - The state of the volume (`creating` | `available` | `in-use` | `deleting` | `error`).
* `tags` - One or more tags associated with the volume.
* `volume_id` - The ID of the volume.
* `request_id` -  The ID of the request.

See detailed information in [Read Volumes](http://docs.outscale.com/api_fcu/definitions/Volume.html#_api_fcu-volume).
See more information in [Read Volume](https://docs-beta.outscale.com/oapi#outscale-api-volume)
