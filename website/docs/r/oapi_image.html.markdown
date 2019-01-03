---
layout: "outscale"
page_title: "OUTSCALE: outscale_image"
sidebar_current: "docs-outscale-resource-image"
description: |-
  Creates an Outscale machine image (OMI) from an existing virtual machine (VM) which is either running or stopped.
---

# outscale_image

Creates an Outscale machine image (OMI) from an existing virtual machine (VM) which is either running or stopped. This action also creates a snapshot of the root volume of the VM, as well as a snapshot of each Block Storage Unit (BSU) volume attached to the VM.

## Example Usage

```hcl
resource "outscale_vm" "basic" {
  image_id = "ami-b4bd8de2"
  vm_type = "t2.micro"
  keypair_name = "terraform-basic"
  security_group_ids = ["sg-6ed31f3e"]
}

resource "outscale_image" "foo" {
  image_name = "tf-testing-%d"
  vm_id = "${outscale_vm.basic.id}"
  vm_id = "i-b69de1d9"
  no_reboot = "true"
  description = "terraform testing"
}
```

## Argument Reference

The following arguments are supported:

* `architecture` - (Optional) The architecture of the OMI (by default, i386).
* `block_device_mappings` - (Optional) One or more block device mappings.
* `description` - (Optional) A description for the new OMI.
* `file_location` - The pre-signed URL of the OMI manifest file, or the full path to the OMI stored in an OSU bucket. If you specify this parameter, a copy of the OMI is created in your account.
* `image_name` - (Required) A unique name for the new OMI. Constraints: 3–128 alphanumeric characters, underscores (_), spaces ( ), parentheses (()), slashes (/), periods (.), or dashes (-).
* `no_reboot` - If false, the VM shuts down before creating the OMI and then reboots. If true, the VM does not.
* `root_device_name` - (Optional) The name of the root device.
* `vm_id` - (Required) The ID of the VM from which you want to create the OMI.

Nested block_device_mappings blocks have the following structure:

* `device_name` - (Optional) The name of the device.
* `virtual_device_name` - (Optional) The name of the virtual device (ephemeralN).
* `no_device` - Suppresses the device which is included in the block device mapping of the OMI.
* `bsu` - (Optional) One or more parameters used to automatically set up volumes when the instance is launched.

Nested bsu blocks have the following structure:

* `delete_on_vm_deletion` - (Optional) By default or if true, the volume is deleted when terminating the instance. If false, the volume is not deleted when terminating the instance.
* `iops` - The number of I/O operations per second (IOPS). This parameter must be specified only if you create an io1 volume. The maximum number of IOPS allowed for io1 volumes is 13 000.
* `snapshot_id` - (Optional) The ID of the snapshot used to create the volume.
* `volume_size` - The size of the volume, in Gibibytes (GiB). If you specify a snapshot ID, the volume size must be at least equal to the snapshot size. If you specify a snapshot ID but no volume size, the volume is created with a size similar to the snapshot one.
* `volume_type` - The type of the volume (`standard` | `io1` | `gp2` | `sc1` | `st1`).

# Attributes

* `account_alias` - The account alias of the owner of the OMI.
* `account_id` - The account ID of the owner of the OMI.
* `architecture` - The architecture of the OMI (by default, i386).
* `block_device_mappings` - One or more block device mappings.
* `creation_date` - The date and time at which the OMI was created.
* `description` - The description of the OMI.
* `file_location` - The location where the OMI file is stored on Object Storage Unit (OSU).
* `image_id` - The ID of the OMI.
* `image_name` - The name of the  OMI.
* `image_type` - The type of the OMI.
* `permissions_to_launch` - Information about the users who have permissions for the resource.
* `product_codes` - The product code associated with the OMI (001 Linux/Unix | 002 Windows | 003 MapR | 004 Linux/Oracle | 005 Windows 10).
* `root_device_name` - The name of the root device.
* `root_device_type` - The type of root device used by the OMI (always bsu).
* `state` - The state of the OMI.
* `state_comment` - Information about the change of state.
* `is_public` - If true, the OMI has public launch permissions.
* `tags` - One or more tags associated with the OMI.
* `request_id` - The ID of the request.

See detailed information in [CreateImage](http://docs.outscale.com/api_fcu/operations/Action_CreateImage_get.html#_api_fcu-action_createimage_get).

See detailed information in [ReadImages](http://docs.outscale.com/api_fcu/operations/Action_DescribeImages_get.html#_api_fcu-action_describeimages_get).

