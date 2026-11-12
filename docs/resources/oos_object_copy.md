---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_object_copy"
subcategory: "OOS API"
sidebar_current: "outscale-oos-object-copy"
description: |-
  Manages an object copy.
---

# outscale_oos_object_copy Resource

Manages an object copy.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/Introduction-to-OOS.html).  
For more information on this resource actions, see this [User Guide's section](https://docs.outscale.com/en/userguide/Copying-an-Object-to-a-Bucket.html).

## Example Usage

### Copy an object and preserve its metadata and tags

```hcl
resource "outscale_oos_bucket" "bucket01" {}

resource "outscale_oos_object" "source01" {
  bucket       = outscale_oos_bucket.bucket01.id
  key          = "source.txt"
  content      = "Source object content"
  content_type = "text/plain"

  metadata = {
    owner = "terraform"
  }

  tags = {
    environment = "example"
  }
}

resource "outscale_oos_object_copy" "copy01" {
  bucket = outscale_oos_bucket.bucket01.id
  key    = "destination.txt"

  copy_source = {
    path = "${outscale_oos_bucket.bucket01.id}/${outscale_oos_object.source01.key}"
  }

  metadata_directive = "COPY"
  tagging_directive  = "COPY"
}
```

### Replace metadata and tags on the copy

```hcl
resource "outscale_oos_object_copy" "copy02" {
  bucket = outscale_oos_bucket.bucket01.id
  key    = "destination-with-metadata.txt"

  copy_source = {
    path = "${outscale_oos_bucket.bucket01.id}/${outscale_oos_object.source01.key}"
  }

  metadata_directive = "REPLACE"
  tagging_directive  = "REPLACE"
  content_type       = "text/plain"

  metadata = {
    owner = "copied-object"
  }

  tags = {
    environment = "copy"
  }
}
```

## Argument Reference

The following arguments are supported:

* `acl` - (Optional) The permissions you grant for your object (`private` | `public-read` | `public-read-write` | `authenticated-read`). The `acl` and `grant` arguments are mutually exclusive.
* `bucket` - (Required) The name of the bucket where you want the object to be copied.
* `cache_control` - (Optional) How you want the object to be handled in terms of caching (`max-age` | `max-stale` | `min-fresh` | `no-cache` | `no-store` | `no-transform` | `only-if-cached` | `stale-if-error`).
* `content_disposition` - (Optional) How you want the object to be displayed when accessed via a browser:<br />- `inline`: When possible, prompts the browser to display the content within the browser window itself. This is useful for images, PDFs, and other media types that browsers can easily render.<br />-`attachment`: Regardless of file type, prompts the browser to download the content instead of displaying it inline. The file is thus saved locally.<br />- `attachment; filename=<NAME_OF_DOWNLOADED_FILE>`: Regardless of file type, prompts the browser to download the content instead of displaying it inline. The file is thus saved locally with the specified filename.
* `content_encoding` - (Optional) The encoding format of the object (`gzip` | `compress` | `deflate` | `identity` | `br`).
* `content_language` - (Optional, Computed) The language the content of the object is in, in language code (ISO 639 format).
* `content_type` - (Optional, Computed) The MIME (Multipurpose Internet Mail Extensions) type of the object.<br />An inaccurately set or absent content-type option can cause objects to be misinterpreted or mishandled by browsers. As a result, you may encounter difficulties accessing or viewing your objects using your preferred browser.
* `copy_source` - (Required) Information about the source object you want to copy.
    * `if_match` - (Optional) Copies the object only if its ETag matches the ETag you specify.
    * `if_modified_since` - (Optional) Copies the object only if it has been modified after the date and time you specify, in UTC format.
    * `if_none_match` - (Optional) Copies the object only if its ETag does not match the ETag you specify.
    * `if_unmodified_since` - (Optional) Copies the object only if it has not been modified since the date and time you specify, in UTC format.
    * `path` - (Required) The path in OOS to the source object you want to copy.
* `encryption_type` - (Optional, Computed) The server-side encryption algorithm to use. (`AES256` | `aws:kms`). Encryption must be enabled on the bucket first via an [`oos_bucket_encryption`](https://registry.terraform.io/providers/outscale/outscale/latest/docs/resources/oos_bucket_encryption) resource.
* `expires` - (Optional) The date and time, in UTC format, at which you consider the object can no longer be cached, and is considered stale.
* `force_delete` - (Optional, Computed) (Default: `false`) Whether Terraform permanently deletes all versions and delete markers for the destination object key when destroying the resource. This includes versions not created by this resource and does not bypass an unexpired object lock.
* `grant` - (Optional) The permissions you grant for your object. The `acl` and `grant` arguments are mutually exclusive.
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
* `key` - (Required) The name you want to give to your object in the bucket.
* `metadata` - (Optional) A map of additional metadata entries you can freely specify.
* `metadata_directive` - (Optional) Whether you want the metadata to be copied from the source objects or replaced with the metadata provided in the `metadata` argument (`COPY` or `REPLACE`).
* `tagging_directive` - (Optional) Whether you want the tag set the to be copied from the source object or replaced with the tag-set provided in the `tags` argument (`COPY` or `REPLACE`).
* `tags` - (Optional) Tags you want to have on the object.

## Attribute Reference

The following attributes are exported:

* `content_length` - The size of the object, in bytes.
* `copy_source` - The path in OOS to the source object of the copy.
    * `version_id` - The version ID of the source object.
* `etag` - The ETag of the object.
* `expiration` - Information about the lifecycle rule applied to the object, if there is one.
    * `expiry-date` - The date and time at which the object is set to expire.
    * `rule-id` - The name given to the lifecycle rule.
* `id` - The Terraform resource identifier, equal to the object key.
* `is_latest` - Whether the version of the object is the latest or not (`true` | `false`).
* `last_modified` - The date and time when the object was uploaded or copied into the bucket.
* `permissions` - Permissions and owner information for the object.
    * `grants` - Information about a granted permission for the object.
        * `grantee` - A grantee.
            * `display_name` - The display name of the grantee.
            * `email_address` - The email address of the grantee.
            * `id` - The S3 user ID of the grantee.
            * `type` - If the permission is for an individual user, `CanonicalUser`. If the permission is for a group of users, `Group`.
            * `uri` - The grantee group (`AllUsers` | `AuthenticatedUsers`).
        * `permission` - The granted permission (`READ` | `WRITE` | `READ_ACP` | `WRITE_ACP` | `FULL_CONTROL`).
    * `owner` - Information about the owner of the object.
        * `display_name` - The OUTSCALE account ID of the owner of the object.
        * `id` - The S3 user ID of the owner of the object.
* `version_id` - The version ID of the object.

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 20 minutes.
* `read` - Defaults to 20 minutes.
* `update` - Defaults to 20 minutes.
* `delete` - Defaults to 60 minutes.
