---
layout: "outscale"
page_title: "OUTSCALE: outscale_access_key"
sidebar_current: "outscale-access-key"
description: |-
  [Provides information about a specific access key.]
---

# outscale_access_key Data Source

Provides information about a specific access key.
For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Access-Keys.html).
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-accesskey).

## Example Usage

```hcl
data "outscale_access_key" "access_key01" { 
    filter {
        name   = "access_key_ids"
        values = ["ABCDEFGHIJ0123456789"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `access_key_ids` - (Optional) The IDs of the access keys.
    * `states` - (Optional) The states of the access keys (`ACTIVE` \| `INACTIVE`).

## Attribute Reference

The following attributes are exported:

* `access_key_id` - The ID of the access key.
* `creation_date` - The date and time of creation of the access key.
* `last_modification_date` - The date and time of the last modification of the access key.
* `state` - The state of the access key (`ACTIVE` if the key is valid for API calls, or `INACTIVE` if not).
