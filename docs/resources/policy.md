---
layout: "outscale"
page_title: "OUTSCALE: outscale_policy"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-policy"
description: |-
  [Manages a policy.]
---

# outscale_policy Resource

Manages a policy.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Policies.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api.html#3ds-outscale-api-policy).

## Example Usage

 ```hcl
resource "outscale_policy" "policy-1"  {
    policy_name = "terraform-policy-1"
    description = "test-terraform"
    document    = file("policy.json")
    path        = "/"
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) A description for the policy.
* `document` - (Required) The policy document, corresponding to a JSON string that contains the policy. This policy document can contain a maximum of 5120 non-whitespace characters. For more information, see [EIM Reference Information](https://docs.outscale.com/en/userguide/EIM-Reference-Information.html) and [EIM Policy Generator](https://docs.outscale.com/en/userguide/EIM-Policy-Generator.html).
* `path` - (Optional) The path of the policy.
* `policy_name` - (Required) The name of the policy.

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

## Import

A policy can be imported using its ORN. For example:

```console

$ terraform import outscale_policy.policy1 orn

```