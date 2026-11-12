---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_object_retention"
subcategory: "OOS API"
sidebar_current: "outscale-oos-object-retention"
description: |-
  Manages an object retention configuration.
---

# outscale_oos_object_retention Resource

Manages an object retention configuration (also known as specific object lock).

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/Introduction-to-OOS.html).  
For more information on this resource actions, see this [User Guide's section](https://docs.outscale.com/en/userguide/Configuring-a-Specific-Object-Lock-on-an-Object.html).

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

* `bucket` - (Required) The name of the bucket you want to apply retention to. The bucket must have object lock enabled.
* `key` - (Required) The name of the object you want to lock.
* `mode` - (Optional, Computed) (Default: `COMPLIANCE`) The retention mode for the object. If set to `COMPLIANCE`, this prevents any user, including the root user, from overwriting or deleting the object until the retention period expires. OOS currently does not support `GOVERNANCE` mode.
* `retain_until_date` - (Required) The date, in ISO 8601 format, until which the object will be locked. For example, if set to `2050-01-01T00:00:00Z`, the object will remain locked until the very start of January 1, 2050, at midnight UTC.<br />Keep in mind that in COMPLIANCE mode:<br />- You cannot bypass the retention period until its expiration. This means you will be billed for the entire retention period you set.<br />- The retention period cannot be earlier than the default object lock of the bucket.
* `version_id` - (Optional) The version ID of the object.

## Attribute Reference

The following attributes are exported:

* `id` - The Terraform resource identifier, in the form `bucket/key`, with `/version_id` when object is versioned.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 20 minutes.
* `read` - Defaults to 20 minutes.
* `update` - Defaults to 20 minutes.
