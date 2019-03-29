---
layout: "outscale"
page_title: "OUTSCALE: outscale_groups_for_user"
sidebar_current: "docs-outscale-datasource-groups-for-user"
description: |-
  Lists the groups a specified user belongs to.
---

# outscale_groups_for_user

Lists the groups a specified user belongs to.

## Example Usage

```hcl
resource "outscale_group" "group" {
    group_name = "test-group-1"
    path = "/"
}

data "outscale_groups_for_user" "groups_ds" {
    user_name = "test-group"
}
```

## Argument Reference

The following arguments are supported:

* `user_name` - The name of the user.

## Filters

None.

## Attributes Reference

The following attributes are exported:

* `groups.N` - (Optional) One or more groups the specified user belongs to.
  + `group_name` - The name of the group.
  + `arn` - The unique identifier of the group (between 20 and 2048 characters).
  + `group_id` - The ID of the group.
  + `path` - The path to the group.
* `request_id` - (Optional) The ID of the request.

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_ListGroupsForUser_get.html#_api_eim-action_listgroupsforuser_get)
