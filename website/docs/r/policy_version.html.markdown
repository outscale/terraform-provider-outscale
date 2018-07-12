---
layout: "outscale"
page_title: "OUTSCALE: outscale_policy_version"
sidebar_current: "docs-outscale-resource-policy-version"
description: |-
  Creates a new version of a specified managed policy.
---

# outscale_policy_version

Creates a new version of a specified managed policy.

A managed policy can have up to five versions.

## Example Usage

```hcl
resource "outscale_policy" "outscale_policy" {
    path = "/"
    policy_name = "test-name-1"
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

resource "outscale_policy_version" "policy" {
    policy_arn = "${outscale_policy.outscale_policy.arn}"
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

* `policy_arn` - The unique identifier (ARN) of the policy.
* `policy_document` - The policy document, providing a description of the policy as a json string.
* `set_as_default` - If set to true, the new policy version is the default version and becomes the operative one.

## Attributes Reference

The following attributes are exported:

* `document` - (Optional) The policy document as a json string.
* `is_default_version` - (Optional) If true, the version is the default one.
* `version_id` - (Optional) The ID of the version.
* `request_id` - (Optional) The ID of the request

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_CreatePolicyVersion_get.html#_api_eim-action_createpolicyversion_get)
