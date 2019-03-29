---
layout: "outscale"
page_title: "Outscale: outscale_user_api_keys"
sidebar_current: "docs-outscale-resource-user-api-keys"
description: |-
  Provides an EIM access key. This is a set of credentials that allow API requests to be made as an EIM user.
---

# outscale_user_api_keys

Provides an EIM access key. This is a set of credentials that allow API requests to be made as an EIM user.

## Example Usage

```hcl
resource "outscale_user" "a_user" {
        user_name = "user-test"
}
resource "outscale_user_api_keys" "a_key" {
        user_name = "${outscale_user.a_user.user_name}"
}
```

## Argument Reference

The following arguments are supported:

* `user_name` - (Optional) The user name of the owner of the key to be created. If you do not specify a user name, Outscale determines one based on the access key ID that sent the request.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `access_key_id` - The access key ID.
* `secret_access_key` - The secret key that enables you to send requests.
* `status` - The state of the access key (active if the key is valid for API calls, or inactive if not).
* `request_id` - The ID of the request.

See more detailed information [Outscale API Documetaion - EIM API Access Keys](http://docs.outscale.com/api_eim/index.html#_access_keys)