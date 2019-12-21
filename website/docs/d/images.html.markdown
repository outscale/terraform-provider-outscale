---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_image"
sidebar_current: "outscale-image"
description: |-
  [Provides information about images.]
---

# outscale_image Data Source

Provides information about images.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+OMIs).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-image).

## Example Usage

```hcl

data "outscale_images" "images01" {
  filter {
    name   = "image_ids"
    values = ["ami-12345678", "ami-12345679"]
  }
}


```

## Argument Reference

The following arguments are supported:

* `filter` - One or more filters.
  * `account_aliases` - (Optional) The account aliases of the owners of the OMIs.
  * `account_ids` - (Optional) The account IDs of the owners of the OMIs. By default, all the OMIs for which you have launch permissions are described.
  * `architectures` - (Optional) The architectures of the OMIs (`i386` \| `x86_64`).
  * `block_device_mapping_delete_on_vm_deletion` - (Optional) Indicates whether the block device mapping is deleted when terminating the VM.
  * `block_device_mapping_device_names` - (Optional) The device names for the volumes.
  * `block_device_mapping_snapshot_ids` - (Optional) The IDs of the snapshots used to create the volumes.
  * `block_device_mapping_volume_sizes` - (Optional) The sizes of the volumes, in gibibytes (GiB).
  * `block_device_mapping_volume_types` - (Optional) The types of volumes (`standard` \| `gp2` \| `io1`).
  * `descriptions` - (Optional) The descriptions of the OMIs, provided when they were created.
  * `file_locations` - (Optional) The locations where the OMI files are stored on Object Storage Unit (OSU).
  * `image_ids` - (Optional) The IDs of the OMIs.
  * `image_names` - (Optional) The names of the OMIs, provided when they were created.
  * `permissions_to_launch_account_ids` - (Optional) The account IDs of the users who have launch permissions for the OMIs.
  * `permissions_to_launch_global_permission` - (Optional) If `true`, lists all public OMIs. If `false`, lists all private OMIs.
  * `root_device_names` - (Optional) The device names of the root devices (for example, `/dev/sda1`).
  * `root_device_types` - (Optional) The types of root device used by the OMIs (always `ebs`).
  * `states` - (Optional) The states of the OMIs.
  * `tag_keys` - (Optional) The keys of the tags associated with the OMIs.
  * `tag_values` - (Optional) The values of the tags associated with the OMIs.
  * `tags` - (Optional) The key/value combination of the tags associated with the OMIs, in the following format: "Filters":{"Tags":["TAGKEY=TAGVALUE"]}.
  * `virtualization_types` - (Optional) The virtualization types (always `hvm`).

## Attribute Reference

The following attributes are exported:

* `images` - Information about one or more OMIs.
  * `account_alias` - The account alias of the owner of the OMI.
  * `account_id` - The account ID of the owner of the OMI.
  * `architecture` - The architecture of the OMI (by default, `i386`).
  * `block_device_mappings` - One or more block device mappings.
    * `bsu` - Information about the BSU volume to create.
      * `delete_on_vm_deletion` - Set to `true` by default, which means that the volume is deleted when the VM is terminated. If set to `false`, the volume is not deleted when the VM is terminated.
      * `iops` - The number of I/O operations per second (IOPS). This parameter must be specified only if you create an `io1` volume. The maximum number of IOPS allowed for `io1` volumes is `13000`.
      * `snapshot_id` - The ID of the snapshot used to create the volume.
      * `volume_size` - The size of the volume, in gibibytes (GiB).<br />
If you specify a snapshot ID, the volume size must be at least equal to the snapshot size.<br />
If you specify a snapshot ID but no volume size, the volume is created with a size similar to the snapshot one.
      * `volume_type` - The type of the volume (`standard` \| `io1` \| `gp2`). If not specified in the request, a `standard` volume is created.<br />
For more information about volume types, see [Volume Types and IOPS](https://wiki.outscale.net/display/EN/About+Volumes#AboutVolumes-VolumeTypesVolumeTypesandIOPS).
    * `device_name` - The name of the device.
    * `virtual_device_name` - The name of the virtual device (ephemeralN).
  * `creation_date` - The date and time at which the OMI was created.
  * `description` - The description of the OMI.
  * `file_location` - The location where the OMI file is stored on Object Storage Unit (OSU).
  * `image_id` - The ID of the OMI.
  * `image_name` - The name of the OMI.
  * `image_type` - The type of the OMI.
  * `permissions_to_launch` - Information about the users who have permissions for the resource.
    * `account_ids` - The account ID of one or more users who have permissions for the resource.
    * `global_permission` - If `true`, the resource is public. If `false`, the resource is private.
  * `product_codes` - The product code associated with the OMI (`0001` Linux/Unix \| `0002` Windows \| `0004` Linux/Oracle \| `0005` Windows 10).
  * `root_device_name` - The name of the root device.
  * `root_device_type` - The type of root device used by the OMI (always `bsu`).
  * `state` - The state of the OMI.
  * `state_comment` - Information about the change of state.
    * `state_code` - The code of the change of state.
    * `state_message` - A message explaining the change of state.
  * `tags` - One or more tags associated with the OMI.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
