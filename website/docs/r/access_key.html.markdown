---
layout: "outscale"
page_title: "OUTSCALE: outscale_access_key"
sidebar_current: "outscale-access-key"
description: |-
  [Manages an access key.]
---

# outscale_access_key Resource

Manages an access key.
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

* `expiration_date` - (Optional) The date and time at which you want the access key to expire, in ISO 8601 format (for example, `2017-06-14` or `2017-06-14T00:00:00Z`). If not specified, the access key has no expiration date.
* `state` - (Optional) The state for the access key (`ACTIVE` | `INACTIVE`).

## Attribute Reference

The following attributes are exported:

* `access_key_id` - The ID of the secret access key.
* `creation_date` - The date and time of creation of the secret access key.
* `expiration_date` - The date at which the access key expires.
* `last_modification_date` - The date and time of the last modification of the secret access key.
* `secret_key` - The secret access key that enables you to send requests.
* `state` - The state of the secret access key (`ACTIVE` if the key is valid for API calls, or `INACTIVE` if not).

## Import

An access key can be imported using its ID. For example:

```console

$ terraform import outscale_access_key.ImportedAccessKey ABCDEFGHIJ0123456789

```