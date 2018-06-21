---
layout: "outscale"
page_title: "OUTSCALE: outscale_group"
sidebar_current: "docs-outscale-resource-group"
description: |-
  Creates a group to which you can add EIM users.
  You can also add an inline policy or attach a managed policy to the group, which is applied to all its users.
---

# outscale_group

Creates a group to which you can add EIM users.
You can also add an inline policy or attach a managed policy to the group, which is applied to all its users.

## Example Usage

```hcl
resource "outscale_group" "group" {
    group_name = "test-group-1"
    path = "/"
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - The name of the group to be created.
* `path` - (Optional) The path to the group, set to a slash  if not specified.

## Attributes Reference

The following attributes are exported:

* `group` - Information about the group.
  + `arn` - The unique identifier of the group (between 20 and 2048 characters).
  + `group_id` - The ID of the group.
  + `group_name` - The name of the group.
  + `path` - The path to the group.
* `users.N` - (Optional) The list of users in the group.
  + `arn` - The unique identifier of the user (between 20 and 2048 characters).
  + `path` - The path to the user.
  + `user_id` - The ID of the user.
  + `user_name` - The name of the user.
* `request_id` - (Optional) The ID of the request.

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_CreateGroup_get.html#_api_eim-action_creategroup_get)
