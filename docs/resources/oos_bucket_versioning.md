---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_bucket_versioning"
subcategory: "OOS API"
sidebar_current: "outscale-oos-bucket-versioning"
description: |-
  [TODO]
---

# outscale_oos_bucket_versioning Resource

<!-- Generated template; complete TODOs. -->

TODO

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/TODO).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/TODO).

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

* `bucket` - (Required) TODO
* `status` - (Required) TODO

## Attribute Reference

The following attributes are exported:

* `id` - The Terraform resource identifier, equal to the name of the bucket whose versioning configuration is managed.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 20 minutes.
* `read` - Defaults to 20 minutes.
* `update` - Defaults to 20 minutes.
* `delete` - Defaults to 60 minutes.
