---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_image_launch_permission"
sidebar_current: "outscale-image-launch-permission"
description: |-
  [Manages an image launch permission.]
---

# outscale_image_launch_permission Resource

Manages an image launch permission.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+OMIs#AboutOMIs-OMIsPermissions,CopiesandExportstoOSU).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#updateimage).

## Example Usage

```hcl

# Add permissions

resource "outscale_image_launch_permission" "image01a" {
	image_id = "ami-12345678"
	permission_additions  {
		account_ids = ["012345678910"]
	}
}

# Remove permissions

resource "outscale_image_launch_permission" "image01b" {
	image_id = "ami-12345678"
	permission_removals  {
		account_ids = ["012345678910"]
	}
}


```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required) The ID of the OMI you want to modify.
* `permission_additions` - (Optional) Information about the users you want to give permissions for the resource.
  * `global_permission` - (Optional) If `true`, the resource is public. If `false`, the resource is private.
  * `account_ids` - (Optional) The account ID of one or more users you want to give permissions to.
* `permission_removals` - (Optional) Information about the users you want to remove permissions for the resource.
  * `global_permission` - (Optional) If `true`, the resource is public. If `false`, the resource is private.
  * `account_ids` - (Optional) The account ID of one or more users you want to remove permissions from.

## Attribute Reference

The following attributes are exported:

* `description` - A description of the OMI.
* `image_id` - The ID of the OMI you want to modify.
* `permission` - Information about the permissions for the resource.
  * `global_permission` - If `true`, the resource is public. If `false`, the resource is private.
  * `accounts_ids` - The account ID of one or more users who have permissions for the resource.
