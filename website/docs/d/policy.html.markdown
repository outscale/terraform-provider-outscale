---
layout: "outscale"
page_title: "OUTSCALE: outscale_policy"
sidebar_current: "docs-outscale-datasource-policy"
description: |-
  Retrieves information about a specified managed policy (default version, number of roles, users or groups the policy is attached to).
---

# outscale_policy

Retrieves information about a specified managed policy (default version, number of roles, users or groups the policy is attached to).

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

data "outscale_policy" "policy_ds" {
    policy_arn = "${outscale_policy.policy.arn}"
}
```

## Argument Reference

The following arguments are supported:

* `policy_arn` - The unique ressource identifier for the policy.

## Filters

None.

## Attributes Reference

The following attributes are exported:

* `arn` - (Optional) The unique identifier of the resource (between 20 and 2048 characters).
* `attachment_count` - (Optional) The number of resources attached to the policy.
* `default_version_id` - (Optional) The ID of the policy default version.
* `description` - (Optional) A friendly name for the policy (between 0 and 1000 characters).
* `is_attachable` - (Optional) Indicates whether the policy can be attached to a group, a role or an EIM user.
* `path` - (Optional) The path to the policy.
* `policy_id` - (Optional) The ID of the policy (between 16 and 32 characters).
* `policy_name` - (Optional) The name of the policy (between 1 and 128 characters).
* `request_id` - (Optional) The ID of the request.

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_GetPolicy_get.html#_api_eim-action_getpolicy_get)
