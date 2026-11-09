---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_object"
subcategory: "OOS API"
sidebar_current: "outscale-oos-object"
description: |-
  [TODO]
---

# outscale_oos_object Resource

<!-- Generated template; complete TODOs. -->

TODO

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/TODO).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/TODO).

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
      ids = ["OOS_USER_ID"]
    }
  }
}
```

## Argument Reference

The following arguments are supported:

* `acl` - (Optional) TODO
* `bucket` - (Required) TODO
* `cache_control` - (Optional) TODO
* `content` - (Optional) The object contents supplied directly as a non-empty UTF-8 string.
* `content_base64` - (Optional) The object contents encoded as a non-empty Base64 string.
* `content_disposition` - (Optional) TODO
* `content_encoding` - (Optional) TODO
* `content_language` - (Optional, Computed) TODO
* `content_type` - (Optional, Computed) TODO
* `encryption_type` - (Optional) TODO
* `expires` - (Optional) TODO
* `force_delete` - (Optional, Computed) (Default: `false`) Whether Terraform permanently deletes all versions and delete markers for this object key when destroying the resource. This includes versions not created by this resource and does not bypass an unexpired object lock.
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
* `key` - (Required) TODO
* `metadata` - (Optional) TODO
* `source` - (Optional) The path to a local file whose contents are uploaded as the object. 
* `tags` - (Optional) TODO

## Attribute Reference

The following attributes are exported:

* `content_length` - TODO
* `etag` - TODO
* `expiration` - TODO
* `id` - The Terraform resource identifier, equal to the object key.
* `is_latest` - TODO
* `last_modified` - TODO
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
* `version_id` - TODO

## Timeouts

The `timeouts` block enables you to configure [timeouts](https://developer.hashicorp.com/terraform/language/resources/configure#define-operation-timeouts) for certain actions:

* `create` - Defaults to 20 minutes.
* `read` - Defaults to 20 minutes.
* `update` - Defaults to 20 minutes.
* `delete` - Defaults to 60 minutes.
