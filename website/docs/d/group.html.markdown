---
layout: "outscale"
page_title: "OUTSCALE: outscale_group"
sidebar_current: "docs-outscale-datasource-group"
description: |-
  Retrieves the list of users that are in a specified group or lists the groups a user belongs to.
---

# outscale_group

Retrieves the list of users that are in a specified group or lists the groups a user belongs to.

## Example Usage

```hcl
resource "outscale_group" "group" {
    group_name = "test-group-1"
    path = "/"
}
data "outscale_group" "group_ds" {
    group_name = "${outscale_group.group.group_name}"
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - The name of the group.

## Filters

None.

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

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_GetGroup_get.html#_api_eim-action_getgroup_get)
