---
layout: "outscale"
page_title: "OUTSCALE: outscale_policies"
subcategory: "Policy"
sidebar_current: "outscale-policies"
description: |-
  [Provides information about policies.]
---

# outscale_policies Data Source

Provides information about policies.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Policies.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api.html#3ds-outscale-api-policy).

## Example Usage

```hcl
data "outscale_policies" "user_policies" {
    filter {
        name   = "only_linked"
        values = [true]
    }
    filter {
        name   = "path_prefix"
        values = ["/"]
    }
    filter {
        name   = "scope"
        values = ["LOCAL"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filters` - One or more filters.
    * `only_linked` - (Optional) If set to true, lists only the policies attached to a user.
    * `path_prefix` - (Optional) The path prefix you can use to filter the results. If not specified, it is set to a slash (`/`).
    * `scope` - (Optional) The scope of the policies. A policy can either be created by Outscale (`OWS`), and therefore applies to all accounts, or be created by its users (`LOCAL`).
* `first_item` - (Optional) The item starting the list of policies requested.
* `results_per_page` - (Optional) The maximum number of items that can be returned in a single response (by default, `100`).

## Attribute Reference

The following attributes are exported:

* `has_more_items` - If true, there are more items to return using the `first_item` parameter in a new request.
* `max_results_limit` - Indicates maximum results defined for the operation.
* `max_results_truncated` - If true, indicates whether requested page size is more than allowed.
* `policies` - Information about one or more policies.
    * `creation_date` - The date and time (UTC) at which the policy was created.
    * `description` - A friendly name for the policy (between 0 and 1000 characters).
    * `is_linkable` - Indicates whether the policy can be linked to a group or an EIM user.
    * `last_modification_date` - The date and time (UTC) at which the policy was last modified.
    * `orn` - The OUTSCALE Resource Name (ORN) of the policy. For more information, see [Resource Identifiers](https://docs.outscale.com/en/userguide/Resource-Identifiers.html).
    * `path` - The path to the policy.
    * `policy_default_version_id` - The ID of the policy default version.
    * `policy_id` - The ID of the policy.
    * `policy_name` - The name of the policy.
    * `resources_count` - The number of resources attached to the policy.
