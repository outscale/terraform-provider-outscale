---
layout: "outscale"
page_title: "OUTSCALE: outscale_reserved_vms_oferr_purchase"
sidebar_current: "docs-outscale-resource-reserved-vms-oferr-purchase"
description: |-
  Purchases reserved instances for your account in a specified Availability Zone (AZ). To guarantee your reservation, resource availability in the AZ is verified beforehand.
---

# outscale_reserved_vms_oferr_purchase

Purchases reserved instances for your account in a specified Availability Zone (AZ). To guarantee your reservation, resource availability in the AZ is verified beforehand.

## Example Usage

```hcl
		resource "outscale_reserved_vms_offer_purchase" "test" {
			instance_count = 1
			reserved_instances_offering_id = ""
		}

```

## Argument Reference

The following arguments are supported:

* `Instance_count` - (Required)	The number of reserved instances you want to purchase.
* `reserved_instances_offering_id` - (Required)	The ID of the reserved instances offering you want to purchase. 


## Attributes Reference

* `reserved_nstances_id` -	The ID of the reservation for your account.
* `reserved_instances_offering_id` -	The ID of the reserved instances offering you want to purchase. 
* `request_id` -	The ID of the request.


See detailed information in [Describe Reserved VMS Oferr Purchase](http://docs.outscale.com/api_fcu/operations/Action_PurchaseReservedInstancesOffering_get.html#_api_fcu-action_purchasereservedinstancesoffering_get).

