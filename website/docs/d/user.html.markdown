---
layout: "outscale"
page_title: "OUTSCALE: outscale_user"
sidebar_current: "docs-outscale-datasource-user"
description: |-
  Gets information about a specified user, including its creation date, path, unique ID and Outscale Resource Name (ORN)
---

# outscale_user

Gets information about a specified user, including its creation date, path, unique ID and Outscale Resource Name (ORN).
If you do not specify a user name, this action returns information about the user who sent the request.

## Example Usage

```hcl
resource "outscale_user" "user" {
    user_name = "test-user"
    path = "/"
}

data "outscale_user" "outscale_user" {
    user_name = "${outscale_user.user.user_name}"
}
```

## Argument Reference

The following arguments are supported:

* `user_name` - (Optional) The name of the user.

## Filters

None.

## Attributes Reference

The following attributes are exported:

* `arn` - The unique identifier of the resource (between 20 and 2048 characters).
* `path` - The path to the user.
* `user_id` - The ID of the user.
* `user_name` - The name of the user.
* `request_id` - (Optional) The ID of the request.

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_GetUser_get.html#_api_eim-action_getuser_get)
