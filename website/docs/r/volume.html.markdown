---
layout: "outscale"
page_title: "OUTSCALE: outscale_volume"
sidebar_current: "outscale-volume"
description: |-
  [Manages a volume.]
---

# outscale_volume Resource

Manages a volume.

~> **Note** When updating an existing volume linked to a running virtual machine (VM), the VM will temporarily be stopped. For more information, see [About Instance Lifecycle](https://docs.outscale.com/en/userguide/About-Instance-Lifecycle.html).

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Volumes.html).
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
* `size` - (Optional) The size of the volume, in gibibytes (GiB). The maximum allowed size for a volume is 14901 GiB. This parameter is required if the volume is not created from a snapshot (`snapshot_id` unspecified).<br />If you update an existing volume, this value must be equal to or greater than the current size of the volume. You might need to reconfigure the file system, for more information see [Increasing the Size of a Volume](https://docs.outscale.com/en/userguide/Increasing-the-Size-of-a-Volume.html).
* `snapshot_id` - (Optional) The ID of the snapshot from which you want to create the volume.
* `subregion_name` - (Required) The Subregion in which you want to create the volume.
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.
* `volume_type` - (Optional) The type of volume you want to create (`io1` | `gp2` | `standard`). If not specified, a `standard` volume is created.<br />If you create or update an existing volume to an `io1` volume, you must also specify the `iops` parameter.<br />For more information about volume types, see [Volume Types and IOPS](https://docs.outscale.com/en/userguide/About-Volumes.html#_volume_types_and_iops).

## Attribute Reference

The following attributes are exported:

* `iops` - The number of I/O operations per second (IOPS):<br />- For `io1` volumes, the number of provisioned IOPS.<br />- For `gp2` volumes, the baseline performance of the volume.
* `linked_volumes` - Information about your volume attachment.
    * `delete_on_vm_deletion` - If true, the volume is deleted when terminating the VM. If false, the volume is not deleted when terminating the VM.
    * `device_name` - The name of the device.
    * `state` - The state of the attachment of the volume (`attaching` \| `detaching` \| `attached` \| `detached`).
    * `vm_id` - The ID of the VM.
    * `volume_id` - The ID of the volume.
* `size` - The size of the volume, in gibibytes (GiB).
* `snapshot_id` - The snapshot from which the volume was created.
* `state` - The state of the volume (`creating` \| `available` \| `in-use` \| `updating` \| `deleting` \| `error`).
* `subregion_name` - The Subregion in which the volume was created.
* `tags` - One or more tags associated with the volume.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
* `volume_id` - The ID of the volume.
* `volume_type` - The type of the volume (`standard` \| `gp2` \| `io1`).

## Import

A volume can be imported using its ID. For example:

```console

$ terraform import outscale_volume.ImportedVolume vol-12345678

```