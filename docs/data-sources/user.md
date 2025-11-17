---
layout: "outscale"
page_title: "OUTSCALE: outscale_user"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-user"
description: |-
  [Provides information about a user.]
---

# outscale_user Data Source

Provides information about a user.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-EIM-Users.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api.html#createuser).

## Example Usage

```hcl
data "outscale_user" "user01" {
    filter {
        name   = "user_ids"
        values = ["XXXXXXXXXXXXXXXX"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - One or more filters.
    * `user_ids` - (Optional) The IDs of the users.
* `first_item` - (Optional) The item starting the list of users requested.
* `results_per_page` - (Optional) The maximum number of items that can be returned in a single response (by default, `100`).

## Attribute Reference

The following attributes are exported:

* `creation_date` - The date and time (UTC) of creation of the EIM user.
* `has_more_items` - If true, there are more items to return using the `first_item` parameter in a new request.
* `last_modification_date` - The date and time (UTC) of the last modification of the EIM user.
* `max_results_limit` - Indicates maximum results defined for the operation.
* `max_results_truncated` - If true, indicates whether requested page size is more than allowed.
* `path` - The path to the EIM user.
* `user_email` - The email address of the EIM user.
* `user_id` - The ID of the EIM user.
* `user_name` - The name of the EIM user.
