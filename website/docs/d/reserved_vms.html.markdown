---
layout: "outscale"
page_title: "OUTSCALE: outscale_reserved_vms"
sidebar_current: "docs-outscale-datasource-vms"
description: |-
  Describes Reserved Instances page
---

# outscale_reserved_vms

Describes Reserved Instances page

## Example Usage

```hcl
data "outscale_reserved_vms" "test" {
    filter {
		name = "availability-zone"
		values = ["${outscale_reserved_vms_oferr.availability_zone}"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional)	The Availability Zone where the reserved instances can be launched.		 
* `filter` - (Optional)	One or more filters.	  
* `offering_type` - (Optional)	The type of reserved instances offering (always One Shot).		 
* `reserved_instances_id`	The ID of one or more reservations for your account.

See detailed information in [Outscale Reserved VMS](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Filters

Use the Filter.N parameter to filter the described Reserved VMS on the following properties:

* `availability-zone`: -	The Availability Zone for the reservation. This filter is similar to the AvailabilityZone parameter.	
* `duration`: -	The duration of the reservation, expressed in seconds (2592000 seconds (1 month) | 31536000 seconds (1 year) | 94608000 seconds (3 years)).	
* `end`: -	The time when the reservation expires (for example, 2020-02-09T08:01:42.000Z).	
* `fixed-price`: -	The purchase price of the reserved instance. The currency depends on the Region of your account.	
* `instance-type`: -	The instance type of the reserved instance.	
* `product-description`: -	The product type of the reserved instance (Linux/UNIX | Windows | MapR).	
* `reserved-instances-id`: -	The ID of the reservation.	
* `start`: -	The time when you purchased the reservation (for example, 2016-02-09T08:01:42.000Z).
* `state`: -	The state of the reservation (active | retired).
* `usage-price`: -	The usage price of the reserved instance, per hour. The currency depends on the Region of your account.	


## Attributes Reference

The following attributes are exported:

* `request_id` -	The ID of the request.
* `reserved_instances_set`	Information about the reserved instances offerings.

See detailed information in [Describe Reserved VMS](http://docs.outscale.com/api_fcu/operations/Action_DescribeReservedInstances_get.html#_api_fcu-action_describereservedinstances_get).