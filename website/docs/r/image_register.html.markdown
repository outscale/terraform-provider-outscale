---
layout: "outscale"
page_title: "OUTSCALE: outscale_image_register"
sidebar_current: "docs-outscale-resource-image-register"
description: |-
  Registers an Outscale Machine Image (OMI) to finalize its creation process.
---

# outscale_image_register

Registers an Outscale Machine Image (OMI) to finalize its creation process.. [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
resource "outscale_vm" "outscale_vm" {
    count = 1
    image_id                    = "ami-880caa66"
    instance_type               = "c4.large"
}
resource "outscale_image_register" "outscale_image_register" {
    name        = "image_%d"
    instance_id = "${outscale_vm.outscale_vm.id}"
}
```

## Argument Reference

The following arguments are supported:

* `architecture` - (Optional) The architecture of the OMI (set to i386 by default).
* `block_device_mapping` - (Optional) One or more Block Device Mapping entries.
* `description` - (Optional) A description for the OMI.
* `image_location` - (Optional) The pre-signed URL of the OMI manifest file, or the full path to the OMI stored in an OSU bucket. If you specify this parameter, a copy of the OMI is created in your account.
* `name` - (Required) A unique name for the OMI.
* `toot_device_name` - (Optional) The name of the root device (for example, /dev/sda1)

## Attributes Reference

* `imageId` - The ID of the newly registered OMI.
* `requestId` - The ID of the request.

See detailed information in [Outscale Image Tasks](http://docs.outscale.com/api_fcu/operations/Action_RegisterImage_get.html#_api_fcu-action_registerimage_get).
