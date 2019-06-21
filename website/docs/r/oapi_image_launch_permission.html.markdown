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
    image_id  = "ami-3e158364"
    vm_type   = "t2.micro"
}
	
resource "outscale_image" "outscale_image" {
    image_name = "terraform test-123"
    vm_id      = "${outscale_vm.outscale_instance.id}"
    no_reboot  = "true"
}

resource "outscale_image_launch_permission" "outscale_image_launch_permission" {
    image_id    = "${outscale_image.outscale_image.image_id}"
    permission_additions {
		account_ids = ["201920914784"]
	}
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required) The ID of the OMI.
* `permission` - (Optional) Permissions to access the OMI.
    * `global_permission` - (Optional) 
    * `accounts_ids` - (Optional)

## Attributes

* `image_id` - The ID of the OMI.
* `description` - A description of the OMI.
* `permissions_to_launch` - One or more launch permissions.
    * `global_permission` - The name of the group (all if public).
    * `account_ids` - The accounts ID of the user.
* `request_id` - The ID of tue request.
