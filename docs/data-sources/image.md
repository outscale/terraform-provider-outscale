---
layout: "outscale"
page_title: "OUTSCALE: outscale_image"
subcategory: "Compute"
sidebar_current: "outscale-image"
description: |-
  [Provides information about an image.]
---

# outscale_image Data Source

Provides information about an image.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-OMIs.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-image).

## Example Usage

```hcl
data "outscale_image" "image01" {
    filter {
        name   = "image_ids"
        values = ["ami-12345678"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `account_aliases` - (Optional) The account aliases of the owners of the OMIs.
    * `account_ids` - (Optional) The account IDs of the owners of the OMIs. By default, all the OMIs for which you have launch permissions are described.
    * `architectures` - (Optional) The architectures of the OMIs (`i386` \| `x86_64`).
    * `block_device_mapping_delete_on_vm_deletion` - (Optional) Whether the volumes are deleted or not when terminating the VM.
    * `block_device_mapping_device_names` - (Optional) The device names for the volumes.
    * `block_device_mapping_snapshot_ids` - (Optional) The IDs of the snapshots used to create the volumes.
    * `block_device_mapping_volume_sizes` - (Optional) The sizes of the volumes, in gibibytes (GiB).
    * `block_device_mapping_volume_types` - (Optional) The types of volumes (`standard` \| `gp2` \| `io1`).
    * `boot_modes` - (Optional) The boot modes compatible with the OMIs. Possible values: `uefi` | `legacy`.
    * `descriptions` - (Optional) The descriptions of the OMIs, provided when they were created.
    * `file_locations` - (Optional) The locations of the buckets where the OMI files are stored.
    * `hypervisors` - (Optional) The hypervisor type of the OMI (always `xen`).
    * `image_ids` - (Optional) The IDs of the OMIs.
    * `image_names` - (Optional) The names of the OMIs, provided when they were created.
    * `permissions_to_launch_account_ids` - (Optional) The account IDs which have launch permissions for the OMIs.
    * `permissions_to_launch_global_permission` - (Optional) If true, lists all public OMIs. If false, lists all private OMIs.
    * `product_code_names` - (Optional) The names of the product codes associated with the OMI.
    * `product_codes` - (Optional) The product codes associated with the OMI.
    * `root_device_names` - (Optional) The name of the root device. This value must be /dev/sda1.
    * `root_device_types` - (Optional) The types of root device used by the OMIs (`bsu` or `ebs`).
    * `secure_boot` - (Optional) Whether secure boot is activated or not.
    * `states` - (Optional) The states of the OMIs (`pending` \| `available` \| `failed`).
    * `tag_keys` - (Optional) The keys of the tags associated with the OMIs.
    * `tag_values` - (Optional) The values of the tags associated with the OMIs.
    * `tags` - (Optional) The key/value combinations of the tags associated with the OMIs, in the following format: `TAGKEY=TAGVALUE`.
    * `virtualization_types` - (Optional) The virtualization types (always `hvm`).

## Attribute Reference

The following attributes are exported:

* `account_alias` - The account alias of the owner of the OMI.
* `account_id` - The account ID of the owner of the OMI.
* `architecture` - The architecture of the OMI.
* `block_device_mappings` - One or more block device mappings.
    * `bsu` - Information about the BSU volume to create.
        * `delete_on_vm_deletion` - By default or if set to true, the volume is deleted when terminating the VM. If false, the volume is not deleted when terminating the VM.
        * `iops` - The number of I/O operations per second (IOPS). This parameter must be specified only if you create an `io1` volume. The maximum number of IOPS allowed for `io1` volumes is `13000` with a maximum performance ratio of 300 IOPS per gibibyte.
        * `snapshot_id` - The ID of the snapshot used to create the volume.
        * `volume_size` - The size of the volume, in gibibytes (GiB).<br />
If you specify a snapshot ID, the volume size must be at least equal to the snapshot size.<br />
If you specify a snapshot ID but no volume size, the volume is created with a size similar to the snapshot one.
        * `volume_type` - The type of the volume (`standard` \| `io1` \| `gp2`). If not specified in the request, a `standard` volume is created.<br />
For more information about volume types, see [About Volumes > Volume Types and IOPS](https://docs.outscale.com/en/userguide/About-Volumes.html#_volume_types_and_iops).
    * `device_name` - The device name for the volume. For a root device, you must use `/dev/sda1`. For other volumes, you must use `/dev/sdX`, `/dev/sdXX`, `/dev/xvdX`, or `/dev/xvdXX` (where the first `X` is a letter between `b` and `z`, and the second `X` is a letter between `a` and `z`).
    * `virtual_device_name` - The name of the virtual device (`ephemeralN`).
* `boot_modes` - The boot modes compatible with the OMI. Possible values: `uefi` | `legacy`.
* `creation_date` - The date and time (UTC) at which the OMI was created.
* `description` - The description of the OMI.
* `file_location` - The location from which the OMI files were created.
* `image_id` - The ID of the OMI.
* `image_name` - The name of the OMI.
* `image_type` - The type of the OMI.
* `permissions_to_launch` - Permissions for the resource.
    * `account_ids` - One or more account IDs that the permission is associated with.
    * `global_permission` - A global permission for all accounts.<br />
(Request) Set this parameter to true to make the resource public (if the parent parameter is `Additions`) or to make the resource private (if the parent parameter is `Removals`).<br />
(Response) If true, the resource is public. If false, the resource is private.
* `product_codes` - The product codes associated with the OMI.
* `root_device_name` - The name of the root device.
* `root_device_type` - The type of root device used by the OMI (always `bsu`).
* `secure_boot` - Whether secure boot is activated or not.
* `state_comment` - Information about the change of state.
    * `state_code` - The code of the change of state.
    * `state_message` - A message explaining the change of state.
* `state` - The state of the OMI (`pending` \| `available` \| `failed`).
* `tags` - One or more tags associated with the OMI.
    * `key` - The key of the tag, with a minimum of 1 character.
    * `value` - The value of the tag, between 0 and 255 characters.
