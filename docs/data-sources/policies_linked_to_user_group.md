---
layout: "outscale"
page_title: "OUTSCALE: outscale_policies_linked_to_user_group"
subcategory: "Identity Access Management (IAM)"
sidebar_current: "outscale-policies-linked-to-user-group"
description: |-
  [Provides information about a policy linked to a user group.]
---

# outscale_policies_linked_to_user_group Data Source

Provides information about a link policy to user group.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Policies.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api.html#readmanagedpolicieslinkedtousergroup).

## Example Usage

```hcl
data "outscale_policies_linked_to_user_group" "managed_policies_linked_to_user_group" {
    user_group_name = "user_group_name"
    filter {
        name   = "user_group_ids"
        values = ["XXXXXXX","YYYYYY"]
    }
    filter {
        name   = "path_prefix"
        values = ["/test/"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - One or more filters.
    * `path_prefix` - (Optional) The path prefix of the groups. If not specified, it is set to a slash (`/`).
    * `user_group_ids` - (Optional) The IDs of the user groups.
* `first_item` - (Optional) The item starting the list of policies requested.
* `results_per_page` - (Optional) The maximum number of items that can be returned in a single response (by default, `100`).
* `user_group_name` - (Required) The name of the group.

## Attribute Reference

The following attributes are exported:

* `creation_date` - The date and time (UTC) at which the linked policy was created.
* `has_more_items` - If true, there are more items to return using the `first_item` parameter in a new request.
* `last_modification_date` - The date and time (UTC) at which the linked policy was last modified.
* `max_results_limit` - Indicates maximum results defined for the operation.
* `max_results_truncated` - If true, indicates whether requested page size is more than allowed.
* `orn` - The OUTSCALE Resource Name (ORN) of the policy. For more information, see [Resource Identifiers](https://docs.outscale.com/en/userguide/Resource-Identifiers.html).
* `policy_id` - The ID of the policy.
* `policy_name` - The name of the policy.
