---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_bucket_lifecycle"
subcategory: "OOS API"
sidebar_current: "outscale-oos-bucket-lifecycle"
description: |-
  Manages a bucket lifecycle configuration.
---

# outscale_oos_bucket_lifecycle Resource

Manages a bucket lifecycle configuration.

~> **Important** The deletion of expired objects via lifecycle configuration is asynchronous. This means a delay is to be expected between when the object is marked as expired and when it is effectively removed from the bucket. The actual removal process happens in the background and is managed by OUTSCALE Object Storage (OOS). You are no longer charged for storing objects once they have expired.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/Introduction-to-OOS.html).  
For more information on this resource actions, see this [User Guide's section](https://docs.outscale.com/en/userguide/Managing-the-Lifecycle-of-Your-Buckets.html).

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

For more examples of `filter` usage, see the [Object filtering](#object-filtering) section below.

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket to which you want to add the lifecycle configuration.
* `rule` - (Required) A rule in the lifecycle configuration.
    * `abort_incomplete_multipart_upload` - (Optional) Specifies when an incomplete multipart upload is aborted.
        * `days_after_initiation` - (Optional) The number of days after which an incomplete multipart upload is aborted.
    * `expiration` - (Optional) Specifies the expiration date for the objects.
        * `date` - (Optional) The lifetime of the objects, using a timestamp. The date must be in ISO 8601 (for example `2050-01-01T00:00:00Z`). The action is applied once the specified date is reached.<br />- The time is always midnight UTC. Specifying a different time (for example `2050-01-01T16:00:00Z`) will not return an error, but the action will still automatically be applied at midnight UTC on the specified date.<br />- You can specify a lifecycle action and set it to a past date. This will immediately qualify all relevant objects for that action.
        * `days` - (Optional) The lifetime of the objects, in days.<br />- When you create a lifecycle rule, the specified number of days is added to the creation time of your object and then rounded to the nearest midnight UTC. For example, if you create an object on February 1st, 2024 at 06:00 UTC and set a 1-day lifecycle rule, the object will expire on February 3rd, 2024 at exactly midnight UTC.
        * `expired_object_delete_marker` - (Optional) If `true`, the expired delete markers of an object are deleted when no older versions of the object are attached to them. This ensures that you do not accumulate unnecessary delete markers in a bucket.
    * `filter` - (Required) Specifies the objects the rule applies to, using `and`, `prefix`, `tags` keys. See the [Object filtering](#object-filtering) section below for examples of usage.
    * `id` - (Optional, Computed) A unique name for the rule. If omitted, the provider generates a unique name for it.
    * `noncurrent_version_expiration` - (Optional) If specified, the noncurrent versions of your objects are automatically deleted once they reach the duration specified in their lifecycle.
        * `noncurrent_days` - (Required) The lifetime, in days, of noncurrent objects.
    * `status` - (Required) If `Enabled`, the rule is applied. If `Disabled`, the rule is not applied.

## Object filtering

~> **Important:** If a lifecycle rule uses tag filtering, you will not be able to define the `expired_object_delete_marker` and `abort_incomplete_multipart_upload` rules.

The objects are identified with a prefix and/or one or more tags.

To have no filtering, use an empty filter:

```hcl
    filter = {}
```

To filter with a prefix, use the following structure:

```hcl
    filter = {
      prefix = "testprefix/"
    }
```

To filter with one or more tags, use the following structure, with an empty prefix string:
```hcl
    filter = {
      and = {
        prefix = ""
        tags = {
          testkey = "testvalue"
        }
      }
    }
```

To filter with a prefix and one or more tags, use the following structure:
```hcl
    filter = {
      and = {
        prefix = "testprefix/"
        tags = {
          testkey = "testvalue"
        }
      }
    }
```

## Attribute Reference

The following attributes are exported:

* `id` - The Terraform resource identifier, equal to the name of the bucket whose lifecycle configuration is managed.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 20 minutes.
* `read` - Defaults to 20 minutes.
* `update` - Defaults to 20 minutes.
* `delete` - Defaults to 60 minutes.
