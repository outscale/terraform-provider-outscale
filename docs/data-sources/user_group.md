---
layout: "outscale"
page_title: "OUTSCALE: outscale_user_group"
subcategory: "User Group"
sidebar_current: "outscale-user-group"
description: |-
  [Provides information about a user group.]
---

# outscale_user_group Data Source

Provides information about a user group.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-EIM-Groups.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api.html#createusergroup).

## Example Usage

```hcl
data "outscale_user_group" "user_group01" {
   user_group_name = "user_group_name"
   path            = "/"
}
```

## Argument Reference

The following arguments are supported:

* `path` - (Optional) The path to the group. If not specified, it is set to a slash (`/`).
* `user_group_name` - (Required) The name of the group.

## Attribute Reference

The following attributes are exported:

* `creation_date` - The date and time (UTC) of creation of the EIM user.
* `last_modification_date` - The date and time (UTC) of the last modification of the EIM user.
* `name` - The name of the user group.
* `orn` - The Outscale Resource Name (ORN) of the user group. For more information, see [Resource Identifiers](https://docs.outscale.com/en/userguide/Resource-Identifiers.html).
* `path` - The path to the EIM user.
* `user_email` - The email address of the EIM user.
* `user_group_id` - The ID of the user group.
* `user_id` - The ID of the EIM user.
* `user_name` - The name of the EIM user.
