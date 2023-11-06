---
layout: "outscale"
page_title: "OUTSCALE: outscale_access_keys"
sidebar_current: "outscale-access-keys"
description: |-
  [Provides information about access keys.]
---

# outscale_access_keys Data Source

Provides information about access keys.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Access-Keys.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-accesskey).

## Example Usage

```hcl
data "outscale_access_keys" "access_keys01" { 
    filter {
        name   = "access_key_ids"
        values = ["ABCDEFGHIJ0123456789", "0123456789ABCDEFGHIJ"]
    }
}
```

## Argument Reference

The following arguments are supported:

* `filter` - (Optional) A combination of a filter name and one or more filter values. You can specify this argument for as many filter names as you need. The filter name can be any of the following:
    * `access_key_ids` - (Optional) The IDs of the access keys.
    * `states` - (Optional) The states of the access keys (`ACTIVE` \| `INACTIVE`).
* `user_name` - (Optional) The name of the EIM user. By default, the user who sends the request (which can be the root account).

## Attribute Reference

The following attributes are exported:

* `access_keys` - A list of access keys.
    * `access_key_id` - The ID of the access key.
    * `creation_date` - The date and time (UTC) of creation of the access key.
    * `expiration_date` - The date (UTC) at which the access key expires.
    * `last_modification_date` - The date and time (UTC) of the last modification of the access key.
    * `state` - The state of the access key (`ACTIVE` if the key is valid for API calls, or `INACTIVE` if not).
