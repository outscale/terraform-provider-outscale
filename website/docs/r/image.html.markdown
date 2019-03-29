---
layout: "outscale"
page_title: "OUTSCALE: outscale_image"
sidebar_current: "docs-outscale-resource-image"
description: |-
  Creates an OMI from an existing instance which is either running or stopped.
---

# outscale_image

Creates an OMI from an existing instance which is either running or stopped.

## Example Usage

```hcl
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
	security_group = ["sg-6ed31f3e"]
}

resource "outscale_image" "foo" {
	name = "tf-outscale-image-name"
	instance_id = "${outscale_vm.basic.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A unique name for the new OMI.
* `instance_id` - (Required) The ID of the instance from which you want to create the OMI.
* `description` - (Optional) A description for the new OMI.
* `root_device_name` - (Optional) The name of the root device (for example, /dev/sda1, or /dev/xvda).
* `image_type` - (Optional) Keyword to choose what virtualization mode created instances will use. Can be either "paravirtual" (the default) or "hvm". The choice of * virtualization type changes the set of further arguments that are required, as described below.
* `architecture` - (Optional) Machine architecture for created instances. Defaults to "x86_64".
* `block_device_mapping` - (Optional) One or more block device mapping entries.
* `dry_run` - If set to true, checks whether you have the required permissions to perform the action.
* `no_reboot` - If set to false, the instance shuts down before creating the OMI and then reboots. If set to true, the instance does not.
* `architecture` - The architecture of the OMI.
* `creation_date` - The date and time at which the OMI was created.
* `image_location` - The location where the OMI is stored on Object Storage Unit (OSU).
* `image_owner_alias` - The account alias of the owner of the OMI.
* `image_owner_id` - The account ID of the owner of the OMI.
* `image_state` - The current state of the OMI (if available, the image is registered and can be used to launch an instance).
* `image_type` - The type of the OMI.
* `is_public` - If true, the OMI has public launch permissions.
* `product_codes` - One or more product codes associated with the OMI.
* `root_device_name` - The name of the device to which the root device is plugged in (for example, /dev/sda1).
* `root_device_type` - The type of root device used by the OMI.
* `state_reason` - The reason for the OMI state change.
* `tag_set` - One or more tags associated with the OMI.

Nested block_device_mapping blocks have the following structure:

* `device_name` - (Optional) The name of the device.
* `virtual_name` - (Optional) The name of the virtual device (ephemeralN).
* `no_device` - Suppresses the device which is included in the block device mapping of the OMI.
* `ebs` - (Optional) One or more parameters used to automatically set up volumes when the instance is launched.

Nested ebs blocks have the following structure:

* `delete_on_termination` - (Optional) By default or if true, the volume is deleted when terminating the instance. If false, the volume is not deleted when terminating the instance.
* `iops` - The number of IOPS supported by the volume.
* `snapshot_id` - (Optional) The ID of the snapshot used to create the volume.
* `volume_size` - The size of the volume, in Gibibytes (GiB).
If you specify a snapshot ID, the volume size must be at least equal to the snapshot size.
If you specify a snapshot ID but no volume size, the volume is created with a size similar to the snapshot one.
* `volume_type` - The type of the volume (`standard` | `io1` | `gp2` | `sc1` | `st1`).

The timeouts block allows you to specify timeouts for certain actions:

* `create` - (Defaults to 40 mins) Used when creating the Image
* `update` - (Defaults to 40 mins) Used when updating the Image
* `delete` - (Defaults to 90 mins) Used when deregistering the Image


# Attributes

* `image_id` - The ID of the new OMI.
* `request_id` - The ID of the request.


See detailed information in [Create Image](http://docs.outscale.com/api_fcu/operations/Action_CreateImage_get.html#_api_fcu-action_createimage_get).
See detailed information in [Describe Images](http://docs.outscale.com/api_fcu/operations/Action_DescribeImages_get.html#_api_fcu-action_describeimages_get).
See detailed information in [Register Image](http://docs.outscale.com/api_fcu/operations/Action_RegisterImage_get.html#_api_fcu-action_registerimage_get).
