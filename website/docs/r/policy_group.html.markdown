---
layout: "outscale"
page_title: "OUTSCALE: outscale_policy_group"
sidebar_current: "docs-outscale-datasource-policy-group"
description: |-
  Creates or updates an inline policy included in a specified group.
---

# outscale_policy_group

Creates or updates an inline policy included in a specified group.

The policy is automatically applied to the all the users of the group after its creation.

## Example Usage

```hcl
resource "outscale_group" "group" {
    group_name = "test_group_1"
    path = "/"
}

resource "outscale_policy_group" "foo" {
    policy_name = "foo_policy_1"
    group_name = "${outscale_group.group.group_name}"
    policy_document = <<EOF
{
    "Version": "2012-10-17",
    "Statement": {
        "Effect": "Allow",
        "Action": "*",
        "Resource": "*"
    }
}
EOF
}
```

## Argument Reference

The following arguments are supported:

* `group_name` - The name of the group.
* `policy_document` - The policy document, providing a description of the policy as a json string.
* `policy_name` - The name of the policy.

## Attributes Reference

The following attributes are exported:

* `group_name` - The name of the group.
* `policy_document` - The policy document, providing a description of the policy as a json string.
* `policy_name` - The name of the policy.
* `request_id` - The ID of the request

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_PutGroupPolicy_get.html#_api_eim-action_putgrouppolicy_get)
