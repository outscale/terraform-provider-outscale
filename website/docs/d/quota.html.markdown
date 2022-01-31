---
layout: "outscale"
page_title: "OUTSCALE: outscale_quota"
sidebar_current: "outscale-quota"
description: |-
  [Provides information about a specific quota.]
---

# outscale_quota Data Source

Provides information about a specific quota.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Your-Account.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#readquotas).

## Example Usage

```hcl
data "outscale_quota" "load_balancer_listeners_quota01" {
  filter {
    name   = "collections"
    values = ["LBU"]
  }
  filter {
    name   = "quota_names"
    values = ["lb_listeners_limit"]
  }
  filter {
    name   = "quota_types"
    values = ["global"]
  }
  filter {
    name   = "short_descriptions"
    values = ["Load Balancer Listeners Limit"]
  }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `collections` - (Optional) The group names of the quotas.
    * `quota_names` - (Optional) The names of the quotas.
    * `quota_types` - (Optional) The resource IDs if these are resource-specific quotas, `global` if they are not.
    * `short_descriptions` - (Optional) The description of the quotas.

## Attribute Reference

The following attributes are exported:

* `account_id` - The account ID of the owner of the quotas.
* `description` - The description of the quota.
* `max_value` - The maximum value of the quota for the OUTSCALE user account (if there is no limit, `0`).
* `name` - The unique name of the quota.
* `quota_collection` - The group name of the quota.
* `quota_type` - The resource ID if it is a resource-specific quota, `global` if it is not.
* `short_description` - The description of the quota.
* `used_value` - The limit value currently used by the OUTSCALE user account.
