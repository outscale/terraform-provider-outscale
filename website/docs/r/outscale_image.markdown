---
layout: "outscale"
page_title: "OUTSCALE: outscale_image"
sidebar_current: "docs-outscale-resource-image"
description: |-
The Outscale Image resource allows the creation and management of a completely-custom Outscale VM Image.

If you just want to duplicate an existing Outscale Image, possibly copying it to another region, it's better to use outscale_image_copy instead.

If you just want to share an existing Outscale Image with another Outscale account, it's better to use outscale_image_launch_permission instead.
---

## Example Usage

```hcl
resource "outscale_vm" "basic" {
	image_id = "ami-8a6a0120"
	instance_type = "t2.micro"
	key_name = "terraform-basic"
	security_group = ["sg-6ed31f3e"]
}

resource "outscale_image" "foo" {
	name = "tf-testing-%d"
	instance_id = "${outscale_vm.basic.id}"
}
```

## Argument Reference

The following arguments are supported:

* `name` - (Required) A region-unique name for the Image.
* `instance_id` - (Required) Instance ID for the Image.
* `description` - (Optional) A longer, human-readable description for the Image.
* `root_device_name` - (Optional) The name of the root device (for example, /dev/sda1, or /dev/xvda).
* `image_type` - (Optional) Keyword to choose what virtualization mode created instances will use. Can be either "paravirtual" (the default) or "hvm". The choice of * virtualization type changes the set of further arguments that are required, as described below.
* `architecture` - (Optional) Machine architecture for created instances. Defaults to "x86_64".
* `block_device_mapping` - (Optional) Nested block describing an ephemeral block device that should be attached to created instances. The structure of this block is described below.


Nested block_device_mapping blocks have the following structure:

* `device_name` - (Optional) The path at which the device is exposed to created instances.
* `virtual_name` - (Optional) A name for the ephemeral device, of the form "ephemeralN" where N is a volume number starting from zero.
* `ebs` - (Optional) Nested block describing an EBS block device that should be attached to created instances. The structure of this block is described below.

Nested ebs blocks have the following structure:

* `delete_on_termination` - (Optional) Boolean controlling whether the EBS volumes created to support each created instance will be deleted once that instance is terminated.
* `iops` - (Required only when volume_type is "io1") Number of I/O operations per second the created volumes will support.
* `snapshot_id` - (Optional) The id of an EBS snapshot that will be used to initialize the created EBS volumes. If set, the volume_size attribute must be at least as large as the referenced snapshot.
* `volume_size` - (Required unless snapshot_id is set) The size of created volumes in GiB. If snapshot_id is set and volume_size is omitted then the volume will have the same size as the selected snapshot.
* `volume_type` - (Optional) The type of EBS volume to create. Can be one of "standard" (the default), "io1" or "gp2".

### Timeouts

The timeouts block allows you to specify timeouts for certain actions:

* `create` - (Defaults to 40 mins) Used when creating the Image
* `update` - (Defaults to 40 mins) Used when updating the Image
* `delete` - (Defaults to 90 mins) Used when deregistering the Image

## Attributes Reference

The following attributes are exported:

* `image_id` - The ID of the created Image.
* `snapshot_id` - The Snapshot ID for the root volume (for EBS-backed AMIs)