---
layout: "outscale"
page_title: "OUTSCALE: outscale_lin_api_access_services"
sidebar_current: "docs-outscale-datasource-lin-api-access-services"
description: |-
  Describes Outscale services available to create VPC endpoints.
---

# outscale_lin_api_access_services

Describes Outscale services available to create VPC endpoints.

## Example Usage

```hcl
data "outscale_lin_api_access_services" "test" {
}
```

## Argument Reference

None

## Attributes Reference

The following attributes are exported:

* `service_name_set.N` - The names of the services you can use for VPC endpoints.

[See detailed information](http://docs.outscale.com/api_fcu/operations/Action_DescribeNetworkInterfaces_get.html#_api_fcu-action_describenetworkinterfaces_get).
