---
layout: "outscale"
page_title: "OUTSCALE: outscale_reserved_vms_oferr"
sidebar_current: "docs-outscale-datasource-vms-oferr"
description: |-
  Describes Reserved Instances Offerings
---

# outscale_reserved_vms_oferr

Describes Reserved Instances Offerings

## Example Usage

```hcl
data "outscale_reserved_vms_oferr" "test" {
    filter {
		name = "availability-zone"
		values = ["${outscale_reserved_vms_oferr.availability_zone}"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `availability_zone` - (Optional)	The Availability Zone where the reserved instances can be launched.		 
* `filter` - (Optional)	One or more filters.	false	Filter	 
* `instance_tenancy` - (Optional)	The tenancy of the reserved instances (default | dedicated).		 
* `instance_type` - (Optional)	The instance type of the reserved instances, including the custom instance types. For more information, see the DescribeInstanceTypes method.		 
* `offering_type` - (Optional)	The type of reserved instances offering (always One Shot).		 
* `product_description` - (Optional)	The product type of the reserved instance (Linux/UNIX | Windows | MapR).		 
* `reserved_instances_offering_id` - (Optional)	One or more offering IDs of reserved instances.		 

See detailed information in [Outscale VMS Oferr](https://wiki.outscale.net/display/DOCU/Getting+Information+About+Your+Instances).

## Filters

Use the Filter.N parameter to filter the described Reserved VMS Oferring on the following properties:

* `availability-zone`: -	The name of the Availability Zone where the reserved instance can be launched.
* `duration`: -	The duration of the reservation, in seconds (2592000 seconds (1 month) | 31536000 seconds (1 year) | 94608000 seconds (3 years)).
* `fixed-price`: -	The purchase price of the reserved instance. The currency depends on the Region of your account.
* `instance-type`: -	The instance type of the reserved instance.
* `product-description`: -	The product type of the reserved instance (Linux/UNIX | Windows | MapR).
* `reserved-instances-offering-id`: -	The offering ID of the reserved instance.
* `usage-price`: -	The usage price of the reserved instance, per hour. The currency depends on the Region of your account.


## Attributes Reference

The following attributes are exported:

* `request_id` -	The ID of the request
* `reserved_instances_offerings_set`	Information about the reserved instances offerings.

See detailed information in [Describe Volume]http://docs.outscale.com/api_fcu/operations/Action_DescribeVolumes_get.html#_api_fcu-action_describevolumes_get.