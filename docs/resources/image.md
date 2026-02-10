---
layout: "outscale"
page_title: "OUTSCALE: outscale_image"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-image"
description: |-
  [Manages an image.]
---

# outscale_image Resource

Manages an image.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-OMIs.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-image).

## Example Usage

### Create an image

```hcl
resource "outscale_image" "image01" {
    image_name = "terraform-omi-create"
    vm_id      = var.vm_id
    no_reboot  = "true"
}
```

### Import an image
~> **Important** Make sure the manifest file is still valid.

```hcl
resource "outscale_image" "image02" {
    description   = "Terraform register OMI"
    image_name    = "terraform-omi-register"
    file_location = "<URL>"
}
```

### Copy an image

```hcl
resource "outscale_image" "image03" {
    description        = "Terraform copy OMI"
    image_name         = "terraform-omi-copy"
    source_image_id    = "ami-12345678"
    source_region_name = "eu-west-2"
}
```

### Create an image with a Block Storage Unit (BSU) volume

```hcl
resource "outscale_image" "image04" {
    image_name = "terraform-omi-bsu"
    block_device_mappings {
        device_name = "/dev/sda1" 
        bsu {
        snapshot_id           = "snap-12345678"
        volume_size           = "120"
        volume_type           = "io1"
        iops                  = 150
        delete_on_vm_deletion = "true"
        }
    }
    root_device_name = "/dev/sda1"
    description      = "Terraform OMI with BSU"
}
```

## Argument Reference

The following arguments are supported:

* `architecture` - (Optional) **When registering from a snapshot:** The architecture of the OMI (`i386` or `x86_64`).
* `block_device_mappings` - (Optional) **(required) When registering from a snapshot:** One or more block device mappings.
    * `bsu` - Information about the BSU volume to create.
        * `delete_on_vm_deletion` - (Optional) By default or if set to true, the volume is deleted when terminating the VM. If false, the volume is not deleted when terminating the VM.
        * `iops` - (Optional) The number of I/O operations per second (IOPS). This parameter must be specified only if you create an `io1` volume. The maximum number of IOPS allowed for `io1` volumes is `13000` with a maximum performance ratio of 300 IOPS per gibibyte.
        * `snapshot_id` - (Optional) The ID of the snapshot used to create the volume.
        * `volume_size` - (Optional) The size of the volume, in gibibytes (GiB).<br />
If you specify a snapshot ID, the volume size must be at least equal to the snapshot size.<br />
If you specify a snapshot ID but no volume size, the volume is created with a size similar to the snapshot one.
        * `volume_type` - (Optional) The type of the volume (`standard` \| `io1` \| `gp2`). If not specified in the request, a `standard` volume is created.<br />
For more information about volume types, see [About Volumes > Volume Types and IOPS](https://docs.outscale.com/en/userguide/About-Volumes.html#_volume_types_and_iops).
    * `device_name` - (Optional) The device name for the volume. For a root device, you must use `/dev/sda1`. For other volumes, you must use `/dev/sdX`, `/dev/sdXX`, `/dev/xvdX`, or `/dev/xvdXX` (where the first `X` is a letter between `b` and `z`, and the second `X` is a letter between `a` and `z`).
    * `virtual_device_name` - (Optional) The name of the virtual device (`ephemeralN`).
* `boot_modes` - (Optional) The boot modes compatible with the OMI. Possible values: `uefi` | `legacy`.
* `description` - (Optional) A description for the new OMI.
* `file_location` - (Optional) **(required) When registering from a bucket by using a manifest file:** The pre-signed URL of the manifest file for the OMI you want to register. For more information, see [Creating a Pre-signed URL](https://docs.outscale.com/en/userguide/Creating-a-Pre-Signed-URL.html).
* `image_name` - (Optional) A unique name for the new OMI.<br />
Constraints: 3-128 alphanumeric characters, underscores (`_`), spaces (` `), parentheses (`()`), slashes (`/`), periods (`.`), or dashes (`-`).
* `no_reboot` - (Optional) **When creating from a VM:** If false, the VM shuts down before creating the OMI and then reboots. If true, the VM does not.
* `product_codes` - (Optional) The product codes associated with the OMI.
* `root_device_name` - (Optional) **(required) When registering from a snapshot:** The name of the root device for the new OMI.
* `source_image_id` - (Optional) **(required) When copying an OMI:** The ID of the OMI you want to copy.
* `source_region_name` - (Optional) **(required) When copying an OMI:** The name of the source Region (always the same as the Region of your account).
* `tags` - (Optional) A tag to add to this resource. You can specify this argument several times.
    * `key` - (Required) The key of the tag, with a minimum of 1 character.
    * `value` - (Required) The value of the tag, between 0 and 255 characters.
* `tpm_mandatory` - (Optional) By default or if set to false, a virtual Trusted Platform Module (vTPM) is not mandatory on VMs created from this OMI. If true, VMs created from this OMI must have a vTPM enabled.
* `vm_id` - (Optional) **(required) When creating from a VM:** The ID of the VM from which you want to create the OMI.

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
* `tpm_mandatory` - If true, a virtual Trusted Platform Module (vTPM) is mandatory for VMs created from this OMI. If false, a vTPM is not mandatory.

## Import

An image can be imported using its ID. For example:

```console

$ terraform import outscale_image.ImportedImage ami-12345678

```