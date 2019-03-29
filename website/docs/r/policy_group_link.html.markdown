---
layout: "outscale"
page_title: "OUTSCALE: outscale_policy_group_link"
sidebar_current: "docs-outscale-resource-policy-group-link"
description: |-
  Attaches a managed policy to a specific group. This policy applies to all the users contained in this group.
---

# outscale_policy_group_link

Attaches a managed policy to a specific group. This policy applies to all the users contained in this group.

## Example Usage

```hcl
resource "outscale_group" "group" {
    group_name = "tf-acc-group-gpa-basic-1"
}

resource "outscale_policy" "policy" {
    policy_name = "tf-acc-policy-gpa-basic-1"
    policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "iam:ChangePassword"
      ],
      "Resource": "*",
      "Effect": "Allow"
    }
  ]
}
EOF
}

resource "outscale_policy_group_link" "test-attach" {
    group_name = "${outscale_group.group.group_name}"
    policy_arn = "${outscale_policy.policy.arn}"
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - The friendly name given to the group you want to attach the policy to (between 1 and 128 characters).
* `policy_arn` - The unique resource identifier for the policy (between 20 and 2048 characters).

## Attributes Reference

The following attributes are exported:

* `policy_name` - The name of the policy.
* `group_name` - The friendly name given to the group you want to attach the policy to (between 1 and 128 characters).
* `policy_arn` - The unique resource identifier for the policy (between 20 and 2048 characters).
* `request_id` - The ID of the request.

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_AttachGroupPolicy_get.html#_api_eim-action_attachgrouppolicy_get)
