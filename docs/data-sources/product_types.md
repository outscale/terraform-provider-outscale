---
layout: "outscale"
page_title: "OUTSCALE: outscale_product_types"
sidebar_current: "outscale-product-types"
description: |-
  [Provides information about product types.]
---

# outscale_product_types Data Source

Provides information about product types.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/Software-Licenses.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-producttype).

## Example Usage

### Read specific product types
```hcl
data "outscale_product_types" "product_types01" {
  filter {
    name   = "product_type_ids"
    values = ["0001", "0002"]
  }    
}
```

### Read all product types
```hcl
data "outscale_product_types" "all_product_types" {
}
```


## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `product_type_ids` - (Optional) The IDs of the product types.

## Attribute Reference

The following attributes are exported:

* `product_types` - Information about one or more product types.
    * `description` - The description of the product type.
    * `product_type_id` - The ID of the product type.
    * `vendor` - The vendor of the product type.
