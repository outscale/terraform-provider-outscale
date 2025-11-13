---
layout: "outscale"
page_title: "OUTSCALE: outscale_user_groups_per_user"
subcategory: "User Group"
sidebar_current: "outscale-user-groups-per-user"
description: |-
  [Provides information about  groups that a specified user belongs to.]
---

# outscale_user_groups_per_user Data Source

Provides information about  groups that a specified user belongs to.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-EIM-Groups.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api.html#readusergroupsperuser).

## Example Usage

```hcl
data "outscale_user_groups_per_user" "user_groups_per_user01" {
    user_name = "user_name"
    user_path = "/"
}
```

## Argument Reference

The following arguments are supported:

* `user_name` - (Required) The name of the user.
* `user_path` - (Optional) The path to the user (by default, `/`).

## Attribute Reference

The following attributes are exported:

* `creation_date` - The date and time (UTC) of creation of the user group.
* `last_modification_date` - The date and time (UTC) of the last modification of the user group.
* `name` - The name of the user group.
* `orn` - The Outscale Resource Name (ORN) of the user group. For more information, see [Resource Identifiers](https://docs.outscale.com/en/userguide/Resource-Identifiers.html).
* `path` - The path to the user group.
* `user_group_id` - The ID of the user group.
