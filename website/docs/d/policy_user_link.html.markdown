---
layout: "outscale"
page_title: "OUTSCALE: outscale_policy_user_link"
sidebar_current: "docs-outscale-datasource-policy-user-link"
description: |-
  Attaches a managed policy to a specific user.
---

# outscale_policy_user_link

Attaches a managed policy to a specific user.

## Example Usage

```hcl
resource "outscale_user" "user" {
    user_name = "test-user-1"
}

resource "outscale_policy" "policy" {
    policy_name = "policy2"
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

data "outscale_policy_user_link" "outscale_policy_user_link" {
    user_name = "${outscale_user.user.user_name}"
}
```

## Argument Reference

The following arguments are supported:

* `path_prefix` - (Optional) The path prefix of the policies, set to a slash   if not specified..
* `user_name` - The name of the user.

## Filters

None.

## Attributes Reference

The following attributes are exported:

* `attached_policies.N` - One or more policies attached to the specified user.
  + `policy_arn` - The unique resource identifier for the policy (between 20 and 2048 characters).
  + `policy_name` - The name of the attached policy (between 1 and 128 characters).
* `request_id` - (Optional) The ID of the request.

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_ListAttachedUserPolicies_get.html#_api_eim-action_listattacheduserpolicies_get)
