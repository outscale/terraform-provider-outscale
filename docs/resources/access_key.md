---
layout: "outscale"
page_title: "OUTSCALE: outscale_access_key"
subcategory: "OUTSCALE API"
sidebar_current: "outscale-access-key"
description: |-
  [Manages an access key.]
---

# outscale_access_key Resource

Manages an access key.

!> When creating an access key, the secret key is stored in the Terraform state. For security reasons, it is strongly recommended to create access keys using the API rather than the Terraform resource. For more information on how to create access keys using the OUTSCALE API, see the [API documentation](https://docs.outscale.com/api#createaccesskey).


For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Access-Keys.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-accesskey).

## Example Usage

### Creating an access key for yourself

```hcl
resource "outscale_access_key" "access_key01" {
    state           = "ACTIVE"
    expiration_date = "2028-01-01"
}
```

### Creating an access key for another user

```hcl
resource "outscale_access_key" "access_key_eim01" {
    user_name       = outscale_user.user-1.user_name
    state           = "ACTIVE"
    expiration_date = "2028-01-01"
    depends_on      = [outscale_user.user-1]
}
```

## Argument Reference

The following arguments are supported:

* `expiration_date` - (Optional) The date and time, or the date, at which you want the access key to expire, in ISO 8601 format (for example, `2020-06-14T00:00:00.000Z`, or `2020-06-14`). To remove an existing expiration date, use the method without specifying this parameter.
* `state` - (Optional) The state for the access key (`ACTIVE` | `INACTIVE`).
* `user_name` - (Optional) The name of the EIM user that owns the key to be created. If you do not specify a user name, this action creates an access key for the user who sends the request (which can be the root account).

## Attribute Reference

The following attributes are exported:

* `access_key_id` - The ID of the access key.
* `creation_date` - The date and time (UTC) at which the access key was created.
* `expiration_date` - The date and time (UTC) at which the access key expires.
* `last_modification_date` - The date and time (UTC) at which the access key was last modified.
* `secret_key` - The secret key that enables you to send requests.
* `state` - The state of the access key (`ACTIVE` if the key is valid for API calls, or `INACTIVE` if not).

## Import

An access key can be imported using its ID. For example:

```console

$ terraform import outscale_access_key.ImportedAccessKey ABCDEFGHIJ0123456789

```