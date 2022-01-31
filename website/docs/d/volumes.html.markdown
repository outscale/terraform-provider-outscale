---
layout: "outscale"
page_title: "OUTSCALE: outscale_volumes"
sidebar_current: "outscale-volumes"
description: |-
  [Provides information about volumes.]
---

# outscale_volumes Data Source

Provides information about volumes.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Volumes.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-volume).

## Example Usage

```hcl
data "outscale_volumes" "outscale_volumes01" {
  filter {
    name   = "volume_ids"
    values = ["vol-12345678", "vol-12345679"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `creation_dates` - (Optional) The dates and times at which the volumes were created.
    * `link_volume_delete_on_vm_deletion` - (Optional) Whether the volumes are deleted or not when terminating the VMs.
    * `link_volume_device_names` - (Optional) The VM device names.
    * `link_volume_link_dates` - (Optional) The dates and times at which the volumes were created.
    * `link_volume_link_states` - (Optional) The attachment states of the volumes (`attaching` \| `detaching` \| `attached` \| `detached`).
    * `link_volume_vm_ids` - (Optional) One or more IDs of VMs.
    * `snapshot_ids` - (Optional) The snapshots from which the volumes were created.
    * `subregion_names` - (Optional) The names of the Subregions in which the volumes were created.
    * `tag_keys` - (Optional) The keys of the tags associated with the volumes.
    * `tag_values` - (Optional) The values of the tags associated with the volumes.
    * `tags` - (Optional) The key/value combination of the tags associated with the volumes, in the following format: &quot;Filters&quot;:{&quot;Tags&quot;:[&quot;TAGKEY=TAGVALUE&quot;]}.
    * `volume_ids` - (Optional) The IDs of the volumes.
    * `volume_sizes` - (Optional) The sizes of the volumes, in gibibytes (GiB).
    * `volume_states` - (Optional) The states of the volumes (`creating` \| `available` \| `in-use` \| `updating` \| `deleting` \| `error`).
    * `volume_types` - (Optional) The types of the volumes (`standard` \| `gp2` \| `io1`).

## Attribute Reference

The following attributes are exported:

* `volumes` - Information about one or more volumes.
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
