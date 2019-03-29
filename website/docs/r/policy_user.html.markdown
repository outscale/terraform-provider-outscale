---
layout: "outscale"
page_title: "OUTSCALE: outscale_policy_user"
sidebar_current: "docs-outscale-resource-policy-user"
description: |-
  Creates or updates an inline policy included in a specified user.
---

# outscale_policy_user

Creates or updates an inline policy included in a specified user.

The policy is automatically attahed to the user after its creation.

## Example Usage

```hcl
resource "outscale_user" "user" {
    user_name = "test_user_1"
    path = "/"
}

resource "outscale_policy_user" "foo" {
    policy_name = "foo_policy_1"
    user_name = "${outscale_user.user.user_name}"
    policy_document = "{\"Version\":\"2012-10-17\",\"Statement\":\"Effect\":\"Allow\",\"Action\":\"*\",\"Resource\":\"*\"}}"

    depends_on = ["outscale_user.user"]
}
```

## Argument Reference

The following arguments are supported:

* `user_name` - The policy document, providing a description of the policy as a json string.
* `policy_document` - The name of the policy.
* `policy_name` - The name of the user.

## Attributes Reference

The following attributes are exported:

* `user_name` - The policy document, providing a description of the policy as a json string.
* `policy_document` - The name of the policy.
* `policy_name` - The name of the user.
* `request_id` - The ID of the request.

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_PutUserPolicy_get.html#_api_eim-action_putuserpolicy_get)
