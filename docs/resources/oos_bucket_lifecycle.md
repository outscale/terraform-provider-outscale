---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_bucket_lifecycle"
subcategory: "OOS API"
sidebar_current: "outscale-oos-bucket-lifecycle"
description: |-
  [TODO]
---

# outscale_oos_bucket_lifecycle Resource

<!-- Generated template; complete TODOs. -->

TODO

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/TODO).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/TODO).

## Example Usage

### Configure expiration for all objects

```hcl
resource "outscale_oos_bucket" "bucket01" {}

resource "outscale_oos_bucket_versioning" "versioning01" {
  bucket = outscale_oos_bucket.bucket01.id
  status = "Enabled"
}

resource "outscale_oos_bucket_lifecycle" "lifecycle01" {
  bucket = outscale_oos_bucket.bucket01.id

  depends_on = [outscale_oos_bucket_versioning.versioning01]

  rule {
    id     = "expire"
    status = "Enabled"

    expiration = {
      days = 2
    }

    filter = {}
  }
}
```

### Configure expiration with a prefix and tags filter

```hcl
resource "outscale_oos_bucket" "bucket02" {}

resource "outscale_oos_bucket_versioning" "versioning02" {
  bucket = outscale_oos_bucket.bucket02.id
  status = "Enabled"
}

resource "outscale_oos_bucket_lifecycle" "lifecycle02" {
  bucket = outscale_oos_bucket.bucket02.id

  depends_on = [outscale_oos_bucket_versioning.versioning02]

  rule {
    id     = "expire"
    status = "Enabled"

    expiration = {
      days = 2
    }

    filter = {
      and = {
        prefix = "logs/"
        tags = {
          test = "test"
        }
      }
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) TODO
* `rule` - (Required) TODO
    * `abort_incomplete_multipart_upload` - (Optional) TODO
        * `days_after_initiation` - (Optional) TODO
    * `expiration` - (Optional) TODO
        * `date` - (Optional) TODO
        * `days` - (Optional) TODO
        * `expired_object_delete_marker` - (Optional) TODO
    * `filter` - (Required) TODO
        * `and` - (Optional) TODO
            * `prefix` - (Required) TODO
            * `tags` - (Required) TODO
        * `prefix` - (Optional) TODO
    * `id` - (Optional, Computed) TODO
    * `noncurrent_version_expiration` - (Optional) TODO
        * `noncurrent_days` - (Required) TODO
    * `status` - (Required) TODO

## Attribute Reference

The following attributes are exported:

* `id` - The Terraform resource identifier, equal to the name of the bucket whose lifecycle configuration is managed.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 20 minutes.
* `read` - Defaults to 20 minutes.
* `update` - Defaults to 20 minutes.
* `delete` - Defaults to 60 minutes.
