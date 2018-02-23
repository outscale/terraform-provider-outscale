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

* `name` - (Required) A region-unique name for the AMI.
* `description` - (Optional) A longer, human-readable description for the AMI.
* `root_device_name` - (Optional) The name of the root device (for example, /dev/sda1, or /dev/xvda).
* `virtualization_type` - (Optional) Keyword to choose what virtualization mode created instances will use. Can be either "paravirtual" (the default) or "hvm". The choice of * virtualization type changes the set of further arguments that are required, as described below.
* `architecture` - (Optional) Machine architecture for created instances. Defaults to "x86_64".
* `ebs_block_device` - (Optional) Nested block describing an EBS block device that should be attached to created instances. The structure of this block is described below.
* `ephemeral_block_device` - (Optional) Nested block describing an ephemeral block device that should be attached to created instances. The structure of this block is described below.