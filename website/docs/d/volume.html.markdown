---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_volume"
sidebar_current: "outscale-volume"
description: |-
  [Provides information about a specific volume.]
---

# outscale_volume Data Source

Provides information about a specific volume.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Volumes).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-volume).

## Example Usage

```hcl

data "outscale_volume" "outscale_volume01" {
  filter {
    name   = "volume_ids"
    values = ["vol-12345678"]
  }
}


```

## Argument Reference

The following arguments are supported:

* `filter` - One or more filters.
  * `creation_dates` - (Optional) The dates and times at which the volumes were created.
  * `link_volume_delete_on_vm_deletion` - (Optional) Indicates whether the volumes are deleted when terminating the VMs.
  * `link_volume_device_names` - (Optional) The VM device names.
  * `link_volume_link_dates` - (Optional) The dates and times at which the volumes were created.
  * `link_volume_link_states` - (Optional) The attachment states of the volumes (`attaching` \| `detaching` \| `attached` \| `detached`).
  * `link_volume_vm_ids` - (Optional) One or more IDs of VMs.
  * `snapshot_ids` - (Optional) The snapshots from which the volumes were created.
  * `subregion_names` - (Optional) The names of the Subregions in which the volumes were created.
  * `tag_keys` - (Optional) The keys of the tags associated with the volumes.
  * `tag_values` - (Optional) The values of the tags associated with the volumes.
  * `tags` - (Optional) The key/value combination of the tags associated with the volumes, in the following format: "Filters":{"Tags":["TAGKEY=TAGVALUE"]}.
  * `volume_ids` - (Optional) The IDs of the volumes.
  * `volume_sizes` - (Optional) The sizes of the volumes, in gibibytes (GiB).
  * `volume_states` - (Optional) The states of the volumes (`creating` \| `available` \| `in-use` \| `deleting` \| `error`).
  * `volume_types` - (Optional) The types of the volumes (`standard` \| `gp2` \| `io1`).

## Attribute Reference

The following attributes are exported:

* `volumes` - Information about one or more volumes.
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
