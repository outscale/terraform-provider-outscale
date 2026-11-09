---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_bucket"
subcategory: "OOS API"
sidebar_current: "outscale-oos-bucket"
description: |-
  [TODO]
---

# outscale_oos_bucket Resource

<!-- Generated template; complete TODOs. -->

TODO

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/TODO).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/TODO).

## Example Usage

### Create a private bucket

Omitting `name` lets the provider generate a unique bucket name.

```hcl
resource "outscale_oos_bucket" "bucket01" {
  acl = "private"
}
```

### Create a bucket with default object retention

```hcl
resource "outscale_oos_bucket" "bucket02" {
  object_lock = {
    enabled = true

    default_retention = {
      mode = "COMPLIANCE"
      days = 1
    }
  }
}
```

~> **Important:** Objects uploaded to this bucket are protected for one day. An unexpired retention period prevents their deletion, even when `force_delete` is enabled.

### Create a bucket with grant 

```hcl
resource "outscale_oos_bucket" "bucket03" {
  grant = {
    full_control = {
      ids = ["OOS_USER_ID"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `acl` - (Optional) TODO
* `force_delete` - (Optional, Computed) (Default: `false`) Whether to empty the bucket before Terraform deletes it. When `true`, all objects, object versions, and delete markers are permanently deleted, including those not managed by Terraform. This does not bypass an unexpired object lock.
* `grant` - (Optional) TODO
    * `full_control` - (Optional) TODO
        * `email_addresses` - (Optional) TODO
        * `ids` - (Optional) TODO
    * `read` - (Optional) TODO
        * `email_addresses` - (Optional) TODO
        * `ids` - (Optional) TODO
    * `read_acp` - (Optional) TODO
        * `email_addresses` - (Optional) TODO
        * `ids` - (Optional) TODO
    * `write` - (Optional) TODO
        * `email_addresses` - (Optional) TODO
        * `ids` - (Optional) TODO
    * `write_acp` - (Optional) TODO
        * `email_addresses` - (Optional) TODO
        * `ids` - (Optional) TODO
* `name` - (Optional, Computed) The bucket name. If omitted, the provider generates a unique name prefixed with `bucket-`.
* `object_lock` - (Optional, Computed) TODO
    * `default_retention` - (Optional) TODO
        * `days` - (Optional) TODO
        * `mode` - (Optional, Computed) (Default: `COMPLIANCE`) TODO
        * `years` - (Optional) TODO
    * `enabled` - (Optional, Computed) (Default: `false`) TODO

## Attribute Reference

The following attributes are exported:

* `creation_date` - TODO
* `id` - The Terraform resource identifier, equal to the bucket name.
* `location_name` - TODO
* `permissions` - TODO
    * `grants` - TODO
        * `grantee` - TODO
            * `display_name` - TODO
            * `email_address` - TODO
            * `id` - TODO
            * `type` - TODO
            * `uri` - TODO
        * `permission` - TODO
    * `owner` - TODO
        * `display_name` - TODO
        * `id` - TODO

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 20 minutes.
* `read` - Defaults to 20 minutes.
* `update` - Defaults to 20 minutes.
* `delete` - Defaults to 60 minutes.
