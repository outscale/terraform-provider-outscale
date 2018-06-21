---
layout: "outscale"
page_title: "OUTSCALE: outscale_groups"
sidebar_current: "docs-outscale-datasource-groups"
description: |-
  Lists the groups with a specified path prefix.
---

# outscale_groups

Lists the groups with a specified path prefix.
If you do not specify any path prefix, this action returns a list of all groups (or an empty list if there are none).

## Example Usage

```hcl
resource "outscale_group" "group" {
    group_name = "test-group-1"
    path = "/"
}

data "outscale_groups" "groups_ds" {
    path_prefix = "/"
}
```

## Argument Reference

The following arguments are supported:

* `path_prefix` - The path prefix, set to a slash  if not specified.

## Filters

None.

## Attributes Reference

The following attributes are exported:

* `group` - (Optional) Information about the group.
  + `group_name` - The name of the group.
  + `arn` - The unique identifier of the group (between 20 and 2048 characters).
  + `group_id` - The ID of the group.
  + `path` - The path to the group.
* `request_id` - (Optional) The ID of the request.

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_ListGroups_get.html#_api_eim-action_listgroups_get)
