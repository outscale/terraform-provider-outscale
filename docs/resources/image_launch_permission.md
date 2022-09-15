---
layout: "outscale"
page_title: "OUTSCALE: outscale_image_launch_permission"
sidebar_current: "outscale-image-launch-permission"
description: |-
  [Manages an image launch permission.]
---

# outscale_image_launch_permission Resource

Manages an image launch permission.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-OMIs.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#updateimage).

## Example Usage

### Add permissions

```hcl
resource "outscale_image_launch_permission" "image01" {
	image_id = "ami-12345678"
	permission_additions  {
		account_ids = ["012345678910"]
	}
}
```

### Remove permissions

```hcl
resource "outscale_image_launch_permission" "image02" {
	image_id = "ami-12345678"
	permission_removals  {
		account_ids = ["012345678910"]
	}
}
```

## Argument Reference

The following arguments are supported:

* `image_id` - (Required) The ID of the OMI you want to modify.
* `permission_additions` - (Optional) Information about the users to whom you want to give permissions for the resource.
    * `account_ids` - (Optional) The account ID of one or more users to whom you want to give permissions.
    * `global_permission` - (Optional) If true, the resource is public. If false, the resource is private.
* `permission_removals` - (Optional) Information about the users from whom you want to remove permissions for the resource.
    * `account_ids` - (Optional) The account ID of one or more users from whom you want to remove permissions.
    * `global_permission` - (Optional) If true, the resource is public. If false, the resource is private.

## Attribute Reference

The following attributes are exported:

* `description` - The description of the OMI.
* `image_id` - The ID of the OMI.
* `permissions_to_launch` - Information about the users who have permissions for the resource.
    * `account_ids` - The account ID of one or more users who have permissions for the resource.
    * `global_permission` - If true, the resource is public. If false, the resource is private.

