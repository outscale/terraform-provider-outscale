---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_bucket_encryption"
subcategory: "OOS API"
sidebar_current: "outscale-oos-bucket-encryption"
description: |-
  Manages a bucket encryption configuration.
---

# outscale_oos_bucket_encryption Resource

Manages a bucket encryption configuration.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/Introduction-to-OOS.html).  
For more information on this resource actions, see this [User Guide's section](https://docs.outscale.com/en/userguide/Using-Bucket-Encryption.html).

## Example Usage

### Enable default bucket encryption

```hcl
resource "outscale_oos_bucket" "bucket01" {}

resource "outscale_oos_bucket_encryption" "encryption01" {
  bucket          = outscale_oos_bucket.bucket01.id
  encryption_type = "AES256"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) The name of the bucket for which you want to enable encryption.
* `encryption_type` - (Optional, Computed) (Default: `AES256`) The server-side encryption algorithm to use:<br />- `AES256` is the default encryption algorithm. This encryption algorithm uses symmetric encryption and a 256-bit key. It is currently fully supported and recommended for use.<br />- `aws:kms` currently corresponds to `AES256` and is used only if a custom cryptographic key is used. Since the cryptographic key for your OUTSCALE account is provided and managed by 3DS OUTSCALE, we recommend using `AES256` for the time being.

## Attribute Reference

The following attributes are exported:

* `bucket_key_enabled` - Whether the cryptographic key is used per bucket (`true`) or per account (always `false`). 3DS OUTSCALE provides a default main cryptographic key for your account to encrypt your data.
* `id` - The Terraform resource identifier, equal to the name of the bucket whose encryption configuration is managed.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 20 minutes.
* `read` - Defaults to 20 minutes.
* `update` - Defaults to 20 minutes.
* `delete` - Defaults to 60 minutes.
