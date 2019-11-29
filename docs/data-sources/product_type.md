---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_product_type"
sidebar_current: "docs-outscale-datasource-product-type"
description: |-
  [Provides information about a specific product type.]
---

# outscale_product_type Data Source

Provides information about a specific product type.
For more information on this resource, see the [User Guide](?).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-producttype).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `filters` - One or more filters.
  * `product_type_ids` - (Optional) The IDs of the product types.

## Attribute Reference

The following attributes are exported:

* `product_types` - Information about one or more product types.
  * `description` - The description of the product type.
  * `product_type_id` - The ID of the product type.
  * `vendor` - The vendor of the product type.
