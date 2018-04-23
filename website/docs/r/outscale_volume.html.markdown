---
layout: "outscale"
page_title: "OUTSCALE: outscale_volume"
sidebar_current: "docs-outscale-resource-volume"
description: |-
  Provides an Outscale Volume resource. This allows volumes to be created, deleted, described and imported.
---

# outscale_volume

  Provides an Outscale Volume resource. This allows volumes to be created, deleted, described and imported. Instances also support [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
resource "outscale_volume" "example" {
    availability_zone = "eu-west-2a"
    volume_type = "gp2"
    size = 40
    tags {
        Name = "External Volume"
    }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Required) The Availability Zone in which you want to create the volume.
* `iops` - (Optional) The number of I/O operations per second, with a maximum ratio of 30 IOPS/GiB (only for io1 volumes).
* `size` - (Optional) The size of the volume, in Gibibytes (GiB). The maximum allowed size for a volume is 14,901 GiB.
* `snapshot_id` - (Optional) The ID of the snapshot from which you want to create the volume.
* `volume_type` - (Optional) The type of volume you want to create (io1 | gp2 | standard | sc1 | st1).

See detailed information in [Outscale Images](http://docs.outscale.com/api_fcu/operations/Action_CreateImage_get.html#_api_fcu-action_createimage_get).


## Attributes Reference

The following attributes are exported:

* `attachment_set` - Information about your volume attachment.
* `availability_zone` - The Availability Zone where the volume is.
* `iops` - The number of I/O operations per second (only for io1 and gp2 volumes).
* `size` - The size of the volume, in Gibibytes (GiB).
* `snapshot_id` - The ID of the snapshot from which the volume was created.
* `status` - The state of the volume (creating| available| in-use| deleting| error).
* `tag_set` - One or more tags associated with the volume.
* `volume_id` - The ID of the volume.
* `volume_type` - The type of the volume (`standard` | `gp2` | `io1` | `sc1` | `st1`).

See detailed information in [Volume Description](http://docs.outscale.com/api_fcu/definitions/Volume.html#_api_fcu-volume).
