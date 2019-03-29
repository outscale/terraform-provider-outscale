---
layout: "outscale"
page_title: "OUTSCALE: outscale_users"
sidebar_current: "docs-outscale-datasource-users"
description: |-
  Lists all users that have a specified path prefix.
---

# outscale_users

Lists all users that have a specified path prefix.
If you do not specify a path prefix, this action returns a list of all users in the account (or an empty list if there are none).

## Example Usage

```hcl
resource "outscale_user" "user" {
    user_name = "test-user"
    path = "/"
}

data "outscale_users" "outscale_users" {
    path_prefix = "${outscale_user.user.path}"
}
```

## Argument Reference

The following arguments are supported:

* `path_prefix` - (Optional) The path prefix.

## Filters

None.

## Attributes Reference

The following attributes are exported:

* `users.N` - A list of users.
  + `arn` - The unique identifier of the resource (between 20 and 2048 characters).
  + `path` - The path to the user.
  + `user_id` - The ID of the user.
  + `user_name` - The name of the user.
* `request_id` - (Optional) The ID of the request.

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_ListUsers_get.html#_api_eim-action_listusers_get)
