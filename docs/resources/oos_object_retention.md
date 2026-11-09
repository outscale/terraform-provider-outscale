---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_object_retention"
subcategory: "OOS API"
sidebar_current: "outscale-oos-object-retention"
description: |-
  [TODO]
---

# outscale_oos_object_retention Resource

<!-- Generated template; complete TODOs. -->

TODO

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/TODO).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/TODO).

## Example Usage

### Apply retention to an object version

```hcl
resource "outscale_oos_bucket" "bucket01" {
  object_lock = {
    enabled = true
  }
}

resource "outscale_oos_object" "object01" {
  bucket  = outscale_oos_bucket.bucket01.id
  key     = "retained-object.txt"
  content = "Object content"
}

resource "outscale_oos_object_retention" "retention01" {
  bucket            = outscale_oos_object.object01.bucket
  key               = outscale_oos_object.object01.key
  version_id        = outscale_oos_object.object01.version_id
  mode              = "COMPLIANCE"
  retain_until_date = "YYYY-MM-DDTHH:MM:SSZ"
}
```

## Argument Reference

The following arguments are supported:

* `bucket` - (Required) TODO
* `key` - (Required) TODO
* `mode` - (Optional, Computed) (Default: `COMPLIANCE`) TODO
* `retain_until_date` - (Required) TODO
* `version_id` - (Optional) TODO

## Attribute Reference

The following attributes are exported:

* `id` - The Terraform resource identifier, in the form `bucket/key`, with `/version_id` when object is versioned.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 20 minutes.
* `read` - Defaults to 20 minutes.
* `update` - Defaults to 20 minutes.
