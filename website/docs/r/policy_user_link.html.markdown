---
layout: "outscale"
page_title: "OUTSCALE: outscale_policy_user_link"
sidebar_current: "docs-outscale-resource-policy-user-link"
description: |-
  Attaches a managed policy to a specific user.
---

# outscale_policy_user_link

Attaches a managed policy to a specific user.

## Example Usage

```hcl
resource "outscale_user" "user" {
    user_name = "test-user-%s"
}

resource "outscale_policy" "policy" {
    policy_name = "%s"
    policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "eim:ChangePassword"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "outscale_policy_user_link" "test-attach" {
    user_name = "${outscale_user.user.user_name}"
    policy_arn = "${outscale_policy.policy.arn}"
}
```

## Argument Reference

The following arguments are supported:

* `policy_arn` - The unique resource identifier for the policy (between 20 and 2048 characters).
* `user_name` - The friendly name of the user you want to attach the policy to (between 1 and 64 characters).

## Attributes Reference

The following attributes are exported:

* `policy_document` - The policy document, providing a description of the policy
* `policy_arn` - The unique resource identifier for the policy (between 20 and 2048 characters).
* `user_name` - The friendly name of the user you want to attach the policy to (between 1 and 64 characters).
* `request_id` - (Optional) The ID of the request.

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_GetUserPolicy_get.html#_api_eim-action_getuserpolicy_get)
