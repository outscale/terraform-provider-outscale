---
layout: "outscale"
page_title: "OUTSCALE: outscale_account_consumption"
sidebar_current: "docs-outscale-datasource-account-consumption"
description: |-
  Displays information about the consumption of your account for each billable resource within the specified time period.
---

# outscale_account_consumption

Displays information about the consumption of your account for each billable resource within the specified time period.

## Example Usage

```hcl
data "outscale_account_consumption" "test" {
  from_date = "2018-02-01"
  to_date = "2018-07-01"
}
```

## Argument Reference

The following arguments are supported:

* `from_date` (Required) - The beginning of the time period, in ISO-8601 format with the date only (for example, 2017-06-14 or 2017-06-14T00:00:00Z).
* `to_date` (Required) - The end of the time period, in ISO-8601 format with the date only (for example, 2017-06-30 or 2017-06-30T00:00:00Z)

## Attributes Reference

The following attributes are exported:

* `entries.N` - Information about the resources consumed during the specified time period.
  * `category` - The category of resource (for example, network).
  * `operation` - The API call that triggered the resource consumption (for example, RunInstances or CreateVolume).
  * `service` - The service of the API call (TinaOS-FCU, TinaOS-LBU, TinaOS-OSU or TinaOS-DirectLink).
  * `title` - A description of the consumed resource.
  * `type` - The type of resource, depending on the API call.
  * `value` - The consumed amount for the resource. The unit depends on the resource type. For more information, see the Title element.
* `request_id` - The ID of the request.

[See detailed description](http://docs.outscale.com/api_icu/operations/Action_ReadConsumptionAccount_get.html#_api_icu-action_readconsumptionaccount_get)
