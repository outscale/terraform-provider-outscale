---
layout: "outscale"
page_title: "OUTSCALE: outscale_catalog"
sidebar_current: "docs-outscale-datasource-catalog"
description: |-
  Returns the price list of Outscale products and services for the current Region.
---

# outscale_catalog

Returns the price list of Outscale products and services for the current Region.

## Example Usage

```hcl
data "outscale_catalog" "test" {}
```

## Argument Reference

No arguments are supported

## Attributes Reference

The following attributes are exported:

* `catalog.N` - A list of access keys and their metadata.
  * `atrributes.N` - One or more catalog attributes (for example, currency).
    * `key` - The key of the catalog attribute.
    * `value` - The value of the catalog attribute.
  * `entries.N` - One or more catalog entries.
    * `atrributes.N` - One or more catalog attributes (for example, currency).
      * `key` - The key of the catalog attribute.
      * `value` - The value of the catalog attribute.
    * `key` - The identifier of the catalog entry.
    * `value` - The value of the catalog attribute.
    * `title` - The description of the catalog entry.
* `request_id` - The ID of the Request.

[See detailed description](http://docs.outscale.com/api_icu/operations/Action_ReadCatalog_get.html)
