---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_image"
sidebar_current: "outscale-image"
description: |-
  [Manages an image.]
---

# outscale_image Resource

Manages an image.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+OMIs).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-image).

## Example Usage

```hcl

# Create an image

resource "outscale_image" "image01" {
  image_name = "terraform-OMI-create"
  vm_id      = var.vm_id
  no_reboot  = "true"
}

# Import an image
## Make sure the manifest is still valide. You can find its location URL in Cockpit.

resource "outscale_image" "image02" {
  description   = "Terraform register OMI"
  image_name    = "terraform-OMI-register"
  file_location = "http://osu.eu-west-2.outscale.com/new-export-omi/omi-for-terraform/ami-b7d7f165/manifest?AWSAccessKeyId=XXXXXXXXXXX&Expires=XXXXXXXXX&Signature=XXXXXXXXXX"
}

# Copy an image

resource "outscale_image" "image03" {
  description        = "Terraform copy OMI"
  image_name         = "terraform-OMI-copy"
  source_image_id    = "ami-249bc1f9"
  source_region_name = "eu-west-2"
}


```

## Argument Reference

The following arguments are supported:

* `architecture` - (Optional) The architecture of the OMI (by default, `i386`).
* `block_device_mappings` - (Optional) One or more block device mappings.
  * `bsu` - Information about the BSU volume to create.
    * `delete_on_vm_deletion` - (Optional) Set to `true` by default, which means that the volume is deleted when the VM is terminated. If set to `false`, the volume is not deleted when the VM is terminated.
    * `iops` - (Optional) The number of I/O operations per second (IOPS). This parameter must be specified only if you create an `io1` volume. The maximum number of IOPS allowed for `io1` volumes is `13000`.
    * `snapshot_id` - (Optional) The ID of the snapshot used to create the volume.
    * `volume_size` - (Optional) The size of the volume, in gibibytes (GiB).<br />
If you specify a snapshot ID, the volume size must be at least equal to the snapshot size.<br />
If you specify a snapshot ID but no volume size, the volume is created with a size similar to the snapshot one.
    * `volume_type` - (Optional) The type of the volume (`standard` \| `io1` \| `gp2`). If not specified in the request, a `standard` volume is created.<br />
For more information about volume types, see [Volume Types and IOPS](https://wiki.outscale.net/display/EN/About+Volumes#AboutVolumes-VolumeTypesVolumeTypesandIOPS).
  * `device_name` - (Optional) The name of the device.
  * `virtual_device_name` - (Optional) The name of the virtual device (ephemeralN).
* `description` - (Optional) A description for the new OMI.
* `file_location` - (Optional) The pre-signed URL of the OMI manifest file, or the full path to the OMI stored in an OSU bucket. If you specify this parameter, a copy of the OMI is created in your account.
* `image_name` - (Optional) A unique name for the new OMI.<br />
Constraints: 3-128 alphanumeric characters, underscores (_), spaces ( ), parentheses (()), slashes (/), periods (.), or dashes (-).
* `no_reboot` - (Optional) If `false`, the VM shuts down before creating the OMI and then reboots. If `true`, the VM does not.
* `root_device_name` - (Optional) The name of the root device.
* `source_image_id` - (Optional) The ID of the OMI you want to copy.
* `source_region_name` - (Optional) The name of the source Region, which must be the same as the Region of your account.
* `vm_id` - (Optional) The ID of the VM from which you want to create the OMI.

## Attribute Reference

The following attributes are exported:

* `image` - Information about the OMI.
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
