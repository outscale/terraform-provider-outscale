---
layout: "outscale"
page_title: "OUTSCALE: outscale_oos_object_copy"
subcategory: "OOS API"
sidebar_current: "outscale-oos-object-copy"
description: |-
  [TODO]
---

# outscale_oos_object_copy Resource

<!-- Generated template; complete TODOs. -->

TODO

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/TODO).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/TODO).

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

* `acl` - (Optional) TODO
* `bucket` - (Required) TODO
* `cache_control` - (Optional, Computed) TODO
* `content_disposition` - (Optional, Computed) TODO
* `content_encoding` - (Optional, Computed) TODO
* `content_language` - (Optional, Computed) TODO
* `content_type` - (Optional, Computed) TODO
* `copy_source` - (Required) TODO
    * `if_match` - (Optional) TODO
    * `if_modified_since` - (Optional) TODO
    * `if_none_match` - (Optional) TODO
    * `if_unmodified_since` - (Optional) TODO
    * `path` - (Required) TODO
* `encryption_type` - (Optional, Computed) TODO
* `expires` - (Optional) TODO
* `force_delete` - (Optional, Computed) (Default: `false`) Whether Terraform permanently deletes all versions and delete markers for the destination object key when destroying the resource. This includes versions not created by this resource and does not bypass an unexpired object lock.
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
* `metadata` - (Optional, Computed) TODO
* `metadata_directive` - (Optional) TODO
* `tagging_directive` - (Optional) TODO
* `tags` - (Optional, Computed) TODO

## Attribute Reference

The following attributes are exported:

* `content_length` - TODO
* `copy_source` - TODO
    * `version_id` - TODO
* `etag` - TODO
* `expiration` - TODO
* `id` - The Terraform resource identifier, equal to the destination object key.
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
