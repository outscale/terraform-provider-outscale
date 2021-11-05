---
layout: "outscale"
page_title: "OUTSCALE: outscale_regions"
sidebar_current: "outscale-regions"
description: |-
  [Provides information about Regions.]
---

# outscale_regions Data Source

Provides information about Regions.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Regions%2C+Endpoints%2C+and+Availability+Zones).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#readregions).

## Example Usage

```hcl
data "outscale_regions" "all_regions" {
  
}
```

## Argument Reference

No argument is supported.

## Attribute Reference

The following attributes are exported:

* `regions` - Information about one or more Regions.
    * `endpoint` - The hostname of the gateway to access the Region.
    * `region_name` - The administrative name of the Region.
