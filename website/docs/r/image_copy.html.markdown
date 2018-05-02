---
layout: "outscale"
page_title: "OUTSCALE: outscale_image_copy"
sidebar_current: "docs-outscale-resource-image-copy"
description: |-
  Copies an Outscale machine image (OMI) to your account, from an account in the same Region.
---

# outscale_image_copy

Copies an Outscale machine image (OMI) to your account, from an account in the same Region.

## Example Usage

```hcl
resource "outscale_vm" "outscale_vm" {
    image_id                    = "ami-880caa66"
    instance_type               = "c4.large"

}

resource "outscale_image" "outscale_image" {
    name        = "image_${outscale_vm.outscale_vm.id}"
    instance_id = "${outscale_vm.outscale_vm.id}"
    #no_reboot   = "false"                 # default value
}

resource "outscale_image_copy" "test" {
		source_image_id = "${outscale_image.outscale_image.image_id}"
		source_region= "eu-west-2"
}
```

## Argument Reference

The following arguments are supported:

* `client_token` - (false) A unique identifier which enables you to manage the idempotency..
* `description` - (false) A description for the new OMI (if different from the source OMI one)..
* `name` - (false) The name of the new OMI (if different from the source OMI one)..
* `source_image_id` - (true) The ID of the OMI you want to copy..
* `source_region` - (true) The name of the source Region..

## Argument Reference

The following arguments are supported:

* `ClientToken` -	A unique identifier which enables you to manage the idempotency.	false	string
* `Description` -	A description for the new OMI (if different from the source OMI one).	false	string
* `Name` -	The name of the new OMI (if different from the source OMI one).	false	string
* `SourceImageId` -	The ID of the OMI you want to copy.	true	string
* `SourceRegion` -	The name of the source Region.	true	string
* `requestId` -	The ID of the request


See detailed information in [Create Copy Image](http://docs.outscale.com/api_fcu/operations/Action_CopyImage_get.html#_api_fcu-action_copyimage_get).

