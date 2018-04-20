---
layout: "outscale"
page_title: "OUTSCALE: outscale_image_register"
sidebar_current: "docs-outscale-resource-image-register"
description: |-
	Registers an Outscale Mmachine image (OMI) to finalize its creation process.
---

# outscale_image

Registers an Outscale Mmachine image (OMI) to finalize its creation process.

## Example Usage

```hcl
resource "outscale_vm" "outscale_vm" {
    count = 1

    image_id                    = "ami-880caa66"
    instance_type               = "c4.large"
}

resource "outscale_image_register" "test" {
    count = 1
    name        = "image_${outscale_vm.outscale_vm.id}"
    instance_id = "${outscale_vm.outscale_vm.id}"
	}
```

## Argument Reference

The following arguments are supported:

* `architecture` - (Optional) Machine architecture for created instances. Defaults to "x86_64".
* `block_device_mapping` - (Optional) One or more block device mapping entries.
* `description` - (Optional) A description for the new OMI.
* `image_location` - The location where the OMI is stored on Object Storage Unit (OSU).
* `name` - (Required) A unique name for the new OMI.
* `root_device_name` - (Optional) The name of the root device (for example, /dev/sda1, or /dev/xvda).

## Attribute Reference

* `imageId`	- The ID of the newly registered OMI.	false	string
* `architecture` - Machine architecture for created instances. Defaults to "x86_64".
* `block_device_mapping` - One or more block device mapping entries.
* `description` - A description for the new OMI.
* `image_location` - The location where the OMI is stored on Object Storage Unit (OSU).
* `name` - (Required) A unique name for the new OMI.
* `root_device_name` - The name of the root device (for example, /dev/sda1, or /dev/xvda).
* `requestId` -	The ID of the request

See detailed information in [Register Image](http://docs.outscale.com/api_fcu/operations/Action_RegisterImage_get.html#_api_fcu-action_registerimage_get).
