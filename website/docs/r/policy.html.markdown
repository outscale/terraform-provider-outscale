---
layout: "outscale"
page_title: "OUTSCALE: outscale_policy"
sidebar_current: "docs-outscale-resource-policy"
description: |-
  Creates a new managed policy to apply to a role, a user or a group.
---

# outscale_policy

This action creates a policy version and sets v1 as the default one.

## Example Usage

```hcl
resource "outscale_policy" "policy" {
    path = "/"
    policy_name = "test-name"
    policy_document = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Action": [
        "ec2:Describe*"
      ],
      "Effect": "Allow",
      "Resource": "*"
    }
  ]
}
EOF
}
```

## Argument Reference

The following arguments are supported:

* `description` - (Optional) A friendly description of the policy, in which you can write information about the permissions contained in the policy.
* `path` - (Optional) The path to the policy.
* `policy_document` - The policy document, corresponding to a JSON string that contains the policy.
* `policy_name` - The name of the policy.

## Attributes Reference

The following attributes are exported:

* `arn` - The unique identifier of the resource (between 20 and 2048 characters).
* `attachment_count` - The number of resources attached to the policy.
* `default_version_id` - The ID of the policy default version.
* `description` - A friendly name for the policy (between 0 and 1000 characters).
* `is_attachable` - Indicates whether the policy can be attached to a group, a role or an EIM user.
* `path` - The path to the policy.
* `policy_id` - The ID of the policy (between 16 and 32 characters).
* `policy_name` - The name of the policy (between 1 and 128 characters).
* `request_id` - The ID of the request.


[See detailed description](http://docs.outscale.com/api_eim/operations/Action_CreatePolicy_get.html#_api_eim-action_createpolicy_get)
