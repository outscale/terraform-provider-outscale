---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_bucket_cors"
subcategory: "OOS API"
sidebar_current: "outscale-oos-bucket-cors"
description: |-
  [TODO]
---

# outscale_oos_bucket_cors Resource

<!-- Generated template; complete TODOs. -->

TODO

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/TODO).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/TODO).

## Example Usage

### Configure CORS rules on a bucket

```hcl
resource "outscale_oos_bucket" "bucket01" {}

resource "outscale_oos_bucket_cors" "cors01" {
  bucket = outscale_oos_bucket.bucket01.id

  rule {
    allowed_headers = ["Authorization", "Content-Type"]
    allowed_methods = ["GET", "PUT"]
    allowed_origins = ["https://www.example.com"]
    expose_headers  = ["Content-Type"]
    max_age_seconds = 600
  }

  rule {
    allowed_headers = ["*"]
    allowed_methods = ["HEAD"]
    allowed_origins = ["https://www.example.org"]
    max_age_seconds = 1200
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) TODO
* `rule` - (Required) TODO
    * `allowed_headers` - (Required) TODO
    * `allowed_methods` - (Required) TODO
    * `allowed_origins` - (Required) TODO
    * `expose_headers` - (Optional) TODO
    * `max_age_seconds` - (Required) TODO

## Attribute Reference

The following attributes are exported:

* `id` - The Terraform resource identifier, equal to the name of the bucket whose CORS configuration is managed.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 20 minutes.
* `read` - Defaults to 20 minutes.
* `update` - Defaults to 20 minutes.
* `delete` - Defaults to 60 minutes.
