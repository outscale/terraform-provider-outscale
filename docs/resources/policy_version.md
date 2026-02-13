---
layout: "outscale"
page_title: "OUTSCALE: outscale_policy_version"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-policy-version"
description: |-
  [Manages a policy version.]
---

# outscale_policy_version Resource

Manages a policy version.

~> **Note** At creation, the initial version of a policy is set to 'V1' by default.


For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/Editing-Managed-Policies-Using-Policy-Versions.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api.html#createpolicyversion).

## Example Usage

```hcl
resource "outscale_policy_version" "Policy2-version-02" {
    policy_orn     = outscale_policy.policy-2.orn
    document       = file("policy.json")
    set_as_default = true
}
```

## Argument Reference

The following arguments are supported:

* `document` - (Required) The policy document, corresponding to a JSON string that contains the policy. This policy document can contain a maximum of 5120 non-whitespace characters. For more information, see [EIM Reference Information](https://docs.outscale.com/en/userguide/EIM-Reference-Information.html) and [EIM Policy Generator](https://docs.outscale.com/en/userguide/EIM-Policy-Generator.html).
* `policy_orn` - (Required) The OUTSCALE Resource Name (ORN) of the policy. For more information, see [Resource Identifiers](https://docs.outscale.com/en/userguide/Resource-Identifiers.html).
* `set_as_default` - (Optional) If set to true, the new policy version is set as the default version, meaning it becomes the active one. Otherwise, the new policy version is not actually active until the `default_version_id` is specified in the `outscale_user` or `outscale_user_group` resources.

## Attribute Reference

The following attributes are exported:

* `body` - The policy document, corresponding to a JSON string that contains the policy. For more information, see [EIM Reference Information](https://docs.outscale.com/en/userguide/EIM-Reference-Information.html) and [EIM Policy Generator](https://docs.outscale.com/en/userguide/EIM-Policy-Generator.html).
* `creation_date` - The date and time (UTC) at which the version was created.
* `default_version` - If true, the version is the default one.
* `version_id` - The ID of the version.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 10 minutes.
* `read` - Defaults to 5 minutes.
* `delete` - Defaults to 5 minutes.
