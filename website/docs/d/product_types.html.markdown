---
layout: "outscale"
page_title: "OUTSCALE: outscale_product_types"
sidebar_current: "docs-outscale-datasource-product-types"
description: |-
Describes one or more product types.
---

# outscale_product_types

Describes one or more product types.

## Example Usage

```hcl
data "outscale_product_types" "test" {}
```

## Argument Reference

You can use the Filter.N parameter to filter the product types on the description property:

* `Filter.N` (Optional)One or more filters.

## Filters

Use the Filter.N parameter to filter the described instances on the following properties:

* `product_type_set` The type of the product.

## Attributes Reference

The following attributes are exported:

* `product_type_set.N` - Information about one or more product types, each containing the following attributes:
  - `description` - The description of the product type
  - `product_type_id` - The ID of the product type.
  - `vendor` - The vendor of the product type.
* `request_id` - The ID of the request.

[See detailed description](http://docs.outscale.com/api_fcu/operations/Action_DescribeProductTypes_get.html#_api_fcu-action_describeproducttypes_get)
