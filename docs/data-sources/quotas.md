---
layout: "outscale"
page_title: "3DS OUTSCALE: outscale_quota"
sidebar_current: "docs-outscale-datasource-quota"
description: |-
  [Provides information about quotas.]
---

# outscale_quota Data Source

Provides information about quotas.
For more information on this resource, see the [User Guide](https://wiki.outscale.net/display/EN/About+Your+3DS+OUTSCALE+Account#AboutYour3DSOUTSCALEAccount-quotasQuotasandConsumption).
For more information on this resource actions, see the [API documentation](https://docs-beta.outscale.com/#3ds-outscale-api-quota).

## Example Usage

```hcl
[exemple de code]
```

## Argument Reference

The following arguments are supported:

* `filters` - One or more filters.
  * `collections` - (Optional) The group names of the quotas.
  * `quota_names` - (Optional) The names of the quotas.
  * `quota_types` - (Optional) The resource IDs if these are resource-specific quotas, `global` if they are not.
  * `short_descriptions` - (Optional) The description of the quotas.

## Attribute Reference

The following attributes are exported:

* `quota_types` - Information about one or more quotas.
  * `quota_type` - The resource ID if it is a resource-specific quota, `global` if it is not.
  * `quotas` - One or more quotas associated with the user.
    * `account_id` - The account ID of the owner of the quotas.
    * `description` - The description of the quota.
    * `max_value` - The maximum value of the quota for the 3DS OUTSCALE user account (if there is no limit, `0`).
    * `name` - The unique name of the quota.
    * `quota_collection` - The group name of the quota.
    * `short_description` - The description of the quota.
    * `used_value` - The limit value currently used by the 3DS OUTSCALE user account.
