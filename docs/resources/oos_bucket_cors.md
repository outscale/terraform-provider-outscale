---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_bucket_cors"
subcategory: "OOS API"
sidebar_current: "outscale-oos-bucket-cors"
description: |-
  [Manages a bucket CORS configuration.]
---

# outscale_oos_bucket_cors Resource

Manages a bucket CORS (cross-origin resource sharing) configuration.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/Introduction-to-OOS.html).  
For more information on this resource actions, see this [User Guide's section](https://docs.outscale.com/en/userguide/Managing-CORS-on-Your-Buckets.html).

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

* `bucket` - (Required) The name of the bucket to which you want to apply the CORS configuration.
* `rule` - (Required) A rule in the CORS configuration.
    * `allowed_headers` - (Required) The headers specified in the Access-Control-Request-Headers header. These headers are allowed in preflight OPTIONS requests. In response to preflight OPTIONS requests, allowed requested headers are returned.
    * `allowed_methods` - (Required) One or more HTTP methods (`GET` | `PUT` | `HEAD` | `POST` | `DELETE`) that you allow the origin to execute.
    * `allowed_origins` - (Required) One or more origins you want to be able to access the bucket from.
    * `expose_headers` - (Optional) One or more headers in the response that you want to make accessible.
    * `max_age_seconds` - (Required) The duration in seconds that the preflight response for the specified source can be cached.

## Attribute Reference

The following attributes are exported:

* `id` - The Terraform resource identifier, equal to the name of the bucket whose CORS configuration is managed.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 20 minutes.
* `read` - Defaults to 20 minutes.
* `update` - Defaults to 20 minutes.
* `delete` - Defaults to 60 minutes.
