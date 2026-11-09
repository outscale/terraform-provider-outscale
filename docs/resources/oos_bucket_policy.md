---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_bucket_policy"
subcategory: "OOS API"
sidebar_current: "outscale-oos-bucket-policy"
description: |-
  [TODO]
---

# outscale_oos_bucket_policy Resource

<!-- Generated template; complete TODOs. -->

TODO

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/TODO).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/TODO).

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

* `bucket` - (Required) TODO
* `policy` - (Required) TODO

## Attribute Reference

The following attributes are exported:

* `id` - The Terraform resource identifier, equal to the name of the bucket whose policy is managed.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 20 minutes.
* `read` - Defaults to 20 minutes.
* `update` - Defaults to 20 minutes.
* `delete` - Defaults to 60 minutes.
