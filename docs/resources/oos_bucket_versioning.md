---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_bucket_versioning"
subcategory: "OOS API"
sidebar_current: "outscale-oos-bucket-versioning"
description: |-
  Manages a bucket versioning configuration.
---

# outscale_oos_bucket_versioning Resource

Manages a bucket versioning configuration.

!> **Warning** Versioning and encryption settings operate independently and are not retroactive. If a bucket contains multiple versions of the same object, any versions uploaded before encryption was enabled will not be encrypted automatically. Encryption only applies to objects or versions uploaded after encryption is enabled.<br /><br />To encrypt older versions, you must re-upload them or apply encryption manually.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/Introduction-to-OOS.html).  
For more information on this resource actions, see this [User Guide's section](https://docs.outscale.com/en/userguide/Using-Bucket-Versioning.html).

## Example Usage

### Enable bucket versioning

```hcl
resource "outscale_oos_bucket" "bucket01" {}

resource "outscale_oos_bucket_versioning" "versioning01" {
  bucket = outscale_oos_bucket.bucket01.id
  status = "Enabled"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket for which you want to enable or disable the versioning.
* `status` - (Required) The versioning state you want to apply to your bucket (`Enabled` | `Suspended`).

## Attribute Reference

The following attributes are exported:

* `id` - The Terraform resource identifier, equal to the name of the bucket whose versioning configuration is managed.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 20 minutes.
* `read` - Defaults to 20 minutes.
* `update` - Defaults to 20 minutes.
* `delete` - Defaults to 60 minutes.
