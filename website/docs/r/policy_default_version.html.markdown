---
layout: "outscale"
page_title: "OUTSCALE: outscale_policy_default_version"
sidebar_current: "docs-outscale-resource-policy-default-version"
description: |-
  Sets a specified version of a managed policy as the default (operative) one.
---

# outscale_policy_default_version

Sets a specified version of a managed policy as the default (operative) one.

You can modify the default version of a policy at any time.

## Example Usage

```hcl
resource "outscale_policy" "outscale_policy" {
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

resource "outscale_policy_default_version" "outscale_policy_default_version" {
    policy_arn = "${outscale_policy.outscale_policy.arn}"
    version_id = "v1"
}
```

## Argument Reference

The following arguments are supported:

* `policy_arn` - The unique identifier (ARN) of the policy.
* `version_id` - The ID of the version.

## Attributes Reference

The following attributes are exported:

* `policy_arn` - The unique identifier (ARN) of the policy.
* `version_id` - The ID of the version.

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_SetDefaultPolicyVersion_get.html#_api_eim-action_setdefaultpolicyversion_get)
