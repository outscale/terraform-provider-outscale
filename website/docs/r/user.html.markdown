---
layout: "outscale"
page_title: "OUTSCALE: outscale_user"
sidebar_current: "docs-outscale-resource-user"
description: |-
  Creates a new user for your account.
---

# outscale_user

Creates a new user for your account.

## Example Usage

```hcl
resource "outscale_user" "user" {
    user_name = "test-user"
    path = "/"
}
```

## Argument Reference

The following arguments are supported:

* `path` - (Optional) The path for the user name. If you do not specify a path, it is set to a slash .
* `user_name` - The name of the user to be created.

## Attributes Reference

The following attributes are exported:

* `arn` - The unique identifier of the resource (between 20 and 2048 characters).
* `path` - The path to the user.
* `user_id` - The ID of the user.
* `user_name` - The name of the user.
* `request_id` - (Optional) The ID of the request.

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_CreateUser_get.html#_api_eim-action_createuser_get)
