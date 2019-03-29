---
layout: "outscale"
page_title: "OUTSCALE: outscale_policy_version"
sidebar_current: "docs-outscale-datasource-policy-version"
description: |-
  Gets information about a specified version of a managed policy.
---

# outscale_policy_version

Gets information about a specified version of a managed policy.

## Example Usage

```hcl
resource "outscale_policy" "outscale_policy" {
    path = "/"
    policy_name = "test-name-%s"
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

data "outscale_policy_version" "policy_version_ds" {
    policy_arn = "${outscale_policy.outscale_policy.arn}",
    version_id = "${outscale_policy_version.policy.id}",
}
```

## Argument Reference

The following arguments are supported:

* `policy_arn` - The unique identifier (ARN) of the policy.
* `version_id` - The ID of the version.

## Attributes Reference

The following attributes are exported:

* `document` - (Optional) The policy document as a json string.
* `is_default_version` - (Optional) If true, the version is the default one.
* `version_id` - (Optional) The ID of the version.
* `request_id` - (Optional) The ID of the request

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_GetPolicyVersion_get.html#_api_eim-action_getpolicyversion_get)
