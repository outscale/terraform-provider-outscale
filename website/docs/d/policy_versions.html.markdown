---
layout: "outscale"
page_title: "OUTSCALE: outscale_policy_versions"
sidebar_current: "docs-outscale-datasource-policy-versions"
description: |-
  Lists information about all the policy versions of a specified managed policy.
---

# outscale_policy_versions

Lists information about all the policy versions of a specified managed policy.

## Example Usage

```hcl
resource "outscale_policy" "outscale_policy" {
    path = "/"
    policy_name = "test-name-2"
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

data "outscale_policy_versions" "policy_versions_ds" {
    policy_arn = "${outscale_policy.outscale_policy.arn}"
}
```

## Argument Reference

The following arguments are supported:

* `policy_arn` - The Outscale Resource Name (ORN) of the policy.

## Attributes Reference

The following attributes are exported:

* `versions.N` - A list of all the versions of the policy.
  + `document` - (Optional) The policy document as a json string.
  + `is_default_version` - (Optional) If true, the version is the default one.
  + `version_id` - (Optional) The ID of the version.
  * `request_id` - (Optional) The ID of the request

[See detailed description](http://docs.outscale.com/api_eim/operations/Action_ListPolicyVersions_get.html#_api_eim-action_listpolicyversions_get)
