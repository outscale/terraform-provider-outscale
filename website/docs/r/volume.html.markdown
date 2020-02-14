---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_volume"
sidebar_current: "outscale-volume"
description: |-
  [Manages a volume.]
---

# outscale_volume Resource

Manages a volume.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Volumes).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-volume).

## Example Usage

```hcl

resource "outscale_volume" "volume01" {
  subregion_name = "${var.region}a"
  size           = 10
  iops           = 100
  volume_type    = "io1"
}


```

## Argument Reference

The following arguments are supported:

* `iops` - (Optional) The number of I/O operations per second (IOPS). This parameter must be specified only if you create an `io1` volume. The maximum number of IOPS allowed for `io1` volumes is `13000`.
* `size` - (Optional) The size of the volume, in gibibytes (GiB). The maximum allowed size for a volume is 14,901 GiB.
* `snapshot_id` - (Optional) The ID of the snapshot from which you want to create the volume.
* `subregion_name` - (Required) The Subregion in which you want to create the volume.
* `volume_type` - (Optional) The type of volume you want to create (`io1` \| `gp2` \| `standard`). If not specified, a `standard` volume is created.<br />
For more information about volume types, see [Volume Types and IOPS](https://wiki.outscale.net/display/EN/About+Volumes#AboutVolumes-VolumeTypesVolumeTypesandIOPS).
* `tags` - One or more tags to add to this resource.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
    
## Attribute Reference

The following attributes are exported:

* `volume` - Information about the volume.
  * `iops` - The number of I/O operations per second (IOPS):  
    For `io1` volumes, the number of provisioned IOPS.  
    For `gp2` volumes, the baseline performance of the volume.
  * `linked_volumes` - Information about your volume attachment.
    * `delete_on_vm_deletion` - If `true`, the volume is deleted when the VM is terminated.
    * `device_name` - The name of the device.
    * `state` - The state of the attachment of the volume (`attaching` \| `detaching` \| `attached` \| `detached`).
    * `vm_id` - The ID of the VM.
    * `volume_id` - The ID of the volume.
  * `size` - The size of the volume, in gibibytes (GiB).
  * `snapshot_id` - The snapshot from which the volume was created.
  * `state` - The state of the volume (`creating` \| `available` \| `in-use` \| `deleting` \| `error`).
  * `subregion_name` - The Subregion in which the volume was created.
  * `tags` - One or more tags associated with the volume.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
  * `volume_id` - The ID of the volume.
  * `volume_type` - The type of the volume (`standard` \| `gp2` \| `io1`).
