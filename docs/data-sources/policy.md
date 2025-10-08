---
layout: "outscale"
page_title: "OUTSCALE: outscale_policy"
subcategory: "Policy"
sidebar_current: "outscale-policy"
description: |-
  [Provides information about a policy.]
---

# outscale_policy Data Source

Provides information about a policy.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Policies.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api.html#3ds-outscale-api-policy).

## Example Usage

```hcl
data "outscale_policy" "user_policy01" {
    policy_orn = "orn:ows:idauth::012345678910:policy/example/example-user-policy"
}
```

## Argument Reference

The following arguments are supported:

* `policy_orn` - (Required) The OUTSCALE Resource Name (ORN) of the policy. For more information, see [Resource Identifiers](https://docs.outscale.com/en/userguide/Resource-Identifiers.html).

## Attribute Reference

The following attributes are exported:

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
