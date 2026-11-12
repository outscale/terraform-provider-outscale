---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_bucket"
subcategory: "OOS API"
sidebar_current: "outscale-oos-bucket"
description: |-
  [Manages a bucket.]
---

# outscale_oos_bucket Resource

Manages a bucket.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/Introduction-to-OOS.html).  
For more information on this resource actions, see this [User Guide's section](https://docs.outscale.com/en/userguide/Managing-Your-Buckets.html).

## Example Usage

### Create a private bucket

Omitting `name` lets the provider generate a unique bucket name prefixed with `bucket-`.

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
      ids = ["S3_USER_ID"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `acl` - (Optional) The permissions you grant for your bucket (`private` | `public-read` | `public-read-write` | `authenticated-read`). The `acl` and `grant` arguments are mutually exclusive.
* `force_delete` - (Optional, Computed) (Default: `false`) Whether to empty the bucket before Terraform deletes it. When `true`, all objects, object versions, and delete markers are permanently deleted, including those not managed by Terraform. This does not bypass an unexpired object lock.
* `grant` - (Optional) The permissions you grant for your bucket. The `acl` and `grant` arguments are mutually exclusive.
    * `full_control` - (Optional) Users to whom you grant `full-control` permission.
        * `email_addresses` - (Optional) One or more user email addresses.
        * `ids` - (Optional) One or more S3 user IDs (not to be confused with OUTSCALE account IDs).
    * `read` - (Optional) Users to whom you grant `read` permission.
        * `email_addresses` - (Optional) One or more user email addresses.
        * `ids` - (Optional) One or more S3 user IDs (not to be confused with OUTSCALE account IDs).
    * `read_acp` - (Optional) Users to whom you grant `read_acp` permission.
        * `email_addresses` - (Optional) One or more user email addresses.
        * `ids` - (Optional) One or more S3 user IDs (not to be confused with OUTSCALE account IDs).
    * `write` - (Optional) Users to whom you grant `write` permission.
        * `email_addresses` - (Optional) One or more user email addresses.
        * `ids` - (Optional) One or more S3 user IDs (not to be confused with OUTSCALE account IDs).
    * `write_acp` - (Optional) Users to whom you grant `write_acp` permission.
        * `email_addresses` - (Optional) One or more user email addresses.
        * `ids` - (Optional) One or more S3 user IDs (not to be confused with OUTSCALE account IDs).
* `name` - (Optional, Computed) The bucket name. If omitted, the provider generates a unique name prefixed with `bucket-`.
* `object_lock` - (Optional, Computed) Object lock configuration for the bucket.
    * `default_retention` - (Optional) The default retention parameters for new objects placed in the bucket. The `days` and `years` arguments are mutually exclusive.
        * `days` - (Optional) The days after which your objects will no longer be protected. The minimum is `1`, and the expiration is calculated from the exact time each object is uploaded to the bucket. For example, if the default retention for your bucket is `1` day and you upload an object at 06:00 UTC, the object will be locked until the next day, at 06:00 UTC.
        * `mode` - (Optional, Computed) (Default: `COMPLIANCE`) The retention mode for the bucket. If set to `COMPLIANCE`, this prevents any user, including the root user, from overwriting or deleting objects within the bucket until the retention period expires. OOS currently does not support `GOVERNANCE` mode.
        * `years` - (Optional) The years after which your objects will no longer be protected.
    * `enabled` - (Optional, Computed) (Default: `false`) Whether object lock is enabled on the bucket.

## Attribute Reference

The following attributes are exported:

* `creation_date` - The date and time of creation of the bucket (for example, `2010-10-01 12:34:56.78 +0000 UTC`).
* `id` - The Terraform resource identifier, equal to the bucket name.
* `location_name` - Specifies the location of the bucket.
* `permissions` - Permissions and owner information for the bucket.
    * `grants` - Information about a granted permission for the bucket.
        * `grantee` - A grantee.
            * `display_name` - The display name of the grantee.
            * `email_address` - The email address of the grantee.
            * `id` - The S3 user ID of the grantee.
            * `type` - If the permission is for an individual user, `CanonicalUser`. If the permission is for a group of users, `Group`.
            * `uri` - The grantee group (`AllUsers` | `AuthenticatedUsers`).
        * `permission` - The granted permission (`READ` | `WRITE` | `READ_ACP` | `WRITE_ACP` | `FULL_CONTROL`).
    * `owner` - Information about the owner of the bucket.
        * `display_name` - The OUTSCALE account ID of the owner of the bucket.
        * `id` - The S3 user ID of the owner of the bucket.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 20 minutes.
* `read` - Defaults to 20 minutes.
* `update` - Defaults to 20 minutes.
* `delete` - Defaults to 60 minutes.
