---
layout: "outscale"
page_title: "OUTSCALE: outcale_regions"
sidebar_current: "docs-outscale-datasource-regions"
description: |-
    Describes available Regions.

---

# outscale_regions

Describes available Regions.

## Example Usage

```hcl
data "outscale_regions" "outscale_regions" {
    region_name = ["eu-west-2", "us-west-1", "us-east-2"]
}
```

## Argument Reference

The following arguments are supported:
	 
* `region_name` -	The name of one or more Regions.

See detailed information in [Outscale Regions](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Filters

Use the Filter.N parameter to filter the described Regions on the following properties:

* `region-name`: - The name of the Region. This filter is similar to the RegionName.N parameter.
* `endpoint`: -	The complete URL of the gateway to access the Region.	


## Attributes Reference

The following attributes are exported:

* `region_info`	Information about one or more Regions.	false	Region
* `request_id`	The ID of the request


See detailed information in [Describe Regions](http://docs.outscale.com/api_fcu/operations/Action_DescribeRegions_get.html#_api_fcu-action_describeregions_get).
