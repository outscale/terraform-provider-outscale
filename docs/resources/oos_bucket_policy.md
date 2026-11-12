---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_bucket_policy"
subcategory: "OOS API"
sidebar_current: "outscale-oos-bucket-policy"
description: |-
  Manages a bucket policy.
---

# outscale_oos_bucket_policy Resource

Manages a bucket policy.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/Introduction-to-OOS.html).  
For more information on this resource actions, see this [User Guide's section](https://docs.outscale.com/en/userguide/Managing-a-Bucket-Policy.html).

## Example Usage

### Deny object downloads using a bucket policy

```hcl
resource "outscale_oos_bucket" "bucket01" {}

resource "outscale_oos_bucket_policy" "policy01" {
  bucket = outscale_oos_bucket.bucket01.id

  policy = jsonencode({
    Version = "2012-10-17"
    Statement = [{
      Effect    = "Deny"
      Principal = "*"
      Action    = "s3:GetObject"
      Resource  = "arn:aws:s3:::${outscale_oos_bucket.bucket01.id}/*"
    }]
  })
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket to which you want to apply the policy.
* `policy` - (Required) The bucket policy, encoded as a JSON string (use the [jsonencode](https://developer.hashicorp.com/terraform/language/functions/jsonencode) function to format it).
    * `Version` - (Required) The version of the policy language (use `2012-10-17`).
    * `Statement` - (Required) Defines the specific permissions you wish to grant regarding your resources (bucket and objects).
        * `Effect` - (Required) Specifies whether to `Allow` or `Deny` access to the resources defined in the policy.
        * `Principal` - (Required) Indicates which accounts are affected by the action (`"*"` for all, or `{"AWS": "arn:aws:iam::123456789012:root"}` for a specific account).
        * `Action` or `NotAction` - (Required) One or more operations that can or cannot be performed. For more information on operations you can specify in your policy, see [Bucket Policy Supported Actions List](https://docs.outscale.com/en/userguide/Bucket-Policy-Supported-Actions-List.html).
        * `Resource` - (Required) Specifies the bucket and objects to which the permissions apply, identified by Amazon Resource Names (ARNs).

## Attribute Reference

The following attributes are exported:

* `id` - The Terraform resource identifier, equal to the name of the bucket whose policy is managed.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 20 minutes.
* `read` - Defaults to 20 minutes.
* `update` - Defaults to 20 minutes.
* `delete` - Defaults to 60 minutes.
