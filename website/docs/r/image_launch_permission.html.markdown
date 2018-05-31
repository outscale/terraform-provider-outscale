---
layout: "outscale"
page_title: "OUTSCALE: outscale_image_launch_permission"
sidebar_current: "docs-outscale-resource-image-launch-permission"
description: |-
	Modifies the specified attribute of an Outscale machine image (OMI).
---

# outscale_image_launch_permission

Modifies the specified attribute of an Outscale machine image (OMI).
You can specify only one attribute at a time. You can modify the permissions to access the OMI by adding or removing user IDs or groups. You can share an OMI with a user that is in the same Region. The user can create a copy of the OMI you shared, obtaining all the rights for the copy of the OMI

## Example Usage

```hcl
resource "outscale_vm" "outscale_instance" {
    count = 1
    image_id                    = "ami-880caa66"
    instance_type               = "c4.large"
}

resource "outscale_image" "outscale_image" {
    name        = "terraform test-123"
    instance_id = "${outscale_vm.outscale_instance.id}"
		no_reboot   = "true"
}

resource "outscale_image_launch_permission" "outscale_image_launch_permission" {
    image_id    = "${outscale_image.outscale_image.image_id}"
    launch_permission {
        add {
            user_id = "520679080430"
			}
		remove {
            user_id = "520679080430"
            }
		}
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - The ID of the OMI.
* `launch_permission` - Permissions to access the OMI.

## Attributes

* `request_id` - The ID of tue request.

[See detailed information](http://docs.outscale.com/api_fcu/operations/Action_ModifyImageAttribute_get.html#_api_fcu-action_modifyimageattribute_get).
