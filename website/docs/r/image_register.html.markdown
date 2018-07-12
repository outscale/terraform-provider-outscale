---
layout: "outscale"
page_title: "OUTSCALE: outscale_image_register"
sidebar_current: "docs-outscale-resource-image-register"
description: |-
  Registers an Outscale Machine Image (OMI) to finalize its creation process.
---

# outscale_image_register

Registers an Outscale Machine Image (OMI) to finalize its creation process.. [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
resource "outscale_volume" "outscale_volume" {
  availability_zone = "eu-west-2a"
  size              = 40
}

resource "outscale_snapshot" "outscale_snapshot" {
  volume_id = "${outscale_volume.outscale_volume.volume_id}"
}

resource "outscale_image_register" "outscale_image_register" {
  name = "registeredImageFromSnapshot-%d"

  root_device_name = "/dev/sda1"

  block_device_mapping {
    snapshot_id = "${outscale_snapshot.outscale_snapshot.snapshot_id}"
    device_name = "/dev/sda1"
  }
}
```

## Argument Reference

The following arguments are supported:

* `architecture` - (Optional) The architecture of the OMI (set to i386 by default).
* `block_device_mapping` - (Optional) One or more Block Device Mapping entries.
  * `device_name` - (Optional) The name of the device. To modify the deleteOnTermination attribute of a volume, this parameter is required.
  * `no_device` - (Optional) Suppresses the device which is included in the block device mapping of the OMI.
  * `virtual_name` - (Optional) The name of the virtual device (ephemeralN).
  * `delete_on_termination` - (Optional) By default or if true, the volume is deleted when terminating the instance. If false, the volume is not deleted when terminating the instance.
  * `iops` - (Optional) The number of IOPS supported by the volume.
  * `snapshot_id` - (Optional) The ID of the snapshot used to create the volume.
  * `volume_size` - (Optional) The size of the volume, in Gibibytes (GiB). If you specify a snapshot ID, the volume size must be at least equal to the snapshot size. If you specify a snapshot ID but no volume size, the volume is created with a size similar to the snapshot one.
  * `volume_type` - The type of the volume (standard | io1 | gp2).
* `description` - (Optional) A description for the OMI.
* `image_location` - (Optional) The pre-signed URL of the OMI manifest file, or the full path to the OMI stored in an OSU bucket. If you specify this parameter, a copy of the OMI is created in your account.
* `name` - (Required) A unique name for the OMI.
* `toot_device_name` - (Optional) The name of the root device (for example, /dev/sda1)

## Attributes Reference

* `imageId` - The ID of the newly registered OMI.
* `requestId` - The ID of the request.

See detailed information in [Outscale Image Tasks](http://docs.outscale.com/api_fcu/operations/Action_RegisterImage_get.html#_api_fcu-action_registerimage_get).
