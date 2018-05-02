---
layout: "outscale"
page_title: "OUTSCALE: outcale_sub_regions"
sidebar_current: "docs-outscale-datasource-sub-regions"
description: |-
    Describes one or more Regions you can access.

---

# outscale_sub_regions

Describes one or more Regions you can access.

## Example Usage

```hcl
data "outscale_sub_regions" "by_name" {
	zone_name = ["eu-west-2a"]
}
```

## Argument Reference

The following arguments are supported:
	 
* `zone-name` -	The name of one or more Availability Zones.

See detailed information in [Outscale Sub Region](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Filters

Use the Filter.N parameter to filter the described Sub Regions on the following properties:

* `region-name`	The name of the Region containing the Availability Zones.
* `state`	The state of the Availability Zone (available | information | impaired | unavailable).	
* `zone-name`	The name of the Availability Zone. This filter is similar to the ZoneName.N parameter.	


## Attributes Reference

The following attributes are exported:

* `availability_zone_info` - Information about the Availability Zones.


See detailed information in [Describe Sub Regions](http://docs.outscale.com/api_fcu/operations/Action_DescribeRegions_get.html#_api_fcu-action_describeregions_get).
