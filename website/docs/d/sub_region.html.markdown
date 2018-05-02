---
layout: "outscale"
page_title: "OUTSCALE: outcale_sub_region"
sidebar_current: "docs-outscale-datasource-sub-region"
description: |-
    Describes one or more Regions you can access.

---

# outscale_sub_region

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

Use the Filter.N parameter to filter the described Sub Region on the following properties:

* `region-name`	The name of the Region containing the Availability Zones.
* `state`	The state of the Availability Zone (available | information | impaired | unavailable).	
* `zone-name`	The name of the Availability Zone. This filter is similar to the ZoneName.N parameter.	


## Attributes Reference

The following attributes are exported:

* `region_name` - The name of the Region containing the Availability Zone.	false	string
* `zone_name` - The name of the Availability Zone.	false	string
* `zone_state` - The state of Availability Zone (always available if the user has access to the Region containing the Availability Zone).
* `request_d` - The ID of the rerquest


See detailed information in [Describe Sub Region](http://docs.outscale.com/api_fcu/operations/Action_DescribeRegions_get.html#_api_fcu-action_describeregions_get).
