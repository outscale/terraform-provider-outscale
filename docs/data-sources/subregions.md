---
layout: "outscale"
page_title: "OUTSCALE: outscale_subregions"
sidebar_current: "outscale-subregions"
description: |-
  [Provides information about subregions.]
---

# outscale_subregions Data Source

Provides information about subregions.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Regions-Endpoints-and-Availability-Zones.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#readsubregions).

## Example Usage

### List a specific Subregion in the current Region

```hcl
data "outscale_subregions" "subregions01" {
  filter {
    name   = "subregion_names"
    values = ["eu-west-2a"]
  }
}
```

### List two specific Subregions in the current Region

```hcl
data "outscale_subregions" "subregions02" {
  filter {
    name   = "subregion_names"
    values = ["eu-west-2a", "eu-west-2b"]
  }
}
```
### List all accessible Subregions in the current Region

```hcl
data "outscale_subregions" "all-subregions" {
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `subregion_names` - (Optional) The names of the Subregions.

## Attribute Reference

The following attributes are exported:

* `subregions` - Information about one or more Subregions.
    * `region_name` - The name of the Region containing the Subregion.
    * `state` - The state of the Subregion (`available` \| `information` \| `impaired` \| `unavailable`).
    * `subregion_name` - The name of the Subregion.
