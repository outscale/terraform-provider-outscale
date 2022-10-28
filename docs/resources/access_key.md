---
layout: "outscale"
page_title: "OUTSCALE: outscale_access_key"
sidebar_current: "outscale-access-key"
description: |-
  [Manages an access key.]
---

# outscale_access_key Resource

Manages an access key.

!> *Warning* When creating an access key, the secret key is stored in the Terraform state. For security reasons, it is strongly recommended to create access keys using the API rather than the Terraform resource. For more information on how to create access keys using the OUTSCALE API, see the [API documentation](https://docs.outscale.com/api#createaccesskey).


For more information on this resource, see the [User Guide](https://docs.outscale.com/en/userguide/About-Access-Keys.html).  
For more information on this resource actions, see the [API documentation](https://docs.outscale.com/api#3ds-outscale-api-accesskey).

## Example Usage

```hcl
resource "outscale_access_key" "access_key01" {
    state           = "ACTIVE"
    expiration_date = "2023-01-01"
}
```

## Argument Reference

The following arguments are supported:

* `expiration_date` - (Optional) The date and time at which you want the access key to expire, in ISO 8601 format (for example, `2017-06-14` or `2017-06-14T00:00:00Z`). To remove an existing expiration date, use the method without specifying this parameter.
* `state` - (Optional) The state for the access key (`ACTIVE` | `INACTIVE`).

## Attribute Reference

The following attributes are exported:

* `access_key_id` - The ID of the access key.
* `creation_date` - The date and time (UTC) of creation of the access key.
* `expiration_date` - The date and time (UTC) at which the access key expires.
* `last_modification_date` - The date and time (UTC) of the last modification of the access key.
* `secret_key` - The access key that enables you to send requests.
* `state` - The state of the access key (`ACTIVE` if the key is valid for API calls, or `INACTIVE` if not).

## Import

An access key can be imported using its ID. For example:

```console

$ terraform import outscale_access_key.ImportedAccessKey ABCDEFGHIJ0123456789

```