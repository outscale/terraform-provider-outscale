---
layout: "outscale"
page_title: "OUTSCALE: outscale_image_copy"
sidebar_current: "docs-outscale-resource-image-copy"
description: |-
  Copies an Outscale Machine Image (OMI) to your account, from an account in the same Region.
---

# outscale_image_copy

Copies an Outscale Machine Image (OMI) to your account, from an account in the same Region.. [provisioning](/docs/provisioners/index.html).

## Example Usage

```hcl
resource "outscale_vm" "outscale_vm" {
    count = 1

    image_id                    = "ami-880caa66"
    instance_type               = "c4.large"
    key_name                    = "integ_sut_keypair"
    security_group              = ["sg-c73d3b6b"]

    provisioner "local-exec" {
        command = "date; who -b"
    }
}

resource "outscale_image" "outscale_image" {
    name        = "image_${outscale_vm.outscale_vm.id}"
    instance_id = "${outscale_vm.outscale_vm.id}"
    #no_reboot   = "false"                 # default value
}

resource "outscale_image_copy" "outscale_image_copy" {
    count = 1

    #source_image_id = ""
}
```

## Argument Reference

The following arguments are supported:
	

* `client_token` - (optional)	A unique identifier which enables you to manage the idempotency.
* `description` - (optional)	A description for the new OMI (if different from the source OMI one).
* `name` - (optional)	The name of the new OMI (if different from the source OMI one).
* `source_image_id` - (required)	The ID of the OMI you want to copy.
* `source_region` - (required)	The name of the source Region.


## Attributes Reference

* `image_id`	The ID of new OMI.
* `client_token`	A unique identifier which enables you to manage the idempotency.
* `description`	A description for the new OMI (if different from the source OMI one).
* `name`	The name of the new OMI (if different from the source OMI one).	
* `source_image_id`	The ID of the OMI you want to copy.	
* `source_region`	The name of the source Region.
* `request_id`	The ID of the request.

See detailed information in [Outscale Image Tasks](http://docs.outscale.com/api_fcu/operations/Action_CopyImage_get.html#_api_fcu-action_copyimage_get).

