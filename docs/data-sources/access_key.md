---
layout: "outscale"
page_title: "OUTSCALE: outscale_access_key"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-access-key"
description: |-
  [Provides information about an access key.]
---

# outscale_access_key Data Source

Provides information about an access key.

For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Access-Keys.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-accesskey).

## Example Usage

### Get one of your own access keys (root account or user)

```hcl
data "outscale_access_key" "access_key01" { 
    filter {
        name   = "access_key_ids"
        values = ["ABCDEFGHIJ0123456789"]
    }
}
```

### Get the access key of another user

```hcl
data "outscale_access_key" "access_key01" {
    user_name = "user_name"
    filter {
        name   = "access_key_ids"
        values = ["XXXXXXXXX"]
    }
    filter {
        name   = "states"
        values = ["ACTIVE"]
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

* `access_key_id` - The ID of the access key.
* `creation_date` - The date and time (UTC) at which the access key was created.
* `expiration_date` - The date and time (UTC) at which the access key expires.
* `last_modification_date` - The date and time (UTC) at which the access key was last modified.
* `state` - The state of the access key (`ACTIVE` if the key is valid for API calls, or `INACTIVE` if not).
