---
layout: "outscale"
page_title: "OUTSCALE: outscale_group_user"
sidebar_current: "docs-outscale-resource-group-user"
description: |-
  Adds a user to a specified group.
---

# outscale_group_user

Adds a user to a specified group.

## Example Usage

```hcl
resource "outscale_group" "group" {
    group_name = "test-group-1"
    path = "/"
}

resource "outscale_user" "user" {
    user_name = "test-user-1"
    path = "/"
}

resource "outscale_group_user" "team" {
    user_name = "${outscale_user.user.user_name}"
    group_name = "${outscale_group.group.group_name}"
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - The name of the group you want to add a user to.
* `user_name` - The name of the user you want to add to the group.

## Attributes Reference

The following attributes are exported:

* `groups,N` - One or more groups the specified user belongs to.
  + `arn` - The unique identifier of the group (between 20 and 2048 characters).
  + `group_id` - The ID of the group.
  + `group_name` - The name of the group.
  + `path` - The path to the group.
* `request_id` - (Optional) The ID of the request.

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_ListGroupsForUser_get.html#_api_eim-action_listgroupsforuser_get)
