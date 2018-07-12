---
layout: "outscale"
page_title: "OUTSCALE: outscale_policy_group_link"
sidebar_current: "docs-outscale-datasource-policy-group-link"
description: |-
  Lists the managed policies attached to a specified group.
---

# outscale_policy_group_link

Lists the managed policies attached to a specified group.

## Example Usage

```hcl
resource "outscale_group" "group" {
    group_name = "tf-acc-group-gpa-basic-1"
}

resource "outscale_policy" "policy" {
    policy_name = "tf-acc-policy-gpa-basic-2"
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

data "outscale_policy_group_link" "outscale_policy_group_link" {
    group_name = "${outscale_group.group.group_name}"
    path_prefix = "${outscale_policy.policy.arn}"
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - The name of the group.
* `path_prefix` - The path prefix of the policies, set to a slash (/) if not specified.

## Filters

None.

## Attributes Reference

The following attributes are exported:

* `attached_policies.N` - One or more policies attached to the specified group.
  * `policy_arn` - The unique resource identifier for the policy (between 20 and 2048 characters).
  * `policy_name` - The name of the attached policy (between 1 and 128 characters).
* `request_id` - The ID of the request.

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_ListAttachedGroupPolicies_get.html#_api_eim-action_listattachedgrouppolicies_get)
