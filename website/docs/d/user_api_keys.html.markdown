---
layout: "outscale"
page_title: "Outscale: outscale_user_api_keys"
sidebar_current: "docs-outscale-resource-user-api-keys"
description: |-
    Returns information about the access key IDs of a specified user.
    If the user does not have any access key ID, this action returns an empty list.
---

# outscale_user_api_keys

Returns information about the access key IDs of a specified user.
If the user does not have any access key ID, this action returns an empty list.

## Example Usage

```hcl
resource "outscale_user" "a_user" {
        user_name = "%s"
}
resource "outscale_user_api_keys" "a_key" {
        user_name = "${outscale_user.a_user.user_name}"
}

data "outscale_user_api_keys" "test_key" {
        user_name = "${outscale_user_api_keys.a_key.user_name}"
}
```

## Argument Reference

The following arguments are supported:

* `user_name` - (Optional) The name of the user.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `acess_key_metadata`  - A list of access keys and their metadata.
* `request_id` - The ID of the request.

### Access Key Metadata List

 The `acess_key_metadata` attribute contains a list with its elements with the following atrributes. 

* `access_key_id` - The access key ID.
* `secret_access_key` - The secret key that enables you to send requests.
* `status` - The state of the access key (active if the key is valid for API calls, or inactive if not).

See more detailed information [Outscale API Documetaion - EIM API Access Keys](http://docs.outscale.com/api_eim/index.html#_access_keys)