---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_object"
subcategory: "OOS API"
sidebar_current: "outscale-oos-object"
description: |-
  Manages an object.
---

# outscale_oos_object Resource

Manages an object.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/Introduction-to-OOS.html).  
For more information on this resource actions, see this [User Guide's section](https://docs.outscale.com/en/userguide/Working-with-Objects.html).

## Example Usage

### Create an object with metadata and tags

```hcl
resource "outscale_oos_bucket" "bucket01" {}

resource "outscale_oos_object" "object01" {
  bucket        = outscale_oos_bucket.bucket01.id
  key           = "path/object.txt"
  content       = "Object content"
  acl           = "private"
  cache_control = "no-cache"
  content_type  = "text/plain"

  metadata = {
    owner = "terraform"
  }

  tags = {
    environment = "example"
  }
}
```

### Upload a local file

When using `source`, the provider automatically handles large files with [multipart upload](https://docs.outscale.com/en/userguide/Using-Multipart-Upload.html).

```hcl
resource "outscale_oos_object" "object02" {
  bucket       = outscale_oos_bucket.bucket01.id
  key          = "path/object.txt"
  source       = "object.txt"
  content_type = "text/plain"
}
```

### Create an object with grant

```hcl
resource "outscale_oos_object" "object04" {
  bucket  = outscale_oos_bucket.bucket01.id
  key     = "path/granted-object.txt"
  content = "Object with grant"

  grant = {
    read = {
      ids = ["S3_USER_ID"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `acl` - (Optional) The permissions you grant for your object (`private` | `public-read` | `public-read-write` | `authenticated-read`). The `acl` and `grant` arguments are mutually exclusive.
* `bucket` - (Required) The bucket in which you want to upload the object.<br />If encryption is enabled on a bucket, the new objects uploaded will automatically be encrypted based on the encryption configuration of the bucket. Encryption does not apply retroactively. For more information, see Enabling or Disabling Bucket Encryption.
* `cache_control` - (Optional) How you want the object to be handled in terms of caching (`max-age` | `max-stale` | `min-fresh` | `no-cache` | `no-store` | `no-transform` | `only-if-cached` | `stale-if-error`).
* `content` - (Optional) The object contents supplied directly as a non-empty UTF-8 string.
* `content_base64` - (Optional) The object contents encoded as a non-empty Base64 string.
* `content_disposition` - (Optional) How you want the object to be displayed when accessed via a browser:<br />- `inline`: When possible, prompts the browser to display the content within the browser window itself. This is useful for images, PDFs, and other media types that browsers can easily render.<br />-`attachment`: Regardless of file type, prompts the browser to download the content instead of displaying it inline. The file is thus saved locally.<br />- `attachment; filename=<NAME_OF_DOWNLOADED_FILE>`: Regardless of file type, prompts the browser to download the content instead of displaying it inline. The file is thus saved locally with the specified filename.
* `content_encoding` - (Optional) The encoding format of the object (`gzip` | `compress` | `deflate` | `identity` | `br`).
* `content_language` - (Optional, Computed) The language the content of the object is in, in language code (ISO 639 format).
* `content_type` - (Optional, Computed) The MIME (Multipurpose Internet Mail Extensions) type of the object.<br />An inaccurately set or absent content-type option can cause objects to be misinterpreted or mishandled by browsers. As a result, you may encounter difficulties accessing or viewing your objects using your preferred browser.
* `encryption_type` - (Optional) The server-side encryption algorithm to use. (`AES256` | `aws:kms`). Encryption must be enabled on the bucket first via an [`oos_bucket_encryption`](https://registry.terraform.io/providers/outscale/outscale/latest/docs/resources/oos_bucket_encryption) resource.
* `expires` - (Optional) The date and time, in UTC format, at which you consider the object can no longer be cached, and is considered stale.
* `force_delete` - (Optional, Computed) (Default: `false`) Whether Terraform permanently deletes all versions and delete markers for this object key when destroying the resource. This includes versions not created by this resource and does not bypass an unexpired object lock.
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
* `source` - (Optional) The path to a local file whose contents are uploaded as the object. When using this argument, the provider automatically handles large files with [multipart upload](https://docs.outscale.com/en/userguide/Using-Multipart-Upload.html).
* `tags` - (Optional) Tags you want to have on the object.

## Attribute Reference

The following attributes are exported:

* `content_length` - The size of the object, in bytes.
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
