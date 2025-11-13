---
layout: "outscale"
page_title: "OUTSCALE: outscale_user_groups"
subcategory: "User Group"
sidebar_current: "outscale-user-groups"
description: |-
  [Provides information about user groups.]
---

# outscale_user_groups Data Source

Provides information about user groups.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-EIM-Groups.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api.html#createusergroup).

## Example Usage

```hcl
data "outscale_user_groups" "usegroups01" {
    filter {
        name   = "user_group_ids"
        values = ["XXXXXXXXX","YYYYYYYYYY"]
    }
    filter {
        name   = "path_prefix"
        values = ["/"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - One or more filters.
    * `path_prefix` - (Optional) The path prefix of the groups. If not specified, it is set to a slash (`/`).
    * `user_group_ids` - (Optional) The IDs of the user groups.
* `first_item` - (Optional) The item starting the list of groups requested.
* `results_per_page` - (Optional) The maximum number of items that can be returned in a single response (by default, `100`).

## Attribute Reference

The following attributes are exported:

* `has_more_items` - If true, there are more items to return using the `first_item` parameter in a new request.
* `max_results_limit` - Indicates maximum results defined for the operation.
* `max_results_truncated` - If true, indicates whether requested page size is more than allowed.
* `user_groups` - A list of user groups.
    * `creation_date` - The date and time (UTC) of creation of the user group.
    * `last_modification_date` - The date and time (UTC) of the last modification of the user group.
    * `name` - The name of the user group.
    * `orn` - The Outscale Resource Name (ORN) of the user group. For more information, see [Resource Identifiers](https://docs.outscale.com/en/userguide/Resource-Identifiers.html).
    * `path` - The path to the user group.
    * `user_group_id` - The ID of the user group.
